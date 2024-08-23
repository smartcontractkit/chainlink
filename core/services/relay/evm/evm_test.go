package evm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestRelayerOpts_Validate(t *testing.T) {
	type fields struct {
		DS                   sqlutil.DataSource
		CSAETHKeystore       evm.CSAETHKeystore
		CapabilitiesRegistry coretypes.CapabilitiesRegistry
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
				CSAETHKeystore:       nil,
				CapabilitiesRegistry: nil,
			},
			wantErrContains: `nil DataSource
nil Keystore`,
		},
		{
			name: "missing ds, keystore",
			fields: fields{
				DS: nil,
			},
			wantErrContains: `nil DataSource
nil Keystore`,
		},
		{
			name: "missing ds, keystore, capabilitiesRegistry",
			fields: fields{
				DS: nil,
			},
			wantErrContains: `nil DataSource
nil Keystore
nil CapabilitiesRegistry`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := evm.RelayerOpts{
				DS:                   tt.fields.DS,
				CSAETHKeystore:       tt.fields.CSAETHKeystore,
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
