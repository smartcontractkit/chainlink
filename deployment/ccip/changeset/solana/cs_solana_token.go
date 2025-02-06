package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/deployment"

	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
)

var _ deployment.ChangeSet[DeploySolanaTokenConfig] = DeploySolanaToken
var _ deployment.ChangeSet[MintSolanaTokenConfig] = MintSolanaToken
var _ deployment.ChangeSet[CreateSolanaTokenATAConfig] = CreateSolanaTokenATA

type DeploySolanaTokenConfig struct {
	ChainSelector    uint64
	TokenProgramName string
	TokenDecimals    uint8
}

func NewTokenInstruction(chain deployment.SolChain, cfg DeploySolanaTokenConfig) ([]solana.Instruction, solana.PrivateKey, error) {
	tokenprogramID, err := GetTokenProgramID(cfg.TokenProgramName)
	if err != nil {
		return nil, nil, err
	}
	tokenAdminPubKey := chain.DeployerKey.PublicKey()
	mint, _ := solana.NewRandomPrivateKey()
	mintPublicKey := mint.PublicKey() // this is the token address
	instructions, err := solTokenUtil.CreateToken(
		context.Background(),
		tokenprogramID,
		mintPublicKey,
		tokenAdminPubKey,
		cfg.TokenDecimals,
		chain.Client,
		deployment.SolDefaultCommitment,
	)
	if err != nil {
		return nil, nil, err
	}
	return instructions, mint, nil
}

func DeploySolanaToken(e deployment.Environment, cfg DeploySolanaTokenConfig) (deployment.ChangesetOutput, error) {
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	instructions, mint, err := NewTokenInstruction(chain, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	err = chain.Confirm(instructions, solCommomUtil.AddSigners(mint))
	if err != nil {
		e.Logger.Errorw("Failed to confirm instructions for link token deployment", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}

	newAddresses := deployment.NewMemoryAddressBook()
	tv := deployment.NewTypeAndVersion(deployment.ContractType(cfg.TokenProgramName), deployment.Version1_0_0)
	err = newAddresses.Save(cfg.ChainSelector, mint.PublicKey().String(), tv)
	if err != nil {
		e.Logger.Errorw("Failed to save link token", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}

	e.Logger.Infow("Deployed contract", "Contract", tv.String(), "addr", mint.PublicKey().String(), "chain", chain.String())

	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

// TODO: there is no validation done around if the token is already deployed
// https://smartcontract-it.atlassian.net/browse/INTAUTO-439
type MintSolanaTokenConfig struct {
	ChainSelector   uint64
	TokenProgram    string
	TokenPubkey     solana.PublicKey
	AmountToAddress map[string]uint64 // address -> amount
}

func MintSolanaToken(e deployment.Environment, cfg MintSolanaTokenConfig) (deployment.ChangesetOutput, error) {
	// get chain
	chain := e.SolChains[cfg.ChainSelector]
	// get addresses
	tokenAddress := cfg.TokenPubkey
	// get token program id
	tokenprogramID, err := GetTokenProgramID(cfg.TokenProgram)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// get mint instructions
	instructions := []solana.Instruction{}
	for toAddress, amount := range cfg.AmountToAddress {
		toAddressBase58 := solana.MustPublicKeyFromBase58(toAddress)
		// get associated token account for toAddress
		ata, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenprogramID, tokenAddress, toAddressBase58)
		mintToI, err := solTokenUtil.MintTo(amount, tokenprogramID, tokenAddress, ata, chain.DeployerKey.PublicKey())
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

type CreateSolanaTokenATAConfig struct {
	ChainSelector uint64
	TokenPubkey   solana.PublicKey
	TokenProgram  string
	ATAList       []string // addresses to create ATAs for
}

func CreateSolanaTokenATA(e deployment.Environment, cfg CreateSolanaTokenATAConfig) (deployment.ChangesetOutput, error) {
	chain := e.SolChains[cfg.ChainSelector]

	tokenprogramID, err := GetTokenProgramID(cfg.TokenProgram)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// create instructions for each ATA
	instructions := []solana.Instruction{}
	for _, ata := range cfg.ATAList {
		createATAIx, _, err := solTokenUtil.CreateAssociatedTokenAccount(
			tokenprogramID,
			cfg.TokenPubkey,
			solana.MustPublicKeyFromBase58(ata),
			chain.DeployerKey.PublicKey(),
		)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		instructions = append(instructions, createATAIx)
	}

	// confirm instructions
	err = chain.Confirm(instructions)
	if err != nil {
		e.Logger.Errorw("Failed to confirm instructions for ATA creation", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{}, nil
}
