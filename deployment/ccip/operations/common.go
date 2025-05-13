package operations

import (
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type OpDependencies struct {
	Env          deployment.Environment
	CurrentState stateview.CCIPOnChainState
}
