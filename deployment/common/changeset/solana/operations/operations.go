package operations

import (
	"github.com/gagliardetto/solana-go"

	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type (
	Deps struct {
		Chain cldfsol.Chain
	}

	DeployInput struct {
		ChainSel     uint64
		ProgramName  string
		Size         int
		Overallocate bool
		// IsUpgrade writes the program binary to a buffer account instead of deploying a new
		// program. It is mutually exclusive with Overallocate.
		IsUpgrade bool
	}

	DeployOutput struct {
		// ProgramID is the address of the deployed program, or of the buffer account holding the
		// new binary when the input asked for an upgrade.
		ProgramID solana.PublicKey
	}
)

func Deploy(b operations.Bundle, deps Deps, in DeployInput) (DeployOutput, error) {
	var out DeployOutput

	b.Logger.Infof("deploying program %q, size %d, chain sel %d, is upgrade %t", in.ProgramName, in.Size, in.ChainSel, in.IsUpgrade)
	programID, err := deps.Chain.DeployProgram(b.Logger, cldfsol.ProgramInfo{
		Name:  in.ProgramName,
		Bytes: in.Size,
	}, in.IsUpgrade, in.Overallocate)
	if err != nil {
		return out, err
	}

	out.ProgramID = solana.MustPublicKeyFromBase58(programID)

	return out, nil
}
