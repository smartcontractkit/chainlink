package cre

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

// ExecuteVaultNodeSettingsConsensusSmokeTest runs the focused coverage for the
// VaultNodeSettingsConsensusEnabled topology: a minimal allowlist CRUD sanity check (proving the
// gateway -> vault OCR path works end to end under consensus-committed DON settings) plus the
// DON settings consensus docker-log assertion. The full allowlist and JWT suites already run on
// the default vault topology within the same suite bucket; repeating them here cannot fit the
// bucket's 7m go test timeout (Test_CRE_V2_Suite_Bucket_B on this topology previously timed out
// with don_settings_consensus still pending and ~20 parallel subtests never started).
func ExecuteVaultNodeSettingsConsensusSmokeTest(t *testing.T, fixture *vaultScenarioFixture, testEnv *ttypes.TestEnvironment) {
	t.Helper()

	gwURL := fixture.GatewayURL.String()
	vaultParsedPublicKey := mustVaultPublicKey(t, fixture.VaultPublicKey)

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

	auth := newAllowlistVaultRequestAuth(owner, sc, wfReg)

	t.Run("allowlist_crud_sanity", func(t *testing.T) {
		secretID := uniqueVaultSecretID("consensussmoke")
		createEnc, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-consensus-smoke-create", vaultParsedPublicKey, workflowOwnerAddress)
		require.NoError(t, err)
		updateEnc, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-consensus-smoke-update", vaultParsedPublicKey, workflowOwnerAddress)
		require.NoError(t, err)
		namespaces := []string{"main"}

		executeVaultSecretsCreateWithAuth(t, auth, createEnc, secretID, owner, gwURL, namespaces)
		executeVaultSecretsListWithAuth(t, auth, []string{secretID}, owner, gwURL, "main")
		executeVaultSecretsUpdateWithAuth(t, auth, updateEnc, secretID, owner, gwURL, namespaces)
		executeVaultSecretsDeleteWithAuth(t, auth, secretID, owner, gwURL, namespaces)
	})

	t.Run("don_settings_consensus", func(t *testing.T) {
		assertVaultDONSettingsQuorumInDockerLogs(t)
	})
}

// assertVaultDONSettingsQuorumInDockerLogs polls chainlink-related container logs until it observes
// a 'DON settings committed/updated from per-field observation quorum' log line emitted by the
// vault OCR state transition when it commits DON-wide settings via per-field 2f+1 consensus.
// These lines are logged at info level: containers run at the default info log level, so
// debug-level quorum detail lines are not visible in docker logs.
//
// The commit line is logged ONCE (the initial seed commit) shortly after the vault OCR DON starts
// producing rounds, so a bounded --tail window cannot be used: background log volume (head
// tracking, per-round protocol lines; hundreds of lines per second per node) scrolls a one-shot
// line out of any --tail window within minutes. Each container is therefore scanned from its full
// log the first time it is seen and incrementally (--since the end of the last successful scan)
// afterwards, so the line stays in scope no matter when it was written or how large the logs grow.
func assertVaultDONSettingsQuorumInDockerLogs(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("docker log scan skipped in -short mode")
	}
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		t.Skip("docker not in PATH; skipping DON settings consensus log scan")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	scannedThrough := make(map[string]time.Time)
	for {
		psOut, err := exec.CommandContext(ctx, dockerBin, "ps", "--format", "{{.Names}}").Output()
		if err != nil {
			select {
			case <-ctx.Done():
				require.Fail(t, "timed out waiting for docker while scanning for DON settings quorum commit log")
			case <-ticker.C:
			}
			continue
		}
		for name := range strings.SplitSeq(strings.TrimSpace(string(psOut)), "\n") {
			if name == "" {
				continue
			}
			ln := strings.ToLower(name)
			if !strings.Contains(ln, "chainlink") && !strings.Contains(ln, "ocr") && !strings.Contains(ln, "capabilit") {
				continue
			}
			// Sample the window end before scanning: lines written after this instant are covered
			// by the next pass, lines written before it are covered by this scan.
			windowEnd := time.Now().UTC()
			logArgs := []string{"logs", name}
			if since, ok := scannedThrough[name]; ok {
				logArgs = append(logArgs, "--since", since.Format(time.RFC3339))
			}
			logs, err := exec.CommandContext(ctx, dockerBin, logArgs...).CombinedOutput()
			if err != nil {
				// Window not marked scanned, so the same range is retried on the next pass.
				continue
			}
			scannedThrough[name] = windowEnd
			logsStr := string(logs)
			if strings.Contains(logsStr, "DON settings committed from per-field observation quorum") ||
				strings.Contains(logsStr, "DON settings updated from per-field observation quorum") {
				framework.L.Info().Str("container", name).Msg("observed vault OCR DON settings quorum commit log")
				return
			}
		}
		select {
		case <-ctx.Done():
			require.Fail(t, "timed out waiting for DON settings quorum commit log line in docker logs (is the local CRE stack running with VaultNodeSettingsConsensusEnabled?)")
		case <-ticker.C:
		}
	}
}
