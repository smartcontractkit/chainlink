package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink/deployment"
)

type AptosCCIPChainState struct {
	AptosMCMSObjAddr aptos.AccountAddress
	AptosCCIPObjAddr aptos.AccountAddress
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
		objAddress := &aptos.AccountAddress{}
		err := objAddress.ParseStringRelaxed(addrStr)
		if err != nil {
			return chainState, fmt.Errorf("failed to parse address %s for %s: %w", addrStr, typeAndVersion.Type, err)
		}
		// Set address based on type
		switch typeAndVersion.Type {
		case AptosMCMSType:
			chainState.AptosMCMSObjAddr = *objAddress
		case AptosCCIPType:
			chainState.AptosCCIPObjAddr = *objAddress
		}
	}
	return chainState, nil
}
