package loggers

import (
	"github.com/sirupsen/logrus"
	"github.com/smartcontractkit/libocr/commontypes"
)

var _ commontypes.Logger = LogrusLogger{}

type LogrusLogger struct {
	logger *logrus.Logger
}

func MakeLogrusLogger() LogrusLogger {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	return LogrusLogger{
		logger,
	}
}

func (l LogrusLogger) Trace(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (l LogrusLogger) Debug(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l LogrusLogger) Info(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l LogrusLogger) Warn(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Warn(msg)
}

func (l LogrusLogger) Error(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Error(msg)
}

func (l LogrusLogger) Critical(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Error("CRITICAL: " + msg)
}
