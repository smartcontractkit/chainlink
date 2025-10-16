package solana

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gagliardetto/solana-go"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf_solana_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana/provider"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

var DefaultSolanaPrivateKey = solana.MustPrivateKeyFromBase58("4u2itaM9r5kxsmoti3GMSDZrQEFpX14o6qPWY9ZrrYTR6kduDBr4YAZJsjawKzGP3wDzyXqterFmfcLUmSBro5AT")

type Deployer struct {
	provider   infra.Provider
	testLogger zerolog.Logger
}

func NewDeployer(testLogger zerolog.Logger, provider *infra.Provider) *Deployer {
	return &Deployer{
		provider:   *provider,
		testLogger: testLogger,
	}
}

type Blockchain struct {
	testLogger    zerolog.Logger
	chainSelector uint64
	ctfOutput     *blockchain.Output
	SolClient     *solrpc.Client
	SolanaChainID string
	PrivateKey    solana.PrivateKey
	ArtifactsDir  string
}

func (s *Blockchain) ChainSelector() uint64 {
	return s.chainSelector
}
func (s *Blockchain) ChainID() uint64 {
	return 0 // Solana doesn't use numeric chain IDs
}

func (s *Blockchain) CtfOutput() *blockchain.Output {
	return s.ctfOutput
}

func (e *Blockchain) Is(chainFamily string) bool {
	return strings.EqualFold(e.ctfOutput.Family, chainFamily)
}

func (e *Blockchain) ChainFamily() string {
	return e.ctfOutput.Family
}

func (s *Blockchain) Fund(ctx context.Context, address string, amount uint64) error {
	recipient := solana.MustPublicKeyFromBase58(address)
	s.testLogger.Info().Msgf("Attempting to fund Solana account %s", recipient.String())

	err := libfunding.SendFundsSol(ctx, s.testLogger, s.SolClient, libfunding.FundsToSendSol{
		Recipent:   recipient,
		PrivateKey: s.PrivateKey,
		Amount:     amount,
	})
	if err != nil {
		return fmt.Errorf("failed to fund Solana account for a node: %w", err)
	}
	s.testLogger.Info().Msgf("Successfully funded Solana account %s", recipient.String())

	return nil
}

func (s *Blockchain) ToCldfChain() (cldf_chain.BlockChain, error) {
	// bcNode := s.CtfOutput().Nodes[0]

	// return &blockchains.CldfChainConfig{
	// 	WSRPCs: []blockchains.RPCs{{
	// 		External: bcNode.ExternalWSUrl, Internal: bcNode.InternalWSUrl,
	// 	}},
	// 	HTTPRPCs: []blockchains.RPCs{{
	// 		External: bcNode.ExternalHTTPUrl, Internal: bcNode.InternalHTTPUrl,
	// 	}},
	// 	ChainType:      strings.ToUpper(s.CtfOutput().Family),
	// 	ChainID:        s.SolanaChainID,
	// 	SolDeployerKey: s.PrivateKey,
	// 	SolArtifactDir: s.ArtifactsDir,
	// }, nil

	solArtifactPath := s.ArtifactsDir
	if solArtifactPath == "" {
		s.testLogger.Info().Msg("Creating tmp directory for generated solana programs and keypairs")
		solArtifactPath, err := os.MkdirTemp("", "solana-artifacts")
		s.testLogger.Info().Msgf("Solana programs tmp dir at %s", solArtifactPath)
		if err != nil {
			return nil, err
		}
	}

	sc := solrpc.New(s.CtfOutput().Nodes[0].ExternalHTTPUrl)
	return cldf_solana.Chain{
		Selector:    s.ChainSelector(),
		Client:      sc,
		DeployerKey: &s.PrivateKey,
		KeypairPath: solArtifactPath + "/deploy-keypair.json",
		URL:         s.CtfOutput().Nodes[0].ExternalHTTPUrl,
		WSURL:       s.CtfOutput().Nodes[0].ExternalWSUrl,
		Confirm: func(instructions []solanago.Instruction, opts ...solCommonUtil.TxModifier) error {
			_, err := solCommonUtil.SendAndConfirm(
				context.Background(), sc, instructions, s.PrivateKey, solrpc.CommitmentConfirmed, opts...,
			)
			return err
		},
		ProgramsPath: solArtifactPath,
	}, nil
}

func (s *Deployer) Deploy(input *blockchain.Input) (blockchains.Blockchain, error) {
	if s.provider.IsCRIB() {
		return nil, errors.New("CRIB deployment for Solana is not supported yet")
	}

	err := initSolanaInput(input)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to init Solana input")
	}

	bcOut, err := blockchain.NewBlockchainNetwork(input)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
	}

	sel, ok := chainselectors.SolanaChainIdToChainSelector()[input.ChainID]
	if !ok {
		return nil, fmt.Errorf("selector not found for solana chainID '%s'", input.ChainID)
	}
	// shouldn't be empty, since we call initSolana before wrap, but just in case
	setErr := setDefaultPrivateKeyIfEmpty()
	if setErr != nil {
		return nil, fmt.Errorf("set default private key solana failed: %w", setErr)
	}

	envp := os.Getenv("SOLANA_PRIVATE_KEY")
	pk, err := solana.PrivateKeyFromBase58(envp)
	if err != nil {
		return nil, errors.New("failed to decode private key for solana")
	}

	if err := cldf_solana_provider.WritePrivateKeyToPath(filepath.Join(input.ContractsDir, "deploy-keypair.json"), pk); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to save private key for solana")
	}

	solClient := rpc.New(bcOut.Nodes[0].ExternalHTTPUrl)

	return &Blockchain{
		SolClient:     solClient,
		SolanaChainID: input.ChainID,
		chainSelector: sel,
		PrivateKey:    pk,
		ArtifactsDir:  input.ContractsDir,
	}, nil
}

// }

// return &cre.Blockchain{
// 	CtfOutput: bcOut,
// 	SolClient: solClient,
// 	SolChain: &cre.SolChain{
// 		ChainSelector: sel,
// 		ChainID:       input.ChainID,
// 		PrivateKey:    pk,
// 		ArtifactsDir:  input.ContractsDir,
// 	},
// 	Funder: &Funder{
// 		solClient:      solClient,
// 		mainPrivateKey: pk,
// 		testLogger:     s.testLogger,
// 	},
// }, nil
// }

// type Funder struct {
// 	solClient      *rpc.Client
// 	mainPrivateKey solana.PrivateKey
// 	testLogger     zerolog.Logger
// }

// func (f *Funder) Fund(ctx context.Context, address string, amount uint64, fundingPrivateKey []byte) error {
// 	recipient := solana.MustPublicKeyFromBase58(address)
// 	f.testLogger.Info().Msgf("Attempting to fund Solana account %s", recipient.String())

// 	err := libfunding.SendFundsSol(ctx, f.testLogger, f.solClient, libfunding.FundsToSendSol{
// 		Recipent:   recipient,
// 		PrivateKey: solana.PrivateKey(fundingPrivateKey),
// 		Amount:     amount,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to fund Solana account for a node: %w", err)
// 	}
// 	f.testLogger.Info().Msgf("Successfully funded Solana account %s", recipient.String())

// 	return nil
// }

// func (f *Funder) Prepare(ctx context.Context, requiredTotal uint64) ([]byte, error) {
// 	private, pkErr := solana.NewRandomPrivateKey()
// 	if pkErr != nil {
// 		return nil, pkgerrors.Wrap(pkErr, "failed to generate private key for solana")
// 	}
// 	public := private.PublicKey()
// 	fundingErr := libfunding.SendFundsSol(ctx, zerolog.Logger{}, f.solClient, libfunding.FundsToSendSol{
// 		Recipent:   public,
// 		PrivateKey: f.mainPrivateKey,
// 		Amount:     requiredTotal,
// 	})
// 	if fundingErr != nil {
// 		return nil, pkgerrors.Wrapf(fundingErr, " failed to fund funding accounts on chain on Solana")
// 	}

// 	return private, nil
// }

var once = &sync.Once{}

func initSolanaInput(bi *blockchain.Input) error {
	err := setDefaultPrivateKeyIfEmpty()
	if err != nil {
		return errors.New("failed to set default solana private key")
	}
	bi.PublicKey = DefaultSolanaPrivateKey.PublicKey().String()
	bi.ContractsDir = getSolProgramsPath(bi.ContractsDir)

	if bi.SolanaPrograms != nil {
		var err2 error
		once.Do(func() {
			if hasSolanaArtifacts(bi.ContractsDir) {
				return
			}
			// TODO PLEX-1718 use latest contracts sha for now. Derive commit sha from go.mod once contracts are in a separate go module
			err2 = memory.DownloadSolanaProgramArtifacts(context.Background(), bi.ContractsDir, logger.Nop(), "b0f7cd3fbdbb")
		})
		if err2 != nil {
			return fmt.Errorf("failed to download solana artifacts: %w", err2)
		}
	}

	return nil
}

func hasSolanaArtifacts(dir string) bool {
	ents, err := os.ReadDir(dir)
	if err != nil { // dir missing or unreadable -> treat as not present
		return false
	}
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if strings.HasSuffix(n, ".so") || strings.HasSuffix(n, ".json") {
			return true
		}
	}
	return false
}

func getSolProgramsPath(path string) string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, path)
}

func setDefaultPrivateKeyIfEmpty() error {
	if os.Getenv("SOLANA_PRIVATE_KEY") == "" {
		setErr := os.Setenv("SOLANA_PRIVATE_KEY", DefaultSolanaPrivateKey.String())
		if setErr != nil {
			return fmt.Errorf("failed to set SOLANA_PRIVATE_KEY environment variable: %w", setErr)
		}
		framework.L.Info().Msgf("Set SOLANA_PRIVATE_KEY environment variable to default value: %s", os.Getenv("PRIVATE_KEY"))
	}

	return nil
}
