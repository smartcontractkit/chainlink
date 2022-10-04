// Docs prints core node documentation and/or a list of errors.
// The docs are Markdown generated from Toml - see config.GenerateDocs.
package main

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/core/config/v2/docs"
)

func main() {
	s, err := docs.GenerateDocs()
	fmt.Print(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid config docs: %v\n", err)
		os.Exit(1)
	}
}
