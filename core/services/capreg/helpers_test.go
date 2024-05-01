package capreg

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/stretchr/testify/require"
)

func Test_filterRelevantDONs(t *testing.T) {
	type args struct {
		localP2PID []byte
		s          State
	}

	peerID := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1337)).PeerID()

	tests := []struct {
		name             string
		args             args
		wantRelevantDONs map[uint32]DON
	}{
		{
			"no relevant DONS",
			args{
				localP2PID: peerID[:],
				s: State{
					DONs: []DON{
						{
							ID:    1,
							Nodes: [][]byte{[]byte("0xdeadbeef")},
						},
					},
				},
			},
			map[uint32]DON{},
		},
		{
			"relevant DONS",
			args{
				localP2PID: peerID[:],
				s: State{
					DONs: []DON{
						{
							ID:    1,
							Nodes: [][]byte{[]byte("0xdeadbeef"), []byte("superdeadbeef"), peerID[:]},
						},
						{
							ID:    2,
							Nodes: [][]byte{[]byte("0xdeadbeef"), []byte("superdeadbeef"), []byte("notme")},
						},
						{
							ID:    3,
							Nodes: [][]byte{[]byte("0xdeadbeef"), []byte("superdeadbeef"), []byte("notme"), peerID[:]},
						},
					},
				},
			},
			map[uint32]DON{
				1: {
					ID:    1,
					Nodes: [][]byte{[]byte("0xdeadbeef"), []byte("superdeadbeef"), peerID[:]},
				},
				3: {
					ID:    3,
					Nodes: [][]byte{[]byte("0xdeadbeef"), []byte("superdeadbeef"), []byte("notme"), peerID[:]},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRelevantDONs := filterRelevantDONs(tt.args.localP2PID, tt.args.s)
			require.Equal(t, tt.wantRelevantDONs, actualRelevantDONs)
		})
	}
}

func Test_nodesChanged(t *testing.T) {
	type args struct {
		oldNodes [][]byte
		newNodes [][]byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"no change",
			args{
				oldNodes: [][]byte{[]byte("0xdeadbeef"), []byte("superdeadbeef"), []byte("notme")},
				newNodes: [][]byte{[]byte("0xdeadbeef"), []byte("superdeadbeef"), []byte("notme")},
			},
			false,
		},
		{
			"change",
			args{
				oldNodes: [][]byte{[]byte("0xdeadbeef"), []byte("superdeadbeef"), []byte("notme")},
				newNodes: [][]byte{[]byte("0xdeadbeef"), []byte("superdeadbeef"), []byte("notme"), []byte("newguy")},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nodesChanged(tt.args.oldNodes, tt.args.newNodes)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_capabilityDiff(t *testing.T) {
	type args struct {
		oldCCs []CapabilityConfiguration
		newCCs []CapabilityConfiguration
	}
	tests := []struct {
		name                         string
		args                         args
		wantRemovedCapabilities      []CapabilityID
		wantNewOrUpdatedCapabilities []CapabilityID
	}{
		{
			"no change",
			args{
				oldCCs: []CapabilityConfiguration{
					{CapabilityID: "c1", OnchainConfigVersion: 1, OnchainConfig: []byte("onchain1"), OffchainConfigVersion: 1, OffchainConfig: []byte("offchain1")},
					{CapabilityID: "d2", OnchainConfigVersion: 2, OnchainConfig: []byte("onchain2"), OffchainConfigVersion: 2, OffchainConfig: []byte("offchain2")},
					{CapabilityID: "e3", OnchainConfigVersion: 3, OnchainConfig: []byte("onchain3"), OffchainConfigVersion: 3, OffchainConfig: []byte("offchain3")},
				},
				newCCs: []CapabilityConfiguration{
					{CapabilityID: "c1", OnchainConfigVersion: 1, OnchainConfig: []byte("onchain1"), OffchainConfigVersion: 1, OffchainConfig: []byte("offchain1")},
					{CapabilityID: "d2", OnchainConfigVersion: 2, OnchainConfig: []byte("onchain2"), OffchainConfigVersion: 2, OffchainConfig: []byte("offchain2")},
					{CapabilityID: "e3", OnchainConfigVersion: 3, OnchainConfig: []byte("onchain3"), OffchainConfigVersion: 3, OffchainConfig: []byte("offchain3")},
				},
			},
			nil,
			nil,
		},
		{
			"removed capabilities",
			args{
				oldCCs: []CapabilityConfiguration{
					{CapabilityID: "c1", OnchainConfigVersion: 1, OnchainConfig: []byte("onchain1"), OffchainConfigVersion: 1, OffchainConfig: []byte("offchain1")},
					{CapabilityID: "d2", OnchainConfigVersion: 2, OnchainConfig: []byte("onchain2"), OffchainConfigVersion: 2, OffchainConfig: []byte("offchain2")},
					{CapabilityID: "e3", OnchainConfigVersion: 3, OnchainConfig: []byte("onchain3"), OffchainConfigVersion: 3, OffchainConfig: []byte("offchain3")},
				},
				newCCs: []CapabilityConfiguration{
					{CapabilityID: "c1", OnchainConfigVersion: 1, OnchainConfig: []byte("onchain1"), OffchainConfigVersion: 1, OffchainConfig: []byte("offchain1")},
				},
			},
			[]CapabilityID{"d2", "e3"},
			nil,
		},
		{
			"new or updated capabilities",
			args{
				oldCCs: []CapabilityConfiguration{
					{CapabilityID: "c1", OnchainConfigVersion: 1, OnchainConfig: []byte("onchain1"), OffchainConfigVersion: 1, OffchainConfig: []byte("offchain1")},
					{CapabilityID: "d2", OnchainConfigVersion: 2, OnchainConfig: []byte("onchain2"), OffchainConfigVersion: 2, OffchainConfig: []byte("offchain2")},
					{CapabilityID: "e3", OnchainConfigVersion: 3, OnchainConfig: []byte("onchain3"), OffchainConfigVersion: 3, OffchainConfig: []byte("offchain3")},
				},
				newCCs: []CapabilityConfiguration{
					{CapabilityID: "c1", OnchainConfigVersion: 1, OnchainConfig: []byte("onchain1"), OffchainConfigVersion: 1, OffchainConfig: []byte("offchain1")},
					{CapabilityID: "d2", OnchainConfigVersion: 2, OnchainConfig: []byte("onchain2"), OffchainConfigVersion: 2, OffchainConfig: []byte("offchain2")},
					{CapabilityID: "e3", OnchainConfigVersion: 3, OnchainConfig: []byte("onchain3"), OffchainConfigVersion: 3, OffchainConfig: []byte("offchain3")},
					{CapabilityID: "f4", OnchainConfigVersion: 4, OnchainConfig: []byte("onchain4"), OffchainConfigVersion: 4, OffchainConfig: []byte("offchain4")},
					{CapabilityID: "g5", OnchainConfigVersion: 5, OnchainConfig: []byte("onchain5"), OffchainConfigVersion: 5, OffchainConfig: []byte("offchain5")},
				},
			},
			nil,
			[]CapabilityID{"f4", "g5"},
		},
		{
			"removed and new or updated capabilities",
			args{
				oldCCs: []CapabilityConfiguration{
					{CapabilityID: "c1", OnchainConfigVersion: 1, OnchainConfig: []byte("onchain1"), OffchainConfigVersion: 1, OffchainConfig: []byte("offchain1")},
					{CapabilityID: "d2", OnchainConfigVersion: 2, OnchainConfig: []byte("onchain2"), OffchainConfigVersion: 2, OffchainConfig: []byte("offchain2")},
					{CapabilityID: "e3", OnchainConfigVersion: 3, OnchainConfig: []byte("onchain3"), OffchainConfigVersion: 3, OffchainConfig: []byte("offchain3")},
				},
				newCCs: []CapabilityConfiguration{
					{CapabilityID: "c1", OnchainConfigVersion: 1, OnchainConfig: []byte("onchain1"), OffchainConfigVersion: 1, OffchainConfig: []byte("offchain1")},
					{CapabilityID: "f4", OnchainConfigVersion: 4, OnchainConfig: []byte("onchain4"), OffchainConfigVersion: 4, OffchainConfig: []byte("offchain4")},
					{CapabilityID: "g5", OnchainConfigVersion: 5, OnchainConfig: []byte("onchain5"), OffchainConfigVersion: 5, OffchainConfig: []byte("offchain5")},
				},
			},
			[]CapabilityID{"d2", "e3"},
			[]CapabilityID{"f4", "g5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRemovedCapabilities, gotNewOrUpdatedCapabilities := capabilityDiff(tt.args.oldCCs, tt.args.newCCs)
			require.Equal(t, tt.wantRemovedCapabilities, gotRemovedCapabilities)
			require.Equal(t, tt.wantNewOrUpdatedCapabilities, gotNewOrUpdatedCapabilities)
		})
	}
}
