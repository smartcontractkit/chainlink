package fakes

import (
	"crypto/ed25519"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

const (
	testWorkflowID    = "ffffaabbccddeeff00112233aabbccddeeff00112233aabbccddeeff00112233"
	testExecutionID   = "aabbccddeeff00112233aabbccddeeff00112233aabbccddeeff00112233eeee"
	testWorkflowOwner = "1100000000000000000000000000000000000000"
	testWorkflowName  = "00112233445566778899"
	testRefID         = "0011"
)

func Test_Simple_EVMEncoder(t *testing.T) {
	t.Parallel()
	nSigners := 4
	signers := make([]ocr2key.KeyBundle, 0, nSigners)
	for range nSigners {
		signers = append(signers, ocr2key.MustNewInsecure(SeedForKeys(), corekeys.EVM))
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
	outputs, capErr := fakeConsensusNoDAG.Report(t.Context(), metadata, input)

	require.NoError(t, capErr)
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

func testReportMetadata() capabilities.RequestMetadata {
	return capabilities.RequestMetadata{
		WorkflowID:               testWorkflowID,
		WorkflowOwner:            testWorkflowOwner,
		WorkflowExecutionID:      testExecutionID,
		WorkflowName:             testWorkflowName,
		WorkflowDonID:            1,
		WorkflowDonConfigVersion: 1,
		ReferenceID:              testRefID,
	}
}

// stellarForwarderDigest is the digest the Stellar CRE forwarder verifies:
// keccak256(keccak256(raw_report) ‖ report_context).
func stellarForwarderDigest(rawReport, reportContext []byte) []byte {
	inner := crypto.Keccak256(rawReport)
	return crypto.Keccak256(append(inner, reportContext...))
}

// splitStellarSig splits public_key(32) ‖ signature(64), the layout the
// Stellar forwarder's Ed25519Signature carries.
func splitStellarSig(t *testing.T, b []byte) (ed25519.PublicKey, []byte) {
	t.Helper()
	require.Len(t, b, ed25519.PublicKeySize+ed25519.SignatureSize)
	return ed25519.PublicKey(b[:ed25519.PublicKeySize]), b[ed25519.PublicKeySize:]
}

func Test_Report_StellarEncoder(t *testing.T) {
	t.Parallel()
	evmSigners := []ocr2key.KeyBundle{
		ocr2key.MustNewInsecure(SeedForKeys(), corekeys.EVM),
		ocr2key.MustNewInsecure(SeedForKeys(), corekeys.EVM),
		ocr2key.MustNewInsecure(SeedForKeys(), corekeys.EVM),
	}
	fc := NewFakeConsensusNoDAG(evmSigners, logger.Test(t))

	payload := []byte("stellar_payload")
	for _, name := range []string{"stellar", "Stellar"} {
		out, capErr := fc.Report(t.Context(), testReportMetadata(), &sdkpb.ReportRequest{EncodedPayload: payload, EncoderName: name})
		require.NoError(t, capErr, name)
		resp := out.Response

		// Shared metadata header + raw payload, as parsed by the Stellar forwarder.
		require.Len(t, resp.RawReport, 109+len(payload))
		require.Equal(t, byte(1), resp.RawReport[0])
		require.Equal(t, payload, resp.RawReport[109:])
		require.Len(t, resp.ReportContext, 96)

		// One ed25519 signature per default signer, from distinct Stellar keys,
		// each valid under the forwarder's digest.
		stellarSigners := fc.signersFor(encoderStellar)
		require.Len(t, resp.Sigs, len(evmSigners))
		digest := stellarForwarderDigest(resp.RawReport, resp.ReportContext)
		seen := map[string]bool{}
		for i, sig := range resp.Sigs {
			pub, edSig := splitStellarSig(t, sig.Signature)
			require.Equal(t, []byte(stellarSigners[i].PublicKey()), []byte(pub))
			require.True(t, ed25519.Verify(pub, digest, edSig), "signature %d must verify under the forwarder digest", i)
			require.False(t, seen[string(pub)], "stellar signers must be distinct")
			seen[string(pub)] = true
		}
	}
}

func Test_Report_StellarEncoder_CustomSigners(t *testing.T) {
	t.Parallel()
	custom := deterministicSigners(corekeys.Stellar, 2)
	fc := NewFakeConsensusNoDAG(nil, logger.Test(t), WithEncoderSigners("Stellar", custom))

	out, capErr := fc.Report(t.Context(), testReportMetadata(), &sdkpb.ReportRequest{EncodedPayload: []byte("x"), EncoderName: "stellar"})
	require.NoError(t, capErr)
	require.Len(t, out.Response.Sigs, len(custom))
	digest := stellarForwarderDigest(out.Response.RawReport, out.Response.ReportContext)
	for i, sig := range out.Response.Sigs {
		pub, edSig := splitStellarSig(t, sig.Signature)
		require.Equal(t, []byte(custom[i].PublicKey()), []byte(pub))
		require.True(t, ed25519.Verify(pub, digest, edSig))
	}
}

func Test_Report_DefaultStellarSignersWithoutDefaults(t *testing.T) {
	t.Parallel()
	// No default signers still yields one Stellar signer: the forwarder
	// rejects reports without signatures.
	fc := NewFakeConsensusNoDAG(nil, logger.Test(t))
	out, capErr := fc.Report(t.Context(), testReportMetadata(), &sdkpb.ReportRequest{EncodedPayload: []byte("x"), EncoderName: "stellar"})
	require.NoError(t, capErr)
	require.Len(t, out.Response.Sigs, 1)
}

func Test_Report_EncoderNamesAreCaseInsensitive(t *testing.T) {
	t.Parallel()
	signer := ocr2key.MustNewInsecure(SeedForKeys(), corekeys.EVM)
	fc := NewFakeConsensusNoDAG([]ocr2key.KeyBundle{signer}, logger.Test(t))
	for _, name := range []string{"evm", "EVM", "Evm", "solana", "Solana", "SOLANA"} {
		_, capErr := fc.Report(t.Context(), testReportMetadata(), &sdkpb.ReportRequest{EncodedPayload: []byte("x"), EncoderName: name})
		require.NoError(t, capErr, name)
	}

	_, capErr := fc.Report(t.Context(), testReportMetadata(), &sdkpb.ReportRequest{EncodedPayload: []byte("x"), EncoderName: "nope"})
	require.Error(t, capErr)
	require.Contains(t, capErr.Error(), "unsupported encoder name: nope")
}
