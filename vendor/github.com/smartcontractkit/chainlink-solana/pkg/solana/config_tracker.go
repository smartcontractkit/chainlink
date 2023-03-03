package solana

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
)

type ConfigTracker struct {
	stateCache *StateCache
	reader     client.Reader
}

func (c *ConfigTracker) Notify() <-chan struct{} {
	return nil // not using websocket, config changes will be handled by polling in libocr
}

// LatestConfigDetails returns information about the latest configuration,
// but not the configuration itself.
func (c *ConfigTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	state, err := c.stateCache.ReadState()
	return state.Config.LatestConfigBlockNumber, state.Config.LatestConfigDigest, err
}

func ConfigFromState(state State) (types.ContractConfig, error) {
	pubKeys := []types.OnchainPublicKey{}
	accounts := []types.Account{}
	oracles, err := state.Oracles.Data()
	if err != nil {
		return types.ContractConfig{}, err
	}
	for _, o := range oracles {
		o := o //  https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		pubKeys = append(pubKeys, o.Signer.Key[:])
		accounts = append(accounts, types.Account(o.Transmitter.String()))
	}

	onchainConfigStruct := median.OnchainConfig{
		Min: state.Config.MinAnswer.BigInt(),
		Max: state.Config.MaxAnswer.BigInt(),
	}

	onchainConfig, err := median.StandardOnchainConfigCodec{}.Encode(onchainConfigStruct)
	if err != nil {
		return types.ContractConfig{}, err
	}
	offchainConfig, err := state.OffchainConfig.Data()
	if err != nil {
		return types.ContractConfig{}, err
	}

	return types.ContractConfig{
		ConfigDigest:          state.Config.LatestConfigDigest,
		ConfigCount:           uint64(state.Config.ConfigCount),
		Signers:               pubKeys,
		Transmitters:          accounts,
		F:                     state.Config.F,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: state.OffchainConfig.Version,
		OffchainConfig:        offchainConfig,
	}, nil
}

// LatestConfig returns the latest configuration.
func (c *ConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	state, err := c.stateCache.ReadState()
	if err != nil {
		return types.ContractConfig{}, err
	}
	return ConfigFromState(state)
}

// LatestBlockHeight returns the height of the most recent block in the chain.
func (c *ConfigTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return c.reader.SlotHeight() // this returns the latest slot height through CommitmentProcessed
}
