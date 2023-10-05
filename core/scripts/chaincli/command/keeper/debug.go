package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// jobCmd represents the command to run the service
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug an upkeep",
	Long:  `This command debugs an upkeep on the povided registry to figure out why it is not performing`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)
		hdlr.Debug(cmd.Context(), args)
	},
}
