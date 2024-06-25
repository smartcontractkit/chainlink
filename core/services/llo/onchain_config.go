package llo

type OnchainConfig struct{}

type OnchainConfigCodec interface {
	Encode(OnchainConfig) ([]byte, error)
	Decode([]byte) (OnchainConfig, error)
}

var _ OnchainConfigCodec = &JSONOnchainConfigCodec{}

// TODO: Replace this with protobuf, if it is actually used for something
type JSONOnchainConfigCodec struct{}

func (c *JSONOnchainConfigCodec) Encode(OnchainConfig) ([]byte, error) {
	return nil, nil
}

func (c *JSONOnchainConfigCodec) Decode([]byte) (OnchainConfig, error) {
	return OnchainConfig{}, nil
}
