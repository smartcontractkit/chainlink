package envvar

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/config/parse"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	// LogLevel reprents a parseable version of the `LOG_LEVEL`env var.
	LogLevel = New("LogLevel", parse.LogLevel)
	// RootDir reprents a parseable version of the `ROOT`env var.
	RootDir = New("RootDir", parse.HomeDir)
	// JSONConsole reprents a parseable version of the `JSON_CONSOLE`env var.
	JSONConsole = New("JSONConsole", parse.Bool)
	// LogFileMaxSize reprents a parseable version of the `LOG_FILE_MAX_SIZE`env var.
	LogFileMaxSize = New("LogFileMaxSize", parse.FileSize)
	// LogFileMaxAge reprents a parseable version of the `LOG_FILE_MAX_AGE`env var.
	LogFileMaxAge = New("LogFileMaxAge", parse.Int64)
	// LogFileMaxBackups reprents a parseable version of the `LOG_FILE_MAX_BACKUPS`env var.
	LogFileMaxBackups = New("LogFileMaxBackups", parse.Int64)
	// LogUnixTS reprents a parseable version of the `LOG_UNIX_TS`env var.
	LogUnixTS = New("LogUnixTS", parse.Bool)
)

// EnvVar is an environment variable which
type EnvVar struct {
	name  string
	parse func(string) (interface{}, error)

	envVarName   string
	defaultValue string
	hasDefault   bool
}

// New creates a new EnvVar for the given name and parse func.
// name must match the ConfigSchema field.
func New(name string, parse func(string) (interface{}, error)) *EnvVar {
	e := &EnvVar{name: name, parse: parse, envVarName: Name(name)}
	e.defaultValue, e.hasDefault = DefaultValue(name)
	return e
}

// Parse attempts to parse the value returned from the environment, falling back to the default value when empty or invalid.
func (e *EnvVar) Parse() (v interface{}, invalid string) {
	var err error
	v, invalid, err = e.ParseFrom(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// ParseFrom attempts to parse the value returned from calling get with the env var name, falling back to the default
// value when empty or invalid.
func (e *EnvVar) ParseFrom(get func(string) string) (v interface{}, invalid string, err error) {
	str := get(e.envVarName)
	if str != "" {
		v, err = e.parse(str)
		if err == nil {
			return
		}
		var df interface{} = e.defaultValue
		if !e.hasDefault {
			df = ZeroValue(e.name)
		}
		invalid = fmt.Sprintf(`Invalid value provided for %s, "%s" - falling back to default "%s": %v`, e.name, str, df, err)
	}

	if !e.hasDefault {
		v = ZeroValue(e.name)
		return
	}

	v, err = e.parse(e.defaultValue)
	err = errors.Wrapf(err, `Invalid default for %s, "%s"`, e.name, e.defaultValue)
	return
}

func (e *EnvVar) ParseString() (v string, invalid string) {
	var i interface{}
	i, invalid = e.Parse()
	return i.(string), invalid
}

func (e *EnvVar) ParseBool() (v bool, invalid string) {
	var i interface{}
	i, invalid = e.Parse()
	return i.(bool), invalid
}

// ParseInt64 parses value into `int64`
func (e *EnvVar) ParseInt64() (v int64, invalid string) {
	var i interface{}
	i, invalid = e.Parse()
	return i.(int64), invalid
}

// ParseFileSize parses value into `utils.FileSize`
func (e *EnvVar) ParseFileSize() (v utils.FileSize, invalid string) {
	var i interface{}
	i, invalid = e.Parse()
	return i.(utils.FileSize), invalid
}

func (e *EnvVar) ParseLogLevel() (v zapcore.Level, invalid string) {
	var i interface{}
	i, invalid = e.Parse()
	var ll zapcore.Level
	switch v := i.(type) {
	case zapcore.Level:
		ll = v
	case *zapcore.Level:
		ll = *v
	}
	return ll, invalid
}
