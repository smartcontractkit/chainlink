package workflow

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/scylladb/go-reflectx"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	creconfig "github.com/smartcontractkit/cre-tools/pkg/blockchain/config"
	capreg "github.com/smartcontractkit/cre-tools/pkg/blockchain/contracts/capability_registry"
	workflowreg "github.com/smartcontractkit/cre-tools/pkg/blockchain/contracts/workflow_registry"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	cap_reg_v2 "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	wf_reg_v2_op "github.com/smartcontractkit/chainlink/deployment/cre/workflow_registry/v2/changeset/operations/contracts"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/stagegen"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

// must match nubmer of events we track in core/services/workflows/syncer/handler.go
const NumberOfTrackedWorkflowRegistryEvents = 6

// NewWorkflowRegistryClient creates a new cre-tools workflow registry client.
// Caller is responsible for calling Close() on the returned client.
func NewWorkflowRegistryClient(sc *seth.Client, workflowRegistryAddr common.Address) (*workflowreg.Client, error) {
	clientConfig := creconfig.ClientConfig{
		SethClient:      sc,
		ContractAddress: workflowRegistryAddr.Hex(),
	}
	wrc, err := workflowreg.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow registry client: %w", err)
	}
	return wrc, nil
}

// NewCapabilityRegistryClient creates a new cre-tools capability registry client.
// Caller is responsible for calling Close() on the returned client.
func NewCapabilityRegistryClient(sc *seth.Client, capabilityRegistryAddr common.Address) (*capreg.Client, error) {
	clientConfig := creconfig.ClientConfig{
		SethClient:      sc,
		ContractAddress: capabilityRegistryAddr.Hex(),
	}
	client, err := capreg.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create capability registry client: %w", err)
	}
	return client, nil
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

	allowedDonIDs := make([]uint32, len(input.AllowedDonIDs))
	for i, donID := range input.AllowedDonIDs {
		allowedDonIDs[i] = libc.MustSafeUint32FromUint64(donID)
	}

	switch input.ContractVersion.Version.Major() {
	case 2:
		chain, ok := input.CldEnv.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return nil, fmt.Errorf("chain %d not found in environment", input.ChainSelector)
		}
		contract, err := workflow_registry_wrapper_v2.NewWorkflowRegistry(
			input.ContractAddress, chain.Client,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create WorkflowRegistry instance")
		}
		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*input.CldEnv,
			nil,
			nil,
			contract.Address(),
			cap_reg_v2.ConfigureForwarderDescription,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create transaction strategy")
		}
		updateSignersReport, err := operations.ExecuteOperation(
			input.CldEnv.OperationsBundle,
			wf_reg_v2_op.UpdateAllowedSignersOp,
			wf_reg_v2_op.WorkflowRegistryOpDeps{
				Env:      input.CldEnv,
				Registry: contract,
				Strategy: strategy,
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
				Env:      input.CldEnv,
				Registry: contract,
				Strategy: strategy,
			},
			wf_reg_v2_op.SetDONLimitOpInput{
				ChainSelector:    input.ChainSelector,
				DONFamily:        config.DefaultDONFamily,
				DONLimit:         libc.MustSafeUint32(1000),
				UserDefaultLimit: libc.MustSafeUint32(100),
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
	default:
		report, err := operations.ExecuteSequence(
			input.CldEnv.OperationsBundle,
			ks_contracts_op.ConfigWorkflowRegistrySeq,
			ks_contracts_op.ConfigWorkflowRegistrySeqDeps{
				Env: input.CldEnv,
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
}

// WaitForAllNodesToHaveExpectedFiltersRegistered manually checks if all WorkflowRegistry filters used by the LogPoller are registered for all nodes. We want to see if this will help with the flakiness.
func WaitForAllNodesToHaveExpectedFiltersRegistered(ctx context.Context, singleFileLogger logger.Logger, testLogger zerolog.Logger, registryChainID uint64, dons *cre.Dons, nodeSet []*cre.NodeSet) error {
	for donIdx, don := range dons.List() {
		if !flags.HasFlag(don.Flags, cre.WorkflowDON) {
			continue
		}

		workerNodes, wErr := don.Workers()
		if wErr != nil {
			return errors.Wrap(wErr, "failed to find worker nodes")
		}

		results := make(map[int]bool)
		tickInterval := 5 * time.Second
		timeoutDuration := 2 * time.Minute

		checkCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
		ticker := time.NewTicker(tickInterval)

	INNER_LOOP:
		for {
			select {
			case <-checkCtx.Done():
				cancel()
				ticker.Stop()
				if errors.Is(checkCtx.Err(), context.DeadlineExceeded) {
					return fmt.Errorf("timed out after %.2f seconds waiting for all nodes to have expected filters registered", timeoutDuration.Seconds())
				}
				return fmt.Errorf("context cancelled while waiting for all nodes to have expected filters registered: %w", checkCtx.Err())
			case <-ticker.C:
				if len(results) == len(workerNodes) {
					testLogger.Info().Msgf("All %d nodes in DON %d have expected filters registered", len(workerNodes), don.ID)
					cancel()
					ticker.Stop()

					break INNER_LOOP
				}

				for _, workerNode := range workerNodes {
					if _, ok := results[workerNode.Index]; ok {
						continue
					}

					testLogger.Info().Msgf("Checking if all WorkflowRegistry filters are registered for worker node %d", workerNode.Index)
					allFilters, filtersErr := getAllFilters(checkCtx, singleFileLogger, big.NewInt(libc.MustSafeInt64(registryChainID)), workerNode.Index, nodeSet[donIdx].DbInput.Port)
					if filtersErr != nil {
						cancel()
						ticker.Stop()
						return errors.Wrap(filtersErr, "failed to get filters")
					}

					for _, filter := range allFilters {
						if strings.Contains(filter.Name, "WorkflowRegistry") {
							if len(filter.EventSigs) == NumberOfTrackedWorkflowRegistryEvents {
								testLogger.Debug().Msgf("Found all WorkflowRegistry filters for node %d", workerNode.Index)
								results[workerNode.Index] = true
								continue
							}

							testLogger.Debug().Msgf("Found only %d WorkflowRegistry filters for node %d", len(filter.EventSigs), workerNode.Index)
						}
					}
				}

				// return if we have results for all nodes, don't wait for next tick
				if len(results) == len(workerNodes) {
					testLogger.Info().Msgf("All %d nodes in DON %d have expected filters registered", len(workerNodes), don.ID)
					cancel()
					ticker.Stop()

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
