package chainreadertest

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
)

var Codec = staticCodec{
	staticCodecConfig: staticCodecConfig{
		n:        3,
		itemType: "itemType",
		maxSize:  37,
	},
}

type staticCodecConfig struct {
	n        int
	itemType string
	maxSize  int
}

type staticCodec struct {
	staticCodecConfig
}

var _ testtypes.CodecEvaluator = staticCodec{}

func (c staticCodec) GetMaxEncodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return c.maxSize, nil
}

func (c staticCodec) GetMaxDecodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return c.maxSize, nil
}

func (c staticCodec) Encode(ctx context.Context, item any, itemType string) ([]byte, error) {
	return nil, errors.New("staticCoded.Encode not used for these test")
}

func (c staticCodec) Decode(ctx context.Context, raw []byte, into any, itemType string) error {
	return errors.New("staticCodec.Decode not used for these test")
}

func (c staticCodec) Evaluate(ctx context.Context, other types.Codec) error {
	maxSize, err := c.GetMaxEncodingSize(ctx, c.n, c.itemType)
	if err != nil {
		return err
	}
	if maxSize != c.maxSize {
		return errors.New("unexpected max size")
	}
	return nil
}
