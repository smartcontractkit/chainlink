package logger

import (
	"log"
	"path"

	"go.uber.org/zap"
)

var logger *Logger

func init() {
	zl, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	SetLogger(NewLogger(zl))
}

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) Write(b []byte) (n int, err error) {
	l.Info(string(b))
	return len(b), nil
}

func NewLogger(zl *zap.Logger) *Logger {
	return &Logger{zl.Sugar()}
}

func NewLoggerFromDir(dir string) *Logger {
	config := generateConfig(dir)
	zl, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return NewLogger(zl)
}

func SetLogger(l *Logger) {
	if logger != nil {
		defer logger.Sync()
	}
	logger = l
}

func SetLoggerDir(dir string) {
	SetLogger(NewLoggerFromDir(dir))
}

func generateConfig(dir string) zap.Config {
	config := zap.NewProductionConfig()
	destination := path.Join(dir, "log.jsonl")
	config.OutputPaths = []string{"stderr", destination}
	config.ErrorOutputPaths = []string{"stderr", destination}
	return config
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Panic(args ...interface{}) {
	logger.Panic(args)
}

func Sync() error {
	return logger.Sync()
}
