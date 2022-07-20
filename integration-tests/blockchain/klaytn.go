package blockchain

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-env/environment"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// Handles specific issues with the Klaytn EVM chain: https://docs.klaytn.com/

// KlaytnMultinodeClient represents a multi-node, EVM compatible client for the Klaytn network
type KlaytnMultinodeClient struct {
	*EthereumMultinodeClient
}

// KlaytnClient represents a single node, EVM compatible client for the Klaytn network
type KlaytnClient struct {
	*EthereumClient
}

// NewKlaytnClient returns an instantiated instance of the Klaytn client that has connected to the server
func NewKlaytnClient(networkSettings *EVMNetwork) (EVMClient, error) {
	client, err := NewEthereumClient(networkSettings)
	if err != nil {
		return nil, err
	}
	log.Info().Str("Network Name", client.GetNetworkName()).Msg("Using custom Klaytn client")
	return &KlaytnClient{client.(*EthereumClient)}, err
}

func NewKlaytnMultiNodeClientSetup(networkSettings *EVMNetwork) func(*environment.Environment) (EVMClient, error) {
	return func(env *environment.Environment) (EVMClient, error) {
		multiNodeClient := &EthereumMultinodeClient{}
		networkSettings.URLs = append(networkSettings.URLs, env.URLs[networkSettings.Name]...)
		for idx, networkURL := range networkSettings.URLs {
			networkSettings.URL = networkURL
			ec, err := NewKlaytnClient(networkSettings)
			if err != nil {
				return nil, err
			}
			ec.SetID(idx)
			multiNodeClient.Clients = append(multiNodeClient.Clients, ec)
		}
		multiNodeClient.DefaultClient = multiNodeClient.Clients[0]
		log.Info().
			Interface("URLs", networkSettings.URLs).
			Msg("Connected multi-node client")
		return &KlaytnMultinodeClient{multiNodeClient}, nil
	}
}

// Fund overrides ethereum's fund to account for Klaytn's gas specifications
// https://docs.klaytn.com/klaytn/design/transaction-fees#unit-price
func (k *KlaytnClient) Fund(
	toAddress string,
	amount *big.Float,
) error {
	privateKey, err := crypto.HexToECDSA(k.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return err
	}
	// Don't bump gas for Klaytn
	gasPrice, err := k.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	nonce, err := k.GetNonce(context.Background(), k.DefaultWallet.address)
	if err != nil {
		return err
	}
	log.Warn().
		Str("Network Name", k.NetworkConfig.Name).
		Msg("Setting GasTipCap = SuggestedGasPrice for Klaytn network")
	// https://docs.klaytn.com/klaytn/design/transaction-fees#gas
	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(k.GetChainID()), &types.DynamicFeeTx{
		ChainID:   k.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     utils.EtherToWei(amount),
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       22000,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "KLAY").
		Str("From", k.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Funding Address")
	if err := k.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}
	return k.ProcessTransaction(tx)
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (k *KlaytnClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := k.TransactionOpts(k.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}

	// Don't bump gas for Klaytn
	// https://docs.klaytn.com/klaytn/design/transaction-fees#unit-price
	log.Warn().
		Str("Network Name", k.NetworkConfig.Name).
		Msg("Setting GasTipCap = SuggestedGasPrice for Klaytn network")
	opts.GasTipCap = nil
	opts.GasPrice = nil

	contractAddress, transaction, contractInstance, err := deployer(opts, k.Client)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := k.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", k.DefaultWallet.Address()).
		Str("Total Gas Cost (KLAY)", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", k.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}
