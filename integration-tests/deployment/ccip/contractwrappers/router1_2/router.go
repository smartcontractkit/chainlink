package router1_2

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

type Router struct {
	*router.Router
}

func (r Router) GetOffRamps(opts *bind.CallOpts) ([]stirng, error) {
	offRamps, err := r.Router.GetOffRamps(opts)
	if err != nil {
		return nil, err
	}
	converted := make([]view.RouterOffRamp, len(offRamps))
	for i, offRamp := range offRamps {
		converted[i] = view.RouterOffRamp{
			SourceChainSelector: offRamp.SourceChainSelector,
			OffRamp:             offRamp.OffRamp,
		}
	}
	return converted, nil
}

func New(r *router.Router) *Router {
	return &Router{r}
}
