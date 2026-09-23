package ton

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"

	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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
	OffRamp          address.Address
	Router           address.Address
	OnRamp           address.Address
	FeeQuoter        address.Address

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
	if !state.OnRamp.IsAddrNone() {
		err := ab.Save(chainSelector, state.OnRamp.String(), cldf.NewTypeAndVersion(shared.OnRamp, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	if !state.FeeQuoter.IsAddrNone() {
		err := ab.Save(chainSelector, state.FeeQuoter.String(), cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_6_0))
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadOnchainState(e cldf.Environment) (map[uint64]CCIPChainState, error) {
	chains := make(map[uint64]CCIPChainState)
	if e.DataStore == nil {
		return chains, fmt.Errorf("TON state loading requires an environment datastore")
	}
	for chainSelector, chain := range e.BlockChains.TonChains() {
		refs := e.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
		chainState, err := loadChainState(chain, refs)
		if err != nil {
			return chains, err
		}
		chains[chainSelector] = chainState
	}
	return chains, nil
}

// loadChainState loads all state for a TON chain from datastore refs.
func loadChainState(chain cldf_ton.Chain, refs []datastore.AddressRef) (CCIPChainState, error) {
	_ = chain // TODO: Use chain to access the client if needed
	state := CCIPChainState{}

	// Most programs upgraded in place, but some are not so we always want to
	// load the latest version
	versions := make(map[cldf.ContractType]semver.Version)
	for _, ref := range refs {
		if ref.Version == nil {
			return state, fmt.Errorf("datastore ref for %s has no version", ref.Address)
		}
		address, err := address.ParseAddr(ref.Address)
		if err != nil {
			return state, err
		}
		contractType := cldf.ContractType(ref.Type)

		switch contractType {
		case commontypes.LinkToken:
			state.LinkTokenAddress = *address
		case shared.TonReceiver:
			state.ReceiverAddress = *address
		case shared.OffRamp:
			state.OffRamp = *address
		case shared.Router:
			state.Router = *address
		case shared.OnRamp:
			state.OnRamp = *address
		case shared.FeeQuoter:
			state.FeeQuoter = *address
		default:
			log.Warn().Str("address", ref.Address).Str("type", string(contractType)).Msg("Unknown TON address type")
			continue
		}

		existingVersion, ok := versions[contractType]
		if ok {
			log.Warn().Str("existingVersion", existingVersion.String()).Str("type", string(contractType)).Msg("Duplicate address type found")
		}
		versions[contractType] = *ref.Version
	}

	return state, nil
}
