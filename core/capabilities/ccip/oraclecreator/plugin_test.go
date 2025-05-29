package oraclecreator

import (
	"math/big"
	"strings"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

func TestPluginTypeToTelemetryType(t *testing.T) {
	tests := []struct {
		name       string
		pluginType types.PluginType
		expected   synchronization.TelemetryType
		expectErr  bool
	}{
		{
			name:       "CCIPCommit plugin type",
			pluginType: types.PluginTypeCCIPCommit,
			expected:   synchronization.OCR3CCIPCommit,
			expectErr:  false,
		},
		{
			name:       "CCIPExec plugin type",
			pluginType: types.PluginTypeCCIPExec,
			expected:   synchronization.OCR3CCIPExec,
			expectErr:  false,
		},
		{
			name:       "Unknown plugin type",
			pluginType: types.PluginType(99),
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pluginTypeToTelemetryType(tt.pluginType)
			if tt.expectErr {
				require.Error(t, err, "Expected an error for plugin type %d", tt.pluginType)
				require.Equal(t, synchronization.TelemetryType(""), result, "Expected empty result for plugin type %d", tt.pluginType)
			} else {
				require.NoError(t, err, "Unexpected error for plugin type %d", tt.pluginType)
				require.Equal(t, tt.expected, result, "Unexpected telemetry type for plugin type %d", tt.pluginType)
			}
		})
	}
}

func TestPluginOracleCreator_getTransmitterFromPublicConfig(t *testing.T) {
	key := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))
	poc := &pluginOracleCreator{p2pID: key}

	peerID := key.PeerID().String()
	transmitAccount := ocrtypes.Account("test-transmit-account")

	tests := []struct {
		name          string
		publicConfig  ocr3confighelper.PublicConfig
		expected      ocrtypes.Account
		expectErr     bool
		expectedError string
	}{
		{
			name: "matching peerID",
			publicConfig: ocr3confighelper.PublicConfig{
				OracleIdentities: []confighelper.OracleIdentity{
					{PeerID: strings.TrimPrefix(peerID, "p2p_"), TransmitAccount: transmitAccount},
					{PeerID: "another-peer", TransmitAccount: ocrtypes.Account("another-account")},
				},
			},
			expected:  transmitAccount,
			expectErr: false,
		},
		{
			name: "non-matching peerID",
			publicConfig: ocr3confighelper.PublicConfig{
				OracleIdentities: []confighelper.OracleIdentity{
					{PeerID: "another-peer", TransmitAccount: ocrtypes.Account("another-account")},
				},
			},
			expected:      ocrtypes.Account(""),
			expectErr:     true,
			expectedError: "no transmitter found for my peer id " + peerID + " in public config",
		},
		{
			name: "empty OracleIdentities",
			publicConfig: ocr3confighelper.PublicConfig{
				OracleIdentities: []confighelper.OracleIdentity{},
			},
			expected:      ocrtypes.Account(""),
			expectErr:     true,
			expectedError: "no transmitter found for my peer id " + peerID + " in public config",
		},
		{
			name:          "nil OracleIdentities",
			publicConfig:  ocr3confighelper.PublicConfig{OracleIdentities: nil},
			expected:      ocrtypes.Account(""),
			expectErr:     true,
			expectedError: "no transmitter found for my peer id " + peerID + " in public config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := poc.getTransmitterFromPublicConfig(tt.publicConfig)
			if tt.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
				require.Equal(t, tt.expected, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}
