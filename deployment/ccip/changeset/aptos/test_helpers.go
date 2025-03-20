package aptos

import (
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/stretchr/testify/assert"
)

const (
	mockMCMSAddress = "0x3f20aa841a0eb5c038775bdb868924770df1ce377cc0013b3ba4ac9fd69a4f90"
	mockCCIPAddress = "0xdccda7ae3917747b973a9d08609c328a37f9f8b44e4e00291be2fb58ae932bac"
	mockAddress     = "0x13a9f1a109368730f2e355d831ba8fbf5942fb82321863d55de54cb4ebe5d18f"
	mockBadAddress  = "0xinvalid"
)

func getTestAddressBook(addrByChain map[uint64]map[string]deployment.TypeAndVersion) deployment.AddressBook {
	ab := deployment.NewMemoryAddressBook()
	for chain, addrTypeAndVersion := range addrByChain {
		for addr, typeAndVersion := range addrTypeAndVersion {
			ab.Save(chain, addr, typeAndVersion)
		}
	}
	return ab
}

func mustParseAddress(t *testing.T, addr string) aptos.AccountAddress {
	t.Helper()
	var address aptos.AccountAddress
	err := address.ParseStringRelaxed(addr)
	assert.NoError(t, err)
	return address
}

// TODO: Unused function, to be used in changesets tests
func getMockChainContractParams(t *testing.T, chainSelector uint64) ChainContractParams {
	mockParsedAddress := mustParseAddress(t, mockAddress)
	return ChainContractParams{
		FeeQuoterParams: FeeQuoterParams{
			MaxFeeJuelsPerMsg:            100,
			LinkToken:                    mockParsedAddress,
			TokenPriceStalenessThreshold: 100,
			FeeTokens:                    []aptos.AccountAddress{mockParsedAddress},
		},
		OffRampParams: OffRampParams{
			ChainSelector:                    chainSelector,
			PermissionlessExecutionThreshold: 7,
			IsRMNVerificationDisabled:        false,
		},
		OnRampParams: OnRampParams{
			ChainSelector:             chainSelector,
			AllowlistAdmin:            mockParsedAddress,
			DestChainSelectors:        []uint64{},
			DestChainEnabled:          []bool{},
			DestChainAllowlistEnabled: []bool{},
		},
	}
}
