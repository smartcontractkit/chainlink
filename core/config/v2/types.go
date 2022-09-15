package v2

import (
	"fmt"
	"net"
	"net/url"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Core holds the core configuration. See chainlink.Config for more information.
type Core struct {
	// General/misc
	ExplorerURL         *models.URL
	InsecureFastScrypt  *bool
	RootDir             *string
	ShutdownGracePeriod *models.Duration

	Feature *Feature

	Database *Database

	TelemetryIngress *TelemetryIngress

	Log *Log

	WebServer *WebServer

	JobPipeline *JobPipeline

	FluxMonitor *FluxMonitor

	OCR2 *OCR2

	OCR *OCR

	P2P *P2P

	Keeper *Keeper

	AutoPprof *AutoPprof

	Pyroscope *Pyroscope

	Sentry *Sentry
}

type Secrets struct {
	DatabaseURL       *models.URL
	DatabaseBackupURL *models.URL

	ExplorerAccessKey *string
	ExplorerSecret    *string

	KeystorePassword *string
	VRFPassword      *string
}

func (s *Secrets) ValidateConfig() (err error) {
	if s.DatabaseURL == nil || (*url.URL)(s.DatabaseURL).String() == "" {
		err = multierr.Append(err, ErrEmpty{Name: "DatabaseURL", Msg: "must be provided and non-empty"})
	} else {
		if verr := config.ValidateDBURL((url.URL)(*s.DatabaseURL)); verr != nil {
			err = multierr.Append(err, ErrInvalid{Name: "DatabaseURL", Value: "*****", Msg: dbURLPasswordComplexity(verr)})
		}
	}
	if s.DatabaseBackupURL != nil {
		if verr := config.ValidateDBURL((url.URL)(*s.DatabaseBackupURL)); verr != nil {
			err = multierr.Append(err, ErrInvalid{Name: "DatabaseBackupURL", Value: "*****", Msg: dbURLPasswordComplexity(verr)})
		}
	}
	if s.KeystorePassword == nil || *s.KeystorePassword == "" {
		err = multierr.Append(err, ErrEmpty{Name: "KeystorePassword", Msg: "must be provided and non-empty"})
	}
	return err
}

func dbURLPasswordComplexity(err error) string {
	return fmt.Sprintf("missing or insufficiently complex password: %s. Database should be secured by a password matching the following complexity requirements: "+utils.PasswordComplexityRequirements, err)
}

func (s *Secrets) String() string {
	return "<hidden>"
}

func (s *Secrets) GoString() string {
	return "<hidden>"
}

func (s *Secrets) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}

func (s *Secrets) MarshalText() ([]byte, error) {
	return []byte("<hidden>"), nil
}

type Feature struct {
	FeedsManager *bool
	LogPoller    *bool
	UICSAKeys    *bool
}

type Database struct {
	DefaultIdleInTxSessionTimeout *models.Duration
	DefaultLockTimeout            *models.Duration
	DefaultQueryTimeout           *models.Duration
	MigrateOnStartup              *bool
	ORMMaxIdleConns               *int64
	ORMMaxOpenConns               *int64

	Backup *DatabaseBackup

	Listener *DatabaseListener

	Lock *DatabaseLock
}

type DatabaseListener struct {
	MaxReconnectDuration *models.Duration
	MinReconnectInterval *models.Duration
	FallbackPollInterval *models.Duration
}

type DatabaseLock struct {
	LeaseDuration        *models.Duration
	LeaseRefreshInterval *models.Duration
}

func (l *DatabaseLock) ValidateConfig() (err error) {
	if l.LeaseRefreshInterval.Duration() > l.LeaseDuration.Duration()/2 {
		err = multierr.Append(err, fmt.Errorf("LeaseRefreshInterval (%s) must be less than or equal to half of LeaseDuration (%s)", l.LeaseRefreshInterval.String(), l.LeaseDuration.String()))
	}
	return
}

type DatabaseBackup struct {
	Dir              *string
	Frequency        *models.Duration
	Mode             *config.DatabaseBackupMode
	OnVersionUpgrade *bool
	// DatabaseBackupURL stored in Secrets
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
	DatabaseQueries *bool
	FileDir         *string
	FileMaxSize     *utils.FileSize
	FileMaxAgeDays  *int64
	FileMaxBackups  *int64
	JSONConsole     *bool
	UnixTS          *bool
}

type WebServer struct {
	AllowOrigins            *string
	BridgeResponseURL       *models.URL
	HTTPWriteTimeout        *models.Duration
	HTTPPort                *uint16
	SecureCookies           *bool
	SessionTimeout          *models.Duration
	SessionReaperExpiration *models.Duration

	MFA *WebServerMFA

	RateLimit *WebServerRateLimit

	TLS *WebServerTLS
}

type WebServerMFA struct {
	RPID     *string
	RPOrigin *string
}

type WebServerRateLimit struct {
	Authenticated         *int64
	AuthenticatedPeriod   *models.Duration
	Unauthenticated       *int64
	UnauthenticatedPeriod *models.Duration
}

type WebServerTLS struct {
	CertPath      *string
	ForceRedirect *bool
	Host          *string
	HTTPSPort     *uint16
	KeyPath       *string
}

type JobPipeline struct {
	DefaultHTTPRequestTimeout *models.Duration
	ExternalInitiatorsEnabled *bool
	HTTPRequestMaxSize        *utils.FileSize
	MaxRunDuration            *models.Duration
	ReaperInterval            *models.Duration
	ReaperThreshold           *models.Duration
	ResultWriteQueueDepth     *uint32
}

type FluxMonitor struct {
	DefaultTransactionQueueDepth *uint32
	SimulateTransactions         *bool
}

type OCR2 struct {
	Enabled                            *bool //TODO disabled https://app.shortcut.com/chainlinklabs/story/33615/
	ContractConfirmations              *uint32
	BlockchainTimeout                  *models.Duration
	ContractPollInterval               *models.Duration
	ContractSubscribeInterval          *models.Duration
	ContractTransmitterTransmitTimeout *models.Duration
	DatabaseTimeout                    *models.Duration
	KeyBundleID                        *models.Sha256Hash
}

type OCR struct {
	Enabled                      *bool //TODO disabled https://app.shortcut.com/chainlinklabs/story/33615/
	ObservationTimeout           *models.Duration
	BlockchainTimeout            *models.Duration
	ContractPollInterval         *models.Duration
	ContractSubscribeInterval    *models.Duration
	DefaultTransactionQueueDepth *uint32
	// Optional
	KeyBundleID          *models.Sha256Hash
	SimulateTransactions *bool
	TransmitterAddress   *ethkey.EIP55Address
}

type P2P struct {
	// V1 and V2
	IncomingMessageBufferSize *int64
	OutgoingMessageBufferSize *int64
	TraceLogging              *bool

	V1 *P2PV1

	V2 *P2PV2
}

func (p *P2P) NetworkStack() ocrnetworking.NetworkingStack {
	switch {
	case p.V1 != nil && p.V2 != nil:
		return ocrnetworking.NetworkingStackV1V2
	case p.V2 != nil:
		return ocrnetworking.NetworkingStackV2
	case p.V1 != nil:
		return ocrnetworking.NetworkingStackV1
	}
	return ocrnetworking.NetworkingStack(0)
}

type P2PV1 struct {
	AnnounceIP                       *net.IP
	AnnouncePort                     *uint16
	BootstrapCheckInterval           *models.Duration
	DefaultBootstrapPeers            *[]string
	DHTAnnouncementCounterUserPrefix *uint32
	DHTLookupInterval                *int64
	ListenIP                         *net.IP
	ListenPort                       *uint16
	NewStreamTimeout                 *models.Duration
	PeerID                           *p2pkey.PeerID
	PeerstoreWriteInterval           *models.Duration
}

func (p *P2PV1) ValidateConfig() (err error) {
	if p.AnnouncePort != nil && p.AnnounceIP == nil {
		err = multierr.Append(err, ErrMissing{Name: "AnnounceIP", Msg: fmt.Sprintf("required when AnnouncePort is set: %d", *p.AnnouncePort)})
	}
	return
}

type P2PV2 struct {
	AnnounceAddresses    *[]string
	DefaultBootstrappers *[]ocrcommontypes.BootstrapperLocator
	DeltaDial            *models.Duration
	DeltaReconcile       *models.Duration
	ListenAddresses      *[]string
}

type Keeper struct {
	DefaultTransactionQueueDepth *uint32
	GasPriceBufferPercent        *uint32
	GasTipCapBufferPercent       *uint32
	BaseFeeBufferPercent         *uint32
	MaximumGracePeriod           *int64
	RegistryCheckGasOverhead     *uint32
	RegistryPerformGasOverhead   *uint32
	RegistryMaxPerformDataSize   *uint32
	RegistrySyncInterval         *models.Duration
	RegistrySyncUpkeepQueueSize  *uint32
	TurnLookBack                 *int64
	TurnFlagEnabled              *bool
	UpkeepCheckGasPriceEnabled   *bool
}

type AutoPprof struct {
	Enabled              *bool //TODO Disabled https://app.shortcut.com/chainlinklabs/story/33615/
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

type Pyroscope struct {
	AuthToken     *string
	ServerAddress *string
	Environment   *string
}

type Sentry struct {
	Debug       *bool
	DSN         *string
	Environment *string
	Release     *string
}
