package main

import (
	"github.com/spf13/cobra"
)

type Product string

var (
	CCIP Product = "ccip"
)

func main() {
	rootCmd := &cobra.Command{}
	// add product command
	pCmd := &cobra.Command{Use: string(CCIP), Short: "CCIP HomeContract deployments"}
	pCmd.AddCommand(ccipDeploy)
	rootCmd.AddCommand(pCmd)
}

var ccipDeploy = &cobra.Command{
	Use:   "deploy-ccip-home",
	Short: "deploy CCIP contracts on home chain",
	Long: `Deploys Capability Registry, CCIP Home and RMNHome contracts on home chain. This is required for starting chainlink nodes
with Capability Registry enabled.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
