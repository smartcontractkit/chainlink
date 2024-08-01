package median

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

const contractName = "median"

type Plugin struct {
	loop.Plugin
	stop services.StopChan
}

func NewPlugin(lggr logger.Logger) *Plugin {
	return &Plugin{Plugin: loop.Plugin{Logger: lggr}, stop: make(services.StopChan)}
}

func (p *Plugin) NewMedianFactory(ctx context.Context, provider types.MedianProvider, dataSource, juelsPerFeeCoin, gasPriceSubunits median.DataSource, errorLog loop.ErrorLog) (loop.ReportingPluginFactory, error) {
	var ctxVals loop.ContextValues
	ctxVals.SetValues(ctx)
	lggr := logger.With(p.Logger, ctxVals.Args()...)

	// We omit gas price in observation to maintain backwards compatibility in libocr (with older nodes).
	// Once all chainlink nodes have updated to libocr version >= fd3cab206b2c
	// the IncludeGasPriceSubunitsInObservation field can be removed

	_, isZeroDataSource := gasPriceSubunits.(*ZeroDataSource)

	includeGasPriceSubunitsInObservation := !isZeroDataSource

	factory := median.NumericalMedianFactory{
		DataSource:                           dataSource,
		JuelsPerFeeCoinDataSource:            juelsPerFeeCoin,
		GasPriceSubunitsDataSource:           gasPriceSubunits,
		IncludeGasPriceSubunitsInObservation: includeGasPriceSubunitsInObservation,
		Logger: logger.NewOCRWrapper(lggr, true, func(msg string) {
			ctx, cancelFn := p.stop.NewCtx()
			defer cancelFn()
			if err := errorLog.SaveError(ctx, msg); err != nil {
				lggr.Errorw("Unable to save error", "err", msg)
			}
		}),
		OnchainConfigCodec: provider.OnchainConfigCodec(),
	}

	if cr := provider.ChainReader(); cr != nil {
		factory.ContractTransmitter = &contractReaderContract{contractReader: cr, lggr: lggr}
	} else {
		factory.ContractTransmitter = provider.MedianContract()
	}

	if codec := provider.Codec(); codec != nil {
		factory.ReportCodec = &reportCodec{codec: codec}
	} else {
		lggr.Info("No codec provided, defaulting back to median specific ReportCodec")
		factory.ReportCodec = provider.ReportCodec()
	}

	s := &reportingPluginFactoryService{lggr: logger.Named(lggr, "ReportingPluginFactory"), ReportingPluginFactory: factory}

	p.SubService(s)

	return s, nil
}

type ZeroDataSource struct{}

func (d *ZeroDataSource) Observe(ctx context.Context, reportTimestamp ocrtypes.ReportTimestamp) (*big.Int, error) {
	return new(big.Int), nil
}

type reportingPluginFactoryService struct {
	services.StateMachine
	lggr logger.Logger
	ocrtypes.ReportingPluginFactory
}

func (r *reportingPluginFactoryService) Name() string { return r.lggr.Name() }

func (r *reportingPluginFactoryService) Start(ctx context.Context) error {
	return r.StartOnce("ReportingPluginFactory", func() error { return nil })
}

func (r *reportingPluginFactoryService) Close() error {
	return r.StopOnce("ReportingPluginFactory", func() error { return nil })
}

func (r *reportingPluginFactoryService) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

// contractReaderContract adapts a [types.ContractReader] to [median.MedianContract].
type contractReaderContract struct {
	contractReader types.ContractReader
	lggr           logger.Logger
}

type latestTransmissionDetailsResponse struct {
	ConfigDigest    ocrtypes.ConfigDigest
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp time.Time
}

type latestRoundRequested struct {
	ConfigDigest ocrtypes.ConfigDigest
	Epoch        uint32
	Round        uint8
}

func (c *contractReaderContract) LatestTransmissionDetails(ctx context.Context) (configDigest ocrtypes.ConfigDigest, epoch uint32, round uint8, latestAnswer *big.Int, latestTimestamp time.Time, err error) {
	var resp latestTransmissionDetailsResponse

	err = c.contractReader.GetLatestValue(ctx, contractName, "LatestTransmissionDetails", primitives.Unconfirmed, nil, &resp)
	if err != nil {
		if !errors.Is(err, types.ErrNotFound) {
			return
		}
		// If there's nothing transmitted yet, an implementation will not have emitted an event,
		// or may not find details of a latest transmission on-chain if it's a function call.
		// A zeroed out latestTransmissionDetailsResponse tells later parts of the system that there's no data yet.
		c.lggr.Warn("LatestTransmissionDetails not found", "err", err)
	}

	// Depending on if there is a LatestAnswer or not, and the implementation of the ContractReader,
	// it's possible that this will be unset. The desired behaviour in that case is to have a zero value.
	if resp.LatestAnswer == nil {
		resp.LatestAnswer = new(big.Int)
	}

	return resp.ConfigDigest, resp.Epoch, resp.Round, resp.LatestAnswer, resp.LatestTimestamp, nil
}

func (c *contractReaderContract) LatestRoundRequested(ctx context.Context, lookback time.Duration) (configDigest ocrtypes.ConfigDigest, epoch uint32, round uint8, err error) {
	var resp latestRoundRequested

	err = c.contractReader.GetLatestValue(ctx, contractName, "LatestRoundRequested", primitives.Unconfirmed, nil, &resp)
	if err != nil {
		if !errors.Is(err, types.ErrNotFound) {
			return
		}
		// If there's nothing on-chain yet, an implementation will not have emitted an event,
		// or may not find details of a latest transmission on-chain if it's a function call.
		// A zeroed out LatestRoundRequested tells later parts of the system that there's no data yet.
		c.lggr.Warn("LatestRoundRequested not found", "err", err)
	}

	return resp.ConfigDigest, resp.Epoch, resp.Round, nil
}
