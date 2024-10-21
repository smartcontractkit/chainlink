package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/changeset"
)

var ccipHomeDeploy = &cobra.Command{
	Use:   "deploy-ccip-home",
	Short: "deploy CCIP contracts on home chain",
	Long: `Deploys Capability Registry, CCIP Home and RMNHome contracts on home chain. This is required for starting chainlink nodes
with Capability Registry enabled.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cribEnv == nil {
			return fmt.Errorf("cribEnv is nil")
		}
		// locate the home chain
		homeChainSel, err := cribEnvConfig.HomeChainSelector()
		if err != nil {
			return err
		}
		if homeChainSel == 0 {
			return fmt.Errorf("homeChainSel should not be 0")
		}
		changeSet, err := changeset.CapRegChangeSet(*cribEnv, homeChainSel)
		if err != nil {
			return err
		}
		addressesForChain, err := changeSet.AddressBook.AddressesForChain(homeChainSel)
		if err != nil {
			return err
		}
		for addr, typeAndVersion := range addressesForChain {
			if typeAndVersion.Type == ccipdeployment.CapabilitiesRegistry {
				fmt.Printf("CapReg: %s\n", addr)
			} else if typeAndVersion.Type == ccipdeployment.CCIPHome {
				fmt.Printf("CCIPHome: %s\n", addr)
			} else if typeAndVersion.Type == ccipdeployment.RMNHome {
				fmt.Printf("RMNHome: %s\n", addr)
			} else {
				return fmt.Errorf("unknown contract type: %s", typeAndVersion.Type)
			}
		}
		// save address book
		addrs, err := changeSet.AddressBook.Addresses()
		if err != nil {
			return err
		}
		addrBytes, err := json.MarshalIndent(addrs, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile("address_book.json", addrBytes, 0644)
	},
}
