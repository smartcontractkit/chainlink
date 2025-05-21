package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/verify"
)

func main() {
	var contractAddress string
	var feedId string
	var rpcUrl string
	var untilSuccessful bool

	var rootCmd = &cobra.Command{
		Use:   "verify",
		Short: "Gets price from Feeds Consumer contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify.ProofOfReserve(rpcUrl, contractAddress, feedId, untilSuccessful, 10*time.Minute)
		},
	}

	rootCmd.Flags().StringVarP(&contractAddress, "contract-address", "a", "0x9A9f2CCfdE556A7E9Ff0848998Aa4a0CFD8863AE", "Contract address")
	rootCmd.Flags().StringVarP(&feedId, "feed-id", "f", "0x018e16c39e0003200000000000000000", "Feed ID")
	rootCmd.Flags().StringVarP(&rpcUrl, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	rootCmd.Flags().BoolVarP(&untilSuccessful, "until-successful", "u", true, "Run until successful")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
