package aptos

import (
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/stretchr/testify/assert"
)

const (
	mockMCMSAddress = "0x3f20aa841a0eb5c038775bdb868924770df1ce377cc0013b3ba4ac9fd69a4f90"
	mockCCIPAddress = "0xdccda7ae3917747b973a9d08609c328a37f9f8b44e4e00291be2fb58ae932bac"
	mockAddress     = "0x13a9f1a109368730f2e355d831ba8fbf5942fb82321863d55de54cb4ebe5d18f"
	mockLinkAddress = "0xa"
	mockBadAddress  = "0xinvalid"

	sepChainSelector     = 11155111
	sepMockOnRampAddress = "0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59"
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

func GetMockChainContractParams(t *testing.T, chainSelector uint64) config.ChainContractParams {
	mockParsedAddress := mustParseAddress(t, mockAddress)
	mockParsedLinkAddress := mustParseAddress(t, mockLinkAddress)

	return config.ChainContractParams{
		FeeQuoterParams: config.FeeQuoterParams{
			MaxFeeJuelsPerMsg:            1000000,
			LinkToken:                    mockParsedLinkAddress,
			TokenPriceStalenessThreshold: 1000000,
			FeeTokens:                    []aptos.AccountAddress{mockParsedLinkAddress},
		},
		OffRampParams: config.OffRampParams{
			ChainSelector:                    chainSelector,
			PermissionlessExecutionThreshold: uint32(60 * 60 * 8),
			IsRMNVerificationDisabled:        []bool{false},
			SourceChainSelectors:             []uint64{sepChainSelector},
			SourceChainIsEnabled:             []bool{true},
			SourceChainsOnRamp:               [][]byte{common.HexToAddress(sepMockOnRampAddress).Bytes()},
		},
		OnRampParams: config.OnRampParams{
			ChainSelector:  chainSelector,
			AllowlistAdmin: mockParsedAddress,
			FeeAggregator:  mockParsedAddress,
		},
	}
}

func getMockMCMSConfig(t *testing.T) types.MCMSWithTimelockConfigV2 {
	return types.MCMSWithTimelockConfigV2{
		Canceller:        proposalutils.SingleGroupMCMSV2(t),
		Proposer:         proposalutils.SingleGroupMCMSV2(t),
		Bypasser:         proposalutils.SingleGroupMCMSV2(t),
		TimelockMinDelay: big.NewInt(0),
	}
}
