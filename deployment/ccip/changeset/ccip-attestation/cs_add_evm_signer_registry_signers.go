package ccip_attestation

import (
	"errors"
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcms_types "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/operations/contract"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	evm_datastore_utils "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/utils/datastore"
	// Register the EVM MCMS reader used for owner and proposal resolution.
	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v1_0_0/adapters"
	ccip_changesets "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	datastore_utils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
)

// EVMSignerRegistryAddSignersChangeset adds missing signer entries to existing
// SignerRegistry 1.0.0 contracts on supported Base chains.
var EVMSignerRegistryAddSignersChangeset = cldf.CreateChangeSet(
	addSigners,
	verifyAddSigners,
)

// Signer is one SignerRegistry entry. NewEVMAddress is optional and starts the
// signer with an in-progress key rotation.
type Signer struct {
	EVMAddress    common.Address `json:"evmAddress" yaml:"evmAddress"`
	NewEVMAddress common.Address `json:"newEVMAddress" yaml:"newEVMAddress"`
}

// AddSignersConfig defines the signer entries that should be present on each
// supported Base chain and the settings for any required MCMS proposal.
type AddSignersConfig struct {
	SignersByChain map[uint64][]Signer `json:"signersByChain" yaml:"signersByChain"`
	MCMS           mcms.Input          `json:"mcms" yaml:"mcms"`
}

type addSignersPlan struct {
	chain           cldf_evm.Chain
	registryAddress common.Address
	signers         []Signer
	governed        bool
}

func verifyAddSigners(e cldf.Environment, cfg AddSignersConfig) error {
	plans, err := prepareAddSigners(e, cfg)
	if err != nil {
		return err
	}
	return validateAddSignersProposal(e, cfg.MCMS, plans)
}

func addSigners(e cldf.Environment, cfg AddSignersConfig) (cldf.ChangesetOutput, error) {
	plans, err := prepareAddSigners(e, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := validateAddSignersProposal(e, cfg.MCMS, plans); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	reports := make([]cldf_ops.Report[any, any], 0, len(plans))
	batchOps := make([]mcms_types.BatchOperation, 0, len(plans))
	for _, plan := range plans {
		e.Logger.Infow(
			"Adding SignerRegistry entries",
			"chainSelector", plan.chain.Selector,
			"registry", plan.registryAddress,
			"signerCount", len(plan.signers),
		)
		input := contract.FunctionInput[[]Signer]{
			Address:       plan.registryAddress,
			ChainSelector: plan.chain.Selector,
			Args:          plan.signers,
		}
		// Live-state reconciliation provides idempotency. Force execution prevents a
		// prior successful report from suppressing a legitimate add after removal.
		report, executeErr := cldf_ops.ExecuteOperation(
			e.OperationsBundle,
			addSignersOperation,
			plan.chain,
			input,
			cldf_ops.WithForceExecute[contract.FunctionInput[[]Signer], cldf_evm.Chain](),
		)
		if report.ID != "" {
			reports = append(reports, report.ToGenericReport())
		}
		if executeErr != nil {
			return cldf.ChangesetOutput{Reports: reports}, fmt.Errorf(
				"failed to add signers to SignerRegistry on chain %d: %w", plan.chain.Selector, executeErr,
			)
		}

		batchOp, batchErr := contract.NewBatchOperationFromWrites([]contract.WriteOutput{report.Output})
		if batchErr != nil {
			return cldf.ChangesetOutput{Reports: reports}, fmt.Errorf(
				"failed to create SignerRegistry batch operation on chain %d: %w", plan.chain.Selector, batchErr,
			)
		}
		if len(batchOp.Transactions) > 0 {
			batchOps = append(batchOps, batchOp)
		}
	}

	return ccip_changesets.NewOutputBuilder(e, ccip_changesets.GetRegistry()).
		WithReports(reports).
		WithBatchOps(batchOps).
		Build(cfg.MCMS)
}

func prepareAddSigners(
	e cldf.Environment,
	cfg AddSignersConfig,
) ([]addSignersPlan, error) {
	selectors, err := validateAddSignersConfig(e, cfg)
	if err != nil {
		return nil, err
	}

	plans := make([]addSignersPlan, 0, len(selectors))
	for _, selector := range selectors {
		chain := e.BlockChains.EVMChains()[selector]
		registryAddress, err := resolveSignerRegistry(e, selector)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve SignerRegistry on chain %d: %w", selector, err)
		}
		registry, err := signer_registry.NewSignerRegistry(registryAddress, chain.Client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind SignerRegistry %s on chain %d: %w", registryAddress, selector, err)
		}

		callOpts := &bind.CallOpts{Context: e.GetContext()}
		currentSigners, err := registry.GetSigners(callOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to read SignerRegistry signers on chain %d: %w", selector, err)
		}
		maxSigners, err := registry.GetMaxSigners(callOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to read SignerRegistry capacity on chain %d: %w", selector, err)
		}
		owner, err := registry.Owner(callOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to read SignerRegistry owner on chain %d: %w", selector, err)
		}
		if maxSigners == nil || !maxSigners.IsUint64() {
			return nil, fmt.Errorf("SignerRegistry on chain %d returned invalid capacity %v", selector, maxSigners)
		}
		currentCount := uint64(len(currentSigners))
		maxCapacity := maxSigners.Uint64()
		if currentCount > maxCapacity {
			return nil, fmt.Errorf(
				"SignerRegistry on chain %d has %d signers, exceeding its capacity %d",
				selector, currentCount, maxCapacity,
			)
		}

		governed, err := validateSignerRegistryOwner(e, cfg.MCMS, chain, owner)
		if err != nil {
			return nil, err
		}
		additions, err := reconcileSigners(selector, currentSigners, cfg.SignersByChain[selector])
		if err != nil {
			return nil, err
		}
		if len(additions) == 0 {
			e.Logger.Infow(
				"SignerRegistry additions already applied; skipping",
				"chainSelector", selector,
				"registry", registryAddress,
			)
			continue
		}
		additionCount := uint64(len(additions))
		if additionCount > maxCapacity-currentCount {
			return nil, fmt.Errorf(
				"adding %d signers on chain %d exceeds SignerRegistry capacity: %d current, %d maximum",
				len(additions), selector, len(currentSigners), maxCapacity,
			)
		}

		plans = append(plans, addSignersPlan{
			chain:           chain,
			registryAddress: registryAddress,
			signers:         additions,
			governed:        governed,
		})
	}

	return plans, nil
}

func validateAddSignersConfig(
	e cldf.Environment,
	cfg AddSignersConfig,
) ([]uint64, error) {
	if len(cfg.SignersByChain) == 0 {
		return nil, errors.New("no signer additions provided")
	}
	if e.DataStore == nil {
		return nil, errors.New("environment DataStore is required")
	}

	selectors := make([]uint64, 0, len(cfg.SignersByChain))
	for selector, signers := range cfg.SignersByChain {
		if selector != BaseMainnetSelector && selector != BaseSepoliaSelector {
			return nil, fmt.Errorf("chain selector %d is not a supported Base chain", selector)
		}
		chain, exists := e.BlockChains.EVMChains()[selector]
		if !exists {
			return nil, fmt.Errorf("EVM chain selector %d not found in environment", selector)
		}
		if chain.DeployerKey == nil {
			return nil, fmt.Errorf("EVM chain selector %d has no deployer key", selector)
		}
		if err := validateSigners(signers); err != nil {
			return nil, fmt.Errorf("invalid signer additions for chain %d: %w", selector, err)
		}
		selectors = append(selectors, selector)
	}
	slices.Sort(selectors)
	return selectors, nil
}

func resolveSignerRegistry(e cldf.Environment, selector uint64) (common.Address, error) {
	return datastore_utils.FindAndFormatRef(
		e.DataStore,
		datastore.AddressRef{
			Type:    datastore.ContractType(shared.EVMSignerRegistry),
			Version: &deployment.Version1_0_0,
		},
		selector,
		evm_datastore_utils.ToNonZeroEVMAddress,
	)
}

func validateSignerRegistryOwner(
	e cldf.Environment,
	mcmsCfg mcms.Input,
	chain cldf_evm.Chain,
	owner common.Address,
) (bool, error) {
	if owner == chain.DeployerKey.From {
		return false, nil
	}

	reader, ok := ccip_changesets.GetRegistry().GetMCMSReader(chain_selectors.FamilyEVM)
	if !ok {
		return false, errors.New("no MCMS reader registered for EVM chains")
	}
	timelockRef, err := reader.GetTimelockRef(e, chain.Selector, mcmsCfg)
	if err != nil {
		return false, fmt.Errorf("failed to resolve MCMS timelock owner on chain %d: %w", chain.Selector, err)
	}
	timelockAddress, err := evm_datastore_utils.ToNonZeroEVMAddress(timelockRef)
	if err != nil {
		return false, fmt.Errorf("failed to resolve MCMS timelock owner on chain %d: %w", chain.Selector, err)
	}
	if owner != timelockAddress {
		return false, fmt.Errorf(
			"SignerRegistry on chain %d has unsupported owner %s; expected deployer %s or MCMS timelock %s",
			chain.Selector, owner.Hex(), chain.DeployerKey.From.Hex(), timelockAddress.Hex(),
		)
	}
	return true, nil
}

func reconcileSigners(
	selector uint64,
	current []signer_registry.ISignerRegistrySigner,
	requested []Signer,
) ([]Signer, error) {
	activeSigners := make(map[common.Address]signer_registry.ISignerRegistrySigner, len(current))
	usedKeys := make(map[common.Address]struct{}, len(current)*2)
	for _, signer := range current {
		activeSigners[signer.EvmAddress] = signer
		usedKeys[signer.EvmAddress] = struct{}{}
		if signer.NewEVMAddress != (common.Address{}) {
			usedKeys[signer.NewEVMAddress] = struct{}{}
		}
	}

	additions := make([]Signer, 0, len(requested))
	for _, signer := range requested {
		if existing, exists := activeSigners[signer.EVMAddress]; exists {
			if signer.NewEVMAddress == (common.Address{}) || signer.NewEVMAddress == existing.NewEVMAddress {
				continue
			}
			return nil, fmt.Errorf(
				"signer %s on chain %d already exists with pending key %s; add-signers cannot change it to %s",
				signer.EVMAddress.Hex(), selector, existing.NewEVMAddress.Hex(), signer.NewEVMAddress.Hex(),
			)
		}
		if _, exists := usedKeys[signer.EVMAddress]; exists {
			return nil, fmt.Errorf("active key %s is already in use on chain %d", signer.EVMAddress.Hex(), selector)
		}
		if signer.NewEVMAddress != (common.Address{}) {
			if _, exists := usedKeys[signer.NewEVMAddress]; exists {
				return nil, fmt.Errorf("pending key %s is already in use on chain %d", signer.NewEVMAddress.Hex(), selector)
			}
		}
		additions = append(additions, signer)
	}
	return additions, nil
}

func validateAddSignersProposal(
	e cldf.Environment,
	mcmsCfg mcms.Input,
	plans []addSignersPlan,
) error {
	hasGovernedWrite := false
	for _, plan := range plans {
		if plan.governed {
			hasGovernedWrite = true
			break
		}
	}
	if !hasGovernedWrite {
		return nil
	}
	if err := mcmsCfg.Validate(); err != nil {
		return fmt.Errorf("invalid MCMS proposal configuration: %w", err)
	}

	reader, ok := ccip_changesets.GetRegistry().GetMCMSReader(chain_selectors.FamilyEVM)
	if !ok {
		return errors.New("no MCMS reader registered for EVM chains")
	}
	for _, plan := range plans {
		if !plan.governed {
			continue
		}
		if _, err := reader.GetChainMetadata(e, plan.chain.Selector, mcmsCfg); err != nil {
			return fmt.Errorf("failed to validate MCMS metadata on chain %d: %w", plan.chain.Selector, err)
		}
	}
	return nil
}
