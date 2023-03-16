package chainlink

import (
	_ "embed"
	"math"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/kylelemons/godebug/diff"
	"github.com/shopspring/decimal"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	relayutils "github.com/smartcontractkit/chainlink-relay/pkg/utils"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/chains/starknet"
	legacy "github.com/smartcontractkit/chainlink/core/config"
	config "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	//go:embed testdata/config-full.toml
	fullTOML string
	//go:embed testdata/config-multi-chain.toml
	multiChainTOML string

	multiChain = Config{
		Core: config.Core{
			RootDir: ptr("my/root/dir"),

			AuditLogger: audit.AuditLoggerConfig{
				Enabled:      ptr(true),
				ForwardToUrl: mustURL("http://localhost:9898"),
				Headers: ptr([]audit.ServiceHeader{
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

			Database: config.Database{
				Listener: config.DatabaseListener{
					FallbackPollInterval: models.MustNewDuration(2 * time.Minute),
				},
			},
			Log: config.Log{
				Level:       ptr(config.LogLevel(zapcore.PanicLevel)),
				JSONConsole: ptr(true),
			},
			JobPipeline: config.JobPipeline{
				HTTPRequest: config.JobPipelineHTTPRequest{
					DefaultTimeout: models.MustNewDuration(30 * time.Second),
				},
			},
			OCR2: config.OCR2{
				Enabled:         ptr(true),
				DatabaseTimeout: models.MustNewDuration(20 * time.Second),
			},
			OCR: config.OCR{
				Enabled:           ptr(true),
				BlockchainTimeout: models.MustNewDuration(5 * time.Second),
			},
			P2P: config.P2P{
				IncomingMessageBufferSize: ptr[int64](999),
			},
			Keeper: config.Keeper{
				GasPriceBufferPercent: ptr[uint16](10),
			},
			AutoPprof: config.AutoPprof{
				CPUProfileRate: ptr[int64](7),
			},
		},
		EVM: []*evmcfg.EVMConfig{
			{
				ChainID: utils.NewBigI(1),
				Chain: evmcfg.Chain{
					FinalityDepth: ptr[uint32](26),
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
				ChainID: utils.NewBigI(42),
				Chain: evmcfg.Chain{
					GasEstimator: evmcfg.GasEstimator{
						PriceDefault: assets.NewWeiI(math.MaxInt64),
					},
				},
				Nodes: []*evmcfg.Node{
					{
						Name:  ptr("primary"),
						WSURL: mustURL("wss://web.socket/test"),
					},
				}},
			{
				ChainID: utils.NewBigI(137),
				Chain: evmcfg.Chain{
					GasEstimator: evmcfg.GasEstimator{
						Mode: ptr("FixedPrice"),
					},
				},
				Nodes: []*evmcfg.Node{
					{
						Name:  ptr("primary"),
						WSURL: mustURL("wss://web.socket/test"),
					},
				}},
		},
		Solana: []*solana.SolanaConfig{
			{
				ChainID: ptr("mainnet"),
				Chain: solcfg.Chain{
					MaxRetries: ptr[int64](12),
				},
				Nodes: []*solcfg.Node{
					{Name: ptr("primary"), URL: relayutils.MustParseURL("http://mainnet.solana.com")},
				},
			},
			{
				ChainID: ptr("testnet"),
				Chain: solcfg.Chain{
					OCR2CachePollPeriod: relayutils.MustNewDuration(time.Minute),
				},
				Nodes: []*solcfg.Node{
					{Name: ptr("primary"), URL: relayutils.MustParseURL("http://testnet.solana.com")},
				},
			},
		},
		Starknet: []*starknet.StarknetConfig{
			{
				ChainID: ptr("foobar"),
				Chain: stkcfg.Chain{
					TxSendFrequency: relayutils.MustNewDuration(time.Hour),
				},
				Nodes: []*stkcfg.Node{
					{Name: ptr("primary"), URL: relayutils.MustParseURL("http://stark.node")},
				},
			},
		},
	}
)

func TestConfig_Marshal(t *testing.T) {
	second := models.MustMakeDuration(time.Second)
	minute := models.MustMakeDuration(time.Minute)
	hour := models.MustMakeDuration(time.Hour)
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
	mustAddress := func(s string) *ethkey.EIP55Address {
		a, err := ethkey.NewEIP55Address(s)
		require.NoError(t, err)
		return &a
	}
	selectionMode := client.NodeSelectionMode_HighestHead

	global := Config{
		Core: config.Core{
			ExplorerURL:         mustURL("http://explorer.url"),
			InsecureFastScrypt:  ptr(true),
			RootDir:             ptr("test/root/dir"),
			ShutdownGracePeriod: models.MustNewDuration(10 * time.Second),
		},
	}

	full := global

	serviceHeaders := []audit.ServiceHeader{
		{Header: "Authorization", Value: "token"},
		{Header: "X-SomeOther-Header", Value: "value with spaces | and a bar+*"},
	}
	full.AuditLogger = audit.AuditLoggerConfig{
		Enabled:        ptr(true),
		ForwardToUrl:   mustURL("http://localhost:9898"),
		Headers:        ptr(serviceHeaders),
		JsonWrapperKey: ptr("event"),
	}

	full.Feature = config.Feature{
		FeedsManager: ptr(true),
		LogPoller:    ptr(true),
		UICSAKeys:    ptr(true),
	}
	full.Database = config.Database{
		DefaultIdleInTxSessionTimeout: models.MustNewDuration(time.Minute),
		DefaultLockTimeout:            models.MustNewDuration(time.Hour),
		DefaultQueryTimeout:           models.MustNewDuration(time.Second),
		LogQueries:                    ptr(true),
		MigrateOnStartup:              ptr(true),
		MaxIdleConns:                  ptr[int64](7),
		MaxOpenConns:                  ptr[int64](13),
		Listener: config.DatabaseListener{
			MaxReconnectDuration: models.MustNewDuration(time.Minute),
			MinReconnectInterval: models.MustNewDuration(5 * time.Minute),
			FallbackPollInterval: models.MustNewDuration(2 * time.Minute),
		},
		Lock: config.DatabaseLock{
			Enabled:              ptr(false),
			LeaseDuration:        &minute,
			LeaseRefreshInterval: &second,
		},
		Backup: config.DatabaseBackup{
			Dir:              ptr("test/backup/dir"),
			Frequency:        &hour,
			Mode:             &legacy.DatabaseBackupModeFull,
			OnVersionUpgrade: ptr(true),
		},
	}
	full.TelemetryIngress = config.TelemetryIngress{
		UniConn:      ptr(true),
		Logging:      ptr(true),
		ServerPubKey: ptr("test-pub-key"),
		URL:          mustURL("https://prom.test"),
		BufferSize:   ptr[uint16](1234),
		MaxBatchSize: ptr[uint16](4321),
		SendInterval: models.MustNewDuration(time.Minute),
		SendTimeout:  models.MustNewDuration(5 * time.Second),
		UseBatchSend: ptr(true),
	}
	full.Log = config.Log{
		Level:       ptr(config.LogLevel(zapcore.DPanicLevel)),
		JSONConsole: ptr(true),
		UnixTS:      ptr(true),
		File: config.LogFile{
			Dir:        ptr("log/file/dir"),
			MaxSize:    ptr[utils.FileSize](100 * utils.GB),
			MaxAgeDays: ptr[int64](17),
			MaxBackups: ptr[int64](9),
		},
	}
	full.WebServer = config.WebServer{
		AllowOrigins:            ptr("*"),
		BridgeResponseURL:       mustURL("https://bridge.response"),
		BridgeCacheTTL:          models.MustNewDuration(10 * time.Second),
		HTTPWriteTimeout:        models.MustNewDuration(time.Minute),
		HTTPPort:                ptr[uint16](56),
		SecureCookies:           ptr(true),
		SessionTimeout:          models.MustNewDuration(time.Hour),
		SessionReaperExpiration: models.MustNewDuration(7 * 24 * time.Hour),
		MFA: config.WebServerMFA{
			RPID:     ptr("test-rpid"),
			RPOrigin: ptr("test-rp-origin"),
		},
		RateLimit: config.WebServerRateLimit{
			Authenticated:         ptr[int64](42),
			AuthenticatedPeriod:   models.MustNewDuration(time.Second),
			Unauthenticated:       ptr[int64](7),
			UnauthenticatedPeriod: models.MustNewDuration(time.Minute),
		},
		TLS: config.WebServerTLS{
			CertPath:      ptr("tls/cert/path"),
			Host:          ptr("tls-host"),
			KeyPath:       ptr("tls/key/path"),
			HTTPSPort:     ptr[uint16](6789),
			ForceRedirect: ptr(true),
		},
	}
	full.JobPipeline = config.JobPipeline{
		ExternalInitiatorsEnabled: ptr(true),
		MaxRunDuration:            models.MustNewDuration(time.Hour),
		MaxSuccessfulRuns:         ptr[uint64](123456),
		ReaperInterval:            models.MustNewDuration(4 * time.Hour),
		ReaperThreshold:           models.MustNewDuration(7 * 24 * time.Hour),
		ResultWriteQueueDepth:     ptr[uint32](10),
		HTTPRequest: config.JobPipelineHTTPRequest{
			MaxSize:        ptr[utils.FileSize](100 * utils.MB),
			DefaultTimeout: models.MustNewDuration(time.Minute),
		},
	}
	full.FluxMonitor = config.FluxMonitor{
		DefaultTransactionQueueDepth: ptr[uint32](100),
		SimulateTransactions:         ptr(true),
	}
	full.OCR2 = config.OCR2{
		Enabled:                            ptr(true),
		ContractConfirmations:              ptr[uint32](11),
		BlockchainTimeout:                  models.MustNewDuration(3 * time.Second),
		ContractPollInterval:               models.MustNewDuration(time.Hour),
		ContractSubscribeInterval:          models.MustNewDuration(time.Minute),
		ContractTransmitterTransmitTimeout: models.MustNewDuration(time.Minute),
		DatabaseTimeout:                    models.MustNewDuration(8 * time.Second),
		KeyBundleID:                        ptr(models.MustSha256HashFromHex("7a5f66bbe6594259325bf2b4f5b1a9c9")),
	}
	full.OCR = config.OCR{
		Enabled:                      ptr(true),
		ObservationTimeout:           models.MustNewDuration(11 * time.Second),
		BlockchainTimeout:            models.MustNewDuration(3 * time.Second),
		ContractPollInterval:         models.MustNewDuration(time.Hour),
		ContractSubscribeInterval:    models.MustNewDuration(time.Minute),
		DefaultTransactionQueueDepth: ptr[uint32](12),
		KeyBundleID:                  ptr(models.MustSha256HashFromHex("acdd42797a8b921b2910497badc50006")),
		SimulateTransactions:         ptr(true),
		TransmitterAddress:           ptr(ethkey.MustEIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")),
	}
	full.P2P = config.P2P{
		IncomingMessageBufferSize: ptr[int64](13),
		OutgoingMessageBufferSize: ptr[int64](17),
		PeerID:                    mustPeerID("12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw"),
		TraceLogging:              ptr(true),
		V1: config.P2PV1{
			Enabled:                          ptr(false),
			AnnounceIP:                       mustIP("1.2.3.4"),
			AnnouncePort:                     ptr[uint16](1234),
			BootstrapCheckInterval:           models.MustNewDuration(time.Minute),
			DefaultBootstrapPeers:            &[]string{"foo", "bar", "should", "these", "be", "typed"},
			DHTAnnouncementCounterUserPrefix: ptr[uint32](4321),
			DHTLookupInterval:                ptr[int64](9),
			ListenIP:                         mustIP("4.3.2.1"),
			ListenPort:                       ptr[uint16](9),
			NewStreamTimeout:                 models.MustNewDuration(time.Second),
			PeerstoreWriteInterval:           models.MustNewDuration(time.Minute),
		},
		V2: config.P2PV2{
			Enabled:           ptr(true),
			AnnounceAddresses: &[]string{"a", "b", "c"},
			DefaultBootstrappers: &[]ocrcommontypes.BootstrapperLocator{
				{PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", Addrs: []string{"foo:42", "bar:10"}},
				{PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", Addrs: []string{"test:99"}},
			},
			DeltaDial:       models.MustNewDuration(time.Minute),
			DeltaReconcile:  models.MustNewDuration(time.Second),
			ListenAddresses: &[]string{"foo", "bar"},
		},
	}
	full.Keeper = config.Keeper{
		DefaultTransactionQueueDepth: ptr[uint32](17),
		GasPriceBufferPercent:        ptr[uint16](12),
		GasTipCapBufferPercent:       ptr[uint16](43),
		BaseFeeBufferPercent:         ptr[uint16](89),
		MaxGracePeriod:               ptr[int64](31),
		TurnLookBack:                 ptr[int64](91),
		Registry: config.KeeperRegistry{
			CheckGasOverhead:    ptr[uint32](90),
			PerformGasOverhead:  ptr[uint32](math.MaxUint32),
			SyncInterval:        models.MustNewDuration(time.Hour),
			SyncUpkeepQueueSize: ptr[uint32](31),
			MaxPerformDataSize:  ptr[uint32](5000),
		},
	}
	full.AutoPprof = config.AutoPprof{
		Enabled:              ptr(true),
		ProfileRoot:          ptr("prof/root"),
		PollInterval:         models.MustNewDuration(time.Minute),
		GatherDuration:       models.MustNewDuration(12 * time.Second),
		GatherTraceDuration:  models.MustNewDuration(13 * time.Second),
		MaxProfileSize:       ptr[utils.FileSize](utils.GB),
		CPUProfileRate:       ptr[int64](7),
		MemProfileRate:       ptr[int64](9),
		BlockProfileRate:     ptr[int64](5),
		MutexProfileFraction: ptr[int64](2),
		MemThreshold:         ptr[utils.FileSize](utils.GB),
		GoroutineThreshold:   ptr[int64](999),
	}
	full.Pyroscope = config.Pyroscope{
		ServerAddress: ptr("http://localhost:4040"),
		Environment:   ptr("tests"),
	}
	full.Sentry = config.Sentry{
		Debug:       ptr(true),
		DSN:         ptr("sentry-dsn"),
		Environment: ptr("dev"),
		Release:     ptr("v1.2.3"),
	}
	full.EVM = []*evmcfg.EVMConfig{
		{
			ChainID: utils.NewBigI(1),
			Enabled: ptr(false),
			Chain: evmcfg.Chain{
				BalanceMonitor: evmcfg.BalanceMonitor{
					Enabled: ptr(true),
				},
				BlockBackfillDepth:   ptr[uint32](100),
				BlockBackfillSkip:    ptr(true),
				ChainType:            ptr("Optimism"),
				FinalityDepth:        ptr[uint32](42),
				FlagsContractAddress: mustAddress("0xae4E781a6218A8031764928E88d457937A954fC3"),

				GasEstimator: evmcfg.GasEstimator{
					Mode:               ptr("L2Suggested"),
					EIP1559DynamicFees: ptr(true),
					BumpPercent:        ptr[uint16](10),
					BumpThreshold:      ptr[uint32](6),
					BumpTxDepth:        ptr[uint16](6),
					BumpMin:            assets.NewWeiI(100),
					FeeCapDefault:      assets.NewWeiI(math.MaxInt64),
					LimitDefault:       ptr[uint32](12),
					LimitMax:           ptr[uint32](17),
					LimitMultiplier:    mustDecimal("1.234"),
					LimitTransfer:      ptr[uint32](100),
					TipCapDefault:      assets.NewWeiI(2),
					TipCapMin:          assets.NewWeiI(1),
					PriceDefault:       assets.NewWeiI(math.MaxInt64),
					PriceMax:           assets.NewWei(utils.HexToBig("FFFFFFFFFFFF")),
					PriceMin:           assets.NewWeiI(13),

					LimitJobType: evmcfg.GasLimitJobType{
						OCR:    ptr[uint32](1001),
						DR:     ptr[uint32](1002),
						VRF:    ptr[uint32](1003),
						FM:     ptr[uint32](1004),
						Keeper: ptr[uint32](1005),
					},

					BlockHistory: evmcfg.BlockHistoryEstimator{
						BatchSize:                 ptr[uint32](17),
						BlockHistorySize:          ptr[uint16](12),
						CheckInclusionBlocks:      ptr[uint16](18),
						CheckInclusionPercentile:  ptr[uint16](19),
						EIP1559FeeCapBufferBlocks: ptr[uint16](13),
						TransactionPercentile:     ptr[uint16](15),
					},
				},

				KeySpecific: []evmcfg.KeySpecific{
					{
						Key: mustAddress("0x2a3e23c6f242F5345320814aC8a1b4E58707D292"),
						GasEstimator: evmcfg.KeySpecificGasEstimator{
							PriceMax: assets.NewWei(utils.HexToBig("FFFFFFFFFFFFFFFFFFFFFFFF")),
						},
					},
				},

				LinkContractAddress:      mustAddress("0x538aAaB4ea120b2bC2fe5D296852D948F07D849e"),
				LogBackfillBatchSize:     ptr[uint32](17),
				LogPollInterval:          &minute,
				LogKeepBlocksDepth:       ptr[uint32](100000),
				MinContractPayment:       assets.NewLinkFromJuels(math.MaxInt64),
				MinIncomingConfirmations: ptr[uint32](13),
				NonceAutoSync:            ptr(true),
				NoNewHeadsThreshold:      &minute,
				OperatorFactoryAddress:   mustAddress("0xa5B85635Be42F21f94F28034B7DA440EeFF0F418"),
				RPCDefaultBatchSize:      ptr[uint32](17),
				RPCBlockQueryDelay:       ptr[uint16](10),

				Transactions: evmcfg.Transactions{
					MaxInFlight:          ptr[uint32](19),
					MaxQueued:            ptr[uint32](99),
					ReaperInterval:       &minute,
					ReaperThreshold:      &minute,
					ResendAfterThreshold: &hour,
					ForwardersEnabled:    ptr(true),
				},

				HeadTracker: evmcfg.HeadTracker{
					HistoryDepth:     ptr[uint32](15),
					MaxBufferSize:    ptr[uint32](17),
					SamplingInterval: &hour,
				},

				NodePool: evmcfg.NodePool{
					PollFailureThreshold: ptr[uint32](5),
					PollInterval:         &minute,
					SelectionMode:        &selectionMode,
					SyncThreshold:        ptr[uint32](13),
				},
				OCR: evmcfg.OCR{
					ContractConfirmations:              ptr[uint16](11),
					ContractTransmitterTransmitTimeout: &minute,
					DatabaseTimeout:                    &second,
					ObservationGracePeriod:             &second,
				},
				OCR2: evmcfg.OCR2{
					Automation: evmcfg.Automation{
						GasLimit: ptr[uint32](540),
					},
				},
			},
			Nodes: []*evmcfg.Node{
				{
					Name:    ptr("foo"),
					HTTPURL: mustURL("https://foo.web"),
					WSURL:   mustURL("wss://web.socket/test"),
				},
				{
					Name:    ptr("bar"),
					HTTPURL: mustURL("https://bar.com"),
					WSURL:   mustURL("wss://web.socket/test"),
				},
				{
					Name:     ptr("broadcast"),
					HTTPURL:  mustURL("http://broadcast.mirror"),
					SendOnly: ptr(true),
				},
			}},
	}
	full.Solana = []*solana.SolanaConfig{
		{
			ChainID: ptr("mainnet"),
			Enabled: ptr(false),
			Chain: solcfg.Chain{
				BalancePollPeriod:       relayutils.MustNewDuration(time.Minute),
				ConfirmPollPeriod:       relayutils.MustNewDuration(time.Second),
				OCR2CachePollPeriod:     relayutils.MustNewDuration(time.Minute),
				OCR2CacheTTL:            relayutils.MustNewDuration(time.Hour),
				TxTimeout:               relayutils.MustNewDuration(time.Hour),
				TxRetryTimeout:          relayutils.MustNewDuration(time.Minute),
				TxConfirmTimeout:        relayutils.MustNewDuration(time.Second),
				SkipPreflight:           ptr(true),
				Commitment:              ptr("banana"),
				MaxRetries:              ptr[int64](7),
				FeeEstimatorMode:        ptr("fixed"),
				ComputeUnitPriceMax:     ptr[uint64](1000),
				ComputeUnitPriceMin:     ptr[uint64](10),
				ComputeUnitPriceDefault: ptr[uint64](100),
				FeeBumpPeriod:           relayutils.MustNewDuration(time.Minute),
			},
			Nodes: []*solcfg.Node{
				{Name: ptr("primary"), URL: relayutils.MustParseURL("http://solana.web")},
				{Name: ptr("foo"), URL: relayutils.MustParseURL("http://solana.foo")},
				{Name: ptr("bar"), URL: relayutils.MustParseURL("http://solana.bar")},
			},
		},
	}
	full.Starknet = []*starknet.StarknetConfig{
		{
			ChainID: ptr("foobar"),
			Enabled: ptr(true),
			Chain: stkcfg.Chain{
				OCR2CachePollPeriod: relayutils.MustNewDuration(6 * time.Hour),
				OCR2CacheTTL:        relayutils.MustNewDuration(3 * time.Minute),
				RequestTimeout:      relayutils.MustNewDuration(time.Minute + 3*time.Second),
				TxTimeout:           relayutils.MustNewDuration(13 * time.Second),
				TxSendFrequency:     relayutils.MustNewDuration(42 * time.Second),
				TxMaxBatchSize:      ptr[int64](17),
			},
			Nodes: []*stkcfg.Node{
				{Name: ptr("primary"), URL: relayutils.MustParseURL("http://stark.node")},
			},
		},
	}

	for _, tt := range []struct {
		name   string
		config Config
		exp    string
	}{
		{"empty", Config{}, ``},
		{"global", global, `ExplorerURL = 'http://explorer.url'
InsecureFastScrypt = true
RootDir = 'test/root/dir'
ShutdownGracePeriod = '10s'
`},
		{"AuditLogger", Config{Core: config.Core{AuditLogger: full.AuditLogger}}, `[AuditLogger]
Enabled = true
ForwardToUrl = 'http://localhost:9898'
JsonWrapperKey = 'event'
Headers = ['Authorization: token', 'X-SomeOther-Header: value with spaces | and a bar+*']
`},
		{"Feature", Config{Core: config.Core{Feature: full.Feature}}, `[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true
`},
		{"Database", Config{Core: config.Core{Database: full.Database}}, `[Database]
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
		{"TelemetryIngress", Config{Core: config.Core{TelemetryIngress: full.TelemetryIngress}}, `[TelemetryIngress]
UniConn = true
Logging = true
ServerPubKey = 'test-pub-key'
URL = 'https://prom.test'
BufferSize = 1234
MaxBatchSize = 4321
SendInterval = '1m0s'
SendTimeout = '5s'
UseBatchSend = true
`},
		{"Log", Config{Core: config.Core{Log: full.Log}}, `[Log]
Level = 'crit'
JSONConsole = true
UnixTS = true

[Log.File]
Dir = 'log/file/dir'
MaxSize = '100.00gb'
MaxAgeDays = 17
MaxBackups = 9
`},
		{"WebServer", Config{Core: config.Core{WebServer: full.WebServer}}, `[WebServer]
AllowOrigins = '*'
BridgeResponseURL = 'https://bridge.response'
BridgeCacheTTL = '10s'
HTTPWriteTimeout = '1m0s'
HTTPPort = 56
SecureCookies = true
SessionTimeout = '1h0m0s'
SessionReaperExpiration = '168h0m0s'

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
`},
		{"FluxMonitor", Config{Core: config.Core{FluxMonitor: full.FluxMonitor}}, `[FluxMonitor]
DefaultTransactionQueueDepth = 100
SimulateTransactions = true
`},
		{"JobPipeline", Config{Core: config.Core{JobPipeline: full.JobPipeline}}, `[JobPipeline]
ExternalInitiatorsEnabled = true
MaxRunDuration = '1h0m0s'
MaxSuccessfulRuns = 123456
ReaperInterval = '4h0m0s'
ReaperThreshold = '168h0m0s'
ResultWriteQueueDepth = 10

[JobPipeline.HTTPRequest]
DefaultTimeout = '1m0s'
MaxSize = '100.00mb'
`},
		{"OCR", Config{Core: config.Core{OCR: full.OCR}}, `[OCR]
Enabled = true
ObservationTimeout = '11s'
BlockchainTimeout = '3s'
ContractPollInterval = '1h0m0s'
ContractSubscribeInterval = '1m0s'
DefaultTransactionQueueDepth = 12
KeyBundleID = 'acdd42797a8b921b2910497badc5000600000000000000000000000000000000'
SimulateTransactions = true
TransmitterAddress = '0xa0788FC17B1dEe36f057c42B6F373A34B014687e'
`},
		{"OCR2", Config{Core: config.Core{OCR2: full.OCR2}}, `[OCR2]
Enabled = true
ContractConfirmations = 11
BlockchainTimeout = '3s'
ContractPollInterval = '1h0m0s'
ContractSubscribeInterval = '1m0s'
ContractTransmitterTransmitTimeout = '1m0s'
DatabaseTimeout = '8s'
KeyBundleID = '7a5f66bbe6594259325bf2b4f5b1a9c900000000000000000000000000000000'
`},
		{"P2P", Config{Core: config.Core{P2P: full.P2P}}, `[P2P]
IncomingMessageBufferSize = 13
OutgoingMessageBufferSize = 17
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw'
TraceLogging = true

[P2P.V1]
Enabled = false
AnnounceIP = '1.2.3.4'
AnnouncePort = 1234
BootstrapCheckInterval = '1m0s'
DefaultBootstrapPeers = ['foo', 'bar', 'should', 'these', 'be', 'typed']
DHTAnnouncementCounterUserPrefix = 4321
DHTLookupInterval = 9
ListenIP = '4.3.2.1'
ListenPort = 9
NewStreamTimeout = '1s'
PeerstoreWriteInterval = '1m0s'

[P2P.V2]
Enabled = true
AnnounceAddresses = ['a', 'b', 'c']
DefaultBootstrappers = ['12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw@foo:42/bar:10', '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw@test:99']
DeltaDial = '1m0s'
DeltaReconcile = '1s'
ListenAddresses = ['foo', 'bar']
`},
		{"Keeper", Config{Core: config.Core{Keeper: full.Keeper}}, `[Keeper]
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
		{"AutoPprof", Config{Core: config.Core{AutoPprof: full.AutoPprof}}, `[AutoPprof]
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
		{"Pyroscope", Config{Core: config.Core{Pyroscope: full.Pyroscope}}, `[Pyroscope]
ServerAddress = 'http://localhost:4040'
Environment = 'tests'
`},
		{"Sentry", Config{Core: config.Core{Sentry: full.Sentry}}, `[Sentry]
Debug = true
DSN = 'sentry-dsn'
Environment = 'dev'
Release = 'v1.2.3'
`},
		{"EVM", Config{EVM: full.EVM}, `[[EVM]]
ChainID = '1'
Enabled = false
BlockBackfillDepth = 100
BlockBackfillSkip = true
ChainType = 'Optimism'
FinalityDepth = 42
FlagsContractAddress = '0xae4E781a6218A8031764928E88d457937A954fC3'
LinkContractAddress = '0x538aAaB4ea120b2bC2fe5D296852D948F07D849e'
LogBackfillBatchSize = 17
LogPollInterval = '1m0s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 13
MinContractPayment = '9.223372036854775807 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
OperatorFactoryAddress = '0xa5B85635Be42F21f94F28034B7DA440EeFF0F418'
RPCDefaultBatchSize = 17
RPCBlockQueryDelay = 10

[EVM.Transactions]
ForwardersEnabled = true
MaxInFlight = 19
MaxQueued = 99
ReaperInterval = '1m0s'
ReaperThreshold = '1m0s'
ResendAfterThreshold = '1h0m0s'

[EVM.BalanceMonitor]
Enabled = true

[EVM.GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '9.223372036854775807 ether'
PriceMax = '281.474976710655 micro'
PriceMin = '13 wei'
LimitDefault = 12
LimitMax = 17
LimitMultiplier = '1.234'
LimitTransfer = 100
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

[EVM.HeadTracker]
HistoryDepth = 15
MaxBufferSize = 17
SamplingInterval = '1h0m0s'

[[EVM.KeySpecific]]
Key = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292'

[EVM.KeySpecific.GasEstimator]
PriceMax = '79.228162514264337593543950335 gether'

[EVM.NodePool]
PollFailureThreshold = 5
PollInterval = '1m0s'
SelectionMode = 'HighestHead'
SyncThreshold = 13

[EVM.OCR]
ContractConfirmations = 11
ContractTransmitterTransmitTimeout = '1m0s'
DatabaseTimeout = '1s'
ObservationGracePeriod = '1s'

[EVM.OCR2]
[EVM.OCR2.Automation]
GasLimit = 540

[[EVM.Nodes]]
Name = 'foo'
WSURL = 'wss://web.socket/test'
HTTPURL = 'https://foo.web'

[[EVM.Nodes]]
Name = 'bar'
WSURL = 'wss://web.socket/test'
HTTPURL = 'https://bar.com'

[[EVM.Nodes]]
Name = 'broadcast'
HTTPURL = 'http://broadcast.mirror'
SendOnly = true
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
Enabled = true
OCR2CachePollPeriod = '6h0m0s'
OCR2CacheTTL = '3m0s'
RequestTimeout = '1m3s'
TxTimeout = '13s'
TxSendFrequency = '42s'
TxMaxBatchSize = 17

[[Starknet.Nodes]]
Name = 'primary'
URL = 'http://stark.node'
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
		for n := range got.EVM[c].Nodes {
			if got.EVM[c].Nodes[n].WSURL == nil {
				got.EVM[c].Nodes[n].WSURL = new(models.URL)
			}
			if got.EVM[c].Nodes[n].SendOnly == nil {
				got.EVM[c].Nodes[n].SendOnly = ptr(true)
			}
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
		{name: "invalid", toml: invalidTOML, exp: `invalid configuration: 4 errors:
	- Database.Lock.LeaseRefreshInterval: invalid value (6s): must be less than or equal to half of LeaseDuration (10s)
	- EVM: 8 errors:
		- 1.ChainID: invalid value (1): duplicate - must be unique
		- 0.Nodes.1.Name: invalid value (foo): duplicate - must be unique
		- 3.Nodes.4.WSURL: invalid value (ws://dupe.com): duplicate - must be unique
		- 0: 3 errors:
			- GasEstimator.BumpTxDepth: invalid value (11): must be less than or equal to Transactions.MaxInFlight
			- GasEstimator: 6 errors:
				- BumpPercent: invalid value (1): may not be less than Geth's default of 10
				- TipCapDefault: invalid value (3 wei): must be greater than or equal to TipCapMinimum
				- FeeCapDefault: invalid value (3 wei): must be greater than or equal to TipCapDefault
				- PriceMin: invalid value (10 gwei): must be less than or equal to PriceDefault
				- PriceMax: invalid value (10 gwei): must be greater than or equal to PriceDefault
				- BlockHistory.BlockHistorySize: invalid value (0): must be greater than or equal to 1 with BlockHistory Mode
			- Nodes: 2 errors:
				- 0: 2 errors:
					- WSURL: missing: required for primary nodes
					- HTTPURL: missing: required for all nodes
				- 1.HTTPURL: missing: required for all nodes
		- 1: 6 errors:
			- ChainType: invalid value (Foo): must not be set with this chain id
			- Nodes: missing: must have at least one node
			- ChainType: invalid value (Foo): must be one of arbitrum, metis, optimism, xdai, optimismBedrock or omitted
			- HeadTracker.HistoryDepth: invalid value (30): must be equal to or reater than FinalityDepth
			- GasEstimator: 2 errors:
				- FeeCapDefault: invalid value (101 wei): must be equal to PriceMax (99 wei) since you are using FixedPrice estimation with gas bumping disabled in EIP1559 mode - PriceMax will be used as the FeeCap for transactions instead of FeeCapDefault
				- PriceMax: invalid value (1 gwei): must be greater than or equal to PriceDefault
			- KeySpecific.Key: invalid value (0xde709f2102306220921060314715629080e2fb77): duplicate - must be unique
		- 2: 5 errors:
			- ChainType: invalid value (Arbitrum): only "optimism" can be used with this chain id
			- Nodes: missing: must have at least one node
			- ChainType: invalid value (Arbitrum): must be one of arbitrum, metis, optimism, xdai, optimismBedrock or omitted
			- FinalityDepth: invalid value (0): must be greater than or equal to 1
			- MinIncomingConfirmations: invalid value (0): must be greater than or equal to 1
		- 3.Nodes: 5 errors:
				- 0: 3 errors:
					- Name: missing: required for all nodes
					- WSURL: missing: required for primary nodes
					- HTTPURL: empty: required for all nodes
				- 1: 3 errors:
					- Name: missing: required for all nodes
					- WSURL: invalid value (http): must be ws or wss
					- HTTPURL: missing: required for all nodes
				- 2: 3 errors:
					- Name: empty: required for all nodes
					- WSURL: missing: required for primary nodes
					- HTTPURL: invalid value (ws): must be http or https
				- 3.HTTPURL: missing: required for all nodes
				- 4.HTTPURL: missing: required for all nodes
		- 4: 2 errors:
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
			- Nodes: missing: must have at least one node`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var c Config
			require.NoError(t, config.DecodeTOML(strings.NewReader(tt.toml), &c))
			c.setDefaults()
			assertValidationError(t, &c, tt.exp)
		})
	}
}

func mustURL(s string) *models.URL {
	var u models.URL
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

func TestNewGeneralConfig_Logger(t *testing.T) {
	const (
		secrets   = "Secrets:\n"
		input     = "Input Configuration:\n"
		effective = "Effective Configuration, with defaults applied:\n"
	)
	tests := []struct {
		name         string
		inputConfig  string
		inputSecrets string

		wantConfig    string
		wantEffective string
		wantSecrets   string
	}{
		{name: "empty", wantEffective: emptyEffectiveTOML},
		{name: "full", inputSecrets: secretsFullTOML, inputConfig: fullTOML,
			wantConfig: fullTOML, wantEffective: fullTOML, wantSecrets: secretsFullRedactedTOML},
		{name: "multi-chain", inputSecrets: secretsMultiTOML, inputConfig: multiChainTOML,
			wantConfig: multiChainTOML, wantEffective: multiChainEffectiveTOML, wantSecrets: secretsMultiRedactedTOML},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lggr, observed := logger.TestLoggerObserved(t, zapcore.InfoLevel)
			opts := GeneralConfigOpts{SkipEnv: true}
			require.NoError(t, opts.ParseTOML(tt.inputConfig, tt.inputSecrets))
			c, err := opts.New(lggr)
			require.NoError(t, err)
			c.LogConfiguration(lggr.Info)

			inputLogs := observed.FilterMessageSnippet(secrets).All()
			if assert.Len(t, inputLogs, 1) {
				got := strings.TrimPrefix(inputLogs[0].Message, secrets)
				assert.Equal(t, tt.wantSecrets, got)
			}

			inputLogs = observed.FilterMessageSnippet(input).All()
			if assert.Len(t, inputLogs, 1) {
				got := strings.TrimPrefix(inputLogs[0].Message, input)
				assert.Equal(t, tt.wantConfig, got)
			}

			inputLogs = observed.FilterMessageSnippet(effective).All()
			if assert.Len(t, inputLogs, 1) {
				got := strings.TrimPrefix(inputLogs[0].Message, effective)
				assert.Equal(t, tt.wantEffective, got)
			}
		})
	}
}

func TestNewGeneralConfig_ParsingError_InvalidSyntax(t *testing.T) {
	invalidTOML := "{ bad syntax {"
	var opts GeneralConfigOpts
	err := opts.ParseTOML(invalidTOML, secretsFullTOML)
	assert.EqualError(t, err, "failed to decode config TOML: toml: invalid character at start of key: {")
}

func TestNewGeneralConfig_ParsingError_DuplicateField(t *testing.T) {
	invalidTOML := `Dev = false
Dev = true`
	var opts GeneralConfigOpts
	err := opts.ParseTOML(invalidTOML, secretsFullTOML)
	assert.EqualError(t, err, "failed to decode config TOML: toml: key Dev is already defined")
}

func TestNewGeneralConfig_SecretsOverrides(t *testing.T) {
	// Provide a keystore password file and an env var with DB URL
	const PWD_OVERRIDE = "great_password"
	const DBURL_OVERRIDE = "http://user@db"

	t.Setenv("CL_DATABASE_URL", DBURL_OVERRIDE)

	// Check for two overrides
	var opts GeneralConfigOpts
	require.NoError(t, opts.ParseTOML(fullTOML, secretsFullTOML))
	c, err := opts.New(logger.TestLogger(t))
	assert.NoError(t, err)
	c.SetPasswords(ptr(PWD_OVERRIDE), nil)
	assert.Equal(t, PWD_OVERRIDE, c.KeystorePassword())
	dbURL := c.DatabaseURL()
	assert.Equal(t, DBURL_OVERRIDE, (&dbURL).String())
}

func TestSecrets_Validate(t *testing.T) {
	for _, tt := range []struct {
		name string
		toml string
		exp  string
	}{
		{name: "partial",
			toml: `Explorer.AccessKey = "access_key"
Explorer.Secret = "secret"`,
			exp: `invalid secrets: 2 errors:
	- Database.URL: empty: must be provided and non-empty
	- Password.Keystore: empty: must be provided and non-empty`},

		{name: "invalid-urls",
			toml: `[Database]
URL = "postgresql://user:passlocalhost:5432/asdf"
BackupURL = "foo-bar?password=asdf"`,
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
		{name: "duplicate-mercury-credentials-disallowed",
			toml: `
[Mercury]
[[Mercury.Credentials]]
URL = "http://example.com/foo"
[[Mercury.Credentials]]
URL = "http://example.COM/foo"`,
			exp: `invalid secrets: 3 errors:
	- Database.URL: empty: must be provided and non-empty
	- Password.Keystore: empty: must be provided and non-empty
	- Mercury.Credentials: may not contain duplicate URLs`},
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
	c.EVM = evmcfg.EVMConfigs{{ChainID: utils.NewBigI(99999133712345)}}
	c.Solana = solana.SolanaConfigs{{ChainID: ptr("unknown solana chain")}}
	c.Starknet = starknet.StarknetConfigs{{ChainID: ptr("unknown starknet chain")}}
	c.setDefaults()
	if s, err := c.TOMLString(); assert.NoError(t, err) {
		t.Log(s, err)
	}
	cfgtest.AssertFieldsNotNil(t, c.Core)
}

func Test_validateEnv(t *testing.T) {
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("DATABASE_URL", "foo")
	assert.ErrorContains(t, validateEnv(), `invalid environment: environment variables DATABASE_URL and CL_DATABASE_URL must be equal, or CL_DATABASE_URL must not be set`)

	t.Setenv("CL_DATABASE_URL", "foo")
	assert.NoError(t, validateEnv())

	t.Setenv("CL_DATABASE_URL", "bar")
	t.Setenv("GAS_UPDATER_ENABLED", "true")
	t.Setenv("ETH_GAS_BUMP_TX_DEPTH", "7")
	assert.ErrorContains(t, validateEnv(), `invalid environment: 3 errors:
	- environment variables DATABASE_URL and CL_DATABASE_URL must be equal, or CL_DATABASE_URL must not be set
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
				c.SetFrom(&f)
			}
			ts, err := c.TOMLString()
			require.NoError(t, err)
			assert.Equal(t, tt.exp, ts)
		})
	}
}
