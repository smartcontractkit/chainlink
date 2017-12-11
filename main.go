package main

import (
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/web"
)

func main() {
	config := services.NewConfig()
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
