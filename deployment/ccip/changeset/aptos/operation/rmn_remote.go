package operation

import (
	"github.com/aptos-labs/aptos-go-sdk"
	fastcurseapi "github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type CurseArgs struct {
	CCIPAddress aptos.AccountAddress
	Subjects    []fastcurseapi.Subject
}

var CurseOp = operations.NewOperation(
	"curse-op",
	Version1_0_0,
	"Curses subjects with RMNRemote",
	func(b operations.Bundle, chain cldf_aptos.Chain, in CurseArgs) (sequences.OnChainOutput, error) {
		return sequences.OnChainOutput{}, nil
	},
)

var UncurseOp = operations.NewOperation(
	"uncurse-op",
	Version1_0_0,
	"Uncurses subjects with RMNRemote",
	func(b operations.Bundle, chain cldf_aptos.Chain, in CurseArgs) (sequences.OnChainOutput, error) {
		return sequences.OnChainOutput{}, nil
	},
)
