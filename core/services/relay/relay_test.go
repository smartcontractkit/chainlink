package relay

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestIdentifier_UnmarshalString(t *testing.T) {
	type fields struct {
		Network Network
		ChainID ChainID
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		want    fields
		args    args
		wantErr bool
	}{
		{name: "evm",
			args:    args{s: "evm.1"},
			wantErr: false,
			want:    fields{Network: EVM, ChainID: "1"},
		},
		{name: "bad network",
			args:    args{s: "notANetwork.1"},
			wantErr: true,
		},
		{name: "bad pattern",
			args:    args{s: "evm_1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ID{}
			err := i.UnmarshalString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Identifier.UnmarshalString() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want.Network, i.Network)
			assert.Equal(t, tt.want.ChainID, i.ChainID)
		})
	}
}

func TestNewID(t *testing.T) {
	rid := NewID(EVM, "chain id")
	assert.Equal(t, EVM, rid.Network)
	assert.Equal(t, "chain id", rid.ChainID)
}

type staticMedianProvider struct {
	types.MedianProvider
}

func (s staticMedianProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return nil
}

func (s staticMedianProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return nil
}

func (s staticMedianProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return nil
}

func (s staticMedianProvider) ReportCodec() median.ReportCodec {
	return nil
}

func (s staticMedianProvider) MedianContract() median.MedianContract {
	return nil
}

func (s staticMedianProvider) OnchainConfigCodec() median.OnchainConfigCodec {
	return nil
}

type staticFunctionsProvider struct {
	types.FunctionsProvider
}

type staticMercuryProvider struct {
	types.MercuryProvider
}

type staticAutomationProvider struct {
	types.AutomationProvider
}

type mockRelayer struct {
	types.Relayer
}

func (m *mockRelayer) NewMedianProvider(rargs types.RelayArgs, pargs types.PluginArgs) (types.MedianProvider, error) {
	return staticMedianProvider{}, nil
}

func (m *mockRelayer) NewFunctionsProvider(rargs types.RelayArgs, pargs types.PluginArgs) (types.FunctionsProvider, error) {
	return staticFunctionsProvider{}, nil
}

func (m *mockRelayer) NewMercuryProvider(rargs types.RelayArgs, pargs types.PluginArgs) (types.MercuryProvider, error) {
	return staticMercuryProvider{}, nil
}

func (m *mockRelayer) NewAutomationProvider(rargs types.RelayArgs, pargs types.PluginArgs) (types.AutomationProvider, error) {
	return staticAutomationProvider{}, nil
}

type mockRelayerExt struct {
	loop.RelayerExt
}

func isType[T any](p any) bool {
	_, ok := p.(T)
	return ok
}

func TestRelayerServerAdapter(t *testing.T) {
	r := &mockRelayer{}
	sa := NewServerAdapter(r, mockRelayerExt{})

	testCases := []struct {
		ProviderType string
		Test         func(p any) bool
		Error        string
	}{
		{
			ProviderType: string(types.Median),
			Test:         isType[types.MedianProvider],
		},
		{
			ProviderType: string(types.Functions),
			Test:         isType[types.FunctionsProvider],
		},
		{
			ProviderType: string(types.Mercury),
			Test:         isType[types.MercuryProvider],
		},
		{
			ProviderType: string(types.CCIPCommit),
			Error:        "provider type not supported",
		},
		{
			ProviderType: string(types.CCIPExecution),
			Error:        "provider type not supported",
		},
		{
			ProviderType: "unknown",
			Error:        "provider type not recognized",
		},
		{
			ProviderType: string(types.GenericPlugin),
			Error:        "unexpected call to NewPluginProvider",
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		pp, err := sa.NewPluginProvider(
			ctx,
			types.RelayArgs{ProviderType: tc.ProviderType},
			types.PluginArgs{},
		)

		if tc.Error != "" {
			assert.ErrorContains(t, err, tc.Error)
		} else {
			assert.NoError(t, err)
			assert.True(t, tc.Test(pp))
		}
	}
}
