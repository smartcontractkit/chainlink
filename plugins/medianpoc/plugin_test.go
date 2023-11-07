package medianpoc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type mockErrorLog struct {
	types.ErrorLog
}

type mockOffchainConfigDigester struct {
	ocrtypes.OffchainConfigDigester
}

type mockContractTransmitter struct {
	ocrtypes.ContractTransmitter
}

type mockContractConfigTracker struct {
	ocrtypes.ContractConfigTracker
}

type mockReportCodec struct {
	median.ReportCodec
}

type mockMedianContract struct {
	median.MedianContract
}

type mockOnchainConfigCodec struct {
	median.OnchainConfigCodec
}

type provider struct {
	types.Service
}

func (p provider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return mockOffchainConfigDigester{}
}

func (p provider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return mockContractTransmitter{}
}

func (p provider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return mockContractConfigTracker{}
}

func (p provider) ReportCodec() median.ReportCodec {
	return mockReportCodec{}
}

func (p provider) MedianContract() median.MedianContract {
	return mockMedianContract{}
}

func (p provider) OnchainConfigCodec() median.OnchainConfigCodec {
	return mockOnchainConfigCodec{}
}

func TestNewPlugin(t *testing.T) {
	lggr := logger.TestLogger(t)
	p := NewPlugin(lggr)

	defaultSpec := "default-spec"
	juelsPerFeeCoinSpec := "jpfc-spec"
	config := types.ReportingPluginServiceConfig{
		PluginConfig: fmt.Sprintf(
			`{"pipelines": {"__DEFAULT_PIPELINE__": "%s", "juelsPerFeeCoinPipeline": "%s"}}`,
			defaultSpec,
			juelsPerFeeCoinSpec,
		),
	}
	pr := &mockPipelineRunner{}
	prov := provider{}

	f, err := p.newFactory(
		context.Background(),
		config,
		prov,
		pr,
		nil,
		mockErrorLog{},
	)
	require.NoError(t, err)

	ds := f.DataSource.(*DataSource)
	assert.Equal(t, defaultSpec, ds.spec)
	jpfcDs := f.JuelsPerFeeCoinDataSource.(*DataSource)
	assert.Equal(t, juelsPerFeeCoinSpec, jpfcDs.spec)
}
