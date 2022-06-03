package toml

import (
	"net"
	"net/url"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//TODO doc
type CoreConfig struct {
	// General/misc
	Dev                 *bool
	ExplorerURL         *URL
	InsecureFastScrypt  *bool
	ReaperExpiration    *models.Duration
	Root                *string
	ShutdownGracePeriod *models.Duration

	Database *DatabaseConfig

	TelemetryIngress *TelemetryIngressConfig

	Log *LogConfig

	WebServer *WebServerConfig

	// Feeds manager
	FeatureFeedsManager *bool
	FeatureUICSAKeys    *bool

	// LogPoller
	FeatureLogPoller *bool

	// Job Pipeline and tasks
	JobPipeline *JobPipelineConfig

	// Flux Monitor
	FMDefaultTransactionQueueDepth *uint32
	FMSimulateTransactions         *bool

	// OCR V2
	FeatureOffchainReporting2 *bool
	OCR2                      *OCR2Config

	// OCR V1
	FeatureOffchainReporting *bool
	OCR                      *OCRConfig

	// P2P Networking
	P2P *P2PConfig

	// Keeper
	Keeper *KeeperConfig

	// Debugging
	AutoPprof *AutoPprofConfig
}

type SecretsConfig struct {
	DatabaseURL       *URL
	ExplorerAccessKey string `toml:",omitempty"`
	ExplorerSecret    string `toml:",omitempty"`
	//TODO more?
}

type DatabaseConfig struct {
	ListenerMaxReconnectDuration  *models.Duration
	ListenerMinReconnectInterval  *models.Duration
	Migrate                       *bool
	ORMMaxIdleConns               *int64
	ORMMaxOpenConns               *int64
	TriggerFallbackDBPollInterval *models.Duration

	// Database Global Lock
	Lock *DatabaseLockConfig

	// Database Autobackups
	Backup *DatabaseBackupConfig
}

type DatabaseLockConfig struct {
	Mode                  *string
	AdvisoryCheckInterval *models.Duration
	AdvisoryID            *int64
	LeaseDuration         *models.Duration
	LeaseRefreshInterval  *models.Duration
}

type DatabaseBackupConfig struct {
	Dir              *string
	Frequency        *models.Duration
	Mode             *config.DatabaseBackupMode
	OnVersionUpgrade *bool
	URL              *URL
}

type TelemetryIngressConfig struct {
	UniConn      *bool
	Logging      *bool
	ServerPubKey *string
	URL          *URL
	BufferSize   *uint16
	MaxBatchSize *uint16
	SendInterval *models.Duration
	SendTimeout  *models.Duration
	UseBatchSend *bool
}

type LogConfig struct {
	JSONConsole    *bool
	FileDir        *string
	Level          *zapcore.Level //TODO is this actually an exceptional case to leave as env var?
	SQL            *bool
	FileMaxSize    *utils.FileSize
	FileMaxAgeDays *int64
	FileMaxBackups *int64
	UnixTS         *bool
}

type WebServerConfig struct {
	// Web Server
	AllowOrigins                   *string
	AuthenticatedRateLimit         *int64
	AuthenticatedRateLimitPeriod   *models.Duration
	BridgeResponseURL              *URL
	HTTPWriteTimeout               *models.Duration
	Port                           *uint16
	SecureCookies                  *bool
	SessionTimeout                 *models.Duration
	UnAuthenticatedRateLimit       *int64
	UnAuthenticatedRateLimitPeriod *models.Duration

	MFA *WebserverMFAConfig

	TLS *WebserverTLSConfig
}

type WebserverMFAConfig struct {
	RPID     *string
	RPOrigin *string
}

type WebserverTLSConfig struct {
	CertPath *string
	Host     *string
	KeyPath  *string
	Port     *uint16
	Redirect *bool
}

type JobPipelineConfig struct {
	DefaultHTTPLimit          *int64
	DefaultHTTPTimeout        *models.Duration
	FeatureExternalInitiators *bool
	MaxRunDuration            *models.Duration
	ReaperInterval            *models.Duration
	ReaperThreshold           *models.Duration
	ResultWriteQueueDepth     *uint32
}

type OCR2Config struct {
	// Global defaults
	ContractConfirmations              *uint32
	BlockchainTimeout                  *models.Duration
	ContractPollInterval               *models.Duration
	ContractSubscribeInterval          *models.Duration
	ContractTransmitterTransmitTimeout *models.Duration
	DatabaseTimeout                    *models.Duration
	KeyBundleID                        *models.Sha256Hash
	MonitoringEndpoint                 *string
}

type OCRConfig struct {
	// Global defaults
	ObservationTimeout           *models.Duration
	BlockchainTimeout            *models.Duration
	ContractPollInterval         *models.Duration
	ContractSubscribeInterval    *models.Duration
	DefaultTransactionQueueDepth *uint32
	// Optional
	KeyBundleID          *models.Sha256Hash
	MonitoringEndpoint   *string
	SimulateTransactions *bool
	TraceLogging         *bool
	TransmitterAddress   *ethkey.EIP55Address
}

type P2PConfig struct {
	// V1 and V2
	IncomingMessageBufferSize *int64
	OutgoingMessageBufferSize *int64

	// V1 Only
	V1 *P2PV1Config

	// V2 Only
	V2 *P2PV2Config
}

type P2PV1Config struct {
	AnnounceIP                       *net.IP
	AnnouncePort                     *uint16
	BootstrapCheckInterval           *models.Duration
	BootstrapPeers                   *[]string
	DHTAnnouncementCounterUserPrefix *uint32
	DHTLookupInterval                *int64
	ListenIP                         *net.IP
	ListenPort                       *uint16
	NewStreamTimeout                 *models.Duration
	PeerID                           *p2pkey.PeerID
	PeerstoreWriteInterval           *models.Duration
}

type P2PV2Config struct {
	AnnounceAddresses *[]string
	Bootstrappers     *[]string
	DeltaDial         *models.Duration
	DeltaReconcile    *models.Duration
	ListenAddresses   *[]string
}

type KeeperConfig struct {
	CheckUpkeepGasPriceFeatureEnabled *bool
	DefaultTransactionQueueDepth      *uint32
	GasPriceBufferPercent             *uint32
	GasTipCapBufferPercent            *uint32
	BaseFeeBufferPercent              *uint32
	MaximumGracePeriod                *int64
	RegistryCheckGasOverhead          *utils.Big
	RegistryPerformGasOverhead        *utils.Big
	RegistrySyncInterval              *models.Duration
	RegistrySyncUpkeepQueueSize       *uint32
	TurnLookBack                      *int64
	TurnFlagEnabled                   *bool
}

type AutoPprofConfig struct {
	Enabled              *bool
	ProfileRoot          *string
	PollInterval         *models.Duration
	GatherDuration       *models.Duration
	GatherTraceDuration  *models.Duration
	MaxProfileSize       *utils.FileSize
	CPUProfileRate       *int64 // runtime.SetCPUProfileRate
	MemProfileRate       *int64 // runtime.MemProfileRate
	BlockProfileRate     *int64 // runtime.SetBlockProfileRate
	MutexProfileFraction *int64 // runtime.SetMutexProfileFraction
	MemThreshold         *utils.FileSize
	GoroutineThreshold   *int64
}

// URL extends url.URL to implement encoding.TextMarshaler.
type URL url.URL

func (u *URL) MarshalText() ([]byte, error) {
	return []byte((*url.URL)(u).String()), nil
}

func (u *URL) UnmarshalText(input []byte) error {
	v, err := url.Parse(string(input))
	if err != nil {
		return err
	}
	*u = URL(*v)
	return nil
}
