package read

import (
	"github.com/ethereum/go-ethereum/common"

	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
)

func TestCreateQueryName(t *testing.T) {
	eventQueries := []EventQuery{
		{
			Filter: query.KeyFilter{Key: "key1"},
			EventBinding: &EventBinding{
				contractName: "ContractA",
				eventName:    "EventA",
				hash:         common.HexToHash("0x1"),
			},
			SequenceDataType: "dataType1",
			Address:          common.HexToAddress("0x123"),
		},
		{
			Filter: query.KeyFilter{Key: "key2"},
			EventBinding: &EventBinding{
				contractName: "ContractB",
				eventName:    "EventB",
				hash:         common.HexToHash("0x2"),
			},
			SequenceDataType: "dataType2",
			Address:          common.HexToAddress("0x456"),
		},
		{
			Filter: query.KeyFilter{Key: "key1"},
			EventBinding: &EventBinding{
				contractName: "ContractA",
				eventName:    "EventA1",
				hash:         common.HexToHash("0x1"),
			},
			SequenceDataType: "dataType1",
			Address:          common.HexToAddress("0x123"),
		},
	}

	expectedQueryName := "ContractA-0x0000000000000000000000000000000000000123-EventA-EventA1-ContractB-0x0000000000000000000000000000000000000456-EventB"
	queryName := createQueryName(eventQueries)

	assert.Equal(t, expectedQueryName, queryName)
}

func TestValidateEventQueries(t *testing.T) {
	tests := []struct {
		name          string
		eventQueries  []EventQuery
		expectedError string
	}{
		{
			name: "valid event queries",
			eventQueries: []EventQuery{
				{
					EventBinding:     &EventBinding{hash: common.HexToHash("0x1")},
					SequenceDataType: "dataType1",
				},
				{
					EventBinding:     &EventBinding{hash: common.HexToHash("0x2")},
					SequenceDataType: "dataType2",
				},
			},
			expectedError: "",
		},
		{
			name: "nil event binding",
			eventQueries: []EventQuery{
				{
					EventBinding:     nil,
					SequenceDataType: "dataType1",
				},
			},
			expectedError: "event binding is nil",
		},
		{
			name: "nil sequence data type",
			eventQueries: []EventQuery{
				{
					EventBinding:     &EventBinding{hash: common.HexToHash("0x1")},
					SequenceDataType: nil,
				},
			},
			expectedError: "sequence data type is nil",
		},
		{
			name: "duplicate event query",
			eventQueries: []EventQuery{
				{
					EventBinding:     &EventBinding{hash: common.HexToHash("0x1")},
					SequenceDataType: "dataType1",
				},
				{
					EventBinding:     &EventBinding{hash: common.HexToHash("0x1")},
					SequenceDataType: "dataType2",
				},
			},
			expectedError: "duplicate event query for event signature 0x0000000000000000000000000000000000000000000000000000000000000001, event name ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEventQueries(tt.eventQueries)
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
