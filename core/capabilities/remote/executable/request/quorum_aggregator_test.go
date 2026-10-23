package request

import (
	"crypto/rand"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_quorumAggregator_VerifyAttestation(t *testing.T) {
	t.Parallel()

	const workflowExecutionID = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac0ce0abbadeed"
	const referenceID = "step1"
	spendUnit, spendValue := "testunit", "42"

	val, err := values.NewMap(map[string]any{"response": "attested"})
	require.NoError(t, err)
	valueProto := values.ProtoMap(val)
	valueBytes, err := proto.Marshal(valueProto)
	require.NoError(t, err)

	configDigest := ocrtypes.ConfigDigest{1, 2, 3, 4, 5}
	seqNr := uint64(100)

	kb1, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	kb2, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)

	validResp := commoncap.CapabilityResponse{
		Metadata: commoncap.ResponseMetadata{
			Metering: []commoncap.MeteringNodeDetail{
				{SpendUnit: spendUnit, SpendValue: spendValue},
			},
		},
		Payload: &anypb.Any{TypeUrl: "type.googleapis.com/values.v1.Map", Value: valueBytes},
	}

	reportData, err := commoncap.ResponseToReportData(workflowExecutionID, referenceID, valueBytes, validResp.Metadata)
	require.NoError(t, err)

	sig1, err := kb1.Sign3(configDigest, seqNr, reportData[:])
	require.NoError(t, err)
	sig2, err := kb2.Sign3(configDigest, seqNr, reportData[:])
	require.NoError(t, err)

	signers := [][]byte{kb1.PublicKey(), kb2.PublicKey()}

	validResp.OCRAttestation = &commoncap.OCRAttestation{
		ConfigDigest:   configDigest,
		SequenceNumber: seqNr,
		Sigs: []commoncap.AttributedSignature{
			{Signer: 0, Signature: sig1},
			{Signer: 1, Signature: sig2},
		},
	}

	a := &quorumAggregator{
		lggr:                          logger.Test(t),
		signers:                       signers,
		workflowExecutionID:           workflowExecutionID,
		referenceID:                   referenceID,
		requiredResponseConfirmations: 2,
	}

	t.Run("not enough signers returns error", func(t *testing.T) {
		t.Parallel()
		aBad := &quorumAggregator{
			workflowExecutionID:           workflowExecutionID,
			referenceID:                   referenceID,
			lggr:                          logger.Test(t),
			requiredResponseConfirmations: 2,
		}
		err := aBad.verifyAttestation(validResp)
		require.Error(t, err)
		require.Contains(t, err.Error(), "number of configured OCR signers is less than required confirmations: got 0, need at least 2")
	})

	t.Run("not enough signatures returns error", func(t *testing.T) {
		t.Parallel()
		respFewSigs := commoncap.CapabilityResponse{
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{{SpendUnit: spendUnit, SpendValue: spendValue}},
			},
			Payload: &anypb.Any{TypeUrl: "type.googleapis.com/values.v1.Map", Value: valueBytes},
			OCRAttestation: &commoncap.OCRAttestation{
				ConfigDigest:   configDigest,
				SequenceNumber: seqNr,
				Sigs:           []commoncap.AttributedSignature{{Signer: 0, Signature: sig1}},
			},
		}
		err := a.verifyAttestation(respFewSigs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not enough signatures")
	})

	t.Run("invalid signer index returns error", func(t *testing.T) {
		t.Parallel()
		respBadSigner := commoncap.CapabilityResponse{
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{{SpendUnit: spendUnit, SpendValue: spendValue}},
			},
			Payload: &anypb.Any{TypeUrl: "type.googleapis.com/values.v1.Map", Value: valueBytes},
			OCRAttestation: &commoncap.OCRAttestation{
				ConfigDigest:   configDigest,
				SequenceNumber: seqNr,
				Sigs: []commoncap.AttributedSignature{
					{Signer: 0, Signature: sig1},
					{Signer: 99, Signature: sig2},
				},
			},
		}
		err := a.verifyAttestation(respBadSigner)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signer index")
	})

	t.Run("duplicate signature returns error", func(t *testing.T) {
		t.Parallel()
		respDupSig := commoncap.CapabilityResponse{
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{{SpendUnit: spendUnit, SpendValue: spendValue}},
			},
			Payload: &anypb.Any{TypeUrl: "type.googleapis.com/values.v1.Map", Value: valueBytes},
			OCRAttestation: &commoncap.OCRAttestation{
				ConfigDigest:   configDigest,
				SequenceNumber: seqNr,
				Sigs: []commoncap.AttributedSignature{
					{Signer: 0, Signature: sig1},
					{Signer: 0, Signature: sig1},
				},
			},
		}
		err := a.verifyAttestation(respDupSig)
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate signature")
	})

	t.Run("invalid signature returns error", func(t *testing.T) {
		t.Parallel()
		badSig := make([]byte, 65)
		_, err := rand.Read(badSig)
		require.NoError(t, err)
		respBadSig := commoncap.CapabilityResponse{
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{{SpendUnit: spendUnit, SpendValue: spendValue}},
			},
			Payload: &anypb.Any{TypeUrl: "type.googleapis.com/values.v1.Map", Value: valueBytes},
			OCRAttestation: &commoncap.OCRAttestation{
				ConfigDigest:   configDigest,
				SequenceNumber: seqNr,
				Sigs: []commoncap.AttributedSignature{
					{Signer: 0, Signature: sig1},
					{Signer: 1, Signature: badSig},
				},
			},
		}
		err = a.verifyAttestation(respBadSig)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signature")
	})

	t.Run("wrong payload bytes produces invalid signature", func(t *testing.T) {
		t.Parallel()
		wrongBytes := []byte("tampered")
		respWrongPayload := commoncap.CapabilityResponse{
			Metadata: commoncap.ResponseMetadata{
				Metering: []commoncap.MeteringNodeDetail{{SpendUnit: spendUnit, SpendValue: spendValue}},
			},
			Payload: &anypb.Any{TypeUrl: "x", Value: wrongBytes},
			OCRAttestation: &commoncap.OCRAttestation{
				ConfigDigest:   configDigest,
				SequenceNumber: seqNr,
				Sigs: []commoncap.AttributedSignature{
					{Signer: 0, Signature: sig1},
					{Signer: 1, Signature: sig2},
				},
			},
		}
		err := a.verifyAttestation(respWrongPayload)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signature")
	})

	t.Run("valid attestation succeeds", func(t *testing.T) {
		t.Parallel()
		err := a.verifyAttestation(validResp)
		require.NoError(t, err)
	})
}

func Test_quorumAggregator_hasValidAttestation_minResponsesToAggregateAboveFPlus1(t *testing.T) {
	t.Parallel()

	const workflowExecutionID = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac0ce0abbadeed"
	const referenceID = "step1"

	val, err := values.NewMap(map[string]any{"response": "attested"})
	require.NoError(t, err)
	valueProto := values.ProtoMap(val)
	valueBytes, err := proto.Marshal(valueProto)
	require.NoError(t, err)

	configDigest := ocrtypes.ConfigDigest{1, 2, 3, 4, 5}
	seqNr := uint64(100)

	kb1, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	kb2, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)

	resp := commoncap.CapabilityResponse{
		Metadata: commoncap.ResponseMetadata{
			Metering: []commoncap.MeteringNodeDetail{
				{SpendUnit: "testunit", SpendValue: "42"},
			},
		},
		Payload: &anypb.Any{TypeUrl: "type.googleapis.com/values.v1.Map", Value: valueBytes},
	}

	reportData, err := commoncap.ResponseToReportData(workflowExecutionID, referenceID, valueBytes, resp.Metadata)
	require.NoError(t, err)

	sig1, err := kb1.Sign3(configDigest, seqNr, reportData[:])
	require.NoError(t, err)
	sig2, err := kb2.Sign3(configDigest, seqNr, reportData[:])
	require.NoError(t, err)

	// report_attestation.go only ever produces F+1 signatures (it stops collecting as
	// soon as it clears the Byzantine fault threshold), so with remoteDonF=1 this
	// attestation carries exactly F+1=2 signatures - the most it could ever have.
	resp.OCRAttestation = &commoncap.OCRAttestation{
		ConfigDigest:   configDigest,
		SequenceNumber: seqNr,
		Sigs: []commoncap.AttributedSignature{
			{Signer: 0, Signature: sig1},
			{Signer: 1, Signature: sig2},
		},
	}

	signers := [][]byte{kb1.PublicKey(), kb2.PublicKey()}

	t.Run("requiredResponseConfirmations at F+1 uses attestation", func(t *testing.T) {
		t.Parallel()
		a := &quorumAggregator{
			lggr:                          logger.Test(t),
			signers:                       signers,
			workflowExecutionID:           workflowExecutionID,
			referenceID:                   referenceID,
			remoteDonF:                    1,
			requiredResponseConfirmations: 2, // F+1
		}
		require.True(t, a.hasValidAttestation(resp))
	})

	t.Run("requiredResponseConfirmations above F+1 falls back to identical responses instead of erroring", func(t *testing.T) {
		t.Parallel()
		a := &quorumAggregator{
			lggr:                          logger.Test(t),
			signers:                       signers,
			workflowExecutionID:           workflowExecutionID,
			referenceID:                   referenceID,
			remoteDonF:                    1,
			requiredResponseConfirmations: 3, // minResponsesToAggregate configured above F+1
		}
		// The attestation is otherwise perfectly valid, but since it can never carry
		// more than F+1 signatures, hasValidAttestation must decline to use it (rather
		// than calling verifyAttestation and logging a spurious "likely a bug" error).
		require.False(t, a.hasValidAttestation(resp))
	})
}

func newAggForQuorumTest(t *testing.T, remoteNodeCount, required, replies int, responseIDCount map[[32]byte]int) *quorumAggregator {
	t.Helper()

	return &quorumAggregator{
		lggr:                          logger.Test(t),
		requiredResponseConfirmations: required,
		remoteNodeCount:               remoteNodeCount,
		responseIDCount:               responseIDCount,
		replies:                       replies,
	}
}

func Test_quorumAggregator_quorumStillPossible(t *testing.T) {
	t.Parallel()

	// 7-node DON, F+1 = 3: six unique responses, one pending → max 2 matches possible.
	t.Run("7 DON 6 unique 1 pending unreachable", func(t *testing.T) {
		t.Parallel()
		counts := make(map[[32]byte]int)
		for i := range 6 {
			counts[[32]byte{byte(i)}] = 1
		}
		a := newAggForQuorumTest(t, 7, 3, 6, counts)
		pending := a.pending()
		require.False(t, a.quorumStillPossible(pending))
	})

	t.Run("7 DON 5 unique 2 pending still possible", func(t *testing.T) {
		t.Parallel()
		counts := make(map[[32]byte]int)
		for i := range 5 {
			counts[[32]byte{byte(i)}] = 1
		}
		a := newAggForQuorumTest(t, 7, 3, 5, counts)
		pending := a.pending()
		require.True(t, a.quorumStillPossible(pending))
	})

	t.Run("7 DON 7 unique all received unreachable", func(t *testing.T) {
		t.Parallel()
		counts := make(map[[32]byte]int)
		for i := range 7 {
			counts[[32]byte{byte(i)}] = 1
		}
		a := newAggForQuorumTest(t, 7, 3, 7, counts)
		pending := a.pending()
		require.False(t, a.quorumStillPossible(pending))
	})

	t.Run("7 DON 2 matching 5 unique still possible with 2 pending", func(t *testing.T) {
		t.Parallel()
		counts := map[[32]byte]int{
			{1}: 2,
			{2}: 1,
			{3}: 1,
			{4}: 1,
			{5}: 1,
			{6}: 1,
		}
		a := newAggForQuorumTest(t, 7, 3, 6, counts)
		pending := a.pending()
		require.True(t, a.quorumStillPossible(pending))
	})
}

func Test_quorumAggregator_trySendQuorumUnreachableError(t *testing.T) {
	t.Parallel()

	counts := make(map[[32]byte]int)
	for i := range 6 {
		counts[[32]byte{byte(i)}] = 1
	}

	a := newAggForQuorumTest(t, 7, 3, 6, counts)

	respErr := a.trySendQuorumUnreachableError()
	require.Error(t, respErr)

	var capErr caperrors.Error
	require.ErrorAs(t, respErr, &capErr)
	require.Equal(t, caperrors.OriginSystem, capErr.Origin())
	require.Equal(t, caperrors.ConsensusFailed, capErr.Code())
	assert.Contains(t, capErr.Error(),
		"[100]ConsensusFailed: response quorum unreachable: not enough matching capability responses: received 6/7 peer responses with 6 unique payloads; best match count 1, need 3 (1 responses pending)")

	// Same unwrap chain as capability_executor after client.Execute wraps with %w.
	wrapped := fmt.Errorf("error executing request: %w", respErr)
	var capErrFromWrapped caperrors.Error
	require.ErrorAs(t, wrapped, &capErrFromWrapped)
	require.Equal(t, caperrors.ConsensusFailed, capErrFromWrapped.Code())
	assert.Contains(t, capErrFromWrapped.Error(),
		"[100]ConsensusFailed: response quorum unreachable: not enough matching capability responses: received 6/7 peer responses with 6 unique payloads; best match count 1, need 3 (1 responses pending)")
}

func Test_quorumAggregator_OnUnmarshalable(t *testing.T) {
	t.Parallel()

	// n=4, F=1 → required=2.
	a := newQuorumAggregator(logger.Test(t), nil, "weid", "ref", 4, 1, 0)

	// Garbage payloads surface an error to the OnMessage caller without
	// resolving the request or counting toward the error quorum.
	for i := range 2 {
		resp, err := a.OnUnmarshalable(p2ptypes.PeerID{}, errors.New("bad payload"))
		require.Nil(t, resp)
		var obs observationError
		require.ErrorAs(t, err, &obs)
		assert.Contains(t, err.Error(), "failed to unmarshal capability response: bad payload")
		// counted as a reply (pending) but not as an error
		assert.Equal(t, 4-(i+1), a.pending())
		assert.Equal(t, 0, a.totalErrorCount)
	}
}

func Test_quorumAggregator_OnError(t *testing.T) {
	t.Parallel()

	t.Run("payload not available never resolves the request", func(t *testing.T) {
		t.Parallel()
		// n=4, F=1 → required=2 → unexpected-state threshold = n-required+1 = 3.
		a := newQuorumAggregator(logger.Test(t), nil, "weid", "ref", 4, 1, 0)

		msg := commoncap.ErrResponsePayloadNotAvailable.Error()
		for range 2 {
			resp, err := a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, msg)
			require.Nil(t, resp)
			require.NoError(t, err)
		}

		// The threshold-th payload-not-available response surfaces the
		// "unexpected state" error to the OnMessage caller without resolving.
		resp, err := a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, msg)
		require.Nil(t, resp)
		var obs observationError
		require.ErrorAs(t, err, &obs)
		assert.Contains(t, err.Error(), "unexpected state: received 3 payload not available responses, while max allowed is 2. This means a bug in the code, please investigate")
		assert.Equal(t, 0, a.totalErrorCount)
		assert.Equal(t, 1, a.pending())
	})

	t.Run("F+1 identical errors resolve the request", func(t *testing.T) {
		t.Parallel()
		// n=4, F=1 → required=2.
		a := newQuorumAggregator(logger.Test(t), nil, "weid", "ref", 4, 1, 0)

		resp, err := a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, "boom")
		require.Nil(t, resp)
		require.NoError(t, err)

		resp, err = a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, "boom")
		require.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "boom")
	})

	t.Run("too many different errors resolve the request early", func(t *testing.T) {
		t.Parallel()
		// n=4, F=1 → required=2 → total error threshold = n-required+1 = 3.
		a := newQuorumAggregator(logger.Test(t), nil, "weid", "ref", 4, 1, 0)

		for _, msg := range []string{"e1", "e2"} {
			resp, err := a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, msg)
			require.Nil(t, resp)
			require.NoError(t, err)
		}

		resp, err := a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, "e3")
		require.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "received 3 errors, last error")
	})
}
