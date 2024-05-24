package command

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// BootstrapNodeCmd launches a chainlink node with a bootstrap job
var BootstrapNodeCmd = &cobra.Command{
	Use:   "bootstrap [ui-port] [p2pv2-port]",
	Short: "Setup a bootstrap node.",
	Long:  `This commands launches a chainlink node inside the docker container and sets up the bootstrap job`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		baseHandler := handler.NewBaseHandler(cfg)

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatal("failed to get force flag: ", err)
		}

		uiPort, err := cmd.Flags().GetInt("ui-port")
		if err != nil {
			log.Fatal("failed to get ui-port flag: ", err)
		}

		p2pv2Port, err := cmd.Flags().GetInt("p2pv2-port")
		if err != nil {
			log.Fatal("failed to get p2pv2-port flag: ", err)
		}

		baseHandler.StartBootstrapNode(cmd.Context(), cfg.RegistryAddress, uiPort, p2pv2Port, force)
	},
}

func init() {
	BootstrapNodeCmd.Flags().BoolP("force", "f", false, "Specify if existing containers should be forcefully removed")
	BootstrapNodeCmd.Flags().Int("ui-port", 5688, "Chainlink node UI listen port")
	BootstrapNodeCmd.Flags().Int("p2pv2-port", 8000, "Chainlink node P2P listen port")
}
