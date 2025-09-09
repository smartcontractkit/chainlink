package workflow

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/scylladb/go-reflectx"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	wf_reg_v2_op "github.com/smartcontractkit/chainlink/deployment/cre/workflow_registry/v2/changeset/operations/contracts"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

// must match nubmer of events we track in core/services/workflows/syncer/handler.go
const NumberOfTrackedWorkflowRegistryEvents = 6

func WaitForWorkflowRegistryFiltersRegistration(
	testLogger zerolog.Logger,
	singleFileLogger logger.Logger,
	infraType infra.Type,
	registryChainID uint64,
	topology *cre.DonTopology,
	nodeSetInput []*cre.CapabilitiesAwareNodeSet,
) error {
	// we currently have no way of checking if filters were registered, when code runs in CRIB
	// as we don't have a way to get its database connection string
	if infraType == infra.CRIB {
		return nil
	}

	return waitForAllNodesToHaveExpectedFiltersRegistered(singleFileLogger, testLogger, registryChainID, topology, nodeSetInput)
}

func ConfigureWorkflowRegistry(
	ctx context.Context,
	testLogger zerolog.Logger,
	singleFileLogger logger.Logger,
	input *cre.WorkflowRegistryInput,
) (*cre.WorkflowRegistryOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if input.Out != nil && input.Out.UseCache {
		return input.Out, nil
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	// we need to filter out all chains from the environment struct
	// that do not have at least one contract deployed to avoid validation errors
	// when configuring workflow registry contract :shrug: :shrug: :shrug:
	allAddresses, addrErr := input.CldEnv.ExistingAddresses.Addresses() //nolint:staticcheck // ignore SA1019 as ExistingAddresses is deprecated but still used
	if addrErr != nil {
		return nil, errors.Wrap(addrErr, "failed to get addresses from address book")
	}

	chainsWithContracts := make(map[uint64]bool)
	for chainSelector, addresses := range allAddresses {
		chainsWithContracts[chainSelector] = len(addresses) > 0
	}

	addresses, addrErr1 := input.CldEnv.DataStore.Addresses().Fetch()
	if addrErr1 != nil {
		return nil, errors.Wrap(addrErr1, "failed to get addresses from datastore")
	}

	for _, addr := range addresses {
		chainsWithContracts[addr.ChainSelector] = true
	}

	nonEmptyBlockchains := make(map[uint64]cldf_chain.BlockChain, 0)
	for chainSelector, chain := range input.CldEnv.BlockChains.EVMChains() {
		if chainsWithContracts[chain.Selector] {
			nonEmptyBlockchains[chainSelector] = chain
		}
	}
	for chainSelector, chain := range input.CldEnv.BlockChains.SolanaChains() {
		if chainsWithContracts[chain.Selector] {
			nonEmptyBlockchains[chainSelector] = chain
		}
	}

	nonEmptyChainsCLDEnvironment := &cldf.Environment{
		Logger:            singleFileLogger,
		ExistingAddresses: input.CldEnv.ExistingAddresses, //nolint:staticcheck // ignore SA1019 as ExistingAddresses is deprecated but still used
		GetContext: func() context.Context {
			return ctx
		},
		DataStore:   input.CldEnv.DataStore,
		BlockChains: cldf_chain.NewBlockChains(nonEmptyBlockchains),
	}
	nonEmptyChainsCLDEnvironment.OperationsBundle = operations.NewBundle(nonEmptyChainsCLDEnvironment.GetContext, singleFileLogger, operations.NewMemoryReporter())

	allowedDonIDs := make([]uint32, len(input.AllowedDonIDs))
	for i, donID := range input.AllowedDonIDs {
		allowedDonIDs[i] = libc.MustSafeUint32FromUint64(donID)
	}

	if input.ContractVersion.Version.Equal(semver.MustParse(config.WorkflowRegistryV2Semver)) {
		updateSignersReport, err := operations.ExecuteOperation(
			input.CldEnv.OperationsBundle,
			wf_reg_v2_op.UpdateAllowedSignersOp,
			wf_reg_v2_op.WorkflowRegistryOpDeps{
				Env: nonEmptyChainsCLDEnvironment,
			},
			wf_reg_v2_op.UpdateAllowedSignersOpInput{
				ChainSelector: input.ChainSelector,
				Signers:       input.WorkflowOwners,
				Allowed:       true,
			},
		)
		if err != nil || !updateSignersReport.Output.Success {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to update allowed signers on workflow registry %s", input.ContractVersion.Version))
		}

		donLimitReport, err := operations.ExecuteOperation(
			input.CldEnv.OperationsBundle,
			wf_reg_v2_op.SetDONLimitOp,
			wf_reg_v2_op.WorkflowRegistryOpDeps{
				Env: nonEmptyChainsCLDEnvironment,
			},
			wf_reg_v2_op.SetDONLimitOpInput{
				ChainSelector: input.ChainSelector,
				DONFamily:     config.DefaultDONFamily,
				Limit:         libc.MustSafeUint32(len(allowedDonIDs)),
				Enabled:       true,
			},
		)
		if err != nil || !donLimitReport.Output.Success {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to set DON Limit on workflow registry %s", input.ContractVersion.Version))
		}

		return &cre.WorkflowRegistryOutput{
			ChainSelector:  input.ChainSelector,
			AllowedDonIDs:  allowedDonIDs,
			WorkflowOwners: input.WorkflowOwners,
		}, nil
	}

	report, err := operations.ExecuteSequence(
		input.CldEnv.OperationsBundle,
		ks_contracts_op.ConfigWorkflowRegistrySeq,
		ks_contracts_op.ConfigWorkflowRegistrySeqDeps{
			Env: nonEmptyChainsCLDEnvironment,
		},
		ks_contracts_op.ConfigWorkflowRegistrySeqInput{
			ContractAddress:       input.ContractAddress,
			RegistryChainSelector: input.ChainSelector,
			AllowedDonIDs:         allowedDonIDs,
			WorkflowOwners:        input.WorkflowOwners,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure workflow registry")
	}

	input.Out = &cre.WorkflowRegistryOutput{
		ChainSelector:  report.Output.RegistryChainSelector,
		AllowedDonIDs:  report.Output.AllowedDonIDs,
		WorkflowOwners: report.Output.WorkflowOwners,
	}
	return input.Out, nil
}

// waitForAllNodesToHaveExpectedFiltersRegistered manually checks if all WorkflowRegistry filters used by the LogPoller are registered for all nodes. We want to see if this will help with the flakiness.
func waitForAllNodesToHaveExpectedFiltersRegistered(singeFileLogger logger.Logger, testLogger zerolog.Logger, homeChainID uint64, donTopology *cre.DonTopology, nodeSetInput []*cre.CapabilitiesAwareNodeSet) error {
	for donIdx, don := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(don.Flags, cre.WorkflowDON) {
			continue
		}

		workderNodes, workersErr := node.FindManyWithLabel(don.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
		if workersErr != nil {
			return errors.Wrap(workersErr, "failed to find worker nodes")
		}

		results := make(map[int]bool)
		ticker := 5 * time.Second
		timeout := 2 * time.Minute

	INNER_LOOP:
		for {
			select {
			case <-time.After(timeout):
				return fmt.Errorf("timed out, when waiting for %.2f seconds, waiting for all nodes to have expected filters registered", timeout.Seconds())
			case <-time.Tick(ticker):
				if len(results) == len(workderNodes) {
					testLogger.Info().Msgf("All %d nodes in DON %d have expected filters registered", len(workderNodes), don.ID)
					break INNER_LOOP
				}

				for _, workerNode := range workderNodes {
					nodeIndex, nodeIndexErr := node.FindLabelValue(workerNode, node.IndexKey)
					if nodeIndexErr != nil {
						return errors.Wrap(nodeIndexErr, "failed to find node index")
					}

					nodeIndexInt, nodeIdxErr := strconv.Atoi(nodeIndex)
					if nodeIdxErr != nil {
						return errors.Wrap(nodeIdxErr, "failed to convert node index to int")
					}

					if _, ok := results[nodeIndexInt]; ok {
						continue
					}

					testLogger.Info().Msgf("Checking if all WorkflowRegistry filters are registered for worker node %d", nodeIndexInt)
					allFilters, filtersErr := getAllFilters(context.Background(), singeFileLogger, big.NewInt(libc.MustSafeInt64(homeChainID)), nodeIndexInt, nodeSetInput[donIdx].DbInput.Port)
					if filtersErr != nil {
						return errors.Wrap(filtersErr, "failed to get filters")
					}

					for _, filter := range allFilters {
						if strings.Contains(filter.Name, "WorkflowRegistry") {
							if len(filter.EventSigs) == NumberOfTrackedWorkflowRegistryEvents {
								testLogger.Debug().Msgf("Found all WorkflowRegistry filters for node %d", nodeIndexInt)
								results[nodeIndexInt] = true
								continue
							}

							testLogger.Debug().Msgf("Found only %d WorkflowRegistry filters for node %d", len(filter.EventSigs), nodeIndexInt)
						}
					}
				}

				// return if we have results for all nodes, don't wait for next tick
				if len(results) == len(workderNodes) {
					testLogger.Info().Msgf("All %d nodes in DON %d have expected filters registered", len(workderNodes), don.ID)
					break INNER_LOOP
				}
			}
		}
	}

	return nil
}

// StartS3 starts MiniIO as S3 Provider, if input is not nil. It's purpose is to store workflow-related artifacts.
func StartS3(testLogger zerolog.Logger, input *s3provider.Input, stageGen *stagegen.StageGen) (*s3provider.Output, error) {
	var s3ProviderOutput *s3provider.Output
	if input != nil {
		fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting MinIO")))
		var s3ProviderErr error
		s3ProviderOutput, s3ProviderErr = s3provider.NewMinioFactory().NewFrom(input)
		if s3ProviderErr != nil {
			return nil, errors.Wrap(s3ProviderErr, "minio provider creation failed")
		}
		testLogger.Debug().Msgf("S3Provider.Output value: %#v", s3ProviderOutput)
		fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("MinIO started in %.2f seconds", stageGen.Elapsed().Seconds())))
	}

	return s3ProviderOutput, nil
}

func newORM(logger logger.Logger, chainID *big.Int, nodeIndex, externalPort int) (logpoller.ORM, *sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", externalPort, postgres.User, postgres.Password, fmt.Sprintf("db_%d", nodeIndex))
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, db, err
	}

	db.MapperFunc(reflectx.CamelToSnakeASCII)
	return logpoller.NewORM(chainID, db, logger), db, nil
}

func getAllFilters(ctx context.Context, logger logger.Logger, chainID *big.Int, nodeIndex, externalPort int) (map[string]logpoller.Filter, error) {
	orm, db, err := newORM(logger, chainID, nodeIndex, externalPort)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	return orm.LoadFilters(ctx)
}
