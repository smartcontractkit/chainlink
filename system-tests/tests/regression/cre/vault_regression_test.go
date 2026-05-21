package cre

import (
	"encoding/json"
	"net/url"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"

	smoke_vault_tests "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

// Test_CRE_V2_Vault_IdentifierValidation_Regression tests that the vault gateway correctly rejects
// requests with invalid identifiers (those containing characters outside alphanumeric+underscore set).
// This is a regression test for error cases not covered by smoke tests.
func Test_CRE_V2_Vault_IdentifierValidation_Regression(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
	ExecuteVaultIdentifierValidationRegressionTest(t, testEnv)
}

// ExecuteVaultIdentifierValidationRegressionTest verifies that the gateway rejects requests whose
// secret identifiers contain characters outside the allowed alphanumeric+underscore set.
// All four management request types (create, update, delete, list) are exercised for invalid key,
// invalid namespace, and invalid owner. This regression test complements the positive-path coverage
// provided by smoke tests.
func ExecuteVaultIdentifierValidationRegressionTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	owner := sc.MustGetRootKeyAddress().Hex()

	wfRegAddr := crecontracts.MustGetAddressFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.Blockchains[0].ChainSelector(),
		keystone_changeset.WorkflowRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()],
		"",
	)
	wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegAddr), sc.Client)
	require.NoError(t, err)
	require.NoError(t, creworkflow.LinkOwner(sc, common.HexToAddress(wfRegAddr), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()]))

	// Get gateway URL
	gatewayURL := mustVaultGatewayURL(t, testEnv)

	// Fetch and parse vault public key
	vaultPublicKeyStr := smoke_vault_tests.FetchVaultPublicKey(t, gatewayURL.String())
	vaultParsedPublicKey := smoke_vault_tests.MustVaultPublicKey(t, vaultPublicKeyStr)
	enc, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-basic", vaultParsedPublicKey, sc.MustGetRootKeyAddress())
	require.NoError(t, err)

	// Test identifier validation via gateway requests
	executeVaultSecretsIdentifierValidationTest(t, enc, owner, gatewayURL.String(), sc, wfReg)

	framework.L.Info().Msg("Vault identifier validation regression test completed")
}

// executeVaultSecretsIdentifierValidationTest verifies that the gateway rejects requests whose
// secret identifiers contain characters outside the allowed alphanumeric+underscore set.
// All four management request types (create, update, delete, list) are exercised for invalid key,
// invalid namespace, and invalid owner.
func executeVaultSecretsIdentifierValidationTest(t *testing.T, encryptedSecret string, owner, gatewayURL string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	t.Helper()

	const (
		validKey         = "validkey"
		invalidKey       = "invalid-key-with-hyphens"   // hyphen not in [a-zA-Z0-9_]
		invalidOwner     = "invalid-owner-with-hyphens" // hyphen not in [a-zA-Z0-9_]
		validNamespace   = "main"
		invalidNamespace = "invalid-namespace-hyphens" // hyphen not in [a-zA-Z0-9_]
	)

	sendWriteAndAssert := func(t *testing.T, method, caseName string, secret *vault_helpers.EncryptedSecret) {
		t.Helper()
		uniqueRequestID := uuid.New().String()
		var body []byte
		var err error
		switch method {
		case vaulttypes.MethodSecretsCreate:
			body, err = json.Marshal(vault_helpers.CreateSecretsRequest{RequestId: uniqueRequestID, EncryptedSecrets: []*vault_helpers.EncryptedSecret{secret}})
		case vaulttypes.MethodSecretsUpdate:
			body, err = json.Marshal(vault_helpers.UpdateSecretsRequest{RequestId: uniqueRequestID, EncryptedSecrets: []*vault_helpers.EncryptedSecret{secret}})
		case vaulttypes.MethodSecretsDelete:
			body, err = json.Marshal(vault_helpers.DeleteSecretsRequest{RequestId: uniqueRequestID, Ids: []*vault_helpers.SecretIdentifier{secret.Id}})
		}
		require.NoError(t, err)
		bodyJSON := json.RawMessage(body)
		req := jsonrpc.Request[json.RawMessage]{Version: jsonrpc.JsonRpcVersion, ID: uniqueRequestID, Method: method, Params: &bodyJSON}
		smoke_vault_tests.AllowlistRequest(t, owner, req, sethClient, wfRegistryContract)
		reqBody, err := json.Marshal(req)
		require.NoError(t, err)
		_, respBody := smoke_vault_tests.SendVaultRequestToGateway(t, gatewayURL, reqBody)
		require.Contains(t, string(respBody), "alphanumeric", "[%s] expected alphanumeric rejection for %s", method, caseName)
		framework.L.Info().Msgf("[%s] %s correctly rejected: %s", method, caseName, string(respBody))
	}

	type writeCase struct {
		name         string
		key, own, ns string
	}
	writeCases := []writeCase{
		{"invalid key", invalidKey, owner, validNamespace},
		{"invalid namespace", validKey, owner, invalidNamespace},
		{"invalid owner", validKey, invalidOwner, validNamespace},
	}

	for _, op := range []string{vaulttypes.MethodSecretsCreate, vaulttypes.MethodSecretsUpdate, vaulttypes.MethodSecretsDelete} {
		framework.L.Info().Msgf("Testing identifier validation for %s request...", op)
		for _, tc := range writeCases {
			sendWriteAndAssert(t, op, tc.name, &vault_helpers.EncryptedSecret{
				Id:             &vault_helpers.SecretIdentifier{Key: tc.key, Owner: tc.own, Namespace: tc.ns},
				EncryptedValue: encryptedSecret,
			})
		}
	}

	framework.L.Info().Msg("Testing identifier validation for list request...")
	uniqueRequestID := uuid.New().String()
	body, err := json.Marshal(vault_helpers.ListSecretIdentifiersRequest{RequestId: uniqueRequestID, Owner: owner, Namespace: invalidNamespace})
	require.NoError(t, err)
	bodyJSON := json.RawMessage(body)
	req := jsonrpc.Request[json.RawMessage]{Version: jsonrpc.JsonRpcVersion, ID: uniqueRequestID, Method: vaulttypes.MethodSecretsList, Params: &bodyJSON}
	smoke_vault_tests.AllowlistRequest(t, owner, req, sethClient, wfRegistryContract)
	reqBody, err := json.Marshal(req)
	require.NoError(t, err)
	_, respBody := smoke_vault_tests.SendVaultRequestToGateway(t, gatewayURL, reqBody)
	require.Contains(t, string(respBody), "alphanumeric", "[list] expected alphanumeric rejection for %s", "invalid namespace")
	framework.L.Info().Msgf("[list] %s correctly rejected: %s", "invalid namespace", string(respBody))

	framework.L.Info().Msg("All identifier validation checks passed")
}

func mustVaultGatewayURL(t *testing.T, testEnv *ttypes.TestEnvironment) *url.URL {
	t.Helper()

	framework.L.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, testEnv.Dons.GatewayConnectors.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.Protocol + "://" + testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.Host + ":" + strconv.Itoa(testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.ExternalPort) + testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.Path)
	require.NoError(t, err, "failed to parse gateway URL")
	framework.L.Info().Msgf("Gateway URL: %s", gatewayURL.String())
	return gatewayURL
}
