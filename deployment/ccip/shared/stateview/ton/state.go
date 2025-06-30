package ton

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"

	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/xssnick/tonutils-go/address"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// TonCCIPChainState holds a Go binding for all the currently deployed CCIP contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type CCIPChainState struct {
	LinkTokenAddress address.Address
	CCIPAddress      address.Address
	OffRamp          address.Address
	Router           address.Address

	// dummy receiver address
	ReceiverAddress address.Address
}

func SaveOnchainState(chainSelector uint64, state CCIPChainState, e cldf.Environment) error {
	ab := e.ExistingAddresses
	if !state.LinkTokenAddress.IsAddrNone() {
		err := ab.Save(chainSelector, state.LinkTokenAddress.String(), cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	if !state.CCIPAddress.IsAddrNone() {
		err := ab.Save(chainSelector, state.CCIPAddress.String(), cldf.NewTypeAndVersion(shared.TonCCIP, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	if !state.ReceiverAddress.IsAddrNone() {
		err := ab.Save(chainSelector, state.ReceiverAddress.String(), cldf.NewTypeAndVersion(shared.TonReceiver, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	if !state.OffRamp.IsAddrNone() {
		err := ab.Save(chainSelector, state.OffRamp.String(), cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	if !state.Router.IsAddrNone() {
		err := ab.Save(chainSelector, state.Router.String(), cldf.NewTypeAndVersion(shared.Router, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadOnchainState(e cldf.Environment) (map[uint64]CCIPChainState, error) {
	chains := make(map[uint64]CCIPChainState)
	for chainSelector, chain := range e.BlockChains.TonChains() {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return chains, err
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}
		chainState, err := loadChainState(chain, addresses)
		if err != nil {
			return chains, err
		}
		chains[chainSelector] = chainState
	}
	return chains, nil
}

// loadChainState Loads all state for a TonChain into state
func loadChainState(chain cldf_ton.Chain, addressTypes map[string]cldf.TypeAndVersion) (CCIPChainState, error) {
	_ = chain // TODO: Use chain to access the client if needed
	state := CCIPChainState{}

	// Most programs upgraded in place, but some are not so we always want to
	// load the latest version
	versions := make(map[cldf.ContractType]semver.Version)
	for addressStr, tvStr := range addressTypes {
		address, err := address.ParseAddr(addressStr)
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
