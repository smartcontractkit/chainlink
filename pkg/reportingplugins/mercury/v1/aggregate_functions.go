package mercury_v1

import (
	"fmt"
	"sort"
)

// GetConsensusLatestBlock gets the latest block that has at least f+1 votes
// Assumes that LatestBlocks are ordered by block number desc
func GetConsensusLatestBlock(paos []PAO, f int) (hash []byte, num int64, ts uint64, err error) {
	// observed blocks grouped by their block number
	groupingsM := make(map[int64][]Block)
	for _, pao := range paos {
		if blocks := pao.GetLatestBlocks(); len(blocks) > 0 {
			for _, block := range blocks {
				groupingsM[block.Num] = append(groupingsM[block.Num], block)
			}
		} else { // DEPRECATED
			// TODO: Remove this handling after deployment (https://smartcontract-it.atlassian.net/browse/MERC-2272)
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
			groupingsM[blockNum] = append(groupingsM[blockNum], NewBlock(blockNum, blockHash, blockTs))
		}
	}

	// sort by latest block number desc
	groupings := make([][]Block, len(groupingsM))
	{
		i := 0
		for _, blocks := range groupingsM {
			groupings[i] = blocks
			i++
		}
	}

	// each grouping will have all blocks with the same block number, sorted desc
	sort.Slice(groupings, func(i, j int) bool {
		return groupings[i][0].Num > groupings[j][0].Num
	})

	// take highest block number with at least f+1 in agreement on everything
	for _, blocks := range groupings {
		m := map[Block]int{}
		maxCnt := 0
		// count unique blocks
		for _, b := range blocks {
			m[b]++
			if cnt := m[b]; cnt > maxCnt {
				maxCnt = cnt
			}
		}
		if maxCnt >= f+1 {
			// at least one set of blocks has f+1 in agreement

			// take the blocks with highest count
			var usableBlocks []Block
			for b, cnt := range m {
				if cnt == maxCnt {
					usableBlocks = append(usableBlocks, b)
				}
			}
			sort.Slice(usableBlocks, func(i, j int) bool {
				return usableBlocks[j].less(usableBlocks[i])
			})

			return usableBlocks[0].HashBytes(), usableBlocks[0].Num, usableBlocks[0].Ts, nil

		}
		// this grouping does not have any identical blocks with at least f+1 in agreement, try next block number down
	}

	return nil, 0, 0, fmt.Errorf("cannot come to consensus on latest block number, got observations: %#v", paos)
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
