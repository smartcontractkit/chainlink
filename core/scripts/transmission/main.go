package main

import (
	"os"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func main() {
	e := helpers.SetupEnv(false)

	switch os.Args[1] {
	case "prepare-setGreeting":
		prepareSetGreeting(e)
	case "prepare-vrfRequest":
		prepareVRFRequest(e)
	default:
		panic("unrecognized subcommand: " + os.Args[1])
	}
}
