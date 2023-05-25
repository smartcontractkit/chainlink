package config

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
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
	DatabaseV1
	Ethereum
	Explorer
	FeatureFlags
	FluxMonitor
	Insecure
	JobPipeline
	Keeper
	Keystore
	Logging
	OCR1Config
	OCR2Config
	P2PNetworking
	P2PV1Networking
	P2PV2Networking
	Prometheus
	Pyroscope
	Secrets
	Sentry
	TelemetryIngress
	Web
	audit.Config

	Database() Database
}
type DatabaseBackupMode string

var (
	DatabaseBackupModeNone DatabaseBackupMode = "none"
	DatabaseBackupModeLite DatabaseBackupMode = "lite"
	DatabaseBackupModeFull DatabaseBackupMode = "full"
)
