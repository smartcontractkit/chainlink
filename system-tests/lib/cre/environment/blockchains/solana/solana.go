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
	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
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

func (s *Blockchain) IsFamily(chainFamily string) bool {
	return strings.EqualFold(s.ctfOutput.Family, chainFamily)
}

func (s *Blockchain) ChainFamily() string {
	return s.ctfOutput.Family
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
	if s.ArtifactsDir == "" {
		s.testLogger.Info().Msg("Creating tmp directory for generated solana programs and keypairs")
		var err error
		s.ArtifactsDir, err = os.MkdirTemp("", "solana-artifacts")
		s.testLogger.Info().Msgf("Solana programs tmp dir at %s", s.ArtifactsDir)
		if err != nil {
			return nil, err
		}
	}

	if len(s.CtfOutput().Nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for chain %s-%s", s.ChainFamily(), s.SolanaChainID)
	}

	sc := solrpc.New(s.CtfOutput().Nodes[0].ExternalHTTPUrl)
	return cldf_solana.Chain{
		Selector:    s.ChainSelector(),
		Client:      sc,
		DeployerKey: &s.PrivateKey,
		KeypairPath: s.ArtifactsDir + "/deploy-keypair.json",
		URL:         s.CtfOutput().Nodes[0].ExternalHTTPUrl,
		WSURL:       s.CtfOutput().Nodes[0].ExternalWSUrl,
		Confirm: func(instructions []solana.Instruction, opts ...solCommonUtil.TxModifier) error {
			_, err := solCommonUtil.SendAndConfirm(
				context.Background(), sc, instructions, s.PrivateKey, solrpc.CommitmentConfirmed, opts...,
			)
			return err
		},
		ProgramsPath: s.ArtifactsDir,
	}, nil
}

func (s *Deployer) Deploy(ctx context.Context, input *blockchain.Input) (blockchains.Blockchain, error) {
	if s.provider.IsCRIB() {
		return nil, errors.New("CRIB deployment for Solana is not supported yet")
	}

	err := initSolanaInput(input)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to init Solana input")
	}

	bcOut, err := blockchain.NewWithContext(ctx, input)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
	}

	sel, ok := chainselectors.SolanaChainIdToChainSelector()[input.ChainID]
	if !ok {
		return nil, fmt.Errorf("selector not found for solana chainID '%s'", input.ChainID)
	}

	envp := os.Getenv("SOLANA_PRIVATE_KEY")
	pk, err := solana.PrivateKeyFromBase58(envp)
	if err != nil {
		return nil, errors.New("failed to decode private key for solana")
	}

	if err := cldf_solana_provider.WritePrivateKeyToPath(filepath.Join(input.ContractsDir, "deploy-keypair.json"), pk); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to save private key for solana")
	}

	solClient := solrpc.New(bcOut.Nodes[0].ExternalHTTPUrl)

	return &Blockchain{
		SolClient:     solClient,
		SolanaChainID: input.ChainID,
		chainSelector: sel,
		PrivateKey:    pk,
		ArtifactsDir:  input.ContractsDir,
		ctfOutput:     bcOut,
	}, nil
}

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
			err2 = solutils.DownloadChainlinkSolanaProgramArtifacts(context.Background(), bi.ContractsDir, "b0f7cd3fbdbb", logger.Nop())
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
