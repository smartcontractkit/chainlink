package resolver

import "strings"

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

type LogLevelConfig struct {
	HeadTracker *LogLevel
	FluxMonitor *LogLevel
	Keeper      *LogLevel
}

type LogLevelConfigResolver struct {
	cfg LogLevelConfig
}

func NewLogLevelConfig(cfg LogLevelConfig) *LogLevelConfigResolver {
	return &LogLevelConfigResolver{cfg}
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
	return &SetServicesLogLevelsPayloadResolver{cfg, inputErrs}
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
	return &SetServicesLogLevelsSuccessResolver{cfg}
}

func (r *SetServicesLogLevelsSuccessResolver) Config() *LogLevelConfigResolver {
	return NewLogLevelConfig(r.cfg)
}
