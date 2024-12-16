package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartcontractkit/chainlink/v2/tools/gomod-local-update/internal/updater"
)

const (
	goBinaryName = "gomod-local-update"
)

var version = "dev"
var usage = fmt.Sprintf(`%s version %%s

Usage:
  cd /path/to/go/module
  %s [flags]
`, goBinaryName, goBinaryName)

func main() {
	cfg, err := updater.ParseFlags(os.Args[1:], version)
	if err != nil {
		fmt.Fprintf(os.Stderr, usage, version)
		log.Fatal(err)
	}

	if cfg.ShowVersion {
		fmt.Printf("%s version %s\n", goBinaryName, version)
		os.Exit(0)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, usage, version)
		log.Fatal(err)
	}

	u := updater.New(
		cfg,
		updater.NewSystemOperator(),
	)

	if err := u.Run(); err != nil {
		log.Fatal(err)
	}
}
