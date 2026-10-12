package changeset

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	router "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/router"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

// genericExtraArgsV2Tag is the CCIP extraArgs v2 tag (GENERIC_EXTRA_ARGS_V2_TAG), followed by
// abi.encode(gasLimit uint256, allowOutOfOrderExecution bool).
const genericExtraArgsV2Tag = "0x181dcf10"

// proposal description labels for the test / remainder / full variants.
const (
	bridgeLabelTest      = "TEST"
	bridgeLabelRemainder = "REMAINDER"
	bridgeLabelFull      = "FULL"
)

// wethABI covers the WETH9 functions the bridge needs: deposit(), approve(), and the read calls used
// for balance preconditions. The proposal-decode ABI (WETH9) lives in the treasury ABI registry.
var wethABI = func() abi.ABI {
	const abiJSON = `[
		{"name":"deposit","type":"function","stateMutability":"payable","inputs":[],"outputs":[]},
		{"name":"approve","type":"function","stateMutability":"nonpayable","inputs":[{"name":"guy","type":"address"},{"name":"wad","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]},
		{"name":"balanceOf","type":"function","stateMutability":"view","inputs":[{"name":"","type":"address"}],"outputs":[{"name":"","type":"uint256"}]}
	]`
	a, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		panic(fmt.Sprintf("failed to parse WETH ABI: %v", err))
	}
	return a
}()

// routerABI is the parsed CCIP Router ABI, used to pack ccipSend calldata for the timelock batch.
var routerABI = func() abi.ABI {
	a, err := abi.JSON(strings.NewReader(router.RouterABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse Router ABI: %v", err))
	}
	return a
}()

// extraArgsV2Args encodes (gasLimit uint256, allowOutOfOrderExecution bool) for CCIP extraArgs v2.
var extraArgsV2Args = func() abi.Arguments {
	uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		panic(err)
	}
	boolTy, err := abi.NewType("bool", "", nil)
	if err != nil {
		panic(err)
	}
	return abi.Arguments{{Type: uint256Ty}, {Type: boolTy}}
}()

// addressArg encodes a single address (used for message.receiver and message.data).
var addressArg = func() abi.Arguments {
	addrTy, err := abi.NewType("address", "", nil)
	if err != nil {
		panic(err)
	}
	return abi.Arguments{{Type: addrTy}}
}()

// bufferedFee applies CCIPFeeBufferBps to a raw CCIP fee quote: fee * bps / 10000. For native fee,
// the Router consumes the whole msg.value, so the amount above the live fee is knowingly donated.
// This donation is acceptable as it indirectly goes back to Finance via the CCIP onRamps.
func bufferedFee(fee *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(fee, CCIPFeeBufferBps), big.NewInt(bpsDenominator))
}

func buildExtraArgsV2(gasLimit uint64) ([]byte, error) {
	enc, err := extraArgsV2Args.Pack(new(big.Int).SetUint64(gasLimit), true)
	if err != nil {
		return nil, fmt.Errorf("failed to encode extraArgs: %w", err)
	}
	return append(common.FromHex(genericExtraArgsV2Tag), enc...), nil
}

// buildCCIPMessage constructs the EVM2Any message: WETH as the bridged token, the destination
// EtherSenderReceiver as the receiver, the final recipient EOA abi-encoded in data, native fee token.
func buildCCIPMessage(weth, esr, recipient common.Address, amount *big.Int, gasLimit uint64) (router.ClientEVM2AnyMessage, error) {
	extraArgs, err := buildExtraArgsV2(gasLimit)
	if err != nil {
		return router.ClientEVM2AnyMessage{}, err
	}
	receiverBytes, err := addressArg.Pack(esr)
	if err != nil {
		return router.ClientEVM2AnyMessage{}, fmt.Errorf("failed to encode receiver: %w", err)
	}
	dataBytes, err := addressArg.Pack(recipient)
	if err != nil {
		return router.ClientEVM2AnyMessage{}, fmt.Errorf("failed to encode data: %w", err)
	}
	return router.ClientEVM2AnyMessage{
		Receiver:     receiverBytes,
		Data:         dataBytes,
		TokenAmounts: []router.ClientEVMTokenAmount{{Token: weth, Amount: amount}},
		FeeToken:     common.Address{}, // address(0) => native fee
		ExtraArgs:    extraArgs,
	}, nil
}

// BridgeCCIPSequenceInput is the input for the bridge sequence.
type BridgeCCIPSequenceInput struct {
	Config types.CCIPBridgeConfig `json:"config"`
}

// BridgeCCIPSequenceOutput is the output of the bridge sequence.
type BridgeCCIPSequenceOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal `json:"mcms_timelock_proposals,omitempty"`
	Description           string                  `json:"description"`
}

// BridgeCCIPSequence builds one MCMS timelock proposal (or two, in test mode) that bridge
// native/wrapped-native value from the source Timelock via CCIP.
var BridgeCCIPSequence = operations.NewSequence(
	"bridge-ccip-sequence",
	semver.MustParse("1.0.0"),
	"Builds MCMS timelock proposals to bridge native/wrapped-native value via CCIP",
	func(b operations.Bundle, deps VaultDeps, input BridgeCCIPSequenceInput) (BridgeCCIPSequenceOutput, error) {
		cfg := input.Config
		out := BridgeCCIPSequenceOutput{Description: cfg.Description}

		if cfg.TestNeeded {
			testTransfers, remainderTransfers, err := splitTransfersForTest(cfg.Transfers)
			if err != nil {
				return BridgeCCIPSequenceOutput{}, err
			}
			testProposal, err := buildBridgeProposal(b, deps, cfg, testTransfers, bridgeLabelTest)
			if err != nil {
				return BridgeCCIPSequenceOutput{}, fmt.Errorf("test proposal: %w", err)
			}
			remainderProposal, err := buildBridgeProposal(b, deps, cfg, remainderTransfers, bridgeLabelRemainder)
			if err != nil {
				return BridgeCCIPSequenceOutput{}, fmt.Errorf("remainder proposal: %w", err)
			}
			out.MCMSTimelockProposals = []mcms.TimelockProposal{*testProposal, *remainderProposal}
			return out, nil
		}

		proposal, err := buildBridgeProposal(b, deps, cfg, cfg.Transfers, bridgeLabelFull)
		if err != nil {
			return BridgeCCIPSequenceOutput{}, err
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
		return out, nil
	},
)

// splitTransfersForTest returns a test slice (each amount = TestBridgeAmount) and a remainder slice
// (each amount = amount - TestBridgeAmount). Requires amount > TestBridgeAmount for every transfer.
func splitTransfersForTest(transfers []types.CCIPBridgeTransfer) (test, remainder []types.CCIPBridgeTransfer, err error) {
	for _, t := range transfers {
		if t.Amount == nil || t.Amount.Cmp(TestBridgeAmount) <= 0 {
			return nil, nil, fmt.Errorf("dest %d: amount must be greater than test amount %s", t.ToChainSelector, TestBridgeAmount)
		}
		testT := t
		testT.Amount = new(big.Int).Set(TestBridgeAmount)

		remT := t
		remT.Amount = new(big.Int).Sub(t.Amount, TestBridgeAmount)

		test = append(test, testT)
		remainder = append(remainder, remT)
	}
	return test, remainder, nil
}

// buildBridgeProposal resolves source contracts, quotes CCIP fees, performs balance preconditions,
// assembles the (deposit?/approve/ccipSend) timelock batch on the source chain, and builds the proposal.
func buildBridgeProposal(
	b operations.Bundle,
	deps VaultDeps,
	cfg types.CCIPBridgeConfig,
	transfers []types.CCIPBridgeTransfer,
	label string,
) (*mcms.TimelockProposal, error) {
	ctx := b.GetContext()
	src := cfg.FromChainSelector
	sourceChain := deps.Chain

	timelockAddr, err := GetContractAddress(deps.DataStore, src, commontypes.RBACTimelock)
	if err != nil {
		return nil, fmt.Errorf("timelock not found for source chain %d: %w", src, err)
	}
	mcmType := commontypes.ProposerManyChainMultisig
	if cfg.MCMSConfig.MCMSAction == mcmstypes.TimelockActionBypass {
		mcmType = commontypes.BypasserManyChainMultisig
	}
	mcmAddr, err := GetContractAddress(deps.DataStore, src, mcmType)
	if err != nil {
		return nil, fmt.Errorf("%s not found for source chain %d: %w", mcmType, src, err)
	}
	wethAddrStr, err := GetContractAddress(deps.DataStore, src, WETHContractType)
	if err != nil {
		return nil, fmt.Errorf("WETH not found for source chain %d: %w", src, err)
	}
	routerAddrStr, err := GetContractAddress(deps.DataStore, src, RouterContractType)
	if err != nil {
		return nil, fmt.Errorf("Router not found for source chain %d: %w", src, err)
	}

	wethAddr := common.HexToAddress(wethAddrStr)
	routerAddr := common.HexToAddress(routerAddrStr)
	timelockAddress := common.HexToAddress(timelockAddr)

	routerContract, err := router.NewRouter(routerAddr, sourceChain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind router %s: %w", routerAddr.Hex(), err)
	}

	gasLimit := cfg.GasLimit
	if gasLimit == 0 {
		gasLimit = defaultBridgeGasLimit
	}

	var txs []mcmstypes.Transaction
	totalFee := big.NewInt(0)
	totalAmount := big.NewInt(0)

	for _, t := range transfers {
		esrStr, err := GetContractAddress(deps.DataStore, t.ToChainSelector, EtherSenderReceiverContractType)
		if err != nil {
			return nil, fmt.Errorf("EtherSenderReceiver not found for dest chain %d: %w", t.ToChainSelector, err)
		}
		esrAddr := common.HexToAddress(esrStr)
		recipient := common.HexToAddress(t.DestinationAddress)

		msg, err := buildCCIPMessage(wethAddr, esrAddr, recipient, t.Amount, gasLimit)
		if err != nil {
			return nil, err
		}

		supported, err := routerContract.IsChainSupported(&bind.CallOpts{Context: ctx}, t.ToChainSelector)
		if err != nil {
			return nil, fmt.Errorf("isChainSupported(%d): %w", t.ToChainSelector, err)
		}
		if !supported {
			return nil, fmt.Errorf("router on source chain %d does not support destination %d", src, t.ToChainSelector)
		}

		fee, err := routerContract.GetFee(&bind.CallOpts{Context: ctx}, t.ToChainSelector, msg)
		if err != nil {
			return nil, fmt.Errorf("getFee(%d): %w", t.ToChainSelector, err)
		}
		feeWithBuffer := bufferedFee(fee)
		if t.MaxCCIPFee != nil && feeWithBuffer.Cmp(t.MaxCCIPFee) > 0 {
			return nil, fmt.Errorf("dest %d: buffered fee %s exceeds max_ccip_fee %s", t.ToChainSelector, feeWithBuffer, t.MaxCCIPFee)
		}

		// native: wrap the amount first (Timelock holds native). wrappedNative: skip (holds WETH).
		if cfg.TokenType == types.TokenTypeNative {
			depositData, err := wethABI.Pack("deposit")
			if err != nil {
				return nil, fmt.Errorf("failed to pack deposit: %w", err)
			}
			depositTx, err := cldfproposalutils.TransactionForChain(src, wethAddr.Hex(), depositData, t.Amount, "WETHDeposit", []string{"vault", "bridge", "wrap"})
			if err != nil {
				return nil, fmt.Errorf("failed to build deposit tx: %w", err)
			}
			txs = append(txs, depositTx)
		}

		approveData, err := wethABI.Pack("approve", routerAddr, t.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to pack approve: %w", err)
		}
		approveTx, err := cldfproposalutils.TransactionForChain(src, wethAddr.Hex(), approveData, big.NewInt(0), "WETHApprove", []string{"vault", "bridge", "approve"})
		if err != nil {
			return nil, fmt.Errorf("failed to build approve tx: %w", err)
		}
		txs = append(txs, approveTx)

		sendData, err := routerABI.Pack("ccipSend", t.ToChainSelector, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to pack ccipSend: %w", err)
		}
		sendTx, err := cldfproposalutils.TransactionForChain(src, routerAddr.Hex(), sendData, feeWithBuffer, "CCIPSend", []string{"vault", "bridge", "ccip-send"})
		if err != nil {
			return nil, fmt.Errorf("failed to build ccipSend tx: %w", err)
		}
		txs = append(txs, sendTx)

		totalFee.Add(totalFee, feeWithBuffer)
		totalAmount.Add(totalAmount, t.Amount)
	}

	if err := assertTimelockBalances(ctx, sourceChain.Client, cfg.TokenType, timelockAddress, wethAddr, totalAmount, totalFee, src); err != nil {
		return nil, err
	}

	timelockAddressByChain := map[uint64]string{src: timelockAddr}
	mcmAddressByChain := map[uint64]string{src: mcmAddr}
	inspectorPerChain := map[uint64]mcmssdk.Inspector{src: mcmsevmsdk.NewInspector(sourceChain.Client)}
	batches := []mcmstypes.BatchOperation{{
		ChainSelector: mcmstypes.ChainSelector(src),
		Transactions:  txs,
	}}

	description := cfg.Description
	if description == "" {
		description = "Vault CCIP Bridge"
	}
	description = fmt.Sprintf("%s [%s]", description, label)

	proposal, err := proposeutils.BuildProposalFromBatchesV2(
		deps.Environment,
		timelockAddressByChain,
		mcmAddressByChain,
		inspectorPerChain,
		batches,
		description,
		*cfg.MCMSConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build MCMS proposal: %w", err)
	}

	b.Logger.Infow("Built CCIP bridge proposal",
		"label", label,
		"source", src,
		"destinations", len(transfers),
		"total_amount_wei", totalAmount.String(),
		"total_buffered_fee_wei", totalFee.String(),
		"operations", len(proposal.Operations))

	return proposal, nil
}

// bridgeBalanceClient is the minimal read interface the balance preconditions need: native balance
// plus the contract-call surface (bind.ContractCaller) used for WETH.balanceOf.
type bridgeBalanceClient interface {
	bind.ContractCaller
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

// assertTimelockBalances verifies the source Timelock can cover the planned bridge:
//   - native:        native >= sum of (amount) + sum of(fee) (wrap consumes native, fee paid in native)
//   - wrappedNative: WETH   >= sum of(amount) AND native >= sum of(fee) (fee still paid in native)
func assertTimelockBalances(
	ctx context.Context,
	client bridgeBalanceClient,
	tokenType string,
	timelock, weth common.Address,
	totalAmount, totalFee *big.Int,
	src uint64,
) error {
	nativeBal, err := client.BalanceAt(ctx, timelock, nil)
	if err != nil {
		return fmt.Errorf("failed to read timelock native balance on chain %d: %w", src, err)
	}

	if tokenType == types.TokenTypeNative {
		need := new(big.Int).Add(totalAmount, totalFee)
		if nativeBal.Cmp(need) < 0 {
			return fmt.Errorf("timelock native balance %s < required %s (amounts + fees) on chain %d", nativeBal, need, src)
		}
		return nil
	}

	// wrappedNative
	if nativeBal.Cmp(totalFee) < 0 {
		return fmt.Errorf("timelock native balance %s < required fees %s on chain %d", nativeBal, totalFee, src)
	}
	wethBal, err := readWETHBalance(ctx, client, weth, timelock)
	if err != nil {
		return fmt.Errorf("failed to read timelock WETH balance on chain %d: %w", src, err)
	}
	if wethBal.Cmp(totalAmount) < 0 {
		return fmt.Errorf("timelock WETH balance %s < required %s on chain %d", wethBal, totalAmount, src)
	}
	return nil
}

// readWETHBalance reads WETH.balanceOf(owner) via a bound contract call.
func readWETHBalance(ctx context.Context, backend bind.ContractCaller, weth, owner common.Address) (*big.Int, error) {
	bc := bind.NewBoundContract(weth, wethABI, backend, nil, nil)
	var out []interface{}
	if err := bc.Call(&bind.CallOpts{Context: ctx}, &out, "balanceOf", owner); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("empty balanceOf result")
	}
	bal, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected balanceOf return type %T", out[0])
	}
	return bal, nil
}
