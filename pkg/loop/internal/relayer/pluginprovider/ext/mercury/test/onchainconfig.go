package mercury_common_test

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

// OnchainConfigCodec is a static implementation of OnchainConfigCodec for testing
var OnchainConfigCodec = staticOnchainConfigCodec{
	onChainConfigCodecParameters: StaticOnChainConfigCodecFixtures,
}

type OnchainConfigCodecEvaluator interface {
	mercury_types.OnchainConfigCodec
	testtypes.Evaluator[mercury_types.OnchainConfigCodec]
}

type onChainConfigCodecParameters struct {
	Encoded []byte
	Decoded mercury_types.OnchainConfig
}

var StaticOnChainConfigCodecFixtures = onChainConfigCodecParameters{
	Encoded: []byte("on chain config to be encoded"),
	Decoded: mercury_types.OnchainConfig{
		Min: big.NewInt(1),
		Max: big.NewInt(100),
	},
}

type staticOnchainConfigCodec struct {
	onChainConfigCodecParameters
}

var _ OnchainConfigCodecEvaluator = staticOnchainConfigCodec{}

func (staticOnchainConfigCodec) Encode(ctx context.Context, c mercury_types.OnchainConfig) ([]byte, error) {
	if !reflect.DeepEqual(c, StaticOnChainConfigCodecFixtures.Decoded) {
		return nil, fmt.Errorf("expected OnchainConfig %v but got %v", StaticOnChainConfigCodecFixtures.Decoded, c)
	}

	return StaticOnChainConfigCodecFixtures.Encoded, nil
}

func (staticOnchainConfigCodec) Decode(context.Context, []byte) (mercury_types.OnchainConfig, error) {
	return StaticOnChainConfigCodecFixtures.Decoded, nil
}

func (staticOnchainConfigCodec) Evaluate(ctx context.Context, other mercury_types.OnchainConfigCodec) error {
	encoded, err := other.Encode(ctx, StaticOnChainConfigCodecFixtures.Decoded)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}
	if !reflect.DeepEqual(encoded, StaticOnChainConfigCodecFixtures.Encoded) {
		return fmt.Errorf("expected encoded %x but got %x", StaticOnChainConfigCodecFixtures.Encoded, encoded)
	}

	decoded, err := other.Decode(ctx, StaticOnChainConfigCodecFixtures.Encoded)
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}
	if !reflect.DeepEqual(decoded, StaticOnChainConfigCodecFixtures.Decoded) {
		return fmt.Errorf("expected decoded %v but got %v", StaticOnChainConfigCodecFixtures.Decoded, decoded)
	}

	return nil
}
