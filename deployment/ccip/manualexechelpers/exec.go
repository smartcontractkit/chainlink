package manualexechelpers

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm/manualexeclib"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
)

// 14 days is the default lookback but it can be overridden.
const DefaultLookback = 14 * 24 * time.Hour

var (
	blockTimeSecondsPerChain = map[uint64]uint64{
		// simchains
		chainsel.GETH_TESTNET.Selector:  1,
		chainsel.GETH_DEVNET_2.Selector: 1,
		chainsel.GETH_DEVNET_3.Selector: 1,

		// arb
		chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector:         1,
		chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector: 1,

		// op
		chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector:         1,
		chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector: 1,

		// base
		chainsel.ETHEREUM_MAINNET_BASE_1.Selector:         1,
		chainsel.ETHEREUM_TESTNET_SEPOLIA_BASE_1.Selector: 1,

		// matic
		chainsel.POLYGON_MAINNET.Selector:        2,
		chainsel.POLYGON_TESTNET_MUMBAI.Selector: 2,
		chainsel.POLYGON_TESTNET_AMOY.Selector:   2,

		// bsc
		chainsel.BINANCE_SMART_CHAIN_MAINNET.Selector: 2,
		chainsel.BINANCE_SMART_CHAIN_TESTNET.Selector: 2,

		// eth
		chainsel.ETHEREUM_MAINNET.Selector:         12,
		chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: 12,
		chainsel.ETHEREUM_TESTNET_HOLESKY.Selector: 12,

		// avax
		chainsel.AVALANCHE_MAINNET.Selector:      3,
		chainsel.AVALANCHE_TESTNET_FUJI.Selector: 3,
	}
	// used if the chain isn't in the map above.
	defaultBlockTimeSeconds uint64 = 2
)

// getStartBlock gets the starting block of a filter logs query based on the current head block and a lookback duration.
// block time is used to calculate the number of blocks to go back.
func getStartBlock(srcChainSel uint64, currentHead uint64, lookbackDuration time.Duration) uint64 {
	blockTimeSeconds := blockTimeSecondsPerChain[srcChainSel]
	if blockTimeSeconds == 0 {
		blockTimeSeconds = defaultBlockTimeSeconds
	}

	toSub := uint64(lookbackDuration.Seconds()) / blockTimeSeconds
	if toSub > currentHead {
		return 1 // start from genesis - might happen for simchains.
	}

	start := currentHead - toSub
	return start
}

func getCommitRootAcceptedEvent(
	ctx context.Context,
	lggr logger.Logger,
	env deployment.Environment,
	state ccipchangeset.CCIPOnChainState,
	srcChainSel uint64,
	destChainSel uint64,
	msgSeqNr uint64,
	lookbackDuration time.Duration,
) (offramp.InternalMerkleRoot, error) {
	hdr, err := env.Chains[destChainSel].Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return offramp.InternalMerkleRoot{}, fmt.Errorf("failed to get header: %w", err)
	}

	start := getStartBlock(srcChainSel, hdr.Number.Uint64(), lookbackDuration)
	lggr.Debugw("Getting commit root accepted event", "startBlock", start)
	iter, err := state.Chains[destChainSel].OffRamp.FilterCommitReportAccepted(
		&bind.FilterOpts{
			Start: start,
		},
	)
	if err != nil {
		return offramp.InternalMerkleRoot{}, fmt.Errorf("failed to filter commit report accepted: %w", err)
	}

	var countMerkleRoots int
	var countNoRoots int
	for iter.Next() {
		if len(iter.Event.BlessedMerkleRoots) == 0 && len(iter.Event.UnblessedMerkleRoots) == 0 {
			countNoRoots++
			continue
		}

		countMerkleRoots++
		for _, root := range iter.Event.BlessedMerkleRoots {
			if root.SourceChainSelector == srcChainSel {
				lggr.Debugw("checking commit root", "minSeqNr", root.MinSeqNr, "maxSeqNr", root.MaxSeqNr, "txHash", iter.Event.Raw.TxHash.String())
				if msgSeqNr >= root.MinSeqNr && msgSeqNr <= root.MaxSeqNr {
					lggr.Debugw("found commit root", "root", root, "txHash", iter.Event.Raw.TxHash.String())
					return root, nil
				}
			}
		}

		for _, root := range iter.Event.UnblessedMerkleRoots {
			if root.SourceChainSelector == srcChainSel {
				lggr.Debugw("checking commit root", "minSeqNr", root.MinSeqNr, "maxSeqNr", root.MaxSeqNr, "txHash", iter.Event.Raw.TxHash.String())
				if msgSeqNr >= root.MinSeqNr && msgSeqNr <= root.MaxSeqNr {
					lggr.Debugw("found commit root", "root", root, "txHash", iter.Event.Raw.TxHash.String())
					return root, nil
				}
			}
		}
	}

	lggr.Debugw("didn't find commit root", "countMerkleRoots", countMerkleRoots, "countNoRoots", countNoRoots)

	return offramp.InternalMerkleRoot{}, errors.New("commit root not found")
}

func getCCIPMessageSentEvents(
	ctx context.Context,
	lggr logger.Logger,
	env deployment.Environment,
	state ccipchangeset.CCIPOnChainState,
	srcChainSel uint64,
	destChainSel uint64,
	merkleRoot offramp.InternalMerkleRoot,
	lookbackDuration time.Duration,
) ([]onramp.OnRampCCIPMessageSent, error) {
	hdr, err := env.Chains[srcChainSel].Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get header: %w", err)
	}

	start := getStartBlock(srcChainSel, hdr.Number.Uint64(), lookbackDuration)

	var seqNrs []uint64
	for i := merkleRoot.MinSeqNr; i <= merkleRoot.MaxSeqNr; i++ {
		seqNrs = append(seqNrs, i)
	}

	lggr.Debugw("would query with",
		"seqNrs", seqNrs,
		"minSeqNr", merkleRoot.MinSeqNr,
		"maxSeqNr", merkleRoot.MaxSeqNr,
		"startBlock", start,
	)

	iter, err := state.Chains[srcChainSel].OnRamp.FilterCCIPMessageSent(
		&bind.FilterOpts{
			Start: start,
		},
		[]uint64{destChainSel},
		seqNrs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to filter ccip message sent: %w", err)
	}

	var ret []onramp.OnRampCCIPMessageSent
	var count int
	var toDestCount int
	for iter.Next() {
		count++
		if iter.Event.DestChainSelector == destChainSel {
			toDestCount++
			lggr.Debugw("checking message",
				"seqNr", iter.Event.SequenceNumber,
				"destChain", iter.Event.DestChainSelector,
				"txHash", iter.Event.Raw.TxHash.String())
			if iter.Event.SequenceNumber >= merkleRoot.MinSeqNr &&
				iter.Event.SequenceNumber <= merkleRoot.MaxSeqNr {
				ret = append(ret, *iter.Event)
			}
		}
	}

	lggr.Debugw("found messages", "count", count, "toDestCount", toDestCount)

	if len(ret) != len(seqNrs) {
		return nil, fmt.Errorf("not all messages found, got: %d, expected: %d", len(ret), len(seqNrs))
	}

	return ret, nil
}

func manuallyExecuteSingle(
	ctx context.Context,
	lggr logger.Logger,
	state ccipchangeset.CCIPOnChainState,
	env deployment.Environment,
	srcChainSel uint64,
	destChainSel uint64,
	msgSeqNr uint64,
	lookbackDuration time.Duration,
	reExecuteIfFailed bool,
	extraDataCodec ccipcommon.ExtraDataCodec,
) error {
	onRampAddress := state.Chains[srcChainSel].OnRamp.Address()

	execState, err := state.Chains[destChainSel].OffRamp.GetExecutionState(&bind.CallOpts{
		Context: ctx,
	}, srcChainSel, msgSeqNr)
	if err != nil {
		return fmt.Errorf("failed to get execution state: %w", err)
	}

	if execState == testhelpers.EXECUTION_STATE_SUCCESS ||
		(execState == testhelpers.EXECUTION_STATE_FAILURE && !reExecuteIfFailed) {
		lggr.Debugw("message already executed", "execState", execState)
		return nil
	}

	lggr.Debugw("addresses",
		"offRampAddress", state.Chains[destChainSel].OffRamp.Address(),
		"onRampAddress", onRampAddress,
		"execState", execState,
	)
	merkleRoot, err := getCommitRootAcceptedEvent(
		ctx,
		lggr,
		env,
		state,
		srcChainSel,
		destChainSel,
		msgSeqNr,
		lookbackDuration,
	)
	if err != nil {
		return fmt.Errorf("failed to get merkle root: %w", err)
	}

	lggr.Debugw("merkle root",
		"merkleRoot", hexutil.Encode(merkleRoot.MerkleRoot[:]),
		"minSeqNr", merkleRoot.MinSeqNr,
		"maxSeqNr", merkleRoot.MaxSeqNr,
		"sourceChainSel", merkleRoot.SourceChainSelector,
	)

	ccipMessageSentEvents, err := getCCIPMessageSentEvents(
		ctx,
		lggr,
		env,
		state,
		srcChainSel,
		destChainSel,
		merkleRoot,
		lookbackDuration,
	)
	if err != nil {
		return fmt.Errorf("failed to get ccip message sent event: %w", err)
	}

	messageHashes, err := manualexeclib.GetMessageHashes(
		ctx,
		lggr,
		onRampAddress,
		ccipMessageSentEvents,
		extraDataCodec,
	)
	if err != nil {
		return fmt.Errorf("failed to get message hashes: %w", err)
	}

	hashes, flags, err := manualexeclib.GetMerkleProof(
		lggr,
		merkleRoot,
		messageHashes,
		msgSeqNr,
	)
	if err != nil {
		return fmt.Errorf("failed to get merkle proof: %w", err)
	}

	lggr.Debugw("got hashes and flags", "hashes", hashes, "flags", flags)

	// since we're only executing one message, we need to only include that message
	// in the report.
	var filteredMsgSentEvents []onramp.OnRampCCIPMessageSent
	for i, event := range ccipMessageSentEvents {
		if event.Message.Header.SequenceNumber == msgSeqNr && event.Message.Header.SourceChainSelector == srcChainSel {
			filteredMsgSentEvents = append(filteredMsgSentEvents, ccipMessageSentEvents[i])
		}
	}

	// sanity check, should not be possible at this point.
	if len(filteredMsgSentEvents) == 0 {
		return fmt.Errorf("no message found for seqNr %d", msgSeqNr)
	}

	execReport, err := manualexeclib.CreateExecutionReport(
		srcChainSel,
		onRampAddress,
		filteredMsgSentEvents,
		hashes,
		flags,
		extraDataCodec,
	)
	if err != nil {
		return fmt.Errorf("failed to create execution report: %w", err)
	}

	txOpts := &bind.TransactOpts{
		From:   env.Chains[destChainSel].DeployerKey.From,
		Nonce:  nil,
		Signer: env.Chains[destChainSel].DeployerKey.Signer,
		Value:  big.NewInt(0),
		// We manually set the gas limit here because estimateGas doesn't take into account
		// internal reverts (such as those that could happen on ERC165 interface checks).
		// This is just a big value for now, we can investigate something more efficient later.
		GasLimit: 1e6,
	}
	tx, err := state.Chains[destChainSel].OffRamp.ManuallyExecute(
		txOpts,
		[]offramp.InternalExecutionReport{execReport},
		[][]offramp.OffRampGasLimitOverride{
			{
				{
					ReceiverExecutionGasLimit: big.NewInt(200_000),
					TokenGasOverrides:         nil,
				},
			},
		},
	)
	_, err = deployment.ConfirmIfNoErrorWithABI(env.Chains[destChainSel], tx, offramp.OffRampABI, err)
	if err != nil {
		return fmt.Errorf("failed to execute message: %w", err)
	}

	lggr.Debugw("successfully manually executed msg", "msgSeqNr", msgSeqNr)

	return nil
}

// ManuallyExecuteAll will manually execute the provided messages if they were not already executed.
// At the moment offchain token data (i.e USDC/Lombard attestations) is not supported.
func ManuallyExecuteAll(
	ctx context.Context,
	lggr logger.Logger,
	state ccipchangeset.CCIPOnChainState,
	env deployment.Environment,
	srcChainSel uint64,
	destChainSel uint64,
	msgSeqNrs []int64,
	lookbackDuration time.Duration,
	reExecuteIfFailed bool,
) error {
	extraDataCodec := ccipcommon.NewExtraDataCodec(ccipcommon.NewExtraDataCodecParams(ccipevm.ExtraDataDecoder{}, ccipsolana.ExtraDataDecoder{}))
	for _, seqNr := range msgSeqNrs {
		err := manuallyExecuteSingle(
			ctx,
			lggr,
			state,
			env,
			srcChainSel,
			destChainSel,
			uint64(seqNr), //nolint:gosec // seqNr is never <= 0.
			lookbackDuration,
			reExecuteIfFailed,
			extraDataCodec,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckAlreadyExecuted will check the execution state of the provided messages and log if they were already executed.
func CheckAlreadyExecuted(
	ctx context.Context,
	lggr logger.Logger,
	state ccipchangeset.CCIPOnChainState,
	srcChainSel uint64,
	destChainSel uint64,
	msgSeqNrs []int64,
) error {
	for _, seqNr := range msgSeqNrs {
		execState, err := state.Chains[destChainSel].OffRamp.GetExecutionState(
			&bind.CallOpts{Context: ctx},
			srcChainSel,
			uint64(seqNr), //nolint:gosec // seqNr is never <= 0.
		)
		if err != nil {
			return fmt.Errorf("failed to get execution state: %w", err)
		}

		if execState == testhelpers.EXECUTION_STATE_SUCCESS || execState == testhelpers.EXECUTION_STATE_FAILURE {
			lggr.Debugw("message already executed", "execState", execState, "msgSeqNr", seqNr)
		} else {
			lggr.Debugw("message not executed", "execState", execState, "msgSeqNr", seqNr)
		}
	}

	return nil
}
