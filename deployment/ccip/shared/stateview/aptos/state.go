package aptos

import (
	"errors"
	"fmt"

	aptos2 "github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	deployment2 "github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

type AptosCCIPChainState struct {
	MCMSAddress      aptos2.AccountAddress
	CCIPAddress      aptos2.AccountAddress
	LinkTokenAddress aptos2.AccountAddress

	// Test contracts
	TestRouterAddress aptos2.AccountAddress
	ReceiverAddress   aptos2.AccountAddress
}

// LoadOnchainStateAptos loads chain state for Aptos chains from env
func LoadOnchainStateAptos(env deployment2.Environment) (map[uint64]AptosCCIPChainState, error) {
	aptosChains := make(map[uint64]AptosCCIPChainState)
	for chainSelector := range env.AptosChains {
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, deployment.ErrChainNotFound) {
				return aptosChains, err
			}
			addresses = make(map[string]deployment.TypeAndVersion)
		}
		chainState, err := loadAptosChainStateFromAddresses(addresses)
		if err != nil {
			return aptosChains, err
		}
		aptosChains[chainSelector] = chainState
	}
	return aptosChains, nil
}

func loadAptosChainStateFromAddresses(addresses map[string]deployment.TypeAndVersion) (AptosCCIPChainState, error) {
	chainState := AptosCCIPChainState{}
	for addrStr, typeAndVersion := range addresses {
		// Parse address
		address := &aptos2.AccountAddress{}
		err := address.ParseStringRelaxed(addrStr)
		if err != nil {
			return chainState, fmt.Errorf("failed to parse address %s for %s: %w", addrStr, typeAndVersion.Type, err)
		}
		// Set address based on type
		switch typeAndVersion.Type {
		case shared.AptosMCMSType:
			chainState.MCMSAddress = *address
		case shared.AptosCCIPType:
			chainState.CCIPAddress = *address
		case types.LinkToken:
			chainState.LinkTokenAddress = *address
		case shared.AptosReceiverType:
			chainState.ReceiverAddress = *address
		}
	}
	return chainState, nil
}

func GetOfframpDynamicConfig(c deployment2.AptosChain, ccipAddress aptos2.AccountAddress) (module_offramp.DynamicConfig, error) {
	offrampBind := ccip_offramp.Bind(ccipAddress, c.Client)
	return offrampBind.Offramp().GetDynamicConfig(&bind.CallOpts{})
}
