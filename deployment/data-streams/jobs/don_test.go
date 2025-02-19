package jobs

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDonJobSpec_MarshalTOML(t *testing.T) {
	trueVal := true
	ocrKeyBundle := "ocr-bundle-123"

	testCases := []struct {
		name       string
		spec       DonJobSpec
		wantSubstr []string
	}{
		{
			name: "with fields populated",
			spec: DonJobSpec{
				Base: Base{
					Name:          "Test-DON",
					Type:          "don",
					SchemaVersion: 1,
					ExternalJobID: uuid.New(),
				},
				ContractID:                        "contract-123",
				TransmitterID:                     "tx-123",
				ForwardingAllowed:                 &trueVal,
				P2PV2Bootstrappers:                []string{"bootstrap1", "bootstrap2"},
				OCRKeyBundleID:                    &ocrKeyBundle,
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
			wantSubstr: []string{
				"contractID",
				"contract-123",
				"transmitterID",
				"tx-123",
				"p2pv2Bootstrappers",
				"bootstrap1",
				"ocrKeyBundleID",
				"ocr-bundle-123",
				"maxTaskDuration",
				"contractConfigTrackerPollInterval",
				"relay",
				"testrelay",
				"pluginType",
				"testplugin",
				"chainID",
				"chain",
				"fromBlock",
				"100",
				"lloConfigMode",
				"mode",
				"lloDonID",
				"200",
				"channelDefinitionsContractAddress",
				"0xabc",
				"channelDefinitionsContractFromBlock",
				"50",
				"donID",
				"300",
				"servers",
				"server1",
				"http://localhost",
			},
		},
		{
			name: "empty minimal fields",
			spec: DonJobSpec{
				Base: Base{
					Name:          "Empty-DON-Test",
					Type:          "don",
					SchemaVersion: 1,
					ExternalJobID: uuid.New(),
				},
				ContractID:  "contract-empty",
				RelayConfig: RelayConfigDon{},
				PluginConfig: PluginConfigDon{
					Servers: map[string]string{},
				},
			},
			wantSubstr: []string{
				"contractID",
				"contract-empty",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tomlBytes, err := tc.spec.MarshalTOML()
			require.NoError(t, err)

			result := string(tomlBytes)
			for _, substr := range tc.wantSubstr {
				require.Contains(t, result, substr, "result %q does not contain expected substring %q", result, substr)
			}
		})
	}
}
