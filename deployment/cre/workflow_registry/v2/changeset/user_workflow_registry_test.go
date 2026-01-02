package changeset

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	evmChain "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

func TestUserWorkflowOperations(t *testing.T) {
	t.Parallel()

	testWorkflowID := "1234567891234567891234567891234567891234567891234567891234567891"
	testWorkflowName := "Test Workflow"
	testDONFamily := "zone-a"
	testURL := "https://example.com"

	t.Run("link-owner allowlist-request upsert pause activate delete unlink-owner", func(t *testing.T) {
		fixture := setupTest(t)

		chain := fixture.rt.Environment().BlockChains.EVMChains()[fixture.selector]
		deployerKey := chain.DeployerKey

		t.Log("Testing link owner...")
		validity, proof, signature := generateAndSignOwnershipProof(
			t,
			common.HexToAddress(fixture.workflowRegistryAddress),
			deployerKey.From.Hex(),
			chain,
			deployerKey.From.Hex(),
			"123",
			"12",
			"WorkflowRegistry 2.0.0",
			0, // 0 for linking
		)
		linkOwnerInput := UserLinkOwnerInput{
			ValidityTimestamp:         validity,
			Proof:                     common.Bytes2Hex(proof.Bytes()),
			Signature:                 common.Bytes2Hex(signature),
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
		}
		linkOwnerChangeset := UserLinkOwner{}
		err := linkOwnerChangeset.VerifyPreconditions(fixture.rt.Environment(), linkOwnerInput)
		require.NoError(t, err, "link owner preconditions should pass")
		_, err = linkOwnerChangeset.Apply(fixture.rt.Environment(), linkOwnerInput)
		require.NoError(t, err, "link owner apply should pass")

		t.Log("Testing allowlist request...")
		allowlistInput := UserAllowlistRequestInput{
			ExpiryTimestamp:           mustConvertInt64ToUint32(time.Now().Add(48 * time.Hour).Unix()),
			RequestDigest:             generateRandom32BytesString(t),
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
		}
		allowlistChangeset := UserAllowlistRequest{}
		err = allowlistChangeset.VerifyPreconditions(fixture.rt.Environment(), allowlistInput)
		require.NoError(t, err, "allowlist request preconditions should pass")
		_, err = allowlistChangeset.Apply(fixture.rt.Environment(), allowlistInput)
		require.NoError(t, err, "allowlist request apply should pass")

		t.Log("Testing user workflow upsert...")
		upsertInput := UserWorkflowUpsertInput{
			WorkflowID:                testWorkflowID,
			WorkflowName:              testWorkflowName,
			WorkflowTag:               testWorkflowName,
			WorkflowStatus:            0,
			DonFamily:                 testDONFamily,
			BinaryURL:                 testURL,
			ConfigURL:                 testURL,
			Attributes:                "",
			KeepAlive:                 false,
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
		}
		changeset := UserWorkflowUpsert{}
		err = changeset.VerifyPreconditions(fixture.rt.Environment(), upsertInput)
		require.NoError(t, err, "preconditions should pass")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), upsertInput)
		require.NoError(t, err, "user workflow upsert apply should pass")
		assert.NotNil(t, csOutput, "user workflow upsert apply should pass")

		t.Log("Testing user workflow pause...")
		pauseInput := UserWorkflowPauseInput{
			WorkflowID:                testWorkflowID,
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
		}
		pauseChangeset := UserWorkflowPause{}
		err = pauseChangeset.VerifyPreconditions(fixture.rt.Environment(), pauseInput)
		require.NoError(t, err, "user workflow pause preconditions should pass")
		_, err = pauseChangeset.Apply(fixture.rt.Environment(), pauseInput)
		require.NoError(t, err, "user workflow pause apply should pass")

		t.Log("Testing user workflow activate...")
		activateInput := UserWorkflowActivateInput{
			WorkflowID:                testWorkflowID,
			DonFamily:                 testDONFamily,
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
		}
		activateChangeset := UserWorkflowActivate{}
		err = activateChangeset.VerifyPreconditions(fixture.rt.Environment(), activateInput)
		require.NoError(t, err, "user workflow activate preconditions should pass")
		_, err = activateChangeset.Apply(fixture.rt.Environment(), activateInput)
		require.NoError(t, err, "user workflow activate apply should pass")

		t.Log("Testing user batch workflow pause...")
		batchPauseInput := UserWorkflowBatchPauseInput{
			WorkflowIDs:               testWorkflowID,
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
		}
		batchPauseChangeset := UserWorkflowBatchPause{}
		err = batchPauseChangeset.VerifyPreconditions(fixture.rt.Environment(), batchPauseInput)
		require.NoError(t, err, "user workflow pause preconditions should pass")
		_, err = batchPauseChangeset.Apply(fixture.rt.Environment(), batchPauseInput)
		require.NoError(t, err, "user workflow pause apply should pass")

		t.Log("Testing user workflow delete...")
		deleteInput := UserWorkflowDeleteInput{
			WorkflowID:                testWorkflowID,
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
		}
		deleteChangeset := UserWorkflowDelete{}
		err = deleteChangeset.VerifyPreconditions(fixture.rt.Environment(), deleteInput)
		require.NoError(t, err, "user workflow delete preconditions should pass")
		_, err = deleteChangeset.Apply(fixture.rt.Environment(), deleteInput)
		require.NoError(t, err, "user workflow delete apply should pass")

		t.Log("Testing unlink owner...")
		validity, _, signature = generateAndSignOwnershipProof(
			t,
			common.HexToAddress(fixture.workflowRegistryAddress),
			deployerKey.From.Hex(),
			chain,
			deployerKey.From.Hex(),
			"123",
			"12",
			"WorkflowRegistry 2.0.0",
			1, // 1 for unlinking
		)
		unlinkOwnerInput := UserUnlinkOwnerInput{
			Address:                   deployerKey.From,
			ValidityTimestamp:         validity,
			Signature:                 common.Bytes2Hex(signature),
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
		}
		unlinkOwnerChangeset := UserUnlinkOwner{}
		err = unlinkOwnerChangeset.VerifyPreconditions(fixture.rt.Environment(), unlinkOwnerInput)
		require.NoError(t, err, "unlink owner preconditions should pass")
		_, err = unlinkOwnerChangeset.Apply(fixture.rt.Environment(), unlinkOwnerInput)
		require.NoError(t, err, "unlink owner apply should pass")
	})
}

func TestUserWorkflowOperationsMCMS(t *testing.T) {
	t.Parallel()

	testWorkflowID := "1234567891234567891234567891234567891234567891234567891234567891"
	testWorkflowName := "Test Workflow"
	testDONFamily := "zone-a"
	testURL := "http://example.com"

	t.Run("link-owner with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		chain := fixture.rt.Environment().BlockChains.EVMChains()[fixture.selector]
		deployerKey := chain.DeployerKey

		t.Log("Testing link owner with MCMS...")
		validity, proof, signature := generateAndSignOwnershipProof(
			t,
			common.HexToAddress(fixture.workflowRegistryAddress),
			deployerKey.From.Hex(),
			chain,
			deployerKey.From.Hex(),
			"123",
			"12",
			"WorkflowRegistry 2.0.0",
			0, // 0 for linking
		)
		linkOwnerInput := UserLinkOwnerInput{
			ValidityTimestamp:         validity,
			Proof:                     common.Bytes2Hex(proof.Bytes()),
			Signature:                 common.Bytes2Hex(signature),
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		}
		linkOwnerChangeset := UserLinkOwner{}
		err := linkOwnerChangeset.VerifyPreconditions(fixture.rt.Environment(), linkOwnerInput)
		require.NoError(t, err, "link owner with MCMS preconditions should pass")
		_, err = linkOwnerChangeset.Apply(fixture.rt.Environment(), linkOwnerInput)
		require.NoError(t, err, "link owner with MCMS apply should pass")
	})

	t.Run("upsert workflow with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)
		upsertInput := UserWorkflowUpsertInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
			WorkflowID:     testWorkflowID,
			WorkflowName:   testWorkflowName,
			WorkflowTag:    testWorkflowName,
			WorkflowStatus: 0,
			DonFamily:      testDONFamily,
			BinaryURL:      testURL,
			ConfigURL:      testURL,
			Attributes:     "",
			KeepAlive:      false,
		}
		changeset := UserWorkflowUpsert{}
		err := changeset.VerifyPreconditions(fixture.rt.Environment(), upsertInput)
		require.NoError(t, err, "MCMS preconditions should pass")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), upsertInput)
		require.NoError(t, err, "user workflow upsert apply should pass")
		assert.NotNil(t, csOutput, "user workflow upsert apply should pass")
		assert.NotNil(t, csOutput.Reports, "user workflow upsert apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from user workflow upsert")
	})

	t.Run("pause workflow with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)
		pauseInput := UserWorkflowPauseInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
			WorkflowID: testWorkflowID,
		}
		pauseChangeset := UserWorkflowPause{}
		err := pauseChangeset.VerifyPreconditions(fixture.rt.Environment(), pauseInput)
		require.NoError(t, err, "user workflow pause with MCMS preconditions should pass")
		csOutput, err := pauseChangeset.Apply(fixture.rt.Environment(), pauseInput)
		require.NoError(t, err, "user workflow pause with MCMS apply should pass")
		assert.NotNil(t, csOutput.Reports, "user workflow pause apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from user workflow pause")
	})

	t.Run("batch pause workflow with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)
		testWorkflowID2 := "1234567891234567891234567891234567891234567891234567891234567892"
		pauseInput := UserWorkflowBatchPauseInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
			WorkflowIDs: strings.Join([]string{testWorkflowID, testWorkflowID2}, ","),
		}
		pauseChangeset := UserWorkflowBatchPause{}
		err := pauseChangeset.VerifyPreconditions(fixture.rt.Environment(), pauseInput)
		require.NoError(t, err, "user workflow pause with MCMS preconditions should pass")
		csOutput, err := pauseChangeset.Apply(fixture.rt.Environment(), pauseInput)
		require.NoError(t, err, "user workflow pause with MCMS apply should pass")
		assert.NotNil(t, csOutput.Reports, "user workflow pause apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from user workflow pause")
	})

	t.Run("batch pause workflow with MCMS - duplicate workflowIDs", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)
		testWorkflowID2 := "1234567891234567891234567891234567891234567891234567891234567891"
		pauseInput := UserWorkflowBatchPauseInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
			WorkflowIDs: strings.Join([]string{testWorkflowID, testWorkflowID2}, ","),
		}
		pauseChangeset := UserWorkflowBatchPause{}
		err := pauseChangeset.VerifyPreconditions(fixture.rt.Environment(), pauseInput)
		require.Errorf(t, err, "user workflow pause with MCMS preconditions should fail due to duplicate workflow IDs")
		require.ErrorContains(t, err, "duplicate workflow ID detected", "error should mention duplicate workflow IDs")
	})

	t.Run("batch pause workflow with MCMS - invalid workflowIDs", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)
		testWorkflowID2 := "123456789123456789123456789123456789123456789123456789123456789Z" // invalid hex character 'Z'
		pauseInput := UserWorkflowBatchPauseInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
			WorkflowIDs: strings.Join([]string{testWorkflowID, testWorkflowID2}, ","),
		}
		pauseChangeset := UserWorkflowBatchPause{}
		err := pauseChangeset.VerifyPreconditions(fixture.rt.Environment(), pauseInput)
		require.Errorf(t, err, "user workflow pause with MCMS preconditions should fail due to duplicate workflow IDs")
		require.ErrorContains(t, err, "invalid workflow ID", "error should mention invalid workflow ID")
	})

	t.Run("activate workflow with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)
		activateInput := UserWorkflowActivateInput{
			WorkflowID:                testWorkflowID,
			DonFamily:                 testDONFamily,
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		}
		activateChangeset := UserWorkflowActivate{}
		err := activateChangeset.VerifyPreconditions(fixture.rt.Environment(), activateInput)
		require.NoError(t, err, "user workflow activate with MCMS preconditions should pass")
		csOutput, err := activateChangeset.Apply(fixture.rt.Environment(), activateInput)
		require.NoError(t, err, "user workflow activate with MCMS apply should pass")
		assert.NotNil(t, csOutput.Reports, "user workflow activate apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from user workflow activate")
	})

	t.Run("delete workflow with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)
		deleteInput := UserWorkflowDeleteInput{
			WorkflowID:                testWorkflowID,
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		}
		deleteChangeset := UserWorkflowDelete{}
		err := deleteChangeset.VerifyPreconditions(fixture.rt.Environment(), deleteInput)
		require.NoError(t, err, "user workflow delete with MCMS preconditions should pass")
		csOutput, err := deleteChangeset.Apply(fixture.rt.Environment(), deleteInput)
		require.NoError(t, err, "user workflow delete with MCMS apply should pass")
		assert.NotNil(t, csOutput.Reports, "user workflow delete apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from user workflow delete")
	})

	t.Run("unlink owner with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)
		chain := fixture.rt.Environment().BlockChains.EVMChains()[fixture.selector]
		deployerKey := chain.DeployerKey

		validity, _, signature := generateAndSignOwnershipProof(
			t,
			common.HexToAddress(fixture.workflowRegistryAddress),
			deployerKey.From.Hex(),
			chain,
			deployerKey.From.Hex(),
			"123",
			"12",
			"WorkflowRegistry 2.0.0",
			1, // 1 for unlinking
		)
		unlinkOwnerInput := UserUnlinkOwnerInput{
			Address:                   deployerKey.From,
			ValidityTimestamp:         validity,
			Signature:                 common.Bytes2Hex(signature),
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		}
		unlinkOwnerChangeset := UserUnlinkOwner{}
		err := unlinkOwnerChangeset.VerifyPreconditions(fixture.rt.Environment(), unlinkOwnerInput)
		require.NoError(t, err, "unlink owner with MCMS preconditions should pass")
		_, err = unlinkOwnerChangeset.Apply(fixture.rt.Environment(), unlinkOwnerInput)
		require.NoError(t, err, "unlink owner with MCMS apply should pass")
	})

	t.Run("allowlist request with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)
		allowlistInput := UserAllowlistRequestInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
			ExpiryTimestamp: mustConvertInt64ToUint32(time.Now().Add(48 * time.Hour).Unix()),
			RequestDigest:   generateRandom32BytesString(t),
		}
		allowlistChangeset := UserAllowlistRequest{}
		err := allowlistChangeset.VerifyPreconditions(fixture.rt.Environment(), allowlistInput)
		require.NoError(t, err, "allowlist request with MCMS preconditions should pass")
		_, err = allowlistChangeset.Apply(fixture.rt.Environment(), allowlistInput)
		require.NoError(t, err, "allowlist request with MCMS apply should pass")
	})
}

// Generate ownership proof based on test data and sign using private key of an allowed test signer.
// Make this signature recoverable by the WorkflowRegistry contract.
func generateAndSignOwnershipProof(
	t *testing.T,
	wrContractAddress common.Address,
	testAddress string, chain evmChain.Chain, signerAddress, orgID, nonce, version string, requestType uint8,
) (*big.Int, common.Hash, []byte) {
	ownershipProofHash := GenerateOwnershipProofHash(testAddress, orgID, nonce)
	ownershipProofHashBytes := common.HexToHash(ownershipProofHash)
	validityTimestamp := time.Now().Add(24 * time.Hour)
	digest, err := PreparePayloadForSigning(OwnershipProofSignaturePayload{
		RequestType:              requestType, // 0 for linking, 1 for unlinking
		WorkflowOwnerAddress:     common.HexToAddress(testAddress),
		ChainID:                  strconv.FormatUint(chainselectors.GETH_TESTNET.EvmChainID, 10),
		WorkflowRegistryContract: wrContractAddress,
		Version:                  version,
		ValidityTimestamp:        validityTimestamp,
		OwnershipProofHash:       ownershipProofHashBytes,
	})
	require.NoError(t, err, "failed to prepare payload for signing")

	sig, err := chain.SignHash(digest)
	require.NoError(t, err, "failed to sign the digest")

	recoveredPub, err := crypto.SigToPub(digest, sig)
	require.NoError(t, err, "failed to recover public key from signature")

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPub)
	require.Equal(
		t,
		common.HexToAddress(signerAddress),
		recoveredAddr,
		"recovered address should match the signer address",
	)

	// small signature fix: ECRECOVER() in Solidity expects recovery ID to be 27 or 28, not 0 or 1
	if sig[64] < 27 {
		sig[64] += 27
	}

	validityTimestampBigInt := big.NewInt(validityTimestamp.Unix())
	return validityTimestampBigInt, ownershipProofHashBytes, sig
}

func generateRandom32BytesString(t *testing.T) string {
	var b [32]byte
	_, err := rand.Read(b[:])
	require.NoError(t, err, "failed to generate random 32 bytes")
	return common.Bytes2Hex(b[:])
}

func mustConvertInt64ToUint32(value int64) uint32 {
	if value < 0 || value > math.MaxUint32 {
		panic(fmt.Sprintf("value %d out of range for uint32", value))
	}
	return uint32(value)
}
