package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	module_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

type CCIPChainState struct {
	MCMSAddress      aptos.AccountAddress
	CCIPAddress      aptos.AccountAddress
	LinkTokenAddress aptos.AccountAddress

	// Test contracts
	TestRouterAddress aptos.AccountAddress
	ReceiverAddress   aptos.AccountAddress
}

// LoadOnchainStateAptos loads chain state for Aptos chains from env
func LoadOnchainStateAptos(env cldf.Environment) (map[uint64]CCIPChainState, error) {
	aptosChains := make(map[uint64]CCIPChainState)
	for chainSelector := range env.BlockChains.AptosChains() {
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return aptosChains, err
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}
		chainState, err := loadAptosChainStateFromAddresses(addresses)
		if err != nil {
			return aptosChains, err
		}
		aptosChains[chainSelector] = chainState
	}
	return aptosChains, nil
}

func loadAptosChainStateFromAddresses(addresses map[string]cldf.TypeAndVersion) (CCIPChainState, error) {
	chainState := CCIPChainState{}
	for addrStr, typeAndVersion := range addresses {
		// Parse address
		address := &aptos.AccountAddress{}
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

func GetOfframpDynamicConfig(c cldf_aptos.Chain, ccipAddress aptos.AccountAddress) (module_offramp.DynamicConfig, error) {
	offrampBind := ccip_offramp.Bind(ccipAddress, c.Client)
	return offrampBind.Offramp().GetDynamicConfig(&bind.CallOpts{})
}
