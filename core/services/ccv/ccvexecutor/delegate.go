package ccvexecutor

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccv/executor"
	"github.com/smartcontractkit/chainlink-ccv/integration/pkg/constructors"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ccv/ccvcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

type Delegate struct {
	delegateLogger logger.Logger
	lggr           logger.Logger
	// Houses secrets that are needed by the executor (e.g. indexer API keys).
	ccvConfig config.CCV
	// TODO: EVM specific (!)
	chainServices []commontypes.ChainService
	ethKs         keystore.Eth

	isNewlyCreatedJob bool
}

func NewDelegate(lggr logger.Logger, ccvConfig config.CCV, ethKs keystore.Eth, chainServices []commontypes.ChainService) *Delegate {
	return &Delegate{
		delegateLogger: logger.Named(lggr, "CCVExecutorDelegate"),
		lggr:           lggr,
		ccvConfig:      ccvConfig,
		chainServices:  chainServices,
		ethKs:          ethKs,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.CCVExecutor
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {
	d.isNewlyCreatedJob = true
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) (services []job.ServiceCtx, err error) {
	d.delegateLogger.Infow("Creating services for CCV executor job", "jobID", spec.ID)

	var decodedCfg executor.Configuration
	err = toml.Unmarshal([]byte(spec.CCVExecutorSpec.ExecutorConfig), &decodedCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal executorConfig into the executor config struct: %w", err)
	}

	err = decodedCfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate executor config: %w", err)
	}

	err = decodedCfg.Monitoring.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate executor monitoring config: %w", err)
	}

	// Chains in the executor configuration should dictate what we end up verifying for.
	chainsInConfig := make([]protocol.ChainSelector, 0, len(decodedCfg.ChainConfiguration))
	for chainSelStr := range decodedCfg.ChainConfiguration {
		parsed, err2 := strconv.ParseUint(chainSelStr, 10, 64)
		if err2 != nil {
			return nil, fmt.Errorf("failed to parse chain selector string from executor config (%s): %w", chainSelStr, err)
		}
		chainsInConfig = append(chainsInConfig, protocol.ChainSelector(parsed))
	}

	legacyChains, missingChains, err := ccvcommon.GetLegacyChains(ctx, d.lggr, d.chainServices, chainsInConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get legacy chains: %w", err)
	}

	roundRobins := make(map[protocol.ChainSelector]keys.RoundRobin, len(legacyChains))
	fromAddresses := make(map[protocol.ChainSelector][]common.Address, len(legacyChains))

	// A chain that the keystore cannot serve is unusable for the executor, because it cannot
	// transmit on it. Drop the chain and let the monitor report it, instead of stopping the job.
	var unusable []protocol.ChainSelector
	for chainSel := range legacyChains {
		id, err3 := chainselectors.GetChainIDFromSelector(uint64(chainSel))
		if err3 != nil {
			d.lggr.Warnw("skipping chain: failed to get chain ID from selector", "chainSelector", chainSel, "err", err3)
			unusable = append(unusable, chainSel)
			continue
		}
		chainID, ok := new(big.Int).SetString(id, 10)
		if !ok {
			d.lggr.Warnw("skipping chain: failed to convert chain ID to big.Int", "chainSelector", chainSel, "chainID", id)
			unusable = append(unusable, chainSel)
			continue
		}

		addressesForChain, err3 := d.ethKs.EnabledAddressesForChain(ctx, chainID)
		if err3 != nil {
			d.lggr.Warnw("skipping chain: failed to get addresses from eth keystore",
				"chainSelector", chainSel, "chainID", chainID.String(), "err", err3)
			unusable = append(unusable, chainSel)
			continue
		}
		if len(addressesForChain) == 0 {
			d.lggr.Warnw("skipping chain: no enabled address in the eth keystore",
				"chainSelector", chainSel, "chainID", chainID.String())
			unusable = append(unusable, chainSel)
			continue
		}

		roundRobins[chainSel] = NewRoundRobin(d.ethKs, chainID)
		fromAddresses[chainSel] = addressesForChain
	}
	for _, chainSel := range unusable {
		delete(legacyChains, chainSel)
	}
	missingChains = append(missingChains, unusable...)

	if len(legacyChains) == 0 {
		return nil, fmt.Errorf("no usable chain remains out of the %d chains in the configuration: %v",
			len(chainsInConfig), chainsInConfig)
	}

	missingChainsMonitor, err := ccvcommon.NewMissingChainsMonitor(
		d.lggr,
		"CCVExecutor/"+decodedCfg.ExecutorID,
		missingChains,
		ccvcommon.DefaultMissingChainsReportInterval,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create missing chains monitor: %w", err)
	}
	if missingChainsMonitor != nil {
		services = append(services, missingChainsMonitor)
	}

	// TODO: pass secrets as a separate param in the constructor.
	ec, err := constructors.NewExecutorCoordinator(
		logger.Named(
			logger.Named(d.lggr, "CCVExecutorCoordinator"),
			decodedCfg.ExecutorID,
		),
		decodedCfg,
		legacyChains,
		roundRobins,
		fromAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor coordinator: %w", err)
	}
	services = append(services, ec)

	return services, nil
}

func (d *Delegate) AfterJobCreated(spec job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, spec job.Job) error {
	return nil
}

// TODO: this is evm specific, shouldn't be.
type roundRobin struct {
	ks      keystore.Eth
	chainID *big.Int
}

func NewRoundRobin(ks keystore.Eth, chainID *big.Int) *roundRobin {
	return &roundRobin{
		ks:      ks,
		chainID: chainID,
	}
}

func (r *roundRobin) GetNextAddress(ctx context.Context, addresses ...common.Address) (common.Address, error) {
	return r.ks.GetRoundRobinAddress(ctx, r.chainID, addresses...)
}
