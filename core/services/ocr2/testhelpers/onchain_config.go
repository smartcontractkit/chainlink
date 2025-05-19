package testhelpers

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
)

func GenerateDefaultOCR2OnchainConfig(minValue *big.Int, maxValue *big.Int) ([]byte, error) {
	return median.StandardOnchainConfigCodec{}.Encode(context.Background(), median.OnchainConfig{Min: minValue, Max: maxValue})
}
