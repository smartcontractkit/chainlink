package main

import (
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/web"
)

func main() {
	config := store.NewConfig()
	logger.SetLoggerDir(config.RootDir)
	defer logger.Sync()
	app := services.NewApplication(config)

	services.Authenticate(app.Store)
	r := web.Router(app)

	if err := app.Start(); err != nil {
		logger.Fatal(err)
	}
	defer app.Stop()
	logger.Fatal(r.Run())
}
