package changeset

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

type CCIPUnwrapTokenChainConfig struct {
	ExternalAddressTransfer string `json:"external_address_transfer,omitempty"`
}

// types
type CCIPUnwrapTokenConfig struct {
	Chains      map[uint64]CCIPUnwrapTokenChainConfig `json:"chains"`
	TestEnabled bool                                  `json:"test_enabled"`

	MCMSConfig *cldfproposalutils.TimelockConfig `json:"mcms_config"`
}

// Given that we interact with native BTC chains, 0.000062 BTC (approximately $5 at 16 Oct 2026) could serve as our minimum test amount.
// 0.000062
var TEST_AMOUNT_WEI = big.NewInt(62000000000000)

type ccipFeeCollectorChangeset struct{}

var CCIPFeeCollector cldf.ChangeSetV2[CCIPUnwrapTokenConfig] = ccipFeeCollectorChangeset{}

func (c ccipFeeCollectorChangeset) VerifyPreconditions(e cldf.Environment, cfg CCIPUnwrapTokenConfig) error {
	if cfg.MCMSConfig == nil {
		return fmt.Errorf("MCMS Config cannot be nil")
	}

	for chainSelector, chainConfig := range cfg.Chains {
		address := chainConfig.ExternalAddressTransfer
		if address == "" {
			continue
		}

		if !common.IsHexAddress(address) {
			return fmt.Errorf(
				"invalid transfer address %q for chain %d",
				address, chainSelector,
			)
		}

		recipient := common.HexToAddress(address)
		if recipient == (common.Address{}) {
			return fmt.Errorf(
				"transfer address cannot be zero for chain %d",
				chainSelector,
			)
		}

		whitelist, err := GetWhitelistedAddresses(
			e,
			[]uint64{chainSelector},
		)
		if err != nil {
			return fmt.Errorf(
				"failed to read whitelist for chain %d: %w",
				chainSelector, err,
			)
		}

		found := false
		for _, entry := range whitelist[chainSelector] {
			if common.IsHexAddress(entry.Address) &&
				common.HexToAddress(entry.Address) == recipient {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf(
				"address %s is not whitelisted for chain %d",
				recipient.Hex(), chainSelector,
			)
		}
	}

	return nil
}

// OKR
// Withdrawing and converting CCIP revenue across 80 blockchains currently relies on heavy manual effort.
// This initiative would use Chainlink Deployments (CLD)
// 1) to provide the Treasury team a self-service option to pull revenue from CCIP v1.6 and v2.0 onRamps
// 2) Create changeset for CLD to automatically unwrap collected fees (which is by default in wrapped native) into
// native gas tokens across all 80 CCIP chains and transfer them to external addresses for mainnet bridging and liquidations
// - eliminating manual processing time.
func (c ccipFeeCollectorChangeset) Apply(e cldf.Environment, cfg CCIPUnwrapTokenConfig) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Unwraping wrapped balances",
		"amountChains", len(cfg.Chains),
	)

	evmChains := e.BlockChains.EVMChains()

	// Pick a deterministic primary chain for deps (only needed for direct execution, not MCMS).
	var (
		primaryChainSelector uint64
		primaryChainSet      bool
	)
	for chainSelector := range cfg.Chains {
		if !primaryChainSet || chainSelector < primaryChainSelector {
			primaryChainSelector = chainSelector
			primaryChainSet = true
		}
	}

	primaryChain := evmChains[primaryChainSelector]

	vaultDeps := VaultDeps{
		Auth:        primaryChain.DeployerKey,
		Chain:       primaryChain,
		Environment: e,
		DataStore:   e.DataStore,
	}

	var proposals []mcms.TimelockProposal
	var reports []operations.Report[any, any]

	if cfg.TestEnabled {
		testSeqReport, err := operations.ExecuteSequence(
			e.OperationsBundle,
			UnwrapWrappedTokeAndTransferSequence,
			vaultDeps,
			UnwrapWrappedTokenAndTransferSequenceInput{
				Chains:         cfg.Chains,
				MCMSConfig:     cfg.MCMSConfig,
				IsTestSequence: true,
			})

		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute CCIP bridge sequence: %w", err)
		}

		testProposal := testSeqReport.Output.MCMSTimelockProposals
		testReport := testSeqReport.ExecutionReports

		proposals = append(proposals, testProposal)
		reports = append(reports, testReport...)
	}

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		UnwrapWrappedTokeAndTransferSequence,
		vaultDeps,
		UnwrapWrappedTokenAndTransferSequenceInput{
			Chains:         cfg.Chains,
			MCMSConfig:     cfg.MCMSConfig,
			IsTestEnabled:  cfg.TestEnabled,
			IsTestSequence: false,
		})

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute CCIP bridge sequence: %w", err)
	}
	proposals = append(proposals, seqReport.Output.MCMSTimelockProposals)
	reports = append(reports, seqReport.ExecutionReports...)

	logger.Infow("Token wrapping proposal generation succeed",
		"proposals", len(proposals))

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: proposals,
		Reports:               reports,
	}, nil
}

type UnwrapWrappedTokenAndTransferSequenceInput struct {
	Chains         map[uint64]CCIPUnwrapTokenChainConfig `json:"chains"`
	MCMSConfig     *cldfproposalutils.TimelockConfig     `json:"mcms_config"`
	IsTestEnabled  bool                                  `json:"is_test_enabled"`
	IsTestSequence bool                                  `json:"is_test_sequence"`
}

type UnwrapWrappedTokenAndTransferSequenceOutput struct {
	MCMSTimelockProposals mcms.TimelockProposal `json:"mcms_timelock_proposals"`
}

// Main Sequence
// Must return the generated proposals
// Operation 1. Unwrap
// Operation 2. Send to x address
var UnwrapWrappedTokeAndTransferSequence = operations.NewSequence(
	"unwrap-wrapped-token-sequence",
	semver.MustParse("1.0.0"),
	"Unwrap the Timelock wrapped tokens and transfer to a designated address",
	func(b operations.Bundle, deps VaultDeps, input UnwrapWrappedTokenAndTransferSequenceInput) (UnwrapWrappedTokenAndTransferSequenceOutput, error) {
		out := UnwrapWrappedTokenAndTransferSequenceOutput{}
		evmChains := deps.Environment.BlockChains.EVMChains()

		var batches []mcmstypes.BatchOperation

		timelockAddresses := make(map[uint64]string)
		mcmAddressByChain := make(map[uint64]string)

		for chainSelector, chainCfg := range input.Chains {
			var mcmsTransactions []mcmstypes.Transaction
			_, ok := evmChains[chainSelector]
			if !ok {
				return UnwrapWrappedTokenAndTransferSequenceOutput{}, fmt.Errorf("chain not found in environment: %d", chainSelector)
			}

			unwrapOpReport, err := operations.ExecuteOperation(
				b,
				UnwrapWrappedTokenOperation,
				deps,
				UnwrapWrappedTokenInput{
					ChainSelector:   chainSelector,
					MCMSConfig:      input.MCMSConfig,
					ExternalAddress: chainCfg.ExternalAddressTransfer,
					IsTestEnabled:   input.IsTestEnabled,
					IsTestOperation: input.IsTestSequence,
				})
			if err != nil {
				return UnwrapWrappedTokenAndTransferSequenceOutput{}, fmt.Errorf("chain %d: unwrap operation failed: %w", chainSelector, err)
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
						Amount:          unwrapOpReport.Output.Amount,
					},
				)
				if err != nil {
					return UnwrapWrappedTokenAndTransferSequenceOutput{}, fmt.Errorf("chain %d: transfer operation failed: %w", chainSelector, err)
				}
				mcmsTransactions = append(mcmsTransactions, nativeTxOpReport.Output.MCMSTransaction)
			}

			batches = append(batches, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(chainSelector),
				Transactions:  mcmsTransactions,
			})
			mcmsTransactions = append(mcmsTransactions, unwrapOpReport.Output.MCMSTransaction)
			timelockAddresses[chainSelector] = unwrapOpReport.Output.TimelockAddr
			mcmAddressByChain[chainSelector] = unwrapOpReport.Output.MCMSAddr
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(deps.Environment,
			timelockAddresses, mcmAddressByChain, nil, batches, "TODO", *input.MCMSConfig,
		)
		if err != nil {
			return UnwrapWrappedTokenAndTransferSequenceOutput{}, fmt.Errorf("failed to build timelock proposal: %w", err)
		}

		out.MCMSTimelockProposals = *proposal

		return out, nil
	},
)

// Operation 1. Unwrap
// Map the input config
// Create batches
// Create proposal
type UnwrapWrappedTokenInput struct {
	ChainSelector   uint64                            `json:"chain_selector"`
	ExternalAddress string                            `json:"external_address"`
	MCMSConfig      *cldfproposalutils.TimelockConfig `json:"mcms_config"`
	IsTestEnabled   bool                              `json:"is_test_enabled"`
	IsTestOperation bool                              `json:"is_test_operation"`
}

type UnwrapWrappedTokenOutput struct {
	TimelockAddr    string
	MCMSAddr        string
	Amount          *big.Int
	MCMSTransaction mcmstypes.Transaction
}

var UnwrapWrappedTokenOperation = operations.NewOperation(
	"TODO",
	semver.MustParse("1.0.0"),
	"TODO",
	func(b operations.Bundle, deps VaultDeps, input UnwrapWrappedTokenInput) (UnwrapWrappedTokenOutput, error) {
		b.Logger.Infow("Starting Unwrapping native token operation",
			"chain", input.ChainSelector)
		out := UnwrapWrappedTokenOutput{}
		evmChains := deps.Environment.BlockChains.EVMChains()

		// Get the address from datastore based on the chain
		_, exists := evmChains[input.ChainSelector]
		if !exists {
			return UnwrapWrappedTokenOutput{}, fmt.Errorf("chain %d not found in environmnent, please check on chain config yaml", input.ChainSelector)
		}

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

		var wrappedAmount *big.Int
		if input.IsTestOperation {
			wrappedAmount = TEST_AMOUNT_WEI
		} else {
			var err error
			wrappedAmount, err = getAddressTokenBalance(timelockAddr, wrappedTokenAddress)
			if err != nil {
				return UnwrapWrappedTokenOutput{}, fmt.Errorf("chain %d failed to get Timelock balances: %w", input.ChainSelector, err)
			}

			if input.IsTestEnabled {
				if wrappedAmount.Cmp(TEST_AMOUNT_WEI) < 0 || wrappedAmount.Cmp(TEST_AMOUNT_WEI) == 0 {
					return UnwrapWrappedTokenOutput{}, fmt.Errorf("chain %d: balance is lower than TEST_AMOUNT_WEI", input.ChainSelector)
				}

				wrappedAmount = new(big.Int).Sub(wrappedAmount, TEST_AMOUNT_WEI)
			}
		}

		data, err := encodeErc20WithdrawCallData(wrappedAmount)
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
		out.Amount = wrappedAmount

		return out, nil
	},
)

// Operation 2. Send token to address
type TransferUnwrappedTokenInput struct {
	ChainSelector   uint64 `json:"chain_selector"`
	ExternalAddress string
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
	return erc20ABI.Pack("withdraw", amount)
}

// TODO: fetch real balances on the timelock
func getAddressTokenBalance(address string, tokenAddress string) (*big.Int, error) {
	return big.NewInt(1), nil
}
