package cmd

import (
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func NewAppWithOptsForTest(s *Shell, opts chainlink.GeneralConfigOpts) *cli.App {
	return newAppWithOpts(s, opts)
}
