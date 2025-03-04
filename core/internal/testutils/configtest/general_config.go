package configtest

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"

	"github.com/smartcontractkit/chainlink-integrations/evm/client"
	evmclient "github.com/smartcontractkit/chainlink-integrations/evm/client"
	"github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

const DefaultPeerID = "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X"

// NewTestGeneralConfig returns a new chainlink.GeneralConfig with default test overrides and one chain with evmclient.NullClientChainID.
func NewTestGeneralConfig(t testing.TB) chainlink.GeneralConfig { return NewGeneralConfig(t, nil) }

// NewGeneralConfig returns a new chainlink.GeneralConfig with overrides.
// The default test overrides are applied before overrideFn, and include one chain with evmclient.NullClientChainID.
func NewGeneralConfig(t testing.TB, overrideFn func(*chainlink.Config, *chainlink.Secrets)) chainlink.GeneralConfig {
	tempDir := t.TempDir()
	g, err := chainlink.GeneralConfigOpts{
		OverrideFn: func(c *chainlink.Config, s *chainlink.Secrets) {
			overrides(c, s)
			c.RootDir = &tempDir
			if fn := overrideFn; fn != nil {
				fn(c, s)
			}
		},
	}.New()
	require.NoError(t, err)
	return g
}

// overrides applies some test config settings and adds a default chain with evmclient.NullClientChainID.
func overrides(c *chainlink.Config, s *chainlink.Secrets) {
	s.Password.Keystore = models.NewSecret("dummy-to-pass-validation")

	c.Insecure.OCRDevelopmentMode = ptr(true)
	c.InsecureFastScrypt = ptr(true)
	c.ShutdownGracePeriod = commonconfig.MustNewDuration(testutils.DefaultWaitTimeout)

	c.Database.DriverName = pgcommon.DriverTxWrappedPostgres
	c.Database.Lock.Enabled = ptr(false)
	c.Database.MaxIdleConns = ptr[int64](20)
	c.Database.MaxOpenConns = ptr[int64](20)
	c.Database.MigrateOnStartup = ptr(false)
	c.Database.DefaultLockTimeout = commonconfig.MustNewDuration(1 * time.Minute)

	c.JobPipeline.ReaperInterval = commonconfig.MustNewDuration(0)
	c.JobPipeline.VerboseLogging = ptr(true)

	c.Mercury.VerboseLogging = ptr(true)

	c.P2P.V2.Enabled = ptr(false)

	c.WebServer.SessionTimeout = commonconfig.MustNewDuration(2 * time.Minute)
	c.WebServer.BridgeResponseURL = commonconfig.MustParseURL("http://localhost:6688")
	testIP := net.ParseIP("127.0.0.1")
	c.WebServer.ListenIP = &testIP
	c.WebServer.TLS.ListenIP = &testIP

	chainID := big.NewI(evmclient.NullClientChainID)

	chainCfg := toml.Defaults(chainID)
	chainCfg.LogPollInterval = commonconfig.MustNewDuration(1 * time.Second) // speed it up from the standard 15s for tests

	c.EVM = append(c.EVM, &toml.EVMConfig{
		ChainID: chainID,
		Chain:   chainCfg,
		Nodes: toml.EVMNodes{
			&toml.Node{
				Name:     ptr("test"),
				WSURL:    &commonconfig.URL{},
				HTTPURL:  &commonconfig.URL{},
				SendOnly: new(bool),
				Order:    ptr[int32](100),
			},
		},
	})
}

// NewGeneralConfigSimulated returns a new chainlink.GeneralConfig with overrides, including the simulated EVM chain.
// The default test overrides are applied before overrideFn.
// The simulated chain (testutils.SimulatedChainID) replaces the null chain (evmclient.NullClientChainID).
func NewGeneralConfigSimulated(t testing.TB, overrideFn func(*chainlink.Config, *chainlink.Secrets)) chainlink.GeneralConfig {
	return NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulated(c, s)
		if fn := overrideFn; fn != nil {
			fn(c, s)
		}
	})
}

// simulated is a config override func that appends the simulated EVM chain (testutils.SimulatedChainID),
// or replaces the null chain (client.NullClientChainID) if that is the only entry.
func simulated(c *chainlink.Config, s *chainlink.Secrets) {
	chainID := big.New(testutils.SimulatedChainID)
	enabled := true
	cfg := toml.EVMConfig{
		ChainID: chainID,
		Chain:   toml.Defaults(chainID),
		Enabled: &enabled,
		Nodes:   toml.EVMNodes{&validTestNode},
	}
	if len(c.EVM) == 1 && c.EVM[0].ChainID.Cmp(big.NewI(client.NullClientChainID)) == 0 {
		c.EVM[0] = &cfg // replace null, if only entry
	} else {
		c.EVM = append(c.EVM, &cfg)
	}
}

var validTestNode = toml.Node{
	Name:     ptr("simulated-node"),
	WSURL:    commonconfig.MustParseURL("WSS://simulated-wss.com/ws"),
	HTTPURL:  commonconfig.MustParseURL("http://simulated.com"),
	SendOnly: nil,
	Order:    ptr(int32(1)),
}

func ptr[T any](v T) *T { return &v }
