package main

import (
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	configlib "github.com/smartcontractkit/chainlink-go/config"
	"github.com/smartcontractkit/chainlink-go/web"
)

func main() {
	config := configlib.New()
	logger.SetLoggerDir(config.RootDir)
	defer logger.Sync()
	store := services.NewStore(config)

	services.Authenticate(store)
	r := web.Router(store)
	err := store.Start()
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Close()
	logger.Fatal(r.Run())
}
