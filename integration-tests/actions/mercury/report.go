package mercury

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func getReportTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "feedId", Type: mustNewType("bytes32")},
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "benchmarkPrice", Type: mustNewType("int192")},
		{Name: "bid", Type: mustNewType("int192")},
		{Name: "ask", Type: mustNewType("int192")},
		{Name: "currentBlockNum", Type: mustNewType("uint64")},
		{Name: "currentBlockHash", Type: mustNewType("bytes32")},
		{Name: "validFromBlockNum", Type: mustNewType("uint64")},
	})
}

var ReportTypes = getReportTypes()

func ValidateReport(reportBlob []byte) error {
	r := map[string]interface{}{}
	err := ReportTypes.UnpackIntoMap(r, []byte(reportBlob))
	if err != nil {
		return err
	}

	feedIdInterface, ok := r["feedId"]
	if !ok {
		return errors.Errorf("unpacked report has no 'feedId'")
	}
	feedID, ok := feedIdInterface.([32]byte)
	if !ok {
		return errors.Errorf("cannot cast feedId to [32]byte, type is %T", feedID)
	}
	log.Trace().Str("FeedID", string(feedID[:])).Msg("Feed ID")

	benchmarkPriceInterface, ok := r["benchmarkPrice"]
	if !ok {
		return errors.Errorf("unpacked report has no 'benchmarkPrice'")
	}
	benchmarkPrice, ok := benchmarkPriceInterface.(*big.Int)
	if !ok {
		return errors.Errorf("cannot cast 'benchmarkPrice' to *big.Int, type is %T", benchmarkPrice)
	}
	log.Trace().Int64("benchmarkPrice", benchmarkPrice.Int64()).Msg("Benchmark price")

	bidInterface, ok := r["bid"]
	if !ok {
		return errors.Errorf("unpacked report has no 'bid'")
	}
	bidPrice, ok := bidInterface.(*big.Int)
	if !ok {
		return errors.Errorf("cannot cast 'bid' to *big.Int, type is %T", bidPrice)
	}
	log.Trace().Int64("bid", benchmarkPrice.Int64()).Msg("Bid price")

	askInterface, ok := r["ask"]
	if !ok {
		return errors.Errorf("unpacked report has no 'ask'")
	}
	askPrice, ok := askInterface.(*big.Int)
	if !ok {
		return errors.Errorf("cannot cast 'bid' to *big.Int, type is %T", askPrice)
	}
	log.Trace().Int64("ask", benchmarkPrice.Int64()).Msg("Ask price")

	currentBlockNumInterface, ok := r["currentBlockNum"]
	if !ok {
		return errors.Errorf("unpacked report has no 'currentBlockNum'")
	}
	currentBlockNumber, ok := currentBlockNumInterface.(uint64)
	if !ok {
		return errors.Errorf("cannot cast 'currentBlockNum' to uint64, type is %T", currentBlockNumber)
	}
	log.Trace().Uint64("currentBlockNumber", currentBlockNumber).Msg("Observation current block number")

	validFromBlockNumInterface, ok := r["validFromBlockNum"]
	if !ok {
		return errors.Errorf("unpacked report has no 'validFromBlockNum'")
	}
	validFromBlockNum, ok := validFromBlockNumInterface.(uint64)
	if !ok {
		return errors.Errorf("cannot cast 'validFromBlockNum' to uint64, type is %T", validFromBlockNum)
	}
	log.Trace().Uint64("validFromBlockNum", currentBlockNumber).Msg("Valid from block number")

	currentBlockHashInterface, ok := r["currentBlockHash"]
	if !ok {
		return errors.Errorf("unpacked report has no 'currentBlockHash'")
	}
	currentBlockHash, ok := currentBlockHashInterface.([32]uint8)
	if !ok {
		return errors.Errorf("cannot cast 'currentBlockHash' to uint64, type is %v", currentBlockHash)
	}
	log.Trace().Any("currentBlockHash", currentBlockHash).Msg("currentBlockHash")

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
