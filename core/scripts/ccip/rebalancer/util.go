package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rebalancer/arb"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rebalancer/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/arm_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	router2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/arbitrum_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/arbitrum_l2_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
)

type universe struct {
	L1 struct {
		Arm                  common.Address
		ArmProxy             common.Address
		TokenPool            common.Address
		Rebalancer           common.Address
		BridgeAdapterAddress common.Address
	}
	L2 struct {
		Arm                  common.Address
		ArmProxy             common.Address
		TokenPool            common.Address
		Rebalancer           common.Address
		BridgeAdapterAddress common.Address
	}
}

func deployUniverse(
	env multienv.Env,
	l1ChainID, l2ChainID uint64,
	l1TokenAddress, l2TokenAddress common.Address,
) universe {
	validateEnv(env, l1ChainID, l2ChainID, false)

	l1Client, l2Client := env.Clients[l1ChainID], env.Clients[l2ChainID]
	l1Transactor, l2Transactor := env.Transactors[l1ChainID], env.Transactors[l2ChainID]
	l1ChainSelector, l2ChainSelector := mustGetChainByEvmID(l1ChainID).Selector, mustGetChainByEvmID(l2ChainID).Selector

	// L1 deploys
	// deploy arm and arm proxy.
	// required by the token pool.
	l1Arm, l1ArmProxy := deployArm(l1Transactor, l1Client, l1ChainID)

	l1RouterAddress, _, _, err := router2.DeployRouter(l1Transactor, l1Client, common.Address{}, l1ArmProxy.Address())
	helpers.PanicErr(err)

	// deploy token pool targeting l1TokenAddress
	l1TokenPool, l1Rebalancer := deployTokenPoolAndRebalancer(l1Transactor, l1Client, l1TokenAddress, l1ArmProxy.Address(), l1ChainSelector, l1RouterAddress)

	// deploy the L1 bridge adapter to point to the token address
	_, tx, _, err := arbitrum_l1_bridge_adapter.DeployArbitrumL1BridgeAdapter(
		l1Transactor,
		l1Client,
		arb.ArbitrumContracts[l1ChainID]["L1GatewayRouter"],
		arb.ArbitrumContracts[l1ChainID]["L1Outbox"],
	)
	helpers.PanicErr(err)
	l1BridgeAdapterAddress := helpers.ConfirmContractDeployed(context.Background(), l1Client, tx, int64(l1ChainID))

	// L2 deploys
	// deploy arm and arm proxy.
	// required by the token pool.
	l2Arm, l2ArmProxy := deployArm(l2Transactor, l2Client, l2ChainID)

	l2RouterAddress, _, _, err := router2.DeployRouter(l2Transactor, l2Client, common.Address{}, l2ArmProxy.Address())
	helpers.PanicErr(err)

	// deploy token pool targeting l2TokenAddress
	l2TokenPool, l2Rebalancer := deployTokenPoolAndRebalancer(l2Transactor, l2Client, l2TokenAddress, l2ArmProxy.Address(), l2ChainSelector, l2RouterAddress)

	// deploy the L2 bridge adapter to point to the token address
	_, tx, _, err = arbitrum_l2_bridge_adapter.DeployArbitrumL2BridgeAdapter(
		l2Transactor,
		l2Client,
		arb.ArbitrumContracts[l2ChainID]["L2GatewayRouter"],
	)
	helpers.PanicErr(err)
	l2BridgeAdapterAddress := helpers.ConfirmContractDeployed(context.Background(), l2Client, tx, int64(l2ChainID))

	// link the l1 and l2 rebalancers together via the SetCrossChainRebalancer function
	tx, err = l1Rebalancer.SetCrossChainRebalancer(l1Transactor, rebalancer.IRebalancerCrossChainRebalancerArgs{
		RemoteRebalancer:    l2Rebalancer.Address(),
		LocalBridge:         l1BridgeAdapterAddress,
		RemoteToken:         l2TokenAddress,
		RemoteChainSelector: l2ChainSelector,
		Enabled:             true,
	})
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), l1Client, tx, int64(l1ChainID), "setting cross chain rebalancer on L1 rebalancer")
	// assertion
	onchainRebalancer, err := l1Rebalancer.GetCrossChainRebalancer(nil, l2ChainSelector)
	helpers.PanicErr(err)
	if onchainRebalancer.RemoteRebalancer != l2Rebalancer.Address() ||
		onchainRebalancer.LocalBridge != l1BridgeAdapterAddress {
		panic(fmt.Sprintf("onchain rebalancer address does not match, expected %s got %s, or local bridge does not match, expected %s got %s",
			l2Rebalancer.Address().Hex(),
			onchainRebalancer.RemoteRebalancer.Hex(),
			l1BridgeAdapterAddress.Hex(),
			onchainRebalancer.LocalBridge.Hex()))
	}

	tx, err = l2Rebalancer.SetCrossChainRebalancer(l2Transactor, rebalancer.IRebalancerCrossChainRebalancerArgs{
		RemoteRebalancer:    l1Rebalancer.Address(),
		LocalBridge:         l2BridgeAdapterAddress,
		RemoteToken:         l1TokenAddress,
		RemoteChainSelector: l1ChainSelector,
		Enabled:             true,
	})
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), l2Client, tx, int64(l2ChainID), "setting cross chain rebalancer on L2 rebalancer")
	// assertion
	onchainRebalancer, err = l2Rebalancer.GetCrossChainRebalancer(nil, l1ChainSelector)
	helpers.PanicErr(err)
	if onchainRebalancer.RemoteRebalancer != l1Rebalancer.Address() ||
		onchainRebalancer.LocalBridge != l2BridgeAdapterAddress {
		panic(fmt.Sprintf("onchain rebalancer address does not match, expected %s got %s, or local bridge does not match, expected %s got %s",
			l1Rebalancer.Address().Hex(),
			onchainRebalancer.RemoteRebalancer.Hex(),
			l2BridgeAdapterAddress.Hex(),
			onchainRebalancer.LocalBridge.Hex()))
	}

	fmt.Println("Deployments complete\n",
		"L1 Arm:", l1Arm.Address().Hex(), "\n",
		"L1 Arm Proxy:", l1ArmProxy.Address().Hex(), "\n",
		"L1 Token Pool:", l1TokenPool.Address().Hex(), "\n",
		"L1 Rebalancer:", l1Rebalancer.Address().Hex(), "\n",
		"L1 Bridge Adapter:", l1BridgeAdapterAddress.Hex(), "\n",
		"L2 Arm:", l2Arm.Address().Hex(), "\n",
		"L2 Arm Proxy:", l2ArmProxy.Address().Hex(), "\n",
		"L2 Token Pool:", l2TokenPool.Address().Hex(), "\n",
		"L2 Rebalancer:", l2Rebalancer.Address().Hex(), "\n",
		"L2 Bridge Adapter:", l2BridgeAdapterAddress.Hex(),
	)

	return universe{
		L1: struct {
			Arm                  common.Address
			ArmProxy             common.Address
			TokenPool            common.Address
			Rebalancer           common.Address
			BridgeAdapterAddress common.Address
		}{
			Arm:                  l1Arm.Address(),
			ArmProxy:             l1ArmProxy.Address(),
			TokenPool:            l1TokenPool.Address(),
			Rebalancer:           l1Rebalancer.Address(),
			BridgeAdapterAddress: l1BridgeAdapterAddress,
		},
		L2: struct {
			Arm                  common.Address
			ArmProxy             common.Address
			TokenPool            common.Address
			Rebalancer           common.Address
			BridgeAdapterAddress common.Address
		}{
			Arm:                  l2Arm.Address(),
			ArmProxy:             l2ArmProxy.Address(),
			TokenPool:            l2TokenPool.Address(),
			Rebalancer:           l2Rebalancer.Address(),
			BridgeAdapterAddress: l2BridgeAdapterAddress,
		},
	}
}

func deployTokenPoolAndRebalancer(
	transactor *bind.TransactOpts,
	client *ethclient.Client,
	tokenAddress,
	armProxyAddress common.Address,
	chainID uint64,
	router common.Address,
) (*lock_release_token_pool.LockReleaseTokenPool, *rebalancer.Rebalancer) {
	_, tx, _, err := lock_release_token_pool.DeployLockReleaseTokenPool(
		transactor,
		client,
		tokenAddress,
		[]common.Address{},
		armProxyAddress,
		true,
		router,
	)
	helpers.PanicErr(err)
	tokenPoolAddress := helpers.ConfirmContractDeployed(context.Background(), client, tx, int64(chainID))
	tokenPool, err := lock_release_token_pool.NewLockReleaseTokenPool(tokenPoolAddress, client)
	helpers.PanicErr(err)

	_, tx, _, err = rebalancer.DeployRebalancer(
		transactor,
		client,
		tokenAddress,
		chainID,
		tokenPoolAddress,
	)
	helpers.PanicErr(err)
	rebalancerAddress := helpers.ConfirmContractDeployed(context.Background(), client, tx, int64(chainID))
	rebalancer, err := rebalancer.NewRebalancer(rebalancerAddress, client)
	helpers.PanicErr(err)
	tx, err = tokenPool.SetRebalancer(transactor, rebalancerAddress)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), client, tx, int64(chainID),
		"setting rebalancer on token pool")
	onchainRebalancer, err := tokenPool.GetRebalancer(nil)
	helpers.PanicErr(err)
	if onchainRebalancer != rebalancerAddress {
		panic(fmt.Sprintf("onchain rebalancer address does not match, expected %s got %s",
			rebalancerAddress.Hex(), onchainRebalancer.Hex()))
	}
	return tokenPool, rebalancer
}

func deployArm(
	transactor *bind.TransactOpts,
	client *ethclient.Client,
	chainID uint64) (*mock_arm_contract.MockARMContract, *arm_proxy_contract.ARMProxyContract) {
	_, tx, _, err := mock_arm_contract.DeployMockARMContract(transactor, client)
	helpers.PanicErr(err)
	armAddress := helpers.ConfirmContractDeployed(context.Background(), client, tx, int64(chainID))
	arm, err := mock_arm_contract.NewMockARMContract(armAddress, client)
	helpers.PanicErr(err)

	_, tx, _, err = arm_proxy_contract.DeployARMProxyContract(transactor, client, arm.Address())
	helpers.PanicErr(err)
	armProxyAddress := helpers.ConfirmContractDeployed(context.Background(), client, tx, int64(chainID))
	armProxy, err := arm_proxy_contract.NewARMProxyContract(armProxyAddress, client)
	helpers.PanicErr(err)

	return arm, armProxy
}

// sum of MaxDurationQuery/MaxDurationObservation/DeltaGrace must be less than DeltaProgress
func setConfig(
	e multienv.Env,
	args setConfigArgs,
) {
	validateEnv(e, args.l1ChainID, args.l2ChainID, false)

	l1Transactor, l2Transactor := e.Transactors[args.l1ChainID], e.Transactors[args.l2ChainID]

	// lengths of all the arrays must be equal
	if len(args.signers) != len(args.offchainPubKeys) ||
		len(args.signers) != len(args.configPubKeys) ||
		len(args.signers) != len(args.l1Transmitters) ||
		len(args.signers) != len(args.l2Transmitters) {
		panic("lengths of all the arrays must be equal")
	}

	l1Rebalancer, err := rebalancer.NewRebalancer(args.l1RebalancerAddress, e.Clients[args.l1ChainID])
	helpers.PanicErr(err)
	l2Rebalancer, err := rebalancer.NewRebalancer(args.l2RebalancerAddress, e.Clients[args.l2ChainID])
	helpers.PanicErr(err)

	// set config on L2 first then L1
	var (
		l1Oracles []confighelper2.OracleIdentityExtra
		l2Oracles []confighelper2.OracleIdentityExtra
	)
	for i := 0; i < len(args.signers); i++ {
		l1Oracles = append(l1Oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OffchainPublicKey: args.offchainPubKeys[i],
				OnchainPublicKey:  args.signers[i].Bytes(),
				PeerID:            args.peerIDs[i],
				TransmitAccount:   types.Account(args.l1Transmitters[i].Hex()),
			},
			ConfigEncryptionPublicKey: args.configPubKeys[i],
		})
		l2Oracles = append(l2Oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OffchainPublicKey: args.offchainPubKeys[i],
				OnchainPublicKey:  args.signers[i].Bytes(),
				PeerID:            args.peerIDs[i],
				TransmitAccount:   types.Account(args.l2Transmitters[i].Hex()),
			},
			ConfigEncryptionPublicKey: args.configPubKeys[i],
		})
	}
	var schedule []int
	for range l1Oracles {
		schedule = append(schedule, 1)
	}
	offchainConfig, onchainConfig := []byte{}, []byte{}
	f := uint8(1)
	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		args.deltaProgress,
		args.deltaResend,
		args.deltaInitial,
		args.deltaRound,
		args.deltaGrace,
		args.deltaCertifiedCommitRequest,
		args.deltaStage,
		args.rMax,
		schedule,
		l2Oracles,
		offchainConfig,
		args.maxDurationQuery,
		args.maxDurationObservation,
		args.maxDurationShouldAcceptAttestedReport,
		args.maxDurationShouldTransmitAcceptedReport,
		int(f),
		onchainConfig)
	helpers.PanicErr(err)
	tx, err := l2Rebalancer.SetOCR3Config(l2Transactor, args.signers, args.l2Transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Clients[args.l2ChainID], tx, int64(args.l2ChainID), "setting OCR3 config on L2 rebalancer")

	fmt.Println("sleeping a bit before setting config on L1")
	time.Sleep(1 * time.Minute)

	// set config on L1
	offchainConfig, onchainConfig = []byte{}, []byte{}
	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err = ocr3confighelper.ContractSetConfigArgsForTests(
		args.deltaProgress,
		args.deltaResend,
		args.deltaInitial,
		args.deltaRound,
		args.deltaGrace,
		args.deltaCertifiedCommitRequest,
		args.deltaStage,
		args.rMax,
		schedule,
		l1Oracles,
		offchainConfig,
		args.maxDurationQuery,
		args.maxDurationObservation,
		args.maxDurationShouldAcceptAttestedReport,
		args.maxDurationShouldTransmitAcceptedReport,
		int(f),
		onchainConfig)
	helpers.PanicErr(err)
	tx, err = l1Rebalancer.SetOCR3Config(l1Transactor, args.signers, args.l1Transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Clients[args.l1ChainID], tx, int64(args.l1ChainID), "setting OCR3 config on L1 rebalancer")
}

func validateEnv(env multienv.Env, l1ChainID, l2ChainID uint64, websocket bool) {
	_, ok := env.Clients[l1ChainID]
	if !ok {
		panic("L1 client not found")
	}
	_, ok = env.Clients[l2ChainID]
	if !ok {
		panic("L2 client not found")
	}
	_, ok = env.Transactors[l1ChainID]
	if !ok {
		panic("L1 transactor not found")
	}
	_, ok = env.Transactors[l2ChainID]
	if !ok {
		panic("L2 transactor not found")
	}
	if websocket {
		_, ok = env.WSURLs[l1ChainID]
		if !ok {
			panic("L1 websocket URL not found")
		}
		_, ok = env.WSURLs[l2ChainID]
		if !ok {
			panic("L2 websocket URL not found")
		}
	}
}
