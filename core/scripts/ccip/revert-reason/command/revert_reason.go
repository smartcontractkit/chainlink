package command

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/revert-reason/config"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/revert-reason/handler"
)

// RevertReasonCmd takes in a failed tx hash and tries to give you the reason
var RevertReasonCmd = &cobra.Command{
	Use:   "reason <tx hash or error string>",
	Short: "Revert reason for failed TX.",
	Long:  `Given a failed TX tries to find the revert reason. args = tx hex address`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		baseHandler := handler.NewBaseHandler(cfg)

		decodeFromError, err := cmd.Flags().GetBool("from-error")
		if err != nil {
			log.Fatal("failed to get withdraw flag: ", err)
		}

		if decodeFromError {
			result, err := baseHandler.RevertReasonFromErrorCodeString(args[0])
			if err != nil {
				log.Fatal("failed to decode error code string: ", err)
			}
			fmt.Print(result)
		} else {
			result, err := baseHandler.RevertReasonFromTx(args[0])
			if err != nil {
				log.Fatal("failed to decode error code string: ", err)
			}
			fmt.Print(result)
		}
	},
}
