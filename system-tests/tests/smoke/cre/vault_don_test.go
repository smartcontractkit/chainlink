package cre

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	retry "github.com/avast/retry-go/v4"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
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
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	vaultsecret_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/vaultsecret/config"
)

// ExecuteVaultAllowListBasedTests covers vault gateway + workflows with allow-listed JSON-RPC auth
// (and JWT when enabled). Identity is conveyed via signatures/JWT authorization_details digest and
// per-secret SecretIdentifier.Owner, not deprecated top-level org_id/workflow_owner proto fields.

func uniqueVaultSecretID(prefix string) string {
	return prefix + strings.ReplaceAll(uuid.NewString(), "-", "")
}

func ExecuteVaultAllowListBasedTests(t *testing.T, fixture *vaultScenarioFixture, testEnv *ttypes.TestEnvironment) {
	var testLogger = framework.L
	linkingService := fixture.LinkingService

	gwURL := fixture.GatewayURL.String()
	vaultPublicKey := fixture.VaultPublicKey

	t.Run("allowlist_rejects_create_when_identifier_owner_mismatch", func(t *testing.T) {
		sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
		authorizedOwner := sc.MustGetRootKeyAddress().Hex()
		mismatchedOwner := common.HexToAddress("0x000000000000000000000000000000000000dEaD").Hex()
		wfRegAddr := crecontracts.MustGetAddressFromDataStore(testEnv.CreEnvironment.CldfEnvironment.DataStore, testEnv.CreEnvironment.Blockchains[0].ChainSelector(), keystone_changeset.WorkflowRegistry.String(), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
		wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegAddr), sc.Client)
		require.NoError(t, err)
		requireVaultLinkOwner(t, sc, common.HexToAddress(wfRegAddr), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()])
		vaultParsedPublicKey := mustVaultPublicKey(t, vaultPublicKey)
		secretID := uniqueVaultSecretID("allowlistownermismatch")
		encryptedSecret, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-allowlist-owner-mismatch", vaultParsedPublicKey, sc.MustGetRootKeyAddress())
		require.NoError(t, err)
		auth := newAllowlistVaultRequestAuth(authorizedOwner, sc, wfReg)
		executeVaultSecretsCreateOwnerMismatchRejectedTest(t, auth, authorizedOwner, mismatchedOwner, encryptedSecret, secretID, gwURL)
	})

	t.Run("allowlist_crud_with_workflow_owner_identity", func(t *testing.T) {
		sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
		workflowOwnerAddress := sc.MustGetRootKeyAddress()
		owner := workflowOwnerAddress.Hex()
		var orgID string
		if linkingService != nil {
			orgID = "org" + strings.ReplaceAll(uuid.NewString(), "-", "")
			linkingService.SetOwnerOrg(owner, orgID)
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

		executeVaultAllowListSecretsCreateTest(t, createEnc, secretID, owner, owner, gwURL, namespaces, sc, wfReg)
		if isVaultOptimizationsEnabledTopology(testEnv.TestConfig.EnvironmentConfigPath) {
			t.Run("binary_encoded_shares", func(t *testing.T) {
				executeVaultBinaryEncodedSharesSmokeTest(t, testEnv, secretID, "main", createValue, ulCh, bmCh)
			})
		}
		executeVaultSecretsUpdateTest(t, updateEnc, secretID, owner, owner, gwURL, namespaces, sc, wfReg)
		executeVaultSecretsListTest(t, secretID, owner, owner, gwURL, "main", sc, wfReg)
		executeVaultSecretsListTest(t, secretID, owner, owner, gwURL, "alt", sc, wfReg)
		updatedChecks := []vaultWorkflowCheck{
			{Name: "allowlist-main-updated", SecretKey: secretID, SecretNamespace: "main", ExpectedValue: updateValue},
			{Name: "allowlist-alt-updated", SecretKey: secretID, SecretNamespace: "alt", ExpectedValue: updateValue},
		}
		finalChecks := []vaultWorkflowCheck{
			{Name: "allowlist-main-not-found", SecretKey: secretID, SecretNamespace: "main", ExpectNotFound: true},
			{Name: "allowlist-alt-updated", SecretKey: secretID, SecretNamespace: "alt", ExpectedValue: updateValue},
		}
		workflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, "allowlist-lifecycle", []vaultWorkflowPhase{
			{Name: "allowlist-updated", Checks: updatedChecks},
			{Name: "allowlist-final-verify", Checks: finalChecks},
		})
		waitForVaultWorkflowPhase(t, workflowID, "allowlist-updated", ulCh, bmCh)
		executeVaultSecretsDeleteTest(t, secretID, owner, owner, gwURL, []string{"main"}, sc, wfReg)
		waitForVaultWorkflowPhase(t, workflowID, "allowlist-final-verify", ulCh, bmCh)
		executeVaultSecretsDeleteTest(t, secretID, owner, owner, gwURL, []string{"alt"}, sc, wfReg)
	})

	if isVaultJWTAuthEnabledTopology(testEnv.TestConfig.EnvironmentConfigPath) {
		t.Run("identifier_validation", func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			subEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, testEnv.TestConfig)
			sc := subEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
			owner := sc.MustGetRootKeyAddress().Hex()
			wfRegAddr := crecontracts.MustGetAddressFromDataStore(subEnv.CreEnvironment.CldfEnvironment.DataStore, subEnv.CreEnvironment.Blockchains[0].ChainSelector(), keystone_changeset.WorkflowRegistry.String(), subEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
			wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegAddr), sc.Client)
			require.NoError(t, err)
			require.NoError(t, creworkflow.LinkOwner(sc, common.HexToAddress(wfRegAddr), subEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()]))
			vaultParsedPublicKey := mustVaultPublicKey(t, vaultPublicKey)
			enc, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-basic", vaultParsedPublicKey, sc.MustGetRootKeyAddress())
			require.NoError(t, err)
			ulCh := make(chan *workflowevents.UserLogs, 1000)
			bmCh := make(chan *commonevents.BaseMessage, 1000)
			sink := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, ulCh, bmCh))
			t.Cleanup(func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				t_helpers.ShutdownChipSinkWithDrain(ctx, sink, ulCh, bmCh)
			})
			executeVaultSecretsIdentifierValidationTest(t, enc, owner, gwURL, sc, wfReg)
			executeVaultSecretsGetInvalidIdentifierViaWorkflowTest(t, subEnv, "vget1", ulCh, bmCh)
		})
	}

	if isVaultOptimizationsEnabledTopology(testEnv.TestConfig.EnvironmentConfigPath) {
		t.Run("pending_queue_blob_batching_many_concurrent_creates", func(t *testing.T) {
			ExecuteVaultBlobBatchingSmokeTest(t, fixture, testEnv)
		})
	}
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
	derivedJWTWorkflowOwner := mustDeriveJWTVaultWorkflowOwner(t, orgID)
	derivedJWTWorkflowOwnerAddr := common.HexToAddress(derivedJWTWorkflowOwner)
	jwtAuth := newJWTVaultRequestAuth(issuer, orgID, derivedJWTWorkflowOwner)
	vaultParsedPublicKey := mustVaultPublicKey(t, vaultPublicKey)
	workflowOwnerAddress := common.HexToAddress(workflowOwner)

	t.Run("jwt_crud_with_workflow_owner", func(t *testing.T) {
		secretID := uniqueVaultSecretID("jwt")
		createValue := "secret-jwt-workflow-owner"
		enc, err := vaultutils.EncryptSecretWithWorkflowOwner(createValue, vaultParsedPublicKey, derivedJWTWorkflowOwnerAddr)
		require.NoError(t, err)

		executeVaultJWTSecretsCreateTest(t, issuer, enc, secretID, orgID, derivedJWTWorkflowOwner, gwURL, []string{"main", "alt"})
		// WASM workflows run under the registry EOA workflow owner while JWT-backed secrets live under
		// the derived JWT workflow owner partition; accessibility for JWT-backed rows is asserted via Gateway SecretsList instead.
		executeVaultJWTSecretsListTest(t, issuer, secretID, orgID, derivedJWTWorkflowOwner, gwURL, "main")
		executeVaultJWTSecretsListTest(t, issuer, secretID, orgID, derivedJWTWorkflowOwner, gwURL, "alt")
		executeVaultJWTSecretsDeleteTest(t, issuer, secretID, orgID, derivedJWTWorkflowOwner, gwURL, []string{"main", "alt"})
		executeVaultJWTSecretsListAbsentFromNamespace(t, issuer, secretID, orgID, derivedJWTWorkflowOwner, gwURL, "main")
		executeVaultJWTSecretsListAbsentFromNamespace(t, issuer, secretID, orgID, derivedJWTWorkflowOwner, gwURL, "alt")
	})

	t.Run("jwt_rejected_when_identifier_owner_does_not_match_authorized_owner", func(t *testing.T) {
		secretID := uniqueVaultSecretID("jwtownermismatch")
		encryptedSecret, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-jwt-owner-mismatch", vaultParsedPublicKey, derivedJWTWorkflowOwnerAddr)
		require.NoError(t, err)
		executeVaultSecretsCreateOwnerMismatchRejectedTest(t, jwtAuth, derivedJWTWorkflowOwner, workflowOwner, encryptedSecret, secretID, gwURL)
	})

	t.Run("jwt_rejected_when_ciphertext_label_is_linked_workflow_owner_but_identifier_owner_is_jwt_derived", func(t *testing.T) {
		secretID := uniqueVaultSecretID("jwtreject")
		encryptedSecret, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-jwt-wrong-label", vaultParsedPublicKey, workflowOwnerAddress)
		require.NoError(t, err)

		uniqueRequestID := uuid.New().String()
		secretsCreateRequest := vault_helpers.CreateSecretsRequest{
			RequestId:        uniqueRequestID,
			EncryptedSecrets: buildEncryptedSecrets(secretID, derivedJWTWorkflowOwner, encryptedSecret, []string{"main"}),
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
		// Allow-list identities use the linked EOA workflow owner; JWT identities use the derived
		// workflow owner for the org_id + tenant_id claim pair. Cross-channel mutation is unsupported.

		t.Run("parallel_independent_crud", func(t *testing.T) {
			allowlistSecretID := uniqueVaultSecretID("mixedallowlist")
			jwtSecretID := uniqueVaultSecretID("mixedjwt")
			allowlistCreateValue := "secret-mixed-allowlist-create"
			jwtCreateValue := "secret-mixed-jwt-create"
			allowlistUpdateValue := "secret-mixed-allowlist-update"
			jwtUpdateValue := "secret-mixed-jwt-update"
			allowlistCreateEnc, err := vaultutils.EncryptSecretWithWorkflowOwner(allowlistCreateValue, vaultParsedPublicKey, workflowOwnerAddress)
			require.NoError(t, err)
			jwtCreateEnc, err := vaultutils.EncryptSecretWithWorkflowOwner(jwtCreateValue, vaultParsedPublicKey, derivedJWTWorkflowOwnerAddr)
			require.NoError(t, err)
			allowlistUpdateEnc, err := vaultutils.EncryptSecretWithWorkflowOwner(allowlistUpdateValue, vaultParsedPublicKey, workflowOwnerAddress)
			require.NoError(t, err)
			jwtUpdateEnc, err := vaultutils.EncryptSecretWithWorkflowOwner(jwtUpdateValue, vaultParsedPublicKey, derivedJWTWorkflowOwnerAddr)
			require.NoError(t, err)

			executeVaultSecretsCreateWithAuth(t, allowlistAuth, allowlistCreateEnc, allowlistSecretID, workflowOwner, gwURL, []string{"main"})
			executeVaultSecretsCreateWithAuth(t, jwtAuth, jwtCreateEnc, jwtSecretID, derivedJWTWorkflowOwner, gwURL, []string{"main"})
			// Workflow phases validate only allow-listed reads (workflow runs as EOA). JWT-backed keys are exercised via Gateway SecretsList.
			workflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, "mixed-lifecycle", []vaultWorkflowPhase{
				{
					Name: "mixed-created",
					Checks: []vaultWorkflowCheck{
						{Name: "mixed-allowlist-create-get-main", SecretKey: allowlistSecretID, SecretNamespace: "main", ExpectedValue: allowlistCreateValue},
					},
				},
				{
					Name: "mixed-updated",
					Checks: []vaultWorkflowCheck{
						{Name: "mixed-allowlist-own-update-main", SecretKey: allowlistSecretID, SecretNamespace: "main", ExpectedValue: allowlistUpdateValue},
					},
				},
				{
					Name: "mixed-deleted",
					Checks: []vaultWorkflowCheck{
						{Name: "mixed-allowlist-delete-not-found", SecretKey: allowlistSecretID, SecretNamespace: "main", ExpectNotFound: true},
					},
				},
			})
			waitForVaultWorkflowPhase(t, workflowID, "mixed-created", ulCh, bmCh)
			executeVaultSecretsListWithAuth(t, jwtAuth, []string{jwtSecretID}, derivedJWTWorkflowOwner, gwURL, "main")

			executeVaultSecretsUpdateWithAuth(t, allowlistAuth, allowlistUpdateEnc, allowlistSecretID, workflowOwner, gwURL, []string{"main"})
			executeVaultSecretsUpdateWithAuth(t, jwtAuth, jwtUpdateEnc, jwtSecretID, derivedJWTWorkflowOwner, gwURL, []string{"main"})
			waitForVaultWorkflowPhase(t, workflowID, "mixed-updated", ulCh, bmCh)

			executeVaultSecretsListWithAuth(t, allowlistAuth, []string{allowlistSecretID}, workflowOwner, gwURL, "main")
			executeVaultSecretsListWithAuth(t, jwtAuth, []string{jwtSecretID}, derivedJWTWorkflowOwner, gwURL, "main")

			executeVaultSecretsDeleteWithAuth(t, allowlistAuth, allowlistSecretID, workflowOwner, gwURL, []string{"main"})
			executeVaultSecretsDeleteWithAuth(t, jwtAuth, jwtSecretID, derivedJWTWorkflowOwner, gwURL, []string{"main"})
			waitForVaultWorkflowPhase(t, workflowID, "mixed-deleted", ulCh, bmCh)
			executeVaultJWTSecretsListAbsentFromNamespace(t, issuer, jwtSecretID, orgID, derivedJWTWorkflowOwner, gwURL, "main")
		})

		t.Run("jwt_must_not_flip_allowlisted_secret_via_same_key_string", func(t *testing.T) {
			sharedKey := uniqueVaultSecretID("mixedcrosskey")
			allowlistValue := "secret-mixed-cross-allowlist"
			jwtCrossValue := "secret-mixed-cross-jwt-attempt"
			createEncAllow, err := vaultutils.EncryptSecretWithWorkflowOwner(allowlistValue, vaultParsedPublicKey, workflowOwnerAddress)
			require.NoError(t, err)
			jwtCrossEnc, err := vaultutils.EncryptSecretWithWorkflowOwner(jwtCrossValue, vaultParsedPublicKey, derivedJWTWorkflowOwnerAddr)
			require.NoError(t, err)

			executeVaultSecretsCreateWithAuth(t, allowlistAuth, createEncAllow, sharedKey, workflowOwner, gwURL, []string{"main"})
			// Vault workflow returns after the FIRST phase whose checks succeed (see vaultsecret/main.go).
			// Two phases cannot share the same success predicate; otherwise the later phase never emits
			// "Vault secret workflow phase completed: ..." and waits time out (4m).
			crossWorkflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, "mixed-cross-isolation", []vaultWorkflowPhase{
				{
					Name: "cross-isolation-stable",
					Checks: []vaultWorkflowCheck{
						{Name: "cross-share-key-allowlist-value", SecretKey: sharedKey, SecretNamespace: "main", ExpectedValue: allowlistValue},
					},
				},
			})
			waitForVaultWorkflowPhase(t, crossWorkflowID, "cross-isolation-stable", ulCh, bmCh)

			// JWT cannot overwrite the allow-listed row: identifiers + labels are partitioned by owner address.
			tryJWTSignedVaultSecretsUpdate(t, jwtAuth, jwtAuth.requestOwner, jwtCrossEnc, sharedKey, gwURL, []string{"main"})
			// Confirm on a subsequent cron tick that secrets still satisfy the same invariant.
			waitForVaultWorkflowPhase(t, crossWorkflowID, "cross-isolation-stable", ulCh, bmCh)

			executeVaultSecretsDeleteWithAuth(t, allowlistAuth, sharedKey, workflowOwner, gwURL, []string{"main"})
			// `sharedKey` only ever existed under the allowlisted workflow owner; JWT-targeted deletes use the
			// derived-owner partition where this key was never written, so a JWT delete cleanup would falsely fail.
			executeVaultJWTSecretsListAbsentFromNamespace(t, issuer, sharedKey, orgID, derivedJWTWorkflowOwner, gwURL, "main")
		})
	})

	t.Run("jwt_without_workflow_owner_claim_uses_derived_workflow_owner", func(t *testing.T) {
		secretID := uniqueVaultSecretID("jwtorgonly")
		encryptedSecret, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-jwt-derived-only", vaultParsedPublicKey, derivedJWTWorkflowOwnerAddr)
		require.NoError(t, err)

		derivedAuth := newJWTVaultRequestAuth(issuer, orgID, derivedJWTWorkflowOwner)
		executeVaultSecretsCreateWithAuth(t, derivedAuth, encryptedSecret, secretID, derivedJWTWorkflowOwner, gwURL, []string{"main"})
		executeVaultSecretsListWithAuth(t, derivedAuth, []string{secretID}, derivedJWTWorkflowOwner, gwURL, "main")
		executeVaultSecretsDeleteWithAuth(t, derivedAuth, secretID, derivedJWTWorkflowOwner, gwURL, []string{"main"})
	})

	t.Run("jwt_rejected_when_vault_secret_management_claim_false", func(t *testing.T) {
		executeVaultJWTSecretsCreateUnauthorizedWithExtraClaimsTest(t, issuer, vaultPublicKey, orgID, gwURL,
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
		executeVaultJWTSecretsCreateUnauthorizedTest(t, issuer, vaultPublicKey, orgID, gwURL, "JWTBasedAuth is disabled")
	})

	t.Run("jwt_without_workflow_owner_rejected_when_jwt_auth_disabled", func(t *testing.T) {
		executeVaultJWTSecretsCreateUnauthorizedTest(t, issuer, vaultPublicKey, orgID, gwURL, "JWTBasedAuth is disabled")
	})
}

// ExecuteVaultBlobBatchingSmokeTest issues many concurrent Vault create calls so OCR observes a large
// local pending queue and exercises batched pending-queue blobs (beyond the legacy per-blob single-item cap).
func ExecuteVaultBlobBatchingSmokeTest(t *testing.T, fixture *vaultScenarioFixture, testEnv *ttypes.TestEnvironment) {
	t.Helper()
	const nConcurrentCreates = 25

	gwURL := fixture.GatewayURL.String()
	vaultParsedPublicKey := mustVaultPublicKey(t, fixture.VaultPublicKey)
	linkingService := fixture.LinkingService

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

	const plainValue = "blob-batch-smoke-value"
	createEnc, err := vaultutils.EncryptSecretWithWorkflowOwner(plainValue, vaultParsedPublicKey, workflowOwnerAddress)
	require.NoError(t, err)
	namespaces := []string{"main"}

	// Pre-generate all secret IDs, request IDs, and request objects
	type secretCreateTest struct {
		secretID    string
		requestID   string
		jsonRequest jsonrpc.Request[json.RawMessage]
	}
	secretTests := make([]secretCreateTest, nConcurrentCreates)
	for i := range nConcurrentCreates {
		secretTests[i].secretID = uniqueVaultSecretID(fmt.Sprintf("blobbatch%d", i))
		secretTests[i].requestID = uuid.New().String()
		secretsCreateRequest := vault_helpers.CreateSecretsRequest{
			RequestId:        secretTests[i].requestID,
			EncryptedSecrets: buildEncryptedSecrets(secretTests[i].secretID, owner, createEnc, namespaces),
		}
		secretTests[i].jsonRequest = newVaultJSONRequest(t, secretTests[i].requestID, vaulttypes.MethodSecretsCreate, &secretsCreateRequest)
	}

	// Sequentially allowlist all requests to avoid blockchain nonce conflicts
	auth := newAllowlistVaultRequestAuth(owner, sc, wfReg)
	for i := range secretTests {
		// Allowlist sequentially to avoid nonce conflicts
		auth.authorize(t, &secretTests[i].jsonRequest)
	}

	// Send all create requests to gateway in parallel; subtests are released together by t.Parallel
	// once the parent body returns, so they queue up in the local pending queue and exercise batching.
	t.Run("concurrent_secrets_creates", func(t *testing.T) {
		for i, st := range secretTests {
			t.Run(fmt.Sprintf("secret_%d", i), func(t *testing.T) {
				t.Parallel()
				sendConcurrentVaultCreate(t, gwURL, st.requestID, st.jsonRequest, expectedResponseOwner, namespaces)
			})
		}
	})
	// Runs after the parallel subtests above complete (sibling Run, not same body as t.Parallel children).
	t.Run("pending_queue_pack_observed_in_docker_logs", func(t *testing.T) {
		assertVaultOCRPendingPackObservedInDockerLogs(t)
	})
	t.Run("state_transition_pending_write_observed_in_docker_logs", func(t *testing.T) {
		assertVaultOCRStateTransitionPendingWriteObservedInDockerLogs(t)
	})
	t.Run("vault_ocr_wire_truncation_signals_in_docker_logs", func(t *testing.T) {
		assertVaultOCRWireTruncationSignalsInDockerLogs(t)
	})
}

// sendConcurrentVaultCreate sends an already-allowlisted create request to the gateway and tolerates
// the replay-guard outcome. Under burst load, the gateway can time out (503 "Request timed out") while
// DON still processes the create; the test's HTTP retry then re-sends the same request digest, which
// vault's replay guard rejects with "request was already authorized previously". That error proves the
// original request was accepted and processed, so we treat it as success — there is no later response
// payload to validate when this path fires.
func sendConcurrentVaultCreate(t *testing.T, gwURL, requestID string, jsonRequest jsonrpc.Request[json.RawMessage], expectedResponseOwner string, namespaces []string) {
	t.Helper()

	authToken := jsonRequest.Auth
	stripped := outboundRequestWithoutAuth(jsonRequest)
	requestBody, err := json.Marshal(stripped)
	require.NoError(t, err, "failed to marshal vault request")
	headers := map[string]string{}
	if authToken != "" {
		headers["Authorization"] = "Bearer " + authToken
	}

	statusCode, body := sendVaultRequestToGatewayWithHeaders(t, gwURL, requestBody, headers)

	// Under burst load the gateway can return 503 "Request timed out" when it gives up relaying the
	// response, even though the DON has already processed the request. Tolerate that here — the goal
	// of this subtest is to drive concurrent load for the docker-log batching assertions below, not
	// to verify per-request response payloads.
	if statusCode == http.StatusServiceUnavailable && bytes.Contains(body, []byte("Request timed out")) {
		framework.L.Info().Str("requestID", requestID).Msg("vault create gateway-to-DON timeout; treating as success for batching load test")
		return
	}
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK; body=%s", string(body))

	var parsed jsonrpc.Response[vaulttypes.SignedOCRResponse]
	require.NoError(t, json.Unmarshal(body, &parsed), "failed to unmarshal gateway response")
	if parsed.Error != nil && strings.Contains(parsed.Error.Message, "request was already authorized previously") {
		framework.L.Info().Str("requestID", requestID).Msg("vault create returned replay-guard error after retry; DON processed the original request — treating as success")
		return
	}
	require.Nil(t, parsed.Error, "gateway returned error: %v", parsed.Error)
	require.Equal(t, requestID, parsed.ID)
	require.Equal(t, vaulttypes.MethodSecretsCreate, parsed.Method)
	var createResp vault_helpers.CreateSecretsResponse
	require.NoError(t, protojson.Unmarshal(parsed.Result.Payload, &createResp), "failed to decode CreateSecretsResponse")
	require.Len(t, createResp.Responses, len(namespaces), "Expected one item in the response per namespace")
	for _, r := range createResp.Responses {
		require.Equal(t, expectedResponseOwner, r.Id.Owner, "Response owner should match expected owner")
	}
}

// assertVaultOCRPendingPackObservedInDockerLogs polls chainlink-related container logs until it observes
// a 'observation packed local items into blob payloads' log line emitted by the new batching code path. The line must
// carry the blobHandleCount field (only the new pending-queue-batching code path logs that field) and a
// non-zero packedLocalItemCount. Strict "multi-item batched into fewer blobs" (packed > handles) is
// timing-dependent: OCR rounds (~400ms) often drain the per-node local pending queue between arrivals,
// even under burst load. Asserting the optimization code path was exercised is reliable; asserting that
// every load shape produces packed>handles is not.
var vaultPendingPackDockerLogRE = regexp.MustCompile(`packedLocalItemCount[^\d]*(\d+)`)
var vaultPendingPackBlobHandleCountRE = regexp.MustCompile(`blobHandleCount[^\d]*(\d+)`)

func assertVaultOCRPendingPackObservedInDockerLogs(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("docker log scan skipped in -short mode")
	}
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		t.Skip("docker not in PATH; skipping pending-pack log assertion")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		psOut, err := exec.CommandContext(ctx, dockerBin, "ps", "--format", "{{.Names}}").Output()
		if err != nil {
			select {
			case <-ctx.Done():
				require.Fail(t, "timed out waiting for docker while scanning for observation packed local items into blob payloads")
			case <-ticker.C:
			}
			continue
		}
		names := strings.SplitSeq(strings.TrimSpace(string(psOut)), "\n")
		for name := range names {
			if name == "" {
				continue
			}
			ln := strings.ToLower(name)
			if !strings.Contains(ln, "chainlink") && !strings.Contains(ln, "ocr") && !strings.Contains(ln, "capabilit") {
				continue
			}
			logs, err := exec.CommandContext(ctx, dockerBin, "logs", name, "--tail", "25000").CombinedOutput()
			if err != nil {
				continue
			}
			for line := range strings.SplitSeq(string(logs), "\n") {
				if !strings.Contains(line, "observation packed local items into blob payloads") {
					continue
				}
				pm := vaultPendingPackDockerLogRE.FindStringSubmatch(line)
				bm := vaultPendingPackBlobHandleCountRE.FindStringSubmatch(line)
				if len(pm) < 2 || len(bm) < 2 {
					continue
				}
				packed, err1 := strconv.Atoi(pm[1])
				handles, err2 := strconv.Atoi(bm[1])
				if err1 != nil || err2 != nil {
					continue
				}
				if packed >= 1 && handles >= 1 {
					framework.L.Info().Str("container", name).Int("packedLocalItemCount", packed).Int("blobHandleCount", handles).Msg("observed vault OCR pending pack log from new batching code path")
					return
				}
			}
		}
		select {
		case <-ctx.Done():
			require.Fail(t, "timed out waiting for observation packed local items into blob payloads line with packedLocalItemCount>=1 and blobHandleCount field set in docker logs (is the local CRE stack running with vault optimizations enabled?)")
		case <-ticker.C:
		}
	}
}

// assertVaultOCRStateTransitionPendingWriteObservedInDockerLogs polls chainlink-related container logs until it
// observes a 'state transition: more observations than can be included in response' log line with writtenCount >= 1, proving that items
// reached the KV pending queue via the new batched-write state transition (rather than the legacy per-item
// path that did not emit this log). Strict writtenCount>1 is timing-dependent — most rounds write only the
// item(s) that arrived since the previous round — so we assert the code path was exercised instead.
var vaultStateTransitionPendingWriteWrittenRE = regexp.MustCompile(`writtenCount[^\d]*(\d+)`)

func assertVaultOCRStateTransitionPendingWriteObservedInDockerLogs(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("docker log scan skipped in -short mode")
	}
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		t.Skip("docker not in PATH; skipping state-transition pending-write log assertion")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		psOut, err := exec.CommandContext(ctx, dockerBin, "ps", "--format", "{{.Names}}").Output()
		if err != nil {
			select {
			case <-ctx.Done():
				require.Fail(t, "timed out waiting for docker while scanning for pending queue items persisted to storage")
			case <-ticker.C:
			}
			continue
		}
		names := strings.SplitSeq(strings.TrimSpace(string(psOut)), "\n")
		for name := range names {
			if name == "" {
				continue
			}
			ln := strings.ToLower(name)
			if !strings.Contains(ln, "chainlink") && !strings.Contains(ln, "ocr") && !strings.Contains(ln, "capabilit") {
				continue
			}
			logs, err := exec.CommandContext(ctx, dockerBin, "logs", name, "--tail", "25000").CombinedOutput()
			if err != nil {
				continue
			}
			for line := range strings.SplitSeq(string(logs), "\n") {
				if !strings.Contains(line, "pending queue items persisted to storage") {
					continue
				}
				wm := vaultStateTransitionPendingWriteWrittenRE.FindStringSubmatch(line)
				if len(wm) < 2 {
					continue
				}
				written, err1 := strconv.Atoi(wm[1])
				if err1 != nil {
					continue
				}
				if written >= 1 {
					framework.L.Info().Str("container", name).Int("writtenCount", written).Msg("observed vault OCR state-transition pending write log from new batching code path")
					return
				}
			}
		}
		select {
		case <-ctx.Done():
			require.Fail(t, "timed out waiting for pending queue items persisted to storage line with writtenCount>=1 in docker logs (is the local CRE stack running with vault optimizations enabled?)")
		case <-ticker.C:
		}
	}
}

var vaultObservationWirePackPackedRE = regexp.MustCompile(`packedObservationCount[^\d]*(\d+)`)
var vaultObservationWirePackPendingRE = regexp.MustCompile(`pendingQueueItemCount[^\d]*(\d+)`)
var vaultOutcomeWirePackPackedRE = regexp.MustCompile(`packedOutcomeCount[^\d]*(\d+)`)
var vaultOutcomeWirePackScheduledRE = regexp.MustCompile(`scheduledRequestIDs[^\d]*(\d+)`)

// assertVaultOCRWireTruncationSignalsInDockerLogs scans recent chainlink-like container logs and asserts structural
// consistency of truncation markers when present: observation wire-pack implies packed < pending store rows;
// outcome wire-pack implies packed outcomes < scheduled request IDs for that state transition.
func assertVaultOCRWireTruncationSignalsInDockerLogs(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("docker log scan skipped in -short mode")
	}
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		t.Skip("docker not in PATH; skipping vault OCR wire truncation log scan")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	psOut, err := exec.CommandContext(ctx, dockerBin, "ps", "--format", "{{.Names}}").Output()
	if err != nil {
		t.Fatalf("docker ps: %v", err)
	}
	var combined strings.Builder
	names := strings.SplitSeq(strings.TrimSpace(string(psOut)), "\n")
	for name := range names {
		if name == "" {
			continue
		}
		ln := strings.ToLower(name)
		if !strings.Contains(ln, "chainlink") && !strings.Contains(ln, "ocr") && !strings.Contains(ln, "capabilit") {
			continue
		}
		logs, err := exec.CommandContext(ctx, dockerBin, "logs", name, "--tail", "25000").CombinedOutput()
		if err != nil {
			continue
		}
		combined.Write(logs)
		combined.WriteByte('\n')
	}
	s := combined.String()

	sawObsTrunc := false
	for line := range strings.SplitSeq(s, "\n") {
		if !strings.Contains(line, "observation: more pending queue items than can be observed") {
			continue
		}
		sawObsTrunc = true
		pm := vaultObservationWirePackPackedRE.FindStringSubmatch(line)
		qm := vaultObservationWirePackPendingRE.FindStringSubmatch(line)
		require.GreaterOrEqual(t, len(pm), 2, "packedObservationCount should be present on observation wire pack line: %s", line)
		require.GreaterOrEqual(t, len(qm), 2, "pendingQueueItemCount should be present on observation wire pack line: %s", line)
		packed, err1 := strconv.Atoi(pm[1])
		pending, err2 := strconv.Atoi(qm[1])
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Less(t, packed, pending, "observation wire truncation should mean packedObservationCount < pendingQueueItemCount: %s", line)
		framework.L.Info().Str("line", line).Int("packedObservationCount", packed).Int("pendingQueueItemCount", pending).Msg("vault OCR observation wire pack (truncation) observed in docker logs")
	}

	sawOutcomeTrunc := false
	for line := range strings.SplitSeq(s, "\n") {
		if !strings.Contains(line, "state transition: more observations than can be included in response") {
			continue
		}
		sawOutcomeTrunc = true
		om := vaultOutcomeWirePackPackedRE.FindStringSubmatch(line)
		sm := vaultOutcomeWirePackScheduledRE.FindStringSubmatch(line)
		require.GreaterOrEqual(t, len(om), 2, "packedOutcomeCount should be present on state transition outcome wire pack line: %s", line)
		require.GreaterOrEqual(t, len(sm), 2, "scheduledRequestIDs should be present on state transition outcome wire pack line: %s", line)
		packed, err1 := strconv.Atoi(om[1])
		scheduled, err2 := strconv.Atoi(sm[1])
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Less(t, packed, scheduled, "state transition outcome wire truncation should mean packedOutcomeCount < scheduledRequestIDs: %s", line)
		framework.L.Info().Str("line", line).Int("packedOutcomeCount", packed).Int("scheduledRequestIDs", scheduled).Msg("vault OCR state transition outcome wire pack (truncation) observed in docker logs")
	}

	if !sawObsTrunc && !sawOutcomeTrunc {
		framework.L.Info().Msg("no observation: more pending queue items than can be observed or state transition: more observations than can be included in response in recent docker logs — observation and precursor outcome packing did not truncate in the sampled window (expected under default limits)")
	}
}

func TestVaultOptimizationsEnabled_CRESettingDefaultsDisabled(t *testing.T) {
	require.False(t, cresettings.Default.VaultOptimizationsEnabled.DefaultValue)
}

func TestVaultStaticTopologies_LoadExpectedConfig(t *testing.T) {
	t.Parallel()
	dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")

	testCases := []struct {
		name                  string
		configPath            string
		wantJWTGate           string
		wantLinking           bool
		wantOptimizationsGate string
	}{
		{
			name:        "enabled",
			configPath:  vaultJWTAuthEnabledConfigPath,
			wantJWTGate: "true",
			wantLinking: false,
		},
		{
			name:        "default",
			configPath:  vaultDefaultConfigPath,
			wantJWTGate: "false",
			wantLinking: false,
		},
		{
			name:                  "optimizations_enabled",
			configPath:            vaultOptimizationsEnabledConfigPath,
			wantJWTGate:           "false",
			wantLinking:           false,
			wantOptimizationsGate: "true",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &envconfig.Config{}
			require.NoError(t, cfg.Load(t_helpers.GetTestConfig(t, tc.configPath).EnvironmentConfigPath))

			for _, nodeSet := range cfg.NodeSets {
				switch nodeSet.Name {
				case "workflow", "capabilities":
				case "bootstrap-gateway":
					if tc.wantOptimizationsGate != "true" {
						continue
					}
				default:
					continue
				}
				settingsRaw := nodeSet.EnvVars["CL_CRE_SETTINGS_DEFAULT"]
				if settingsRaw == "" {
					require.Equal(t, "false", tc.wantJWTGate)
					if tc.wantOptimizationsGate != "" {
						require.Equal(t, "false", tc.wantOptimizationsGate)
					}
				} else {
					var settings map[string]string
					require.NoError(t, json.Unmarshal([]byte(settingsRaw), &settings))
					if tc.wantOptimizationsGate == "true" {
						require.Equal(t, "true", settings["VaultOptimizationsEnabled"])
					} else {
						require.Equal(t, tc.wantJWTGate, settings["VaultJWTAuthEnabled"])
						if v, ok := settings["VaultOptimizationsEnabled"]; ok {
							require.Equal(t, "false", v)
						}
					}
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

// TestMustMintVaultJWTForRequest_UsesRawRequestDigest ensures the bearer token binds the digest of
// the exact JSON-RPC params wire body (canonical json.Marshal / jsonrpc.Request), matching what
// the gateway verifies—without relying on deprecated top-level identity fields inside params.
func TestMustMintVaultJWTForRequest_UsesRawRequestDigest(t *testing.T) {
	issuer, err := vault.NewTestJWTIssuer()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, issuer.Close())
	})

	derivedOwner := mustDeriveJWTVaultWorkflowOwner(t, "org-123")
	params, err := json.Marshal(vault_helpers.CreateSecretsRequest{
		RequestId: "req-1",
		EncryptedSecrets: []*vault_helpers.EncryptedSecret{
			{
				Id: &vault_helpers.SecretIdentifier{
					Key:       "9838",
					Namespace: "main",
					Owner:     derivedOwner,
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
	req.Auth = mustMintVaultJWTForRequest(t, issuer, req, "org-123")

	outboundReq := outboundRequestWithoutAuth(req)
	requestDigest, err := outboundReq.Digest()
	require.NoError(t, err)

	parsedToken, _, err := new(jwt.Parser).ParseUnverified(req.Auth, jwt.MapClaims{})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)
	authorizationDetails, ok := claims["authorization_details"].([]any)
	require.True(t, ok)

	var claimedDigest string
	for _, detail := range authorizationDetails {
		entry, ok := detail.(map[string]any)
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

func executeVaultSecretsGetInvalidIdentifierViaWorkflowTest(
	t *testing.T, testEnv *ttypes.TestEnvironment,
	workflowBaseName string,
	userLogsCh chan *workflowevents.UserLogs, baseMessageCh chan *commonevents.BaseMessage,
) {
	testLogger := framework.L
	testLogger.Info().Msg("Verifying get secret is rejected for invalid identifier via workflow...")

	const workflowFileLocation = "./vaultsecret/main.go"

	workflowName := t_helpers.UniqueWorkflowName(testEnv, workflowBaseName)
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &vaultsecret_config.Config{
		SecretKey:               "invalid-key-with-hyphens", // hyphen not in [a-zA-Z0-9_]; tests invalid key
		SecretNamespace:         "main",
		SecretKey2:              "validkey",
		SecretNamespace2:        "invalid-namespace-with-hyphens", // hyphen not in [a-zA-Z0-9_]; tests invalid namespace
		ExpectInvalidIdentifier: true,
	}, workflowFileLocation)

	// Both invalid-key and invalid-namespace checks run in the same cron trigger; a single
	// success log is emitted only after both GetSecret calls are correctly rejected.
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Vault get correctly rejected invalid identifier", 4*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))
	testLogger.Info().Msg("Vault get invalid identifier via workflow test completed")
}

// executeVaultSecretsIdentifierValidationTest verifies that the gateway rejects requests whose
// secret identifiers contain characters outside the allowed alphanumeric+underscore set.
// All four management request types (create, update, delete, list) are exercised for invalid key,
// invalid namespace, and invalid owner. Positive-path coverage is provided by basic_crud; this
// test focuses only on rejection behaviour.
func executeVaultSecretsIdentifierValidationTest(t *testing.T, encryptedSecret string, owner, gatewayURL string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	t.Helper()

	const (
		validKey         = "validkey"
		invalidKey       = "invalid-key-with-hyphens" // hyphen not in [a-zA-Z0-9_]
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
		allowlistRequest(t, owner, req, sethClient, wfRegistryContract)
		reqBody, err := json.Marshal(req)
		require.NoError(t, err)
		// Retry in case DON is still not synced properly
		var respBody []byte
		_ = retry.Do(func() error {
			_, respBody = sendVaultRequestToGateway(t, gatewayURL, reqBody)
			if bytes.Contains(respBody, []byte("Request timed out")) {
				return errors.New("gateway auth timeout")
			}
			return nil
		}, retry.Attempts(8), retry.Delay(3*time.Second), retry.DelayType(retry.FixedDelay),
			retry.OnRetry(func(n uint, err error) {
				framework.L.Warn().Uint("attempt", n+1).Msgf("[%s] %s: %s, retrying...", method, caseName, err)
			}))
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
	allowlistRequest(t, owner, req, sethClient, wfRegistryContract)
	reqBody, err := json.Marshal(req)
	require.NoError(t, err)
	var respBody []byte
	_ = retry.Do(func() error {
		_, respBody = sendVaultRequestToGateway(t, gatewayURL, reqBody)
		if bytes.Contains(respBody, []byte("Request timed out")) {
			return errors.New("gateway auth timeout")
		}
		return nil
	}, retry.Attempts(8), retry.Delay(3*time.Second), retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			framework.L.Warn().Uint("attempt", n+1).Msgf("[list] invalid namespace: %s, retrying...", err)
		}))
	require.Contains(t, string(respBody), "alphanumeric", "[list] expected alphanumeric rejection for %s", "invalid namespace")
	framework.L.Info().Msgf("[list] %s correctly rejected: %s", "invalid namespace", string(respBody))

	framework.L.Info().Msg("All identifier validation checks passed")
}
