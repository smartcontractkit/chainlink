package examples

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	creenv "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/environment"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/deploy"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/pkg/fake"
)

var rpcURL string

var DeployPermissionlessFeedsConsumerCmd = &cobra.Command{
	Use:   "deploy-permissionless-feeds-consumer",
	Short: "Deploy a Permissionless Feeds Consumer contract",
	Long:  `Deploy a Permissionless Feeds Consumer contract to the specified blockchain network using the provided RPC URL.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resolvedRPC, err := resolveRPC(cmd, rpcURL)
		if err != nil {
			return err
		}

		address, deployErr := deploy.PermissionlessFeedsConsumer(resolvedRPC)
		if deployErr != nil {
			return errors.Wrap(deployErr, "failed to deploy Permissionless Feeds Consumer contract")
		}

		fmt.Printf("\033[35m\nDeployed Permissionless Feeds Consumer contract to: %s\033[0m\n\n", address.Hex())

		return nil
	},
}

var DeployBalanceReaderCmd = &cobra.Command{
	Use:   "deploy-balance-reader",
	Short: "Deploy a Balance Reader contract",
	Long:  `Deploy a Balance Reader contract to the specified blockchain network using the provided RPC URL.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resolvedRPC, err := resolveRPC(cmd, rpcURL)
		if err != nil {
			return err
		}

		address, deployErr := deploy.BalanceReader(resolvedRPC)
		if deployErr != nil {
			return errors.Wrap(deployErr, "failed to deploy Balance Reader contract")
		}

		fmt.Printf("\033[35m\nDeployed Balance Reader contract to: %s\033[0m\n\n", address.Hex())

		return nil
	},
}

var DeployFakePriceProviderCmd = &cobra.Command{
	Use:   "deploy-fake-price-provider",
	Short: "Deploy a fake price provider locally",
	Long:  `Deploy a fake price provider service locally that can be used for testing workflows. Returns the URL where the service is accessible.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		authKey, _ := cmd.Flags().GetString("auth-key")
		port, _ := cmd.Flags().GetInt("port")
		feedIDs, _ := cmd.Flags().GetStringSlice("feed-ids")

		_, err := fake.DeployPriceProvider(authKey, port, feedIDs, "")
		if err != nil {
			return errors.Wrap(err, "failed to deploy fake price provider")
		}

		// Keep the service running
		select {}
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

func resolveRPC(cmd *cobra.Command, defaultRPC string) (string, error) {
	if cmd.Flags().Changed("rpc-url") {
		return defaultRPC, nil
	}

	resolver, err := creenv.TryLoadLocalCREStateResolver()
	if err != nil {
		return "", errors.Wrap(err, "failed to load local CRE state")
	}

	if resolver == nil {
		return defaultRPC, nil
	}

	resolvedRPC, err := resolver.RegistryRPC()
	if err != nil {
		return defaultRPC, nil
	}

	return resolvedRPC, nil
}

func init() {
	DeployPermissionlessFeedsConsumerCmd.Flags().StringVarP(&rpcURL, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	DeployBalanceReaderCmd.Flags().StringVarP(&rpcURL, "rpc-url", "r", "http://localhost:8545", "RPC URL")

	DeployFakePriceProviderCmd.Flags().String("auth-key", "Bearer test-auth-key", "Authentication key for the price provider")
	DeployFakePriceProviderCmd.Flags().Int("port", 80, "Port to run the fake price provider on")
	DeployFakePriceProviderCmd.Flags().StringSlice("feed-ids", []string{"0x1234567890123456789012345678901234567890123456789012345678901234"}, "Feed IDs to provide prices for")

	contractsCmd.AddCommand(DeployPermissionlessFeedsConsumerCmd)
	contractsCmd.AddCommand(DeployBalanceReaderCmd)

	ExamplesCmd.AddCommand(contractsCmd)
	ExamplesCmd.AddCommand(DeployFakePriceProviderCmd)
}
