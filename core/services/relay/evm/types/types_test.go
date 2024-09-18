package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// ChainID   *big.Big   `json:"chainID"`
// FromBlock uint64     `json:"fromBlock"`

// // Contract-specific
// EffectiveTransmitterAddress null.String    `json:"effectiveTransmitterAddress"`
// SendingKeys                 pq.StringArray `json:"sendingKeys"`

// // Mercury-specific
// FeedID *common.Hash `json:"feedID"`
func Test_RelayConfig(t *testing.T) {
	cid := testutils.NewRandomEVMChainID()
	fromBlock := uint64(2222)
	feedID := utils.NewHash()
	rawToml := fmt.Sprintf(`
ChainID = "%s"
FromBlock = %d
FeedID = "0x%x"
`, cid, fromBlock, feedID[:])

	var rc RelayConfig
	err := toml.Unmarshal([]byte(rawToml), &rc)
	require.NoError(t, err)

	assert.Equal(t, cid.String(), rc.ChainID.String())
	assert.Equal(t, fromBlock, rc.FromBlock)
	assert.Equal(t, feedID.Hex(), rc.FeedID.Hex())
}

func Test_ChainReaderConfig(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		expected  ChainReaderConfig
	}{
		{
			name: "Valid JSON",
			jsonInput: `
{
   "contracts":{
      "Contract1":{
		 "contractABI":"[  {    \"anonymous\": false,    \"inputs\": [      {        \"indexed\": true,        \"internalType\": \"address\",        \"name\": \"requester\",        \"type\": \"address\"      },      {        \"indexed\": false,        \"internalType\": \"bytes32\",        \"name\": \"configDigest\",        \"type\": \"bytes32\"      },      {        \"indexed\": false,        \"internalType\": \"uint32\",        \"name\": \"epoch\",        \"type\": \"uint32\"      },      {        \"indexed\": false,        \"internalType\": \"uint8\",        \"name\": \"round\",        \"type\": \"uint8\"      }    ],    \"name\": \"RoundRequested\",    \"type\": \"event\"  },  {    \"inputs\": [],    \"name\": \"latestTransmissionDetails\",    \"outputs\": [      {        \"internalType\": \"bytes32\",        \"name\": \"configDigest\",        \"type\": \"bytes32\"      },      {        \"internalType\": \"uint32\",        \"name\": \"epoch\",        \"type\": \"uint32\"      },      {        \"internalType\": \"uint8\",        \"name\": \"round\",        \"type\": \"uint8\"      },      {        \"internalType\": \"int192\",        \"name\": \"latestAnswer_\",        \"type\": \"int192\"      },      {        \"internalType\": \"uint64\",        \"name\": \"latestTimestamp_\",        \"type\": \"uint64\"      }    ],    \"stateMutability\": \"view\",    \"type\": \"function\"  }]",
         "contractPollingFilter":{
            "genericEventNames":[
               "event1",
               "event2"
            ],
            "pollingFilter":{
               "topic2":[
                  "0x1abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52"
               ],
               "topic3":[
                  "0x2abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52"
               ],
               "topic4":[
                  "0x3abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52"
               ],
               "retention":"1m0s",
               "maxLogsKept":100,
               "logsPerBlock":10
            }
         },
         "configs":{
            "config1":"{\"cacheEnabled\":true,\"chainSpecificName\":\"specificName1\",\"inputModifications\":[{\"Fields\":[\"ts\"],\"Type\":\"epoch to time\"},{\"Fields\":{\"a\":\"b\"},\"Type\":\"rename\"}],\"outputModifications\":[{\"Fields\":[\"ts\"],\"Type\":\"epoch to time\"},{\"Fields\":{\"c\":\"d\"},\"Type\":\"rename\"}],\"eventDefinitions\":{\"genericTopicNames\":{\"TopicKey1\":\"TopicVal1\"},\"genericDataWordDefs\":{\"DataWordKey\": \"DataWordKey\"},\"pollingFilter\":{\"topic2\":[\"0x4abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52\"],\"topic3\":[\"0x5abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52\"],\"topic4\":[\"0x6abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52\"],\"retention\":\"1m0s\",\"maxLogsKept\":100,\"logsPerBlock\":10}},\"confidenceConfirmations\":{\"0.0\":0,\"1.0\":-1}}"
         }
      }
   }
}
`, expected: ChainReaderConfig{
				Contracts: map[string]ChainContractReader{
					"Contract1": {
						ContractABI: "[  {    \"anonymous\": false,    \"inputs\": [      {        \"indexed\": true,        \"internalType\": \"address\",        \"name\": \"requester\",        \"type\": \"address\"      },      {        \"indexed\": false,        \"internalType\": \"bytes32\",        \"name\": \"configDigest\",        \"type\": \"bytes32\"      },      {        \"indexed\": false,        \"internalType\": \"uint32\",        \"name\": \"epoch\",        \"type\": \"uint32\"      },      {        \"indexed\": false,        \"internalType\": \"uint8\",        \"name\": \"round\",        \"type\": \"uint8\"      }    ],    \"name\": \"RoundRequested\",    \"type\": \"event\"  },  {    \"inputs\": [],    \"name\": \"latestTransmissionDetails\",    \"outputs\": [      {        \"internalType\": \"bytes32\",        \"name\": \"configDigest\",        \"type\": \"bytes32\"      },      {        \"internalType\": \"uint32\",        \"name\": \"epoch\",        \"type\": \"uint32\"      },      {        \"internalType\": \"uint8\",        \"name\": \"round\",        \"type\": \"uint8\"      },      {        \"internalType\": \"int192\",        \"name\": \"latestAnswer_\",        \"type\": \"int192\"      },      {        \"internalType\": \"uint64\",        \"name\": \"latestTimestamp_\",        \"type\": \"uint64\"      }    ],    \"stateMutability\": \"view\",    \"type\": \"function\"  }]",
						ContractPollingFilter: ContractPollingFilter{
							GenericEventNames: []string{"event1", "event2"},
							PollingFilter: PollingFilter{
								Topic2:       evmtypes.HashArray{common.HexToHash("0x1abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52")},
								Topic3:       evmtypes.HashArray{common.HexToHash("0x2abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52")},
								Topic4:       evmtypes.HashArray{common.HexToHash("0x3abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52")},
								Retention:    models.Interval(time.Minute * 1),
								MaxLogsKept:  100,
								LogsPerBlock: 10,
							},
						},
						Configs: map[string]*ChainReaderDefinition{
							"config1": {
								CacheEnabled:      true,
								ChainSpecificName: "specificName1",
								ReadType:          Method,
								InputModifications: codec.ModifiersConfig{
									&codec.EpochToTimeModifierConfig{
										Fields: []string{"ts"},
									},
									&codec.RenameModifierConfig{
										Fields: map[string]string{
											"a": "b",
										},
									},
								},
								OutputModifications: codec.ModifiersConfig{
									&codec.EpochToTimeModifierConfig{
										Fields: []string{"ts"},
									},
									&codec.RenameModifierConfig{
										Fields: map[string]string{
											"c": "d",
										},
									},
								},
								ConfidenceConfirmations: map[string]int{"0.0": 0, "1.0": -1},
								EventDefinitions: &EventDefinitions{
									GenericTopicNames:    map[string]string{"TopicKey1": "TopicVal1"},
									GenericDataWordNames: map[string]string{"DataWordKey": "DataWordKey"},
									PollingFilter: &PollingFilter{
										Topic2:       evmtypes.HashArray{common.HexToHash("0x4abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52")},
										Topic3:       evmtypes.HashArray{common.HexToHash("0x5abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52")},
										Topic4:       evmtypes.HashArray{common.HexToHash("0x6abbe4784b1fb071039bb9cb50b82978fb5d3ab98fb512c032e75786b93e2c52")},
										Retention:    models.Interval(time.Minute * 1),
										MaxLogsKept:  100,
										LogsPerBlock: 10,
									},
								},
							},
						},
					},
				},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config ChainReaderConfig
			config.Contracts = make(map[string]ChainContractReader)
			require.Nil(t, json.Unmarshal([]byte(tt.jsonInput), &config))
			require.Equal(t, tt.expected, config)
		})
	}
}
