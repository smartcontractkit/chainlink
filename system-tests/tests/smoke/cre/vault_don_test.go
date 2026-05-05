package cre

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"

	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func uniqueVaultSecretID(prefix string) string {
	return prefix + strings.ReplaceAll(uuid.NewString(), "-", "")
}

func ExecuteVaultAllowListBasedTests(t *testing.T, fixture *vaultScenarioFixture, testEnv *ttypes.TestEnvironment) {
	var testLogger = framework.L
	linkingService := fixture.LinkingService

	gwURL := fixture.GatewayURL.String()
	vaultPublicKey := fixture.VaultPublicKey

	t.Run("allowlist_crud_with_workflow_owner_identity", func(t *testing.T) {
		sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
		workflowOwnerAddress := sc.MustGetRootKeyAddress()
		owner := workflowOwnerAddress.Hex()
		expectedResponseOwner := owner
		orgID := ""
		orgIDAsSecretOwnerEnabled := isVaultJWTAuthEnabledTopology(testEnv.TestConfig.EnvironmentConfigPath)
		if linkingService != nil {
			orgID = "org" + strings.ReplaceAll(uuid.NewString(), "-", "")
			linkingService.SetOwnerOrg(owner, orgID)
			if orgIDAsSecretOwnerEnabled {
				expectedResponseOwner = orgID
			}
		}
		if orgIDAsSecretOwnerEnabled {
			require.NotEmpty(t, orgID, "JWT auth enabled topology must link the workflow owner to an org ID")
		}
		wfRegAddr := crecontracts.MustGetAddressFromDataStore(testEnv.CreEnvironment.CldfEnvironment.DataStore, testEnv.CreEnvironment.Blockchains[0].ChainSelector(), keystone_changeset.WorkflowRegistry.String(), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
		wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegAddr), sc.Client)
		require.NoError(t, err)
		requireVaultLinkOwner(t, sc, common.HexToAddress(wfRegAddr), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()])
		vaultParsedPublicKey := mustVaultPublicKey(t, vaultPublicKey)
		secretID := uniqueVaultSecretID("allowlist")
		createValue := "secret-basic-create"
		updateValue := "secret-basic-update"
		createEnc, err := vaultutils.EncryptSecretWithWorkflowOwner(createValue, vaultParsedPublicKey, workflowOwnerAddress)
		require.NoError(t, err)
		updateEnc, err := vaultutils.EncryptSecretWithWorkflowOwner(updateValue, vaultParsedPublicKey, workflowOwnerAddress)
		require.NoError(t, err)
		ulCh := make(chan *workflowevents.UserLogs, 1000)
		bmCh := make(chan *commonevents.BaseMessage, 1000)
		sink := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, ulCh, bmCh))
		t.Cleanup(func() {
			// can't use t.Context() here because it will have been cancelled before the cleanup function is called
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			t_helpers.ShutdownChipSinkWithDrain(ctx, sink, ulCh, bmCh)
		})
		namespaces := []string{"main", "alt"}

		executeVaultAllowListSecretsCreateTest(t, createEnc, secretID, owner, expectedResponseOwner, gwURL, namespaces, sc, wfReg)
		var orgIDLabelSecretID string
		var orgIDLabelCreateValue string
		if orgIDAsSecretOwnerEnabled {
			orgIDLabelSecretID = secretID + "orgidlabel"
			orgIDLabelCreateValue = "secret-basic-create-org-id-label"
			orgIDLabelCreateEnc, orgErr := vaultutils.EncryptSecretWithOrgID(orgIDLabelCreateValue, vaultParsedPublicKey, orgID)
			require.NoError(t, orgErr)
			executeVaultAllowListSecretsCreateTest(t, orgIDLabelCreateEnc, orgIDLabelSecretID, owner, expectedResponseOwner, gwURL, namespaces, sc, wfReg)
		}
		executeVaultSecretsUpdateTest(t, updateEnc, secretID, owner, expectedResponseOwner, gwURL, namespaces, sc, wfReg)
		executeVaultSecretsListTest(t, secretID, owner, expectedResponseOwner, gwURL, "main", sc, wfReg)
		executeVaultSecretsListTest(t, secretID, owner, expectedResponseOwner, gwURL, "alt", sc, wfReg)
		updatedChecks := []vaultWorkflowCheck{
			{Name: "allowlist-main-updated", SecretKey: secretID, SecretNamespace: "main", ExpectedValue: updateValue},
			{Name: "allowlist-alt-updated", SecretKey: secretID, SecretNamespace: "alt", ExpectedValue: updateValue},
		}
		finalChecks := []vaultWorkflowCheck{
			{Name: "allowlist-main-not-found", SecretKey: secretID, SecretNamespace: "main", ExpectNotFound: true},
			{Name: "allowlist-alt-updated", SecretKey: secretID, SecretNamespace: "alt", ExpectedValue: updateValue},
		}
		if orgIDAsSecretOwnerEnabled {
			orgIDChecks := []vaultWorkflowCheck{
				{Name: "allowlist-org-id-label-main", SecretKey: orgIDLabelSecretID, SecretNamespace: "main", ExpectedValue: orgIDLabelCreateValue},
				{Name: "allowlist-org-id-label-alt", SecretKey: orgIDLabelSecretID, SecretNamespace: "alt", ExpectedValue: orgIDLabelCreateValue},
			}
			updatedChecks = append(updatedChecks, orgIDChecks...)
			finalChecks = append(finalChecks, orgIDChecks...)
		}
		workflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, "allowlist-lifecycle", []vaultWorkflowPhase{
			{Name: "allowlist-updated", Checks: updatedChecks},
			{Name: "allowlist-final-verify", Checks: finalChecks},
		})
		waitForVaultWorkflowPhase(t, workflowID, "allowlist-updated", ulCh, bmCh)
		executeVaultSecretsDeleteTest(t, secretID, owner, expectedResponseOwner, gwURL, []string{"main"}, sc, wfReg)
		waitForVaultWorkflowPhase(t, workflowID, "allowlist-final-verify", ulCh, bmCh)
		executeVaultSecretsDeleteTest(t, secretID, owner, expectedResponseOwner, gwURL, []string{"alt"}, sc, wfReg)
		if orgIDAsSecretOwnerEnabled {
			executeVaultSecretsDeleteTest(t, orgIDLabelSecretID, owner, expectedResponseOwner, gwURL, namespaces, sc, wfReg)
		}
	})

	if !isVaultJWTAuthEnabledTopology(testEnv.TestConfig.EnvironmentConfigPath) {
		return
	}

	t.Run("allowlist_crud_with_org_id_identity", func(t *testing.T) {
		require.NotNil(t, linkingService, "JWT auth enabled topology must include a linking service")

		sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
		workflowOwnerAddress := sc.MustGetRootKeyAddress()
		workflowOwner := workflowOwnerAddress.Hex()
		orgID := "org" + strings.ReplaceAll(uuid.NewString(), "-", "")
		linkingService.SetOwnerOrg(workflowOwner, orgID)

		wfRegAddr := crecontracts.MustGetAddressFromDataStore(testEnv.CreEnvironment.CldfEnvironment.DataStore, testEnv.CreEnvironment.Blockchains[0].ChainSelector(), keystone_changeset.WorkflowRegistry.String(), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
		wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegAddr), sc.Client)
		require.NoError(t, err)
		requireVaultLinkOwner(t, sc, common.HexToAddress(wfRegAddr), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()])

		vaultParsedPublicKey := mustVaultPublicKey(t, vaultPublicKey)
		secretID := uniqueVaultSecretID("allowlistorgid")
		createEnc, err := vaultutils.EncryptSecretWithOrgID("secret-org-id-owner-create", vaultParsedPublicKey, orgID)
		require.NoError(t, err)
		updateEnc, err := vaultutils.EncryptSecretWithOrgID("secret-org-id-owner-update", vaultParsedPublicKey, orgID)
		require.NoError(t, err)

		allowlistAuth := newAllowlistVaultRequestAuth(workflowOwner, sc, wfReg)
		namespaces := []string{"main"}
		executeVaultSecretsCreateWithAuthExpectOwnersAndIdentifierOwner(t, allowlistAuth, orgID, createEnc, secretID, []string{orgID}, gwURL, namespaces)
		executeVaultSecretsUpdateWithAuthAndIdentifierOwner(t, allowlistAuth, orgID, updateEnc, secretID, orgID, gwURL, namespaces)
		executeVaultSecretsListWithAuthAndOwner(t, allowlistAuth, orgID, []string{secretID}, orgID, gwURL, "main")
		executeVaultSecretsDeleteWithAuthAndIdentifierOwner(t, allowlistAuth, orgID, secretID, orgID, gwURL, namespaces)
	})
}

func ExecuteVaultMixedAuthTest(t *testing.T, fixture *vaultScenarioFixture, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	issuer := fixture.Issuer
	linkingService := fixture.LinkingService

	gatewayURL := fixture.GatewayURL
	vaultPublicKey := fixture.VaultPublicKey

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	workflowOwner := sc.MustGetRootKeyAddress().Hex()
	orgID := "org" + strings.ReplaceAll(uuid.NewString(), "-", "")
	linkingService.SetOwnerOrg(workflowOwner, orgID)

	wfRegAddr := crecontracts.MustGetAddressFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.Blockchains[0].ChainSelector(),
		keystone_changeset.WorkflowRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()],
		"",
	)
	wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegAddr), sc.Client)
	require.NoError(t, err)
	requireVaultLinkOwner(t, sc, common.HexToAddress(wfRegAddr), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()])

	allowlistAuth := newAllowlistVaultRequestAuth(workflowOwner, sc, wfReg)

	ulCh := make(chan *workflowevents.UserLogs, 1000)
	bmCh := make(chan *commonevents.BaseMessage, 1000)
	sink := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, ulCh, bmCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, sink, ulCh, bmCh)
	})

	gwURL := gatewayURL.String()
	jwtAuth := newJWTVaultRequestAuth(issuer, orgID, workflowOwner)
	vaultParsedPublicKey := mustVaultPublicKey(t, vaultPublicKey)
	workflowOwnerAddress := common.HexToAddress(workflowOwner)

	t.Run("jwt_crud_with_workflow_owner", func(t *testing.T) {
		secretID := uniqueVaultSecretID("jwt")
		createValue := "secret-jwt-workflow-owner"
		enc, err := vaultutils.EncryptSecretWithOrgID(createValue, vaultParsedPublicKey, orgID)
		require.NoError(t, err)

		executeVaultJWTSecretsCreateTest(t, issuer, enc, secretID, orgID, workflowOwner, gwURL, []string{"main", "alt"})
		workflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, "jwt-lifecycle", []vaultWorkflowPhase{
			{
				Name: "jwt-created",
				Checks: []vaultWorkflowCheck{
					{Name: "jwt-create-get-main", SecretKey: secretID, SecretNamespace: "main", ExpectedValue: createValue},
					{Name: "jwt-create-get-alt", SecretKey: secretID, SecretNamespace: "alt", ExpectedValue: createValue},
				},
			},
			{
				Name: "jwt-deleted",
				Checks: []vaultWorkflowCheck{
					{Name: "jwt-delete-main-not-found", SecretKey: secretID, SecretNamespace: "main", ExpectNotFound: true},
					{Name: "jwt-delete-alt-not-found", SecretKey: secretID, SecretNamespace: "alt", ExpectNotFound: true},
				},
			},
		})
		waitForVaultWorkflowPhase(t, workflowID, "jwt-created", ulCh, bmCh)
		executeVaultJWTSecretsListTest(t, issuer, secretID, orgID, workflowOwner, gwURL, "main")
		executeVaultJWTSecretsListTest(t, issuer, secretID, orgID, workflowOwner, gwURL, "alt")
		executeVaultJWTSecretsDeleteTest(t, issuer, secretID, orgID, workflowOwner, gwURL, []string{"main", "alt"})
		waitForVaultWorkflowPhase(t, workflowID, "jwt-deleted", ulCh, bmCh)
	})

	t.Run("jwt_rejected_when_secret_labeled_as_workflow_owner", func(t *testing.T) {
		secretID := uniqueVaultSecretID("jwtreject")
		encryptedSecret, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-jwt-wrong-label", vaultParsedPublicKey, workflowOwnerAddress)
		require.NoError(t, err)

		uniqueRequestID := uuid.New().String()
		secretsCreateRequest := vault_helpers.CreateSecretsRequest{
			RequestId:        uniqueRequestID,
			EncryptedSecrets: buildEncryptedSecrets(secretID, orgID, encryptedSecret, []string{"main"}),
		}
		jsonRequest := newVaultJSONRequest(t, uniqueRequestID, vaulttypes.MethodSecretsCreate, &secretsCreateRequest)
		jwtAuth.apply(t, &jsonRequest)

		jsonResponse := sendVaultJWTRequestToGatewayExpectError(t, gwURL, jsonRequest, http.StatusBadRequest)
		require.Equal(t, uniqueRequestID, jsonResponse.ID)
		require.NotNil(t, jsonResponse.Error)
		require.Equal(t, jsonrpc.ErrInvalidParams, jsonResponse.Error.Code)
		require.Contains(t, jsonResponse.Error.Error(), "doesn't have owner as the label")
	})

	t.Run("mixed_allowlist_and_jwt_auth", func(t *testing.T) {
		t.Run("cross_auth_create_update_list_and_delete", func(t *testing.T) {
			allowlistSecretID := uniqueVaultSecretID("mixedallowlist")
			jwtSecretID := uniqueVaultSecretID("mixedjwt")
			allowlistCreateValue := "secret-mixed-allowlist-create"
			jwtCreateValue := "secret-mixed-jwt-create"
			allowlistUpdateValue := "secret-mixed-allowlist-update"
			jwtUpdateValue := "secret-mixed-jwt-update"
			allowlistCreateEnc, err := vaultutils.EncryptSecretWithOrgID(allowlistCreateValue, vaultParsedPublicKey, orgID)
			require.NoError(t, err)
			jwtCreateEnc, err := vaultutils.EncryptSecretWithOrgID(jwtCreateValue, vaultParsedPublicKey, orgID)
			require.NoError(t, err)
			allowlistUpdateEnc, err := vaultutils.EncryptSecretWithOrgID(allowlistUpdateValue, vaultParsedPublicKey, orgID)
			require.NoError(t, err)
			jwtUpdateEnc, err := vaultutils.EncryptSecretWithOrgID(jwtUpdateValue, vaultParsedPublicKey, orgID)
			require.NoError(t, err)

			executeVaultSecretsCreateWithAuth(t, allowlistAuth, allowlistCreateEnc, allowlistSecretID, orgID, gwURL, []string{"main"})
			executeVaultSecretsCreateWithAuth(t, jwtAuth, jwtCreateEnc, jwtSecretID, orgID, gwURL, []string{"main"})
			workflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, "mixed-lifecycle", []vaultWorkflowPhase{
				{
					Name: "mixed-created",
					Checks: []vaultWorkflowCheck{
						{Name: "mixed-allowlist-create-get-main", SecretKey: allowlistSecretID, SecretNamespace: "main", ExpectedValue: allowlistCreateValue},
						{Name: "mixed-jwt-create-get-main", SecretKey: jwtSecretID, SecretNamespace: "main", ExpectedValue: jwtCreateValue},
					},
				},
				{
					Name: "mixed-updated",
					Checks: []vaultWorkflowCheck{
						{Name: "mixed-jwt-update-get-main", SecretKey: allowlistSecretID, SecretNamespace: "main", ExpectedValue: jwtUpdateValue},
						{Name: "mixed-allowlist-update-get-main", SecretKey: jwtSecretID, SecretNamespace: "main", ExpectedValue: allowlistUpdateValue},
					},
				},
				{
					Name: "mixed-deleted",
					Checks: []vaultWorkflowCheck{
						{Name: "mixed-allowlist-delete-not-found", SecretKey: allowlistSecretID, SecretNamespace: "main", ExpectNotFound: true},
						{Name: "mixed-jwt-delete-not-found", SecretKey: jwtSecretID, SecretNamespace: "main", ExpectNotFound: true},
					},
				},
			})
			waitForVaultWorkflowPhase(t, workflowID, "mixed-created", ulCh, bmCh)

			executeVaultSecretsUpdateWithAuth(t, jwtAuth, jwtUpdateEnc, allowlistSecretID, orgID, gwURL, []string{"main"})
			executeVaultSecretsUpdateWithAuth(t, allowlistAuth, allowlistUpdateEnc, jwtSecretID, orgID, gwURL, []string{"main"})
			waitForVaultWorkflowPhase(t, workflowID, "mixed-updated", ulCh, bmCh)

			executeVaultSecretsListWithAuth(t, allowlistAuth, []string{allowlistSecretID, jwtSecretID}, orgID, gwURL, "main")
			executeVaultSecretsListWithAuth(t, jwtAuth, []string{allowlistSecretID, jwtSecretID}, orgID, gwURL, "main")

			executeVaultSecretsDeleteWithAuth(t, allowlistAuth, allowlistSecretID, orgID, gwURL, []string{"main"})
			executeVaultSecretsDeleteWithAuth(t, jwtAuth, jwtSecretID, orgID, gwURL, []string{"main"})
			waitForVaultWorkflowPhase(t, workflowID, "mixed-deleted", ulCh, bmCh)
		})
	})

	t.Run("jwt_without_workflow_owner_uses_org_id_identity", func(t *testing.T) {
		secretID := uniqueVaultSecretID("jwtorgonly")
		encryptedSecret, err := vaultutils.EncryptSecretWithOrgID("secret-jwt-org-only", vaultParsedPublicKey, orgID)
		require.NoError(t, err)

		orgOnlyJWTAuth := newJWTVaultRequestAuth(issuer, orgID, "")
		executeVaultSecretsCreateWithAuth(t, orgOnlyJWTAuth, encryptedSecret, secretID, orgID, gwURL, []string{"main"})
		executeVaultSecretsListWithAuth(t, orgOnlyJWTAuth, []string{secretID}, orgID, gwURL, "main")
		executeVaultSecretsDeleteWithAuth(t, orgOnlyJWTAuth, secretID, orgID, gwURL, []string{"main"})
	})

	t.Run("jwt_rejected_when_vault_secret_management_claim_false", func(t *testing.T) {
		executeVaultJWTSecretsCreateUnauthorizedWithExtraClaimsTest(t, issuer, vaultPublicKey, orgID, workflowOwner, gwURL,
			map[string]any{vaultcap.ClaimVaultSecretManagementEnabled: "false"},
			vaultcap.ErrVaultSecretManagementNotEnabled.Error(),
		)
	})
}

func ExecuteVaultJWTDisabledTest(t *testing.T, fixture *vaultScenarioFixture) {
	t.Helper()
	issuer := fixture.Issuer
	gatewayURL := fixture.GatewayURL
	vaultPublicKey := fixture.VaultPublicKey

	orgID := "org" + strings.ReplaceAll(uuid.NewString(), "-", "")
	gwURL := gatewayURL.String()

	t.Run("jwt_with_workflow_owner_rejected_when_jwt_auth_disabled", func(t *testing.T) {
		executeVaultJWTSecretsCreateUnauthorizedTest(t, issuer, vaultPublicKey, orgID, "0x1234567890abcdef1234567890abcdef12345678", gwURL, "JWTBasedAuth is disabled")
	})

	t.Run("jwt_without_workflow_owner_rejected_when_jwt_auth_disabled", func(t *testing.T) {
		executeVaultJWTSecretsCreateUnauthorizedTest(t, issuer, vaultPublicKey, orgID, "", gwURL, "JWTBasedAuth is disabled")
	})
}

func TestVaultStaticTopologies_LoadExpectedConfig(t *testing.T) {
	t.Parallel()
	dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")

	testCases := []struct {
		name        string
		configPath  string
		wantJWTGate string
		wantOrgGate string
		wantLinking bool
	}{
		{
			name:        "enabled",
			configPath:  vaultJWTAuthEnabledConfigPath,
			wantJWTGate: "true",
			wantOrgGate: "true",
			wantLinking: false,
		},
		{
			name:        "default",
			configPath:  vaultDefaultConfigPath,
			wantJWTGate: "false",
			wantOrgGate: "false",
			wantLinking: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &envconfig.Config{}
			require.NoError(t, cfg.Load(t_helpers.GetTestConfig(t, tc.configPath).EnvironmentConfigPath))

			for _, nodeSet := range cfg.NodeSets {
				if nodeSet.Name != "workflow" && nodeSet.Name != "capabilities" {
					continue
				}
				settingsRaw := nodeSet.EnvVars["CL_CRE_SETTINGS_DEFAULT"]
				if settingsRaw == "" {
					require.Equal(t, "false", tc.wantJWTGate)
					require.Equal(t, "false", tc.wantOrgGate)
				} else {
					var settings map[string]string
					require.NoError(t, json.Unmarshal([]byte(settingsRaw), &settings))
					require.Equal(t, tc.wantJWTGate, settings["VaultJWTAuthEnabled"])
					require.Equal(t, tc.wantOrgGate, settings["VaultOrgIdAsSecretOwnerEnabled"])
				}

				for _, nodeSpec := range nodeSet.NodeSpecs {
					if tc.wantLinking {
						require.Contains(t, nodeSpec.Node.UserConfigOverrides, "[CRE.Linking]")
						require.Contains(t, nodeSpec.Node.UserConfigOverrides, dockerHost+":18124")
						continue
					}
					require.Empty(t, nodeSpec.Node.UserConfigOverrides)
				}
			}
		})
	}
}

func TestMustMintVaultJWTForRequest_UsesRawRequestDigest(t *testing.T) {
	issuer, err := vault.NewTestJWTIssuer()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, issuer.Close())
	})

	params, err := json.Marshal(vault_helpers.CreateSecretsRequest{
		RequestId: "req-1",
		EncryptedSecrets: []*vault_helpers.EncryptedSecret{
			{
				Id: &vault_helpers.SecretIdentifier{
					Key:       "9838",
					Namespace: "main",
					Owner:     "org-123",
				},
				EncryptedValue: "cipher+/==",
			},
		},
	})
	require.NoError(t, err)

	rawParams := json.RawMessage(params)
	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-1",
		Method:  vaulttypes.MethodSecretsCreate,
		Params:  &rawParams,
	}
	req.Auth = mustMintVaultJWTForRequest(t, issuer, req, "org-123", "0xAbCdEf0123456789AbCdEf0123456789AbCdEf01")

	outboundReq := outboundRequestWithoutAuth(req)
	requestDigest, err := outboundReq.Digest()
	require.NoError(t, err)

	parsedToken, _, err := new(jwt.Parser).ParseUnverified(req.Auth, jwt.MapClaims{})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)
	authorizationDetails, ok := claims["authorization_details"].([]interface{})
	require.True(t, ok)

	var claimedDigest string
	for _, detail := range authorizationDetails {
		entry, ok := detail.(map[string]interface{})
		require.True(t, ok)
		if entry["type"] == "request_digest" {
			claimedDigest, ok = entry["value"].(string)
			require.True(t, ok)
			break
		}
	}

	require.NotEmpty(t, claimedDigest)
	require.Equal(t, requestDigest, claimedDigest)
}
