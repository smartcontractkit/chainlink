package envvar

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/config/parse"
)

var (
	LogLevel    = New("LogLevel", parse.LogLevel)
	RootDir     = New("RootDir", parse.HomeDir)
	JSONConsole = New("JSONConsole", parse.Bool)
	LogToDisk   = New("LogToDisk", parse.Bool)
	LogUnixTS   = New("LogUnixTS", parse.Bool)
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
