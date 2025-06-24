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
	}

	DeployOutput struct {
		ProgramID solana.PublicKey
	}
)

func Deploy(b operations.Bundle, deps Deps, in DeployInput) (DeployOutput, error) {
	var out DeployOutput

	b.Logger.Infof("deploying program %q, size %d, chain sel %d", in.ProgramName, in.Size, in.ChainSel)
	programID, err := deps.Chain.DeployProgram(b.Logger, cldfsol.ProgramInfo{
		Name:  in.ProgramName,
		Bytes: in.Size,
	}, false, in.Overallocate)
	if err != nil {
		return out, err
	}

	out.ProgramID = solana.MustPublicKeyFromBase58(programID)

	return out, nil
}
