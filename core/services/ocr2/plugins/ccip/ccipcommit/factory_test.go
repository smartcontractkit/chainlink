package ccipcommit

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccip2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdataprovidermocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/ccipdataprovider/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	dbMocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdb/mocks"
)

// Assert that NewReportingPlugin keeps retrying until it succeeds.
//
// NewReportingPlugin makes several calls (e.g. CommitStoreReader.ChangeConfig) that can fail. We use mocks to cause the
// first call to each of these functions to fail, then all subsequent calls succeed. We assert that NewReportingPlugin
// retries a sufficient number of times to get through the transient errors and eventually succeed.
func TestNewReportingPluginRetriesUntilSuccess(t *testing.T) {
	ctx := tests.Context(t)
	commitConfig := CommitPluginStaticConfig{}
	commitConfig.lggr = logger.TestLogger(t)
	commitConfig.metricsCollector = ccip2.NoopMetricsCollector

	// For this unit test, ensure that there is no delay between retries
	commitConfig.newReportingPluginRetryConfig = ccipdata.RetryConfig{
		InitialDelay: 0 * time.Nanosecond,
		MaxDelay:     0 * time.Nanosecond,
	}

	// Set up the OffRampReader mock
	mockCommitStore := new(mocks.CommitStoreReader)

	// The first call is set to return an error, the following calls return a nil error
	mockCommitStore.
		On("ChangeConfig", mock.Anything, mock.Anything, mock.Anything).
		Return(ccip.Address(""), errors.New("")).
		Once()
	mockCommitStore.
		On("ChangeConfig", mock.Anything, mock.Anything, mock.Anything).
		Return(ccip.Address("0x7c6e4F0BDe29f83BC394B75a7f313B7E5DbD2d77"), nil).
		Times(5)

	mockCommitStore.
		On("OffchainConfig", mock.Anything).
		Return(ccip.CommitOffchainConfig{}, errors.New("")).
		Once()
	mockCommitStore.
		On("OffchainConfig", mock.Anything).
		Return(ccip.CommitOffchainConfig{}, nil).
		Times(3)

	mockCommitStore.
		On("GasPriceEstimator", mock.Anything).
		Return(nil, errors.New("")).
		Once()
	mockCommitStore.
		On("GasPriceEstimator", mock.Anything).
		Return(nil, nil).
		Times(2)

	commitConfig.commitStore = mockCommitStore

	mockPriceService := new(dbMocks.PriceService)

	mockPriceService.
		On("UpdateDynamicConfig", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("")).
		Once()
	mockPriceService.
		On("UpdateDynamicConfig", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	commitConfig.priceService = mockPriceService

	priceRegistryProvider := new(ccipdataprovidermocks.PriceRegistry)
	priceRegistryProvider.
		On("NewPriceRegistryReader", mock.Anything, mock.Anything).
		Return(nil, errors.New("")).
		Once()
	priceRegistryProvider.
		On("NewPriceRegistryReader", mock.Anything, mock.Anything).
		Return(nil, nil).
		Once()
	commitConfig.priceRegistryProvider = priceRegistryProvider

	commitConfig.lggr, _ = logger.NewLogger()

	factory := NewCommitReportingPluginFactory(commitConfig)
	reportingConfig := types.ReportingPluginConfig{}
	reportingConfig.OnchainConfig = []byte{1, 2, 3}
	reportingConfig.OffchainConfig = []byte{1, 2, 3}

	// Assert that NewReportingPlugin succeeds despite many transient internal failures (mocked out above)
	_, _, err := factory.NewReportingPlugin(ctx, reportingConfig)
	assert.NoError(t, err)
}
