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
		fmt.Fprintln(os.Stderr, "\nIf you find any unexpected errors or bugs with these checks, reach out to the #ci-cd Slack channel")
		os.Exit(1)
	}
}
