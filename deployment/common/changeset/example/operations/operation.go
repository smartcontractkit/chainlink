package example

import (
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var DeployLinkOp = operations.NewOperation(
	"deploy-link-token-op",
	semver.MustParse("1.0.0"),
	"Deploy LINK Contract Operation",
	func(b operations.Bundle, deps EthereumDeps, input operations.EmptyInput) (common.Address, error) {
		linkToken, err := cldf.DeployContract[*link_token.LinkToken](b.Logger, deps.Chain, deps.AB,
			func(chain deployment.Chain) cldf.ContractDeploy[*link_token.LinkToken] {
				linkTokenAddr, tx, linkToken, err2 := link_token.DeployLinkToken(
					chain.DeployerKey,
					chain.Client,
				)
				return cldf.ContractDeploy[*link_token.LinkToken]{
					Address:  linkTokenAddr,
					Contract: linkToken,
					Tx:       tx,
					Tv:       deployment.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0),
					Err:      err2,
				}
			})
		if err != nil {
			b.Logger.Errorw("Failed to deploy link token", "chain", deps.Chain.String(), "err", err)
			return common.Address{}, err
		}
		return linkToken.Address, nil
	})

type GrantMintRoleConfig struct {
	ContractAddress common.Address
	To              common.Address
}

var GrantMintOp = operations.NewOperation(
	"grant-mint-role-op",
	semver.MustParse("1.0.0"),
	"Grant Mint Role Operation",
	func(b operations.Bundle, deps EthereumDeps, input GrantMintRoleConfig) (any, error) {
		contract, err := link_token.NewLinkToken(input.ContractAddress, deps.Chain.Client)
		if err != nil {
			return nil, err
		}
		tx, err := contract.GrantMintRole(deps.Auth, input.To)
		_, err = deployment.ConfirmIfNoError(deps.Chain, tx, err)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

type MintLinkConfig struct {
	Amount          *big.Int
	To              common.Address
	ContractAddress common.Address
}

var MintLinkOp = operations.NewOperation(
	"mint-link-token-op",
	semver.MustParse("1.0.0"),
	"Mint LINK Operation",
	func(ctx operations.Bundle, deps EthereumDeps, input MintLinkConfig) (any, error) {
		contract, err := link_token.NewLinkToken(input.ContractAddress, deps.Chain.Client)
		if err != nil {
			return nil, err
		}
		tx, err := contract.Mint(deps.Auth, input.To, input.Amount)
		_, err = deployment.ConfirmIfNoError(deps.Chain, tx, err)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
