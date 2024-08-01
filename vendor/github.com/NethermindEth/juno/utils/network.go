package utils

import (
	"encoding"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/spf13/pflag"
)

var ErrUnknownNetwork = errors.New("unknown network (known: mainnet, goerli, goerli2, integration)")

type Network int

// The following are necessary for Cobra and Viper, respectively, to unmarshal log level
// CLI/config parameters properly.
var (
	_ pflag.Value              = (*Network)(nil)
	_ encoding.TextUnmarshaler = (*Network)(nil)
)

const (
	MAINNET Network = iota
	GOERLI
	GOERLI2
	INTEGRATION
)

func (n Network) String() string {
	switch n {
	case MAINNET:
		return "mainnet"
	case GOERLI:
		return "goerli"
	case GOERLI2:
		return "goerli2"
	case INTEGRATION:
		return "integration"
	default:
		// Should not happen.
		panic(ErrUnknownNetwork)
	}
}

func (n *Network) Set(s string) error {
	switch s {
	case "MAINNET", "mainnet":
		*n = MAINNET
	case "GOERLI", "goerli":
		*n = GOERLI
	case "GOERLI2", "goerli2":
		*n = GOERLI2
	case "INTEGRATION", "integration":
		*n = INTEGRATION
	default:
		return ErrUnknownNetwork
	}
	return nil
}

func (n *Network) Type() string {
	return "Network"
}

func (n *Network) UnmarshalText(text []byte) error {
	return n.Set(string(text))
}

func (n Network) URL() string {
	switch n {
	case GOERLI:
		return "https://alpha4.starknet.io/feeder_gateway/"
	case MAINNET:
		return "https://alpha-mainnet.starknet.io/feeder_gateway/"
	case GOERLI2:
		return "https://alpha4-2.starknet.io/feeder_gateway/"
	case INTEGRATION:
		return "https://external.integration.starknet.io/feeder_gateway/"
	default:
		// Should not happen.
		panic(ErrUnknownNetwork)
	}
}

func (n Network) ChainID() *felt.Felt {
	switch n {
	case GOERLI, INTEGRATION:
		return new(felt.Felt).SetBytes([]byte("SN_GOERLI"))
	case MAINNET:
		return new(felt.Felt).SetBytes([]byte("SN_MAIN"))
	case GOERLI2:
		return new(felt.Felt).SetBytes([]byte("SN_GOERLI2"))
	default:
		// Should not happen.
		panic(ErrUnknownNetwork)
	}
}
