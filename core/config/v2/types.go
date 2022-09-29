package v2

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"go.uber.org/multierr"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var ErrUnsupported = errors.New("unsupported with config v2")

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

var (
	//go:embed docs/core.toml
	defaultsTOML string
	defaults     Core
)

func init() {
	if err := cfgtest.DocDefaultsOnly(strings.NewReader(defaultsTOML), &defaults); err != nil {
		log.Fatalf("Failed to initialize defaults from docs: %v", err)
	}
}

func CoreDefaults() (c Core) {
	c.SetFrom(&defaults)
	return
}

// SetFrom updates c with any non-nil values from f.
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

	if f.Feature != nil {
		if c.Feature == nil {
			c.Feature = &Feature{}
		}
		c.Feature.setFrom(f.Feature)
	}

	if f.Database != nil {
		if c.Database == nil {
			c.Database = &Database{}
		}
		c.Database.setFrom(f.Database)

	}

	if f.TelemetryIngress != nil {
		if c.TelemetryIngress == nil {
			c.TelemetryIngress = &TelemetryIngress{}
		}
		c.TelemetryIngress.setFrom(f.TelemetryIngress)
	}

	if f.Log != nil {
		if c.Log == nil {
			c.Log = &Log{}
		}
		c.Log.setFrom(f.Log)
	}

	if f.WebServer != nil {
		if c.WebServer == nil {
			c.WebServer = &WebServer{}
		}
		c.WebServer.setFrom(f.WebServer)
	}

	if f.JobPipeline != nil {
		if c.JobPipeline == nil {
			c.JobPipeline = &JobPipeline{}
		}
		c.JobPipeline.setFrom(f.JobPipeline)
	}

	if f.FluxMonitor != nil {
		if c.FluxMonitor == nil {
			c.FluxMonitor = &FluxMonitor{}
		}
		c.FluxMonitor.setFrom(f.FluxMonitor)
	}

	if f.OCR2 != nil {
		if c.OCR2 == nil {
			c.OCR2 = &OCR2{}
		}
		c.OCR2.setFrom(f.OCR2)
	}

	if f.OCR != nil {
		if c.OCR == nil {
			c.OCR = &OCR{}
		}
		c.OCR.setFrom(f.OCR)
	}

	if f.P2P != nil {
		if c.P2P == nil {
			c.P2P = &P2P{}
		}
		c.P2P.setFrom(f.P2P)
	}

	if f.Keeper != nil {
		if c.Keeper == nil {
			c.Keeper = &Keeper{}
		}
		c.Keeper.setFrom(f.Keeper)
	}

	if f.AutoPprof != nil {
		if c.AutoPprof == nil {
			c.AutoPprof = &AutoPprof{}
		}
		c.AutoPprof.setFrom(f.AutoPprof)
	}

	if f.Pyroscope != nil {
		if c.Pyroscope == nil {
			c.Pyroscope = &Pyroscope{}
		}
		c.Pyroscope.setFrom(f.Pyroscope)
	}

	if f.Sentry != nil {
		if c.Sentry == nil {
			c.Sentry = &Sentry{}
		}
		c.Sentry.setFrom(f.Sentry)
	}
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
	MigrateOnStartup              *bool
	ORMMaxIdleConns               *int64
	ORMMaxOpenConns               *int64

	Backup *DatabaseBackup

	Listener *DatabaseListener

	Lock *DatabaseLock
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
	if v := f.MigrateOnStartup; v != nil {
		d.MigrateOnStartup = v
	}
	if v := f.ORMMaxIdleConns; v != nil {
		d.ORMMaxIdleConns = v
	}
	if v := f.ORMMaxOpenConns; v != nil {
		d.ORMMaxOpenConns = v
	}

	if f.Backup != nil {
		if d.Backup == nil {
			d.Backup = &DatabaseBackup{}
		}
		d.Backup.setFrom(f.Backup)
	}
	if f.Listener != nil {
		if d.Listener == nil {
			d.Listener = &DatabaseListener{}
		}
		d.Listener.setFrom(f.Listener)
	}

	if f.Lock != nil {
		if d.Lock == nil {
			d.Lock = &DatabaseLock{}
		}
		d.Lock.setFrom(f.Lock)
	}
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

type Log struct {
	DatabaseQueries *bool
	FileDir         *string
	FileMaxSize     *utils.FileSize
	FileMaxAgeDays  *int64
	FileMaxBackups  *int64
	JSONConsole     *bool
	UnixTS          *bool
}

func (l *Log) setFrom(f *Log) {
	if v := f.DatabaseQueries; v != nil {
		l.DatabaseQueries = v
	}
	if v := f.FileDir; v != nil {
		l.FileDir = v
	}
	if v := f.FileMaxSize; v != nil {
		l.FileMaxSize = v
	}
	if v := f.FileMaxAgeDays; v != nil {
		l.FileMaxAgeDays = v
	}
	if v := f.FileMaxBackups; v != nil {
		l.FileMaxBackups = v
	}
	if v := f.JSONConsole; v != nil {
		l.JSONConsole = v
	}
	if v := f.UnixTS; v != nil {
		l.UnixTS = v
	}
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

func (w *WebServer) setFrom(f *WebServer) {
	if v := f.AllowOrigins; v != nil {
		w.AllowOrigins = v
	}
	if v := f.BridgeResponseURL; v != nil {
		w.BridgeResponseURL = v
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
	if f.MFA != nil {
		if w.MFA == nil {
			w.MFA = &WebServerMFA{}
		}
		w.MFA.setFrom(f.MFA)
	}
	if f.RateLimit != nil {
		if w.RateLimit == nil {
			w.RateLimit = &WebServerRateLimit{}
		}
		w.RateLimit.setFrom(f.RateLimit)
	}
	if f.TLS != nil {
		if w.TLS == nil {
			w.TLS = &WebServerTLS{}
		}
		w.TLS.setFrom(f.TLS)
	}
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
	DefaultHTTPRequestTimeout *models.Duration //TODO HTTPRequestTimeout/HTTPRequest.Timeout https://app.shortcut.com/chainlinklabs/story/54384/standardize-toml-field-names
	ExternalInitiatorsEnabled *bool
	HTTPRequestMaxSize        *utils.FileSize
	MaxRunDuration            *models.Duration
	ReaperInterval            *models.Duration
	ReaperThreshold           *models.Duration
	ResultWriteQueueDepth     *uint32
}

func (j *JobPipeline) setFrom(f *JobPipeline) {
	if v := f.DefaultHTTPRequestTimeout; v != nil {
		j.DefaultHTTPRequestTimeout = v
	}
	if v := f.ExternalInitiatorsEnabled; v != nil {
		j.ExternalInitiatorsEnabled = v
	}
	if v := f.HTTPRequestMaxSize; v != nil {
		j.HTTPRequestMaxSize = v
	}
	if v := f.MaxRunDuration; v != nil {
		j.MaxRunDuration = v
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
	if v := f.TraceLogging; v != nil {
		p.TraceLogging = v
	}
	if f.V1 != nil {
		if p.V1 == nil {
			p.V1 = &P2PV1{}
		}
		p.V1.setFrom(f.V1)
	}
	if f.V2 != nil {
		if p.V2 == nil {
			p.V2 = &P2PV2{}
		}
		p.V2.setFrom(f.V2)
	}
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
	PeerID                           *p2pkey.PeerID
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
	if v := f.PeerID; v != nil {
		p.PeerID = v
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
	if v := f.MaximumGracePeriod; v != nil {
		k.MaximumGracePeriod = v
	}
	if v := f.RegistryCheckGasOverhead; v != nil {
		k.RegistryCheckGasOverhead = v
	}
	if v := f.RegistryPerformGasOverhead; v != nil {
		k.RegistryPerformGasOverhead = v
	}
	if v := f.RegistryMaxPerformDataSize; v != nil {
		k.RegistryMaxPerformDataSize = v
	}
	if v := f.RegistrySyncInterval; v != nil {
		k.RegistrySyncInterval = v
	}
	if v := f.RegistrySyncUpkeepQueueSize; v != nil {
		k.RegistrySyncUpkeepQueueSize = v
	}
	if v := f.TurnLookBack; v != nil {
		k.TurnLookBack = v
	}
	if v := f.TurnFlagEnabled; v != nil {
		k.TurnFlagEnabled = v
	}
	if v := f.UpkeepCheckGasPriceEnabled; v != nil {
		k.UpkeepCheckGasPriceEnabled = v
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
	//TODO enabled?
	AuthToken     *string //TODO move to secrets? https://app.shortcut.com/chainlinklabs/story/54383/document-secrets-toml
	ServerAddress *string
	Environment   *string
}

func (p *Pyroscope) setFrom(f *Pyroscope) {
	if v := f.AuthToken; v != nil {
		p.AuthToken = v
	}
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
