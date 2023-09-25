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
	expectedBody, err := hexutil.Decode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000f80000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc08610000000000000000")
	require.NoError(t, err)

	parsedBody, err := decodeUSDCMessageSent(expectedBody)
	require.NoError(t, err)

	expectedPostParse := "0x0000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc0861"

	require.Equal(t, expectedPostParse, hexutil.Encode(parsedBody))

	sourceChainEventsMock := ccipdata.MockReader{}
	sourceChainEventsMock.On("GetLastUSDCMessagePriorToLogIndexInTx", mock.Anything, mock.Anything, mock.Anything).Return(expectedBody, nil)

	usdcService := NewUSDCTokenDataReader(logger.TestLogger(t), &sourceChainEventsMock, mockUSDCTokenAddress, mockMsgTransmitter, mockOnRampAddress, nil)

	// Make the first call and assert the underlying function is called
	body, err := usdcService.getUSDCMessageBody(context.Background(), internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{})
	require.NoError(t, err)
	require.Equal(t, body, parsedBody)

	sourceChainEventsMock.AssertNumberOfCalls(t, "GetLastUSDCMessagePriorToLogIndexInTx", 1)

	// Make another call and assert that the cache is used
	body, err = usdcService.getUSDCMessageBody(context.Background(), internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{})
	require.NoError(t, err)
	require.Equal(t, body, parsedBody)
	sourceChainEventsMock.AssertNumberOfCalls(t, "GetLastUSDCMessagePriorToLogIndexInTx", 1)
}
