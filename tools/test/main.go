// Package main launches the tools/test CLI.
//
//go:generate go run . --sync-skills
package main

import (
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/cmd"
)

func main() {
	cmd.Execute()
}
