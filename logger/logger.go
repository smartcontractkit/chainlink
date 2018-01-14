package logger

import (
	"log"
	"path"

	"go.uber.org/zap"
)

var logger *Logger

type Logger struct {
	*zap.SugaredLogger
}

func init() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	logger = &Logger{zapLogger.Sugar()}
}

func NewLogger(dir string) *Logger {
	config := generateConfig(dir)
	zap, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return &Logger{zap.Sugar()}
}

func SetLoggerDir(dir string) {
	defer logger.Sync()
	logger = NewLogger(dir)
}

func generateConfig(dir string) zap.Config {
	config := zap.NewProductionConfig()
	destination := path.Join(dir, "log.jsonl")
	config.OutputPaths = []string{"stdout", destination}
	return config
}

func (l *Logger) Write(b []byte) (n int, err error) {
	l.Info(string(b))
	return len(b), nil
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

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Panic(args ...interface{}) {
	logger.Panic(args)
}

func Sync() error {
	return logger.Sync()
}
