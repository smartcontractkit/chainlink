package config

import (
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

// nolint
var (
	ErrEnvUnset = errors.New("env var unset")
)

type LogfFn func(string, ...any)

type AppConfig interface {
	AppID() uuid.UUID
	RootDir() string
	ShutdownGracePeriod() time.Duration
	InsecureFastScrypt() bool
	DefaultChainID() *big.Int

	Validate() error
	ValidateDB() error
	LogConfiguration(log LogfFn)
	SetLogLevel(lvl zapcore.Level) error
	SetLogSQL(logSQL bool)
	SetPasswords(keystore, vrf *string)

	FeatureFlags
	OCR2Config
	Pyroscope
	Secrets

	Database() Database
	AuditLogger() AuditLogger
	Keeper() Keeper
	TelemetryIngress() TelemetryIngress
	Sentry() Sentry
	JobPipeline() JobPipeline
	Log() Log
	FluxMonitor() FluxMonitor
	WebServer() WebServer
	AutoPprof() AutoPprof
	Insecure() Insecure
	Explorer() Explorer
	Password() Password
	Prometheus() Prometheus
	P2P() P2P
	OCR() OCR
}

type DatabaseBackupMode string

var (
	DatabaseBackupModeNone DatabaseBackupMode = "none"
	DatabaseBackupModeLite DatabaseBackupMode = "lite"
	DatabaseBackupModeFull DatabaseBackupMode = "full"
)
