package fakes

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

const (
	testWorkflowID    = "ffffaabbccddeeff00112233aabbccddeeff00112233aabbccddeeff00112233"
	testExecutionID   = "aabbccddeeff00112233aabbccddeeff00112233aabbccddeeff00112233eeee"
	testWorkflowOwner = "1100000000000000000000000000000000000000"
	testWorkflowName  = "00112233445566778899"
	testRefID         = "0011"
)

func Test_Simple_EVMEncoder(t *testing.T) {
	nSigners := 4
	signers := []ocr2key.KeyBundle{}
	for range nSigners {
		signers = append(signers, ocr2key.MustNewInsecure(SeedForKeys(), chaintype.EVM))
	}

	metadata := capabilities.RequestMetadata{
		WorkflowID:               testWorkflowID,
		WorkflowOwner:            testWorkflowOwner,
		WorkflowExecutionID:      testExecutionID,
		WorkflowName:             testWorkflowName,
		WorkflowDonID:            1,
		WorkflowDonConfigVersion: 1,
		ReferenceID:              testRefID,
	}

	input := &sdkpb.ReportRequest{
		EncodedPayload: []byte("test_observation_value"),
		EncoderName:    "evm",
	}
	fakeConsensusNoDAG := NewFakeConsensusNoDAG(signers, logger.Test(t))
	outputs, err := fakeConsensusNoDAG.Report(t.Context(), metadata, input)

	require.NoError(t, err)
	require.Len(t, outputs.Response.Sigs, nSigners)

	// validate signatures
	digest, err := ocr2types.BytesToConfigDigest(outputs.Response.ConfigDigest)
	require.NoError(t, err)
	fullHash := ocr2key.ReportToSigData3(digest, outputs.Response.SeqNr, outputs.Response.RawReport)
	for idx, sig := range outputs.Response.Sigs {
		signerPubkey, err2 := crypto.SigToPub(fullHash, sig.Signature)
		require.NoError(t, err2)
		recoveredAddr := crypto.PubkeyToAddress(*signerPubkey)
		expectedAddr := common.BytesToAddress(signers[idx].PublicKey())
		require.Equal(t, expectedAddr, recoveredAddr)
	}
}
