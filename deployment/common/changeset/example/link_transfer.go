package example

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

const MaxTimelockDelay = 24 * 7 * time.Hour

type TransferConfig struct {
	To    common.Address
	Value *big.Int
}

type LinkTransferConfig struct {
	Transfers  map[uint64][]TransferConfig
	From       common.Address
	McmsConfig *proposalutils.TimelockConfig
}

func getDeployer(e cldf.Environment, chain uint64, mcmConfig *proposalutils.TimelockConfig) *bind.TransactOpts {
	if mcmConfig == nil {
		return e.BlockChains.EVMChains()[chain].DeployerKey
	}

	return cldf.SimTransactOpts()
}

// Validate checks that the LinkTransferConfig is valid.
func (cfg LinkTransferConfig) Validate(e cldf.Environment) error {
	ctx := e.GetContext()
	// Check that Transfers map has at least one chainSel
	if len(cfg.Transfers) == 0 {
		return errors.New("transfers map must have at least one chainSel")
	}

	// Check transfers config values.
	for chainSel, transfers := range cfg.Transfers {
		selector, err := chain_selectors.GetSelectorFamily(chainSel)
		if err != nil {
			return fmt.Errorf("invalid chain selector: %w", err)
		}
		if selector != chain_selectors.FamilyEVM {
			return fmt.Errorf("chain selector %d is not an EVM chain", chainSel)
		}
		chain, ok := e.BlockChains.EVMChains()[chainSel]
		if !ok {
			return fmt.Errorf("chain with selector %d not found", chainSel)
		}
		addrs, err := e.ExistingAddresses.AddressesForChain(chainSel)
		if err != nil {
			return fmt.Errorf("error getting addresses for chain %d: %w", chainSel, err)
		}
		if len(transfers) == 0 {
			return fmt.Errorf("transfers for chainSel %d must have at least one LinkTransfer", chainSel)
		}
		totalAmount := big.NewInt(0)
		linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addrs)
		if err != nil {
			return fmt.Errorf("error loading link token state during validation: %w", err)
		}
		for _, transfer := range transfers {
			if transfer.To == (common.Address{}) {
				return errors.New("'to' address for transfers must be set")
			}
			if transfer.Value == nil {
				return errors.New("value for transfers must be set")
			}
			if transfer.Value.Cmp(big.NewInt(0)) == 0 {
				return errors.New("value for transfers must be non-zero")
			}
			if transfer.Value.Cmp(big.NewInt(0)) == -1 {
				return errors.New("value for transfers must be positive")
			}
			totalAmount.Add(totalAmount, transfer.Value)
		}
		// check that from address has enough funds for the transfers
		balance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, cfg.From)
		if err != nil {
			return fmt.Errorf("error getting balance of sender: %w", err)
		}
		if balance.Cmp(totalAmount) < 0 {
			return fmt.Errorf("sender does not have enough funds for transfers for chain selector %d, required: %s, available: %s", chainSel, totalAmount.String(), balance.String())
		}
	}

	if cfg.McmsConfig == nil {
		return nil
	}

	// Upper bound for min delay (7 days)
	if cfg.McmsConfig.MinDelay > MaxTimelockDelay {
		return errors.New("minDelay must be less than 7 days")
	}

	return nil
}

// initStatePerChain initializes the state for each chain selector on the provided config
func initStatePerChain(cfg *LinkTransferConfig, e cldf.Environment) (
	linkStatePerChain map[uint64]*changeset.LinkTokenState,
	mcmsStatePerChain map[uint64]*changeset.MCMSWithTimelockState,
	err error) {
	linkStatePerChain = map[uint64]*changeset.LinkTokenState{}
	mcmsStatePerChain = map[uint64]*changeset.MCMSWithTimelockState{}
	// Load state for each chain
	chainSelectors := []uint64{}
	for chainSelector := range cfg.Transfers {
		chainSelectors = append(chainSelectors, chainSelector)
	}
	linkStatePerChain, err = changeset.MaybeLoadLinkTokenState(e, chainSelectors)
	if err != nil {
		return nil, nil, err
	}
	mcmsStatePerChain, err = changeset.MaybeLoadMCMSWithTimelockState(e, chainSelectors)
	if err != nil {
		return nil, nil, err
	}
	return linkStatePerChain, mcmsStatePerChain, nil
}

// transferOrBuildTx transfers the LINK tokens or builds the tx for the MCMS proposal
func transferOrBuildTx(
	e cldf.Environment,
	linkState *changeset.LinkTokenState,
	transfer TransferConfig,
	opts *bind.TransactOpts,
	chain cldf_evm.Chain,
	mcmsConfig *proposalutils.TimelockConfig) (*ethTypes.Transaction, error) {
	tx, err := linkState.LinkToken.Transfer(opts, transfer.To, transfer.Value)
	if err != nil {
		return nil, fmt.Errorf("error packing transfer tx data: %w", err)
	}
	// only wait for tx if we are not using MCMS
	if mcmsConfig == nil {
		if _, err := cldf.ConfirmIfNoError(chain, tx, err); err != nil {
			e.Logger.Errorw("Failed to confirm transfer tx", "chain", chain.String(), "err", err)
			return nil, err
		}
	}
	return tx, nil
}

// LinkTransferV2 is an reimplementation of LinkTransfer that uses the new MCMS SDK.
func LinkTransferV2(e cldf.Environment, cfg *LinkTransferConfig) (cldf.ChangesetOutput, error) {
	err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid LinkTransferConfig: %w", err)
	}

	proposerAddressPerChain := map[uint64]string{}
	inspectorPerChain := map[uint64]sdk.Inspector{}
	timelockAddressesPerChain := map[uint64]string{}
	linkStatePerChain, mcmsStatePerChain, err := initStatePerChain(cfg, e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	allBatches := []mcmstypes.BatchOperation{}
	for chainSelector := range cfg.Transfers {
		chain := e.BlockChains.EVMChains()[chainSelector]
		linkAddress := linkStatePerChain[chainSelector].LinkToken.Address()
		mcmsState := mcmsStatePerChain[chainSelector]
		linkState := linkStatePerChain[chainSelector]

		proposerAddressPerChain[chainSelector] = mcmsState.ProposerMcm.Address().Hex()
		inspectorPerChain[chainSelector] = evm.NewInspector(chain.Client)

		timelockAddress := mcmsState.Timelock.Address().Hex()
		timelockAddressesPerChain[chainSelector] = timelockAddress

		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  []mcmstypes.Transaction{},
		}

		opts := getDeployer(e, chainSelector, cfg.McmsConfig)
		totalAmount := big.NewInt(0)
		for _, transfer := range cfg.Transfers[chainSelector] {
			tx, err := transferOrBuildTx(e, linkState, transfer, opts, chain, cfg.McmsConfig)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			op := evm.NewTransaction(linkAddress, tx.Data(), big.NewInt(0), string(types.LinkToken), []string{})
			batch.Transactions = append(batch.Transactions, op)
			totalAmount.Add(totalAmount, transfer.Value)
		}

		allBatches = append(allBatches, batch)
	}

	if cfg.McmsConfig != nil {
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			e,
			timelockAddressesPerChain,
			proposerAddressPerChain,
			inspectorPerChain,
			allBatches,
			"LINK Value transfer proposal",
			*cfg.McmsConfig,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}
