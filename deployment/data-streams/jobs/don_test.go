package jobs

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
)

const donSpecTOML1 = `name = 'Test-DON'
type = 'don'
schemaVersion = 1
externalJobID = '00000000-0000-0000-0000-000000000000'
contractID = 'contract-123'
transmitterID = 'tx-123'
forwardingAllowed = true
p2pv2Bootstrappers = ['bootstrap1', 'bootstrap2']
ocrKeyBundleID = 'ocr-bundle-123'
maxTaskDuration = 10000000000
contractConfigTrackerPollInterval = 60000000000
relay = 'testrelay'
pluginType = 'testplugin'

[relayConfig]
chainID = 'chain'
fromBlock = 100
lloConfigMode = 'mode'
lloDonID = 200

[pluginConfig]
channelDefinitionsContractAddress = '0xabc'
channelDefinitionsContractFromBlock = 50
donID = 300
servers = {server1 = 'http://localhost'}
`

const donSpecTOML2 = `name = 'Empty-DON-Test'
type = 'don'
schemaVersion = 1
externalJobID = '00000000-0000-0000-0000-000000000000'
contractID = 'contract-empty'

[relayConfig]
chainID = ''

[pluginConfig]
channelDefinitionsContractAddress = ''
channelDefinitionsContractFromBlock = 0
donID = 0
servers = {}
`

func TestDonJobSpec_MarshalTOML(t *testing.T) {
	testCases := []struct {
		name string
		spec DonJobSpec
		want string
	}{
		{
			name: "with fields populated",
			spec: DonJobSpec{
				Base: Base{
					Name:          "Test-DON",
					Type:          "don",
					SchemaVersion: 1,
					ExternalJobID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				},
				ContractID:                        "contract-123",
				TransmitterID:                     "tx-123",
				ForwardingAllowed:                 pointer.To(true),
				P2PV2Bootstrappers:                []string{"bootstrap1", "bootstrap2"},
				OCRKeyBundleID:                    pointer.To("ocr-bundle-123"),
				MaxTaskDuration:                   10 * time.Second,
				ContractConfigTrackerPollInterval: 1 * time.Minute,
				Relay:                             "testrelay",
				PluginType:                        "testplugin",
				RelayConfig: RelayConfigDon{
					ChainID:       "chain",
					FromBlock:     100,
					LLOConfigMode: "mode",
					LLODonID:      200,
				},
				PluginConfig: PluginConfigDon{
					ChannelDefinitionsContractAddress:   "0xabc",
					ChannelDefinitionsContractFromBlock: 50,
					DonID:                               300,
					Servers:                             map[string]string{"server1": "http://localhost"},
				},
			},
			want: donSpecTOML1,
		},
		{
			name: "empty minimal fields",
			spec: DonJobSpec{
				Base: Base{
					Name:          "Empty-DON-Test",
					Type:          "don",
					SchemaVersion: 1,
					ExternalJobID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				},
				ContractID:  "contract-empty",
				RelayConfig: RelayConfigDon{},
				PluginConfig: PluginConfigDon{
					Servers: map[string]string{},
				},
			},
			want: donSpecTOML2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tomlBytes, err := tc.spec.MarshalTOML()
			require.NoError(t, err)
			got := string(tomlBytes)
			require.Equal(t, tc.want, got)
		})
	}
}
