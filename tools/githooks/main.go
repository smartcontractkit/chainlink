package main

import (
	"context"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/cmd"
)

func main() {
	if err := cmd.Execute(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
