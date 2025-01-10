package feeds

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func Test_NewChainType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    string
		want    ChainType
		wantErr error
	}{
		{
			name: "EVM Chain Type",
			give: "EVM",
			want: ChainTypeEVM,
		},
		{
			name: "Starknet Chain Type",
			give: "STARKNET",
			want: ChainTypeStarknet,
		},
		{
			name: "Tron Chain Type",
			give: "TRON",
			want: ChainTypeTron,
		},
		{
			name:    "Invalid Chain Type",
			give:    "",
			want:    ChainTypeUnknown,
			wantErr: errors.New("invalid chain type"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ct, err := NewChainType(tt.give)

			assert.Equal(t, tt.want, ct)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func Test_ToPluginType(t *testing.T) {
	t.Parallel()

	pt, err := ToPluginType("commit")
	require.NoError(t, err)
	assert.Equal(t, pt, PluginTypeCommit)

	pt, err = ToPluginType("execute")
	require.NoError(t, err)
	assert.Equal(t, pt, PluginTypeExecute)

	pt, err = ToPluginType("median")
	require.NoError(t, err)
	assert.Equal(t, pt, PluginTypeMedian)

	pt, err = ToPluginType("mercury")
	require.NoError(t, err)
	assert.Equal(t, pt, PluginTypeMercury)

	pt, err = ToPluginType("rebalancer")
	require.NoError(t, err)
	assert.Equal(t, pt, PluginTypeRebalancer)

	pt, err = ToPluginType("xxx")
	require.Error(t, err)
	assert.Equal(t, pt, PluginTypeUnknown)
	assert.EqualError(t, err, "unknown plugin type")
}

func Test_FromPluginType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "commit", FromPluginTypeInput(PluginTypeCommit))
	assert.Equal(t, "execute", FromPluginTypeInput(PluginTypeExecute))
	assert.Equal(t, "median", FromPluginTypeInput(PluginTypeMedian))
	assert.Equal(t, "mercury", FromPluginTypeInput(PluginTypeMercury))
	assert.Equal(t, "rebalancer", FromPluginTypeInput(PluginTypeRebalancer))
	assert.Equal(t, "unknown", FromPluginTypeInput(PluginTypeUnknown))
}

func Test_FluxMonitorConfig_Value(t *testing.T) {
	t.Parallel()

	cfg := FluxMonitorConfig{Enabled: true}
	want := `{"enabled":true}`

	val, err := cfg.Value()
	require.NoError(t, err)

	actual, ok := val.([]byte)
	require.True(t, ok)

	assert.Equal(t, want, string(actual))
}

func Test_FluxMonitorConfig_Scan(t *testing.T) {
	t.Parallel()

	var (
		give = `{"enabled":true}`
		want = FluxMonitorConfig{Enabled: true}
	)

	var actual FluxMonitorConfig
	err := actual.Scan([]byte(give))
	require.NoError(t, err)

	assert.Equal(t, want, actual)
}

func Test_OCR1Config_Value(t *testing.T) {
	t.Parallel()

	var (
		multiaddr   = "multiaddr"
		p2pPeerID   = "peerid"
		keyBundleID = "ocrkeyid"
	)

	tests := []struct {
		name string
		give OCR1Config
		want string
	}{
		{
			name: "all fields populated",
			give: OCR1Config{
				Enabled:     true,
				IsBootstrap: false,
				Multiaddr:   null.StringFrom(multiaddr),
				P2PPeerID:   null.StringFrom(p2pPeerID),
				KeyBundleID: null.StringFrom(keyBundleID),
			},
			want: `{"enabled":true,"is_bootstrap":false,"multiaddr":"multiaddr","p2p_peer_id":"peerid","key_bundle_id":"ocrkeyid"}`,
		},
		{
			name: "bootstrap fields populated",
			give: OCR1Config{
				Enabled:     true,
				IsBootstrap: true,
				Multiaddr:   null.StringFrom(multiaddr),
				P2PPeerID:   null.StringFromPtr(nil),
				KeyBundleID: null.StringFromPtr(nil),
			},
			want: `{"enabled":true,"is_bootstrap":true,"multiaddr":"multiaddr","p2p_peer_id":null,"key_bundle_id":null}`,
		},
		{
			name: "multiaddr field populated",
			give: OCR1Config{
				Enabled:     true,
				IsBootstrap: false,
				Multiaddr:   null.StringFromPtr(nil),
				P2PPeerID:   null.StringFrom(p2pPeerID),
				KeyBundleID: null.StringFrom(keyBundleID),
			},
			want: `{"enabled":true,"is_bootstrap":false,"multiaddr":null,"p2p_peer_id":"peerid","key_bundle_id":"ocrkeyid"}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.give.Value()
			require.NoError(t, err)

			actual, ok := val.([]byte)
			require.True(t, ok)

			assert.Equal(t, tt.want, string(actual))
		})
	}
}

func Test_OCR1Config_Scan(t *testing.T) {
	t.Parallel()

	var (
		multiaddr   = "multiaddr"
		p2pPeerID   = "peerid"
		keyBundleID = "ocrkeyid"
	)

	tests := []struct {
		name string
		give string
		want OCR1Config
	}{
		{
			name: "all fields populated",
			give: `{"enabled":true,"is_bootstrap":false,"multiaddr":"multiaddr","p2p_peer_id":"peerid","key_bundle_id":"ocrkeyid"}`,
			want: OCR1Config{
				Enabled:     true,
				IsBootstrap: false,
				Multiaddr:   null.StringFrom(multiaddr),
				P2PPeerID:   null.StringFrom(p2pPeerID),
				KeyBundleID: null.StringFrom(keyBundleID),
			},
		},
		{
			name: "bootstrap fields populated",
			give: `{"enabled":true,"is_bootstrap":true,"multiaddr":"multiaddr","p2p_peer_id":null,"key_bundle_id":null}`,
			want: OCR1Config{
				Enabled:     true,
				IsBootstrap: true,
				Multiaddr:   null.StringFrom(multiaddr),
				P2PPeerID:   null.StringFromPtr(nil),
				KeyBundleID: null.StringFromPtr(nil),
			},
		},
		{
			name: "multiaddr field populated",
			give: `{"enabled":true,"is_bootstrap":false,"multiaddr":null,"p2p_peer_id":"peerid","key_bundle_id":"ocrkeyid"}`,
			want: OCR1Config{
				Enabled:     true,
				IsBootstrap: false,
				Multiaddr:   null.StringFromPtr(nil),
				P2PPeerID:   null.StringFrom(p2pPeerID),
				KeyBundleID: null.StringFrom(keyBundleID),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			var actual OCR1Config
			err := actual.Scan([]byte(tt.give))
			require.NoError(t, err)

			assert.Equal(t, tt.want, actual)
		})
	}
}

func Test_Plugins_Value(t *testing.T) {
	t.Parallel()

	var (
		give = Plugins{
			Commit:     true,
			Execute:    true,
			Median:     false,
			Mercury:    true,
			Rebalancer: false,
		}
		want = `{"commit":true,"execute":true,"median":false,"mercury":true,"rebalancer":false}`
	)

	val, err := give.Value()
	require.NoError(t, err)

	actual, ok := val.([]byte)
	require.True(t, ok)

	assert.Equal(t, want, string(actual))
}

func Test_Plugins_Scan(t *testing.T) {
	t.Parallel()

	var (
		give = `{"commit":true,"execute":true,"median":false,"mercury":true,"rebalancer":false}`
		want = Plugins{
			Commit:     true,
			Execute:    true,
			Median:     false,
			Mercury:    true,
			Rebalancer: false,
		}
	)

	var actual Plugins
	err := actual.Scan(give)
	require.NoError(t, err)

	assert.Equal(t, want, actual)
}

func Test_OCR2Config_Value(t *testing.T) {
	t.Parallel()

	var (
		give = OCR2ConfigModel{
			Enabled:          true,
			IsBootstrap:      false,
			Multiaddr:        null.StringFrom("multiaddr"),
			ForwarderAddress: null.StringFrom("forwarderaddress"),
			P2PPeerID:        null.StringFrom("peerid"),
			KeyBundleID:      null.StringFrom("ocrkeyid"),
			Plugins: Plugins{
				Commit:     true,
				Execute:    true,
				Median:     false,
				Mercury:    true,
				Rebalancer: false,
			},
		}
		want = `{"enabled":true,"is_bootstrap":false,"multiaddr":"multiaddr","forwarder_address":"forwarderaddress","p2p_peer_id":"peerid","key_bundle_id":"ocrkeyid","plugins":{"commit":true,"execute":true,"median":false,"mercury":true,"rebalancer":false}}`
	)

	val, err := give.Value()
	require.NoError(t, err)

	actual, ok := val.([]byte)
	require.True(t, ok)

	assert.Equal(t, want, string(actual))
}

func Test_OCR2Config_Scan(t *testing.T) {
	t.Parallel()

	var (
		give = `{"enabled":true,"is_bootstrap":false,"multiaddr":"multiaddr","forwarder_address":"forwarderaddress","p2p_peer_id":"peerid","key_bundle_id":"ocrkeyid","plugins":{"commit":true,"execute":true,"median":false,"mercury":true,"rebalancer":false}}`
		want = OCR2ConfigModel{
			Enabled:          true,
			IsBootstrap:      false,
			Multiaddr:        null.StringFrom("multiaddr"),
			ForwarderAddress: null.StringFrom("forwarderaddress"),
			P2PPeerID:        null.StringFrom("peerid"),
			KeyBundleID:      null.StringFrom("ocrkeyid"),
			Plugins: Plugins{
				Commit:     true,
				Execute:    true,
				Median:     false,
				Mercury:    true,
				Rebalancer: false,
			},
		}
	)

	var actual OCR2ConfigModel
	err := actual.Scan([]byte(give))
	require.NoError(t, err)

	assert.Equal(t, want, actual)
}

func Test_JobProposal_CanEditDefinition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status SpecStatus
		want   bool
	}{
		{
			name:   "pending",
			status: SpecStatusPending,
			want:   true,
		},
		{
			name:   "cancelled",
			status: SpecStatusCancelled,
			want:   true,
		},
		{
			name:   "approved",
			status: SpecStatusApproved,
			want:   false,
		},
		{
			name:   "rejected",
			status: SpecStatusRejected,
			want:   false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jp := &JobProposalSpec{Status: tc.status}
			assert.Equal(t, tc.want, jp.CanEditDefinition())
		})
	}
}

// Test_toMetrics tests the toMetrics method
func Test_toMetrics(t *testing.T) {
	t.Parallel()

	jpCounts := JobProposalCounts{
		Cancelled: 0,
		Pending:   1,
		Approved:  2,
		Rejected:  3,
		Deleted:   4,
		Revoked:   5,
	}

	metrics := jpCounts.toMetrics()

	assert.Equal(t, metrics, map[JobProposalStatus]float64{
		JobProposalStatusCancelled: 0,
		JobProposalStatusPending:   1,
		JobProposalStatusApproved:  2,
		JobProposalStatusRejected:  3,
		JobProposalStatusDeleted:   4,
		JobProposalStatusRevoked:   5,
	})
}
