package usdc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	mockOnRampAddress    = utils.RandomAddress()
	mockUSDCTokenAddress = utils.RandomAddress()
	mockMsgTransmitter   = utils.RandomAddress()
)

func TestUSDCReader_callAttestationApi(t *testing.T) {
	t.Skipf("Skipping test because it uses the real USDC attestation API")
	usdcMessageHash := "912f22a13e9ccb979b621500f6952b2afd6e75be7eadaed93fc2625fe11c52a2"
	attestationURI, err := url.ParseRequestURI("https://iris-api-sandbox.circle.com")
	require.NoError(t, err)
	usdcService := NewUSDCTokenDataReader(logger.TestLogger(t), nil, mockUSDCTokenAddress, mockMsgTransmitter, mockOnRampAddress, attestationURI)

	attestation, err := usdcService.callAttestationApi(context.Background(), [32]byte(common.FromHex(usdcMessageHash)))
	require.NoError(t, err)

	require.Equal(t, attestationStatusPending, attestation.Status)
	require.Equal(t, "PENDING", attestation.Attestation)
}

func TestUSDCReader_callAttestationApiMock(t *testing.T) {
	response := attestationResponse{
		Status:      attestationStatusSuccess,
		Attestation: "720502893578a89a8a87982982ef781c18b193",
	}

	ts := getMockUSDCEndpoint(t, response)
	defer ts.Close()
	attestationURI, err := url.ParseRequestURI(ts.URL)
	require.NoError(t, err)

	usdcService := NewUSDCTokenDataReader(logger.TestLogger(t), nil, mockUSDCTokenAddress, mockMsgTransmitter, mockOnRampAddress, attestationURI)
	attestation, err := usdcService.callAttestationApi(context.Background(), utils.RandomBytes32())
	require.NoError(t, err)

	require.Equal(t, response.Status, attestation.Status)
	require.Equal(t, response.Attestation, attestation.Attestation)
}

func TestUSDCReader_callAttestationApiMockError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()
	attestationURI, err := url.ParseRequestURI(ts.URL)
	require.NoError(t, err)

	usdcService := NewUSDCTokenDataReader(logger.TestLogger(t), nil, mockUSDCTokenAddress, mockMsgTransmitter, mockOnRampAddress, attestationURI)
	_, err = usdcService.callAttestationApi(context.Background(), utils.RandomBytes32())
	require.Error(t, err)
}

func getMockUSDCEndpoint(t *testing.T, response attestationResponse) *httptest.Server {
	responseBytes, err := json.Marshal(response)
	require.NoError(t, err)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(responseBytes)
		require.NoError(t, err)
	}))
}

// Asserts the hard coded event signature matches Keccak256("MessageSent(bytes)")
func TestGetUSDCReaderSourceLPFilters(t *testing.T) {
	usdcService := NewUSDCTokenDataReader(logger.TestLogger(t), nil, mockUSDCTokenAddress, mockMsgTransmitter, mockOnRampAddress, nil)

	filters := usdcService.GetSourceLogPollerFilters()

	require.Equal(t, 1, len(filters))
	filter := filters[0]
	require.Equal(t, logpoller.FilterName(MESSAGE_SENT_FILTER_NAME, mockMsgTransmitter.Hex()), filter.Name)
	hash, err := utils.Keccak256([]byte("MessageSent(bytes)"))
	require.NoError(t, err)
	require.Equal(t, hash, filter.EventSigs[0].Bytes())
	require.Equal(t, mockMsgTransmitter, filter.Addresses[0])
}

func TestGetUSDCMessageBody(t *testing.T) {
	expectedBody, err := hexutil.Decode("0x0000000000000001000000020000000000048D71000000000000000000000000EB08F243E5D3FCFF26A9E38AE5520A669F4019D000000000000000000000000023A04D5935ED8BC8E3EB78DB3541F0ABFB001C6E0000000000000000000000006CB3ED9B441EB674B58495C8B3324B59FAFF5243000000000000000000000000000000005425890298AED601595A70AB815C96711A31BC65000000000000000000000000AB4F961939BFE6A93567CC57C59EED7084CE2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035E08285CFED1EF159236728F843286C55FC0861")
	require.NoError(t, err)
	expectedBodyHash := utils.Keccak256Fixed(expectedBody)

	sourceChainEventsMock := ccipdata.MockReader{}
	sourceChainEventsMock.On("GetLastUSDCMessagePriorToLogIndexInTx", mock.Anything, mock.Anything, mock.Anything).Return(expectedBody, nil)

	usdcService := NewUSDCTokenDataReader(logger.TestLogger(t), &sourceChainEventsMock, mockUSDCTokenAddress, mockMsgTransmitter, mockOnRampAddress, nil)

	// Make the first call and assert the underlying function is called
	body, err := usdcService.getUSDCMessageBodyHash(context.Background(), internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{})
	require.NoError(t, err)
	require.Equal(t, body, expectedBodyHash)

	sourceChainEventsMock.AssertNumberOfCalls(t, "GetLastUSDCMessagePriorToLogIndexInTx", 1)

	// Make another call and assert that the cache is used
	body, err = usdcService.getUSDCMessageBodyHash(context.Background(), internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{})
	require.NoError(t, err)
	require.Equal(t, body, expectedBodyHash)
	sourceChainEventsMock.AssertNumberOfCalls(t, "GetLastUSDCMessagePriorToLogIndexInTx", 1)
}
