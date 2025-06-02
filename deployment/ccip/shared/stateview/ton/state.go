package ton

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"

	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	tonaddress "github.com/xssnick/tonutils-go/address"
)

// TonCCIPChainState holds a Go binding for all the currently deployed CCIP contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type CCIPChainState struct {
	LinkTokenAddress tonaddress.Address
	CCIPAddress      tonaddress.Address
	OffRamp          tonaddress.Address
	Router           tonaddress.Address

	// dummy receiver address
	ReceiverAddress tonaddress.Address
}

func SaveOnchainStateTon(chainSelector uint64, tonState CCIPChainState, e cldf.Environment) error {
	ab := e.ExistingAddresses
	if !tonState.LinkTokenAddress.IsAddrNone() {
		err := ab.Save(chainSelector, tonState.LinkTokenAddress.String(), cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	if !tonState.CCIPAddress.IsAddrNone() {
		err := ab.Save(chainSelector, tonState.CCIPAddress.String(), cldf.NewTypeAndVersion(shared.TonCCIP, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	if !tonState.ReceiverAddress.IsAddrNone() {
		err := ab.Save(chainSelector, tonState.ReceiverAddress.String(), cldf.NewTypeAndVersion(shared.TonReceiver, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	if !tonState.OffRamp.IsAddrNone() {
		err := ab.Save(chainSelector, tonState.OffRamp.String(), cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	if !tonState.Router.IsAddrNone() {
		err := ab.Save(chainSelector, tonState.Router.String(), cldf.NewTypeAndVersion(shared.Router, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadOnchainStateTon(e cldf.Environment) (map[uint64]CCIPChainState, error) {
	tonChains := make(map[uint64]CCIPChainState)
	for chainSelector, chain := range e.BlockChains.TonChains() {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return tonChains, err
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}
		chainState, err := LoadChainStateTon(chain, addresses)
		if err != nil {
			return tonChains, err
		}
		tonChains[chainSelector] = chainState
	}
	return tonChains, nil
}

// LoadChainStateTon Loads all state for a TonChain into state
func LoadChainStateTon(chain cldf_ton.Chain, addresses map[string]cldf.TypeAndVersion) (CCIPChainState, error) {
	state := CCIPChainState{}

	// Most programs upgraded in place, but some are not so we always want to
	// load the latest version
	versions := make(map[cldf.ContractType]semver.Version)
	for addressStr, tvStr := range addresses {
		address, err := tonaddress.ParseAddr(addressStr)
		if err != nil {
			return state, err
		}

		switch tvStr.Type {
		case commontypes.LinkToken:
			state.LinkTokenAddress = *address
		case shared.TonCCIP:
			state.CCIPAddress = *address
		case shared.TonReceiver:
			state.ReceiverAddress = *address
		case shared.OffRamp:
			state.OffRamp = *address
		case shared.Router:
			state.Router = *address
		default:
			log.Warn().Str("address", addressStr).Str("type", string(tvStr.Type)).Msg("Unknown TON address type")
			continue
		}

		existingVersion, ok := versions[tvStr.Type]
		if ok {
			log.Warn().Str("existingVersion", existingVersion.String()).Str("type", string(tvStr.Type)).Msg("Duplicate address type found")
		}
		versions[tvStr.Type] = tvStr.Version
	}

	return state, nil
}
