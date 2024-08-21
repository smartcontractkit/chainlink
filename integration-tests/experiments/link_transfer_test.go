package experiments

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestLinkTransfer(t *testing.T) {
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig([]string{"Smoke"}, tc.Keeper)
	require.NoError(t, err, "Failed to get config")

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "Error deploying test environment")

	evmNetwork, err := env.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	var one int64 = 1
	config.Seth.EphemeralAddrs = &one

	sethClient, err := utils.TestAwareSethClient(t, config, evmNetwork)
	require.NoError(t, err, "Error getting seth client")

	linkTokenContract, err := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	toTransfer := big.NewInt(100_000_000)
	transferredSoFar := big.NewInt(0)

	for i := 0; i < 1500; i++ {
		l.Warn().Msgf("Transfer attempt: %d", i+1)
		transferredSoFar = transferredSoFar.Add(transferredSoFar, toTransfer)
		err = linkTokenContract.Transfer(sethClient.Addresses[1].Hex(), toTransfer)
		require.NoError(t, err, "Error transferring LINK")

		balance, err := linkTokenContract.BalanceOf(context.Background(), sethClient.Addresses[1].Hex())
		require.NoError(t, err, "Error getting LINK balance")

		if balance.Cmp(transferredSoFar) < 0 {
			require.True(t, false, "Incorrect LINK balance. Expected at least: %s. Got: %s", transferredSoFar.String(), balance.String())
		}
	}
}

func TestLinkTransfer_NoSeth(t *testing.T) {
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig([]string{"Smoke"}, tc.Keeper)
	require.NoError(t, err, "Failed to get config")

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "Error deploying test environment")

	evmNetwork, err := env.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	var one int64 = 1
	config.Seth.EphemeralAddrs = &one

	sethClient, err := utils.TestAwareSethClient(t, config, evmNetwork)
	require.NoError(t, err, "Error getting seth client")

	rpcClient, err := rpc.DialOptions(context.Background(), sethClient.Cfg.FirstNetworkURL())
	require.NoError(t, err, "Error dialing RPC client")
	client := ethclient.NewClient(rpcClient)

	privateKey, err := crypto.HexToECDSA(sethClient.Cfg.Network.PrivateKeys[0])
	require.NoError(t, err, "Error getting private key")

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err, "Error creating transactor")

	gasPrice, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err, "Error getting gas price")

	opts.GasPrice = gasPrice

	require.NoError(t, err, "Error creating transactor")
	_, tx, linkTokenContract, err := link_token_interface.DeployLinkToken(opts, client)

	_, err = bind.WaitDeployed(context.Background(), client, tx)
	require.NoError(t, err, "Error waiting for contract to be deployed")

	toTransfer := big.NewInt(100_000_000)
	transferredSoFar := big.NewInt(0)

	for i := 0; i < 1500; i++ {
		l.Warn().Msgf("Transfer attempt: %d", i+1)
		l.Info().Msgf("Transferring %s LINK", toTransfer.String())
		transferredSoFar = transferredSoFar.Add(transferredSoFar, toTransfer)

		opts, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
		require.NoError(t, err, "Error creating transactor")

		gasPrice, err = client.SuggestGasPrice(context.Background())
		require.NoError(t, err, "Error getting gas price")

		opts.GasPrice = gasPrice

		tx, err := linkTokenContract.Transfer(opts, sethClient.Addresses[1], toTransfer)
		require.NoError(t, err, "Error transferring LINK")

		receipt, err := bind.WaitMined(context.Background(), client, tx)
		require.NoError(t, err, "Error waiting for transaction to be mined")
		require.Equal(t, uint64(1), receipt.Status, "Transaction failed")

		callOpts := bind.CallOpts{Context: context.Background(), From: sethClient.Addresses[0]}

		balance, err := linkTokenContract.BalanceOf(&callOpts, sethClient.Addresses[1])
		require.NoError(t, err, "Error getting LINK balance")

		if balance.Cmp(transferredSoFar) < 0 {
			require.True(t, false, "Incorrect LINK balance. Expected at least: %s. Got: %s", transferredSoFar.String(), balance.String())
		} else {
			l.Info().Msgf("LINK Balance: %s", balance.String())
		}
	}
}
