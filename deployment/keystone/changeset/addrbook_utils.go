package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ccipowner "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	capReg "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	feeds_consumer "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	keystoneForwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3Capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
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
		internal.CapabilitiesRegistry,
		capReg.NewCapabilitiesRegistry,
	)
}

// ocr3FromAddrBook retrieves OCR3Capability contracts for the given chain.
func ocr3FromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*ocr3Capability.OCR3Capability, error) {
	return getContractsFromAddrBook[ocr3Capability.OCR3Capability](
		addrBook,
		chain,
		internal.OCR3Capability,
		ocr3Capability.NewOCR3Capability,
	)
}

// forwardersFromAddrBook retrieves KeystoneForwarder contracts for the given chain.
func forwardersFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*keystoneForwarder.KeystoneForwarder, error) {
	return getContractsFromAddrBook[keystoneForwarder.KeystoneForwarder](
		addrBook,
		chain,
		internal.KeystoneForwarder,
		keystoneForwarder.NewKeystoneForwarder,
	)
}

// feedsConsumersFromAddrBook retrieves FeedsConsumer contracts for the given chain.
func feedsConsumersFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*feeds_consumer.KeystoneFeedsConsumer, error) {
	return getContractsFromAddrBook[feeds_consumer.KeystoneFeedsConsumer](
		addrBook,
		chain,
		internal.FeedConsumer,
		feeds_consumer.NewKeystoneFeedsConsumer,
	)
}

// proposersFromAddrBook retrieves ManyChainMultiSig proposer contracts for the given chain.
func proposersFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*ccipowner.ManyChainMultiSig, error) {
	return getContractsFromAddrBook[ccipowner.ManyChainMultiSig](
		addrBook,
		chain,
		internal.ProposerManyChainMultiSig,
		ccipowner.NewManyChainMultiSig,
	)
}

// timelocksFromAddrBook retrieves RBACTimelock contracts for the given chain.
func timelocksFromAddrBook(addrBook deployment.AddressBook, chain deployment.Chain) ([]*ccipowner.RBACTimelock, error) {
	return getContractsFromAddrBook[ccipowner.RBACTimelock](
		addrBook,
		chain,
		internal.RBACTimelock,
		ccipowner.NewRBACTimelock,
	)
}
