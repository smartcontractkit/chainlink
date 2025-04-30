package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"

	solToken "github.com/gagliardetto/solana-go/programs/token"

	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
)

// use this changest to deploy a token, create ATAs and mint the token to those ATAs
var _ deployment.ChangeSet[DeploySolanaTokenConfig] = DeploySolanaToken

// use this changeset to mint the token to an address
var _ deployment.ChangeSet[MintSolanaTokenConfig] = MintSolanaToken

// use this changeset to create ATAs for a token
var _ deployment.ChangeSet[CreateSolanaTokenATAConfig] = CreateSolanaTokenATA

// use this changeset to set the authority of a token
var _ deployment.ChangeSet[SetTokenAuthorityConfig] = SetTokenAuthority

func getMintIxs(e deployment.Environment, chain deployment.SolChain, tokenprogramID, mint solana.PublicKey, amountToAddress map[string]uint64) ([]solana.Instruction, error) {
	instructions := []solana.Instruction{}
	for toAddress, amount := range amountToAddress {
		e.Logger.Infof("Minting %d to %s", amount, toAddress)
		toAddressBase58 := solana.MustPublicKeyFromBase58(toAddress)
		// get associated token account for toAddress
		ata, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenprogramID, mint, toAddressBase58)
		mintToI, err := solTokenUtil.MintTo(amount, tokenprogramID, mint, ata, chain.DeployerKey.PublicKey())
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, mintToI)
	}
	return instructions, nil
}

func createATAIx(e deployment.Environment, chain deployment.SolChain, tokenprogramID, mint solana.PublicKey, ataList []string) ([]solana.Instruction, error) {
	instructions := []solana.Instruction{}
	for _, ata := range ataList {
		e.Logger.Infof("Creating ATA for account %s for token %s", ata, mint.String())
		createATAIx, _, err := solTokenUtil.CreateAssociatedTokenAccount(
			tokenprogramID,
			mint,
			solana.MustPublicKeyFromBase58(ata),
			chain.DeployerKey.PublicKey(),
		)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, createATAIx)
	}
	return instructions, nil
}

// TODO: add option to set token mint authority by taking in its public key
// might need to take authority private key if it needs to sign that
type DeploySolanaTokenConfig struct {
	ChainSelector       uint64
	TokenProgramName    deployment.ContractType
	TokenDecimals       uint8
	TokenSymbol         string
	ATAList             []string          // addresses to create ATAs for
	MintAmountToAddress map[string]uint64 // address -> amount
}

func NewTokenInstruction(chain deployment.SolChain, cfg DeploySolanaTokenConfig) ([]solana.Instruction, solana.PrivateKey, error) {
	tokenprogramID, err := GetTokenProgramID(cfg.TokenProgramName)
	if err != nil {
		return nil, nil, err
	}
	// token mint authority
	// can accept a private key in config and pass in pub key here and private key as signer
	tokenAdminPubKey := chain.DeployerKey.PublicKey()
	mintPrivKey, _ := solana.NewRandomPrivateKey()
	mint := mintPrivKey.PublicKey() // this is the token address
	instructions, err := solTokenUtil.CreateToken(
		context.Background(),
		tokenprogramID,
		mint,
		tokenAdminPubKey,
		cfg.TokenDecimals,
		chain.Client,
		cldf.SolDefaultCommitment,
	)
	if err != nil {
		return nil, nil, err
	}
	return instructions, mintPrivKey, nil
}

func DeploySolanaToken(e deployment.Environment, cfg DeploySolanaTokenConfig) (deployment.ChangesetOutput, error) {
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	tokenprogramID, err := GetTokenProgramID(cfg.TokenProgramName)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// create token ix
	instructions, mintPrivKey, err := NewTokenInstruction(chain, cfg)
	mint := mintPrivKey.PublicKey()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// ata ix
	ataIxs, err := createATAIx(e, chain, tokenprogramID, mint, cfg.ATAList)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	instructions = append(instructions, ataIxs...)

	// mint ix
	mintIxs, err := getMintIxs(e, chain, tokenprogramID, mint, cfg.MintAmountToAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	instructions = append(instructions, mintIxs...)

	err = chain.Confirm(instructions, solCommonUtil.AddSigners(mintPrivKey))
	if err != nil {
		e.Logger.Errorw("Failed to confirm instructions for token deployment", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}

	newAddresses := deployment.NewMemoryAddressBook()
	tv := deployment.NewTypeAndVersion(deployment.ContractType(cfg.TokenProgramName), deployment.Version1_0_0)
	tv.AddLabel(cfg.TokenSymbol)
	err = newAddresses.Save(cfg.ChainSelector, mint.String(), tv)
	if err != nil {
		e.Logger.Errorw("Failed to save token", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}

	e.Logger.Infow("Deployed contract", "Contract", tv.String(), "addr", mint.String(), "chain", chain.String())

	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

type MintSolanaTokenConfig struct {
	ChainSelector   uint64
	TokenPubkey     string
	AmountToAddress map[string]uint64 // address -> amount
}

func (cfg MintSolanaTokenConfig) Validate(e deployment.Environment) error {
	chain := e.SolChains[cfg.ChainSelector]
	tokenAddress := solana.MustPublicKeyFromBase58(cfg.TokenPubkey)
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	tokenprogramID, err := chainState.TokenToTokenProgram(tokenAddress)
	if err != nil {
		return err
	}

	accountInfo, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), tokenAddress, &rpc.GetAccountInfoOpts{
		Commitment: cldf.SolDefaultCommitment,
	})
	if err != nil {
		fmt.Println("error getting account info", err)
		return err
	}
	if accountInfo == nil || accountInfo.Value == nil {
		return fmt.Errorf("token address %s not found", tokenAddress.String())
	}
	if accountInfo.Value.Owner != tokenprogramID {
		return fmt.Errorf("token address %s is not owned by the SPL token program", tokenAddress.String())
	}
	return nil
}

func MintSolanaToken(e deployment.Environment, cfg MintSolanaTokenConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// get chain
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	// get addresses
	tokenAddress := solana.MustPublicKeyFromBase58(cfg.TokenPubkey)
	// get token program id
	tokenprogramID, _ := chainState.TokenToTokenProgram(tokenAddress)

	// get mint instructions
	instructions, err := getMintIxs(e, chain, tokenprogramID, tokenAddress, cfg.AmountToAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// confirm instructions
	err = chain.Confirm(instructions)
	if err != nil {
		e.Logger.Errorw("Failed to confirm instructions for token minting", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Minted tokens on", "chain", cfg.ChainSelector, "for token", tokenAddress.String())

	return deployment.ChangesetOutput{}, nil
}

type CreateSolanaTokenATAConfig struct {
	ChainSelector uint64
	TokenPubkey   solana.PublicKey
	TokenProgram  deployment.ContractType
	ATAList       []string // addresses to create ATAs for
}

func CreateSolanaTokenATA(e deployment.Environment, cfg CreateSolanaTokenATAConfig) (deployment.ChangesetOutput, error) {
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	tokenprogramID, err := chainState.TokenToTokenProgram(cfg.TokenPubkey)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// create instructions for each ATA
	instructions, err := createATAIx(e, chain, tokenprogramID, cfg.TokenPubkey, cfg.ATAList)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// confirm instructions
	err = chain.Confirm(instructions)
	if err != nil {
		e.Logger.Errorw("Failed to confirm instructions for ATA creation", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Created ATAs on", "chain", cfg.ChainSelector, "for token", cfg.TokenPubkey.String(), "numATAs", len(cfg.ATAList))

	return deployment.ChangesetOutput{}, nil
}

type SetTokenAuthorityConfig struct {
	ChainSelector uint64
	AuthorityType solToken.AuthorityType
	TokenPubkey   solana.PublicKey
	NewAuthority  solana.PublicKey
}

func SetTokenAuthority(e deployment.Environment, cfg SetTokenAuthorityConfig) (deployment.ChangesetOutput, error) {
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	tokenprogramID, err := chainState.TokenToTokenProgram(cfg.TokenPubkey)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	ix, err := solToken.NewSetAuthorityInstruction(
		cfg.AuthorityType,
		cfg.NewAuthority,
		cfg.TokenPubkey,
		chain.DeployerKey.PublicKey(),
		solana.PublicKeySlice{},
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	tokenIx := &solTokenUtil.TokenInstruction{Instruction: ix, Program: tokenprogramID}

	// confirm instructions
	if err = chain.Confirm([]solana.Instruction{tokenIx}); err != nil {
		e.Logger.Errorw("Failed to confirm instructions for SetTokenAuthority", "chain", chain.String(), "err", err)
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Set token authority on", "chain", cfg.ChainSelector, "for token", cfg.TokenPubkey.String(), "newAuthority", cfg.NewAuthority.String(), "authorityType", cfg.AuthorityType)

	return deployment.ChangesetOutput{}, nil
}
