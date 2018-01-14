package store

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/models"
)

type Store struct {
	*models.ORM
	Config    Config
	KeyStore  *KeyStore
	sigs      chan os.Signal
	Exiter    func(int)
	TxManager *TxManager
}

func NewStore(config Config) *Store {
	err := os.MkdirAll(config.RootDir, os.FileMode(0700))
	if err != nil {
		logger.Fatal(err)
	}
	orm := models.NewORM(config.RootDir)
	ethrpc, err := rpc.Dial(config.EthereumURL)
	if err != nil {
		logger.Fatal(err)
	}
	keyStore := NewKeyStore(config.KeysDir())
	store := &Store{
		ORM:      orm,
		Config:   config,
		KeyStore: keyStore,
		Exiter:   os.Exit,
		TxManager: &TxManager{
			Config:    config,
			EthClient: &EthClient{ethrpc},
			KeyStore:  keyStore,
			ORM:       orm,
		},
	}
	return store
}

func (s *Store) Start() {
	s.sigs = make(chan os.Signal, 1)
	signal.Notify(s.sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-s.sigs
		s.Close()
		s.Exiter(1)
	}()
}
