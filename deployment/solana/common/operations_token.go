package deployment_solana_common

import (
	"context"

	solana_go "github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	deployment_solana "github.com/smartcontractkit/chainlink/deployment/solana/extension"
)

type DeploySolanaTokenInput struct {
	TokenProgramName string
	TokenDecimals    uint8
}

var DeploySolanaTokenOp = deployment.NewOperation(
	"1.0.0",
	"Deploy a Solana token",
	func(ctx deployment.OpContext, deps deployment_solana.SolanaDeps, input DeploySolanaTokenInput) (*deployment_solana.SolanaTxResult, error) {

		// Build instruction
		tokenprogramID, err := solana.GetTokenProgramID(input.TokenProgramName)
		if err != nil {
			return nil, err
		}

		tokenAdminPubKey := deps.Auth.PublicKey()
		mint, _ := solana_go.NewRandomPrivateKey()
		mintPublicKey := mint.PublicKey() // this is the token address
		instructions, err := solTokenUtil.CreateToken(
			context.Background(),
			tokenprogramID,
			mintPublicKey,
			tokenAdminPubKey,
			input.TokenDecimals,
			deps.Client,
			deployment.SolDefaultCommitment,
		)
		if err != nil {
			return nil, err
		}

		// Send and Confirm
		txResult, err := deps.SendAndConfirm(
			context.Background(),
			deps.Client,
			instructions,
			deps.Auth,
			solRpc.CommitmentConfirmed,
			solCommomUtil.AddSigners(mint),
		)
		if err != nil {
			return nil, err
		}

		return &deployment_solana.SolanaTxResult{
			Address: mintPublicKey,
			Receipt: *txResult,
		}, nil
	},
)

type MintSolanaTokenInput struct {
	TokenProgram    string
	TokenPubkey     solana_go.PublicKey
	AmountToAddress map[string]uint64 // address -> amount
}

var MintSolanaTokenOp = deployment.NewOperation(
	"1.0.0",
	"Mint a Solana token",
	func(ctx deployment.OpContext, deps deployment_solana.SolanaDeps, input MintSolanaTokenInput) (*deployment_solana.SolanaTxResult, error) {
		// get chain
		// chain := e.SolChains[cfg.ChainSelector]
		// get addresses
		tokenAddress := input.TokenPubkey
		// get token program id
		tokenprogramID, err := solana.GetTokenProgramID(input.TokenProgram)
		if err != nil {
			return nil, err
		}
		// get mint instructions
		instructions := []solana_go.Instruction{}
		for toAddress, amount := range input.AmountToAddress {
			toAddressBase58 := solana_go.MustPublicKeyFromBase58(toAddress)
			// get associated token account for toAddress
			ata, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenprogramID, tokenAddress, toAddressBase58)
			mintToI, err := solTokenUtil.MintTo(amount, tokenprogramID, tokenAddress, ata, deps.Auth.PublicKey())
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, mintToI)
		}
		// Send and Confirm
		txResult, err := deps.SendAndConfirm(
			context.Background(),
			deps.Client,
			instructions,
			deps.Auth,
			solRpc.CommitmentConfirmed,
		)
		if err != nil {
			return nil, err
		}
		return &deployment_solana.SolanaTxResult{
			Address: tokenAddress,
			Receipt: *txResult,
		}, nil
	},
)

type CreateSolanaTokenATAInput struct {
	TokenPubkey  solana_go.PublicKey
	TokenProgram string
	ATAList      []string // addresses to create ATAs for
}

var CreateSolanaTokenATAOp = deployment.NewOperation(
	"1.0.0",
	"Create Solana token ATA",
	func(ctx deployment.OpContext, deps deployment_solana.SolanaDeps, input CreateSolanaTokenATAInput) (*deployment_solana.SolanaTxResult, error) {
		tokenprogramID, err := solana.GetTokenProgramID(input.TokenProgram)
		if err != nil {
			return nil, err
		}

		// create instructions for each ATA
		instructions := []solana_go.Instruction{}
		for _, ata := range input.ATAList {
			createATAIx, _, err := solTokenUtil.CreateAssociatedTokenAccount(
				tokenprogramID,
				input.TokenPubkey,
				solana_go.MustPublicKeyFromBase58(ata),
				deps.Auth.PublicKey(),
			)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, createATAIx)
		}

		// confirm instructions
		txResult, err := deps.SendAndConfirm(
			context.Background(),
			deps.Client,
			instructions,
			deps.Auth,
			solRpc.CommitmentConfirmed,
		)
		if err != nil {
			return nil, err
		}

		return &deployment_solana.SolanaTxResult{
			Address: input.TokenPubkey,
			Receipt: *txResult,
		}, nil
	},
)
