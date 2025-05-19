package evm

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint:revive // dot-imports
)

// ClientWithContractHistory makes it possible to modify client.Client CallContract so that it returns historical data.
type ClientWithContractHistory struct {
	client.Client
	valsWithCall
	HT    logpoller.HeadTracker
	codec clcommontypes.RemoteCodec
}

func (cwh *ClientWithContractHistory) Init(_ context.Context, config types.ChainReaderConfig) error {
	cwh.valsWithCall = make(map[int64]valWithCall)
	parsedTypes := codec.ParsedTypes{
		EncoderDefs: make(map[string]types.CodecEntry),
		DecoderDefs: make(map[string]types.CodecEntry),
	}

	// setup codec for method calls
	for contractName, contractCfg := range config.Contracts {
		contractAbi, err := abi.JSON(strings.NewReader(contractCfg.ContractABI))
		if err != nil {
			return err
		}

		for genericName, readDef := range contractCfg.Configs {
			if readDef.ReadType == types.Event {
				continue
			}

			injectEVMSpecificCodecModifiers(readDef)

			inputMod, err := readDef.InputModifications.ToModifier(codec.DecoderHooks...)
			if err != nil {
				return err
			}

			outputMod, err := readDef.OutputModifications.ToModifier(codec.DecoderHooks...)
			if err != nil {
				return err
			}

			method := contractAbi.Methods[readDef.ChainSpecificName]
			input, output := types.NewCodecEntry(method.Inputs, method.ID, inputMod), types.NewCodecEntry(method.Outputs, nil, outputMod)

			if err = input.Init(); err != nil {
				return err
			}
			if err = output.Init(); err != nil {
				return err
			}

			parsedTypes.EncoderDefs[codec.WrapItemType(contractName, genericName, true)] = input
			parsedTypes.DecoderDefs[codec.WrapItemType(contractName, genericName, false)] = output
		}
	}

	parsedCodec, err := parsedTypes.ToCodec()
	if err != nil {
		return err
	}

	cwh.codec = parsedCodec

	return nil
}

func (cwh *ClientWithContractHistory) SetUintLatestValue(ctx context.Context, val uint64, forCall ExpectedGetLatestValueArgs) error {
	latestBlock, _, err := cwh.HT.LatestAndFinalizedBlock(ctx)
	if err != nil {
		return err
	}

	cwh.valsWithCall[latestBlock.Number] = valWithCall{
		ExpectedGetLatestValueArgs: forCall,
		val:                        val,
	}
	return nil
}

// findValClosestToBlock returns the value that was added in block that is <= blockNumber, otherwise return empty valWithCall
func (cwh *ClientWithContractHistory) findValClosestToBlock(blockNumber *big.Int) (val valWithCall) {
	var valIsInBlock int64 = math.MaxInt64
	for block, v := range cwh.valsWithCall {
		if block <= blockNumber.Int64() && block <= valIsInBlock {
			valIsInBlock, val = block, v
		}
	}
	return val
}

// CallContract does the standard CallContract if blockNumber is latest, otherwise returns mocked historical value since SimulatedBackend doesn't support historical calls.
func (cwh *ClientWithContractHistory) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	latestBlock, _, err := cwh.HT.LatestAndFinalizedBlock(ctx)
	if err != nil {
		return nil, err
	}

	if blockNumber == nil || blockNumber.Int64() == latestBlock.BlockNumber() {
		return cwh.Client.CallContract(ctx, msg, blockNumber)
	}

	// hardcoded to work with the simplest test case of uint64 return value with no params, since we just want to test proper finality handling.
	valAndCall := cwh.findValClosestToBlock(blockNumber)
	if valAndCall.ContractName == "" || valAndCall.ReadName == "" {
		return nil, fmt.Errorf("no recorded values found for historical contract call in block : %s, stored historical values are: \n\t%s", blockNumber, cwh.valsWithCall.String())
	}

	// encode the expected call to compare with the actual call
	dataToCmp, err := cwh.codec.Encode(ctx, valAndCall.Params, codec.WrapItemType(valAndCall.ContractName, valAndCall.ReadName, true))
	if err != nil {
		return nil, err
	}

	if string(msg.Data) == string(dataToCmp) {
		// format the value to match the standard evm rpc response
		paddedHexString := fmt.Sprintf("%064s", hexutil.EncodeUint64(valAndCall.val)[2:])
		return common.Hex2Bytes(paddedHexString), nil
	}

	return nil, fmt.Errorf("found historical contract value, but it doesn't match the contract call, stored historical values are: \n\t%s", cwh.valsWithCall.String())
}

// valsWithCall is used to mock historical data since SimulatedBackend doesn't support it
type valsWithCall map[int64]valWithCall

func (v valsWithCall) String() string {
	result := "valsWithCall{\n"
	for key, val := range v {
		result += fmt.Sprintf("  Blocknumber: %d, Value: %s\n", key, val.String())
	}
	result += "}"
	return result
}

type valWithCall struct {
	ExpectedGetLatestValueArgs
	val uint64
}

func (v valWithCall) String() string {
	return fmt.Sprintf("{ExpectedArgs: {%s}, Value: %d}", v.ExpectedGetLatestValueArgs.String(), v.val)
}
