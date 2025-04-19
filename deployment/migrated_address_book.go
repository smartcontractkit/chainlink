package deployment

import "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

// Implementations have been migrated to the chainlink-deployments-framework
// Using type alias here to avoid updating all the references in the codebase.
// This file will be removed in the future once we migrate all the code
type (
	ContractType     = deployment.ContractType
	TypeAndVersion   = deployment.TypeAndVersion
	AddressBook      = deployment.AddressBook
	AddressesByChain = deployment.AddressesByChain
	AddressBookMap   = deployment.AddressBookMap
	LabelSet         = deployment.LabelSet
)

var (
	TypeAndVersionFromString     = deployment.TypeAndVersionFromString
	NewTypeAndVersion            = deployment.NewTypeAndVersion
	MustTypeAndVersionFromString = deployment.MustTypeAndVersionFromString
	NewMemoryAddressBook         = deployment.NewMemoryAddressBook
	NewMemoryAddressBookFromMap  = deployment.NewMemoryAddressBookFromMap
	SearchAddressBook            = deployment.SearchAddressBook
	AddressBookContains          = deployment.AddressBookContains
	EnsureDeduped                = deployment.EnsureDeduped
	NewLabelSet                  = deployment.NewLabelSet

	ErrInvalidChainSelector = deployment.ErrInvalidChainSelector
	ErrInvalidAddress       = deployment.ErrInvalidAddress
	ErrChainNotFound        = deployment.ErrChainNotFound
)
