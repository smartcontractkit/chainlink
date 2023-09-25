package usdc_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata/usdc"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	mockOnRampAddress    = utils.RandomAddress()
	mockUSDCTokenAddress = utils.RandomAddress()
	mockMsgTransmitter   = utils.RandomAddress()
)

type attestationResponse struct {
	Status      string `json:"status"`
	Attestation string `json:"attestation"`
}

type messageAndAttestation struct {
	Message     []byte
	Attestation []byte
}

func (m messageAndAttestation) AbiString() string {
	return `
	[{
		"components": [
			{"name": "message", "type": "bytes"},
			{"name": "attestation", "type": "bytes"}
		],
		"type": "tuple"
	}]`
}

func (m messageAndAttestation) Validate() error {
	return nil
}

type usdcPayload []byte

func (d usdcPayload) AbiString() string {
	return `[{"type": "bytes"}]`
}

func (d usdcPayload) Validate() error {
	return nil
}

func TestUSDCReader_ReadTokenData(t *testing.T) {
	response := attestationResponse{
		Status:      "complete",
		Attestation: "0x9049623e91719ef2aa63c55f357be2529b0e7122ae552c18aff8db58b4633c4d3920ff03d3a6d1ddf11f06bf64d7fd60d45447ac81f527ba628877dc5ca759651b08ffae25a6d3b1411749765244f0a1c131cbfe04430d687a2e12fd9d2e6dc08e118ad95d94ad832332cf3c4f7a4f3da0baa803b7be024b02db81951c0f0714de1b",
	}
	abiEncodedMessageBody, err := hexutil.Decode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000f80000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc08610000000000000000")
	require.NoError(t, err)
	rawMessageBody, err := abihelpers.DecodeAbiStruct[usdcPayload](abiEncodedMessageBody)
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		messageHash := utils.Keccak256Fixed(rawMessageBody)
		expectedUrl := "/v1/attestations/0x" + hex.EncodeToString(messageHash[:])
		require.Equal(t, expectedUrl, r.URL.Path)

		responseBytes, err2 := json.Marshal(response)
		require.NoError(t, err2)

		_, err2 = w.Write(responseBytes)
		require.NoError(t, err2)
	}))

	defer ts.Close()

	seqNum := uint64(23825)
	txHash := utils.RandomBytes32()
	logIndex := int64(4)

	eventsClient := ccipdata.MockReader{}
	eventsClient.On("GetSendRequestsBetweenSeqNums",
		mock.Anything,
		mockOnRampAddress,
		seqNum,
		seqNum,
		0,
	).Return([]ccipdata.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested]{
		{
			Data: evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested{
				Raw: types.Log{
					TxHash: txHash,
					Index:  uint(logIndex),
				},
			},
		},
	}, nil)

	eventsClient.On("GetLastUSDCMessagePriorToLogIndexInTx",
		mock.Anything,
		logIndex,
		common.Hash(txHash),
	).Return(abiEncodedMessageBody, nil)
	attestationURI, err := url.ParseRequestURI(ts.URL)
	require.NoError(t, err)

	usdcService := usdc.NewUSDCTokenDataReader(logger.TestLogger(t), &eventsClient, mockUSDCTokenAddress, mockMsgTransmitter, mockOnRampAddress, attestationURI)
	msgAndAttestation, err := usdcService.ReadTokenData(context.Background(), internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{
			SequenceNumber: seqNum,
		},
		TxHash:   txHash,
		LogIndex: uint(logIndex),
	})
	require.NoError(t, err)

	attestationBytes, err := hex.DecodeString(strings.TrimPrefix(response.Attestation, "0x"))
	require.NoError(t, err)

	encodeAbiStruct, err := abihelpers.EncodeAbiStruct[messageAndAttestation](messageAndAttestation{
		Message:     rawMessageBody,
		Attestation: attestationBytes,
	})
	require.NoError(t, err)

	require.Equal(t, encodeAbiStruct, msgAndAttestation)
}
