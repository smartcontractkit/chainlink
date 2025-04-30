package changeset

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type DeployerGroup struct {
	e                 deployment.Environment
	state             CCIPOnChainState
	mcmConfig         *proposalutils.TimelockConfig
	deploymentContext *DeploymentContext
	txDecoder         *proposalutils.TxCallDecoder
	describeContext   *proposalutils.ArgumentContext
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
	return proposalutils.TransactionForChain(selector, d.Tx.To().Hex(), d.Tx.Data(), d.Tx.Value(), "", []string{})
}

type SolanaDescribedTransaction struct {
	Tx           solana.Instruction
	ProgramID    string
	ContractType deployment.ContractType
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

type DeployerGroupWithContext interface {
	WithDeploymentContext(description string) *DeployerGroup
}

type deployerGroupBuilder struct {
	e               deployment.Environment
	state           CCIPOnChainState
	mcmConfig       *proposalutils.TimelockConfig
	txDecoder       *proposalutils.TxCallDecoder
	describeContext *proposalutils.ArgumentContext
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
func NewDeployerGroup(e deployment.Environment, state CCIPOnChainState, mcmConfig *proposalutils.TimelockConfig) DeployerGroupWithContext {
	addresses, _ := e.ExistingAddresses.Addresses()
	return &deployerGroupBuilder{
		e:               e,
		mcmConfig:       mcmConfig,
		state:           state,
		txDecoder:       proposalutils.NewTxCallDecoder(nil),
		describeContext: proposalutils.NewArgumentContext(addresses),
	}
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
	txOpts := d.e.Chains[chain].DeployerKey
	if d.mcmConfig != nil {
		txOpts = deployment.SimTransactOpts()
		txOpts = &bind.TransactOpts{
			From:       d.state.Chains[chain].Timelock.Address(),
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
		nonce, err := d.e.Chains[chain].Client.PendingNonceAt(context.Background(), txOpts.From)
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
		if abiStr, ok := d.state.Chains[chain].ABIByAddress[tx.To().Hex()]; ok {
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

type DeployerForSVM func(solana.PublicKey) (solana.Instruction, string, deployment.ContractType, error)

func (d *DeployerGroup) GetDeployerForSVM(chain uint64) (func(DeployerForSVM) (solana.Instruction, error), error) {
	var authority solana.PublicKey = d.e.SolChains[chain].DeployerKey.PublicKey()

	if d.mcmConfig != nil {
		addresses, err := addressForChain(d.e, chain)
		if err != nil {
			return nil, fmt.Errorf("failed to load addresses for chain %d: %w", chain, err)
		}
		mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(d.e.SolChains[chain], addresses)
		if err != nil {
			return nil, fmt.Errorf("failed to load mcm state: %w", err)
		}
		timelockSignerPDA := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
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

func (d *DeployerGroup) Enact() (deployment.ChangesetOutput, error) {
	if d.mcmConfig != nil {
		return d.enactMcms()
	}

	return d.enactDeployer()
}

func ValidateMCMS(env deployment.Environment, selector uint64, mcmConfig *proposalutils.TimelockConfig) error {
	family, err := chain_selectors.GetSelectorFamily(selector)
	if err != nil {
		return fmt.Errorf("failed to get chain selector family: %w", err)
	}

	switch family {
	case chain_selectors.FamilyEVM:
		state, err := LoadOnchainState(env)
		if err != nil {
			return fmt.Errorf("failed to load onchain state: %w", err)
		}
		mcmsState, ok := state.EVMMCMSStateByChain()[selector]
		if !ok {
			return fmt.Errorf("failed to get mcms state for chain %d", selector)
		}
		if err := mcmConfig.Validate(env.Chains[selector], mcmsState); err != nil {
			return fmt.Errorf("mcm config is invalid for chain %d: %w", selector, err)
		}
	case chain_selectors.FamilySolana:
		if err := mcmConfig.ValidateSolana(env, selector); err != nil {
			return fmt.Errorf("mcm config is invalid for chain %d: %w", selector, err)
		}
	default:
		return fmt.Errorf("unsupported chain family: %s", family)
	}
	return nil
}

func (d *DeployerGroup) enactMcms() (deployment.ChangesetOutput, error) {
	contexts := d.getContextChainInOrder()
	proposals := make([]mcmslib.TimelockProposal, 0, len(contexts))
	describedProposals := make([]string, 0, len(contexts))
	for _, dc := range contexts {
		batches := make([]mcmstypes.BatchOperation, 0, len(dc.transactions))
		describedBatches := make([][]string, 0, len(dc.transactions))
		for selector, txs := range dc.transactions {
			err := ValidateMCMS(d.e, selector, d.mcmConfig)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate mcms state: %w", err)
			}
			mcmTransactions := make([]mcmstypes.Transaction, len(txs))
			describedTxs := make([]string, len(txs))
			for i, tx := range txs {
				var err error
				mcmTransactions[i], err = tx.ToMCMS(selector)
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to build mcms transaction: %w", err)
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

		timelocks := BuildTimelockAddressPerChain(d.e, d.state)
		mcmContractByAction, err := BuildMcmAddressesPerChainByAction(d.e, d.state, d.mcmConfig)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get proposer mcms for chain: %w", err)
		}
		inspectors, err := proposalutils.McmsInspectors(d.e)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get mcms inspector for chain: %w", err)
		}

		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			d.e,
			timelocks,
			mcmContractByAction,
			inspectors,
			batches,
			dc.description,
			*d.mcmConfig,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal %w", err)
		}
		describedProposal := proposalutils.DescribeTimelockProposal(proposal, describedBatches)

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

	return deployment.ChangesetOutput{
		MCMSTimelockProposals:      proposals,
		DescribedTimelockProposals: describedProposals,
	}, nil
}

func getBatchCountForChain(chain mcmstypes.ChainSelector, timelockProposal *mcmslib.TimelockProposal) uint64 {
	batches := make([]mcmstypes.BatchOperation, 0)
	for _, batchOperation := range timelockProposal.Operations {
		if batchOperation.ChainSelector == chain {
			batches = append(batches, batchOperation)
		}
	}

	return uint64(len(batches))
}

func (d *DeployerGroup) enactDeployer() (deployment.ChangesetOutput, error) {
	contexts := d.getContextChainInOrder()
	for _, c := range contexts {
		g := errgroup.Group{}
		for selector, txs := range c.transactions {
			selector, txs := selector, txs
			g.Go(func() error {
				for _, tx := range txs {
					if evmTx, ok := tx.(EvmDescribedTransaction); ok {
						err := d.e.Chains[selector].Client.SendTransaction(context.Background(), evmTx.Tx)
						if err != nil {
							return fmt.Errorf("failed to send transaction: %w", err)
						}
						// TODO how to pass abi here to decode error reason
						_, err = deployment.ConfirmIfNoError(d.e.Chains[selector], evmTx.Tx, err)
						if err != nil {
							return fmt.Errorf("waiting for tx to be mined failed: %w", err)
						}
						d.e.Logger.Infow("Transaction sent", "chain", selector, "tx", evmTx.Tx.Hash().Hex(), "description", tx.Describe())
					} else if solanaTx, ok := tx.(SolanaDescribedTransaction); ok {
						err := d.e.SolChains[selector].Confirm([]solana.Instruction{solanaTx.Tx})
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
			return deployment.ChangesetOutput{}, err
		}
	}
	return deployment.ChangesetOutput{}, nil
}

func BuildTimelockPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*proposalutils.TimelockExecutionContracts {
	timelocksPerChain := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	for _, chain := range e.Chains {
		timelocksPerChain[chain.Selector] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[chain.Selector].Timelock,
			CallProxy: state.Chains[chain.Selector].CallProxy,
		}
	}
	return timelocksPerChain
}

func BuildTimelockAddressPerChain(e deployment.Environment, onchainState CCIPOnChainState) map[uint64]string {
	addressPerChain := make(map[uint64]string)
	for _, chain := range e.Chains {
		addressPerChain[chain.Selector] = onchainState.Chains[chain.Selector].Timelock.Address().Hex()
	}

	// TODO: This should come from the Solana chain state which should be enhanced to contain timlock and MCMS address
	for selector, chain := range e.SolChains {
		addresses, _ := addressForChain(e, selector)
		mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
		addressPerChain[selector] = mcmsSolana.ContractAddress(mcmState.TimelockProgram, mcmsSolana.PDASeed(mcmState.TimelockSeed))
	}

	return addressPerChain
}

func BuildMcmAddressesPerChainByAction(e deployment.Environment, onchainState CCIPOnChainState, mcmCfg *proposalutils.TimelockConfig) (map[uint64]string, error) {
	if mcmCfg == nil {
		return nil, errors.New("mcm config is nil, cannot get mcms address")
	}
	addressPerChain := make(map[uint64]string)
	for _, chain := range e.Chains {
		mcmContract, err := mcmCfg.MCMBasedOnAction(onchainState.EVMMCMSStateByChain()[chain.Selector])
		if err != nil {
			return nil, fmt.Errorf("failed to get mcms for action %s: %w", mcmCfg.MCMSAction, err)
		}
		addressPerChain[chain.Selector] = mcmContract.Address().Hex()
	}

	// TODO: This should come from the Solana chain state which should be enhanced to contain timlock and MCMS address
	for selector, chain := range e.SolChains {
		addresses, _ := addressForChain(e, selector)
		mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
		address, err := mcmCfg.MCMBasedOnActionSolana(*mcmState)
		if err != nil {
			return nil, fmt.Errorf("failed to get mcms for action %s: %w", mcmCfg.MCMSAction, err)
		}
		addressPerChain[selector] = address
	}

	return addressPerChain, nil
}

func addressForChain(e deployment.Environment, selector uint64) (map[string]deployment.TypeAndVersion, error) {
	return e.ExistingAddresses.AddressesForChain(selector) //nolint:staticcheck // Uncomment above once datastore is updated to contains addresses
}
