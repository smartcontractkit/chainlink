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

func getExpectedFilters(logEmitters []*contracts.LogEmitter, cfg *Config) []ExpectedFilter {
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

var nodeHasExpectedFilters = func(expectedFilters []ExpectedFilter, logger core_logger.SugaredLogger, chainID *big.Int, postgresDb *ctf_test_env.PostgresDb) (bool, error) {
	orm, db, err := NewOrm(logger, chainID, postgresDb)
	if err != nil {
		return false, err
	}

	defer db.Close()
	knownFilters, err := orm.LoadFilters()
	if err != nil {
		return false, err
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
			return false, fmt.Errorf("no filter found for emitter %s and topic %s", expectedFilter.emitterAddress.String(), expectedFilter.topic.Hex())
		}
	}

	return true, nil
}

var randomWait = func(minMilliseconds, maxMilliseconds int) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomMilliseconds := rand.Intn(maxMilliseconds-minMilliseconds+1) + minMilliseconds
	time.Sleep(time.Duration(randomMilliseconds) * time.Millisecond)
}

type LogEmitterChannel struct {
	logsEmitted int
	err         error
	// unused
	// currentIndex int
}

func getIntSlice(length int) []int {
	result := make([]int, length)
	for i := 0; i < length; i++ {
		result[i] = i
	}

	return result
}

func getStringSlice(length int) []string {
	result := make([]string, length)
	for i := 0; i < length; i++ {
		result[i] = "amazing event"
	}

	return result
}

var emitEvents = func(ctx context.Context, l zerolog.Logger, logEmitter *contracts.LogEmitter, cfg *Config, wg *sync.WaitGroup, results chan LogEmitterChannel) {
	address := (*logEmitter).Address().String()
	localCounter := 0
	select {
	case <-ctx.Done():
		l.Warn().Str("Emitter address", address).Msg("Context cancelled, not emitting events")
		return
	default:
		defer wg.Done()
		for i := 0; i < cfg.LoopedConfig.ExecutionCount; i++ {
			for _, event := range cfg.General.EventsToEmit {
				l.Debug().Str("Emitter address", address).Str("Event type", event.Name).Str("index", fmt.Sprintf("%d/%d", (i+1), cfg.LoopedConfig.ExecutionCount)).Msg("Emitting log from emitter")
				var err error
				switch event.Name {
				case "Log1":
					_, err = (*logEmitter).EmitLogInts(getIntSlice(cfg.General.EventsPerTx))
				case "Log2":
					_, err = (*logEmitter).EmitLogIntsIndexed(getIntSlice(cfg.General.EventsPerTx))
				case "Log3":
					_, err = (*logEmitter).EmitLogStrings(getStringSlice(cfg.General.EventsPerTx))
				default:
					err = fmt.Errorf("unknown event name: %s", event.Name)
				}

				if err != nil {
					results <- LogEmitterChannel{
						logsEmitted: 0,
						err:         err,
					}
					return
				}
				localCounter += cfg.General.EventsPerTx

				randomWait(cfg.LoopedConfig.FuzzConfig.MinEmitWaitTimeMs, cfg.LoopedConfig.FuzzConfig.MaxEmitWaitTimeMs)
			}

			if (i+1)%10 == 0 {
				l.Info().Str("Emitter address", address).Str("Index", fmt.Sprintf("%d/%d", i+1, cfg.LoopedConfig.ExecutionCount)).Msg("Emitted all three events")
			}
		}

		l.Info().Str("Emitter address", address).Int("Total logs emitted", localCounter).Msg("Finished emitting events")

		results <- LogEmitterChannel{
			logsEmitted: localCounter,
			err:         nil,
		}
	}
}

var chainHasFinalisedEndBlock = func(l zerolog.Logger, evmClient blockchain.EVMClient, endBlock int64) (bool, error) {
	effectiveEndBlock := endBlock + 1
	lastFinalisedBlockHeader, err := evmClient.GetLatestFinalizedBlockHeader(context.Background())
	if err != nil {
		return false, err
	}

	l.Info().Int64("Last finalised block header", lastFinalisedBlockHeader.Number.Int64()).Int64("End block", effectiveEndBlock).Int64("Blocks left till end block", effectiveEndBlock-lastFinalisedBlockHeader.Number.Int64()).Msg("Waiting for the finalized block to move beyond end block")

	return lastFinalisedBlockHeader.Number.Int64() > effectiveEndBlock, nil
}

var logPollerHasFinalisedEndBlock = func(endBlock int64, chainID *big.Int, l zerolog.Logger, coreLogger core_logger.SugaredLogger, nodes *test_env.ClCluster) (bool, error) {
	wg := &sync.WaitGroup{}

	type boolQueryResult struct {
		nodeName     string
		hasFinalised bool
		err          error
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
					nodeName:     clNode.ContainerName,
					hasFinalised: latestBlock.FinalizedBlockNumber > endBlock,
					err:          nil,
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
				l.Warn().Str("Node name", r.nodeName).Msg("CL node has not finalised end block yet")
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

var clNodesHaveExpectedLogCount = func(startBlock, endBlock int64, chainID *big.Int, expectedLogCount int, expectedFilters []ExpectedFilter, l zerolog.Logger, coreLogger core_logger.SugaredLogger, nodes *test_env.ClCluster) (bool, error) {
	wg := &sync.WaitGroup{}

	type logQueryResult struct {
		nodeName         string
		logCount         int
		hasExpectedCount bool
		err              error
	}

	queryCh := make(chan logQueryResult, len(nodes.Nodes)-1)
	ctx, cancelFn := context.WithCancel(context.Background())

	for i := 1; i < len(nodes.Nodes); i++ {
		wg.Add(1)

		go func(clNode *test_env.ClNode, r chan logQueryResult) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				orm, db, err := NewOrm(coreLogger, chainID, clNode.PostgresDb)
				if err != nil {
					r <- logQueryResult{
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
						r <- logQueryResult{
							nodeName:         clNode.ContainerName,
							logCount:         0,
							hasExpectedCount: false,
							err:              err,
						}
					}

					foundLogsCount += len(logs)
				}

				r <- logQueryResult{
					nodeName:         clNode.ContainerName,
					logCount:         foundLogsCount,
					hasExpectedCount: foundLogsCount >= expectedLogCount,
					err:              err,
				}
			}
		}(nodes.Nodes[i], queryCh)
	}

	var err error
	allFoundCh := make(chan bool, 1)

	go func() {
		foundMap := make(map[string]bool, 0)
		for r := range queryCh {
			if r.err != nil {
				err = r.err
				cancelFn()
				return
			}

			foundMap[r.nodeName] = r.hasExpectedCount
			if r.hasExpectedCount {
				l.Info().Str("Node name", r.nodeName).Int("Logs count", r.logCount).Msg("Expected log count found in CL node")
			} else {
				l.Warn().Str("Node name", r.nodeName).Str("Found/Expected logs", fmt.Sprintf("%d/%d", r.logCount, expectedLogCount)).Int("Missing logs", expectedLogCount-r.logCount).Msg("Too low log count found in CL node")
			}

			if len(foundMap) == len(nodes.Nodes)-1 {
				allFound := true
				for _, v := range foundMap {
					if !v {
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
	close(queryCh)

	return <-allFoundCh, err
}

type MissingLogs map[string][]geth_types.Log

func (m *MissingLogs) IsEmpty() bool {
	for _, v := range *m {
		if len(v) > 0 {
			return false
		}
	}

	return true
}

var getMissingLogs = func(startBlock, endBlock int64, logEmitters []*contracts.LogEmitter, evmClient blockchain.EVMClient, clnodeCluster *test_env.ClCluster, l zerolog.Logger, coreLogger core_logger.SugaredLogger, cfg *Config) (MissingLogs, error) {
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

				l.Info().Str("Node name", nodeName).Msg("Fetching log poller logs")
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
						l.Debug().Str("Event name", event.Name).Str("Emitter address", address.String()).Msg("Fetching single emitter's logs")
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

						l.Debug().Str("Event name", event.Name).Str("Emitter address", address.String()).Int("Log count", len(result)).Msg("Logs found per node")
					}
				}

				l.Warn().Int("Count", len(logs)).Str("Node name", nodeName).Msg("Fetched log poller logs")

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

	l.Info().Msg("Started comparison of logs from EVM node and CL nodes. This may take a while if there's a lot of logs")
	missingCh := make(chan missingLogResult, len(clnodeCluster.Nodes)-1)
	evmLogCount := len(allLogsInEVMNode)
	for i := 1; i < len(clnodeCluster.Nodes); i++ {
		wg.Add(1)

		go func(i int, result chan missingLogResult) {
			defer wg.Done()
			nodeName := clnodeCluster.Nodes[i].ContainerName
			l.Info().Str("Node name", nodeName).Str("Progress", fmt.Sprintf("0/%d", evmLogCount)).Msg("Comparing single CL node's logs with EVM logs")

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
					l.Info().Str("Node name", nodeName).Str("Progress", fmt.Sprintf("%d/%d", i, evmLogCount)).Msg("Comparing single CL node's logs with EVM logs")
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

	expectedTotalLogsEmitted := getExpectedLogCount(cfg)
	if int64(len(allLogsInEVMNode)) != expectedTotalLogsEmitted {
		l.Warn().Str("Actual/Expected", fmt.Sprintf("%d/%d", expectedTotalLogsEmitted, len(allLogsInEVMNode))).Msg("Some of the test logs were not found in EVM node. This is a bug in the test")
	}

	return missingLogs, nil
}

var printMissingLogsByType = func(missingLogs map[string][]geth_types.Log, l zerolog.Logger, cfg *Config) {
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

	for k, v := range missingByType {
		l.Warn().Str("Event name", k).Int("Missing count", v).Msg("Missing logs by type")
	}
}

var getEVMLogs = func(startBlock, endBlock int64, logEmitters []*contracts.LogEmitter, evmClient blockchain.EVMClient, l zerolog.Logger, cfg *Config) ([]geth_types.Log, error) {
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

	l.Warn().Int("Count", len(allLogsInEVMNode)).Msg("Logs in EVM node")

	return allLogsInEVMNode, nil
}

func executeGenerator(t *testing.T, cfg *Config, logEmitters []*contracts.LogEmitter) (int, error) {
	if cfg.General.Generator == GeneratorType_WASP {
		return runWaspGenerator(t, cfg, logEmitters)
	}

	return runLoopedGenerator(t, cfg, logEmitters)
}

func runWaspGenerator(t *testing.T, cfg *Config, logEmitters []*contracts.LogEmitter) (int, error) {
	l := logging.GetTestLogger(t)

	var RPSprime int64

	// if LPS is set, we need to calculate based on countract count and events per transaction
	if cfg.Wasp.Load.LPS > 0 {
		RPSprime = cfg.Wasp.Load.LPS / int64(cfg.General.Contracts) / int64(cfg.General.EventsPerTx) / int64(len(cfg.General.EventsToEmit))

		if RPSprime < 1 {
			return 0, fmt.Errorf("invalid load configuration, effective RPS would have been zero. Adjust LPS, contracts count, events per tx or events to emit")
		}
	}

	// if RPS is set simply split it between contracts
	if cfg.Wasp.Load.RPS > 0 {
		RPSprime = cfg.Wasp.Load.RPS / int64(cfg.General.Contracts)
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
			RateLimitUnitDuration: cfg.Wasp.Load.RateLimitUnitDuration.Duration(),
			CallTimeout:           cfg.Wasp.Load.CallTimeout.Duration(),
			Schedule: wasp.Plain(
				RPSprime,
				cfg.Wasp.Load.Duration.Duration(),
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

func runLoopedGenerator(t *testing.T, cfg *Config, logEmitters []*contracts.LogEmitter) (int, error) {
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
		for emitter := range emitterCh {
			if emitter.err != nil {
				emitErr = emitter.err
				cancelFn()
				return
			}
			aggrChan <- emitter.logsEmitted
		}
	}()

	wg.Wait()
	close(emitterCh)

	for i := 0; i < len(logEmitters); i++ {
		total += <-aggrChan
	}

	if emitErr != nil {
		return 0, emitErr
	}

	return int(total), nil
}

func getExpectedLogCount(cfg *Config) int64 {
	if cfg.General.Generator == GeneratorType_WASP {
		if cfg.Wasp.Load.RPS != 0 {
			return cfg.Wasp.Load.RPS * int64(cfg.Wasp.Load.Duration.Duration().Seconds()) * int64(cfg.General.EventsPerTx)
		}
		return cfg.Wasp.Load.LPS * int64(cfg.Wasp.Load.Duration.Duration().Seconds())
	}

	return int64(len(cfg.General.EventsToEmit) * cfg.LoopedConfig.ExecutionCount * cfg.General.Contracts * cfg.General.EventsPerTx)
}

var chaosPauseSyncFn = func(l zerolog.Logger, testEnv *test_env.CLClusterTestEnv) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomBool := rand.Intn(2) == 0

	randomNode := testEnv.ClCluster.Nodes[rand.Intn(len(testEnv.ClCluster.Nodes)-1)+1]
	var component ctf_test_env.EnvComponent

	if randomBool {
		component = randomNode.EnvComponent
	} else {
		component = randomNode.PostgresDb.EnvComponent
	}

	pauseTimeSec := rand.Intn(20-5) + 5
	l.Info().Str("Container", component.ContainerName).Int("Pause time", pauseTimeSec).Msg("Pausing component")
	pauseTimeDur := time.Duration(pauseTimeSec) * time.Second
	err := component.ChaosPause(l, pauseTimeDur)
	l.Info().Str("Container", component.ContainerName).Msg("Component unpaused")

	if err != nil {
		return err
	}

	return nil
}

var executeChaosExperiment = func(l zerolog.Logger, testEnv *test_env.CLClusterTestEnv, cfg *Config, errorCh chan error) {
	if cfg.ChaosConfig == nil || cfg.ChaosConfig.ExperimentCount == 0 {
		errorCh <- nil
		return
	}

	chaosChan := make(chan error, cfg.ChaosConfig.ExperimentCount)

	wg := &sync.WaitGroup{}

	go func() {
		// if we wanted to have more than 1 container paused, we'd need to make sure we aren't trying to pause an already paused one
		guardChan := make(chan struct{}, 1)

		for i := 0; i < cfg.ChaosConfig.ExperimentCount; i++ {
			i := i
			wg.Add(1)
			guardChan <- struct{}{}
			go func() {
				defer func() {
					<-guardChan
					wg.Done()
					l.Info().Str("Current/Total", fmt.Sprintf("%d/%d", i, cfg.ChaosConfig.ExperimentCount)).Msg("Done with experiment")
				}()
				chaosChan <- chaosPauseSyncFn(l, testEnv)
			}()
		}

		wg.Wait()

		close(chaosChan)
	}()

	go func() {
		for err := range chaosChan {
			// This will receive errors until chaosChan is closed
			if err != nil {
				// If an error is encountered, log it, send it to the error channel, and return from the function
				l.Err(err).Msg("Error encountered during chaos experiment")
				errorCh <- err
				return // Return on actual error
			}
			// No need for an else block here, because if err is nil (which happens when the channel is closed),
			// the loop will exit and the following log and nil send will execute.
		}

		// After the loop exits, which it will do when chaosChan is closed, log that all experiments are finished.
		l.Info().Msg("All chaos experiments finished")
		errorCh <- nil // Only send nil once, after all errors have been handled and the channel is closed
	}()
}

var GetFinalityDepth = func(chainId int64) (int64, error) {
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

var GetEndBlockToWaitFor = func(endBlock, chainId int64, cfg *Config) (int64, error) {
	if cfg.General.UseFinalityTag {
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
	defaultOCRRegistryConfig = contracts.KeeperRegistrySettings{
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

func setupLogPollerTestDocker(
	t *testing.T,
	registryVersion ethereum.KeeperRegistryVersion,
	registryConfig contracts.KeeperRegistrySettings,
	upkeepsNeeded int,
	lpPollingInterval time.Duration,
	finalityTagEnabled bool,
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
	network := networks.MustGetSelectedNetworksFromEnv()[0]

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
		WithWaitingForFinalization().
		WithEthereumChainConfig(ctf_test_env.EthereumChainConfig{
			SecondsPerSlot: 8,
			SlotsPerEpoch:  2,
		}).
		Build()
	require.NoError(t, err, "Error building ethereum network config")

	env, err = test_env.NewCLTestEnvBuilder().
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
	ocrConfig, err := actions.BuildAutoOCR2ConfigVarsLocal(l, workerNodes, registryConfig, registrar.Address(), 30*time.Second, registry.RegistryOwnerAddress())
	require.NoError(t, err, "Error building OCR config vars")
	err = registry.SetConfig(automationDefaultRegistryConfig, ocrConfig)
	require.NoError(t, err, "Registry config should be set successfully")
	require.NoError(t, env.EVMClient.WaitForEvents(), "Waiting for config to be set")

	return env.EVMClient, nodeClients, env.ContractDeployer, linkToken, registry, registrar, env
}
