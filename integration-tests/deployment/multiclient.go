package deployment

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pelletier/go-toml/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

const (
	RPC_RETRY_ATTEMPTS = 10
	RPC_RETRY_DELAY    = 1000 * time.Millisecond
)

// MultiClient should comply with the OnchainClient interface
var _ OnchainClient = &MultiClient{}

type MultiClient struct {
	*ethclient.Client
	backup []*ethclient.Client
	// we will use Seth only for gas estimations, confirming and tracing the transactions, but for sending transactions we will use pure ethclient
	// so that MultiClient conforms to the OnchainClient interface
	SethClient   *seth.Client
	EvmKMSClient *evmKMSClient
	chainId      uint64
}

type RPC struct {
	RPCName string `toml:"rpc_name"`
	HTTPURL string `toml:"http_url"`
	WSURL   string `toml:"ws_url"`
}

type Config struct {
	EnvConfig EnvConfig `toml:"env_config"`
}

type EnvConfig struct {
	TestWalletKey        string       `toml:"test_wallet_key"`
	KmsDeployerKeyId     string       `toml:"kms_deployer_key_id"`
	KmsDeployerKeyRegion string       `toml:"kms_deployer_key_region"`
	AwsProfileName       string       `toml:"aws_profile_name"`
	EvmNetworks          []EvmNetwork `toml:"evm_networks"`
	// Seth-related
	GethWrappersDirs []string `toml:"geth_wrappers_dirs"`
	SethConfigFile   string   `toml:"seth_config_file"`
}

type EvmNetwork struct {
	ChainID         uint64 `toml:"chain_id"`
	EtherscanAPIKey string `toml:"etherscan_api_key"`
	EtherscanUrl    string `toml:"etherscan_url"`
	RPCs            []RPC  `toml:"rpcs"`
}

func initRpcClients(rpcs []RPC) (*ethclient.Client, []*ethclient.Client) {
	if len(rpcs) == 0 {
		panic("No RPCs provided")
	}
	clients := make([]*ethclient.Client, 0, len(rpcs))

	for _, rpc := range rpcs {
		client, err := ethclient.Dial(rpc.HTTPURL)
		if err != nil {
			panic(err)
		}
		clients = append(clients, client)
	}
	return clients[0], clients[1:]
}

func NewMultiClientWithSeth(rpcs []RPC, chainId uint64, config Config) *MultiClient {
	mainClient, backupClients := initRpcClients(rpcs)
	mc := &MultiClient{
		Client:  mainClient,
		backup:  backupClients,
		chainId: chainId,
	}

	sethClient, err := buildSethClient(rpcs[0].HTTPURL, chainId, config)
	if err != nil {
		panic(err)
	}

	mc.SethClient = sethClient
	mc.EvmKMSClient = initialiseKMSClient(config)

	return mc
}

func buildSethClient(rpc string, chainId uint64, config Config) (*seth.Client, error) {
	var sethClient *seth.Client
	var err error

	// if config path is provided use the TOML file to configure Seth to provide maximum flexibility
	if config.EnvConfig.SethConfigFile != "" {
		sethConfig, readErr := readSethConfigFromFile(config.EnvConfig.SethConfigFile)
		if readErr != nil {
			return nil, readErr
		}

		sethClient, err = seth.NewClientBuilderWithConfig(sethConfig).
			UseNetworkWithChainId(chainId).
			WithRpcUrl(rpc).
			WithPrivateKeys([]string{config.EnvConfig.TestWalletKey}).
			Build()
	} else {
		// if full flexibility is not needed we create a client with reasonable defaults
		// if you need to further tweak them, please refer to https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/seth/README.md
		sethClient, err = seth.NewClientBuilder().
			WithRpcUrl(rpc).
			WithPrivateKeys([]string{config.EnvConfig.TestWalletKey}).
			WithProtections(true, true, seth.MustMakeDuration(1*time.Minute)).
			WithGethWrappersFolders(config.EnvConfig.GethWrappersDirs).
			// Fast priority will add a 20% buffer on top of what the node suggests
			// we will use last 20 block to estimate block congestion and further bump gas price suggested by the node
			WithGasPriceEstimations(true, 20, seth.Priority_Fast).
			Build()
	}

	return sethClient, err
}

func readSethConfigFromFile(configPath string) (*seth.Config, error) {
	d, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var sethConfig seth.Config
	err = toml.Unmarshal(d, &sethConfig)
	if err != nil {
		return nil, err
	}

	return &sethConfig, nil
}

func initialiseKMSClient(config Config) *evmKMSClient {
	if config.EnvConfig.KmsDeployerKeyId != "" && config.EnvConfig.KmsDeployerKeyRegion != "" {
		var awsSessionFn AwsSessionFn
		if config.EnvConfig.AwsProfileName != "" {
			awsSessionFn = awsSessionFromProfileFn
		} else {
			awsSessionFn = awsSessionFromEnvVarsFn
		}
		return NewEVMKMSClient(kms.New(awsSessionFn(config)), config.EnvConfig.KmsDeployerKeyId)
	}
	return nil
}

func (mc *MultiClient) GetKMSKey() *bind.TransactOpts {
	kmsTxOpts, err := mc.EvmKMSClient.GetKMSTransactOpts(context.Background(), big.NewInt(int64(mc.chainId)))
	if err != nil {
		panic(err)
	}
	// nonce needs to be `nil` so that RPC node sets it, otherwise Seth would set it to whatever it was, when we requested the key
	return mc.SethClient.NewTXOpts(seth.WithNonce(nil), seth.WithFrom(kmsTxOpts.From), seth.WithSignerFn(kmsTxOpts.Signer))
}

func (mc *MultiClient) GetTestWalletKey() *bind.TransactOpts {
	// nonce needs to be `nil` so that RPC node sets it, otherwise Seth would set it to whatever it was, when we requested the key
	return mc.SethClient.NewTXOpts(seth.WithNonce(nil))
}

func (mc *MultiClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	err := mc.retryWithBackups(func(client *ethclient.Client) error {
		var err error
		receipt, err = client.TransactionReceipt(ctx, txHash)
		return err
	})
	return receipt, err
}

func (mc *MultiClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return mc.retryWithBackups(func(client *ethclient.Client) error {
		return client.SendTransaction(ctx, tx)
	})
}

func (mc *MultiClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	var code []byte
	err := mc.retryWithBackups(func(client *ethclient.Client) error {
		var err error
		code, err = client.CodeAt(ctx, account, blockNumber)
		return err
	})
	return code, err
}

func (mc *MultiClient) NonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var count uint64
	err := mc.retryWithBackups(func(client *ethclient.Client) error {
		var err error
		count, err = client.NonceAt(ctx, account, nil)
		return err
	})
	return count, err
}

func (mc *MultiClient) retryWithBackups(op func(*ethclient.Client) error) error {
	var err error
	for _, client := range append([]*ethclient.Client{mc.Client}, mc.backup...) {
		err2 := retry.Do(func() error {
			err = op(client)
			if err != nil {
				fmt.Printf("  [MultiClient RPC] Retrying with new client, error: %v\n", err)
				return err
			}
			return nil
		}, retry.Attempts(RPC_RETRY_ATTEMPTS), retry.Delay(RPC_RETRY_DELAY))
		if err2 == nil {
			return nil
		}
	}
	return err
}
