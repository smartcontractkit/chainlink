package changeset

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	aptosBind "github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	module_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const (
	AptosMCMSType     deployment.ContractType = "AptosManyChainMultisig"
	AptosCCIPType     deployment.ContractType = "AptosCCIP"
	AptosReceiverType deployment.ContractType = "AptosReceiver"
)

type AptosCCIPChainState struct {
	MCMSAddress      aptos.AccountAddress
	CCIPAddress      aptos.AccountAddress
	LinkTokenAddress aptos.AccountAddress

	// Test contracts
	TestRouterAddress aptos.AccountAddress
	ReceiverAddress   aptos.AccountAddress
}

// LoadOnchainStateAptos loads chain state for Aptos chains from env
func LoadOnchainStateAptos(env deployment.Environment) (map[uint64]AptosCCIPChainState, error) {
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
		address := &aptos.AccountAddress{}
		err := address.ParseStringRelaxed(addrStr)
		if err != nil {
			return chainState, fmt.Errorf("failed to parse address %s for %s: %w", addrStr, typeAndVersion.Type, err)
		}
		// Set address based on type
		switch typeAndVersion.Type {
		case AptosMCMSType:
			chainState.MCMSAddress = *address
		case AptosCCIPType:
			chainState.CCIPAddress = *address
		case commontypes.LinkToken:
			chainState.LinkTokenAddress = *address
		case AptosReceiverType:
			chainState.ReceiverAddress = *address
		}
	}
	return chainState, nil
}

func getOfframpDynamicConfig(c deployment.AptosChain, ccipAddress aptos.AccountAddress) (module_offramp.DynamicConfig, error) {
	offrampBind := ccip_offramp.Bind(ccipAddress, c.Client)
	return offrampBind.Offramp().GetDynamicConfig(&aptosBind.CallOpts{})
}
