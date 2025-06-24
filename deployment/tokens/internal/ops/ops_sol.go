package ops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf_chain_utils "github.com/smartcontractkit/chainlink-deployments-framework/chain/utils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

const (
	// linkTokenDecimalsSolana is the number of decimals for the Solana LINK token.
	linkTokenDecimalsSolana = 9
)

// OpSolDeployLinkToken deploys a LINK token contract on Solana.
type OpSolDeployLinkTokenDeps struct {
	Client      *rpc.Client
	ConfirmFunc func(instructions []solana.Instruction, opts ...common.TxModifier) error
}

// OpSolDeployLinkTokenInput represents the input parameters for the OpSolDeployLinkToken operation.
type OpSolDeployLinkTokenInput struct {
	// ChainSelector is the unique identifier for the chain where the operation will be executed.
	// It is used as part of the unique cache key for the report and in logging, but is not used
	// in the operation logic itself.
	ChainSelector uint64 `json:"chainSelector"`
	// TokenAdminPublicKey is the public key of the admin account for the token.
	TokenAdminPublicKey solana.PublicKey `json:"tokenAdminPublicKey"`
}

// OpSolDeployLinkTokenOutput represents the output of the OpSolDeployLinkToken operation.
type OpSolDeployLinkTokenOutput struct {
	// MintPublicKey is represents the token's address
	MintPublicKey solana.PublicKey `json:"mintPublicKey"`
	Type          string           `json:"type"`
	Version       string           `json:"version"`
}

// OpSolDeployLinkToken is an operation that deploys the LINK token contract on the Solana
// blockchain. The token is deployed as a Token2022 token with 9 decimals.
var OpSolDeployLinkToken = operations.NewOperation(
	"sol-deploy-link-token",
	semver.MustParse("1.0.0"),
	"Deploy Solana LINK Token Contract",
	func(b operations.Bundle, deps OpSolDeployLinkTokenDeps, in OpSolDeployLinkTokenInput) (OpSolDeployLinkTokenOutput, error) {
		out := OpSolDeployLinkTokenOutput{}

		chainInfo, err := cldf_chain_utils.ChainInfo(in.ChainSelector)
		if err != nil {
			b.Logger.Errorw("Failed to get chain info",
				"chainSelector", in.ChainSelector,
				"err", err,
			)

			return out, err
		}

		// Generate the publicKey of the new token mint
		mint, err := solana.NewRandomPrivateKey()
		if err != nil {
			b.Logger.Errorw("Failed to generate mint public key",
				"chainSelector", in.ChainSelector,
				"chainName", chainInfo.ChainName,
				"err", err,
			)

			return out, fmt.Errorf("failed to generate mint public key: %w", err)
		}

		mintPublicKey := mint.PublicKey() // This is the token address

		// Create the token
		instructions, err := tokens.CreateToken(
			b.GetContext(),
			solana.TokenProgramID,
			mintPublicKey,
			in.TokenAdminPublicKey,
			linkTokenDecimalsSolana,
			deps.Client,
			cldf_solana.SolDefaultCommitment,
		)
		if err != nil {
			b.Logger.Errorw("Failed to generate instructions for link token deployment",
				"chainSelector", in.ChainSelector,
				"chainName", chainInfo.ChainName,
				"err", err,
			)

			return out, fmt.Errorf("failed to generate instructions for link token deployment: %w", err)
		}

		// Confirm the transaction
		if err = deps.ConfirmFunc(instructions, common.AddSigners(mint)); err != nil {
			b.Logger.Errorw("Failed to confirm instructions for link token deployment",
				"chainSelector", in.ChainSelector,
				"chainName", chainInfo.ChainName,
				"err", err,
			)

			return out, err
		}

		return OpSolDeployLinkTokenOutput{
			MintPublicKey: mintPublicKey,
			Type:          LinkTokenTypeAndVersion1.Type.String(),
			Version:       LinkTokenTypeAndVersion1.Version.String(),
		}, nil
	})
