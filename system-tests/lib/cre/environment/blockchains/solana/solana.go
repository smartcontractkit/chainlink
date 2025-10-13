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
	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_solana_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana/provider"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

var DefaultSolanaPrivateKey = solana.MustPrivateKeyFromBase58("4u2itaM9r5kxsmoti3GMSDZrQEFpX14o6qPWY9ZrrYTR6kduDBr4YAZJsjawKzGP3wDzyXqterFmfcLUmSBro5AT")

type Deployer struct {
	provider infra.Provider
}

func NewDeployer(provider *infra.Provider) *Deployer {
	return &Deployer{
		provider: *provider,
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

	sel, ok := chainselectors.SolanaChainIdToChainSelector()[bcOut.ChainID]
	if !ok {
		return nil, fmt.Errorf("selector not found for solana chainID '%s'", bcOut.ChainID)
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

	return &cre.Blockchain{
		CtfOutput: bcOut,
		SolClient: rpc.New(bcOut.Nodes[0].ExternalHTTPUrl),
		SolChain: &cre.SolChain{
			ChainSelector: sel,
			ChainID:       bcOut.ChainID,
			PrivateKey:    pk,
			ArtifactsDir:  input.ContractsDir,
		},
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
