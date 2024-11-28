package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	ccipowner "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone"

	capReg "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	keystoneForwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	ocr3Capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
)

// contractConstructor is a function type that takes an address and a client,
// returning the contract instance and an error.
type contractConstructor[T any] func(address common.Address, client bind.ContractBackend) (*T, error)

// getContractFromAddrBook is a generic function to retrieve a single contract instance
// of a specific type from the address book. It returns the first matching instance found
// or an error if none are found.
func getContractFromAddrBook[T any](
	addrBook deployment.AddressBook,
	chain deployment.Chain,
	desiredType deployment.ContractType,
	constructor contractConstructor[T],
) (*T, error) {
	chainAddresses, err := addrBook.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chain.Selector, err)
	}

	for addr, typeAndVersion := range chainAddresses {
		if typeAndVersion.Type == desiredType {
			address := common.HexToAddress(addr)
			contractInstance, err := constructor(address, chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to construct %s at address %s: %w", desiredType, addr, err)
			}
			return contractInstance, nil
		}
	}

	return nil, fmt.Errorf("no %s found for chain %d", desiredType, chain.Selector)
}

// capRegistryFromAddrBook returns the CapabilitiesRegistry contract for the given chain and address book.
func capRegistryFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) (*capReg.CapabilitiesRegistry, error) {
	return getContractFromAddrBook[capReg.CapabilitiesRegistry](
		addrBook,
		chain,
		keystone.CapabilitiesRegistry,
		capReg.NewCapabilitiesRegistry,
	)
}

// ocr3FromAddrBook returns the OCR3Capability contract for the given chain and address book.
func ocr3FromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) (*ocr3Capability.OCR3Capability, error) {
	return getContractFromAddrBook[ocr3Capability.OCR3Capability](
		addrBook,
		chain,
		keystone.OCR3Capability,
		ocr3Capability.NewOCR3Capability,
	)
}

// forwarderFromAddrBook returns the KeystoneForwarder contract for the given chain and address book.
func forwarderFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) (*keystoneForwarder.KeystoneForwarder, error) {
	return getContractFromAddrBook[keystoneForwarder.KeystoneForwarder](
		addrBook,
		chain,
		keystone.KeystoneForwarder,
		keystoneForwarder.NewKeystoneForwarder,
	)
}

// proposerFromAddrBook returns the ManyChainMultiSig proposer contract for the given chain and address book.
func proposerFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) (*ccipowner.ManyChainMultiSig, error) {
	return getContractFromAddrBook[ccipowner.ManyChainMultiSig](
		addrBook,
		chain,
		keystone.ProposerManyChainMultiSig,
		ccipowner.NewManyChainMultiSig,
	)
}

// timelockFromAddrBook returns the RBACTimelock contract for the given chain and address book.
func timelockFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) (*ccipowner.RBACTimelock, error) {
	return getContractFromAddrBook[ccipowner.RBACTimelock](
		addrBook,
		chain,
		keystone.RBACTimelock,
		ccipowner.NewRBACTimelock,
	)
}
