package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

// CustomMetadata is a placeholder type for testing purposes.
type CustomMetadata struct {
	Field string `json:"field"`
}

// Clone creates a deep copy of CustomMetadata.
func (cm CustomMetadata) Clone() CustomMetadata {
	return CustomMetadata{
		Field: cm.Field,
	}
}

func TestToDefault(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() MutableDataStore[CustomMetadata, CustomMetadata]
		expected MutableDataStore[DefaultMetadata, DefaultMetadata]
	}{
		{
			name: "successful conversion",
			setup: func() MutableDataStore[CustomMetadata, CustomMetadata] {
				ds := NewMemoryDataStore[CustomMetadata, CustomMetadata]()

				err := ds.Addresses().Add(AddressRef{
					Address:       "addr1",
					Type:          "type1",
					Version:       semver.MustParse("1.0.0"),
					ChainSelector: 1,
					Qualifier:     "qualifier1",
					Labels:        NewLabelSet("label1", "label2"),
				})
				require.NoError(t, err)

				err = ds.ContractMetadata().Add(ContractMetadata[CustomMetadata]{
					ChainSelector: 1,
					Address:       "contract1",
					Metadata:      CustomMetadata{Field: "value1"},
				})
				require.NoError(t, err)

				err = ds.EnvMetadata().Set(EnvMetadata[CustomMetadata]{
					Domain:      "domain1",
					Environment: "env1",
					Metadata:    CustomMetadata{Field: "envValue1"},
				})
				require.NoError(t, err)

				return ds
			},
			expected: &MemoryDataStore[DefaultMetadata, DefaultMetadata]{
				AddressRefStore: &MemoryAddressRefStore{
					Records: []AddressRef{
						{
							Address:       "addr1",
							Type:          "type1",
							Version:       semver.MustParse("1.0.0"),
							ChainSelector: 1,
							Qualifier:     "qualifier1",
							Labels:        NewLabelSet("label1", "label2"),
						},
					},
				},
				ContractMetadataStore: &MemoryContractMetadataStore[DefaultMetadata]{
					Records: []ContractMetadata[DefaultMetadata]{
						{
							ChainSelector: 1,
							Address:       "contract1",
							Metadata: DefaultMetadata{
								Data: `{"field":"value1"}`,
							},
						},
					},
				},
				EnvMetadataStore: &MemoryEnvMetadataStore[DefaultMetadata]{
					Record: &EnvMetadata[DefaultMetadata]{
						Domain:      "domain1",
						Environment: "env1",
						Metadata: DefaultMetadata{
							Data: `{"field":"envValue1"}`,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataStore := tt.setup()

			// Test ToDefault
			defaultStore, err := ToDefault(dataStore.Seal())
			require.NoError(t, err)

			require.Equal(t, tt.expected, defaultStore)

		})
	}
}

func TestFromDefault(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() MutableDataStore[DefaultMetadata, DefaultMetadata]
		expected DataStore[CustomMetadata, CustomMetadata]
	}{
		{
			name: "successful conversion",
			setup: func() MutableDataStore[DefaultMetadata, DefaultMetadata] {
				ds := NewMemoryDataStore[DefaultMetadata, DefaultMetadata]()

				err := ds.Addresses().Add(AddressRef{
					Address:       "addr1",
					Type:          "type1",
					Version:       semver.MustParse("1.0.0"),
					ChainSelector: 1,
					Qualifier:     "qualifier1",
					Labels:        NewLabelSet("label1", "label2"),
				})
				require.NoError(t, err)

				err = ds.ContractMetadata().Add(ContractMetadata[DefaultMetadata]{
					ChainSelector: 1,
					Address:       "contract1",
					Metadata: DefaultMetadata{
						Data: `{"field":"value1"}`,
					},
				})
				require.NoError(t, err)

				err = ds.EnvMetadata().Set(EnvMetadata[DefaultMetadata]{
					Domain:      "domain1",
					Environment: "env1",
					Metadata: DefaultMetadata{
						Data: `{"field":"envValue1"}`,
					},
				})
				require.NoError(t, err)

				return ds
			},
			expected: &sealedMemoryDataStore[CustomMetadata, CustomMetadata]{
				AddressRefStore: &MemoryAddressRefStore{
					Records: []AddressRef{
						{
							Address:       "addr1",
							Type:          "type1",
							Version:       semver.MustParse("1.0.0"),
							ChainSelector: 1,
							Qualifier:     "qualifier1",
							Labels:        NewLabelSet("label1", "label2"),
						},
					},
				},
				ContractMetadataStore: &MemoryContractMetadataStore[CustomMetadata]{
					Records: []ContractMetadata[CustomMetadata]{
						{
							ChainSelector: 1,
							Address:       "contract1",
							Metadata:      CustomMetadata{Field: "value1"},
						},
					},
				},
				EnvMetadataStore: &MemoryEnvMetadataStore[CustomMetadata]{
					Record: &EnvMetadata[CustomMetadata]{
						Domain:      "domain1",
						Environment: "env1",
						Metadata:    CustomMetadata{Field: "envValue1"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataStore := tt.setup()

			// Test FromDefault
			customStore, err := FromDefault[CustomMetadata, CustomMetadata](dataStore.Seal())
			require.NoError(t, err)

			require.Equal(t, tt.expected, customStore)
		})
	}
}
