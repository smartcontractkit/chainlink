package evm_test

import (
	"testing"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"

	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestRelayerOpts_Validate(t *testing.T) {
	cfg := configtest.NewTestGeneralConfig(t)
	type fields struct {
		DB               *sqlx.DB
		QConfig          pg.QConfig
		CSAETHKeystore   evm.CSAETHKeystore
		EventBroadcaster pg.EventBroadcaster
	}
	tests := []struct {
		name            string
		fields          fields
		wantErrContains string
	}{
		{
			name: "all invalid",
			fields: fields{
				DB:               nil,
				QConfig:          nil,
				CSAETHKeystore:   nil,
				EventBroadcaster: nil,
			},
			wantErrContains: `nil DB
nil QConfig
nil Keystore
nil Eventbroadcaster`,
		},
		{
			name: "missing db, keystore",
			fields: fields{
				DB:               nil,
				QConfig:          cfg.Database(),
				CSAETHKeystore:   nil,
				EventBroadcaster: pg.NewNullEventBroadcaster(),
			},
			wantErrContains: `nil DB
nil Keystore`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := evm.RelayerOpts{
				DB:               tt.fields.DB,
				QConfig:          tt.fields.QConfig,
				CSAETHKeystore:   tt.fields.CSAETHKeystore,
				EventBroadcaster: tt.fields.EventBroadcaster,
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
