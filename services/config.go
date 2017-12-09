package services

import (
	"log"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

type Config struct {
	RootDir string
}

func NewConfig(dir string) Config {
	dir, err := homedir.Expand(dir)
	if err != nil {
		log.Fatal(err)
	}
	if err = os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		log.Fatal(err)
	}
	return Config{dir}
}

func (self Config) KeysDir() string {
	return path.Join(self.RootDir, "keys")
}
