package mercury

import (
	"math/big"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// TODO: add other fields, bid, ask etc
func ValidateReport(r map[string]interface{}) error {
	feedIdInterface, ok := r["feedId"]
	if !ok {
		return errors.Errorf("unpacked report has no 'feedId'")
	}
	feedID, ok := feedIdInterface.([32]byte)
	if !ok {
		return errors.Errorf("cannot cast feedId to [32]byte, type is %T", feedID)
	}
	log.Trace().Str("FeedID", string(feedID[:])).Msg("Feed ID")

	priceInterface, ok := r["median"]
	if !ok {
		return errors.Errorf("unpacked report has no 'median'")
	}
	medianPrice, ok := priceInterface.(*big.Int)
	if !ok {
		return errors.Errorf("cannot cast median to *big.Int, type is %T", medianPrice)
	}
	log.Trace().Int64("Price", medianPrice.Int64()).Msg("Median price")

	observationsBlockNumberInterface, ok := r["observationsBlocknumber"]
	if !ok {
		return errors.Errorf("unpacked report has no 'observationsBlocknumber'")
	}
	observationsBlockNumber, ok := observationsBlockNumberInterface.(uint64)
	if !ok {
		return errors.Errorf("cannot cast observationsBlocknumber to uint64, type is %T", observationsBlockNumber)
	}
	log.Trace().Uint64("Block", observationsBlockNumber).Msg("Observation block number")

	observationsTimestampInterface, ok := r["observationsTimestamp"]
	if !ok {
		return errors.Errorf("unpacked report has no 'observationsTimestamp'")
	}
	observationsTimestamp, ok := observationsTimestampInterface.(uint32)
	if !ok {
		return errors.Errorf("cannot cast observationsTimestamp to uint32, type is %T", observationsTimestamp)
	}
	log.Trace().Uint32("Timestamp", observationsTimestamp).Msg("Observation timestamp")

	return nil
}
