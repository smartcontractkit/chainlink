package mercury

import (
	"fmt"
	"math/big"
	"sort"
)

var Zero = big.NewInt(0)

// NOTE: All aggregate functions assume at least one element in the passed slice
// The passed slice might be mutated (sorted)

// GetConsensusTimestamp gets the median timestamp
func GetConsensusTimestamp(paos []PAO) uint32 {
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].GetTimestamp() < paos[j].GetTimestamp()
	})
	return paos[len(paos)/2].GetTimestamp()
}

// GetConsensusBenchmarkPrice gets the median benchmark price
func GetConsensusBenchmarkPrice(paos []PAO, f int) (*big.Int, error) {
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

type PAOBid interface {
	GetBid() (*big.Int, bool)
}

// GetConsensusBid gets the median bid
func GetConsensusBid(paos []PAOBid, f int) (*big.Int, error) {
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

type PAOAsk interface {
	GetAsk() (*big.Int, bool)
}

// GetConsensusAsk gets the median ask
func GetConsensusAsk(paos []PAOAsk, f int) (*big.Int, error) {
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

type PAOMaxFinalizedTimestamp interface {
	GetMaxFinalizedTimestamp() (int64, bool)
}

// GetConsensusMaxFinalizedTimestamp returns the highest count with > f observations
func GetConsensusMaxFinalizedTimestamp(paos []PAOMaxFinalizedTimestamp, f int) (int64, error) {
	var validTimestampCount int
	timestampFrequencyMap := map[int64]int{}
	for _, pao := range paos {
		ts, valid := pao.GetMaxFinalizedTimestamp()
		if valid {
			validTimestampCount++
			timestampFrequencyMap[ts]++
		}
	}

	// check if we have enough valid timestamps at all
	if validTimestampCount < f+1 {
		return 0, fmt.Errorf("fewer than f+1 observations have a valid maxFinalizedTimestamp (got: %d/%d)", validTimestampCount, len(paos))
	}

	var maxTs int64 = -2 // -1 is smallest valid amount
	for ts, cnt := range timestampFrequencyMap {
		// ignore any timestamps with <= f observations
		if cnt > f && ts > maxTs {
			maxTs = ts
		}
	}

	if maxTs < -1 {
		return 0, fmt.Errorf("no valid maxFinalizedTimestamp with at least f+1 votes (got counts: %v)", timestampFrequencyMap)
	}

	return maxTs, nil
}

type PAOLinkFee interface {
	GetLinkFee() (*big.Int, bool)
}

// GetConsensusLinkFee gets the median link fee
func GetConsensusLinkFee(paos []PAOLinkFee, f int) (*big.Int, error) {
	var validLinkFees []*big.Int
	for _, pao := range paos {
		fee, valid := pao.GetLinkFee()
		if valid && fee.Sign() >= 0 {
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

type PAONativeFee interface {
	GetNativeFee() (*big.Int, bool)
}

// GetConsensusNativeFee gets the median native fee
func GetConsensusNativeFee(paos []PAONativeFee, f int) (*big.Int, error) {
	var validNativeFees []*big.Int
	for _, pao := range paos {
		fee, valid := pao.GetNativeFee()
		if valid && fee.Sign() >= 0 {
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
