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
	EVMEnabled() bool
	EVMRPCEnabled() bool
	CosmosEnabled() bool
	SolanaEnabled() bool
	StarkNetEnabled() bool

	Validate() error
	ValidateDB() error
	LogConfiguration(log LogfFn)
	SetLogLevel(lvl zapcore.Level) error
	SetLogSQL(logSQL bool)
	SetPasswords(keystore, vrf *string)

	AuditLogger() AuditLogger
	AutoPprof() AutoPprof
	Database() Database
	Explorer() Explorer
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
}

type DatabaseBackupMode string

var (
	DatabaseBackupModeNone DatabaseBackupMode = "none"
	DatabaseBackupModeLite DatabaseBackupMode = "lite"
	DatabaseBackupModeFull DatabaseBackupMode = "full"
)
