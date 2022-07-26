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
	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	relayutils "github.com/smartcontractkit/chainlink-relay/pkg/utils"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	tercfg "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	legacy "github.com/smartcontractkit/chainlink/core/config"
	config "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
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
	multiChain     = Config{
		Core: config.Core{
			RootDir: ptr("my/root/dir"),

			Database: &config.Database{
				Listener: &config.DatabaseListener{
					FallbackPollInterval: models.MustNewDuration(2 * time.Minute),
				},
			},
			Log: &config.Log{
				JSONConsole: ptr(true),
			},
			JobPipeline: &config.JobPipeline{
				DefaultHTTPRequestTimeout: models.MustNewDuration(30 * time.Second),
			},
			OCR2: &config.OCR2{
				DatabaseTimeout: models.MustNewDuration(20 * time.Second),
			},
			OCR: &config.OCR{
				BlockchainTimeout: models.MustNewDuration(5 * time.Second),
			},
			P2P: &config.P2P{
				IncomingMessageBufferSize: ptr[int64](999),
			},
			Keeper: &config.Keeper{
				GasPriceBufferPercent: ptr[uint32](10),
			},
			AutoPprof: &config.AutoPprof{
				CPUProfileRate: ptr[int64](7),
			},
		},
		EVM: []*EVMConfig{
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
					GasEstimator: &evmcfg.GasEstimator{
						PriceDefault: utils.NewBigI(math.MaxInt64).Wei(),
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
					GasEstimator: &evmcfg.GasEstimator{
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
		Solana: []*SolanaConfig{
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
		Terra: []*TerraConfig{
			{
				ChainID: ptr("Columbus-5"),
				Chain: tercfg.Chain{
					MaxMsgsPerBatch: ptr[int64](13),
				},
				Nodes: []*tercfg.Node{
					{Name: ptr("primary"), TendermintURL: relayutils.MustParseURL("http://columbus.terra.com")},
				}},
			{
				ChainID: ptr("Bombay-12"),
				Chain: tercfg.Chain{
					BlocksUntilTxTimeout: ptr[int64](20),
				},
				Nodes: []*tercfg.Node{
					{Name: ptr("primary"), TendermintURL: relayutils.MustParseURL("http://bombay.terra.com")},
				}},
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

	global := Config{
		Core: config.Core{
			ExplorerURL:         mustURL("http://explorer.url"),
			InsecureFastScrypt:  ptr(true),
			RootDir:             ptr("test/root/dir"),
			ShutdownGracePeriod: models.MustNewDuration(10 * time.Second),
		},
	}

	full := global
	full.Feature = &config.Feature{
		FeedsManager: ptr(true),
		LogPoller:    ptr(true),
		UICSAKeys:    ptr(true),
	}
	full.Database = &config.Database{
		DefaultIdleInTxSessionTimeout: models.MustNewDuration(time.Minute),
		DefaultLockTimeout:            models.MustNewDuration(time.Hour),
		DefaultQueryTimeout:           models.MustNewDuration(time.Second),

		MigrateOnStartup: ptr(true),
		ORMMaxIdleConns:  ptr[int64](7),
		ORMMaxOpenConns:  ptr[int64](13),
		Listener: &config.DatabaseListener{
			MaxReconnectDuration: models.MustNewDuration(time.Minute),
			MinReconnectInterval: models.MustNewDuration(5 * time.Minute),
			FallbackPollInterval: models.MustNewDuration(2 * time.Minute),
		},
		Lock: &config.DatabaseLock{
			LeaseDuration:        &minute,
			LeaseRefreshInterval: &second,
		},
		Backup: &config.DatabaseBackup{
			Dir:              ptr("test/backup/dir"),
			Frequency:        &hour,
			Mode:             &legacy.DatabaseBackupModeFull,
			OnVersionUpgrade: ptr(true),
			URL:              mustURL("http://test.back.up/fake"),
		},
	}
	full.TelemetryIngress = &config.TelemetryIngress{
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
	full.Log = &config.Log{
		JSONConsole:     ptr(true),
		FileDir:         ptr("log/file/dir"),
		DatabaseQueries: ptr(true),
		FileMaxSize:     ptr[utils.FileSize](100 * utils.GB),
		FileMaxAgeDays:  ptr[int64](17),
		FileMaxBackups:  ptr[int64](9),
		UnixTS:          ptr(true),
	}
	full.WebServer = &config.WebServer{
		AllowOrigins:            ptr("*"),
		BridgeResponseURL:       mustURL("https://bridge.response"),
		HTTPWriteTimeout:        models.MustNewDuration(time.Minute),
		HTTPPort:                ptr[uint16](56),
		SecureCookies:           ptr(true),
		SessionTimeout:          models.MustNewDuration(time.Hour),
		SessionReaperExpiration: models.MustNewDuration(7 * 24 * time.Hour),
		MFA: &config.WebServerMFA{
			RPID:     ptr("test-rpid"),
			RPOrigin: ptr("test-rp-origin"),
		},
		RateLimit: &config.WebServerRateLimit{
			Authenticated:         ptr[int64](42),
			AuthenticatedPeriod:   models.MustNewDuration(time.Second),
			Unauthenticated:       ptr[int64](7),
			UnauthenticatedPeriod: models.MustNewDuration(time.Minute),
		},
		TLS: &config.WebServerTLS{
			CertPath:      ptr("tls/cert/path"),
			Host:          ptr("tls-host"),
			KeyPath:       ptr("tls/key/path"),
			HTTPSPort:     ptr[uint16](6789),
			ForceRedirect: ptr(true),
		},
	}
	full.JobPipeline = &config.JobPipeline{
		HTTPRequestMaxSize:        ptr[utils.FileSize](100 * utils.MB),
		DefaultHTTPRequestTimeout: models.MustNewDuration(time.Minute),
		ExternalInitiatorsEnabled: ptr(true),
		MaxRunDuration:            models.MustNewDuration(time.Hour),
		ReaperInterval:            models.MustNewDuration(4 * time.Hour),
		ReaperThreshold:           models.MustNewDuration(7 * 24 * time.Hour),
		ResultWriteQueueDepth:     ptr[uint32](10),
	}
	full.FluxMonitor = &config.FluxMonitor{
		DefaultTransactionQueueDepth: ptr[uint32](100),
		SimulateTransactions:         ptr(true),
	}
	full.OCR2 = &config.OCR2{
		Enabled:                            ptr(true),
		ContractConfirmations:              ptr[uint32](11),
		BlockchainTimeout:                  models.MustNewDuration(3 * time.Second),
		ContractPollInterval:               models.MustNewDuration(time.Hour),
		ContractSubscribeInterval:          models.MustNewDuration(time.Minute),
		ContractTransmitterTransmitTimeout: models.MustNewDuration(time.Minute),
		DatabaseTimeout:                    models.MustNewDuration(8 * time.Second),
		KeyBundleID:                        ptr(models.MustSha256HashFromHex("7a5f66bbe6594259325bf2b4f5b1a9c9")),
	}
	full.OCR = &config.OCR{
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
	full.P2P = &config.P2P{
		IncomingMessageBufferSize: ptr[int64](13),
		OutgoingMessageBufferSize: ptr[int64](17),
		TraceLogging:              ptr(true),
		V1: &config.P2PV1{
			AnnounceIP:                       mustIP("1.2.3.4"),
			AnnouncePort:                     ptr[uint16](1234),
			BootstrapCheckInterval:           models.MustNewDuration(time.Minute),
			DefaultBootstrapPeers:            &[]string{"foo", "bar", "should", "these", "be", "typed"},
			DHTAnnouncementCounterUserPrefix: ptr[uint32](4321),
			DHTLookupInterval:                ptr[int64](9),
			ListenIP:                         mustIP("4.3.2.1"),
			ListenPort:                       ptr[uint16](9),
			NewStreamTimeout:                 models.MustNewDuration(time.Second),
			PeerID:                           mustPeerID("12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw"),
			PeerstoreWriteInterval:           models.MustNewDuration(time.Minute),
		},
		V2: &config.P2PV2{
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
	full.Keeper = &config.Keeper{
		DefaultTransactionQueueDepth: ptr[uint32](17),
		GasPriceBufferPercent:        ptr[uint32](12),
		GasTipCapBufferPercent:       ptr[uint32](43),
		BaseFeeBufferPercent:         ptr[uint32](89),
		MaximumGracePeriod:           ptr[int64](31),
		RegistryCheckGasOverhead:     utils.NewBigI(90),
		RegistryPerformGasOverhead:   utils.NewBig(new(big.Int).SetUint64(math.MaxUint64)),
		RegistrySyncInterval:         models.MustNewDuration(time.Hour),
		RegistrySyncUpkeepQueueSize:  ptr[uint32](31),
		TurnLookBack:                 ptr[int64](91),
		TurnFlagEnabled:              ptr(true),
		UpkeepCheckGasPriceEnabled:   ptr(true),
	}
	full.AutoPprof = &config.AutoPprof{
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
	full.Sentry = &config.Sentry{
		Debug:       ptr(true),
		DSN:         ptr("sentry-dsn"),
		Environment: ptr("dev"),
		Release:     ptr("v1.2.3"),
	}
	full.EVM = []*EVMConfig{
		{
			ChainID: utils.NewBigI(1),
			Enabled: ptr(false),
			Chain: evmcfg.Chain{
				BalanceMonitor: &evmcfg.BalanceMonitor{
					Enabled:    ptr(true),
					BlockDelay: ptr[uint16](17),
				},
				BlockBackfillDepth:   ptr[uint32](100),
				BlockBackfillSkip:    ptr(true),
				ChainType:            ptr("Optimism"),
				FinalityDepth:        ptr[uint32](42),
				FlagsContractAddress: mustAddress("0xae4E781a6218A8031764928E88d457937A954fC3"),

				GasEstimator: &evmcfg.GasEstimator{
					Mode:               ptr("L2Suggested"),
					EIP1559DynamicFees: ptr(true),
					BumpPercent:        ptr[uint16](10),
					BumpThreshold:      ptr[uint32](6),
					BumpTxDepth:        ptr[uint16](6),
					BumpMin:            utils.NewBigI(100).Wei(),
					FeeCapDefault:      utils.NewBigI(math.MaxInt64).Wei(),
					LimitDefault:       ptr[uint32](12),
					LimitMultiplier:    mustDecimal("1.234"),
					LimitTransfer:      ptr[uint32](100),
					LimitOCRJobType:    ptr[uint32](1001),
					LimitDRJobType:     ptr[uint32](1002),
					LimitVRFJobType:    ptr[uint32](1003),
					LimitFMJobType:     ptr[uint32](1004),
					LimitKeeperJobType: ptr[uint32](1005),
					TipCapDefault:      utils.NewBigI(2).Wei(),
					TipCapMinimum:      utils.NewBigI(1).Wei(),
					PriceDefault:       utils.NewBigI(math.MaxInt64).Wei(),
					PriceMax:           utils.NewBig(utils.HexToBig("FFFFFFFFFFFF")).Wei(),
					PriceMin:           utils.NewBigI(13).Wei(),

					BlockHistory: &evmcfg.BlockHistoryEstimator{
						BatchSize:                 ptr[uint32](17),
						BlockDelay:                ptr[uint16](10),
						BlockHistorySize:          ptr[uint16](12),
						EIP1559FeeCapBufferBlocks: ptr[uint16](13),
						TransactionPercentile:     ptr[uint16](15),
					},
				},

				KeySpecific: []evmcfg.KeySpecific{
					{
						Key: mustAddress("0x2a3e23c6f242F5345320814aC8a1b4E58707D292"),
						GasEstimator: &evmcfg.KeySpecificGasEstimator{
							PriceMax: utils.NewBig(utils.HexToBig("FFFFFFFFFFFFFFFFFFFFFFFF")).Wei(),
						},
					},
				},

				LinkContractAddress:  mustAddress("0x538aAaB4ea120b2bC2fe5D296852D948F07D849e"),
				LogBackfillBatchSize: ptr[uint32](17),
				LogPollInterval:      &minute,

				MaxInFlightTransactions:  ptr[uint32](19),
				MaxQueuedTransactions:    ptr[uint32](99),
				MinIncomingConfirmations: ptr[uint32](13),
				MinimumContractPayment:   assets.NewLinkFromJuels(math.MaxInt64),

				NonceAutoSync: ptr(true),

				OperatorFactoryAddress: mustAddress("0xa5B85635Be42F21f94F28034B7DA440EeFF0F418"),

				RPCDefaultBatchSize:    ptr[uint32](17),
				TxReaperInterval:       &minute,
				TxReaperThreshold:      &minute,
				TxResendAfterThreshold: &hour,
				UseForwarders:          ptr(true),

				HeadTracker: &evmcfg.HeadTracker{
					BlockEmissionIdleWarningThreshold: &hour,
					HistoryDepth:                      ptr[uint32](15),
					MaxBufferSize:                     ptr[uint32](17),
					SamplingInterval:                  &hour,
				},

				NodePool: &evmcfg.NodePool{
					NoNewHeadsThreshold:  &minute,
					PollFailureThreshold: ptr[uint32](5),
					PollInterval:         &minute,
				},
				OCR: &evmcfg.OCR{
					ContractConfirmations:              ptr[uint16](11),
					ContractTransmitterTransmitTimeout: &minute,
					DatabaseTimeout:                    &second,
					ObservationTimeout:                 &second,
					ObservationGracePeriod:             &second,
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
	full.Solana = []*SolanaConfig{
		{
			ChainID: ptr("mainnet"),
			Enabled: ptr(false),
			Chain: solcfg.Chain{
				BalancePollPeriod:   relayutils.MustNewDuration(time.Minute),
				ConfirmPollPeriod:   relayutils.MustNewDuration(time.Second),
				OCR2CachePollPeriod: relayutils.MustNewDuration(time.Minute),
				OCR2CacheTTL:        relayutils.MustNewDuration(time.Hour),
				TxTimeout:           relayutils.MustNewDuration(time.Hour),
				TxRetryTimeout:      relayutils.MustNewDuration(time.Minute),
				TxConfirmTimeout:    relayutils.MustNewDuration(time.Second),
				SkipPreflight:       ptr(true),
				Commitment:          ptr("banana"),
				MaxRetries:          ptr[int64](7),
			},
			Nodes: []*solcfg.Node{
				{Name: ptr("primary"), URL: relayutils.MustParseURL("http://solana.web")},
				{Name: ptr("foo"), URL: relayutils.MustParseURL("http://solana.foo")},
				{Name: ptr("bar"), URL: relayutils.MustParseURL("http://solana.bar")},
			},
		},
	}
	full.Terra = []*TerraConfig{
		{
			ChainID: ptr("Bombay-12"),
			Enabled: ptr(true),
			Chain: tercfg.Chain{
				BlockRate:             relayutils.MustNewDuration(time.Minute),
				BlocksUntilTxTimeout:  ptr[int64](12),
				ConfirmPollPeriod:     relayutils.MustNewDuration(time.Second),
				FallbackGasPriceULuna: mustDecimal("0.001"),
				FCDURL:                relayutils.MustParseURL("http://terra.com"),
				GasLimitMultiplier:    mustDecimal("1.2"),
				MaxMsgsPerBatch:       ptr[int64](17),
				OCR2CachePollPeriod:   relayutils.MustNewDuration(time.Minute),
				OCR2CacheTTL:          relayutils.MustNewDuration(time.Hour),
				TxMsgTimeout:          relayutils.MustNewDuration(time.Second),
			},
			Nodes: []*tercfg.Node{
				{Name: ptr("primary"), TendermintURL: relayutils.MustParseURL("http://tender.mint")},
				{Name: ptr("foo"), TendermintURL: relayutils.MustParseURL("http://foo.url")},
				{Name: ptr("bar"), TendermintURL: relayutils.MustParseURL("http://bar.web")},
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
		{"Feature", Config{Core: config.Core{Feature: full.Feature}}, `
[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true
`},
		{"Database", Config{Core: config.Core{Database: full.Database}}, `
[Database]
DefaultIdleInTxSessionTimeout = '1m0s'
DefaultLockTimeout = '1h0m0s'
DefaultQueryTimeout = '1s'
MigrateOnStartup = true
ORMMaxIdleConns = 7
ORMMaxOpenConns = 13

[Database.Backup]
Dir = 'test/backup/dir'
Frequency = '1h0m0s'
Mode = 'full'
OnVersionUpgrade = true
URL = 'http://test.back.up/fake'

[Database.Listener]
MaxReconnectDuration = '1m0s'
MinReconnectInterval = '5m0s'
FallbackPollInterval = '2m0s'

[Database.Lock]
LeaseDuration = '1m0s'
LeaseRefreshInterval = '1s'
`},
		{"TelemetryIngress", Config{Core: config.Core{TelemetryIngress: full.TelemetryIngress}}, `
[TelemetryIngress]
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
		{"Log", Config{Core: config.Core{Log: full.Log}}, `
[Log]
DatabaseQueries = true
FileDir = 'log/file/dir'
FileMaxSize = '100.00gb'
FileMaxAgeDays = 17
FileMaxBackups = 9
JSONConsole = true
UnixTS = true
`},
		{"WebServer", Config{Core: config.Core{WebServer: full.WebServer}}, `
[WebServer]
AllowOrigins = '*'
BridgeResponseURL = 'https://bridge.response'
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
		{"FluxMonitor", Config{Core: config.Core{FluxMonitor: full.FluxMonitor}}, `
[FluxMonitor]
DefaultTransactionQueueDepth = 100
SimulateTransactions = true
`},
		{"JobPipeline", Config{Core: config.Core{JobPipeline: full.JobPipeline}}, `
[JobPipeline]
DefaultHTTPRequestTimeout = '1m0s'
ExternalInitiatorsEnabled = true
HTTPRequestMaxSize = '100.00mb'
MaxRunDuration = '1h0m0s'
ReaperInterval = '4h0m0s'
ReaperThreshold = '168h0m0s'
ResultWriteQueueDepth = 10
`},
		{"OCR", Config{Core: config.Core{OCR: full.OCR}}, `
[OCR]
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
		{"OCR2", Config{Core: config.Core{OCR2: full.OCR2}}, `
[OCR2]
Enabled = true
ContractConfirmations = 11
BlockchainTimeout = '3s'
ContractPollInterval = '1h0m0s'
ContractSubscribeInterval = '1m0s'
ContractTransmitterTransmitTimeout = '1m0s'
DatabaseTimeout = '8s'
KeyBundleID = '7a5f66bbe6594259325bf2b4f5b1a9c900000000000000000000000000000000'
`},
		{"P2P", Config{Core: config.Core{P2P: full.P2P}}, `
[P2P]
IncomingMessageBufferSize = 13
OutgoingMessageBufferSize = 17
TraceLogging = true

[P2P.V1]
AnnounceIP = '1.2.3.4'
AnnouncePort = 1234
BootstrapCheckInterval = '1m0s'
DefaultBootstrapPeers = ['foo', 'bar', 'should', 'these', 'be', 'typed']
DHTAnnouncementCounterUserPrefix = 4321
DHTLookupInterval = 9
ListenIP = '4.3.2.1'
ListenPort = 9
NewStreamTimeout = '1s'
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw'
PeerstoreWriteInterval = '1m0s'

[P2P.V2]
AnnounceAddresses = ['a', 'b', 'c']
DefaultBootstrappers = ['12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw@foo:42/bar:10', '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw@test:99']
DeltaDial = '1m0s'
DeltaReconcile = '1s'
ListenAddresses = ['foo', 'bar']
`},
		{"Keeper", Config{Core: config.Core{Keeper: full.Keeper}}, `
[Keeper]
DefaultTransactionQueueDepth = 17
GasPriceBufferPercent = 12
GasTipCapBufferPercent = 43
BaseFeeBufferPercent = 89
MaximumGracePeriod = 31
RegistryCheckGasOverhead = '90'
RegistryPerformGasOverhead = '18446744073709551615'
RegistrySyncInterval = '1h0m0s'
RegistrySyncUpkeepQueueSize = 31
TurnLookBack = 91
TurnFlagEnabled = true
UpkeepCheckGasPriceEnabled = true
`},
		{"AutoPprof", Config{Core: config.Core{AutoPprof: full.AutoPprof}}, `
[AutoPprof]
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
		{"Sentry", Config{Core: config.Core{Sentry: full.Sentry}}, `
[Sentry]
Debug = true
DSN = 'sentry-dsn'
Environment = 'dev'
Release = 'v1.2.3'
`},
		{"EVM", Config{EVM: full.EVM}, `
[[EVM]]
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
MaxInFlightTransactions = 19
MaxQueuedTransactions = 99
MinIncomingConfirmations = 13
MinimumContractPayment = '9.223372036854775807 link'
NonceAutoSync = true
OperatorFactoryAddress = '0xa5B85635Be42F21f94F28034B7DA440EeFF0F418'
RPCDefaultBatchSize = 17
TxReaperInterval = '1m0s'
TxReaperThreshold = '1m0s'
TxResendAfterThreshold = '1h0m0s'
UseForwarders = true

[EVM.BalanceMonitor]
Enabled = true
BlockDelay = 17

[EVM.GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '9.223372036854775807 ether'
PriceMax = '281.474976710655 micro'
PriceMin = '13 wei'
LimitDefault = 12
LimitMultiplier = '1.234'
LimitTransfer = 100
LimitOCRJobType = 1001
LimitDRJobType = 1002
LimitVRFJobType = 1003
LimitFMJobType = 1004
LimitKeeperJobType = 1005
BumpMin = '100 wei'
BumpPercent = 10
BumpThreshold = 6
BumpTxDepth = 6
EIP1559DynamicFees = true
FeeCapDefault = '9.223372036854775807 ether'
TipCapDefault = '2 wei'
TipCapMinimum = '1 wei'

[EVM.GasEstimator.BlockHistory]
BatchSize = 17
BlockDelay = 10
BlockHistorySize = 12
EIP1559FeeCapBufferBlocks = 13
TransactionPercentile = 15

[EVM.HeadTracker]
BlockEmissionIdleWarningThreshold = '1h0m0s'
HistoryDepth = 15
MaxBufferSize = 17
SamplingInterval = '1h0m0s'

[[EVM.KeySpecific]]
Key = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292'

[EVM.KeySpecific.GasEstimator]
PriceMax = '79.228162514264337593543950335 gether'

[EVM.NodePool]
NoNewHeadsThreshold = '1m0s'
PollFailureThreshold = 5
PollInterval = '1m0s'

[EVM.OCR]
ContractConfirmations = 11
ContractTransmitterTransmitTimeout = '1m0s'
DatabaseTimeout = '1s'
ObservationTimeout = '1s'
ObservationGracePeriod = '1s'

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
		{"Solana", Config{Solana: full.Solana}, `
[[Solana]]
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
		{"Terra", Config{Terra: full.Terra}, `
[[Terra]]
ChainID = 'Bombay-12'
Enabled = true
BlockRate = '1m0s'
BlocksUntilTxTimeout = 12
ConfirmPollPeriod = '1s'
FallbackGasPriceULuna = '0.001'
FCDURL = 'http://terra.com'
GasLimitMultiplier = '1.2'
MaxMsgsPerBatch = 17
OCR2CachePollPeriod = '1m0s'
OCR2CacheTTL = '1h0m0s'
TxMsgTimeout = '1s'

[[Terra.Nodes]]
Name = 'primary'
TendermintURL = 'http://tender.mint'

[[Terra.Nodes]]
Name = 'foo'
TendermintURL = 'http://foo.url'

[[Terra.Nodes]]
Name = 'bar'
TendermintURL = 'http://bar.web'
`},
		{"full", full, fullTOML},
		{"multi-chain", multiChain, multiChainTOML},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.config.TOMLString()
			require.NoError(t, err)
			assert.Equal(t, tt.exp, s, diff.Diff(tt.exp, s))

			var got Config
			d := toml.NewDecoder(strings.NewReader(s)).DisallowUnknownFields()
			require.NoError(t, d.Decode(&got))
			ts, err := got.TOMLString()
			require.NoError(t, err)
			assert.Equal(t, tt.config, got, diff.Diff(s, ts))
		})
	}
}

func TestConfig_full(t *testing.T) {
	var got Config
	d := toml.NewDecoder(strings.NewReader(fullTOML)).DisallowUnknownFields()
	require.NoError(t, d.Decode(&got))
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
	var invalid Config
	d := toml.NewDecoder(strings.NewReader(invalidTOML)).DisallowUnknownFields()
	require.NoError(t, d.Decode(&invalid))
	if err := invalid.Validate(); assert.Error(t, err) {
		got := err.Error()
		exp := `3 errors:
	1) EVM: 3 errors:
		1) ChainID: invalid value 1: duplicate - must be unique
		2) 0: Nodes: 3 errors:
				1) Name: invalid value foo: duplicate - must be unique
				2) 0: HTTPURL: missing: required for all nodes
				3) 1: 2 errors:
					1) WSURL: missing: required for SendOnly nodes
					2) HTTPURL: missing: required for all nodes
		3) 1: Chain: KeySpecific: duplicate address: 0xde709f2102306220921060314715629080e2fb77
	2) Solana: 2 errors:
		1) ChainID: invalid value mainnet: duplicate - must be unique
		2) 1: Nodes: 3 errors:
				1) Name: invalid value bar: duplicate - must be unique
				2) 0: URL: missing: required for all nodes
				3) 1: URL: missing: required for all nodes
	3) Terra: 2 errors:
		1) ChainID: invalid value Bombay-12: duplicate - must be unique
		2) 0: Nodes: 3 errors:
				1) Name: invalid value test: duplicate - must be unique
				2) 0: TendermintURL: missing: required for all nodes
				3) 1: TendermintURL: missing: required for all nodes`
		assert.Equal(t, exp, got, diff.Diff(exp, got))
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

func ptr[T any](v T) *T {
	return &v
}

var (
	//go:embed testdata/config-multi-chain-effective.toml
	multiChainEffectiveTOML string
)

func TestNewGeneralConfig_Logger(t *testing.T) {
	const (
		input     = "Input Configuration:\n"
		effective = "Effective Configuration, with defaults applied:\n"
	)
	tests := []struct {
		name          string
		inputConfig   string
		wantConfig    string
		wantEffective string
	}{
		{name: "empty"},
		{name: "full", inputConfig: fullTOML, wantConfig: fullTOML, wantEffective: fullTOML},
		{name: "multi-chain", inputConfig: multiChainTOML, wantConfig: multiChainTOML, wantEffective: multiChainEffectiveTOML},
		// TODO: more test cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lggr, observed := logger.TestLoggerObserved(t, zapcore.InfoLevel)
			c, err := NewGeneralConfig(tt.inputConfig)
			require.NoError(t, err)
			c.LogConfiguration(lggr.Info)
			inputLogs := observed.FilterMessageSnippet(input).All()
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
