package mercury

import (
	"fmt"
	"math/big"
	"sort"
)

// NOTE: All aggregate functions assume at least one element in the passed slice
// The passed slice might be mutated (sorted)

// GetConsensusTimestamp gets the median timestamp
func GetConsensusTimestamp(paos []ParsedAttributedObservation) uint32 {
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].GetTimestamp() < paos[j].GetTimestamp()
	})
	return paos[len(paos)/2].GetTimestamp()
}

// GetConsensusBenchmarkPrice gets the median benchmark price
func GetConsensusBenchmarkPrice(paos []ParsedAttributedObservation, f int) (*big.Int, error) {
	var validBenchmarkPrices []*big.Int
	for _, pao := range paos {
		bmPrice, valid := pao.GetBenchmarkPrice()
		if valid {
			validBenchmarkPrices = append(validBenchmarkPrices, bmPrice)
		}
	}

	if len(validBenchmarkPrices) < f+1 {
		return nil, fmt.Errorf("fewer than f+1 observations have a valid price (got: %d/%d)", len(validBenchmarkPrices), len(paos))
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
		bid, valid := pao.GetBid()
		if valid {
			validBids = append(validBids, bid)
		}
	}
	if len(validBids) < f+1 {
		return nil, fmt.Errorf("fewer than f+1 observations have a valid price (got: %d/%d)", len(validBids), len(paos))
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
		ask, valid := pao.GetAsk()
		if valid {
			validAsks = append(validAsks, ask)
		}
	}
	if len(validAsks) < f+1 {
		return nil, fmt.Errorf("fewer than f+1 observations have a valid price (got: %d/%d)", len(validAsks), len(paos))
	}
	sort.Slice(validAsks, func(i, j int) bool {
		return validAsks[i].Cmp(validAsks[j]) < 0
	})

	return validAsks[len(validAsks)/2], nil
}

func GetConsensusMaxFinalizedTimestamp(paos []ParsedAttributedObservation, f int) (uint32, error) {
	var validTimestampCount int
	timestampFrequencyMap := map[uint32]int{}
	for _, pao := range paos {
		ts, valid := pao.GetMaxFinalizedTimestamp()
		if valid {
			validTimestampCount++
			timestampFrequencyMap[ts]++
		}
	}

	// check if we have enough valid timestamps
	if validTimestampCount < f+1 {
		return 0, fmt.Errorf("fewer than f+1 observations have a valid maxFinalizedTimestamp (got: %d/%d)", validTimestampCount, len(paos))
	}

	var timestampFrequencyMaxCnt int
	for _, cnt := range timestampFrequencyMap {
		if cnt > timestampFrequencyMaxCnt {
			timestampFrequencyMaxCnt = cnt
		}
	}

	// check if we have enough valid timestamps with the max frequency
	if timestampFrequencyMaxCnt < f+1 {
		return 0, fmt.Errorf("no valid maxFinalizedTimestamp with at least f+1 votes (got counts: %v)", timestampFrequencyMap)
	}

	// select timestamps with the max frequency (in case there are more than one)
	// sort them deterministically
	var validTimestampsWithMaxFrequency []uint32
	for ts, cnt := range timestampFrequencyMap {
		if cnt == timestampFrequencyMaxCnt {
			validTimestampsWithMaxFrequency = append(validTimestampsWithMaxFrequency, ts)
		}
	}
	sort.Slice(validTimestampsWithMaxFrequency, func(i, j int) bool {
		return validTimestampsWithMaxFrequency[i] < validTimestampsWithMaxFrequency[j]
	})

	return validTimestampsWithMaxFrequency[0], nil
}

// GetConsensusLinkFee gets the median link fee
func GetConsensusLinkFee(paos []ParsedAttributedObservation, f int) (*big.Int, error) {
	var validLinkFees []*big.Int
	for _, pao := range paos {
		fee, valid := pao.GetLinkFee()
		if valid {
			validLinkFees = append(validLinkFees, fee)
		}
	}
	if len(validLinkFees) < f+1 {
		return nil, fmt.Errorf("fewer than f+1 observations have a valid linkFee (got: %d/%d)", len(validLinkFees), len(paos))
	}
	sort.Slice(validLinkFees, func(i, j int) bool {
		return validLinkFees[i].Cmp(validLinkFees[j]) < 0
	})

	return validLinkFees[len(validLinkFees)/2], nil
}

// GetConsensusNativeFee gets the median native fee
func GetConsensusNativeFee(paos []ParsedAttributedObservation, f int) (*big.Int, error) {
	var validNativeFees []*big.Int
	for _, pao := range paos {
		fee, valid := pao.GetNativeFee()
		if valid {
			validNativeFees = append(validNativeFees, fee)
		}
	}
	if len(validNativeFees) < f+1 {
		return nil, fmt.Errorf("fewer than f+1 observations have a valid nativeFee (got: %d/%d)", len(validNativeFees), len(paos))
	}
	sort.Slice(validNativeFees, func(i, j int) bool {
		return validNativeFees[i].Cmp(validNativeFees[j]) < 0
	})

	return validNativeFees[len(validNativeFees)/2], nil
}
