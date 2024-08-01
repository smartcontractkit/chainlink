package ocr3types

import (
	"context"
	"encoding/binary"
	"math/rand"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type CountingMercuryPlugin struct {
	logger            commontypes.Logger
	initializedReport bool
}

func (p *CountingMercuryPlugin) Observation(ctx context.Context, repts types.ReportTimestamp, previousReport types.Report) (types.Observation, error) {
	return []byte{byte(rand.Int() % 2)}, nil
}

func (p *CountingMercuryPlugin) Report(repts types.ReportTimestamp, previousReport types.Report, aos []types.AttributedObservation) (bool, types.Report, error) {
	report := make([]byte, 4)
	if len(previousReport) == 0 {
		if p.initializedReport {
			panic("this should not happen")
		}
		return true, report, nil
	}

	if len(previousReport) != 0 {
		p.initializedReport = true
	}

	shouldReport := false
	for _, ao := range aos {
		if len(ao.Observation) != 1 {
			p.logger.Warn("invalid ao", nil)
			continue
		}

		if ao.Observation[0] > 0 {
			shouldReport = !shouldReport
		}
	}

	count := binary.BigEndian.Uint32(previousReport)
	if shouldReport {
		count++
	}
	binary.BigEndian.PutUint32(report, count)

	return shouldReport, report, nil
}

func (p *CountingMercuryPlugin) Close() error {
	return nil
}

type CountingMercuryPluginFactory struct{}

func (fac *CountingMercuryPluginFactory) NewMercuryPlugin(_ MercuryPluginConfig) (MercuryPlugin, MercuryPluginInfo, error) {
	return &CountingMercuryPlugin{},
		MercuryPluginInfo{
			"CountingMercuryPlugin", MercuryPluginLimits{
				1,
				4,
			},
		},
		nil
}
