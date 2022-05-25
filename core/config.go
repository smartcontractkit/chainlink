package main

import (
	"net"
	"net/url"
	"time"

	"go.uber.org/zap/zapcore"

	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Config struct {
	// General/misc
	Dev                  bool `toml:",omitempty"`
	ExplorerURL          *URL
	FlagsContractAddress string   `toml:",omitempty"`
	InsecureFastScrypt   bool     `toml:",omitempty"`
	ReaperExpiration     Duration `toml:",omitempty"`
	RootDir              string   `toml:",omitempty"`
	ShutdownGracePeriod  Duration `toml:",omitempty"`

	Database *DatabaseConfig

	TelemetryIngress *TelemetryIngressConfig

	Log *LogConfig

	WebServer *WebServerConfig

	// Feeds manager
	FeatureFeedsManager bool `toml:",omitempty"`
	FeatureUICSAKeys    bool `toml:",omitempty"`

	// LogPoller
	FeatureLogPoller bool `toml:",omitempty"`

	// Job Pipeline and tasks
	JobPipeline *JobPipelineConfig

	// Flux Monitor
	FMDefaultTransactionQueueDepth uint32 `toml:",omitempty"`
	FMSimulateTransactions         bool   `toml:",omitempty"`

	// OCR V2
	FeatureOffchainReporting2 bool `toml:",omitempty"`
	OCR2                      *OCR2Config

	// OCR V1
	FeatureOffchainReporting bool `toml:",omitempty"`
	OCR                      *OCRConfig

	// P2P Networking
	P2P *P2PConfig

	// Keeper
	Keeper *KeeperConfig

	// Debugging
	AutoPprof *AutoPprofConfig

	//TODO do we need [EVM] wide config too?
	EVM map[string]EVMConfig `toml:",omitempty"`

	Solana map[string]SolanaConfig `toml:",omitempty"`

	Terra map[string]TerraConfig `toml:",omitempty"`
}

type Secrets struct {
	DatabaseURL       *URL
	ExplorerAccessKey string `toml:",omitempty"`
	ExplorerSecret    string `toml:",omitempty"`
	//TODO more?
}

type EVMConfig struct {
	evmtypes.ChainTOMLCfg
	Nodes map[string]EVMNode
}

type EVMNode struct {
	WSURL    *URL
	HTTPURL  *URL
	SendOnly bool `toml:",omitempty"`
}

type SolanaConfig struct {
	SolanaChainCfg
	Nodes map[string]solanaNode
}

type SolanaChainCfg struct {
	BalancePollPeriod   Duration
	ConfirmPollPeriod   Duration
	OCR2CachePollPeriod Duration
	OCR2CacheTTL        Duration
	TxTimeout           Duration
	TxRetryTimeout      Duration
	TxConfirmTimeout    Duration
	SkipPreflight       bool   `toml:",omitempty"`
	Commitment          string `toml:",omitempty"`
	MaxRetries          int    `toml:",omitempty"`
}

type solanaNode struct {
	URL *URL
}

type TerraConfig struct {
	TerraChainCfg
	Nodes map[string]TerraNode
}

type TerraChainCfg struct {
	BlockRate             Duration
	BlocksUntilTxTimeout  int
	ConfirmPollPeriod     Duration
	FallbackGasPriceULuna string //TODO decimal number type?
	FCDURL                *URL
	GasLimitMultiplier    float64
	MaxMsgsPerBatch       int64
	OCR2CachePollPeriod   Duration
	OCR2CacheTTL          Duration
	TxMsgTimeout          Duration
}

type TerraNode struct {
	TendermintURL *URL
}

type DatabaseConfig struct {
	ListenerMaxReconnectDuration  Duration `toml:",omitempty"`
	ListenerMinReconnectInterval  Duration `toml:",omitempty"`
	Migrate                       bool     `toml:",omitempty"`
	ORMMaxIdleConns               int      `toml:",omitempty"`
	ORMMaxOpenConns               int      `toml:",omitempty"`
	TriggerFallbackDBPollInterval Duration `toml:",omitempty"`
	// Database Global Lock
	AdvisoryLockCheckInterval Duration `toml:",omitempty"`
	AdvisoryLockID            int64    `toml:",omitempty"`
	LockingMode               string   `toml:",omitempty"`
	LeaseLockDuration         Duration `toml:",omitempty"`
	LeaseLockRefreshInterval  Duration `toml:",omitempty"`
	// Database Autobackups
	BackupDir              string   `toml:",omitempty"`
	BackupFrequency        Duration `toml:",omitempty"`
	BackupMode             string   `toml:",omitempty"`
	BackupOnVersionUpgrade bool     `toml:",omitempty"`
	BackupURL              *URL
}

type TelemetryIngressConfig struct {
	UniConn      bool   `toml:",omitempty"`
	Logging      bool   `toml:",omitempty"`
	ServerPubKey string `toml:",omitempty"`
	URL          *URL
	BufferSize   uint     `toml:",omitempty"`
	MaxBatchSize uint     `toml:",omitempty"`
	SendInterval Duration `toml:",omitempty"`
	SendTimeout  Duration `toml:",omitempty"`
	UseBatchSend bool     `toml:",omitempty"`
}

type LogConfig struct {
	JSONConsole bool           `toml:",omitempty"`
	FileDir     string         `toml:",omitempty"`
	Level       zapcore.Level  `toml:",omitempty"`
	SQL         bool           `toml:",omitempty"`
	FileMaxSize utils.FileSize `toml:",omitempty"`
	//TODO units?
	FileMaxAge     int64 `toml:",omitempty"`
	FileMaxBackups int64 `toml:",omitempty"`
	UnixTS         bool  `toml:",omitempty"`
}

type WebServerConfig struct {
	// Web Server
	AllowOrigins                   string   `toml:",omitempty"`
	AuthenticatedRateLimit         int64    `toml:",omitempty"`
	AuthenticatedRateLimitPeriod   Duration `toml:",omitempty"`
	BridgeResponseURL              *URL
	HTTPWriteTimeout               Duration `toml:",omitempty"`
	Port                           uint16   `toml:",omitempty"`
	SecureCookies                  bool     `toml:",omitempty"`
	SessionTimeout                 Duration `toml:",omitempty"`
	UnAuthenticatedRateLimit       int64    `toml:",omitempty"`
	UnAuthenticatedRateLimitPeriod Duration `toml:",omitempty"`

	// Web Server MFA
	RPID     string `toml:",omitempty"`
	RPOrigin string `toml:",omitempty"`

	// Web Server TLS
	TLSCertPath string `toml:",omitempty"`
	TLSHost     string `toml:",omitempty"`
	TLSKeyPath  string `toml:",omitempty"`
	TLSPort     uint16 `toml:",omitempty"`
	TLSRedirect bool   `toml:",omitempty"`
}

type JobPipelineConfig struct {
	DefaultHTTPLimit          int64    `toml:",omitempty"`
	DefaultHTTPTimeout        Duration `toml:",omitempty"`
	FeatureExternalInitiators bool     `toml:",omitempty"`
	MaxRunDuration            Duration `toml:",omitempty"`
	ReaperInterval            Duration `toml:",omitempty"`
	ReaperThreshold           Duration `toml:",omitempty"`
	ResultWriteQueueDepth     uint64   `toml:",omitempty"`
}

type OCR2Config struct {
	// Global defaults
	ContractConfirmations              uint     `toml:",omitempty"`
	BlockchainTimeout                  Duration `toml:",omitempty"`
	ContractPollInterval               Duration `toml:",omitempty"`
	ContractSubscribeInterval          Duration `toml:",omitempty"`
	ContractTransmitterTransmitTimeout Duration `toml:",omitempty"`
	DatabaseTimeout                    Duration `toml:",omitempty"`
	KeyBundleID                        string   `toml:",omitempty"`
	MonitoringEndpoint                 string   `toml:",omitempty"`
}

type OCRConfig struct {
	// Global defaults
	ObservationTimeout           Duration `toml:",omitempty"`
	BlockchainTimeout            Duration `toml:",omitempty"`
	ContractPollInterval         Duration `toml:",omitempty"`
	ContractSubscribeInterval    Duration `toml:",omitempty"`
	DefaultTransactionQueueDepth uint32   `toml:",omitempty"`
	// Optional
	KeyBundleID          string `toml:",omitempty"`
	MonitoringEndpoint   string `toml:",omitempty"`
	SimulateTransactions bool   `toml:",omitempty"`
	TraceLogging         bool   `toml:",omitempty"`
	TransmitterAddress   string `toml:",omitempty"`
	// DEPRECATED
	OutgoingMessageBufferSize int      `toml:",omitempty"`
	IncomingMessageBufferSize int      `toml:",omitempty"`
	DHTLookupInterval         int      `toml:",omitempty"`
	BootstrapCheckInterval    Duration `toml:",omitempty"`
	NewStreamTimeout          Duration `toml:",omitempty"`
}

type P2PConfig struct {
	// V1 and V2
	NetworkingStack           ocrnetworking.NetworkingStack `toml:",omitempty"`
	IncomingMessageBufferSize int                           `toml:",omitempty"`
	OutgoingMessageBufferSize int                           `toml:",omitempty"`
	// V1 Only
	AnnounceIP                       net.IP        `toml:",omitempty"`
	AnnouncePort                     uint16        `toml:",omitempty"`
	BootstrapCheckInterval           Duration      `toml:",omitempty"`
	BootstrapPeers                   []string      `toml:",omitempty"`
	DHTAnnouncementCounterUserPrefix uint32        `toml:",omitempty"`
	DHTLookupInterval                int           `toml:",omitempty"`
	ListenIP                         net.IP        `toml:",omitempty"`
	ListenPort                       uint16        `toml:",omitempty"`
	NewStreamTimeout                 Duration      `toml:",omitempty"`
	PeerID                           p2pkey.PeerID `toml:",omitempty"`
	PeerstoreWriteInterval           Duration      `toml:",omitempty"`
	// V2 Only
	V2AnnounceAddresses []string `toml:",omitempty"`
	V2Bootstrappers     []string `toml:",omitempty"`
	V2DeltaDial         Duration `toml:",omitempty"`
	V2DeltaReconcile    Duration `toml:",omitempty"`
	V2ListenAddresses   []string `toml:",omitempty"`
}

type KeeperConfig struct {
	CheckUpkeepGasPriceFeatureEnabled bool     `toml:",omitempty"`
	DefaultTransactionQueueDepth      uint32   `toml:",omitempty"`
	GasPriceBufferPercent             uint32   `toml:",omitempty"`
	GasTipCapBufferPercent            uint32   `toml:",omitempty"`
	BaseFeeBufferPercent              uint32   `toml:",omitempty"`
	MaximumGracePeriod                int64    `toml:",omitempty"`
	RegistryCheckGasOverhead          uint64   `toml:",omitempty"`
	RegistryPerformGasOverhead        uint64   `toml:",omitempty"`
	RegistrySyncInterval              Duration `toml:",omitempty"`
	RegistrySyncUpkeepQueueSize       uint32   `toml:",omitempty"`
	TurnLookBack                      int64    `toml:",omitempty"`
	TurnFlagEnabled                   bool     `toml:",omitempty"`
}

type AutoPprofConfig struct {
	Enabled              bool           `toml:",omitempty"`
	ProfileRoot          string         `toml:",omitempty"`
	PollInterval         Duration       `toml:",omitempty"`
	GatherDuration       Duration       `toml:",omitempty"`
	GatherTraceDuration  Duration       `toml:",omitempty"`
	MaxProfileSize       utils.FileSize `toml:",omitempty"`
	CPUProfileRate       int            `toml:",omitempty"`
	MemProfileRate       int            `toml:",omitempty"`
	BlockProfileRate     int            `toml:",omitempty"`
	MutexProfileFraction int            `toml:",omitempty"`
	MemThreshold         utils.FileSize `toml:",omitempty"`
	GoroutineThreshold   int            `toml:",omitempty"`
}

// Duration extends time.Duration with encoding/text methods.
type Duration time.Duration

// MarshalText implements the text.Marshaler interface.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

// UnmarshalText implements the text.Unmarshaler interface.
func (d *Duration) UnmarshalText(input []byte) error {
	v, err := time.ParseDuration(string(input))
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

// URL extends url.URL with encoding/text methods.
type URL url.URL

// MarshalText implements the text.Marshaler interface.
func (u *URL) MarshalText() ([]byte, error) {
	return []byte((*url.URL)(u).String()), nil
}

// UnmarshalText implements the text.Unmarshaler interface.
func (u *URL) UnmarshalText(input []byte) error {
	v, err := url.Parse(string(input))
	if err != nil {
		return err
	}
	*u = URL(*v)
	return nil
}
