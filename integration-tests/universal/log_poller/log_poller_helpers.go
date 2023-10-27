package logpoller

import (
	"bytes"
	"context"
	"errors"
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
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_blockchain "github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	cltypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/stretchr/testify/require"

	"github.com/scylladb/go-reflectx"
	it_utils "github.com/smartcontractkit/chainlink/integration-tests/utils"
	utils2 "github.com/smartcontractkit/chainlink/integration-tests/utils"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	lpEvm "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	le "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	core_logger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

func ExecuteBasicLogPollerTest(t *testing.T, cfg *Config) {
	l := logging.GetTestLogger(t)
	coreLogger := core_logger.TestLogger(t) //needed by ORM ¯\_(ツ)_/¯

	if cfg.General.EventsToEmit == nil || len(cfg.General.EventsToEmit) == 0 {
		l.Warn().Msg("No events to emit specified, using all events from log emitter contract")
		for _, event := range EmitterABI.Events {
			cfg.General.EventsToEmit = append(cfg.General.EventsToEmit, event)
		}
	}

	l.Info().Msg("Starting basic log poller test")

	var (
		err      error
		testName = "basic-log-poller"
	)

	chainClient, _, contractDeployer, linkToken, registry, registrar, testEnv := setupLogPollerTestDocker(
		t, testName, ethereum.RegistryVersion_2_1, defaultOCRRegistryConfig, false, time.Duration(500*time.Millisecond), 500, 10, false,
	)

	upKeepsNeeded := cfg.General.Contracts * len(cfg.General.EventsToEmit)
	_, upkeepIDs := actions.DeployConsumers(
		t,
		registry,
		registrar,
		linkToken,
		contractDeployer,
		chainClient,
		upKeepsNeeded,
		big.NewInt(automationDefaultLinkFunds),
		automationDefaultUpkeepGasLimit,
		true,
		false,
	)

	// Deploy Log Emitter contracts
	logEmitters := make([]*contracts.LogEmitter, 0)
	for i := 0; i < cfg.General.Contracts; i++ {
		logEmitter, err := testEnv.ContractDeployer.DeployLogEmitterContract()
		logEmitters = append(logEmitters, &logEmitter)
		require.NoError(t, err, "Error deploying log emitter contract")
		l.Info().Str("Contract address", logEmitter.Address().Hex()).Msg("Log emitter contract deployed")
	}

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	for i := 0; i < len(upkeepIDs); i++ {
		emitterAddress := (*logEmitters[i%cfg.General.Contracts]).Address()
		upkeepID := upkeepIDs[i]
		topicId := cfg.General.EventsToEmit[i%len(cfg.General.EventsToEmit)].ID

		l.Info().Int("Upkeep id", int(upkeepID.Int64())).Str("Emitter address", emitterAddress.String()).Str("Topic", topicId.Hex()).Msg("Registering log trigger for log emitter")
		err = registerSingleTopicFilter(registry, upkeepID, emitterAddress, topicId)
		require.NoError(t, err, "Error registering log trigger for log emitter")
	}

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error encountered when waiting for setting trigger config for upkeeps")

	// Make sure that all nodes have expected filters registered before starting to emit events
	expectedFilters := getExpectedFilters(logEmitters, cfg)
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
			nodeName := testEnv.ClCluster.Nodes[i].ContainerName
			l.Info().Str("Node name", nodeName).Msg("Fetching filters from log poller's DB")
			orm, err := NewOrm(coreLogger, testEnv.EVMClient.GetChainID(), testEnv.ClCluster.Nodes[i].PostgresDb)
			require.NoError(t, err, "Error creating ORM")

			hasFilters, err := nodeHasExpectedFilters(expectedFilters, orm)
			if err != nil {
				l.Warn().Err(err).Msg("Error checking if node has expected filters. Retrying...")
				return
			}

			g.Expect(hasFilters).To(gomega.BeTrue(), "Not all expected filters were found in the DB")
		}
	}, "30s", "1s").Should(gomega.Succeed())
	l.Info().Msg("All nodes have expected filters registered")

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := testEnv.EVMClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")
	startBlock := int64(sb)

	l.Info().Msg("Starting event emission")
	totalLogsEmitted, err := executeGenerator(t, cfg, logEmitters)
	require.NoError(t, err, "Error executing event generator")
	expectedLogsEmitted := getExpectedLogCount(cfg)
	l.Info().Int("Total logs emitted", totalLogsEmitted).Int64("Expected total logs emitted", expectedLogsEmitted).Msg("Finished emitting events")
	// l.Info().Int("Actual total logs emitted", totalLogsEmitted).Int("Expected total logs emitted", len(logEmitters)*len(cfg.General.EventsToEmit)*cfg.LoopedConfig.ExecutionCount).Msg("Finished emitting events")

	// Save block number after finishing to emit events, so that we can later use it when querying logs
	eb, err := testEnv.EVMClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")
	// +10 to be safe, but this should be done fluently, but in such a way that can handle reorgs, so making sure that each trx is included in finalised block might not suffice
	// as some trx might be included in a block that was pruned
	endBlock := int64(eb) + 10

	waitDuration := "30s"
	l.Warn().Str("Duration", waitDuration).Msg("Waiting for logs to be processed by all nodes and for chain to advance beyond finality")

	// Wait until last block in which events were emitted has been finalised
	// how long should we wait here until all logs are processed? wait for block X to be processed by all nodes?
	gom.Eventually(func(g gomega.Gomega) {
		hasAdvanced, err := chainHasAdvancedBeyondFinality(testEnv.EVMClient, endBlock)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if chain has advanced beyond finality. Retrying...")
		}
		g.Expect(hasAdvanced).To(gomega.BeTrue(), "Chain has not advanced beyond finality")
	}, waitDuration, "1s").Should(gomega.Succeed())

	// Wait until all CL nodes have exactly the same logs emitted by test contracts as the EVM node has
	logConsistencyWaitDuration := "1m"
	l.Warn().Str("Duration", logConsistencyWaitDuration).Msg("Waiting for CL nodes to have all the logs that EVM node has")

	gom.Eventually(func(g gomega.Gomega) {
		missingLogs, err := getMissingLogs(startBlock, endBlock, logEmitters, testEnv.EVMClient, testEnv.ClCluster, l, coreLogger, cfg)
		if err != nil {
			l.Warn().Err(err).Msg("Error getting missing logs. Retrying...")
		}

		if !missingLogs.IsEmpty() {
			printMissingLogsByType(missingLogs, l, cfg)
		}
		g.Expect(missingLogs.IsEmpty()).To(gomega.BeTrue(), "Some CL nodes were missing logs")
	}, logConsistencyWaitDuration, "1s").Should(gomega.Succeed())
}

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

// this is not really possible, log trigger doesn't support multiple topics, even if log poller does
var registerMultipleTopicsFilter = func(registry contracts.KeeperRegistry, upkeepID *big.Int, emitterAddress common.Address, topics []abi.Event) error {
	if len(topics) > 4 {
		return errors.New("Cannot register more than 4 topics")
	}

	var getTopic = func(topics []abi.Event, i int) common.Hash {
		if i > len(topics)-1 {
			return bytes0
		}

		return topics[i].ID
	}

	var getFilterSelector = func(topics []abi.Event) (uint8, error) {
		switch len(topics) {
		case 0:
			return 0, errors.New("Cannot register filter with 0 topics")
		case 1:
			return 0, nil
		case 2:
			return 1, nil
		case 3:
			return 3, nil
		case 4:
			return 7, nil
		default:
			return 0, errors.New("Cannot register filter with more than 4 topics")
		}
	}

	filterSelector, err := getFilterSelector(topics)
	if err != nil {
		return err
	}

	logTriggerConfigStruct := automation_utils_2_1.LogTriggerConfig{
		ContractAddress: emitterAddress,
		FilterSelector:  filterSelector,
		Topic0:          getTopic(topics, 0),
		Topic1:          getTopic(topics, 1),
		Topic2:          getTopic(topics, 2),
		Topic3:          getTopic(topics, 3),
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

func NewOrm(logger core_logger.SugaredLogger, chainID *big.Int, postgresDb *ctf_test_env.PostgresDb) (*lpEvm.DbORM, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", postgresDb.ExternalPort, "postgres", "mysecretpassword", "testdb")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.MapperFunc(reflectx.CamelToSnakeASCII)
	return lpEvm.NewORM(chainID, db, logger, pg.NewQConfig(false)), nil
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

var nodeHasExpectedFilters = func(expectedFilters []ExpectedFilter, orm *lpEvm.DbORM) (bool, error) {
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
			return false, fmt.Errorf("No filter found for emitter %s and topic %s", expectedFilter.emitterAddress.String(), expectedFilter.topic.Hex())
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
	logsEmitted  int
	err          error
	currentIndex int
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
	defer wg.Done()
	address := (*logEmitter).Address().String()
	localCounter := 0
	select {
	case <-ctx.Done():
		return
	default:
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
					err = fmt.Errorf("Unknown event name: %s", event.Name)
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

		results <- LogEmitterChannel{
			logsEmitted: localCounter,
			err:         nil,
		}
	}
}

var chainHasAdvancedBeyondFinality = func(evmClint ctf_blockchain.EVMClient, endBlock int64) (bool, error) {
	lastFinalisedBlockHeader, err := evmClint.GetLatestFinalizedBlockHeader(context.Background())
	if err != nil {
		return false, err
	}

	return lastFinalisedBlockHeader.Number.Int64() > endBlock+1, nil
}

// TODO this fails, because not all logs emited by contracts are even in the EVM node
var nodeDbHasExpectedLogCount = func(startBlock, endBlock int64, chainID *big.Int, expectedLogCount int, expectedFilters []ExpectedFilter, logger core_logger.SugaredLogger, postgresDB *ctf_test_env.PostgresDb) (bool, error) {
	orm, err := NewOrm(logger, chainID, postgresDB)
	if err != nil {
		return false, err
	}

	foundLogsCount := 0

	for _, filter := range expectedFilters {
		logs, err := orm.SelectLogs(startBlock, endBlock, filter.emitterAddress, filter.topic)
		if err != nil {
			return false, err
		}

		foundLogsCount += len(logs)
	}

	return foundLogsCount == expectedLogCount, nil
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

var getMissingLogs = func(startBlock, endBlock int64, logEmitters []*contracts.LogEmitter, evmClient ctf_blockchain.EVMClient, clnodeCluster *test_env.ClCluster, l zerolog.Logger, coreLogger core_logger.SugaredLogger, cfg *Config) (MissingLogs, error) {
	wg := sync.WaitGroup{}

	type dbQueryResult struct {
		err      error
		nodeName string
		logs     []logpoller.Log
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	resultCh := make(chan dbQueryResult, len(clnodeCluster.Nodes))

	for i := 1; i < len(clnodeCluster.Nodes); i++ {
		wg.Add(1)

		go func(ctx context.Context, i int, r chan dbQueryResult) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				nodeName := clnodeCluster.Nodes[i].ContainerName

				l.Info().Str("Node name", nodeName).Msg("Start fetching logs from log poller's DB")
				orm, err := NewOrm(coreLogger, evmClient.GetChainID(), clnodeCluster.Nodes[i].PostgresDb)
				if err != nil {
					r <- dbQueryResult{
						err:      err,
						nodeName: nodeName,
						logs:     []logpoller.Log{},
					}
				}

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

				l.Warn().Int("Total per node", len(logs)).Str("Node name", nodeName).Msg("Total logs per node")

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
				dbError = r.err
				cancelFn()
				return
			}
			allLogPollerLogs[r.nodeName] = r.logs
		}
	}()

	wg.Wait()
	close(resultCh)

	if dbError != nil {
		return nil, dbError
	}

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

	for i := 1; i < len(clnodeCluster.Nodes); i++ {
		nodeName := clnodeCluster.Nodes[i].ContainerName
		l.Info().Str("Node name", nodeName).Int("Log count", len(allLogPollerLogs[nodeName])).Msg("CL node log count")

		missingLogs[nodeName] = make([]geth_types.Log, 0)
		for _, evmLog := range allLogsInEVMNode {
			logFound := false
			for _, logPollerLog := range allLogPollerLogs[nodeName] {
				if logPollerLog.BlockNumber == int64(evmLog.BlockNumber) && logPollerLog.TxHash == evmLog.TxHash && bytes.Equal(logPollerLog.Data, evmLog.Data) && logPollerLog.LogIndex == int64(evmLog.Index) &&
					logPollerLog.Address == evmLog.Address && logPollerLog.BlockHash == evmLog.BlockHash && bytes.Equal(logPollerLog.Topics[0][:], evmLog.Topics[0].Bytes()) {
					logFound = true
					continue
				}
			}

			if !logFound {
				missingLogs[nodeName] = append(missingLogs[nodeName], evmLog)
			}
		}
	}

	expectedTotalLogsEmitted := getExpectedLogCount(cfg)
	if int64(len(allLogsInEVMNode)) != expectedTotalLogsEmitted {
		l.Warn().Int64("Expected", expectedTotalLogsEmitted).Int("Actual", len(allLogsInEVMNode)).Msg("Total logs emitted by contracts found in EVM node")
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
			if _, ok := missingByType[humanName]; ok {
				missingByType[humanName] += 1
			} else {
				missingByType[humanName] = 1
			}
		}
	}

	for k, v := range missingByType {
		l.Warn().Str("Event name", k).Int("Missing count", v).Msg("Missing logs by type")
	}
}

func executeGenerator(t *testing.T, cfg *Config, logEmitters []*contracts.LogEmitter) (int, error) {
	if cfg.General.Generator == GeneratorType_WASP {
		return runWaspGenerator(t, cfg, logEmitters)
	}

	return runLoopedGenerator(t, cfg, logEmitters)
}

func runWaspGenerator(t *testing.T, cfg *Config, logEmitters []*contracts.LogEmitter) (int, error) {
	l := logging.GetTestLogger(t)

	RPSprime := cfg.Wasp.Load.RPS / int64(cfg.General.Contracts) / int64(len(cfg.General.EventsToEmit))

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
		})
		p.Add(g, err)
	}

	_, err := p.Run(true)

	if err != nil {
		return 0, err
	}

	total := 0
	for _, g := range p.Generators {
		if v, ok := g.GetData().OKData.Data[0].(int); ok {
			total += v
		}
	}

	return total, nil
}

func runLoopedGenerator(t *testing.T, cfg *Config, logEmitters []*contracts.LogEmitter) (int, error) {
	l := logging.GetTestLogger(t)

	// Start emitting events in parallel, each contract is emitting events in a separate goroutine
	// We will stop as soon as we encounter an error
	var wg sync.WaitGroup
	emitterCh := make(chan LogEmitterChannel, len(logEmitters))

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	for i := 0; i < len(logEmitters); i++ {
		wg.Add(1)
		go emitEvents(ctx, l, logEmitters[i], cfg, &wg, emitterCh)
	}

	var emitErr error
	totalLogsEmitted := 0

	go func() {
		for emitter := range emitterCh {
			if emitter.err != nil {
				emitErr = emitter.err
				cancelFn()
				return
			}
			totalLogsEmitted += emitter.logsEmitted
		}
	}()

	wg.Wait()

	close(emitterCh)

	if emitErr != nil {
		return 0, emitErr
	}

	return totalLogsEmitted, nil
}

func getExpectedLogCount(cfg *Config) int64 {
	if cfg.General.Generator == GeneratorType_WASP {
		return cfg.Wasp.Load.RPS * int64(cfg.Wasp.Load.Duration.Duration().Seconds()) * int64(cfg.General.EventsPerTx)
	}

	return int64(len(cfg.General.EventsToEmit) * cfg.LoopedConfig.ExecutionCount * cfg.General.Contracts * cfg.General.EventsPerTx)
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
	testName string,
	registryVersion ethereum.KeeperRegistryVersion,
	registryConfig contracts.KeeperRegistrySettings,
	statefulDb bool,
	lpPollingInterval time.Duration,
	blockBackfillDepth uint32,
	finalityDepth uint32,
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
	require.False(t, finalityTagEnabled && finalityDepth > 0, "Cannot use finality tag and finality depth at the same time")
	l := logging.GetTestLogger(t)
	// Add registry version to config
	registryConfig.RegistryVersion = registryVersion
	network := networks.SelectedNetwork

	// build the node config
	clNodeConfig := node.NewConfig(node.NewBaseConfig())
	syncInterval := models.MustMakeDuration(5 * time.Minute)
	clNodeConfig.Feature.LogPoller = it_utils.Ptr[bool](true)
	clNodeConfig.OCR2.Enabled = it_utils.Ptr[bool](true)
	clNodeConfig.Keeper.TurnLookBack = it_utils.Ptr[int64](int64(0))
	clNodeConfig.Keeper.Registry.SyncInterval = &syncInterval
	clNodeConfig.Keeper.Registry.PerformGasOverhead = it_utils.Ptr[uint32](uint32(150000))
	clNodeConfig.P2P.V2.Enabled = it_utils.Ptr[bool](true)
	clNodeConfig.P2P.V2.AnnounceAddresses = &[]string{"0.0.0.0:6690"}
	clNodeConfig.P2P.V2.ListenAddresses = &[]string{"0.0.0.0:6690"}

	//launch the environment
	var env *test_env.CLClusterTestEnv
	var err error
	chainlinkNodeFunding := 1.0
	l.Debug().Msgf("Funding amount: %f", chainlinkNodeFunding)
	clNodesCount := 5

	var logPolllerSettingsFn = func(chain *evmcfg.Chain) *evmcfg.Chain {
		chain.LogPollInterval = models.MustNewDuration(lpPollingInterval)
		chain.BlockBackfillDepth = utils2.Ptr[uint32](blockBackfillDepth)
		chain.FinalityDepth = utils2.Ptr[uint32](finalityDepth)
		chain.FinalityTagEnabled = utils2.Ptr[bool](finalityTagEnabled)
		return chain
	}

	var evmClientSettingsFn = func(network *blockchain.EVMNetwork) *blockchain.EVMNetwork {
		network.FinalityDepth = uint64(finalityDepth)
		network.FinalityTag = finalityTagEnabled
		return network
	}

	env, err = test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithGeth().
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

	linkToken, err := env.ContractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Error deploying LINK token")

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
