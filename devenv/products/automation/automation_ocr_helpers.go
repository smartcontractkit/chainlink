package automation

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	pkg_errors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
)

// DeployLegacyConsumers deploys and registers keeper consumers. If ephemeral addresses are enabled, it will deploy and register the consumers from ephemeral addresses, but each upkpeep will be registered with root key address as the admin. Which means
// that functions like setting upkeep configuration, pausing, unpausing, etc. will be done by the root key address. It deploys multicall contract and sends link funds to each deployment address.
func DeployLegacyConsumers(t *testing.T, chainClient *seth.Client, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, linkToken contracts.LinkToken, numberOfUpkeeps int, linkFundsForEachUpkeep *big.Int, upkeepGasLimit uint32, isLogTrigger bool, isMercury bool, isBillingTokenNative bool, wethToken contracts.WETHToken) ([]contracts.KeeperConsumer, []*big.Int) {
	// Fund deployers with LINK, no need to do this for Native token
	if !isBillingTokenNative {
		err := DeployMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfUpkeeps, linkFundsForEachUpkeep)
		require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")
	}

	upkeeps := DeployKeeperConsumers(t, chainClient, numberOfUpkeeps, isLogTrigger, isMercury)
	require.Len(t, upkeeps, numberOfUpkeeps, "Number of upkeeps should match")
	upkeepsAddresses := []string{}
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIDs := RegisterUpkeepContracts(
		t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, isLogTrigger, isMercury, isBillingTokenNative, wethToken,
	)
	require.Len(t, upkeepIDs, numberOfUpkeeps, "Number of upkeepIds should match")
	return upkeeps, upkeepIDs
}

// DeployConsumers deploys and registers keeper consumers. If ephemeral addresses are enabled, it will deploy and register the consumers from ephemeral addresses, but each upkpeep will be registered with root key address as the admin. Which means
// that functions like setting upkeep configuration, pausing, unpausing, etc. will be done by the root key address. It deploys multicall contract and sends link funds to each deployment address.
func DeployConsumers(t *testing.T, chainClient *seth.Client, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, linkToken contracts.LinkToken, numberOfUpkeeps int, linkFundsForEachUpkeep *big.Int, upkeepGasLimit uint32, isLogTrigger bool, isMercury bool, isBillingTokenNative bool, wethToken contracts.WETHToken, config Automation) ([]contracts.KeeperConsumer, []*big.Int) {
	// Fund deployers with LINK, no need to do this for Native token
	if !isBillingTokenNative {
		err := SetupMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfUpkeeps, linkFundsForEachUpkeep, config)
		require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")
	}

	upkeeps := SetupKeeperConsumers(t, chainClient, numberOfUpkeeps, isLogTrigger, isMercury, config)
	require.Len(t, upkeeps, numberOfUpkeeps, "Number of upkeeps should match")
	upkeepsAddresses := []string{}
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIDs := RegisterUpkeepContracts(
		t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, isLogTrigger, isMercury, isBillingTokenNative, wethToken,
	)
	require.Len(t, upkeepIDs, numberOfUpkeeps, "Number of upkeepIDs should match")
	return upkeeps, upkeepIDs
}

func SetupMultiCallAddress(chainClient *seth.Client, config Automation) (common.Address, error) {
	if config.DeployedContracts.MultiCall != "" {
		return common.HexToAddress(config.DeployedContracts.MultiCall), nil
	}

	multicallAddress, err := contracts.DeployMultiCallContract(chainClient)
	if err != nil {
		return common.Address{}, pkg_errors.Wrap(err, "Error deploying multicall contract")
	}
	return multicallAddress, nil
}

// SetupMultiCallAndFundDeploymentAddresses setups multicall contract and sends link funds to each deployment address
func SetupMultiCallAndFundDeploymentAddresses(
	chainClient *seth.Client,
	linkToken contracts.LinkToken,
	numberOfUpkeeps int,
	linkFundsForEachUpkeep *big.Int,
	config Automation,
) error {
	concurrency, err := GetAndAssertCorrectConcurrency(chainClient, 1)
	if err != nil {
		return err
	}

	operationsPerAddress := numberOfUpkeeps / concurrency

	multicallAddress, err := SetupMultiCallAddress(chainClient, config)
	if err != nil {
		return pkg_errors.Wrap(err, "Error deploying multicall contract")
	}

	return SendLinkFundsToDeploymentAddresses(chainClient, concurrency, numberOfUpkeeps, operationsPerAddress, multicallAddress, linkFundsForEachUpkeep, linkToken)
}

// DeployMultiCallAndFundDeploymentAddresses deploys multicall contract and sends link funds to each deployment address
func DeployMultiCallAndFundDeploymentAddresses(
	chainClient *seth.Client,
	linkToken contracts.LinkToken,
	numberOfUpkeeps int,
	linkFundsForEachUpkeep *big.Int,
) error {
	concurrency, err := GetAndAssertCorrectConcurrency(chainClient, 1)
	if err != nil {
		return err
	}

	operationsPerAddress := numberOfUpkeeps / concurrency

	multicallAddress, err := contracts.DeployMultiCallContract(chainClient)
	if err != nil {
		return pkg_errors.Wrap(err, "Error deploying multicall contract")
	}

	return SendLinkFundsToDeploymentAddresses(chainClient, concurrency, numberOfUpkeeps, operationsPerAddress, multicallAddress, linkFundsForEachUpkeep, linkToken)
}

// SendLinkFundsToDeploymentAddresses sends LINK token to all addresses, but the root one, from the root address. It uses
// Multicall contract to batch all transfers in a single transaction. It also checks if the funds were transferred correctly.
// It's primary use case is to fund addresses that will be used for Upkeep registration (as that requires LINK balance) during
// Automation/Keeper test setup.
func SendLinkFundsToDeploymentAddresses(
	chainClient *seth.Client,
	concurrency,
	totalUpkeeps,
	operationsPerAddress int,
	multicallAddress common.Address,
	linkAmountPerUpkeep *big.Int,
	linkToken contracts.LinkToken,
) error {
	const maxBatchSize = 75 // keep multicall tx gas comfortably below the block limit
	var generateCallData = func(receiver common.Address, amount *big.Int) ([]byte, error) {
		abi, err := link_token_interface.LinkTokenMetaData.GetAbi()
		if err != nil {
			return nil, err
		}
		data, err := abi.Pack("transfer", receiver, amount)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	toTransferToMultiCallContract := big.NewInt(0).Mul(linkAmountPerUpkeep, big.NewInt(int64(totalUpkeeps+concurrency)))
	toTransferPerClient := big.NewInt(0).Mul(linkAmountPerUpkeep, big.NewInt(int64(operationsPerAddress+1)))

	// As a hack we use the geth wrapper directly, because we need to access receipt to get block number, which we will use to query the balance
	// This is needed as querying with 'latest' block number very rarely, but still, return stale balance. That's happening even though we wait for
	// the transaction to be mined.
	linkInstance, err := link_token_interface.NewLinkToken(common.HexToAddress(linkToken.Address()), contracts.MustNewWrappedContractBackend(nil, chainClient))
	if err != nil {
		return err
	}
	tx, err := chainClient.Decode(linkInstance.Transfer(chainClient.NewTXOpts(), multicallAddress, toTransferToMultiCallContract))
	if err != nil {
		return err
	}

	if tx.Receipt == nil {
		return pkg_errors.New("transaction receipt for LINK transfer to multicall contract is nil")
	}

	multiBalance, err := linkInstance.BalanceOf(&bind.CallOpts{From: chainClient.Addresses[0], BlockNumber: tx.Receipt.BlockNumber}, multicallAddress)
	if err != nil {
		return pkg_errors.Wrapf(err, "error getting LINK balance of multicall contract")
	}

	// Old code that's querying latest block
	// err := linkToken.Transfer(multicallAddress.Hex(), toTransferToMultiCallContract)
	// if err != nil {
	//	return errors.Wrapf(err, "Error transferring LINK to multicall contract")
	//}
	//
	// balance, err := linkToken.BalanceOf(context.Background(), multicallAddress.Hex())
	// if err != nil {
	//	return errors.Wrapf(err, "Error getting LINK balance of multicall contract")
	//}

	if multiBalance.Cmp(toTransferToMultiCallContract) < 0 {
		return fmt.Errorf("incorrect LINK balance of multicall contract. Expected at least: %s. Got: %s", toTransferToMultiCallContract.String(), multiBalance.String())
	}

	// Transfer LINK to ephemeral keys
	multiCallData := make([][]byte, 0)
	for i := 1; i <= concurrency; i++ {
		data, err := generateCallData(chainClient.Addresses[i], toTransferPerClient)
		if err != nil {
			return pkg_errors.Wrapf(err, "error generating call data for LINK transfer")
		}
		multiCallData = append(multiCallData, data)
	}

	call := []contracts.Call{}
	for _, d := range multiCallData {
		data := contracts.Call{Target: common.HexToAddress(linkToken.Address()), AllowFailure: false, CallData: d}
		call = append(call, data)
	}

	multiCallABI, err := abi.JSON(strings.NewReader(contracts.MultiCallABI))
	if err != nil {
		return pkg_errors.Wrapf(err, "error getting Multicall contract ABI")
	}
	boundContract := bind.NewBoundContract(multicallAddress, multiCallABI, chainClient.Client, chainClient.Client, chainClient.Client)
	var lastReceipt *gethtypes.Receipt
	for start := 0; start < len(call); start += maxBatchSize {
		end := start + maxBatchSize
		if end > len(call) {
			end = len(call)
		}
		chunk := make([]contracts.Call, end-start)
		copy(chunk, call[start:end])
		// call aggregate3 to group a safe number of transfers per transaction
		ephemeralTx, err := chainClient.Decode(boundContract.Transact(chainClient.NewTXOpts(), "aggregate3", chunk))
		if err != nil {
			return pkg_errors.Wrapf(err, "error calling Multicall contract")
		}
		if ephemeralTx.Receipt == nil {
			return pkg_errors.New("transaction receipt for LINK transfer to ephemeral keys is nil")
		}
		lastReceipt = ephemeralTx.Receipt
	}

	if lastReceipt == nil {
		return pkg_errors.New("multicall transfer batch did not execute")
	}

	for i := 1; i <= concurrency; i++ {
		ephemeralBalance, err := linkInstance.BalanceOf(&bind.CallOpts{From: chainClient.Addresses[0], BlockNumber: lastReceipt.BlockNumber}, chainClient.Addresses[i])
		// Old code that's querying latest block, for now we prefer to use block number from the transaction receipt
		// balance, err := linkToken.BalanceOf(context.Background(), chainClient.Addresses[i].Hex())
		if err != nil {
			return pkg_errors.Wrapf(err, "error getting LINK balance of ephemeral key %d", i)
		}
		if ephemeralBalance.Cmp(toTransferPerClient) < 0 {
			return fmt.Errorf("incorrect LINK balance after transfer. Ephemeral key %d. Expected: %s. Got: %s", i, toTransferPerClient.String(), ephemeralBalance.String())
		}
	}

	return nil
}
