package cre

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

const vaultNodesetName = "capabilities"

// ExecutionTimestampsEnabled must be on for direct reads: getSecretsDirectlyEnabled gates on
// the execution timestamp, which is only assigned when this setting allows (base_engine.go).
const vaultDirectReadEnabledTOML = `[PerWorkflow]
FeatureVaultGetSecretsDirectlyActivePeriod = '[2000-01-01 00:00:00 +0000 UTC,2100-01-01 00:00:00 +0000 UTC]'
ExecutionTimestampsEnabled = 'true'`

var vaultDirectFlagRE = regexp.MustCompile(`get_secrets_directly:\s*true`)

// vaultDirectReadEnv is what the direct read tests share: a verifier workflow of
// their own, so the flag override doesn't leak into other vault tests.
type vaultDirectReadEnv struct {
	testEnv  *ttypes.TestEnvironment
	verifier *vaultVerifierHandle
	owner    string
	gwURL    string
	sc       *seth.Client
	wfReg    *workflow_registry_v2_wrapper.WorkflowRegistry
	ulCh     chan *workflowevents.UserLogs
	bmCh     chan *commonevents.BaseMessage
}

func newVaultDirectReadEnv(t *testing.T, fixture *vaultScenarioFixture, testEnv *ttypes.TestEnvironment, workflowBaseName string) *vaultDirectReadEnv {
	t.Helper()

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	owner := sc.MustGetRootKeyAddress().Hex()
	if fixture.LinkingService != nil {
		fixture.LinkingService.SetOwnerOrg(owner, uniqueVaultSecretID("org"))
	}
	wfRegVersion := testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()]
	wfRegAddr := crecontracts.MustGetAddressFromDataStore(testEnv.CreEnvironment.CldfEnvironment.DataStore, testEnv.CreEnvironment.Blockchains[0].ChainSelector(), keystone_changeset.WorkflowRegistry.String(), wfRegVersion, "")
	wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(common.HexToAddress(wfRegAddr), sc.Client)
	require.NoError(t, err)
	requireVaultLinkOwner(t, sc, common.HexToAddress(wfRegAddr), wfRegVersion)

	ulCh := make(chan *workflowevents.UserLogs, 1000)
	bmCh := make(chan *commonevents.BaseMessage, 1000)
	sink := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(framework.L, ulCh, bmCh))
	t.Cleanup(func() {
		// can't use t.Context() here because it will have been cancelled before the cleanup function is called
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, sink, ulCh, bmCh)
	})

	return &vaultDirectReadEnv{
		testEnv:  testEnv,
		verifier: deployVaultVerifierWorkflow(t, testEnv, fixture.TriggerAuth, workflowBaseName),
		owner:    owner,
		gwURL:    fixture.GatewayURL.String(),
		sc:       sc,
		wfReg:    wfReg,
		ulCh:     ulCh,
		bmCh:     bmCh,
	}
}

func (e *vaultDirectReadEnv) enableDirectRead(t *testing.T) *t_helpers.CRESettingsHandle {
	t.Helper()
	return t_helpers.ApplyCRESettings(t, e.testEnv, t_helpers.Workflow(strings.TrimPrefix(e.verifier.WorkflowID, "0x"), vaultDirectReadEnabledTOML))
}

func (e *vaultDirectReadEnv) createSecret(t *testing.T, fixture *vaultScenarioFixture, secretID, value string, namespaces ...string) {
	t.Helper()
	enc, err := vaultutils.EncryptSecretWithWorkflowOwner(value, mustVaultPublicKey(t, fixture.VaultPublicKey), common.HexToAddress(e.owner))
	require.NoError(t, err)
	executeVaultAllowListSecretsCreateTest(t, enc, secretID, e.owner, e.owner, e.gwURL, namespaces, e.sc, e.wfReg)
}

func (e *vaultDirectReadEnv) updateSecret(t *testing.T, fixture *vaultScenarioFixture, secretID, value string, namespaces ...string) {
	t.Helper()
	enc, err := vaultutils.EncryptSecretWithWorkflowOwner(value, mustVaultPublicKey(t, fixture.VaultPublicKey), common.HexToAddress(e.owner))
	require.NoError(t, err)
	executeVaultSecretsUpdateTest(t, enc, secretID, e.owner, e.owner, e.gwURL, namespaces, e.sc, e.wfReg)
}

func (e *vaultDirectReadEnv) awaitPhase(t *testing.T, phase vaultWorkflowPhase) {
	t.Helper()
	triggerAndAwaitVaultWorkflowPhase(t, e.verifier, phase, e.ulCh, e.bmCh)
}

// vaultGetSecretsLogged reports whether a Vault node logged a GetSecrets request for secretID,
// with the direct read flag set or not.
func vaultGetSecretsLogged(t *testing.T, testEnv *ttypes.TestEnvironment, secretID string, direct bool) bool {
	t.Helper()
	return t_helpers.ContainerLogLineMatchesForNodeset(t, testEnv, vaultNodesetName, func(line string) bool {
		return strings.Contains(line, "received workflow get secrets request") &&
			strings.Contains(line, secretID) &&
			vaultDirectFlagRE.MatchString(line) == direct
	})
}

func requireVaultGetSecretsLoggedEventually(t *testing.T, testEnv *ttypes.TestEnvironment, secretID string, direct bool) {
	t.Helper()
	require.Eventually(t, func() bool {
		return vaultGetSecretsLogged(t, testEnv, secretID, direct)
	}, time.Minute, 2*time.Second, "expected a vault GetSecrets request for %s with direct=%t", secretID, direct)
}

// ExecuteVaultDirectReadTests covers GetSecrets served from the Vault nodes' local
// replicated state. The subtests run in order and share the flag override.
func ExecuteVaultDirectReadTests(t *testing.T, fixture *vaultScenarioFixture, testEnv *ttypes.TestEnvironment) {
	env := newVaultDirectReadEnv(t, fixture, testEnv, "direct-read-verifier")
	flag := env.enableDirectRead(t)

	secretA := uniqueVaultSecretID("directa")
	secretB := uniqueVaultSecretID("directb")

	t.Run("happy_path", func(t *testing.T) {
		env.createSecret(t, fixture, secretA, "direct-a-v1", "main", "alt")
		env.createSecret(t, fixture, secretB, "direct-b-v1", "main")

		env.awaitPhase(t, vaultWorkflowPhase{
			Name: "direct-single",
			Checks: []vaultWorkflowCheck{
				{Name: "a-main", SecretKey: secretA, SecretNamespace: "main", ExpectedValue: "direct-a-v1"},
				{Name: "a-alt", SecretKey: secretA, SecretNamespace: "alt", ExpectedValue: "direct-a-v1"},
				{Name: "b-main", SecretKey: secretB, SecretNamespace: "main", ExpectedValue: "direct-b-v1"},
			},
		})
		env.awaitPhase(t, vaultWorkflowPhase{
			Name:  "direct-batch",
			Batch: true,
			Checks: []vaultWorkflowCheck{
				{Name: "a-main", SecretKey: secretA, SecretNamespace: "main", ExpectedValue: "direct-a-v1"},
				{Name: "b-main", SecretKey: secretB, SecretNamespace: "main", ExpectedValue: "direct-b-v1"},
				{Name: "a-alt", SecretKey: secretA, SecretNamespace: "alt", ExpectedValue: "direct-a-v1"},
			},
		})

		requireVaultGetSecretsLoggedEventually(t, testEnv, secretA, true)
		requireVaultGetSecretsLoggedEventually(t, testEnv, secretB, true)
		t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, vaultNodesetName, "served direct get secrets request", time.Minute, 2*time.Second)
	})

	t.Run("read_after_write", func(t *testing.T) {
		env.updateSecret(t, fixture, secretA, "direct-a-v2", "main")
		env.awaitPhase(t, vaultWorkflowPhase{
			Name:   "direct-updated",
			Checks: []vaultWorkflowCheck{{Name: "a-main", SecretKey: secretA, SecretNamespace: "main", ExpectedValue: "direct-a-v2"}},
		})

		executeVaultSecretsDeleteTest(t, secretB, env.owner, env.owner, env.gwURL, []string{"main"}, env.sc, env.wfReg)
		env.awaitPhase(t, vaultWorkflowPhase{
			Name:   "direct-deleted",
			Checks: []vaultWorkflowCheck{{Name: "b-main", SecretKey: secretB, SecretNamespace: "main", ExpectNotFound: true}},
		})
	})

	t.Run("errors_match_ocr_path", func(t *testing.T) {
		executeVaultSecretsGetInvalidIdentifierViaWorkflowTest(t, env.verifier, env.ulCh, env.bmCh)
		executeVaultSecretsGetBatchTooBigViaWorkflowTest(t, env.verifier, env.ulCh, env.bmCh)
	})

	t.Run("flag_off_uses_ocr", func(t *testing.T) {
		flag.Reset(t)

		secretC := uniqueVaultSecretID("directc")
		env.createSecret(t, fixture, secretC, "direct-c-v1", "main")
		env.awaitPhase(t, vaultWorkflowPhase{
			Name:   "ocr-read",
			Checks: []vaultWorkflowCheck{{Name: "c-main", SecretKey: secretC, SecretNamespace: "main", ExpectedValue: "direct-c-v1"}},
		})

		// Only positive checks: settings rollout is best-effort, so early retries may still see the old flag.
		requireVaultGetSecretsLoggedEventually(t, testEnv, secretC, false)
	})
}
