package v2

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var ErrUnsupported = errors.New("unsupported with config v2")

// Core holds the core configuration. See chainlink.Config for more information.
type Core struct {
	// General/misc
	AppID               uuid.UUID `toml:"-"` // random or test
	DevMode             bool      `toml:"-"` // from environment
	ExplorerURL         *models.URL
	InsecureFastScrypt  *bool
	RootDir             *string
	ShutdownGracePeriod *models.Duration

	Feature          Feature                 `toml:",omitempty"`
	Database         Database                `toml:",omitempty"`
	TelemetryIngress TelemetryIngress        `toml:",omitempty"`
	AuditLogger      audit.AuditLoggerConfig `toml:",omitempty"`
	Log              Log                     `toml:",omitempty"`
	WebServer        WebServer               `toml:",omitempty"`
	JobPipeline      JobPipeline             `toml:",omitempty"`
	FluxMonitor      FluxMonitor             `toml:",omitempty"`
	OCR2             OCR2                    `toml:",omitempty"`
	OCR              OCR                     `toml:",omitempty"`
	P2P              P2P                     `toml:",omitempty"`
	Keeper           Keeper                  `toml:",omitempty"`
	AutoPprof        AutoPprof               `toml:",omitempty"`
	Pyroscope        Pyroscope               `toml:",omitempty"`
	Sentry           Sentry                  `toml:",omitempty"`
	Insecure         Insecure                `toml:",omitempty"`
}

var (
	//go:embed docs/core.toml
	defaultsTOML string
	defaults     Core
)

func init() {
	if err := cfgtest.DocDefaultsOnly(strings.NewReader(defaultsTOML), &defaults, DecodeTOML); err != nil {
		log.Fatalf("Failed to initialize defaults from docs: %v", err)
	}
}

func CoreDefaults() (c Core) {
	c.SetFrom(&defaults)
	c.Database.Dialect = dialects.Postgres // not user visible - overridden for tests only
	return
}

// SetFrom updates c with any non-nil values from f. (currently TOML field only!)
func (c *Core) SetFrom(f *Core) {
	if v := f.ExplorerURL; v != nil {
		c.ExplorerURL = f.ExplorerURL
	}
	if v := f.InsecureFastScrypt; v != nil {
		c.InsecureFastScrypt = v
	}
	if v := f.RootDir; v != nil {
		c.RootDir = v
	}
	if v := f.ShutdownGracePeriod; v != nil {
		c.ShutdownGracePeriod = v
	}

	c.Feature.setFrom(&f.Feature)
	c.Database.setFrom(&f.Database)
	c.TelemetryIngress.setFrom(&f.TelemetryIngress)
	c.AuditLogger.SetFrom(&f.AuditLogger)
	c.Log.setFrom(&f.Log)

	c.WebServer.setFrom(&f.WebServer)
	c.JobPipeline.setFrom(&f.JobPipeline)

	c.FluxMonitor.setFrom(&f.FluxMonitor)
	c.OCR2.setFrom(&f.OCR2)
	c.OCR.setFrom(&f.OCR)
	c.P2P.setFrom(&f.P2P)
	c.Keeper.setFrom(&f.Keeper)

	c.AutoPprof.setFrom(&f.AutoPprof)
	c.Pyroscope.setFrom(&f.Pyroscope)
	c.Sentry.setFrom(&f.Sentry)
	c.Insecure.setFrom(&f.Insecure)
}

type Secrets struct {
	Database   DatabaseSecrets   `toml:",omitempty"`
	Explorer   ExplorerSecrets   `toml:",omitempty"`
	Password   Passwords         `toml:",omitempty"`
	Pyroscope  PyroscopeSecrets  `toml:",omitempty"`
	Prometheus PrometheusSecrets `toml:",omitempty"`
}

func dbURLPasswordComplexity(err error) string {
	return fmt.Sprintf("missing or insufficiently complex password: %s. Database should be secured by a password matching the following complexity requirements: "+utils.PasswordComplexityRequirements, err)
}

type DatabaseSecrets struct {
	URL                  *models.SecretURL
	BackupURL            *models.SecretURL
	AllowSimplePasswords bool
}

func (d *DatabaseSecrets) ValidateConfig() (err error) {
	if d.URL == nil || (*url.URL)(d.URL).String() == "" {
		err = multierr.Append(err, ErrEmpty{Name: "URL", Msg: "must be provided and non-empty"})
	} else if !d.AllowSimplePasswords {
		if verr := config.ValidateDBURL((url.URL)(*d.URL)); verr != nil {
			err = multierr.Append(err, ErrInvalid{Name: "URL", Value: "*****", Msg: dbURLPasswordComplexity(verr)})
		}
	}
	if d.BackupURL != nil && !d.AllowSimplePasswords {
		if verr := config.ValidateDBURL((url.URL)(*d.BackupURL)); verr != nil {
			err = multierr.Append(err, ErrInvalid{Name: "BackupURL", Value: "*****", Msg: dbURLPasswordComplexity(verr)})
		}
	}
	return err
}

type ExplorerSecrets struct {
	AccessKey *models.Secret
	Secret    *models.Secret
}

type Passwords struct {
	Keystore *models.Secret
	VRF      *models.Secret
}

func (p *Passwords) ValidateConfig() (err error) {
	if p.Keystore == nil || *p.Keystore == "" {
		err = multierr.Append(err, ErrEmpty{Name: "Keystore", Msg: "must be provided and non-empty"})
	}
	return err
}

type PyroscopeSecrets struct {
	AuthToken *models.Secret
}

type PrometheusSecrets struct {
	AuthToken *models.Secret
}
type Feature struct {
	FeedsManager *bool
	LogPoller    *bool
	UICSAKeys    *bool
}

func (f *Feature) setFrom(f2 *Feature) {
	if v := f2.FeedsManager; v != nil {
		f.FeedsManager = v
	}
	if v := f2.LogPoller; v != nil {
		f.LogPoller = v
	}
	if v := f2.UICSAKeys; v != nil {
		f.UICSAKeys = v
	}
}

type Database struct {
	DefaultIdleInTxSessionTimeout *models.Duration
	DefaultLockTimeout            *models.Duration
	DefaultQueryTimeout           *models.Duration
	Dialect                       dialects.DialectName `toml:"-"`
	LogQueries                    *bool
	MaxIdleConns                  *int64
	MaxOpenConns                  *int64
	MigrateOnStartup              *bool

	Backup   DatabaseBackup   `toml:",omitempty"`
	Listener DatabaseListener `toml:",omitempty"`
	Lock     DatabaseLock     `toml:",omitempty"`
}

func (d *Database) LockingMode() string {
	if *d.Lock.Enabled {
		return "lease"
	}
	return "none"
}

func (d *Database) setFrom(f *Database) {
	if v := f.DefaultIdleInTxSessionTimeout; v != nil {
		d.DefaultIdleInTxSessionTimeout = v
	}
	if v := f.DefaultLockTimeout; v != nil {
		d.DefaultLockTimeout = v
	}
	if v := f.DefaultQueryTimeout; v != nil {
		d.DefaultQueryTimeout = v
	}
	if v := f.LogQueries; v != nil {
		d.LogQueries = v
	}
	if v := f.MigrateOnStartup; v != nil {
		d.MigrateOnStartup = v
	}
	if v := f.MaxIdleConns; v != nil {
		d.MaxIdleConns = v
	}
	if v := f.MaxOpenConns; v != nil {
		d.MaxOpenConns = v
	}

	d.Backup.setFrom(&f.Backup)
	d.Listener.setFrom(&f.Listener)
	d.Lock.setFrom(&f.Lock)
}

type DatabaseListener struct {
	MaxReconnectDuration *models.Duration
	MinReconnectInterval *models.Duration
	FallbackPollInterval *models.Duration
}

func (d *DatabaseListener) setFrom(f *DatabaseListener) {
	if v := f.MaxReconnectDuration; v != nil {
		d.MaxReconnectDuration = v
	}
	if v := f.MinReconnectInterval; v != nil {
		d.MinReconnectInterval = v
	}
	if v := f.FallbackPollInterval; v != nil {
		d.FallbackPollInterval = v
	}
}

type DatabaseLock struct {
	Enabled              *bool
	LeaseDuration        *models.Duration
	LeaseRefreshInterval *models.Duration
}

func (l *DatabaseLock) ValidateConfig() (err error) {
	if l.LeaseRefreshInterval.Duration() > l.LeaseDuration.Duration()/2 {
		err = multierr.Append(err, ErrInvalid{Name: "LeaseRefreshInterval", Value: l.LeaseRefreshInterval.String(),
			Msg: fmt.Sprintf("must be less than or equal to half of LeaseDuration (%s)", l.LeaseDuration.String())})
	}
	return
}

func (l *DatabaseLock) setFrom(f *DatabaseLock) {
	if v := f.Enabled; v != nil {
		l.Enabled = v
	}
	if v := f.LeaseDuration; v != nil {
		l.LeaseDuration = v
	}
	if v := f.LeaseRefreshInterval; v != nil {
		l.LeaseRefreshInterval = v
	}
}

// DatabaseBackup
//
// Note: url is stored in Secrets.DatabaseBackupURL
type DatabaseBackup struct {
	Dir              *string
	Frequency        *models.Duration
	Mode             *config.DatabaseBackupMode
	OnVersionUpgrade *bool
}

func (d *DatabaseBackup) setFrom(f *DatabaseBackup) {
	if v := f.Dir; v != nil {
		d.Dir = v
	}
	if v := f.Frequency; v != nil {
		d.Frequency = v
	}
	if v := f.Mode; v != nil {
		d.Mode = v
	}
	if v := f.OnVersionUpgrade; v != nil {
		d.OnVersionUpgrade = v
	}
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

func (t *TelemetryIngress) setFrom(f *TelemetryIngress) {
	if v := f.UniConn; v != nil {
		t.UniConn = v
	}
	if v := f.Logging; v != nil {
		t.Logging = v
	}
	if v := f.ServerPubKey; v != nil {
		t.ServerPubKey = v
	}
	if v := f.URL; v != nil {
		t.URL = v
	}
	if v := f.BufferSize; v != nil {
		t.BufferSize = v
	}
	if v := f.MaxBatchSize; v != nil {
		t.MaxBatchSize = v
	}
	if v := f.SendInterval; v != nil {
		t.SendInterval = v
	}
	if v := f.SendTimeout; v != nil {
		t.SendTimeout = v
	}
	if v := f.UseBatchSend; v != nil {
		t.UseBatchSend = v
	}
}

// LogLevel replaces dpanic with crit/CRIT
type LogLevel zapcore.Level

func (l LogLevel) String() string {
	zl := zapcore.Level(l)
	if zl == zapcore.DPanicLevel {
		return "crit"
	}
	return zl.String()
}

func (l LogLevel) CapitalString() string {
	zl := zapcore.Level(l)
	if zl == zapcore.DPanicLevel {
		return "CRIT"
	}
	return zl.CapitalString()
}

func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	switch string(text) {
	case "crit", "CRIT":
		*l = LogLevel(zapcore.DPanicLevel)
		return nil
	}
	return (*zapcore.Level)(l).UnmarshalText(text)
}

type Log struct {
	Level       *LogLevel
	JSONConsole *bool
	UnixTS      *bool

	File LogFile `toml:",omitempty"`
}

func (l *Log) setFrom(f *Log) {
	if v := f.Level; v != nil {
		l.Level = v
	}
	if v := f.JSONConsole; v != nil {
		l.JSONConsole = v
	}
	if v := f.UnixTS; v != nil {
		l.UnixTS = v
	}
	l.File.setFrom(&f.File)
}

type LogFile struct {
	Dir        *string
	MaxSize    *utils.FileSize
	MaxAgeDays *int64
	MaxBackups *int64
}

func (l *LogFile) setFrom(f *LogFile) {
	if v := f.Dir; v != nil {
		l.Dir = v
	}
	if v := f.MaxSize; v != nil {
		l.MaxSize = v
	}
	if v := f.MaxAgeDays; v != nil {
		l.MaxAgeDays = v
	}
	if v := f.MaxBackups; v != nil {
		l.MaxBackups = v
	}
}

type WebServer struct {
	AllowOrigins            *string
	BridgeResponseURL       *models.URL
	BridgeCacheTTL          *models.Duration
	HTTPWriteTimeout        *models.Duration
	HTTPPort                *uint16
	SecureCookies           *bool
	SessionTimeout          *models.Duration
	SessionReaperExpiration *models.Duration

	MFA       WebServerMFA       `toml:",omitempty"`
	RateLimit WebServerRateLimit `toml:",omitempty"`
	TLS       WebServerTLS       `toml:",omitempty"`
}

func (w *WebServer) setFrom(f *WebServer) {
	if v := f.AllowOrigins; v != nil {
		w.AllowOrigins = v
	}
	if v := f.BridgeResponseURL; v != nil {
		w.BridgeResponseURL = v
	}
	if v := f.BridgeCacheTTL; v != nil {
		w.BridgeCacheTTL = v
	}
	if v := f.HTTPWriteTimeout; v != nil {
		w.HTTPWriteTimeout = v
	}
	if v := f.HTTPPort; v != nil {
		w.HTTPPort = v
	}
	if v := f.SecureCookies; v != nil {
		w.SecureCookies = v
	}
	if v := f.SessionTimeout; v != nil {
		w.SessionTimeout = v
	}
	if v := f.SessionReaperExpiration; v != nil {
		w.SessionReaperExpiration = v
	}

	w.MFA.setFrom(&f.MFA)
	w.RateLimit.setFrom(&f.RateLimit)
	w.TLS.setFrom(&f.TLS)
}

type WebServerMFA struct {
	RPID     *string
	RPOrigin *string
}

func (w *WebServerMFA) setFrom(f *WebServerMFA) {
	if v := f.RPID; v != nil {
		w.RPID = v
	}
	if v := f.RPOrigin; v != nil {
		w.RPOrigin = v
	}
}

type WebServerRateLimit struct {
	Authenticated         *int64
	AuthenticatedPeriod   *models.Duration
	Unauthenticated       *int64
	UnauthenticatedPeriod *models.Duration
}

func (w *WebServerRateLimit) setFrom(f *WebServerRateLimit) {
	if v := f.Authenticated; v != nil {
		w.Authenticated = v
	}
	if v := f.AuthenticatedPeriod; v != nil {
		w.AuthenticatedPeriod = v
	}
	if v := f.Unauthenticated; v != nil {
		w.Unauthenticated = v
	}
	if v := f.UnauthenticatedPeriod; v != nil {
		w.UnauthenticatedPeriod = v
	}
}

type WebServerTLS struct {
	CertPath      *string
	ForceRedirect *bool
	Host          *string
	HTTPSPort     *uint16
	KeyPath       *string
}

func (w *WebServerTLS) setFrom(f *WebServerTLS) {
	if v := f.CertPath; v != nil {
		w.CertPath = v
	}
	if v := f.ForceRedirect; v != nil {
		w.ForceRedirect = v
	}
	if v := f.Host; v != nil {
		w.Host = v
	}
	if v := f.HTTPSPort; v != nil {
		w.HTTPSPort = v
	}
	if v := f.KeyPath; v != nil {
		w.KeyPath = v
	}
}

type JobPipeline struct {
	ExternalInitiatorsEnabled *bool
	MaxRunDuration            *models.Duration
	MaxSuccessfulRuns         *uint64
	ReaperInterval            *models.Duration
	ReaperThreshold           *models.Duration
	ResultWriteQueueDepth     *uint32

	HTTPRequest JobPipelineHTTPRequest `toml:",omitempty"`
}

func (j *JobPipeline) setFrom(f *JobPipeline) {
	if v := f.ExternalInitiatorsEnabled; v != nil {
		j.ExternalInitiatorsEnabled = v
	}
	if v := f.MaxRunDuration; v != nil {
		j.MaxRunDuration = v
	}
	if v := f.MaxSuccessfulRuns; v != nil {
		j.MaxSuccessfulRuns = v
	}
	if v := f.ReaperInterval; v != nil {
		j.ReaperInterval = v
	}
	if v := f.ReaperThreshold; v != nil {
		j.ReaperThreshold = v
	}
	if v := f.ResultWriteQueueDepth; v != nil {
		j.ResultWriteQueueDepth = v
	}
	j.HTTPRequest.setFrom(&f.HTTPRequest)

}

type JobPipelineHTTPRequest struct {
	DefaultTimeout *models.Duration
	MaxSize        *utils.FileSize
}

func (j *JobPipelineHTTPRequest) setFrom(f *JobPipelineHTTPRequest) {
	if v := f.DefaultTimeout; v != nil {
		j.DefaultTimeout = v
	}
	if v := f.MaxSize; v != nil {
		j.MaxSize = v
	}
}

type FluxMonitor struct {
	DefaultTransactionQueueDepth *uint32
	SimulateTransactions         *bool
}

func (m *FluxMonitor) setFrom(f *FluxMonitor) {
	if v := f.DefaultTransactionQueueDepth; v != nil {
		m.DefaultTransactionQueueDepth = v
	}
	if v := f.SimulateTransactions; v != nil {
		m.SimulateTransactions = v
	}
}

type OCR2 struct {
	Enabled                            *bool
	ContractConfirmations              *uint32
	BlockchainTimeout                  *models.Duration
	ContractPollInterval               *models.Duration
	ContractSubscribeInterval          *models.Duration
	ContractTransmitterTransmitTimeout *models.Duration
	DatabaseTimeout                    *models.Duration
	KeyBundleID                        *models.Sha256Hash
	CaptureEATelemetry                 *bool
}

func (o *OCR2) setFrom(f *OCR2) {
	if v := f.Enabled; v != nil {
		o.Enabled = v
	}
	if v := f.ContractConfirmations; v != nil {
		o.ContractConfirmations = v
	}
	if v := f.BlockchainTimeout; v != nil {
		o.BlockchainTimeout = v
	}
	if v := f.ContractPollInterval; v != nil {
		o.ContractPollInterval = v
	}
	if v := f.ContractSubscribeInterval; v != nil {
		o.ContractSubscribeInterval = v
	}
	if v := f.ContractTransmitterTransmitTimeout; v != nil {
		o.ContractTransmitterTransmitTimeout = v
	}
	if v := f.DatabaseTimeout; v != nil {
		o.DatabaseTimeout = v
	}
	if v := f.KeyBundleID; v != nil {
		o.KeyBundleID = v
	}
	if v := f.CaptureEATelemetry; v != nil {
		o.CaptureEATelemetry = v
	}
}

type OCR struct {
	Enabled                      *bool
	ObservationTimeout           *models.Duration
	BlockchainTimeout            *models.Duration
	ContractPollInterval         *models.Duration
	ContractSubscribeInterval    *models.Duration
	DefaultTransactionQueueDepth *uint32
	// Optional
	KeyBundleID          *models.Sha256Hash
	SimulateTransactions *bool
	TransmitterAddress   *ethkey.EIP55Address
	CaptureEATelemetry   *bool
}

func (o *OCR) setFrom(f *OCR) {
	if v := f.Enabled; v != nil {
		o.Enabled = v
	}
	if v := f.ObservationTimeout; v != nil {
		o.ObservationTimeout = v
	}
	if v := f.BlockchainTimeout; v != nil {
		o.BlockchainTimeout = v
	}
	if v := f.ContractPollInterval; v != nil {
		o.ContractPollInterval = v
	}
	if v := f.ContractSubscribeInterval; v != nil {
		o.ContractSubscribeInterval = v
	}
	if v := f.DefaultTransactionQueueDepth; v != nil {
		o.DefaultTransactionQueueDepth = v
	}
	if v := f.KeyBundleID; v != nil {
		o.KeyBundleID = v
	}
	if v := f.SimulateTransactions; v != nil {
		o.SimulateTransactions = v
	}
	if v := f.TransmitterAddress; v != nil {
		o.TransmitterAddress = v
	}
	if v := f.CaptureEATelemetry; v != nil {
		o.CaptureEATelemetry = v
	}
}

type P2P struct {
	IncomingMessageBufferSize *int64
	OutgoingMessageBufferSize *int64
	PeerID                    *p2pkey.PeerID
	TraceLogging              *bool

	V1 P2PV1 `toml:",omitempty"`
	V2 P2PV2 `toml:",omitempty"`
}

func (p *P2P) NetworkStack() ocrnetworking.NetworkingStack {
	v1, v2 := *p.V1.Enabled, *p.V2.Enabled
	switch {
	case v1 && v2:
		return ocrnetworking.NetworkingStackV1V2
	case v2:
		return ocrnetworking.NetworkingStackV2
	case v1:
		return ocrnetworking.NetworkingStackV1
	}
	return ocrnetworking.NetworkingStack(0)
}

func (p *P2P) setFrom(f *P2P) {
	if v := f.IncomingMessageBufferSize; v != nil {
		p.IncomingMessageBufferSize = v
	}
	if v := f.OutgoingMessageBufferSize; v != nil {
		p.OutgoingMessageBufferSize = v
	}
	if v := f.PeerID; v != nil {
		p.PeerID = v
	}
	if v := f.TraceLogging; v != nil {
		p.TraceLogging = v
	}

	p.V1.setFrom(&f.V1)
	p.V2.setFrom(&f.V2)
}

type P2PV1 struct {
	Enabled                          *bool
	AnnounceIP                       *net.IP
	AnnouncePort                     *uint16
	BootstrapCheckInterval           *models.Duration
	DefaultBootstrapPeers            *[]string
	DHTAnnouncementCounterUserPrefix *uint32
	DHTLookupInterval                *int64
	ListenIP                         *net.IP
	ListenPort                       *uint16
	NewStreamTimeout                 *models.Duration
	PeerstoreWriteInterval           *models.Duration
}

func (p *P2PV1) ValidateConfig() (err error) {
	//TODO or empty?
	if p.AnnouncePort != nil && p.AnnounceIP == nil {
		err = multierr.Append(err, ErrMissing{Name: "AnnounceIP", Msg: fmt.Sprintf("required when AnnouncePort is set: %d", *p.AnnouncePort)})
	}
	return
}

func (p *P2PV1) setFrom(f *P2PV1) {
	if v := f.Enabled; v != nil {
		p.Enabled = v
	}
	if v := f.AnnounceIP; v != nil {
		p.AnnounceIP = v
	}
	if v := f.AnnouncePort; v != nil {
		p.AnnouncePort = v
	}
	if v := f.BootstrapCheckInterval; v != nil {
		p.BootstrapCheckInterval = v
	}
	if v := f.DefaultBootstrapPeers; v != nil {
		p.DefaultBootstrapPeers = v
	}
	if v := f.DHTAnnouncementCounterUserPrefix; v != nil {
		p.DHTAnnouncementCounterUserPrefix = v
	}
	if v := f.DHTLookupInterval; v != nil {
		p.DHTLookupInterval = v
	}
	if v := f.ListenIP; v != nil {
		p.ListenIP = v
	}
	if v := f.ListenPort; v != nil {
		p.ListenPort = v
	}
	if v := f.NewStreamTimeout; v != nil {
		p.NewStreamTimeout = v
	}
	if v := f.PeerstoreWriteInterval; v != nil {
		p.PeerstoreWriteInterval = v
	}
}

type P2PV2 struct {
	Enabled              *bool
	AnnounceAddresses    *[]string
	DefaultBootstrappers *[]ocrcommontypes.BootstrapperLocator
	DeltaDial            *models.Duration
	DeltaReconcile       *models.Duration
	ListenAddresses      *[]string
}

func (p *P2PV2) setFrom(f *P2PV2) {
	if v := f.Enabled; v != nil {
		p.Enabled = v
	}
	if v := f.AnnounceAddresses; v != nil {
		p.AnnounceAddresses = v
	}
	if v := f.DefaultBootstrappers; v != nil {
		p.DefaultBootstrappers = v
	}
	if v := f.DeltaDial; v != nil {
		p.DeltaDial = v
	}
	if v := f.DeltaReconcile; v != nil {
		p.DeltaReconcile = v
	}
	if v := f.ListenAddresses; v != nil {
		p.ListenAddresses = v
	}
}

type Keeper struct {
	DefaultTransactionQueueDepth *uint32
	GasPriceBufferPercent        *uint16
	GasTipCapBufferPercent       *uint16
	BaseFeeBufferPercent         *uint16
	MaxGracePeriod               *int64
	TurnLookBack                 *int64

	Registry KeeperRegistry `toml:",omitempty"`
}

func (k *Keeper) setFrom(f *Keeper) {
	if v := f.DefaultTransactionQueueDepth; v != nil {
		k.DefaultTransactionQueueDepth = v
	}
	if v := f.GasPriceBufferPercent; v != nil {
		k.GasPriceBufferPercent = v
	}
	if v := f.GasTipCapBufferPercent; v != nil {
		k.GasTipCapBufferPercent = v
	}
	if v := f.BaseFeeBufferPercent; v != nil {
		k.BaseFeeBufferPercent = v
	}
	if v := f.MaxGracePeriod; v != nil {
		k.MaxGracePeriod = v
	}
	if v := f.TurnLookBack; v != nil {
		k.TurnLookBack = v
	}

	k.Registry.setFrom(&f.Registry)

}

type KeeperRegistry struct {
	CheckGasOverhead    *uint32
	PerformGasOverhead  *uint32
	MaxPerformDataSize  *uint32
	SyncInterval        *models.Duration
	SyncUpkeepQueueSize *uint32
}

func (k *KeeperRegistry) setFrom(f *KeeperRegistry) {
	if v := f.CheckGasOverhead; v != nil {
		k.CheckGasOverhead = v
	}
	if v := f.PerformGasOverhead; v != nil {
		k.PerformGasOverhead = v
	}
	if v := f.MaxPerformDataSize; v != nil {
		k.MaxPerformDataSize = v
	}
	if v := f.SyncInterval; v != nil {
		k.SyncInterval = v
	}
	if v := f.SyncUpkeepQueueSize; v != nil {
		k.SyncUpkeepQueueSize = v
	}
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

func (p *AutoPprof) setFrom(f *AutoPprof) {
	if v := f.Enabled; v != nil {
		p.Enabled = v
	}
	if v := f.ProfileRoot; v != nil {
		p.ProfileRoot = v
	}
	if v := f.PollInterval; v != nil {
		p.PollInterval = v
	}
	if v := f.GatherDuration; v != nil {
		p.GatherDuration = v
	}
	if v := f.GatherTraceDuration; v != nil {
		p.GatherTraceDuration = v
	}
	if v := f.MaxProfileSize; v != nil {
		p.MaxProfileSize = v
	}
	if v := f.CPUProfileRate; v != nil {
		p.CPUProfileRate = v
	}
	if v := f.MemProfileRate; v != nil {
		p.MemProfileRate = v
	}
	if v := f.BlockProfileRate; v != nil {
		p.BlockProfileRate = v
	}
	if v := f.MutexProfileFraction; v != nil {
		p.MutexProfileFraction = v
	}
	if v := f.MemThreshold; v != nil {
		p.MemThreshold = v
	}
	if v := f.GoroutineThreshold; v != nil {
		p.GoroutineThreshold = v
	}
}

type Pyroscope struct {
	ServerAddress *string
	Environment   *string
}

func (p *Pyroscope) setFrom(f *Pyroscope) {
	if v := f.ServerAddress; v != nil {
		p.ServerAddress = v
	}
	if v := f.Environment; v != nil {
		p.Environment = v
	}
}

type Sentry struct {
	Debug       *bool
	DSN         *string
	Environment *string
	Release     *string
}

func (s *Sentry) setFrom(f *Sentry) {
	if v := f.Debug; v != nil {
		s.Debug = f.Debug
	}
	if v := f.DSN; v != nil {
		s.DSN = f.DSN
	}
	if v := f.Environment; v != nil {
		s.Environment = f.Environment
	}
	if v := f.Release; v != nil {
		s.Release = f.Release
	}
}

type Insecure struct {
	DevWebServer         *bool
	OCRDevelopmentMode   *bool
	InfiniteDepthQueries *bool
	DisableRateLimiting  *bool
}

func (ins *Insecure) ValidateConfig() (err error) {
	if build.Dev {
		return
	}
	if ins.DevWebServer != nil && *ins.DevWebServer {
		err = multierr.Append(err, ErrInvalid{Name: "DevWebServer", Value: *ins.DevWebServer, Msg: "insecure configs are not allowed on secure builds"})
	}
	if ins.OCRDevelopmentMode != nil && *ins.OCRDevelopmentMode {
		err = multierr.Append(err, ErrInvalid{Name: "OCRDevelopmentMode", Value: *ins.OCRDevelopmentMode, Msg: "insecure configs are not allowed on secure builds"})
	}
	if ins.InfiniteDepthQueries != nil && *ins.InfiniteDepthQueries {
		err = multierr.Append(err, ErrInvalid{Name: "InfiniteDepthQueries", Value: *ins.InfiniteDepthQueries, Msg: "insecure configs are not allowed on secure builds"})
	}
	if ins.DisableRateLimiting != nil && *ins.DisableRateLimiting {
		err = multierr.Append(err, ErrInvalid{Name: "DisableRateLimiting", Value: *ins.DisableRateLimiting, Msg: "insecure configs are not allowed on secure builds"})
	}
	return err
}

func (ins *Insecure) setFrom(f *Insecure) {
	if v := f.DevWebServer; v != nil {
		ins.DevWebServer = f.DevWebServer
	}
	if v := f.InfiniteDepthQueries; v != nil {
		ins.InfiniteDepthQueries = f.InfiniteDepthQueries
	}
	if v := f.DisableRateLimiting; v != nil {
		ins.DisableRateLimiting = f.DisableRateLimiting
	}
	if v := f.OCRDevelopmentMode; v != nil {
		ins.OCRDevelopmentMode = f.OCRDevelopmentMode
	}
}
