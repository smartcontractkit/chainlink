package logpoller

import (
	"bytes"
	"context"
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
	"github.com/rs/zerolog"
	"github.com/scylladb/go-reflectx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"

	ac "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_compatible_utils"
	le "github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/log_emitter"
	cltypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctf_concurrency "github.com/smartcontractkit/chainlink-testing-framework/lib/concurrency"
	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	lp_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/log_poller"
	core_logger "github.com/smartcontractkit/chainlink/v2/core/logger"
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

// NewORM returns a new logpoller.orm instance
func NewORM(logger core_logger.SugaredLogger, chainID *big.Int, postgresDb *ctf_test_env.PostgresDb) (logpoller.ORM, *sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", postgresDb.ExternalPort, postgresDb.User, postgresDb.Password, postgresDb.DbName)
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

// GetExpectedFilters returns a slice of ExpectedFilter structs based on the provided log emitters and config
func GetExpectedFilters(logEmitters []*contracts.LogEmitter, cfg *lp_config.Config) []ExpectedFilter {
	expectedFilters := make([]ExpectedFilter, 0)
	for _, emitter := range logEmitters {
		for _, event := range cfg.General.EventsToEmit {
			expectedFilters = append(expectedFilters, ExpectedFilter{
				emitterAddress: (*emitter).Address(),
				topic:          event.ID,
			})
		}
	}

	return expectedFilters
}

// NodeHasExpectedFilters returns true if the provided node has all the expected filters registered
func NodeHasExpectedFilters(ctx context.Context, expectedFilters []ExpectedFilter, logger core_logger.SugaredLogger, chainID *big.Int, postgresDb *ctf_test_env.PostgresDb) (bool, string, error) {
	orm, db, err := NewORM(logger, chainID, postgresDb)
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
	for i := 0; i < length; i++ {
		result[i] = i
	}

	return result
}

// getStringSlice returns a slice of strings of the provided length
func getStringSlice(length int) []string {
	result := make([]string, length)
	for i := 0; i < length; i++ {
		result[i] = "amazing event"
	}

	return result
}

// LogPollerHasFinalisedEndBlock returns true if all CL nodes have finalised processing the provided end block
func LogPollerHasFinalisedEndBlock(endBlock int64, chainID *big.Int, l zerolog.Logger, coreLogger core_logger.SugaredLogger, nodes *test_env.ClCluster) (bool, error) {
	wg := &sync.WaitGroup{}

	type boolQueryResult struct {
		nodeName       string
		hasFinalised   bool
		finalizedBlock int64
		err            error
	}

	endBlockCh := make(chan boolQueryResult, len(nodes.Nodes)-1)
	ctx, cancelFn := context.WithCancel(context.Background())

	for i := 1; i < len(nodes.Nodes); i++ {
		wg.Add(1)

		go func(clNode *test_env.ClNode, r chan boolQueryResult) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				orm, db, err := NewORM(coreLogger, chainID, clNode.PostgresDb)
				if err != nil {
					r <- boolQueryResult{
						nodeName:     clNode.ContainerName,
						hasFinalised: false,
						err:          err,
					}
				}

				defer db.Close()

				latestBlock, err := orm.SelectLatestBlock(ctx)
				if err != nil {
					r <- boolQueryResult{
						nodeName:     clNode.ContainerName,
						hasFinalised: false,
						err:          err,
					}
				}

				r <- boolQueryResult{
					nodeName:       clNode.ContainerName,
					finalizedBlock: latestBlock.FinalizedBlockNumber,
					hasFinalised:   latestBlock.FinalizedBlockNumber > endBlock,
					err:            nil,
				}

			}
		}(nodes.Nodes[i], endBlockCh)
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

			if len(foundMap) == len(nodes.Nodes)-1 {
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

// ClNodesHaveExpectedLogCount returns true if all CL nodes have the expected log count in the provided block range and matching the provided filters
func ClNodesHaveExpectedLogCount(startBlock, endBlock int64, chainID *big.Int, expectedLogCount int, expectedFilters []ExpectedFilter, l zerolog.Logger, coreLogger core_logger.SugaredLogger, nodes *test_env.ClCluster) (bool, error) {
	wg := &sync.WaitGroup{}

	type logQueryResult struct {
		nodeName         string
		logCount         int
		hasExpectedCount bool
		err              error
	}

	resultChan := make(chan logQueryResult, len(nodes.Nodes)-1)
	ctx, cancelFn := context.WithCancel(context.Background())

	for i := 1; i < len(nodes.Nodes); i++ {
		wg.Add(1)

		go func(clNode *test_env.ClNode, resultChan chan logQueryResult) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				orm, db, err := NewORM(coreLogger, chainID, clNode.PostgresDb)
				if err != nil {
					resultChan <- logQueryResult{
						nodeName:         clNode.ContainerName,
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
							nodeName:         clNode.ContainerName,
							logCount:         0,
							hasExpectedCount: false,
							err:              err,
						}
					}

					foundLogsCount += len(logs)
				}

				resultChan <- logQueryResult{
					nodeName:         clNode.ContainerName,
					logCount:         foundLogsCount,
					hasExpectedCount: foundLogsCount >= expectedLogCount,
					err:              nil,
				}
			}
		}(nodes.Nodes[i], resultChan)
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

			if len(foundMap) == len(nodes.Nodes)-1 {
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

// GetMissingLogs returns a map of CL node name to missing logs in that node compared to EVM node to which the provided evm client is connected
func GetMissingLogs(
	startBlock, endBlock int64,
	logEmitters []*contracts.LogEmitter,
	client *seth.Client,
	clnodeCluster *test_env.ClCluster,
	l zerolog.Logger,
	coreLogger core_logger.SugaredLogger,
	cfg *lp_config.Config,
) (MissingLogs, error) {
	wg := &sync.WaitGroup{}

	type dbQueryResult struct {
		err      error
		nodeName string
		logs     []logpoller.Log
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	resultCh := make(chan dbQueryResult, len(clnodeCluster.Nodes)-1)

	for i := 1; i < len(clnodeCluster.Nodes); i++ {
		wg.Add(1)

		go func(ctx context.Context, i int, r chan dbQueryResult) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				l.Warn().Msg("Context cancelled. Terminating fetching logs from log poller's DB")
				return
			default:
				nodeName := clnodeCluster.Nodes[i].ContainerName

				l.Debug().Str("Node name", nodeName).Msg("Fetching log poller logs")
				orm, db, err := NewORM(coreLogger, big.NewInt(client.ChainID), clnodeCluster.Nodes[i].PostgresDb)
				if err != nil {
					r <- dbQueryResult{
						err:      err,
						nodeName: nodeName,
						logs:     []logpoller.Log{},
					}
				}

				defer db.Close()
				logs := make([]logpoller.Log, 0)

				for j := 0; j < len(logEmitters); j++ {
					address := (*logEmitters[j]).Address()

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

	allLogsInEVMNode, err := getEVMLogs(ctx, startBlock, endBlock, logEmitters, client, l, cfg)
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

	missingCh := make(chan missingLogResult, len(clnodeCluster.Nodes)-1)
	for i := 1; i < len(clnodeCluster.Nodes); i++ {
		wg.Add(1)

		go func(i int, result chan missingLogResult) {
			defer wg.Done()
			nodeName := clnodeCluster.Nodes[i].ContainerName
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
						logPollerLog.Address == evmLog.Address && logPollerLog.BlockHash == evmLog.BlockHash && bytes.Equal(logPollerLog.Topics[0][:], evmLog.Topics[0].Bytes()) {
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

	expectedTotalLogsEmitted := GetExpectedLogCount(cfg)
	if int64(len(allLogsInEVMNode)) != expectedTotalLogsEmitted {
		l.Warn().
			Str("Actual/Expected", fmt.Sprintf("%d/%d", expectedTotalLogsEmitted, len(allLogsInEVMNode))).
			Msg("Actual number of logs found on EVM nodes differs from expected ones. Most probably this is a bug in the test")
	}

	return missingLogs, nil
}

// PrintMissingLogsInfo prints various useful information about the missing logs
func PrintMissingLogsInfo(missingLogs map[string][]geth_types.Log, l zerolog.Logger, cfg *lp_config.Config) {
	var findHumanName = func(topic common.Hash) string {
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
func getEVMLogs(ctx context.Context, startBlock, endBlock int64, logEmitters []*contracts.LogEmitter, client *seth.Client, l zerolog.Logger, cfg *lp_config.Config) ([]geth_types.Log, error) {
	allLogsInEVMNode := make([]geth_types.Log, 0)
	for j := 0; j < len(logEmitters); j++ {
		address := (*logEmitters[j]).Address()
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

// ExecuteGenerator executes the configured generator and returns the total number of logs emitted
func ExecuteGenerator(t *testing.T, cfg *lp_config.Config, client *seth.Client, logEmitters []*contracts.LogEmitter) (int, error) {
	if *cfg.General.Generator == lp_config.GeneratorType_WASP {
		return runWaspGenerator(t, cfg, logEmitters)
	}

	return runLoopedGenerator(t, cfg, client, logEmitters)
}

// runWaspGenerator runs the wasp generator and returns the total number of logs emitted
func runWaspGenerator(t *testing.T, cfg *lp_config.Config, logEmitters []*contracts.LogEmitter) (int, error) {
	l := logging.GetTestLogger(t)

	var RPSprime int64

	// if LPS is set, we need to calculate based on countract count and events per transaction
	if *cfg.Wasp.LPS > 0 {
		RPSprime = *cfg.Wasp.LPS / int64(*cfg.General.Contracts) / int64(*cfg.General.EventsPerTx) / int64(len(cfg.General.EventsToEmit))

		if RPSprime < 1 {
			return 0, fmt.Errorf("invalid load configuration, effective RPS would have been zero. Adjust LPS, contracts count, events per tx or events to emit")
		}
	}

	// if RPS is set simply split it between contracts
	if *cfg.Wasp.RPS > 0 {
		RPSprime = *cfg.Wasp.RPS / int64(*cfg.General.Contracts)
	}

	counter := &Counter{
		mu:    &sync.Mutex{},
		value: 0,
	}

	p := wasp.NewProfile()

	for _, logEmitter := range logEmitters {
		g, err := wasp.NewGenerator(&wasp.Config{
			T:                     t,
			LoadType:              wasp.RPS,
			GenName:               fmt.Sprintf("log_poller_gen_%s", (*logEmitter).Address().String()),
			RateLimitUnitDuration: cfg.Wasp.RateLimitUnitDuration.Duration,
			CallTimeout:           cfg.Wasp.CallTimeout.Duration,
			Schedule: wasp.Plain(
				RPSprime,
				cfg.Wasp.Duration.Duration,
			),
			Gun: NewLogEmitterGun(
				logEmitter,
				cfg.General.EventsToEmit,
				*cfg.General.EventsPerTx,
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
	emitter      *contracts.LogEmitter
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
func runLoopedGenerator(t *testing.T, cfg *lp_config.Config, client *seth.Client, logEmitters []*contracts.LogEmitter) (int, error) {
	l := logging.GetTestLogger(t)

	tasks := make([]logEmissionTask, 0)
	for i := 0; i < *cfg.LoopedConfig.ExecutionCount; i++ {
		for _, logEmitter := range logEmitters {
			tasks = append(tasks, logEmissionTask{
				emitter:      logEmitter,
				eventsToEmit: cfg.General.EventsToEmit,
				eventsPerTx:  *cfg.General.EventsPerTx,
			})
		}
	}

	l.Info().Int("Total tasks", len(tasks)).Msg("Starting to emit events")

	var atomicCounter = atomic.Int32{}

	var emitAllEventsFn = func(resultCh chan emittedLogsData, errorCh chan error, _ int, task logEmissionTask) {
		current := atomicCounter.Add(1)

		address := (*task.emitter).Address().String()

		for _, event := range cfg.General.EventsToEmit {
			l.Debug().Str("Emitter address", address).Str("Event type", event.Name).Str("index", fmt.Sprintf("%d/%d", current, *cfg.LoopedConfig.ExecutionCount)).Msg("Emitting log from emitter")
			var err error
			switch event.Name {
			case "Log1":
				_, err = client.Decode((*task.emitter).EmitLogIntsFromKey(getIntSlice(*cfg.General.EventsPerTx), client.AnySyncedKey()))
			case "Log2":
				_, err = client.Decode((*task.emitter).EmitLogIntsIndexedFromKey(getIntSlice(*cfg.General.EventsPerTx), client.AnySyncedKey()))
			case "Log3":
				_, err = client.Decode((*task.emitter).EmitLogStringsFromKey(getStringSlice(*cfg.General.EventsPerTx), client.AnySyncedKey()))
			case "Log4":
				_, err = client.Decode((*task.emitter).EmitLogIntMultiIndexedFromKey(1, 1, *cfg.General.EventsPerTx, client.AnySyncedKey()))
			default:
				err = fmt.Errorf("unknown event name: %s", event.Name)
			}

			if err != nil {
				errorCh <- err
				return
			}
			randomWait(*cfg.LoopedConfig.MinEmitWaitTimeMs, *cfg.LoopedConfig.MaxEmitWaitTimeMs)

			if (current)%10 == 0 {
				l.Info().Str("Emitter address", address).Str("Index", fmt.Sprintf("%d/%d", current, *cfg.LoopedConfig.ExecutionCount)).Msgf("Emitted all %d events", len(cfg.General.EventsToEmit))
			}
		}

		resultCh <- emittedLogsData{
			*cfg.General.EventsPerTx * len(cfg.General.EventsToEmit),
		}
	}

	executor := ctf_concurrency.NewConcurrentExecutor[emittedLogsData, emittedLogsData, logEmissionTask](l)
	r, err := executor.Execute(int(*client.Cfg.EphemeralAddrs), tasks, emitAllEventsFn)

	if err != nil {
		return 0, err
	}

	var total int
	for _, result := range r {
		total += result.count
	}

	return total, nil
}

// GetExpectedLogCount returns the expected number of logs to be emitted based on the provided config
func GetExpectedLogCount(cfg *lp_config.Config) int64 {
	if *cfg.General.Generator == lp_config.GeneratorType_WASP {
		if *cfg.Wasp.RPS != 0 {
			return *cfg.Wasp.RPS * int64(cfg.Wasp.Duration.Seconds()) * int64(*cfg.General.EventsPerTx)
		}
		return *cfg.Wasp.LPS * int64(cfg.Wasp.Duration.Duration.Seconds())
	}

	return int64(len(cfg.General.EventsToEmit) * *cfg.LoopedConfig.ExecutionCount * *cfg.General.Contracts * *cfg.General.EventsPerTx)
}

type PauseData struct {
	StartBlock      uint64
	EndBlock        uint64
	TargetComponent string
	ContaineName    string
}

var ChaosPauses = []PauseData{}

// chaosPauseSyncFn pauses ranom container of the provided type for a random amount of time between 5 and 20 seconds
func chaosPauseSyncFn(l zerolog.Logger, client *seth.Client, cluster *test_env.ClCluster, targetComponent string) ChaosPauseData {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	randomNode := cluster.Nodes[rand.Intn(len(cluster.Nodes)-1)+1]
	var component ctf_test_env.EnvComponent

	switch strings.ToLower(targetComponent) {
	case "chainlink":
		component = randomNode.EnvComponent
	case "postgres":
		component = randomNode.PostgresDb.EnvComponent
	default:
		return ChaosPauseData{Err: fmt.Errorf("unknown component %s", targetComponent)}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pauseStartBlock, err := client.Client.BlockNumber(ctx)
	if err != nil {
		return ChaosPauseData{Err: err}
	}
	pauseTimeSec := rand.Intn(20-5) + 5
	l.Info().Str("Container", component.ContainerName).Int("Pause time", pauseTimeSec).Msg("Pausing component")
	pauseTimeDur := time.Duration(pauseTimeSec) * time.Second
	err = component.ChaosPause(l, pauseTimeDur)
	if err != nil {
		return ChaosPauseData{Err: err}
	}
	l.Info().Str("Container", component.ContainerName).Msg("Component unpaused")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pauseEndBlock, err := client.Client.BlockNumber(ctx)
	if err != nil {
		return ChaosPauseData{Err: err}
	}

	return ChaosPauseData{PauseData: PauseData{
		StartBlock:      pauseStartBlock,
		EndBlock:        pauseEndBlock,
		TargetComponent: targetComponent,
		ContaineName:    component.ContainerName,
	}}
}

type ChaosPauseData struct {
	Err       error
	PauseData PauseData
}

// ExecuteChaosExperiment executes the configured chaos experiment, which consist of pausing CL node or Postgres containers
func ExecuteChaosExperiment(l zerolog.Logger, testEnv *test_env.CLClusterTestEnv, sethClient *seth.Client, testConfig *tc.TestConfig, errorCh chan error) {
	if testConfig == nil || testConfig.LogPoller.ChaosConfig == nil || *testConfig.LogPoller.ChaosConfig.ExperimentCount == 0 {
		errorCh <- nil
		return
	}

	chaosChan := make(chan ChaosPauseData, *testConfig.LogPoller.ChaosConfig.ExperimentCount)
	wg := &sync.WaitGroup{}

	go func() {
		// if we wanted to have more than 1 container paused, we'd need to make sure we aren't trying to pause an already paused one
		guardChan := make(chan struct{}, 1)

		for i := 0; i < *testConfig.LogPoller.ChaosConfig.ExperimentCount; i++ {
			i := i
			wg.Add(1)
			guardChan <- struct{}{}
			go func() {
				defer func() {
					<-guardChan
					wg.Done()
					current := i + 1
					l.Info().Str("Current/Total", fmt.Sprintf("%d/%d", current, *testConfig.LogPoller.ChaosConfig.ExperimentCount)).Msg("Done with experiment")
				}()
				chaosChan <- chaosPauseSyncFn(l, sethClient, testEnv.ClCluster, *testConfig.LogPoller.ChaosConfig.TargetComponent)
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

// GetEndBlockToWaitFor returns the end block to wait for based on chain id and finality tag provided in config
func GetEndBlockToWaitFor(endBlock int64, network blockchain.EVMNetwork, cfg *lp_config.Config) (int64, error) {
	if *cfg.General.UseFinalityTag {
		return endBlock + 1, nil
	}

	if network.FinalityDepth > math.MaxInt64 {
		return -1, fmt.Errorf("finality depth overflows int64: %d", network.FinalityDepth)
	}
	return endBlock + int64(network.FinalityDepth), nil
}

const (
	defaultAmountOfUpkeeps = 2
)

var (
	DefaultOCRRegistryConfig = contracts.KeeperRegistrySettings{
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
)

// SetupLogPollerTestDocker starts the DON and private Ethereum network
func SetupLogPollerTestDocker(
	t *testing.T,
	registryVersion ethereum.KeeperRegistryVersion,
	registryConfig contracts.KeeperRegistrySettings,
	upkeepsNeeded int,
	finalityTagEnabled bool,
	testConfig *tc.TestConfig,
	logScannerSettings test_env.ChainlinkNodeLogScannerSettings,
) (
	*seth.Client,
	[]*nodeclient.ChainlinkClient,
	contracts.LinkToken,
	contracts.KeeperRegistry,
	contracts.KeeperRegistrar,
	*test_env.CLClusterTestEnv,
	*blockchain.EVMNetwork,
) {
	l := logging.GetTestLogger(t)

	// Add registry version to config
	registryConfig.RegistryVersion = registryVersion
	network := networks.MustGetSelectedNetworkConfig(testConfig.Network)[0]

	//launch the environment
	var env *test_env.CLClusterTestEnv
	chainlinkNodeFunding := 0.5
	l.Debug().Msgf("Funding amount: %f", chainlinkNodeFunding)
	clNodesCount := 5

	var evmNetworkExtraSettingsFn = func(network *blockchain.EVMNetwork) *blockchain.EVMNetwork {
		// we need it, because by default finality depth is 0 for our simulated network
		if network.Simulated && !finalityTagEnabled {
			network.FinalityDepth = 10
		}
		network.FinalityTag = finalityTagEnabled
		return network
	}

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, testConfig)
	require.NoError(t, err, "Error building ethereum network config")

	env, err = test_env.NewCLTestEnvBuilder().
		WithTestConfig(testConfig).
		WithTestInstance(t).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithCLNodes(clNodesCount).
		WithEVMNetworkOptions(evmNetworkExtraSettingsFn).
		WithChainlinkNodeLogScanner(logScannerSettings).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "Error deploying test environment")

	evmNetwork, err := env.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	chainClient, err := utils.TestAwareSethClient(t, testConfig, evmNetwork)
	require.NoError(t, err, "Error getting seth client")

	err = actions.FundChainlinkNodesFromRootAddress(l, chainClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()), big.NewFloat(chainlinkNodeFunding))
	require.NoError(t, err, "Failed to fund the nodes")

	nodeClients := env.ClCluster.NodeAPIs()
	workerNodes := nodeClients[1:]

	linkToken, err := contracts.DeployLinkTokenContract(l, chainClient)
	require.NoError(t, err, "Error deploying LINK token")

	wethToken, err := contracts.DeployWETHTokenContract(l, chainClient)
	require.NoError(t, err, "Error deploying weth token contract")

	// This feed is used for both eth/usd and link/usd
	ethUSDFeed, err := contracts.DeployMockETHUSDFeed(chainClient, registryConfig.FallbackLinkPrice)
	require.NoError(t, err, "Error deploying eth usd feed contract")

	linkBalance, err := linkToken.BalanceOf(context.Background(), chainClient.MustGetRootKeyAddress().Hex())
	require.NoError(t, err, "Error getting LINK balance")

	l.Info().Str("Balance", big.NewInt(0).Div(linkBalance, big.NewInt(1e18)).String()).Msg("LINK balance")
	minLinkBalanceSingleNode := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(9))
	minLinkBalance := big.NewInt(0).Mul(minLinkBalanceSingleNode, big.NewInt(int64(upkeepsNeeded)))
	if linkBalance.Cmp(minLinkBalance) < 0 {
		require.FailNowf(t, "Not enough LINK", "Not enough LINK to run the test. Need at least %s. but has only %s", big.NewInt(0).Div(minLinkBalance, big.NewInt(1e18)).String(), big.NewInt(0).Div(linkBalance, big.NewInt(1e18)).String())
	}

	registry, registrar := actions.DeployAutoOCRRegistryAndRegistrar(
		t,
		chainClient,
		registryVersion,
		registryConfig,
		linkToken,
		wethToken,
		ethUSDFeed,
	)

	// Fund the registry with LINK
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(defaultAmountOfUpkeeps))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	err = actions.CreateOCRKeeperJobsLocal(l, nodeClients, registry.Address(), network.ChainID, 0, registryVersion)
	require.NoError(t, err, "Error creating OCR Keeper Jobs")
	ocrConfig, err := actions.BuildAutoOCR2ConfigVarsLocal(l, workerNodes, registryConfig, registrar.Address(), 30*time.Second, registry.RegistryOwnerAddress(), registry.ChainModuleAddress(), registry.ReorgProtectionEnabled())
	require.NoError(t, err, "Error building OCR config vars")
	err = registry.SetConfigTypeSafe(ocrConfig)
	require.NoError(t, err, "Registry config should be set successfully")

	return chainClient, nodeClients, linkToken, registry, registrar, env, &network
}

// UploadLogEmitterContracts uploads the configured number of log emitter contracts
func UploadLogEmitterContracts(l zerolog.Logger, t *testing.T, client *seth.Client, testConfig *tc.TestConfig) []*contracts.LogEmitter {
	logEmitters := make([]*contracts.LogEmitter, 0)
	for i := 0; i < *testConfig.LogPoller.General.Contracts; i++ {
		logEmitter, err := contracts.DeployLogEmitterContract(l, client)
		logEmitters = append(logEmitters, &logEmitter)
		require.NoError(t, err, "Error deploying log emitter contract")
		l.Info().Str("Contract address", logEmitter.Address().Hex()).Msg("Log emitter contract deployed")
		time.Sleep(200 * time.Millisecond)
	}

	return logEmitters
}

// AssertUpkeepIdsUniqueness asserts that the provided upkeep IDs are unique
func AssertUpkeepIdsUniqueness(upkeepIDs []*big.Int) error {
	upKeepIdSeen := make(map[int64]bool)
	for _, upkeepID := range upkeepIDs {
		if _, ok := upKeepIdSeen[upkeepID.Int64()]; ok {
			return fmt.Errorf("Duplicate upkeep ID %d", upkeepID.Int64())
		}
		upKeepIdSeen[upkeepID.Int64()] = true
	}

	return nil
}

// AssertContractAddressUniquneness asserts that the provided contract addresses are unique
func AssertContractAddressUniquneness(logEmitters []*contracts.LogEmitter) error {
	contractAddressSeen := make(map[string]bool)
	for _, logEmitter := range logEmitters {
		address := (*logEmitter).Address().String()
		if _, ok := contractAddressSeen[address]; ok {
			return fmt.Errorf("Duplicate contract address %s", address)
		}
		contractAddressSeen[address] = true
	}

	return nil
}

// RegisterFiltersAndAssertUniquness registers the configured log filters and asserts that the filters are unique
// meaning that for each log emitter address and topic there is only one filter
func RegisterFiltersAndAssertUniquness(l zerolog.Logger, registry contracts.KeeperRegistry, upkeepIDs []*big.Int, logEmitters []*contracts.LogEmitter, cfg *lp_config.Config, upKeepsNeeded int) error {
	uniqueFilters := make(map[string]bool)

	upkeepIdIndex := 0
	for i := 0; i < len(logEmitters); i++ {
		for j := 0; j < len(cfg.General.EventsToEmit); j++ {
			emitterAddress := (*logEmitters[i]).Address()
			topicId := cfg.General.EventsToEmit[j].ID

			upkeepID := upkeepIDs[upkeepIdIndex]
			l.Debug().Int("Upkeep id", int(upkeepID.Int64())).Str("Emitter address", emitterAddress.String()).Str("Topic", topicId.Hex()).Msg("Registering log trigger for log emitter")
			err := registerSingleTopicFilter(registry, upkeepID, emitterAddress, topicId)
			randomWait(150, 300)
			if err != nil {
				return fmt.Errorf("%w: Error registering log trigger for log emitter %s", err, emitterAddress.String())
			}

			if i%10 == 0 {
				l.Info().Msgf("Registered log trigger for topic %d for log emitter %d/%d", j, i, len(logEmitters))
			}

			key := fmt.Sprintf("%s-%s", emitterAddress.String(), topicId.Hex())
			if _, ok := uniqueFilters[key]; ok {
				return fmt.Errorf("Duplicate filter %s", key)
			}
			uniqueFilters[key] = true
			upkeepIdIndex++
		}
	}

	if upKeepsNeeded != len(uniqueFilters) {
		return fmt.Errorf("Number of unique filters should be equal to number of upkeeps. Expected %d. Got %d", upKeepsNeeded, len(uniqueFilters))
	}

	return nil
}

// FluentlyCheckIfAllNodesHaveLogCount checks if all CL nodes have the expected log count for the provided block range and expected filters
// It will retry until the provided duration is reached or until all nodes have the expected log count
func FluentlyCheckIfAllNodesHaveLogCount(duration string, startBlock, endBlock int64, expectedLogCount int, expectedFilters []ExpectedFilter, l zerolog.Logger, coreLogger core_logger.SugaredLogger, testEnv *test_env.CLClusterTestEnv, chainId int64) (bool, error) {
	logCountWaitDuration, err := time.ParseDuration(duration)
	if err != nil {
		return false, err
	}
	endTime := time.Now().Add(logCountWaitDuration)

	// not using gomega here, because I want to see which logs were missing
	allNodesLogCountMatches := false
	for time.Now().Before(endTime) {
		logCountMatches, clErr := ClNodesHaveExpectedLogCount(startBlock, endBlock, big.NewInt(chainId), expectedLogCount, expectedFilters, l, coreLogger, testEnv.ClCluster)
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
