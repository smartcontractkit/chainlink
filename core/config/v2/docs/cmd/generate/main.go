// Docs prints core node documentation and/or a list of errors.
// The docs are Markdown generated from Toml - see config.GenerateConfig & config.GenerateSecrets.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/smartcontractkit/chainlink/core/config/v2/docs"
)

var outDir = flag.String("o", "", "output directory")

func main() {
	flag.Parse()

	c, err := docs.GenerateConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid config docs: %v\n", err)
		os.Exit(1)
	}
	s, err := docs.GenerateSecrets()
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid secrets docs: %v\n", err)
		os.Exit(1)
	}
	if err = os.WriteFile(path.Join(*outDir, "CONFIG.md"), []byte(c), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write config docs: %v\n", err)
		os.Exit(1)
	}
	if err = os.WriteFile(path.Join(*outDir, "SECRETS.md"), []byte(s), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write secrets docs: %v\n", err)
		os.Exit(1)
	}
}
