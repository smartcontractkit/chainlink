package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/deployment"
)

func TestAddressRefByChainSelector(t *testing.T) {
	var (
		recordOne = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 2,
			Type:          "typeX",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressRef
		giveChain      uint64
		expectedResult []AddressRef
	}{
		{
			name: "success: returns record with given chain",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveChain:      2,
			expectedResult: []AddressRef{recordTwo},
		},
		{
			name: "success: returns no record with given chain",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveChain:      5,
			expectedResult: []AddressRef{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := AddressRefByChainSelector(tt.giveChain)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}

func TestAddressRefByType(t *testing.T) {
	var (
		recordOne = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 2,
			Type:          "typeX",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressRef
		giveType       deployment.ContractType
		expectedResult []AddressRef
	}{
		{
			name: "success: returns record with given type",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveType: "typeX",
			expectedResult: []AddressRef{
				recordTwo,
			},
		},
		{
			name: "success: returns no record with given type",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveType:       "typeL",
			expectedResult: []AddressRef{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := AddressRefByType(tt.giveType)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}

func TestAddressRefByVersion(t *testing.T) {
	var (
		recordOne = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 2,
			Type:          "typeX",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressRef
		giveVersion    *semver.Version
		expectedResult []AddressRef
	}{
		{
			name: "success: returns record with given version",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveVersion: semver.MustParse("0.5.0"),
			expectedResult: []AddressRef{
				recordOne,
				recordTwo,
			},
		},
		{
			name: "success: returns no record with given version",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveVersion:    semver.MustParse("0.6.0"),
			expectedResult: []AddressRef{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := AddressRefByVersion(tt.giveVersion)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}

func TestAddressRefByQualifier(t *testing.T) {
	var (
		recordOne = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 1,
			Type:          "type1",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual1",
			Labels: deployment.NewLabelSet(
				"label1", "label2", "label3",
			),
		}

		recordTwo = AddressRef{
			Address:       "0x2324224",
			ChainSelector: 2,
			Type:          "typeX",
			Version:       semver.MustParse("0.5.0"),
			Qualifier:     "qual2",
			Labels: deployment.NewLabelSet(
				"label13", "label23", "label33",
			),
		}
	)

	tests := []struct {
		name           string
		givenState     []AddressRef
		giveQualifier  string
		expectedResult []AddressRef
	}{
		{
			name: "success: returns record with given qualifier",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveQualifier: "qual1",
			expectedResult: []AddressRef{
				recordOne,
			},
		},
		{
			name: "success: returns no record with given qualifier",
			givenState: []AddressRef{
				recordOne,
				recordTwo,
			},
			giveQualifier:  "qual32",
			expectedResult: []AddressRef{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := AddressRefByQualifier(tt.giveQualifier)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}

func TestContractMetadataByChainSelector(t *testing.T) {
	var (
		recordOne = ContractMetadata{
			ChainSelector: 1,
			Address:       "0x2324224",
			Metadata:      "metadata1",
		}

		recordTwo = ContractMetadata{
			ChainSelector: 2,
			Address:       "0x2324224",
			Metadata:      "metadata2",
		}
	)

	tests := []struct {
		name           string
		givenState     []ContractMetadata
		giveChain      uint64
		expectedResult []ContractMetadata
	}{
		{
			name: "success: returns record with given chain",
			givenState: []ContractMetadata{
				recordOne,
				recordTwo,
			},
			giveChain:      2,
			expectedResult: []ContractMetadata{recordTwo},
		},
		{
			name: "success: returns no record with given chain",
			givenState: []ContractMetadata{
				recordOne,
				recordTwo,
			},
			giveChain:      5,
			expectedResult: []ContractMetadata(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := ContractMetadataByChainSelector(tt.giveChain)
			filteredRecords := filter(tt.givenState)
			assert.Equal(t, tt.expectedResult, filteredRecords)
		})
	}
}
