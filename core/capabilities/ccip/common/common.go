package common

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
)

// HashedCapabilityID returns the hashed capability id in a manner equivalent to the capability registry.
func HashedCapabilityID(capabilityLabelledName, capabilityVersion string) (r [32]byte, err error) {
	// TODO: investigate how to avoid parsing the ABI everytime.
	tabi := `[{"type": "string"}, {"type": "string"}]`
	abiEncoded, err := utils.ABIEncode(tabi, capabilityLabelledName, capabilityVersion)
	if err != nil {
		return r, fmt.Errorf("failed to ABI encode capability version and labelled name: %w", err)
	}

	h := crypto.Keccak256(abiEncoded)
	copy(r[:], h)
	return r, nil
}

// PluginConfig is a struct that holds all the necessary information for a CCIP plugin to function.
type PluginConfig struct {
	CommitPluginCodec   cciptypes.CommitPluginCodec
	ExecutePluginCodec  cciptypes.ExecutePluginCodec
	MessageHasher       func(lggr logger.Logger) cciptypes.MessageHasher
	TokenDataEncoder    cciptypes.TokenDataEncoder
	GasEstimateProvider cciptypes.EstimateProvider
	RMNCrypto           func(lggr logger.Logger) cciptypes.RMNCrypto
	// PriceOnlyCommitFn optional method override for price only commit reports.
	PriceOnlyCommitFn    string
	GetChainReaderConfig func(lggr logger.Logger,
		chainID string,
		destChainID string,
		homeChainID string,
		ofc OffChainConfig,
		chainSelector cciptypes.ChainSelector,
	) ([]byte, error)
	GetChainWriter func(
		ctx context.Context,
		chainID string,
		relayer loop.Relayer,
		transmitters map[types.RelayID][]string,
		execBatchGasLimit uint64,
		chainFamily string,
		offrampProgramAddress []byte,
	) (types.ContractWriter, error)
}

type OffChainConfig struct {
	CommitOffchainConfig *pluginconfig.CommitOffchainConfig
	ExecOffchainConfig   *pluginconfig.ExecuteOffchainConfig
}

func (ofc OffChainConfig) CommitEmpty() bool {
	return ofc.CommitOffchainConfig == nil
}

func (ofc OffChainConfig) ExecEmpty() bool {
	return ofc.ExecOffchainConfig == nil
}

func (ofc OffChainConfig) Commit() *pluginconfig.CommitOffchainConfig {
	return ofc.CommitOffchainConfig
}

func (ofc OffChainConfig) Exec() *pluginconfig.ExecuteOffchainConfig {
	return ofc.ExecOffchainConfig
}

// Exactly one of both plugins should be empty at any given time.
func (ofc OffChainConfig) IsValid() bool {
	return (ofc.CommitEmpty() && !ofc.ExecEmpty()) || (!ofc.CommitEmpty() && ofc.ExecEmpty())
}

// CreatePluginConfig creates a PluginConfig for the given chain family.
func CreatePluginConfig(chainFamily string) (PluginConfig, error) {
	extraDataCodec := NewExtraDataCodec(ccipevm.ExtraDataDecoder{}, ccipsolana.ExtraDataDecoder{})
	switch chainFamily {
	case chainsel.FamilyEVM:
		return PluginConfig{
			CommitPluginCodec:  ccipevm.NewCommitPluginCodecV1(),
			ExecutePluginCodec: ccipevm.NewExecutePluginCodecV1(extraDataCodec),
			MessageHasher: func(lggr logger.Logger) cciptypes.MessageHasher {
				return ccipevm.NewMessageHasherV1(lggr, extraDataCodec)
			},
			TokenDataEncoder:     ccipevm.NewEVMTokenDataEncoder(),
			GasEstimateProvider:  ccipevm.NewGasEstimateProvider(),
			RMNCrypto:            func(lggr logger.Logger) cciptypes.RMNCrypto { return ccipevm.NewEVMRMNCrypto(lggr) },
			GetChainReaderConfig: ccipevm.GetEVMChainReaderConfig,
			GetChainWriter:       ccipevm.GetEVMChainWriter,
		}, nil
	case chainsel.FamilySolana:
		return PluginConfig{
			CommitPluginCodec:  ccipsolana.NewCommitPluginCodecV1(),
			ExecutePluginCodec: ccipsolana.NewExecutePluginCodecV1(extraDataCodec),
			MessageHasher: func(lggr logger.Logger) cciptypes.MessageHasher {
				return ccipsolana.NewMessageHasherV1(lggr, extraDataCodec)
			},
			TokenDataEncoder:     ccipsolana.NewSolanaTokenDataEncoder(),
			GasEstimateProvider:  ccipsolana.NewGasEstimateProvider(),
			RMNCrypto:            func(lggr logger.Logger) cciptypes.RMNCrypto { return nil },
			GetChainReaderConfig: ccipsolana.GetSolanaChainReaderConfig,
			GetChainWriter:       ccipsolana.GetSolanaChainWriter,
			PriceOnlyCommitFn:    consts.MethodCommitPriceOnly,
		}, nil
	}

	return PluginConfig{}, fmt.Errorf("unsupported chain family: %s", chainFamily)
}
