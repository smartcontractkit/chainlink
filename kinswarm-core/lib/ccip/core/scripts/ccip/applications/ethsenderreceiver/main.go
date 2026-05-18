package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/liquiditymanager/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ether_sender_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

var (
	senderABI = abihelpers.MustParseABI(ether_sender_receiver.EtherSenderReceiverABI)
)

// nolint
func main() {
	switch os.Args[1] {
	case "inject-eth-liquidity":
		cmd := flag.NewFlagSet("inject-eth-liquidity", flag.ExitOnError)
		chainID := cmd.Uint64("chain-id", 0, "Chain ID")
		lrtpAddress := cmd.String("lrtp-address", "", "Lock/Release token pool address")
		amount := cmd.String("amount", "", "Amount to inject in wei")

		helpers.ParseArgs(cmd, os.Args[2:], "chain-id", "lrtp-address", "amount")

		env := multienv.New(false, false)
		lrtp, err := lock_release_token_pool.NewLockReleaseTokenPool(common.HexToAddress(*lrtpAddress), env.Clients[*chainID])
		helpers.PanicErr(err)

		balance, err := env.Clients[*chainID].BalanceAt(context.Background(), env.Transactors[*chainID].From, nil)
		helpers.PanicErr(err)

		if balance.Cmp(decimal.RequireFromString(*amount).BigInt()) < 0 {
			panic(fmt.Sprintf("Insufficient balance to inject: %s < %s, please get more ETH", balance, *amount))
		}

		wethAddress, err := lrtp.GetToken(nil)
		helpers.PanicErr(err)

		weth, err := weth9.NewWETH9(wethAddress, env.Clients[*chainID])
		helpers.PanicErr(err)

		tx, err := weth.Deposit(&bind.TransactOpts{
			From:   env.Transactors[*chainID].From,
			Signer: env.Transactors[*chainID].Signer,
			Value:  decimal.RequireFromString(*amount).BigInt(),
		})
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), env.Clients[*chainID], tx, int64(*chainID), "depositing ETH to get WETH, amount:", *amount, "wei")

		// token.Transfer the weth into the pool since we're not the owner
		tx, err = weth.Transfer(env.Transactors[*chainID], common.HexToAddress(*lrtpAddress), decimal.RequireFromString(*amount).BigInt())
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), env.Clients[*chainID], tx, int64(*chainID), "transferring WETH to LRTP, amount:", *amount, "wei")
	case "deploy":
		cmd := flag.NewFlagSet("deploy", flag.ExitOnError)
		chainID := cmd.Uint64("chain-id", 0, "Chain ID")
		routerAddress := cmd.String("router-address", "", "Router address")

		helpers.ParseArgs(cmd, os.Args[2:], "chain-id", "router-address")
		env := multienv.New(false, false)
		_, tx, _, err := ether_sender_receiver.DeployEtherSenderReceiver(
			env.Transactors[*chainID],
			env.Clients[*chainID],
			common.HexToAddress(*routerAddress),
		)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), env.Clients[*chainID], tx, int64(*chainID))
	case "get-fee":
		cmd := flag.NewFlagSet("get-fee", flag.ExitOnError)
		chainID := cmd.Uint64("chain-id", 0, "Chain ID")
		destChainID := cmd.Uint64("dest-chain-id", 0, "Destination chain ID")
		senderReceiverAddress := cmd.String("sender-receiver-address", "", "Sender receiver address")
		// message data
		destReceiver := cmd.String("dest-receiver-address", "", "Destination receiver address")
		destEOA := cmd.String("dest-eoa-address", "", "Destination EOA address")
		tokenAddress := cmd.String("token-address", "", "Token address")
		feeToken := cmd.String("fee-token", "", "Fee token address")
		amount := cmd.String("amount", "", "Amount")

		helpers.ParseArgs(cmd, os.Args[2:],
			"chain-id",
			"dest-chain-id",
			"sender-receiver-address",
			"dest-receiver-address",
			"dest-eoa-address",
			"token-address",
			"fee-token",
			"amount",
		)

		destChain, ok := chainsel.ChainByEvmChainID(*destChainID)
		if !ok {
			panic(fmt.Sprintf("Unknown chain ID: %d", *destChainID))
		}

		env := multienv.New(false, false)
		senderReceiver, err := ether_sender_receiver.NewEtherSenderReceiver(common.HexToAddress(*senderReceiverAddress), env.Clients[*chainID])
		helpers.PanicErr(err)

		receiverBytes, err := utils.ABIEncode(`[{"type": "address"}]`, common.HexToAddress(*destReceiver))
		helpers.PanicErr(err)
		destEOABytes, err := utils.ABIEncode(`[{"type": "address"}]`, common.HexToAddress(*destEOA))
		helpers.PanicErr(err)
		fmt.Println("receiver bytes:", hexutil.Encode(receiverBytes),
			"dest eoa bytes:", hexutil.Encode(destEOABytes), "fee token:", common.HexToAddress(*feeToken))

		msg := ether_sender_receiver.ClientEVM2AnyMessage{
			Receiver: receiverBytes,
			Data:     destEOABytes,
			TokenAmounts: []ether_sender_receiver.ClientEVMTokenAmount{
				{
					Token:  common.HexToAddress(*tokenAddress),
					Amount: decimal.RequireFromString(*amount).BigInt(),
				},
			},
			FeeToken: common.HexToAddress(*feeToken),
			// ExtraArgs: nil, // will be filled in by the contract
		}
		fee, err := senderReceiver.GetFee(nil, destChain.Selector, msg)
		helpers.PanicErr(err)

		fmt.Println("fee is:", fee, "juels/wei")
	case "ccip-send":
		cmd := flag.NewFlagSet("ccip-send", flag.ExitOnError)
		chainID := cmd.Uint64("chain-id", 0, "Chain ID")
		destChainID := cmd.Uint64("dest-chain-id", 0, "Destination chain ID")
		senderReceiverAddress := cmd.String("sender-receiver-address", "", "Sender receiver address")
		// message data
		destReceiver := cmd.String("dest-receiver-address", "", "Destination receiver address")
		feeToken := cmd.String("fee-token", "", "Fee token address")
		amount := cmd.String("amount", "", "Amount")
		gasLimit := cmd.Int64("gas-limit", 25_000, "Gas limit for the ccipReceive on destination chain")

		helpers.ParseArgs(cmd, os.Args[2:],
			"chain-id",
			"dest-chain-id",
			"sender-receiver-address",
			"dest-receiver-address",
			"fee-token",
			"amount",
		)

		destChain, ok := chainsel.ChainByEvmChainID(*destChainID)
		if !ok {
			panic(fmt.Sprintf("Unknown chain ID: %d", *destChainID))
		}

		env := multienv.New(false, false)
		senderReceiver, err := ether_sender_receiver.NewEtherSenderReceiver(common.HexToAddress(*senderReceiverAddress), env.Clients[*chainID])
		helpers.PanicErr(err)

		receiverBytes, err := utils.ABIEncode(`[{"type": "address"}]`, common.HexToAddress(*destReceiver))
		helpers.PanicErr(err)
		feeTok := common.HexToAddress(*feeToken)
		extraArgsV1Selector := hexutil.MustDecode("0x97a657c9")
		extraArgsV1, err := utils.ABIEncode(`[{"type": "uint256"}]`, big.NewInt(*gasLimit))
		helpers.PanicErr(err)
		extraArgsV1 = append(extraArgsV1Selector, extraArgsV1...)
		fmt.Println("extra args v1:", hexutil.Encode(extraArgsV1))
		msg := ether_sender_receiver.ClientEVM2AnyMessage{
			Receiver: receiverBytes,
			// Data: , // will be filled in by the contract
			TokenAmounts: []ether_sender_receiver.ClientEVMTokenAmount{
				{
					// Token: , // will be filled in by the contract
					Amount: decimal.RequireFromString(*amount).BigInt(),
				},
			},
			FeeToken:  feeTok,
			ExtraArgs: extraArgsV1,
		}
		fee, err := senderReceiver.GetFee(nil, destChain.Selector, msg)
		helpers.PanicErr(err)

		fmt.Println("fee is:", fee, "juels/wei")

		if (feeTok == common.Address{}) {
			totalValue := new(big.Int).Add(fee, decimal.RequireFromString(*amount).BigInt())
			ethBalance, err := env.Clients[*chainID].BalanceAt(context.Background(), env.Transactors[*chainID].From, nil)
			helpers.PanicErr(err)
			if ethBalance.Cmp(totalValue) < 0 {
				panic(fmt.Sprintf("Insufficient balance to send: %s < %s, please get more ETH", ethBalance, totalValue))
			}

			fmt.Println("Sending with value:", totalValue.String(), "total balance:", ethBalance.String())
			packed, err := senderABI.Pack("ccipSend", destChain.Selector, msg)
			helpers.PanicErr(err)
			nonce, err := env.Clients[*chainID].PendingNonceAt(context.Background(), env.Transactors[*chainID].From)
			helpers.PanicErr(err)
			gasPrice, err := env.Clients[*chainID].SuggestGasPrice(context.Background())
			helpers.PanicErr(err)
			toAddr := common.HexToAddress(*senderReceiverAddress)
			rawTx := types.NewTx(&types.LegacyTx{
				Nonce:    nonce,
				GasPrice: gasPrice,
				Gas:      500_000,
				To:       &toAddr,
				Value:    totalValue,
				Data:     packed,
			})
			signedTx, err := env.Transactors[*chainID].Signer(env.Transactors[*chainID].From, rawTx)
			helpers.PanicErr(err)
			err = env.Clients[*chainID].SendTransaction(context.Background(), signedTx)
			helpers.PanicErr(err)
			helpers.ConfirmTXMined(context.Background(), env.Clients[*chainID], signedTx, int64(*chainID), "ccip send native, msg value:", totalValue.String(), "fee:", fee.String())
		} else {
			// non-native fee token, so approve first then send.
			erc20Token, err := erc20.NewERC20(feeTok, env.Clients[*chainID])
			helpers.PanicErr(err)

			// check if we have enough to provide allowance
			balance, err := erc20Token.BalanceOf(nil, env.Transactors[*chainID].From)
			helpers.PanicErr(err)
			if balance.Cmp(fee) < 0 {
				panic(fmt.Sprintf("Insufficient balance to provide allowance: %s < %s, please get more tokens", balance, fee))
			}

			// approve
			tx, err := erc20Token.Approve(env.Transactors[*chainID], common.HexToAddress(*senderReceiverAddress), fee)
			helpers.PanicErr(err)
			helpers.ConfirmTXMined(context.Background(), env.Clients[*chainID], tx, int64(*chainID),
				"approving sender receiver to spend fee token, approval amount:", fee.String())

			// send message cross chain
			tx, err = senderReceiver.CcipSend(&bind.TransactOpts{
				From:   env.Transactors[*chainID].From,
				Signer: env.Transactors[*chainID].Signer,
				Value:  decimal.RequireFromString(*amount).BigInt(),
			}, destChain.Selector, msg)
			helpers.PanicErr(err)
			helpers.ConfirmTXMined(context.Background(), env.Clients[*chainID], tx, int64(*chainID), "ccip send non-native")
		}
	}
}
