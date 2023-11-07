package evm

import "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

type parsedTypes struct {
	encoderDefs map[string]*types.CodecEntry
	decoderDefs map[string]*types.CodecEntry
	// TODO transforms...
}
