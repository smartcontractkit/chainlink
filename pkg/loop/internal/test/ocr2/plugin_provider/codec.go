package pluginprovider

import (
	"context"
	"errors"
)

type staticCodec struct{}

func (c staticCodec) GetMaxEncodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return 0, errors.New("not used for these test")
}

func (c staticCodec) GetMaxDecodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return 0, errors.New("not used for these test")
}

func (c staticCodec) Encode(ctx context.Context, item any, itemType string) ([]byte, error) {
	return nil, errors.New("not used for these test")
}

func (c staticCodec) Decode(ctx context.Context, raw []byte, into any, itemType string) error {
	return errors.New("not used for these test")
}
