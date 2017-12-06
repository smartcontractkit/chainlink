package logger

import (
	"log"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
)

var logger *Logger

type Logger struct {
	*zap.SugaredLogger
}

func init() {
	logger = NewLogger("production")
}

func NewLogger(env string) *Logger {
	config := generateConfig(env)
	zap, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return &Logger{zap.Sugar()}
}

func SetLogger(newLogger *Logger) {
	defer logger.Sync()
	logger = newLogger
}

func generateConfig(env string) zap.Config {
	config := zap.NewProductionConfig()
	dir, err := homedir.Expand("~/.chainlink")
	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir(dir, os.FileMode(0700))
	destination := path.Join(dir, "log."+env+".jsonl")
	config.OutputPaths = []string{"stdout", destination}
	return config
}

func (self *Logger) Write(b []byte) (n int, err error) {
	self.Info(string(b))
	return len(b), nil
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
