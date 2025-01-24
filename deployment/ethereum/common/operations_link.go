package deployment_common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	deployment_ethereum "github.com/smartcontractkit/chainlink/deployment/ethereum/extension"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

var newLinkToken = deployment_ethereum.NewContractCtorFn(link_token.NewLinkToken)

// Deployment Op
var DeployLinkOp = deployment_ethereum.NewEthDeployOperationFromBindingNoParams(link_token.DeployLinkToken, "v1")

// Transfer Op
type TransferLinkInput struct {
	contractAddress common.Address
	To              common.Address
	Amount          *big.Int
}

func (i TransferLinkInput) GetOrderedParams() []any {
	return []any{i.To, i.Amount}
}
func (i TransferLinkInput) Address() common.Address {
	return i.contractAddress
}

var TransferLinkOp = deployment_ethereum.NewEthOperationFromBinding[TransferLinkInput](newLinkToken, "v1", "transfer")
