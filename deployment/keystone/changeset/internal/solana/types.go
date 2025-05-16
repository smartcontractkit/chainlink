package solana

import (
	"github.com/gagliardetto/solana-go"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type DeployResponse struct {
	Address solana.PublicKey
	Tx      solana.Signature
	Tv      cldf.TypeAndVersion
}

type DeployRequest struct {
	Chain cldf.SolChain
}
