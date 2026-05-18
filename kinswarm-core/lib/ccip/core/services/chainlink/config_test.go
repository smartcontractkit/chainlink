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
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	commoncfg "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"

	"github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	legacy "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
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
			RootDir: ptr("my/root/dir"),
			AuditLogger: toml.AuditLogger{
				Enabled:      ptr(true),
				ForwardToUrl: mustURL("http://localhost:9898"),
				Headers: ptr([]models.ServiceHeader{
					{
						Header: "Authorization",
						Value:  "token",
					},
					{
						Header: "X-SomeOther-Header",
						Value:  "value with spaces | and a bar+*",
					},
				}),
				JsonWrapperKey: ptr("event"),
			},
			Database: toml.Database{
				Listener: toml.DatabaseListener{
					FallbackPollInterval: commoncfg.MustNewDuration(2 * time.Minute),
				},
			},
			Log: toml.Log{
				Level:       ptr(toml.LogLevel(zapcore.PanicLevel)),
				JSONConsole: ptr(true),
			},
			JobPipeline: toml.JobPipeline{
				HTTPRequest: toml.JobPipelineHTTPRequest{
					DefaultTimeout: commoncfg.MustNewDuration(30 * time.Second),
				},
			},
			OCR2: toml.OCR2{
				Enabled:         ptr(true),
				DatabaseTimeout: commoncfg.MustNewDuration(20 * time.Second),
			},
			OCR: toml.OCR{
				Enabled:           ptr(true),
				BlockchainTimeout: commoncfg.MustNewDuration(5 * time.Second),
			},
			P2P: toml.P2P{
				IncomingMessageBufferSize: ptr[int64](999),
			},
			Keeper: toml.Keeper{
				GasPriceBufferPercent: ptr[uint16](10),
			},
			AutoPprof: toml.AutoPprof{
				CPUProfileRate: ptr[int64](7),
			},
		},
		EVM: []*evmcfg.EVMConfig{
			{
				ChainID: ubig.NewI(1),
				Chain: evmcfg.Chain{
					FinalityDepth:        ptr[uint32](26),
					FinalityTagEnabled:   ptr[bool](false),
					FinalizedBlockOffset: ptr[uint32](12),
				},
				Nodes: []*evmcfg.Node{
					{
						Name:  ptr("primary"),
						WSURL: mustURL("wss://web.socket/mainnet"),
					},
					{
						Name:     ptr("secondary"),
						HTTPURL:  mustURL("http://broadcast.mirror"),
						SendOnly: ptr(true),
					},
				}},
			{
				ChainID: ubig.NewI(42),
				Chain: evmcfg.Chain{
					GasEstimator: evmcfg.GasEstimator{
						PriceDefault: assets.NewWeiI(math.MaxInt64),
					},
				},
				Nodes: []*evmcfg.Node{
					{
						Name:  ptr("foo"),
						WSURL: mustURL("wss://web.socket/test/foo"),
					},
				}},
			{
				ChainID: ubig.NewI(137),
				Chain: evmcfg.Chain{
					GasEstimator: evmcfg.GasEstimator{
						Mode: ptr("FixedPrice"),
					},
				},
				Nodes: []*evmcfg.Node{
					{
						Name:  ptr("bar"),
						WSURL: mustURL("wss://web.socket/test/bar"),
					},
				}},
		},
		Cosmos: []*coscfg.TOMLConfig{
			{
				ChainID: ptr("Ibiza-808"),
				Chain: coscfg.Chain{
					MaxMsgsPerBatch: ptr[int64](13),
				},
				Nodes: []*coscfg.Node{
					{Name: ptr("primary"), TendermintURL: commoncfg.MustParseURL("http://columbus.cosmos.com")},
				}},
			{
				ChainID: ptr("Malaga-420"),
				Chain: coscfg.Chain{
					BlocksUntilTxTimeout: ptr[int64](20),
				},
				Nodes: []*coscfg.Node{
					{Name: ptr("secondary"), TendermintURL: commoncfg.MustParseURL("http://bombay.cosmos.com")},
				}},
		},
		Solana: []*solcfg.TOMLConfig{
			{
				ChainID: ptr("mainnet"),
				Chain: solcfg.Chain{
					MaxRetries: ptr[int64](12),
				},
				Nodes: []*solcfg.Node{
					{Name: ptr("primary"), URL: commoncfg.MustParseURL("http://mainnet.solana.com")},
				},
			},
			{
				ChainID: ptr("testnet"),
				Chain: solcfg.Chain{
					OCR2CachePollPeriod: commoncfg.MustNewDuration(time.Minute),
				},
				Nodes: []*solcfg.Node{
					{Name: ptr("secondary"), URL: commoncfg.MustParseURL("http://testnet.solana.com")},
				},
			},
		},
		Starknet: []*stkcfg.TOMLConfig{
			{
				ChainID: ptr("foobar"),
				Chain: stkcfg.Chain{
					ConfirmationPoll: commoncfg.MustNewDuration(time.Hour),
				},
				FeederURL: commoncfg.MustParseURL("http://feeder.url"),
				Nodes: []*stkcfg.Node{
					{Name: ptr("primary"), URL: commoncfg.MustParseURL("http://stark.node"), APIKey: ptr("key")},
				},
			},
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
	selectionMode := client.NodeSelectionModeHighestHead

	global := Config{
		Core: toml.Core{
			InsecureFastScrypt:  ptr(true),
			RootDir:             ptr("test/root/dir"),
			ShutdownGracePeriod: commoncfg.MustNewDuration(10 * time.Second),
			Insecure: toml.Insecure{
				DevWebServer:         ptr(false),
				OCRDevelopmentMode:   ptr(false),
				InfiniteDepthQueries: ptr(false),
				DisableRateLimiting:  ptr(false),
			},
			Tracing: toml.Tracing{
				Enabled:         ptr(true),
				CollectorTarget: ptr("localhost:4317"),
				NodeID:          ptr("clc-ocr-sol-devnet-node-1"),
				SamplingRatio:   ptr(1.0),
				Mode:            ptr("tls"),
				TLSCertPath:     ptr("/path/to/cert.pem"),
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
		Enabled:        ptr(true),
		ForwardToUrl:   mustURL("http://localhost:9898"),
		Headers:        ptr(serviceHeaders),
		JsonWrapperKey: ptr("event"),
	}

	full.Feature = toml.Feature{
		FeedsManager:       ptr(true),
		LogPoller:          ptr(true),
		UICSAKeys:          ptr(true),
		CCIP:               ptr(true),
		MultiFeedsManagers: ptr(true),
	}
	full.Database = toml.Database{
		DefaultIdleInTxSessionTimeout: commoncfg.MustNewDuration(time.Minute),
		DefaultLockTimeout:            commoncfg.MustNewDuration(time.Hour),
		DefaultQueryTimeout:           commoncfg.MustNewDuration(time.Second),
		LogQueries:                    ptr(true),
		MigrateOnStartup:              ptr(true),
		MaxIdleConns:                  ptr[int64](7),
		MaxOpenConns:                  ptr[int64](13),
		Listener: toml.DatabaseListener{
			MaxReconnectDuration: commoncfg.MustNewDuration(time.Minute),
			MinReconnectInterval: commoncfg.MustNewDuration(5 * time.Minute),
			FallbackPollInterval: commoncfg.MustNewDuration(2 * time.Minute),
		},
		Lock: toml.DatabaseLock{
			Enabled:              ptr(false),
			LeaseDuration:        &minute,
			LeaseRefreshInterval: &second,
		},
		Backup: toml.DatabaseBackup{
			Dir:              ptr("test/backup/dir"),
			Frequency:        &hour,
			Mode:             &legacy.DatabaseBackupModeFull,
			OnVersionUpgrade: ptr(true),
		},
	}
	full.TelemetryIngress = toml.TelemetryIngress{
		UniConn:      ptr(false),
		Logging:      ptr(true),
		BufferSize:   ptr[uint16](1234),
		MaxBatchSize: ptr[uint16](4321),
		SendInterval: commoncfg.MustNewDuration(time.Minute),
		SendTimeout:  commoncfg.MustNewDuration(5 * time.Second),
		UseBatchSend: ptr(true),
		Endpoints: []toml.TelemetryIngressEndpoint{{
			Network:      ptr("EVM"),
			ChainID:      ptr("1"),
			ServerPubKey: ptr("test-pub-key"),
			URL:          mustURL("prom.test")},
		},
	}

	full.Log = toml.Log{
		Level:       ptr(toml.LogLevel(zapcore.DPanicLevel)),
		JSONConsole: ptr(true),
		UnixTS:      ptr(true),
		File: toml.LogFile{
			Dir:        ptr("log/file/dir"),
			MaxSize:    ptr[utils.FileSize](100 * utils.GB),
			MaxAgeDays: ptr[int64](17),
			MaxBackups: ptr[int64](9),
		},
	}
	full.WebServer = toml.WebServer{
		AuthenticationMethod:    ptr("local"),
		AllowOrigins:            ptr("*"),
		BridgeResponseURL:       mustURL("https://bridge.response"),
		BridgeCacheTTL:          commoncfg.MustNewDuration(10 * time.Second),
		HTTPWriteTimeout:        commoncfg.MustNewDuration(time.Minute),
		HTTPPort:                ptr[uint16](56),
		SecureCookies:           ptr(true),
		SessionTimeout:          commoncfg.MustNewDuration(time.Hour),
		SessionReaperExpiration: commoncfg.MustNewDuration(7 * 24 * time.Hour),
		HTTPMaxSize:             ptr(utils.FileSize(uint64(32770))),
		StartTimeout:            commoncfg.MustNewDuration(15 * time.Second),
		ListenIP:                mustIP("192.158.1.37"),
		MFA: toml.WebServerMFA{
			RPID:     ptr("test-rpid"),
			RPOrigin: ptr("test-rp-origin"),
		},
		LDAP: toml.WebServerLDAP{
			ServerTLS:                   ptr(true),
			SessionTimeout:              commoncfg.MustNewDuration(15 * time.Minute),
			QueryTimeout:                commoncfg.MustNewDuration(2 * time.Minute),
			BaseUserAttr:                ptr("uid"),
			BaseDN:                      ptr("dc=custom,dc=example,dc=com"),
			UsersDN:                     ptr("ou=users"),
			GroupsDN:                    ptr("ou=groups"),
			ActiveAttribute:             ptr("organizationalStatus"),
			ActiveAttributeAllowedValue: ptr("ACTIVE"),
			AdminUserGroupCN:            ptr("NodeAdmins"),
			EditUserGroupCN:             ptr("NodeEditors"),
			RunUserGroupCN:              ptr("NodeRunners"),
			ReadUserGroupCN:             ptr("NodeReadOnly"),
			UserApiTokenEnabled:         ptr(false),
			UserAPITokenDuration:        commoncfg.MustNewDuration(240 * time.Hour),
			UpstreamSyncInterval:        commoncfg.MustNewDuration(0 * time.Second),
			UpstreamSyncRateLimit:       commoncfg.MustNewDuration(2 * time.Minute),
		},
		RateLimit: toml.WebServerRateLimit{
			Authenticated:         ptr[int64](42),
			AuthenticatedPeriod:   commoncfg.MustNewDuration(time.Second),
			Unauthenticated:       ptr[int64](7),
			UnauthenticatedPeriod: commoncfg.MustNewDuration(time.Minute),
		},
		TLS: toml.WebServerTLS{
			CertPath:      ptr("tls/cert/path"),
			Host:          ptr("tls-host"),
			KeyPath:       ptr("tls/key/path"),
			HTTPSPort:     ptr[uint16](6789),
			ForceRedirect: ptr(true),
			ListenIP:      mustIP("192.158.1.38"),
		},
	}
	full.JobPipeline = toml.JobPipeline{
		ExternalInitiatorsEnabled: ptr(true),
		MaxRunDuration:            commoncfg.MustNewDuration(time.Hour),
		MaxSuccessfulRuns:         ptr[uint64](123456),
		ReaperInterval:            commoncfg.MustNewDuration(4 * time.Hour),
		ReaperThreshold:           commoncfg.MustNewDuration(7 * 24 * time.Hour),
		ResultWriteQueueDepth:     ptr[uint32](10),
		VerboseLogging:            ptr(false),
		HTTPRequest: toml.JobPipelineHTTPRequest{
			MaxSize:        ptr[utils.FileSize](100 * utils.MB),
			DefaultTimeout: commoncfg.MustNewDuration(time.Minute),
		},
	}
	full.FluxMonitor = toml.FluxMonitor{
		DefaultTransactionQueueDepth: ptr[uint32](100),
		SimulateTransactions:         ptr(true),
	}
	full.OCR2 = toml.OCR2{
		Enabled:                            ptr(true),
		ContractConfirmations:              ptr[uint32](11),
		BlockchainTimeout:                  commoncfg.MustNewDuration(3 * time.Second),
		ContractPollInterval:               commoncfg.MustNewDuration(time.Hour),
		ContractSubscribeInterval:          commoncfg.MustNewDuration(time.Minute),
		ContractTransmitterTransmitTimeout: commoncfg.MustNewDuration(time.Minute),
		DatabaseTimeout:                    commoncfg.MustNewDuration(8 * time.Second),
		KeyBundleID:                        ptr(models.MustSha256HashFromHex("7a5f66bbe6594259325bf2b4f5b1a9c9")),
		CaptureEATelemetry:                 ptr(false),
		CaptureAutomationCustomTelemetry:   ptr(true),
		DefaultTransactionQueueDepth:       ptr[uint32](1),
		SimulateTransactions:               ptr(false),
		TraceLogging:                       ptr(false),
	}
	full.OCR = toml.OCR{
		Enabled:                      ptr(true),
		ObservationTimeout:           commoncfg.MustNewDuration(11 * time.Second),
		BlockchainTimeout:            commoncfg.MustNewDuration(3 * time.Second),
		ContractPollInterval:         commoncfg.MustNewDuration(time.Hour),
		ContractSubscribeInterval:    commoncfg.MustNewDuration(time.Minute),
		DefaultTransactionQueueDepth: ptr[uint32](12),
		KeyBundleID:                  ptr(models.MustSha256HashFromHex("acdd42797a8b921b2910497badc50006")),
		SimulateTransactions:         ptr(true),
		TransmitterAddress:           ptr(types.MustEIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")),
		CaptureEATelemetry:           ptr(false),
		TraceLogging:                 ptr(false),
	}
	full.P2P = toml.P2P{
		IncomingMessageBufferSize: ptr[int64](13),
		OutgoingMessageBufferSize: ptr[int64](17),
		PeerID:                    mustPeerID("12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw"),
		TraceLogging:              ptr(true),
		V2: toml.P2PV2{
			Enabled:           ptr(false),
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
		Peering: toml.P2P{
			IncomingMessageBufferSize: ptr[int64](13),
			OutgoingMessageBufferSize: ptr[int64](17),
			PeerID:                    mustPeerID("12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw"),
			TraceLogging:              ptr(true),
			V2: toml.P2PV2{
				Enabled:           ptr(false),
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
		ExternalRegistry: toml.ExternalRegistry{
			Address:   ptr(""),
			ChainID:   ptr("1"),
			NetworkID: ptr("evm"),
		},
		Dispatcher: toml.Dispatcher{
			SupportedVersion:   ptr(1),
			ReceiverBufferSize: ptr(10000),
			RateLimit: toml.DispatcherRateLimit{
				GlobalRPS:      ptr(800.0),
				GlobalBurst:    ptr(1000),
				PerSenderRPS:   ptr(10.0),
				PerSenderBurst: ptr(50),
			},
		},
		GatewayConnector: toml.GatewayConnector{
			ChainIDForNodeKey:         ptr("11155111"),
			NodeAddress:               ptr("0x68902d681c28119f9b2531473a417088bf008e59"),
			DonID:                     ptr("example_don"),
			WSHandshakeTimeoutMillis:  ptr[uint32](100),
			AuthMinChallengeLen:       ptr[int](10),
			AuthTimestampToleranceSec: ptr[uint32](10),
			Gateways: []toml.ConnectorGateway{
				{ID: ptr("example_gateway"), URL: ptr("wss://localhost:8081/node")},
			},
		},
	}
	full.Keeper = toml.Keeper{
		DefaultTransactionQueueDepth: ptr[uint32](17),
		GasPriceBufferPercent:        ptr[uint16](12),
		GasTipCapBufferPercent:       ptr[uint16](43),
		BaseFeeBufferPercent:         ptr[uint16](89),
		MaxGracePeriod:               ptr[int64](31),
		TurnLookBack:                 ptr[int64](91),
		Registry: toml.KeeperRegistry{
			CheckGasOverhead:    ptr[uint32](90),
			PerformGasOverhead:  ptr[uint32](math.MaxUint32),
			SyncInterval:        commoncfg.MustNewDuration(time.Hour),
			SyncUpkeepQueueSize: ptr[uint32](31),
			MaxPerformDataSize:  ptr[uint32](5000),
		},
	}
	full.AutoPprof = toml.AutoPprof{
		Enabled:              ptr(true),
		ProfileRoot:          ptr("prof/root"),
		PollInterval:         commoncfg.MustNewDuration(time.Minute),
		GatherDuration:       commoncfg.MustNewDuration(12 * time.Second),
		GatherTraceDuration:  commoncfg.MustNewDuration(13 * time.Second),
		MaxProfileSize:       ptr[utils.FileSize](utils.GB),
		CPUProfileRate:       ptr[int64](7),
		MemProfileRate:       ptr[int64](9),
		BlockProfileRate:     ptr[int64](5),
		MutexProfileFraction: ptr[int64](2),
		MemThreshold:         ptr[utils.FileSize](utils.GB),
		GoroutineThreshold:   ptr[int64](999),
	}
	full.Pyroscope = toml.Pyroscope{
		ServerAddress: ptr("http://localhost:4040"),
		Environment:   ptr("tests"),
	}
	full.Sentry = toml.Sentry{
		Debug:       ptr(true),
		DSN:         ptr("sentry-dsn"),
		Environment: ptr("dev"),
		Release:     ptr("v1.2.3"),
	}
	full.Telemetry = toml.Telemetry{
		Enabled:            ptr(true),
		CACertFile:         ptr("cert-file"),
		Endpoint:           ptr("example.com/collector"),
		InsecureConnection: ptr(true),
		ResourceAttributes: map[string]string{"Baz": "test", "Foo": "bar"},
		TraceSampleRatio:   ptr(0.01),
	}
	full.EVM = []*evmcfg.EVMConfig{
		{
			ChainID: ubig.NewI(1),
			Enabled: ptr(false),
			Chain: evmcfg.Chain{
				AutoCreateKey: ptr(false),
				BalanceMonitor: evmcfg.BalanceMonitor{
					Enabled: ptr(true),
				},
				BlockBackfillDepth:   ptr[uint32](100),
				BlockBackfillSkip:    ptr(true),
				ChainType:            chaintype.NewConfig("Optimism"),
				FinalityDepth:        ptr[uint32](42),
				FinalityTagEnabled:   ptr[bool](false),
				FlagsContractAddress: mustAddress("0xae4E781a6218A8031764928E88d457937A954fC3"),
				FinalizedBlockOffset: ptr[uint32](16),

				GasEstimator: evmcfg.GasEstimator{
					Mode:               ptr("SuggestedPrice"),
					EIP1559DynamicFees: ptr(true),
					BumpPercent:        ptr[uint16](10),
					BumpThreshold:      ptr[uint32](6),
					BumpTxDepth:        ptr[uint32](6),
					BumpMin:            assets.NewWeiI(100),
					FeeCapDefault:      assets.NewWeiI(math.MaxInt64),
					LimitDefault:       ptr[uint64](12),
					LimitMax:           ptr[uint64](17),
					LimitMultiplier:    mustDecimal("1.234"),
					LimitTransfer:      ptr[uint64](100),
					EstimateLimit:      ptr(false),
					TipCapDefault:      assets.NewWeiI(2),
					TipCapMin:          assets.NewWeiI(1),
					PriceDefault:       assets.NewWeiI(math.MaxInt64),
					PriceMax:           assets.NewWei(mustHexToBig(t, "FFFFFFFFFFFF")),
					PriceMin:           assets.NewWeiI(13),

					LimitJobType: evmcfg.GasLimitJobType{
						OCR:    ptr[uint32](1001),
						DR:     ptr[uint32](1002),
						VRF:    ptr[uint32](1003),
						FM:     ptr[uint32](1004),
						Keeper: ptr[uint32](1005),
						OCR2:   ptr[uint32](1006),
					},

					BlockHistory: evmcfg.BlockHistoryEstimator{
						BatchSize:                 ptr[uint32](17),
						BlockHistorySize:          ptr[uint16](12),
						CheckInclusionBlocks:      ptr[uint16](18),
						CheckInclusionPercentile:  ptr[uint16](19),
						EIP1559FeeCapBufferBlocks: ptr[uint16](13),
						TransactionPercentile:     ptr[uint16](15),
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
				LogBackfillBatchSize:         ptr[uint32](17),
				LogPollInterval:              &minute,
				LogKeepBlocksDepth:           ptr[uint32](100000),
				LogPrunePageSize:             ptr[uint32](10000),
				BackupLogPollerBlockDelay:    ptr[uint64](532),
				MinContractPayment:           commonassets.NewLinkFromJuels(math.MaxInt64),
				MinIncomingConfirmations:     ptr[uint32](13),
				NonceAutoSync:                ptr(true),
				NoNewHeadsThreshold:          &minute,
				OperatorFactoryAddress:       mustAddress("0xa5B85635Be42F21f94F28034B7DA440EeFF0F418"),
				LogBroadcasterEnabled:        ptr(true),
				RPCDefaultBatchSize:          ptr[uint32](17),
				RPCBlockQueryDelay:           ptr[uint16](10),
				NoNewFinalizedHeadsThreshold: &hour,

				Transactions: evmcfg.Transactions{
					MaxInFlight:          ptr[uint32](19),
					MaxQueued:            ptr[uint32](99),
					ReaperInterval:       &minute,
					ReaperThreshold:      &minute,
					ResendAfterThreshold: &hour,
					ForwardersEnabled:    ptr(true),
					AutoPurge: evmcfg.AutoPurgeConfig{
						Enabled: ptr(false),
					},
				},

				HeadTracker: evmcfg.HeadTracker{
					HistoryDepth:            ptr[uint32](15),
					MaxBufferSize:           ptr[uint32](17),
					SamplingInterval:        &hour,
					FinalityTagBypass:       ptr[bool](false),
					MaxAllowedFinalityDepth: ptr[uint32](1500),
				},

				NodePool: evmcfg.NodePool{
					PollFailureThreshold:       ptr[uint32](5),
					PollInterval:               &minute,
					SelectionMode:              &selectionMode,
					SyncThreshold:              ptr[uint32](13),
					LeaseDuration:              &zeroSeconds,
					NodeIsSyncingEnabled:       ptr(true),
					FinalizedBlockPollInterval: &second,
					EnforceRepeatableRead:      ptr(true),
					DeathDeclarationDelay:      &minute,
					NewHeadsPollInterval:       &zeroSeconds,
					Errors: evmcfg.ClientErrors{
						NonceTooLow:                       ptr[string]("(: |^)nonce too low"),
						NonceTooHigh:                      ptr[string]("(: |^)nonce too high"),
						ReplacementTransactionUnderpriced: ptr[string]("(: |^)replacement transaction underpriced"),
						LimitReached:                      ptr[string]("(: |^)limit reached"),
						TransactionAlreadyInMempool:       ptr[string]("(: |^)transaction already in mempool"),
						TerminallyUnderpriced:             ptr[string]("(: |^)terminally underpriced"),
						InsufficientEth:                   ptr[string]("(: |^)insufficient eth"),
						TxFeeExceedsCap:                   ptr[string]("(: |^)tx fee exceeds cap"),
						L2FeeTooLow:                       ptr[string]("(: |^)l2 fee too low"),
						L2FeeTooHigh:                      ptr[string]("(: |^)l2 fee too high"),
						L2Full:                            ptr[string]("(: |^)l2 full"),
						TransactionAlreadyMined:           ptr[string]("(: |^)transaction already mined"),
						Fatal:                             ptr[string]("(: |^)fatal"),
						ServiceUnavailable:                ptr[string]("(: |^)service unavailable"),
						TooManyResults:                    ptr[string]("(: |^)too many results"),
					},
				},
				OCR: evmcfg.OCR{
					ContractConfirmations:              ptr[uint16](11),
					ContractTransmitterTransmitTimeout: &minute,
					DatabaseTimeout:                    &second,
					DeltaCOverride:                     commoncfg.MustNewDuration(time.Hour),
					DeltaCJitterOverride:               commoncfg.MustNewDuration(time.Second),
					ObservationGracePeriod:             &second,
				},
				OCR2: evmcfg.OCR2{
					Automation: evmcfg.Automation{
						GasLimit: ptr[uint32](540),
					},
				},
				Workflow: evmcfg.Workflow{
					GasLimitDefault: ptr[uint64](400000),
				},
			},
			Nodes: []*evmcfg.Node{
				{
					Name:    ptr("foo"),
					HTTPURL: mustURL("https://foo.web"),
					WSURL:   mustURL("wss://web.socket/test/foo"),
				},
				{
					Name:    ptr("bar"),
					HTTPURL: mustURL("https://bar.com"),
					WSURL:   mustURL("wss://web.socket/test/bar"),
				},
				{
					Name:     ptr("broadcast"),
					HTTPURL:  mustURL("http://broadcast.mirror"),
					SendOnly: ptr(true),
				},
			}},
	}
	full.Solana = []*solcfg.TOMLConfig{
		{
			ChainID: ptr("mainnet"),
			Enabled: ptr(false),
			Chain: solcfg.Chain{
				BalancePollPeriod:       commoncfg.MustNewDuration(time.Minute),
				ConfirmPollPeriod:       commoncfg.MustNewDuration(time.Second),
				OCR2CachePollPeriod:     commoncfg.MustNewDuration(time.Minute),
				OCR2CacheTTL:            commoncfg.MustNewDuration(time.Hour),
				TxTimeout:               commoncfg.MustNewDuration(time.Hour),
				TxRetryTimeout:          commoncfg.MustNewDuration(time.Minute),
				TxConfirmTimeout:        commoncfg.MustNewDuration(time.Second),
				SkipPreflight:           ptr(true),
				Commitment:              ptr("banana"),
				MaxRetries:              ptr[int64](7),
				FeeEstimatorMode:        ptr("fixed"),
				ComputeUnitPriceMax:     ptr[uint64](1000),
				ComputeUnitPriceMin:     ptr[uint64](10),
				ComputeUnitPriceDefault: ptr[uint64](100),
				FeeBumpPeriod:           commoncfg.MustNewDuration(time.Minute),
				BlockHistoryPollPeriod:  commoncfg.MustNewDuration(time.Minute),
			},
			Nodes: []*solcfg.Node{
				{Name: ptr("primary"), URL: commoncfg.MustParseURL("http://solana.web")},
				{Name: ptr("foo"), URL: commoncfg.MustParseURL("http://solana.foo")},
				{Name: ptr("bar"), URL: commoncfg.MustParseURL("http://solana.bar")},
			},
		},
	}
	full.Starknet = []*stkcfg.TOMLConfig{
		{
			ChainID: ptr("foobar"),
			Enabled: ptr(true),
			Chain: stkcfg.Chain{
				OCR2CachePollPeriod: commoncfg.MustNewDuration(6 * time.Hour),
				OCR2CacheTTL:        commoncfg.MustNewDuration(3 * time.Minute),
				RequestTimeout:      commoncfg.MustNewDuration(time.Minute + 3*time.Second),
				TxTimeout:           commoncfg.MustNewDuration(13 * time.Second),
				ConfirmationPoll:    commoncfg.MustNewDuration(42 * time.Second),
			},
			FeederURL: commoncfg.MustParseURL("http://feeder.url"),
			Nodes: []*stkcfg.Node{
				{Name: ptr("primary"), URL: commoncfg.MustParseURL("http://stark.node"), APIKey: ptr("key")},
			},
		},
	}
	full.Cosmos = []*coscfg.TOMLConfig{
		{
			ChainID: ptr("Malaga-420"),
			Enabled: ptr(true),
			Chain: coscfg.Chain{
				Bech32Prefix:         ptr("wasm"),
				BlockRate:            commoncfg.MustNewDuration(time.Minute),
				BlocksUntilTxTimeout: ptr[int64](12),
				ConfirmPollPeriod:    commoncfg.MustNewDuration(time.Second),
				FallbackGasPrice:     mustDecimal("0.001"),
				GasToken:             ptr("ucosm"),
				GasLimitMultiplier:   mustDecimal("1.2"),
				MaxMsgsPerBatch:      ptr[int64](17),
				OCR2CachePollPeriod:  commoncfg.MustNewDuration(time.Minute),
				OCR2CacheTTL:         commoncfg.MustNewDuration(time.Hour),
				TxMsgTimeout:         commoncfg.MustNewDuration(time.Second),
			},
			Nodes: []*coscfg.Node{
				{Name: ptr("primary"), TendermintURL: commoncfg.MustParseURL("http://tender.mint")},
				{Name: ptr("foo"), TendermintURL: commoncfg.MustParseURL("http://foo.url")},
				{Name: ptr("bar"), TendermintURL: commoncfg.MustParseURL("http://bar.web")},
			},
		},
	}
	full.Mercury = toml.Mercury{
		Cache: toml.MercuryCache{
			LatestReportTTL:      commoncfg.MustNewDuration(100 * time.Second),
			MaxStaleAge:          commoncfg.MustNewDuration(101 * time.Second),
			LatestReportDeadline: commoncfg.MustNewDuration(102 * time.Second),
		},
		TLS: toml.MercuryTLS{
			CertFile: ptr("/path/to/cert.pem"),
		},
		Transmitter: toml.MercuryTransmitter{
			TransmitQueueMaxSize: ptr(uint32(123)),
			TransmitTimeout:      commoncfg.MustNewDuration(234 * time.Second),
		},
		VerboseLogging: ptr(true),
	}

	for _, tt := range []struct {
		name   string
		config Config
		exp    string
	}{
		{"empty", Config{}, ``},
		{"global", global, `InsecureFastScrypt = true
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
		{"FluxMonitor", Config{Core: toml.Core{FluxMonitor: full.FluxMonitor}}, `[FluxMonitor]
DefaultTransactionQueueDepth = 100
SimulateTransactions = true
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
DefaultTransactionQueueDepth = 1
SimulateTransactions = false
TraceLogging = false
`},
		{"P2P", Config{Core: toml.Core{P2P: full.P2P}}, `[P2P]
IncomingMessageBufferSize = 13
OutgoingMessageBufferSize = 17
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw'
TraceLogging = true

[P2P.V2]
Enabled = false
AnnounceAddresses = ['a', 'b', 'c']
DefaultBootstrappers = ['12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw@foo:42/bar:10', '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw@test:99']
DeltaDial = '1m0s'
DeltaReconcile = '1s'
ListenAddresses = ['foo', 'bar']
`},
		{"Keeper", Config{Core: toml.Core{Keeper: full.Keeper}}, `[Keeper]
DefaultTransactionQueueDepth = 17
GasPriceBufferPercent = 12
GasTipCapBufferPercent = 43
BaseFeeBufferPercent = 89
MaxGracePeriod = 31
TurnLookBack = 91

[Keeper.Registry]
CheckGasOverhead = 90
PerformGasOverhead = 4294967295
MaxPerformDataSize = 5000
SyncInterval = '1h0m0s'
SyncUpkeepQueueSize = 31
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
FinalityTagEnabled = false
FlagsContractAddress = '0xae4E781a6218A8031764928E88d457937A954fC3'
LinkContractAddress = '0x538aAaB4ea120b2bC2fe5D296852D948F07D849e'
LogBackfillBatchSize = 17
LogPollInterval = '1m0s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 10000
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
ForwardersEnabled = true
MaxInFlight = 19
MaxQueued = 99
ReaperInterval = '1m0s'
ReaperThreshold = '1m0s'
ResendAfterThreshold = '1h0m0s'

[EVM.Transactions.AutoPurge]
Enabled = false

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

[[EVM.KeySpecific]]
Key = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292'

[EVM.KeySpecific.GasEstimator]
PriceMax = '79.228162514264337593543950335 gether'

[EVM.NodePool]
PollFailureThreshold = 5
PollInterval = '1m0s'
SelectionMode = 'HighestHead'
SyncThreshold = 13
LeaseDuration = '0s'
NodeIsSyncingEnabled = true
FinalizedBlockPollInterval = '1s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

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

[[EVM.Nodes]]
Name = 'foo'
WSURL = 'wss://web.socket/test/foo'
HTTPURL = 'https://foo.web'

[[EVM.Nodes]]
Name = 'bar'
WSURL = 'wss://web.socket/test/bar'
HTTPURL = 'https://bar.com'

[[EVM.Nodes]]
Name = 'broadcast'
HTTPURL = 'http://broadcast.mirror'
SendOnly = true
`},
		{"Cosmos", Config{Cosmos: full.Cosmos}, `[[Cosmos]]
ChainID = 'Malaga-420'
Enabled = true
Bech32Prefix = 'wasm'
BlockRate = '1m0s'
BlocksUntilTxTimeout = 12
ConfirmPollPeriod = '1s'
FallbackGasPrice = '0.001'
GasToken = 'ucosm'
GasLimitMultiplier = '1.2'
MaxMsgsPerBatch = 17
OCR2CachePollPeriod = '1m0s'
OCR2CacheTTL = '1h0m0s'
TxMsgTimeout = '1s'

[[Cosmos.Nodes]]
Name = 'primary'
TendermintURL = 'http://tender.mint'

[[Cosmos.Nodes]]
Name = 'foo'
TendermintURL = 'http://foo.url'

[[Cosmos.Nodes]]
Name = 'bar'
TendermintURL = 'http://bar.web'
`},
		{"Solana", Config{Solana: full.Solana}, `[[Solana]]
ChainID = 'mainnet'
Enabled = false
BalancePollPeriod = '1m0s'
ConfirmPollPeriod = '1s'
OCR2CachePollPeriod = '1m0s'
OCR2CacheTTL = '1h0m0s'
TxTimeout = '1h0m0s'
TxRetryTimeout = '1m0s'
TxConfirmTimeout = '1s'
SkipPreflight = true
Commitment = 'banana'
MaxRetries = 7
FeeEstimatorMode = 'fixed'
ComputeUnitPriceMax = 1000
ComputeUnitPriceMin = 10
ComputeUnitPriceDefault = 100
FeeBumpPeriod = '1m0s'
BlockHistoryPollPeriod = '1m0s'

[[Solana.Nodes]]
Name = 'primary'
URL = 'http://solana.web'

[[Solana.Nodes]]
Name = 'foo'
URL = 'http://solana.foo'

[[Solana.Nodes]]
Name = 'bar'
URL = 'http://solana.bar'
`},
		{"Starknet", Config{Starknet: full.Starknet}, `[[Starknet]]
ChainID = 'foobar'
FeederURL = 'http://feeder.url'
Enabled = true
OCR2CachePollPeriod = '6h0m0s'
OCR2CacheTTL = '3m0s'
RequestTimeout = '1m3s'
TxTimeout = '13s'
ConfirmationPoll = '42s'

[[Starknet.Nodes]]
Name = 'primary'
URL = 'http://stark.node'
APIKey = 'key'
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
TransmitQueueMaxSize = 123
TransmitTimeout = '3m54s'
`},
		{"full", full, fullTOML},
		{"multi-chain", multiChain, multiChainTOML},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.config.TOMLString()
			require.NoError(t, err)
			assert.Equal(t, tt.exp, s, diff.Diff(tt.exp, s))

			var got Config

			require.NoError(t, config.DecodeTOML(strings.NewReader(s), &got))
			ts, err := got.TOMLString()

			require.NoError(t, err)
			assert.Equal(t, tt.config, got, diff.Diff(s, ts))
		})
	}
}

func TestConfig_full(t *testing.T) {
	var got Config
	require.NoError(t, config.DecodeTOML(strings.NewReader(fullTOML), &got))
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
			got.EVM[c].Workflow.GasLimitDefault = ptr(uint64(400000))
		}
		for n := range got.EVM[c].Nodes {
			if got.EVM[c].Nodes[n].WSURL == nil {
				got.EVM[c].Nodes[n].WSURL = new(commoncfg.URL)
			}
			if got.EVM[c].Nodes[n].SendOnly == nil {
				got.EVM[c].Nodes[n].SendOnly = ptr(true)
			}
			if got.EVM[c].Nodes[n].Order == nil {
				got.EVM[c].Nodes[n].Order = ptr(int32(100))
			}
		}
		if got.EVM[c].Transactions.AutoPurge.Threshold == nil {
			got.EVM[c].Transactions.AutoPurge.Threshold = ptr(uint32(0))
		}
		if got.EVM[c].Transactions.AutoPurge.MinAttempts == nil {
			got.EVM[c].Transactions.AutoPurge.MinAttempts = ptr(uint32(0))
		}
		if got.EVM[c].Transactions.AutoPurge.DetectionApiUrl == nil {
			got.EVM[c].Transactions.AutoPurge.DetectionApiUrl = new(commoncfg.URL)
		}
	}

	cfgtest.AssertFieldsNotNil(t, got)
}

//go:embed testdata/config-invalid.toml
var invalidTOML string

func TestConfig_Validate(t *testing.T) {
	for _, tt := range []struct {
		name string
		toml string
		exp  string
	}{
		{name: "invalid", toml: invalidTOML, exp: `invalid configuration: 8 errors:
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
			- ChainType: invalid value (Foo): must be one of arbitrum, astar, celo, gnosis, hedera, kroma, mantle, metis, optimismBedrock, sei, scroll, wemix, xlayer, zkevm, zksync, zircuit or omitted
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
			- ChainType: invalid value (Arbitrum): must be one of arbitrum, astar, celo, gnosis, hedera, kroma, mantle, metis, optimismBedrock, sei, scroll, wemix, xlayer, zkevm, zksync, zircuit or omitted
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
	- Cosmos: 5 errors:
		- 1.ChainID: invalid value (Malaga-420): duplicate - must be unique
		- 0.Nodes.1.Name: invalid value (test): duplicate - must be unique
		- 0.Nodes: 2 errors:
				- 0.TendermintURL: missing: required for all nodes
				- 1.TendermintURL: missing: required for all nodes
		- 1.Nodes: missing: must have at least one node
		- 2: 2 errors:
			- ChainID: missing: required for all chains
			- Nodes: missing: must have at least one node
	- Solana: 5 errors:
		- 1.ChainID: invalid value (mainnet): duplicate - must be unique
		- 1.Nodes.1.Name: invalid value (bar): duplicate - must be unique
		- 0.Nodes: missing: must have at least one node
		- 1.Nodes: 2 errors:
				- 0.URL: missing: required for all nodes
				- 1.URL: missing: required for all nodes
		- 2: 2 errors:
			- ChainID: missing: required for all chains
			- Nodes: missing: must have at least one node
	- Starknet: 3 errors:
		- 0.Nodes.1.Name: invalid value (primary): duplicate - must be unique
		- 0.ChainID: missing: required for all chains
		- 1: 2 errors:
			- ChainID: missing: required for all chains
			- Nodes: missing: must have at least one node
	- Aptos.0.Enabled: invalid value (1): expected bool`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var c Config
			require.NoError(t, config.DecodeTOML(strings.NewReader(tt.toml), &c))
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
	const PWD_OVERRIDE = "great_password"
	const DBURL_OVERRIDE = "http://user@db"

	t.Setenv("CL_DATABASE_URL", DBURL_OVERRIDE)

	// Check for two overrides
	opts := GeneralConfigOpts{
		ConfigStrings:  []string{fullTOML},
		SecretsStrings: []string{secretsFullTOML},
	}
	c, err := opts.New()
	assert.NoError(t, err)
	c.SetPasswords(ptr(PWD_OVERRIDE), nil)
	assert.Equal(t, PWD_OVERRIDE, c.Password().Keystore())
	dbURL := c.Database().URL()
	assert.Equal(t, DBURL_OVERRIDE, (&dbURL).String())
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
			require.NoError(t, config.DecodeTOML(strings.NewReader(tt.toml), &s))
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
	c.EVM = evmcfg.EVMConfigs{{ChainID: ubig.NewI(99999133712345)}}
	c.Cosmos = coscfg.TOMLConfigs{{ChainID: ptr("unknown cosmos chain")}}
	c.Solana = solcfg.TOMLConfigs{{ChainID: ptr("unknown solana chain")}}
	c.Starknet = stkcfg.TOMLConfigs{{ChainID: ptr("unknown starknet chain")}}
	c.setDefaults()
	if s, err := c.TOMLString(); assert.NoError(t, err) {
		t.Log(s, err)
	}
	cfgtest.AssertFieldsNotNil(t, c.Core)
}

func Test_validateEnv(t *testing.T) {
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("DATABASE_URL", "foo")
	assert.ErrorContains(t, validateEnv(), `invalid environment: 2 errors:
	- environment variable DATABASE_URL must not be set: unsupported with config v2
	- environment variable LOG_LEVEL must not be set: unsupported with config v2`)

	t.Setenv("GAS_UPDATER_ENABLED", "true")
	t.Setenv("ETH_GAS_BUMP_TX_DEPTH", "7")
	assert.ErrorContains(t, validateEnv(), `invalid environment: 4 errors:
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
				require.NoError(t, config.DecodeTOML(strings.NewReader(fs), &f))
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
						Enabled:     ptr(true),
						Mode:        ptr("unencrypted"),
						TLSCertPath: ptr("/path/to/cert.pem"),
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

func ptr[T any](t T) *T { return &t }

func mustHexToBig(t *testing.T, hx string) *big.Int {
	n, err := hex.ParseBig(hx)
	require.NoError(t, err)
	return n
}
