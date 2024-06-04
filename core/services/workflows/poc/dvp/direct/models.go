package direct

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

// These would be generated from the capability's json schema, for now
// I'm assuming that we'll use an HTTP trigger from SWIFT, but we could make a trigger for them

type HttpTrigger struct {
	Body []byte
}

func NewHttpTrigger(ref, typeName string) *capabilities.RemoteTrigger[*HttpTrigger] {
	return &capabilities.RemoteTrigger[*HttpTrigger]{
		// TODO this would be what we can use to distinguish between different triggers
		// to allow data normalization
		RefName:  ref,
		TypeName: typeName,
	}
}

type KvStoreSetValuesRequest struct {
	Values map[string]any
}

type KvStoreGetValuesRequest struct {
	Keys []string
}

type KvStoreGetValuesResponse struct {
	Values map[string]any
}

type KeyValueStore interface {
	AddSetValuesTarget(ref string, wb *workflow.Builder[capabilities.ConsensusResult[*KvStoreSetValuesRequest]]) error
	AddGetValues(ref string, wb *workflow.Builder[*KvStoreGetValuesRequest]) (*workflow.Builder[*KvStoreGetValuesResponse], error)
}

func NewKeyValueStore(typeName string) KeyValueStore {
	return &keyValueStore{typeName: typeName}
}

type keyValueStore struct {
	typeName string
}

func (k *keyValueStore) AddSetValuesTarget(ref string, wb *workflow.Builder[capabilities.ConsensusResult[*KvStoreSetValuesRequest]]) error {
	return workflow.AddTarget[*KvStoreSetValuesRequest](wb, &capabilities.RemoteTarget[*KvStoreSetValuesRequest]{
		RefName:  ref,
		TypeName: k.typeName,
	})
}

func (k *keyValueStore) AddGetValues(ref string, wb *workflow.Builder[*KvStoreGetValuesRequest]) (*workflow.Builder[*KvStoreGetValuesResponse], error) {
	return workflow.AddStep[*KvStoreGetValuesRequest, *KvStoreGetValuesResponse](wb, &capabilities.RemoteAction[*KvStoreGetValuesRequest, *KvStoreGetValuesResponse]{
		RefName:  ref,
		TypeName: k.typeName,
	})
}

type ChainWriteRequest struct {
	Body []byte
}

type ChainWriter interface {
	AddWriteTarget(ref string, wb *workflow.Builder[capabilities.ConsensusResult[*ChainWriteRequest]]) error
}

func NewChainWriter(typeName string) ChainWriter {
	return &chainWriter{typeName: typeName}
}

type chainWriter struct {
	typeName string
}

func (c *chainWriter) AddWriteTarget(ref string, wb *workflow.Builder[capabilities.ConsensusResult[*ChainWriteRequest]]) error {
	return workflow.AddTarget[*ChainWriteRequest](wb, &capabilities.RemoteTarget[*ChainWriteRequest]{
		RefName:  ref,
		TypeName: c.typeName,
	})
}

type EncodeRequest struct {
	// There's actually versioned bytes that we would use, but keep it simple for now
	CborBody []byte
}

type EncodeResponse struct {
	Encoded []byte
}

type Codec interface {
	AddEncode(ref string, wb *workflow.Builder[*EncodeRequest]) (*workflow.Builder[*EncodeResponse], error)
	// Other methods would exist too.
}

func NewCodec(typeName string) Codec {
	return &codec{typeName: typeName}
}

type codec struct {
	typeName string
}

func (c *codec) AddEncode(ref string, wb *workflow.Builder[*EncodeRequest]) (*workflow.Builder[*EncodeResponse], error) {
	return workflow.AddStep[*EncodeRequest, *EncodeResponse](wb, &capabilities.RemoteAction[*EncodeRequest, *EncodeResponse]{
		RefName:  ref,
		TypeName: c.typeName,
	})
}
