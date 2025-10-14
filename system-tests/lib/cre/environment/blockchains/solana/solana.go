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
	"github.com/gagliardetto/solana-go/rpc"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_solana_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana/provider"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
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

func (s *Deployer) Deploy(input *blockchain.Input) (*cre.Blockchain, error) {
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

	return &cre.Blockchain{
		CtfOutput: bcOut,
		SolClient: solClient,
		SolChain: &cre.SolChain{
			ChainSelector: sel,
			ChainID:       input.ChainID,
			PrivateKey:    pk,
			ArtifactsDir:  input.ContractsDir,
		},
		Funder: &Funder{
			solClient:      solClient,
			mainPrivateKey: pk,
			testLogger:     s.testLogger,
		},
	}, nil
}

type Funder struct {
	solClient      *rpc.Client
	mainPrivateKey solana.PrivateKey
	testLogger     zerolog.Logger
}

func (f *Funder) Fund(ctx context.Context, address string, amount uint64, fundingPrivateKey []byte) error {
	recipient := solana.MustPublicKeyFromBase58(address)
	f.testLogger.Info().Msgf("Attempting to fund Solana account %s", recipient.String())

	err := libfunding.SendFundsSol(ctx, f.testLogger, f.solClient, libfunding.FundsToSendSol{
		Recipent:   recipient,
		PrivateKey: solana.PrivateKey(fundingPrivateKey),
		Amount:     amount,
	})
	if err != nil {
		return fmt.Errorf("failed to fund Solana account for a node: %w", err)
	}
	f.testLogger.Info().Msgf("Successfully funded Solana account %s", recipient.String())

	return nil
}

func (f *Funder) Prepare(ctx context.Context, requiredTotal uint64) ([]byte, error) {
	private, pkErr := solana.NewRandomPrivateKey()
	if pkErr != nil {
		return nil, pkgerrors.Wrap(pkErr, "failed to generate private key for solana")
	}
	public := private.PublicKey()
	fundingErr := libfunding.SendFundsSol(ctx, zerolog.Logger{}, f.solClient, libfunding.FundsToSendSol{
		Recipent:   public,
		PrivateKey: f.mainPrivateKey,
		Amount:     requiredTotal,
	})
	if fundingErr != nil {
		return nil, pkgerrors.Wrapf(fundingErr, " failed to fund funding accounts on chain on Solana")
	}

	return private, nil
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
