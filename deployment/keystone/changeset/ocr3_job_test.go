package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveOCR3Job(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg OCR3JobConfig
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				cfg: OCR3JobConfig{
					JobName:              "OCR3 MultiChain Capability",
					ChainID:              "11155111",
					OCR2EVMKeyBundleID:   "0x123",
					TransmitterID:        "0x456",
					OCR2AptosKeyBundleID: "0x789",
					ContractID:           "0x5Cf6c5D188Ed72428aEC2E53afA5149F7B50f38D",
					P2Pv2Bootstrappers: []string{
						"12D3KooWB78smeNRVYBVsZVQUXTZ9jemuW7Y7yvwSMov56NZW2z1@cl-keystone-one-bt-0:5001",
						"12D3KooWBeFB8CGj1Pqph86ugvp7JCSuN88ffx1gtAS9857x9hoE@cl-keystone-one-bt-2:5001",
						"12D3KooWDjdJ2kMS5CDZe41orMhpnFEZ2Ht13WY37ecWZRYjAVvi@cl-keystone-one-bt-3:5001",
						"12D3KooWPBeMpPvDNmBdAjab7aLCEafZs5Q4EHcp4CAqZSLzScUL@cl-keystone-one-bt-1:5001",
					},
					ExternalJobID: "f1ac5211-ab79-4c31-ba1c-0997b72db466",
				},
			},
			want: `type = "offchainreporting2"
schemaVersion = 1
name = "OCR3 MultiChain Capability"
externalJobID = "f1ac5211-ab79-4c31-ba1c-0997b72db466"
contractID = "0x5Cf6c5D188Ed72428aEC2E53afA5149F7B50f38D"
ocrKeyBundleID = "0x123"
p2pv2Bootstrappers = [
  "12D3KooWB78smeNRVYBVsZVQUXTZ9jemuW7Y7yvwSMov56NZW2z1@cl-keystone-one-bt-0:5001",
  "12D3KooWBeFB8CGj1Pqph86ugvp7JCSuN88ffx1gtAS9857x9hoE@cl-keystone-one-bt-2:5001",
  "12D3KooWDjdJ2kMS5CDZe41orMhpnFEZ2Ht13WY37ecWZRYjAVvi@cl-keystone-one-bt-3:5001",
  "12D3KooWPBeMpPvDNmBdAjab7aLCEafZs5Q4EHcp4CAqZSLzScUL@cl-keystone-one-bt-1:5001",
]
relay = "evm"
pluginType = "plugin"
transmitterID = "0x456"

[relayConfig]
chainID = "11155111"

[pluginConfig]
command = "chainlink-ocr3-capability"
ocrVersion = 3
pluginName = "ocr-capability"
providerType = "ocr3-capability"
telemetryType = "plugin"

[onchainSigningStrategy]
strategyName = 'multi-chain'
[onchainSigningStrategy.config]
evm = "0x123"
aptos = "0x789"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ResolveOCR3Job(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveOCR3Job() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
