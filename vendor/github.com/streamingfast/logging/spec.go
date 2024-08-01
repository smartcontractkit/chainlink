package logging

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"go.uber.org/zap/zapcore"
)

type levelSpec struct {
	key      string
	level    zapcore.Level
	trace    bool
	ordering int
}

func (l *levelSpec) String() string {
	return fmt.Sprintf("key %q (level %s, trace %t, order %d)", l.key, l.level, l.trace, l.ordering)
}

type logLevelSpec struct {
	incrementingIndex int
	byKey             map[string][]*levelSpec
	byLevel           map[zapcore.Level][]*levelSpec
}

func (l *logLevelSpec) String() string {
	if len(l.byKey) == 0 {
		return "<empty>"
	}

	byKeySpecs := make([]string, 0, len(l.byKey))
	for key, keySpecs := range l.byKey {
		keySpecStrings := make([]string, len(keySpecs))
		for i, keySpec := range keySpecs {
			keySpecStrings[i] = keySpec.String()
		}

		byKeySpecs = append(byKeySpecs, fmt.Sprintf("%s => [%s]", key, strings.Join(keySpecStrings, ", ")))
	}

	byLevelSpecs := make([]string, 0, len(l.byLevel))
	for level, levelSpecs := range l.byLevel {
		keySpecStrings := make([]string, len(levelSpecs))
		for i, keySpec := range levelSpecs {
			keySpecStrings[i] = keySpec.String()
		}

		byLevelSpecs = append(byLevelSpecs, fmt.Sprintf("%s => [%s]", level, strings.Join(keySpecStrings, ", ")))
	}

	return fmt.Sprintf("By Key (%s), By Level (%s)", strings.Join(byKeySpecs, " | "), strings.Join(byLevelSpecs, " | "))
}

func envGetFromMap(mappings map[string]string) func(string) string {
	return func(key string) string {
		return mappings[key]
	}
}

func newLogLevelSpecFromEnv(prefix string) *logLevelSpec {
	return newLogLevelSpec(envGetFromMap(map[string]string{
		"TRACE":   os.Getenv(prefix + "TRACE"),
		"DEBUG":   os.Getenv(prefix + "DEBUG"),
		"INFO":    os.Getenv(prefix + "INFO"),
		"WARNING": os.Getenv(prefix + "WARNING"),
		"WARN":    os.Getenv(prefix + "WARN"),
		"ERROR":   os.Getenv(prefix + "ERROR"),
		"DLOG":    os.Getenv(prefix + "DLOG"),
	}))
}

func newLogLevelSpec(envGet func(string) string) *logLevelSpec {
	spec := &logLevelSpec{byKey: map[string][]*levelSpec{}, byLevel: map[zapcore.Level][]*levelSpec{}}

	// Ordering is important, a DEBUG will overrides ERROR if there is conflict(s)
	spec.fillEnvFlat(zapcore.ErrorLevel, false, "ERROR", envGet)
	spec.fillEnvFlat(zapcore.WarnLevel, false, "WARNING", envGet)
	spec.fillEnvFlat(zapcore.WarnLevel, false, "WARN", envGet)
	spec.fillEnvFlat(zapcore.InfoLevel, false, "INFO", envGet)
	spec.fillEnvFlat(zapcore.DebugLevel, false, "DEBUG", envGet)
	spec.fillEnvFlat(zapcore.DebugLevel, true, "TRACE", envGet)

	input := strings.TrimSpace(envGet("DLOG"))
	if input != "" {
		spec.fillKeyValue(input)
	}

	return spec
}

func (s *logLevelSpec) add(key string, level zapcore.Level, trace bool) {
	s.incrementingIndex++
	s.byKey[key] = append(s.byKey[key], &levelSpec{key, level, trace, s.incrementingIndex})
}

func (s *logLevelSpec) fillEnvFlat(level zapcore.Level, trace bool, key string, envGet func(string) string) {
	input := strings.TrimSpace(envGet(key))
	if input != "" {
		s.fillFlat(level, trace, input)
	}
}

// fillKeyValue parse the input received against the following format:
//
// ```
// <key1>,<key2>,...
// ```
//
// How actually the key are interpreted is not the responsibility of the
// spec. The <key> can be any string as long as it does not contain "," and "="
// while.
//
func (s *logLevelSpec) fillFlat(level zapcore.Level, trace bool, input string) {
	for _, key := range strings.Split(input, ",") {
		key = strings.TrimSpace(key)
		if key != "" {
			s.add(key, level, trace)
		}
	}
}

// fillKeyValue parse the input received against the following format:
//
// ```
// <key1>=<level>,<key2>=<level>,...
// ```
//
// How actually the key are interpreted is not the responsibility of the
// spec. The <key> can be any string as long as it does not contain "," and "="
// while the <level> should be one of 'error', 'warn' (or 'warning'), 'info', 'debug'
// or 'trace'.
//
func (s *logLevelSpec) fillKeyValue(input string) {
	for _, keyValue := range strings.Split(input, ",") {
		keyValue = strings.TrimSpace(keyValue)
		if keyValue == "" {
			continue
		}

		parts := strings.SplitN(keyValue, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" || value == "" {
			continue
		}

		level, trace, ok := valueToLevelAndTrace(value)
		if !ok {
			continue
		}

		s.add(key, level, trace)
	}
}

// sortedSpecs sorts the current `spec` instance by their `ordering` value, The `ordering`
// value represents the order in which this spec was added relative to all others. It means
// that an low ordering value has less precedence than a high ordering value.
//
// The actual semantic of the ordering is left to the one constructing the overall `logLevelSpec`,
// it could be for example that a level read from a config file has more precedence than one
// created from and environment variable.
func (s *logLevelSpec) sortedSpecs() (specs []*levelSpec) {
	specs = make([]*levelSpec, 0, len(s.byKey))
	for _, specSlice := range s.byKey {
		specs = append(specs, specSlice...)
	}

	sort.Slice(specs, func(i, j int) bool {
		return specs[i].ordering < specs[j].ordering
	})

	return
}

func valueToLevelAndTrace(input string) (zapcore.Level, bool, bool) {
	switch strings.ToLower(input) {
	case "trace":
		return zapcore.DebugLevel, true, true
	case "debug":
		return zapcore.DebugLevel, false, true
	case "info":
		return zapcore.InfoLevel, false, true
	case "warn", "warning":
		return zapcore.WarnLevel, false, true
	case "error":
		return zapcore.ErrorLevel, false, true
	}

	// Invalid case, the actual level is there but should not be considered
	return zapcore.PanicLevel, false, false
}
