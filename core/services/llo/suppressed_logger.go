package llo

import "github.com/smartcontractkit/chainlink-common/pkg/logger"

// Suppressed logger swallows debug and/or info levels
// Useful for OCR to calm down its verbosity

var _ logger.Logger = &SuppressedLogger{}

func NewSuppressedLogger(lggr logger.Logger, debug, info bool) logger.Logger {
	return &SuppressedLogger{
		Logger:     lggr,
		DebugLevel: debug,
		InfoLevel:  info,
	}
}

type SuppressedLogger struct {
	logger.Logger
	DebugLevel bool
	InfoLevel  bool
}

func (s *SuppressedLogger) Debug(args ...interface{}) {
	if s.DebugLevel {
		s.Logger.Debug(args...)
	}
}
func (s *SuppressedLogger) Info(args ...interface{}) {
	if s.InfoLevel {
		s.Logger.Info(args...)
	}
}
func (s *SuppressedLogger) Debugf(format string, values ...interface{}) {
	if s.DebugLevel {
		s.Logger.Debugf(format, values...)
	}
}
func (s *SuppressedLogger) Infof(format string, values ...interface{}) {
	if s.InfoLevel {
		s.Logger.Infof(format, values...)
	}
}
func (s *SuppressedLogger) Debugw(msg string, keysAndValues ...interface{}) {
	if s.DebugLevel {
		s.Logger.Debugw(msg, keysAndValues...)
	}
}
func (s *SuppressedLogger) Infow(msg string, keysAndValues ...interface{}) {
	if s.InfoLevel {
		s.Logger.Infow(msg, keysAndValues...)
	}
}
