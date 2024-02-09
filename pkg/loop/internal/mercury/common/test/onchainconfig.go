package mercury_common_test

import (
	"fmt"
	"math/big"
	"reflect"

	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

type OnChainConfigCodecParameters struct {
	Encoded []byte
	Decoded mercury_types.OnchainConfig
}

var StaticOnChainConfigCodecFixtures = OnChainConfigCodecParameters{
	Encoded: []byte("on chain config to be encoded"),
	Decoded: mercury_types.OnchainConfig{
		Min: big.NewInt(1),
		Max: big.NewInt(100),
	},
}

type StaticOnchainConfigCodec struct{}

var _ mercury_types.OnchainConfigCodec = StaticOnchainConfigCodec{}

func (StaticOnchainConfigCodec) Encode(c mercury_types.OnchainConfig) ([]byte, error) {
	if !reflect.DeepEqual(c, StaticOnChainConfigCodecFixtures.Decoded) {
		return nil, fmt.Errorf("expected OnchainConfig %v but got %v", StaticOnChainConfigCodecFixtures.Decoded, c)
	}

	return StaticOnChainConfigCodecFixtures.Encoded, nil
}

func (StaticOnchainConfigCodec) Decode([]byte) (mercury_types.OnchainConfig, error) {
	return StaticOnChainConfigCodecFixtures.Decoded, nil
}
