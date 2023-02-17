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

// GetConsensusCurrentBlock gets the most common (mode) block hash/number.
// In the event of a tie, use the lowest numerical value
func GetConsensusCurrentBlock(paos []ParsedAttributedObservation, f int) (hash []byte, num int64, err error) {
	// pick the most common blockhash with at least f+1 votes

	// cast to string for map key; this case is optimised by the go compiler:
	// https://github.com/golang/go/commit/f5f5a8b6209f84961687d993b93ea0d397f5d5bf#diff-3437cd20ec7506421b8d8b653efa9bfe
	m := map[string]int{}
	maxCnt := 0
	for _, pao := range paos {
		h := pao.CurrentBlockHash
		m[string(h)]++
		if cnt := m[string(h)]; cnt > maxCnt {
			maxCnt = cnt
		}
	}

	if maxCnt < f+1 {
		return nil, 0, errors.New("no block hash with at least f+1 votes")
	}
	// guaranteed to be at least one hash after this

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

	// get block numbers for observations with matching hash
	var matchingPaos []ParsedAttributedObservation
	for _, pao := range paos {
		if string(pao.CurrentBlockHash) == string(hash) {
			matchingPaos = append(matchingPaos, pao)
		}
	}

	// pick the most common block number with at least f+1 votes
	m2 := map[int64]int{}
	maxCnt = 0
	for _, pao := range matchingPaos {
		n := pao.CurrentBlockNum
		m2[n]++
		if cnt := m2[n]; cnt > maxCnt {
			maxCnt = cnt
		}
	}

	var nums []int64
	for num, cnt := range m2 {
		if cnt == maxCnt {
			nums = append(nums, num)
		}
	}

	if maxCnt < f+1 {
		return nil, 0, errors.Errorf("no block number matching hash 0x%x with at least f+1 votes", hash)
	}
	// guaranteed to be at least one num after this

	// determistic tie-break for number
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	num = nums[0]

	return hash, num, nil
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
