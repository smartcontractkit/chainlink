package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ccipowner "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone"
	capReg "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	keystoneForwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	ocr3Capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
)

// contractConstructor is a function type that takes an address and a client,
// returning the contract instance and an error.
type contractConstructor[T any] func(address common.Address, client bind.ContractBackend) (*T, error)

// getContractsFromAddrBook retrieves a list of contract instances of a specified type from the address book.
// It uses the provided constructor to initialize matching contracts for the given chain.
func getContractsFromAddrBook[T any](
	addrBook deployment.AddressBook,
	chain deployment.Chain,
	desiredType deployment.ContractType,
	constructor contractConstructor[T],
) ([]*T, error) {
	chainAddresses, err := addrBook.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chain.Selector, err)
	}

	var contracts []*T
	for addr, typeAndVersion := range chainAddresses {
		if typeAndVersion.Type == desiredType {
			address := common.HexToAddress(addr)
			contractInstance, err := constructor(address, chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to construct %s at %s: %w", desiredType, addr, err)
			}
			contracts = append(contracts, contractInstance)
		}
	}

	if len(contracts) == 0 {
		return nil, fmt.Errorf("no %s found for chain %d", desiredType, chain.Selector)
	}

	return contracts, nil
}

// capRegistriesFromAddrBook retrieves CapabilitiesRegistry contracts for the given chain.
func capRegistriesFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*capReg.CapabilitiesRegistry, error) {
	return getContractsFromAddrBook[capReg.CapabilitiesRegistry](
		addrBook,
		chain,
		keystone.CapabilitiesRegistry,
		capReg.NewCapabilitiesRegistry,
	)
}

// ocr3FromAddrBook retrieves OCR3Capability contracts for the given chain.
func ocr3FromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*ocr3Capability.OCR3Capability, error) {
	return getContractsFromAddrBook[ocr3Capability.OCR3Capability](
		addrBook,
		chain,
		keystone.OCR3Capability,
		ocr3Capability.NewOCR3Capability,
	)
}

// forwardersFromAddrBook retrieves KeystoneForwarder contracts for the given chain.
func forwardersFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*keystoneForwarder.KeystoneForwarder, error) {
	return getContractsFromAddrBook[keystoneForwarder.KeystoneForwarder](
		addrBook,
		chain,
		keystone.KeystoneForwarder,
		keystoneForwarder.NewKeystoneForwarder,
	)
}

// feedsConsumersFromAddrBook retrieves FeedsConsumer contracts for the given chain.
func feedsConsumersFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*feeds_consumer.KeystoneFeedsConsumer, error) {
	return getContractsFromAddrBook[feeds_consumer.KeystoneFeedsConsumer](
		addrBook,
		chain,
		keystone.FeedConsumer,
		feeds_consumer.NewKeystoneFeedsConsumer,
	)
}

// proposersFromAddrBook retrieves ManyChainMultiSig proposer contracts for the given chain.
func proposersFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*ccipowner.ManyChainMultiSig, error) {
	return getContractsFromAddrBook[ccipowner.ManyChainMultiSig](
		addrBook,
		chain,
		keystone.ProposerManyChainMultiSig,
		ccipowner.NewManyChainMultiSig,
	)
}

// timelocksFromAddrBook retrieves RBACTimelock contracts for the given chain.
func timelocksFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*ccipowner.RBACTimelock, error) {
	return getContractsFromAddrBook[ccipowner.RBACTimelock](
		addrBook,
		chain,
		keystone.RBACTimelock,
		ccipowner.NewRBACTimelock,
	)
}
