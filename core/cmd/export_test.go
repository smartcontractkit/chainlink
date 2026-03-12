package cmd

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/urfave/cli"
)

func NewAppWithOptsForTest(s *Shell, opts chainlink.GeneralConfigOpts) *cli.App {
	return newAppWithOpts(s, opts)
}
