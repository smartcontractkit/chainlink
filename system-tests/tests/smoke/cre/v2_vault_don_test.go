package cre

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	ctfblockchain "github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	vaultsecret_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/vaultsecret/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crevault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/vault"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func ExecuteVaultTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	var testLogger = framework.L

	testLogger.Info().Msgf("Ensuring DKG result packages are present...")
	require.Eventually(t, func() bool {
		for _, nodeSet := range testEnv.Config.NodeSets {
			if slices.Contains(nodeSet.Capabilities, cre.VaultCapability) {
				for i, node := range nodeSet.NodeSpecs {
					if !slices.Contains(node.Roles, cre.BootstrapNode) {
						packageCount, err := vault.GetResultPackageCount(t.Context(), i, nodeSet.DbInput.Port)
						if err != nil || packageCount != 1 {
							return false
						}
					}
				}
				return true
			}
		}
		return false
	}, time.Second*300, time.Second*5)

	testLogger.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, testEnv.Dons.GatewayConnectors.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.Protocol + "://" + testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.Host + ":" + strconv.Itoa(testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.ExternalPort) + testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.Path)
	require.NoError(t, err, "failed to parse gateway URL")
	testLogger.Info().Msgf("Gateway URL: %s", gatewayURL.String())

	vaultPublicKey := FetchVaultPublicKey(t, gatewayURL.String())
	updateVaultCapabilityConfigInRegistry(t, testEnv, vaultPublicKey)

	gwURL := gatewayURL.String()

	t.Run("basic_crud", func(t *testing.T) {
		if parallelEnabled && fanoutEnabled {
			t.Parallel()
		}
		subEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, testEnv.TestConfig)
		sc := subEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
		owner := sc.MustGetRootKeyAddress().Hex()
		wfRegAddr := crecontracts.MustGetAddressFromDataStore(subEnv.CreEnvironment.CldfEnvironment.DataStore, subEnv.CreEnvironment.Blockchains[0].ChainSelector(), keystone_changeset.WorkflowRegistry.String(), subEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
		wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegAddr), sc.Client)
		require.NoError(t, err)
		require.NoError(t, creworkflow.LinkOwner(sc, common.HexToAddress(wfRegAddr), subEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()]))
		secretID := strconv.Itoa(rand.Intn(10000))
		enc, err := crevault.EncryptSecret("secret-basic", vaultPublicKey, sc.MustGetRootKeyAddress())
		require.NoError(t, err)
		ulCh := make(chan *workflowevents.UserLogs, 1000)
		bmCh := make(chan *commonevents.BaseMessage, 1000)
		sink := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, ulCh, bmCh))
		t.Cleanup(func() { sink.Shutdown(t.Context()); close(ulCh); close(bmCh) })
		namespaces := []string{"main", "alt"}

		executeVaultSecretsCreateTest(t, enc, secretID, owner, gwURL, namespaces, sc, wfReg)
		executeVaultSecretsGetViaWorkflowTest(t, subEnv, "bget1", secretID, "main", ulCh, bmCh)
		executeVaultSecretsGetViaWorkflowTest(t, subEnv, "bgeta1", secretID, "alt", ulCh, bmCh)
		executeVaultSecretsUpdateTest(t, enc, secretID, owner, gwURL, namespaces, sc, wfReg)
		executeVaultSecretsGetViaWorkflowTest(t, subEnv, "bget2", secretID, "main", ulCh, bmCh)
		executeVaultSecretsGetViaWorkflowTest(t, subEnv, "bgeta2", secretID, "alt", ulCh, bmCh)
		executeVaultSecretsListTest(t, secretID, owner, gwURL, "main", sc, wfReg)
		executeVaultSecretsListTest(t, secretID, owner, gwURL, "alt", sc, wfReg)
		executeVaultSecretsDeleteTest(t, secretID, owner, gwURL, []string{"main"}, sc, wfReg)
		executeVaultSecretsGetNotFoundViaWorkflowTest(t, subEnv, "bdel1", secretID, "main", ulCh, bmCh)
		executeVaultSecretsGetViaWorkflowTest(t, subEnv, "bgeta3", secretID, "alt", ulCh, bmCh)
		executeVaultSecretsDeleteTest(t, secretID, owner, gwURL, []string{"alt"}, sc, wfReg)
		executeVaultSecretsGetNotFoundViaWorkflowTest(t, subEnv, "bdela1", secretID, "alt", ulCh, bmCh)
	})
}

func executeVaultSecretsCreateTest(t *testing.T, encryptedSecret, secretID, owner, gatewayURL string, namespaces []string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	framework.L.Info().Msgf("Creating secrets (namespaces=%v)...", namespaces)

	uniqueRequestID := uuid.New().String()

	encryptedSecrets := make([]*vault_helpers.EncryptedSecret, 0, len(namespaces))
	for _, namespace := range namespaces {
		encryptedSecrets = append(encryptedSecrets, &vault_helpers.EncryptedSecret{
			Id: &vault_helpers.SecretIdentifier{
				Key:       secretID,
				Owner:     owner,
				Namespace: namespace,
			},
			EncryptedValue: encryptedSecret,
		})
	}

	secretsCreateRequest := vault_helpers.CreateSecretsRequest{
		RequestId:        uniqueRequestID,
		EncryptedSecrets: encryptedSecrets,
	}
	secretsCreateRequestBody, err := json.Marshal(secretsCreateRequest) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsCreateRequestBodyJSON := json.RawMessage(secretsCreateRequestBody)
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsCreate,
		Params:  &secretsCreateRequestBodyJSON,
	}
	allowlistRequest(t, owner, jsonRequest, sethClient, wfRegistryContract)

	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}
	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsCreate, jsonResponse.Method)

	signedOCRResponse := jsonResponse.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponse.String())

	// TODO: Verify the authenticity of this signed report, by ensuring that the signatures indeed match the payload
	createSecretsResponse := vault_helpers.CreateSecretsResponse{}
	err = protojson.Unmarshal(signedOCRResponse.Payload, &createSecretsResponse)
	require.NoError(t, err, "failed to decode payload into CreateSecretsResponse proto")
	framework.L.Info().Msgf("CreateSecretsResponse decoded as: %s", createSecretsResponse.String())

	require.Len(t, createSecretsResponse.Responses, len(namespaces), "Expected one item in the response per namespace")
	respByNs := make(map[string]*vault_helpers.CreateSecretResponse, len(namespaces))
	for _, r := range createSecretsResponse.GetResponses() {
		respByNs[r.GetId().GetNamespace()] = r
	}
	for _, namespace := range namespaces {
		result, ok := respByNs[namespace]
		require.True(t, ok, "missing response for namespace %s", namespace)
		require.Empty(t, result.GetError())
		require.Equal(t, secretID, result.GetId().Key)
		require.Equal(t, owner, result.GetId().Owner)
	}

	framework.L.Info().Msgf("Secrets created successfully (namespaces=%v)", namespaces)
}

func executeVaultSecretsGetViaWorkflowTest(
	t *testing.T, testEnv *ttypes.TestEnvironment,
	workflowBaseName, secretKey, secretNamespace string,
	userLogsCh chan *workflowevents.UserLogs, baseMessageCh chan *commonevents.BaseMessage,
) {
	testLogger := framework.L
	testLogger.Info().Msgf("Verifying secret retrieval via workflow (key=%s, namespace=%s)...", secretKey, secretNamespace)

	workflowName := t_helpers.UniqueWorkflowName(testEnv, workflowBaseName)
	cfg := &vaultsecret_config.Config{
		SecretKey:       secretKey,
		SecretNamespace: secretNamespace,
	}
	const workflowFileLocation = "./vaultsecret/main.go"
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, cfg, workflowFileLocation)

	expectedLog := "Vault secret retrieved successfully via workflow"
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectedLog, 4*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))
	testLogger.Info().Msg("Vault secret get via workflow test completed")
}

func executeVaultSecretsGetNotFoundViaWorkflowTest(
	t *testing.T, testEnv *ttypes.TestEnvironment,
	workflowBaseName, secretKey, secretNamespace string,
	userLogsCh chan *workflowevents.UserLogs, baseMessageCh chan *commonevents.BaseMessage,
) {
	testLogger := framework.L
	testLogger.Info().Msgf("Verifying secret is NOT retrievable via workflow after deletion (key=%s, namespace=%s)...", secretKey, secretNamespace)

	workflowName := t_helpers.UniqueWorkflowName(testEnv, workflowBaseName)
	cfg := &vaultsecret_config.Config{
		SecretKey:       secretKey,
		SecretNamespace: secretNamespace,
		ExpectNotFound:  true,
	}
	const workflowFileLocation = "./vaultsecret/main.go"
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, cfg, workflowFileLocation)

	expectedLog := "Vault secret correctly not found after deletion"
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectedLog, 4*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))
	testLogger.Info().Msg("Vault secret not-found via workflow test completed")
}

func executeVaultSecretsUpdateTest(t *testing.T, encryptedSecret, secretID, owner, gatewayURL string, namespaces []string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	framework.L.Info().Msgf("Updating secrets (namespaces=%v)...", namespaces)
	uniqueRequestID := uuid.New().String()

	encryptedSecrets := make([]*vault_helpers.EncryptedSecret, 0, len(namespaces)+1)
	for _, namespace := range namespaces {
		encryptedSecrets = append(encryptedSecrets, &vault_helpers.EncryptedSecret{
			Id: &vault_helpers.SecretIdentifier{
				Key:       secretID,
				Owner:     owner,
				Namespace: namespace,
			},
			EncryptedValue: encryptedSecret,
		})
	}
	encryptedSecrets = append(encryptedSecrets, &vault_helpers.EncryptedSecret{
		Id: &vault_helpers.SecretIdentifier{
			Key:       "invalid",
			Owner:     owner,
			Namespace: namespaces[0],
		},
		EncryptedValue: encryptedSecret,
	})

	secretsUpdateRequest := vault_helpers.UpdateSecretsRequest{
		RequestId:        uniqueRequestID,
		EncryptedSecrets: encryptedSecrets,
	}
	secretsUpdateRequestBody, err := json.Marshal(secretsUpdateRequest) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsUpdateRequestBodyJSON := json.RawMessage(secretsUpdateRequestBody)
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsUpdate,
		Params:  &secretsUpdateRequestBodyJSON,
	}
	allowlistRequest(t, owner, jsonRequest, sethClient, wfRegistryContract)

	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}

	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsUpdate, jsonResponse.Method)

	signedOCRResponse := jsonResponse.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponse.String())

	// TODO: Verify the authenticity of this signed report, by ensuring that the signatures indeed match the payload

	updateSecretsResponse := vault_helpers.UpdateSecretsResponse{}
	err = protojson.Unmarshal(signedOCRResponse.Payload, &updateSecretsResponse)
	require.NoError(t, err, "failed to decode payload into UpdateSecretsResponse proto")
	framework.L.Info().Msgf("UpdateSecretsResponse decoded as: %s", updateSecretsResponse.String())

	require.Len(t, updateSecretsResponse.Responses, len(namespaces)+1, "Expected one updated item per namespace plus one invalid item")
	var foundInvalid bool
	updateRespByNs := make(map[string]*vault_helpers.UpdateSecretResponse, len(namespaces))
	for _, r := range updateSecretsResponse.GetResponses() {
		if r.GetId().GetKey() == "invalid" {
			require.Contains(t, r.Error, "key does not exist")
			foundInvalid = true
			continue
		}
		updateRespByNs[r.GetId().GetNamespace()] = r
	}
	require.True(t, foundInvalid, "expected an error response for the 'invalid' key")
	for _, namespace := range namespaces {
		result, ok := updateRespByNs[namespace]
		require.True(t, ok, "missing update response for namespace %s", namespace)
		require.Empty(t, result.GetError())
		require.Equal(t, secretID, result.GetId().Key)
		require.Equal(t, owner, result.GetId().Owner)
	}

	framework.L.Info().Msgf("Secrets updated successfully (namespaces=%v)", namespaces)
}

func executeVaultSecretsListTest(t *testing.T, secretID, owner, gatewayURL, namespace string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	framework.L.Info().Msgf("Listing secrets (namespace=%s)...", namespace)
	uniqueRequestID := uuid.New().String()
	secretsListRequest := vault_helpers.ListSecretIdentifiersRequest{
		RequestId: uniqueRequestID,
		Owner:     owner,
		Namespace: namespace,
	}
	secretsListRequestBody, err := json.Marshal(secretsListRequest) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsUpdateRequestBodyJSON := json.RawMessage(secretsListRequestBody)
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsList,
		Params:  &secretsUpdateRequestBodyJSON,
	}
	allowlistRequest(t, owner, jsonRequest, sethClient, wfRegistryContract)

	// Ensure that multiple requests can be allowlisted
	uniqueRequestIDTwo := uuid.New().String()
	secretsListRequestTwo := vault_helpers.ListSecretIdentifiersRequest{
		RequestId: uniqueRequestIDTwo,
		Owner:     owner,
		Namespace: namespace,
	}
	secretsListRequestBodyTwo, err := json.Marshal(secretsListRequestTwo) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsUpdateRequestBodyJSONTwo := json.RawMessage(secretsListRequestBodyTwo)
	jsonRequestTwo := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestIDTwo,
		Method:  vaulttypes.MethodSecretsList,
		Params:  &secretsUpdateRequestBodyJSONTwo,
	}
	allowlistRequest(t, owner, jsonRequestTwo, sethClient, wfRegistryContract)

	// Request 1
	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}

	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsList, jsonResponse.Method)

	signedOCRResponse := jsonResponse.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponse.String())

	// Request 2
	requestBodyTwo, err := json.Marshal(jsonRequestTwo)
	require.NoError(t, err, "failed to marshal secrets request")
	statusCodeTwo, httpResponseBodyTwo := sendVaultRequestToGateway(t, gatewayURL, requestBodyTwo)
	require.Equal(t, http.StatusOK, statusCodeTwo, "Gateway endpoint should respond with 200 OK")
	var jsonResponseTwo jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBodyTwo, &jsonResponseTwo)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponseTwo)
	if jsonResponseTwo.Error != nil {
		require.Empty(t, jsonResponseTwo.Error.Error())
	}
	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponseTwo.Version)
	require.Equal(t, uniqueRequestIDTwo, jsonResponseTwo.ID)
	require.Equal(t, vaulttypes.MethodSecretsList, jsonResponseTwo.Method)
	signedOCRResponseTwo := jsonResponseTwo.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponseTwo.String())

	// TODO: Verify the authenticity of this signed report, by ensuring that the signatures indeed match the payload

	listSecretsResponse := vault_helpers.ListSecretIdentifiersResponse{}
	err = protojson.Unmarshal(signedOCRResponse.Payload, &listSecretsResponse)
	require.NoError(t, err, "failed to decode payload into ListSecretIdentifiersResponse proto")
	framework.L.Info().Msgf("ListSecretIdentifiersResponse decoded as: %s", listSecretsResponse.String())

	require.True(t, listSecretsResponse.Success, err)
	require.GreaterOrEqual(t, len(listSecretsResponse.Identifiers), 1, "Expected at least one item in the response")
	var keys = make([]string, 0, len(listSecretsResponse.Identifiers))
	for _, identifier := range listSecretsResponse.Identifiers {
		keys = append(keys, identifier.Key)
		require.Equal(t, owner, identifier.Owner)
		require.Equal(t, namespace, identifier.Namespace)
	}
	require.Contains(t, keys, secretID)
	framework.L.Info().Msgf("Secrets listed successfully (namespace=%s)", namespace)
}

func executeVaultSecretsDeleteTest(t *testing.T, secretID, owner, gatewayURL string, namespaces []string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	framework.L.Info().Msgf("Deleting secrets (namespaces=%v)...", namespaces)
	uniqueRequestID := uuid.New().String()

	deleteIDs := make([]*vault_helpers.SecretIdentifier, 0, len(namespaces)+1)
	for _, namespace := range namespaces {
		deleteIDs = append(deleteIDs, &vault_helpers.SecretIdentifier{
			Key:       secretID,
			Owner:     owner,
			Namespace: namespace,
		})
	}
	deleteIDs = append(deleteIDs, &vault_helpers.SecretIdentifier{
		Key:       "invalid",
		Owner:     owner,
		Namespace: namespaces[0],
	})

	secretsDeleteRequest := vault_helpers.DeleteSecretsRequest{
		RequestId: uniqueRequestID,
		Ids:       deleteIDs,
	}
	secretsDeleteRequestBody, err := json.Marshal(secretsDeleteRequest) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsDeleteRequestBodyJSON := json.RawMessage(secretsDeleteRequestBody)
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsDelete,
		Params:  &secretsDeleteRequestBodyJSON,
	}
	allowlistRequest(t, owner, jsonRequest, sethClient, wfRegistryContract)

	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")
	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}

	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsDelete, jsonResponse.Method)

	signedOCRResponse := jsonResponse.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponse.String())

	// TODO: Verify the authenticity of this signed report, by ensuring that the signatures indeed match the payload

	deleteSecretsResponse := vault_helpers.DeleteSecretsResponse{}
	err = protojson.Unmarshal(signedOCRResponse.Payload, &deleteSecretsResponse)
	require.NoError(t, err, "failed to decode payload into DeleteSecretResponse proto")
	framework.L.Info().Msgf("DeleteSecretResponse decoded as: %s", deleteSecretsResponse.String())

	require.Len(t, deleteSecretsResponse.Responses, len(namespaces)+1, "Expected one deleted item per namespace plus one invalid item")
	var foundDeleteInvalid bool
	deleteRespByNs := make(map[string]*vault_helpers.DeleteSecretResponse, len(namespaces))
	for _, r := range deleteSecretsResponse.GetResponses() {
		if r.GetId().GetKey() == "invalid" {
			require.Contains(t, r.Error, "key does not exist")
			foundDeleteInvalid = true
			continue
		}
		deleteRespByNs[r.GetId().GetNamespace()] = r
	}
	require.True(t, foundDeleteInvalid, "expected an error response for the 'invalid' key")
	for _, namespace := range namespaces {
		result, ok := deleteRespByNs[namespace]
		require.True(t, ok, "missing delete response for namespace %s", namespace)
		require.True(t, result.Success, result.Error)
		require.Equal(t, owner, result.Id.Owner)
		require.Equal(t, secretID, result.Id.Key)
	}

	framework.L.Info().Msgf("Secrets deleted successfully (namespaces=%v)", namespaces)
}

// updateVaultCapabilityConfigInRegistry updates the on-chain capabilities registry
// so that the vault@1.0.0 capability config includes DefaultConfig with VaultPublicKey
// and Threshold. This is required for workflows that call runtime.GetSecret().
// Uses the original deployer key (not per-test key) since the registry is owned by the deployer.
func updateVaultCapabilityConfigInRegistry(t *testing.T, testEnv *ttypes.TestEnvironment, vaultPublicKey string) {
	t.Helper()
	testLogger := framework.L
	testLogger.Info().Msg("Updating vault capability config in capabilities registry with VaultPublicKey...")

	capRegAddr := crecontracts.MustGetAddressFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.RegistryChainSelector,
		keystone_changeset.CapabilitiesRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()],
		"",
	)

	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0])
	sethClient := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient

	deployerClient, err := seth.NewClientBuilder().
		WithRpcUrl(sethClient.URL).
		WithPrivateKeys([]string{ctfblockchain.DefaultAnvilPrivateKey}).
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	require.NoError(t, err, "failed to create deployer seth client")

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(capRegAddr), deployerClient.Client,
	)
	require.NoError(t, err, "failed to create capabilities registry wrapper")

	allDONs, err := capReg.GetDONs(&bind.CallOpts{}, big.NewInt(0), big.NewInt(100))
	require.NoError(t, err, "failed to get DONs from registry")

	var don *capabilities_registry_v2.CapabilitiesRegistryDONInfo
	for i := range allDONs {
		for _, cc := range allDONs[i].CapabilityConfigurations {
			if cc.CapabilityId == "vault@1.0.0" {
				don = &allDONs[i]
				break
			}
		}
		if don != nil {
			break
		}
	}
	require.NotNil(t, don, "could not find a DON with vault@1.0.0 capability in the registry")
	testLogger.Info().Msgf("Found vault capability on DON %q (ID=%d)", don.Name, don.Id)

	newConfigs := make([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, 0, len(don.CapabilityConfigurations))
	for _, cc := range don.CapabilityConfigurations {
		if cc.CapabilityId == "vault@1.0.0" {
			existingConfig := &capabilitiespb.CapabilityConfig{}
			if len(cc.Config) > 0 {
				require.NoError(t, proto.Unmarshal(cc.Config, existingConfig), "failed to unmarshal existing vault capability config")
			}

			vaultCfg := map[string]interface{}{
				"VaultPublicKey": vaultPublicKey,
				"Threshold":      1,
			}
			valueMap, wrapErr := values.WrapMap(vaultCfg)
			require.NoError(t, wrapErr, "failed to wrap vault config values")

			existingConfig.DefaultConfig = values.ProtoMap(valueMap)

			configBytes, marshalErr := proto.Marshal(existingConfig)
			require.NoError(t, marshalErr, "failed to marshal updated vault capability config")

			cc.Config = configBytes
			testLogger.Info().Msg("Injected VaultPublicKey and Threshold into vault@1.0.0 capability config")
		}
		newConfigs = append(newConfigs, cc)
	}

	updateParams := capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams{
		Name:                     don.Name,
		Nodes:                    don.NodeP2PIds,
		CapabilityConfigurations: newConfigs,
		IsPublic:                 don.IsPublic,
		F:                        don.F,
		Config:                   don.Config,
	}

	_, err = deployerClient.Decode(capReg.UpdateDONByName(deployerClient.NewTXOpts(), don.Name, updateParams))
	require.NoError(t, err, "UpdateDONByName tx failed")

	testLogger.Info().Msg("Waiting for registry syncer to propagate the on-chain config change...")
	time.Sleep(15 * time.Second) // registry syncer polls every 12s; one tick + margin
}

func allowlistRequest(t *testing.T, owner string, request jsonrpc.Request[json.RawMessage], sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	requestDigest, err := request.Digest()
	require.NoError(t, err, "failed to get digest for request")
	requestDigestBytes, err := hex.DecodeString(requestDigest)
	require.NoError(t, err, "failed to decode digest")
	reqDigestBytes := [32]byte(requestDigestBytes)
	_, err = wfRegistryContract.AllowlistRequest(sethClient.NewTXOpts(), reqDigestBytes, uint32(time.Now().Add(1*time.Hour).Unix())) //nolint:gosec // disable G115
	require.NoError(t, err, "failed to allowlist request")

	framework.L.Info().Msgf("Allowlisting request digest at contract %s, for owner: %s, digestHexStr: %s", wfRegistryContract.Address().Hex(), owner, requestDigest)
	allowedList, err := wfRegistryContract.GetAllowlistedRequests(&bind.CallOpts{}, big.NewInt(0), big.NewInt(100))
	require.NoError(t, err, "failed to validate allowlisted request")
	for _, req := range allowedList {
		if req.RequestDigest == reqDigestBytes {
			framework.L.Info().Msgf("Request digest found in allowlist")
		}
		framework.L.Info().Msgf("Allowlisted request digestHexStr: %s, owner: %s, expiry: %d", hex.EncodeToString(req.RequestDigest[:]), req.Owner.Hex(), req.ExpiryTimestamp)
	}
}
