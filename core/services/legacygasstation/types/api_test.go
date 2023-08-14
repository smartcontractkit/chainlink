package types

import (
	"bytes"
	"encoding/base64"
	"math/big"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
)

func TestSendTransactionRequestTransform(t *testing.T) {
	requestBody := `{
		"from": "0x8b94A8792dcbb482F2A49569AcE5E7c29fF5c93d",
		"target": "0x61aF15229af1CEd7Ca4a80f623F85a8a91420C04",
		"target_name": "sampleERC20",
		"version": "1",
		"nonce": 41348595284370659424721411682442273173799930917222432033995375917025790037177,
		"receiver": "0xbF7872b0253227F1251cf26dD9E892a3a4158FaF",
		"amount": 10000000000000000000000,
		"chain_id": 5,
		"destination_chain_id": 20000,
		"valid_until_time": 1682449148,
		"signature": "a9VQaaVBf5W2O/rppOutrjsoq9Sk7m+aVoBuT/2o2ykT3hzzHQmtDmELLr/noQeUPqHdSWDPh1xL540G/FNm+xs="
	}`

	httpRequest, err := http.NewRequest("POST", "/send_transaction", strings.NewReader(requestBody))
	require.NoError(t, err)
	var request SendTransactionRequest
	err = binding.JSON.Bind(httpRequest, &request)
	require.NoError(t, err)

	require.Equal(t, "0x8b94A8792dcbb482F2A49569AcE5E7c29fF5c93d", request.From.String())
	require.Equal(t, "0x61aF15229af1CEd7Ca4a80f623F85a8a91420C04", request.Target.String())
	require.Equal(t, "sampleERC20", request.TargetName)
	require.Equal(t, "1", request.Version)
	nonce, ok := big.NewInt(0).SetString("41348595284370659424721411682442273173799930917222432033995375917025790037177", 10)
	require.True(t, ok)
	require.True(t, nonce.Cmp(request.Nonce) == 0)
	require.Equal(t, "0xbF7872b0253227F1251cf26dD9E892a3a4158FaF", request.Receiver.String())
	amount, ok := big.NewInt(0).SetString("10000000000000000000000", 10)
	require.True(t, ok)
	require.True(t, amount.Cmp(request.Amount) == 0)
	require.Equal(t, uint64(5), request.SourceChainID)
	require.Equal(t, uint64(20000), request.DestinationChainID)
	validUntilTime, ok := big.NewInt(0).SetString("1682449148", 10)
	require.True(t, ok)
	require.True(t, validUntilTime.Cmp(request.ValidUntilTime) == 0)
	require.True(t, amount.Cmp(request.Amount) == 0)
	decodedBytes, err := base64.StdEncoding.DecodeString("a9VQaaVBf5W2O/rppOutrjsoq9Sk7m+aVoBuT/2o2ykT3hzzHQmtDmELLr/noQeUPqHdSWDPh1xL540G/FNm+xs=")
	if err != nil {
		panic(err)
	}
	require.True(t, bytes.Equal(decodedBytes, request.Signature))
}
