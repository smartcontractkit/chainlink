package services

import (
	"log"
	"os"
	"path"

	"github.com/caarlos0/env"
	homedir "github.com/mitchellh/go-homedir"
)

type Config struct {
	RootDir           string `env:"ROOT" envDefault:"~/.chainlink"`
	BasicAuthUsername string `env:"USERNAME" envDefault:"chainlink"`
	BasicAuthPassword string `env:"PASSWORD" envDefault:"twochains"`
}

func New() Config {
	config := Config{}
	env.Parse(&config)
	dir, err := homedir.Expand(config.RootDir)
	if err != nil {
		log.Fatal(err)
	}
	if err = os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		log.Fatal(err)
	}
	config.RootDir = dir
	return config
}

func (self Config) KeysDir() string {
	return path.Join(self.RootDir, "keys")
}
