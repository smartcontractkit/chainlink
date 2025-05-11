package metadata

import (
	"encoding/json"
	"errors"
	"fmt"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

type DataStreamsMutableDataStore = ds.MutableDataStore[SerializedContractMetadata, ds.DefaultMetadata]

// SerializedContractMetadata provides a generic container for contract metadata
// that can be serialized/deserialized to/from JSON
type SerializedContractMetadata struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

// Clone creates a copy of the SerializedContractMetadata
func (s SerializedContractMetadata) Clone() SerializedContractMetadata {
	contentCopy := make([]byte, len(s.Content))
	copy(contentCopy, s.Content)

	return SerializedContractMetadata{
		Type:    s.Type,
		Content: contentCopy,
	}
}

// NewSerializedContractMetadata serializes any contract metadata
func NewSerializedContractMetadata[T any](metadata GenericContractMetadata[T]) (*SerializedContractMetadata, error) {
	content, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return &SerializedContractMetadata{
		Content: content,
	}, nil
}

// DeserializeMetadata converts a SerializedContractMetadata to a GenericContractMetadata of type T
func DeserializeMetadata[T any](serialized SerializedContractMetadata) (*GenericContractMetadata[T], error) {
	if len(serialized.Content) == 0 {
		return nil, errors.New("empty content in serialized metadata")
	}

	var result GenericContractMetadata[T]
	if err := json.Unmarshal(serialized.Content, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contract metadata: %w", err)
	}

	return &result, nil
}

type CommonContractMetadata struct {
	DeployBlock uint64 `json:"deployBlock"`
}

// GenericContractMetadata is a generic container for any view type
// Use as Content for SerializedContractMetadata
type GenericContractMetadata[T any] struct {
	Metadata CommonContractMetadata `json:"metadata"`
	// View is intended to be populated with the contract's view usually after state change to have an off chain representation of the contract
	// Tip: The view generated should only depend on the contract state retrieved on-chain and not any external information.
	// This will keep the view generation deterministic and consistent. Any external information not available in the contract
	// but useful should be stored in some other way which makes re-generating a view easy
	View T `json:"view"`
}
