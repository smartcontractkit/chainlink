package v2

import (
	"net"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//TODO doc
type Core struct {
	// General/misc
	Dev                 *bool
	ExplorerURL         *models.URL
	InsecureFastScrypt  *bool
	ReaperExpiration    *models.Duration
	Root                *string
	ShutdownGracePeriod *models.Duration

	Database *Database

	TelemetryIngress *TelemetryIngress

	Log *Log

	WebServer *WebServer

	// Feeds manager
	FeatureFeedsManager *bool
	FeatureUICSAKeys    *bool

	// LogPoller
	FeatureLogPoller *bool

	// Job Pipeline and tasks
	JobPipeline *JobPipeline

	// Flux Monitor
	FMDefaultTransactionQueueDepth *uint32
	FMSimulateTransactions         *bool

	// OCR V2
	FeatureOffchainReporting2 *bool
	OCR2                      *OCR2

	// OCR V1
	FeatureOffchainReporting *bool
	OCR                      *OCR

	// P2P Networking
	P2P *P2P

	// Keeper
	Keeper *Keeper

	// Debugging
	AutoPprof *AutoPprof
}

type Secrets struct {
	DatabaseURL       *models.URL
	ExplorerAccessKey string `toml:",omitempty"`
	ExplorerSecret    string `toml:",omitempty"`
	//TODO more?
}

type Database struct {
	ListenerMaxReconnectDuration  *models.Duration
	ListenerMinReconnectInterval  *models.Duration
	Migrate                       *bool
	ORMMaxIdleConns               *int64
	ORMMaxOpenConns               *int64
	TriggerFallbackDBPollInterval *models.Duration

	// Database Global Lock
	Lock *DatabaseLock

	// Database Autobackups
	Backup *DatabaseBackup
}

type DatabaseLock struct {
	Mode                  *string
	AdvisoryCheckInterval *models.Duration
	AdvisoryID            *int64
	LeaseDuration         *models.Duration
	LeaseRefreshInterval  *models.Duration
}

type DatabaseBackup struct {
	Dir              *string
	Frequency        *models.Duration
	Mode             *config.DatabaseBackupMode
	OnVersionUpgrade *bool
	URL              *models.URL
}

type TelemetryIngress struct {
	UniConn      *bool
	Logging      *bool
	ServerPubKey *string
	URL          *models.URL
	BufferSize   *uint16
	MaxBatchSize *uint16
	SendInterval *models.Duration
	SendTimeout  *models.Duration
	UseBatchSend *bool
}

type Log struct {
	JSONConsole    *bool
	FileDir        *string
	Level          *zapcore.Level //TODO is this actually an exceptional case to leave as env var?
	SQL            *bool
	FileMaxSize    *utils.FileSize
	FileMaxAgeDays *int64
	FileMaxBackups *int64
	UnixTS         *bool
}

type WebServer struct {
	// Web Server
	AllowOrigins                   *string
	AuthenticatedRateLimit         *int64
	AuthenticatedRateLimitPeriod   *models.Duration
	BridgeResponseURL              *models.URL
	HTTPWriteTimeout               *models.Duration
	Port                           *uint16
	SecureCookies                  *bool
	SessionTimeout                 *models.Duration
	UnAuthenticatedRateLimit       *int64
	UnAuthenticatedRateLimitPeriod *models.Duration

	MFA *WebServerMFA

	TLS *WebServerTLS
}

type WebServerMFA struct {
	RPID     *string
	RPOrigin *string
}

type WebServerTLS struct {
	CertPath *string
	Host     *string
	KeyPath  *string
	Port     *uint16
	Redirect *bool
}

type JobPipeline struct {
	DefaultHTTPLimit          *int64
	DefaultHTTPTimeout        *models.Duration
	FeatureExternalInitiators *bool
	MaxRunDuration            *models.Duration
	ReaperInterval            *models.Duration
	ReaperThreshold           *models.Duration
	ResultWriteQueueDepth     *uint32
}

type OCR2 struct {
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

type OCR struct {
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

type P2P struct {
	// V1 and V2
	IncomingMessageBufferSize *int64
	OutgoingMessageBufferSize *int64

	// V1 Only
	V1 *P2PV1

	// V2 Only
	V2 *P2PV2
}

type P2PV1 struct {
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

type P2PV2 struct {
	AnnounceAddresses *[]string
	Bootstrappers     *[]string
	DeltaDial         *models.Duration
	DeltaReconcile    *models.Duration
	ListenAddresses   *[]string
}

type Keeper struct {
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

type AutoPprof struct {
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
