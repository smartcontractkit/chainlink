package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartcontractkit/chainlink/tools/gomod-required-updater/internal/updater"
)

var version = "dev"
var usage = `gomod-required-updater version %s

Usage:
  gomod-required-updater [flags]

Examples:
  # Update modules specified in config file
  gomod-required-updater -config modules.toml

  # Update specific modules
  gomod-required-updater -module github.com/org/repo1 -module github.com/org/repo2

  # Dry run with specific modules
  gomod-required-updater -dry-run -module github.com/org/repo1
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
		updater.NewGitOperator(),
		updater.NewSystemOperator(),
		cfg,
	)

	if err := u.Run(); err != nil {
		log.Fatal(err)
	}
}