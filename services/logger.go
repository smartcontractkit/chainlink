package services

import (
	"log"
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
	logger = newLogger
}

func generateConfig(env string) zap.Config {
	config := zap.NewProductionConfig()
	dir, err := homedir.Expand("~/.chainlink")
	if err != nil {
		log.Fatal(err)
	}
	destination := path.Join(dir, "log."+env+".jsonl")
	config.OutputPaths = []string{"stdout", destination}
	return config
}

func (self *Logger) Write(b []byte) (n int, err error) {
	self.Info(string(b))
	return len(b), nil
}

func GetLogger() *Logger {
	return logger
}
