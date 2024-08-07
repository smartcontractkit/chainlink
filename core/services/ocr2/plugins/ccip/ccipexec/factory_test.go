package ccipexec

import (
	"errors"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccip2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdataprovidermocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/ccipdataprovider/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
)

// Assert that NewReportingPlugin keeps retrying until it succeeds.
//
// NewReportingPlugin makes several calls (e.g. OffRampReader.ChangeConfig()) that can fail. We use mocks to cause the
// first call to each of these functions to fail, then all subsequent calls succeed. We assert that NewReportingPlugin
// retries a sufficient number of times to get through the transient errors and eventually succeed.
func TestNewReportingPluginRetriesUntilSuccess(t *testing.T) {
	execConfig := ExecutionPluginStaticConfig{}
	execConfig.lggr = logger.TestLogger(t)
	execConfig.metricsCollector = ccip2.NoopMetricsCollector

	// For this unit test, ensure that there is no delay between retries
	execConfig.newReportingPluginRetryConfig = ccipdata.RetryConfig{
		InitialDelay: 0 * time.Nanosecond,
		MaxDelay:     0 * time.Nanosecond,
	}

	// Set up the OffRampReader mock
	mockOffRampReader := new(mocks.OffRampReader)

	// The first call is set to return an error, the following calls return a nil error
	mockOffRampReader.On("ChangeConfig", mock.Anything, mock.Anything, mock.Anything).Return(ccip.Address(""), ccip.Address(""), errors.New("")).Once()
	mockOffRampReader.On("ChangeConfig", mock.Anything, mock.Anything, mock.Anything).Return(ccip.Address("addr1"), ccip.Address("addr2"), nil).Times(5)

	mockOffRampReader.On("OffchainConfig", mock.Anything).Return(ccip.ExecOffchainConfig{}, errors.New("")).Once()
	mockOffRampReader.On("OffchainConfig", mock.Anything).Return(ccip.ExecOffchainConfig{}, nil).Times(3)

	mockOffRampReader.On("GasPriceEstimator", mock.Anything).Return(nil, errors.New("")).Once()
	mockOffRampReader.On("GasPriceEstimator", mock.Anything).Return(nil, nil).Times(2)

	mockOffRampReader.On("OnchainConfig", mock.Anything).Return(ccip.ExecOnchainConfig{}, errors.New("")).Once()
	mockOffRampReader.On("OnchainConfig", mock.Anything).Return(ccip.ExecOnchainConfig{}, nil).Times(1)

	execConfig.offRampReader = mockOffRampReader

	// Set up the PriceRegistry mock
	priceRegistryProvider := new(ccipdataprovidermocks.PriceRegistry)
	priceRegistryProvider.On("NewPriceRegistryReader", mock.Anything, mock.Anything).Return(nil, errors.New("")).Once()
	priceRegistryProvider.On("NewPriceRegistryReader", mock.Anything, mock.Anything).Return(nil, nil).Once()
	execConfig.priceRegistryProvider = priceRegistryProvider

	execConfig.lggr, _ = logger.NewLogger()

	factory := NewExecutionReportingPluginFactory(execConfig)
	reportingConfig := types.ReportingPluginConfig{}
	reportingConfig.OnchainConfig = []byte{1, 2, 3}
	reportingConfig.OffchainConfig = []byte{1, 2, 3}

	// Assert that NewReportingPlugin succeeds despite many transient internal failures (mocked out above)
	_, _, err := factory.NewReportingPlugin(reportingConfig)
	assert.Equal(t, nil, err)
}
