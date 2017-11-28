package logger

import (
	homedir "github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"log"
	"path"
)

var logger *Logger

type Logger struct {
	*zap.Logger
}

func generateConfig() zap.Config {
	config := zap.NewProductionConfig()
	dir, err := homedir.Expand("~/.chainlink")
	if err != nil {
		log.Fatal(err)
	}
	destination := path.Join(dir, "log.jsonl")
	config.OutputPaths = []string{"stdout", destination}
	return config
}

func init() {
	config := generateConfig()
	zap, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	logger = &Logger{zap}
}

func ForGin() *Logger {
	config := generateConfig()
	config.DisableCaller = true
	zap, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return &Logger{zap}
}

func (self *Logger) Write(b []byte) (n int, err error) {
	self.Info(string(b))
	return len(b), nil
}

func Get() *Logger {
	return logger
}
