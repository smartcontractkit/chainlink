package logpoller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	geth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/scylladb/go-reflectx"
	"github.com/stretchr/testify/require"

	ac "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_compatible_utils"
	le "github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/log_emitter"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	cltypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	common_logger "github.com/smartcontractkit/chainlink-common/pkg/logger"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
	"github.com/smartcontractkit/chainlink/devenv/products/automation/concurrency"
)

var (
	EmitterABI, _     = abi.JSON(strings.NewReader(le.LogEmitterABI))
	automatoinConvABI = cltypes.MustGetABI(ac.AutomationCompatibleUtilsABI)
	bytes0            = [32]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	} // bytes representation of 0x0000000000000000000000000000000000000000000000000000000000000000

)

var registerSingleTopicFilter = func(registry contracts.KeeperRegistry, upkeepID *big.Int, emitterAddress common.Address, topic common.Hash) error {
	logTriggerConfigStruct := ac.IAutomationV21PlusCommonLogTriggerConfig{
		ContractAddress: emitterAddress,
		FilterSelector:  0,
		Topic0:          topic,
		Topic1:          bytes0,
		Topic2:          bytes0,
		Topic3:          bytes0,
	}
	encodedLogTriggerConfig, err := automatoinConvABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
	if err != nil {
		return err
	}

	err = registry.SetUpkeepTriggerConfig(upkeepID, encodedLogTriggerConfig)
	if err != nil {
		return err
	}

	return nil
}

// Currently Unused November 8, 2023, Might be useful in the near future so keeping it here for now
// this is not really possible, log trigger doesn't support multiple topics, even if log poller does
// var registerMultipleTopicsFilter = func(registry contracts.KeeperRegistry, upkeepID *big.Int, emitterAddress common.Address, topics []abi.Event) error {
// 	if len(topics) > 4 {
// 		return errors.New("Cannot register more than 4 topics")
// 	}

// 	var getTopic = func(topics []abi.Event, i int) common.Hash {
// 		if i > len(topics)-1 {
// 			return bytes0
// 		}

// 		return topics[i].ID
// 	}

// 	var getFilterSelector = func(topics []abi.Event) (uint8, error) {
// 		switch len(topics) {
// 		case 0:
// 			return 0, errors.New("Cannot register filter with 0 topics")
// 		case 1:
// 			return 0, nil
// 		case 2:
// 			return 1, nil
// 		case 3:
// 			return 3, nil
// 		case 4:
// 			return 7, nil
// 		default:
// 			return 0, errors.New("Cannot register filter with more than 4 topics")
// 		}
// 	}

// 	filterSelector, err := getFilterSelector(topics)
// 	if err != nil {
// 		return err
// 	}

// 	logTriggerConfigStruct := automation_convenience.LogTriggerConfig{
// 		ContractAddress: emitterAddress,
// 		FilterSelector:  filterSelector,
// 		Topic0:          getTopic(topics, 0),
// 		Topic1:          getTopic(topics, 1),
// 		Topic2:          getTopic(topics, 2),
// 		Topic3:          getTopic(topics, 3),
// 	}
// 	encodedLogTriggerConfig, err := automatoinConvABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
// 	if err != nil {
// 		return err
// 	}

// 	err = registry.SetUpkeepTriggerConfig(upkeepID, encodedLogTriggerConfig)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func newORM(logger common_logger.Logger, chainID *big.Int, nodeIndex, externalPort int) (logpoller.ORM, *sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", externalPort, postgres.User, postgres.Password, fmt.Sprintf("db_%d", nodeIndex))
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, db, err
	}

	db.MapperFunc(reflectx.CamelToSnakeASCII)
	return logpoller.NewORM(chainID, db, logger), db, nil
}

type ExpectedFilter struct {
	emitterAddress common.Address
	topic          common.Hash
}

// getExpectedFilters returns a slice of ExpectedFilter structs based on the provided log emitters and config
func getExpectedFilters(logEmitters []contracts.LogEmitter, cfg *Config) []ExpectedFilter {
	expectedFilters := make([]ExpectedFilter, 0)
	for _, emitter := range logEmitters {
		for _, event := range cfg.General.EventsToEmit {
			expectedFilters = append(expectedFilters, ExpectedFilter{
				emitterAddress: emitter.Address(),
				topic:          event.ID,
			})
		}
	}

	return expectedFilters
}

// nodeHasExpectedFilters returns true if the provided node has all the expected filters registered
func nodeHasExpectedFilters(ctx context.Context, expectedFilters []ExpectedFilter, logger common_logger.SugaredLogger, chainID *big.Int, nodeIndex, dbExternalPort int) (bool, string, error) {
	orm, db, err := newORM(logger, chainID, nodeIndex, dbExternalPort)
	if err != nil {
		return false, "", err
	}

	defer db.Close()
	knownFilters, err := orm.LoadFilters(ctx)
	if err != nil {
		return false, "", err
	}

	for _, expectedFilter := range expectedFilters {
		filterFound := false
		for _, knownFilter := range knownFilters {
			if bytes.Equal(expectedFilter.emitterAddress.Bytes(), knownFilter.Addresses[0].Bytes()) && bytes.Equal(expectedFilter.topic.Bytes(), knownFilter.EventSigs[0].Bytes()) {
				filterFound = true
				break
			}
		}

		if !filterFound {
			return false, fmt.Sprintf("no filter found for emitter %s and topic %s", expectedFilter.emitterAddress.String(), expectedFilter.topic.Hex()), nil
		}
	}

	return true, "", nil
}

// randomWait waits for a random amount of time between minMilliseconds and maxMilliseconds
func randomWait(minMilliseconds, maxMilliseconds int) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomMilliseconds := rand.Intn(maxMilliseconds-minMilliseconds+1) + minMilliseconds
	time.Sleep(time.Duration(randomMilliseconds) * time.Millisecond)
}

// getIntSlice returns a slice of ints of the provided length
func getIntSlice(length int) []int {
	result := make([]int, length)
	for i := range length {
		result[i] = i
	}

	return result
}

// getStringSlice returns a slice of strings of the provided length
func getStringSlice(length int) []string {
	result := make([]string, length)
	for i := range length {
		result[i] = "amazing event"
	}

	return result
}

// logPollerHasFinalisedEndBlock returns true if all CL nodes have finalised processing the provided end block
func logPollerHasFinalisedEndBlock(endBlock int64, chainID *big.Int, l zerolog.Logger, coreLogger common_logger.SugaredLogger, testEnv *logPollerEnvironment) (bool, error) {
	wg := &sync.WaitGroup{}

	type boolQueryResult struct {
		nodeName       string
		hasFinalised   bool
		finalizedBlock int64
		err            error
	}

	endBlockCh := make(chan boolQueryResult, len(testEnv.nodes.NodeSpecs)-1)
	ctx, cancelFn := context.WithCancel(context.Background())

	for i := 1; i < len(testEnv.nodes.NodeSpecs); i++ {
		wg.Add(1)

		go func(clNode *clnode.Output, idx int, r chan boolQueryResult) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				orm, db, err := newORM(coreLogger, chainID, idx, testEnv.dbPort())
				if err != nil {
					r <- boolQueryResult{
						nodeName:     clNode.Node.ContainerName,
						hasFinalised: false,
						err:          err,
					}
				}

				defer db.Close()

				latestBlock, err := orm.SelectLatestBlock(ctx)
				if err != nil {
					r <- boolQueryResult{
						nodeName:     clNode.Node.ContainerName,
						hasFinalised: false,
						err:          err,
					}
				}

				r <- boolQueryResult{
					nodeName:       clNode.Node.ContainerName,
					finalizedBlock: latestBlock.FinalizedBlockNumber,
					hasFinalised:   latestBlock.FinalizedBlockNumber > endBlock,
					err:            nil,
				}
			}
		}(testEnv.nodes.Out.CLNodes[i], i, endBlockCh)
	}

	var err error
	allFinalisedCh := make(chan bool, 1)

	go func() {
		foundMap := make(map[string]bool, 0)
		for r := range endBlockCh {
			if r.err != nil {
				err = r.err
				cancelFn()
				return
			}

			foundMap[r.nodeName] = r.hasFinalised
			if r.hasFinalised {
				l.Info().Str("Node name", r.nodeName).Msg("CL node has finalised end block")
			} else {
				l.Warn().Int64("Has", r.finalizedBlock).Int64("Want", endBlock).Str("Node name", r.nodeName).Msg("CL node has not finalised end block yet")
			}

			if len(foundMap) == len(testEnv.nodes.NodeSpecs)-1 {
				allFinalised := true
				for _, v := range foundMap {
					if !v {
						allFinalised = false
						break
					}
				}

				allFinalisedCh <- allFinalised
				return
			}
		}
	}()

	wg.Wait()
	close(endBlockCh)

	return <-allFinalisedCh, err
}

// nodesHaveExpectedLogCount returns true if all CL nodes have the expected log count in the provided block range and matching the provided filters
func nodesHaveExpectedLogCount(startBlock, endBlock int64, chainID *big.Int, expectedLogCount int, expectedFilters []ExpectedFilter, l zerolog.Logger, coreLogger common_logger.SugaredLogger, testEnv *logPollerEnvironment) (bool, error) {
	wg := &sync.WaitGroup{}

	type logQueryResult struct {
		nodeName         string
		logCount         int
		hasExpectedCount bool
		err              error
	}

	resultChan := make(chan logQueryResult, len(testEnv.nodes.NodeSpecs)-1)
	ctx, cancelFn := context.WithCancel(context.Background())

	for i := 1; i < len(testEnv.nodes.NodeSpecs); i++ {
		wg.Add(1)

		go func(clNode *clnode.Output, idx int, resultChan chan logQueryResult) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				orm, db, err := newORM(coreLogger, chainID, idx, testEnv.dbPort())
				if err != nil {
					resultChan <- logQueryResult{
						nodeName:         clNode.Node.ContainerName,
						logCount:         0,
						hasExpectedCount: false,
						err:              err,
					}
				}

				defer db.Close()
				foundLogsCount := 0

				for _, filter := range expectedFilters {
					logs, err := orm.SelectLogs(ctx, startBlock, endBlock, filter.emitterAddress, filter.topic)
					if err != nil {
						resultChan <- logQueryResult{
							nodeName:         clNode.Node.ContainerName,
							logCount:         0,
							hasExpectedCount: false,
							err:              err,
						}
					}

					foundLogsCount += len(logs)
				}

				resultChan <- logQueryResult{
					nodeName:         clNode.Node.ContainerName,
					logCount:         foundLogsCount,
					hasExpectedCount: foundLogsCount >= expectedLogCount,
					err:              nil,
				}
			}
		}(testEnv.nodes.Out.CLNodes[i], i, resultChan)
	}

	var err error
	allFoundCh := make(chan bool, 1)

	go func() {
		foundMap := make(map[string]bool, 0)
		for r := range resultChan {
			if r.err != nil {
				err = r.err
				cancelFn()
				return
			}

			foundMap[r.nodeName] = r.hasExpectedCount
			if r.hasExpectedCount {
				l.Debug().
					Str("Node name", r.nodeName).
					Int("Logs count", r.logCount).
					Msg("Expected log count found in CL node")
			} else {
				l.Debug().
					Str("Node name", r.nodeName).
					Str("Found/Expected logs", fmt.Sprintf("%d/%d", r.logCount, expectedLogCount)).
					Int("Missing logs", expectedLogCount-r.logCount).
					Msg("Too low log count found in CL node")
			}

			if len(foundMap) == len(testEnv.nodes.NodeSpecs)-1 {
				allFound := true
				for _, hadAllLogs := range foundMap {
					if !hadAllLogs {
						allFound = false
						break
					}
				}

				allFoundCh <- allFound
				return
			}
		}
	}()

	wg.Wait()
	close(resultChan)

	return <-allFoundCh, err
}

type MissingLogs map[string][]geth_types.Log

// IsEmpty returns true if there are no missing logs
func (m *MissingLogs) IsEmpty() bool {
	for _, v := range *m {
		if len(v) > 0 {
			return false
		}
	}

	return true
}

// missingLogs returns a map of CL node name to missing logs in that node compared to EVM node to which the provided evm client is connected
func missingLogs(
	startBlock, endBlock int64,
	testEnv *logPollerEnvironment,
	l zerolog.Logger,
	coreLogger common_logger.SugaredLogger,
	cfg *Config,
) (MissingLogs, error) {
	wg := &sync.WaitGroup{}

	type dbQueryResult struct {
		err      error
		nodeName string
		logs     []logpoller.Log
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	resultCh := make(chan dbQueryResult, len(testEnv.nodes.NodeSpecs)-1)

	for i := 1; i < len(testEnv.nodes.NodeSpecs); i++ {
		wg.Add(1)

		go func(ctx context.Context, i int, r chan dbQueryResult) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				l.Warn().Msg("Context cancelled. Terminating fetching logs from log poller's DB")
				return
			default:
				nodeName := testEnv.nodes.Out.CLNodes[i].Node.ContainerName

				l.Debug().Str("Node name", nodeName).Msg("Fetching log poller logs")
				orm, db, err := newORM(coreLogger, big.NewInt(testEnv.chainClient.ChainID), i, testEnv.dbPort())
				if err != nil {
					r <- dbQueryResult{
						err:      err,
						nodeName: nodeName,
						logs:     []logpoller.Log{},
					}
				}

				defer db.Close()
				logs := make([]logpoller.Log, 0)

				for j := range testEnv.logEmitters {
					address := testEnv.logEmitters[j].Address()

					for _, event := range cfg.General.EventsToEmit {
						l.Trace().Str("Event name", event.Name).Str("Emitter address", address.String()).Msg("Fetching single emitter's logs")
						result, err := orm.SelectLogs(ctx, startBlock, endBlock, address, event.ID)
						if err != nil {
							r <- dbQueryResult{
								err:      err,
								nodeName: nodeName,
								logs:     []logpoller.Log{},
							}
						}

						sort.Slice(result, func(i, j int) bool {
							return result[i].BlockNumber < result[j].BlockNumber
						})

						logs = append(logs, result...)

						l.Trace().Str("Event name", event.Name).Str("Emitter address", address.String()).Int("Log count", len(result)).Msg("Logs found per node")
					}
				}

				l.Info().Int("Count", len(logs)).Str("Node name", nodeName).Msg("Fetched log poller logs")

				r <- dbQueryResult{
					err:      nil,
					nodeName: nodeName,
					logs:     logs,
				}
			}
		}(ctx, i, resultCh)
	}

	allLogPollerLogs := make(map[string][]logpoller.Log, 0)
	missingLogs := map[string][]geth_types.Log{}
	var dbError error

	go func() {
		for r := range resultCh {
			if r.err != nil {
				l.Err(r.err).Str("Node name", r.nodeName).Msg("Error fetching logs from log poller's DB")
				dbError = r.err
				cancelFn()
				return
			}
			// use channel for aggregation and then for := range over it after closing resultCh?
			allLogPollerLogs[r.nodeName] = r.logs
		}
	}()

	wg.Wait()
	close(resultCh)

	if dbError != nil {
		return nil, dbError
	}

	allLogsInEVMNode, err := getEVMLogs(ctx, startBlock, endBlock, testEnv.logEmitters, testEnv.chainClient, l, cfg)
	if err != nil {
		return nil, err
	}

	wg = &sync.WaitGroup{}

	type missingLogResult struct {
		nodeName string
		logs     []geth_types.Log
	}

	evmLogCount := len(allLogsInEVMNode)
	l.Info().Int("Log count", evmLogCount).Msg("Started comparison of logs from EVM node and CL nodes. This may take a while if there's a lot of logs")

	missingCh := make(chan missingLogResult, len(testEnv.nodes.NodeSpecs)-1)
	for i := 1; i < len(testEnv.nodes.NodeSpecs); i++ {
		wg.Add(1)

		go func(i int, result chan missingLogResult) {
			defer wg.Done()
			nodeName := testEnv.nodes.Out.CLNodes[i].Node.ContainerName
			l.Debug().Str("Node name", nodeName).Str("Progress", fmt.Sprintf("0/%d", evmLogCount)).Msg("Comparing single CL node's logs with EVM logs")

			missingLogs := make([]geth_types.Log, 0)
			for i, evmLog := range allLogsInEVMNode {
				logFound := false
				if evmLog.BlockNumber > math.MaxInt64 {
					panic(fmt.Errorf("block number overflows int64: %d", evmLog.BlockNumber))
				}
				if evmLog.Index > math.MaxInt64 {
					panic(fmt.Errorf("index overflows int64: %d", evmLog.Index))
				}
				for _, logPollerLog := range allLogPollerLogs[nodeName] {
					if logPollerLog.BlockNumber == int64(evmLog.BlockNumber) && logPollerLog.TxHash == evmLog.TxHash && bytes.Equal(logPollerLog.Data, evmLog.Data) && logPollerLog.LogIndex == int64(evmLog.Index) &&
						logPollerLog.Address == evmLog.Address && logPollerLog.BlockHash == evmLog.BlockHash && bytes.Equal(logPollerLog.Topics[0], evmLog.Topics[0].Bytes()) {
						logFound = true
						continue
					}
				}

				if i%10000 == 0 && i != 0 {
					l.Debug().Str("Node name", nodeName).Str("Progress", fmt.Sprintf("%d/%d", i, evmLogCount)).Msg("Comparing single CL node's logs with EVM logs")
				}

				if !logFound {
					missingLogs = append(missingLogs, evmLog)
				}
			}

			if len(missingLogs) > 0 {
				l.Warn().Int("Count", len(missingLogs)).Str("Node name", nodeName).Msg("Some EMV logs were missing from CL node")
			} else {
				l.Info().Str("Node name", nodeName).Str("Missing/Total logs", fmt.Sprintf("%d/%d", len(missingLogs), evmLogCount)).Msg("All EVM logs were found in CL node")
			}

			result <- missingLogResult{
				nodeName: nodeName,
				logs:     missingLogs,
			}
		}(i, missingCh)
	}

	wg.Wait()
	close(missingCh)

	for v := range missingCh {
		if len(v.logs) > 0 {
			missingLogs[v.nodeName] = v.logs
		}
	}

	expectedTotalLogsEmitted := getExpectedLogCount(cfg)
	if int64(len(allLogsInEVMNode)) != expectedTotalLogsEmitted {
		l.Warn().
			Str("Actual/Expected", fmt.Sprintf("%d/%d", expectedTotalLogsEmitted, len(allLogsInEVMNode))).
			Msg("Actual number of logs found on EVM nodes differs from expected ones. Most probably this is a bug in the test")
	}

	return missingLogs, nil
}

// printMissingLogsInfo prints various useful information about the missing logs
func printMissingLogsInfo(missingLogs map[string][]geth_types.Log, l zerolog.Logger, cfg *Config) {
	findHumanName := func(topic common.Hash) string {
		for _, event := range cfg.General.EventsToEmit {
			if event.ID == topic {
				return event.Name
			}
		}

		return "Unknown event"
	}

	missingByType := make(map[string]int)
	for _, logs := range missingLogs {
		for _, v := range logs {
			humanName := findHumanName(v.Topics[0])
			missingByType[humanName]++
		}
	}

	l.Debug().Msg("Missing log by event name")
	for k, v := range missingByType {
		l.Debug().Str("Event name", k).Int("Missing count", v).Msg("Missing logs by type")
	}

	missingByBlock := make(map[uint64]int)
	for _, logs := range missingLogs {
		for _, l := range logs {
			missingByBlock[l.BlockNumber]++
		}
	}

	l.Debug().Msg("Missing logs by block")
	for k, v := range missingByBlock {
		l.Debug().Uint64("Block number", k).Int("Missing count", v).Msg("Missing logs by block")
	}

	missingByEmitter := make(map[string]int)
	for _, logs := range missingLogs {
		for _, l := range logs {
			missingByEmitter[l.Address.String()]++
		}
	}

	l.Debug().Msg("Missing logs by emitter")
	for k, v := range missingByEmitter {
		l.Debug().Str("Emitter address", k).Int("Missing count", v).Msg("Missing logs by emitter")
	}
}

// getEVMLogs returns a slice of all logs emitted by the provided log emitters in the provided block range,
// which are present in the EVM node to which the provided evm client is connected
func getEVMLogs(ctx context.Context, startBlock, endBlock int64, logEmitters []contracts.LogEmitter, client *seth.Client, l zerolog.Logger, cfg *Config) ([]geth_types.Log, error) {
	allLogsInEVMNode := make([]geth_types.Log, 0)
	for j := range logEmitters {
		address := logEmitters[j].Address()
		for _, event := range cfg.General.EventsToEmit {
			l.Debug().Str("Event name", event.Name).Str("Emitter address", address.String()).Msg("Fetching logs from EVM node")
			logsInEVMNode, err := client.Client.FilterLogs(ctx, geth.FilterQuery{
				Addresses: []common.Address{(address)},
				Topics:    [][]common.Hash{{event.ID}},
				FromBlock: big.NewInt(startBlock),
				ToBlock:   big.NewInt(endBlock),
			})
			if err != nil {
				return nil, err
			}

			sort.Slice(logsInEVMNode, func(i, j int) bool {
				return logsInEVMNode[i].BlockNumber < logsInEVMNode[j].BlockNumber
			})

			allLogsInEVMNode = append(allLogsInEVMNode, logsInEVMNode...)
			l.Debug().Str("Event name", event.Name).Str("Emitter address", address.String()).Int("Log count", len(logsInEVMNode)).Msg("Logs found in EVM node")
		}
	}

	l.Info().Int("Count", len(allLogsInEVMNode)).Msg("Logs in EVM node")

	return allLogsInEVMNode, nil
}

// executeGenerator executes the configured generator and returns the total number of logs emitted
func executeGenerator(t *testing.T, cfg *Config, client *seth.Client, logEmitters []contracts.LogEmitter) (int, error) {
	if cfg.General.Generator == GeneratorType_WASP {
		return runWaspGenerator(t, cfg, logEmitters)
	}

	return runLoopedGenerator(cfg, client, logEmitters)
}

// runWaspGenerator runs the wasp generator and returns the total number of logs emitted
func runWaspGenerator(t *testing.T, cfg *Config, logEmitters []contracts.LogEmitter) (int, error) {
	l := framework.L

	var RPSprime int64

	// if LPS is set, we need to calculate based on countract count and events per transaction
	if cfg.Wasp.LPS > 0 {
		RPSprime = cfg.Wasp.LPS / int64(cfg.General.Contracts) / int64(cfg.General.EventsPerTx) / int64(len(cfg.General.EventsToEmit))

		if RPSprime < 1 {
			return 0, errors.New("invalid load configuration, effective RPS would have been zero. Adjust LPS, contracts count, events per tx or events to emit")
		}
	}

	// if RPS is set simply split it between contracts
	if cfg.Wasp.RPS > 0 {
		RPSprime = cfg.Wasp.RPS / int64(cfg.General.Contracts)
	}

	counter := &Counter{
		mu:    sync.Mutex{},
		value: 0,
	}

	p := wasp.NewProfile()

	for _, logEmitter := range logEmitters {
		g, err := wasp.NewGenerator(&wasp.Config{
			T:                     t,
			LoadType:              wasp.RPS,
			GenName:               "log_poller_gen_" + logEmitter.Address().String(),
			RateLimitUnitDuration: cfg.Wasp.RateLimitUnitDuration,
			CallTimeout:           cfg.Wasp.CallTimeout,
			Schedule: wasp.Plain(
				RPSprime,
				cfg.Wasp.Duration,
			),
			Gun: NewLogEmitterGun(
				logEmitter,
				cfg.General.EventsToEmit,
				cfg.General.EventsPerTx,
				l,
			),
			SharedData: counter,
		})
		p.Add(g, err)
	}

	_, err := p.Run(true)
	if err != nil {
		return 0, err
	}

	return counter.value, nil
}

type logEmissionTask struct {
	emitter      contracts.LogEmitter
	eventsToEmit []abi.Event
	eventsPerTx  int
}

type emittedLogsData struct {
	count int
}

func (d emittedLogsData) GetResult() emittedLogsData {
	return d
}

// runLoopedGenerator runs the looped generator and returns the total number of logs emitted
func runLoopedGenerator(cfg *Config, client *seth.Client, logEmitters []contracts.LogEmitter) (int, error) {
	l := framework.L

	tasks := make([]logEmissionTask, 0)
	for i := 0; i < cfg.LoopedConfig.ExecutionCount; i++ {
		for _, logEmitter := range logEmitters {
			tasks = append(tasks, logEmissionTask{
				emitter:      logEmitter,
				eventsToEmit: cfg.General.EventsToEmit,
				eventsPerTx:  cfg.General.EventsPerTx,
			})
		}
	}

	l.Info().Int("Total tasks", len(tasks)).Msg("Starting to emit events")

	// we need to sync nodes manually, because we are not using ephemeral addresses
	if err := client.NonceManager.UpdateNonces(); err != nil {
		return 0, err
	}

	atomicCounter := atomic.Int32{}

	emitAllEventsFn := func(resultCh chan emittedLogsData, errorCh chan error, _ int, task logEmissionTask) {
		current := atomicCounter.Add(1)

		address := task.emitter.Address().String()

		for _, event := range cfg.General.EventsToEmit {
			l.Debug().Str("Emitter address", address).Str("Event type", event.Name).Str("index", fmt.Sprintf("%d/%d", current, cfg.LoopedConfig.ExecutionCount)).Msg("Emitting log from emitter")
			var err error
			switch event.Name {
			case "Log1":
				_, err = client.Decode(task.emitter.EmitLogIntsFromKey(getIntSlice(cfg.General.EventsPerTx), client.AnySyncedKey()))
			case "Log2":
				_, err = client.Decode(task.emitter.EmitLogIntsIndexedFromKey(getIntSlice(cfg.General.EventsPerTx), client.AnySyncedKey()))
			case "Log3":
				_, err = client.Decode(task.emitter.EmitLogStringsFromKey(getStringSlice(cfg.General.EventsPerTx), client.AnySyncedKey()))
			case "Log4":
				_, err = client.Decode(task.emitter.EmitLogIntMultiIndexedFromKey(1, 1, cfg.General.EventsPerTx, client.AnySyncedKey()))
			default:
				err = fmt.Errorf("unknown event name: %s", event.Name)
			}

			if err != nil {
				errorCh <- err
				return
			}
			randomWait(cfg.LoopedConfig.MinEmitWaitTimeMs, cfg.LoopedConfig.MaxEmitWaitTimeMs)

			if (current)%10 == 0 {
				l.Info().Str("Emitter address", address).Str("Index", fmt.Sprintf("%d/%d", current, cfg.LoopedConfig.ExecutionCount)).Msgf("Emitted all %d events", len(cfg.General.EventsToEmit))
			}
		}

		resultCh <- emittedLogsData{
			cfg.General.EventsPerTx * len(cfg.General.EventsToEmit),
		}
	}

	executor := concurrency.NewConcurrentExecutor[emittedLogsData, emittedLogsData, logEmissionTask](l)
	r, err := executor.Execute(len(client.Cfg.Network.PrivateKeys)-1, tasks, emitAllEventsFn)
	if err != nil {
		return 0, err
	}

	var total int
	for _, result := range r {
		total += result.count
	}

	return total, nil
}

// getExpectedLogCount returns the expected number of logs to be emitted based on the provided config
func getExpectedLogCount(cfg *Config) int64 {
	if cfg.General.Generator == GeneratorType_WASP {
		if cfg.Wasp.RPS != 0 {
			return cfg.Wasp.RPS * int64(cfg.Wasp.Duration.Seconds()) * int64(cfg.General.EventsPerTx)
		}
		return cfg.Wasp.LPS * int64(cfg.Wasp.Duration.Seconds())
	}

	return int64(len(cfg.General.EventsToEmit) * cfg.LoopedConfig.ExecutionCount * cfg.General.Contracts * cfg.General.EventsPerTx)
}

type PauseData struct {
	StartBlock      uint64
	EndBlock        uint64
	TargetComponent string
	ContaineName    string
}

var ChaosPauses = []PauseData{}

// chaosPauseSyncFn pauses ranom container of the provided type for a random amount of time between 5 and 20 seconds
func chaosPauseSyncFn(ctx context.Context, dtc *chaos.DockerChaos, l zerolog.Logger, client *seth.Client, nodes *nodeset.Input, targetComponent string) ChaosPauseData {
	// var component ctf_test_env.EnvComponent
	var containerName string

	switch strings.ToLower(targetComponent) {
	case "chainlink":
		// component = randomNode.EnvComponent
		rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNode := nodes.Out.CLNodes[rand.Intn(len(nodes.Out.CLNodes)-1)+1]
		containerName = randomNode.Node.ContainerName
	case "postgres":
		containerName = nodes.DbInput.Name
		// component = randomNode.PostgresDb.EnvComponent
	default:
		return ChaosPauseData{Err: fmt.Errorf("unknown component %s", targetComponent)}
	}

	callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	pauseStartBlock, err := client.Client.BlockNumber(callCtx)
	cancel()
	if err != nil {
		return ChaosPauseData{Err: err}
	}
	pauseTimeSec := rand.Intn(20-5) + 5
	l.Info().Str("Container", containerName).Int("Pause time", pauseTimeSec).Msg("Pausing component")
	pauseTimeDur := time.Duration(pauseTimeSec) * time.Second

	err = dtc.Chaos(containerName, chaos.CmdPause, "")
	if err != nil {
		return ChaosPauseData{Err: fmt.Errorf("failed to pause docker container: %s, %w", containerName, err)}
	}
	time.Sleep(pauseTimeDur)
	if err := dtc.RemoveAll(); err != nil {
		return ChaosPauseData{Err: fmt.Errorf("failed to unpause docker container %s: %w", containerName, err)}
	}
	l.Info().Str("Container", containerName).Msg("Component unpaused")

	callCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	pauseEndBlock, err := client.Client.BlockNumber(callCtx)
	cancel()
	if err != nil {
		return ChaosPauseData{Err: err}
	}

	return ChaosPauseData{PauseData: PauseData{
		StartBlock:      pauseStartBlock,
		EndBlock:        pauseEndBlock,
		TargetComponent: targetComponent,
		ContaineName:    containerName,
	}}
}

type ChaosPauseData struct {
	Err       error
	PauseData PauseData
}

// executeChaosExperiment executes the configured chaos experiment, which consist of pausing CL node or Postgres containers
func executeChaosExperiment(ctx context.Context, l zerolog.Logger, nodes *nodeset.Input, sethClient *seth.Client, config *Config, errorCh chan error) {
	if config == nil || config.ChaosConfig == nil || config.ChaosConfig.ExperimentCount == 0 {
		errorCh <- nil
		return
	}

	dtc, err := chaos.NewDockerChaos(ctx)
	if err != nil {
		errorCh <- fmt.Errorf("failed to created docker-tc container: %w", err)
		return
	}

	chaosChan := make(chan ChaosPauseData, config.ChaosConfig.ExperimentCount)
	wg := &sync.WaitGroup{}

	go func() {
		// if we wanted to have more than 1 container paused, we'd need to make sure we aren't trying to pause an already paused one
		guardChan := make(chan struct{}, 1)

		for i := 0; i < config.ChaosConfig.ExperimentCount; i++ {
			i := i
			wg.Add(1)
			guardChan <- struct{}{}
			go func() {
				defer func() {
					<-guardChan
					wg.Done()
					current := i + 1
					l.Info().Str("Current/Total", fmt.Sprintf("%d/%d", current, config.ChaosConfig.ExperimentCount)).Msg("Done with experiment")
				}()
				chaosChan <- chaosPauseSyncFn(ctx, dtc, l, sethClient, nodes, config.ChaosConfig.TargetComponent)
				time.Sleep(10 * time.Second)
			}()
		}

		wg.Wait()

		close(chaosChan)
	}()

	go func() {
		var pauseData []PauseData
		for result := range chaosChan {
			if result.Err != nil {
				l.Err(result.Err).Msg("Error encountered during chaos experiment")
				errorCh <- result.Err
				return // Return on actual error
			}

			pauseData = append(pauseData, result.PauseData)
		}

		l.Info().Msg("All chaos experiments finished")
		errorCh <- nil // Only send nil once, after all errors have been handled and the channel is closed

		for _, p := range pauseData {
			l.Debug().Str("Target component", p.TargetComponent).Str("Container", p.ContaineName).Str("Block range", fmt.Sprintf("%d - %d", p.StartBlock, p.EndBlock)).Msgf("Details of executed chaos pause")
		}
	}()
}

// getEndBlockToWaitFor returns the end block to wait for based on chain id and finality tag provided in config
func getEndBlockToWaitFor(endBlock int64, config *automation.Automation) (int64, error) {
	if config.EVMNetworkSettings.FinalityTagEnabled != nil && *config.EVMNetworkSettings.FinalityTagEnabled {
		return endBlock + 1, nil
	}

	return endBlock + int64(*config.EVMNetworkSettings.FinalityDepth), nil //nolint:gosec // disable G115
}

const (
	defaultAmountOfUpkeeps = 2
)

var DefaultOCRRegistryConfig = contracts.KeeperRegistrySettings{
	PaymentPremiumPPB:    uint32(200000000),
	FlatFeeMicroLINK:     uint32(0),
	BlockCountPerTurn:    big.NewInt(10),
	CheckGasLimit:        uint32(2500000),
	StalenessSeconds:     big.NewInt(90000),
	GasCeilingMultiplier: uint16(1),
	MinUpkeepSpend:       big.NewInt(0),
	MaxPerformGas:        uint32(5000000),
	FallbackGasPrice:     big.NewInt(2e11),
	FallbackLinkPrice:    big.NewInt(2e18),
	MaxCheckDataSize:     uint32(5000),
	MaxPerformDataSize:   uint32(5000),
}

// uploadLogEmitterContracts uploads the configured number of log emitter contracts
func uploadLogEmitterContracts(l zerolog.Logger, t *testing.T, client *seth.Client, config *Config) []contracts.LogEmitter {
	logEmitters := make([]contracts.LogEmitter, 0)
	for i := 0; i < config.General.Contracts; i++ {
		logEmitter, err := contracts.DeployLogEmitterContract(l, client)
		logEmitters = append(logEmitters, logEmitter)
		require.NoError(t, err, "Error deploying log emitter contract")
		l.Info().Str("Contract address", logEmitter.Address().Hex()).Msg("Log emitter contract deployed")
		time.Sleep(200 * time.Millisecond)
	}

	return logEmitters
}

// assertUpkeepIDsUniqueness asserts that the provided upkeep IDs are unique
func assertUpkeepIDsUniqueness(upkeepIDs []*big.Int) error {
	upKeepIDSeen := make(map[int64]bool)
	for _, upkeepID := range upkeepIDs {
		if _, ok := upKeepIDSeen[upkeepID.Int64()]; ok {
			return fmt.Errorf("duplicate upkeep ID %d", upkeepID.Int64())
		}
		upKeepIDSeen[upkeepID.Int64()] = true
	}

	return nil
}

// assertContractAddressUniquneness asserts that the provided contract addresses are unique
func assertContractAddressUniquneness(logEmitters []contracts.LogEmitter) error {
	contractAddressSeen := make(map[string]bool)
	for _, logEmitter := range logEmitters {
		address := logEmitter.Address().String()
		if _, ok := contractAddressSeen[address]; ok {
			return fmt.Errorf("duplicate contract address %s", address)
		}
		contractAddressSeen[address] = true
	}

	return nil
}

// registerFiltersAndAssertUniquness registers the configured log filters and asserts that the filters are unique
// meaning that for each log emitter address and topic there is only one filter
func registerFiltersAndAssertUniquness(l zerolog.Logger, registry contracts.KeeperRegistry, upkeepIDs []*big.Int, logEmitters []contracts.LogEmitter, cfg *Config, upKeepsNeeded int) error {
	uniqueFilters := make(map[string]bool)

	upkeepIDIndex := 0
	for i := range logEmitters {
		for j := 0; j < len(cfg.General.EventsToEmit); j++ {
			emitterAddress := logEmitters[i].Address()
			topicID := cfg.General.EventsToEmit[j].ID

			upkeepID := upkeepIDs[upkeepIDIndex]
			l.Debug().Int("Upkeep id", int(upkeepID.Int64())).Str("Emitter address", emitterAddress.String()).Str("Topic", topicID.Hex()).Msg("Registering log trigger for log emitter")
			err := registerSingleTopicFilter(registry, upkeepID, emitterAddress, topicID)
			randomWait(150, 300)
			if err != nil {
				return fmt.Errorf("%w: Error registering log trigger for log emitter %s", err, emitterAddress.String())
			}

			if i%10 == 0 {
				l.Info().Msgf("Registered log trigger for topic %d for log emitter %d/%d", j, i+1, len(logEmitters))
			}

			key := fmt.Sprintf("%s-%s", emitterAddress.String(), topicID.Hex())
			if _, ok := uniqueFilters[key]; ok {
				return fmt.Errorf("duplicate filter %s", key)
			}
			uniqueFilters[key] = true
			upkeepIDIndex++
		}
	}

	if upKeepsNeeded != len(uniqueFilters) {
		return fmt.Errorf("number of unique filters should be equal to number of upkeeps. Expected %d. Got %d", upKeepsNeeded, len(uniqueFilters))
	}

	return nil
}

// checkIfAllNodesHaveLogCount checks if all CL nodes have the expected log count for the provided block range and expected filters
// It will retry until the provided duration is reached or until all nodes have the expected log count
func checkIfAllNodesHaveLogCount(duration string, startBlock, endBlock int64, expectedLogCount int, expectedFilters []ExpectedFilter, l zerolog.Logger, coreLogger common_logger.SugaredLogger, testEnv *logPollerEnvironment) (bool, error) {
	logCountWaitDuration, err := time.ParseDuration(duration)
	if err != nil {
		return false, err
	}
	endTime := time.Now().Add(logCountWaitDuration)

	// not using gomega here, because I want to see which logs were missing
	allNodesLogCountMatches := false
	for time.Now().Before(endTime) {
		logCountMatches, clErr := nodesHaveExpectedLogCount(startBlock, endBlock, big.NewInt(testEnv.chainClient.ChainID), expectedLogCount, expectedFilters, l, coreLogger, testEnv)
		if clErr != nil {
			l.Warn().
				Err(clErr).
				Msg("Error checking if CL nodes have expected log count. Retrying...")
		}
		if logCountMatches {
			allNodesLogCountMatches = true
			break
		}
		l.Warn().
			Msg("At least one CL node did not have expected log count. Retrying...")
		time.Sleep(10 * time.Second)
	}

	return allNodesLogCountMatches, nil
}

type logPollerEnvironment struct {
	logger zerolog.Logger
	nodes  *nodeset.Input

	chainClient *seth.Client

	config *automation.Automation

	logEmitters   []contracts.LogEmitter
	upkeepIDs     []*big.Int
	upKeepsNeeded int

	registry  contracts.KeeperRegistry
	registrar contracts.KeeperRegistrar
	linkToken contracts.LinkToken
}

func newLpTestEnvironment(chainClient *seth.Client, config *automation.Automation, in *de.Cfg) (*logPollerEnvironment, error) {
	lpTestEnv := &logPollerEnvironment{
		chainClient: chainClient,
		config:      config,
		nodes:       in.NodeSets[0],
	}

	if err := lpTestEnv.loadContracts(); err != nil {
		return nil, err
	}

	return lpTestEnv, nil
}

func (l *logPollerEnvironment) dbPort() int {
	if l.nodes.DbInput.Port != 0 {
		return l.nodes.DbInput.Port
	}

	return postgres.ExposedStaticPort
}

func (l *logPollerEnvironment) loadContracts() error {
	if err := l.loadLINK(l.config.DeployedContracts.LinkToken); err != nil {
		return fmt.Errorf("error loading link token contract: %w", err)
	}

	if err := l.loadRegistry(l.config.DeployedContracts.Registry, l.config.DeployedContracts.ChainModule); err != nil {
		return fmt.Errorf("error loading registry contract: %w", err)
	}

	if l.registry.RegistryOwnerAddress().String() != l.chainClient.MustGetRootKeyAddress().String() {
		return errors.New("registry owner address is not the root key address")
	}

	if err := l.loadRegistrar(l.config.DeployedContracts.Registrar); err != nil {
		return fmt.Errorf("error loading registrar contract: %w", err)
	}

	return nil
}

func (l *logPollerEnvironment) loadLINK(address string) error {
	linkToken, err := contracts.LoadLinkTokenContract(l.logger, l.chainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	l.linkToken = linkToken
	l.logger.Info().Str("LINK Token Address", l.linkToken.Address()).Msg("Successfully loaded LINK Token")
	return nil
}

func (l *logPollerEnvironment) loadRegistry(registryAddress, chainModuleAddress string) error {
	registry, err := contracts.LoadKeeperRegistry(l.logger, l.chainClient, common.HexToAddress(registryAddress), contracts.RegistryVersion_2_1, common.HexToAddress(chainModuleAddress))
	if err != nil {
		return err
	}
	l.registry = registry
	l.logger.Info().Str("ChainModule Address", chainModuleAddress).Str("Registry Address", l.registry.Address()).Msg("Successfully loaded Registry")
	return nil
}

func (l *logPollerEnvironment) loadRegistrar(address string) error {
	if l.registry == nil {
		return errors.New("registry must be deployed or loaded before registrar")
	}
	// l.RegistrarSettings.RegistryAddr = l.registry.Address()
	registrar, err := contracts.LoadKeeperRegistrar(l.chainClient, common.HexToAddress(address), contracts.RegistryVersion_2_1)
	if err != nil {
		return err
	}
	l.logger.Info().Str("Registrar Address", registrar.Address()).Msg("Successfully loaded Registrar")
	l.registrar = registrar
	return nil
}

// waitForAllNodesToHaveExpectedFiltersRegisteredOrFail waits until all nodes have expected filters registered until timeout
func waitForAllNodesToHaveExpectedFiltersRegisteredOrFail(ctx context.Context, l zerolog.Logger, coreLogger common_logger.SugaredLogger, t *testing.T, testEnv *logPollerEnvironment, expectedFilters []ExpectedFilter) {
	// Make sure that all nodes have expected filters registered before starting to emit events

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		hasFilters := false
		for i := 1; i < len(testEnv.nodes.NodeSpecs); i++ {
			nodeName := testEnv.nodes.Out.CLNodes[i].Node.ContainerName
			l.Info().
				Str("Node name", nodeName).
				Msg("Fetching filters from log poller's DB")
			var message string
			var err error

			hasFilters, message, err = nodeHasExpectedFilters(ctx, expectedFilters, coreLogger, big.NewInt(testEnv.chainClient.ChainID), i, testEnv.dbPort())
			if !hasFilters || err != nil {
				if message == "" {
					message = err.Error()
				}
				l.Warn().
					Str("Details", message).
					Msg("Some filters were missing, but we will retry")
				break
			}
		}
		g.Expect(hasFilters).To(gomega.BeTrue(), "Not all expected filters were found in the DB")
	}, "5m", "10s").Should(gomega.Succeed())

	l.Info().
		Msg("All nodes have expected filters registered")
	l.Info().
		Int("Count", len(expectedFilters)).
		Msg("Expected filters count")
}

// waitUntilNodesHaveTheSameLogsAsEvm checks whether all CL nodes have the same number of logs as EVM node
// if not, then it prints missing logs and wait for some time and checks again
func waitUntilNodesHaveTheSameLogsAsEvm(l zerolog.Logger, coreLogger common_logger.SugaredLogger, t *testing.T, allNodesLogCountMatches bool, lpTestEnv *logPollerEnvironment, config *Config, startBlock, endBlock int64, waitDuration string) {
	logCountWaitDuration, err := time.ParseDuration(waitDuration)
	require.NoError(t, err, "Error parsing log count wait duration")

	allNodesHaveAllExpectedLogs := false
	if !allNodesLogCountMatches {
		missingLogs, err := missingLogs(startBlock, endBlock, lpTestEnv, l, coreLogger, config)
		if err == nil {
			if !missingLogs.IsEmpty() {
				printMissingLogsInfo(missingLogs, l, config)
			} else {
				allNodesHaveAllExpectedLogs = true
				l.Info().Msg("All CL nodes have all the logs that EVM node has")
			}
		}
	}

	require.True(t, allNodesLogCountMatches, "Not all CL nodes had expected log count after %s", logCountWaitDuration)

	// Wait until all CL nodes have exactly the same logs emitted by test contracts as the EVM node has
	// but only in the rare case that first attempt to do it failed (basically here want to know not only
	// if log count matches, but whether details of every single log match)
	if !allNodesHaveAllExpectedLogs {
		logConsistencyWaitDuration := "5m"
		l.Info().
			Str("Duration", logConsistencyWaitDuration).
			Msg("Waiting for CL nodes to have all the logs that EVM node has")

		gom := gomega.NewGomegaWithT(t)
		gom.Eventually(func(g gomega.Gomega) {
			missingLogs, err := missingLogs(startBlock, endBlock, lpTestEnv, l, coreLogger, config)
			if err != nil {
				l.Warn().
					Err(err).
					Msg("Error getting missing logs. Retrying...")
			}

			if !missingLogs.IsEmpty() {
				printMissingLogsInfo(missingLogs, l, config)
			}
			g.Expect(missingLogs.IsEmpty()).To(gomega.BeTrue(), "Some CL nodes were missing logs")
		}, logConsistencyWaitDuration, "10s").Should(gomega.Succeed())
	}
}
