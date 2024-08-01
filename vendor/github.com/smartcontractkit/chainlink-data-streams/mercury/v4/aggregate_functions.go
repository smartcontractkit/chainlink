package v4

import "fmt"

type PAOMarketStatus interface {
	GetMarketStatus() (uint32, bool)
}

// GetConsensusMarketStatus gets the most common status, provided that it is at least F+1.
func GetConsensusMarketStatus(paos []PAOMarketStatus, f int) (uint32, error) {
	marketStatusCounts := make(map[uint32]int)
	for _, pao := range paos {
		marketStatus, valid := pao.GetMarketStatus()
		if valid {
			marketStatusCounts[marketStatus]++
		}
	}

	var mostCommonMarketStatus uint32
	var mostCommonCount int
	for marketStatus, count := range marketStatusCounts {
		if count > mostCommonCount {
			mostCommonMarketStatus = marketStatus
			mostCommonCount = count
		} else if count == mostCommonCount {
			// For stability, always prefer the smaller enum value in case of ties.
			// In practice this will prefer CLOSED over OPEN.
			if marketStatus < mostCommonMarketStatus {
				mostCommonMarketStatus = marketStatus
			}
		}
	}

	if mostCommonCount < f+1 {
		return 0, fmt.Errorf("market status has fewer than f+1 observations (status %d got %d/%d)", mostCommonMarketStatus, mostCommonCount, len(paos))
	}

	return mostCommonMarketStatus, nil
}
