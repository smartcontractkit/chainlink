package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ToDefault is a utility function that converts a DataStore with domain specific
// metadata types into a DataStore with DefaultMetadata types.
// NOTE: It is assumed that the domain specific metadata types are JSON serializable.
func ToDefault[CM Cloneable[CM], EM Cloneable[EM]](
	dataStore DataStore[CM, EM],
) (MutableDataStore[DefaultMetadata, DefaultMetadata], error) {
	converted := NewMemoryDataStore[DefaultMetadata, DefaultMetadata]()

	// Copy all addressRef over to the new data store, no conversion is needed
	addressRefs, err := dataStore.Addresses().Fetch()
	if err != nil {
		return nil, fmt.Errorf("error fetching AddressRefs: %w", err)
	}

	for _, ar := range addressRefs {
		err := converted.Addresses().Add(ar)
		if err != nil {
			return nil, fmt.Errorf("error adding AddressRef: %w", err)
		}
	}

	// Copy all contractMetadata over to the new data store and convert the metadata
	// to a JSON string. This is done by marshaling the metadata into a JSON string.
	contractMetadata, err := dataStore.ContractMetadata().Fetch()
	if err != nil {
		return nil, fmt.Errorf("error fetching ContractMetadata: %w", err)
	}

	for _, cm := range contractMetadata {
		jsonData, err := json.Marshal(cm.Metadata)
		if err != nil {
			return nil, fmt.Errorf("error marshaling ContractMetadata: %w", err)
		}

		err = converted.ContractMetadata().Add(ContractMetadata[DefaultMetadata]{
			ChainSelector: cm.ChainSelector,
			Address:       cm.Address,
			Metadata: DefaultMetadata{
				Data: string(jsonData),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error adding ContractMetadata: %w", err)
		}
	}

	// Fetch the EnvMetadata and check if it was set.
	envMetadata, err := dataStore.EnvMetadata().Get()
	if err != nil {
		if errors.Is(err, ErrEnvMetadataNotSet) {
			// If the env metadata was not set, Get() will return ErrEnvMetadataNotSet.
			// In this case, we don't need to do anything.
			return converted, nil
		}
		return nil, err
	}

	// Convert the EnvMetadata to a JSON string. This is done by marshaling the metadata
	// into a JSON string.
	jsonData, err := json.Marshal(envMetadata.Metadata)
	if err != nil {
		return nil, fmt.Errorf("error marshaling EnvMetadata: %w", err)
	}

	// Set the EnvMetadata in the new data store with the JSON string.
	err = converted.EnvMetadata().Set(EnvMetadata[DefaultMetadata]{
		Domain:      envMetadata.Domain,
		Environment: envMetadata.Environment,
		Metadata: DefaultMetadata{
			Data: string(jsonData),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error updating EnvMetadata: %w", err)
	}

	return converted, nil
}

// FromDefault is a utility function that converts a DataStore with DefaultMetadata types
// into a DataStore with domain specific metadata types.
// NOTE: It is assumed that the domain specific metadata types are JSON deserializable.
func FromDefault[CM Cloneable[CM], EM Cloneable[EM]](
	defaultStore DataStore[DefaultMetadata, DefaultMetadata],
) (DataStore[CM, EM], error) {
	converted := NewMemoryDataStore[CM, EM]()

	// Copy all addressRef over to the new data store, no conversion is needed
	addressRefs, err := defaultStore.Addresses().Fetch()
	if err != nil {
		return nil, fmt.Errorf("error fetching AddressRefs: %w", err)
	}

	for _, ar := range addressRefs {
		err := converted.Addresses().Add(ar)
		if err != nil {
			return nil, fmt.Errorf("error adding AddressRef: %w", err)
		}
	}

	// Copy all contractMetadata over to the new data store and convert the metadata
	// to to the domain specific type. This is done by unmarshaling the JSON string
	// representing the metadata into the concrete type.
	contractMetadata, err := defaultStore.ContractMetadata().Fetch()
	if err != nil {
		return nil, fmt.Errorf("error fetching ContractMetadata: %w", err)
	}

	for _, cm := range contractMetadata {
		var metadata CM
		err := json.Unmarshal([]byte(cm.Metadata.Data), &metadata)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling ContractMetadata: %w", err)
		}

		err = converted.ContractMetadata().Add(ContractMetadata[CM]{
			ChainSelector: cm.ChainSelector,
			Address:       cm.Address,
			Metadata:      metadata,
		})
		if err != nil {
			return nil, fmt.Errorf("error adding ContractMetadata: %w", err)
		}
	}

	// Fetch the EnvMetadata and check if it was set.
	envMetadata, err := defaultStore.EnvMetadata().Get()
	if err != nil {
		if errors.Is(err, ErrEnvMetadataNotSet) {
			// If the env metadata was not set, Get() will return ErrEnvMetadataNotSet.
			// In this case, we don't need to do anything.
			return converted.Seal(), nil
		}
		return nil, err
	}

	// Convert the EnvMetadata to the domain specific type. This is done by unmarshaling
	// the JSON string representing the metadata into the concrete type.
	var metadata EM
	err = json.Unmarshal([]byte(envMetadata.Metadata.Data), &metadata)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling EnvMetadata: %w", err)
	}

	// Set the EnvMetadata in the new data store with the domain specific type.
	err = converted.EnvMetadata().Set(EnvMetadata[EM]{
		Domain:      envMetadata.Domain,
		Environment: envMetadata.Environment,
		Metadata:    metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("error updating EnvMetadata: %w", err)
	}

	return converted.Seal(), nil
}
