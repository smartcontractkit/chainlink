package changeset

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

// Given that we interact with native BTC chains, 0.000062 BTC (approximately $5 at 16 Oct 2026) could serve as our minimum test amount.
// 0.000062
var testAmountWei = big.NewInt(62000000000000)

type ccipUnwrapperChangeset struct{}

var CCIPUnwrapperChangeset cldf.ChangeSetV2[vaulttypes.CCIPUnwrapTokenConfig] = ccipUnwrapperChangeset{}

func (c ccipUnwrapperChangeset) VerifyPreconditions(e cldf.Environment, cfg vaulttypes.CCIPUnwrapTokenConfig) error {
	return validateUnwrapChangesetConfig(e, cfg)
}

func (c ccipUnwrapperChangeset) Apply(e cldf.Environment, cfg vaulttypes.CCIPUnwrapTokenConfig) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Unwraping wrapped balances",
		"amountChains", len(cfg.Chains),
	)

	evmChains := e.BlockChains.EVMChains()

	// Pick a deterministic primary chain for deps (only needed for direct execution, not MCMS).
	primaryChainSelector, found := lowestChainSelector(cfg.Chains)
	if !found {
		return cldf.ChangesetOutput{}, errors.New("no chains configured")
	}

	primaryChain := evmChains[primaryChainSelector]

	vaultDeps := VaultDeps{
		Auth:        primaryChain.DeployerKey,
		Chain:       primaryChain,
		Environment: e,
		DataStore:   e.DataStore,
	}

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		UnwrapWrappedTokeAndTransferSequence,
		vaultDeps,
		UnwrapWrappedTokenAndTransferSequenceInput{
			Chains:     cfg.Chains,
			MCMSConfig: cfg.MCMSConfig,
			TestNeeded: cfg.TestEnabled,
		})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute unwrap sequence: %w", err)
	}

	logger.Infow("Unwrap wrapped token proposal generation succeeded",
		"proposals", len(seqReport.Output.MCMSTimelockProposals))

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
		Reports:               seqReport.ExecutionReports,
	}, nil
}

type UnwrapWrappedTokenAndTransferSequenceInput struct {
	Chains     map[uint64]vaulttypes.CCIPUnwrapTokenChainConfig `json:"chains"`
	MCMSConfig *cldfproposalutils.TimelockConfig                `json:"mcms_config"`
	TestNeeded bool                                             `json:"test_needed"`
}

type UnwrapWrappedTokenAndTransferSequenceOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal `json:"mcms_timelock_proposals"`
}

const (
	unwrapProposalDescription = "Vault Unwrap Wrapped Token"
	unwrapLabelTest           = "TEST"
	unwrapLabelRemainder      = "REMAINDER"
	unwrapLabelFull           = "FULL"
)

// Main Sequence
var UnwrapWrappedTokeAndTransferSequence = operations.NewSequence(
	"unwrap-wrapped-token-sequence",
	semver.MustParse("1.0.0"),
	"Unwrap the Timelock wrapped tokens and transfer to a designated address",
	func(b operations.Bundle, deps VaultDeps, input UnwrapWrappedTokenAndTransferSequenceInput) (UnwrapWrappedTokenAndTransferSequenceOutput, error) {
		out := UnwrapWrappedTokenAndTransferSequenceOutput{}

		// Read every chain's wrapped balance once, up front, so the test and remainder
		// proposals are built from the same snapshot and always sum to the total balance.
		balances := make(map[uint64]*big.Int, len(input.Chains))
		for chainSelector := range input.Chains {
			balance, err := readTimelockWrappedBalance(b, deps, chainSelector)
			if err != nil {
				return UnwrapWrappedTokenAndTransferSequenceOutput{}, fmt.Errorf("chain %d: failed to read wrapped balance: %w", chainSelector, err)
			}
			balances[chainSelector] = balance

			// Without a test split there is no other balance check; avoid generating a no-op
			// proposal that withdraws and transfers zero.
			if !input.TestNeeded && balance.Sign() == 0 {
				return UnwrapWrappedTokenAndTransferSequenceOutput{}, fmt.Errorf("chain %d: wrapped balance is zero, nothing to unwrap", chainSelector)
			}
		}

		if input.TestNeeded {
			testAmounts := make(map[uint64]*big.Int, len(balances))
			remainderAmounts := make(map[uint64]*big.Int, len(balances))
			for chainSelector, balance := range balances {
				test, remainder, err := splitUnwrapAmount(balance)
				if err != nil {
					return UnwrapWrappedTokenAndTransferSequenceOutput{}, fmt.Errorf("chain %d: %w", chainSelector, err)
				}
				testAmounts[chainSelector] = test
				remainderAmounts[chainSelector] = remainder
			}

			testProposal, err := buildUnwrapProposal(b, deps, input, testAmounts, unwrapLabelTest)
			if err != nil {
				return UnwrapWrappedTokenAndTransferSequenceOutput{}, fmt.Errorf("test proposal: %w", err)
			}
			remainderProposal, err := buildUnwrapProposal(b, deps, input, remainderAmounts, unwrapLabelRemainder)
			if err != nil {
				return UnwrapWrappedTokenAndTransferSequenceOutput{}, fmt.Errorf("remainder proposal: %w", err)
			}
			out.MCMSTimelockProposals = []mcms.TimelockProposal{*testProposal, *remainderProposal}
			return out, nil
		}

		fullProposal, err := buildUnwrapProposal(b, deps, input, balances, unwrapLabelFull)
		if err != nil {
			return UnwrapWrappedTokenAndTransferSequenceOutput{}, err
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*fullProposal}
		return out, nil
	},
)

// splitUnwrapAmount splits a wrapped-token balance into the test amount and the remainder,
// mirroring splitTransfersForTest in bridge_operations.go. Requires balance > TEST_AMOUNT_WEI
// so the remainder is strictly positive.
func splitUnwrapAmount(balance *big.Int) (test, remainder *big.Int, err error) {
	if balance == nil || balance.Cmp(testAmountWei) <= 0 {
		return nil, nil, fmt.Errorf("wrapped balance %s must be greater than test amount %s", balance, testAmountWei)
	}
	return new(big.Int).Set(testAmountWei), new(big.Int).Sub(balance, testAmountWei), nil
}

// readTimelockWrappedBalance resolves the chain's Timelock and wrapped-native token from the
// datastore and reads the Timelock's wrapped balance on-chain.
func readTimelockWrappedBalance(b operations.Bundle, deps VaultDeps, chainSelector uint64) (*big.Int, error) {
	evmChains := deps.Environment.BlockChains.EVMChains()
	chain, ok := evmChains[chainSelector]
	if !ok {
		return nil, fmt.Errorf("chain not found in environment: %d", chainSelector)
	}

	timelockAddr, err := getRequiredContractAddress(deps.DataStore, chainSelector, commontypes.RBACTimelock)
	if err != nil {
		return nil, fmt.Errorf("timelock not found: %w", err)
	}
	wrappedTokenAddr, err := getRequiredContractAddress(deps.DataStore, chainSelector, "WETH")
	if err != nil {
		return nil, fmt.Errorf("wrapped token address not found in datastore: %w", err)
	}

	return readWETHBalance(b.GetContext(), chain.Client,
		common.HexToAddress(wrappedTokenAddr), common.HexToAddress(timelockAddr))
}

// buildUnwrapProposal assembles the per-chain (withdraw + optional native transfer) batches for
// the given amounts and builds a labeled timelock proposal.
func buildUnwrapProposal(
	b operations.Bundle,
	deps VaultDeps,
	input UnwrapWrappedTokenAndTransferSequenceInput,
	amounts map[uint64]*big.Int,
	label string,
) (*mcms.TimelockProposal, error) {
	var batches []mcmstypes.BatchOperation
	timelockAddresses := make(map[uint64]string)
	mcmAddressByChain := make(map[uint64]string)

	for chainSelector, chainCfg := range input.Chains {
		amount, ok := amounts[chainSelector]
		if !ok {
			return nil, fmt.Errorf("chain %d: no unwrap amount resolved", chainSelector)
		}

		var mcmsTransactions []mcmstypes.Transaction

		unwrapOpReport, err := operations.ExecuteOperation(
			b,
			UnwrapWrappedTokenOperation,
			deps,
			UnwrapWrappedTokenInput{
				ChainSelector: chainSelector,
				MCMSConfig:    input.MCMSConfig,
				Amount:        amount,
			})
		if err != nil {
			return nil, fmt.Errorf("chain %d: unwrap operation failed: %w", chainSelector, err)
		}
		mcmsTransactions = append(mcmsTransactions, unwrapOpReport.Output.MCMSTransaction)

		if chainCfg.ExternalAddressTransfer != "" {
			nativeTxOpReport, err := operations.ExecuteOperation(
				b,
				TransferUnwrappedTokenOperation,
				deps,
				TransferUnwrappedTokenInput{
					ChainSelector:   chainSelector,
					ExternalAddress: chainCfg.ExternalAddressTransfer,
					Amount:          amount,
				},
			)
			if err != nil {
				return nil, fmt.Errorf("chain %d: transfer operation failed: %w", chainSelector, err)
			}
			mcmsTransactions = append(mcmsTransactions, nativeTxOpReport.Output.MCMSTransaction)
		}

		batches = append(batches, mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  mcmsTransactions,
		})
		timelockAddresses[chainSelector] = unwrapOpReport.Output.TimelockAddr
		mcmAddressByChain[chainSelector] = unwrapOpReport.Output.MCMSAddr
	}

	proposal, err := proposeutils.BuildProposalFromBatchesV2(deps.Environment,
		timelockAddresses, mcmAddressByChain, nil, batches,
		fmt.Sprintf("%s [%s]", unwrapProposalDescription, label), *input.MCMSConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build timelock proposal: %w", err)
	}
	return proposal, nil
}

// Operation 1. Unwrap
type UnwrapWrappedTokenInput struct {
	ChainSelector uint64                            `json:"chain_selector"`
	MCMSConfig    *cldfproposalutils.TimelockConfig `json:"mcms_config"`
	Amount        *big.Int                          `json:"amount"`
}

type UnwrapWrappedTokenOutput struct {
	TimelockAddr    string                `json:"timelock_address"`
	MCMSAddr        string                `json:"mcms_address"`
	Amount          *big.Int              `json:"amount"`
	MCMSTransaction mcmstypes.Transaction `json:"mcms_transaction"`
}

var UnwrapWrappedTokenOperation = operations.NewOperation(
	"unwrap-wrapped-token",
	semver.MustParse("1.0.0"),
	"Unwraps wrapped native token held by the chain Timelock (WETH withdraw)",
	func(b operations.Bundle, deps VaultDeps, input UnwrapWrappedTokenInput) (UnwrapWrappedTokenOutput, error) {
		b.Logger.Infow("Starting unwrapping wrapped token operation",
			"chain", input.ChainSelector,
			"amount", input.Amount)
		out := UnwrapWrappedTokenOutput{}

		timelockAddr, err := getRequiredContractAddress(deps.DataStore, input.ChainSelector, commontypes.RBACTimelock)
		if err != nil {
			return UnwrapWrappedTokenOutput{}, fmt.Errorf("timelock not found for chain %d: %w", input.ChainSelector, err)
		}

		mcmType := commontypes.ProposerManyChainMultisig
		if input.MCMSConfig.MCMSAction == mcmstypes.TimelockActionBypass {
			mcmType = commontypes.BypasserManyChainMultisig
		}
		mcmAddr, err := GetContractAddress(deps.DataStore, input.ChainSelector, mcmType)
		if err != nil {
			return UnwrapWrappedTokenOutput{}, fmt.Errorf("%s not found for chain %d: %w", mcmType.String(), input.ChainSelector, err)
		}

		wrappedTokenAddress, err := getRequiredContractAddress(deps.DataStore, input.ChainSelector, "WETH")
		if err != nil {
			return UnwrapWrappedTokenOutput{}, fmt.Errorf("wrapped token address not found for chain %d in datastore: %w", input.ChainSelector, err)
		}

		data, err := encodeErc20WithdrawCallData(input.Amount)
		if err != nil {
			return UnwrapWrappedTokenOutput{}, fmt.Errorf("chain %d failed to encode withdraw calldata: %w", input.ChainSelector, err)
		}

		unwrapTx, err := cldfproposalutils.TransactionForChain(
			input.ChainSelector,
			wrappedTokenAddress,
			data,
			big.NewInt(0),
			"Withdraw Wrapped Token",
			[]string{"vault", "unwrap-token"},
		)
		if err != nil {
			return UnwrapWrappedTokenOutput{}, fmt.Errorf("chain %d failed to create unwrap transaction: %w", input.ChainSelector, err)
		}

		out.TimelockAddr = timelockAddr
		out.MCMSAddr = mcmAddr
		out.MCMSTransaction = unwrapTx
		out.Amount = input.Amount

		return out, nil
	},
)

// Operation 2. Send token to address
type TransferUnwrappedTokenInput struct {
	ChainSelector   uint64   `json:"chain_selector"`
	ExternalAddress string   `json:"external_address"`
	Amount          *big.Int `json:"amount"`
}

type TransferUnwrappedTokenOutput struct {
	MCMSTransaction mcmstypes.Transaction `json:"mcms_transaction"`
}

var TransferUnwrappedTokenOperation = operations.NewOperation(
	"native-transfer",
	semver.MustParse("1.0.0"),
	"Operation to create native transaction after unwrap wrapped token",
	func(b operations.Bundle, deps VaultDeps, input TransferUnwrappedTokenInput) (TransferUnwrappedTokenOutput, error) {
		txn, err := cldfproposalutils.TransactionForChain(
			input.ChainSelector,
			input.ExternalAddress,
			[]byte{},
			input.Amount,
			"NativeTransfer",
			[]string{"vault", "native-transfer"},
		)
		if err != nil {
			return TransferUnwrappedTokenOutput{}, fmt.Errorf("failed to create transaction for chain %d: %w", input.ChainSelector, err)
		}

		return TransferUnwrappedTokenOutput{
			MCMSTransaction: txn,
		}, nil
	},
)

// Utils.
func encodeErc20WithdrawCallData(amount *big.Int) ([]byte, error) {
	return wethABI.Pack("withdraw", amount)
}
