package sui

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

type CCIPChainState struct {
	CCIPRouterAddress          string
	CCIPAddress                string
	CCIPObjectRef              string
	MCMsAddress                string
	TokenPoolAddress           string
	LockReleaseAddress         string
	LockReleaseStateId         string
	FeeQuoterCapId             string
	OnRampAddress              string
	OnRampStateObjectId        string
	OffRampAddress             string
	OffRampOwnerCapId          string
	OffRampStateObjectId       string
	LinkTokenAddress           string
	LinkTokenCoinMetadataId    string
	LinkTokenTreasuryCapId     string
	CCIPBurnMintTokenPool      string
	CCIPBurnMintTokenPoolState string
}

// LoadOnchainStatesui loads chain state for sui chains from env
func LoadOnchainStatesui(env cldf.Environment) (map[uint64]CCIPChainState, error) {
	rawChains := env.BlockChains.SuiChains()
	suiChains := make(map[uint64]CCIPChainState)

	for chainSelector := range rawChains {
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty state
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chainSelector, err)
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}

		chainState, err := loadsuiChainStateFromAddresses(addresses)
		if err != nil {
			return nil, fmt.Errorf("failed to load chain state for chain %d: %w", chainSelector, err)
		}

		suiChains[chainSelector] = chainState
	}

	return suiChains, nil
}

func loadsuiChainStateFromAddresses(addresses map[string]cldf.TypeAndVersion) (CCIPChainState, error) {
	chainState := CCIPChainState{}
	for addr, typeAndVersion := range addresses {
		// Parse addresss

		switch typeAndVersion.Type {
		case shared.SuiCCIPRouterType:
			chainState.CCIPRouterAddress = addr

		case shared.SuiCCIPType:
			chainState.CCIPAddress = addr

		case shared.SuiLockReleaseTPType:
			chainState.LockReleaseAddress = addr

		case shared.SuiLockReleaseTPStateType:
			chainState.LockReleaseStateId = addr

		case shared.SuiMCMSType:
			chainState.MCMsAddress = addr

		case shared.SuiTokenPoolType:
			chainState.TokenPoolAddress = addr

		case shared.SuiCCIPObjectRefType:
			chainState.CCIPObjectRef = addr

		case shared.SuiFeeQuoterCapType:
			chainState.FeeQuoterCapId = addr

		case shared.SuiOnRampType:
			chainState.OnRampAddress = addr

		case shared.SuiOnRampStateObjectIdType:
			chainState.OnRampStateObjectId = addr

		case shared.SuiOffRampType:
			chainState.OffRampAddress = addr

		case shared.SuiOffRampStateObjectIdType:
			chainState.OffRampStateObjectId = addr

		case shared.SuiOffRampOwnerCapObjectIdType:
			chainState.OffRampOwnerCapId = addr

		case shared.SuiLinkTokenType:
			chainState.LinkTokenAddress = addr

		case shared.SuiLinkTokenObjectMetadataId:
			chainState.LinkTokenCoinMetadataId = addr

		case shared.SuiLinkTokenTreasuryCapId:
			chainState.LinkTokenTreasuryCapId = addr

		case shared.SuiBnMTokenPoolType:
			chainState.CCIPBurnMintTokenPool = addr

		case shared.SuiBnMTokenPoolStateType:
			chainState.CCIPBurnMintTokenPoolState = addr
		}
		// Set address based on type

	}
	return chainState, nil
}
