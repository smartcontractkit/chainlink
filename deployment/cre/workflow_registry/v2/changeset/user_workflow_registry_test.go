package changeset

import (
	"math/big"
	"strconv"
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
	testURL := "http://example.com"

	t.Run("link-owner upsert pause activate delete unlink-owner", func(t *testing.T) {
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
			ValidityTimestamp: validity,
			Proof:             common.Bytes2Hex(proof.Bytes()),
			Signature:         common.Bytes2Hex(signature),
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
			},
		}
		linkOwnerChangeset := UserLinkOwner{}
		err := linkOwnerChangeset.VerifyPreconditions(fixture.rt.Environment(), linkOwnerInput)
		require.NoError(t, err, "link owner preconditions should pass")
		_, err = linkOwnerChangeset.Apply(fixture.rt.Environment(), linkOwnerInput)
		require.NoError(t, err, "link owner apply should pass")

		t.Log("Testing user workflow upsert...")
		upsertInput := UserWorkflowUpsertInput{
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
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
		err = changeset.VerifyPreconditions(fixture.rt.Environment(), upsertInput)
		require.NoError(t, err, "preconditions should pass")
		t.Log("User workflow upsert preconditions passed")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), upsertInput)
		require.NoError(t, err, "user workflow upsert apply should pass")
		assert.NotNil(t, csOutput, "user workflow upsert apply should pass")

		t.Log("Testing user workflow pause...")
		pauseInput := UserWorkflowPauseInput{
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
			},
			WorkflowID: testWorkflowID,
		}
		pauseChangeset := UserWorkflowPause{}
		err = pauseChangeset.VerifyPreconditions(fixture.rt.Environment(), pauseInput)
		require.NoError(t, err, "user workflow pause preconditions should pass")
		_, err = pauseChangeset.Apply(fixture.rt.Environment(), pauseInput)
		require.NoError(t, err, "user workflow pause apply should pass")
		t.Log("User workflow paused successfully")

		t.Log("Testing user workflow activate...")
		activateInput := UserWorkflowActivateInput{
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
			},
			WorkflowID: testWorkflowID,
			DonFamily:  testDONFamily,
		}
		activateChangeset := UserWorkflowActivate{}
		err = activateChangeset.VerifyPreconditions(fixture.rt.Environment(), activateInput)
		require.NoError(t, err, "user workflow activate preconditions should pass")
		_, err = activateChangeset.Apply(fixture.rt.Environment(), activateInput)
		require.NoError(t, err, "user workflow activate apply should pass")
		t.Log("User workflow activated successfully")

		t.Log("Testing user workflow delete...")
		deleteInput := UserWorkflowDeleteInput{
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
			},
			WorkflowID: testWorkflowID,
		}
		deleteChangeset := UserWorkflowDelete{}
		err = deleteChangeset.VerifyPreconditions(fixture.rt.Environment(), deleteInput)
		require.NoError(t, err, "user workflow delete preconditions should pass")
		_, err = deleteChangeset.Apply(fixture.rt.Environment(), deleteInput)
		require.NoError(t, err, "user workflow delete apply should pass")
		t.Log("User workflow deleted successfully")

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
			Address:           deployerKey.From,
			ValidityTimestamp: validity,
			Signature:         common.Bytes2Hex(signature),
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
			},
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
			ValidityTimestamp: validity,
			Proof:             common.Bytes2Hex(proof.Bytes()),
			Signature:         common.Bytes2Hex(signature),
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				MCMSConfig: &contracts.MCMSConfig{
					MinDelay: 1 * time.Second,
					TimelockQualifierPerChain: map[uint64]string{
						fixture.selector: "",
					},
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

		t.Log("Testing user workflow upsert with MCMS preconditions...")
		upsertInput := UserWorkflowUpsertInput{
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				MCMSConfig: &contracts.MCMSConfig{
					MinDelay: 1 * time.Second,
					TimelockQualifierPerChain: map[uint64]string{
						fixture.selector: "",
					},
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
		t.Log("User workflow upsert with MCMS preconditions passed")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), upsertInput)
		require.NoError(t, err, "user workflow upsert apply should pass")
		assert.NotNil(t, csOutput, "user workflow upsert apply should pass")
		assert.NotNil(t, csOutput.Reports, "user workflow upsert apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from user workflow upsert")
	})

	t.Run("pause workflow with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Testing user workflow pause with MCMS...")
		pauseInput := UserWorkflowPauseInput{
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				MCMSConfig: &contracts.MCMSConfig{
					MinDelay: 1 * time.Second,
					TimelockQualifierPerChain: map[uint64]string{
						fixture.selector: "",
					},
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
		t.Log("User workflow paused with MCMS successfully")
	})

	t.Run("activate workflow with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Testing user workflow activate with MCMS...")
		activateInput := UserWorkflowActivateInput{
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				MCMSConfig: &contracts.MCMSConfig{
					MinDelay: 1 * time.Second,
					TimelockQualifierPerChain: map[uint64]string{
						fixture.selector: "",
					},
				},
			},
			WorkflowID: testWorkflowID,
			DonFamily:  testDONFamily,
		}
		activateChangeset := UserWorkflowActivate{}
		err := activateChangeset.VerifyPreconditions(fixture.rt.Environment(), activateInput)
		require.NoError(t, err, "user workflow activate with MCMS preconditions should pass")
		csOutput, err := activateChangeset.Apply(fixture.rt.Environment(), activateInput)
		require.NoError(t, err, "user workflow activate with MCMS apply should pass")
		assert.NotNil(t, csOutput.Reports, "user workflow activate apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from user workflow activate")
		t.Log("User workflow activated with MCMS successfully")
	})

	t.Run("delete workflow with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Testing user workflow delete with MCMS...")
		deleteInput := UserWorkflowDeleteInput{
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				MCMSConfig: &contracts.MCMSConfig{
					MinDelay: 1 * time.Second,
					TimelockQualifierPerChain: map[uint64]string{
						fixture.selector: "",
					},
				},
			},
			WorkflowID: testWorkflowID,
		}
		deleteChangeset := UserWorkflowDelete{}
		err := deleteChangeset.VerifyPreconditions(fixture.rt.Environment(), deleteInput)
		require.NoError(t, err, "user workflow delete with MCMS preconditions should pass")
		csOutput, err := deleteChangeset.Apply(fixture.rt.Environment(), deleteInput)
		require.NoError(t, err, "user workflow delete with MCMS apply should pass")
		assert.NotNil(t, csOutput.Reports, "user workflow delete apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from user workflow delete")
		t.Log("User workflow deleted with MCMS successfully")
	})

	t.Run("unlink owner with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		chain := fixture.rt.Environment().BlockChains.EVMChains()[fixture.selector]
		deployerKey := chain.DeployerKey

		t.Log("Testing unlink owner with MCMS...")
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
			Address:           deployerKey.From,
			ValidityTimestamp: validity,
			Signature:         common.Bytes2Hex(signature),
			CommonWorkflowInput: CommonWorkflowInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				MCMSConfig: &contracts.MCMSConfig{
					MinDelay: 1 * time.Second,
					TimelockQualifierPerChain: map[uint64]string{
						fixture.selector: "",
					},
				},
			},
		}
		unlinkOwnerChangeset := UserUnlinkOwner{}
		err := unlinkOwnerChangeset.VerifyPreconditions(fixture.rt.Environment(), unlinkOwnerInput)
		require.NoError(t, err, "unlink owner with MCMS preconditions should pass")
		_, err = unlinkOwnerChangeset.Apply(fixture.rt.Environment(), unlinkOwnerInput)
		require.NoError(t, err, "unlink owner with MCMS apply should pass")
	})
}

// Generate ownersip proof based on test data and sign using private key of an allowed test signer.
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
