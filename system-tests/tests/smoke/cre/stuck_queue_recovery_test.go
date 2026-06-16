package cre

import (
	"context"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/require"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
)

// ExecuteVaultStuckQueueRecoverySmokeTest wedges OCR by submitting one create whose ciphertext is
// accepted at ingress but produces a round-2 store-backed observation larger than
// VaultMaxObservationSizeLimit (requires VaultOptimizationsEnabled=false). After the stuck-round
// threshold fires and the pending queue is purged, a normal-sized create must succeed.
func ExecuteVaultStuckQueueRecoverySmokeTest(t *testing.T, fixture *vaultScenarioFixture, testEnv *ttypes.TestEnvironment) {
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

	t.Run("large_ciphertext_wedge_times_out", func(t *testing.T) {
		wedgeSecretID := uniqueVaultSecretID("stuckwedge")
		wedgeEnc := mustEncryptStuckQueueWedgeSecret(t, vaultParsedPublicKey, workflowOwnerAddress)
		submitStuckQueueWedgeCreateExpectingTimeout(t, auth, wedgeEnc, wedgeSecretID, owner, gwURL)
	})

	t.Run("stuck_threshold_and_pending_queue_purge_observed_in_docker_logs", func(t *testing.T) {
		assertVaultOCRStuckQueueRecoveryObservedInDockerLogs(t)
	})

	t.Run("ocr_liveness_after_recovery", func(t *testing.T) {
		recoverySecretID := uniqueVaultSecretID("afterrecovery")
		recoveryPlaintext := "secret-after-stuck-queue-recovery"
		recoveryEnc, err := encryptVaultSecretForOwner(t, recoveryPlaintext, vaultParsedPublicKey, workflowOwnerAddress)
		require.NoError(t, err)
		executeVaultAllowListSecretsCreateTest(t, recoveryEnc, recoverySecretID, owner, owner, gwURL, []string{"main"}, sc, wfReg)
		executeVaultSecretsListTest(t, recoverySecretID, owner, owner, gwURL, "main", sc, wfReg)
	})
}

func mustEncryptStuckQueueWedgeSecret(t *testing.T, pk *tdh2easy.PublicKey, owner common.Address) string {
	t.Helper()
	plaintextSize := pickStuckQueueWedgePlaintextSize(t, pk, owner)
	return encryptVaultSecretOfByteSize(t, pk, owner, plaintextSize)
}

func submitStuckQueueWedgeCreateExpectingTimeout(
	t *testing.T,
	auth vaultRequestAuth,
	encryptedSecret, secretID, owner, gatewayURL string,
) {
	t.Helper()

	requestID := uuid.New().String()
	secretsCreateRequest := vault_helpers.CreateSecretsRequest{
		RequestId:        requestID,
		EncryptedSecrets: buildEncryptedSecrets(secretID, owner, encryptedSecret, []string{"main"}),
	}
	jsonRequest := newVaultJSONRequest(t, requestID, vaulttypes.MethodSecretsCreate, &secretsCreateRequest)
	auth.apply(t, &jsonRequest)

	jsonResponse := sendVaultSignedOCRRequestToGateway(t, gatewayURL, jsonRequest)
	if jsonResponse.ID == "" {
		framework.L.Info().Str("requestID", requestID).Msg("stuck queue wedge create: gateway-to-DON timeout (expected)")
		return
	}
	if jsonResponse.Error != nil {
		framework.L.Info().Str("requestID", requestID).Str("error", jsonResponse.Error.Message).Msg("stuck queue wedge create: gateway error (acceptable while wedged)")
		return
	}
	require.Fail(t, "expected stuck queue wedge create to time out or error before OCR completed", "requestID=%s method=%s", requestID, jsonResponse.Method)
}

var (
	vaultStuckThresholdDockerLogRE = regexp.MustCompile(`pending queue skipped: stuck-round threshold reached`)
	vaultStuckCountDockerLogRE     = regexp.MustCompile(`count[^\d]*(\d+)`)
)

func assertVaultOCRStuckQueueRecoveryObservedInDockerLogs(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("docker log scan skipped in -short mode")
	}
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		t.Skip("docker not in PATH; skipping stuck-queue recovery log assertion")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	sawPendingSeed := false
	sawStuckThreshold := false

	for {
		logs := combinedVaultOCRContainerLogs(ctx, t, dockerBin)
		if !sawPendingSeed {
			for line := range strings.SplitSeq(logs, "\n") {
				if !strings.Contains(line, "pending queue items persisted to storage") {
					continue
				}
				wm := vaultStateTransitionPendingWriteWrittenRE.FindStringSubmatch(line)
				if len(wm) < 2 {
					continue
				}
				written, err1 := strconv.Atoi(wm[1])
				if err1 == nil && written >= 1 {
					sawPendingSeed = true
					framework.L.Info().Int("writtenCount", written).Msg("stuck queue recovery: observed KV pending seed write")
					break
				}
			}
		}

		if !sawStuckThreshold {
			for line := range strings.SplitSeq(logs, "\n") {
				if !vaultStuckThresholdDockerLogRE.MatchString(line) {
					continue
				}
				cm := vaultStuckCountDockerLogRE.FindStringSubmatch(line)
				if len(cm) >= 2 {
					count, err1 := strconv.Atoi(cm[1])
					if err1 == nil && count >= 3 {
						sawStuckThreshold = true
						framework.L.Info().Int("count", count).Msg("stuck queue recovery: observed stuck-round threshold log")
						break
					}
				}
				sawStuckThreshold = true
				framework.L.Info().Str("line", line).Msg("stuck queue recovery: observed stuck-round threshold log")
				break
			}
		}

		if sawStuckThreshold {
			for line := range strings.SplitSeq(logs, "\n") {
				if !strings.Contains(line, "pending queue items persisted to storage") {
					continue
				}
				wm := vaultStateTransitionPendingWriteWrittenRE.FindStringSubmatch(line)
				if len(wm) < 2 {
					continue
				}
				written, err1 := strconv.Atoi(wm[1])
				if err1 == nil && written == 0 {
					framework.L.Info().Msg("stuck queue recovery: observed pending queue purge (writtenCount=0)")
					return
				}
			}
		}

		select {
		case <-ctx.Done():
			require.Fail(t, "timed out waiting for stuck-queue recovery logs",
				"sawPendingSeed=%v sawStuckThreshold=%v (need pending write writtenCount>=1, stuck threshold log, then purge writtenCount=0)",
				sawPendingSeed, sawStuckThreshold)
		case <-ticker.C:
		}
	}
}

func combinedVaultOCRContainerLogs(ctx context.Context, t *testing.T, dockerBin string) string {
	t.Helper()
	psOut, err := exec.CommandContext(ctx, dockerBin, "ps", "--format", "{{.Names}}").Output()
	if err != nil {
		return ""
	}
	var combined strings.Builder
	for name := range strings.SplitSeq(strings.TrimSpace(string(psOut)), "\n") {
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
	return combined.String()
}
