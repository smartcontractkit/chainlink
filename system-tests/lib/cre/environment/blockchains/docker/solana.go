package docker

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"

	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

type SolanaDeployer struct {
	commonLogger logger.Logger
}

func NewSolanaDeployer(commonLogger logger.Logger) *SolanaDeployer {
	return &SolanaDeployer{
		commonLogger: commonLogger,
	}
}

func (s *SolanaDeployer) Deploy(input *blockchain.Input) (*creblockchains.DeployedBLockchain, error) {
	err := initSolanaInput(input)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to init Solana input")
	}

	bcOut, err := blockchain.NewBlockchainNetwork(input)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
	}

	w, wrapErr := creblockchains.WrapSolana(input, bcOut)
	if wrapErr != nil {
		return nil, pkgerrors.Wrap(wrapErr, "failed to wrap Solana")
	}

	chainsConfigs := make([]devenv.ChainConfig, 0)
	cfg, cfgErr := cre.ChainConfigFromWrapped(w)
	if cfgErr != nil {
		return nil, pkgerrors.Wrap(cfgErr, "failed to wrap blockchain output to chain config")
	}
	chainsConfigs = append(chainsConfigs, cfg)

	cldfBlockchain, err := devenv.NewChains(s.commonLogger, chainsConfigs)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create chains")
	}

	return &creblockchains.DeployedBLockchain{
		CldfBlockchain: maps.Collect(cldfBlockchain.All())[w.ChainSelector],
		Blockchain:     w,
	}, nil
}

var once = &sync.Once{}

func initSolanaInput(bi *blockchain.Input) error {
	err := creblockchains.SetDefaultSolanaPrivateKeyIfEmpty(creblockchains.DefaultSolanaPrivateKey)
	if err != nil {
		return errors.New("failed to set default solana private key")
	}
	bi.PublicKey = creblockchains.DefaultSolanaPrivateKey.PublicKey().String()
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
