package commit

import (
	"context"
	"fmt"
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/ccipocr3/internal/libs/slicelib"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"
)

// observeLatestCommittedSeqNums finds the maximum committed sequence numbers for each source chain.
// If we cannot observe the dest we return an empty slice and no error.
func observeLatestCommittedSeqNums(
	ctx context.Context,
	lggr logger.Logger,
	ccipReader cciptypes.CCIPReader,
	readableChains mapset.Set[cciptypes.ChainSelector],
	destChain cciptypes.ChainSelector,
	knownSourceChains []cciptypes.ChainSelector,
) ([]cciptypes.SeqNumChain, error) {
	sort.Slice(knownSourceChains, func(i, j int) bool { return knownSourceChains[i] < knownSourceChains[j] })
	latestCommittedSeqNumsObservation := make([]cciptypes.SeqNumChain, 0)
	if readableChains.Contains(destChain) {
		lggr.Debugw("reading latest committed sequence from destination")
		onChainLatestCommittedSeqNums, err := ccipReader.NextSeqNum(ctx, knownSourceChains)
		if err != nil {
			return latestCommittedSeqNumsObservation, fmt.Errorf("get next seq nums: %w", err)
		}
		lggr.Debugw("observed latest committed sequence numbers on destination", "latestCommittedSeqNumsObservation", onChainLatestCommittedSeqNums)
		for i, ch := range knownSourceChains {
			latestCommittedSeqNumsObservation = append(latestCommittedSeqNumsObservation, cciptypes.NewSeqNumChain(ch, onChainLatestCommittedSeqNums[i]))
		}
	}
	return latestCommittedSeqNumsObservation, nil
}

// observeNewMsgs finds the new messages for each supported chain based on the provided max sequence numbers.
// If latestCommitSeqNums is empty (first ever OCR round), it will return an empty slice.
func observeNewMsgs(
	ctx context.Context,
	lggr logger.Logger,
	ccipReader cciptypes.CCIPReader,
	msgHasher cciptypes.MessageHasher,
	readableChains mapset.Set[cciptypes.ChainSelector],
	latestCommittedSeqNums []cciptypes.SeqNumChain,
	msgScanBatchSize int,
) ([]cciptypes.CCIPMsg, error) {
	// Find the new msgs for each supported chain based on the discovered max sequence numbers.
	newMsgsPerChain := make([][]cciptypes.CCIPMsg, len(latestCommittedSeqNums))
	eg := new(errgroup.Group)

	for chainIdx, seqNumChain := range latestCommittedSeqNums {
		if !readableChains.Contains(seqNumChain.ChainSel) {
			lggr.Debugw("reading chain is not supported", "chain", seqNumChain.ChainSel)
			continue
		}

		seqNumChain := seqNumChain
		chainIdx := chainIdx
		eg.Go(func() error {
			minSeqNum := seqNumChain.SeqNum + 1
			maxSeqNum := minSeqNum + cciptypes.SeqNum(msgScanBatchSize)
			lggr.Debugw("scanning for new messages",
				"chain", seqNumChain.ChainSel, "minSeqNum", minSeqNum, "maxSeqNum", maxSeqNum)

			newMsgs, err := ccipReader.MsgsBetweenSeqNums(
				ctx, seqNumChain.ChainSel, cciptypes.NewSeqNumRange(minSeqNum, maxSeqNum))
			if err != nil {
				return fmt.Errorf("get messages between seq nums: %w", err)
			}

			if len(newMsgs) > 0 {
				lggr.Debugw("discovered new messages", "chain", seqNumChain.ChainSel, "newMsgs", len(newMsgs))
			} else {
				lggr.Debugw("no new messages discovered", "chain", seqNumChain.ChainSel)
			}

			for i := range newMsgs {
				h, err := msgHasher.Hash(ctx, newMsgs[i])
				if err != nil {
					return fmt.Errorf("hash message: %w", err)
				}
				newMsgs[i].MsgHash = h // populate msgHash field
			}

			newMsgsPerChain[chainIdx] = newMsgs
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("wait for new msg observations: %w", err)
	}

	observedNewMsgs := make([]cciptypes.CCIPMsg, 0)
	for chainIdx := range latestCommittedSeqNums {
		observedNewMsgs = append(observedNewMsgs, newMsgsPerChain[chainIdx]...)
	}
	return observedNewMsgs, nil
}

func observeTokenPrices(
	ctx context.Context,
	tokenPricesReader cciptypes.TokenPricesReader,
	tokens []types.Account,
) ([]cciptypes.TokenPrice, error) {
	tokenPrices, err := tokenPricesReader.GetTokenPricesUSD(ctx, tokens)
	if err != nil {
		return nil, fmt.Errorf("get token prices: %w", err)
	}

	if len(tokenPrices) != len(tokens) {
		return nil, fmt.Errorf("internal critical error token prices length mismatch: got %d, want %d",
			len(tokenPrices), len(tokens))
	}

	tokenPricesUSD := make([]cciptypes.TokenPrice, 0, len(tokens))
	for i, token := range tokens {
		tokenPricesUSD = append(tokenPricesUSD, cciptypes.NewTokenPrice(token, tokenPrices[i]))
	}

	return tokenPricesUSD, nil
}

func observeGasPrices(ctx context.Context, ccipReader cciptypes.CCIPReader, chains []cciptypes.ChainSelector) ([]cciptypes.GasPriceChain, error) {
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

	gasPricesGwei := make([]cciptypes.GasPriceChain, 0, len(chains))
	for i, chain := range chains {
		gasPricesGwei = append(gasPricesGwei, cciptypes.NewGasPriceChain(gasPrices[i].Int, chain))
	}

	return gasPricesGwei, nil
}

// newMsgsConsensus comes in consensus on the observed messages for each source chain. Generates one merkle root
// for each source chain based on the consensus on the messages.
func newMsgsConsensus(
	lggr logger.Logger,
	maxSeqNums []cciptypes.SeqNumChain,
	observations []cciptypes.CommitPluginObservation,
	fChainCfg map[cciptypes.ChainSelector]int,
) ([]cciptypes.MerkleRootChain, error) {
	maxSeqNumsPerChain := make(map[cciptypes.ChainSelector]cciptypes.SeqNum)
	for _, seqNumChain := range maxSeqNums {
		maxSeqNumsPerChain[seqNumChain.ChainSel] = seqNumChain.SeqNum
	}

	// Gather all messages from all observations.
	msgsFromObservations := make([]cciptypes.CCIPMsgBaseDetails, 0)
	for _, obs := range observations {
		msgsFromObservations = append(msgsFromObservations, obs.NewMsgs...)
	}
	lggr.Debugw("total observed messages across all followers", "msgs", len(msgsFromObservations))

	// Filter out messages less than or equal to the max sequence numbers.
	msgsFromObservations = slicelib.Filter(msgsFromObservations, func(msg cciptypes.CCIPMsgBaseDetails) bool {
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
		func(msg cciptypes.CCIPMsgBaseDetails) cciptypes.ChainSelector { return msg.SourceChain },
	)

	// Come to consensus on the observed messages by source chain.
	consensusBySourceChain := make(map[cciptypes.ChainSelector]observedMsgsConsensus)
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

	merkleRoots := make([]cciptypes.MerkleRootChain, 0)
	for sourceChain, consensus := range consensusBySourceChain {
		merkleRoots = append(
			merkleRoots,
			cciptypes.NewMerkleRootChain(sourceChain, consensus.seqNumRange, consensus.merkleRoot),
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
	chainSel cciptypes.ChainSelector,
	observedMsgs []cciptypes.CCIPMsgBaseDetails,
	fChainCfg map[cciptypes.ChainSelector]int,
) (observedMsgsConsensus, error) {
	fChain, ok := fChainCfg[chainSel]
	if !ok {
		return observedMsgsConsensus{}, fmt.Errorf("fchain not found for chain %d", chainSel)
	}
	lggr.Debugw("observed messages consensus",
		"chain", chainSel, "fChain", fChain, "observedMsgs", len(observedMsgs))

	// First come to consensus about the (sequence number, msg hash) pairs.
	// For each sequence number consider the Hash with the most votes.
	msgSeqNumToHashCounts := make(map[cciptypes.SeqNum]map[string]int) // seqNum -> msgHash -> count
	for _, msg := range observedMsgs {
		if _, exists := msgSeqNumToHashCounts[msg.SeqNum]; !exists {
			msgSeqNumToHashCounts[msg.SeqNum] = make(map[string]int)
		}
		msgSeqNumToHashCounts[msg.SeqNum][msg.MsgHash.String()]++
	}
	lggr.Debugw("observed message counts", "chain", chainSel, "msgSeqNumToHashCounts", msgSeqNumToHashCounts)

	msgObservationsCount := make(map[cciptypes.SeqNum]int)
	msgSeqNumToHash := make(map[cciptypes.SeqNum]cciptypes.Bytes32)
	for seqNum, hashCounts := range msgSeqNumToHashCounts {
		if len(hashCounts) == 0 {
			lggr.Fatalw("hash counts should never be empty", "seqNum", seqNum)
			continue
		}

		// Find the MsgHash with the most votes for each sequence number.
		hashesSlice := make([]string, 0, len(hashCounts))
		for h := range hashCounts {
			hashesSlice = append(hashesSlice, h)
		}
		// determinism in case we have the same count for different hashes
		sort.Slice(hashesSlice, func(i, j int) bool { return hashesSlice[i] < hashesSlice[j] })

		maxCnt := hashCounts[hashesSlice[0]]
		mostVotedHash := hashesSlice[0]
		for _, h := range hashesSlice[1:] {
			cnt := hashCounts[h]
			if cnt > maxCnt {
				maxCnt = cnt
				mostVotedHash = h
			}
		}

		msgObservationsCount[seqNum] = maxCnt
		hashBytes, err := cciptypes.NewBytes32FromString(mostVotedHash)
		if err != nil {
			return observedMsgsConsensus{}, fmt.Errorf("critical issue converting hash '%s' to bytes32: %w",
				mostVotedHash, err)
		}
		msgSeqNumToHash[seqNum] = hashBytes
	}
	lggr.Debugw("observed message consensus", "chain", chainSel, "msgSeqNumToHash", msgSeqNumToHash)

	// Filter out msgs not observed by at least 2f_chain+1 followers.
	msgSeqNumsQuorum := mapset.NewSet[cciptypes.SeqNum]()
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
	seqNumConsensusRange := cciptypes.NewSeqNumRange(msgSeqNumsQuorumSlice[0], msgSeqNumsQuorumSlice[0])
	for _, seqNum := range msgSeqNumsQuorumSlice[1:] {
		if seqNum != seqNumConsensusRange.End()+1 {
			break // Found a gap in the sequence numbers.
		}
		seqNumConsensusRange.SetEnd(seqNum)
	}

	treeLeaves := make([][32]byte, 0)
	for seqNum := seqNumConsensusRange.Start(); seqNum <= seqNumConsensusRange.End(); seqNum++ {
		msgHash, ok := msgSeqNumToHash[seqNum]
		if !ok {
			return observedMsgsConsensus{}, fmt.Errorf("msg hash not found for seq num %d", seqNum)
		}
		treeLeaves = append(treeLeaves, msgHash)
	}

	lggr.Debugw("constructing merkle tree", "chain", chainSel, "treeLeaves", len(treeLeaves))
	tree, err := merklemulti.NewTree(hashutil.NewKeccak(), treeLeaves)
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
func maxSeqNumsConsensus(lggr logger.Logger, fChain int, observations []cciptypes.CommitPluginObservation) []cciptypes.SeqNumChain {
	observedSeqNumsPerChain := make(map[cciptypes.ChainSelector][]cciptypes.SeqNum)
	for _, obs := range observations {
		for _, maxSeqNum := range obs.MaxSeqNums {
			if _, exists := observedSeqNumsPerChain[maxSeqNum.ChainSel]; !exists {
				observedSeqNumsPerChain[maxSeqNum.ChainSel] = make([]cciptypes.SeqNum, 0)
			}
			observedSeqNumsPerChain[maxSeqNum.ChainSel] = append(observedSeqNumsPerChain[maxSeqNum.ChainSel], maxSeqNum.SeqNum)
		}
	}

	seqNums := make([]cciptypes.SeqNumChain, 0, len(observedSeqNumsPerChain))
	for ch, observedSeqNums := range observedSeqNumsPerChain {
		if len(observedSeqNums) < 2*fChain+1 {
			lggr.Warnw("not enough observations for chain", "chain", ch, "observedSeqNums", observedSeqNums)
			continue
		}

		sort.Slice(observedSeqNums, func(i, j int) bool { return observedSeqNums[i] < observedSeqNums[j] })
		seqNums = append(seqNums, cciptypes.NewSeqNumChain(ch, observedSeqNums[fChain]))
	}

	sort.Slice(seqNums, func(i, j int) bool { return seqNums[i].ChainSel < seqNums[j].ChainSel })
	return seqNums
}

// tokenPricesConsensus returns the median price for tokens that have at least 2f_chain+1 observations.
func tokenPricesConsensus(
	observations []cciptypes.CommitPluginObservation,
	fChain int,
) ([]cciptypes.TokenPrice, error) {
	pricesPerToken := make(map[types.Account][]cciptypes.BigInt)
	for _, obs := range observations {
		for _, price := range obs.TokenPrices {
			if _, exists := pricesPerToken[price.TokenID]; !exists {
				pricesPerToken[price.TokenID] = make([]cciptypes.BigInt, 0)
			}
			pricesPerToken[price.TokenID] = append(pricesPerToken[price.TokenID], price.Price)
		}
	}

	// Keep the median
	consensusPrices := make([]cciptypes.TokenPrice, 0)
	for token, prices := range pricesPerToken {
		if len(prices) < 2*fChain+1 {
			continue
		}
		consensusPrices = append(consensusPrices, cciptypes.NewTokenPrice(token, slicelib.BigIntSortedMiddle(prices).Int))
	}

	sort.Slice(consensusPrices, func(i, j int) bool { return consensusPrices[i].TokenID < consensusPrices[j].TokenID })
	return consensusPrices, nil
}

func gasPricesConsensus(lggr logger.Logger, observations []cciptypes.CommitPluginObservation, fChain int) []cciptypes.GasPriceChain {
	// Group the observed gas prices by chain.
	gasPricePerChain := make(map[cciptypes.ChainSelector][]cciptypes.BigInt)
	for _, obs := range observations {
		for _, gasPrice := range obs.GasPrices {
			if _, exists := gasPricePerChain[gasPrice.ChainSel]; !exists {
				gasPricePerChain[gasPrice.ChainSel] = make([]cciptypes.BigInt, 0)
			}
			gasPricePerChain[gasPrice.ChainSel] = append(gasPricePerChain[gasPrice.ChainSel], gasPrice.GasPrice)
		}
	}

	// Keep the median
	consensusGasPrices := make([]cciptypes.GasPriceChain, 0)
	for chain, gasPrices := range gasPricePerChain {
		if len(gasPrices) < 2*fChain+1 {
			lggr.Warnw("not enough gas price observations", "chain", chain, "gasPrices", gasPrices)
			continue
		}

		consensusGasPrices = append(
			consensusGasPrices,
			cciptypes.NewGasPriceChain(slicelib.BigIntSortedMiddle(gasPrices).Int, chain),
		)
	}

	sort.Slice(consensusGasPrices, func(i, j int) bool { return consensusGasPrices[i].ChainSel < consensusGasPrices[j].ChainSel })
	return consensusGasPrices
}

// fChainConsensus comes to consensus on the plugin config based on the observations.
// We cannot trust the state of a single follower, so we need to come to consensus on the config.
func fChainConsensus(
	observations []cciptypes.CommitPluginObservation, // observations from all followers
) map[cciptypes.ChainSelector]int {
	// Come to consensus on fChain.
	// Use the fChain observed by most followers for each chain.
	fChainCounts := make(map[cciptypes.ChainSelector]map[int]int) // {chain: {fChain: count}}
	for _, obs := range observations {
		for chain, fChain := range obs.FChain {
			if _, exists := fChainCounts[chain]; !exists {
				fChainCounts[chain] = make(map[int]int)
			}
			fChainCounts[chain][fChain]++
		}
	}
	consensusFChain := make(map[cciptypes.ChainSelector]int)
	for chain, counts := range fChainCounts {
		maxCount := 0
		for fChain, count := range counts {
			if count > maxCount {
				maxCount = count
				consensusFChain[chain] = fChain
			}
		}
	}

	return consensusFChain
}

// validateObservedSequenceNumbers checks if the sequence numbers of the provided messages are unique for each chain and
// that they match the observed max sequence numbers.
func validateObservedSequenceNumbers(msgs []cciptypes.CCIPMsgBaseDetails, maxSeqNums []cciptypes.SeqNumChain) error {
	// If the observer did not include sequence numbers it means that it's not a destination chain reader.
	// In that case we cannot do any msg sequence number validations.
	if len(maxSeqNums) == 0 {
		return nil
	}

	// MaxSeqNums must be unique for each chain.
	maxSeqNumsMap := make(map[cciptypes.ChainSelector]cciptypes.SeqNum)
	for _, maxSeqNum := range maxSeqNums {
		if _, exists := maxSeqNumsMap[maxSeqNum.ChainSel]; exists {
			return fmt.Errorf("duplicate max sequence number for chain %d", maxSeqNum.ChainSel)
		}
		maxSeqNumsMap[maxSeqNum.ChainSel] = maxSeqNum.SeqNum
	}

	seqNums := make(map[cciptypes.ChainSelector]mapset.Set[cciptypes.SeqNum], len(msgs))
	hashes := mapset.NewSet[string]()
	for _, msg := range msgs {
		if msg.MsgHash.IsEmpty() {
			return fmt.Errorf("observed msg hash must not be empty")
		}

		if _, exists := seqNums[msg.SourceChain]; !exists {
			seqNums[msg.SourceChain] = mapset.NewSet[cciptypes.SeqNum]()
		}

		// The same sequence number must not appear more than once for the same chain and must be valid.
		if seqNums[msg.SourceChain].Contains(msg.SeqNum) {
			return fmt.Errorf("duplicate sequence number %d for chain %d", msg.SeqNum, msg.SourceChain)
		}
		seqNums[msg.SourceChain].Add(msg.SeqNum)

		// The observed msg hash cannot appear twice for different msgs.
		if hashes.Contains(msg.MsgHash.String()) {
			return fmt.Errorf("duplicate msg hash %s", msg.MsgHash.String())
		}
		hashes.Add(msg.MsgHash.String())

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
	msgs []cciptypes.CCIPMsgBaseDetails,
	seqNums []cciptypes.SeqNumChain,
	nodeSupportedChains mapset.Set[cciptypes.ChainSelector],
	destChain cciptypes.ChainSelector,
) error {

	if len(seqNums) > 0 && !nodeSupportedChains.Contains(destChain) {
		return fmt.Errorf("observer must be a writer if it observes sequence numbers")
	}

	if len(msgs) == 0 {
		return nil
	}

	for _, msg := range msgs {
		// Observer must be able to read the chain that the message is coming from.
		if !nodeSupportedChains.Contains(msg.SourceChain) {
			return fmt.Errorf("observer not allowed to read chain %d", msg.SourceChain)
		}
	}

	return nil
}

func validateObservedTokenPrices(tokenPrices []cciptypes.TokenPrice) error {
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

func validateObservedGasPrices(gasPrices []cciptypes.GasPriceChain) error {
	// Duplicate gas prices must not appear for the same chain and must not be empty.
	gasPriceChains := mapset.NewSet[cciptypes.ChainSelector]()
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
	seqNumRange cciptypes.SeqNumRange
	merkleRoot  [32]byte
}

func (o observedMsgsConsensus) isEmpty() bool {
	return o.seqNumRange.Start() == 0 && o.seqNumRange.End() == 0 && o.merkleRoot == [32]byte{}
}
