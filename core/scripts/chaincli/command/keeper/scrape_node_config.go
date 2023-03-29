package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// scrapeNodes represents the command to run the service
var scrapeNodes = &cobra.Command{
	Use:   "scrape-node-config",
	Short: "ScrapeNodes automation OCR2 node configs",
	Long:  `This command scrape OCR2 node configs`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewBaseHandler(cfg)
		hdlr.ScrapeNodes()
	},
}
