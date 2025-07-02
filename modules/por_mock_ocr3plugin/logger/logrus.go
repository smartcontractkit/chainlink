package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/smartcontractkit/libocr/commontypes"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	return &Logger{
		logger,
	}
}

func (l *Logger) Trace(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (l *Logger) Debug(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l *Logger) Info(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l *Logger) Warn(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Warn(msg)
}

func (l *Logger) Error(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Error(msg)
}

func (l *Logger) Critical(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Error("CRITICAL: " + msg)
}
