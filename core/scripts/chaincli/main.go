package main

import (
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/command"
	_ "github.com/smartcontractkit/chainlink/core/scripts/chaincli/command/keeper"
)

func main() {
	command.Execute()
}
