package evm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestRelayerOpts_Validate(t *testing.T) {
	type fields struct {
		DS                   sqlutil.DataSource
		CSAKeystore          core.Keystore
		EVMKeystore          keys.ChainStore
		CapabilitiesRegistry core.CapabilitiesRegistry
	}
	tests := []struct {
		name            string
		fields          fields
		wantErrContains string
	}{
		{
			name: "all invalid",
			fields: fields{
				DS:                   nil,
				EVMKeystore:          nil,
				CSAKeystore:          nil,
				CapabilitiesRegistry: nil,
			},
			wantErrContains: `nil DataSource
nil CSAKeystore
nil EVMKeystore`,
		},
		{
			name: "missing ds, keystore",
			fields: fields{
				DS: nil,
			},
			wantErrContains: `nil DataSource
nil CSAKeystore
nil EVMKeystore`,
		},
		{
			name: "missing ds, keystore, capabilitiesRegistry",
			fields: fields{
				DS: nil,
			},
			wantErrContains: `nil DataSource
nil CSAKeystore
nil EVMKeystore
nil CapabilitiesRegistry`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := evm.RelayerOpts{
				DS:                   tt.fields.DS,
				EVMKeystore:          tt.fields.EVMKeystore,
				CSAKeystore:          tt.fields.CSAKeystore,
				CapabilitiesRegistry: tt.fields.CapabilitiesRegistry,
			}
			err := c.Validate()
			if tt.wantErrContains != "" {
				assert.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
