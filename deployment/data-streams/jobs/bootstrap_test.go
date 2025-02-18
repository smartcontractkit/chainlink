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

func Test_Bootstrap(t *testing.T) {
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
