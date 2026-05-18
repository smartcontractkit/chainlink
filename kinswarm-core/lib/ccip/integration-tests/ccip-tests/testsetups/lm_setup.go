package testsetups

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"

	ocrconfighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"golang.org/x/crypto/curve25519"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	integrationactions "github.com/smartcontractkit/chainlink/integration-tests/actions"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/testhelpers/integration"
)

type LMTestSetupOutputs struct {
	CCIPTestSetUpOutputs
	LMModules map[int64]*actions.LMCommon
}

// TODO - Copied over from ccip tests as such. Refactor and remove unused code
func (o *LMTestSetupOutputs) CreateLMEnvironment(
	lggr *zerolog.Logger,
	envName string,
	reportPath string,
) map[int64]blockchain.EVMClient {
	t := o.Cfg.Test
	testConfig := o.Cfg
	var (
		ccipEnv  *actions.CCIPTestEnv
		k8Env    *environment.Environment
		err      error
		chains   []blockchain.EVMClient
		local    *test_env.CLClusterTestEnv
		deployCL func() error
	)

	envConfig := createEnvironmentConfig(t, envName, testConfig, reportPath)

	configureCLNode := !testConfig.useExistingDeployment() || pointer.GetString(testConfig.EnvInput.EnvToConnect) != ""
	namespace := ""
	if testConfig.TestGroupInput.LoadProfile != nil {
		namespace = testConfig.TestGroupInput.LoadProfile.TestRunName
	}
	require.False(t, testConfig.localCluster() && testConfig.ExistingCLCluster(),
		"local cluster and existing cluster cannot be true at the same time")
	// if it's a new deployment, deploy the env
	// Or if EnvToConnect is given connect to that k8 environment
	if configureCLNode {
		if !testConfig.ExistingCLCluster() {
			// if it's a local cluster, deploy the local cluster in docker
			if testConfig.localCluster() {
				local, deployCL = DeployLocalCluster(t, testConfig)
				ccipEnv = &actions.CCIPTestEnv{
					LocalCluster: local,
				}
				namespace = "local-docker-deployment"
			} else {
				// Otherwise, deploy the k8s env
				lggr.Info().Msg("Deploying test environment")
				// deploy the env if configureCLNode is true
				k8Env = DeployEnvironments(t, envConfig, testConfig)
				ccipEnv = &actions.CCIPTestEnv{K8Env: k8Env}
				namespace = ccipEnv.K8Env.Cfg.Namespace
			}
		} else {
			// if there is already a cluster, use the existing cluster to connect to the nodes
			ccipEnv = &actions.CCIPTestEnv{}
			mockserverURL := pointer.GetString(testConfig.EnvInput.Mockserver)
			require.NotEmpty(t, mockserverURL, "mockserver URL cannot be nil")
			ccipEnv.MockServer = ctfClient.NewMockserverClient(&ctfClient.MockserverConfig{
				LocalURL:   mockserverURL,
				ClusterURL: mockserverURL,
			})
		}
		ccipEnv.CLNodeWithKeyReady, _ = errgroup.WithContext(o.SetUpContext)
		o.Env = ccipEnv
		if ccipEnv.K8Env != nil && ccipEnv.K8Env.WillUseRemoteRunner() {
			return nil
		}
	} else {
		// if configureCLNode is false it means we don't need to deploy any additional pods,
		// use a placeholder env to create just the remote runner in it.
		if value, set := os.LookupEnv(config.EnvVarJobImage); set && value != "" {
			k8Env = environment.New(envConfig)
			err = k8Env.Run()
			require.NoErrorf(t, err, "error creating environment remote runner")
			o.Env = &actions.CCIPTestEnv{K8Env: k8Env}
			if k8Env.WillUseRemoteRunner() {
				return nil
			}
		}
	}
	chainByChainID := make(map[int64]blockchain.EVMClient)
	if pointer.GetBool(testConfig.TestGroupInput.LocalCluster) {
		require.NotNil(t, ccipEnv.LocalCluster, "Local cluster shouldn't be nil")
		for _, n := range ccipEnv.LocalCluster.EVMNetworks {
			if evmClient, err := blockchain.NewEVMClientFromNetwork(*n, *lggr); err == nil {
				chainByChainID[evmClient.GetChainID().Int64()] = evmClient
				chains = append(chains, evmClient)
			} else {
				lggr.Error().Err(err).Msgf("EVMClient for chainID %d not found", n.ChainID)
			}
		}
	} else {
		for _, n := range testConfig.SelectedNetworks {
			if _, ok := chainByChainID[n.ChainID]; ok {
				continue
			}
			var ec blockchain.EVMClient
			if k8Env == nil {
				ec, err = blockchain.ConnectEVMClient(n, *lggr)
			} else {
				log.Info().Interface("urls", k8Env.URLs).Msg("URLs")
				ec, err = blockchain.NewEVMClient(n, k8Env, *lggr)
			}
			require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
			chains = append(chains, ec)
			chainByChainID[n.ChainID] = ec
		}
	}
	if configureCLNode {
		ccipEnv.CLNodeWithKeyReady.Go(func() error {
			var totalNodes int
			if !o.Cfg.ExistingCLCluster() {
				if ccipEnv.LocalCluster != nil {
					err = deployCL()
					if err != nil {
						return err
					}
				}
				err = ccipEnv.ConnectToDeployedNodes()
				if err != nil {
					return fmt.Errorf("error connecting to chainlink nodes: %w", err)
				}
				totalNodes = pointer.GetInt(testConfig.EnvInput.NewCLCluster.NoOfNodes)
			} else {
				totalNodes = pointer.GetInt(testConfig.EnvInput.ExistingCLCluster.NoOfNodes)
				err = ccipEnv.ConnectToExistingNodes(o.Cfg.EnvInput)
				if err != nil {
					return fmt.Errorf("error deploying and connecting to chainlink nodes: %w", err)
				}
			}
			err = ccipEnv.SetUpNodeKeysAndFund(lggr, big.NewFloat(testConfig.TestGroupInput.NodeFunding), chains)
			if err != nil {
				return fmt.Errorf("error setting up nodes and keys %w", err)
			}
			// first node is the bootstrapper
			ccipEnv.CommitNodeStartIndex = 1
			ccipEnv.ExecNodeStartIndex = 1
			ccipEnv.NumOfCommitNodes = testConfig.TestGroupInput.NoOfCommitNodes
			ccipEnv.NumOfExecNodes = ccipEnv.NumOfCommitNodes
			if !pointer.GetBool(testConfig.TestGroupInput.CommitAndExecuteOnSameDON) {
				if len(ccipEnv.CLNodesWithKeys) < 11 {
					return fmt.Errorf("not enough CL nodes for separate commit and execution nodes")
				}
				if testConfig.TestGroupInput.NoOfCommitNodes >= totalNodes {
					return fmt.Errorf("number of commit nodes can not be greater than total number of nodes in DON")
				}
				// first two nodes are reserved for bootstrap commit and bootstrap exec
				ccipEnv.CommitNodeStartIndex = 2
				ccipEnv.ExecNodeStartIndex = 2 + testConfig.TestGroupInput.NoOfCommitNodes
				ccipEnv.NumOfExecNodes = totalNodes - (2 + testConfig.TestGroupInput.NoOfCommitNodes)
				if ccipEnv.NumOfExecNodes < 4 {
					return fmt.Errorf("insufficient number of exec nodes")
				}
			}
			ccipEnv.NumOfAllowedFaultyExec = (ccipEnv.NumOfExecNodes - 1) / 3
			ccipEnv.NumOfAllowedFaultyCommit = (ccipEnv.NumOfCommitNodes - 1) / 3
			return nil
		})
	}

	t.Cleanup(func() {
		if configureCLNode {
			if ccipEnv.LocalCluster != nil {
				err := ccipEnv.LocalCluster.Terminate()
				require.NoError(t, err, "Local cluster termination shouldn't fail")
				//require.NoError(t, o.Reporter.SendReport(t, namespace, false), "Aggregating and sending report shouldn't fail")
				return
			}
			if pointer.GetBool(testConfig.TestGroupInput.KeepEnvAlive) || testConfig.ExistingCLCluster() {
				//require.NoError(t, o.Reporter.SendReport(t, namespace, true), "Aggregating and sending report shouldn't fail")
				return
			}
			lggr.Info().Msg("Tearing down the environment")
			err = integrationactions.TeardownSuite(t, nil, ccipEnv.K8Env, ccipEnv.CLNodes, o.Reporter, zapcore.DPanicLevel, o.Cfg.EnvInput)
			require.NoError(t, err, "Environment teardown shouldn't fail")
		} else {
			//just send the report
			require.NoError(t, o.Reporter.SendReport(t, namespace, true), "Aggregating and sending report shouldn't fail")
		}
	})
	return chainByChainID
}

func (o *LMTestSetupOutputs) DeployLMChainContracts(
	lggr *zerolog.Logger,
	networkCfg blockchain.EVMNetwork,
	lmCommon actions.LMCommon,
	l2ChainID int64,
) error {
	var k8Env *environment.Environment
	ccipEnv := o.Env
	chainClient := lmCommon.ChainClient
	if ccipEnv != nil {
		k8Env = ccipEnv.K8Env
	}
	if k8Env != nil && chainClient.NetworkSimulated() {
		networkCfg.URLs = k8Env.URLs[chainClient.GetNetworkConfig().Name]
	}

	chain, err := blockchain.ConcurrentEVMClient(networkCfg, k8Env, chainClient, *lggr)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create chain client for %s: %w", networkCfg.Name, err))
	}

	chain.ParallelTransactions(true)
	//defer chain.Close()

	cd, err := contracts.NewCCIPContractsDeployer(lggr, chain)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create contract deployer: %w", err))
	}

	// Deploy Wrapped Native contract only on private geth networks
	if lmCommon.ChainSelectror == chainselectors.GETH_TESTNET.Selector ||
		lmCommon.ChainSelectror == chainselectors.GETH_DEVNET_2.Selector {
		lggr.Info().Msg("Deploying Wrapped Native contract")
		wrapperNative, err := cd.DeployWrappedNative()
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to deploy Wrapped Native contract: %w", err))
		}
		lggr.Info().Str("Address", wrapperNative.String()).Msg("Deployed Wrapped Native contract")
		lmCommon.WrapperNative = wrapperNative
	}

	// Deploy Bridge Adapter contracts
	switch lmCommon.ChainSelectror {
	case chainselectors.GETH_TESTNET.Selector:
		lggr.Info().Msg("Deploying Mock L1 Bridge Adapter contract")
		bridgeAdapter, err := cd.DeployMockL1BridgeAdapter(*lmCommon.WrapperNative, true)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to deploy Mock L1 Bridge Adapter contract: %w", err))
		}
		lggr.Info().Str("Address", bridgeAdapter.EthAddress.String()).Msg("Deployed Mock L1 Bridge Adapter contract")
		lmCommon.BridgeAdapterAddr = bridgeAdapter.EthAddress
	case chainselectors.GETH_DEVNET_2.Selector:
		lggr.Info().Msg("Deploying Mock L2 Bridge Adapter contract")
		bridgeAdapter, err := cd.DeployMockL2BridgeAdapter()
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to deploy Mock L2 Bridge Adapter contract: %w", err))
		}
		lggr.Info().Str("Address", bridgeAdapter.EthAddress.String()).Msg("Deployed Mock L2 Bridge Adapter contract")
		lmCommon.BridgeAdapterAddr = bridgeAdapter.EthAddress
	case chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector:
		if l2ChainID == int64(chainselectors.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.EvmChainID) {
			wethAddress := common.HexToAddress("0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9")
			lmCommon.WrapperNative = &wethAddress
			lggr.Info().Msg("Deploying Arbitrum L1 Bridge Adapter contract")
			bridgeAdapter, err := cd.DeployArbitrumL1BridgeAdapter(
				common.HexToAddress("0xcE18836b233C83325Cc8848CA4487e94C6288264"),
				common.HexToAddress("0x65f07C7D521164a4d5DaC6eB8Fac8DA067A3B78F"),
			)
			if err != nil {
				return errors.WithStack(fmt.Errorf("failed to deploy Arbitrum L1 Bridge Adapter contract: %w", err))
			}
			lggr.Info().Str("Address", bridgeAdapter.EthAddress.String()).Msg("Deployed Arbitrum L1 Bridge Adapter contract")
			lmCommon.BridgeAdapterAddr = bridgeAdapter.EthAddress
		}
		if l2ChainID == int64(chainselectors.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.EvmChainID) {
			wethAddress := common.HexToAddress("0x7b79995e5f793a07bc00c21412e50ecae098e7f9")
			lmCommon.WrapperNative = &wethAddress
			lggr.Info().Msg("Deploying Optimism L1 Bridge Adapter contract")
			bridgeAdapter, err := cd.DeployOptimismL1BridgeAdapter(
				common.HexToAddress("0xFBb0621E0B23b5478B630BD55a5f21f67730B0F1"),
				*lmCommon.WrapperNative,
				common.HexToAddress("0x16Fc5058F25648194471939df75CF27A2fdC48BC"),
			)
			if err != nil {
				return errors.WithStack(fmt.Errorf("failed to deploy Optimism L1 Bridge Adapter contract: %w", err))
			}
			lggr.Info().Str("Address", bridgeAdapter.EthAddress.String()).Msg("Deployed Optimism L1 Bridge Adapter contract")
			lmCommon.BridgeAdapterAddr = bridgeAdapter.EthAddress
		}
	case chainselectors.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector:
		wethAddress := common.HexToAddress("0x980B62Da83eFf3D4576C647993b0c1D7faf17c73")
		lmCommon.WrapperNative = &wethAddress
		lggr.Info().Msg("Deploying Arbitrum L2 Bridge Adapter contract")
		bridgeAdapter, err := cd.DeployArbitrumL2BridgeAdapter(common.HexToAddress("0x9fDD1C4E4AA24EEc1d913FABea925594a20d43C7"))
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to deploy Arbitrum L2 Bridge Adapter contract: %w", err))
		}
		lggr.Info().Str("Address", bridgeAdapter.EthAddress.String()).Msg("Deployed Arbitrum L2 Bridge Adapter contract")
		lmCommon.BridgeAdapterAddr = bridgeAdapter.EthAddress
	case chainselectors.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector:
		wethAddress := common.HexToAddress("0x4200000000000000000000000000000000000006")
		lmCommon.WrapperNative = &wethAddress
		lggr.Info().Msg("Deploying Optimism L2 Bridge Adapter contract")
		bridgeAdapter, err := cd.DeployOptimismL2BridgeAdapter(*lmCommon.WrapperNative)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to deploy Optimism L2 Bridge Adapter contract: %w", err))
		}
		lggr.Info().Str("Address", bridgeAdapter.EthAddress.String()).Msg("Deployed Optimism L2 Bridge Adapter contract")
		lmCommon.BridgeAdapterAddr = bridgeAdapter.EthAddress
	}

	// Deploy Mock ARM contract
	lggr.Info().Msg("Deploying Mock ARM contract")
	mockRMNContract, err := cd.DeployMockRMNContract()
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to deploy Mock ARM contract: %w", err))
	}
	lggr.Info().Str("Address", mockRMNContract.String()).Msg("Deployed Mock ARM contract")
	lmCommon.MockArm = mockRMNContract

	// Deploy ARM Proxy contract
	lggr.Info().Msg("Deploying ARM Proxy contract")
	RMNProxyContract, err := cd.DeployArmProxy(*mockRMNContract)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to deploy ARM Proxy contract: %w", err))
	}
	lggr.Info().Str("Address", RMNProxyContract.EthAddress.String()).Msg("Deployed ARM Proxy contract")
	lmCommon.ArmProxy = RMNProxyContract

	// Deploy CCIP Router contract
	lggr.Info().Msg("Deploying CCIP Router contract")
	ccipRouterContract, err := cd.DeployRouter(common.Address{}, *lmCommon.ArmProxy.EthAddress)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to deploy CCIP Router contract: %w", err))
	}
	lggr.Info().Str("Address", ccipRouterContract.EthAddress.String()).Msg("Deployed CCIP Router contract")
	lmCommon.CcipRouter = ccipRouterContract

	// Deploy Lock Release Token contract
	lggr.Info().Msg("Deploying Lock Release Token contract")
	lockReleaseTokenPool, err := cd.DeployLockReleaseTokenPoolContract(lmCommon.WrapperNative.String(), *lmCommon.MockArm, lmCommon.CcipRouter.EthAddress)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to deploy Lock Release Token contract: %w", err))
	}
	lggr.Info().Str("Address", lockReleaseTokenPool.EthAddress.String()).Msg("Deployed Lock Release Token contract")
	lmCommon.TokenPool = lockReleaseTokenPool

	// Deploy Liquidity Manager contract
	lggr.Info().Msg("Deploying Liquidity Manager contract")
	liquidityManager, err := cd.DeployLiquidityManager(*lmCommon.WrapperNative, lmCommon.ChainSelectror, lmCommon.TokenPool.EthAddress, lmCommon.MinimumLiquidity)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to deploy Liquidity Manager contract: %w", err))
	}
	lggr.Info().Str("Address", liquidityManager.EthAddress.String()).Msg("Deployed Liquidity Manager contract")
	lmCommon.LM = liquidityManager

	// Set Liquidity Manager on Token Pool
	lggr.Info().Msg("Setting Liquidity Manager on Token Pool")
	err = lockReleaseTokenPool.SetRebalancer(*liquidityManager.EthAddress)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to set Liquidity Manager on Token Pool: %w", err))
	}
	lggr.Info().Msg("Set Liquidity Manager on Token Pool")

	err = chain.WaitForEvents()
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to wait for events: %w", err))
	}

	// Verify on chain rebalancer from token pool matches deployed Liquidity Manager
	onchainRebalancer, err := lockReleaseTokenPool.GetRebalancer()
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to get rebalancer from Token Pool: %w", err))
	}
	if onchainRebalancer != *liquidityManager.EthAddress {
		return errors.WithStack(fmt.Errorf("onchainRebalancer doesn not match the deployed Liquidity Manager"))
	}

	lggr.Debug().Interface("lmCommon", lmCommon).Msg("lmCommon")
	o.LMModules[chainClient.GetChainID().Int64()] = &lmCommon

	return nil
}

func stripKeyPrefix(key string) string {
	chunks := strings.Split(key, "_")
	if len(chunks) == 3 {
		return chunks[2]
	}
	return key
}

func (o *LMTestSetupOutputs) SetOCR3Config(chainId int64) error {
	clNodesWithKeys := o.Env.CLNodesWithKeys[strconv.FormatInt(chainId, 10)]
	donNodes := clNodesWithKeys[1:]
	oracleIdentities := make([]ocrconfighelper2.OracleIdentityExtra, 0)
	var onChainKeys []ocrtypes2.OnchainPublicKey
	var transmitters []common.Address
	var schedule []int

	for i, nodeWithKeys := range donNodes {
		ocr2Key := nodeWithKeys.KeysBundle.OCR2Key.Data
		offChainPubKeyTemp, err := hex.DecodeString(stripKeyPrefix(ocr2Key.Attributes.OffChainPublicKey))
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to decode offchain public key: %w", err))
		}
		formattedOnChainPubKey := stripKeyPrefix(ocr2Key.Attributes.OnChainPublicKey)
		cfgPubKeyTemp, err := hex.DecodeString(stripKeyPrefix(ocr2Key.Attributes.ConfigPublicKey))
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to decode config public key: %w", err))
		}
		cfgPubKeyBytes := [ed25519.PublicKeySize]byte{}
		copy(cfgPubKeyBytes[:], cfgPubKeyTemp)
		offChainPubKey := [curve25519.PointSize]byte{}
		copy(offChainPubKey[:], offChainPubKeyTemp)
		ethAddress := nodeWithKeys.KeysBundle.EthAddress
		p2pKeys := nodeWithKeys.KeysBundle.P2PKeys
		peerID := p2pKeys.Data[0].Attributes.PeerID
		oracleIdentities = append(oracleIdentities, ocrconfighelper2.OracleIdentityExtra{
			OracleIdentity: ocrconfighelper2.OracleIdentity{
				OffchainPublicKey: offChainPubKey,
				OnchainPublicKey:  common.HexToAddress(formattedOnChainPubKey).Bytes(),
				PeerID:            peerID,
				TransmitAccount:   ocrtypes2.Account(ethAddress),
			},
			ConfigEncryptionPublicKey: cfgPubKeyBytes,
		})
		onChainKeys = append(onChainKeys, oracleIdentities[i].OnchainPublicKey)
		transmitters = append(transmitters, common.HexToAddress(ethAddress))
		schedule = append(schedule, 1)

	}
	signers, err := evm.OnchainPublicKeyToAddress(onChainKeys)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to convert onchain public keys to addresses: %w", err))
	}

	offchainConfig, onchainConfig := []byte{}, []byte{}
	f := uint8(1)
	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		2*time.Minute,
		2*time.Minute,
		20*time.Second,
		2*time.Second,
		20*time.Second,
		10*time.Second,
		40*time.Second,
		3,
		schedule,
		oracleIdentities,
		offchainConfig,
		50*time.Millisecond,
		1*time.Minute,
		1*time.Minute,
		1*time.Second,
		int(f),
		onchainConfig,
	)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to set OCR3 config args for tests: %w", err))
	}
	err = o.LMModules[chainId].LM.SetOCR3Config(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to set OCR3 config: %w", err))
	}
	return nil
}

func (o *LMTestSetupOutputs) FundPool(chainId int64, lggr *zerolog.Logger, fundingAmount *big.Int) error {
	token, err := erc20.NewERC20(*o.LMModules[chainId].WrapperNative, o.LMModules[chainId].ChainClient.Backend())
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create ERC20 contract instance: %w", err))
	}
	balance, err := token.BalanceOf(nil, common.HexToAddress(o.LMModules[chainId].ChainClient.GetDefaultWallet().Address()))
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to get token pool balance: %w", err))
	}
	lggr.Debug().Str("balance", balance.String()).Msg("weth balance of transactor")
	symbol, err := token.Symbol(nil)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to get token symbol: %w", err))
	}
	if symbol == "WETH" {
		weth, err := weth9.NewWETH9(*o.LMModules[chainId].WrapperNative, o.LMModules[chainId].ChainClient.Backend())
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to create WETH contract instance: %w", err))
		}
		nativeBalance, err := o.LMModules[chainId].ChainClient.BalanceAt(context.Background(), common.HexToAddress(o.LMModules[chainId].ChainClient.GetDefaultWallet().Address()))
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to get native balance: %w", err))
		}
		lggr.Debug().Str("nativeBalance", nativeBalance.String()).Msg("nativeBalance")
		if nativeBalance.Cmp(fundingAmount) < 0 {
			return errors.WithStack(fmt.Errorf("not enough native balance"))
		}
		lggr.Info().Msg("Depositing tokenpool funding to WETH contract")
		txOpts, err := o.LMModules[chainId].ChainClient.TransactionOpts(o.LMModules[chainId].ChainClient.GetDefaultWallet())
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to get transaction options: %w", err))
		}
		txOpts.Value = fundingAmount
		tx, err := weth.Deposit(txOpts)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to deposit to WETH contract: %w", err))
		}
		receipt, err := bind.WaitMined(context.Background(), o.LMModules[chainId].ChainClient.DeployBackend(), tx)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to wait for transaction receipt: %w", err))
		}

		lggr.Info().Str("tx hash", receipt.TxHash.String()).Msg("Deposited tokenpool funding to WETH contract")
	}
	lggr.Info().Msg("Funding token pool")
	txOpts, err := o.LMModules[chainId].ChainClient.TransactionOpts(o.LMModules[chainId].ChainClient.GetDefaultWallet())
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to get transaction options: %w", err))

	}
	tx, err := token.Transfer(txOpts, o.LMModules[chainId].TokenPool.EthAddress, fundingAmount)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to transfer to token pool: %w", err))
	}
	receipt, err := bind.WaitMined(context.Background(), o.LMModules[chainId].ChainClient.DeployBackend(), tx)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to wait for transaction receipt: %w", err))
	}
	lggr.Info().Str("tx hash", receipt.TxHash.String()).Msg("Funded token pool")

	balance, err = token.BalanceOf(nil, o.LMModules[chainId].TokenPool.EthAddress)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to get token pool balance: %w", err))
	}
	lggr.Debug().Str("balance", balance.String()).Msg("weth balance of token pool")

	return nil
}

func (o *LMTestSetupOutputs) FundLM(chainId int64, lggr *zerolog.Logger, fundingAmount *big.Int) error {
	transactor, err := o.LMModules[chainId].ChainClient.TransactionOpts(o.LMModules[chainId].ChainClient.GetDefaultWallet())
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to get transaction options: %w", err))
	}
	cl := o.LMModules[chainId].ChainClient.Backend()

	nonce, err := cl.PendingNonceAt(context.Background(), transactor.From)
	if err != nil {
		return err
	}

	gasPrice, err := cl.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	gasEstimate, err := cl.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  transactor.From,
		To:    o.LMModules[chainId].LM.EthAddress,
		Value: fundingAmount,
	})
	if err != nil {
		return err
	}

	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasEstimate,
			To:       o.LMModules[chainId].LM.EthAddress,
			Value:    fundingAmount,
		},
	)
	signedTx, err := transactor.Signer(transactor.From, tx)
	if err != nil {
		return err
	}
	lggr.Info().Msg("Funding Liquidity Manager")
	err = cl.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(context.Background(), o.LMModules[chainId].ChainClient.DeployBackend(), signedTx)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to wait for transaction receipt: %w", err))
	}
	lggr.Info().Str("tx hash", receipt.TxHash.String()).Msg("Funded Liquidity Manager")
	return nil
}

func (o *LMTestSetupOutputs) AddJobs(chainId int64, lggr *zerolog.Logger) error {
	// Add bootstrap job
	clNodesWithKeys := o.Env.CLNodesWithKeys[strconv.FormatInt(chainId, 10)]
	bootstrapNode := clNodesWithKeys[0]
	bootstrapSpec, err := integrationtesthelpers.NewBootsrapJobSpec(&integrationtesthelpers.LMJobSpecParams{
		ChainID:            uint64(chainId),
		ContractID:         o.LMModules[chainId].LM.EthAddress.String(),
		CfgTrackerInterval: 15 * time.Second,
	})
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create bootstrap job spec: %w", err))
	}
	lggr.Info().Msg("Adding bootstrap job")
	j, err := bootstrapNode.Node.MustCreateJob(bootstrapSpec)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create bootstrap job: %w", err))
	}
	lggr.Info().Str("jobId", j.Data.ID).Msg("Bootstrap job added")

	P2Pv2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapNode.KeysBundle.P2PKeys.Data[0].Attributes.PeerID, bootstrapNode.Node.InternalIP(), 6690)

	// Add LM jobs
	donNodes := clNodesWithKeys[1:]

	for _, node := range donNodes {
		lmJobSpec, err := integrationtesthelpers.NewJobSpec(&integrationtesthelpers.LMJobSpecParams{
			ChainID:                 uint64(chainId),
			ContractID:              o.LMModules[chainId].LM.EthAddress.String(),
			OCRKeyBundleID:          node.KeysBundle.OCR2Key.Data.ID,
			TransmitterID:           node.KeysBundle.EthAddress,
			P2PV2Bootstrappers:      pq.StringArray{P2Pv2Bootstrapper},
			CfgTrackerInterval:      15 * time.Second,
			LiquidityManagerAddress: *o.LMModules[chainId].LM.EthAddress,
			NetworkSelector:         o.LMModules[chainId].ChainSelectror,
			Type:                    "ping-pong",
		})
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to create LM job spec: %w", err))
		}
		lggr.Debug().Interface("lmJobSpec", lmJobSpec).Msg("lmJobSpec")
		lggr.Info().Str("Node URL", node.Node.URL()).Msg("Adding LM job")
		j, err := node.Node.MustCreateJob(lmJobSpec)
		if err != nil {
			return errors.WithStack(fmt.Errorf("failed to create LM job: %w", err))
		}
		lggr.Info().Str("jobId", j.Data.ID).Msg("LM job added")

	}
	return nil
}

func LMDefaultTestSetup(
	t *testing.T,
	lggr *zerolog.Logger,
	envName string,
	testConfig *CCIPTestConfig,
) *LMTestSetupOutputs {
	var (
		err error
	)
	reportPath := "tmp_laneconfig"
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()
	lmModules := make(map[int64]*actions.LMCommon)
	setUpArgs := &LMTestSetupOutputs{
		CCIPTestSetUpOutputs{
			SetUpContext: parent,
			Cfg:          testConfig,
		},
		lmModules,
	}

	chainByChainID := setUpArgs.CreateLMEnvironment(lggr, envName, reportPath)

	chainAddGrp, _ := errgroup.WithContext(setUpArgs.SetUpContext)
	lggr.Info().Msg("Deploying common contracts")
	chainSelectors := make(map[int64]uint64)

	testConfig.SelectedNetworks, _, err = testConfig.EnvInput.EVMNetworks()
	require.NoError(t, err)

	testConfig.AllNetworks = make(map[string]blockchain.EVMNetwork)
	for _, net := range testConfig.SelectedNetworks {
		testConfig.AllNetworks[net.Name] = net
		if _, exists := chainSelectors[net.ChainID]; !exists {
			chainSelectors[net.ChainID], err = chainselectors.SelectorFromChainId(uint64(net.ChainID))
			require.NoError(t, err)
		}
	}

	l1ChainId := testConfig.SelectedNetworks[0].ChainID
	l2ChainId := testConfig.SelectedNetworks[1].ChainID

	for _, net := range testConfig.AllNetworks {
		chain := chainByChainID[net.ChainID]
		net := net
		net.HTTPURLs = chain.GetNetworkConfig().HTTPURLs
		net.URLs = chain.GetNetworkConfig().URLs
		var selectors []uint64
		for chainId, selector := range chainSelectors {
			if chainId == net.ChainID {
				selectors = append(selectors, selector)
			}
		}
		lmCommon, err := actions.DefaultLMModule(
			chain,
			big.NewInt(0),
			selectors[0],
		)
		require.NoError(t, err)
		chainAddGrp.Go(func() error {
			return setUpArgs.DeployLMChainContracts(lggr, net, *lmCommon, l2ChainId)
		})
	}
	require.NoError(t, chainAddGrp.Wait(), "Deploying common contracts shouldn't fail")

	lggr.Debug().Interface("lmModules", lmModules).Msg("lmModules")

	//Set Cross Chain Rebalancer on L1 Rebalancer
	err = lmModules[l1ChainId].LM.SetCrossChainRebalancer(
		liquiditymanager.ILiquidityManagerCrossChainRebalancerArgs{
			RemoteRebalancer:    *lmModules[l2ChainId].LM.EthAddress,
			LocalBridge:         *lmModules[l1ChainId].BridgeAdapterAddr,
			RemoteToken:         *lmModules[l2ChainId].WrapperNative,
			RemoteChainSelector: lmModules[l2ChainId].ChainSelectror,
			Enabled:             true,
		})
	require.NoError(t, err, "Setting Cross Chain Rebalancer on L1 Rebalancer shouldn't fail")

	//Set Cross Chain Rebalancer on L2 Rebalancer
	err = lmModules[l2ChainId].LM.SetCrossChainRebalancer(
		liquiditymanager.ILiquidityManagerCrossChainRebalancerArgs{
			RemoteRebalancer:    *lmModules[l1ChainId].LM.EthAddress,
			LocalBridge:         *lmModules[l2ChainId].BridgeAdapterAddr,
			RemoteToken:         *lmModules[l1ChainId].WrapperNative,
			RemoteChainSelector: lmModules[l1ChainId].ChainSelectror,
			Enabled:             true,
		})
	require.NoError(t, err, "Setting Cross Chain Rebalancer on L1 Rebalancer shouldn't fail")

	// Wait for setting cross chain balancers on both chains to confirm
	err = lmModules[l1ChainId].ChainClient.WaitForEvents()
	require.NoError(t, err, "Waiting for events to confirm on L1 chain shouldn't fail")

	err = lmModules[l2ChainId].ChainClient.WaitForEvents()
	require.NoError(t, err, "Waiting for events to confirm on L2 chain shouldn't fail")

	// Verify that onchain rebalancer matches the deployed Liquidity Manager
	onchainRebalancerL1, err := lmModules[l1ChainId].TokenPool.GetRebalancer()
	require.NoError(t, err, "Getting rebalancer from Token Pool shouldn't fail")

	onchainRebalancerL2, err := lmModules[l2ChainId].TokenPool.GetRebalancer()
	require.NoError(t, err, "Getting rebalancer from Token Pool shouldn't fail")

	if onchainRebalancerL1.String() != lmModules[l1ChainId].LM.EthAddress.String() ||
		onchainRebalancerL2.String() != lmModules[l2ChainId].LM.EthAddress.String() {
		lggr.Debug().
			Str("onchainRebalancerL1", onchainRebalancerL1.String()).
			Str("onchainRebalancerL2", onchainRebalancerL2.String()).
			Str("L2 LM", lmModules[l2ChainId].LM.EthAddress.String()).
			Str("L1 LM", lmModules[l1ChainId].LM.EthAddress.String()).
			Msg("Onchain rebalancer mismatch")
		t.Fatalf("Onchain rebalancer mismatch")
	}

	// Fund L1 Token Pool
	err = setUpArgs.FundPool(l1ChainId, lggr, big.NewInt(1000000000))
	require.NoError(t, err, "Funding L1 Token Pool shouldn't fail")

	//Fund L1 LM
	err = setUpArgs.FundLM(l1ChainId, lggr, big.NewInt(1000000000))
	require.NoError(t, err, "Funding L1 LM shouldn't fail")

	err = lmModules[l1ChainId].ChainClient.WaitForEvents()
	require.NoError(t, err, "Waiting for events to confirm on L1 chain shouldn't fail")

	// Fund L2 Token Pool
	err = setUpArgs.FundPool(l2ChainId, lggr, big.NewInt(1000000000))
	require.NoError(t, err, "Funding L2 Token Pool shouldn't fail")

	//Fund L2 LM
	err = setUpArgs.FundLM(l2ChainId, lggr, big.NewInt(1000000000))
	require.NoError(t, err, "Funding L2 LM shouldn't fail")

	err = lmModules[l2ChainId].ChainClient.WaitForEvents()
	require.NoError(t, err, "Waiting for events to confirm on L2 chain shouldn't fail")

	liquidity, err := setUpArgs.LMModules[l1ChainId].LM.GetLiquidity()
	require.NoError(t, err, "Getting liquidity from L1 LM shouldn't fail")
	lggr.Debug().Interface("liquidity", liquidity).Msg("Liquidity")
	require.Equal(t, big.NewInt(1000000000), liquidity, "Liquidity should match")

	liquidity, err = setUpArgs.LMModules[l2ChainId].LM.GetLiquidity()
	require.NoError(t, err, "Getting liquidity from L1 LM shouldn't fail")
	lggr.Debug().Interface("liquidity", liquidity).Msg("Liquidity")
	require.Equal(t, big.NewInt(1000000000), liquidity, "Liquidity should match")

	err = setUpArgs.Env.CLNodeWithKeyReady.Wait()
	require.NoError(t, err, "Waiting for CL nodes to be ready shouldn't fail")

	err = setUpArgs.AddJobs(l1ChainId, lggr)
	require.NoError(t, err, "Adding jobs on L1 chain shouldn't fail")

	// Set Config on L2 Chain
	err = setUpArgs.SetOCR3Config(l2ChainId)
	require.NoError(t, err, "Setting OCR3 config on L2 chain shouldn't fail")

	// TODO: Remove this sleep when it is no longer needed
	time.Sleep(30 * time.Second)

	// Set Config on L1 Chain
	err = setUpArgs.SetOCR3Config(l1ChainId)
	require.NoError(t, err, "Setting OCR3 config on L1 chain shouldn't fail")

	defer lmModules[l1ChainId].ChainClient.Close()
	defer lmModules[l2ChainId].ChainClient.Close()

	return setUpArgs
}
