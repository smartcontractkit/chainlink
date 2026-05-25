package evm

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	linkviewv10 "github.com/smartcontractkit/cld-changesets/pkg/contract/link/view/v10"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/link_token"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

type LinkTokenState struct {
	LinkToken *link_token.LinkToken
}

func (s LinkTokenState) GenerateLinkView() (linkviewv10.LinkTokenView, error) {
	if s.LinkToken == nil {
		return linkviewv10.LinkTokenView{}, errors.New("link token not found")
	}
	return linkviewv10.GenerateLinkTokenView(s.LinkToken)
}

func MaybeLoadLinkTokenChainState(chain cldf_evm.Chain, addresses map[string]cldf.TypeAndVersion) (*LinkTokenState, error) {
	state := LinkTokenState{}
	linkToken := cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0)

	// Convert map keys to a slice
	wantTypes := []cldf.TypeAndVersion{linkToken}

	// Ensure we either have the bundle or not.
	_, err := cldf.EnsureDeduped(addresses, wantTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to check link token on chain %s error: %w", chain.Name(), err)
	}

	for address, tvStr := range addresses {
		if tvStr.Type == linkToken.Type && tvStr.Version.String() == linkToken.Version.String() {
			lt, err := link_token.NewLinkToken(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.LinkToken = lt
		}
	}
	return &state, nil
}
