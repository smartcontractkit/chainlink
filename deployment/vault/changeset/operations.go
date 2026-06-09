package changeset

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

type VaultDeps struct {
	Chain       cldf_evm.Chain
	Auth        *bind.TransactOpts
	DataStore   datastore.DataStore
	Environment cldf.Environment
}

// ValidateWhitelistAddressesInput validates that recipient addresses are whitelisted.
type ValidateWhitelistAddressesInput struct {
	ChainSelector uint64   `json:"chain_selector"`
	Addresses     []string `json:"addresses"`
}

// WhitelistValidationOutput contains whitelist validation results.
type WhitelistValidationOutput struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

// FundTimelockInput funds a timelock contract
type FundTimelockInput struct {
	ChainSelector uint64   `json:"chain_selector"`
	Amount        *big.Int `json:"amount"`
}

// FundTimelockOutput contains funding transaction details
type FundTimelockOutput struct {
	ChainSelector   uint64      `json:"chain_selector"`
	TimelockAddress string      `json:"timelock_address"`
	Amount          *big.Int    `json:"amount"`
	TxHash          common.Hash `json:"tx_hash"`
}

// ExecuteNativeTransferInput executes a single native transfer
type ExecuteNativeTransferInput struct {
	ChainSelector uint64               `json:"chain_selector"`
	Transfer      types.NativeTransfer `json:"transfer"`
}

// ExecuteNativeTransferOutput contains transfer transaction details
type ExecuteNativeTransferOutput struct {
	ChainSelector uint64      `json:"chain_selector"`
	To            string      `json:"to"`
	Amount        *big.Int    `json:"amount"`
	TxHash        common.Hash `json:"tx_hash"`
}

var ValidateWhitelistAddressesOp = operations.NewOperation(
	"validate-whitelist-addresses",
	semver.MustParse("1.0.0"),
	"Validates that recipient addresses are whitelisted",
	func(b operations.Bundle, deps VaultDeps, input ValidateWhitelistAddressesInput) (WhitelistValidationOutput, error) {
		b.Logger.Infow("Validating addresses against whitelist",
			"chain", input.ChainSelector,
			"addresses", len(input.Addresses))

		output := WhitelistValidationOutput{Valid: true, Errors: []string{}}

		whitelistMetadata, err := getChainWhitelistMutable(deps.DataStore, input.ChainSelector)
		if err != nil {
			return output, fmt.Errorf("failed to get whitelist for chain %d: %w", input.ChainSelector, err)
		}

		for _, address := range input.Addresses {
			found := false
			addrHex := common.HexToAddress(address).Hex()
			for _, whitelistedAddr := range whitelistMetadata.Addresses {
				if whitelistedAddr.Address == addrHex {
					found = true
					break
				}
			}

			if !found {
				output.Valid = false
				output.Errors = append(output.Errors, fmt.Sprintf("address %s not whitelisted on chain %d", addrHex, input.ChainSelector))
			}
		}

		if !output.Valid {
			return output, fmt.Errorf("validation failed: %v", output.Errors)
		}

		b.Logger.Infow("Whitelist validation completed successfully",
			"chain", input.ChainSelector,
			"addresses", len(input.Addresses))

		return output, nil
	},
)

func executeWhitelistValidation(
	b operations.Bundle,
	deps VaultDeps,
	chainSelector uint64,
	addresses []string,
) (WhitelistValidationOutput, error) {
	report, err := operations.ExecuteOperation(b, ValidateWhitelistAddressesOp, deps, ValidateWhitelistAddressesInput{
		ChainSelector: chainSelector,
		Addresses:     addresses,
	})
	if err != nil {
		return WhitelistValidationOutput{}, err
	}
	return report.Output, nil
}

// FundTimelockOp funds Timelock with native tokens
var FundTimelockOp = operations.NewOperation(
	"fund-timelock",
	semver.MustParse("1.0.0"),
	"Funds Timelock with native tokens",
	func(b operations.Bundle, deps VaultDeps, input FundTimelockInput) (FundTimelockOutput, error) {
		timelockAddr, err := GetContractAddress(deps.DataStore, input.ChainSelector, commontypes.RBACTimelock)
		if err != nil {
			return FundTimelockOutput{}, fmt.Errorf("timelock not found for chain %d: %w", input.ChainSelector, err)
		}

		timelockAddress := common.HexToAddress(timelockAddr)

		b.Logger.Infow("Funding timelock with native tokens",
			"chain", input.ChainSelector,
			"timelock", timelockAddress.Hex(),
			"amount", input.Amount.String(),
			"from", deps.Auth.From.Hex())

		nonce, err := deps.Chain.Client.PendingNonceAt(b.GetContext(), deps.Auth.From)
		if err != nil {
			return FundTimelockOutput{}, fmt.Errorf("failed to get nonce for chain %d: %w", input.ChainSelector, err)
		}

		tx := &ethTypes.DynamicFeeTx{
			Nonce:     nonce,
			To:        &timelockAddress,
			Value:     input.Amount,
			Gas:       50000, // Higher for timelock
			GasFeeCap: big.NewInt(20_000_000_000),
			GasTipCap: big.NewInt(1_000_000_000),
			Data:      nil,
		}

		signedTx, err := deps.Chain.DeployerKey.Signer(deps.Auth.From, ethTypes.NewTx(tx))
		if err != nil {
			return FundTimelockOutput{}, fmt.Errorf("failed to sign funding transaction for chain %d: %w", input.ChainSelector, err)
		}

		err = deps.Chain.Client.SendTransaction(b.GetContext(), signedTx)
		if err != nil {
			return FundTimelockOutput{}, fmt.Errorf("failed to send funding transaction for chain %d: %w", input.ChainSelector, err)
		}

		_, err = deps.Chain.Confirm(signedTx)
		if err != nil {
			return FundTimelockOutput{}, fmt.Errorf("failed to confirm funding transaction for chain %d (tx %s): %w", input.ChainSelector, signedTx.Hash().Hex(), err)
		}

		output := FundTimelockOutput{
			ChainSelector:   input.ChainSelector,
			TimelockAddress: timelockAddress.Hex(),
			Amount:          input.Amount,
			TxHash:          signedTx.Hash(),
		}

		b.Logger.Infow("Timelock funded successfully",
			"chain", input.ChainSelector,
			"timelock", timelockAddress.Hex(),
			"amount", input.Amount.String(),
			"tx", signedTx.Hash().Hex())

		return output, nil
	},
)

var ExecuteNativeTransferOp = operations.NewOperation(
	"execute-native-transfer",
	semver.MustParse("1.0.0"),
	"Executes a single native token transfer",
	func(b operations.Bundle, deps VaultDeps, input ExecuteNativeTransferInput) (ExecuteNativeTransferOutput, error) {
		recipientAddress := common.HexToAddress(input.Transfer.To)

		b.Logger.Infow("Executing native transfer",
			"chain", input.ChainSelector,
			"to", recipientAddress.Hex(),
			"amount", input.Transfer.Amount.String())

		nonce, err := deps.Chain.Client.PendingNonceAt(b.GetContext(), deps.Auth.From)
		if err != nil {
			return ExecuteNativeTransferOutput{}, fmt.Errorf("failed to get nonce for chain %d: %w", input.ChainSelector, err)
		}

		tx := &ethTypes.DynamicFeeTx{
			Nonce:     nonce,
			To:        &recipientAddress,
			Value:     input.Transfer.Amount,
			Gas:       21000,
			GasFeeCap: big.NewInt(20_000_000_000),
			GasTipCap: big.NewInt(1_000_000_000),
			Data:      nil,
		}

		signedTx, err := deps.Chain.DeployerKey.Signer(deps.Auth.From, ethTypes.NewTx(tx))
		if err != nil {
			return ExecuteNativeTransferOutput{}, fmt.Errorf("failed to sign transfer to %s on chain %d: %w", recipientAddress.Hex(), input.ChainSelector, err)
		}

		err = deps.Chain.Client.SendTransaction(b.GetContext(), signedTx)
		if err != nil {
			return ExecuteNativeTransferOutput{}, fmt.Errorf("failed to send transfer to %s on chain %d: %w", recipientAddress.Hex(), input.ChainSelector, err)
		}

		_, err = deps.Chain.Confirm(signedTx)
		if err != nil {
			return ExecuteNativeTransferOutput{}, fmt.Errorf("failed to confirm transfer to %s on chain %d (tx %s): %w", recipientAddress.Hex(), input.ChainSelector, signedTx.Hash().Hex(), err)
		}

		output := ExecuteNativeTransferOutput{
			ChainSelector: input.ChainSelector,
			To:            recipientAddress.Hex(),
			Amount:        input.Transfer.Amount,
			TxHash:        signedTx.Hash(),
		}

		b.Logger.Infow("Native transfer completed",
			"chain", input.ChainSelector,
			"to", recipientAddress.Hex(),
			"amount", input.Transfer.Amount.String(),
			"tx", signedTx.Hash().Hex())

		return output, nil
	},
)

// BatchNativeTransferSequenceInput is the input for the batch transfer sequence
type BatchNativeTransferSequenceInput struct {
	TransfersByChain map[uint64][]types.NativeTransfer `json:"transfers_by_chain"`
	MCMSConfig       *cldfproposalutils.TimelockConfig `json:"mcms_config,omitempty"`
	Description      string                            `json:"description"`
}

type BatchNativeTransferSequenceOutput struct {
	ValidationResults     map[uint64]WhitelistValidationOutput     `json:"validation_results"`
	FundingResults        map[uint64]FundTimelockOutput            `json:"funding_results,omitempty"`
	TransferResults       map[uint64][]ExecuteNativeTransferOutput `json:"transfer_results,omitempty"`
	MCMSTimelockProposals []mcms.TimelockProposal                  `json:"mcms_timelock_proposals,omitempty"`
	Description           string                                   `json:"description"`
}

// BatchNativeTransferSequence executes the batch transfer workflow:
// 1. Validate all transfers against whitelist
// 2. Execute directly OR create MCMS proposals based on MCMSConfig
var BatchNativeTransferSequence = operations.NewSequence(
	"batch-native-transfer-sequence",
	semver.MustParse("1.0.0"),
	"Executes batch native transfers with direct and MCMS execution",
	func(b operations.Bundle, deps VaultDeps, input BatchNativeTransferSequenceInput) (BatchNativeTransferSequenceOutput, error) {
		b.Logger.Infow("Starting batch native transfer sequence",
			"chains", len(input.TransfersByChain),
			"mcms_mode", input.MCMSConfig != nil,
			"description", input.Description)

		output := BatchNativeTransferSequenceOutput{
			ValidationResults: make(map[uint64]WhitelistValidationOutput),
			FundingResults:    make(map[uint64]FundTimelockOutput),
			TransferResults:   make(map[uint64][]ExecuteNativeTransferOutput),
			Description:       input.Description,
		}

		b.Logger.Infow("Validating transfers against whitelist")
		for chainSelector, transfers := range input.TransfersByChain {
			validation, err := executeWhitelistValidation(
				b, deps, chainSelector, recipientAddressesFromNativeTransfers(transfers),
			)
			if err != nil {
				return BatchNativeTransferSequenceOutput{}, fmt.Errorf("validation failed for chain %d: %w", chainSelector, err)
			}

			output.ValidationResults[chainSelector] = validation
		}

		if input.MCMSConfig == nil {
			return executeDirectTransfersOperation(b, deps, input, output)
		}
		return generateMCMSProposals(b, deps, input, output)
	},
)

func executeDirectTransfersOperation(b operations.Bundle, deps VaultDeps, input BatchNativeTransferSequenceInput, output BatchNativeTransferSequenceOutput) (BatchNativeTransferSequenceOutput, error) {
	b.Logger.Infow("Executing native transfers directly")

	evmChains := deps.Environment.BlockChains.EVMChains()

	for chainSelector, transfers := range input.TransfersByChain {
		chain, exists := evmChains[chainSelector]
		if !exists {
			return BatchNativeTransferSequenceOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}

		chainDeps := VaultDeps{
			Chain:       chain,
			Auth:        chain.DeployerKey,
			DataStore:   deps.DataStore,
			Environment: deps.Environment,
		}

		transferResults := make([]ExecuteNativeTransferOutput, 0, len(transfers))

		for i, transfer := range transfers {
			transferInput := ExecuteNativeTransferInput{
				ChainSelector: chainSelector,
				Transfer:      transfer,
			}

			transferReport, err := operations.ExecuteOperation(
				b, ExecuteNativeTransferOp, chainDeps, transferInput,
			)
			if err != nil {
				return BatchNativeTransferSequenceOutput{}, fmt.Errorf("transfer %d failed on chain %d: %w", i, chainSelector, err)
			}

			transferResults = append(transferResults, transferReport.Output)
		}

		output.TransferResults[chainSelector] = transferResults
	}

	b.Logger.Infow("Direct transfer execution completed successfully",
		"chains", len(input.TransfersByChain))

	return output, nil
}

func generateMCMSProposals(b operations.Bundle, deps VaultDeps, input BatchNativeTransferSequenceInput, output BatchNativeTransferSequenceOutput) (BatchNativeTransferSequenceOutput, error) {
	b.Logger.Infow("Generating MCMS timelock proposals")

	var batches []mcmstypes.BatchOperation
	timelockAddressByChain := make(map[uint64]string)
	mcmAddressByChain := make(map[uint64]string)
	inspectorPerChain := make(map[uint64]mcmssdk.Inspector)

	evmChains := deps.Environment.BlockChains.EVMChains()

	for chainSelector, transfers := range input.TransfersByChain {
		chain, exists := evmChains[chainSelector]
		if !exists {
			return BatchNativeTransferSequenceOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}

		timelockAddr, err := GetContractAddress(deps.DataStore, chainSelector, commontypes.RBACTimelock)
		if err != nil {
			return BatchNativeTransferSequenceOutput{}, fmt.Errorf("timelock not found for chain %d: %w", chainSelector, err)
		}

		var mcmAddr string
		var contractName string
		if input.MCMSConfig.MCMSAction == mcmstypes.TimelockActionBypass {
			mcmAddr, err = GetContractAddress(deps.DataStore, chainSelector, commontypes.BypasserManyChainMultisig)
			contractName = "bypasser"
		} else {
			mcmAddr, err = GetContractAddress(deps.DataStore, chainSelector, commontypes.ProposerManyChainMultisig)
			contractName = "proposer"
		}
		if err != nil {
			return BatchNativeTransferSequenceOutput{}, fmt.Errorf("%s not found for chain %d: %w", contractName, chainSelector, err)
		}

		timelockAddressByChain[chainSelector] = timelockAddr
		mcmAddressByChain[chainSelector] = mcmAddr
		inspectorPerChain[chainSelector] = mcmsevmsdk.NewInspector(chain.Client)

		var transactions []mcmstypes.Transaction
		for _, transfer := range transfers {
			tx, err := cldfproposalutils.TransactionForChain(
				chainSelector,
				transfer.To,
				[]byte{},
				transfer.Amount,
				"NativeTransfer",
				[]string{"vault", "native-transfer"},
			)
			if err != nil {
				return BatchNativeTransferSequenceOutput{}, fmt.Errorf("failed to create transaction for chain %d: %w", chainSelector, err)
			}

			transactions = append(transactions, tx)
		}

		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  transactions,
		}
		batches = append(batches, batch)
	}

	description := input.Description
	if description == "" {
		description = "Batch Native Token Transfer"
	}

	proposal, err := proposeutils.BuildProposalFromBatchesV2(
		deps.Environment,
		timelockAddressByChain,
		mcmAddressByChain,
		inspectorPerChain,
		batches,
		description,
		*input.MCMSConfig,
	)
	if err != nil {
		return BatchNativeTransferSequenceOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", err)
	}

	output.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}

	b.Logger.Infow("MCMS proposal generation completed successfully",
		"chains", len(input.TransfersByChain),
		"operations_count", len(proposal.Operations))

	return output, nil
}

var erc20ABI = func() abi.ABI {
	const abiJSON = `[{"name":"transfer","type":"function","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]}]`
	a, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		panic(fmt.Sprintf("failed to parse ERC20 ABI: %v", err))
	}
	return a
}()

func encodeERC20TransferCalldata(tr types.ERC20Transfer) ([]byte, error) {
	return erc20ABI.Pack("transfer", common.HexToAddress(tr.Payee), tr.Amount)
}

// TransferERC20SequenceInput is the input for the ERC20 transfer sequence.
type TransferERC20SequenceInput struct {
	TransfersByChain map[uint64][]types.ERC20Transfer  `json:"transfers_by_chain"`
	MCMSConfig       *cldfproposalutils.TimelockConfig `json:"mcms_config"`
	Description      string                            `json:"description"`
}

// TransferERC20SequenceOutput is the output of the ERC20 transfer sequence.
type TransferERC20SequenceOutput struct {
	ValidationResults     map[uint64]WhitelistValidationOutput `json:"validation_results"`
	MCMSTimelockProposals []mcms.TimelockProposal              `json:"mcms_timelock_proposals,omitempty"`
	Description           string                               `json:"description"`
}

// TransferERC20Sequence validates ERC20 transfers against the whitelist and builds MCMS proposals.
var TransferERC20Sequence = operations.NewSequence(
	"transfer-erc20-sequence",
	semver.MustParse("1.0.0"),
	"Validates ERC20 transfers and builds MCMS timelock proposals",
	func(b operations.Bundle, deps VaultDeps, input TransferERC20SequenceInput) (TransferERC20SequenceOutput, error) {
		b.Logger.Infow("Starting ERC20 transfer sequence",
			"chains", len(input.TransfersByChain),
			"description", input.Description)

		output := TransferERC20SequenceOutput{
			ValidationResults: make(map[uint64]WhitelistValidationOutput),
			Description:       input.Description,
		}

		b.Logger.Infow("Validating ERC20 transfers against whitelist")
		for chainSelector, transfers := range input.TransfersByChain {
			validation, err := executeWhitelistValidation(
				b, deps, chainSelector, recipientAddressesFromERC20Transfers(transfers),
			)
			if err != nil {
				return TransferERC20SequenceOutput{}, fmt.Errorf("validation failed for chain %d: %w", chainSelector, err)
			}
			output.ValidationResults[chainSelector] = validation
		}

		if input.MCMSConfig == nil {
			return TransferERC20SequenceOutput{}, errors.New("MCMSConfig is required for ERC20 transfers")
		}

		return generateERC20MCMSProposals(b, deps, input, output)
	},
)

func generateERC20MCMSProposals(b operations.Bundle, deps VaultDeps, input TransferERC20SequenceInput, output TransferERC20SequenceOutput) (TransferERC20SequenceOutput, error) {
	b.Logger.Infow("Generating ERC20 MCMS timelock proposals",
		"chains", len(input.TransfersByChain))

	var batches []mcmstypes.BatchOperation
	timelockAddressByChain := make(map[uint64]string)
	mcmAddressByChain := make(map[uint64]string)
	inspectorPerChain := make(map[uint64]mcmssdk.Inspector)

	evmChains := deps.Environment.BlockChains.EVMChains()

	for chainSelector, transfers := range input.TransfersByChain {
		chain, exists := evmChains[chainSelector]
		if !exists {
			return TransferERC20SequenceOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}

		timelockAddr, err := GetContractAddressWithQualifier(deps.DataStore, chainSelector, commontypes.RBACTimelock, input.MCMSConfig.TimelockQualifierPerChain[chainSelector])
		if err != nil {
			return TransferERC20SequenceOutput{}, fmt.Errorf("timelock not found for chain %d: %w", chainSelector, err)
		}

		var mcmAddr string
		var contractName string
		if input.MCMSConfig.MCMSAction == mcmstypes.TimelockActionBypass {
			mcmAddr, err = GetContractAddressWithQualifier(deps.DataStore, chainSelector, commontypes.BypasserManyChainMultisig, input.MCMSConfig.TimelockQualifierPerChain[chainSelector])
			contractName = "bypasser"
		} else {
			mcmAddr, err = GetContractAddressWithQualifier(deps.DataStore, chainSelector, commontypes.ProposerManyChainMultisig, input.MCMSConfig.TimelockQualifierPerChain[chainSelector])
			contractName = "proposer"
		}
		if err != nil {
			return TransferERC20SequenceOutput{}, fmt.Errorf("%s not found for chain %d: %w", contractName, chainSelector, err)
		}

		timelockAddressByChain[chainSelector] = timelockAddr
		mcmAddressByChain[chainSelector] = mcmAddr
		inspectorPerChain[chainSelector] = mcmsevmsdk.NewInspector(chain.Client)

		var transactions []mcmstypes.Transaction
		for i, tr := range transfers {
			data, err := encodeERC20TransferCalldata(tr)
			if err != nil {
				return TransferERC20SequenceOutput{}, fmt.Errorf("chain %d transfer %d: failed to encode calldata: %w", chainSelector, i, err)
			}

			tx, err := cldfproposalutils.TransactionForChain(
				chainSelector,
				tr.Token,
				data,
				big.NewInt(0),
				"ERC20Transfer",
				[]string{"vault", "erc20-transfer"},
			)
			if err != nil {
				return TransferERC20SequenceOutput{}, fmt.Errorf("chain %d transfer %d: %w", chainSelector, i, err)
			}
			transactions = append(transactions, tx)
		}

		batches = append(batches, mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  transactions,
		})
	}

	description := input.Description
	if description == "" {
		description = "Vault ERC20 Transfer"
	}

	proposal, err := proposeutils.BuildProposalFromBatchesV2(
		deps.Environment,
		timelockAddressByChain,
		mcmAddressByChain,
		inspectorPerChain,
		batches,
		description,
		*input.MCMSConfig,
	)
	if err != nil {
		return TransferERC20SequenceOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", err)
	}

	output.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}

	b.Logger.Infow("ERC20 MCMS proposal generation completed successfully",
		"chains", len(input.TransfersByChain),
		"operations_count", len(proposal.Operations))

	return output, nil
}
