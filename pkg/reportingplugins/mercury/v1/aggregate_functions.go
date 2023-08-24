package mercury_v1

import (
	"fmt"
	"sort"
)

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
func GetConsensusCurrentBlock(paos []PAO, f int) (hash []byte, num int64, ts uint64, err error) {
	m := map[block]int{}
	maxCnt := 0
	var validObsCnt int
	for _, pao := range paos {
		blockHash, valid := pao.GetCurrentBlockHash()
		if !valid {
			continue
		}
		blockNum, valid := pao.GetCurrentBlockNum()
		if !valid {
			continue
		}
		blockTs, valid := pao.GetCurrentBlockTimestamp()
		if !valid {
			continue
		}
		if valid {
			validObsCnt++
			b := block{string(blockHash), blockNum, blockTs}
			m[b]++
			if cnt := m[b]; cnt > maxCnt {
				maxCnt = cnt
			}
		}
	}

	if validObsCnt < f+1 {
		return nil, 0, 0, fmt.Errorf("fewer than f+1 observations have a valid current block (got: %d/%d, f=%d)", validObsCnt, len(paos), f)
	}

	if maxCnt < f+1 {
		return nil, 0, 0, fmt.Errorf("no unique valid block observation with at least f+1 votes (got %d/%d, f=%d)", maxCnt, len(paos), f)
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
func GetConsensusMaxFinalizedBlockNum(paos []PAO, f int) (int64, error) {
	var validPaos []PAO
	for _, pao := range paos {
		_, valid := pao.GetMaxFinalizedBlockNumber()
		if valid {
			validPaos = append(validPaos, pao)
		}
	}
	if len(validPaos) < f+1 {
		return 0, fmt.Errorf("fewer than f+1 observations have a valid maxFinalizedBlockNumber (got: %d/%d, f=%d)", len(validPaos), len(paos), f)
	}
	// pick the most common block number with at least f+1 votes
	m := map[int64]int{}
	maxCnt := 0
	for _, pao := range validPaos {
		n, _ := pao.GetMaxFinalizedBlockNumber()
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
		return 0, fmt.Errorf("no valid maxFinalizedBlockNumber with at least f+1 votes (got counts: %v, f=%d)", m, f)
	}
	// guaranteed to be at least one num after this

	// determistic tie-break for number
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	return nums[0], nil
}
