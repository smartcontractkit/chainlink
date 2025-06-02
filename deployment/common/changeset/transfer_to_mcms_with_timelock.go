package changeset

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

type TransferToMCMSWithTimelockConfig struct {
	ContractsByChain map[uint64][]common.Address
	// MCMSConfig is for the accept ownership proposal
	MCMSConfig proposalutils.TimelockConfig
}

type Ownable interface {
	Owner(opts *bind.CallOpts) (common.Address, error)
	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*gethtypes.Transaction, error)
	AcceptOwnership(opts *bind.TransactOpts) (*gethtypes.Transaction, error)
	Address() common.Address
}

func LoadOwnableContract(addr common.Address, client bind.ContractBackend) (common.Address, Ownable, error) {
	// Just using the ownership interface from here.
	c, err := burn_mint_erc677.NewBurnMintERC677(addr, client)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to create contract: %w", err)
	}
	owner, err := c.Owner(nil)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to get owner of contract %s: %w", c.Address(), err)
	}
	return owner, c, nil
}

func (t TransferToMCMSWithTimelockConfig) Validate(e cldf.Environment) error {
	evmChains := e.BlockChains.EVMChains()
	for chainSelector, contracts := range t.ContractsByChain {
		for _, contract := range contracts {
			// Cannot transfer an unknown address.
			// Note this also assures non-zero addresses.
			if exists, err := cldf.AddressBookContains(e.ExistingAddresses, chainSelector, contract.String()); err != nil || !exists {
				if err != nil {
					return fmt.Errorf("failed to check address book: %w", err)
				}
				return fmt.Errorf("contract %s not found in address book", contract)
			}
			owner, _, err := LoadOwnableContract(contract, evmChains[chainSelector].Client)
			if err != nil {
				return fmt.Errorf("failed to load ownable: %w", err)
			}
			if owner != evmChains[chainSelector].DeployerKey.From {
				return fmt.Errorf("contract %s is not owned by the deployer key", contract)
			}
		}
		// If there is no timelock and mcms proposer on the chain, the transfer will fail.
		if _, err := cldf.SearchAddressBook(e.ExistingAddresses, chainSelector, types.RBACTimelock); err != nil {
			return fmt.Errorf("timelock not present on the chain %w", err)
		}
		if _, err := cldf.SearchAddressBook(e.ExistingAddresses, chainSelector, types.ProposerManyChainMultisig); err != nil {
			return fmt.Errorf("mcms proposer not present on the chain %w", err)
		}
	}

	return nil
}

var _ cldf.ChangeSet[TransferToMCMSWithTimelockConfig] = TransferToMCMSWithTimelockV2

// TransferToMCMSWithTimelockV2 is a reimplementation of TransferToMCMSWithTimelock which uses the new MCMS library.
func TransferToMCMSWithTimelockV2(
	e cldf.Environment,
	cfg TransferToMCMSWithTimelockConfig,
) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	batches := []mcmstypes.BatchOperation{}
	timelockAddressByChain := make(map[uint64]string)
	inspectorPerChain := map[uint64]sdk.Inspector{}
	proposerAddressByChain := make(map[uint64]string)
	evmChains := e.BlockChains.EVMChains()
	for chainSelector, contracts := range cfg.ContractsByChain {
		// Already validated that the timelock/proposer exists.
		timelockAddr, _ := cldf.SearchAddressBook(e.ExistingAddresses, chainSelector, types.RBACTimelock)
		proposerAddr, _ := cldf.SearchAddressBook(e.ExistingAddresses, chainSelector, types.ProposerManyChainMultisig)
		timelockAddressByChain[chainSelector] = timelockAddr
		proposerAddressByChain[chainSelector] = proposerAddr
		inspectorPerChain[chainSelector] = evm.NewInspector(evmChains[chainSelector].Client)

		var ops []mcmstypes.Transaction
		for _, contract := range contracts {
			// Just using the ownership interface.
			// Already validated is ownable.
			owner, c, _ := LoadOwnableContract(contract, evmChains[chainSelector].Client)
			if owner.String() == timelockAddr {
				// Already owned by timelock.
				e.Logger.Infof("contract %s already owned by timelock", contract)
				continue
			}
			tx, err := c.TransferOwnership(evmChains[chainSelector].DeployerKey, common.HexToAddress(timelockAddr))
			_, err = cldf.ConfirmIfNoError(evmChains[chainSelector], tx, err)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of contract %T: %w", contract, err)
			}
			tx, err = c.AcceptOwnership(cldf.SimTransactOpts())
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate accept ownership calldata of %s: %w", contract, err)
			}
			ops = append(ops, mcmstypes.Transaction{
				To:               contract.Hex(),
				Data:             tx.Data(),
				AdditionalFields: json.RawMessage(`{"value": 0}`), // JSON-encoded `{"value": 0}`
			})
		}
		batches = append(batches, mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  ops,
		})
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelockAddressByChain, proposerAddressByChain, inspectorPerChain,
		batches, "Transfer ownership to timelock", cfg.MCMSConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w, batches: %+v", err, batches)
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

var _ cldf.ChangeSet[TransferToDeployerConfig] = TransferToDeployer

type TransferToDeployerConfig struct {
	ContractAddress common.Address
	ChainSel        uint64
}

// TransferToDeployer relies on the deployer key
// still being a timelock admin and transfers the ownership of a contract
// back to the deployer key. It's effectively the rollback function of transferring
// to the timelock.
func TransferToDeployer(e cldf.Environment, cfg TransferToDeployerConfig) (cldf.ChangesetOutput, error) {
	evmChains := e.BlockChains.EVMChains()
	owner, ownable, err := LoadOwnableContract(cfg.ContractAddress, evmChains[cfg.ChainSel].Client)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if owner == evmChains[cfg.ChainSel].DeployerKey.From {
		e.Logger.Infof("Contract %s already owned by deployer", cfg.ContractAddress)
		return cldf.ChangesetOutput{}, nil
	}
	tx, err := ownable.TransferOwnership(cldf.SimTransactOpts(), evmChains[cfg.ChainSel].DeployerKey.From)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	addrs, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSel)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	tls, err := MaybeLoadMCMSWithTimelockChainState(evmChains[cfg.ChainSel], addrs)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	calls := []owner_helpers.RBACTimelockCall{
		{
			Target: ownable.Address(),
			Data:   tx.Data(),
			Value:  big.NewInt(0),
		},
	}
	var salt [32]byte
	binary.BigEndian.PutUint32(salt[:], uint32(time.Now().Unix()))
	tx, err = tls.Timelock.ScheduleBatch(evmChains[cfg.ChainSel].DeployerKey, calls, [32]byte{}, salt, big.NewInt(0))
	if _, err = cldf.ConfirmIfNoErrorWithABI(evmChains[cfg.ChainSel], tx, owner_helpers.RBACTimelockABI, err); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	e.Logger.Infof("scheduled transfer ownership batch with tx %s", tx.Hash().Hex())
	timelockExecutorProxy, err := owner_helpers.NewRBACTimelock(tls.CallProxy.Address(), evmChains[cfg.ChainSel].Client)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error creating timelock executor proxy: %w", err)
	}

	tx, err = timelockExecutorProxy.ExecuteBatch(
		evmChains[cfg.ChainSel].DeployerKey, calls, [32]byte{}, salt)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error executing batch: %w", err)
	}
	if _, err = cldf.ConfirmIfNoErrorWithABI(evmChains[cfg.ChainSel], tx, owner_helpers.RBACTimelockABI, err); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	e.Logger.Infof("executed transfer ownership to deployer key with tx %s", tx.Hash().Hex())

	tx, err = ownable.AcceptOwnership(evmChains[cfg.ChainSel].DeployerKey)
	if _, err = cldf.ConfirmIfNoError(evmChains[cfg.ChainSel], tx, err); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	e.Logger.Infof("deployer key accepted ownership tx %s", tx.Hash().Hex())
	return cldf.ChangesetOutput{}, nil
}

var _ cldf.ChangeSet[RenounceTimelockDeployerConfig] = RenounceTimelockDeployer

type RenounceTimelockDeployerConfig struct {
	ChainSel uint64
}

func (cfg RenounceTimelockDeployerConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(cfg.ChainSel); err != nil {
		return fmt.Errorf("invalid chain selector: %w", err)
	}

	_, ok := e.BlockChains.EVMChains()[cfg.ChainSel]
	if !ok {
		return fmt.Errorf("chain selector: %d not found in environment", cfg.ChainSel)
	}

	// MCMS should already exists
	state, err := MaybeLoadMCMSWithTimelockState(e, []uint64{cfg.ChainSel})
	if err != nil {
		return err
	}

	contract, ok := state[cfg.ChainSel]
	if !ok {
		return fmt.Errorf("mcms contracts not found on chain %d", cfg.ChainSel)
	}
	if contract.Timelock == nil {
		return fmt.Errorf("timelock not found on chain %d", cfg.ChainSel)
	}

	return nil
}

// RenounceTimelockDeployer revokes the deployer key from administering the contract.
func RenounceTimelockDeployer(e cldf.Environment, cfg RenounceTimelockDeployerConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	contracts, err := MaybeLoadMCMSWithTimelockState(e, []uint64{cfg.ChainSel})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	tl := contracts[cfg.ChainSel].Timelock
	admin, err := tl.ADMINROLE(&bind.CallOpts{Context: e.GetContext()})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get admin role: %w", err)
	}

	chain := e.BlockChains.EVMChains()[cfg.ChainSel]
	tx, err := tl.RenounceRole(chain.DeployerKey, admin, chain.DeployerKey.From)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to revoke deployer key: %w", err)
	}
	if _, err := cldf.ConfirmIfNoErrorWithABI(chain, tx, owner_helpers.RBACTimelockABI, err); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	e.Logger.Infof("revoked deployer key from owning contract %s", tl.Address().Hex())
	return cldf.ChangesetOutput{}, nil
}
