package params

import (
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// encodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
// copied from "github.com/terra-money/core/app/params"
type encodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeencodingConfig creates an encodingConfig for an amino based test configuration.
func makeEncodingConfig() encodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	return encodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

// TODO: import as params.MakeEncoding config
var config = makeEncodingConfig()

var initOnce sync.Once

// Initialize the cosmos sdk at most one time
func InitCosmosSdk(bech32Prefix, token string) {
	initOnce.Do(func() { initCosmosSdk(bech32Prefix, token) })
}

func initCosmosSdk(bech32Prefix, token string) {
	// copied from wasmd https://github.com/CosmWasm/wasmd/blob/88e01a98ab8a87b98dc26c03715e6aef5c92781b/app/app.go#L163-L174
	// NOTE: Bech32 is configured globally, blocked on https://github.com/cosmos/cosmos-sdk/issues/13140
	var (
		// bech32PrefixAccAddr defines the Bech32 prefix of an account's address
		bech32PrefixAccAddr = bech32Prefix
		// bech32PrefixAccPub defines the Bech32 prefix of an account's public key
		bech32PrefixAccPub = bech32Prefix + sdk.PrefixPublic
		// bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
		bech32PrefixValAddr = bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
		// bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
		bech32PrefixValPub = bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
		// bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
		bech32PrefixConsAddr = bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
		// bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
		bech32PrefixConsPub = bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
	)

	// Extracted from app.MakeEncodingConfig() to ensure that we only call them once, since they race and can panic.
	std.RegisterLegacyAminoCodec(config.Amino)
	// This registers base sdk, tx and crypto types, see
	// https://github.com/cosmos/cosmos-sdk/blob/47f46643affd7ec7978329c42bac47275ac7e1cc/std/codec.go#L20
	std.RegisterInterfaces(config.InterfaceRegistry)
	// needed for Client.Account() to deserialize authtypes.AccountI
	authtypes.RegisterInterfaces(config.InterfaceRegistry)

	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	sdkConfig.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	sdkConfig.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)
	sdkConfig.Seal()

	for _, d := range []struct {
		denom    string
		decimals int64
	}{
		{token, 0},
		{"m" + token, 3},
		{"u" + token, 6},
		{"n" + token, 9},
	} {
		dec := sdk.NewDecWithPrec(1, d.decimals)
		if err := sdk.RegisterDenom(d.denom, dec); err != nil {
			panic(fmt.Errorf("failed to register denomination %q: %w", d.denom, err))
		}
	}
}

func NewClientContext() client.Context {
	return client.Context{}.
		WithCodec(config.Marshaler).
		WithLegacyAmino(config.Amino).
		WithInterfaceRegistry(config.InterfaceRegistry).
		WithTxConfig(config.TxConfig)
}

func ClientTxConfig() client.TxConfig {
	return config.TxConfig
}
