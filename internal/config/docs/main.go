package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/internal/config"
)

func main() {
	s, err := config.GenerateDocs()
	fmt.Print(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid config docs: %v\n", err)
		os.Exit(1)
	}
}
