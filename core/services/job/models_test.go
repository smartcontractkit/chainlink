package job

import (
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

func TestOCR2OracleSpec_RelayIdentifier(t *testing.T) {
	type fields struct {
		Relay       relay.Network
		ChainID     string
		RelayConfig JSONConfig
	}
	tests := []struct {
		name    string
		fields  fields
		want    relay.ID
		wantErr bool
	}{
		{name: "err no chain id",
			fields:  fields{},
			want:    relay.ID{},
			wantErr: true,
		},
		{
			name: "evm explicitly configured",
			fields: fields{
				Relay:   relay.EVM,
				ChainID: "1",
			},
			want: relay.ID{Network: relay.EVM, ChainID: relay.ChainID("1")},
		},
		{
			name: "evm implicitly configured",
			fields: fields{
				Relay:       relay.EVM,
				RelayConfig: map[string]any{"chainID": 1},
			},
			want: relay.ID{Network: relay.EVM, ChainID: relay.ChainID("1")},
		},
		{
			name: "evm implicitly configured with bad value",
			fields: fields{
				Relay:       relay.EVM,
				RelayConfig: map[string]any{"chainID": float32(1)},
			},
			want:    relay.ID{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &OCR2OracleSpec{
				Relay:       tt.fields.Relay,
				ChainID:     tt.fields.ChainID,
				RelayConfig: tt.fields.RelayConfig,
			}
			got, err := s.RelayID()
			if (err != nil) != tt.wantErr {
				t.Errorf("OCR2OracleSpec.RelayIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OCR2OracleSpec.RelayIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
