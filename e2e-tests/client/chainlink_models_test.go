package client

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func TestOCR2TaskJobSpec_String(t *testing.T) {
	for _, tt := range []struct {
		name string
		spec OCR2TaskJobSpec
		exp  string
	}{
		{
			name: "chain-reader-codec",
			spec: OCR2TaskJobSpec{
				OCR2OracleSpec: job.OCR2OracleSpec{
					RelayConfig: map[string]interface{}{
						"chainID":   1337,
						"fromBlock": 42,
						"chainReader": evmtypes.ChainReaderConfig{
							Contracts: map[string]evmtypes.ChainContractReader{
								"median": {
									ContractABI: `[
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "requester",
        "type": "address"
      }
    ],
    "name": "RoundRequested",
    "type": "event"
  }
]
`,
									Configs: map[string]*evmtypes.ChainReaderDefinition{
										"LatestTransmissionDetails": {
											ChainSpecificName: "latestTransmissionDetails",
											OutputModifications: codec.ModifiersConfig{
												&codec.EpochToTimeModifierConfig{
													Fields: []string{"LatestTimestamp_"},
												},
												&codec.RenameModifierConfig{
													Fields: map[string]string{
														"LatestAnswer_":    "LatestAnswer",
														"LatestTimestamp_": "LatestTimestamp",
													},
												},
											},
										},
										"LatestRoundRequested": {
											ChainSpecificName: "RoundRequested",
											ReadType:          evmtypes.Event,
										},
									},
								},
							},
						},
						"codec": evmtypes.CodecConfig{
							Configs: map[string]evmtypes.ChainCodecConfig{
								"MedianReport": {
									TypeABI: `[
  {
    "Name": "Timestamp",
    "Type": "uint32"
  }
]
`,
								},
							},
						},
					},
					PluginConfig: map[string]interface{}{"juelsPerFeeCoinSource": `		// data source 1
		ds1          [type=bridge name="%s"];
		ds1_parse    [type=jsonparse path="data"];
		ds1_multiply [type=multiply times=2];

		// data source 2
		ds2          [type=http method=GET url="%s"];
		ds2_parse    [type=jsonparse path="data"];
		ds2_multiply [type=multiply times=2];

		ds1 -> ds1_parse -> ds1_multiply -> answer1;
		ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median index=0];
`,
					},
				},
			},
			exp: `
type                                   = ""
name                                   = ""
forwardingAllowed                      = false
relay                                  = ""
schemaVersion                          = 1
contractID                             = ""

[relayConfig]
chainID = 1337
fromBlock = 42

[relayConfig.chainReader]
[relayConfig.chainReader.contracts]
[relayConfig.chainReader.contracts.median]
contractABI = "[\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": true,\n        \"internalType\": \"address\",\n        \"name\": \"requester\",\n        \"type\": \"address\"\n      }\n    ],\n    \"name\": \"RoundRequested\",\n    \"type\": \"event\"\n  }\n]\n"

[relayConfig.chainReader.contracts.median.configs]
LatestRoundRequested = "{\n  \"chainSpecificName\": \"RoundRequested\",\n  \"readType\": \"event\"\n}\n"
LatestTransmissionDetails = "{\n  \"chainSpecificName\": \"latestTransmissionDetails\",\n  \"output_modifications\": [\n    {\n      \"Fields\": [\n        \"LatestTimestamp_\"\n      ],\n      \"Type\": \"epoch to time\"\n    },\n    {\n      \"Fields\": {\n        \"LatestAnswer_\": \"LatestAnswer\",\n        \"LatestTimestamp_\": \"LatestTimestamp\"\n      },\n      \"Type\": \"rename\"\n    }\n  ]\n}\n"

[relayConfig.codec]
[relayConfig.codec.configs]
[relayConfig.codec.configs.MedianReport]
typeABI = "[\n  {\n    \"Name\": \"Timestamp\",\n    \"Type\": \"uint32\"\n  }\n]\n"

`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.spec.String()
			require.NoError(t, err)
			require.Equal(t, tt.exp, got)
		})
	}
}
