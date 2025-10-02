package common

import (
	"sync"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// Ensure ExtraDataCodecRegistry implements the ExtraDataCodecBundle interface from chainlink-common
var _ cciptypes.ExtraDataCodecBundle = (*ExtraDataCodecRegistry)(nil)

// ExtraDataCodecRegistry is a singleton registry that manages SourceChainExtraDataCodec instances
// for different chain families. It implements the ExtraDataCodecBundle interface from chainlink-common
// by delegating to the existing ExtraDataCodec implementation.
//
// Terminology:
//   - "ExtraDataCodecRegistry": refers to the entire singleton registry instance. It both maintains the map of
//     chain family to codec and provides thread-safe access to it.
//   - "ExtraDataCodecBundle": is the interface defined in chainlink-common that the registry implements and that
//     can be called over gRPC.
type ExtraDataCodecRegistry struct {
	extraDataCodec cciptypes.ExtraDataCodecMap
	mu             sync.RWMutex
}

var (
	registryInstance *ExtraDataCodecRegistry
	registryOnce     sync.Once
)

// GetExtraDataCodecRegistry returns the singleton instance of ExtraDataCodecRegistry. This is only called
// in core node.
func GetExtraDataCodecRegistry() *ExtraDataCodecRegistry {
	registryOnce.Do(func() {
		registryInstance = &ExtraDataCodecRegistry{
			extraDataCodec: make(cciptypes.ExtraDataCodecMap),
		}
	})
	return registryInstance
}

// RegisterFamilyNoopCodec registers a chain family with a no-op SourceChainExtraDataCodec if not already registered.
// This is used when we know which chain families we want to support but don't have a specific codec
// implementation initialized for it yet. This should only be called from core node, not over gRPC.
func (r *ExtraDataCodecRegistry) RegisterFamilyNoopCodec(chainFamily string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.extraDataCodec[chainFamily]; !exists {
		r.extraDataCodec[chainFamily] = NoOpSourceChainExtraDataCodec{}
	}
}

// RegisterCodec registers a SourceChainExtraDataCodec for a specific chain family and is only called
// within core node, not over gRPC.
func (r *ExtraDataCodecRegistry) RegisterCodec(chainFamily string, codec SourceChainExtraDataCodec) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.extraDataCodec[chainFamily] = codec
}

// ============ gRPC-compatible implementation of ExtraDataCodecBundle interface ============

// DecodeExtraArgs can be called either over gRPC or not. It is used to decode extra args for a specific
// source chain family
func (r *ExtraDataCodecRegistry) DecodeExtraArgs(
	extraArgs cciptypes.Bytes,
	sourceChainSelector cciptypes.ChainSelector,
) (map[string]any, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.extraDataCodec.DecodeExtraArgs(extraArgs, sourceChainSelector)
}

// DecodeTokenAmountDestExecData can be called either over gRPC or not. It is used to decode dest exec
// data for a specific source chain family.
func (r *ExtraDataCodecRegistry) DecodeTokenAmountDestExecData(
	destExecData cciptypes.Bytes,
	sourceChainSelector cciptypes.ChainSelector,
) (map[string]any, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.extraDataCodec.DecodeTokenAmountDestExecData(destExecData, sourceChainSelector)
}

type NoOpSourceChainExtraDataCodec struct{}

func (n NoOpSourceChainExtraDataCodec) DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error) {
	return make(map[string]any), nil
}

func (n NoOpSourceChainExtraDataCodec) DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error) {
	return make(map[string]any), nil
}
