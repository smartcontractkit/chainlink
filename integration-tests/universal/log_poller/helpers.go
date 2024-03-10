package logpoller

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"strings"
	"sync"
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

	"github.com/smartcontractkit/wasp"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	cltypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	le "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	core_logger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	lp_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/log_poller"
)

var (
	EmitterABI, _      = abi.JSON(strings.NewReader(le.LogEmitterABI))
	automationUtilsABI = cltypes.MustGetABI(automation_utils_2_1.AutomationUtilsABI)
	bytes0             = [32]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	} // bytes representation of 0x0000000000000000000000000000000000000000000000000000000000000000

)

var registerSingleTopicFilter = func(registry contracts.KeeperRegistry, upkeepID *big.Int, emitterAddress common.Address, topic common.Hash) error {
	logTriggerConfigStruct := automation_utils_2_1.LogTriggerConfig{
		ContractAddress: emitterAddress,
		FilterSelector:  0,
		Topic0:          topic,
		Topic1:          bytes0,
		Topic2:          bytes0,
		Topic3:          bytes0,
	}
	encodedLogTriggerConfig, err := automationUtilsABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
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

// 	logTriggerConfigStruct := automation_utils_2_1.LogTriggerConfig{
// 		ContractAddress: emitterAddress,
// 		FilterSelector:  filterSelector,
// 		Topic0:          getTopic(topics, 0),
// 		Topic1:          getTopic(topics, 1),
// 		Topic2:          getTopic(topics, 2),
// 		Topic3:          getTopic(topics, 3),
// 	}
// 	encodedLogTriggerConfig, err := automationUtilsABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
// 	if err != nil {
// 		return err
// 	}

// 	err = registry.SetUpkeepTriggerConfig(upkeepID, encodedLogTriggerConfig)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// NewOrm returns a new logpoller.DbORM instance
func NewOrm(logger core_logger.SugaredLogger, chainID *big.Int, postgresDb *ctf_test_env.PostgresDb) (*logpoller.DbORM, *sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", postgresDb.ExternalPort, postgresDb.User, postgresDb.Password, postgresDb.DbName)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, db, err
	}

	db.MapperFunc(reflectx.CamelToSnakeASCII)
	return logpoller.NewORM(chainID, db, logger, pg.NewQConfig(false)), db, nil
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
func NodeHasExpectedFilters(expectedFilters []ExpectedFilter, logger core_logger.SugaredLogger, chainID *big.Int, postgresDb *ctf_test_env.PostgresDb) (bool, string, error) {
	orm, db, err := NewOrm(logger, chainID, postgresDb)
	if err != nil {
		return false, "", err
	}

	defer db.Close()
	knownFilters, err := orm.LoadFilters()
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

type LogEmitterChannel struct {
	logsEmitted int
	err         error
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

// emitEvents emits events from the provided log emitter concurrently according to the provided config
func emitEvents(ctx context.Context, l zerolog.Logger, logEmitter *contracts.LogEmitter, cfg *lp_config.Config, wg *sync.WaitGroup, results chan LogEmitterChannel) {
	address := (*logEmitter).Address().String()
	localCounter := 0
	defer wg.Done()
	for i := 0; i < *cfg.LoopedConfig.ExecutionCount; i++ {
		for _, event := range cfg.General.EventsToEmit {
			select {
			case <-ctx.Done():
				l.Warn().Str("Emitter address", address).Msg("Context cancelled, not emitting events")
				return
			default:
				l.Debug().Str("Emitter address", address).Str("Event type", event.Name).Str("index", fmt.Sprintf("%d/%d", (i+1), cfg.LoopedConfig.ExecutionCount)).Msg("Emitting log from emitter")
				var err error
				switch event.Name {
				case "Log1":
					_, err = (*logEmitter).EmitLogInts(getIntSlice(*cfg.General.EventsPerTx))
				case "Log2":
					_, err = (*logEmitter).EmitLogIntsIndexed(getIntSlice(*cfg.General.EventsPerTx))
				case "Log3":
					_, err = (*logEmitter).EmitLogStrings(getStringSlice(*cfg.General.EventsPerTx))
				case "Log4":
					_, err = (*logEmitter).EmitLogIntMultiIndexed(1, 1, *cfg.General.EventsPerTx)
				default:
					err = fmt.Errorf("unknown event name: %s", event.Name)
				}

				if err != nil {
					results <- LogEmitterChannel{
						err: err,
					}
					return
				}
				localCounter += *cfg.General.EventsPerTx

				randomWait(*cfg.LoopedConfig.MinEmitWaitTimeMs, *cfg.LoopedConfig.MaxEmitWaitTimeMs)
			}

			if (i+1)%10 == 0 {
				l.Info().Str("Emitter address", address).Str("Index", fmt.Sprintf("%d/%d", i+1, *cfg.LoopedConfig.ExecutionCount)).Msg("Emitted all three events")
			}
		}
	}

	l.Info().Str("Emitter address", address).Int("Total logs emitted", localCounter).Msg("Finished emitting events")

	results <- LogEmitterChannel{
		logsEmitted: localCounter,
		err:         nil,
	}
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
				orm, db, err := NewOrm(coreLogger, chainID, clNode.PostgresDb)
				if err != nil {
					r <- boolQueryResult{
						nodeName:     clNode.ContainerName,
						hasFinalised: false,
						err:          err,
					}
				}

				defer db.Close()

				latestBlock, err := orm.SelectLatestBlock()
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
				orm, db, err := NewOrm(coreLogger, chainID, clNode.PostgresDb)
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
					logs, err := orm.SelectLogs(startBlock, endBlock, filter.emitterAddress, filter.topic)
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
func GetMissingLogs(startBlock, endBlock int64, logEmitters []*contracts.LogEmitter, evmClient blockchain.EVMClient, clnodeCluster *test_env.ClCluster, l zerolog.Logger, coreLogger core_logger.SugaredLogger, cfg *lp_config.Config) (MissingLogs, error) {
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
				orm, db, err := NewOrm(coreLogger, evmClient.GetChainID(), clnodeCluster.Nodes[i].PostgresDb)
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
						result, err := orm.SelectLogs(startBlock, endBlock, address, event.ID)
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

	allLogsInEVMNode, err := getEVMLogs(startBlock, endBlock, logEmitters, evmClient, l, cfg)
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
				l.Info().Str("Node name", nodeName).Msg("All EVM logs were found in CL node")
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
			Msg("Some of the test logs were not found in EVM node. This is a bug in the test")
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
func getEVMLogs(startBlock, endBlock int64, logEmitters []*contracts.LogEmitter, evmClient blockchain.EVMClient, l zerolog.Logger, cfg *lp_config.Config) ([]geth_types.Log, error) {
	allLogsInEVMNode := make([]geth_types.Log, 0)
	for j := 0; j < len(logEmitters); j++ {
		address := (*logEmitters[j]).Address()
		for _, event := range cfg.General.EventsToEmit {
			l.Debug().Str("Event name", event.Name).Str("Emitter address", address.String()).Msg("Fetching logs from EVM node")
			logsInEVMNode, err := evmClient.FilterLogs(context.Background(), geth.FilterQuery{
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
func ExecuteGenerator(t *testing.T, cfg *lp_config.Config, logEmitters []*contracts.LogEmitter) (int, error) {
	if *cfg.General.Generator == lp_config.GeneratorType_WASP {
		return runWaspGenerator(t, cfg, logEmitters)
	}

	return runLoopedGenerator(t, cfg, logEmitters)
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

// runLoopedGenerator runs the looped generator and returns the total number of logs emitted
func runLoopedGenerator(t *testing.T, cfg *lp_config.Config, logEmitters []*contracts.LogEmitter) (int, error) {
	l := logging.GetTestLogger(t)

	// Start emitting events in parallel, each contract is emitting events in a separate goroutine
	// We will stop as soon as we encounter an error
	wg := &sync.WaitGroup{}
	emitterCh := make(chan LogEmitterChannel, len(logEmitters))

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	for i := 0; i < len(logEmitters); i++ {
		wg.Add(1)
		go emitEvents(ctx, l, logEmitters[i], cfg, wg, emitterCh)
	}

	var emitErr error
	total := 0

	aggrChan := make(chan int, len(logEmitters))

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case emitter := <-emitterCh:
				if emitter.err != nil {
					emitErr = emitter.err
					cancelFn()
					return
				}
				aggrChan <- emitter.logsEmitted
			}
		}
	}()

	wg.Wait()
	close(emitterCh)

	if emitErr != nil {
		return 0, emitErr
	}

	for i := 0; i < len(logEmitters); i++ {
		total += <-aggrChan
	}

	return int(total), nil
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
func chaosPauseSyncFn(l zerolog.Logger, testEnv *test_env.CLClusterTestEnv, targetComponent string) ChaosPauseData {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	randomNode := testEnv.ClCluster.Nodes[rand.Intn(len(testEnv.ClCluster.Nodes)-1)+1]
	var component ctf_test_env.EnvComponent

	switch strings.ToLower(targetComponent) {
	case "chainlink":
		component = randomNode.EnvComponent
	case "postgres":
		component = randomNode.PostgresDb.EnvComponent
	default:
		return ChaosPauseData{Err: fmt.Errorf("unknown component %s", targetComponent)}
	}

	ctx := context.Background()
	pauseStartBlock, err := testEnv.EVMClient.LatestBlockNumber(ctx)
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

	pauseEndBlock, err := testEnv.EVMClient.LatestBlockNumber(ctx)
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
func ExecuteChaosExperiment(l zerolog.Logger, testEnv *test_env.CLClusterTestEnv, cfg *lp_config.Config, errorCh chan error) {
	if cfg.ChaosConfig == nil || *cfg.ChaosConfig.ExperimentCount == 0 {
		errorCh <- nil
		return
	}

	chaosChan := make(chan ChaosPauseData, *cfg.ChaosConfig.ExperimentCount)
	wg := &sync.WaitGroup{}

	go func() {
		// if we wanted to have more than 1 container paused, we'd need to make sure we aren't trying to pause an already paused one
		guardChan := make(chan struct{}, 1)

		for i := 0; i < *cfg.ChaosConfig.ExperimentCount; i++ {
			i := i
			wg.Add(1)
			guardChan <- struct{}{}
			go func() {
				defer func() {
					<-guardChan
					wg.Done()
					current := i + 1
					l.Info().Str("Current/Total", fmt.Sprintf("%d/%d", current, cfg.ChaosConfig.ExperimentCount)).Msg("Done with experiment")
				}()
				chaosChan <- chaosPauseSyncFn(l, testEnv, *cfg.ChaosConfig.TargetComponent)
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

// GetFinalityDepth returns the finality depth for the provided chain ID
func GetFinalityDepth(chainId int64) (int64, error) {
	var finalityDepth int64
	switch chainId {
	// Ethereum Sepolia
	case 11155111:
		finalityDepth = 50
	// Polygon Mumbai
	case 80001:
		finalityDepth = 500
	// Simulated network
	case 1337:
		finalityDepth = 10
	default:
		return 0, fmt.Errorf("no known finality depth for chain %d", chainId)
	}

	return finalityDepth, nil
}

// GetEndBlockToWaitFor returns the end block to wait for based on chain id and finality tag provided in config
func GetEndBlockToWaitFor(endBlock, chainId int64, cfg *lp_config.Config) (int64, error) {
	if *cfg.General.UseFinalityTag {
		return endBlock + 1, nil
	}

	finalityDepth, err := GetFinalityDepth(chainId)
	if err != nil {
		return 0, err
	}

	return endBlock + finalityDepth, nil
}

const (
	automationDefaultUpkeepGasLimit  = uint32(2500000)
	automationDefaultLinkFunds       = int64(9e18)
	automationDefaultUpkeepsToDeploy = 10
	automationExpectedData           = "abcdef"
	defaultAmountOfUpkeeps           = 2
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

	automationDefaultRegistryConfig = contracts.KeeperRegistrySettings{
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
	lpPollingInterval time.Duration,
	backupPollingInterval uint64,
	finalityTagEnabled bool,
	testConfig *tc.TestConfig,
) (
	blockchain.EVMClient,
	[]*client.ChainlinkClient,
	contracts.ContractDeployer,
	contracts.LinkToken,
	contracts.KeeperRegistry,
	contracts.KeeperRegistrar,
	*test_env.CLClusterTestEnv,
) {
	l := logging.GetTestLogger(t)

	// Add registry version to config
	registryConfig.RegistryVersion = registryVersion
	network := networks.MustGetSelectedNetworkConfig(testConfig.Network)[0]

	finalityDepth, err := GetFinalityDepth(network.ChainID)
	require.NoError(t, err, "Error getting finality depth")

	// build the node config
	clNodeConfig := node.NewConfig(node.NewBaseConfig())
	syncInterval := *commonconfig.MustNewDuration(5 * time.Minute)
	clNodeConfig.Feature.LogPoller = ptr.Ptr[bool](true)
	clNodeConfig.OCR2.Enabled = ptr.Ptr[bool](true)
	clNodeConfig.Keeper.TurnLookBack = ptr.Ptr[int64](int64(0))
	clNodeConfig.Keeper.Registry.SyncInterval = &syncInterval
	clNodeConfig.Keeper.Registry.PerformGasOverhead = ptr.Ptr[uint32](uint32(150000))
	clNodeConfig.P2P.V2.Enabled = ptr.Ptr[bool](true)
	clNodeConfig.P2P.V2.AnnounceAddresses = &[]string{"0.0.0.0:6690"}
	clNodeConfig.P2P.V2.ListenAddresses = &[]string{"0.0.0.0:6690"}

	//launch the environment
	var env *test_env.CLClusterTestEnv
	chainlinkNodeFunding := 0.5
	l.Debug().Msgf("Funding amount: %f", chainlinkNodeFunding)
	clNodesCount := 5

	var logPolllerSettingsFn = func(chain *evmcfg.Chain) *evmcfg.Chain {
		chain.LogPollInterval = commonconfig.MustNewDuration(lpPollingInterval)
		chain.FinalityDepth = ptr.Ptr[uint32](uint32(finalityDepth))
		chain.FinalityTagEnabled = ptr.Ptr[bool](finalityTagEnabled)
		chain.BackupLogPollerBlockDelay = ptr.Ptr[uint64](backupPollingInterval)
		return chain
	}

	var evmClientSettingsFn = func(network *blockchain.EVMNetwork) *blockchain.EVMNetwork {
		network.FinalityDepth = uint64(finalityDepth)
		network.FinalityTag = finalityTagEnabled
		return network
	}

	ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
	cfg, err := ethBuilder.
		WithConsensusType(ctf_test_env.ConsensusType_PoS).
		WithConsensusLayer(ctf_test_env.ConsensusLayer_Prysm).
		WithExecutionLayer(ctf_test_env.ExecutionLayer_Geth).
		WithEthereumChainConfig(ctf_test_env.EthereumChainConfig{
			SecondsPerSlot: 4,
			SlotsPerEpoch:  2,
		}).
		Build()
	require.NoError(t, err, "Error building ethereum network config")

	env, err = test_env.NewCLTestEnvBuilder().
		WithTestConfig(testConfig).
		WithTestInstance(t).
		WithPrivateEthereumNetwork(cfg).
		WithCLNodes(clNodesCount).
		WithCLNodeConfig(clNodeConfig).
		WithFunding(big.NewFloat(chainlinkNodeFunding)).
		WithChainOptions(logPolllerSettingsFn).
		EVMClientNetworkOptions(evmClientSettingsFn).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "Error deploying test environment")

	env.ParallelTransactions(true)
	nodeClients := env.ClCluster.NodeAPIs()
	workerNodes := nodeClients[1:]

	var linkToken contracts.LinkToken

	switch network.ChainID {
	// Simulated
	case 1337:
		linkToken, err = env.ContractDeployer.DeployLinkTokenContract()
	// Ethereum Sepolia
	case 11155111:
		linkToken, err = env.ContractLoader.LoadLINKToken("0x779877A7B0D9E8603169DdbD7836e478b4624789")
	// Polygon Mumbai
	case 80001:
		linkToken, err = env.ContractLoader.LoadLINKToken("0x326C977E6efc84E512bB9C30f76E30c160eD06FB")
	default:
		panic("Not implemented")
	}
	require.NoError(t, err, "Error loading/deploying LINK token")

	linkBalance, err := env.EVMClient.BalanceAt(context.Background(), common.HexToAddress(linkToken.Address()))
	require.NoError(t, err, "Error getting LINK balance")

	l.Info().Str("Balance", big.NewInt(0).Div(linkBalance, big.NewInt(1e18)).String()).Msg("LINK balance")
	minLinkBalanceSingleNode := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(9))
	minLinkBalance := big.NewInt(0).Mul(minLinkBalanceSingleNode, big.NewInt(int64(upkeepsNeeded)))
	if minLinkBalance.Cmp(linkBalance) < 0 {
		require.FailNowf(t, "Not enough LINK", "Not enough LINK to run the test. Need at least %s", big.NewInt(0).Div(minLinkBalance, big.NewInt(1e18)).String())
	}

	registry, registrar := actions.DeployAutoOCRRegistryAndRegistrar(
		t,
		registryVersion,
		registryConfig,
		linkToken,
		env.ContractDeployer,
		env.EVMClient,
	)

	// Fund the registry with LINK
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(defaultAmountOfUpkeeps))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	err = actions.CreateOCRKeeperJobsLocal(l, nodeClients, registry.Address(), network.ChainID, 0, registryVersion)
	require.NoError(t, err, "Error creating OCR Keeper Jobs")
	ocrConfig, err := actions.BuildAutoOCR2ConfigVarsLocal(l, workerNodes, registryConfig, registrar.Address(), 30*time.Second, registry.RegistryOwnerAddress(), registry.ChainModuleAddress(), registry.ReorgProtectionEnabled())
	require.NoError(t, err, "Error building OCR config vars")
	err = registry.SetConfig(automationDefaultRegistryConfig, ocrConfig)
	require.NoError(t, err, "Registry config should be set successfully")
	require.NoError(t, env.EVMClient.WaitForEvents(), "Waiting for config to be set")

	return env.EVMClient, nodeClients, env.ContractDeployer, linkToken, registry, registrar, env
}

// UploadLogEmitterContractsAndWaitForFinalisation uploads the configured number of log emitter contracts and waits for the upload blocks to be finalised
func UploadLogEmitterContractsAndWaitForFinalisation(l zerolog.Logger, t *testing.T, testEnv *test_env.CLClusterTestEnv, cfg *lp_config.Config) []*contracts.LogEmitter {
	logEmitters := make([]*contracts.LogEmitter, 0)
	for i := 0; i < *cfg.General.Contracts; i++ {
		logEmitter, err := testEnv.ContractDeployer.DeployLogEmitterContract()
		logEmitters = append(logEmitters, &logEmitter)
		require.NoError(t, err, "Error deploying log emitter contract")
		l.Info().Str("Contract address", logEmitter.Address().Hex()).Msg("Log emitter contract deployed")
		time.Sleep(200 * time.Millisecond)
	}
	afterUploadBlock, err := testEnv.EVMClient.LatestBlockNumber(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest block number")

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		targetBlockNumber := int64(afterUploadBlock + 1)
		finalized, err := testEnv.EVMClient.GetLatestFinalizedBlockHeader(testcontext.Get(t))
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if contract were uploaded. Retrying...")
			return
		}
		finalizedBlockNumber := finalized.Number.Int64()

		if finalizedBlockNumber < targetBlockNumber {
			l.Debug().Int64("Finalized block", finalized.Number.Int64()).Int64("After upload block", int64(afterUploadBlock+1)).Msg("Waiting for contract upload to finalise")
		}

		g.Expect(finalizedBlockNumber >= targetBlockNumber).To(gomega.BeTrue(), "Contract upload did not finalize in time")
	}, "2m", "10s").Should(gomega.Succeed())

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
func FluentlyCheckIfAllNodesHaveLogCount(duration string, startBlock, endBlock int64, expectedLogCount int, expectedFilters []ExpectedFilter, l zerolog.Logger, coreLogger core_logger.SugaredLogger, testEnv *test_env.CLClusterTestEnv) (bool, error) {
	logCountWaitDuration, err := time.ParseDuration(duration)
	if err != nil {
		return false, err
	}
	endTime := time.Now().Add(logCountWaitDuration)

	// not using gomega here, because I want to see which logs were missing
	allNodesLogCountMatches := false
	for time.Now().Before(endTime) {
		logCountMatches, clErr := ClNodesHaveExpectedLogCount(startBlock, endBlock, testEnv.EVMClient.GetChainID(), expectedLogCount, expectedFilters, l, coreLogger, testEnv.ClCluster)
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
	}

	return allNodesLogCountMatches, nil
}
