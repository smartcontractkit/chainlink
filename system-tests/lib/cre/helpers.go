package cre

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
)

// ChainConfigFromWrapped converts a single wrapped chain into a devenv.ChainConfig.
func ChainConfigFromWrapped(w *WrappedBlockchainOutput) (devenv.ChainConfig, error) {
	if w == nil || w.BlockchainOutput == nil || len(w.BlockchainOutput.Nodes) == 0 {
		return devenv.ChainConfig{}, errors.New("invalid wrapped blockchain output")
	}
	n := w.BlockchainOutput.Nodes[0]

	cfg := devenv.ChainConfig{
		WSRPCs: []devenv.CribRPCs{{
			External: n.ExternalWSUrl, Internal: n.InternalWSUrl,
		}},
		HTTPRPCs: []devenv.CribRPCs{{
			External: n.ExternalHTTPUrl, Internal: n.InternalHTTPUrl,
		}},
	}

	cfg.ChainType = strings.ToUpper(w.BlockchainOutput.Family)

	// Solana
	if w.SolChain != nil {
		cfg.ChainID = w.SolChain.ChainID
		cfg.SolDeployerKey = w.SolChain.PrivateKey
		cfg.SolArtifactDir = w.SolChain.ArtifactsDir
		return cfg, nil
	}

	if strings.EqualFold(cfg.ChainType, blockchain.FamilyTron) {
		cfg.ChainID = strconv.FormatUint(w.ChainID, 10)
		privateKey, err := crypto.HexToECDSA(w.DeployerPrivateKey)
		if err != nil {
			return devenv.ChainConfig{}, errors.Wrap(err, "failed to parse private key for Tron")
		}

		deployerKey, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(conversions.MustSafeInt64(w.ChainID)))
		if err != nil {
			return devenv.ChainConfig{}, errors.Wrap(err, "failed to create transactor for Tron")
		}
		cfg.DeployerKey = deployerKey
		return cfg, nil
	}

	// EVM
	if w.SethClient == nil {
		return devenv.ChainConfig{}, fmt.Errorf("blockchain output evm family without SethClient for chainID %d", w.ChainID)
	}

	cfg.ChainID = strconv.FormatUint(w.ChainID, 10)
	cfg.ChainName = w.SethClient.Cfg.Network.Name
	// ensure nonce fetched from chain at use time
	cfg.DeployerKey = w.SethClient.NewTXOpts(seth.WithNonce(nil))

	return cfg, nil
}
