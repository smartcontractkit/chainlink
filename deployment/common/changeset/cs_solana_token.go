package changeset

import (
	"context"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink/deployment"

	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
)

type DeploySolanaTokenConfig struct {
	ChainSelector    uint64
	TokenSymbol      string
	TokenProgramName string
	ATAList          []string // addresses to create ATAs for
}

var _ deployment.ChangeSet[*DeploySolanaTokenConfig] = DeploySolanaToken

func DeploySolanaToken(e deployment.Environment, cfg *DeploySolanaTokenConfig) (deployment.ChangesetOutput, error) {
	// validate
	tokenprogramID, err := deployment.GetTokenProgramID(cfg.TokenProgramName)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	tokenAdminPubKey := chain.DeployerKey.PublicKey()
	mint, _ := solana.NewRandomPrivateKey()
	mintPublicKey := mint.PublicKey() // this is the token address

	instructions, err := solTokenUtil.CreateToken(
		context.Background(), tokenprogramID, mintPublicKey, tokenAdminPubKey, TokenDecimalsSolana, chain.Client, solRpc.CommitmentConfirmed,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// these are associated token accounts for the addresses in the list
	// these are the default accounts that created per (token, owner) pair
	// hence they are PDAs and dont need to be stored in the address book
	for _, ata := range cfg.ATAList {
		createATAIx, _, err := solTokenUtil.CreateAssociatedTokenAccount(
			tokenprogramID, mintPublicKey, solana.MustPublicKeyFromBase58(ata), tokenAdminPubKey)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		instructions = append(instructions, createATAIx)
	}

	err = chain.Confirm(instructions, solCommomUtil.AddSigners(mint))
	if err != nil {
		e.Logger.Errorw("Failed to confirm instructions for link token deployment", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}

	// address book update
	// i am storing using TokenSymbol but i dont know how to load state
	// from the address book given the address alone
	newAddresses := deployment.NewMemoryAddressBook()
	tv := deployment.NewTypeAndVersion(deployment.ContractType(cfg.TokenSymbol), deployment.Version1_0_0)
	err = newAddresses.Save(chain.Selector, mintPublicKey.String(), tv)
	if err != nil {
		e.Logger.Errorw("Failed to save link token", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}

	e.Logger.Infow("Deployed contract", "Contract", tv.String(), "addr", mintPublicKey.String(), "chain", chain.String())

	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

type MintSolanaTokenConfig struct {
	ChainSelector uint64
	TokenSymbol   string
	TokenProgram  string
	Amount        uint64
	ToAddressList []string
}

func MintSolanaToken(e deployment.Environment, cfg *MintSolanaTokenConfig) (deployment.ChangesetOutput, error) {
	// get chain
	chain := e.SolChains[cfg.ChainSelector]
	// get addresses
	tokenAddress, err := deployment.FindTokenAddress(e, cfg.ChainSelector, cfg.TokenSymbol)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// get token program id
	tokenprogramID, err := deployment.GetTokenProgramID(cfg.TokenProgram)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// get mint instructions
	instructions := []solana.Instruction{}
	for _, toAddress := range cfg.ToAddressList {
		toAddressBase58 := solana.MustPublicKeyFromBase58(toAddress)
		// get associated token account for toAddress
		ata, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenprogramID, tokenAddress, toAddressBase58)
		mintToI, err := solTokenUtil.MintTo(cfg.Amount, tokenprogramID, tokenAddress, ata, chain.DeployerKey.PublicKey())
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		instructions = append(instructions, mintToI)
	}
	// confirm instructions
	err = chain.Confirm(instructions)
	if err != nil {
		e.Logger.Errorw("Failed to confirm instructions for token minting", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}
