package deployment_common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	deployment_ethereum "github.com/smartcontractkit/chainlink/deployment/ethereum/extension"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

// var newLinkToken = deployment_ethereum.NewContractCtorFn(link_token.NewLinkToken)

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

var TransferLinkOp = deployment_ethereum.NewEthOperationFromBinding[TransferLinkInput](link_token.LinkTokenMetaData, "v1", "transfer")

// Mint Op
type MintLinkInput struct {
	contractAddress common.Address
	To              common.Address
	Amount          *big.Int
}

func (i MintLinkInput) GetOrderedParams() []any {
	return []any{i.To, i.Amount}
}
func (i MintLinkInput) Address() common.Address {
	return i.contractAddress
}

var MintLinkOp = deployment_ethereum.NewEthOperationFromBinding[MintLinkInput](link_token.LinkTokenMetaData, "v1", "mint")

// Grant Mint Role Op
type GrantLinkInput struct {
	contractAddress common.Address
	To              common.Address
}

func (i GrantLinkInput) GetOrderedParams() []any {
	return []any{i.To}
}
func (i GrantLinkInput) Address() common.Address {
	return i.contractAddress
}

var GrantMintLinkOp = deployment_ethereum.NewEthOperationFromBinding[GrantLinkInput](link_token.LinkTokenMetaData, "v1", "grantMintRole")
