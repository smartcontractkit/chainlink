package loghelper

import "github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"

type LoggerWithContext struct {
	logger  types.Logger
	context types.LogFields
}

func MakeLoggerWithContext(logger types.Logger, context types.LogFields) LoggerWithContext {
	return LoggerWithContext{
		logger,
		context,
	}
}

func (l LoggerWithContext) addContextToFieldsIfNotPresent(fields types.LogFields) types.LogFields {
	if fields == nil {
		fields = l.context
	} else {
		for k, v := range l.context {
			if _, ok := fields[k]; !ok {
				fields[k] = v
			}
		}
	}
	return fields
}

func (l LoggerWithContext) Trace(msg string, fields types.LogFields) {
	l.logger.Trace(msg, l.addContextToFieldsIfNotPresent(fields))
}

func (l LoggerWithContext) Debug(msg string, fields types.LogFields) {
	l.logger.Debug(msg, l.addContextToFieldsIfNotPresent(fields))
}

func (l LoggerWithContext) Info(msg string, fields types.LogFields) {
	l.logger.Info(msg, l.addContextToFieldsIfNotPresent(fields))
}

func (l LoggerWithContext) Warn(msg string, fields types.LogFields) {
	l.logger.Warn(msg, l.addContextToFieldsIfNotPresent(fields))
}

func (l LoggerWithContext) Error(msg string, fields types.LogFields) {
	l.logger.Error(msg, l.addContextToFieldsIfNotPresent(fields))
}
