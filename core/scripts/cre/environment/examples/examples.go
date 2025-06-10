package examples

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/deploy"
)

var rpcURL string

var DeployPermissionlessFeedsConsumerCmd = &cobra.Command{
	Use:   "deploy-permissionless-feeds-consumer",
	Short: "Deploy a Permissionless Feeds Consumer contract",
	Long:  `Deploy a Permissionless Feeds Consumer contract to the specified blockchain network using the provided RPC URL.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, deployErr := deploy.PermissionlessFeedsConsumer(rpcURL)
		if deployErr != nil {
			return errors.Wrap(deployErr, "failed to deploy Permissionless Feeds Consumer contract")
		}

		fmt.Printf("\033[35m\nDeployed Permissionless Feeds Consumer contract to: %s\033[0m\n\n", address.Hex())

		return nil
	},
}

var contractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Deploy example contracts",
}

var ExamplesCmd = &cobra.Command{
	Use:   "examples",
	Short: "Deploy various examples",
}

func init() {
	DeployPermissionlessFeedsConsumerCmd.Flags().StringVarP(&rpcURL, "rpc-url", "r", "http://localhost:8545", "RPC URL")

	contractsCmd.AddCommand(DeployPermissionlessFeedsConsumerCmd)
	ExamplesCmd.AddCommand(contractsCmd)
}
