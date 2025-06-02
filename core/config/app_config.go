package config

import (
	"time"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

var (
	ErrEnvUnset = pkgerrors.New("env var unset")
)

type LogfFn func(string, ...any)

type AppConfig interface {
	AppID() uuid.UUID
	RootDir() string
	ShutdownGracePeriod() time.Duration
	InsecureFastScrypt() bool
	InsecurePPROFHeap() bool
	EVMEnabled() bool
	CosmosEnabled() bool
	SolanaEnabled() bool
	StarkNetEnabled() bool
	AptosEnabled() bool
	TronEnabled() bool
	TONEnabled() bool

	Validate() error
	ValidateDB() error
	LogConfiguration(log, warn LogfFn)
	SetLogLevel(lvl zapcore.Level) error
	SetLogSQL(logSQL bool)
	SetPasswords(keystore, vrf *string)

	AuditLogger() AuditLogger
	AutoPprof() AutoPprof
	Capabilities() Capabilities
	Workflows() Workflows
	Database() Database
	Feature() Feature
	FluxMonitor() FluxMonitor
	Insecure() Insecure
	JobPipeline() JobPipeline
	Keeper() Keeper
	Log() Log
	Mercury() Mercury
	OCR() OCR
	OCR2() OCR2
	P2P() P2P
	Password() Password
	Prometheus() Prometheus
	Pyroscope() Pyroscope
	Sentry() Sentry
	TelemetryIngress() TelemetryIngress
	Threshold() Threshold
	WebServer() WebServer
	Tracing() Tracing
	Telemetry() Telemetry
	CRE() CRE
	Billing() Billing
}

type DatabaseBackupMode string

var (
	DatabaseBackupModeNone DatabaseBackupMode = "none"
	DatabaseBackupModeLite DatabaseBackupMode = "lite"
	DatabaseBackupModeFull DatabaseBackupMode = "full"
)
