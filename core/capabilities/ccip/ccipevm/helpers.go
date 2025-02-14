package ccipevm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/oraclecreator"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const (
	svmV1DecodeName       = "decodeSVMExtraArgsV1"
	evmV1DecodeName       = "decodeEVMExtraArgsV1"
	evmV2DecodeName       = "decodeEVMExtraArgsV2"
	evmDestExecDataKey    = "destGasAmount"
	defaultCommitGasLimit = 500_000
)

var (
	abiUint32               = ABITypeOrPanic("uint32")
	TokenDestGasOverheadABI = abi.Arguments{
		{
			Type: abiUint32,
		},
	}
)

func decodeExtraArgsV1V2(extraArgs []byte) (gasLimit *big.Int, err error) {
	if len(extraArgs) < 4 {
		return nil, fmt.Errorf("extra args too short: %d, should be at least 4 (i.e the extraArgs tag)", len(extraArgs))
	}

	var method string
	if bytes.Equal(extraArgs[:4], evmExtraArgsV1Tag) {
		method = evmV1DecodeName
	} else if bytes.Equal(extraArgs[:4], evmExtraArgsV2Tag) {
		method = evmV2DecodeName
	} else {
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs)
	}
	ifaces, err := messageHasherABI.Methods[method].Inputs.UnpackValues(extraArgs[4:])
	if err != nil {
		return nil, fmt.Errorf("abi decode extra args v1: %w", err)
	}
	// gas limit is always the first argument, and allow OOO isn't set explicitly
	// on the message.
	_, ok := ifaces[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("expected *big.Int, got %T", ifaces[0])
	}
	return ifaces[0].(*big.Int), nil
}

// abiEncodeMethodInputs encodes the inputs for a method call.
// example abi: `[{ "name" : "method", "type": "function", "inputs": [{"name": "a", "type": "uint256"}]}]`
func abiEncodeMethodInputs(abiDef abi.ABI, inputs ...interface{}) ([]byte, error) {
	packed, err := abiDef.Pack("method", inputs...)
	if err != nil {
		return nil, err
	}
	return packed[4:], nil // remove the method selector
}

func ABITypeOrPanic(t string) abi.Type {
	abiType, err := abi.NewType(t, "", nil)
	if err != nil {
		panic(err)
	}
	return abiType
}

// Decodes the given bytes into a uint32, based on the encoding of destGasAmount in FeeQuoter.sol
func decodeTokenDestGasOverhead(destExecData []byte) (uint32, error) {
	ifaces, err := TokenDestGasOverheadABI.UnpackValues(destExecData)
	if err != nil {
		return 0, fmt.Errorf("abi decode TokenDestGasOverheadABI: %w", err)
	}
	_, ok := ifaces[0].(uint32)
	if !ok {
		return 0, fmt.Errorf("expected uint32, got %T", ifaces[0])
	}
	return ifaces[0].(uint32), nil
}

func CreatePluginConfig() oraclecreator.Plugin {
	return oraclecreator.Plugin{
		CommitPluginCodec:   NewCommitPluginCodecV1(),
		ExecutePluginCodec:  NewExecutePluginCodecV1(),
		ExtraArgsCodec:      ccipcommon.NewExtraDataCodec(),
		MessageHasher:       func(lggr logger.Logger) cciptypes.MessageHasher { return NewMessageHasherV1(lggr) },
		TokenDataEncoder:    NewEVMTokenDataEncoder(),
		GasEstimateProvider: NewGasEstimateProvider(),
		RMNCrypto:           func(lggr logger.Logger) cciptypes.RMNCrypto { return NewEVMRMNCrypto(lggr) },
		AddressToString: func(addr []byte, checkSum bool) string {
			offRampAddr := common.BytesToAddress(addr).Hex()
			if !checkSum {
				offRampAddr = hexutil.Encode(addr)
			}
			return offRampAddr
		},
		GetChainReaderConfig: getEVMChainReaderConfig,
		GetChainWriter:       getEVMChainWriter,
	}
}

func getEVMChainReaderConfig(
	lggr logger.Logger,
	chainID string,
	destChainID string,
	homeChainID string,
	ofc cctypes.OffChainConfig,
	chainSelector cciptypes.ChainSelector,
) ([]byte, error) {
	var chainReaderConfig evmrelaytypes.ChainReaderConfig
	if chainID == destChainID {
		chainReaderConfig = evmconfig.DestReaderConfig
	} else {
		chainReaderConfig = evmconfig.SourceReaderConfig
	}

	if !ofc.CommitEmpty() && ofc.Commit().PriceFeedChainSelector == chainSelector {
		lggr.Debugw("Adding feed reader config", "chainID", chainID)
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.FeedReaderConfig)
	}

	if isUSDCEnabled(ofc) {
		lggr.Debugw("Adding USDC reader config", "chainID", chainID)
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.USDCReaderConfig)
	}

	if chainID == homeChainID {
		lggr.Debugw("Adding home chain reader config", "chainID", chainID)
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.HomeChainReaderConfigRaw)
	}

	marshaledConfig, err := json.Marshal(chainReaderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain reader config: %w", err)
	}

	return marshaledConfig, nil
}

func getEVMChainWriter(
	ctx context.Context,
	chainID string,
	relayer loop.Relayer,
	transmitters map[types.RelayID][]string,
	execBatchGasLimit uint64,
	chainFamily string,
	offrampProgramAddress []byte,
	destChainSelector uint64,
) (types.ContractWriter, error) {
	var fromAddress common.Address
	transmitter, ok := transmitters[types.NewRelayID(chainFamily, chainID)]
	if ok {
		fromAddress = common.HexToAddress(transmitter[0])
	}

	evmConfig, err := evmconfig.ChainWriterConfigRaw(
		fromAddress,
		defaultCommitGasLimit,
		execBatchGasLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create EVM chain writer config: %w", err)
	}

	chainWriterConfig, err := json.Marshal(evmConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal EVM chain writer config: %w", err)
	}

	cw, err := relayer.NewContractWriter(ctx, chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer for chain %s: %w", chainID, err)
	}

	return cw, nil
}

func isUSDCEnabled(ofc cctypes.OffChainConfig) bool {
	if ofc.ExecEmpty() {
		return false
	}

	return ofc.Exec().IsUSDCEnabled()
}
