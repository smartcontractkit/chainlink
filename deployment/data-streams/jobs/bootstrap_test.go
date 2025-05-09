package jobs

import (
	"testing"

	"github.com/google/uuid"
)

const bootstrapSpecTOML = `name = 'bootstrap 1'
type = 'bootstrap'
schemaVersion = 1
externalJobID = 'f1ac5211-ab79-4c31-ba1c-0997b72db466'
contractID = '0x123'
donID = 1
relay = 'evm'

[relayConfig]
chainID = '42161'
fromBlock = 283806260
`

func TestNewBootstrapSpec(t *testing.T) {
	t.Parallel()

	externalJobID := uuid.New()

	tests := []struct {
		name          string
		contractID    string
		donID         uint64
		donName       string
		relay         RelayType
		relayConfig   RelayConfig
		externalJobID uuid.UUID
		want          *BootstrapSpec
	}{
		{
			name:       "success",
			contractID: "0x01",
			donID:      123,
			donName:    "don-123",
			relay:      RelayTypeEVM,
			relayConfig: RelayConfig{
				ChainID:   "234",
				FromBlock: 345,
			},
			externalJobID: externalJobID,
			want: &BootstrapSpec{
				Base: Base{
					Name:          "don-123 | 123",
					Type:          JobSpecTypeBootstrap,
					SchemaVersion: 1,
					ExternalJobID: externalJobID,
				},
				ContractID: "0x01",
				DonID:      123,
				Relay:      RelayTypeEVM,
				RelayConfig: RelayConfig{
					ChainID:   "234",
					FromBlock: 345,
				}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NewBootstrapSpec(tc.contractID, tc.donID, tc.donName, tc.relay, tc.relayConfig, tc.externalJobID)
			if got == nil {
				t.Fatal("got nil")
			}
			if *got != *tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMarshalTOML(t *testing.T) {
	t.Parallel()

	bootstrapSpec := BootstrapSpec{
		Base: Base{
			Name:          "bootstrap 1",
			Type:          JobSpecTypeBootstrap,
			SchemaVersion: 1,
			ExternalJobID: uuid.MustParse("f1ac5211-ab79-4c31-ba1c-0997b72db466"),
		},
		ContractID: "0x123",
		DonID:      1,
		Relay:      RelayTypeEVM,
		RelayConfig: RelayConfig{
			ChainID:   "42161",
			FromBlock: 283806260,
		},
	}

	tests := []struct {
		name string
		give BootstrapSpec
		want string
	}{
		{
			name: "bootstrap 1",
			give: bootstrapSpec,
			want: bootstrapSpecTOML,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.give.MarshalTOML()
			if err != nil {
				t.Fatal(err)
			}

			if string(got) != test.want {
				t.Errorf("got %s, want %s", got, test.want)
			}
		})
	}
}
