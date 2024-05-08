package capreg

import (
	"errors"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/stretchr/testify/require"
)

func Test_capabilityLauncher_RegisterCapabilityFactory(t *testing.T) {
	type args struct {
		cf func(t *testing.T) CapabilityFactory
	}
	tests := []struct {
		name                 string
		args                 args
		startingCapabilities []CapabilityFactory
		wantErr              bool
	}{
		{
			"success",
			args{
				cf: func(t *testing.T) CapabilityFactory {
					cf := NewMockCapabilityFactory(t)
					cf.On("CapabilityID").Return("a1")
					return cf
				},
			},
			[]CapabilityFactory{},
			false,
		},
		{
			"bad capability id",
			args{
				cf: func(t *testing.T) CapabilityFactory {
					cf := NewMockCapabilityFactory(t)
					cf.On("CapabilityID").Return("")
					return cf
				},
			},
			[]CapabilityFactory{},
			true,
		},
		{
			"capability already registered",
			args{
				cf: func(t *testing.T) CapabilityFactory {
					cf := NewMockCapabilityFactory(t)
					cf.On("CapabilityID").Return("a1")
					return cf
				},
			},
			[]CapabilityFactory{
				func() CapabilityFactory {
					cf := NewMockCapabilityFactory(t)
					cf.On("CapabilityID").Return("a1")
					return cf
				}(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myP2PID := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1337)).PeerID()
			l := NewCapabilityLauncher(myP2PID[:])
			for _, v := range tt.startingCapabilities {
				require.NoError(t, l.RegisterCapabilityFactory(v))
			}
			err := l.RegisterCapabilityFactory(tt.args.cf(t))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_capabilityLauncher_Close(t *testing.T) {
	tests := []struct {
		name         string
		capabilities []CapabilityFactory
		wantErr      bool
	}{
		{
			"no errors closing",
			[]CapabilityFactory{
				func() CapabilityFactory {
					cf := NewMockCapabilityFactory(t)
					cf.On("CapabilityID").Return("a1")
					cf.On("Close").Return(nil)
					return cf
				}(),
			},
			false,
		},
		{
			"errors bubbled up if one capability fails to close",
			[]CapabilityFactory{
				func() CapabilityFactory {
					cf := NewMockCapabilityFactory(t)
					cf.On("CapabilityID").Return("a1")
					cf.On("Close").Return(nil)
					return cf
				}(),
				func() CapabilityFactory {
					cf := NewMockCapabilityFactory(t)
					cf.On("CapabilityID").Return("a2")
					cf.On("Close").Return(errors.New("close error"))
					return cf
				}(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myP2PID := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1337)).PeerID()
			l := NewCapabilityLauncher(myP2PID[:])
			for _, v := range tt.capabilities {
				require.NoError(t, l.RegisterCapabilityFactory(v))
			}
			if tt.wantErr {
				require.Error(t, l.Close())
			} else {
				require.NoError(t, l.Close())
			}
		})
	}
}

func Test_capabilityLauncher_handleDeletedDONs(t *testing.T) {
	type args struct {
		relevantDONs map[uint32]DON
		errs         error
	}
	tests := []struct {
		name    string
		args    args
		myDONs  map[uint32]DON
		wantErr bool
	}{
		{
			"both empty",
			args{
				relevantDONs: map[uint32]DON{},
			},
			map[uint32]DON{},
			false,
		},
		{
			"no deleted dons",
			args{
				relevantDONs: map[uint32]DON{
					1: {ID: 1, CapabilityConfigurations: []CapabilityConfiguration{{CapabilityID: "a1"}}},
					2: {ID: 2, CapabilityConfigurations: []CapabilityConfiguration{{CapabilityID: "a2"}}},
				},
			},
			map[uint32]DON{
				1: {ID: 1, CapabilityConfigurations: []CapabilityConfiguration{{CapabilityID: "a1"}}},
				2: {ID: 2, CapabilityConfigurations: []CapabilityConfiguration{{CapabilityID: "a2"}}},
			},
			false,
		},
		{
			"some deleted dons",
			args{
				relevantDONs: map[uint32]DON{
					1: {ID: 1, CapabilityConfigurations: []CapabilityConfiguration{{CapabilityID: "a1"}}},
				},
			},
			map[uint32]DON{
				1: {ID: 1, CapabilityConfigurations: []CapabilityConfiguration{{CapabilityID: "a1"}}},
				2: {ID: 2, CapabilityConfigurations: []CapabilityConfiguration{{CapabilityID: "a2"}}},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myP2PID := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1337)).PeerID()
			l := NewCapabilityLauncher(myP2PID[:])
			err := l.handleDeletedDONs(testutils.Context(t), tt.args.relevantDONs, tt.args.errs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
