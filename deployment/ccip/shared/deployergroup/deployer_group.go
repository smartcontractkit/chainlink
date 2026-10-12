package deployergroup

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"

	pdasol "github.com/smartcontractkit/cld-changesets/pkg/family/solana"
	"golang.org/x/sync/errgroup"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
)

type DeployerGroup struct {
	e                 cldf.Environment
	state             stateview.CCIPOnChainState
	mcmConfig         *cldfproposalutils.TimelockConfig
	deploymentContext *DeploymentContext
	txDecoder         *shared.TxCallDecoder
	describeContext   *shared.ArgumentContext
}

type DescribedTransaction interface {
	Describe() string
	ToMCMS(selector uint64) (mcmstypes.Transaction, error)
}

type EvmDescribedTransaction struct {
	Tx          *types.Transaction
	Description string
}

func (d EvmDescribedTransaction) Describe() string {
	return d.Description
}

func (d EvmDescribedTransaction) ToMCMS(selector uint64) (mcmstypes.Transaction, error) {
	return cldfproposalutils.TransactionForChain(selector, d.Tx.To().Hex(), d.Tx.Data(), d.Tx.Value(), "", []string{})
}

type SolanaDescribedTransaction struct {
	Tx           solana.Instruction
	ProgramID    string
	ContractType cldf.ContractType
	Description  string
}

func (d SolanaDescribedTransaction) Describe() string {
	return d.Description
}

func (d SolanaDescribedTransaction) ToMCMS(selector uint64) (mcmstypes.Transaction, error) {
	data, err := d.Tx.Data()
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to extract data: %w", err)
	}

	// We clear the signer as we will add back timelock PDA signer
	for _, account := range d.Tx.Accounts() {
		if account.IsSigner {
			account.IsSigner = false
		}
	}
	tx, err := mcmsSolana.NewTransaction(
		d.ProgramID,
		data,
		big.NewInt(0),          // e.g. value
		d.Tx.Accounts(),        // pass along needed accounts
		string(d.ContractType), // some string identifying the target
		[]string{},             // any relevant metadata
	)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}
	return tx, nil
}

type DeploymentContext struct {
	description    string
	transactions   map[uint64][]DescribedTransaction
	previousConfig *DeploymentContext
}

func NewDeploymentContext(description string) *DeploymentContext {
	return &DeploymentContext{
		description:    description,
		transactions:   make(map[uint64][]DescribedTransaction),
		previousConfig: nil,
	}
}

func (d *DeploymentContext) Fork(description string) *DeploymentContext {
	return &DeploymentContext{
		description:    description,
		transactions:   make(map[uint64][]DescribedTransaction),
		previousConfig: d,
	}
}

type WithContext interface {
	WithDeploymentContext(description string) *DeployerGroup
}

// DeployerGroupWithContext is a deprecated alias for WithContext, retained so
// existing external consumers keep compiling.
//
// Deprecated: use WithContext.
type DeployerGroupWithContext = WithContext //nolint:revive // stutter accepted to preserve the pre-rename API

type deployerGroupBuilder struct {
	e               cldf.Environment
	state           stateview.CCIPOnChainState
	mcmConfig       *cldfproposalutils.TimelockConfig
	txDecoder       *shared.TxCallDecoder
	describeContext *shared.ArgumentContext
}

func (d *deployerGroupBuilder) WithDeploymentContext(description string) *DeployerGroup {
	return &DeployerGroup{
		e:                 d.e,
		mcmConfig:         d.mcmConfig,
		state:             d.state,
		txDecoder:         d.txDecoder,
		describeContext:   d.describeContext,
		deploymentContext: NewDeploymentContext(description),
	}
}

// DeployerGroup is an abstraction that lets developers write their changeset
// without needing to know if it's executed using a DeployerKey or an MCMS proposal.
//
// Example usage:
//
//	deployerGroup := NewDeployerGroup(e, state, mcmConfig)
//	selector := 0
//	# Get the right deployer key for the chain
//	deployer := deployerGroup.GetDeployer(selector)
//	state.Chains[selector].RMNRemote.Curse()
//	# Execute the transaction or create the proposal
//	deployerGroup.Enact("Curse RMNRemote")
func NewDeployerGroup(e cldf.Environment, state stateview.CCIPOnChainState, mcmConfig *cldfproposalutils.TimelockConfig) WithContext {
	addresses, _ := e.ExistingAddresses.Addresses()
	d := &deployerGroupBuilder{
		e:               e,
		mcmConfig:       mcmConfig,
		state:           state,
		txDecoder:       shared.NewTxCallDecoder(nil),
		describeContext: shared.NewArgumentContext(addresses),
	}
	// update state if timelock needs to be loaded from datastore with qualifier
	if d.mcmConfig != nil && d.mcmConfig.TimelockQualifierPerChain != nil {
		if e.DataStore == nil {
			panic("datastore is required when using timelock qualifiers")
		}
		for selector := range d.mcmConfig.TimelockQualifierPerChain {
			err := d.state.UpdateMCMSStateWithAddressFromDatastoreForChain(d.e, selector, d.mcmConfig.TimelockQualifierPerChain[selector])
			if err != nil {
				panic(fmt.Errorf("failed to load timelock address for chain %d: %w", selector, err))
			}
		}
	}
	return d
}

func (d *DeployerGroup) WithDeploymentContext(description string) *DeployerGroup {
	return &DeployerGroup{
		e:                 d.e,
		mcmConfig:         d.mcmConfig,
		state:             d.state,
		txDecoder:         d.txDecoder,
		describeContext:   d.describeContext,
		deploymentContext: d.deploymentContext.Fork(description),
	}
}

func (d *DeployerGroup) GetDeployer(chain uint64) (*bind.TransactOpts, error) {
	txOpts := d.e.BlockChains.EVMChains()[chain].DeployerKey
	if d.mcmConfig != nil {
		txOpts = cldf.SimTransactOpts()
		txOpts = &bind.TransactOpts{
			From:       d.state.MustGetEVMChainState(chain).Timelock.Address(),
			Signer:     txOpts.Signer,
			GasLimit:   txOpts.GasLimit,
			GasPrice:   txOpts.GasPrice,
			Nonce:      txOpts.Nonce,
			Value:      txOpts.Value,
			GasFeeCap:  txOpts.GasFeeCap,
			GasTipCap:  txOpts.GasTipCap,
			Context:    txOpts.Context,
			AccessList: txOpts.AccessList,
			NoSend:     txOpts.NoSend,
		}
	}
	sim := &bind.TransactOpts{
		From:       txOpts.From,
		Signer:     txOpts.Signer,
		GasLimit:   txOpts.GasLimit,
		GasPrice:   txOpts.GasPrice,
		Nonce:      txOpts.Nonce,
		Value:      txOpts.Value,
		GasFeeCap:  txOpts.GasFeeCap,
		GasTipCap:  txOpts.GasTipCap,
		Context:    txOpts.Context,
		AccessList: txOpts.AccessList,
		NoSend:     true,
	}
	oldSigner := sim.Signer

	var startingNonce *big.Int
	if txOpts.Nonce != nil {
		startingNonce = new(big.Int).Set(txOpts.Nonce)
	} else {
		nonce, err := d.e.BlockChains.EVMChains()[chain].Client.PendingNonceAt(context.Background(), txOpts.From)
		if err != nil {
			return nil, fmt.Errorf("could not get nonce for deployer: %w", err)
		}
		startingNonce = new(big.Int).SetUint64(nonce)
	}
	dc := d.deploymentContext
	sim.Signer = func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
		txCount, err := d.getTransactionCount(chain)
		if err != nil {
			return nil, err
		}

		currentNonce := big.NewInt(0).Add(startingNonce, txCount)

		tx, err := oldSigner(a, t)
		if err != nil {
			return nil, err
		}
		var description string
		if abiStr, ok := d.state.MustGetEVMChainState(chain).ABIByAddress[tx.To().Hex()]; ok {
			_abi, err := abi.JSON(strings.NewReader(abiStr))
			if err != nil {
				return nil, fmt.Errorf("could not get ABI: %w", err)
			}
			decodedCall, err := d.txDecoder.Analyze(tx.To().String(), &_abi, tx.Data())
			if err != nil {
				d.e.Logger.Errorw("could not analyze transaction",
					"chain", chain, "address", tx.To().Hex(), "nonce", currentNonce, "error", err)
			} else {
				description = decodedCall.Describe(d.describeContext)
			}
		}
		dc.transactions[chain] = append(dc.transactions[chain], EvmDescribedTransaction{Tx: tx, Description: description})
		// Update the nonce to consider the transactions that have been sent
		sim.Nonce = big.NewInt(0).Add(currentNonce, big.NewInt(1))
		return tx, nil
	}
	return sim, nil
}

type DeployerForSVM func(solana.PublicKey) (solana.Instruction, string, cldf.ContractType, error)

func (d *DeployerGroup) GetDeployerForSVM(chain uint64) (func(DeployerForSVM) (solana.Instruction, error), error) {
	authority := d.e.BlockChains.SolanaChains()[chain].DeployerKey.PublicKey()
	if d.mcmConfig != nil {
		qualifier := shared.DefaultMCMSQualifier
		if d.mcmConfig.TimelockQualifierPerChain != nil && d.mcmConfig.TimelockQualifierPerChain[chain] != "" {
			qualifier = d.mcmConfig.TimelockQualifierPerChain[chain]
		}
		mcmState, err := solanastateview.LoadMCMSWithTimelockFromDataStore(&d.e, d.e.BlockChains.SolanaChains()[chain], qualifier)
		if err != nil {
			return nil, fmt.Errorf("failed to load mcm state: %w", err)
		}
		timelockSignerPDA := pdasol.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
		authority = timelockSignerPDA
	}

	return func(deployer DeployerForSVM) (solana.Instruction, error) {
		ix, programID, contractType, err := deployer(authority)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy instruction: %w", err)
		}
		dc := d.deploymentContext
		dc.transactions[chain] = append(dc.transactions[chain], SolanaDescribedTransaction{
			Tx:           ix,
			ProgramID:    programID,
			ContractType: contractType,
		})
		return ix, nil
	}, nil
}

func (d *DeployerGroup) getContextChainInOrder() []*DeploymentContext {
	contexts := make([]*DeploymentContext, 0)
	for c := d.deploymentContext; c != nil; c = c.previousConfig {
		contexts = append(contexts, c)
	}
	slices.Reverse(contexts)
	return contexts
}

func (d *DeployerGroup) getTransactions() map[uint64][]DescribedTransaction {
	transactions := make(map[uint64][]DescribedTransaction)
	for _, c := range d.getContextChainInOrder() {
		for k, v := range c.transactions {
			transactions[k] = append(transactions[k], v...)
		}
	}
	return transactions
}

func (d *DeployerGroup) getTransactionCount(chain uint64) (*big.Int, error) {
	txs := d.getTransactions()
	return big.NewInt(int64(len(txs[chain]))), nil
}

func (d *DeployerGroup) Enact() (cldf.ChangesetOutput, error) {
	if d.mcmConfig != nil {
		return d.enactMcms()
	}

	return d.enactDeployer()
}

func ValidateMCMSWithState(env cldf.Environment, selector uint64, mcmConfig *cldfproposalutils.TimelockConfig, state stateview.CCIPOnChainState) error {
	family, err := chain_selectors.GetSelectorFamily(selector)
	if err != nil {
		return fmt.Errorf("failed to get chain selector family: %w", err)
	}

	switch family {
	case chain_selectors.FamilyEVM:
		mcmsState, ok := state.EVMMCMSStateByChain()[selector]
		if !ok {
			return fmt.Errorf("failed to get mcms state for chain %d", selector)
		}
		if err := mcmConfig.Validate(env.BlockChains.EVMChains()[selector], mcmsState); err != nil {
			return fmt.Errorf("mcm config is invalid for chain %d: %w", selector, err)
		}
	case chain_selectors.FamilySolana:
		if err := stateview.ValidateSolanaTimelockConfig(env, selector, mcmConfig); err != nil {
			return fmt.Errorf("mcm config is invalid for chain %d: %w", selector, err)
		}
	case chain_selectors.FamilyTon:
		// TODO : Implement validation for TON, https://smartcontract-it.atlassian.net/browse/NONEVM-1939
		return nil
	default:
		return fmt.Errorf("unsupported chain family: %s", family)
	}
	return nil
}

func (d *DeployerGroup) enactMcms() (cldf.ChangesetOutput, error) {
	contexts := d.getContextChainInOrder()
	proposals := make([]mcmslib.TimelockProposal, 0, len(contexts))
	describedProposals := make([]string, 0, len(contexts))

	for _, dc := range contexts {
		batches := make([]mcmstypes.BatchOperation, 0, len(dc.transactions))
		describedBatches := make([][]string, 0, len(dc.transactions))
		for selector, txs := range dc.transactions {
			err := ValidateMCMSWithState(d.e, selector, d.mcmConfig, d.state)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate mcms state: %w", err)
			}
			mcmTransactions := make([]mcmstypes.Transaction, len(txs))
			describedTxs := make([]string, len(txs))
			for i, tx := range txs {
				var err error
				mcmTransactions[i], err = tx.ToMCMS(selector)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to build mcms transaction: %w", err)
				}
				describedTxs[i] = tx.Describe()
			}

			batches = append(batches, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(selector),
				Transactions:  mcmTransactions,
			})
			describedBatches = append(describedBatches, describedTxs)
		}

		if len(batches) == 0 {
			d.e.Logger.Warnf("No batch was produced from deployment context skipping proposal: %s", dc.description)
			continue
		}

		timelocks, err := BuildTimelockAddressPerChain(d.e, d.state, d.mcmConfig.TimelockQualifierPerChain)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get timelocks for chain: %w", err)
		}
		mcmContractByAction, err := BuildMcmAddressesPerChainByAction(d.e, d.state, d.mcmConfig, d.mcmConfig.TimelockQualifierPerChain)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get proposer mcms for chain: %w", err)
		}
		inspectors, err := cldfproposalutils.McmsInspectors(d.e)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get mcms inspector for chain: %w", err)
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(
			d.e,
			timelocks,
			mcmContractByAction,
			inspectors,
			batches,
			dc.description,
			*d.mcmConfig,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal %w", err)
		}
		describedProposal := shared.DescribeTimelockProposal(proposal, describedBatches)

		// Update the proposal metadata to incorporate the startingOpCount
		// from the previous proposal
		if len(proposals) > 0 {
			previousProposal := proposals[len(proposals)-1]
			for chain, metadata := range previousProposal.ChainMetadata {
				nextStartingOp := metadata.StartingOpCount + getBatchCountForChain(chain, proposal)
				proposal.ChainMetadata[chain] = mcmstypes.ChainMetadata{
					StartingOpCount: nextStartingOp,
					MCMAddress:      proposal.ChainMetadata[chain].MCMAddress,
				}
			}
		}

		proposals = append(proposals, *proposal)
		describedProposals = append(describedProposals, describedProposal)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals:      proposals,
		DescribedTimelockProposals: describedProposals,
	}, nil
}

func getBatchCountForChain(chain mcmstypes.ChainSelector, timelockProposal *mcmslib.TimelockProposal) uint64 {
	var count uint64
	for _, batchOperation := range timelockProposal.Operations {
		if batchOperation.ChainSelector == chain {
			count++
		}
	}
	return count
}

func (d *DeployerGroup) enactDeployer() (cldf.ChangesetOutput, error) {
	contexts := d.getContextChainInOrder()
	evmChains := d.e.BlockChains.EVMChains()
	for _, c := range contexts {
		g := errgroup.Group{}
		for selector, txs := range c.transactions {
			g.Go(func() error {
				for _, tx := range txs {
					if evmTx, ok := tx.(EvmDescribedTransaction); ok {
						err := evmChains[selector].Client.SendTransaction(context.Background(), evmTx.Tx)
						if err != nil {
							return fmt.Errorf("failed to send transaction: %w", err)
						}
						abiStr, ok := d.state.MustGetEVMChainState(selector).ABIByAddress[evmTx.Tx.To().Hex()]
						if ok {
							_, err = cldf.ConfirmIfNoErrorWithABI(evmChains[selector], evmTx.Tx, abiStr, err)
							if err != nil {
								return fmt.Errorf("waiting for tx to be mined failed: %w", err)
							}
						} else {
							_, err = cldf.ConfirmIfNoError(evmChains[selector], evmTx.Tx, err)
							if err != nil {
								return fmt.Errorf("waiting for tx to be mined failed: %w", err)
							}
						}

						d.e.Logger.Infow("Transaction sent", "chain", selector, "tx", evmTx.Tx.Hash().Hex(), "description", tx.Describe())
					} else if solanaTx, ok := tx.(SolanaDescribedTransaction); ok {
						err := d.e.BlockChains.SolanaChains()[selector].Confirm([]solana.Instruction{solanaTx.Tx})
						if err != nil {
							return fmt.Errorf("waiting for tx to be mined failed: %w", err)
						}
					} else {
						return errors.New("transaction type not supported")
					}
				}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	return cldf.ChangesetOutput{}, nil
}

func BuildTimelockAddressPerChain(e cldf.Environment, onchainState stateview.CCIPOnChainState, addressQualifierPerChain map[uint64]string) (map[uint64]string, error) {
	addressPerChain := make(map[uint64]string)
	for _, chain := range e.BlockChains.EVMChains() {
		addressPerChain[chain.Selector] = onchainState.MustGetEVMChainState(chain.Selector).Timelock.Address().Hex()
	}

	// TODO: expose MCMS addresses in the Solana chain state.
	for selector, chain := range e.BlockChains.SolanaChains() {
		qualifier := shared.DefaultMCMSQualifier
		if addressQualifierPerChain != nil && addressQualifierPerChain[selector] != "" {
			qualifier = addressQualifierPerChain[selector]
		}
		mcmState, err := solanastateview.LoadMCMSWithTimelockFromDataStore(&e, chain, qualifier)
		if err != nil {
			return nil, fmt.Errorf("failed to load addresses for chain %d: %w", selector, err)
		}
		addressPerChain[selector] = mcmsSolana.ContractAddress(mcmState.TimelockProgram, mcmsSolana.PDASeed(mcmState.TimelockSeed))
	}

	return addressPerChain, nil
}

func BuildMcmAddressesPerChainByAction(e cldf.Environment, onchainState stateview.CCIPOnChainState, mcmCfg *cldfproposalutils.TimelockConfig, mcmQualifier map[uint64]string) (map[uint64]string, error) {
	if mcmCfg == nil {
		return nil, errors.New("mcm config is nil, cannot get mcms address")
	}
	addressPerChain := make(map[uint64]string)
	for _, chain := range e.BlockChains.EVMChains() {
		mcmContract, err := mcmCfg.MCMBasedOnAction(onchainState.EVMMCMSStateByChain()[chain.Selector])
		if err != nil {
			return nil, fmt.Errorf("failed to get mcms for action %s (chain: %d): %w", mcmCfg.MCMSAction, chain.Selector, err)
		}
		addressPerChain[chain.Selector] = mcmContract.Address().Hex()
	}

	// TODO: expose MCMS addresses in the Solana chain state.
	for selector, chain := range e.BlockChains.SolanaChains() {
		qualifier := shared.DefaultMCMSQualifier
		if mcmQualifier != nil && mcmQualifier[selector] != "" {
			qualifier = mcmQualifier[selector]
		}
		mcmState, err := solanastateview.LoadMCMSWithTimelockFromDataStore(&e, chain, qualifier)
		if err != nil {
			return nil, fmt.Errorf("failed to load addresses for chain %d: %w", selector, err)
		}
		address, err := mcmCfg.MCMBasedOnActionSolana(*mcmState)
		if err != nil {
			return nil, fmt.Errorf("failed to get mcms for action %s (chain: %d): %w", mcmCfg.MCMSAction, selector, err)
		}
		addressPerChain[selector] = address
	}

	return addressPerChain, nil
}
