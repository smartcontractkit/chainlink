package jobs

type StreamType string

const (
	StreamTypeQuote        = StreamType("quote")
	StreamTypeMedian       = StreamType("median")
	StreamTypeMarketStatus = StreamType("market-status")
	StreamTypeFundingRate  = StreamType("funding-rate")
	StreamTypeRWA          = StreamType("rwa")
	StreamTypeDEX          = StreamType("dex")
	StreamTypeDataLink     = StreamType("data-link")
	// StreamTypeConsolidated is used for the consolidated stream type
	StreamTypeConsolidated = StreamType("consolidated")
)

func (st StreamType) Valid() bool {
	switch st {
	case StreamTypeQuote, StreamTypeMedian:
		return true
	case StreamTypeMarketStatus, StreamTypeDataLink, StreamTypeConsolidated:
		return false // these are not implemented, yet.
	}

	return false
}
