package commit

import (
	"context"
	"fmt"
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/ccipocr3/internal/codec"
	"github.com/smartcontractkit/ccipocr3/internal/libs/hashlib"
	"github.com/smartcontractkit/ccipocr3/internal/libs/merklemulti"
	"github.com/smartcontractkit/ccipocr3/internal/libs/slicelib"
	"github.com/smartcontractkit/ccipocr3/internal/model"
	"github.com/smartcontractkit/ccipocr3/internal/reader"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/sync/errgroup"

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
	msgHasher codec.MessageHasher,
	readableChains mapset.Set[model.ChainSelector],
	maxSeqNumsPerChain []model.SeqNumChain,
	msgScanBatchSize int,
) ([]model.CCIPMsg, error) {
	// Find the new msgs for each supported chain based on the discovered max sequence numbers.
	newMsgsPerChain := make([][]model.CCIPMsg, len(maxSeqNumsPerChain))
	eg := new(errgroup.Group)

	for chainIdx, seqNumChain := range maxSeqNumsPerChain {
		if !readableChains.Contains(seqNumChain.ChainSel) {
			lggr.Debugw("reading chain is not supported", "chain", seqNumChain.ChainSel)
			continue
		}

		seqNumChain := seqNumChain
		chainIdx := chainIdx
		eg.Go(func() error {
			minSeqNum := seqNumChain.SeqNum + 1
			maxSeqNum := minSeqNum + model.SeqNum(msgScanBatchSize)
			lggr.Debugw("scanning for new messages",
				"chain", seqNumChain.ChainSel, "minSeqNum", minSeqNum, "maxSeqNum", maxSeqNum)

			newMsgs, err := ccipReader.MsgsBetweenSeqNums(
				ctx, seqNumChain.ChainSel, model.NewSeqNumRange(minSeqNum, maxSeqNum))
			if err != nil {
				return fmt.Errorf("get messages between seq nums: %w", err)
			}

			if len(newMsgs) > 0 {
				lggr.Debugw("discovered new messages", "chain", seqNumChain.ChainSel, "newMsgs", len(newMsgs))
			} else {
				lggr.Debugw("no new messages discovered", "chain", seqNumChain.ChainSel)
			}

			for _, msg := range newMsgs {
				msgHash, err := msgHasher.Hash(msg)
				if err != nil {
					return fmt.Errorf("hash message: %w", err)
				}

				if msgHash != msg.ID {
					lggr.Warnw("invalid message discovered", "msg", msg, "err", err)
					continue
				}
			}

			newMsgsPerChain[chainIdx] = newMsgs
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("wait for new msg observations: %w", err)
	}

	observedNewMsgs := make([]model.CCIPMsg, 0)
	for chainIdx := range maxSeqNumsPerChain {
		observedNewMsgs = append(observedNewMsgs, newMsgsPerChain[chainIdx]...)
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

func observeGasPrices(ctx context.Context, ccipReader reader.CCIP, chains []model.ChainSelector) ([]model.GasPriceChain, error) {
	if len(chains) == 0 {
		return nil, nil
	}

	gasPrices, err := ccipReader.GasPrices(ctx, chains)
	if err != nil {
		return nil, fmt.Errorf("get gas prices: %w", err)
	}

	if len(gasPrices) != len(chains) {
		return nil, fmt.Errorf("internal critical error gas prices length mismatch: got %d, want %d",
			len(gasPrices), len(chains))
	}

	gasPricesGwei := make([]model.GasPriceChain, 0, len(chains))
	for i, chain := range chains {
		gasPricesGwei = append(gasPricesGwei, model.NewGasPriceChain(gasPrices[i].Int, chain))
	}

	return gasPricesGwei, nil
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

	// First come to consensus about the (sequence number, id) pairs.
	// For each sequence number consider correct the ID with the most votes.
	msgSeqNumToIDCounts := make(map[model.SeqNum]map[string]int) // seqNum -> msgID -> count
	for _, msg := range observedMsgs {
		if _, exists := msgSeqNumToIDCounts[msg.SeqNum]; !exists {
			msgSeqNumToIDCounts[msg.SeqNum] = make(map[string]int)
		}
		msgSeqNumToIDCounts[msg.SeqNum][msg.ID.String()]++
	}
	lggr.Debugw("observed message counts", "chain", chainSel, "msgSeqNumToIdCounts", msgSeqNumToIDCounts)

	msgObservationsCount := make(map[model.SeqNum]int)
	msgSeqNumToID := make(map[model.SeqNum]model.Bytes32)
	for seqNum, idCounts := range msgSeqNumToIDCounts {
		if len(idCounts) == 0 {
			lggr.Errorw("critical error id counts should never be empty", "seqNum", seqNum)
			continue
		}

		// Find the ID with the most votes for each sequence number.
		idsSlice := make([]string, 0, len(idCounts))
		for id := range idCounts {
			idsSlice = append(idsSlice, id)
		}
		// determinism in case we have the same count for different ids
		sort.Slice(idsSlice, func(i, j int) bool { return idsSlice[i] < idsSlice[j] })

		maxCnt := idCounts[idsSlice[0]]
		mostVotedID := idsSlice[0]
		for _, id := range idsSlice[1:] {
			cnt := idCounts[id]
			if cnt > maxCnt {
				maxCnt = cnt
				mostVotedID = id
			}
		}

		msgObservationsCount[seqNum] = maxCnt
		idBytes, err := model.NewBytes32FromString(mostVotedID)
		if err != nil {
			return observedMsgsConsensus{}, fmt.Errorf("critical issue converting id '%s' to bytes32: %w",
				mostVotedID, err)
		}
		msgSeqNumToID[seqNum] = idBytes
	}
	lggr.Debugw("observed message consensus", "chain", chainSel, "msgSeqNumToId", msgSeqNumToID)

	// Filter out msgs not observed by at least 2f_chain+1 followers.
	msgSeqNumsQuorum := mapset.NewSet[model.SeqNum]()
	for seqNum, count := range msgObservationsCount {
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
		consensusMsgID, ok := msgSeqNumToID[msg.SeqNum]
		if !ok || consensusMsgID != msg.ID {
			continue
		}
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
func maxSeqNumsConsensus(lggr logger.Logger, fChain int, observations []model.CommitPluginObservation) []model.SeqNumChain {
	observedSeqNumsPerChain := make(map[model.ChainSelector][]model.SeqNum)
	for _, obs := range observations {
		for _, maxSeqNum := range obs.MaxSeqNums {
			if _, exists := observedSeqNumsPerChain[maxSeqNum.ChainSel]; !exists {
				observedSeqNumsPerChain[maxSeqNum.ChainSel] = make([]model.SeqNum, 0)
			}
			observedSeqNumsPerChain[maxSeqNum.ChainSel] = append(observedSeqNumsPerChain[maxSeqNum.ChainSel], maxSeqNum.SeqNum)
		}
	}

	seqNums := make([]model.SeqNumChain, 0, len(observedSeqNumsPerChain))
	for ch, observedSeqNums := range observedSeqNumsPerChain {
		if len(observedSeqNums) < 2*fChain+1 {
			lggr.Warnw("not enough observations for chain", "chain", ch, "observedSeqNums", observedSeqNums)
			continue
		}

		sort.Slice(observedSeqNums, func(i, j int) bool { return observedSeqNums[i] < observedSeqNums[j] })
		seqNums = append(seqNums, model.NewSeqNumChain(ch, observedSeqNums[fChain]))
	}

	sort.Slice(seqNums, func(i, j int) bool { return seqNums[i].ChainSel < seqNums[j].ChainSel })
	return seqNums
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

func gasPricesConsensus(lggr logger.Logger, observations []model.CommitPluginObservation, fChain int) []model.GasPriceChain {
	// Group the observed gas prices by chain.
	gasPricePerChain := make(map[model.ChainSelector][]model.BigInt)
	for _, obs := range observations {
		for _, gasPrice := range obs.GasPrices {
			if _, exists := gasPricePerChain[gasPrice.ChainSel]; !exists {
				gasPricePerChain[gasPrice.ChainSel] = make([]model.BigInt, 0)
			}
			gasPricePerChain[gasPrice.ChainSel] = append(gasPricePerChain[gasPrice.ChainSel], gasPrice.GasPrice)
		}
	}

	// Keep the median
	consensusGasPrices := make([]model.GasPriceChain, 0)
	for chain, gasPrices := range gasPricePerChain {
		if len(gasPrices) < 2*fChain+1 {
			lggr.Warnw("not enough gas price observations", "chain", chain, "gasPrices", gasPrices)
			continue
		}

		consensusGasPrices = append(
			consensusGasPrices,
			model.NewGasPriceChain(slicelib.BigIntSortedMiddle(gasPrices).Int, chain),
		)
	}

	sort.Slice(consensusGasPrices, func(i, j int) bool { return consensusGasPrices[i].ChainSel < consensusGasPrices[j].ChainSel })
	return consensusGasPrices
}

// pluginConfigConsensus comes to consensus on the plugin config based on the observations.
// We cannot trust the state of a single follower, so we need to come to consensus on the config.
func pluginConfigConsensus(
	baseCfg model.CommitPluginConfig, // the config of the follower calling this function
	observations []model.CommitPluginObservation, // observations from all followers
) model.CommitPluginConfig {
	consensusCfg := baseCfg

	// Come to consensus on fChain.
	// Use the fChain observed by most followers for each chain.
	fChainCounts := make(map[model.ChainSelector]map[int]int) // {chain: {fChain: count}}
	for _, obs := range observations {
		for chain, fChain := range obs.PluginConfig.FChain {
			if _, exists := fChainCounts[chain]; !exists {
				fChainCounts[chain] = make(map[int]int)
			}
			fChainCounts[chain][fChain]++
		}
	}
	consensusFChain := make(map[model.ChainSelector]int)
	for chain, counts := range fChainCounts {
		maxCount := 0
		for fChain, count := range counts {
			if count > maxCount {
				maxCount = count
				consensusFChain[chain] = fChain
			}
		}
	}
	consensusCfg.FChain = consensusFChain

	// Come to consensus on what the feeTokens are.
	// We want to keep the tokens observed by at least 2f_chain+1 followers.
	feeTokensCounts := make(map[types.Account]int)
	for _, obs := range observations {
		for _, token := range obs.PluginConfig.PricedTokens {
			feeTokensCounts[token]++
		}
	}
	consensusFeeTokens := make([]types.Account, 0)
	for token, count := range feeTokensCounts {
		if count >= 2*consensusCfg.FChain[consensusCfg.DestChain]+1 {
			consensusFeeTokens = append(consensusFeeTokens, token)
		}
	}
	consensusCfg.PricedTokens = consensusFeeTokens

	// Come to consensus on reading observers.
	// An observer can read a chain only if at least 2f_chain+1 followers observed that.
	observerReadChainsCounts := make(map[commontypes.OracleID]map[model.ChainSelector]int)
	for _, obs := range observations {
		for observer, info := range obs.PluginConfig.ObserverInfo {
			if _, exists := observerReadChainsCounts[observer]; !exists {
				observerReadChainsCounts[observer] = make(map[model.ChainSelector]int)
			}
			for _, chain := range info.Reads {
				observerReadChainsCounts[observer][chain]++
			}
		}
	}
	consensusObserverInfo := make(map[commontypes.OracleID]model.ObserverInfo)
	for observer, chainCounts := range observerReadChainsCounts {
		observerReadChains := make([]model.ChainSelector, 0)
		for chain, count := range chainCounts {
			if count >= 2*consensusCfg.FChain[consensusCfg.DestChain]+1 {
				observerReadChains = append(observerReadChains, chain)
			}
		}
		observerInfo := consensusCfg.ObserverInfo[observer]
		observerInfo.Reads = observerReadChains
		consensusObserverInfo[observer] = observerInfo
	}

	return consensusCfg
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

func validateObservedGasPrices(gasPrices []model.GasPriceChain) error {
	// Duplicate gas prices must not appear for the same chain and must not be empty.
	gasPriceChains := mapset.NewSet[model.ChainSelector]()
	for _, g := range gasPrices {
		if gasPriceChains.Contains(g.ChainSel) {
			return fmt.Errorf("duplicate gas price for chain %d", g.ChainSel)
		}
		gasPriceChains.Add(g.ChainSel)
		if g.GasPrice.IsEmpty() {
			return fmt.Errorf("gas price must not be empty")
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
