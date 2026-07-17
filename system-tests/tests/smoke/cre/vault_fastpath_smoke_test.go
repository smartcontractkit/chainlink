package cre

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

var (
	vaultFastPathSentToBufferLogRE  = regexp.MustCompile(`sent request to fast-path buffer`)
	vaultSlowPathSentToHandlerLogRE = regexp.MustCompile(`sent request to OCR handler`)
)

// ExecuteVaultFastPathAssertions runs fast-path-specific checks that are orthogonal to the
// allowlist-based CRUD tests already exercised on the fast-path topology. It proves the workflow
// GetSecrets calls actually traversed the fast-path buffer (not the OCR handler) and that the
// fast-path e2e happy path returns the correct secret.
func ExecuteVaultFastPathAssertions(t *testing.T, fixture *vaultScenarioFixture, testEnv *ttypes.TestEnvironment) {
	gwURL := fixture.GatewayURL.String()
	vaultPublicKey := fixture.VaultPublicKey

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	workflowOwnerAddress := sc.MustGetRootKeyAddress()
	owner := workflowOwnerAddress.Hex()

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

	vaultParsedPublicKey := mustVaultPublicKey(t, vaultPublicKey)
	secretID := uniqueVaultSecretID("fastpath")
	createValue := "secret-fast-path"
	createEnc, err := vaultutils.EncryptSecretWithWorkflowOwner(createValue, vaultParsedPublicKey, workflowOwnerAddress)
	require.NoError(t, err)

	ulCh := make(chan *workflowevents.UserLogs, 1000)
	bmCh := make(chan *commonevents.BaseMessage, 1000)
	sink := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(framework.L, ulCh, bmCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, sink, ulCh, bmCh)
	})

	executeVaultAllowListSecretsCreateTest(t, createEnc, secretID, owner, owner, gwURL, []string{"main"}, sc, wfReg)

	t.Run("fast_path_happy_path_returns_secret", func(t *testing.T) {
		workflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, "fast-path-run", []vaultWorkflowPhase{{
			Name: "fast-path-fetch",
			Checks: []vaultWorkflowCheck{{
				Name:            "fast-path-main",
				SecretKey:       secretID,
				SecretNamespace: "main",
				ExpectedValue:   createValue,
			}},
		}})
		waitForVaultWorkflowPhase(t, workflowID, "fast-path-fetch", ulCh, bmCh)
	})

	t.Run("fast_path_uses_buffer_not_ocr_handler", func(t *testing.T) {
		// Re-trigger the same workflow to ensure we observe the fast-path log lines.
		workflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, "fast-path-verify", []vaultWorkflowPhase{{
			Name: "fast-path-verify-fetch",
			Checks: []vaultWorkflowCheck{{
				Name:            "fast-path-verify-main",
				SecretKey:       secretID,
				SecretNamespace: "main",
				ExpectedValue:   createValue,
			}},
		}})
		waitForVaultWorkflowPhase(t, workflowID, "fast-path-verify-fetch", ulCh, bmCh)
		assertVaultFastPathBufferLogObservedInDockerLogs(t)
		assertVaultNoSlowPathGetSecretsLogObserved(t)
	})

	executeVaultSecretsDeleteTest(t, secretID, owner, owner, gwURL, []string{"main"}, sc, wfReg)
}

// assertVaultFastPathBufferLogObservedInDockerLogs polls chainlink/ocr/capability container logs
// until it observes the "sent request to fast-path buffer" debug log emitted by the vault
// capability when VaultFastPathGetSecretsEnabled is on. This proves GetSecrets read requests are
// bypassing the OCR request store.
func assertVaultFastPathBufferLogObservedInDockerLogs(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("docker log scan skipped in -short mode")
	}
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		t.Skip("docker not in PATH; skipping fast-path buffer log assertion")
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
				require.Fail(t, "timed out waiting for docker while scanning for fast-path buffer log")
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
				if vaultFastPathSentToBufferLogRE.MatchString(line) {
					framework.L.Info().Str("container", name).Str("line", line).Msg("observed vault fast-path buffer log")
					return
				}
			}
		}
		select {
		case <-ctx.Done():
			require.Fail(t, "timed out waiting for 'sent request to fast-path buffer' log in docker logs (is the local CRE stack running with VaultFastPathGetSecretsEnabled?)")
		case <-ticker.C:
		}
	}
}

// assertVaultNoSlowPathGetSecretsLogObserved scans container logs for the slow-path "sent request to
// OCR handler" log line. This is a defensive check: if any GetSecrets request mistakenly entered the
// OCR request store on a fast-path topology, it should be detectable here. Note that writes
// (Create/Update/Delete) legitimately use the OCR handler, so we only flag lines that also reference
// the GetSecrets method.
func assertVaultNoSlowPathGetSecretsLogObserved(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("docker log scan skipped in -short mode")
	}
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		t.Skip("docker not in PATH; skipping slow-path log assertion")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	psOut, err := exec.CommandContext(ctx, dockerBin, "ps", "--format", "{{.Names}}").Output()
	require.NoError(t, err, "failed to list docker containers")
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
			if !vaultSlowPathSentToHandlerLogRE.MatchString(line) {
				continue
			}
			// The slow-path log itself is method-agnostic; gate it on the same debug line also
			// carrying a GetSecrets request method.
			if strings.Contains(line, vaulttypes.MethodSecretsGet) || strings.Contains(line, "GetSecrets") {
				require.Fail(t, "observed slow-path 'sent request to OCR handler' log for a GetSecrets request on a fast-path topology", "container=%s line=%s", name, line)
			}
		}
	}
}
