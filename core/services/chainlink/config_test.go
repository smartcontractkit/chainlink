package chainlink

import (
	_ "embed"
	"math"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/kylelemons/godebug/diff"
	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	solcfg "github.com/smartcontractkit/chainlink/core/chains/solana/config"
	tercfg "github.com/smartcontractkit/chainlink/core/chains/terra/config"
	legacy "github.com/smartcontractkit/chainlink/core/config"
	config "github.com/smartcontractkit/chainlink/core/config/v2"
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
				TriggerFallbackDBPollInterval: models.MustNewDuration(2 * time.Minute),
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
		EVM: []EVMConfig{
			{
				ChainID: utils.NewBigI(1),
				Chain: evmcfg.Chain{
					FinalityDepth: ptr[uint32](26),
				},
				Nodes: []evmcfg.Node{
					{
						Name:  ptr("primary"),
						WSURL: mustURL("wss://web.socket/test"),
					},
					{
						Name:     ptr("secondary"),
						HTTPURL:  mustURL("http://broadcast.mirror"),
						SendOnly: ptr(true),
					},
				}},
		},
		Solana: []SolanaConfig{
			{
				ChainID: "mainnet",
				Chain: solcfg.Chain{
					MaxRetries: ptr[int64](12),
				},
				Nodes: []solcfg.Node{
					{Name: "primary", URL: mustURL("http://solana.com")},
				},
			},
		},
		Terra: []TerraConfig{
			{
				ChainID: "Columbus-5",
				Chain: tercfg.Chain{
					MaxMsgsPerBatch: ptr[int64](13),
				},
				Nodes: []tercfg.Node{
					{Name: "primary", TendermintURL: mustURL("http://solana.com")},
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

	global := Config{
		Core: config.Core{
			Dev:                 ptr(true),
			ExplorerURL:         mustURL("http://explorer.url"),
			InsecureFastScrypt:  ptr(true),
			ReaperExpiration:    models.MustNewDuration(7 * 24 * time.Hour),
			RootDir:             ptr("test/root/dir"),
			ShutdownGracePeriod: models.MustNewDuration(10 * time.Second),
		},
	}

	full := global
	full.Feature = &config.Feature{
		FeedsManager:       ptr(true),
		LogPoller:          ptr(true),
		OffchainReporting2: ptr(true),
		OffchainReporting:  ptr(true),
		UICSAKeys:          ptr(true),
	}
	full.Database = &config.Database{
		DefaultIdleInTxSessionTimeout: models.MustNewDuration(time.Minute),
		DefaultLockTimeout:            models.MustNewDuration(time.Hour),
		DefaultQueryTimeout:           models.MustNewDuration(time.Second),

		MigrateOnStartup:              ptr(true),
		ORMMaxIdleConns:               ptr[int64](7),
		ORMMaxOpenConns:               ptr[int64](13),
		TriggerFallbackDBPollInterval: models.MustNewDuration(2 * time.Minute),
		Listener: &config.DatabaseListener{
			MaxReconnectDuration: models.MustNewDuration(time.Minute),
			MinReconnectInterval: models.MustNewDuration(5 * time.Minute),
		},
		Lock: &config.DatabaseLock{
			Mode:                  ptr("advisory"),
			AdvisoryCheckInterval: models.MustNewDuration(5 * time.Minute),
			AdvisoryID:            ptr[int64](345982730592843),
			LeaseDuration:         &minute,
			LeaseRefreshInterval:  &second,
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
		AllowOrigins:     ptr("*"),
		ExternalURL:      mustURL("https://bridge.response"),
		HTTPWriteTimeout: models.MustNewDuration(time.Minute),
		HTTPPort:         ptr[uint16](56),
		SecureCookies:    ptr(true),
		SessionTimeout:   models.MustNewDuration(time.Hour),
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
		HTTPRequestMaxSizeBytes:   ptr[int64](67),
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
		ContractConfirmations:              ptr[uint32](11),
		BlockchainTimeout:                  models.MustNewDuration(3 * time.Second),
		ContractPollInterval:               models.MustNewDuration(time.Hour),
		ContractSubscribeInterval:          models.MustNewDuration(time.Minute),
		ContractTransmitterTransmitTimeout: models.MustNewDuration(time.Minute),
		DatabaseTimeout:                    models.MustNewDuration(8 * time.Second),
		KeyBundleID:                        ptr(models.MustSha256HashFromHex("7a5f66bbe6594259325bf2b4f5b1a9c9")),
		MonitoringEndpoint:                 ptr("test-mon-end"),
	}
	full.OCR = &config.OCR{
		ObservationTimeout:           models.MustNewDuration(11 * time.Second),
		BlockchainTimeout:            models.MustNewDuration(3 * time.Second),
		ContractPollInterval:         models.MustNewDuration(time.Hour),
		ContractSubscribeInterval:    models.MustNewDuration(time.Minute),
		DefaultTransactionQueueDepth: ptr[uint32](12),
		KeyBundleID:                  ptr(models.MustSha256HashFromHex("acdd42797a8b921b2910497badc50006")),
		MonitoringEndpoint:           ptr("test-monitor"),
		SimulateTransactions:         ptr(true),
		TransmitterAddress:           ptr(ethkey.MustEIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")),
	}
	full.P2P = &config.P2P{
		IncomingMessageBufferSize: ptr[int64](13),
		OutgoingMessageBufferSize: ptr[int64](17),
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
			AnnounceAddresses:    &[]string{"a", "b", "c"},
			DefaultBootstrappers: &[]string{"1", "2", "3"},
			DeltaDial:            models.MustNewDuration(time.Minute),
			DeltaReconcile:       models.MustNewDuration(time.Second),
			ListenAddresses:      &[]string{"foo", "bar"},
		},
	}
	full.Keeper = &config.Keeper{
		CheckUpkeepGasPriceFeatureEnabled: ptr(true),
		DefaultTransactionQueueDepth:      ptr[uint32](17),
		GasPriceBufferPercent:             ptr[uint32](12),
		GasTipCapBufferPercent:            ptr[uint32](43),
		BaseFeeBufferPercent:              ptr[uint32](89),
		MaximumGracePeriod:                ptr[int64](31),
		RegistryCheckGasOverhead:          utils.NewBigI(90),
		RegistryPerformGasOverhead:        utils.NewBig(new(big.Int).SetUint64(math.MaxUint64)),
		RegistrySyncInterval:              models.MustNewDuration(time.Hour),
		RegistrySyncUpkeepQueueSize:       ptr[uint32](31),
		TurnLookBack:                      ptr[int64](91),
		TurnFlagEnabled:                   ptr(true),
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
	full.EVM = []EVMConfig{
		{
			ChainID: utils.NewBigI(1),
			Chain: evmcfg.Chain{
				BalanceMonitorEnabled:             ptr(true),
				BlockBackfillDepth:                ptr[uint32](100),
				BlockBackfillSkip:                 ptr(true),
				BlockEmissionIdleWarningThreshold: &hour,
				BlockHistoryEstimator: &evmcfg.BlockHistoryEstimator{
					BatchSize:                 ptr[uint32](17),
					BlockDelay:                ptr[uint16](10),
					BlockHistorySize:          ptr[uint16](12),
					EIP1559FeeCapBufferBlocks: ptr[uint16](13),
					TransactionPercentile:     ptr[uint16](15),
				},
				ChainType:            ptr("Optimism"),
				EIP1559DynamicFees:   ptr(true),
				FinalityDepth:        ptr[uint32](42),
				FlagsContractAddress: mustAddress("0xae4E781a6218A8031764928E88d457937A954fC3"),

				GasBumpPercent:     ptr[uint16](10),
				GasBumpTxDepth:     ptr[uint16](6),
				GasBumpWei:         utils.NewBigI(100),
				GasEstimatorMode:   ptr("L2Suggested"),
				GasFeeCapDefault:   utils.NewBigI(math.MaxInt64),
				GasLimitDefault:    utils.NewBigI(12),
				GasLimitMultiplier: mustDecimal("1.234"),
				GasPriceDefault:    utils.NewBigI(math.MaxInt64),
				GasTipCapDefault:   utils.NewBigI(2),
				GasTipCapMinimum:   utils.NewBigI(1),

				HeadTrackerHistoryDepth:     ptr[uint32](15),
				HeadTrackerMaxBufferSize:    ptr[uint32](17),
				HeadTrackerSamplingInterval: &hour,

				KeySpecific: []evmcfg.KeySpecific{
					{
						Key:            mustAddress("0x2a3e23c6f242F5345320814aC8a1b4E58707D292"),
						MaxGasPriceWei: utils.NewBig(utils.HexToBig("FFFFFFFFFFFFFFFFFFFFFFFF")),
					},
				},

				LinkContractAddress:  mustAddress("0x538aAaB4ea120b2bC2fe5D296852D948F07D849e"),
				LogBackfillBatchSize: ptr[uint32](17),
				LogPollInterval:      &minute,

				MaxGasPriceWei:           utils.NewBig(utils.HexToBig("FFFFFFFFFFFF")),
				MaxInFlightTransactions:  ptr[uint32](19),
				MaxQueuedTransactions:    ptr[uint32](99),
				MinGasPriceWei:           utils.NewBigI(13),
				MinIncomingConfirmations: ptr[uint32](13),
				MinimumContractPayment:   assets.NewLinkFromJuels(math.MaxInt64),

				NonceAutoSync: ptr(true),

				OperatorFactoryAddress: mustAddress("0xa5B85635Be42F21f94F28034B7DA440EeFF0F418"),

				OCRContractConfirmations:              ptr[uint16](11),
				OCRContractTransmitterTransmitTimeout: &minute,
				OCRDatabaseTimeout:                    &second,
				OCRObservationTimeout:                 &second,
				OCRObservationGracePeriod:             &second,
				OCR2ContractConfirmations:             ptr[uint16](7),

				RPCDefaultBatchSize:    ptr[uint32](17),
				TxReaperInterval:       &minute,
				TxReaperThreshold:      &minute,
				TxResendAfterThreshold: &hour,
				UseForwarders:          ptr(true),

				NodePool: &evmcfg.NodePool{
					NoNewHeadsThreshold:  &minute,
					PollFailureThreshold: ptr[uint32](5),
					PollInterval:         &minute,
				},
			},
			Nodes: []evmcfg.Node{
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
	full.Solana = []SolanaConfig{
		{
			ChainID: "mainnet",
			Enabled: ptr(false),
			Chain: solcfg.Chain{
				BalancePollPeriod:   models.MustNewDuration(time.Minute),
				ConfirmPollPeriod:   models.MustNewDuration(time.Second),
				OCR2CachePollPeriod: models.MustNewDuration(time.Minute),
				OCR2CacheTTL:        models.MustNewDuration(time.Hour),
				TxTimeout:           models.MustNewDuration(time.Hour),
				TxRetryTimeout:      models.MustNewDuration(time.Minute),
				TxConfirmTimeout:    models.MustNewDuration(time.Second),
				SkipPreflight:       ptr(true),
				Commitment:          ptr("banana"),
				MaxRetries:          ptr[int64](7),
			},
			Nodes: []solcfg.Node{
				{Name: "primary", URL: mustURL("http://solana.web")},
				{Name: "foo", URL: mustURL("http://solana.foo")},
				{Name: "bar", URL: mustURL("http://solana.bar")},
			},
		},
	}
	full.Terra = []TerraConfig{
		{
			ChainID: "Bombay-12",
			Chain: tercfg.Chain{
				BlockRate:             models.MustNewDuration(time.Minute),
				BlocksUntilTxTimeout:  ptr[int64](12),
				ConfirmPollPeriod:     models.MustNewDuration(time.Second),
				FallbackGasPriceULuna: mustDecimal("0.001"),
				FCDURL:                mustURL("http://terra.com"),
				GasLimitMultiplier:    mustDecimal("1.2"),
				MaxMsgsPerBatch:       ptr[int64](17),
				OCR2CachePollPeriod:   models.MustNewDuration(time.Minute),
				OCR2CacheTTL:          models.MustNewDuration(time.Hour),
				TxMsgTimeout:          models.MustNewDuration(time.Second),
			},
			Nodes: []tercfg.Node{
				{Name: "primary", TendermintURL: mustURL("http://tender.mint")},
				{Name: "foo", TendermintURL: mustURL("http://foo.url")},
				{Name: "bar", TendermintURL: mustURL("http://bar.web")},
			},
		},
	}

	for _, tt := range []struct {
		name   string
		config Config
		exp    string
	}{
		{"empty", Config{}, ``},
		{"global", global, `Dev = true
ExplorerURL = 'http://explorer.url'
InsecureFastScrypt = true
ReaperExpiration = '168h0m0s'
RootDir = 'test/root/dir'
ShutdownGracePeriod = '10s'
`},
		{"Feature", Config{Core: config.Core{Feature: full.Feature}}, `
[Feature]
FeedsManager = true
LogPoller = true
OffchainReporting2 = true
OffchainReporting = true
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
TriggerFallbackDBPollInterval = '2m0s'

[Database.Listener]
MaxReconnectDuration = '1m0s'
MinReconnectInterval = '5m0s'

[Database.Lock]
Mode = 'advisory'
AdvisoryCheckInterval = '5m0s'
AdvisoryID = 345982730592843
LeaseDuration = '1m0s'
LeaseRefreshInterval = '1s'

[Database.Backup]
Dir = 'test/backup/dir'
Frequency = '1h0m0s'
Mode = 'full'
OnVersionUpgrade = true
URL = 'http://test.back.up/fake'
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
ExternalURL = 'https://bridge.response'
HTTPWriteTimeout = '1m0s'
HTTPPort = 56
SecureCookies = true
SessionTimeout = '1h0m0s'

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
HTTPRequestMaxSizeBytes = 67
MaxRunDuration = '1h0m0s'
ReaperInterval = '4h0m0s'
ReaperThreshold = '168h0m0s'
ResultWriteQueueDepth = 10
`},
		{"OCR", Config{Core: config.Core{OCR: full.OCR}}, `
[OCR]
ObservationTimeout = '11s'
BlockchainTimeout = '3s'
ContractPollInterval = '1h0m0s'
ContractSubscribeInterval = '1m0s'
DefaultTransactionQueueDepth = 12
KeyBundleID = 'acdd42797a8b921b2910497badc5000600000000000000000000000000000000'
MonitoringEndpoint = 'test-monitor'
SimulateTransactions = true
TransmitterAddress = '0xa0788FC17B1dEe36f057c42B6F373A34B014687e'
`},
		{"OCR2", Config{Core: config.Core{OCR2: full.OCR2}}, `
[OCR2]
ContractConfirmations = 11
BlockchainTimeout = '3s'
ContractPollInterval = '1h0m0s'
ContractSubscribeInterval = '1m0s'
ContractTransmitterTransmitTimeout = '1m0s'
DatabaseTimeout = '8s'
KeyBundleID = '7a5f66bbe6594259325bf2b4f5b1a9c900000000000000000000000000000000'
MonitoringEndpoint = 'test-mon-end'
`},
		{"P2P", Config{Core: config.Core{P2P: full.P2P}}, `
[P2P]
IncomingMessageBufferSize = 13
OutgoingMessageBufferSize = 17

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
DefaultBootstrappers = ['1', '2', '3']
DeltaDial = '1m0s'
DeltaReconcile = '1s'
ListenAddresses = ['foo', 'bar']
`},
		{"Keeper", Config{Core: config.Core{Keeper: full.Keeper}}, `
[Keeper]
CheckUpkeepGasPriceFeatureEnabled = true
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
BalanceMonitorEnabled = true
BlockBackfillDepth = 100
BlockBackfillSkip = true
BlockEmissionIdleWarningThreshold = '1h0m0s'
ChainType = 'Optimism'
EIP1559DynamicFees = true
FinalityDepth = 42
FlagsContractAddress = '0xae4E781a6218A8031764928E88d457937A954fC3'
GasBumpPercent = 10
GasBumpTxDepth = 6
GasBumpWei = '100'
GasEstimatorMode = 'L2Suggested'
GasFeeCapDefault = '9223372036854775807'
GasLimitDefault = '12'
GasLimitMultiplier = '1.234'
GasPriceDefault = '9223372036854775807'
GasTipCapDefault = '2'
GasTipCapMinimum = '1'
HeadTrackerHistoryDepth = 15
HeadTrackerMaxBufferSize = 17
HeadTrackerSamplingInterval = '1h0m0s'
LinkContractAddress = '0x538aAaB4ea120b2bC2fe5D296852D948F07D849e'
LogBackfillBatchSize = 17
LogPollInterval = '1m0s'
MaxGasPriceWei = '281474976710655'
MaxInFlightTransactions = 19
MaxQueuedTransactions = 99
MinGasPriceWei = '13'
MinIncomingConfirmations = 13
MinimumContractPayment = '9223372036854775807'
NonceAutoSync = true
OCRContractConfirmations = 11
OCRContractTransmitterTransmitTimeout = '1m0s'
OCRDatabaseTimeout = '1s'
OCRObservationTimeout = '1s'
OCRObservationGracePeriod = '1s'
OCR2ContractConfirmations = 7
OperatorFactoryAddress = '0xa5B85635Be42F21f94F28034B7DA440EeFF0F418'
RPCDefaultBatchSize = 17
TxReaperInterval = '1m0s'
TxReaperThreshold = '1m0s'
TxResendAfterThreshold = '1h0m0s'
UseForwarders = true

[EVM.BlockHistoryEstimator]
BatchSize = 17
BlockDelay = 10
BlockHistorySize = 12
EIP1559FeeCapBufferBlocks = 13
TransactionPercentile = 15

[[EVM.KeySpecific]]
Key = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292'
MaxGasPriceWei = '79228162514264337593543950335'

[EVM.NodePool]
NoNewHeadsThreshold = '1m0s'
PollFailureThreshold = 5
PollInterval = '1m0s'

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
			s, err := prettyPrint(tt.config)
			require.NoError(t, err)
			assert.Equal(t, tt.exp, s, diff.Diff(tt.exp, s))

			var got Config
			require.NoError(t, toml.Unmarshal([]byte(s), &got))
			assert.Equal(t, tt.config, got)
		})
	}
}

//TODO TestConfig_Unmarshal

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
