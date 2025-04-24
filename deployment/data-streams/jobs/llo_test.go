package jobs

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
)

const lloSpecTOML1 = `name = 'Test-DON'
type = 'offchainreporting2'
schemaVersion = 1
externalJobID = '00000000-0000-0000-0000-000000000000'
contractID = 'contract-123'
transmitterID = 'tx-123'
forwardingAllowed = true
p2pv2Bootstrappers = [' PeerID@Host1:Port1/Host2:Port2', 'PeerID@Host1:Port1/Host2:Port2']
ocrKeyBundleID = 'ocr-bundle-123'
maxTaskDuration = '10s'
contractConfigTrackerPollInterval = '60s'
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

const lloSpecTOML2 = `name = 'Empty-DON-Test'
type = 'offchainreporting2'
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

func TestLLOJobSpec_MarshalTOML(t *testing.T) {
	testCases := []struct {
		name string
		spec LLOJobSpec
		want string
	}{
		{
			name: "with fields populated",
			spec: LLOJobSpec{
				Base: Base{
					Name:          "Test-DON",
					Type:          "offchainreporting2",
					SchemaVersion: 1,
					ExternalJobID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				},
				ContractID:                        "contract-123",
				TransmitterID:                     "tx-123",
				ForwardingAllowed:                 pointer.To(true),
				P2PV2Bootstrappers:                []string{" PeerID@Host1:Port1/Host2:Port2", "PeerID@Host1:Port1/Host2:Port2"},
				OCRKeyBundleID:                    pointer.To("ocr-bundle-123"),
				MaxTaskDuration:                   TOMLDuration(10 * time.Second),
				ContractConfigTrackerPollInterval: TOMLDuration(60 * time.Second),
				Relay:                             "testrelay",
				PluginType:                        "testplugin",
				RelayConfig: RelayConfigLLO{
					ChainID:       "chain",
					FromBlock:     100,
					LLOConfigMode: "mode",
					LLODonID:      200,
				},
				PluginConfig: PluginConfigLLO{
					ChannelDefinitionsContractAddress:   "0xabc",
					ChannelDefinitionsContractFromBlock: 50,
					DonID:                               300,
					Servers:                             map[string]string{"server1": "http://localhost"},
				},
			},
			want: lloSpecTOML1,
		},
		{
			name: "empty minimal fields",
			spec: LLOJobSpec{
				Base: Base{
					Name:          "Empty-DON-Test",
					Type:          "offchainreporting2",
					SchemaVersion: 1,
					ExternalJobID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				},
				ContractID:  "contract-empty",
				RelayConfig: RelayConfigLLO{},
				PluginConfig: PluginConfigLLO{
					Servers: map[string]string{},
				},
			},
			want: lloSpecTOML2,
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
