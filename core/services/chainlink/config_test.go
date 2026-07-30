package chainlink

import (
	_ "embed"
	"math"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/kylelemons/godebug/diff"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	commoncfg "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/config/configtest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	mercurytransmitter "github.com/smartcontractkit/chainlink-data-streams/llo/transmitter/de"
	"github.com/smartcontractkit/chainlink-framework/multinode"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
	evmcfg "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	//go:embed testdata/config-full.toml
	fullTOML string
	//go:embed testdata/config-multi-chain.toml
	multiChainTOML string

	multiChain = Config{
		Core: toml.Core{
			RootDir: new("my/root/dir"),
			AuditLogger: toml.AuditLogger{
				Enabled:      new(true),
				ForwardToUrl: mustURL("http://localhost:9898"),
				Headers: new([]models.ServiceHeader{
					{
						Header: "Authorization",
						Value:  "token",
					},
					{
						Header: "X-SomeOther-Header",
						Value:  "value with spaces | and a bar+*",
					},
				}),
				JsonWrapperKey: new("event"),
			},
			Database: toml.Database{
				Listener: toml.DatabaseListener{
					FallbackPollInterval: commoncfg.MustNewDuration(2 * time.Minute),
				},
			},
			Log: toml.Log{
				Level:       ptr(toml.LogLevel(zapcore.PanicLevel)),
				JSONConsole: new(true),
			},
			JobPipeline: toml.JobPipeline{
				HTTPRequest: toml.JobPipelineHTTPRequest{
					DefaultTimeout: commoncfg.MustNewDuration(30 * time.Second),
				},
			},
			OCR2: toml.OCR2{
				Enabled:         new(true),
				DatabaseTimeout: commoncfg.MustNewDuration(20 * time.Second),
			},
			OCR: toml.OCR{
				Enabled:           new(true),
				BlockchainTimeout: commoncfg.MustNewDuration(5 * time.Second),
			},
			P2P: toml.P2P{
				IncomingMessageBufferSize: new(int64(999)),
			},
			AutoPprof: toml.AutoPprof{
				CPUProfileRate: new(int64(7)),
			},
			Workflows: toml.Workflows{
				Limits: toml.Limits{
					Global:   new(int32(200)),
					PerOwner: new(int32(200)),
				},
			},
		},
		EVM: []*evmcfg.EVMConfig{
			{
				ChainID: sqlutil.NewI(1),
				Chain: evmcfg.Chain{
					FinalityDepth:        new(uint32(26)),
					SafeDepth:            new(uint32(0)),
					FinalityTagEnabled:   new(true),
					SafeTagSupported:     new(true),
					FinalizedBlockOffset: new(uint32(12)),
				},
				Nodes: []*evmcfg.Node{
					{
						Name:  new("primary"),
						WSURL: mustURL("wss://web.socket/mainnet"),
					},
					{
						Name:     new("secondary"),
						HTTPURL:  mustURL("http://broadcast.mirror"),
						SendOnly: new(true),
					},
				}},
			{
				ChainID: sqlutil.NewI(42),
				Chain: evmcfg.Chain{
					GasEstimator: evmcfg.GasEstimator{
						PriceDefault: assets.NewWeiI(math.MaxInt64),
					},
				},
				Nodes: []*evmcfg.Node{
					{
						Name:  new("foo"),
						WSURL: mustURL("wss://web.socket/test/foo"),
					},
				}},
			{
				ChainID: sqlutil.NewI(137),
				Chain: evmcfg.Chain{
					GasEstimator: evmcfg.GasEstimator{
						Mode: new("FixedPrice"),
					},
				},
				Nodes: []*evmcfg.Node{
					{
						Name:  new("bar"),
						WSURL: mustURL("wss://web.socket/test/bar"),
					},
				}},
		},
	}
)

func TestConfig_Marshal(t *testing.T) {
	zeroSeconds := *commoncfg.MustNewDuration(time.Second * 0)
	second := *commoncfg.MustNewDuration(time.Second)
	minute := *commoncfg.MustNewDuration(time.Minute)
	hour := *commoncfg.MustNewDuration(time.Hour)
	mustPeerID := func(s string) *p2pkey.PeerID {
		id, err := p2pkey.MakePeerID(s)
		require.NoError(t, err)
		return &id
	}
	mustDecimal := func(s string) *decimal.Decimal {
		d, err := decimal.NewFromString(s)
		require.NoError(t, err)
		return &d
	}
	mustAddress := func(s string) *types.EIP55Address {
		a, err := types.NewEIP55Address(s)
		require.NoError(t, err)
		return &a
	}
	selectionMode := multinode.NodeSelectionModeHighestHead

	global := Config{
		Core: toml.Core{
			InsecureFastScrypt:  new(true),
			InsecurePPROFHeap:   new(true),
			RootDir:             new("test/root/dir"),
			ShutdownGracePeriod: commoncfg.MustNewDuration(10 * time.Second),
			Insecure: toml.Insecure{
				DevWebServer:         new(false),
				OCRDevelopmentMode:   new(false),
				InfiniteDepthQueries: new(false),
				DisableRateLimiting:  new(false),
			},
			Tracing: toml.Tracing{
				Enabled:         new(true),
				CollectorTarget: new("localhost:4317"),
				NodeID:          new("clc-ocr-sol-devnet-node-1"),
				SamplingRatio:   new(1.0),
				Mode:            new("tls"),
				TLSCertPath:     new("/path/to/cert.pem"),
				Attributes: map[string]string{
					"test": "load",
					"env":  "dev",
				},
			},
		},
	}

	full := global

	serviceHeaders := []models.ServiceHeader{
		{Header: "Authorization", Value: "token"},
		{Header: "X-SomeOther-Header", Value: "value with spaces | and a bar+*"},
	}
	full.AuditLogger = toml.AuditLogger{
		Enabled:        new(true),
		ForwardToUrl:   mustURL("http://localhost:9898"),
		Headers:        new(serviceHeaders),
		JsonWrapperKey: new("event"),
	}

	full.Feature = toml.Feature{
		FeedsManager:       new(true),
		LogPoller:          new(true),
		UICSAKeys:          new(true),
		CCIP:               new(true),
		MultiFeedsManagers: new(true),
	}
	full.Database = toml.Database{
		DefaultIdleInTxSessionTimeout: commoncfg.MustNewDuration(time.Minute),
		DefaultLockTimeout:            commoncfg.MustNewDuration(time.Hour),
		DefaultQueryTimeout:           commoncfg.MustNewDuration(time.Second),
		LogQueries:                    new(true),
		MigrateOnStartup:              new(true),
		MaxIdleConns:                  new(int64(7)),
		MaxOpenConns:                  new(int64(13)),
		Listener: toml.DatabaseListener{
			MaxReconnectDuration: commoncfg.MustNewDuration(time.Minute),
			MinReconnectInterval: commoncfg.MustNewDuration(5 * time.Minute),
			FallbackPollInterval: commoncfg.MustNewDuration(2 * time.Minute),
		},
		Lock: toml.DatabaseLock{
			Enabled:              new(false),
			LeaseDuration:        &minute,
			LeaseRefreshInterval: &second,
		},
		Backup: toml.DatabaseBackup{
			Dir:              new("test/backup/dir"),
			Frequency:        &hour,
			Mode:             &config.DatabaseBackupModeFull,
			OnVersionUpgrade: new(true),
		},
	}
	full.TelemetryIngress = toml.TelemetryIngress{
		UniConn:            new(false),
		Logging:            new(true),
		BufferSize:         new(uint16(1234)),
		MaxBatchSize:       new(uint16(4321)),
		SendInterval:       commoncfg.MustNewDuration(time.Minute),
		SendTimeout:        commoncfg.MustNewDuration(5 * time.Second),
		UseBatchSend:       new(true),
		ChipIngressEnabled: new(false),
		Endpoints: []toml.TelemetryIngressEndpoint{{
			Network:      new("EVM"),
			ChainID:      new("1"),
			ServerPubKey: new("test-pub-key"),
			URL:          mustURL("prom.test")},
		},
	}

	full.Log = toml.Log{
		Level:       ptr(toml.LogLevel(zapcore.DPanicLevel)),
		JSONConsole: new(true),
		UnixTS:      new(true),
		File: toml.LogFile{
			Dir:        new("log/file/dir"),
			MaxSize:    new((utils.FileSize)(100 * utils.GB)),
			MaxAgeDays: new(int64(17)),
			MaxBackups: new(int64(9)),
		},
	}
	full.WebServer = toml.WebServer{
		AuthenticationMethod:    new("local"),
		AllowOrigins:            new("*"),
		BridgeResponseURL:       mustURL("https://bridge.response"),
		BridgeCacheTTL:          commoncfg.MustNewDuration(10 * time.Second),
		HTTPWriteTimeout:        commoncfg.MustNewDuration(time.Minute),
		HTTPPort:                new(uint16(56)),
		SecureCookies:           new(true),
		SessionTimeout:          commoncfg.MustNewDuration(time.Hour),
		SessionReaperExpiration: commoncfg.MustNewDuration(7 * 24 * time.Hour),
		HTTPMaxSize:             ptr(utils.FileSize(uint64(32770))),
		StartTimeout:            commoncfg.MustNewDuration(15 * time.Second),
		ListenIP:                mustIP("192.158.1.37"),
		MFA: toml.WebServerMFA{
			RPID:     new("test-rpid"),
			RPOrigin: new("test-rp-origin"),
		},
		LDAP: toml.WebServerLDAP{
			ServerTLS:                   new(true),
			SessionTimeout:              commoncfg.MustNewDuration(15 * time.Minute),
			QueryTimeout:                commoncfg.MustNewDuration(2 * time.Minute),
			BaseUserAttr:                new("uid"),
			BaseDN:                      new("dc=custom,dc=example,dc=com"),
			UsersDN:                     new("ou=users"),
			GroupsDN:                    new("ou=groups"),
			ActiveAttribute:             new("organizationalStatus"),
			ActiveAttributeAllowedValue: new("ACTIVE"),
			AdminUserGroupCN:            new("NodeAdmins"),
			EditUserGroupCN:             new("NodeEditors"),
			RunUserGroupCN:              new("NodeRunners"),
			ReadUserGroupCN:             new("NodeReadOnly"),
			UserApiTokenEnabled:         new(false),
			UserAPITokenDuration:        commoncfg.MustNewDuration(240 * time.Hour),
			UpstreamSyncInterval:        commoncfg.MustNewDuration(0 * time.Second),
			UpstreamSyncRateLimit:       commoncfg.MustNewDuration(2 * time.Minute),
		},
		OIDC: toml.WebServerOIDC{
			ClientID:             new("abcd1234"),
			ProviderURL:          new("https://id.provider.com/oauth2/default"),
			RedirectURL:          new("http://localhost:3000/signin"),
			ClaimName:            new("groups"),
			AdminClaim:           new("NodeAdmins"),
			EditClaim:            new("NodeEditors"),
			RunClaim:             new("NodeRunners"),
			ReadClaim:            new("NodeReadOnly"),
			SessionTimeout:       commoncfg.MustNewDuration(15 * time.Minute),
			UserAPITokenEnabled:  new(false),
			UserAPITokenDuration: commoncfg.MustNewDuration(240 * time.Hour),
		},
		RateLimit: toml.WebServerRateLimit{
			Authenticated:         new(int64(42)),
			AuthenticatedPeriod:   commoncfg.MustNewDuration(time.Second),
			Unauthenticated:       new(int64(7)),
			UnauthenticatedPeriod: commoncfg.MustNewDuration(time.Minute),
		},
		TLS: toml.WebServerTLS{
			CertPath:      new("tls/cert/path"),
			Host:          new("tls-host"),
			KeyPath:       new("tls/key/path"),
			HTTPSPort:     new(uint16(6789)),
			ForceRedirect: new(true),
			ListenIP:      mustIP("192.158.1.38"),
		},
	}
	full.JobPipeline = toml.JobPipeline{
		ExternalInitiatorsEnabled: new(true),
		MaxRunDuration:            commoncfg.MustNewDuration(time.Hour),
		MaxSuccessfulRuns:         new(uint64(123456)),
		ReaperInterval:            commoncfg.MustNewDuration(4 * time.Hour),
		ReaperThreshold:           commoncfg.MustNewDuration(7 * 24 * time.Hour),
		ResultWriteQueueDepth:     new(uint32(10)),
		VerboseLogging:            new(false),
		HTTPRequest: toml.JobPipelineHTTPRequest{
			MaxSize:        new((utils.FileSize)(100 * utils.MB)),
			DefaultTimeout: commoncfg.MustNewDuration(time.Minute),
		},
	}
	full.FluxMonitor = toml.FluxMonitor{ //nolint:staticcheck // deprecated config surface must match embedded config-full.toml
		DefaultTransactionQueueDepth: new(uint32(1)),
		SimulateTransactions:         new(false),
	}
	full.OCR2 = toml.OCR2{
		Enabled:                            new(true),
		ContractConfirmations:              new(uint32(11)),
		BlockchainTimeout:                  commoncfg.MustNewDuration(3 * time.Second),
		ContractPollInterval:               commoncfg.MustNewDuration(time.Hour),
		ContractSubscribeInterval:          commoncfg.MustNewDuration(time.Minute),
		ContractTransmitterTransmitTimeout: commoncfg.MustNewDuration(time.Minute),
		DatabaseTimeout:                    commoncfg.MustNewDuration(8 * time.Second),
		KeyBundleID:                        new(corekeys.MustSha256HashFromHex("7a5f66bbe6594259325bf2b4f5b1a9c9")),
		CaptureEATelemetry:                 new(false),
		CaptureAutomationCustomTelemetry:   new(true),
		AllowNoBootstrappers:               new(true),
		DefaultTransactionQueueDepth:       new(uint32(1)),
		SimulateTransactions:               new(false),
		TraceLogging:                       new(false),
		SampleTelemetry:                    new(false),
		KeyValueStoreRootDir:               new("~/.chainlink-data"),
	}
	full.OCR = toml.OCR{
		Enabled:                      new(true),
		ObservationTimeout:           commoncfg.MustNewDuration(11 * time.Second),
		BlockchainTimeout:            commoncfg.MustNewDuration(3 * time.Second),
		ContractPollInterval:         commoncfg.MustNewDuration(time.Hour),
		ContractSubscribeInterval:    commoncfg.MustNewDuration(time.Minute),
		DefaultTransactionQueueDepth: new(uint32(12)),
		KeyBundleID:                  new(corekeys.MustSha256HashFromHex("acdd42797a8b921b2910497badc50006")),
		SimulateTransactions:         new(true),
		TransmitterAddress:           new(types.MustEIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")),
		CaptureEATelemetry:           new(false),
		TraceLogging:                 new(false),
		ConfigLogValidation:          new(false),
	}
	full.P2P = toml.P2P{
		IncomingMessageBufferSize: new(int64(13)),
		OutgoingMessageBufferSize: new(int64(17)),
		PeerID:                    mustPeerID("12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw"),
		TraceLogging:              new(true),
		EnableExperimentalRageP2P: new(true),
		V2: toml.P2PV2{
			Enabled:           new(false),
			AnnounceAddresses: &[]string{"a", "b", "c"},
			DefaultBootstrappers: &[]ocrcommontypes.BootstrapperLocator{
				{PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", Addrs: []string{"foo:42", "bar:10"}},
				{PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", Addrs: []string{"test:99"}},
			},
			DeltaDial:       commoncfg.MustNewDuration(time.Minute),
			DeltaReconcile:  commoncfg.MustNewDuration(time.Second),
			ListenAddresses: &[]string{"foo", "bar"},
		},
	}
	full.Capabilities = toml.Capabilities{
		RateLimit: toml.EngineExecutionRateLimit{
			GlobalRPS:      new(200.00),
			GlobalBurst:    new(200),
			PerSenderRPS:   new(200.0),
			PerSenderBurst: new(200),
		},

		Peering: toml.P2P{
			IncomingMessageBufferSize: new(int64(13)),
			OutgoingMessageBufferSize: new(int64(17)),
			PeerID:                    mustPeerID("12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw"),
			TraceLogging:              new(true),
			EnableExperimentalRageP2P: new(true),
			V2: toml.P2PV2{
				Enabled:           new(false),
				AnnounceAddresses: &[]string{"a", "b", "c"},
				DefaultBootstrappers: &[]ocrcommontypes.BootstrapperLocator{
					{PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", Addrs: []string{"foo:42", "bar:10"}},
					{PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", Addrs: []string{"test:99"}},
				},
				DeltaDial:       commoncfg.MustNewDuration(time.Minute),
				DeltaReconcile:  commoncfg.MustNewDuration(2 * time.Second),
				ListenAddresses: &[]string{"foo", "bar"},
			},
		},
		SharedPeering: toml.SharedPeering{
			Enabled: new(false),
			Bootstrappers: &[]ocrcommontypes.BootstrapperLocator{
				{PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", Addrs: []string{"foo:42", "bar:10"}},
				{PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", Addrs: []string{"test:99"}},
			},
			StreamConfig: toml.StreamConfig{
				IncomingMessageBufferSize:  new(500),
				OutgoingMessageBufferSize:  new(500),
				MaxMessageLenBytes:         new(500000),
				MessageRateLimiterRate:     new(100.0),
				MessageRateLimiterCapacity: new(uint32(500)),
				BytesRateLimiterRate:       new(5000000.0),
				BytesRateLimiterCapacity:   new(uint32(10000000)),
			},
		},
		ExternalRegistry: toml.ExternalRegistry{
			Address:         new(""),
			ChainID:         new("1"),
			NetworkID:       new("evm"),
			ContractVersion: new("1.0.0"),
		},
		WorkflowRegistry: toml.WorkflowRegistry{
			Address:                 new(""),
			ChainID:                 new("1"),
			ContractVersion:         new("1.0.0"),
			NetworkID:               new("evm"),
			MaxBinarySize:           ptr(utils.FileSize(20 * utils.MB)),
			MaxEncryptedSecretsSize: ptr(utils.FileSize(26.4 * utils.KB)),
			MaxConfigSize:           ptr(utils.FileSize(50 * utils.KB)),
			SyncStrategy:            new("event"),
			MaxConcurrency:          new(12),
			MaxActivationRetries:    new(100),
			WorkflowStorage: toml.WorkflowStorage{
				ArtifactStorageHost: new(""),
				URL:                 new(""),
				TLSEnabled:          new(true),
			},
			ModuleCache: toml.ModuleCache{
				Enabled:            new(false),
				DiskMonitorEnabled: new(false),
				IdleEviction:       new(true),
				IdleTimeout:        commoncfg.MustNewDuration(10 * time.Minute),
				MaxLoaded:          new(200),
				CacheDir:           new(""),
			},
			AdditionalSourcesConfig: []toml.AdditionalWorkflowSource{
				{
					URL:        new("localhost:50051"),
					TLSEnabled: new(true),
					Name:       new("test-grpc-source"),
				},
			},
		},
		Dispatcher: toml.Dispatcher{
			SupportedVersion:   new(1),
			ReceiverBufferSize: new(10000),
			RateLimit: toml.DispatcherRateLimit{
				GlobalRPS:      new(800.0),
				GlobalBurst:    new(1000),
				PerSenderRPS:   new(10.0),
				PerSenderBurst: new(50),
			},
			SendToSharedPeer: new(false),
		},
		GatewayConnector: toml.GatewayConnector{
			ChainIDForNodeKey:         new("11155111"),
			NodeAddress:               new("0x68902d681c28119f9b2531473a417088bf008e59"),
			DonID:                     new("example_don"),
			WSHandshakeTimeoutMillis:  new(uint32(100)),
			AuthMinChallengeLen:       new(10),
			AuthTimestampToleranceSec: new(uint32(10)),
			Gateways: []toml.ConnectorGateway{
				{ID: new("example_gateway"), DonID: new("example_gateway_don"), URL: new("wss://localhost:8081/node")},
			},
		},
		Local: toml.LocalCapabilities{
			RegistryBasedLaunchAllowlist: []string{`^cron@1\.0\.0$`, `^http-action@.*$`},
			Capabilities: map[string]toml.CapabilityNodeConfig{
				"http-action@1.0.0": {
					BinaryPathOverride: new("/opt/chainlink/binaries/http_action"),
					Config:             map[string]string{"proxyMode": "gateway", "allowedPorts": "443,8443"},
				},
				"cron@1.0.0": {
					BinaryPathOverride: new("/opt/chainlink/binaries/cron"),
					Config:             map[string]string{"fastestScheduleIntervalSeconds": "60"},
				},
			},
		},
	}
	full.Workflows = toml.Workflows{
		Limits: toml.Limits{
			Global:   new(int32(200)),
			PerOwner: new(int32(200)),
		},
	}
	full.AutoPprof = toml.AutoPprof{
		Enabled:              new(true),
		ProfileRoot:          new("prof/root"),
		PollInterval:         commoncfg.MustNewDuration(time.Minute),
		GatherDuration:       commoncfg.MustNewDuration(12 * time.Second),
		GatherTraceDuration:  commoncfg.MustNewDuration(13 * time.Second),
		MaxProfileSize:       new((utils.FileSize)(utils.GB)),
		CPUProfileRate:       new(int64(7)),
		MemProfileRate:       new(int64(9)),
		BlockProfileRate:     new(int64(5)),
		MutexProfileFraction: new(int64(2)),
		MemThreshold:         new((utils.FileSize)(utils.GB)),
		GoroutineThreshold:   new(int64(999)),
	}
	full.Pyroscope = toml.Pyroscope{
		ServerAddress:        new("http://localhost:4040"),
		Environment:          new("tests"),
		LinkTracesToProfiles: new(true),
	}
	full.Sentry = toml.Sentry{
		Debug:       new(true),
		DSN:         new("sentry-dsn"),
		Environment: new("dev"),
		Release:     new("v1.2.3"),
	}
	full.Telemetry = toml.Telemetry{
		Enabled:                                new(true),
		CACertFile:                             new("cert-file"),
		Endpoint:                               new("example.com/collector"),
		InsecureConnection:                     new(true),
		ResourceAttributes:                     map[string]string{"Baz": "test", "Foo": "bar"},
		TraceSampleRatio:                       new(0.01),
		EmitterBatchProcessor:                  new(true),
		EmitterExportTimeout:                   commoncfg.MustNewDuration(1 * time.Second),
		AuthHeadersTTL:                         commoncfg.MustNewDuration(0 * time.Second),
		ChipIngressEndpoint:                    new("example.com/chip-ingress"),
		ChipIngressInsecureConnection:          new(false),
		ChipIngressBatchEmitterEnabled:         new(true),
		ChipIngressBufferSize:                  new(uint(10000)),
		ChipIngressMaxBatchSize:                new(uint(1000)),
		ChipIngressMaxConcurrentSends:          new(10),
		ChipIngressSendInterval:                commoncfg.MustNewDuration(500 * time.Millisecond),
		ChipIngressSendTimeout:                 commoncfg.MustNewDuration(10 * time.Second),
		ChipIngressDrainTimeout:                commoncfg.MustNewDuration(30 * time.Second),
		ChipIngressMaxGRPCRequestSize:          new(10485760),
		DurableEmitterEnabled:                  new(false),
		DurableEmitterRetransmitBatchSize:      new(500),
		DurableEmitterEventTTL:                 commoncfg.MustNewDuration(1 * time.Hour),
		DurableEmitterMaxQueuePayloadBytes:     new(int64(1073741824)),
		DurableEmitterInsertBatchFlushInterval: commoncfg.MustNewDuration(50 * time.Millisecond),
		HeartbeatInterval:                      commoncfg.MustNewDuration(1 * time.Second),
		LogStreamingEnabled:                    new(false),
		LogLevel:                               new("info"),
		LogBatchProcessor:                      new(true),
		LogExportTimeout:                       commoncfg.MustNewDuration(1 * time.Second),
		LogExportMaxBatchSize:                  new(512),
		LogExportInterval:                      ptrDuration(1 * time.Second),
		LogMaxQueueSize:                        new(2048),
		MetricViewsDenyAttributes:              []string{"event_id"},
		MetricCardinalityLimit:                 new(100000),

		PrometheusBridge: toml.PrometheusBridge{
			Enabled:  new(true),
			Prefixes: []string{"ocr_"},
		},
	}
	full.CRE = toml.CreConfig{
		UseLocalTimeProvider: new(true),
		EnableDKGRecipient:   new(false),
		DebugMode:            new(false),
		Streams: &toml.StreamsConfig{
			WsURL:   new("streams.url"),
			RestURL: new("streams.url"),
		},
		WorkflowFetcher: &toml.WorkflowFetcherConfig{
			URL: new("https://workflow.fetcher.url"),
		},
		Linking: &toml.LinkingConfig{
			URL:            new(""),
			TLSEnabled:     new(true),
			RequestTimeout: commoncfg.MustNewDuration(2 * time.Second),
		},
		ConfidentialRelay: &toml.ConfidentialRelayConfig{
			Enabled:          new(bool),
			TrustEnclaves:    new(bool),
			RequireBFTQuorum: new(bool),
		},
	}
	full.Billing = toml.Billing{
		URL:        new("localhost:4319"),
		TLSEnabled: new(true),
	}
	full.BridgeStatusReporter = toml.BridgeStatusReporter{
		Enabled:              new(false),
		StatusPath:           new("/status"),
		PollingInterval:      commoncfg.MustNewDuration(5 * time.Minute),
		IgnoreInvalidBridges: new(true),
		IgnoreJoblessBridges: new(false),
	}
	enabledOCR2PluginTypes := []string{"median"}
	full.JobSpecReporter = toml.JobSpecReporter{
		Enabled:                new(true),
		PollingInterval:        commoncfg.MustNewDuration(time.Hour),
		EnabledOCR2PluginTypes: &enabledOCR2PluginTypes,
	}
	full.Sharding = toml.Sharding{
		ShardingEnabled:          new(false),
		ArbiterPort:              new(uint16(9876)),
		ArbiterPollInterval:      commoncfg.MustNewDuration(12 * time.Second),
		ArbiterRetryInterval:     commoncfg.MustNewDuration(12 * time.Second),
		ShardIndex:               new(uint16(0)),
		ShardOrchestratorPort:    new(uint16(50051)),
		ShardOrchestratorAddress: &commoncfg.URL{},
	}
	full.LOOPP = toml.LOOPP{
		GRPCServerMaxRecvMsgSize: new((utils.FileSize)(42 * utils.MB)),
	}
	full.JobDistributor = toml.JobDistributor{
		DisplayName: new("test-node"),
	}
	full.EVM = []*evmcfg.EVMConfig{
		{
			ChainID: sqlutil.NewI(1),
			Enabled: new(false),
			Chain: evmcfg.Chain{
				AutoCreateKey: new(false),
				BalanceMonitor: evmcfg.BalanceMonitor{
					Enabled: new(true),
				},
				BlockBackfillDepth:   new(uint32(100)),
				BlockBackfillSkip:    new(true),
				ChainType:            chaintype.NewConfig("Optimism"),
				FinalityDepth:        new(uint32(42)),
				SafeDepth:            new(uint32(0)),
				FinalityTagEnabled:   new(true),
				SafeTagSupported:     new(true),
				FlagsContractAddress: mustAddress("0xae4E781a6218A8031764928E88d457937A954fC3"),
				FinalizedBlockOffset: new(uint32(16)),

				GasEstimator: evmcfg.GasEstimator{
					Mode:               new("SuggestedPrice"),
					EIP1559DynamicFees: new(true),
					BumpPercent:        new(uint16(10)),
					BumpThreshold:      new(uint32(6)),
					BumpTxDepth:        new(uint32(6)),
					BumpMin:            assets.NewWeiI(100),
					FeeCapDefault:      assets.NewWeiI(math.MaxInt64),
					LimitDefault:       new(uint64(12)),
					LimitMax:           new(uint64(17)),
					LimitMultiplier:    mustDecimal("1.234"),
					LimitTransfer:      new(uint64(100)),
					EstimateLimit:      new(false),
					TipCapDefault:      assets.NewWeiI(2),
					TipCapMin:          assets.NewWeiI(1),
					PriceDefault:       assets.NewWeiI(math.MaxInt64),
					PriceMax:           assets.NewWei(mustHexToBig(t, "FFFFFFFFFFFF")),
					PriceMin:           assets.NewWeiI(13),

					LimitJobType: evmcfg.GasLimitJobType{
						OCR:    new(uint32(1001)),
						DR:     new(uint32(1002)),
						VRF:    new(uint32(1003)),
						FM:     new(uint32(1004)),
						Keeper: new(uint32(1005)),
						OCR2:   new(uint32(1006)),
					},

					BlockHistory: evmcfg.BlockHistoryEstimator{
						BatchSize:                 new(uint32(17)),
						BlockHistorySize:          new(uint16(12)),
						CheckInclusionBlocks:      new(uint16(18)),
						CheckInclusionPercentile:  new(uint16(19)),
						EIP1559FeeCapBufferBlocks: new(uint16(13)),
						TransactionPercentile:     new(uint16(15)),
					},
					FeeHistory: evmcfg.FeeHistoryEstimator{
						CacheTimeout: &second,
					},
				},

				KeySpecific: []evmcfg.KeySpecific{
					{
						Key: mustAddress("0x2a3e23c6f242F5345320814aC8a1b4E58707D292"),
						GasEstimator: evmcfg.KeySpecificGasEstimator{
							PriceMax: assets.NewWei(mustHexToBig(t, "FFFFFFFFFFFFFFFFFFFFFFFF")),
						},
					},
				},

				LinkContractAddress:          mustAddress("0x538aAaB4ea120b2bC2fe5D296852D948F07D849e"),
				LogBackfillBatchSize:         new(uint32(17)),
				LogPollInterval:              &minute,
				LogKeepBlocksDepth:           new(uint32(100000)),
				LogPollerSkipEmptyBlocks:     new(false),
				LogPrunePageSize:             new(uint32(0)),
				BackupLogPollerBlockDelay:    new(uint64(532)),
				MinContractPayment:           commonassets.NewLinkFromJuels(math.MaxInt64),
				MinIncomingConfirmations:     new(uint32(13)),
				NonceAutoSync:                new(true),
				NoNewHeadsThreshold:          &minute,
				OperatorFactoryAddress:       mustAddress("0xa5B85635Be42F21f94F28034B7DA440EeFF0F418"),
				LogBroadcasterEnabled:        new(true),
				RPCDefaultBatchSize:          new(uint32(17)),
				RPCBlockQueryDelay:           new(uint16(10)),
				NoNewFinalizedHeadsThreshold: &hour,

				Transactions: evmcfg.Transactions{
					Enabled:              new(true),
					MaxInFlight:          new(uint32(19)),
					MaxQueued:            new(uint32(99)),
					ReaperInterval:       &minute,
					ReaperThreshold:      &minute,
					ResendAfterThreshold: &hour,
					ForwardersEnabled:    new(true),
					AutoPurge: evmcfg.AutoPurgeConfig{
						Enabled: new(false),
					},
					TransactionManagerV2: evmcfg.TransactionManagerV2Config{
						Enabled:                     new(false),
						ReadRequestsToMultipleNodes: new(false),
						Bundles:                     new(false),
					},
					ConfirmationTimeout: &minute,
				},

				HeadTracker: evmcfg.HeadTracker{
					HistoryDepth:            new(uint32(15)),
					MaxBufferSize:           new(uint32(17)),
					SamplingInterval:        &hour,
					FinalityTagBypass:       new(false),
					MaxAllowedFinalityDepth: new(uint32(1500)),
					PersistenceEnabled:      new(false),
					PersistenceBatchSize:    new(int64(100)),
				},

				NodePool: evmcfg.NodePool{
					PollFailureThreshold:                new(uint32(5)),
					PollSuccessThreshold:                new(uint32(0)),
					PollInterval:                        &minute,
					SelectionMode:                       &selectionMode,
					SyncThreshold:                       new(uint32(13)),
					LeaseDuration:                       &zeroSeconds,
					NodeIsSyncingEnabled:                new(true),
					FinalizedBlockPollInterval:          &second,
					HistoricalBalanceCheckAddress:       new(types.MustEIP55Address("0x0000000000000000000000000000000000000000")),
					EnforceRepeatableRead:               new(true),
					DeathDeclarationDelay:               &minute,
					VerifyChainID:                       new(true),
					NewHeadsPollInterval:                &zeroSeconds,
					ExternalRequestMaxResponseSize:      new(uint32(10)),
					FinalizedStateCheckFailureThreshold: new(uint32(0)),
					Errors: evmcfg.ClientErrors{
						NonceTooLow:                       new("(: |^)nonce too low"),
						NonceTooHigh:                      new("(: |^)nonce too high"),
						ReplacementTransactionUnderpriced: new("(: |^)replacement transaction underpriced"),
						LimitReached:                      new("(: |^)limit reached"),
						TransactionAlreadyInMempool:       new("(: |^)transaction already in mempool"),
						TerminallyUnderpriced:             new("(: |^)terminally underpriced"),
						InsufficientEth:                   new("(: |^)insufficient eth"),
						TxFeeExceedsCap:                   new("(: |^)tx fee exceeds cap"),
						L2FeeTooLow:                       new("(: |^)l2 fee too low"),
						L2FeeTooHigh:                      new("(: |^)l2 fee too high"),
						L2Full:                            new("(: |^)l2 full"),
						TransactionAlreadyMined:           new("(: |^)transaction already mined"),
						Fatal:                             new("(: |^)fatal"),
						ServiceUnavailable:                new("(: |^)service unavailable"),
						TooManyResults:                    new("(: |^)too many results"),
						MissingBlocks:                     new("(: |^)missing blocks"),
						FinalizedStateUnavailable:         new("(: |^)(missing trie node|state not available|historical state unavailable)"),
					},
				},
				OCR: evmcfg.OCR{
					ContractConfirmations:              new(uint16(11)),
					ContractTransmitterTransmitTimeout: &minute,
					DatabaseTimeout:                    &second,
					DeltaCOverride:                     commoncfg.MustNewDuration(time.Hour),
					DeltaCJitterOverride:               commoncfg.MustNewDuration(time.Second),
					ObservationGracePeriod:             &second,
				},
				OCR2: evmcfg.OCR2{
					Automation: evmcfg.Automation{
						GasLimit: new(uint32(540)),
					},
				},
				Workflow: evmcfg.Workflow{
					GasLimitDefault:   new(uint64(400000)),
					TxAcceptanceState: ptr(commontypes.Unconfirmed),
					PollPeriod:        commoncfg.MustNewDuration(time.Second * 2),
					AcceptanceTimeout: commoncfg.MustNewDuration(time.Second * 30),
				},
			},
			Nodes: []*evmcfg.Node{
				{
					Name:              new("foo"),
					HTTPURL:           mustURL("https://foo.web"),
					WSURL:             mustURL("wss://web.socket/test/foo"),
					HTTPURLExtraWrite: mustURL("https://foo.web/extra"),
				},
				{
					Name:    new("bar"),
					HTTPURL: mustURL("https://bar.com"),
					WSURL:   mustURL("wss://web.socket/test/bar"),
				},
				{
					Name:     new("broadcast"),
					HTTPURL:  mustURL("http://broadcast.mirror"),
					SendOnly: new(true),
				},
			}},
	}
	full.Mercury = toml.Mercury{
		Cache: toml.MercuryCache{
			LatestReportTTL:      commoncfg.MustNewDuration(100 * time.Second),
			MaxStaleAge:          commoncfg.MustNewDuration(101 * time.Second),
			LatestReportDeadline: commoncfg.MustNewDuration(102 * time.Second),
		},
		TLS: toml.MercuryTLS{
			CertFile: new("/path/to/cert.pem"),
		},
		Transmitter: toml.MercuryTransmitter{
			Protocol:             ptr(mercurytransmitter.MercuryTransmitterProtocolGRPC),
			TransmitQueueMaxSize: new(uint32(123)),
			TransmitTimeout:      commoncfg.MustNewDuration(234 * time.Second),
			TransmitConcurrency:  new(uint32(456)),
			ReaperFrequency:      commoncfg.MustNewDuration(567 * time.Second),
			ReaperMaxAge:         commoncfg.MustNewDuration(678 * time.Hour),
		},
		VerboseLogging: new(true),
	}

	for _, tt := range []struct {
		name   string
		config Config
		exp    string
	}{
		{"empty", Config{}, ``},
		{"global", global, `InsecureFastScrypt = true
InsecurePPROFHeap = true
RootDir = 'test/root/dir'
ShutdownGracePeriod = '10s'

[Insecure]
DevWebServer = false
OCRDevelopmentMode = false
InfiniteDepthQueries = false
DisableRateLimiting = false

[Tracing]
Enabled = true
CollectorTarget = 'localhost:4317'
NodeID = 'clc-ocr-sol-devnet-node-1'
SamplingRatio = 1.0
Mode = 'tls'
TLSCertPath = '/path/to/cert.pem'

[Tracing.Attributes]
env = 'dev'
test = 'load'
`},
		{"AuditLogger", Config{Core: toml.Core{AuditLogger: full.AuditLogger}}, `[AuditLogger]
Enabled = true
ForwardToUrl = 'http://localhost:9898'
JsonWrapperKey = 'event'
Headers = ['Authorization: token', 'X-SomeOther-Header: value with spaces | and a bar+*']
`},
		{"Feature", Config{Core: toml.Core{Feature: full.Feature}}, `[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true
CCIP = true
MultiFeedsManagers = true
`},
		{"Database", Config{Core: toml.Core{Database: full.Database}}, `[Database]
DefaultIdleInTxSessionTimeout = '1m0s'
DefaultLockTimeout = '1h0m0s'
DefaultQueryTimeout = '1s'
LogQueries = true
MaxIdleConns = 7
MaxOpenConns = 13
MigrateOnStartup = true

[Database.Backup]
Dir = 'test/backup/dir'
Frequency = '1h0m0s'
Mode = 'full'
OnVersionUpgrade = true

[Database.Listener]
MaxReconnectDuration = '1m0s'
MinReconnectInterval = '5m0s'
FallbackPollInterval = '2m0s'

[Database.Lock]
Enabled = false
LeaseDuration = '1m0s'
LeaseRefreshInterval = '1s'
`},
		{"TelemetryIngress", Config{Core: toml.Core{TelemetryIngress: full.TelemetryIngress}}, `[TelemetryIngress]
UniConn = false
Logging = true
BufferSize = 1234
MaxBatchSize = 4321
SendInterval = '1m0s'
SendTimeout = '5s'
UseBatchSend = true
ChipIngressEnabled = false

[[TelemetryIngress.Endpoints]]
Network = 'EVM'
ChainID = '1'
URL = 'prom.test'
ServerPubKey = 'test-pub-key'
`},

		{"Log", Config{Core: toml.Core{Log: full.Log}}, `[Log]
Level = 'crit'
JSONConsole = true
UnixTS = true

[Log.File]
Dir = 'log/file/dir'
MaxSize = '100.00gb'
MaxAgeDays = 17
MaxBackups = 9
`},
		{"WebServer", Config{Core: toml.Core{WebServer: full.WebServer}}, `[WebServer]
AuthenticationMethod = 'local'
AllowOrigins = '*'
BridgeResponseURL = 'https://bridge.response'
BridgeCacheTTL = '10s'
HTTPWriteTimeout = '1m0s'
HTTPPort = 56
SecureCookies = true
SessionTimeout = '1h0m0s'
SessionReaperExpiration = '168h0m0s'
HTTPMaxSize = '32.77kb'
StartTimeout = '15s'
ListenIP = '192.158.1.37'

[WebServer.LDAP]
ServerTLS = true
SessionTimeout = '15m0s'
QueryTimeout = '2m0s'
BaseUserAttr = 'uid'
BaseDN = 'dc=custom,dc=example,dc=com'
UsersDN = 'ou=users'
GroupsDN = 'ou=groups'
ActiveAttribute = 'organizationalStatus'
ActiveAttributeAllowedValue = 'ACTIVE'
AdminUserGroupCN = 'NodeAdmins'
EditUserGroupCN = 'NodeEditors'
RunUserGroupCN = 'NodeRunners'
ReadUserGroupCN = 'NodeReadOnly'
UserApiTokenEnabled = false
UserAPITokenDuration = '240h0m0s'
UpstreamSyncInterval = '0s'
UpstreamSyncRateLimit = '2m0s'

[WebServer.OIDC]
ClientID = 'abcd1234'
ProviderURL = 'https://id.provider.com/oauth2/default'
RedirectURL = 'http://localhost:3000/signin'
ClaimName = 'groups'
AdminClaim = 'NodeAdmins'
EditClaim = 'NodeEditors'
RunClaim = 'NodeRunners'
ReadClaim = 'NodeReadOnly'
SessionTimeout = '15m0s'
UserAPITokenEnabled = false
UserAPITokenDuration = '240h0m0s'

[WebServer.MFA]
RPID = 'test-rpid'
RPOrigin = 'test-rp-origin'

[WebServer.RateLimit]
Authenticated = 42
AuthenticatedPeriod = '1s'
Unauthenticated = 7
UnauthenticatedPeriod = '1m0s'

[WebServer.TLS]
CertPath = 'tls/cert/path'
ForceRedirect = true
Host = 'tls-host'
HTTPSPort = 6789
KeyPath = 'tls/key/path'
ListenIP = '192.158.1.38'
`},
		{"JobPipeline", Config{Core: toml.Core{JobPipeline: full.JobPipeline}}, `[JobPipeline]
ExternalInitiatorsEnabled = true
MaxRunDuration = '1h0m0s'
MaxSuccessfulRuns = 123456
ReaperInterval = '4h0m0s'
ReaperThreshold = '168h0m0s'
ResultWriteQueueDepth = 10
VerboseLogging = false

[JobPipeline.HTTPRequest]
DefaultTimeout = '1m0s'
MaxSize = '100.00mb'
`},
		{"OCR", Config{Core: toml.Core{OCR: full.OCR}}, `[OCR]
Enabled = true
ObservationTimeout = '11s'
BlockchainTimeout = '3s'
ContractPollInterval = '1h0m0s'
ContractSubscribeInterval = '1m0s'
DefaultTransactionQueueDepth = 12
KeyBundleID = 'acdd42797a8b921b2910497badc5000600000000000000000000000000000000'
SimulateTransactions = true
TransmitterAddress = '0xa0788FC17B1dEe36f057c42B6F373A34B014687e'
CaptureEATelemetry = false
TraceLogging = false
ConfigLogValidation = false
`},
		{"OCR2", Config{Core: toml.Core{OCR2: full.OCR2}}, `[OCR2]
Enabled = true
ContractConfirmations = 11
BlockchainTimeout = '3s'
ContractPollInterval = '1h0m0s'
ContractSubscribeInterval = '1m0s'
ContractTransmitterTransmitTimeout = '1m0s'
DatabaseTimeout = '8s'
KeyBundleID = '7a5f66bbe6594259325bf2b4f5b1a9c900000000000000000000000000000000'
CaptureEATelemetry = false
CaptureAutomationCustomTelemetry = true
AllowNoBootstrappers = true
DefaultTransactionQueueDepth = 1
SimulateTransactions = false
TraceLogging = false
SampleTelemetry = false
KeyValueStoreRootDir = '~/.chainlink-data'
`},
		{"JobDistributor", Config{Core: toml.Core{JobDistributor: full.JobDistributor}}, `[JobDistributor]
DisplayName = 'test-node'
`},
		{"P2P", Config{Core: toml.Core{P2P: full.P2P}}, `[P2P]
IncomingMessageBufferSize = 13
OutgoingMessageBufferSize = 17
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw'
TraceLogging = true
EnableExperimentalRageP2P = true

[P2P.V2]
Enabled = false
AnnounceAddresses = ['a', 'b', 'c']
DefaultBootstrappers = ['12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw@foo:42/bar:10', '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw@test:99']
DeltaDial = '1m0s'
DeltaReconcile = '1s'
ListenAddresses = ['foo', 'bar']
`},
		{"AutoPprof", Config{Core: toml.Core{AutoPprof: full.AutoPprof}}, `[AutoPprof]
Enabled = true
ProfileRoot = 'prof/root'
PollInterval = '1m0s'
GatherDuration = '12s'
GatherTraceDuration = '13s'
MaxProfileSize = '1.00gb'
CPUProfileRate = 7
MemProfileRate = 9
BlockProfileRate = 5
MutexProfileFraction = 2
MemThreshold = '1.00gb'
GoroutineThreshold = 999
`},
		{"Pyroscope", Config{Core: toml.Core{Pyroscope: full.Pyroscope}}, `[Pyroscope]
ServerAddress = 'http://localhost:4040'
Environment = 'tests'
LinkTracesToProfiles = true
`},
		{"Sentry", Config{Core: toml.Core{Sentry: full.Sentry}}, `[Sentry]
Debug = true
DSN = 'sentry-dsn'
Environment = 'dev'
Release = 'v1.2.3'
`},
		{"EVM", Config{EVM: full.EVM}, `[[EVM]]
ChainID = '1'
Enabled = false
AutoCreateKey = false
BlockBackfillDepth = 100
BlockBackfillSkip = true
ChainType = 'Optimism'
FinalityDepth = 42
SafeDepth = 0
FinalityTagEnabled = true
SafeTagSupported = true
FlagsContractAddress = '0xae4E781a6218A8031764928E88d457937A954fC3'
LinkContractAddress = '0x538aAaB4ea120b2bC2fe5D296852D948F07D849e'
LogBackfillBatchSize = 17
LogPollInterval = '1m0s'
LogPollerSkipEmptyBlocks = false
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 532
MinIncomingConfirmations = 13
MinContractPayment = '9.223372036854775807 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
OperatorFactoryAddress = '0xa5B85635Be42F21f94F28034B7DA440EeFF0F418'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 17
RPCBlockQueryDelay = 10
FinalizedBlockOffset = 16
NoNewFinalizedHeadsThreshold = '1h0m0s'

[EVM.Transactions]
Enabled = true
ForwardersEnabled = true
MaxInFlight = 19
MaxQueued = 99
ReaperInterval = '1m0s'
ReaperThreshold = '1m0s'
ResendAfterThreshold = '1h0m0s'
ConfirmationTimeout = '1m0s'

[EVM.Transactions.AutoPurge]
Enabled = false

[EVM.Transactions.TransactionManagerV2]
Enabled = false
ReadRequestsToMultipleNodes = false
Bundles = false

[EVM.BalanceMonitor]
Enabled = true

[EVM.GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '9.223372036854775807 ether'
PriceMax = '281.474976710655 micro'
PriceMin = '13 wei'
LimitDefault = 12
LimitMax = 17
LimitMultiplier = '1.234'
LimitTransfer = 100
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 10
BumpThreshold = 6
BumpTxDepth = 6
EIP1559DynamicFees = true
FeeCapDefault = '9.223372036854775807 ether'
TipCapDefault = '2 wei'
TipCapMin = '1 wei'

[EVM.GasEstimator.LimitJobType]
OCR = 1001
OCR2 = 1006
DR = 1002
VRF = 1003
FM = 1004
Keeper = 1005

[EVM.GasEstimator.BlockHistory]
BatchSize = 17
BlockHistorySize = 12
CheckInclusionBlocks = 18
CheckInclusionPercentile = 19
EIP1559FeeCapBufferBlocks = 13
TransactionPercentile = 15

[EVM.GasEstimator.FeeHistory]
CacheTimeout = '1s'

[EVM.HeadTracker]
HistoryDepth = 15
MaxBufferSize = 17
SamplingInterval = '1h0m0s'
MaxAllowedFinalityDepth = 1500
FinalityTagBypass = false
PersistenceEnabled = false
PersistenceBatchSize = 100

[[EVM.KeySpecific]]
Key = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292'

[EVM.KeySpecific.GasEstimator]
PriceMax = '79.228162514264337593543950335 gether'

[EVM.NodePool]
PollFailureThreshold = 5
PollSuccessThreshold = 0
PollInterval = '1m0s'
SelectionMode = 'HighestHead'
SyncThreshold = 13
LeaseDuration = '0s'
NodeIsSyncingEnabled = true
FinalizedBlockPollInterval = '1s'
HistoricalBalanceCheckAddress = '0x0000000000000000000000000000000000000000'
FinalizedStateCheckFailureThreshold = 0
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'
VerifyChainID = true
ExternalRequestMaxResponseSize = 10

[EVM.NodePool.Errors]
NonceTooLow = '(: |^)nonce too low'
NonceTooHigh = '(: |^)nonce too high'
ReplacementTransactionUnderpriced = '(: |^)replacement transaction underpriced'
LimitReached = '(: |^)limit reached'
TransactionAlreadyInMempool = '(: |^)transaction already in mempool'
TerminallyUnderpriced = '(: |^)terminally underpriced'
InsufficientEth = '(: |^)insufficient eth'
TxFeeExceedsCap = '(: |^)tx fee exceeds cap'
L2FeeTooLow = '(: |^)l2 fee too low'
L2FeeTooHigh = '(: |^)l2 fee too high'
L2Full = '(: |^)l2 full'
TransactionAlreadyMined = '(: |^)transaction already mined'
Fatal = '(: |^)fatal'
ServiceUnavailable = '(: |^)service unavailable'
TooManyResults = '(: |^)too many results'
MissingBlocks = '(: |^)missing blocks'
FinalizedStateUnavailable = '(: |^)(missing trie node|state not available|historical state unavailable)'

[EVM.OCR]
ContractConfirmations = 11
ContractTransmitterTransmitTimeout = '1m0s'
DatabaseTimeout = '1s'
DeltaCOverride = '1h0m0s'
DeltaCJitterOverride = '1s'
ObservationGracePeriod = '1s'

[EVM.OCR2]
[EVM.OCR2.Automation]
GasLimit = 540

[EVM.Workflow]
GasLimitDefault = 400000
TxAcceptanceState = 2
PollPeriod = '2s'
AcceptanceTimeout = '30s'

[[EVM.Nodes]]
Name = 'foo'
WSURL = 'wss://web.socket/test/foo'
HTTPURL = 'https://foo.web'
HTTPURLExtraWrite = 'https://foo.web/extra'

[[EVM.Nodes]]
Name = 'bar'
WSURL = 'wss://web.socket/test/bar'
HTTPURL = 'https://bar.com'

[[EVM.Nodes]]
Name = 'broadcast'
HTTPURL = 'http://broadcast.mirror'
SendOnly = true
`},
		{"Mercury", Config{Core: toml.Core{Mercury: full.Mercury}}, `[Mercury]
VerboseLogging = true

[Mercury.Cache]
LatestReportTTL = '1m40s'
MaxStaleAge = '1m41s'
LatestReportDeadline = '1m42s'

[Mercury.TLS]
CertFile = '/path/to/cert.pem'

[Mercury.Transmitter]
Protocol = 'grpc'
TransmitQueueMaxSize = 123
TransmitTimeout = '3m54s'
TransmitConcurrency = 456
ReaperFrequency = '9m27s'
ReaperMaxAge = '678h0m0s'
`},
		{"full", full, fullTOML},
		{"multi-chain", multiChain, multiChainTOML},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.config.TOMLString()
			require.NoError(t, err)
			assert.Equal(t, tt.exp, s, diff.Diff(tt.exp, s))

			var got Config

			require.NoError(t, commoncfg.DecodeTOML(strings.NewReader(s), &got))
			ts, err := got.TOMLString()

			require.NoError(t, err)
			assert.Equal(t, tt.config, got, diff.Diff(s, ts))
		})
	}
}

func TestConfig_full(t *testing.T) {
	var got Config
	require.NoError(t, commoncfg.DecodeTOML(strings.NewReader(fullTOML), &got))
	// Except for some EVM node fields.
	for c := range got.EVM {
		addr, err := types.NewEIP55Address("0x2a3e23c6f242F5345320814aC8a1b4E58707D292")
		require.NoError(t, err)
		if got.EVM[c].Workflow.FromAddress == nil {
			got.EVM[c].Workflow.FromAddress = &addr
		}
		if got.EVM[c].Workflow.ForwarderAddress == nil {
			got.EVM[c].Workflow.ForwarderAddress = &addr
		}
		if got.EVM[c].Workflow.GasLimitDefault == nil {
			got.EVM[c].Workflow.GasLimitDefault = new(uint64(400000))
		}
		for n := range got.EVM[c].Nodes {
			if got.EVM[c].Nodes[n].WSURL == nil {
				got.EVM[c].Nodes[n].WSURL = new(commoncfg.URL)
			}
			if got.EVM[c].Nodes[n].SendOnly == nil {
				got.EVM[c].Nodes[n].SendOnly = new(true)
			}
			if got.EVM[c].Nodes[n].Order == nil {
				got.EVM[c].Nodes[n].Order = new(int32(100))
			}
			if got.EVM[c].Nodes[n].HTTPURLExtraWrite == nil {
				got.EVM[c].Nodes[n].HTTPURLExtraWrite = new(commoncfg.URL)
			}
			if got.EVM[c].Nodes[n].IsLoadBalancedRPC == nil {
				got.EVM[c].Nodes[n].IsLoadBalancedRPC = new(false)
			}
		}
		if got.EVM[c].Transactions.TransactionManagerV2.BlockTime == nil {
			got.EVM[c].Transactions.TransactionManagerV2.BlockTime = new(commoncfg.Duration)
		}
		if got.EVM[c].Transactions.TransactionManagerV2.CustomURL == nil {
			got.EVM[c].Transactions.TransactionManagerV2.CustomURL = new(commoncfg.URL)
		}
		if got.EVM[c].Transactions.TransactionManagerV2.DualBroadcast == nil {
			got.EVM[c].Transactions.TransactionManagerV2.DualBroadcast = new(false)
		}
		if got.EVM[c].Transactions.TransactionManagerV2.ReadRequestsToMultipleNodes == nil {
			got.EVM[c].Transactions.TransactionManagerV2.ReadRequestsToMultipleNodes = new(false)
		}
		if got.EVM[c].Transactions.TransactionManagerV2.Bundles == nil {
			got.EVM[c].Transactions.TransactionManagerV2.Bundles = new(false)
		}
		if got.EVM[c].Transactions.TransactionManagerV2.FastlaneAuctionRequestTimeout == nil {
			got.EVM[c].Transactions.TransactionManagerV2.FastlaneAuctionRequestTimeout = new(commoncfg.Duration)
		}
		if got.EVM[c].Transactions.TransactionManagerV2.FeeBoost == nil {
			got.EVM[c].Transactions.TransactionManagerV2.FeeBoost = new(false)
		}
		if got.EVM[c].Transactions.AutoPurge.Threshold == nil {
			got.EVM[c].Transactions.AutoPurge.Threshold = new(uint32(0))
		}
		if got.EVM[c].Transactions.AutoPurge.MinAttempts == nil {
			got.EVM[c].Transactions.AutoPurge.MinAttempts = new(uint32(0))
		}
		if got.EVM[c].Transactions.AutoPurge.DetectionApiUrl == nil {
			got.EVM[c].Transactions.AutoPurge.DetectionApiUrl = new(commoncfg.URL)
		}
		if got.EVM[c].GasEstimator.DAOracle.OracleType == nil {
			oracleType := evmcfg.DAOracleOPStack
			got.EVM[c].GasEstimator.DAOracle.OracleType = &oracleType
		}
		if got.EVM[c].GasEstimator.DAOracle.OracleAddress == nil {
			got.EVM[c].GasEstimator.DAOracle.OracleAddress = new(types.EIP55Address)
		}
		if got.EVM[c].GasEstimator.DAOracle.CustomGasPriceCalldata == nil {
			got.EVM[c].GasEstimator.DAOracle.CustomGasPriceCalldata = new(string)
		}
		if got.EVM[c].GasEstimator.SenderAddress == nil {
			got.EVM[c].GasEstimator.SenderAddress = new(types.EIP55Address)
		}
	}

	configtest.AssertFieldsNotNil(t, got)
}

//go:embed testdata/config-invalid.toml
var invalidTOML string

func TestConfig_Validate(t *testing.T) {
	for _, tt := range []struct {
		name string
		toml string
		exp  string
	}{
		{name: "invalid", toml: invalidTOML, exp: `invalid configuration: 10 errors:
	- P2P.V2.Enabled: invalid value (false): P2P required for OCR or OCR2. Please enable P2P or disable OCR/OCR2.
	- Database.Lock.LeaseRefreshInterval: invalid value (6s): must be less than or equal to half of LeaseDuration (10s)
	- WebServer: 8 errors:
		- LDAP.BaseDN: invalid value (<nil>): LDAP BaseDN can not be empty
		- LDAP.BaseUserAttr: invalid value (<nil>): LDAP BaseUserAttr can not be empty
		- LDAP.UsersDN: invalid value (<nil>): LDAP UsersDN can not be empty
		- LDAP.GroupsDN: invalid value (<nil>): LDAP GroupsDN can not be empty
		- LDAP.AdminUserGroupCN: invalid value (<nil>): LDAP AdminUserGroupCN can not be empty
		- LDAP.RunUserGroupCN: invalid value (<nil>): LDAP ReadUserGroupCN can not be empty
		- LDAP.RunUserGroupCN: invalid value (<nil>): LDAP RunUserGroupCN can not be empty
		- LDAP.ReadUserGroupCN: invalid value (<nil>): LDAP ReadUserGroupCN can not be empty
	- EVM: 10 errors:
		- 1.ChainID: invalid value (1): duplicate - must be unique
		- 0.Nodes.1.Name: invalid value (foo): duplicate - must be unique
		- 3.Nodes.4.WSURL: invalid value (ws://dupe.com): duplicate - must be unique
		- 0: 4 errors:
			- Nodes: missing: 0th node (primary) must have a valid WSURL when LogBroadcaster is enabled
			- GasEstimator.BumpTxDepth: invalid value (11): must be less than or equal to Transactions.MaxInFlight
			- GasEstimator: 6 errors:
				- BumpPercent: invalid value (1): may not be less than Geth's default of 10
				- TipCapDefault: invalid value (3 wei): must be greater than or equal to TipCapMinimum
				- FeeCapDefault: invalid value (3 wei): must be greater than or equal to TipCapDefault
				- PriceMin: invalid value (10 gwei): must be less than or equal to PriceDefault
				- PriceMax: invalid value (10 gwei): must be greater than or equal to PriceDefault
				- BlockHistory.BlockHistorySize: invalid value (0): must be greater than or equal to 1 with BlockHistory Mode
			- Nodes: 2 errors:
				- 0.HTTPURL: missing: required for all nodes
				- 1.HTTPURL: missing: required for all nodes
		- 1: 10 errors:
			- ChainType: invalid value (Foo): must not be set with this chain id
			- Nodes: missing: must have at least one node
			- ChainType: invalid value (Foo): must be one of arbitrum, astar, celo, gnosis, hedera, kroma, mantle, metis, optimismBedrock, sei, scroll, wemix, xlayer, zkevm, zksync, zircuit, tron, rootstock, pharos, jovay or omitted
			- HeadTracker.HistoryDepth: invalid value (30): must be greater than or equal to FinalizedBlockOffset
			- GasEstimator.BumpThreshold: invalid value (0): cannot be 0 if auto-purge feature is enabled for Foo
			- Transactions.AutoPurge.Threshold: missing: needs to be set if auto-purge feature is enabled for Foo
			- Transactions.AutoPurge.MinAttempts: missing: needs to be set if auto-purge feature is enabled for Foo
			- GasEstimator: 2 errors:
				- FeeCapDefault: invalid value (101 wei): must be equal to PriceMax (99 wei) since you are using FixedPrice estimation with gas bumping disabled in EIP1559 mode - PriceMax will be used as the FeeCap for transactions instead of FeeCapDefault
				- PriceMax: invalid value (1 gwei): must be greater than or equal to PriceDefault
			- HeadTracker.MaxAllowedFinalityDepth: invalid value (0): must be greater than or equal to 1
			- KeySpecific.Key: invalid value (0xde709f2102306220921060314715629080e2fb77): duplicate - must be unique
		- 2: 5 errors:
			- ChainType: invalid value (Arbitrum): only "optimismBedrock" can be used with this chain id
			- Nodes: missing: must have at least one node
			- ChainType: invalid value (Arbitrum): must be one of arbitrum, astar, celo, gnosis, hedera, kroma, mantle, metis, optimismBedrock, sei, scroll, wemix, xlayer, zkevm, zksync, zircuit, tron, rootstock, pharos, jovay or omitted
			- FinalityDepth: invalid value (0): must be greater than or equal to 1
			- MinIncomingConfirmations: invalid value (0): must be greater than or equal to 1
		- 3: 3 errors:
			- Nodes: missing: 0th node (primary) must have a valid WSURL when LogBroadcaster is enabled
			- Nodes: missing: 2th node (primary) must have a valid WSURL when LogBroadcaster is enabled
			- Nodes: 5 errors:
				- 0: 2 errors:
					- Name: missing: required for all nodes
					- HTTPURL: empty: required for all nodes
				- 1: 3 errors:
					- Name: missing: required for all nodes
					- WSURL: invalid value (http): must be ws or wss
					- HTTPURL: missing: required for all nodes
				- 2: 2 errors:
					- Name: empty: required for all nodes
					- HTTPURL: invalid value (ws): must be http or https
				- 3.HTTPURL: missing: required for all nodes
				- 4.HTTPURL: missing: required for all nodes
		- 4: 2 errors:
			- ChainID: missing: required for all chains
			- Nodes: missing: must have at least one node
		- 5.Transactions.AutoPurge.DetectionApiUrl: invalid value (): must be set for scroll
		- 6.Nodes: missing: 0th node (primary) must have a valid WSURL when http polling is disabled
	- Cosmos: 4 errors:
		- 1.ChainID: invalid value (Malaga-420): duplicate - must be unique
		- 0.Nodes.1.Name: invalid value (test): duplicate - must be unique
		- 1.Nodes: missing: expected at least one node
		- 2: 2 errors:
			- ChainID: missing: required for all chains
			- Nodes: missing: expected at least one node
	- Solana: 4 errors:
		- 1.ChainID: invalid value (mainnet): duplicate - must be unique
		- 1.Nodes.1.Name: invalid value (bar): duplicate - must be unique
		- 0.Nodes: missing: expected at least one node
		- 2: 2 errors:
			- ChainID: missing: required for all chains
			- Nodes: missing: expected at least one node
	- Starknet: 3 errors:
		- 0.Nodes.1.Name: invalid value (primary): duplicate - must be unique
		- 0.ChainID: missing: required for all chains
		- 1: 2 errors:
			- ChainID: missing: required for all chains
			- Nodes: missing: expected at least one node
	- Aptos: 2 errors:
		- 0.Nodes.1.Name: invalid value (primary): duplicate - must be unique
		- 0: 2 errors:
			- Enabled: invalid value (1): expected bool
			- ChainID: missing: required for all chains
	- Tron: 2 errors:
		- 0.Nodes.1.Name: invalid value (tron-test): duplicate - must be unique
		- 0: 2 errors:
			- Enabled: invalid value (1): expected bool
			- ChainID: missing: required for all chains
	- TON: 2 errors:
		- 0.Nodes.1.Name: invalid value (ton-test): duplicate - must be unique
		- 0: 2 errors:
			- Enabled: invalid value (1): expected bool
			- ChainID: missing: required for all chains`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var c Config
			require.NoError(t, commoncfg.DecodeTOML(strings.NewReader(tt.toml), &c))
			c.setDefaults()
			assertValidationError(t, &c, tt.exp)
		})
	}
}

func mustURL(s string) *commoncfg.URL {
	var u commoncfg.URL
	if err := u.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return &u
}

func mustIP(s string) *net.IP {
	var ip net.IP
	if err := ip.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return &ip
}

var (
	//go:embed testdata/secrets-empty-effective.toml
	emptyEffectiveSecretsTOML string
	//go:embed testdata/config-empty-effective.toml
	emptyEffectiveTOML string
	//go:embed testdata/config-multi-chain-effective.toml
	multiChainEffectiveTOML string

	//go:embed testdata/secrets-full.toml
	secretsFullTOML string
	//go:embed testdata/secrets-full-redacted.toml
	secretsFullRedactedTOML string

	//go:embed testdata/secrets-multi.toml
	secretsMultiTOML string
	//go:embed testdata/secrets-multi-redacted.toml
	secretsMultiRedactedTOML string
)

func Test_generalConfig_LogConfiguration(t *testing.T) {
	const (
		secrets   = "# Secrets:\n"
		input     = "# Input Configuration:\n"
		effective = "# Effective Configuration, with defaults applied:\n"
		warning   = "# Configuration warning:\n"

		deprecated = `` // none
	)
	tests := []struct {
		name         string
		inputConfig  string
		inputSecrets string

		wantConfig    string
		wantEffective string
		wantSecrets   string
		wantWarning   string
	}{
		{name: "empty", wantEffective: emptyEffectiveTOML, wantSecrets: emptyEffectiveSecretsTOML},
		{name: "full", inputSecrets: secretsFullTOML, inputConfig: fullTOML,
			wantConfig: fullTOML, wantEffective: fullTOML, wantSecrets: secretsFullRedactedTOML, wantWarning: deprecated},
		{name: "multi-chain", inputSecrets: secretsMultiTOML, inputConfig: multiChainTOML,
			wantConfig: multiChainTOML, wantEffective: multiChainEffectiveTOML, wantSecrets: secretsMultiRedactedTOML},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lggr, observed := logger.TestLoggerObserved(t, zapcore.InfoLevel)
			opts := GeneralConfigOpts{
				SkipEnv:        true,
				ConfigStrings:  []string{tt.inputConfig},
				SecretsStrings: []string{tt.inputSecrets},
			}
			c, err := opts.New()
			require.NoError(t, err)
			c.LogConfiguration(lggr.Infof, lggr.Warnf)

			inputLogs := observed.FilterMessageSnippet(secrets).All()
			if assert.Len(t, inputLogs, 1) {
				assert.Equal(t, zapcore.InfoLevel, inputLogs[0].Level)
				got := strings.TrimPrefix(inputLogs[0].Message, secrets)
				got = strings.TrimSuffix(got, "\n")
				assert.Equal(t, tt.wantSecrets, got)
			}

			inputLogs = observed.FilterMessageSnippet(input).All()
			if assert.Len(t, inputLogs, 1) {
				assert.Equal(t, zapcore.InfoLevel, inputLogs[0].Level)
				got := strings.TrimPrefix(inputLogs[0].Message, input)
				got = strings.TrimSuffix(got, "\n")
				assert.Equal(t, tt.wantConfig, got)
			}

			inputLogs = observed.FilterMessageSnippet(effective).All()
			if assert.Len(t, inputLogs, 1) {
				assert.Equal(t, zapcore.InfoLevel, inputLogs[0].Level)
				got := strings.TrimPrefix(inputLogs[0].Message, effective)
				got = strings.TrimSuffix(got, "\n")
				assert.Equal(t, tt.wantEffective, got)
			}

			inputLogs = observed.FilterMessageSnippet(warning).All()
			if tt.wantWarning != "" && assert.Len(t, inputLogs, 1) {
				assert.Equal(t, zapcore.WarnLevel, inputLogs[0].Level)
				got := strings.TrimPrefix(inputLogs[0].Message, warning)
				got = strings.TrimSuffix(got, "\n")
				assert.Equal(t, tt.wantWarning, got)
			}
		})
	}
}

func TestNewGeneralConfig_ParsingError_InvalidSyntax(t *testing.T) {
	invalidTOML := "{ bad syntax {"
	opts := GeneralConfigOpts{
		ConfigStrings:  []string{invalidTOML},
		SecretsStrings: []string{secretsFullTOML},
	}
	_, err := opts.New()
	assert.EqualError(t, err, "failed to decode config TOML: toml: invalid character at start of key: {")
}

func TestNewGeneralConfig_ParsingError_DuplicateField(t *testing.T) {
	invalidTOML := `Dev = false
Dev = true`
	opts := GeneralConfigOpts{
		ConfigStrings:  []string{invalidTOML},
		SecretsStrings: []string{secretsFullTOML},
	}
	_, err := opts.New()
	assert.EqualError(t, err, "failed to decode config TOML: toml: key Dev is already defined")
}

func TestNewGeneralConfig_SecretsOverrides(t *testing.T) {
	// Provide a keystore password file and an env var with DB URL
	const pwdOverride = "great_password"
	const dbURLOverride = "http://user@db"

	t.Setenv("CL_DATABASE_URL", dbURLOverride)

	// Check for two overrides
	opts := GeneralConfigOpts{
		ConfigStrings:  []string{fullTOML},
		SecretsStrings: []string{secretsFullTOML},
	}
	c, err := opts.New()
	require.NoError(t, err)
	c.SetPasswords(ptr(pwdOverride), nil)
	assert.Equal(t, pwdOverride, c.Password().Keystore())
	dbURL := c.Database().URL()
	assert.Equal(t, dbURLOverride, (&dbURL).String())
}

func TestSecrets_Validate(t *testing.T) {
	for _, tt := range []struct {
		name string
		toml string
		exp  string
	}{
		{name: "partial",
			toml: `
Database.AllowSimplePasswords = true`,
			exp: `invalid secrets: 2 errors:
	- Database.URL: empty: must be provided and non-empty
	- Password.Keystore: empty: must be provided and non-empty`},

		{name: "invalid-urls",
			toml: `[Database]
URL = "postgresql://user:passlocalhost:5432/asdf"
BackupURL = "foo-bar?password=asdf"
AllowSimplePasswords = false`,
			exp: `invalid secrets: 2 errors:
	- Database: 2 errors:
		- URL: invalid value (*****): missing or insufficiently complex password: DB URL must be authenticated; plaintext URLs are not allowed. Database should be secured by a password matching the following complexity requirements: 
	Must have a length of 16-50 characters
	Must not comprise:
		Leading or trailing whitespace (note that a trailing newline in the password file, if present, will be ignored)
	
		- BackupURL: invalid value (*****): missing or insufficiently complex password: 
	Expected password complexity:
	Must be at least 16 characters long
	Must not comprise:
		Leading or trailing whitespace
		A user's API email
	
	Faults:
		password is less than 16 characters long
	. Database should be secured by a password matching the following complexity requirements: 
	Must have a length of 16-50 characters
	Must not comprise:
		Leading or trailing whitespace (note that a trailing newline in the password file, if present, will be ignored)
	
	- Password.Keystore: empty: must be provided and non-empty`},

		{name: "invalid-urls-allowed",
			toml: `[Database]
URL = "postgresql://user:passlocalhost:5432/asdf"
BackupURL = "foo-bar?password=asdf"
AllowSimplePasswords = true`,
			exp: `invalid secrets: Password.Keystore: empty: must be provided and non-empty`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var s Secrets
			require.NoError(t, commoncfg.DecodeTOML(strings.NewReader(tt.toml), &s))
			assertValidationError(t, &s, tt.exp)
		})
	}
}

func assertValidationError(t *testing.T, invalid interface{ Validate() error }, expMsg string) {
	t.Helper()
	if err := invalid.Validate(); assert.Error(t, err) {
		got := err.Error()
		assert.Equal(t, expMsg, got, diff.Diff(expMsg, got))
	}
}

func TestConfig_setDefaults(t *testing.T) {
	var c Config
	c.EVM = evmcfg.EVMConfigs{{ChainID: sqlutil.NewI(99999133712345)}}
	c.Cosmos = RawConfigs{{"ChainID": new("unknown cosmos chain")}}
	c.Solana = RawConfigs{{"ChainID": new("unknown solana chain")}}
	c.Starknet = RawConfigs{{"ChainID": new("unknown starknet chain")}}
	c.setDefaults()

	s, err := c.TOMLString()
	require.NoError(t, err)
	t.Log(s)

	configtest.AssertFieldsNotNil(t, c.Core)
}

func Test_validateEnv(t *testing.T) {
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("DATABASE_URL", "foo")
	require.ErrorContains(t, validateEnv(), `invalid environment: 2 errors:
	- environment variable DATABASE_URL must not be set: unsupported with config v2
	- environment variable LOG_LEVEL must not be set: unsupported with config v2`)

	t.Setenv("GAS_UPDATER_ENABLED", "true")
	t.Setenv("ETH_GAS_BUMP_TX_DEPTH", "7")
	require.ErrorContains(t, validateEnv(), `invalid environment: 4 errors:
	- environment variable DATABASE_URL must not be set: unsupported with config v2
	- environment variable LOG_LEVEL must not be set: unsupported with config v2
	- environment variable ETH_GAS_BUMP_TX_DEPTH must not be set: unsupported with config v2
	- environment variable GAS_UPDATER_ENABLED must not be set: unsupported with config v2`)
}

func TestConfig_SetFrom(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		exp  string
		from []string
	}{
		{"empty", "", []string{""}},
		{"empty-full", fullTOML, []string{"", fullTOML}},
		{"empty-multi", multiChainTOML, []string{"", multiChainTOML}},
		{"full-empty", fullTOML, []string{fullTOML, ""}},
		{"multi-empty", multiChainTOML, []string{multiChainTOML, ""}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var c Config
			for _, fs := range tt.from {
				var f Config
				require.NoError(t, commoncfg.DecodeTOML(strings.NewReader(fs), &f))
				require.NoError(t, c.SetFrom(&f))
			}
			ts, err := c.TOMLString()

			require.NoError(t, err)
			assert.Equal(t, tt.exp, ts)
		})
	}
}

func TestConfig_warnings(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		expectedErrors []string
	}{
		{
			name:           "No warnings",
			config:         Config{},
			expectedErrors: nil,
		},
		{
			name: "Value warning - unencrypted mode with TLS path set",
			config: Config{
				Core: toml.Core{
					Tracing: toml.Tracing{
						Enabled:     new(true),
						Mode:        new("unencrypted"),
						TLSCertPath: new("/path/to/cert.pem"),
					},
				},
			},
			expectedErrors: []string{"Tracing.TLSCertPath: invalid value (/path/to/cert.pem): must be empty when Tracing.Mode is 'unencrypted'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.warnings()
			if len(tt.expectedErrors) == 0 {
				assert.NoError(t, err)
			} else {
				for _, expectedErr := range tt.expectedErrors {
					assert.Contains(t, err.Error(), expectedErr)
				}
			}
		})
	}
}

//go:fix inline
func ptr[T any](t T) *T { return new(t) }

func mustHexToBig(t *testing.T, hx string) *big.Int {
	n, err := hex.ParseBig(hx)
	require.NoError(t, err)
	return n
}

func TestRawConfig_IsEnabled(t *testing.T) {
	assert.True(t, RawConfig{"Enabled": true}.IsEnabled())
	assert.True(t, RawConfig{"Enabled": nil}.IsEnabled())
	assert.True(t, RawConfig{}.IsEnabled())

	assert.False(t, RawConfig{"Enabled": false}.IsEnabled())
	assert.False(t, RawConfig{"Enabled": "garbage"}.IsEnabled())
}

func TestRawConfig_SetDefaults(t *testing.T) {
	c := RawConfig{"Enabled": true}
	c.SetDefaults()
	require.NotContains(t, c, "Enabled")
	c["Enabled"] = false
	c.SetDefaults()
	require.Contains(t, c, "Enabled")
	require.Equal(t, false, c["Enabled"])
}
