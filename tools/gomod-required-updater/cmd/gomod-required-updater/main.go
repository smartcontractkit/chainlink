package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartcontractkit/chainlink/v2/tools/gomod-required-updater/internal/updater"
)

var version = "dev"
var usage = `gomod-required-updater version %s

Usage:
  cd /path/to/go/module
  gomod-required-updater [flags]
`

func main() {
	cfg, err := updater.ParseFlags(os.Args[1:], version)
	if err != nil {
		fmt.Fprintf(os.Stderr, usage, version)
		log.Fatal(err)
	}

	if cfg.ShowVersion {
		fmt.Printf("gomod-required-updater version %s\n", version)
		os.Exit(0)
	}

	u := updater.New(
		updater.NewModuleOperator(cfg),
		updater.NewSystemOperator(),
		cfg,
	)

	if err := u.Run(); err != nil {
		log.Fatal(err)
	}
}
