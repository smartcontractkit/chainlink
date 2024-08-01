// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger = zap.NewNop()

type loggerConfig struct {
	isRootLogger   bool
	shortName      string
	defaultLevel   *zapcore.Level
	isTraceEnabled *bool
	onUpdate       func(newLogger *zap.Logger)
}

// LoggerOption are option parameters that you can set when creating a `PackageLogger`.
type LoggerOption interface {
	apply(config *loggerConfig)
}

type loggerOptionFunc func(config *loggerConfig)

func (f loggerOptionFunc) apply(config *loggerConfig) {
	f(config)
}

// Deprecated: Use LoggerOnUpdate instead.
func RegisterOnUpdate(onUpdate func(newLogger *zap.Logger)) LoggerOption {
	return LoggerOnUpdate(onUpdate)
}

// LoggerOnUpdate enable you to have a hook function that will receive the new logger
// that is going to be assigned to your logger instance. This is useful in some situation
// where you need to update other instances or re-configuring a bit the logger when
// a new one is attached.
//
// This is called **after** the instance has been re-assigned.
func LoggerOnUpdate(onUpdate func(newLogger *zap.Logger)) LoggerOption {
	return loggerOptionFunc(func(config *loggerConfig) {
		config.onUpdate = onUpdate
	})
}

// LoggerDefaultLevel can be used to set the default level of the logger if nothing else is overriding it.
//
// While the library offers you to set the default level, we recommend to not use this method
// unless you feel is strictly necessary, specially in libraries code. Indeed, setting for example your
// level to `INFO` on the loggers of your library would mean that anyone importing your code
// and instantiating the loggers would automatically see your `INFO` log line which is usually
// disruptive.
//
// Instead of using this, use `logging.WithDefaultSpec` to specify a default level via the logger's short
// name for example.
func LoggerDefaultLevel(level zapcore.Level) LoggerOption {
	return loggerOptionFunc(func(config *loggerConfig) {
		config.defaultLevel = &level
	})
}

func loggerRoot() LoggerOption {
	return loggerOptionFunc(func(config *loggerConfig) {
		config.isRootLogger = true
	})
}

func loggerShortName(shortName string) LoggerOption {
	return loggerOptionFunc(func(config *loggerConfig) {
		config.shortName = shortName
	})
}

func loggerWithTracer(isEnabled *bool) LoggerOption {
	return loggerOptionFunc(func(config *loggerConfig) {
		config.isTraceEnabled = isEnabled
	})
}

type LoggerExtender func(*zap.Logger) *zap.Logger

type loggerFactory func(name string, level zap.AtomicLevel) *zap.Logger

type registryEntry struct {
	isRoot       bool
	packageID    string
	shortName    string
	atomicLevel  zap.AtomicLevel
	traceEnabled *bool
	logPtr       *zap.Logger
	onUpdate     func(newLogger *zap.Logger)
}

func (e *registryEntry) String() string {
	return e.string(false)
}

func (e *registryEntry) string(extended bool) string {
	shortName := "<none>"
	if e.shortName != "" {
		shortName = e.shortName
	}

	loggerPtr := "<nil>"
	levels := ""
	if e.logPtr != nil {
		loggerPtr = fmt.Sprintf("%p", e.logPtr)
		if extended {
			levels = " [" + computeLevelsString(e.logPtr.Core()) + "]"
		}
	}

	traceEnabled := false
	if e.traceEnabled != nil {
		traceEnabled = *e.traceEnabled
	}

	return fmt.Sprintf("%s @ %s (level: %s, trace?: %t, ptr: %s%s)", shortName, e.packageID, e.atomicLevel.Level(), traceEnabled, loggerPtr, levels)
}

var zapLevels = []zapcore.Level{
	zap.DebugLevel,
	zap.InfoLevel,
	zap.WarnLevel,
	zap.ErrorLevel,
	zap.DPanicLevel,
	zap.PanicLevel,
}

func computeLevelsString(core zapcore.Core) string {
	levels := make([]string, len(zapLevels))
	for i, level := range zapLevels {
		state := "Disabled"
		if core.Enabled(level) {
			state = "Enabled"
		}

		levels[i] = fmt.Sprintf("%s => %s", level.String(), state)
	}

	return strings.Join(levels, ", ")
}

// Deprecated: Use `var zlog, _ = logging.PackageLogger(<shortName>, "...")` instead.
func Register(packageID string, zlogPtr **zap.Logger, options ...LoggerOption) {
	if *zlogPtr == nil {
		*zlogPtr = zap.NewNop()
	}
	register(globalRegistry, packageID, *zlogPtr, options...)
}

func register2(registry *registry, shortName string, packageID string, options ...LoggerOption) (*zap.Logger, Tracer) {
	logger := zap.NewNop()
	tracer := boolTracer{new(bool)}

	allOptions := append([]LoggerOption{
		loggerShortName(shortName),
		loggerWithTracer(tracer.value),
	}, options...)

	register(registry, packageID, logger, allOptions...)

	return logger, tracer
}

func register(registry *registry, packageID string, zlogPtr *zap.Logger, options ...LoggerOption) {
	if zlogPtr == nil {
		panic("the zlog pointer (of type **zap.Logger) must be set")
	}

	config := loggerConfig{}
	for _, opt := range options {
		opt.apply(&config)
	}

	defaultLevel := zapcore.ErrorLevel
	if config.defaultLevel != nil {
		defaultLevel = *config.defaultLevel
	}

	entry := &registryEntry{
		isRoot:       config.isRootLogger,
		packageID:    packageID,
		shortName:    config.shortName,
		traceEnabled: config.isTraceEnabled,
		atomicLevel:  zap.NewAtomicLevelAt(defaultLevel),
		logPtr:       zlogPtr,
		onUpdate:     config.onUpdate,
	}

	registry.registerEntry(entry)

	logger := defaultLogger
	if zlogPtr != nil {
		logger = zlogPtr
	}

	// The tracing has already been set, so we can go unspecified here to not change anything
	setLogger(entry, logger, unspecifiedTracing)
}

func Set(logger *zap.Logger, regexps ...string) {
	for name, entry := range globalRegistry.entriesByPackageID {
		if len(regexps) == 0 {
			setLogger(entry, logger, unspecifiedTracing)
		} else {
			for _, re := range regexps {
				regex, err := regexp.Compile(re)
				if (err == nil && regex.MatchString(name)) || (err != nil && name == re) {
					setLogger(entry, logger, unspecifiedTracing)
				}
			}
		}
	}
}

// Extend is different than `Set` by being able to re-configure the existing logger set for
// all registered logger in the registry. This is useful for example to add a field to the
// currently set logger:
//
// ```
// logger.Extend(func (current *zap.Logger) { return current.With("name", "value") }, "github.com/dfuse-io/app.*")
// ```
func Extend(extender LoggerExtender, regexps ...string) {
	extend(extender, unspecifiedTracing, regexps...)
}

func extend(extender LoggerExtender, tracing tracingType, regexps ...string) {
	for name, entry := range globalRegistry.entriesByPackageID {
		if entry.logPtr == nil {
			continue
		}

		if len(regexps) == 0 {
			setLogger(entry, extender(entry.logPtr), tracing)
		} else {
			for _, re := range regexps {
				if regexp.MustCompile(re).MatchString(name) {
					setLogger(entry, extender(entry.logPtr), tracing)
				}
			}
		}
	}
}

// Override sets the given logger on previously registered and next
// registrations.  Useful in tests.
//
// Deprecated: Call `logging.InstantiateLoggers` directly and use the `logging.WithDefaultSpec`
// to configure the various loggers.
func Override(logger *zap.Logger) {
	defaultLogger = logger
	Set(logger)
}

// TestingOverride calls `Override` (or `Set`, see below) with a development
// logger setup correctly with the right level based on some environment variables.
//
// By default, override using a `zap.NewDevelopment` logger (`info`), if
// environment variable `DEBUG` is set to anything or environment variable `TRACE`
// is set to `true`, logger is set in `debug` level.
//
// If `DEBUG` is set to something else than `true` and/or if `TRACE` is set
// to something else than
//
// Deprecated: Call `logging.InstantiateLoggers` directly instead.
func TestingOverride() {
	debug := os.Getenv("DEBUG")
	trace := os.Getenv("TRACE")
	if debug == "" && trace == "" {
		return
	}

	logger, _ := zap.NewDevelopment()

	regex := ""
	if debug != "true" {
		regex = debug
	}

	if regex == "" && trace != "true" {
		regex = trace
	}

	if regex == "" {
		Override(logger)
	} else {
		for _, regexPart := range strings.Split(regex, ",") {
			regexPart = strings.TrimSpace(regexPart)
			if regexPart != "" {
				Set(logger, regexPart)
			}
		}
	}
}

type tracingType uint8

const (
	unspecifiedTracing tracingType = iota
	enableTracing
	disableTracing
)

func setLogger(entry *registryEntry, logger *zap.Logger, tracing tracingType) {
	if entry == nil || logger == nil {
		return
	}

	ve := reflect.ValueOf(entry.logPtr).Elem()
	ve.Set(reflect.ValueOf(logger).Elem())

	if entry.traceEnabled != nil && tracing != unspecifiedTracing {
		switch tracing {
		case enableTracing:
			*entry.traceEnabled = true
		case disableTracing:
			*entry.traceEnabled = false
		}
	}

	if entry.onUpdate != nil {
		entry.onUpdate(logger)
	}
}

type registry struct {
	sync.RWMutex

	name               string
	factory            loggerFactory
	entriesByPackageID map[string]*registryEntry
	entriesByShortName map[string][]*registryEntry

	rootEntry *registryEntry

	dbgLogger *zap.Logger
}

func newRegistry(name string, logger *zap.Logger) *registry {
	registryLogger := logger.Named(name)
	registryLogger.Info("creating registry")

	registry := &registry{
		name:               name,
		entriesByPackageID: make(map[string]*registryEntry),
		entriesByShortName: make(map[string][]*registryEntry),
		dbgLogger:          registryLogger,
	}

	registry.factory = func(name string, level zap.AtomicLevel) *zap.Logger {
		loggerOptions := newInstantiateOptions()

		return newLogger(registry.dbgLogger, name, level, &loggerOptions)
	}

	return registry
}

func debugLoggerForLoggingLibrary() (*zap.Logger, Tracer) {
	registry := newRegistry("logging_dbg", zap.NewNop())
	logger, tracer := packageLogger(registry, "logging", "github.com/streamingfast/logging")

	registry.dbgLogger = logger

	registry.forAllEntries(func(entry *registryEntry) {
		registry.createLoggerForEntry(entry)
	})

	spec := newLogLevelSpecFromEnv("__LOGGING_")
	registry.forAllEntriesMatchingSpec(spec, func(entry *registryEntry, level zapcore.Level, trace bool) {
		registry.setLevelForEntry(entry, level, trace)
	})

	return logger, tracer
}

func (r *registry) registerEntry(entry *registryEntry) {
	if entry == nil {
		panic("refusing to add a nil registry entry")
	}

	id := validateEntryIdentifier("package ID", entry.packageID, false)
	shortName := validateEntryIdentifier("short name", entry.shortName, true)

	if actual := r.entriesByPackageID[id]; actual != nil {
		panic(fmt.Sprintf("packageID %q is already registered", id))
	}

	entry.packageID = id
	entry.shortName = shortName

	r.entriesByPackageID[id] = entry
	if shortName != "" {
		r.entriesByShortName[shortName] = append(r.entriesByShortName[shortName], entry)
	}

	if entry.isRoot {
		if r.rootEntry != nil {
			panic(fmt.Errorf("trying to register a second root logger, existing root logger is registered under %s (%s), trying to now register %s (%s)",
				r.rootEntry.shortName,
				r.rootEntry.packageID,
				entry.shortName,
				entry.packageID,
			))
		}

		r.rootEntry = entry
		r.dbgLogger.Info("registering root logger", zap.Stringer("entry", r.rootEntry))
	}

	r.dbgLogger.Info("registered entry", zap.String("short_name", shortName), zap.String("id", id))
}

func (r *registry) forAllEntries(callback func(entry *registryEntry)) {
	for _, entry := range r.entriesByPackageID {
		callback(entry)
	}
}

// forAllEntriesMatchingSpec iterate sequentially through the sorted spec
func (r *registry) forAllEntriesMatchingSpec(spec *logLevelSpec, callback func(entry *registryEntry, level zapcore.Level, trace bool)) {
	for _, specForKey := range spec.sortedSpecs() {
		if specForKey.key == "true" || specForKey.key == "*" {
			for _, entry := range r.entriesByPackageID {
				callback(entry, specForKey.level, specForKey.trace)
			}

			continue
		}

		r.forEntriesMatchingSpec(specForKey, callback)
	}
}

func (r *registry) forEntriesMatchingSpec(spec *levelSpec, callback func(entry *registryEntry, level zapcore.Level, trace bool)) {
	r.dbgLogger.Debug("looking in short names to find spec key", zap.String("key", spec.key))
	entries, found := r.entriesByShortName[spec.key]
	if found {
		r.dbgLogger.Debug("found logger in short names", zap.Int("count", len(entries)))
		for _, entry := range entries {
			callback(entry, spec.level, spec.trace)
		}
		return
	}

	r.dbgLogger.Debug("looking in package IDs to find spec key", zap.String("key", spec.key))
	entry, found := r.entriesByPackageID[spec.key]
	if found {
		r.dbgLogger.Debug("found logger in package ID", zap.Stringer("entry", entry))
		callback(entry, spec.level, spec.trace)
		return
	}

	r.dbgLogger.Debug("looking in package IDs by regex", zap.String("key", spec.key))
	regex, err := regexp.Compile(spec.key)
	if err != nil {
		r.dbgLogger.Debug("spec key is not a regex, we already matched exact package ID, nothing to do more", zap.Error(err))
		return
	}

	for packageID, entry := range r.entriesByPackageID {
		if regex.MatchString(packageID) {
			callback(entry, spec.level, spec.trace)
		}
	}
}

func (r *registry) createLoggerForEntry(entry *registryEntry) {
	if entry == nil {
		return
	}

	r.dbgLogger.Info("creating logger on entry from registry factory",
		zap.Stringer("to_level", entry.atomicLevel),
		zap.Boolp("trace_enabled", entry.traceEnabled),
		zap.Stringer("entry", entry),
	)

	logger := r.factory(entry.shortName, entry.atomicLevel)

	ve := reflect.ValueOf(entry.logPtr).Elem()
	ve.Set(reflect.ValueOf(logger).Elem())

	if entry.onUpdate != nil {
		entry.onUpdate(logger)
	}
}

func (r *registry) setLevelForEntry(entry *registryEntry, level zapcore.Level, trace bool) {
	if entry == nil {
		return
	}

	r.dbgLogger.Info("setting logger level", zap.Stringer("to_level", level), zap.Bool("trace_enabled", trace), zap.Stringer("entry", entry))
	entry.atomicLevel.SetLevel(level)

	// It's possible for an entry to have no tracer registered, for example if the legacy
	// register method is used. We must protect from this and not set anything.
	if entry.traceEnabled != nil {
		*entry.traceEnabled = trace
	}
}

func (r *registry) dumpRegistryToLogger() {
	r.dbgLogger.Info("dumping registry to logger", zap.Int("entries", len(r.entriesByPackageID)))

	for _, entry := range r.entriesByPackageID {
		r.dbgLogger.Info("registered entry", zap.String("entry", entry.string(true)))
	}

	if r.rootEntry != nil {
		r.dbgLogger.Info("registered root entry", zap.Stringer("entry", r.rootEntry))
	} else {
		r.dbgLogger.Info("no root entry")
	}

	r.dbgLogger.Info("dumping terminated")
}

func validateEntryIdentifier(tag string, rawInput string, allowEmpty bool) string {
	input := strings.TrimSpace(rawInput)
	if input == "" && !allowEmpty {
		panic(fmt.Errorf("the %s %q is invalid, must not be empty", tag, input))
	}

	if input == "true" {
		panic(fmt.Errorf("the %s %q is invalid, the identifier 'true' is reserved", tag, input))
	}

	if input == "*" {
		panic(fmt.Errorf("the %s %q is invalid, the identifier '*' is reserved", tag, input))
	}

	if strings.HasPrefix(input, "-") {
		panic(fmt.Errorf("the %s %q is invalid, must not starts with the '-' character", tag, input))
	}

	if strings.Contains(input, ",") {
		panic(fmt.Errorf("the %s %q is invalid, must not contain the ',' character", tag, input))
	}

	if strings.Contains(input, "=") {
		panic(fmt.Errorf("the %s %q is invalid, must not contain the '=' character", tag, input))
	}

	return input
}
