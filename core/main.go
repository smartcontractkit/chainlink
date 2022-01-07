package main

import (
	"log"
	"os"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/recovery"
	"github.com/smartcontractkit/chainlink/core/wire"
)

func main() {
	recovery.ReportPanics(func() {
		client := wire.InitializeProductionClient()
		Run(client, os.Args...)
	})
}

// Run runs the CLI, providing further command instructions by default.
func Run(client *cmd.Client, args ...string) {
	app := cmd.NewApp(client)
	client.Logger.ErrorIf(app.Run(args), "Error running app")
	if err := client.Logger.Sync(); err != nil {
		log.Fatal(err)
	}
}
