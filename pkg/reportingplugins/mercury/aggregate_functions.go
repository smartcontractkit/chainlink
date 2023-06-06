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
	// pick the most common blockhash with at least f+1 votes
	hash, err = getConsensusCurrentBlockHash(validPaos, f+1)
	if err != nil {
		return hash, 0, 0, errors.Wrap(err, "couldn't get consensus current block")
	}

	// pick the most common block number with at least f+1 votes
	num, err = getConsensusCurrentBlockNum(validPaos, string(hash), f+1)
	if err != nil {
		return hash, num, 0, errors.Wrap(err, "couldn't get consensus current block")
	}

	// pick the most common block timestamp with at least f+1 votes
	ts, err = getConsensusCurrentBlockTimestamp(validPaos, string(hash), num, f+1)
	if err != nil {
		return hash, num, ts, errors.Wrap(err, "couldn't get consensus current block")
	}

	return hash, num, ts, nil
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

func getConsensusCurrentBlockHash(paos []ParsedAttributedObservation, threshold int) (hash []byte, err error) {
	m := map[string]int{}
	maxCnt := 0
	for _, pao := range paos {
		h := pao.CurrentBlockHash
		m[string(h)]++
		if cnt := m[string(h)]; cnt > maxCnt {
			maxCnt = cnt
		}
	}

	if maxCnt < threshold {
		return nil, errors.New("no block hash with at least f+1 votes")
	}

	var hashes []string
	for hash, cnt := range m {
		if cnt == maxCnt {
			hashes = append(hashes, hash)
		}
	}

	// determistic tie-break for hash
	sort.Slice(hashes, func(i, j int) bool {
		return hashes[i] < hashes[j]
	})

	hash = []byte(hashes[0])
	return
}

func getConsensusCurrentBlockNum(paos []ParsedAttributedObservation, blockHash string, threshold int) (num int64, err error) {
	var matchingPaos []ParsedAttributedObservation
	for _, pao := range paos {
		if string(pao.CurrentBlockHash) == blockHash {
			matchingPaos = append(matchingPaos, pao)
		}
	}

	m := map[int64]int{}
	maxCnt := 0
	for _, pao := range matchingPaos {
		n := pao.CurrentBlockNum
		m[n]++
		if cnt := m[n]; cnt > maxCnt {
			maxCnt = cnt
		}
	}

	if maxCnt < threshold {
		return 0, errors.Errorf("no block number matching hash 0x%x with at least f+1 votes", blockHash)
	}

	var nums []int64
	for num, cnt := range m {
		if cnt == maxCnt {
			nums = append(nums, num)
		}
	}

	// determistic tie-break for num
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	num = nums[0]
	return
}

func getConsensusCurrentBlockTimestamp(paos []ParsedAttributedObservation, blockHash string, blockNum int64, threshold int) (ts uint64, err error) {
	var matchingPaos []ParsedAttributedObservation
	for _, pao := range paos {
		if string(pao.CurrentBlockHash) == blockHash && pao.CurrentBlockNum == blockNum {
			matchingPaos = append(matchingPaos, pao)
		}
	}

	m := map[uint64]int{}
	maxCnt := 0
	for _, pao := range matchingPaos {
		n := pao.CurrentBlockTimestamp
		m[n]++
		if cnt := m[n]; cnt > maxCnt {
			maxCnt = cnt
		}
	}

	if maxCnt < threshold {
		return 0, errors.Errorf("no block timestamp matching block hash 0x%x and block number %d with at least f+1 votes", blockHash, blockNum)
	}

	var timestamps []uint64
	for ts, cnt := range m {
		if cnt == maxCnt {
			timestamps = append(timestamps, ts)
		}
	}

	// determistic tie-break for timestamps
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i] < timestamps[j]
	})

	ts = timestamps[0]
	return
}
