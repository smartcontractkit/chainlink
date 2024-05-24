package inprocessprovider_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayerset/inprocessprovider"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestRegisterStandAloneProvider_Median(t *testing.T) {
	s := grpc.NewServer()

	p := ocr2test.AgnosticPluginProvider
	err := inprocessprovider.RegisterStandAloneProvider(s, p, "some-type-we-do-not-support")
	require.ErrorContains(t, err, "unsupported stand alone provider")

	err = inprocessprovider.RegisterStandAloneProvider(s, p, "median")
	require.ErrorContains(t, err, "expected median provider got")

	err = inprocessprovider.RegisterStandAloneProvider(s, testMedianProvider{}, "median")
	require.NoError(t, err)
}

func TestRegisterStandAloneProvider_GenericPlugin(t *testing.T) {
	s := grpc.NewServer()

	err := inprocessprovider.RegisterStandAloneProvider(s, testPluginProvider{}, "plugin")
	require.NoError(t, err)
}

type testMedianProvider struct {
}

func (t testMedianProvider) Name() string {
	return ""
}

func (t testMedianProvider) Start(ctx context.Context) error {
	return nil
}

func (t testMedianProvider) Close() error {
	return nil
}

func (t testMedianProvider) Ready() error {
	return nil
}

func (t testMedianProvider) HealthReport() map[string]error {
	return nil
}

func (t testMedianProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return nil
}

func (t testMedianProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return nil
}

func (t testMedianProvider) ContractTransmitter() libocr.ContractTransmitter {
	return nil
}

func (t testMedianProvider) ChainReader() types.ContractReader {
	return nil
}

func (t testMedianProvider) Codec() types.Codec {
	return nil
}

func (t testMedianProvider) ReportCodec() median.ReportCodec {
	return nil
}

func (t testMedianProvider) MedianContract() median.MedianContract {
	return nil
}

func (t testMedianProvider) OnchainConfigCodec() median.OnchainConfigCodec {
	return nil
}

type testPluginProvider struct {
}

func (t testPluginProvider) Name() string {
	return ""
}

func (t testPluginProvider) Start(ctx context.Context) error {
	return nil
}

func (t testPluginProvider) Close() error {
	return nil
}

func (t testPluginProvider) Ready() error {
	return nil
}

func (t testPluginProvider) HealthReport() map[string]error {
	return nil
}

func (t testPluginProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return nil
}

func (t testPluginProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return nil
}

func (t testPluginProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return nil
}

func (t testPluginProvider) ChainReader() types.ContractReader {
	return nil
}

func (t testPluginProvider) Codec() types.Codec {
	return nil
}
