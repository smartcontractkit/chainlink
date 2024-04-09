package automationv2_1

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/stretchr/testify/require"
)

func TestMulti(t *testing.T) {
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Load", tc.Automation)
	if err != nil {
		t.Fatal(err)
	}

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(network).
		WithMockAdapter().
		// WithCLNodes(6).
		// WithFunding(big.NewFloat(.5)).
		WithStandardCleanup().
		WithSeth().
		Build()
	require.NoError(t, err)

	chainClient, err := env.GetSethClient(1337)
	require.NoError(t, err, "Error getting seth client")

	multicallAddress, err := contracts.DeployMultiCallContract(chainClient)
	require.NoError(t, err, "Error deploying multicall contract")

	linkToken, err := contracts.DeployLinkTokenContract(l, chainClient)
	require.NoError(t, err, "Error deploying link token contract")

	automationDefaultLinkFunds := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(10000))) //10000 LINK

	numberOfClients := int(*chainClient.Cfg.EphemeralAddrs)
	err = linkToken.Transfer(multicallAddress.Hex(), big.NewInt(0).Mul(automationDefaultLinkFunds, big.NewInt(int64(numberOfClients))))
	require.NoError(t, err, "Error transferring LINK to multicall contract")

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

	// Transfer LINK to ephemeral keys
	multiCallData := make([][]byte, 0)
	for i := 1; i <= numberOfClients; i++ {
		data, err := generateCallData(chainClient.Addresses[i], automationDefaultLinkFunds)
		require.NoError(t, err, "Error generating call data for LINK transfer")
		multiCallData = append(multiCallData, data)
	}

	var call []contracts.Call
	for _, d := range multiCallData {
		data := contracts.Call{Target: common.HexToAddress(linkToken.Address()), AllowFailure: false, CallData: d}
		call = append(call, data)
	}

	multiCallABI, err := abi.JSON(strings.NewReader(contracts.MultiCallABI))
	require.NoError(t, err, "Error getting multicall abi")
	boundContract := bind.NewBoundContract(multicallAddress, multiCallABI, chainClient.Client, chainClient.Client, chainClient.Client)
	// call aggregate3 to group all msg call data and send them in a single transaction
	_, err = chainClient.Decode(boundContract.Transact(chainClient.NewTXOpts(), "aggregate3", call))
	require.NoError(t, err, "Error transferring LINK to ephemeral keys")
	l.Info().Msg("Transferred LINK to ephemeral keys")

	for i := 1; i <= numberOfClients; i++ {
		balance, err := linkToken.BalanceOf(context.Background(), chainClient.Addresses[i].Hex())
		require.NoError(t, err, "Error getting balance of ephemeral key")
		require.Equal(t, automationDefaultLinkFunds, balance, "Expected balance to be %v but got %v", automationDefaultLinkFunds, balance)
	}
}
