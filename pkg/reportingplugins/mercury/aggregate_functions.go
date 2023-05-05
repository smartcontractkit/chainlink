package mercury

import (
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
func GetConsensusBenchmarkPrice(paos []ParsedAttributedObservation) *big.Int {
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].BenchmarkPrice.Cmp(paos[j].BenchmarkPrice) < 0
	})

	return paos[len(paos)/2].BenchmarkPrice
}

// GetConsensusBid gets the median bid
func GetConsensusBid(paos []ParsedAttributedObservation) *big.Int {
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Bid.Cmp(paos[j].Bid) < 0
	})

	return paos[len(paos)/2].Bid
}

// GetConsensusAsk gets the median ask
func GetConsensusAsk(paos []ParsedAttributedObservation) *big.Int {
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Ask.Cmp(paos[j].Ask) < 0
	})

	return paos[len(paos)/2].Ask
}

// GetConsensusCurrentBlock gets the most common (mode) block hash/number/timestamps.
// In the event of a tie, use the lowest numerical value
func GetConsensusCurrentBlock(paos []ParsedAttributedObservation, f int) (hash []byte, num int64, ts uint64, err error) {
	// pick the most common blockhash with at least f+1 votes
	hash, err = getConsensusCurrentBlockHash(paos, f+1)
	if err != nil {
		return hash, 0, 0, errors.Wrap(err, "couldn't get consensus current block")
	}

	// pick the most common block number with at least f+1 votes
	num, err = getConsensusCurrentBlockNum(paos, string(hash), f+1)
	if err != nil {
		return hash, num, 0, errors.Wrap(err, "coulnd't get consensus current block")
	}

	// pick the most common block timestamp with at least f+1 votes
	ts, err = getConsensusCurrentBlockTimestamp(paos, string(hash), num, f+1)
	if err != nil {
		return hash, num, ts, errors.Wrap(err, "coulnd't get consensus current block")
	}

	return hash, num, ts, nil
}

// GetConsensusValidFromBlock gets the most common (mode) ValidFromBlockNum
// In the event of a tie, the lower number is chosen
func GetConsensusValidFromBlock(paos []ParsedAttributedObservation, f int) (int64, error) {
	// pick the most common block number with at least f+1 votes
	m := map[int64]int{}
	maxCnt := 0
	for _, pao := range paos {
		n := pao.ValidFromBlockNum
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
		return 0, errors.New("no valid from block number with at least f+1 votes")
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
