package gasprice

import (
	"math/big"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
)

type WrongOffchainConfig struct {
	MaxGasPrice1 []int `json:"maxGasPrice1" cbor:"maxGasPrice1"`
}

func TestGasPrice_Check(t *testing.T) {
	lggr := logger.TestLogger(t)
	uid, _ := new(big.Int).SetString("1843548457736589226156809205796175506139185429616502850435279853710366065936", 10)

	tests := []struct {
		Name                   string
		MaxGasPrice            *big.Int
		CurrentLegacyGasPrice  *big.Int
		CurrentDynamicGasPrice *big.Int
		ExpectedResult         encoding.UpkeepFailureReason
		FailedToGetFee         bool
		NotConfigured          bool
		ParsingFailed          bool
	}{
		{
			Name:           "no offchain config",
			ExpectedResult: encoding.UpkeepFailureReasonNone,
		},
		{
			Name:           "maxGasPrice not configured in offchain config",
			NotConfigured:  true,
			ExpectedResult: encoding.UpkeepFailureReasonNone,
		},
		{
			Name:           "fail to parse offchain config",
			ParsingFailed:  true,
			MaxGasPrice:    big.NewInt(10_000_000_000),
			ExpectedResult: encoding.UpkeepFailureReasonNone,
		},
		{
			Name:           "fail to retrieve current gas price",
			MaxGasPrice:    big.NewInt(8_000_000_000),
			FailedToGetFee: true,
			ExpectedResult: encoding.UpkeepFailureReasonNone,
		},
		{
			Name:                  "current gas price is too high - legacy",
			MaxGasPrice:           big.NewInt(10_000_000_000),
			CurrentLegacyGasPrice: big.NewInt(18_000_000_000),
			ExpectedResult:        encoding.UpkeepFailureReasonGasPriceTooHigh,
		},
		{
			Name:                   "current gas price is too high - dynamic",
			MaxGasPrice:            big.NewInt(10_000_000_000),
			CurrentDynamicGasPrice: big.NewInt(15_000_000_000),
			ExpectedResult:         encoding.UpkeepFailureReasonGasPriceTooHigh,
		},
		{
			Name:                  "current gas price is less than user's max gas price - legacy",
			MaxGasPrice:           big.NewInt(8_000_000_000),
			CurrentLegacyGasPrice: big.NewInt(5_000_000_000),
			ExpectedResult:        encoding.UpkeepFailureReasonNone,
		},
		{
			Name:                   "current gas price is less than user's max gas price - dynamic",
			MaxGasPrice:            big.NewInt(10_000_000_000),
			CurrentDynamicGasPrice: big.NewInt(8_000_000_000),
			ExpectedResult:         encoding.UpkeepFailureReasonNone,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctx := testutils.Context(t)
			ge := gasMocks.NewEvmFeeEstimator(t)
			if test.FailedToGetFee {
				ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
					gas.EvmFee{},
					feeLimit,
					errors.New("failed to retrieve gas price"),
				)
			} else if test.CurrentLegacyGasPrice != nil {
				ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
					gas.EvmFee{
						Legacy: assets.NewWei(test.CurrentLegacyGasPrice),
					},
					feeLimit,
					nil,
				)
			} else if test.CurrentDynamicGasPrice != nil {
				ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
					gas.EvmFee{
						DynamicFeeCap: assets.NewWei(test.CurrentDynamicGasPrice),
						DynamicTipCap: assets.NewWei(big.NewInt(1_000_000_000)),
					},
					feeLimit,
					nil,
				)
			}

			var oc []byte
			if test.ParsingFailed {
				oc, _ = cbor.Marshal(WrongOffchainConfig{MaxGasPrice1: []int{1, 2, 3}})
				if len(oc) > 0 {
					oc[len(oc)-1] = 0x99
				}
			} else if test.NotConfigured {
				oc = []byte{1, 2, 3, 4} // parsing this will set maxGasPrice field to nil
			} else if test.MaxGasPrice != nil {
				oc, _ = cbor.Marshal(UpkeepOffchainConfig{MaxGasPrice: test.MaxGasPrice})
			}
			fr := CheckGasPrice(ctx, uid, oc, ge, lggr)
			assert.Equal(t, test.ExpectedResult, fr)
		})
	}
}
