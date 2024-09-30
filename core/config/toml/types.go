package toml

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/parse"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	configutils "github.com/smartcontractkit/chainlink/v2/core/utils/config"
)

var ErrUnsupported = errors.New("unsupported with config v2")

// Core holds the core configuration. See chainlink.Config for more information.
type Core struct {
	// General/misc
	AppID               uuid.UUID `toml:"-"` // random or test
	InsecureFastScrypt  *bool
	RootDir             *string
	ShutdownGracePeriod *commonconfig.Duration

	Feature          Feature          `toml:",omitempty"`
	Database         Database         `toml:",omitempty"`
	TelemetryIngress TelemetryIngress `toml:",omitempty"`
	AuditLogger      AuditLogger      `toml:",omitempty"`
	Log              Log              `toml:",omitempty"`
	WebServer        WebServer        `toml:",omitempty"`
	JobPipeline      JobPipeline      `toml:",omitempty"`
	FluxMonitor      FluxMonitor      `toml:",omitempty"`
	OCR2             OCR2             `toml:",omitempty"`
	OCR              OCR              `toml:",omitempty"`
	P2P              P2P              `toml:",omitempty"`
	Keeper           Keeper           `toml:",omitempty"`
	AutoPprof        AutoPprof        `toml:",omitempty"`
	Pyroscope        Pyroscope        `toml:",omitempty"`
	Sentry           Sentry           `toml:",omitempty"`
	Insecure         Insecure         `toml:",omitempty"`
	Tracing          Tracing          `toml:",omitempty"`
	Mercury          Mercury          `toml:",omitempty"`
	Capabilities     Capabilities     `toml:",omitempty"`
	Telemetry        Telemetry        `toml:",omitempty"`
}

// SetFrom updates c with any non-nil values from f. (currently TOML field only!)
func (c *Core) SetFrom(f *Core) {
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
	c.Mercury.setFrom(&f.Mercury)
	c.Capabilities.setFrom(&f.Capabilities)

	c.AutoPprof.setFrom(&f.AutoPprof)
	c.Pyroscope.setFrom(&f.Pyroscope)
	c.Sentry.setFrom(&f.Sentry)
	c.Insecure.setFrom(&f.Insecure)
	c.Tracing.setFrom(&f.Tracing)
	c.Telemetry.setFrom(&f.Telemetry)
}

func (c *Core) ValidateConfig() (err error) {
	_, verr := parse.HomeDir(*c.RootDir)
	if verr != nil {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "RootDir", Value: true, Msg: fmt.Sprintf("Failed to expand RootDir. Please use an explicit path: %s", verr)})
	}

	if (*c.OCR.Enabled || *c.OCR2.Enabled) && !*c.P2P.V2.Enabled {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "P2P.V2.Enabled", Value: false, Msg: "P2P required for OCR or OCR2. Please enable P2P or disable OCR/OCR2."})
	}

	if *c.Tracing.Enabled && *c.Telemetry.Enabled {
		if c.Tracing.CollectorTarget == c.Telemetry.Endpoint {
			err = multierr.Append(err, configutils.ErrInvalid{Name: "Tracing.CollectorTarget", Value: *c.Tracing.CollectorTarget, Msg: "Same as Telemetry.Endpoint. Must be different or disabled."})
		}
	}

	return err
}

type Secrets struct {
	Database   DatabaseSecrets          `toml:",omitempty"`
	Password   Passwords                `toml:",omitempty"`
	WebServer  WebServerSecrets         `toml:",omitempty"`
	Pyroscope  PyroscopeSecrets         `toml:",omitempty"`
	Prometheus PrometheusSecrets        `toml:",omitempty"`
	Mercury    MercurySecrets           `toml:",omitempty"`
	Threshold  ThresholdKeyShareSecrets `toml:",omitempty"`
}

func dbURLPasswordComplexity(err error) string {
	return fmt.Sprintf("missing or insufficiently complex password: %s. Database should be secured by a password matching the following complexity requirements: "+utils.PasswordComplexityRequirements, err)
}

type DatabaseSecrets struct {
	URL                  *models.SecretURL
	BackupURL            *models.SecretURL
	AllowSimplePasswords *bool
}

func validateDBURL(dbURI url.URL) error {
	if strings.Contains(dbURI.Redacted(), "_test") {
		return nil
	}

	// url params take priority if present, multiple params are ignored by postgres (it picks the first)
	q := dbURI.Query()
	// careful, this is a raw database password
	pw := q.Get("password")
	if pw == "" {
		// fallback to user info
		userInfo := dbURI.User
		if userInfo == nil {
			return fmt.Errorf("DB URL must be authenticated; plaintext URLs are not allowed")
		}
		var pwSet bool
		pw, pwSet = userInfo.Password()
		if !pwSet {
			return fmt.Errorf("DB URL must be authenticated; password is required")
		}
	}

	return utils.VerifyPasswordComplexity(pw)
}

func (d *DatabaseSecrets) ValidateConfig() (err error) {
	return d.validateConfig(build.Mode())
}

func (d *DatabaseSecrets) validateConfig(buildMode string) (err error) {
	if d.URL == nil || (*url.URL)(d.URL).String() == "" {
		err = multierr.Append(err, configutils.ErrEmpty{Name: "URL", Msg: "must be provided and non-empty"})
	} else if *d.AllowSimplePasswords && buildMode == build.Prod {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "AllowSimplePasswords", Value: true, Msg: "insecure configs are not allowed on secure builds"})
	} else if !*d.AllowSimplePasswords {
		if verr := validateDBURL((url.URL)(*d.URL)); verr != nil {
			err = multierr.Append(err, configutils.ErrInvalid{Name: "URL", Value: "*****", Msg: dbURLPasswordComplexity(verr)})
		}
	}
	if d.BackupURL != nil && !*d.AllowSimplePasswords {
		if verr := validateDBURL((url.URL)(*d.BackupURL)); verr != nil {
			err = multierr.Append(err, configutils.ErrInvalid{Name: "BackupURL", Value: "*****", Msg: dbURLPasswordComplexity(verr)})
		}
	}
	return err
}

func (d *DatabaseSecrets) SetFrom(f *DatabaseSecrets) (err error) {
	err = d.validateMerge(f)
	if err != nil {
		return err
	}

	if v := f.AllowSimplePasswords; v != nil {
		d.AllowSimplePasswords = v
	}
	if v := f.BackupURL; v != nil {
		d.BackupURL = v
	}
	if v := f.URL; v != nil {
		d.URL = v
	}
	return nil
}

func (d *DatabaseSecrets) validateMerge(f *DatabaseSecrets) (err error) {
	if d.AllowSimplePasswords != nil && f.AllowSimplePasswords != nil {
		err = multierr.Append(err, configutils.ErrOverride{Name: "AllowSimplePasswords"})
	}

	if d.BackupURL != nil && f.BackupURL != nil {
		err = multierr.Append(err, configutils.ErrOverride{Name: "BackupURL"})
	}

	if d.URL != nil && f.URL != nil {
		err = multierr.Append(err, configutils.ErrOverride{Name: "URL"})
	}

	return err
}

type Passwords struct {
	Keystore *models.Secret
	VRF      *models.Secret
}

func (p *Passwords) SetFrom(f *Passwords) (err error) {
	err = p.validateMerge(f)
	if err != nil {
		return err
	}

	if v := f.Keystore; v != nil {
		p.Keystore = v
	}
	if v := f.VRF; v != nil {
		p.VRF = v
	}

	return nil
}

func (p *Passwords) validateMerge(f *Passwords) (err error) {
	if p.Keystore != nil && f.Keystore != nil {
		err = multierr.Append(err, configutils.ErrOverride{Name: "Keystore"})
	}

	if p.VRF != nil && f.VRF != nil {
		err = multierr.Append(err, configutils.ErrOverride{Name: "VRF"})
	}

	return err
}

func (p *Passwords) ValidateConfig() (err error) {
	if p.Keystore == nil || *p.Keystore == "" {
		err = multierr.Append(err, configutils.ErrEmpty{Name: "Keystore", Msg: "must be provided and non-empty"})
	}
	return err
}

type PyroscopeSecrets struct {
	AuthToken *models.Secret
}

func (p *PyroscopeSecrets) SetFrom(f *PyroscopeSecrets) (err error) {
	err = p.validateMerge(f)
	if err != nil {
		return err
	}

	if v := f.AuthToken; v != nil {
		p.AuthToken = v
	}

	return nil
}

func (p *PyroscopeSecrets) validateMerge(f *PyroscopeSecrets) (err error) {
	if p.AuthToken != nil && f.AuthToken != nil {
		err = multierr.Append(err, configutils.ErrOverride{Name: "AuthToken"})
	}

	return err
}

type PrometheusSecrets struct {
	AuthToken *models.Secret
}

func (p *PrometheusSecrets) SetFrom(f *PrometheusSecrets) (err error) {
	err = p.validateMerge(f)
	if err != nil {
		return err
	}

	if v := f.AuthToken; v != nil {
		p.AuthToken = v
	}

	return nil
}

func (p *PrometheusSecrets) validateMerge(f *PrometheusSecrets) (err error) {
	if p.AuthToken != nil && f.AuthToken != nil {
		err = multierr.Append(err, configutils.ErrOverride{Name: "AuthToken"})
	}

	return err
}

type Feature struct {
	FeedsManager       *bool
	LogPoller          *bool
	UICSAKeys          *bool
	CCIP               *bool
	MultiFeedsManagers *bool
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
	if v := f2.CCIP; v != nil {
		f.CCIP = v
	}
	if v := f2.MultiFeedsManagers; v != nil {
		f.MultiFeedsManagers = v
	}
}

type Database struct {
	DefaultIdleInTxSessionTimeout *commonconfig.Duration
	DefaultLockTimeout            *commonconfig.Duration
	DefaultQueryTimeout           *commonconfig.Duration
	Dialect                       dialects.DialectName `toml:"-"`
	LogQueries                    *bool
	MaxIdleConns                  *int64
	MaxOpenConns                  *int64
	MigrateOnStartup              *bool

	Backup   DatabaseBackup   `toml:",omitempty"`
	Listener DatabaseListener `toml:",omitempty"`
	Lock     DatabaseLock     `toml:",omitempty"`
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
	MaxReconnectDuration *commonconfig.Duration
	MinReconnectInterval *commonconfig.Duration
	FallbackPollInterval *commonconfig.Duration
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
	LeaseDuration        *commonconfig.Duration
	LeaseRefreshInterval *commonconfig.Duration
}

func (l *DatabaseLock) Mode() string {
	if *l.Enabled {
		return "lease"
	}
	return "none"
}

func (l *DatabaseLock) ValidateConfig() (err error) {
	if l.LeaseRefreshInterval.Duration() > l.LeaseDuration.Duration()/2 {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "LeaseRefreshInterval", Value: l.LeaseRefreshInterval.String(),
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
	Frequency        *commonconfig.Duration
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
	BufferSize   *uint16
	MaxBatchSize *uint16
	SendInterval *commonconfig.Duration
	SendTimeout  *commonconfig.Duration
	UseBatchSend *bool
	Endpoints    []TelemetryIngressEndpoint `toml:",omitempty"`
}

type TelemetryIngressEndpoint struct {
	Network      *string
	ChainID      *string
	URL          *commonconfig.URL
	ServerPubKey *string
}

func (t *TelemetryIngress) setFrom(f *TelemetryIngress) {
	if v := f.UniConn; v != nil {
		t.UniConn = v
	}
	if v := f.Logging; v != nil {
		t.Logging = v
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
	if v := f.Endpoints; v != nil {
		t.Endpoints = v
	}
}

type AuditLogger struct {
	Enabled        *bool
	ForwardToUrl   *commonconfig.URL
	JsonWrapperKey *string
	Headers        *[]models.ServiceHeader
}

func (p *AuditLogger) SetFrom(f *AuditLogger) {
	if v := f.Enabled; v != nil {
		p.Enabled = v
	}
	if v := f.ForwardToUrl; v != nil {
		p.ForwardToUrl = v
	}
	if v := f.JsonWrapperKey; v != nil {
		p.JsonWrapperKey = v
	}
	if v := f.Headers; v != nil {
		p.Headers = v
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
	AuthenticationMethod    *string
	AllowOrigins            *string
	BridgeResponseURL       *commonconfig.URL
	BridgeCacheTTL          *commonconfig.Duration
	HTTPWriteTimeout        *commonconfig.Duration
	HTTPPort                *uint16
	SecureCookies           *bool
	SessionTimeout          *commonconfig.Duration
	SessionReaperExpiration *commonconfig.Duration
	HTTPMaxSize             *utils.FileSize
	StartTimeout            *commonconfig.Duration
	ListenIP                *net.IP

	LDAP      WebServerLDAP      `toml:",omitempty"`
	MFA       WebServerMFA       `toml:",omitempty"`
	RateLimit WebServerRateLimit `toml:",omitempty"`
	TLS       WebServerTLS       `toml:",omitempty"`
}

func (w *WebServer) setFrom(f *WebServer) {
	if v := f.AuthenticationMethod; v != nil {
		w.AuthenticationMethod = v
	}
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
	if v := f.ListenIP; v != nil {
		w.ListenIP = v
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
	if v := f.StartTimeout; v != nil {
		w.StartTimeout = v
	}
	if v := f.HTTPMaxSize; v != nil {
		w.HTTPMaxSize = v
	}

	w.LDAP.setFrom(&f.LDAP)
	w.MFA.setFrom(&f.MFA)
	w.RateLimit.setFrom(&f.RateLimit)
	w.TLS.setFrom(&f.TLS)
}

func (w *WebServer) ValidateConfig() (err error) {
	// Validate LDAP fields when authentication method is LDAPAuth
	if *w.AuthenticationMethod != string(sessions.LDAPAuth) {
		return
	}

	// Assert LDAP fields when AuthMethod set to LDAP
	if *w.LDAP.BaseDN == "" {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "LDAP.BaseDN", Msg: "LDAP BaseDN can not be empty"})
	}
	if *w.LDAP.BaseUserAttr == "" {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "LDAP.BaseUserAttr", Msg: "LDAP BaseUserAttr can not be empty"})
	}
	if *w.LDAP.UsersDN == "" {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "LDAP.UsersDN", Msg: "LDAP UsersDN can not be empty"})
	}
	if *w.LDAP.GroupsDN == "" {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "LDAP.GroupsDN", Msg: "LDAP GroupsDN can not be empty"})
	}
	if *w.LDAP.AdminUserGroupCN == "" {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "LDAP.AdminUserGroupCN", Msg: "LDAP AdminUserGroupCN can not be empty"})
	}
	if *w.LDAP.EditUserGroupCN == "" {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "LDAP.RunUserGroupCN", Msg: "LDAP ReadUserGroupCN can not be empty"})
	}
	if *w.LDAP.RunUserGroupCN == "" {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "LDAP.RunUserGroupCN", Msg: "LDAP RunUserGroupCN can not be empty"})
	}
	if *w.LDAP.ReadUserGroupCN == "" {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "LDAP.ReadUserGroupCN", Msg: "LDAP ReadUserGroupCN can not be empty"})
	}
	return err
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
	AuthenticatedPeriod   *commonconfig.Duration
	Unauthenticated       *int64
	UnauthenticatedPeriod *commonconfig.Duration
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
	ListenIP      *net.IP
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
	if v := f.ListenIP; v != nil {
		w.ListenIP = v
	}
}

type WebServerLDAP struct {
	ServerTLS                   *bool
	SessionTimeout              *commonconfig.Duration
	QueryTimeout                *commonconfig.Duration
	BaseUserAttr                *string
	BaseDN                      *string
	UsersDN                     *string
	GroupsDN                    *string
	ActiveAttribute             *string
	ActiveAttributeAllowedValue *string
	AdminUserGroupCN            *string
	EditUserGroupCN             *string
	RunUserGroupCN              *string
	ReadUserGroupCN             *string
	UserApiTokenEnabled         *bool
	UserAPITokenDuration        *commonconfig.Duration
	UpstreamSyncInterval        *commonconfig.Duration
	UpstreamSyncRateLimit       *commonconfig.Duration
}

func (w *WebServerLDAP) setFrom(f *WebServerLDAP) {
	if v := f.ServerTLS; v != nil {
		w.ServerTLS = v
	}
	if v := f.SessionTimeout; v != nil {
		w.SessionTimeout = v
	}
	if v := f.SessionTimeout; v != nil {
		w.SessionTimeout = v
	}
	if v := f.QueryTimeout; v != nil {
		w.QueryTimeout = v
	}
	if v := f.BaseUserAttr; v != nil {
		w.BaseUserAttr = v
	}
	if v := f.BaseDN; v != nil {
		w.BaseDN = v
	}
	if v := f.UsersDN; v != nil {
		w.UsersDN = v
	}
	if v := f.GroupsDN; v != nil {
		w.GroupsDN = v
	}
	if v := f.ActiveAttribute; v != nil {
		w.ActiveAttribute = v
	}
	if v := f.ActiveAttributeAllowedValue; v != nil {
		w.ActiveAttributeAllowedValue = v
	}
	if v := f.AdminUserGroupCN; v != nil {
		w.AdminUserGroupCN = v
	}
	if v := f.EditUserGroupCN; v != nil {
		w.EditUserGroupCN = v
	}
	if v := f.RunUserGroupCN; v != nil {
		w.RunUserGroupCN = v
	}
	if v := f.ReadUserGroupCN; v != nil {
		w.ReadUserGroupCN = v
	}
	if v := f.UserApiTokenEnabled; v != nil {
		w.UserApiTokenEnabled = v
	}
	if v := f.UserAPITokenDuration; v != nil {
		w.UserAPITokenDuration = v
	}
	if v := f.UpstreamSyncInterval; v != nil {
		w.UpstreamSyncInterval = v
	}
	if v := f.UpstreamSyncRateLimit; v != nil {
		w.UpstreamSyncRateLimit = v
	}
}

type WebServerLDAPSecrets struct {
	ServerAddress     *models.SecretURL
	ReadOnlyUserLogin *models.Secret
	ReadOnlyUserPass  *models.Secret
}

func (w *WebServerLDAPSecrets) setFrom(f *WebServerLDAPSecrets) {
	if v := f.ServerAddress; v != nil {
		w.ServerAddress = v
	}
	if v := f.ReadOnlyUserLogin; v != nil {
		w.ReadOnlyUserLogin = v
	}
	if v := f.ReadOnlyUserPass; v != nil {
		w.ReadOnlyUserPass = v
	}
}

type WebServerSecrets struct {
	LDAP WebServerLDAPSecrets `toml:",omitempty"`
}

func (w *WebServerSecrets) SetFrom(f *WebServerSecrets) error {
	w.LDAP.setFrom(&f.LDAP)
	return nil
}

type JobPipeline struct {
	ExternalInitiatorsEnabled *bool
	MaxRunDuration            *commonconfig.Duration
	MaxSuccessfulRuns         *uint64
	ReaperInterval            *commonconfig.Duration
	ReaperThreshold           *commonconfig.Duration
	ResultWriteQueueDepth     *uint32
	VerboseLogging            *bool

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
	if v := f.VerboseLogging; v != nil {
		j.VerboseLogging = v
	}
	j.HTTPRequest.setFrom(&f.HTTPRequest)
}

type JobPipelineHTTPRequest struct {
	DefaultTimeout *commonconfig.Duration
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
	BlockchainTimeout                  *commonconfig.Duration
	ContractPollInterval               *commonconfig.Duration
	ContractSubscribeInterval          *commonconfig.Duration
	ContractTransmitterTransmitTimeout *commonconfig.Duration
	DatabaseTimeout                    *commonconfig.Duration
	KeyBundleID                        *models.Sha256Hash
	CaptureEATelemetry                 *bool
	CaptureAutomationCustomTelemetry   *bool
	DefaultTransactionQueueDepth       *uint32
	SimulateTransactions               *bool
	TraceLogging                       *bool
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
	if v := f.CaptureAutomationCustomTelemetry; v != nil {
		o.CaptureAutomationCustomTelemetry = v
	}
	if v := f.DefaultTransactionQueueDepth; v != nil {
		o.DefaultTransactionQueueDepth = v
	}
	if v := f.SimulateTransactions; v != nil {
		o.SimulateTransactions = v
	}
	if v := f.TraceLogging; v != nil {
		o.TraceLogging = v
	}
}

type OCR struct {
	Enabled                      *bool
	ObservationTimeout           *commonconfig.Duration
	BlockchainTimeout            *commonconfig.Duration
	ContractPollInterval         *commonconfig.Duration
	ContractSubscribeInterval    *commonconfig.Duration
	DefaultTransactionQueueDepth *uint32
	// Optional
	KeyBundleID          *models.Sha256Hash
	SimulateTransactions *bool
	TransmitterAddress   *types.EIP55Address
	CaptureEATelemetry   *bool
	TraceLogging         *bool
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
	if v := f.TraceLogging; v != nil {
		o.TraceLogging = v
	}
}

type P2P struct {
	IncomingMessageBufferSize *int64
	OutgoingMessageBufferSize *int64
	PeerID                    *p2pkey.PeerID
	TraceLogging              *bool

	V2 P2PV2 `toml:",omitempty"`
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

	p.V2.setFrom(&f.V2)
}

type P2PV2 struct {
	Enabled              *bool
	AnnounceAddresses    *[]string
	DefaultBootstrappers *[]ocrcommontypes.BootstrapperLocator
	DeltaDial            *commonconfig.Duration
	DeltaReconcile       *commonconfig.Duration
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
	SyncInterval        *commonconfig.Duration
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
	PollInterval         *commonconfig.Duration
	GatherDuration       *commonconfig.Duration
	GatherTraceDuration  *commonconfig.Duration
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
	return ins.validateConfig(build.Mode())
}

func (ins *Insecure) validateConfig(buildMode string) (err error) {
	if buildMode == build.Dev {
		return
	}
	if ins.DevWebServer != nil && *ins.DevWebServer {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "DevWebServer", Value: *ins.DevWebServer, Msg: "insecure configs are not allowed on secure builds"})
	}
	// OCRDevelopmentMode is allowed on dev/test builds.
	if ins.OCRDevelopmentMode != nil && *ins.OCRDevelopmentMode && buildMode == build.Prod {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "OCRDevelopmentMode", Value: *ins.OCRDevelopmentMode, Msg: "insecure configs are not allowed on secure builds"})
	}
	if ins.InfiniteDepthQueries != nil && *ins.InfiniteDepthQueries {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "InfiniteDepthQueries", Value: *ins.InfiniteDepthQueries, Msg: "insecure configs are not allowed on secure builds"})
	}
	if ins.DisableRateLimiting != nil && *ins.DisableRateLimiting {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "DisableRateLimiting", Value: *ins.DisableRateLimiting, Msg: "insecure configs are not allowed on secure builds"})
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

type MercuryCache struct {
	LatestReportTTL      *commonconfig.Duration
	MaxStaleAge          *commonconfig.Duration
	LatestReportDeadline *commonconfig.Duration
}

func (mc *MercuryCache) setFrom(f *MercuryCache) {
	if v := f.LatestReportTTL; v != nil {
		mc.LatestReportTTL = v
	}
	if v := f.MaxStaleAge; v != nil {
		mc.MaxStaleAge = v
	}
	if v := f.LatestReportDeadline; v != nil {
		mc.LatestReportDeadline = v
	}
}

type MercuryTLS struct {
	CertFile *string
}

func (m *MercuryTLS) setFrom(f *MercuryTLS) {
	if v := f.CertFile; v != nil {
		m.CertFile = v
	}
}

func (m *MercuryTLS) ValidateConfig() (err error) {
	if *m.CertFile != "" {
		if !isValidFilePath(*m.CertFile) {
			err = multierr.Append(err, configutils.ErrInvalid{Name: "CertFile", Value: *m.CertFile, Msg: "must be a valid file path"})
		}
	}
	return
}

type MercuryTransmitter struct {
	TransmitQueueMaxSize *uint32
	TransmitTimeout      *commonconfig.Duration
}

func (m *MercuryTransmitter) setFrom(f *MercuryTransmitter) {
	if v := f.TransmitQueueMaxSize; v != nil {
		m.TransmitQueueMaxSize = v
	}
	if v := f.TransmitTimeout; v != nil {
		m.TransmitTimeout = v
	}
}

type Mercury struct {
	Cache          MercuryCache       `toml:",omitempty"`
	TLS            MercuryTLS         `toml:",omitempty"`
	Transmitter    MercuryTransmitter `toml:",omitempty"`
	VerboseLogging *bool              `toml:",omitempty"`
}

func (m *Mercury) setFrom(f *Mercury) {
	m.Cache.setFrom(&f.Cache)
	m.TLS.setFrom(&f.TLS)
	m.Transmitter.setFrom(&f.Transmitter)
	if v := f.VerboseLogging; v != nil {
		m.VerboseLogging = v
	}
}

func (m *Mercury) ValidateConfig() (err error) {
	return m.TLS.ValidateConfig()
}

type MercuryCredentials struct {
	// LegacyURL is the legacy base URL for mercury v0.2 API
	LegacyURL *models.SecretURL
	// URL is the base URL for mercury v0.3 API
	URL *models.SecretURL
	// Username is the user id for mercury credential
	Username *models.Secret
	// Password is the user secret key for mercury credential
	Password *models.Secret
}

type MercurySecrets struct {
	Credentials map[string]MercuryCredentials
}

func (m *MercurySecrets) SetFrom(f *MercurySecrets) (err error) {
	err = m.validateMerge(f)
	if err != nil {
		return err
	}

	if m.Credentials != nil && f.Credentials != nil {
		for k, v := range f.Credentials {
			m.Credentials[k] = v
		}
	} else if v := f.Credentials; v != nil {
		m.Credentials = v
	}

	return nil
}

func (m *MercurySecrets) validateMerge(f *MercurySecrets) (err error) {
	if m.Credentials != nil && f.Credentials != nil {
		for k := range f.Credentials {
			if _, exists := m.Credentials[k]; exists {
				err = multierr.Append(err, configutils.ErrOverride{Name: fmt.Sprintf("Credentials[\"%s\"]", k)})
			}
		}
	}

	return err
}

func (m *MercurySecrets) ValidateConfig() (err error) {
	urls := make(map[string]struct{}, len(m.Credentials))
	for name, creds := range m.Credentials {
		if name == "" {
			err = multierr.Append(err, configutils.ErrEmpty{Name: "Name", Msg: "must be provided and non-empty"})
		}
		if creds.URL == nil || creds.URL.URL() == nil {
			err = multierr.Append(err, configutils.ErrMissing{Name: "URL", Msg: "must be provided and non-empty"})
			continue
		}
		if creds.LegacyURL != nil && creds.LegacyURL.URL() == nil {
			err = multierr.Append(err, configutils.ErrMissing{Name: "Legacy URL", Msg: "must be a valid URL"})
			continue
		}
		s := creds.URL.URL().String()
		if _, exists := urls[s]; exists {
			err = multierr.Append(err, configutils.NewErrDuplicate("URL", s))
		}
		urls[s] = struct{}{}
	}
	return err
}

type ExternalRegistry struct {
	Address   *string
	NetworkID *string
	ChainID   *string
}

func (r *ExternalRegistry) setFrom(f *ExternalRegistry) {
	if f.Address != nil {
		r.Address = f.Address
	}

	if f.NetworkID != nil {
		r.NetworkID = f.NetworkID
	}

	if f.ChainID != nil {
		r.ChainID = f.ChainID
	}
}

type Dispatcher struct {
	SupportedVersion   *int
	ReceiverBufferSize *int
	RateLimit          DispatcherRateLimit
}

func (d *Dispatcher) setFrom(f *Dispatcher) {
	d.RateLimit.setFrom(&f.RateLimit)

	if f.ReceiverBufferSize != nil {
		d.ReceiverBufferSize = f.ReceiverBufferSize
	}

	if f.SupportedVersion != nil {
		d.SupportedVersion = f.SupportedVersion
	}
}

type DispatcherRateLimit struct {
	GlobalRPS      *float64
	GlobalBurst    *int
	PerSenderRPS   *float64
	PerSenderBurst *int
}

func (drl *DispatcherRateLimit) setFrom(f *DispatcherRateLimit) {
	if f.GlobalRPS != nil {
		drl.GlobalRPS = f.GlobalRPS
	}
	if f.GlobalBurst != nil {
		drl.GlobalBurst = f.GlobalBurst
	}
	if f.PerSenderRPS != nil {
		drl.PerSenderRPS = f.PerSenderRPS
	}
	if f.PerSenderBurst != nil {
		drl.PerSenderBurst = f.PerSenderBurst
	}
}

type GatewayConnector struct {
	ChainIDForNodeKey         *string
	NodeAddress               *string
	DonID                     *string
	Gateways                  []ConnectorGateway
	WSHandshakeTimeoutMillis  *uint32
	AuthMinChallengeLen       *int
	AuthTimestampToleranceSec *uint32
}

func (r *GatewayConnector) setFrom(f *GatewayConnector) {
	if f.ChainIDForNodeKey != nil {
		r.ChainIDForNodeKey = f.ChainIDForNodeKey
	}

	if f.NodeAddress != nil {
		r.NodeAddress = f.NodeAddress
	}

	if f.DonID != nil {
		r.DonID = f.DonID
	}

	if f.Gateways != nil {
		r.Gateways = f.Gateways
	}

	if !reflect.ValueOf(f.WSHandshakeTimeoutMillis).IsZero() {
		r.WSHandshakeTimeoutMillis = f.WSHandshakeTimeoutMillis
	}

	if f.AuthMinChallengeLen != nil {
		r.AuthMinChallengeLen = f.AuthMinChallengeLen
	}

	if f.AuthTimestampToleranceSec != nil {
		r.AuthTimestampToleranceSec = f.AuthTimestampToleranceSec
	}
}

type ConnectorGateway struct {
	ID  *string
	URL *string
}

type Capabilities struct {
	Peering          P2P              `toml:",omitempty"`
	Dispatcher       Dispatcher       `toml:",omitempty"`
	ExternalRegistry ExternalRegistry `toml:",omitempty"`
	GatewayConnector GatewayConnector `toml:",omitempty"`
}

func (c *Capabilities) setFrom(f *Capabilities) {
	c.Peering.setFrom(&f.Peering)
	c.ExternalRegistry.setFrom(&f.ExternalRegistry)
	c.Dispatcher.setFrom(&f.Dispatcher)
	c.GatewayConnector.setFrom(&f.GatewayConnector)
}

type ThresholdKeyShareSecrets struct {
	ThresholdKeyShare *models.Secret
}

func (t *ThresholdKeyShareSecrets) SetFrom(f *ThresholdKeyShareSecrets) (err error) {
	err = t.validateMerge(f)
	if err != nil {
		return err
	}

	if v := f.ThresholdKeyShare; v != nil {
		t.ThresholdKeyShare = v
	}

	return nil
}

func (t *ThresholdKeyShareSecrets) validateMerge(f *ThresholdKeyShareSecrets) (err error) {
	if t.ThresholdKeyShare != nil && f.ThresholdKeyShare != nil {
		err = multierr.Append(err, configutils.ErrOverride{Name: "ThresholdKeyShare"})
	}

	return err
}

type Tracing struct {
	Enabled         *bool
	CollectorTarget *string
	NodeID          *string
	SamplingRatio   *float64
	Mode            *string
	TLSCertPath     *string
	Attributes      map[string]string `toml:",omitempty"`
}

func (t *Tracing) setFrom(f *Tracing) {
	if v := f.Enabled; v != nil {
		t.Enabled = v
	}
	if v := f.CollectorTarget; v != nil {
		t.CollectorTarget = v
	}
	if v := f.NodeID; v != nil {
		t.NodeID = v
	}
	if v := f.Attributes; v != nil {
		t.Attributes = v
	}
	if v := f.SamplingRatio; v != nil {
		t.SamplingRatio = v
	}
	if v := f.Mode; v != nil {
		t.Mode = v
	}
	if v := f.TLSCertPath; v != nil {
		t.TLSCertPath = v
	}
}

func (t *Tracing) ValidateConfig() (err error) {
	if t.Enabled == nil || !*t.Enabled {
		return err
	}

	if t.SamplingRatio != nil {
		if *t.SamplingRatio < 0 || *t.SamplingRatio > 1 {
			err = multierr.Append(err, configutils.ErrInvalid{Name: "SamplingRatio", Value: *t.SamplingRatio, Msg: "must be between 0 and 1"})
		}
	}

	if t.Mode != nil {
		switch *t.Mode {
		case "tls":
			// TLSCertPath must be set
			if t.TLSCertPath == nil {
				err = multierr.Append(err, configutils.ErrMissing{Name: "TLSCertPath", Msg: "must be set when Tracing.Mode is tls"})
			} else {
				ok := isValidFilePath(*t.TLSCertPath)
				if !ok {
					err = multierr.Append(err, configutils.ErrInvalid{Name: "TLSCertPath", Value: *t.TLSCertPath, Msg: "must be a valid file path"})
				}
			}
		case "unencrypted":
			// no-op
		default:
			// Mode must be either "tls" or "unencrypted"
			err = multierr.Append(err, configutils.ErrInvalid{Name: "Mode", Value: *t.Mode, Msg: "must be either 'tls' or 'unencrypted'"})
		}
	}

	if t.CollectorTarget != nil && t.Mode != nil {
		switch *t.Mode {
		case "tls":
			if !isValidURI(*t.CollectorTarget) {
				err = multierr.Append(err, configutils.ErrInvalid{Name: "CollectorTarget", Value: *t.CollectorTarget, Msg: "must be a valid URI"})
			}
		case "unencrypted":
			// Unencrypted traces can not be sent to external networks
			if !isValidLocalURI(*t.CollectorTarget) {
				err = multierr.Append(err, configutils.ErrInvalid{Name: "CollectorTarget", Value: *t.CollectorTarget, Msg: "must be a valid local URI"})
			}
		default:
			// no-op
		}
	}

	return err
}

type Telemetry struct {
	Enabled            *bool
	CACertFile         *string
	Endpoint           *string
	InsecureConnection *bool
	ResourceAttributes map[string]string `toml:",omitempty"`
	TraceSampleRatio   *float64
}

func (b *Telemetry) setFrom(f *Telemetry) {
	if v := f.Enabled; v != nil {
		b.Enabled = v
	}
	if v := f.CACertFile; v != nil {
		b.CACertFile = v
	}
	if v := f.Endpoint; v != nil {
		b.Endpoint = v
	}
	if v := f.InsecureConnection; v != nil {
		b.InsecureConnection = v
	}
	if v := f.ResourceAttributes; v != nil {
		b.ResourceAttributes = v
	}
	if v := f.TraceSampleRatio; v != nil {
		b.TraceSampleRatio = v
	}
}

func (b *Telemetry) ValidateConfig() (err error) {
	if b.Enabled == nil || !*b.Enabled {
		return nil
	}
	if b.Endpoint == nil || *b.Endpoint == "" {
		err = multierr.Append(err, configutils.ErrMissing{Name: "Endpoint", Msg: "must be set when Telemetry is enabled"})
	}
	if b.InsecureConnection != nil && *b.InsecureConnection {
		if build.IsProd() {
			err = multierr.Append(err, configutils.ErrInvalid{Name: "InsecureConnection", Value: true, Msg: "cannot be used in production builds"})
		}
	} else {
		if b.CACertFile == nil || *b.CACertFile == "" {
			err = multierr.Append(err, configutils.ErrMissing{Name: "CACertFile", Msg: "must be set, unless InsecureConnection is used"})
		}
	}
	if ratio := b.TraceSampleRatio; ratio != nil && (*ratio < 0 || *ratio > 1) {
		err = multierr.Append(err, configutils.ErrInvalid{Name: "TraceSampleRatio", Value: *ratio, Msg: "must be between 0 and 1"})
	}

	return err
}

var hostnameRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`)

// Validates uri is valid external or local URI
func isValidURI(uri string) bool {
	if strings.Contains(uri, "://") {
		_, err := url.ParseRequestURI(uri)
		return err == nil
	}

	return isValidLocalURI(uri)
}

// isValidLocalURI returns true if uri is a valid local URI
// External URIs (e.g. http://) are not valid local URIs, and will return false.
func isValidLocalURI(uri string) bool {
	parts := strings.Split(uri, ":")
	if len(parts) == 2 {
		host, port := parts[0], parts[1]

		// Validating hostname
		if !isValidHostname(host) {
			return false
		}

		// Validating port
		if _, err := net.LookupPort("tcp", port); err != nil {
			return false
		}

		return true
	}
	return false
}

func isValidHostname(hostname string) bool {
	return hostnameRegex.MatchString(hostname)
}

func isValidFilePath(path string) bool {
	return len(path) > 0 && len(path) < 4096
}
