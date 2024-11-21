package llo

import "github.com/smartcontractkit/chainlink-common/pkg/logger"

// Suppressed logger swallows debug/info unless the verbose flag is turned on
// Useful for OCR to calm down its verbosity

var _ logger.Logger = &SuppressedLogger{}

func NewSuppressedLogger(lggr logger.Logger, verbose bool) logger.Logger {
	return &SuppressedLogger{
		Logger:  lggr,
		Verbose: verbose,
	}
}

type SuppressedLogger struct {
	logger.Logger
	Verbose bool
}

func (s *SuppressedLogger) Debug(args ...interface{}) {
	if s.Verbose {
		s.Logger.Debug(args...)
	}
}
func (s *SuppressedLogger) Info(args ...interface{}) {
	if s.Verbose {
		s.Logger.Info(args...)
	}
}
func (s *SuppressedLogger) Debugf(format string, values ...interface{}) {
	if s.Verbose {
		s.Logger.Debugf(format, values...)
	}
}
func (s *SuppressedLogger) Infof(format string, values ...interface{}) {
	if s.Verbose {
		s.Logger.Infof(format, values...)
	}
}
func (s *SuppressedLogger) Debugw(msg string, keysAndValues ...interface{}) {
	if s.Verbose {
		s.Logger.Debugw(msg, keysAndValues...)
	}
}
func (s *SuppressedLogger) Infow(msg string, keysAndValues ...interface{}) {
	if s.Verbose {
		s.Logger.Infow(msg, keysAndValues...)
	}
}
