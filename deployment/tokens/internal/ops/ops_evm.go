package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	cldf_chain_utils "github.com/smartcontractkit/chainlink-deployments-framework/chain/utils"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/link_token"
)

var (
	// LinkToken is the burn/mint link token which is now used in all new deployments.
	// https://github.com/smartcontractkit/chainlink/blob/develop/core/gethwrappers/shared/generated/link_token/link_token.go#L34
	LinkTokenTypeAndVersion1 = cldf.NewTypeAndVersion(
		"LinkToken",
		*semver.MustParse("1.0.0"),
	)
)

// OpEVMDeployLinkTokenDeps defines the dependencies to perform the OpEVMDeployLinkToken
// operation.
type OpEVMDeployLinkTokenDeps struct {
	Auth        *bind.TransactOpts
	Backend     bind.ContractBackend
	ConfirmFunc func(tx *types.Transaction) (uint64, error)
}

// OpEVMDeployLinkTokenInput represents the input parameters for the OpEVMDeployLinkToken operation.
type OpEVMDeployLinkTokenInput struct {
	// ChainSelector is the unique identifier for the chain where the operation will be executed.
	// It is used as part of the unique cache key for the report and in logging, but is not used
	// in the operation logic itself.
	ChainSelector uint64 `json:"chainSelector"`
}

// OpEvmDeployLinkTokenOutput represents the output of the OpEVMDeployLinkToken operation.
type OpEvmDeployLinkTokenOutput struct {
	Address common.Address `json:"address"`
	Type    string         `json:"type"`
	Version string         `json:"version"`
}

// OpEVMDeployLinkToken is an operation that deploys the burn/mint LINK token
// contract on an EVM-compatible blockchain.
var OpEVMDeployLinkToken = operations.NewOperation(
	"evm-deploy-link-token",
	semver.MustParse("1.0.0"),
	"Deploy EVM LINK Token Contract",
	func(b operations.Bundle, deps OpEVMDeployLinkTokenDeps, in OpEVMDeployLinkTokenInput) (OpEvmDeployLinkTokenOutput, error) {
		out := OpEvmDeployLinkTokenOutput{}

		chainInfo, err := cldf_chain_utils.ChainInfo(in.ChainSelector)
		if err != nil {
			b.Logger.Errorw("Failed to get chain info",
				"chainSelector", in.ChainSelector,
				"err", err,
			)

			return out, err
		}

		// Deploy the link token
		addr, tx, _, err := link_token.DeployLinkToken(
			deps.Auth,
			deps.Backend,
		)
		if err != nil {
			b.Logger.Errorw("Failed to deploy link token",
				"chainSelector", in.ChainSelector,
				"chainName", chainInfo.ChainName,
				"err", err,
			)

			return out, err
		}

		// Confirm the transaction
		if _, err = deps.ConfirmFunc(tx); err != nil {
			b.Logger.Errorw("Failed to confirm deployment",
				"chainSelector", in.ChainSelector,
				"chainName", chainInfo.ChainName,
				"contractAddr", addr.String(),
				"err", err,
			)

			return out, err
		}

		return OpEvmDeployLinkTokenOutput{
			Address: addr,
			Type:    LinkTokenTypeAndVersion1.Type.String(),
			Version: LinkTokenTypeAndVersion1.Version.String(),
		}, nil
	})
