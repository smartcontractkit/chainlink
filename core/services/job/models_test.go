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

func TestJSONConfig_BytesWithPreservedJson(t *testing.T) {
	type testCases struct {
		name     string
		Input    JSONConfig
		Expected []byte
	}
	tests := []testCases{
		{
			name: "json",
			Input: JSONConfig{
				"key": "{\"nestedKey\": {\"nestedValue\":123}}",
			},
			Expected: []byte(`{"key":{"nestedKey":{"nestedValue":123}}}`),
		},
		{
			name: "broken json gets treated as a regular string",
			Input: JSONConfig{
				"key": "2324{\"nes4tedKey\":\"nestedValue\"}",
			},
			Expected: []byte(`{"key":"2324{\"nes4tedKey\":\"nestedValue\"}"}`),
		},
		{
			name: "number",
			Input: JSONConfig{
				"key": 1,
			},
			Expected: []byte(`{"key":1}`),
		},
		{
			name: "string",
			Input: JSONConfig{
				"key": "abc",
			},
			Expected: []byte(`{"key":"abc"}`),
		},
		{
			name: "all together",
			Input: JSONConfig{
				"key1": "{\"nestedKey\": {\"nestedValue\":123}}",
				"key2": "2324{\"key\":\"value\"}",
				"key3": 1,
				"key4": "abc",
			},
			Expected: []byte(`{"key1":{"nestedKey":{"nestedValue":123}},"key2":"2324{\"key\":\"value\"}","key3":1,"key4":"abc"}`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.Input.BytesWithPreservedJson()
			if !reflect.DeepEqual(result, tc.Expected) {
				t.Errorf("Input: %v, BytesWithPreservedJson() returned unexpected result. Expected: %s, Got: %s", tc.Input, tc.Expected, result)
			}
		})

	}
}
