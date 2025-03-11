package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// CreatePluginConfig creates a PluginConfig for the given chain family.
func CreatePluginConfig(chainFamily string) (cctypes.PluginConfig, error) {
	switch chainFamily {
	case chainsel.FamilyEVM:
		return cctypes.PluginConfig{
			CommitPluginCodec:   ccipevm.NewCommitPluginCodecV1(),
			ExecutePluginCodec:  ccipevm.NewExecutePluginCodecV1(),
			ExtraArgsCodec:      NewExtraDataCodec(),
			MessageHasher:       func(lggr logger.Logger) cciptypes.MessageHasher { return ccipevm.NewMessageHasherV1(lggr) },
			TokenDataEncoder:    ccipevm.NewEVMTokenDataEncoder(),
			GasEstimateProvider: ccipevm.NewGasEstimateProvider(),
			RMNCrypto:           func(lggr logger.Logger) cciptypes.RMNCrypto { return ccipevm.NewEVMRMNCrypto(lggr) },
			AddressToString: func(addr []byte, checkSum bool) string {
				offRampAddr := common.BytesToAddress(addr).Hex()
				if !checkSum {
					offRampAddr = hexutil.Encode(addr)
				}
				return offRampAddr
			},
			GetChainReaderConfig: ccipevm.GetEVMChainReaderConfig,
			GetChainWriter:       ccipevm.GetEVMChainWriter,
		}, nil
	case chainsel.FamilySolana:
		return cctypes.PluginConfig{
			CommitPluginCodec:    ccipsolana.NewCommitPluginCodecV1(),
			ExecutePluginCodec:   ccipsolana.NewExecutePluginCodecV1(),
			ExtraArgsCodec:       NewExtraDataCodec(),
			MessageHasher:        func(lggr logger.Logger) cciptypes.MessageHasher { return ccipsolana.NewMessageHasherV1(lggr) },
			TokenDataEncoder:     ccipsolana.NewSolanaTokenDataEncoder(),
			GasEstimateProvider:  ccipsolana.NewGasEstimateProvider(),
			RMNCrypto:            func(lggr logger.Logger) cciptypes.RMNCrypto { return nil },
			AddressToString:      func(addr []byte, checkSum bool) string { return solana.PublicKeyFromBytes(addr).String() },
			GetChainReaderConfig: ccipsolana.GetSolanaChainReaderConfig,
			GetChainWriter:       ccipsolana.GetSolanaChainWriter,
		}, nil
	}

	return cctypes.PluginConfig{}, fmt.Errorf("unsupported chain family: %s", chainFamily)
}
