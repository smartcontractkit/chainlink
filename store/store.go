package store

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/models"
)

type Store struct {
	*models.ORM
	Config   Config
	KeyStore *KeyStore
	sigs     chan os.Signal
	Exiter   func(int)
}

func NewStore(config Config) *Store {
	err := os.MkdirAll(config.RootDir, os.FileMode(0700))
	if err != nil {
		logger.Fatal(err)
	}
	orm := models.NewORM(config.RootDir)
	store := &Store{
		ORM:      orm,
		Config:   config,
		KeyStore: NewKeyStore(config.KeysDir()),
		Exiter:   os.Exit,
	}
	return store
}

func (self *Store) Start() {
	self.sigs = make(chan os.Signal, 1)
	signal.Notify(self.sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-self.sigs
		self.Close()
		self.Exiter(1)
	}()
}
