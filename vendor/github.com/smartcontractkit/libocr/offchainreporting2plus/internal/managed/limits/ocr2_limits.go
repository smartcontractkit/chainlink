package limits

import (
	"crypto/ed25519"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr2config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func OCR2Limits(cfg ocr2config.PublicConfig, reportingPluginLimits types.ReportingPluginLimits, maxSigLen int) (types.BinaryNetworkEndpointLimits, error) {
	overflow := false

	// These two helper functions add/multiply together a bunch of numbers and set overflow to true if the result
	// lies outside the range [0; math.MaxInt32]. We compare with int32 rather than int to be independent of
	// the underlying architecture.
	add := func(xs ...int) int {
		sum := big.NewInt(0)
		for _, x := range xs {
			sum.Add(sum, big.NewInt(int64(x)))
		}
		if !(big.NewInt(0).Cmp(sum) <= 0 && sum.Cmp(big.NewInt(int64(math.MaxInt32))) <= 0) {
			overflow = true
		}
		return int(sum.Int64())
	}
	mul := func(xs ...int) int {
		prod := big.NewInt(1)
		for _, x := range xs {
			prod.Mul(prod, big.NewInt(int64(x)))
		}
		if !(big.NewInt(0).Cmp(prod) <= 0 && prod.Cmp(big.NewInt(int64(math.MaxInt32))) <= 0) {
			overflow = true
		}
		return int(prod.Int64())
	}

	const overhead = 256

	maxLenNewEpoch := overhead
	maxLenObserveReq := add(reportingPluginLimits.MaxQueryLength, overhead)
	maxLenObserve := add(reportingPluginLimits.MaxObservationLength, overhead)
	maxLenReportReq := add(mul(add(reportingPluginLimits.MaxObservationLength, ed25519.SignatureSize), cfg.N()), overhead)
	maxLenReport := add(reportingPluginLimits.MaxReportLength, ed25519.SignatureSize, overhead)
	maxLenFinal := add(reportingPluginLimits.MaxReportLength, mul(maxSigLen, cfg.N()), overhead)
	maxLenFinalEcho := maxLenFinal

	maxMessageSize := max(maxLenObserveReq, maxLenObserve, maxLenReportReq, maxLenReport, maxLenFinal, maxLenFinalEcho)

	messagesRate := (1.0*float64(time.Second)/float64(cfg.DeltaResend) +
		1.0*float64(time.Second)/float64(cfg.DeltaProgress) +
		1.0*float64(time.Second)/float64(cfg.DeltaRound) +
		3.0*float64(time.Second)/float64(cfg.DeltaRound) +
		2.0*float64(time.Second)/float64(cfg.DeltaRound)) * 2.0

	messagesCapacity := mul(add(2, 6), 2)

	bytesRate := float64(time.Second)/float64(cfg.DeltaResend)*float64(maxLenNewEpoch) +
		float64(time.Second)/float64(cfg.DeltaProgress)*float64(maxLenNewEpoch) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenObserveReq) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenObserve) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenReportReq) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenReport) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenFinal) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenFinalEcho)

	bytesCapacity := mul(add(maxLenNewEpoch, maxLenObserveReq, maxLenObserve, maxLenReportReq, maxLenReport, maxLenFinal, maxLenFinalEcho), 2)

	if overflow {
		// this should not happen due to us checking the limits in types.go
		return types.BinaryNetworkEndpointLimits{}, fmt.Errorf("int32 overflow while computing bandwidth limits")
	}

	return types.BinaryNetworkEndpointLimits{
		maxMessageSize,
		messagesRate,
		messagesCapacity,
		bytesRate,
		bytesCapacity,
	}, nil
}
