package config

import (
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

	Validate() error
	ValidateDB() error
	LogConfiguration(log LogfFn)
	SetLogLevel(lvl zapcore.Level) error
	SetLogSQL(logSQL bool)
	SetPasswords(keystore, vrf *string)

	AutoPprof
	Ethereum
	Explorer
	FeatureFlags
	Insecure
	Keystore
	OCR1Config
	OCR2Config
	P2PNetworking
	P2PV1Networking
	P2PV2Networking
	Prometheus
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
}

type DatabaseBackupMode string

var (
	DatabaseBackupModeNone DatabaseBackupMode = "none"
	DatabaseBackupModeLite DatabaseBackupMode = "lite"
	DatabaseBackupModeFull DatabaseBackupMode = "full"
)
