package request

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func TestClientRequest_Expired_aggregationGrace(t *testing.T) {
	t.Parallel()

	requestTimeout := 100 * time.Millisecond
	t.Run("not expired before requestTimeout plus grace", func(t *testing.T) {
		t.Parallel()
		c := &ClientRequest{
			createdAt:      time.Now().Add(-requestTimeout - time.Millisecond), // less than defaultResponseAggregationGrace
			requestTimeout: requestTimeout,
		}
		require.False(t, c.Expired())
	})

	t.Run("expired after requestTimeout plus grace", func(t *testing.T) {
		t.Parallel()
		c := &ClientRequest{
			createdAt:      time.Now().Add(-requestTimeout - defaultResponseAggregationGrace - time.Millisecond),
			requestTimeout: requestTimeout,
		}
		require.True(t, c.Expired())
	})
}

func Test_ClientRequest_VerifyAttestation(t *testing.T) {
	const workflowExecutionID = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
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

	c := &ClientRequest{
		lggr:                          logger.Test(t),
		signers:                       signers,
		workflowExecutionID:           workflowExecutionID,
		referenceID:                   referenceID,
		requiredResponseConfirmations: 2,
	}

	t.Run("not enough signers returns error", func(t *testing.T) {
		cBad := &ClientRequest{
			workflowExecutionID:           workflowExecutionID,
			referenceID:                   referenceID,
			lggr:                          logger.Test(t),
			requiredResponseConfirmations: 2,
		}
		err := cBad.verifyAttestation(validResp)
		require.Error(t, err)
		require.Contains(t, err.Error(), "number of configured OCR signers is less than required confirmations: got 0, need at least 2")
	})

	t.Run("not enough signatures returns error", func(t *testing.T) {
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
		err := c.verifyAttestation(respFewSigs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not enough signatures")
	})

	t.Run("invalid signer index returns error", func(t *testing.T) {
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
		err := c.verifyAttestation(respBadSigner)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signer index")
	})

	t.Run("duplicate signature returns error", func(t *testing.T) {
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
		err := c.verifyAttestation(respDupSig)
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate signature")
	})

	t.Run("invalid signature returns error", func(t *testing.T) {
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
		err = c.verifyAttestation(respBadSig)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signature")
	})

	t.Run("wrong payload bytes produces invalid signature", func(t *testing.T) {
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
		err := c.verifyAttestation(respWrongPayload)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signature")
	})

	t.Run("valid attestation succeeds", func(t *testing.T) {
		err := c.verifyAttestation(validResp)
		require.NoError(t, err)
	})
}

func newReqForQuorumTest(t *testing.T, remoteNodeCount, required, responsesReceived int, responseIDCount map[[32]byte]int) *ClientRequest {
	t.Helper()

	responseReceived := make(map[p2ptypes.PeerID]bool, responsesReceived)
	for i := 0; i < responsesReceived; i++ {
		var peer p2ptypes.PeerID
		peer[0] = byte(i)
		responseReceived[peer] = true
	}

	return &ClientRequest{
		lggr:                          logger.Test(t),
		requiredResponseConfirmations: required,
		remoteNodeCount:               remoteNodeCount,
		responseIDCount:               responseIDCount,
		responseReceived:              responseReceived,
	}
}

func TestClientRequest_quorumStillPossible(t *testing.T) {
	t.Parallel()

	// 7-node DON, F+1 = 3: six unique responses, one pending → max 2 matches possible.
	t.Run("7 DON 6 unique 1 pending unreachable", func(t *testing.T) {
		t.Parallel()
		counts := make(map[[32]byte]int)
		for i := 0; i < 6; i++ {
			counts[[32]byte{byte(i)}] = 1
		}
		c := newReqForQuorumTest(t, 7, 3, 6, counts)
		pending := c.pending()
		require.False(t, c.quorumStillPossible(pending))
	})

	t.Run("7 DON 5 unique 2 pending still possible", func(t *testing.T) {
		t.Parallel()
		counts := make(map[[32]byte]int)
		for i := 0; i < 5; i++ {
			counts[[32]byte{byte(i)}] = 1
		}
		c := newReqForQuorumTest(t, 7, 3, 5, counts)
		pending := c.pending()
		require.True(t, c.quorumStillPossible(pending))
	})

	t.Run("7 DON 7 unique all received unreachable", func(t *testing.T) {
		t.Parallel()
		counts := make(map[[32]byte]int)
		for i := 0; i < 7; i++ {
			counts[[32]byte{byte(i)}] = 1
		}
		c := newReqForQuorumTest(t, 7, 3, 7, counts)
		pending := c.pending()
		require.False(t, c.quorumStillPossible(pending))
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
		c := newReqForQuorumTest(t, 7, 3, 6, counts)
		pending := c.pending()
		require.True(t, c.quorumStillPossible(pending))
	})
}

func TestClientRequest_trySendQuorumUnreachableError(t *testing.T) {
	t.Parallel()

	counts := make(map[[32]byte]int)
	for i := 0; i < 6; i++ {
		counts[[32]byte{byte(i)}] = 1
	}

	c := newReqForQuorumTest(t, 7, 3, 6, counts)
	c.responseCh = make(chan clientResponse, 1)

	c.trySendQuorumUnreachableError()
	require.True(t, c.respSent)

	var respErr error
	select {
	case resp := <-c.responseCh:
		require.Error(t, resp.Err)
		respErr = resp.Err
	default:
		t.Fatal("expected error response on channel")
	}

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
