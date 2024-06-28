package evm_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestRelayerOpts_Validate(t *testing.T) {
	type fields struct {
		DS                   sqlutil.DataSource
		CSAETHKeystore       evm.CSAETHKeystore
		CapabilitiesRegistry coretypes.CapabilitiesRegistry
		DBURL                *url.URL
	}
	tests := []struct {
		name            string
		fields          fields
		wantErrContains string
		beforeFn        func(t *testing.T)
	}{
		{
			name: "all invalid",
			fields: fields{
				DS:                   nil,
				CSAETHKeystore:       nil,
				CapabilitiesRegistry: nil,
				DBURL:                nil,
			},
			beforeFn: func(t *testing.T) { t.Setenv(string(env.DatabaseURL), "") },
			wantErrContains: `nil DataSource
nil Keystore
nil CapabilitiesRegistry
no DBURL provided and CL_DATABASE_URL unset`,
		},
		{
			name: "missing ds, keystore, capabilitiesRegistry",
			fields: fields{
				DS:    nil,
				DBURL: nil,
			},
			beforeFn: func(t *testing.T) { t.Setenv(string(env.DatabaseURL), ":/#unparseable") },
			wantErrContains: `nil DataSource
nil Keystore
nil CapabilitiesRegistry
failed to parse CL_DATABASE_URL`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFn(t)
			c := evm.RelayerOpts{
				DS:                   tt.fields.DS,
				CSAETHKeystore:       tt.fields.CSAETHKeystore,
				CapabilitiesRegistry: tt.fields.CapabilitiesRegistry,
			}
			err := c.Validate()
			require.Equal(t, tt.wantErrContains != "", err != nil)
			if tt.wantErrContains != "" {
				assert.Contains(t, err.Error(), tt.wantErrContains)
			}
		})
	}
}
