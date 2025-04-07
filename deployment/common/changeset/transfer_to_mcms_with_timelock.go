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
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/deployment"
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

func (t TransferToMCMSWithTimelockConfig) Validate(e deployment.Environment) error {
	for chainSelector, contracts := range t.ContractsByChain {
		for _, contract := range contracts {
			// Cannot transfer an unknown address.
			// Note this also assures non-zero addresses.
			if exists, err := deployment.AddressBookContains(e.ExistingAddresses, chainSelector, contract.String()); err != nil || !exists {
				if err != nil {
					return fmt.Errorf("failed to check address book: %w", err)
				}
				return fmt.Errorf("contract %s not found in address book", contract)
			}
			owner, _, err := LoadOwnableContract(contract, e.Chains[chainSelector].Client)
			if err != nil {
				return fmt.Errorf("failed to load ownable: %w", err)
			}
			if owner != e.Chains[chainSelector].DeployerKey.From {
				return fmt.Errorf("contract %s is not owned by the deployer key", contract)
			}
		}
		// If there is no timelock and mcms proposer on the chain, the transfer will fail.
		if _, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.RBACTimelock); err != nil {
			return fmt.Errorf("timelock not present on the chain %w", err)
		}
		if _, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.ProposerManyChainMultisig); err != nil {
			return fmt.Errorf("mcms proposer not present on the chain %w", err)
		}
	}

	return nil
}

var _ deployment.ChangeSet[TransferToMCMSWithTimelockConfig] = TransferToMCMSWithTimelock

// TransferToMCMSWithTimelock creates a changeset that transfers ownership of all the
// contracts in the provided configuration to the timelock on the chain and generates
// a corresponding accept ownership proposal to complete the transfer.
// It assumes that DeployMCMSWithTimelockV2 has already been run s.t.
// the timelock and mcmses exist on the chain and that the proposed addresses to transfer ownership
// are currently owned by the deployer key.
// Deprecated: Use TransferToMCMSWithTimelockV2 instead.
func TransferToMCMSWithTimelock(
	e deployment.Environment,
	cfg TransferToMCMSWithTimelockConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []timelock.BatchChainOperation
	timelocksByChain := make(map[uint64]common.Address)
	proposersByChain := make(map[uint64]*owner_helpers.ManyChainMultiSig)
	for chainSelector, contracts := range cfg.ContractsByChain {
		// Already validated that the timelock/proposer exists.
		timelockAddr, _ := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.RBACTimelock)
		proposerAddr, _ := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.ProposerManyChainMultisig)
		timelocksByChain[chainSelector] = common.HexToAddress(timelockAddr)
		proposer, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(proposerAddr), e.Chains[chainSelector].Client)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create proposer mcms: %w", err)
		}
		proposersByChain[chainSelector] = proposer

		var ops []mcms.Operation
		for _, contract := range contracts {
			// Just using the ownership interface.
			// Already validated is ownable.
			owner, c, _ := LoadOwnableContract(contract, e.Chains[chainSelector].Client)
			if owner.String() == timelockAddr {
				// Already owned by timelock.
				e.Logger.Infof("contract %s already owned by timelock", contract)
				continue
			}
			tx, err := c.TransferOwnership(e.Chains[chainSelector].DeployerKey, common.HexToAddress(timelockAddr))
			_, err = deployment.ConfirmIfNoError(e.Chains[chainSelector], tx, err)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of contract %T: %w", contract, err)
			}
			tx, err = c.AcceptOwnership(deployment.SimTransactOpts())
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate accept ownership calldata of %s: %w", contract, err)
			}
			ops = append(ops, mcms.Operation{
				To:    contract,
				Data:  tx.Data(),
				Value: big.NewInt(0),
			})
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chainSelector),
			Batch:           ops,
		})
	}
	proposal, err := proposalutils.BuildProposalFromBatches(
		timelocksByChain, proposersByChain, batches, "Transfer ownership to timelock", cfg.MCMSConfig.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w, batches: %+v", err, batches)
	}

	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{*proposal}}, nil
}

var _ deployment.ChangeSet[TransferToMCMSWithTimelockConfig] = TransferToMCMSWithTimelockV2

// TransferToMCMSWithTimelockV2 is a reimplementation of TransferToMCMSWithTimelock which uses the new MCMS library.
func TransferToMCMSWithTimelockV2(
	e deployment.Environment,
	cfg TransferToMCMSWithTimelockConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	batches := []mcmstypes.BatchOperation{}
	timelockAddressByChain := make(map[uint64]string)
	inspectorPerChain := map[uint64]sdk.Inspector{}
	proposerAddressByChain := make(map[uint64]string)
	for chainSelector, contracts := range cfg.ContractsByChain {
		// Already validated that the timelock/proposer exists.
		timelockAddr, _ := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.RBACTimelock)
		proposerAddr, _ := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.ProposerManyChainMultisig)
		timelockAddressByChain[chainSelector] = timelockAddr
		proposerAddressByChain[chainSelector] = proposerAddr
		inspectorPerChain[chainSelector] = evm.NewInspector(e.Chains[chainSelector].Client)

		var ops []mcmstypes.Transaction
		for _, contract := range contracts {
			// Just using the ownership interface.
			// Already validated is ownable.
			owner, c, _ := LoadOwnableContract(contract, e.Chains[chainSelector].Client)
			if owner.String() == timelockAddr {
				// Already owned by timelock.
				e.Logger.Infof("contract %s already owned by timelock", contract)
				continue
			}
			tx, err := c.TransferOwnership(e.Chains[chainSelector].DeployerKey, common.HexToAddress(timelockAddr))
			_, err = deployment.ConfirmIfNoError(e.Chains[chainSelector], tx, err)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of contract %T: %w", contract, err)
			}
			tx, err = c.AcceptOwnership(deployment.SimTransactOpts())
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate accept ownership calldata of %s: %w", contract, err)
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
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w, batches: %+v", err, batches)
	}

	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

var _ deployment.ChangeSet[TransferToDeployerConfig] = TransferToDeployer

type TransferToDeployerConfig struct {
	ContractAddress common.Address
	ChainSel        uint64
}

// TransferToDeployer relies on the deployer key
// still being a timelock admin and transfers the ownership of a contract
// back to the deployer key. It's effectively the rollback function of transferring
// to the timelock.
func TransferToDeployer(e deployment.Environment, cfg TransferToDeployerConfig) (deployment.ChangesetOutput, error) {
	owner, ownable, err := LoadOwnableContract(cfg.ContractAddress, e.Chains[cfg.ChainSel].Client)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if owner == e.Chains[cfg.ChainSel].DeployerKey.From {
		e.Logger.Infof("Contract %s already owned by deployer", cfg.ContractAddress)
		return deployment.ChangesetOutput{}, nil
	}
	tx, err := ownable.TransferOwnership(deployment.SimTransactOpts(), e.Chains[cfg.ChainSel].DeployerKey.From)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	addrs, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSel)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	tls, err := MaybeLoadMCMSWithTimelockChainState(e.Chains[cfg.ChainSel], addrs)
	if err != nil {
		return deployment.ChangesetOutput{}, err
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
	tx, err = tls.Timelock.ScheduleBatch(e.Chains[cfg.ChainSel].DeployerKey, calls, [32]byte{}, salt, big.NewInt(0))
	if _, err = deployment.ConfirmIfNoErrorWithABI(e.Chains[cfg.ChainSel], tx, owner_helpers.RBACTimelockABI, err); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infof("scheduled transfer ownership batch with tx %s", tx.Hash().Hex())
	timelockExecutorProxy, err := owner_helpers.NewRBACTimelock(tls.CallProxy.Address(), e.Chains[cfg.ChainSel].Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error creating timelock executor proxy: %w", err)
	}

	tx, err = timelockExecutorProxy.ExecuteBatch(
		e.Chains[cfg.ChainSel].DeployerKey, calls, [32]byte{}, salt)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error executing batch: %w", err)
	}
	if _, err = deployment.ConfirmIfNoErrorWithABI(e.Chains[cfg.ChainSel], tx, owner_helpers.RBACTimelockABI, err); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	e.Logger.Infof("executed transfer ownership to deployer key with tx %s", tx.Hash().Hex())

	tx, err = ownable.AcceptOwnership(e.Chains[cfg.ChainSel].DeployerKey)
	if _, err = deployment.ConfirmIfNoError(e.Chains[cfg.ChainSel], tx, err); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infof("deployer key accepted ownership tx %s", tx.Hash().Hex())
	return deployment.ChangesetOutput{}, nil
}

var _ deployment.ChangeSet[RenounceTimelockDeployerConfig] = RenounceTimelockDeployer

type RenounceTimelockDeployerConfig struct {
	ChainSel uint64
}

func (cfg RenounceTimelockDeployerConfig) Validate(e deployment.Environment) error {
	if err := deployment.IsValidChainSelector(cfg.ChainSel); err != nil {
		return fmt.Errorf("invalid chain selector: %w", err)
	}

	_, ok := e.Chains[cfg.ChainSel]
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
func RenounceTimelockDeployer(e deployment.Environment, cfg RenounceTimelockDeployerConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	contracts, err := MaybeLoadMCMSWithTimelockState(e, []uint64{cfg.ChainSel})
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	tl := contracts[cfg.ChainSel].Timelock
	admin, err := tl.ADMINROLE(&bind.CallOpts{Context: e.GetContext()})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get admin role: %w", err)
	}

	chain := e.Chains[cfg.ChainSel]
	tx, err := tl.RenounceRole(chain.DeployerKey, admin, chain.DeployerKey.From)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to revoke deployer key: %w", err)
	}
	if _, err := deployment.ConfirmIfNoErrorWithABI(chain, tx, owner_helpers.RBACTimelockABI, err); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infof("revoked deployer key from owning contract %s", tl.Address().Hex())
	return deployment.ChangesetOutput{}, nil
}
