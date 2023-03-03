package util

import (
	"bytes"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/smartcontractkit/libocr/commontypes"
)

func XXXTestingOnly_MakePrefixLogger(prefix string) Logger {
	logger := MakeLogger()
	logger.logger.Out = prefixWriter{prefix, logger.logger.Out}
	return logger
}

type prefixWriter struct {
	prefix string
	writer io.Writer
}

func (p prefixWriter) Write(po []byte) (n int, err error) {

	if len(po) > 1000 {
		po = bytes.Join([][]byte{po[:1000], []byte("<truncated>\n")}, nil)
	}
	return p.writer.Write(
		bytes.Join([][]byte{[]byte(p.prefix), []byte(" "), po}, nil),
	)
}

type Logger struct {
	logger *logrus.Logger
}

func MakeLogger() Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	return Logger{
		logger,
	}
}

func (l Logger) Trace(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (l Logger) Debug(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l Logger) Info(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l Logger) Warn(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Warn(msg)
}

func (l Logger) Error(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Error(msg)
}

func (l Logger) Critical(msg string, fields commontypes.LogFields) {
	l.logger.WithFields(logrus.Fields(fields)).Error("CRITICAL: " + msg)
}
