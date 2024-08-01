package v1

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	mercurytypes "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"

	"github.com/smartcontractkit/chainlink-data-streams/mercury"
)

// MaxAllowedBlocks indicates the maximum len of LatestBlocks in any given observation.
// observations that violate this will be discarded
const MaxAllowedBlocks = 10

// DataSource implementations must be thread-safe. Observe may be called by many
// different threads concurrently.
type DataSource interface {
	// Observe queries the data source. Returns a value or an error. Once the
	// context is expires, Observe may still do cheap computations and return a
	// result, but should return as quickly as possible.
	//
	// More details: In the current implementation, the context passed to
	// Observe will time out after MaxDurationObservation. However, Observe
	// should *not* make any assumptions about context timeout behavior. Once
	// the context times out, Observe should prioritize returning as quickly as
	// possible, but may still perform fast computations to return a result
	// rather than error. For example, if Observe medianizes a number of data
	// sources, some of which already returned a result to Observe prior to the
	// context's expiry, Observe might still compute their median, and return it
	// instead of an error.
	//
	// Important: Observe should not perform any potentially time-consuming
	// actions like database access, once the context passed has expired.
	Observe(ctx context.Context, repts types.ReportTimestamp, fetchMaxFinalizedBlockNum bool) (v1.Observation, error)
}

var _ ocr3types.MercuryPluginFactory = Factory{}

// Maximum length in bytes of Observation, Report returned by the
// MercuryPlugin. Used for defending against spam attacks.
const maxObservationLength = 4 + // timestamp
	mercury.ByteWidthInt192 + // benchmarkPrice
	mercury.ByteWidthInt192 + // bid
	mercury.ByteWidthInt192 + // ask
	1 + // pricesValid
	8 + // currentBlockNum
	32 + // currentBlockHash
	8 + // currentBlockTimestamp
	1 + // currentBlockValid
	8 + // maxFinalizedBlockNumber
	1 + // maxFinalizedBlockNumberValid
	32 + // [> overapprox. of protobuf overhead <]
	MaxAllowedBlocks*(8+ // num
		32+ // hash
		8+ // ts
		32) // [> overapprox. of protobuf overhead <]

type Factory struct {
	dataSource         DataSource
	logger             logger.Logger
	onchainConfigCodec mercurytypes.OnchainConfigCodec
	reportCodec        v1.ReportCodec
}

func NewFactory(ds DataSource, lggr logger.Logger, occ mercurytypes.OnchainConfigCodec, rc v1.ReportCodec) Factory {
	return Factory{ds, lggr, occ, rc}
}

func (fac Factory) NewMercuryPlugin(configuration ocr3types.MercuryPluginConfig) (ocr3types.MercuryPlugin, ocr3types.MercuryPluginInfo, error) {
	offchainConfig, err := mercury.DecodeOffchainConfig(configuration.OffchainConfig)
	if err != nil {
		return nil, ocr3types.MercuryPluginInfo{}, err
	}

	onchainConfig, err := fac.onchainConfigCodec.Decode(configuration.OnchainConfig)
	if err != nil {
		return nil, ocr3types.MercuryPluginInfo{}, err
	}

	maxReportLength, err := fac.reportCodec.MaxReportLength(configuration.N)
	if err != nil {
		return nil, ocr3types.MercuryPluginInfo{}, err
	}

	r := &reportingPlugin{
		offchainConfig,
		onchainConfig,
		fac.dataSource,
		fac.logger,
		fac.reportCodec,
		configuration.ConfigDigest,
		configuration.F,
		mercury.EpochRound{},
		maxReportLength,
	}

	return r, ocr3types.MercuryPluginInfo{
		Name: "Mercury",
		Limits: ocr3types.MercuryPluginLimits{
			MaxObservationLength: maxObservationLength,
			MaxReportLength:      maxReportLength,
		},
	}, nil
}

var _ ocr3types.MercuryPlugin = (*reportingPlugin)(nil)

type reportingPlugin struct {
	offchainConfig mercury.OffchainConfig
	onchainConfig  mercurytypes.OnchainConfig
	dataSource     DataSource
	logger         logger.Logger
	reportCodec    v1.ReportCodec

	configDigest             types.ConfigDigest
	f                        int
	latestAcceptedEpochRound mercury.EpochRound
	maxReportLength          int
}

func (rp *reportingPlugin) Observation(ctx context.Context, repts types.ReportTimestamp, previousReport types.Report) (types.Observation, error) {
	obs, err := rp.dataSource.Observe(ctx, repts, previousReport == nil)
	if err != nil {
		return nil, fmt.Errorf("DataSource.Observe returned an error: %s", err)
	}

	p := MercuryObservationProto{Timestamp: uint32(time.Now().Unix())}

	var obsErrors []error
	if previousReport == nil {
		// if previousReport we fall back to the observed MaxFinalizedBlockNumber
		if obs.MaxFinalizedBlockNumber.Err != nil {
			obsErrors = append(obsErrors, err)
		} else if obs.CurrentBlockNum.Err == nil && obs.CurrentBlockNum.Val < obs.MaxFinalizedBlockNumber.Val {
			obsErrors = append(obsErrors, fmt.Errorf("failed to observe ValidFromBlockNum; current block number %d (hash: 0x%x) < max finalized block number %d; ignoring observation for out-of-date RPC", obs.CurrentBlockNum.Val, obs.CurrentBlockHash.Val, obs.MaxFinalizedBlockNumber.Val))
		} else {
			p.MaxFinalizedBlockNumber = obs.MaxFinalizedBlockNumber.Val // MaxFinalizedBlockNumber comes as -1 if unset
			p.MaxFinalizedBlockNumberValid = true
		}
	}

	var bpErr, bidErr, askErr error
	if obs.BenchmarkPrice.Err != nil {
		bpErr = fmt.Errorf("failed to observe BenchmarkPrice: %w", obs.BenchmarkPrice.Err)
		obsErrors = append(obsErrors, bpErr)
	} else if benchmarkPrice, err := mercury.EncodeValueInt192(obs.BenchmarkPrice.Val); err != nil {
		bpErr = fmt.Errorf("failed to observe BenchmarkPrice; encoding failed: %w", err)
		obsErrors = append(obsErrors, bpErr)
	} else {
		p.BenchmarkPrice = benchmarkPrice
	}

	if obs.Bid.Err != nil {
		bidErr = fmt.Errorf("failed to observe Bid: %w", obs.Bid.Err)
		obsErrors = append(obsErrors, bidErr)
	} else if bid, err := mercury.EncodeValueInt192(obs.Bid.Val); err != nil {
		bidErr = fmt.Errorf("failed to observe Bid; encoding failed: %w", err)
		obsErrors = append(obsErrors, bidErr)
	} else {
		p.Bid = bid
	}

	if obs.Ask.Err != nil {
		askErr = fmt.Errorf("failed to observe Ask: %w", obs.Ask.Err)
		obsErrors = append(obsErrors, askErr)
	} else if ask, err := mercury.EncodeValueInt192(obs.Ask.Val); err != nil {
		askErr = fmt.Errorf("failed to observe Ask; encoding failed: %w", err)
		obsErrors = append(obsErrors, askErr)
	} else {
		p.Ask = ask
	}

	if bpErr == nil && bidErr == nil && askErr == nil {
		p.PricesValid = true
	}

	if obs.CurrentBlockNum.Err != nil {
		obsErrors = append(obsErrors, fmt.Errorf("failed to observe CurrentBlockNum: %w", obs.CurrentBlockNum.Err))
	} else {
		p.CurrentBlockNum = obs.CurrentBlockNum.Val
	}

	if obs.CurrentBlockHash.Err != nil {
		obsErrors = append(obsErrors, fmt.Errorf("failed to observe CurrentBlockHash: %w", obs.CurrentBlockHash.Err))
	} else {
		p.CurrentBlockHash = obs.CurrentBlockHash.Val
	}

	if obs.CurrentBlockTimestamp.Err != nil {
		obsErrors = append(obsErrors, fmt.Errorf("failed to observe CurrentBlockTimestamp: %w", obs.CurrentBlockTimestamp.Err))
	} else {
		p.CurrentBlockTimestamp = obs.CurrentBlockTimestamp.Val
	}

	if obs.CurrentBlockNum.Err == nil && obs.CurrentBlockHash.Err == nil && obs.CurrentBlockTimestamp.Err == nil {
		p.CurrentBlockValid = true
	}

	if len(obsErrors) > 0 {
		rp.logger.Warnw(fmt.Sprintf("Observe failed %d/6 observations", len(obsErrors)), "err", errors.Join(obsErrors...))
	}

	p.LatestBlocks = make([]*BlockProto, len(obs.LatestBlocks))
	for i, b := range obs.LatestBlocks {
		p.LatestBlocks[i] = &BlockProto{Num: b.Num, Hash: []byte(b.Hash), Ts: b.Ts}
	}
	if len(p.LatestBlocks) == 0 {
		rp.logger.Warn("Observation had no LatestBlocks")
	}

	return proto.Marshal(&p)
}

func parseAttributedObservation(ao types.AttributedObservation) (PAO, error) {
	var pao parsedAttributedObservation
	var obs MercuryObservationProto
	if err := proto.Unmarshal(ao.Observation, &obs); err != nil {
		return parsedAttributedObservation{}, fmt.Errorf("attributed observation cannot be unmarshaled: %s", err)
	}

	pao.Timestamp = obs.Timestamp
	pao.Observer = ao.Observer

	if obs.PricesValid {
		var err error
		pao.BenchmarkPrice, err = mercury.DecodeValueInt192(obs.BenchmarkPrice)
		if err != nil {
			return parsedAttributedObservation{}, fmt.Errorf("benchmarkPrice cannot be converted to big.Int: %s", err)
		}
		pao.Bid, err = mercury.DecodeValueInt192(obs.Bid)
		if err != nil {
			return parsedAttributedObservation{}, fmt.Errorf("bid cannot be converted to big.Int: %s", err)
		}
		pao.Ask, err = mercury.DecodeValueInt192(obs.Ask)
		if err != nil {
			return parsedAttributedObservation{}, fmt.Errorf("ask cannot be converted to big.Int: %s", err)
		}
		pao.PricesValid = true
	}

	if len(obs.LatestBlocks) > 0 {
		if len(obs.LatestBlocks) > MaxAllowedBlocks {
			return parsedAttributedObservation{}, fmt.Errorf("LatestBlocks too large; got: %d, max: %d", len(obs.LatestBlocks), MaxAllowedBlocks)
		}
		for _, b := range obs.LatestBlocks {
			pao.LatestBlocks = append(pao.LatestBlocks, v1.NewBlock(b.Num, b.Hash, b.Ts))

			// Ignore observation if it has duplicate blocks by number or hash
			// for security to avoid the case where one node can "throw" block
			// numbers by including a bunch of duplicates
			nums := make(map[int64]struct{}, len(pao.LatestBlocks))
			hashes := make(map[string]struct{}, len(pao.LatestBlocks))
			for _, block := range pao.LatestBlocks {
				if _, exists := nums[block.Num]; exists {
					return parsedAttributedObservation{}, fmt.Errorf("observation invalid for observer %d; got duplicate block number: %d", ao.Observer, block.Num)
				}
				if _, exists := hashes[block.Hash]; exists {
					return parsedAttributedObservation{}, fmt.Errorf("observation invalid for observer %d; got duplicate block hash: 0x%x", ao.Observer, block.HashBytes())
				}
				nums[block.Num] = struct{}{}
				hashes[block.Hash] = struct{}{}

				if len(block.Hash) != mercury.EvmHashLen {
					return parsedAttributedObservation{}, fmt.Errorf("wrong len for hash: %d (expected: %d)", len(block.Hash), mercury.EvmHashLen)
				}
				if block.Num < 0 {
					return parsedAttributedObservation{}, fmt.Errorf("negative block number: %d", block.Num)
				}
			}

			// sort desc
			sort.SliceStable(pao.LatestBlocks, func(i, j int) bool {
				// NOTE: This ought to be redundant since observing nodes
				// should give us the blocks pre-sorted, but is included here
				// for safety
				return pao.LatestBlocks[j].Less(pao.LatestBlocks[i])
			})
		}
	} else if obs.CurrentBlockValid {
		// DEPRECATED
		// TODO: Remove this handling after deployment (https://smartcontract-it.atlassian.net/browse/MERC-2272)
		if len(obs.CurrentBlockHash) != mercury.EvmHashLen {
			return parsedAttributedObservation{}, fmt.Errorf("wrong len for hash: %d (expected: %d)", len(obs.CurrentBlockHash), mercury.EvmHashLen)
		}
		pao.CurrentBlockHash = obs.CurrentBlockHash
		if obs.CurrentBlockNum < 0 {
			return parsedAttributedObservation{}, fmt.Errorf("negative block number: %d", obs.CurrentBlockNum)
		}
		pao.CurrentBlockNum = obs.CurrentBlockNum
		pao.CurrentBlockTimestamp = obs.CurrentBlockTimestamp
		pao.CurrentBlockValid = true
	}

	if obs.MaxFinalizedBlockNumberValid {
		pao.MaxFinalizedBlockNumber = obs.MaxFinalizedBlockNumber
		pao.MaxFinalizedBlockNumberValid = true
	}

	return pao, nil
}

func parseAttributedObservations(lggr logger.Logger, aos []types.AttributedObservation) []PAO {
	paos := make([]PAO, 0, len(aos))
	for i, ao := range aos {
		pao, err := parseAttributedObservation(ao)
		if err != nil {
			lggr.Warnw("parseAttributedObservations: dropping invalid observation",
				"observer", ao.Observer,
				"error", err,
				"i", i,
			)
			continue
		}
		paos = append(paos, pao)
	}
	return paos
}

func (rp *reportingPlugin) Report(repts types.ReportTimestamp, previousReport types.Report, aos []types.AttributedObservation) (shouldReport bool, report types.Report, err error) {
	paos := parseAttributedObservations(rp.logger, aos)

	if len(paos) == 0 {
		return false, nil, errors.New("got zero valid attributed observations")
	}

	// By assumption, we have at most f malicious oracles, so there should be at least f+1 valid paos
	if !(rp.f+1 <= len(paos)) {
		return false, nil, fmt.Errorf("only received %v valid attributed observations, but need at least f+1 (%v)", len(paos), rp.f+1)
	}

	rf, err := rp.buildReportFields(previousReport, paos)
	if err != nil {
		rp.logger.Errorw("failed to build report fields", "paos", paos, "f", rp.f, "reportFields", rf, "repts", repts, "err", err)
		return false, nil, err
	}

	if rf.CurrentBlockNum < rf.ValidFromBlockNum {
		rp.logger.Debugw("shouldReport: no (overlap)", "currentBlockNum", rf.CurrentBlockNum, "validFromBlockNum", rf.ValidFromBlockNum, "repts", repts)
		return false, nil, nil
	}

	if err = rp.validateReport(rf); err != nil {
		rp.logger.Errorw("shouldReport: no (validation error)", "reportFields", rf, "err", err, "repts", repts, "paos", paos)
		return false, nil, err
	}
	rp.logger.Debugw("shouldReport: yes",
		"timestamp", repts,
	)

	report, err = rp.reportCodec.BuildReport(rf)
	if err != nil {
		rp.logger.Debugw("failed to BuildReport", "paos", paos, "f", rp.f, "reportFields", rf, "repts", repts)
		return false, nil, err
	}
	if !(len(report) <= rp.maxReportLength) {
		return false, nil, fmt.Errorf("report with len %d violates MaxReportLength limit set by ReportCodec (%d)", len(report), rp.maxReportLength)
	} else if len(report) == 0 {
		return false, nil, errors.New("report may not have zero length (invariant violation)")
	}

	return true, report, nil
}

func (rp *reportingPlugin) buildReportFields(previousReport types.Report, paos []PAO) (rf v1.ReportFields, merr error) {
	var err error
	if previousReport != nil {
		var maxFinalizedBlockNumber int64
		maxFinalizedBlockNumber, err = rp.reportCodec.CurrentBlockNumFromReport(previousReport)
		if err != nil {
			merr = errors.Join(merr, err)
		} else {
			rf.ValidFromBlockNum = maxFinalizedBlockNumber + 1
		}
	} else {
		var maxFinalizedBlockNumber int64
		maxFinalizedBlockNumber, err = GetConsensusMaxFinalizedBlockNum(paos, rp.f)
		if err != nil {
			merr = errors.Join(merr, err)
		} else {
			rf.ValidFromBlockNum = maxFinalizedBlockNumber + 1
		}
	}

	mPaos := convert(paos)

	rf.Timestamp = mercury.GetConsensusTimestamp(mPaos)

	rf.BenchmarkPrice, err = mercury.GetConsensusBenchmarkPrice(mPaos, rp.f)
	if err != nil {
		merr = errors.Join(merr, fmt.Errorf("GetConsensusBenchmarkPrice failed: %w", err))
	}

	rf.Bid, err = mercury.GetConsensusBid(convertBid(paos), rp.f)
	if err != nil {
		merr = errors.Join(merr, fmt.Errorf("GetConsensusBid failed: %w", err))
	}

	rf.Ask, err = mercury.GetConsensusAsk(convertAsk(paos), rp.f)
	if err != nil {
		merr = errors.Join(merr, fmt.Errorf("GetConsensusAsk failed: %w", err))
	}

	rf.CurrentBlockHash, rf.CurrentBlockNum, rf.CurrentBlockTimestamp, err = GetConsensusLatestBlock(paos, rp.f)
	if err != nil {
		merr = errors.Join(merr, fmt.Errorf("GetConsensusCurrentBlock failed: %w", err))
	}

	return rf, merr
}

func (rp *reportingPlugin) validateReport(rf v1.ReportFields) error {
	return errors.Join(
		mercury.ValidateBetween("median benchmark price", rf.BenchmarkPrice, rp.onchainConfig.Min, rp.onchainConfig.Max),
		mercury.ValidateBetween("median bid", rf.Bid, rp.onchainConfig.Min, rp.onchainConfig.Max),
		mercury.ValidateBetween("median ask", rf.Ask, rp.onchainConfig.Min, rp.onchainConfig.Max),
		ValidateCurrentBlock(rf),
	)
}

func (rp *reportingPlugin) Close() error {
	return nil
}

// convert funcs are necessary because go is not smart enough to cast
// []interface1 to []interface2 even if interface1 is a superset of interface2
func convert(pao []PAO) (ret []mercury.PAO) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertBid(pao []PAO) (ret []mercury.PAOBid) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertAsk(pao []PAO) (ret []mercury.PAOAsk) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
