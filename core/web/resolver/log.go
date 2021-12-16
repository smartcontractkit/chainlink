package resolver

import (
	"strings"

	"github.com/pkg/errors"
)

type LogLevel string

const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
)

func FromLogLevel(logLvl LogLevel) string {
	switch logLvl {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		return strings.ToLower(string(logLvl))
	}
}

func ToLogLevel(str string) (LogLevel, error) {
	switch str {
	case "debug":
		return LogLevelDebug, nil
	case "info":
		return LogLevelInfo, nil
	case "warn":
		return LogLevelWarn, nil
	case "error":
		return LogLevelError, nil
	default:
		return "", errors.New("invalid log level")
	}
}

type LogLevelConfig struct {
	HeadTracker *LogLevel
	FluxMonitor *LogLevel
	Keeper      *LogLevel
}

type LogLevelConfigResolver struct {
	cfg LogLevelConfig
}

func NewLogLevelConfig(cfg LogLevelConfig) *LogLevelConfigResolver {
	return &LogLevelConfigResolver{cfg: cfg}
}

func (r *LogLevelConfigResolver) HeadTracker() *LogLevel {
	return r.cfg.HeadTracker
}

func (r *LogLevelConfigResolver) FluxMonitor() *LogLevel {
	return r.cfg.FluxMonitor
}

func (r *LogLevelConfigResolver) Keeper() *LogLevel {
	return r.cfg.Keeper
}

// -- SetServiceLogLevel Mutation --

type SetServicesLogLevelsPayloadResolver struct {
	cfg       *LogLevelConfig
	inputErrs map[string]string
}

func NewSetServicesLogLevelsPayload(cfg *LogLevelConfig, inputErrs map[string]string) *SetServicesLogLevelsPayloadResolver {
	return &SetServicesLogLevelsPayloadResolver{cfg: cfg, inputErrs: inputErrs}
}

func (r *SetServicesLogLevelsPayloadResolver) ToSetServicesLogLevelsSuccess() (*SetServicesLogLevelsSuccessResolver, bool) {
	if r.inputErrs != nil {
		return nil, false
	}

	return NewSetServicesLogLevelsSuccess(*r.cfg), true
}

func (r *SetServicesLogLevelsPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type SetServicesLogLevelsSuccessResolver struct {
	cfg LogLevelConfig
}

func NewSetServicesLogLevelsSuccess(cfg LogLevelConfig) *SetServicesLogLevelsSuccessResolver {
	return &SetServicesLogLevelsSuccessResolver{cfg: cfg}
}

func (r *SetServicesLogLevelsSuccessResolver) Config() *LogLevelConfigResolver {
	return NewLogLevelConfig(r.cfg)
}

// SQL Logging config

type SQLLoggingResolver struct {
	enabled bool
}

func NewSQLLogging(enabled bool) *SQLLoggingResolver {
	return &SQLLoggingResolver{enabled: enabled}
}

func (r *SQLLoggingResolver) Enabled() bool {
	return r.enabled
}

// -- GetSQLLogging Query --

type GetSQLLoggingPayloadResolver struct {
	enabled bool
}

func NewGetSQLLoggingPayload(enabled bool) *GetSQLLoggingPayloadResolver {
	return &GetSQLLoggingPayloadResolver{enabled: enabled}
}

func (r *GetSQLLoggingPayloadResolver) ToSQLLogging() (*SQLLoggingResolver, bool) {
	return NewSQLLogging(r.enabled), true
}

// -- SetSQLLogging Mutation --

type SetSQLLoggingPayloadResolver struct {
	enabled bool
}

func NewSetSQLLoggingPayload(enabled bool) *SetSQLLoggingPayloadResolver {
	return &SetSQLLoggingPayloadResolver{enabled: enabled}
}

func (r *SetSQLLoggingPayloadResolver) ToSetSQLLoggingSuccess() (*SetSQLLoggingSuccessResolver, bool) {
	return NewSetSQLLoggingSuccess(r.enabled), true
}

type SetSQLLoggingSuccessResolver struct {
	enabled bool
}

func NewSetSQLLoggingSuccess(enabled bool) *SetSQLLoggingSuccessResolver {
	return &SetSQLLoggingSuccessResolver{enabled: enabled}
}

func (r *SetSQLLoggingSuccessResolver) SQLLogging() *SQLLoggingResolver {
	return NewSQLLogging(r.enabled)
}

// -- GetLogLevel Query --

type GlobalLogLevelResolver struct {
	lvl string
}

func GlobalLogLevel(lvl string) *GlobalLogLevelResolver {
	return &GlobalLogLevelResolver{lvl: lvl}
}

func (r *GlobalLogLevelResolver) Level() (LogLevel, error) {
	return ToLogLevel(r.lvl)
}

type GlobalLogLevelPayloadResolver struct {
	lgLvl string
}

func NewGlobalLogLevelPayload(lgLvl string) *GlobalLogLevelPayloadResolver {
	return &GlobalLogLevelPayloadResolver{lgLvl: lgLvl}
}

func (r *GlobalLogLevelPayloadResolver) ToGlobalLogLevel() (*GlobalLogLevelResolver, bool) {
	return GlobalLogLevel(r.lgLvl), true
}

// -- UpdateGlobalLogLevel Mutation --

type SetGlobalLogLevelPayloadResolver struct {
	lvl       LogLevel
	inputErrs map[string]string
}

func NewSetGlobalLogLevelPayload(lvl LogLevel, inputErrs map[string]string) *SetGlobalLogLevelPayloadResolver {
	return &SetGlobalLogLevelPayloadResolver{lvl: lvl, inputErrs: inputErrs}
}

func (r *SetGlobalLogLevelPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

func (r *SetGlobalLogLevelPayloadResolver) ToSetGlobalLogLevelSuccess() (*SetGlobalLogLevelSuccessResolver, bool) {
	if r.inputErrs != nil {
		return nil, false
	}

	return NewSetGlobalLogLevelSuccess(r.lvl), true
}

type SetGlobalLogLevelSuccessResolver struct {
	lvl LogLevel
}

func NewSetGlobalLogLevelSuccess(lvl LogLevel) *SetGlobalLogLevelSuccessResolver {
	return &SetGlobalLogLevelSuccessResolver{lvl: lvl}
}

func (r *SetGlobalLogLevelSuccessResolver) GlobalLogLevel() *GlobalLogLevelResolver {
	return GlobalLogLevel(FromLogLevel(r.lvl))
}
