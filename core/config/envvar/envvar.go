package envvar

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/config/parse"
)

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
// nolint
var (
	AdvisoryLockID                    = NewInt64("AdvisoryLockID")
	AuthenticatedRateLimitPeriod      = NewDuration("AuthenticatedRateLimitPeriod")
	AutoPprofPollInterval             = NewDuration("AutoPprofPollInterval")
	AutoPprofGatherDuration           = NewDuration("AutoPprofGatherDuration")
	AutoPprofGatherTraceDuration      = NewDuration("AutoPprofGatherTraceDuration")
	DatabaseURL                       = New("DatabaseURL", parse.DatabaseURL)
	BlockBackfillDepth                = NewUint64("BlockBackfillDepth")
	HTTPServerWriteTimeout            = NewDuration("HTTPServerWriteTimeout")
	JobPipelineMaxRunDuration         = NewDuration("JobPipelineMaxRunDuration")
	JobPipelineMaxSuccessfulRuns      = NewUint64("JobPipelineMaxSuccessfulRuns")
	JobPipelineReaperInterval         = NewDuration("JobPipelineReaperInterval")
	JobPipelineReaperThreshold        = NewDuration("JobPipelineReaperThreshold")
	JobPipelineResultWriteQueueDepth  = NewUint64("JobPipelineResultWriteQueueDepth")
	KeeperRegistryCheckGasOverhead    = NewUint32("KeeperRegistryCheckGasOverhead")
	KeeperRegistryPerformGasOverhead  = NewUint32("KeeperRegistryPerformGasOverhead")
	KeeperRegistryMaxPerformDataSize  = NewUint32("KeeperRegistryMaxPerformDataSize")
	KeeperRegistrySyncInterval        = NewDuration("KeeperRegistrySyncInterval")
	KeeperRegistrySyncUpkeepQueueSize = NewUint32("KeeperRegistrySyncUpkeepQueueSize")
	LogLevel                          = New[zapcore.Level]("LogLevel", parse.LogLevel)
	LogSQL                            = NewBool("LogSQL")
	RootDir                           = New[string]("RootDir", parse.HomeDir)
	JSONConsole                       = NewBool("JSONConsole")
	LogFileMaxSize                    = New("LogFileMaxSize", parse.FileSize)
	LogFileMaxAge                     = New("LogFileMaxAge", parse.Int64)
	LogFileMaxBackups                 = New("LogFileMaxBackups", parse.Int64)
	LogUnixTS                         = NewBool("LogUnixTS")
)

// EnvVar is an environment variable parsed as T.
type EnvVar[T any] struct {
	name  string
	parse func(string) (T, error)

	envVarName   string
	defaultValue string
	hasDefault   bool
}

// New creates a new EnvVar for the given name and parse func.
// name must match the ConfigSchema field.
func New[T any](name string, parse func(string) (T, error)) *EnvVar[T] {
	e := &EnvVar[T]{name: name, parse: parse, envVarName: Name(name)}
	e.defaultValue, e.hasDefault = DefaultValue(name)
	return e
}

// Parse attempts to parse the value returned from the environment, falling back to the default value when empty or invalid.
func (e *EnvVar[T]) Parse() (v T, invalid string) {
	var err error
	v, invalid, err = e.ParseFrom(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// ParsePtr attempts to parse the value from the environment, returning nil if the env var was empty or invalid.
func (e *EnvVar[T]) ParsePtr() *T {
	if os.Getenv(e.envVarName) == "" {
		return nil
	}
	v, invalid, err := e.ParseFrom(os.Getenv)
	if err != nil {
		log.Fatal(e.envVarName, err)
	}
	if invalid != "" {
		return nil
	}
	return &v
}

// ParseFrom attempts to parse the value returned from calling get with the env var name, falling back to the default
// value when empty or invalid.
func (e *EnvVar[T]) ParseFrom(get func(string) string) (v T, invalid string, err error) {
	str := get(e.envVarName)
	if str != "" {
		v, err = e.parse(str)
		if err == nil {
			return
		}
		var df interface{} = e.defaultValue
		if !e.hasDefault {
			var t T
			df = t
		}
		invalid = fmt.Sprintf(`Invalid value provided for %s, "%s" - falling back to default "%s": %v`, e.name, str, df, err)
		err = nil
	}

	if !e.hasDefault {
		// zero value
		return
	}

	v, err = e.parse(e.defaultValue)
	err = errors.Wrapf(err, `Invalid default for %s, "%s"`, e.name, e.defaultValue)
	return
}

func NewString(name string) *EnvVar[string] {
	return New[string](name, parse.String)
}

func NewBool(name string) *EnvVar[bool] {
	return New[bool](name, strconv.ParseBool)
}

func NewInt64(name string) *EnvVar[int64] {
	return New[int64](name, parse.Int64)
}

func NewUint64(name string) *EnvVar[uint64] {
	return New[uint64](name, parse.Uint64)
}

func NewUint32(name string) *EnvVar[uint32] {
	return New[uint32](name, parse.Uint32)
}

func NewUint16(name string) *EnvVar[uint16] {
	return New[uint16](name, parse.Uint16)
}

func NewDuration(name string) *EnvVar[time.Duration] {
	return New[time.Duration](name, time.ParseDuration)
}
