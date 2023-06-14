package mercury

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/pkg/errors"
)

// NOTE: All aggregate functions assume at least one element in the passed slice
// The passed slice might be mutated (sorted)

// GetConsensusTimestamp gets the median timestamp
func GetConsensusTimestamp(paos []ParsedAttributedObservation) uint32 {
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Timestamp < paos[j].Timestamp
	})
	return paos[len(paos)/2].Timestamp
}

// GetConsensusBenchmarkPrice gets the median benchmark price
func GetConsensusBenchmarkPrice(paos []ParsedAttributedObservation, f int) (*big.Int, error) {
	var validBenchmarkPrices []*big.Int
	for _, pao := range paos {
		if pao.PricesValid {
			validBenchmarkPrices = append(validBenchmarkPrices, pao.BenchmarkPrice)
		}
	}
	if len(validBenchmarkPrices) < f+1 {
		return nil, errors.New("fewer than f+1 observations have a valid price")
	}
	sort.Slice(validBenchmarkPrices, func(i, j int) bool {
		return validBenchmarkPrices[i].Cmp(validBenchmarkPrices[j]) < 0
	})

	return validBenchmarkPrices[len(validBenchmarkPrices)/2], nil
}

// GetConsensusBid gets the median bid
func GetConsensusBid(paos []ParsedAttributedObservation, f int) (*big.Int, error) {
	var validBids []*big.Int
	for _, pao := range paos {
		if pao.PricesValid {
			validBids = append(validBids, pao.Bid)
		}
	}
	if len(validBids) < f+1 {
		return nil, errors.New("fewer than f+1 observations have a valid price")
	}
	sort.Slice(validBids, func(i, j int) bool {
		return validBids[i].Cmp(validBids[j]) < 0
	})

	return validBids[len(validBids)/2], nil
}

// GetConsensusAsk gets the median ask
func GetConsensusAsk(paos []ParsedAttributedObservation, f int) (*big.Int, error) {
	var validAsks []*big.Int
	for _, pao := range paos {
		if pao.PricesValid {
			validAsks = append(validAsks, pao.Ask)
		}
	}
	if len(validAsks) < f+1 {
		return nil, errors.New("fewer than f+1 observations have a valid price")
	}
	sort.Slice(validAsks, func(i, j int) bool {
		return validAsks[i].Cmp(validAsks[j]) < 0
	})

	return validAsks[len(validAsks)/2], nil
}

type block struct {
	hash string
	num  int64
	ts   uint64
}

func (b1 block) less(b2 block) bool {
	if b1.num == b2.num && b1.ts == b2.ts {
		// tie-break on hash, all else being equal
		return b1.hash < b2.hash
	} else if b1.num == b2.num {
		// if block number is equal and timestamps differ, take the latest timestamp
		return b1.ts > b2.ts
	} else {
		// if block number is different, take the higher block number
		return b1.num > b2.num
	}
}

// GetConsensusCurrentBlock gets the most common (mode) block hash/number/timestamps.
// In the event of a tie, use the lowest numerical value
func GetConsensusCurrentBlock(paos []ParsedAttributedObservation, f int) (hash []byte, num int64, ts uint64, err error) {
	var validPaos []ParsedAttributedObservation
	for _, pao := range paos {
		if pao.CurrentBlockValid {
			validPaos = append(validPaos, pao)
		}
	}
	if len(validPaos) < f+1 {
		return nil, 0, 0, fmt.Errorf("fewer than f+1 observations have a valid current block (got: %d/%d)", len(validPaos), len(paos))
	}

	m := map[block]int{}
	maxCnt := 0
	for _, pao := range validPaos {
		b := block{string(pao.CurrentBlockHash), pao.CurrentBlockNum, pao.CurrentBlockTimestamp}
		m[b]++
		if cnt := m[b]; cnt > maxCnt {
			maxCnt = cnt
		}
	}

	if maxCnt < f+1 {
		return nil, 0, 0, errors.New("no unique block with at least f+1 votes")
	}

	var blocks []block
	for b, cnt := range m {
		if cnt == maxCnt {
			blocks = append(blocks, b)
		}
	}
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].less(blocks[j])
	})

	return []byte(blocks[0].hash), blocks[0].num, blocks[0].ts, nil
}

// GetConsensusMaxFinalizedBlockNum gets the most common (mode)
// ConsensusMaxFinalizedBlockNum In the event of a tie, the lower number is
// chosen
func GetConsensusMaxFinalizedBlockNum(paos []ParsedAttributedObservation, f int) (int64, error) {
	var validPaos []ParsedAttributedObservation
	for _, pao := range paos {
		if pao.MaxFinalizedBlockNumberValid {
			validPaos = append(validPaos, pao)
		}
	}
	if len(validPaos) < f+1 {
		return 0, fmt.Errorf("fewer than f+1 observations have a valid maxFinalizedBlockNumber (got: %d/%d)", len(validPaos), len(paos))
	}
	// pick the most common block number with at least f+1 votes
	m := map[int64]int{}
	maxCnt := 0
	for _, pao := range validPaos {
		n := pao.MaxFinalizedBlockNumber
		m[n]++
		if cnt := m[n]; cnt > maxCnt {
			maxCnt = cnt
		}
	}

	var nums []int64
	for num, cnt := range m {
		if cnt == maxCnt {
			nums = append(nums, num)
		}
	}

	if maxCnt < f+1 {
		return 0, fmt.Errorf("no valid maxFinalizedBlockNumber with at least f+1 votes (got counts: %v)", m)
	}
	// guaranteed to be at least one num after this

	// determistic tie-break for number
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	return nums[0], nil
}
