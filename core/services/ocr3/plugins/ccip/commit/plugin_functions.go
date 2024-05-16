package commit

import (
	"context"
	"fmt"
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/ccipocr3/internal/libs/hashlib"
	"github.com/smartcontractkit/ccipocr3/internal/libs/merklemulti"
	"github.com/smartcontractkit/ccipocr3/internal/libs/slicelib"
	"github.com/smartcontractkit/ccipocr3/internal/model"
	"github.com/smartcontractkit/ccipocr3/internal/reader"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// observeMaxSeqNums finds the maximum committed sequence numbers for each source chain.
// If a sequence number is pending (is not on-chain yet), it will be included in the results.
func observeMaxSeqNums(
	ctx context.Context,
	lggr logger.Logger,
	ccipReader reader.CCIP,
	previousOutcomeBytes []byte,
	readableChains mapset.Set[model.ChainSelector],
	destChain model.ChainSelector,
	knownSourceChains []model.ChainSelector,
) ([]model.SeqNumChain, error) {
	// If there is a previous outcome, start with the sequence numbers of it.
	seqNumPerChain := make(map[model.ChainSelector]model.SeqNum)
	if previousOutcomeBytes != nil {
		lggr.Debugw("observing based on previous outcome")
		prevOutcome, err := model.DecodeCommitPluginOutcome(previousOutcomeBytes)
		if err != nil {
			return nil, fmt.Errorf("decode commit plugin previous outcome: %w", err)
		}
		lggr.Debugw("previous outcome decoded", "outcome", prevOutcome.String())

		for _, seqNumChain := range prevOutcome.MaxSeqNums {
			if seqNumChain.SeqNum > seqNumPerChain[seqNumChain.ChainSel] {
				seqNumPerChain[seqNumChain.ChainSel] = seqNumChain.SeqNum
			}
		}
		lggr.Debugw("discovered sequence numbers from prev outcome", "seqNumPerChain", seqNumPerChain)
	}

	// If reading destination chain is supported find the latest sequence numbers per chain from the onchain state.
	if readableChains.Contains(destChain) {
		lggr.Debugw("reading sequence numbers from destination")
		onChainSeqNums, err := ccipReader.NextSeqNum(ctx, knownSourceChains)
		if err != nil {
			return nil, fmt.Errorf("get next seq nums: %w", err)
		}
		lggr.Debugw("discovered sequence numbers from destination", "onChainSeqNums", onChainSeqNums)

		// Update the seq nums if the on-chain sequence number is greater than previous outcome.
		for i, ch := range knownSourceChains {
			if onChainSeqNums[i] > seqNumPerChain[ch] {
				seqNumPerChain[ch] = onChainSeqNums[i]
				lggr.Debugw("updated sequence number", "chain", ch, "seqNum", onChainSeqNums[i])
			}
		}
	}

	maxChainSeqNums := make([]model.SeqNumChain, 0)
	for ch, seqNum := range seqNumPerChain {
		maxChainSeqNums = append(maxChainSeqNums, model.NewSeqNumChain(ch, seqNum))
	}

	sort.Slice(maxChainSeqNums, func(i, j int) bool { return maxChainSeqNums[i].ChainSel < maxChainSeqNums[j].ChainSel })
	return maxChainSeqNums, nil
}

// observeNewMsgs finds the new messages for each supported chain based on the provided max sequence numbers.
func observeNewMsgs(
	ctx context.Context,
	lggr logger.Logger,
	ccipReader reader.CCIP,
	readableChains mapset.Set[model.ChainSelector],
	maxSeqNumsPerChain []model.SeqNumChain,
	msgScanBatchSize int,
) ([]model.CCIPMsgBaseDetails, error) {
	// Find the new msgs for each supported chain based on the discovered max sequence numbers.
	observedNewMsgs := make([]model.CCIPMsgBaseDetails, 0)
	for _, seqNumChain := range maxSeqNumsPerChain {
		if !readableChains.Contains(seqNumChain.ChainSel) {
			lggr.Debugw("reading chain is not supported", "chain", seqNumChain.ChainSel)
			continue
		}

		minSeqNum := seqNumChain.SeqNum + 1
		maxSeqNum := minSeqNum + model.SeqNum(msgScanBatchSize)
		lggr.Debugw("scanning for new messages",
			"chain", seqNumChain.ChainSel, "minSeqNum", minSeqNum, "maxSeqNum", maxSeqNum)

		newMsgs, err := ccipReader.MsgsBetweenSeqNums(
			ctx, []model.ChainSelector{seqNumChain.ChainSel}, model.NewSeqNumRange(minSeqNum, maxSeqNum))
		if err != nil {
			return nil, fmt.Errorf("get messages between seq nums: %w", err)
		}

		if len(newMsgs) > 0 {
			lggr.Debugw("discovered new messages", "chain", seqNumChain.ChainSel, "newMsgs", len(newMsgs))
		} else {
			lggr.Debugw("no new messages discovered", "chain", seqNumChain.ChainSel)
		}

		for _, msg := range newMsgs {
			observedNewMsgs = append(observedNewMsgs, msg.CCIPMsgBaseDetails)
		}
	}

	return observedNewMsgs, nil
}

func observeTokenPrices(
	ctx context.Context,
	tokenPricesReader reader.TokenPrices,
	tokens []types.Account,
) ([]model.TokenPrice, error) {
	tokenPrices, err := tokenPricesReader.GetTokenPricesUSD(ctx, tokens)
	if err != nil {
		return nil, fmt.Errorf("get token prices: %w", err)
	}

	if len(tokenPrices) != len(tokens) {
		return nil, fmt.Errorf("internal critical error token prices length mismatch: got %d, want %d",
			len(tokenPrices), len(tokens))
	}

	tokenPricesUSD := make([]model.TokenPrice, 0, len(tokens))
	for i, token := range tokens {
		tokenPricesUSD = append(tokenPricesUSD, model.NewTokenPrice(token, tokenPrices[i]))
	}

	return tokenPricesUSD, nil
}

// newMsgsConsensus comes in consensus on the observed messages for each source chain. Generates one merkle root
// for each source chain based on the consensus on the messages.
func newMsgsConsensus(
	lggr logger.Logger,
	maxSeqNums []model.SeqNumChain,
	observations []model.CommitPluginObservation,
	fChainCfg map[model.ChainSelector]int,
) ([]model.MerkleRootChain, error) {
	maxSeqNumsPerChain := make(map[model.ChainSelector]model.SeqNum)
	for _, seqNumChain := range maxSeqNums {
		maxSeqNumsPerChain[seqNumChain.ChainSel] = seqNumChain.SeqNum
	}

	// Gather all messages from all observations.
	msgsFromObservations := make([]model.CCIPMsgBaseDetails, 0)
	for _, obs := range observations {
		msgsFromObservations = append(msgsFromObservations, obs.NewMsgs...)
	}
	lggr.Debugw("total observed messages across all followers", "msgs", len(msgsFromObservations))

	// Filter out messages less than or equal to the max sequence numbers.
	msgsFromObservations = slicelib.Filter(msgsFromObservations, func(msg model.CCIPMsgBaseDetails) bool {
		maxSeqNum, ok := maxSeqNumsPerChain[msg.SourceChain]
		if !ok {
			return false
		}
		return msg.SeqNum > maxSeqNum
	})
	lggr.Debugw("observed messages after filtering", "msgs", len(msgsFromObservations))

	// Group messages by source chain.
	sourceChains, groupedMsgs := slicelib.GroupBy(
		msgsFromObservations,
		func(msg model.CCIPMsgBaseDetails) model.ChainSelector { return msg.SourceChain },
	)

	// Come to consensus on the observed messages by source chain.
	consensusBySourceChain := make(map[model.ChainSelector]observedMsgsConsensus)
	for _, sourceChain := range sourceChains { // note: we iterate using sourceChains slice for deterministic order.
		observedMsgs, ok := groupedMsgs[sourceChain]
		if !ok {
			lggr.Panicw("source chain not found in grouped messages", "sourceChain", sourceChain)
		}

		msgsConsensus, err := newMsgsConsensusForChain(lggr, sourceChain, observedMsgs, fChainCfg)
		if err != nil {
			return nil, fmt.Errorf("calculate observed msgs consensus: %w", err)
		}

		if msgsConsensus.isEmpty() {
			lggr.Debugw("no consensus on observed messages", "sourceChain", sourceChain)
			continue
		}
		consensusBySourceChain[sourceChain] = msgsConsensus
		lggr.Debugw("observed messages consensus", "sourceChain", sourceChain, "consensus", msgsConsensus)
	}

	merkleRoots := make([]model.MerkleRootChain, 0)
	for sourceChain, consensus := range consensusBySourceChain {
		merkleRoots = append(
			merkleRoots,
			model.NewMerkleRootChain(sourceChain, consensus.seqNumRange, consensus.merkleRoot),
		)
	}

	sort.Slice(merkleRoots, func(i, j int) bool { return merkleRoots[i].ChainSel < merkleRoots[j].ChainSel })
	return merkleRoots, nil
}

// Given a list of observed msgs
//   - Keep the messages that were observed by at least 2f_chain+1 followers.
//   - Starting from the first message (min seq num), keep adding the messages to the merkle tree until a gap is found.
func newMsgsConsensusForChain(
	lggr logger.Logger,
	chainSel model.ChainSelector,
	observedMsgs []model.CCIPMsgBaseDetails,
	fChainCfg map[model.ChainSelector]int,
) (observedMsgsConsensus, error) {
	fChain, ok := fChainCfg[chainSel]
	if !ok {
		return observedMsgsConsensus{}, fmt.Errorf("fchain not found for chain %d", chainSel)
	}
	lggr.Debugw("observed messages consensus",
		"chain", chainSel, "fChain", fChain, "observedMsgs", len(observedMsgs))

	// Reach consensus on the observed msgs sequence numbers.
	msgSeqNums := make(map[model.SeqNum]int)
	for _, msg := range observedMsgs {
		msgSeqNums[msg.SeqNum]++
		// TODO: message data might be spoofed, validate the message data
	}
	lggr.Debugw("observed message counts", "chain", chainSel, "msgSeqNums", msgSeqNums)

	// Filter out msgs not observed by at least 2f_chain+1 followers.
	msgSeqNumsQuorum := mapset.NewSet[model.SeqNum]()
	for seqNum, count := range msgSeqNums {
		if count >= 2*fChain+1 {
			msgSeqNumsQuorum.Add(seqNum)
		}
	}
	if msgSeqNumsQuorum.Cardinality() == 0 {
		return observedMsgsConsensus{}, nil
	}

	// Come to consensus on the observed messages sequence numbers range.
	msgSeqNumsQuorumSlice := msgSeqNumsQuorum.ToSlice()
	sort.Slice(msgSeqNumsQuorumSlice, func(i, j int) bool { return msgSeqNumsQuorumSlice[i] < msgSeqNumsQuorumSlice[j] })
	seqNumConsensusRange := model.NewSeqNumRange(msgSeqNumsQuorumSlice[0], msgSeqNumsQuorumSlice[0])
	for _, seqNum := range msgSeqNumsQuorumSlice[1:] {
		if seqNum != seqNumConsensusRange.End()+1 {
			break // Found a gap in the sequence numbers.
		}
		seqNumConsensusRange.SetEnd(seqNum)
	}

	msgsBySeqNum := make(map[model.SeqNum]model.CCIPMsgBaseDetails)
	for _, msg := range observedMsgs {
		msgsBySeqNum[msg.SeqNum] = msg
	}

	treeLeaves := make([][32]byte, 0)
	for seqNum := seqNumConsensusRange.Start(); seqNum <= seqNumConsensusRange.End(); seqNum++ {
		msg, ok := msgsBySeqNum[seqNum]
		if !ok {
			return observedMsgsConsensus{}, fmt.Errorf("msg not found in map for seq num %d", seqNum)
		}
		treeLeaves = append(treeLeaves, msg.ID)
	}

	lggr.Debugw("constructing merkle tree", "chain", chainSel, "treeLeaves", len(treeLeaves))
	tree, err := merklemulti.NewTree(hashlib.NewKeccakCtx(), treeLeaves)
	if err != nil {
		return observedMsgsConsensus{}, fmt.Errorf("construct merkle tree from %d leaves: %w", len(treeLeaves), err)
	}

	return observedMsgsConsensus{
		seqNumRange: seqNumConsensusRange,
		merkleRoot:  tree.Root(),
	}, nil
}

// maxSeqNumsConsensus groups the observed max seq nums across all followers per chain.
// Orders the sequence numbers and selects the one at the index of destination chain fChain.
//
// For example:
//
//	seqNums: [1, 1, 1, 10, 10, 10, 10, 10, 10]
//	fChain: 4
//	result: 10
//
// Selecting seqNums[fChain] ensures:
//   - At least one honest node has seen this value, so adversary cannot bias the value lower which would cause reverts
//   - If an honest oracle reports sorted_min[f] which happens to be stale i.e. that oracle has a delayed view
//     of the chain, then the report will revert onchain but still succeed upon retry
//   - We minimize the risk of naturally hitting the error condition minSeqNum > maxSeqNum due to oracles
//     delayed views of the chain (would be an issue with taking sorted_mins[-f])
func maxSeqNumsConsensus(lggr logger.Logger, fChain int, observations []model.CommitPluginObservation,
) ([]model.SeqNumChain, error) {
	observedSeqNumsPerChain := make(map[model.ChainSelector][]model.SeqNum)
	for _, obs := range observations {
		for _, maxSeqNum := range obs.MaxSeqNums {
			if _, exists := observedSeqNumsPerChain[maxSeqNum.ChainSel]; !exists {
				observedSeqNumsPerChain[maxSeqNum.ChainSel] = make([]model.SeqNum, 0)
			}
			observedSeqNumsPerChain[maxSeqNum.ChainSel] = append(observedSeqNumsPerChain[maxSeqNum.ChainSel], maxSeqNum.SeqNum)
		}
	}

	maxSeqNumsConsensus := make([]model.SeqNumChain, 0, len(observedSeqNumsPerChain))
	for ch, observedSeqNums := range observedSeqNumsPerChain {
		if len(observedSeqNums) < 2*fChain+1 {
			lggr.Warnw("not enough observations for chain", "chain", ch, "observedSeqNums", observedSeqNums)
			continue
		}

		sort.Slice(observedSeqNums, func(i, j int) bool { return observedSeqNums[i] < observedSeqNums[j] })
		maxSeqNumsConsensus = append(maxSeqNumsConsensus, model.NewSeqNumChain(ch, observedSeqNums[fChain]))
	}

	sort.Slice(maxSeqNumsConsensus, func(i, j int) bool { return maxSeqNumsConsensus[i].ChainSel < maxSeqNumsConsensus[j].ChainSel })
	return maxSeqNumsConsensus, nil
}

// tokenPricesConsensus returns the median price for tokens that have at least 2f_chain+1 observations.
func tokenPricesConsensus(
	observations []model.CommitPluginObservation,
	fChain int,
) ([]model.TokenPrice, error) {
	pricesPerToken := make(map[types.Account][]model.BigInt)
	for _, obs := range observations {
		for _, price := range obs.TokenPrices {
			if _, exists := pricesPerToken[price.TokenID]; !exists {
				pricesPerToken[price.TokenID] = make([]model.BigInt, 0)
			}
			pricesPerToken[price.TokenID] = append(pricesPerToken[price.TokenID], price.Price)
		}
	}

	// Keep the median
	consensusPrices := make([]model.TokenPrice, 0)
	for token, prices := range pricesPerToken {
		if len(prices) < 2*fChain+1 {
			continue
		}
		consensusPrices = append(consensusPrices, model.NewTokenPrice(token, slicelib.BigIntSortedMiddle(prices).Int))
	}

	sort.Slice(consensusPrices, func(i, j int) bool { return consensusPrices[i].TokenID < consensusPrices[j].TokenID })
	return consensusPrices, nil
}

// validateObservedSequenceNumbers checks if the sequence numbers of the provided messages are unique for each chain and
// that they match the observed max sequence numbers.
func validateObservedSequenceNumbers(msgs []model.CCIPMsgBaseDetails, maxSeqNums []model.SeqNumChain) error {
	// MaxSeqNums must be unique for each chain.
	maxSeqNumsMap := make(map[model.ChainSelector]model.SeqNum)
	for _, maxSeqNum := range maxSeqNums {
		if _, exists := maxSeqNumsMap[maxSeqNum.ChainSel]; exists {
			return fmt.Errorf("duplicate max sequence number for chain %d", maxSeqNum.ChainSel)
		}
		maxSeqNumsMap[maxSeqNum.ChainSel] = maxSeqNum.SeqNum
	}

	seqNums := make(map[model.ChainSelector]mapset.Set[model.SeqNum], len(msgs))
	for _, msg := range msgs {
		// The same sequence number must not appear more than once for the same chain and must be valid.

		if _, exists := seqNums[msg.SourceChain]; !exists {
			seqNums[msg.SourceChain] = mapset.NewSet[model.SeqNum]()
		}

		if seqNums[msg.SourceChain].Contains(msg.SeqNum) {
			return fmt.Errorf("duplicate sequence number %d for chain %d", msg.SeqNum, msg.SourceChain)
		}
		seqNums[msg.SourceChain].Add(msg.SeqNum)

		// The observed msg sequence number cannot be less than or equal to the max observed sequence number.
		maxSeqNum, exists := maxSeqNumsMap[msg.SourceChain]
		if !exists {
			return fmt.Errorf("max sequence number observation not found for chain %d", msg.SourceChain)
		}
		if msg.SeqNum <= maxSeqNum {
			return fmt.Errorf("max sequence number %d must be greater than observed sequence number %d for chain %d",
				maxSeqNum, msg.SeqNum, msg.SourceChain)
		}
	}

	return nil
}

// validateObserverReadingEligibility checks if the observer is eligible to observe the messages it observed.
func validateObserverReadingEligibility(
	observer commontypes.OracleID,
	msgs []model.CCIPMsgBaseDetails,
	observerCfg map[commontypes.OracleID]model.ObserverInfo,
) error {
	if len(msgs) == 0 {
		return nil
	}

	observerInfo, exists := observerCfg[observer]
	if !exists {
		return fmt.Errorf("observer not found in config")
	}

	observerReadChains := mapset.NewSet(observerInfo.Reads...)

	for _, msg := range msgs {
		// Observer must be able to read the chain that the message is coming from.
		if !observerReadChains.Contains(msg.SourceChain) {
			return fmt.Errorf("observer not allowed to read chain %d", msg.SourceChain)
		}
	}

	return nil
}

func validateObservedTokenPrices(tokenPrices []model.TokenPrice) error {
	tokensWithPrice := mapset.NewSet[types.Account]()
	for _, t := range tokenPrices {
		if tokensWithPrice.Contains(t.TokenID) {
			return fmt.Errorf("duplicate token price for token: %s", t.TokenID)
		}
		tokensWithPrice.Add(t.TokenID)

		if t.Price.IsEmpty() {
			return fmt.Errorf("token price must not be empty")
		}
	}

	return nil
}

func validateObservedGasPrices(gasPrices []model.GasPriceChain, tokenPrices []model.TokenPrice) error {
	// Duplicate gas prices must not appear for the same chain and must not be empty.
	gasPriceChains := mapset.NewSet[model.ChainSelector]()
	for _, g := range gasPrices {
		if gasPriceChains.Contains(g.ChainSel) {
			return fmt.Errorf("duplicate gas price for chain %d", g.ChainSel)
		}
		gasPriceChains.Add(g.ChainSel)
		if g.GasPrice == nil {
			return fmt.Errorf("gas price must not be nil")
		}
	}

	return nil
}

type observedMsgsConsensus struct {
	seqNumRange model.SeqNumRange
	merkleRoot  [32]byte
}

func (o observedMsgsConsensus) isEmpty() bool {
	return o.seqNumRange.Start() == 0 && o.seqNumRange.End() == 0 && o.merkleRoot == [32]byte{}
}
