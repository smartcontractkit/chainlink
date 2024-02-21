package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rebalancer/arb"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rebalancer/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/arbitrum_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/arbitrum_l2_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
)

type setConfigArgs struct {
	l1ChainID                               uint64
	l2ChainID                               uint64
	l1RebalancerAddress                     common.Address
	l2RebalancerAddress                     common.Address
	signers                                 []common.Address
	offchainPubKeys                         []types.OffchainPublicKey
	configPubKeys                           []types.ConfigEncryptionPublicKey
	peerIDs                                 []string
	l1Transmitters                          []common.Address
	l2Transmitters                          []common.Address
	deltaProgress                           time.Duration
	deltaResend                             time.Duration
	deltaInitial                            time.Duration
	deltaRound                              time.Duration
	deltaGrace                              time.Duration
	deltaCertifiedCommitRequest             time.Duration
	deltaStage                              time.Duration
	rMax                                    uint64
	maxDurationQuery                        time.Duration
	maxDurationObservation                  time.Duration
	maxDurationShouldAcceptAttestedReport   time.Duration
	maxDurationShouldTransmitAcceptedReport time.Duration
	f                                       int
}

func main() {
	switch os.Args[1] {
	case "deploy-universe":
		cmd := flag.NewFlagSet("deploy-universe", flag.ExitOnError)
		l1ChainID := cmd.Uint64("l1-chain-id", 0, "L1 Chain ID")
		l2ChainID := cmd.Uint64("l2-chain-id", 0, "L2 Chain ID")
		l1TokenAddress := cmd.String("l1-token-address", "", "L1 Token Address")
		l2TokenAddress := cmd.String("l2-token-address", "", "L2 Token Address")

		helpers.ParseArgs(cmd, os.Args[2:], "l1-chain-id", "l2-chain-id")
		deployUniverse(
			multienv.New(false, false),
			*l1ChainID,
			*l2ChainID,
			common.HexToAddress(*l1TokenAddress),
			common.HexToAddress(*l2TokenAddress),
		)
	case "set-config":
		cmd := flag.NewFlagSet("set-config", flag.ExitOnError)
		l1ChainID := cmd.Uint64("l1-chain-id", 0, "L1 Chain ID")
		l2ChainID := cmd.Uint64("l2-chain-id", 0, "L2 Chain ID")
		l1RebalancerAddress := cmd.String("l1-rebalancer-address", "", "L1 Rebalancer Address")
		l2RebalancerAddress := cmd.String("l2-rebalancer-address", "", "L2 Rebalancer Address")
		// OCR information
		signers := cmd.String("signers", "", "comma separated list of OCR signers (onchain public keys)")
		offchainPubKeys := cmd.String("offchain-pubkeys", "", "comma separated list of OCR3 offchain pubkeys")
		configPubKeys := cmd.String("config-pubkeys", "", "comma separated list of OCR3 config pubkeys")
		peerIDs := cmd.String("peer-ids", "", "comma separated list of OCR3 peer IDs")
		l1Transmitters := cmd.String("l1-transmitters", "", "comma separated list of l1 transmitters")
		l2Transmitters := cmd.String("l2-transmitters", "", "comma separated list of l2 transmitters")
		deltaProgress := cmd.Duration("delta-progress", 2*time.Minute, "delta progress")
		deltaResend := cmd.Duration("delta-resend", 2*time.Minute, "delta resend")
		deltaInitial := cmd.Duration("delta-initial", 20*time.Second, "delta initial")
		deltaRound := cmd.Duration("delta-round", 2*time.Second, "delta round")
		deltaGrace := cmd.Duration("delta-grace", 20*time.Second, "delta grace")
		deltaCertifiedCommitRequest := cmd.Duration("delta-certified-commit-request", 10*time.Second, "delta certified commit request")
		deltaStage := cmd.Duration("delta-stage", 40*time.Second, "delta stage")
		rMax := cmd.Uint64("r-max", 3, "r max")
		maxDurationQuery := cmd.Duration("max-duration-query", 50*time.Millisecond, "max duration query")
		maxDurationObservation := cmd.Duration("max-duration-observation", 1*time.Minute, "max duration observation")
		maxDurationShouldAcceptAttestedReport := cmd.Duration("max-duration-should-accept-attested-report", 1*time.Minute, "max duration should accept attested report")
		maxDurationShouldTransmitAcceptedReport := cmd.Duration("max-duration-should-transmit-accepted-report", 1*time.Second, "max duration should transmit accepted report")
		f := cmd.Int("f", 1, "f")

		helpers.ParseArgs(cmd, os.Args[2:],
			"l1-chain-id",
			"l2-chain-id",
			"l1-rebalancer-address",
			"l2-rebalancer-address",
			"signers",
			"offchain-pubkeys",
			"config-pubkeys",
			"l1-transmitters",
			"l2-transmitters",
		)

		args := setConfigArgs{
			l1ChainID:                               *l1ChainID,
			l2ChainID:                               *l2ChainID,
			l1RebalancerAddress:                     common.HexToAddress(*l1RebalancerAddress),
			l2RebalancerAddress:                     common.HexToAddress(*l2RebalancerAddress),
			signers:                                 parseOnchainPubKeys(*signers),
			offchainPubKeys:                         parseOffchainPubKeys(*offchainPubKeys),
			configPubKeys:                           parseConfigPubKeys(*configPubKeys),
			peerIDs:                                 strings.Split(*peerIDs, ","),
			l1Transmitters:                          helpers.ParseAddressSlice(*l1Transmitters),
			l2Transmitters:                          helpers.ParseAddressSlice(*l2Transmitters),
			deltaProgress:                           *deltaProgress,
			deltaResend:                             *deltaResend,
			deltaInitial:                            *deltaInitial,
			deltaRound:                              *deltaRound,
			deltaGrace:                              *deltaGrace,
			deltaCertifiedCommitRequest:             *deltaCertifiedCommitRequest,
			deltaStage:                              *deltaStage,
			rMax:                                    *rMax,
			maxDurationQuery:                        *maxDurationQuery,
			maxDurationObservation:                  *maxDurationObservation,
			maxDurationShouldAcceptAttestedReport:   *maxDurationShouldAcceptAttestedReport,
			maxDurationShouldTransmitAcceptedReport: *maxDurationShouldTransmitAcceptedReport,
			f:                                       *f,
		}

		setConfig(
			multienv.New(false, false),
			args,
		)
	case "setup-rebalancer-nodes":
		setupRebalancerNodes(multienv.New(true, true))
	case "fund-contracts":
		cmd := flag.NewFlagSet("fund-contracts", flag.ExitOnError)
		l1ChainID := cmd.Uint64("l1-chain-id", 0, "L1 Chain ID")
		l2ChainID := cmd.Uint64("l2-chain-id", 0, "L2 Chain ID")
		l1RebalancerAddress := cmd.String("l1-rebalancer-address", "", "L1 Rebalancer Address")
		l2RebalancerAddress := cmd.String("l2-rebalancer-address", "", "L2 Rebalancer Address")
		l1TokenAddress := cmd.String("l1-token-address", "", "L1 Token Address")
		l2TokenAddress := cmd.String("l2-token-address", "", "L2 Token Address")
		l1TokenPoolAddress := cmd.String("l1-token-pool-address", "", "L1 Token Pool Address")
		l2TokenPoolAddress := cmd.String("l2-token-pool-address", "", "L2 Token Pool Address")
		l1TokenPoolAmount := cmd.String("l1-token-pool-amount", "1000000000", "L1 Token Pool Amount")      // 1 Gwei
		l2TokenPoolAmount := cmd.String("l2-token-pool-amount", "1000000000", "L2 Token Pool Amount")      // 1 Gwei
		l1RebalancerAmount := cmd.String("l1-rebalancer-amount", "70000000000000000", "Rebalancer Amount") // 0.07 ETH
		l2RebalancerAmount := cmd.String("l2-rebalancer-amount", "0", "Rebalancer Amount")

		helpers.ParseArgs(cmd, os.Args[2:],
			"l2-chain-id",
			"l2-rebalancer-address",
			"l2-token-address",
			"l2-token-pool-address",
		)

		env := multienv.New(false, false)
		fundPoolAndRebalancer(
			env,
			*l1ChainID,
			common.HexToAddress(*l1TokenAddress),
			common.HexToAddress(*l1TokenPoolAddress),
			common.HexToAddress(*l1RebalancerAddress),
			decimal.RequireFromString(*l1TokenPoolAmount).BigInt(),
			decimal.RequireFromString(*l1RebalancerAmount).BigInt(),
		)
		fundPoolAndRebalancer(
			env,
			*l2ChainID,
			common.HexToAddress(*l2TokenAddress),
			common.HexToAddress(*l2TokenPoolAddress),
			common.HexToAddress(*l2RebalancerAddress),
			decimal.RequireFromString(*l2TokenPoolAmount).BigInt(),
			decimal.RequireFromString(*l2RebalancerAmount).BigInt(),
		)
	case "get-cross-chain-rebalancers":
		cmd := flag.NewFlagSet("get-cross-chain-rebalancers", flag.ExitOnError)
		chainID := cmd.Uint64("chain-id", 0, "Chain ID")
		rebalancerAddress := cmd.String("rebalancer-address", "", "Rebalancer Address")

		helpers.ParseArgs(cmd, os.Args[2:], "chain-id", "rebalancer-address")

		env := multienv.New(false, false)
		client, ok := env.Clients[*chainID]
		if !ok {
			panic("client for chain id not found, please set appropriate env vars")
		}

		rebal, err := rebalancer.NewRebalancer(common.HexToAddress(*rebalancerAddress), client)
		helpers.PanicErr(err)

		xchainRebalancers, err := rebal.GetAllCrossChainRebalancers(nil)
		helpers.PanicErr(err)
		for _, xchainRebal := range xchainRebalancers {
			fmt.Println("Remote rebalancer address:", xchainRebal.RemoteRebalancer.Hex(), "\n",
				"Remote chain ID:", xchainRebal.RemoteChainSelector, "\n",
				"Remote token address:", xchainRebal.RemoteToken.Hex(), "\n",
				"Local bridge:", xchainRebal.LocalBridge.Hex(), "\n",
				"Enabled:", xchainRebal.Enabled,
			)
			fmt.Println()
		}
	case "arb-finalize-l1":
		cmd := flag.NewFlagSet("arb-finalize-l1", flag.ExitOnError)
		l1ChainID := cmd.Uint64("l1-chain-id", 0, "L1 Chain ID")
		l2ChainID := cmd.Uint64("l2-chain-id", 0, "L2 Chain ID")
		l2TxHash := cmd.String("l2-tx-hash", "", "L2 Tx Hash")
		l1BridgeAdapterAddress := cmd.String("l1-bridge-adapter-address", "", "L1 Bridge Adapter Address")

		helpers.ParseArgs(cmd, os.Args[2:], "l1-chain-id", "l2-chain-id", "l2-tx-hash", "l1-bridge-adapter-address")

		env := multienv.New(false, false)

		arb.FinalizeL1(
			env,
			*l1ChainID,
			*l2ChainID,
			common.HexToAddress(*l1BridgeAdapterAddress),
			common.HexToHash(*l2TxHash))
	case "deploy-arb-l1-adapter":
		cmd := flag.NewFlagSet("deploy-arb-l1-adapter", flag.ExitOnError)
		l1ChainID := cmd.Uint64("l1-chain-id", 0, "L1 Chain ID")
		helpers.ParseArgs(cmd, os.Args[2:], "l1-chain-id")

		env := multienv.New(false, false)
		l1GatewayRouter := arb.ArbitrumContracts[*l1ChainID]["L1GatewayRouter"]
		l1Outbox := arb.ArbitrumContracts[*l1ChainID]["L1Outbox"]
		l2Client := env.Clients[*l1ChainID]
		_, tx, _, err := arbitrum_l1_bridge_adapter.DeployArbitrumL1BridgeAdapter(env.Transactors[*l1ChainID], l2Client, l1GatewayRouter, l1Outbox)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), l2Client, tx, int64(*l1ChainID))
	case "deploy-arb-l2-adapter":
		cmd := flag.NewFlagSet("deploy-arb-l2-adapter", flag.ExitOnError)
		l2ChainID := cmd.Uint64("l2-chain-id", 0, "L2 Chain ID")
		helpers.ParseArgs(cmd, os.Args[2:], "l2-chain-id")

		env := multienv.New(false, false)
		l2GatewayRouter := arb.ArbitrumContracts[*l2ChainID]["L2GatewayRouter"]
		l2Client := env.Clients[*l2ChainID]
		_, tx, _, err := arbitrum_l2_bridge_adapter.DeployArbitrumL2BridgeAdapter(env.Transactors[*l2ChainID], l2Client, l2GatewayRouter)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), l2Client, tx, int64(*l2ChainID))
	case "arb-withdraw-from-l2":
		cmd := flag.NewFlagSet("arb-withdraw-from-l2", flag.ExitOnError)
		l2ChainID := cmd.Uint64("l2-chain-id", 0, "L2 Chain ID")
		l2BridgeAdapterAddress := cmd.String("l2-bridge-adapter-address", "", "L2 Bridge Adapter Address")
		amount := cmd.String("amount", "1", "Amount")
		l1ToAddress := cmd.String("l1-to-address", "", "L1 Address")
		l2TokenAddress := cmd.String("l2-token-address", "", "Token Address")
		l1TokenAddress := cmd.String("l1-token-address", "", "L1 Token Address")

		helpers.ParseArgs(cmd, os.Args[2:],
			"l2-chain-id", "l2-bridge-adapter-address", "l1-to-address", "l2-token-address", "l1-token-address")

		env := multienv.New(false, false)
		arb.WithdrawFromL2(
			env,
			*l2ChainID,
			common.HexToAddress(*l2BridgeAdapterAddress),
			decimal.RequireFromString(*amount).BigInt(),
			common.HexToAddress(*l1ToAddress),
			common.HexToAddress(*l2TokenAddress),
			common.HexToAddress(*l1TokenAddress),
		)
	case "arb-send-to-l2":
		cmd := flag.NewFlagSet("arb-send-to-l2", flag.ExitOnError)
		l1ChainID := cmd.Uint64("l1-chain-id", 0, "L1 Chain ID")
		l2ChainID := cmd.Uint64("l2-chain-id", 0, "L2 Chain ID")
		l1BridgeAdapterAddress := cmd.String("l1-bridge-adapter-address", "", "L1 Bridge Adapter Address")
		amount := cmd.String("amount", "1", "Amount")
		l2ToAddress := cmd.String("l2-to-address", "", "L2 Address")
		l1TokenAddress := cmd.String("l1-token-address", "", "L1 Token Address")
		l1RefundAddress := cmd.String("l1-refund-address", "", "L1 Refund Address")
		l2RefundAddress := cmd.String("l2-refund-address", "", "L2 Refund Address")

		helpers.ParseArgs(cmd, os.Args[2:],
			"l1-chain-id", "l2-chain-id", "l1-bridge-adapter-address", "l2-to-address", "l1-token-address")

		env := multienv.New(false, false)
		arb.SendToL2(
			env,
			*l1ChainID,
			*l2ChainID,
			common.HexToAddress(*l1BridgeAdapterAddress),
			common.HexToAddress(*l1TokenAddress),
			common.HexToAddress(*l1RefundAddress),
			common.HexToAddress(*l2RefundAddress),
			common.HexToAddress(*l2ToAddress),
			decimal.RequireFromString(*amount).BigInt(),
		)
	case "deposit-weth":
		cmd := flag.NewFlagSet("deposit-weth", flag.ExitOnError)
		chainID := cmd.Uint64("chain-id", 0, "Chain ID")
		amount := cmd.String("amount", "1000000000000000000", "Amount")
		wethAddress := cmd.String("weth-address", "", "WETH Address")
		helpers.ParseArgs(cmd, os.Args[2:], "chain-id")

		env := multienv.New(false, false)
		weth, err := weth9.NewWETH9(common.HexToAddress(*wethAddress), env.Clients[*chainID])
		helpers.PanicErr(err)

		tx, err := weth.Deposit(&bind.TransactOpts{
			From:   env.Transactors[*chainID].From,
			Signer: env.Transactors[*chainID].Signer,
			Value:  decimal.RequireFromString(*amount).BigInt(),
		})
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), env.Clients[*chainID], tx, int64(*chainID))
	case "transfer-weth":
		cmd := flag.NewFlagSet("transfer-weth", flag.ExitOnError)
		chainID := cmd.Uint64("chain-id", 0, "Chain ID")
		amount := cmd.String("amount", "1000000000000000000", "Amount")
		wethAddress := cmd.String("weth-address", "", "WETH Address")
		toAddress := cmd.String("to-address", "", "To Address")
		helpers.ParseArgs(cmd, os.Args[2:], "chain-id", "weth-address", "to-address")

		env := multienv.New(false, false)
		weth, err := weth9.NewWETH9(common.HexToAddress(*wethAddress), env.Clients[*chainID])
		helpers.PanicErr(err)

		tx, err := weth.Transfer(env.Transactors[*chainID], common.HexToAddress(*toAddress), decimal.RequireFromString(*amount).BigInt())
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), env.Clients[*chainID], tx, int64(*chainID))
	}
}

func parseOnchainPubKeys(onchainPubKeys string) []common.Address {
	split := strings.Split(onchainPubKeys, ",")
	ocrPubKeys := make([]common.Address, len(split))
	for i, key := range split {
		decoded, err := hex.DecodeString(key)
		helpers.PanicErr(err)
		ocrPubKeys[i] = common.BytesToAddress(decoded)
	}
	return ocrPubKeys
}

func parseOffchainPubKeys(offchainPubKeys string) []types.OffchainPublicKey {
	split := strings.Split(offchainPubKeys, ",")
	ocrPubKeys := make([]types.OffchainPublicKey, len(split))
	for i, key := range split {
		k, err := hex.DecodeString(key)
		helpers.PanicErr(err)
		ocrPubKeys[i] = types.OffchainPublicKey(k)
	}
	return ocrPubKeys
}

func parseConfigPubKeys(configPubKeys string) []types.ConfigEncryptionPublicKey {
	split := strings.Split(configPubKeys, ",")
	ocrPubKeys := make([]types.ConfigEncryptionPublicKey, len(split))
	for i, key := range split {
		k, err := hex.DecodeString(key)
		helpers.PanicErr(err)
		ocrPubKeys[i] = types.ConfigEncryptionPublicKey(k)
	}
	return ocrPubKeys
}
