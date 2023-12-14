package testhelpers

import (
	"math/big"

	"github.com/smartcontractkit/libocr/bigbigendian"
)

func GenerateDefaultOCR2OnchainConfig(minValue *big.Int, maxValue *big.Int) ([]byte, error) {
	serializedConfig := make([]byte, 0)

	s1, err := bigbigendian.SerializeSigned(1, big.NewInt(1)) //version
	if err != nil {
		return nil, err
	}
	serializedConfig = append(serializedConfig, s1...)

	s2, err := bigbigendian.SerializeSigned(24, minValue) //min
	if err != nil {
		return nil, err
	}
	serializedConfig = append(serializedConfig, s2...)

	s3, err := bigbigendian.SerializeSigned(24, maxValue) //max
	if err != nil {
		return nil, err
	}
	serializedConfig = append(serializedConfig, s3...)

	return serializedConfig, nil
}
