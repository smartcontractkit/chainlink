package main

import (
	_ "embed"
	"fmt"
	"math"
	"net"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/kylelemons/godebug/diff"
	"github.com/pelletier/go-toml/v2"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func ExampleConfig() {
	if s, err := prettyPrint(Config{
		RootDir: "my/root/dir",

		Database: &DatabaseConfig{
			URL: mustURL("http://local.postgres"),
		},
		Log: &LogConfig{
			Level: zapcore.WarnLevel,
		},
		JobPipeline: &JobPipelineConfig{
			DefaultHTTPTimeout: Duration(30 * time.Second),
		},
		OCR2: &OCR2Config{
			DatabaseTimeout: Duration(20 * time.Second),
		},
		OCR: &OCRConfig{
			BlockchainTimeout: Duration(5 * time.Second),
		},
		P2P: &P2PConfig{
			AnnouncePort: 999,
		},
		Keeper: &KeeperConfig{
			GasPriceBufferPercent: 10,
		},
		AutoPprof: &AutoPprofConfig{
			CPUProfileRate: 7,
		},
		EVM: map[string]EVMConfig{
			"1": {
				//TODO more fields
				Nodes: map[string]evmNode{
					"primary": {
						WSURL: mustURL("wss://web.socket/test"),
					},
					"secondary": {
						HTTPURL:  mustURL("http://broadcast.mirror"),
						SendOnly: true,
					},
				}},
			//TODO more chains
		},
	}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(s)
	}
	// Output:
	// RootDir = 'my/root/dir'
	//
	// [Database]
	// URL = 'http://local.postgres'
	//
	// [Log]
	// Level = 'warn'
	//
	// [JobPipeline]
	// DefaultHTTPTimeout = '30s'
	//
	// [OCR2]
	// DatabaseTimeout = '20s'
	//
	// [OCR]
	// BlockchainTimeout = '5s'
	//
	// [P2P]
	// AnnouncePort = 999
	//
	// [Keeper]
	// GasPriceBufferPercent = 10
	//
	// [AutoPprof]
	// CPUProfileRate = 7
	//
	// [EVM]
	//
	// [EVM.1]
	//
	// [EVM.1.Nodes]
	//
	// [EVM.1.Nodes.primary]
	// WSURL = 'wss://web.socket/test'
	//
	// [EVM.1.Nodes.secondary]
	// HTTPURL = 'http://broadcast.mirror'
	// SendOnly = true
}

//go:embed testdata/config-full.toml
var fullToml string //TODO copy and add comments for an example-config.toml?

func TestConfig_Marshal(t *testing.T) {
	global := Config{
		Dev:                  true,
		ExplorerAccessKey:    "test-access-key",
		ExplorerSecret:       "test-secret",
		ExplorerURL:          mustURL("http://explorer.url"),
		FlagsContractAddress: "0x1234",
		InsecureFastScrypt:   true,
		ReaperExpiration:     Duration(7 * 24 * time.Hour),
		RootDir:              "test/root/dir",
		ShutdownGracePeriod:  Duration(10 * time.Second),

		FeatureFeedsManager: true,
		FeatureUICSAKeys:    true,

		FeatureLogPoller: true,

		FMDefaultTransactionQueueDepth: 100,
		FMSimulateTransactions:         true,

		FeatureOffchainReporting2: true,
		FeatureOffchainReporting:  true,
	}

	full := global
	full.Database = &DatabaseConfig{
		URL:                           mustURL("http://data.base/url"),
		ListenerMaxReconnectDuration:  Duration(time.Minute),
		ListenerMinReconnectInterval:  Duration(5 * time.Minute),
		Migrate:                       true,
		ORMMaxIdleConns:               7,
		ORMMaxOpenConns:               13,
		TriggerFallbackDBPollInterval: Duration(2 * time.Minute),
		AdvisoryLockCheckInterval:     Duration(5 * time.Minute),
		AdvisoryLockID:                345982730592843,
		LockingMode:                   "advisory",
		BackupDir:                     "test/backup/dir",
		BackupOnVersionUpgrade:        true,
		BackupURL:                     mustURL("http://test.back.up/fake"),
	}
	full.TelemetryIngress = &TelemetryIngressConfig{
		UniConn:      true,
		Logging:      true,
		ServerPubKey: "test-pub-key",
		URL:          mustURL("https://prom.test"),
		BufferSize:   1234,
		MaxBatchSize: 4321,
		SendInterval: Duration(time.Minute),
		SendTimeout:  Duration(5 * time.Second),
		UseBatchSend: true,
	}
	full.Log = &LogConfig{
		JSONConsole:    true,
		FileDir:        "log/file/dir",
		SQL:            true,
		FileMaxSize:    100 * utils.GB,
		FileMaxAge:     17,
		FileMaxBackups: 9,
		UnixTS:         true,
	}
	full.WebServer = &WebServerConfig{
		AllowOrigins:                   "*",
		AuthenticatedRateLimit:         42,
		AuthenticatedRateLimitPeriod:   Duration(time.Second),
		BridgeResponseURL:              mustURL("https://bridge.response"),
		HTTPWriteTimeout:               Duration(time.Minute),
		Port:                           56,
		SecureCookies:                  true,
		SessionTimeout:                 Duration(time.Hour),
		UnAuthenticatedRateLimit:       7,
		UnAuthenticatedRateLimitPeriod: Duration(time.Minute),
		RPID:                           "test-rpid",
		RPOrigin:                       "test-rp-origin",
		TLSCertPath:                    "tls/cert/path",
		TLSHost:                        "tls-host",
		TLSKeyPath:                     "tls/key/path",
		TLSPort:                        6789,
		TLSRedirect:                    true,
	}
	full.JobPipeline = &JobPipelineConfig{
		DefaultHTTPLimit:          67,
		DefaultHTTPTimeout:        Duration(time.Minute),
		FeatureExternalInitiators: true,
		MaxRunDuration:            Duration(time.Hour),
		ReaperInterval:            Duration(4 * time.Hour),
		ReaperThreshold:           Duration(7 * 24 * time.Hour),
		ResultWriteQueueDepth:     10,
	}
	full.OCR2 = &OCR2Config{
		ContractConfirmations:              11,
		BlockchainTimeout:                  Duration(3 * time.Second),
		ContractPollInterval:               Duration(time.Hour),
		ContractSubscribeInterval:          Duration(time.Minute),
		ContractTransmitterTransmitTimeout: Duration(time.Minute),
		DatabaseTimeout:                    Duration(8 * time.Second),
		KeyBundleID:                        "test-bundle-id",
		MonitoringEndpoint:                 "test-mon-end",
	}
	full.OCR = &OCRConfig{
		ContractConfirmations:              11,
		ContractTransmitterTransmitTimeout: Duration(time.Minute),
		DatabaseTimeout:                    Duration(time.Second),
		ObservationGracePeriod:             Duration(time.Minute),
		ObservationTimeout:                 Duration(11 * time.Second),
		BlockchainTimeout:                  Duration(3 * time.Second),
		ContractPollInterval:               Duration(time.Hour),
		ContractSubscribeInterval:          Duration(time.Minute),
		DefaultTransactionQueueDepth:       12,
		KeyBundleID:                        "test-key-bundle-id",
		MonitoringEndpoint:                 "test-monitor",
		SimulateTransactions:               true,
		TransmitterAddress:                 "0x1234abcd",
		OutgoingMessageBufferSize:          7,
		IncomingMessageBufferSize:          3,
		DHTLookupInterval:                  9,
		BootstrapCheckInterval:             Duration(time.Minute),
		NewStreamTimeout:                   Duration(time.Second),
	}
	full.P2P = &P2PConfig{
		NetworkingStack:                  ocrnetworking.NetworkingStackV1V2,
		IncomingMessageBufferSize:        13,
		OutgoingMessageBufferSize:        17,
		AnnounceIP:                       net.ParseIP("1:2:3:4"),
		AnnouncePort:                     1234,
		BootstrapCheckInterval:           Duration(time.Minute),
		BootstrapPeers:                   []string{"foo", "bar", "should", "these", "be", "typed"},
		DHTAnnouncementCounterUserPrefix: 4321,
		DHTLookupInterval:                9,
		ListenIP:                         net.ParseIP("4:3:2:1"),
		ListenPort:                       9,
		NewStreamTimeout:                 Duration(time.Second),
		PeerID:                           "ASDF",
		PeerstoreWriteInterval:           Duration(time.Minute),
		V2AnnounceAddresses:              []string{"a", "b", "c"},
		V2Bootstrappers:                  []string{"1", "2", "3"},
		V2DeltaDial:                      Duration(time.Minute),
		V2DeltaReconcile:                 Duration(time.Second),
		V2ListenAddresses:                []string{"foo", "bar"},
	}
	full.Keeper = &KeeperConfig{
		CheckUpkeepGasPriceFeatureEnabled: true,
		DefaultTransactionQueueDepth:      17,
		GasPriceBufferPercent:             12,
		GasTipCapBufferPercent:            43,
		BaseFeeBufferPercent:              89,
		MaximumGracePeriod:                31,
		RegistryCheckGasOverhead:          90,
		RegistryPerformGasOverhead:        math.MaxUint64,
		RegistrySyncInterval:              Duration(time.Hour),
		RegistrySyncUpkeepQueueSize:       31,
		TurnLookBack:                      91,
		TurnFlagEnabled:                   true,
	}
	full.AutoPprof = &AutoPprofConfig{
		Enabled:              true,
		ProfileRoot:          "prof/root",
		PollInterval:         Duration(time.Minute),
		GatherDuration:       Duration(12 * time.Second),
		GatherTraceDuration:  Duration(13 * time.Second),
		MaxProfileSize:       utils.GB,
		CPUProfileRate:       7,
		MemProfileRate:       9,
		BlockProfileRate:     5,
		MutexProfileFraction: 2,
		MemThreshold:         utils.GB,
		GoroutineThreshold:   999,
	}
	full.EVM = map[string]EVMConfig{
		"1": {
			//TODO more fields
			Nodes: map[string]evmNode{
				"primary": {
					WSURL: mustURL("wss://web.socket/test"),
				},
				"secondary": {
					HTTPURL:  mustURL("http://broadcast.mirror"),
					SendOnly: true,
				},
			}},
		//TODO more chains
	}
	full.Solana = map[string]SolanaConfig{
		//TODO
	}
	full.Terra = map[string]TerraConfig{
		//TODO
	}

	for _, tt := range []struct {
		name   string
		config Config
		exp    string
	}{
		{"empty", Config{}, ``},
		{"global", global, `Dev = true
ExplorerAccessKey = 'test-access-key'
ExplorerSecret = 'test-secret'
ExplorerURL = 'http://explorer.url'
FlagsContractAddress = '0x1234'
InsecureFastScrypt = true
ReaperExpiration = '168h0m0s'
RootDir = 'test/root/dir'
ShutdownGracePeriod = '10s'
FeatureFeedsManager = true
FeatureUICSAKeys = true
FeatureLogPoller = true
FMDefaultTransactionQueueDepth = 100
FMSimulateTransactions = true
FeatureOffchainReporting2 = true
FeatureOffchainReporting = true
`},
		{"Database", Config{Database: full.Database}, `
[Database]
URL = 'http://data.base/url'
ListenerMaxReconnectDuration = '1m0s'
ListenerMinReconnectInterval = '5m0s'
Migrate = true
ORMMaxIdleConns = 7
ORMMaxOpenConns = 13
TriggerFallbackDBPollInterval = '2m0s'
AdvisoryLockCheckInterval = '5m0s'
AdvisoryLockID = 345982730592843
LockingMode = 'advisory'
BackupDir = 'test/backup/dir'
BackupOnVersionUpgrade = true
BackupURL = 'http://test.back.up/fake'
`},
		{"TelemetryIngress", Config{TelemetryIngress: full.TelemetryIngress}, `
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
		{"Log", Config{Log: full.Log}, `
[Log]
JSONConsole = true
FileDir = 'log/file/dir'
SQL = true
FileMaxSize = '100.00gb'
FileMaxAge = 17
FileMaxBackups = 9
UnixTS = true
`},
		{"WebServer", Config{WebServer: full.WebServer}, `
[WebServer]
AllowOrigins = '*'
AuthenticatedRateLimit = 42
AuthenticatedRateLimitPeriod = '1s'
BridgeResponseURL = 'https://bridge.response'
HTTPWriteTimeout = '1m0s'
Port = 56
SecureCookies = true
SessionTimeout = '1h0m0s'
UnAuthenticatedRateLimit = 7
UnAuthenticatedRateLimitPeriod = '1m0s'
RPID = 'test-rpid'
RPOrigin = 'test-rp-origin'
TLSCertPath = 'tls/cert/path'
TLSHost = 'tls-host'
TLSKeyPath = 'tls/key/path'
TLSPort = 6789
TLSRedirect = true
`},
		{"JobPipeline", Config{JobPipeline: full.JobPipeline}, `
[JobPipeline]
DefaultHTTPLimit = 67
DefaultHTTPTimeout = '1m0s'
FeatureExternalInitiators = true
MaxRunDuration = '1h0m0s'
ReaperInterval = '4h0m0s'
ReaperThreshold = '168h0m0s'
ResultWriteQueueDepth = 10
`},
		{"OCR", Config{OCR: full.OCR}, `
[OCR]
ContractConfirmations = 11
ContractTransmitterTransmitTimeout = '1m0s'
DatabaseTimeout = '1s'
ObservationGracePeriod = '1m0s'
ObservationTimeout = '11s'
BlockchainTimeout = '3s'
ContractPollInterval = '1h0m0s'
ContractSubscribeInterval = '1m0s'
DefaultTransactionQueueDepth = 12
KeyBundleID = 'test-key-bundle-id'
MonitoringEndpoint = 'test-monitor'
SimulateTransactions = true
TransmitterAddress = '0x1234abcd'
OutgoingMessageBufferSize = 7
IncomingMessageBufferSize = 3
DHTLookupInterval = 9
BootstrapCheckInterval = '1m0s'
NewStreamTimeout = '1s'
`},
		{"OCR2", Config{OCR2: full.OCR2}, `
[OCR2]
ContractConfirmations = 11
BlockchainTimeout = '3s'
ContractPollInterval = '1h0m0s'
ContractSubscribeInterval = '1m0s'
ContractTransmitterTransmitTimeout = '1m0s'
DatabaseTimeout = '8s'
KeyBundleID = 'test-bundle-id'
MonitoringEndpoint = 'test-mon-end'
`},
		{"P2P", Config{P2P: full.P2P}, `
[P2P]
NetworkingStack = 'V1V2'
IncomingMessageBufferSize = 13
OutgoingMessageBufferSize = 17
AnnouncePort = 1234
BootstrapCheckInterval = '1m0s'
BootstrapPeers = ['foo', 'bar', 'should', 'these', 'be', 'typed']
DHTAnnouncementCounterUserPrefix = 4321
DHTLookupInterval = 9
ListenPort = 9
NewStreamTimeout = '1s'
PeerID = 'ASDF'
PeerstoreWriteInterval = '1m0s'
V2AnnounceAddresses = ['a', 'b', 'c']
V2Bootstrappers = ['1', '2', '3']
V2DeltaDial = '1m0s'
V2DeltaReconcile = '1s'
V2ListenAddresses = ['foo', 'bar']
`},
		{"Keeper", Config{Keeper: full.Keeper}, `
[Keeper]
CheckUpkeepGasPriceFeatureEnabled = true
DefaultTransactionQueueDepth = 17
GasPriceBufferPercent = 12
GasTipCapBufferPercent = 43
BaseFeeBufferPercent = 89
MaximumGracePeriod = 31
RegistryCheckGasOverhead = 90
RegistryPerformGasOverhead = 18446744073709551615
RegistrySyncInterval = '1h0m0s'
RegistrySyncUpkeepQueueSize = 31
TurnLookBack = 91
TurnFlagEnabled = true
`},
		{"AutoPprof", Config{AutoPprof: full.AutoPprof}, `
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
		{"evm", Config{EVM: full.EVM}, `
[EVM]

[EVM.1]

[EVM.1.Nodes]

[EVM.1.Nodes.primary]
WSURL = 'wss://web.socket/test'

[EVM.1.Nodes.secondary]
HTTPURL = 'http://broadcast.mirror'
SendOnly = true
`},
		//TODO solana, terra
		{"full", full, fullToml},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, err := prettyPrint(tt.config)
			require.NoError(t, err)
			assert.Equal(t, tt.exp, s, diff.Diff(tt.exp, s))
		})
	}
}

//TODO TestConfig_Unmarshal

func mustURL(s string) *URL {
	var u URL
	if err := u.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return &u
}

var multiLineBreak = regexp.MustCompile("(\n){2,}")

//TODO hopefully not really necessary...
func prettyPrint(c Config) (string, error) {
	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	// remove runs of line breaks
	s := multiLineBreak.ReplaceAllLiteralString(string(b), "\n")
	// restore them preceding keys
	s = strings.Replace(s, "\n[", "\n\n[", -1)
	s = strings.TrimPrefix(s, "\n")
	return s, nil
}
