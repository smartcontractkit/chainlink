package request

import (
	"crypto/rand"
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

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
