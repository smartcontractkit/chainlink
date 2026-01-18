package llo_e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// LogEvidence captures proof of data at each stage
// Note: MAGIC_NUMBER, MAGIC_NUMBER_FORMAT5, MAGIC_NUMBER_FORMAT7 are defined in llo_e2e.go
type LogEvidence struct {
	Stage      string
	Container  string
	LineNumber int
	LogSnippet string
	KeyValue   string
}

// TestLLOStreamsE2E runs the full LLO Streams Trigger E2E test
//
// This test proves the complete data flow:
//
//	Mock EA (424242) → Stream Jobs → LLO Plugin → OCR Consensus
//	                                                    ↓
//	                            CRE Transmitter → streams-trigger@2.0.0
//	                                                    ↓
//	                            TriggerPublisher → Workflow DON
//	                                                    ↓
//	                            llo-consumer workflow → LLO_E2E_VALUE logs
//
// KEY ASSERTION: The workflow outputs Value=424242, proving it can read the value
func TestLLOStreamsE2E(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	envDir := getEnvironmentDir()
	t.Logf("Environment directory: %s", envDir)

	// Evidence collection for final summary
	var evidence []LogEvidence

	// Step 1: Start environment if not running
	t.Log("Step 1: Checking/Starting CRE Environment...")
	if !isEnvironmentRunning() {
		t.Log("Starting environment...")
		cmd := exec.CommandContext(ctx, "go", "run", ".", "env", "start")
		cmd.Dir = envDir
		cmd.Env = append(os.Environ(), "CTF_CONFIGS=configs/streams-don-to-don-e2e.toml")
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "Failed to start environment: %s", string(out))
	}
	t.Log("✓ Environment running")

	// Step 2: Start Mock EA
	t.Log("Step 2: Starting Mock EA...")
	err := startMockEA(ctx)
	require.NoError(t, err, "Failed to start Mock EA")
	t.Log("✓ Mock EA running")

	// Step 3: Deploy LLO contracts
	t.Log("Step 3: Deploying LLO Contracts...")
	contracts, err := deployLLOContracts(ctx)
	require.NoError(t, err, "Failed to deploy LLO contracts")
	t.Logf("✓ Configurator: %s", contracts.ConfiguratorAddr.Hex())
	t.Logf("✓ ChannelConfigStore: %s", contracts.ChannelConfigStoreAddr.Hex())

	// Step 4: Collect node info
	t.Log("Step 4: Collecting Streams DON Node Information...")
	nodeInfos, err := collectNodeInfo(ctx)
	require.NoError(t, err, "Failed to collect node info")
	require.Len(t, nodeInfos, numStreamsNodes, "Expected %d nodes", numStreamsNodes)
	t.Logf("✓ Collected info from %d nodes", len(nodeInfos))

	// Step 5: Configure OCR
	t.Log("Step 5: Configuring OCR on Configurator Contract...")
	err = configureOCR(ctx, contracts, nodeInfos)
	require.NoError(t, err, "Failed to configure OCR")
	t.Log("✓ OCR configuration set")

	// Step 6: Create bridges
	t.Log("Step 6: Creating Bridges to Mock EA...")
	err = createBridges(ctx, nodeInfos)
	require.NoError(t, err, "Failed to create bridges")
	t.Log("✓ Bridges created")

	// Step 7: Deploy stream jobs
	t.Log("Step 7: Deploying Stream Jobs...")
	err = deployStreamJobs(ctx, nodeInfos)
	require.NoError(t, err, "Failed to deploy stream jobs")
	t.Log("✓ Stream jobs deployed")

	// Step 8: Deploy LLO jobs
	t.Log("Step 8: Deploying LLO Jobs with CRE Transmitter...")
	err = deployLLOJobs(ctx, contracts, nodeInfos)
	require.NoError(t, err, "Failed to deploy LLO jobs")
	t.Log("✓ LLO jobs deployed")

	// Wait for jobs to start, then re-emit config to ensure LogPoller catches it
	t.Log("Step 8b: Waiting for jobs to start and re-emitting OCR config...")
	time.Sleep(10 * time.Second)

	// Re-emit OCR config - this ensures the LogPoller (which starts with the LLO jobs)
	// definitely sees the config event
	err = configureOCR(ctx, contracts, nodeInfos)
	if err != nil {
		t.Logf("Warning: Failed to re-emit OCR config: %v (continuing anyway)", err)
	} else {
		t.Log("✓ OCR config re-emitted for LogPoller")
	}

	// Wait for config to be picked up
	t.Log("Step 8c: Waiting for LLO to pick up config...")
	err = waitForLLOConfigPickup(ctx)
	if err != nil {
		t.Logf("Warning: LLO config pickup not confirmed: %v (continuing anyway)", err)
	} else {
		t.Log("✓ LLO config picked up")
	}

	// Step 9: Verify Mock EA requests and capture evidence
	t.Log("Step 9: Verifying Mock EA is receiving requests...")
	err = verifyMockEARequests(ctx)
	require.NoError(t, err, "Mock EA not receiving requests")
	mockEAEvidence := captureLogEvidence(ctx, "mock-ea", fmt.Sprintf("%d", MAGIC_NUMBER_FORMAT5), "STAGE 1: Mock EA (Format 5)")
	if mockEAEvidence != nil {
		evidence = append(evidence, *mockEAEvidence)
		t.Log("✓ Mock EA is receiving requests")
		t.Log("")
		t.Log("  ┌─────────────────────────────────────────────────────────────────────────────┐")
		t.Logf("  │ 📍 STAGE 1: Mock EA returning MAGIC_NUMBER_FORMAT5=%d                    │", MAGIC_NUMBER_FORMAT5)
		t.Log("  ├─────────────────────────────────────────────────────────────────────────────┤")
		t.Logf("  │ 📄 %s:%d", mockEAEvidence.Container, mockEAEvidence.LineNumber)
		t.Log("  │")
		t.Log("  │ Exact command:")
		t.Logf("  │   docker logs %s 2>&1 | sed -n '%dp'", mockEAEvidence.Container, mockEAEvidence.LineNumber)
		t.Log("  └─────────────────────────────────────────────────────────────────────────────┘")
		t.Log("")
	}

	// Step 10: Wait for LLO reports (evidence captured later after workflow deployed)
	t.Log("Step 10: Waiting for LLO Reports and CRE Transmitter...")
	err = waitForLLOReports(ctx)
	require.NoError(t, err, "Failed to detect LLO reports")
	t.Log("✓ LLO is producing reports")

	// Step 11: Deploy and verify workflow
	t.Log("Step 11: Deploying consumer workflow...")
	err = deployAndVerifyWorkflow(ctx)
	if err != nil {
		t.Logf("Warning: Workflow deployment issue: %v", err)
	}

	// Step 12: THE KEY ASSERTION - Verify workflow received and decoded the MAGIC_NUMBER
	t.Log("Step 12: Verifying workflow decoded MAGIC_NUMBER from Streams DON...")
	reports, err := waitForValueReports(ctx, t, 2*time.Minute)
	require.NoError(t, err, "Failed to find LLO_E2E_VALUE logs")
	require.NotEmpty(t, reports, "No LLO_E2E_VALUE logs found")

	// Capture workflow evidence for Format 5
	format5Evidence := captureLogEvidenceByFormat(ctx, "workflow-node0", 5, "STAGE 3a: Workflow (Format 5 - protobuf)")
	if format5Evidence != nil {
		evidence = append(evidence, *format5Evidence)
		t.Log("")
		t.Log("  ┌─────────────────────────────────────────────────────────────────────────────┐")
		t.Logf("  │ 📍 STAGE 3a: Workflow decoded Format 5 - MAGIC_NUMBER=%d                 │", MAGIC_NUMBER_FORMAT5)
		t.Log("  ├─────────────────────────────────────────────────────────────────────────────┤")
		t.Logf("  │ 📄 %s:%d", format5Evidence.Container, format5Evidence.LineNumber)
		t.Log("  │")
		t.Log("  │ Exact command:")
		t.Logf("  │   docker logs %s 2>&1 | sed -n '%dp'", format5Evidence.Container, format5Evidence.LineNumber)
		t.Log("  │")
		t.Logf("  │ Key evidence: \"LLO_E2E_VALUE[SeqNr=...]: Format=5 Value=%d Match=true\"", MAGIC_NUMBER_FORMAT5)
		t.Log("  └─────────────────────────────────────────────────────────────────────────────┘")
		t.Log("")
	}

	// Capture workflow evidence for Format 7
	format7Evidence := captureLogEvidenceByFormat(ctx, "workflow-node0", 7, "STAGE 3b: Workflow (Format 7 - ABI)")
	if format7Evidence != nil {
		evidence = append(evidence, *format7Evidence)
		t.Log("")
		t.Log("  ┌─────────────────────────────────────────────────────────────────────────────┐")
		t.Logf("  │ 📍 STAGE 3b: Workflow decoded Format 7 - MAGIC_NUMBER=%d                 │", MAGIC_NUMBER_FORMAT7)
		t.Log("  ├─────────────────────────────────────────────────────────────────────────────┤")
		t.Logf("  │ 📄 %s:%d", format7Evidence.Container, format7Evidence.LineNumber)
		t.Log("  │")
		t.Log("  │ Exact command:")
		t.Logf("  │   docker logs %s 2>&1 | sed -n '%dp'", format7Evidence.Container, format7Evidence.LineNumber)
		t.Log("  │")
		t.Logf("  │ Key evidence: \"LLO_E2E_VALUE[SeqNr=...]: Format=7 Value=%d Match=true\"", MAGIC_NUMBER_FORMAT7)
		t.Log("  └─────────────────────────────────────────────────────────────────────────────┘")
		t.Log("")
	}

	// Capture Streams DON evidence NOW (after workflow subscribed, so nSubscribers > 0)
	streamsEvidence := captureLogEvidenceWithSubscribers(ctx, "streams-node0", "STAGE 2: Streams DON (delivering to workflow)")
	if streamsEvidence != nil {
		// Insert at position 1 (after Mock EA, before Workflow)
		evidence = append(evidence[:1], append([]LogEvidence{*streamsEvidence}, evidence[1:]...)...)
		t.Log("")
		t.Log("  ┌─────────────────────────────────────────────────────────────────────────────┐")
		t.Log("  │ 📍 STAGE 2: CRE Transmitter delivering events to Workflow DON              │")
		t.Log("  ├─────────────────────────────────────────────────────────────────────────────┤")
		t.Logf("  │ 📄 %s:%d", streamsEvidence.Container, streamsEvidence.LineNumber)
		t.Log("  │")
		t.Log("  │ Exact command:")
		t.Logf("  │   docker logs %s 2>&1 | sed -n '%dp'", streamsEvidence.Container, streamsEvidence.LineNumber)
		t.Log("  │")
		t.Log("  │ Key evidence: nSubscribers >= 1 (workflow actively receiving)")
		t.Log("  └─────────────────────────────────────────────────────────────────────────────┘")
		t.Log("")
	}

	// KEY ASSERTIONS - Verify the workflow can decode and read the MAGIC_NUMBER
	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════════")
	t.Log("  ASSERTING: Workflow decoded MAGIC_NUMBER from Streams DON")
	t.Log("═══════════════════════════════════════════════════════════════════")

	// Must have at least 2 reports to prove continuous data flow
	require.GreaterOrEqual(t, len(reports), 2,
		"Expected at least 2 LLO_E2E_VALUE reports to prove continuous data flow. Got %d.", len(reports))

	// Count how many reports have the correct value for each format
	matchCountF5 := 0
	matchCountF7 := 0
	for i, report := range reports {
		t.Logf("  Report %d: SeqNr=%d, Format=%d, Value=%d, Expected=%d, Match=%v, Sigs=%d",
			i+1, report.SeqNr, report.Format, report.Value, report.Expected, report.Match, report.Sigs)
		require.Greater(t, report.SeqNr, 0, "SeqNr should be > 0")

		// Check against appropriate magic number based on format
		if report.Format == 5 && report.Value == MAGIC_NUMBER_FORMAT5 {
			matchCountF5++
		} else if report.Format == 7 && report.Value == MAGIC_NUMBER_FORMAT7 {
			matchCountF7++
		} else if report.Format == 0 && report.Value == MAGIC_NUMBER {
			// backward compat for old format without Format field
			matchCountF5++
		}
	}

	// Log format breakdown
	t.Logf("  Format 5 matches: %d (expected value=%d)", matchCountF5, MAGIC_NUMBER_FORMAT5)
	t.Logf("  Format 7 matches: %d (expected value=%d)", matchCountF7, MAGIC_NUMBER_FORMAT7)

	// ASSERT: Both formats must have at least one match
	require.Greater(t, matchCountF5, 0,
		"Expected at least 1 Format 5 report with MAGIC_NUMBER_FORMAT5=%d. Got 0 F5 matches out of %d reports.",
		MAGIC_NUMBER_FORMAT5, len(reports))
	require.Greater(t, matchCountF7, 0,
		"Expected at least 1 Format 7 report with MAGIC_NUMBER_FORMAT7=%d. Got 0 F7 matches out of %d reports.",
		MAGIC_NUMBER_FORMAT7, len(reports))

	// Print final summary with all evidence
	t.Log("")
	t.Log("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Log("║  ✅ E2E TEST PASSED - BOTH FORMATS VERIFIED                                  ║")
	t.Log("╠══════════════════════════════════════════════════════════════════════════════╣")
	t.Logf("║  MAGIC_NUMBER_FORMAT5: %d (CapabilityTrigger - protobuf)                   ║", MAGIC_NUMBER_FORMAT5)
	t.Logf("║  MAGIC_NUMBER_FORMAT7: %d (EVMABIEncodeUnpackedExpr - ABI)                ║", MAGIC_NUMBER_FORMAT7)
	t.Logf("║  Format 5 matches: %d ✓ / Format 7 matches: %d ✓ / Total: %d              ║", matchCountF5, matchCountF7, len(reports))
	t.Logf("║  SeqNr range: %d → %d                                                       ║", reports[0].SeqNr, reports[len(reports)-1].SeqNr)
	t.Log("╠══════════════════════════════════════════════════════════════════════════════╣")
	t.Log("║  Data Flow Verified (Format 5 - protobuf):                                   ║")
	t.Logf("║    Mock EA (%d) → Stream 1 (TEST/USD) → Channel 1 (Format 5) ✓            ║", MAGIC_NUMBER_FORMAT5)
	t.Log("║                 → LLO Plugin → CRE Transmitter → streams-trigger@2.0.0      ║")
	t.Logf("║                 → llo-consumer → Value=%d ✓                               ║", MAGIC_NUMBER_FORMAT5)
	t.Log("║                                                                              ║")
	t.Log("║  Data Flow Verified (Format 7 - ABI):                                        ║")
	t.Logf("║    Mock EA (%d) → Stream 4 (DATA/USD) → Channel 2 (Format 7) ✓            ║", MAGIC_NUMBER_FORMAT7)
	t.Log("║                 → LLO Plugin → CRE Transmitter → streams-trigger@2.0.0      ║")
	t.Logf("║                 → llo-consumer → Value=%d ✓                               ║", MAGIC_NUMBER_FORMAT7)
	t.Log("╚══════════════════════════════════════════════════════════════════════════════╝")
	t.Log("")
	t.Log("📋 PROOF - Exact log lines used to verify each stage:")
	t.Log("")
	for _, e := range evidence {
		cmd := fmt.Sprintf("docker logs %s 2>&1 | sed -n '%dp'", e.Container, e.LineNumber)
		t.Logf("  # %s", e.Stage)
		t.Logf("  $ %s", cmd)
		out, err := exec.CommandContext(ctx, "bash", "-c", cmd).CombinedOutput()
		if err != nil {
			t.Logf("    (error: %v)", err)
		} else {
			t.Logf("    → %s", strings.TrimSpace(string(out)))
		}
		t.Log("")
	}
}

// ValueReport represents a parsed LLO_E2E_VALUE log entry
type ValueReport struct {
	SeqNr    int
	Format   int // 5 = CapabilityTrigger (protobuf), 7 = EVMABIEncodeUnpackedExpr (ABI)
	Value    int
	Expected int
	Match    bool
	Sigs     int // Number of signatures on the report
	Raw      string
}

// waitForValueReports waits for LLO_E2E_VALUE logs from the workflow nodes
func waitForValueReports(ctx context.Context, t *testing.T, timeout time.Duration) ([]ValueReport, error) {
	deadline := time.Now().Add(timeout)
	var reports []ValueReport
	seenSeqNrs := make(map[int]bool)

	for time.Now().Before(deadline) {
		for i := 0; i < numWorkflowNodes; i++ {
			containerName := fmt.Sprintf("workflow-node%d", i)
			// Search all logs - high volume containers may have LLO_E2E_VALUE buried deep
			cmd := exec.CommandContext(ctx, "docker", "logs", containerName)
			out, err := cmd.CombinedOutput()
			if err != nil {
				continue
			}

			logs := string(out)
			for _, line := range strings.Split(logs, "\n") {
				if strings.Contains(line, "LLO_E2E_VALUE[SeqNr=") {
					report := parseValueReport(line)
					if !seenSeqNrs[report.SeqNr] {
						seenSeqNrs[report.SeqNr] = true
						reports = append(reports, report)
						t.Logf("Found LLO_E2E_VALUE: SeqNr=%d Format=%d Value=%d Match=%v Sigs=%d",
							report.SeqNr, report.Format, report.Value, report.Match, report.Sigs)
					}
				}
			}
		}

		// Count matches per format - we need at least one of EACH format to prove both paths work
		matchCountF5 := 0
		matchCountF7 := 0
		for _, r := range reports {
			if r.Format == 5 && r.Value == MAGIC_NUMBER_FORMAT5 {
				matchCountF5++
			} else if r.Format == 7 && r.Value == MAGIC_NUMBER_FORMAT7 {
				matchCountF7++
			} else if r.Format == 0 && r.Value == MAGIC_NUMBER {
				matchCountF5++ // backward compat
			}
		}
		// Exit when we have at least one match from EACH format
		if matchCountF5 >= 1 && matchCountF7 >= 1 {
			return reports, nil
		}

		time.Sleep(5 * time.Second)
	}

	if len(reports) > 0 {
		return reports, nil
	}

	return nil, fmt.Errorf("no LLO_E2E_VALUE logs found within %v", timeout)
}

// parseValueReport parses a LLO_E2E_VALUE log line
// Format: LLO_E2E_VALUE[SeqNr=N]: Format=X Value=Y Expected=Z Match=W Sigs=V
func parseValueReport(line string) ValueReport {
	report := ValueReport{Raw: line}

	// Extract SeqNr
	if idx := strings.Index(line, "SeqNr="); idx != -1 {
		end := strings.Index(line[idx:], "]")
		if end != -1 {
			fmt.Sscanf(line[idx+6:idx+end], "%d", &report.SeqNr)
		}
	}

	// Extract Format (5 or 7)
	if idx := strings.Index(line, "Format="); idx != -1 {
		fmt.Sscanf(line[idx+7:], "%d", &report.Format)
	}

	// Extract Value
	if idx := strings.Index(line, "Value="); idx != -1 {
		fmt.Sscanf(line[idx+6:], "%d", &report.Value)
	}

	// Extract Expected
	if idx := strings.Index(line, "Expected="); idx != -1 {
		fmt.Sscanf(line[idx+9:], "%d", &report.Expected)
	}

	// Extract Match
	if strings.Contains(line, "Match=true") {
		report.Match = true
	}

	// Extract Sigs (signature count)
	if idx := strings.Index(line, "Sigs="); idx != -1 {
		fmt.Sscanf(line[idx+5:], "%d", &report.Sigs)
	}

	return report
}

// TestLLOStreamsE2E_QuickVerify is a quick verification test
// Run with: make test-quick
func TestLLOStreamsE2E_QuickVerify(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Log("Quick verification: Checking if workflow decoded MAGIC_NUMBER...")

	// Capture evidence from all stages
	var evidence []LogEvidence

	mockEAEvidence := captureLogEvidence(ctx, "mock-ea", "424242", "Mock EA")
	if mockEAEvidence != nil {
		evidence = append(evidence, *mockEAEvidence)
	}

	// Capture streams evidence showing active subscribers (nSubscribers > 0)
	streamsEvidence := captureLogEvidenceWithSubscribers(ctx, "streams-node0", "Streams DON (nSubscribers >= 1)")
	if streamsEvidence != nil {
		evidence = append(evidence, *streamsEvidence)
	}

	// Capture workflow evidence for BOTH formats
	format5Evidence := captureLogEvidenceByFormat(ctx, "workflow-node0", 5, "Workflow (Format 5)")
	if format5Evidence != nil {
		evidence = append(evidence, *format5Evidence)
	}

	format7Evidence := captureLogEvidenceByFormat(ctx, "workflow-node0", 7, "Workflow (Format 7)")
	if format7Evidence != nil {
		evidence = append(evidence, *format7Evidence)
	}

	reports, err := waitForValueReportsQuiet(ctx, 1*time.Minute)
	require.NoError(t, err, "No LLO_E2E_VALUE logs found - is the environment running?")
	require.GreaterOrEqual(t, len(reports), 2, "Expected at least 2 LLO_E2E_VALUE reports")

	matchCountF5 := 0
	matchCountF7 := 0
	for _, report := range reports {
		if report.Format == 5 && report.Value == MAGIC_NUMBER_FORMAT5 {
			matchCountF5++
		} else if report.Format == 7 && report.Value == MAGIC_NUMBER_FORMAT7 {
			matchCountF7++
		} else if report.Format == 0 && report.Value == MAGIC_NUMBER {
			matchCountF5++ // backward compat
		}
	}

	// Only show first 3 and last 2 reports (for brevity)
	t.Log("")
	t.Log("  📊 Sample reports (showing first 3 and last 2):")
	showCount := 0
	for i, report := range reports {
		if i < 3 || i >= len(reports)-2 {
			t.Logf("    Report %d: SeqNr=%d, Format=%d, Value=%d, Match=%v, Sigs=%d",
				i+1, report.SeqNr, report.Format, report.Value, report.Match, report.Sigs)
			showCount++
		} else if i == 3 && len(reports) > 5 {
			t.Logf("    ... (%d more reports) ...", len(reports)-5)
		}
	}

	// ASSERT: Both formats must have at least one match
	require.Greater(t, matchCountF5, 0,
		"Expected at least 1 Format 5 report with MAGIC_NUMBER_FORMAT5=%d. Got 0 F5 matches.",
		MAGIC_NUMBER_FORMAT5)
	require.Greater(t, matchCountF7, 0,
		"Expected at least 1 Format 7 report with MAGIC_NUMBER_FORMAT7=%d. Got 0 F7 matches.",
		MAGIC_NUMBER_FORMAT7)

	t.Log("")
	t.Log("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║  ✅ QUICK VERIFY PASSED - Format5=%d matches ✓, Format7=%d matches ✓       ║", matchCountF5, matchCountF7)
	t.Log("╚══════════════════════════════════════════════════════════════════════════════╝")
	// Actually run the verification commands and show output
	t.Log("")
	t.Log("📋 PROOF - Exact log lines used to verify each stage:")
	t.Log("")
	for _, e := range evidence {
		cmd := fmt.Sprintf("docker logs %s 2>&1 | sed -n '%dp'", e.Container, e.LineNumber)
		t.Logf("  # %s", e.Stage)
		t.Logf("  $ %s", cmd)
		out, err := exec.CommandContext(ctx, "bash", "-c", cmd).CombinedOutput()
		if err != nil {
			t.Logf("    (error: %v)", err)
		} else {
			t.Logf("    → %s", strings.TrimSpace(string(out)))
		}
		t.Log("")
	}
}

// waitForValueReportsQuiet waits for LLO_E2E_VALUE logs without printing each one
func waitForValueReportsQuiet(ctx context.Context, timeout time.Duration) ([]ValueReport, error) {
	deadline := time.Now().Add(timeout)
	var reports []ValueReport
	seenSeqNrs := make(map[int]bool)

	for time.Now().Before(deadline) {
		for i := 0; i < numWorkflowNodes; i++ {
			containerName := fmt.Sprintf("workflow-node%d", i)
			cmd := exec.CommandContext(ctx, "docker", "logs", containerName)
			out, err := cmd.CombinedOutput()
			if err != nil {
				continue
			}

			logs := string(out)
			for _, line := range strings.Split(logs, "\n") {
				if strings.Contains(line, "LLO_E2E_VALUE[SeqNr=") {
					report := parseValueReport(line)
					if !seenSeqNrs[report.SeqNr] {
						seenSeqNrs[report.SeqNr] = true
						reports = append(reports, report)
					}
				}
			}
		}

		// Count matches per format - we need at least one of EACH format
		matchCountF5 := 0
		matchCountF7 := 0
		for _, r := range reports {
			if r.Format == 5 && r.Value == MAGIC_NUMBER_FORMAT5 {
				matchCountF5++
			} else if r.Format == 7 && r.Value == MAGIC_NUMBER_FORMAT7 {
				matchCountF7++
			} else if r.Format == 0 && r.Value == MAGIC_NUMBER {
				matchCountF5++ // backward compat
			}
		}
		// Exit when we have at least one match from EACH format
		if matchCountF5 >= 1 && matchCountF7 >= 1 {
			return reports, nil
		}

		time.Sleep(5 * time.Second)
	}

	if len(reports) > 0 {
		return reports, nil
	}

	return nil, fmt.Errorf("no LLO_E2E_VALUE logs found within %v", timeout)
}

// captureLogEvidence finds a log line containing the search term and returns evidence
func captureLogEvidence(ctx context.Context, container, searchTerm, stage string) *LogEvidence {
	cmd := exec.CommandContext(ctx, "docker", "logs", container)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if strings.Contains(line, searchTerm) {
			return &LogEvidence{
				Stage:      stage,
				Container:  container,
				LineNumber: i + 1, // 1-indexed
				LogSnippet: line,
				KeyValue:   searchTerm,
			}
		}
	}
	return nil
}

// captureLogEvidenceByFormat finds LLO_E2E_VALUE logs for a specific report format
func captureLogEvidenceByFormat(ctx context.Context, container string, format int, stage string) *LogEvidence {
	cmd := exec.CommandContext(ctx, "docker", "logs", container)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	// Look for LLO_E2E_VALUE lines with the specific format
	searchPattern := fmt.Sprintf("Format=%d Value=", format)
	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if strings.Contains(line, "LLO_E2E_VALUE") && strings.Contains(line, searchPattern) {
			return &LogEvidence{
				Stage:      stage,
				Container:  container,
				LineNumber: i + 1, // 1-indexed
				LogSnippet: line,
				KeyValue:   searchPattern,
			}
		}
	}
	return nil
}

// captureLogEvidenceWithSubscribers finds CRE Transmitter logs showing active subscribers
func captureLogEvidenceWithSubscribers(ctx context.Context, container, stage string) *LogEvidence {
	cmd := exec.CommandContext(ctx, "docker", "logs", container)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	// Look for ProcessReport logs where nSubscribers >= 1
	// Search backwards to find the most recent one with subscribers
	lines := strings.Split(string(out), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.Contains(line, "ProcessReport pushing event") && strings.Contains(line, `"nSubscribers":`) {
			// Check if nSubscribers is > 0
			// Pattern: "nSubscribers":N where N > 0
			if !strings.Contains(line, `"nSubscribers":0`) {
				return &LogEvidence{
					Stage:      stage,
					Container:  container,
					LineNumber: i + 1, // 1-indexed
					LogSnippet: line,
					KeyValue:   "nSubscribers >= 1",
				}
			}
		}
	}
	return nil
}

// truncateLog truncates a log line to maxLen characters
func truncateLog(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	// Pad to maxLen for alignment
	return s + strings.Repeat(" ", maxLen-len(s))
}

// truncateForCmd truncates a string for use in a shell command
func truncateForCmd(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}
