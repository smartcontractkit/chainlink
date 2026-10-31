package executable_test

import (
	"bytes"
	"context"
	"math/big"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/noop"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	reporttypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	aptoscappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/aptos"
	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	solcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/solana"
	stellarcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/stellar"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	wrDonID            = uint32(1)
	wrDonConfigVersion = uint32(3)
	wrSeqNr            = uint64(42)
)

var wrConfigDigest = ocrtypes.ConfigDigest{1, 2, 3, 4, 5, 6, 7, 8}

type wrRequest interface {
	proto.Message
	GetReport() *sdk.ReportResponse
}

type wrFamily struct {
	chainType    corekeys.ChainType
	capabilityID string
	newRequest   func(report *sdk.ReportResponse) wrRequest
}

var wrFamilies = []wrFamily{
	{corekeys.EVM, "evm:ChainSelector:1@1.0.0", func(r *sdk.ReportResponse) wrRequest {
		return &evmcappb.WriteReportRequest{Receiver: []byte("receiver"), Report: r}
	}},
	{corekeys.Solana, "solana:ChainSelector:1@1.0.0", func(r *sdk.ReportResponse) wrRequest {
		return &solcappb.WriteReportRequest{Receiver: []byte("receiver"), Report: r}
	}},
	{corekeys.Aptos, "aptos:ChainSelector:1@1.0.0", func(r *sdk.ReportResponse) wrRequest {
		return &aptoscappb.WriteReportRequest{Receiver: []byte("receiver"), Report: r}
	}},
	{corekeys.Stellar, "stellar:ChainSelector:1@1.0.0", func(r *sdk.ReportResponse) wrRequest {
		return &stellarcappb.WriteReportRequest{ContractId: "receiver", Report: r}
	}},
}

var evmFamily = wrFamilies[0]

type wrFixture struct {
	family wrFamily
	// keys holds the key bundle of every oracle by chain type.
	keys       []map[string]ocr2key.KeyBundle
	peers      []p2ptypes.PeerID
	signers    []ocrtypes.OnchainPublicKey
	callingDon commoncap.DON
}

func newWRFixture(t *testing.T, family wrFamily, n int, f uint8) wrFixture {
	fx := wrFixture{family: family}
	for range n {
		keys := map[string]ocr2key.KeyBundle{}
		for _, fam := range wrFamilies {
			kb, err := ocr2key.New(fam.chainType)
			require.NoError(t, err)
			keys[string(fam.chainType)] = kb
		}
		signer, err := ocrcommon.MarshalMultichainKeyBundle(keys)
		require.NoError(t, err)
		fx.keys = append(fx.keys, keys)
		fx.signers = append(fx.signers, signer)
		fx.peers = append(fx.peers, NewP2PPeerID(t))
	}
	fx.callingDon = commoncap.DON{ID: wrDonID, ConfigVersion: wrDonConfigVersion, F: f, Members: fx.peers}
	return fx
}

func newRawReport(t *testing.T, donID, configVersion uint32) []byte {
	md := reporttypes.Metadata{
		Version:          1,
		ExecutionID:      strings.Repeat("ab", 32),
		Timestamp:        1,
		DONID:            donID,
		DONConfigVersion: configVersion,
		WorkflowID:       strings.Repeat("cd", 32),
		WorkflowName:     strings.Repeat("ef", 10),
		WorkflowOwner:    strings.Repeat("01", 20),
		ReportID:         "0001",
	}
	raw, err := md.Encode()
	require.NoError(t, err)
	return append(raw, []byte("report body")...)
}

func (fx wrFixture) keyBundle(i int) ocr2key.KeyBundle {
	return fx.keys[i][string(fx.family.chainType)]
}

func sign(t *testing.T, kb ocr2key.KeyBundle, rawReport []byte) []byte {
	sig, err := ocrcommon.NewOCR3OnchainKeyringAdapter(kb).Sign(wrConfigDigest, wrSeqNr, ocr3types.ReportWithInfo[[]byte]{Report: rawReport})
	require.NoError(t, err)
	return sig
}

func verify(kb ocr2key.KeyBundle, rawReport []byte, sig *sdk.AttributedSignature) bool {
	return ocrcommon.NewOCR3OnchainKeyringAdapter(kb).Verify(kb.PublicKey(), wrConfigDigest, wrSeqNr, ocr3types.ReportWithInfo[[]byte]{Report: rawReport}, sig.Signature)
}

// sigsFrom signs rawReport with the keys at the given indexes, attributing each signature to its index.
func (fx wrFixture) sigsFrom(t *testing.T, rawReport []byte, idxs ...int) []*sdk.AttributedSignature {
	sigs := make([]*sdk.AttributedSignature, 0, len(idxs))
	for _, i := range idxs {
		sigs = append(sigs, &sdk.AttributedSignature{Signature: sign(t, fx.keyBundle(i), rawReport), SignerId: uint32(i)}) //nolint:gosec // test index
	}
	return sigs
}

func (fx wrFixture) newPayload(t *testing.T, rawReport []byte, sigs []*sdk.AttributedSignature) []byte {
	wrAny, err := anypb.New(fx.family.newRequest(&sdk.ReportResponse{
		ConfigDigest: wrConfigDigest[:],
		SeqNr:        wrSeqNr,
		RawReport:    rawReport,
		Sigs:         sigs,
	}))
	require.NoError(t, err)
	payload, err := pb.MarshalCapabilityRequest(commoncap.CapabilityRequest{
		Payload: wrAny,
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "workflowID",
			WorkflowExecutionID: "workflowExecutionID",
			WorkflowDonID:       wrDonID,
		},
	})
	require.NoError(t, err)
	return payload
}

func (fx wrFixture) decodeSigs(t *testing.T, payload []byte) []*sdk.AttributedSignature {
	req, err := pb.UnmarshalCapabilityRequest(payload)
	require.NoError(t, err)
	wrReq := fx.family.newRequest(nil)
	require.NoError(t, req.Payload.UnmarshalTo(wrReq))
	return wrReq.GetReport().GetSigs()
}

func (fx wrFixture) requireAllValid(t *testing.T, rawReport []byte, sigs []*sdk.AttributedSignature) {
	seen := map[uint32]bool{}
	for _, sig := range sigs {
		require.True(t, verify(fx.keyBundle(int(sig.SignerId)), rawReport, sig), "signature attributed to signer %d", sig.SignerId)
		require.False(t, seen[sig.SignerId], "duplicate signer %d", sig.SignerId)
		seen[sig.SignerId] = true
	}
}

func signerIDs(sigs []*sdk.AttributedSignature) []uint32 {
	ids := make([]uint32, len(sigs))
	for i, s := range sigs {
		ids[i] = s.SignerId
	}
	return ids
}

func (fx wrFixture) newAccumulatorWithSigners(t *testing.T, signers map[uint32][]ocrtypes.OnchainPublicKey) types.RequestAccumulator {
	factory, err := executable.NewWriteReportSignatureAggregator(signers, logger.Test(t))
	require.NoError(t, err)
	acc, err := factory.NewAccumulator(t.Context(), fx.family.capabilityID, fx.callingDon)
	require.NoError(t, err)
	return acc
}

func (fx wrFixture) newAccumulator(t *testing.T) types.RequestAccumulator {
	return fx.newAccumulatorWithSigners(t, map[uint32][]ocrtypes.OnchainPublicKey{wrDonID: fx.signers})
}

func requireReady(t *testing.T, acc types.RequestAccumulator) []byte {
	payload, ready, err := acc.Payload()
	require.NoError(t, err)
	require.True(t, ready)
	return payload
}

func requireNotReady(t *testing.T, acc types.RequestAccumulator) {
	_, ready, err := acc.Payload()
	require.NoError(t, err)
	require.False(t, ready)
}

func TestWriteReportSignatureAggregator(t *testing.T) {
	t.Parallel()

	rawReport := newRawReport(t, wrDonID, wrDonConfigVersion)

	for _, family := range wrFamilies {
		t.Run(string(family.chainType), func(t *testing.T) {
			t.Parallel()
			fx := newWRFixture(t, family, 4, 1)

			t.Run("not ready before any request", func(t *testing.T) {
				t.Parallel()
				requireNotReady(t, fx.newAccumulator(t))
			})

			t.Run("merges different valid subsets from honest nodes", func(t *testing.T) {
				t.Parallel()
				acc := fx.newAccumulator(t)

				require.NoError(t, acc.Add(fx.peers[0], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 2, 3))))
				require.NoError(t, acc.Add(fx.peers[1], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 0, 1))))

				sigs := fx.decodeSigs(t, requireReady(t, acc))
				require.Equal(t, []uint32{0, 1}, signerIDs(sigs))
				fx.requireAllValid(t, rawReport, sigs)
			})

			t.Run("excludes corrupted signature from the last sender", func(t *testing.T) {
				t.Parallel()
				acc := fx.newAccumulator(t)

				require.NoError(t, acc.Add(fx.peers[0], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 1, 2))))
				corrupted := fx.sigsFrom(t, rawReport, 0, 1)
				corrupted[0].Signature[len(corrupted[0].Signature)-10] ^= 0xff
				corrupted[1].Signature = bytes.Repeat([]byte{1}, len(corrupted[1].Signature))
				require.NoError(t, acc.Add(fx.peers[1], fx.newPayload(t, rawReport, corrupted)))

				sigs := fx.decodeSigs(t, requireReady(t, acc))
				require.Equal(t, []uint32{1, 2}, signerIDs(sigs))
				fx.requireAllValid(t, rawReport, sigs)
			})

			t.Run("falls back to latest payload when quorum has too few valid signatures", func(t *testing.T) {
				t.Parallel()
				acc := fx.newAccumulator(t)

				corrupted := fx.sigsFrom(t, rawReport, 0, 1)
				corrupted[1].Signature[len(corrupted[1].Signature)-5] ^= 0xff
				latest := fx.newPayload(t, rawReport, corrupted)
				require.NoError(t, acc.Add(fx.peers[0], latest))
				require.Equal(t, latest, requireReady(t, acc))
			})

			t.Run("excludes valid signature from a key outside the workflow DON", func(t *testing.T) {
				t.Parallel()
				acc := fx.newAccumulator(t)

				outsider, err := ocr2key.New(family.chainType)
				require.NoError(t, err)
				sigs := []*sdk.AttributedSignature{
					{Signature: sign(t, fx.keyBundle(0), rawReport), SignerId: 0},
					{Signature: sign(t, outsider, rawReport), SignerId: 1},
				}
				require.NoError(t, acc.Add(fx.peers[0], fx.newPayload(t, rawReport, sigs)))
				require.NoError(t, acc.Add(fx.peers[1], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 2))))

				got := fx.decodeSigs(t, requireReady(t, acc))
				require.Equal(t, []uint32{0, 2}, signerIDs(got))
				fx.requireAllValid(t, rawReport, got)
			})

			t.Run("drops signature attributed to another signer", func(t *testing.T) {
				t.Parallel()
				acc := fx.newAccumulator(t)

				sig := sign(t, fx.keyBundle(2), rawReport)
				sigs := []*sdk.AttributedSignature{
					{Signature: sig, SignerId: 0},
					{Signature: sig, SignerId: 2},
				}
				require.NoError(t, acc.Add(fx.peers[0], fx.newPayload(t, rawReport, sigs)))
				require.NoError(t, acc.Add(fx.peers[1], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 3))))

				got := fx.decodeSigs(t, requireReady(t, acc))
				require.Equal(t, []uint32{2, 3}, signerIDs(got))
				fx.requireAllValid(t, rawReport, got)
			})

			t.Run("drops signer id out of range", func(t *testing.T) {
				t.Parallel()
				acc := fx.newAccumulator(t)

				sigs := []*sdk.AttributedSignature{
					{Signature: sign(t, fx.keyBundle(0), rawReport), SignerId: 0},
					{Signature: sign(t, fx.keyBundle(1), rawReport), SignerId: 99},
				}
				require.NoError(t, acc.Add(fx.peers[0], fx.newPayload(t, rawReport, sigs)))
				require.NoError(t, acc.Add(fx.peers[1], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 3))))

				got := fx.decodeSigs(t, requireReady(t, acc))
				require.Equal(t, []uint32{0, 3}, signerIDs(got))
				fx.requireAllValid(t, rawReport, got)
			})

			t.Run("result does not depend on arrival order", func(t *testing.T) {
				t.Parallel()
				payloads := [][]byte{
					fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 3, 1)),
					fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 2, 0)),
				}
				accA := fx.newAccumulator(t)
				require.NoError(t, accA.Add(fx.peers[0], payloads[0]))
				require.NoError(t, accA.Add(fx.peers[1], payloads[1]))
				accB := fx.newAccumulator(t)
				require.NoError(t, accB.Add(fx.peers[1], payloads[1]))
				require.NoError(t, accB.Add(fx.peers[0], payloads[0]))

				require.Equal(t, requireReady(t, accA), requireReady(t, accB))
			})

			t.Run("falls back to latest payload on DON ID mismatch", func(t *testing.T) {
				t.Parallel()
				acc := fx.newAccumulator(t)

				otherDonReport := newRawReport(t, wrDonID+1, wrDonConfigVersion)
				corrupted := fx.sigsFrom(t, otherDonReport, 0, 1)
				corrupted[0].Signature[len(corrupted[0].Signature)-3] ^= 0xff
				latest := fx.newPayload(t, otherDonReport, corrupted)
				require.NoError(t, acc.Add(fx.peers[0], latest))
				require.Equal(t, latest, requireReady(t, acc))
			})

			t.Run("falls back to latest payload on config version mismatch", func(t *testing.T) {
				t.Parallel()
				acc := fx.newAccumulator(t)

				oldVersionReport := newRawReport(t, wrDonID, wrDonConfigVersion-1)
				corrupted := fx.sigsFrom(t, oldVersionReport, 0, 1)
				corrupted[0].Signature[len(corrupted[0].Signature)-3] ^= 0xff
				latest := fx.newPayload(t, oldVersionReport, corrupted)
				require.NoError(t, acc.Add(fx.peers[0], latest))
				require.Equal(t, latest, requireReady(t, acc))
			})

			t.Run("falls back to latest payload when calling DON signers are unknown", func(t *testing.T) {
				t.Parallel()
				acc := fx.newAccumulatorWithSigners(t, map[uint32][]ocrtypes.OnchainPublicKey{})

				latest := fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 0))
				require.NoError(t, acc.Add(fx.peers[0], latest))
				require.Equal(t, latest, requireReady(t, acc))
			})

			t.Run("drops signatures of an undecodable signer", func(t *testing.T) {
				t.Parallel()
				signers := slices.Clone(fx.signers)
				signers[1] = ocrtypes.OnchainPublicKey{1, 2, 3}
				acc := fx.newAccumulatorWithSigners(t, map[uint32][]ocrtypes.OnchainPublicKey{wrDonID: signers})

				require.NoError(t, acc.Add(fx.peers[0], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 0, 1))))
				require.NoError(t, acc.Add(fx.peers[1], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 2))))

				got := fx.decodeSigs(t, requireReady(t, acc))
				require.Equal(t, []uint32{0, 2}, signerIDs(got))
			})
		})
	}

	t.Run("falls back to latest payload when signers have no key of the report family", func(t *testing.T) {
		t.Parallel()
		fx := newWRFixture(t, wrFamilies[1], 4, 1)
		signers := make([]ocrtypes.OnchainPublicKey, len(fx.keys))
		for i, keys := range fx.keys {
			var err error
			signers[i], err = ocrcommon.MarshalMultichainKeyBundle(map[string]ocr2key.KeyBundle{string(corekeys.EVM): keys[string(corekeys.EVM)]})
			require.NoError(t, err)
		}
		acc := fx.newAccumulatorWithSigners(t, map[uint32][]ocrtypes.OnchainPublicKey{wrDonID: signers})

		latest := fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 0))
		require.NoError(t, acc.Add(fx.peers[0], latest))
		require.Equal(t, latest, requireReady(t, acc))
	})

	t.Run("drops high S variant of a valid secp256k1 signature", func(t *testing.T) {
		t.Parallel()
		fx := newWRFixture(t, evmFamily, 4, 1)
		acc := fx.newAccumulator(t)

		sigs := fx.sigsFrom(t, rawReport, 0, 1)
		sigs[0].Signature = highS(sigs[0].Signature)
		require.True(t, verify(fx.keyBundle(0), rawReport, sigs[0]), "high S variant should pass keyring verification")
		require.NoError(t, acc.Add(fx.peers[0], fx.newPayload(t, rawReport, sigs)))
		require.NoError(t, acc.Add(fx.peers[1], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 2))))

		got := fx.decodeSigs(t, requireReady(t, acc))
		require.Equal(t, []uint32{1, 2}, signerIDs(got))
	})
}

func highS(sig []byte) []byte {
	out := bytes.Clone(sig)
	s := new(big.Int).Sub(crypto.S256().Params().N, new(big.Int).SetBytes(sig[32:64]))
	s.FillBytes(out[32:64])
	out[64] ^= 1
	return out
}

type recordingCapability struct {
	abstractTestCapability
	mu       sync.Mutex
	requests []commoncap.CapabilityRequest
}

func (c *recordingCapability) Execute(_ context.Context, req commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.requests = append(c.requests, req)
	return commoncap.CapabilityResponse{}, nil
}

func TestServerRequest_WriteReportSignatureAggregation(t *testing.T) {
	t.Parallel()

	fx := newWRFixture(t, evmFamily, 4, 1)
	rawReport := newRawReport(t, wrDonID, wrDonConfigVersion)
	capabilityPeerID := NewP2PPeerID(t)

	dispatcher := mocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil)

	capability := &recordingCapability{}
	sr, err := request.NewServerRequest(capability, types.MethodExecute, evmFamily.capabilityID, 2, capabilityPeerID, fx.callingDon,
		"requestMessageID", dispatcher, time.Minute, "WriteReport", limits.NewGateLimiter(false), fx.newAccumulator(t), logger.Test(t), noop.Meter{})
	require.NoError(t, err)

	send := func(sender p2ptypes.PeerID, payload []byte) {
		require.NoError(t, sr.OnMessage(t.Context(), &types.MessageBody{
			Sender:          sender[:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("requestMessageID"),
			CapabilityId:    evmFamily.capabilityID,
			CapabilityDonId: 2,
			CallerDonId:     wrDonID,
			Method:          types.MethodExecute,
			Payload:         payload,
		}))
	}

	send(fx.peers[0], fx.newPayload(t, rawReport, fx.sigsFrom(t, rawReport, 0, 1)))

	// the malicious node completes the quorum with a corrupted signature
	corrupted := fx.sigsFrom(t, rawReport, 2, 3)
	corrupted[0].Signature[7] ^= 0xff
	send(fx.peers[1], fx.newPayload(t, rawReport, corrupted))

	capability.mu.Lock()
	defer capability.mu.Unlock()
	require.Len(t, capability.requests, 1)
	var wrReq evmcappb.WriteReportRequest
	require.NoError(t, capability.requests[0].Payload.UnmarshalTo(&wrReq))
	require.Equal(t, []uint32{0, 1}, signerIDs(wrReq.GetReport().GetSigs()))
	fx.requireAllValid(t, rawReport, wrReq.GetReport().GetSigs())
	require.Equal(t, []byte("receiver"), wrReq.GetReceiver())
}
