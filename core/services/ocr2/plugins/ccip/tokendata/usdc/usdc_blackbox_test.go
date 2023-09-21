package usdc_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
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

func TestUSDCReader_ReadTokenData(t *testing.T) {
	response := attestationResponse{
		Status:      "complete",
		Attestation: "720502893578a89a8a87982982ef781c18b193",
	}

	attestationBytes, err := hex.DecodeString(response.Attestation)
	require.NoError(t, err)

	responseBytes, err := json.Marshal(response)
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err = w.Write(responseBytes)
		require.NoError(t, err)
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
	).Return(attestationBytes, nil)
	attestationURI, err := url.ParseRequestURI(ts.URL)
	require.NoError(t, err)

	usdcService := usdc.NewUSDCTokenDataReader(&eventsClient, mockUSDCTokenAddress, mockMsgTransmitter, mockOnRampAddress, attestationURI)
	attestation, err := usdcService.ReadTokenData(context.Background(), internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		InternalEVM2EVMMessage: evm_2_evm_offramp.InternalEVM2EVMMessage{
			SequenceNumber: seqNum,
		},
		TxHash:   txHash,
		LogIndex: uint(logIndex),
	})
	require.NoError(t, err)

	require.Equal(t, attestationBytes, attestation)
}
