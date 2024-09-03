package solclient

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"golang.org/x/sync/errgroup"

	access_controller2 "github.com/smartcontractkit/chainlink-solana/contracts/generated/access_controller"
	ocr_2 "github.com/smartcontractkit/chainlink-solana/contracts/generated/ocr_2"
	store2 "github.com/smartcontractkit/chainlink-solana/contracts/generated/store"
	test_env_sol "github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

// All account sizes are calculated from Rust structures, ex. programs/access-controller/src/lib.rs:L80
// there is some wrapper in "anchor" that creates accounts for programs automatically, but we are doing that explicitly
const (
	Discriminator = 8
	// TokenMintAccountSize default size of data required for a new mint account
	TokenMintAccountSize             = uint64(82)
	TokenAccountSize                 = uint64(165)
	AccessControllerStateAccountSize = uint64(Discriminator + solana.PublicKeyLength + solana.PublicKeyLength + 8 + 32*64)
	StoreAccountSize                 = uint64(Discriminator + solana.PublicKeyLength*3)
	OCRTransmissionsAccountSize      = uint64(Discriminator + 192 + 8192*48)
	OCRProposalAccountSize           = Discriminator + 1 + 32 + 1 + 1 + (1 + 4) + 32 + ProposedOraclesSize + OCROffChainConfigSize
	ProposedOracleSize               = uint64(solana.PublicKeyLength + 20 + 4 + solana.PublicKeyLength)
	ProposedOraclesSize              = ProposedOracleSize*19 + 8
	OCROracle                        = uint64(solana.PublicKeyLength + 20 + solana.PublicKeyLength + solana.PublicKeyLength + 4 + 8)
	OCROraclesSize                   = OCROracle*19 + 8
	OCROffChainConfigSize            = uint64(8 + 4096 + 8)
	OCRConfigSize                    = 32 + 32 + 32 + 32 + 32 + 32 + 16 + 16 + (1 + 1 + 2 + 4 + 4 + 32) + (4 + 32 + 8) + (4 + 4)
	OCRAccountSize                   = Discriminator + 1 + 1 + 2 + 4 + solana.PublicKeyLength + OCRConfigSize + OCROffChainConfigSize + OCROraclesSize
	keypairSuffix                    = "-keypair.json"
)

type Authority struct {
	PublicKey solana.PublicKey
	Nonce     uint8
}

type ContractDeployer struct {
	Client   *Client
	Accounts *Accounts
}

// GenerateAuthorities generates authorities so other contracts can access OCR with on-chain calls when signer needed
func (c *ContractDeployer) GenerateAuthorities(seeds []string) error {
	authorities := make(map[string]*Authority)
	for _, seed := range seeds {
		auth, nonce, err := c.Client.FindAuthorityAddress(seed, c.Accounts.OCR.PublicKey(), c.Client.ProgramWallets["ocr2-keypair.json"].PublicKey())
		if err != nil {
			return err
		}
		authorities[seed] = &Authority{
			PublicKey: auth,
			Nonce:     nonce,
		}
	}
	c.Accounts.Authorities = authorities
	c.Accounts.Owner = c.Client.DefaultWallet
	return nil
}

// addMintInstr adds instruction for creating new mint (token)
func (c *ContractDeployer) addMintInstr(instr *[]solana.Instruction) error {
	accInstr, err := c.Client.CreateAccInstr(c.Accounts.Mint.PublicKey(), TokenMintAccountSize, token.ProgramID)
	if err != nil {
		return err
	}
	*instr = append(
		*instr,
		accInstr,
		token.NewInitializeMintInstruction(
			18,
			c.Accounts.MintAuthority.PublicKey(),
			c.Accounts.MintAuthority.PublicKey(),
			c.Accounts.Mint.PublicKey(),
			solana.SysVarRentPubkey,
		).Build())
	return nil
}

func (c *ContractDeployer) SetupAssociatedAccount() (*solana.PublicKey, *solana.PublicKey, error) {
	vault := c.Accounts.Authorities["vault"]
	payer := c.Client.DefaultWallet
	instr := make([]solana.Instruction, 0)
	ainstr := associatedtokenaccount.NewCreateInstruction(
		c.Client.DefaultWallet.PublicKey(),
		vault.PublicKey,
		c.Accounts.Mint.PublicKey(),
	).Build()
	aaccount := ainstr.Impl.(associatedtokenaccount.Create).AccountMetaSlice[1].PublicKey
	instr = append(instr,
		ainstr,
	)
	c.Accounts.OCRVaultAssociatedPubKey = aaccount
	err := c.Client.TXSync(
		"Setup associated account",
		rpc.CommitmentConfirmed,
		instr,
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
	if err != nil {
		return nil, nil, err
	}
	return &aaccount, &vault.PublicKey, err
}

// AddNewAssociatedAccInstr adds instruction to create new account associated with some mint (token)
func (c *ContractDeployer) AddNewAssociatedAccInstr(acc solana.PublicKey, ownerPubKey solana.PublicKey, assocAccount solana.PublicKey, instr *[]solana.Instruction) error {
	accInstr, err := c.Client.CreateAccInstr(acc, TokenAccountSize, token.ProgramID)
	if err != nil {
		return err
	}
	*instr = append(*instr,
		accInstr,
		token.NewInitializeAccountInstruction(
			acc,
			c.Accounts.Mint.PublicKey(),
			ownerPubKey,
			solana.SysVarRentPubkey,
		).Build(),
		associatedtokenaccount.NewCreateInstruction(
			c.Client.DefaultWallet.PublicKey(),
			assocAccount,
			c.Accounts.Mint.PublicKey(),
		).Build(),
	)
	return nil
}

func (c *ContractDeployer) DeployOCRv2Store(billingAC string) (*Store, error) {
	programWallet := c.Client.ProgramWallets["store-keypair.json"]
	payer := c.Client.DefaultWallet
	accInstruction, err := c.Client.CreateAccInstr(c.Accounts.Store.PublicKey(), StoreAccountSize, programWallet.PublicKey())
	if err != nil {
		return nil, err
	}
	bacPublicKey, err := solana.PublicKeyFromBase58(billingAC)
	if err != nil {
		return nil, err
	}
	err = c.Client.TXSync(
		"Deploy store",
		rpc.CommitmentConfirmed,
		[]solana.Instruction{
			accInstruction,
			store2.NewInitializeInstruction(
				c.Accounts.Store.PublicKey(),
				c.Accounts.Owner.PublicKey(),
				bacPublicKey,
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(c.Accounts.Owner.PublicKey()) {
				return &c.Accounts.Owner.PrivateKey
			}
			if key.Equals(c.Accounts.Store.PublicKey()) {
				return &c.Accounts.Store.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
	if err != nil {
		return nil, err
	}
	return &Store{
		Client:        c.Client,
		Store:         c.Accounts.Store,
		Feed:          c.Accounts.Feed,
		Owner:         c.Accounts.Owner,
		ProgramWallet: programWallet,
	}, nil
}

func (c *ContractDeployer) CreateFeed(desc string, decimals uint8, granularity int, liveLength int) error {
	payer := c.Client.DefaultWallet
	programWallet := c.Client.ProgramWallets["store-keypair.json"]
	feedAccInstruction, err := c.Client.CreateAccInstr(c.Accounts.Feed.PublicKey(), OCRTransmissionsAccountSize, programWallet.PublicKey())
	if err != nil {
		return err
	}
	err = c.Client.TXSync(
		"Create feed",
		rpc.CommitmentConfirmed,
		[]solana.Instruction{
			feedAccInstruction,
			store2.NewCreateFeedInstruction(
				desc,
				decimals,
				uint8(granularity),
				uint32(liveLength),
				c.Accounts.Feed.PublicKey(),
				c.Accounts.Owner.PublicKey(),
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(c.Accounts.Owner.PublicKey()) {
				return &c.Accounts.Owner.PrivateKey
			}
			if key.Equals(c.Accounts.Feed.PublicKey()) {
				return &c.Accounts.Feed.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContractDeployer) DeployLinkTokenContract() (*LinkToken, error) {
	var err error
	payer := c.Client.DefaultWallet

	instr := make([]solana.Instruction, 0)
	if err = c.addMintInstr(&instr); err != nil {
		return nil, err
	}
	err = c.Client.TXAsync(
		"Creating LINK Token and associated accounts",
		instr,
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(c.Accounts.Mint.PublicKey()) {
				return &c.Accounts.Mint.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			if key.Equals(c.Accounts.MintAuthority.PublicKey()) {
				return &c.Accounts.MintAuthority.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
	if err != nil {
		return nil, err
	}
	return &LinkToken{
		Client:        c.Client,
		Mint:          c.Accounts.Mint,
		MintAuthority: c.Accounts.MintAuthority,
	}, nil
}

func (c *ContractDeployer) InitOCR2(billingControllerAddr string, requesterControllerAddr string) (*OCRv2, error) {
	programWallet := c.Client.ProgramWallets["ocr2-keypair.json"]
	payer := c.Client.DefaultWallet
	ocrAccInstruction, err := c.Client.CreateAccInstr(c.Accounts.OCR.PublicKey(), OCRAccountSize, programWallet.PublicKey())
	if err != nil {
		return nil, err
	}
	bacPubKey, err := solana.PublicKeyFromBase58(billingControllerAddr)
	if err != nil {
		return nil, err
	}
	racPubKey, err := solana.PublicKeyFromBase58(requesterControllerAddr)
	if err != nil {
		return nil, err
	}
	assocVault, vault, err := c.SetupAssociatedAccount()
	if err != nil {
		return nil, err
	}
	instr := make([]solana.Instruction, 0)
	instr = append(instr,
		ocrAccInstruction,
		ocr_2.NewInitializeInstructionBuilder().
			SetMinAnswer(ag_binary.Int128{
				Lo: 1,
				Hi: 0,
			}).
			SetMaxAnswer(ag_binary.Int128{
				Lo: 1000000,
				Hi: 0,
			}).
			SetStateAccount(c.Accounts.OCR.PublicKey()).
			SetFeedAccount(c.Accounts.Feed.PublicKey()).
			SetOwnerAccount(c.Accounts.Owner.PublicKey()).
			SetTokenMintAccount(c.Accounts.Mint.PublicKey()).
			SetTokenVaultAccount(*assocVault).
			SetVaultAuthorityAccount(*vault).
			SetRequesterAccessControllerAccount(racPubKey).
			SetBillingAccessControllerAccount(bacPubKey).
			Build())
	err = c.Client.TXSync(
		"Initializing OCRv2",
		rpc.CommitmentConfirmed,
		instr,
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			if key.Equals(c.Accounts.OCR.PublicKey()) {
				return &c.Accounts.OCR.PrivateKey
			}
			if key.Equals(c.Accounts.Owner.PublicKey()) {
				return &c.Accounts.Owner.PrivateKey
			}
			if key.Equals(c.Accounts.OCRVault.PublicKey()) {
				return &c.Accounts.OCRVault.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
	if err != nil {
		return nil, err
	}
	return &OCRv2{
		ContractDeployer:         c,
		Client:                   c.Client,
		State:                    c.Accounts.OCR,
		Authorities:              c.Accounts.Authorities,
		Owner:                    c.Accounts.Owner,
		Proposal:                 c.Accounts.Proposal,
		OCRVaultAssociatedPubKey: *assocVault,
		Mint:                     c.Accounts.Mint,
		ProgramWallet:            programWallet,
	}, nil
}

func (c *ContractDeployer) DeployProgramRemote(programName string, env *environment.Environment) error {
	log.Debug().Str("Program", programName).Msg("Deploying program")
	programPath := filepath.Join("programs", programName)
	programKeyFileName := strings.Replace(programName, ".so", keypairSuffix, -1)
	programKeyFilePath := filepath.Join("programs", programKeyFileName)
	cmd := fmt.Sprintf("solana program deploy --program-id %s %s", programKeyFilePath, programPath)
	pl, err := env.Client.ListPods(env.Cfg.Namespace, "app=sol")
	if err != nil {
		return err
	}
	stdOutBytes, stdErrBytes, _ := env.Client.ExecuteInPod(env.Cfg.Namespace, pl.Items[0].Name, "sol-val", strings.Split(cmd, " "))
	log.Debug().Str("STDOUT", string(stdOutBytes)).Str("STDERR", string(stdErrBytes)).Str("CMD", cmd).Send()
	return nil
}

func (c *ContractDeployer) DeployProgramRemoteLocal(programName string, sol *test_env_sol.Solana) error {
	log.Info().Str("Program", programName).Msg("Deploying program")
	programPath := filepath.Join("programs", programName)
	programKeyFileName := strings.Replace(programName, ".so", keypairSuffix, -1)
	programKeyFilePath := filepath.Join("programs", programKeyFileName)
	cmd := fmt.Sprintf("solana program deploy --program-id %s %s", programKeyFilePath, programPath)
	_, res, err := sol.Container.Exec(context.Background(), strings.Split(cmd, " "))
	if err != nil {
		return err
	}
	out, err := io.ReadAll(res)
	if err != nil {
		return err
	}
	log.Info().Str("Output", string(out)).Msg("Deploying " + programName)
	return nil
}

func (c *ContractDeployer) DeployOCRv2AccessController() (*AccessController, error) {
	programWallet := c.Client.ProgramWallets["access_controller-keypair.json"]
	payer := c.Client.DefaultWallet
	stateAcc := solana.NewWallet()
	accInstruction, err := c.Client.CreateAccInstr(stateAcc.PublicKey(), AccessControllerStateAccountSize, programWallet.PublicKey())
	if err != nil {
		return nil, err
	}
	err = c.Client.TXAsync(
		"Initializing access controller",
		[]solana.Instruction{
			accInstruction,
			access_controller2.NewInitializeInstruction(
				stateAcc.PublicKey(),
				c.Accounts.Owner.PublicKey(),
			).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(c.Accounts.Owner.PublicKey()) {
				return &c.Accounts.Owner.PrivateKey
			}
			if key.Equals(stateAcc.PublicKey()) {
				return &stateAcc.PrivateKey
			}
			if key.Equals(payer.PublicKey()) {
				return &payer.PrivateKey
			}
			return nil
		},
		payer.PublicKey(),
	)
	if err != nil {
		return nil, err
	}
	return &AccessController{
		State:         stateAcc,
		Client:        c.Client,
		Owner:         c.Accounts.Owner,
		ProgramWallet: programWallet,
	}, nil
}

func (c *ContractDeployer) RegisterAnchorPrograms() {
	access_controller2.SetProgramID(c.Client.ProgramWallets["access_controller-keypair.json"].PublicKey())
	store2.SetProgramID(c.Client.ProgramWallets["store-keypair.json"].PublicKey())
	ocr_2.SetProgramID(c.Client.ProgramWallets["ocr2-keypair.json"].PublicKey())
}

func (c *ContractDeployer) ValidateProgramsDeployed() error {
	keys := []solana.PublicKey{}
	names := []string{}
	for i := range c.Client.ProgramWallets {
		keys = append(keys, c.Client.ProgramWallets[i].PublicKey())
		names = append(names, strings.TrimSuffix(i, keypairSuffix))
	}

	res, err := c.Client.RPC.GetMultipleAccountsWithOpts(
		context.Background(),
		keys,
		&rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentConfirmed,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	output := []string{}
	invalid := false
	for i := range res.Value {
		if res.Value[i] == nil {
			invalid = true
			output = append(output, fmt.Sprintf("%s=nil", names[i]))
			continue
		}
		if !res.Value[i].Executable {
			invalid = true
			output = append(output, fmt.Sprintf("%s=notProgram(%s)", names[i], keys[i].String()))
			continue
		}
		output = append(output, fmt.Sprintf("%s=valid(%s)", names[i], keys[i].String()))
	}

	if invalid {
		return fmt.Errorf("Programs not deployed: %s", strings.Join(output, " "))
	}
	return nil
}

func (c *ContractDeployer) LoadPrograms(contractsDir string) error {
	keyFiles, err := c.Client.ListDirFilenamesByExt(contractsDir, ".json")
	if err != nil {
		return err
	}
	log.Debug().Interface("Files", keyFiles).Msg("Program key files")
	for _, kfn := range keyFiles {
		pk, err := solana.PrivateKeyFromSolanaKeygenFile(filepath.Join(contractsDir, kfn))
		if err != nil {
			return err
		}
		w, err := c.Client.LoadWallet(pk.String())
		if err != nil {
			return err
		}
		c.Client.ProgramWallets[kfn] = w
	}
	log.Debug().Interface("Keys", c.Client.ProgramWallets).Msg("Program wallets")
	return nil
}

func (c *ContractDeployer) DeployAnchorProgramsRemote(contractsDir string, env *environment.Environment) error {
	contractBinaries, err := c.Client.ListDirFilenamesByExt(contractsDir, ".so")
	if err != nil {
		return err
	}
	log.Debug().Interface("Binaries", contractBinaries).Msg("Program binaries")
	g := errgroup.Group{}
	for _, bin := range contractBinaries {
		bin := bin
		g.Go(func() error {
			return c.DeployProgramRemote(bin, env)
		})
	}
	return g.Wait()
}

func (c *ContractDeployer) DeployAnchorProgramsRemoteDocker(contractsDir string, sol *test_env_sol.Solana) error {
	contractBinaries, err := c.Client.ListDirFilenamesByExt(contractsDir, ".so")
	if err != nil {
		return err
	}
	log.Debug().Interface("Binaries", contractBinaries).Msg("Program binaries")
	g := errgroup.Group{}
	for _, bin := range contractBinaries {
		bin := bin
		g.Go(func() error {
			return c.DeployProgramRemoteLocal(bin, sol)
		})
	}
	return g.Wait()
}

func (c *Client) FindAuthorityAddress(seed string, statePubKey solana.PublicKey, progPubKey solana.PublicKey) (solana.PublicKey, uint8, error) {
	log.Debug().
		Str("Seed", seed).
		Str("StatePubKey", statePubKey.String()).
		Str("ProgramPubKey", progPubKey.String()).
		Msg("Trying to find program authority")
	auth, nonce, err := solana.FindProgramAddress([][]byte{[]byte(seed), statePubKey.Bytes()}, progPubKey)
	if err != nil {
		return solana.PublicKey{}, 0, err
	}
	log.Debug().Str("Authority", auth.String()).Uint8("Nonce", nonce).Msg("Found authority addr")
	return auth, nonce, err
}

func NewContractDeployer(client *Client, lt *LinkToken) (*ContractDeployer, error) {
	cd := &ContractDeployer{
		Accounts: &Accounts{
			OCR:           solana.NewWallet(),
			Store:         solana.NewWallet(),
			Feed:          solana.NewWallet(),
			Proposal:      solana.NewWallet(),
			Owner:         solana.NewWallet(),
			Mint:          solana.NewWallet(),
			MintAuthority: solana.NewWallet(),
			OCRVault:      solana.NewWallet(),
		},
		Client: client,
	}
	if lt != nil {
		cd.Accounts.Mint = lt.Mint
		cd.Accounts.MintAuthority = lt.MintAuthority
	}
	return cd, nil
}
