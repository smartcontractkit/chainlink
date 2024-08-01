package loghelper

import (
	"context"
	"errors"

	"github.com/smartcontractkit/libocr/commontypes"
)

type LoggerWithContext interface {
	commontypes.Logger
	MakeChild(extraContext commontypes.LogFields) LoggerWithContext
	MakeUpdated(updatedContext commontypes.LogFields) LoggerWithContext
	ErrorIfNotCanceled(msg string, ctx context.Context, fields commontypes.LogFields)
}

type loggerWithContextImpl struct {
	logger  commontypes.Logger
	context commontypes.LogFields
}

// MakeRootLoggerWithContext creates a base logger by wrapping a commontypes.Logger.
// NOTE! Most loggers should extend an existing LoggerWithContext using MakeChild!
func MakeRootLoggerWithContext(logger commontypes.Logger) LoggerWithContext {
	return loggerWithContextImpl{logger, commontypes.LogFields{}}
}

func (l loggerWithContextImpl) Trace(msg string, fields commontypes.LogFields) {
	l.logger.Trace(msg, MergePreserve(l.context, fields))
}

func (l loggerWithContextImpl) Debug(msg string, fields commontypes.LogFields) {
	l.logger.Debug(msg, MergePreserve(l.context, fields))
}

func (l loggerWithContextImpl) Info(msg string, fields commontypes.LogFields) {
	l.logger.Info(msg, MergePreserve(l.context, fields))
}

func (l loggerWithContextImpl) Warn(msg string, fields commontypes.LogFields) {
	l.logger.Warn(msg, MergePreserve(l.context, fields))
}

func (l loggerWithContextImpl) Error(msg string, fields commontypes.LogFields) {
	l.logger.Error(msg, MergePreserve(l.context, fields))
}

func (l loggerWithContextImpl) Critical(msg string, fields commontypes.LogFields) {
	l.logger.Critical(msg, MergePreserve(l.context, fields))
}

func (l loggerWithContextImpl) ErrorIfNotCanceled(msg string, ctx context.Context, fields commontypes.LogFields) {
	if !errors.Is(ctx.Err(), context.Canceled) {
		l.logger.Error(msg, MergePreserve(l.context, fields))
	} else {
		l.logger.Debug("logging as debug due to context cancellation: "+msg, MergePreserve(l.context, fields))
	}
}

// MakeChild is the preferred way to create a new specialized logger.
// It will reuse the base commontypes.Logger and create a new extended context.
func (l loggerWithContextImpl) MakeChild(extra commontypes.LogFields) LoggerWithContext {
	return loggerWithContextImpl{
		l.logger,
		MergePreserve(l.context, extra),
	}
}

// MakeUpdated will reuse the base commontypes.Logger and create a new extended context,
// overwriting any entries in the context with the ones from upserts.
func (l loggerWithContextImpl) MakeUpdated(upserts commontypes.LogFields) LoggerWithContext {
	return loggerWithContextImpl{
		l.logger,
		MergeOverwrite(l.context, upserts),
	}
}

// Helpers

// MergePreserve will create a new LogFields and add all the properties from extras on it.
// Key conflicts are resolved by prefixing the key for the new value with underscores until there's no conflict.
func MergePreserve(extras ...commontypes.LogFields) commontypes.LogFields {
	base := commontypes.LogFields{}
	for _, extra := range extras {
		for k, v := range extra {
			add(base, k, v)
		}
	}
	return base
}

// add (key, val) to base. If base already has key, then the old key will be
// left in place and the new key will be prefixed with underscore.
func add(base commontypes.LogFields, key string, val interface{}) {
	for {
		_, found := base[key]
		if found {
			key = "_" + key
			continue
		}
		base[key] = val
		return
	}
}

// MergeOverwrite will create a new LogFields and add all the properties from upserts on it.
// Key conflicts are resolved by preferring the upserted value.
func MergeOverwrite(upserts ...commontypes.LogFields) commontypes.LogFields {
	base := commontypes.LogFields{}
	for _, logfields := range upserts {
		for k, v := range logfields {
			base[k] = v
		}
	}
	return base
}
