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

type ServiceLogLevel struct {
	Name  string
	Level LogLevel
}

type ServiceLogLevelInput struct {
	Name  string   `json:"name"`
	Level LogLevel `json:"level"`
}

type ServiceLogLevelResolver struct {
	logLvl ServiceLogLevel
}

func NewServiceLogLevel(logLvl ServiceLogLevel) *ServiceLogLevelResolver {
	return &ServiceLogLevelResolver{logLvl}
}

func NewServicesLogLevel(logLvls []ServiceLogLevel) []*ServiceLogLevelResolver {
	var resolvers []*ServiceLogLevelResolver
	for _, logLvl := range logLvls {
		resolvers = append(resolvers, NewServiceLogLevel(logLvl))
	}

	return resolvers
}

func (r *ServiceLogLevelResolver) Name() string {
	return r.logLvl.Name
}

func (r *ServiceLogLevelResolver) Level() LogLevel {
	return r.logLvl.Level
}

// -- SetServiceLogLevel Mutation --

type SetServicesLogLevelsPayloadResolver struct {
	svcLvls   []ServiceLogLevel
	inputErrs map[string]string
}

func NewSetServicesLogLevelsPayload(svcLvls []ServiceLogLevel, inputErrs map[string]string) *SetServicesLogLevelsPayloadResolver {
	return &SetServicesLogLevelsPayloadResolver{svcLvls, inputErrs}
}

func (r *SetServicesLogLevelsPayloadResolver) ToSetServicesLogLevelsSuccess() (*SetServicesLogLevelsSuccessResolver, bool) {
	if r.inputErrs != nil {
		return nil, false
	}

	return NewSetServicesLogLevelsSuccess(r.svcLvls), true
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
	svcLvls []ServiceLogLevel
}

func NewSetServicesLogLevelsSuccess(svcLvls []ServiceLogLevel) *SetServicesLogLevelsSuccessResolver {
	return &SetServicesLogLevelsSuccessResolver{svcLvls}
}

func (r *SetServicesLogLevelsSuccessResolver) LogLevels() []*ServiceLogLevelResolver {
	return NewServicesLogLevel(r.svcLvls)
}
