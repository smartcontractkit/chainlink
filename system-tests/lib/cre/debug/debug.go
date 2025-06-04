package debug

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func PrintTestDebug(ctx context.Context, testName string, l zerolog.Logger, input types.DebugInput) {
	l.Info().Msg("🔍 Debug information from Chainlink Node logs:")

	if err := input.Validate(); err != nil {
		l.Error().Err(err).Msg("Input validation failed. No debug information will be printed")
		return
	}

	if input.InfraInput.InfraType == libtypes.CRIB {
		l.Error().Msg("❌ Debug information is not supported for CRIB")
		return
	}

	var allLogFiles []*os.File

	defer func() {
		for _, f := range allLogFiles {
			if f != nil {
				_ = (*f).Close()
			}
		}
	}()

	for _, debugDon := range input.DebugDons {
		logFiles, err := getLogFileHandles(testName, l, debugDon)
		if err != nil {
			l.Error().Err(err).Msg("Failed to get log file handles. No debug information will be printed")
			return
		}

		allLogFiles = append(allLogFiles, logFiles...)
		workflowNodeSet, err := node.FindManyWithLabel(debugDon.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.WorkerNode}, node.EqualLabels)
		if err != nil {
			// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
			l.Error().Err(err).Msg("Failed to find worker nodes. No debug information will be printed")
			return
		}

		workflowNodeCount := len(workflowNodeSet)

		if flags.HasFlag(debugDon.Flags, types.WorkflowDON) {
			if detectedCount := checkHowManyNodesLogsHaveText(logFiles, workflowNodeCount, "handled event.*WorkflowRegisteredV1"); detectedCount != workflowNodeCount {
				l.Error().Msgf("❌ Workflow registration was not detected by all nodes (%d/%d)", detectedCount, workflowNodeCount)
				return
			}
			l.Info().Msg("✅ Workflow registration was detected")
		}

		if flags.HasFlag(debugDon.Flags, types.WorkflowDON) {
			if detectedCount := checkHowManyNodesLogsHaveText(logFiles, workflowNodeCount, "step request enqueued"); detectedCount != workflowNodeCount {
				l.Error().Msgf("❌ Workflow was not executing on all nodes (%d/%d)", detectedCount, workflowNodeCount)
				return
			}
			l.Info().Msg("✅ Workflow was executing")
		}

		if flags.HasFlag(debugDon.Flags, types.OCR3Capability) {
			if detectedCount := checkHowManyNodesLogsHaveText(logFiles, workflowNodeCount, "✅ committed outcome"); detectedCount != workflowNodeCount {
				l.Error().Msgf("❌ OCR was not executing on all nodes (%d/%d)", detectedCount, workflowNodeCount)
				return
			}
			l.Info().Msg("✅ OCR was executing")
		}

		if flags.HasFlag(debugDon.Flags, types.WriteEVMCapability) {
			if !checkIfAtLeastOneReportWasSent(logFiles, workflowNodeCount) {
				l.Error().Msg("❌ Reports were not sent")
				return
			}

			l.Info().Msg("✅ Reports were sent")

			// debug report transmissions
			ReportTransmissions(ctx, logFiles, l, input.BlockchainOutput.Nodes[0].ExternalWSUrl)
		}

		// Add support for new capabilities here as needed, if there is some specific debug information to be printed
	}
}

// This function is used to go through Chainlink Node logs and look for entries related to report transmissions.
// Once such a log entry is found, it looks for transaction hash and then it tries to decode the transaction and print the result.
func ReportTransmissions(ctx context.Context, logFiles []*os.File, l zerolog.Logger, wsRPCURL string) {
	/*
	 Example log entry:
	 2025-01-28T14:44:48.080Z [DEBUG] Node sent transaction                              multinode@v0.0.0-20250121205514-f73e2f86c23b/transaction_sender.go:180 chainID=1337 logger=EVM.1337.TransactionSender tx={"type":"0x0","chainId":"0x539","nonce":"0x0","to":"0xcf7ed3acca5a467e9e704c703e8d87f634fb0fc9","gas":"0x61a80","gasPrice":"0x3b9aca00","maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x0","input":"0x11289565000000000000000000000000a513e6e4b8f2a923d98304ec87f64353c4d5c853000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000010d010f715db03509d388f706e16137722000e26aa650a64ac826ae8e5679cdf57fd96798ed50000000010000000100000a9c593aaed2f5371a5bc0779d1b8ea6f9c7d37bfcbb876a0a9444dbd36f64306466323239353031f39fd6e51aad88f6f4ce6ab8827279cfffb92266000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001018bfe88407000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000bb5c162c8000000000000000000000000000000000000000000000000000000006798ed37000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000e700d4c57250eac9dc925c951154c90c1b6017944322fb2075055d8bdbe19000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000041561c171b7465e8efef35572ef82adedb49ea71b8344a34a54ce5e853f80ca1ad7d644ebe710728f21ebfc3e2407bd90173244f744faa011c3a57213c8c585de90000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004165e6f3623acc43f163a58761655841bfebf3f6b4ea5f8d34c64188036b0ac23037ebbd3854b204ca26d828675395c4b9079ca068d9798326eb8c93f26570a1080100000000000000000000000000000000000000000000000000000000000000","v":"0xa96","r":"0x168547e96e7088c212f85a4e8dddce044bbb2abfd5ccc8a5451fdfcb812c94e5","s":"0x2a735a3df046632c2aaa7e583fe161113f3345002e6c9137bbfa6800a63f28a4","hash":"0x3fc5508310f8deef09a46ad594dcc5dc9ba415319ef1dfa3136335eb9e87ff4d"} version=2.19.0@05c05a9

	 What we are looking for:
	 "hash":"0x3fc5508310f8deef09a46ad594dcc5dc9ba415319ef1dfa3136335eb9e87ff4d"
	*/
	reportTransmissionTxHashPattern := regexp.MustCompile(`"hash":"(0x[0-9a-fA-F]+)"`)

	// let's be prudent and assume that in extreme scenario when feed price isn't updated, but
	// transmission is still sent, we might have multiple transmissions per node, and if we want
	// to avoid blocking on the channel, we need to have a higher buffer
	resultsCh := make(chan string, len(logFiles)*4)

	wg := &sync.WaitGroup{}
	for _, f := range logFiles {
		wg.Add(1)
		file := f

		go func() {
			defer wg.Done()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				jsonLogLine := scanner.Text()

				if !strings.Contains(jsonLogLine, "Node sent transaction") {
					continue
				}

				match := reportTransmissionTxHashPattern.MatchString(jsonLogLine)
				if match {
					resultsCh <- reportTransmissionTxHashPattern.FindStringSubmatch(jsonLogLine)[1]
				}
			}
		}()
	}

	wg.Wait()
	close(resultsCh)

	if len(resultsCh) == 0 {
		l.Error().Msg("❌ No report transmissions found in Chainlink Node logs.")
		return
	}

	// required as Seth prints transaction traces to stdout with debug level
	_ = os.Setenv(seth.LogLevelEnvVar, "debug")

	sc, err := seth.NewClientBuilder().
		WithRpcUrl(wsRPCURL).
		WithReadOnlyMode().
		// we assume that chainlink-evm is in the same directory as chainlink, so we can use relative path
		WithGethWrappersFolders([]string{"../../../../../chainlink-evm/gethwrappers/keystone/generated"}). // point Seth to the folder with keystone geth wrappers, so that it can load contract ABIs
		Build()

	if err != nil {
		l.Error().Err(err).Msg("Failed to create seth client")
		return
	}

	for txHash := range resultsCh {
		l.Info().Msgf("🔍 Tracing report transmission transaction %s", txHash)
		// set tracing level to all to trace also successful transactions
		sc.Cfg.TracingLevel = seth.TracingLevel_All
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
		tx, _, txErr := sc.Client.TransactionByHash(ctxWithTimeout, common.HexToHash(txHash))
		if txErr != nil {
			cancel()
			l.Warn().Err(txErr).Msgf("Failed to get transaction by hash %s", txHash)
			continue
		}
		cancel()
		_, decodedErr := sc.DecodeTx(tx)

		if decodedErr != nil {
			l.Error().Err(decodedErr).Msgf("Transmission transaction %s failed due to %s", txHash, decodedErr.Error())
			continue
		}
	}
}

func getLogFileHandles(testName string, l zerolog.Logger, debugDon *types.DebugDon) ([]*os.File, error) {
	var logFiles []*os.File

	var belongsToCurrentEnv = func(filePath string) bool {
		for i, containerName := range debugDon.ContainerNames {
			// TODO check corresponding flag when looking for bootstrap node
			// skip the first node, as it's the bootstrap node
			if i == 0 {
				continue
			}

			if strings.EqualFold(filePath, containerName+".log") {
				return true
			}
		}
		return false
	}

	logsDir := "logs/docker-" + testName

	fileWalkErr := filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && belongsToCurrentEnv(info.Name()) {
			file, fileErr := os.Open(path)
			if fileErr != nil {
				return fmt.Errorf("failed to open file %s: %w", path, fileErr)
			}
			logFiles = append(logFiles, file)
		}
		return nil
	})

	expectedLogCount := len(debugDon.ContainerNames) - 1
	if len(logFiles) != expectedLogCount {
		l.Warn().Int("Expected", expectedLogCount).Int("Got", len(logFiles)).Msg("Number of log files does not match number of worker nodes. Some logs might be missing.")
	}

	if fileWalkErr != nil {
		l.Error().Err(fileWalkErr).Msg("Error walking through log files. Will not look for report transmission transaction hashes")
		return nil, fileWalkErr
	}

	return logFiles, nil
}

func checkIfAtLeastOneReportWasSent(logFiles []*os.File, workflowNodeCount int) bool {
	// we are looking for "Node sent transaction" log entry, which might appear various times in the logs
	// but most probably not in the logs of all nodes, since they take turns in sending reports
	// our buffer must be large enough to capture all the possible log entries in order to avoid channel blocking
	bufferSize := workflowNodeCount * 4

	return checkHowManyNodesLogsHaveText(logFiles, bufferSize, "Node sent transaction") > 0
}

func checkHowManyNodesLogsHaveText(logFiles []*os.File, bufferSize int, expectedText string) int {
	wg := &sync.WaitGroup{}

	resultsCh := make(chan struct{}, bufferSize)

	for _, f := range logFiles {
		wg.Add(1)
		file := f

		go func() {
			defer func() {
				wg.Done()
				// reset file pointer to the beginning of the file
				// so that subsequent reads start from the beginning
				_, _ = file.Seek(0, io.SeekStart)
			}()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				jsonLogLine := scanner.Text()

				matched, err := regexp.MatchString(expectedText, jsonLogLine)
				if err != nil {
					return
				}

				if matched {
					resultsCh <- struct{}{}
					return
				}
			}
		}()
	}

	wg.Wait()
	close(resultsCh)

	var found int
	for range resultsCh {
		found++
	}

	return found
}
