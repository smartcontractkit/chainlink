package medianpoc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type mockErrorLog struct {
	core.ErrorLog
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

func (p provider) ChainReader() types.ChainReader {
	return nil
}

func (p provider) Codec() types.Codec {
	return nil
}

func TestNewPlugin(t *testing.T) {
	lggr := logger.TestLogger(t)
	p := NewPlugin(lggr)

	defaultSpec := "default-spec"
	juelsPerFeeCoinSpec := "jpfc-spec"
	config := core.ReportingPluginServiceConfig{
		PluginConfig: fmt.Sprintf(
			`{"pipelines": [{"name": "__DEFAULT_PIPELINE__", "spec": "%s"},{"name": "juelsPerFeeCoinPipeline", "spec": "%s"}]}`,
			defaultSpec,
			juelsPerFeeCoinSpec,
		),
	}
	pr := &mockPipelineRunner{}
	prov := provider{}

	f, err := p.newFactory(
		tests.Context(t),
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
