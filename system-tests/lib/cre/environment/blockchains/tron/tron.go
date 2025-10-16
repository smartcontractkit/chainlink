package tron

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	tron_addr "github.com/fbsobreira/gotron-sdk/pkg/address"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	tronprovider "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron/provider"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

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
	testLogger         zerolog.Logger
	chainSelector      uint64
	chainID            uint64
	ctfOutput          *blockchain.Output
	cldfChain          *cldf_tron.Chain
	DeployerPrivateKey string
}

func (t *Blockchain) ChainSelector() uint64 {
	return t.chainSelector
}
func (t *Blockchain) ChainID() uint64 {
	return t.chainID
}

func (t *Blockchain) CtfOutput() *blockchain.Output {
	return t.ctfOutput
}

func (e *Blockchain) Is(chainFamily string) bool {
	return strings.EqualFold(e.ctfOutput.Family, chainFamily)
}

func (t *Blockchain) Fund(ctx context.Context, address string, amount uint64) error {
	t.testLogger.Info().Msgf("Attempting to fund TRON account %s", address)

	if err := t.lazyInitTronChain(); err != nil {
		return pkgerrors.Wrap(err, "failed to lazy initialize tron chain")
	}

	receiverAddress := tron_addr.EVMAddressToAddress(common.HexToAddress(address))

	tx, err := t.cldfChain.Client.Transfer(t.cldfChain.Address, receiverAddress, libc.MustSafeInt64(amount))
	if err != nil {
		return pkgerrors.Wrapf(err, "failed to create transfer transaction for TRON account %s", address)
	}

	txInfo, err := t.cldfChain.SendAndConfirm(ctx, tx, nil)
	if err != nil {
		return pkgerrors.Wrapf(err, "failed to send and confirm transfer to TRON node %s", address)
	}

	t.testLogger.Info().Msgf("Successfully funded TRON account %s with %d SUN, txHash: %s", receiverAddress.String(), amount, txInfo.ID)

	return nil
}

func (e *Blockchain) ToCldfConfig() (*blockchains.CldfChainConfig, error) {
	return nil, nil
}

func (t *Blockchain) lazyInitTronChain() error {
	if t.cldfChain != nil {
		return nil
	}

	// copied from system-tests/lib/cre/chain.go
	externalHTTPURL := t.ctfOutput.Nodes[0].ExternalHTTPUrl

	signerGen, err := tronprovider.SignerGenCTFDefault()
	if err != nil {
		return fmt.Errorf("failed to create signer generator: %w", err)
	}

	fullNodeURL := strings.Replace(externalHTTPURL, "/jsonrpc", "/wallet", 1)
	solidityNodeURL := strings.Replace(externalHTTPURL, "/jsonrpc", "/walletsolidity", 1)

	tronRPCProvider := tronprovider.NewRPCChainProvider(t.chainSelector, tronprovider.RPCChainProviderConfig{
		FullNodeURL:       fullNodeURL,
		SolidityNodeURL:   solidityNodeURL,
		DeployerSignerGen: signerGen,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	tronChain, err := tronRPCProvider.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize tron chain: %w", err)
	}

	tc, ok := tronChain.(cldf_tron.Chain)
	if !ok {
		return fmt.Errorf("expected cldf_tron.Chain, got %T", tronChain)
	}

	*t.cldfChain = tc

	return nil
}

func (t *Deployer) Deploy(input *blockchain.Input) (blockchains.Blockchain, error) {
	if t.provider.IsCRIB() {
		return nil, errors.New("CRIB deployment for Tron is not supported yet")
	}

	var bcOut *blockchain.Output
	var err error

	if input.Out != nil {
		bcOut = input.Out
	} else {
		bcOut, err = blockchain.NewBlockchainNetwork(input)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
		}
	}

	chainID, err := strconv.ParseUint(bcOut.ChainID, 10, 64)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", bcOut.ChainID)
	}
	selector, err := chainselectors.SelectorFromChainId(chainID)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to get chain selector for chain id %s", bcOut.ChainID)
	}

	// if jsonrpc is not present, add it
	if !strings.HasSuffix(bcOut.Nodes[0].ExternalHTTPUrl, "/jsonrpc") {
		bcOut.Nodes[0].ExternalHTTPUrl += "/jsonrpc"
	}
	if !strings.HasSuffix(bcOut.Nodes[0].InternalHTTPUrl, "/jsonrpc") {
		bcOut.Nodes[0].InternalHTTPUrl += "/jsonrpc"
	}

	return &Blockchain{
		testLogger:         t.testLogger,
		chainSelector:      selector,
		chainID:            chainID,
		ctfOutput:          bcOut,
		DeployerPrivateKey: blockchain.TRONAccounts.PrivateKeys[0],
	}, nil

	// // copied from system-tests/lib/cre/chain.go
	// externalHTTPURL := bcOut.Nodes[0].ExternalHTTPUrl
	// internalHTTPURL := bcOut.Nodes[0].InternalHTTPUrl

	// signerGen, err := tronprovider.SignerGenCTFDefault()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create signer generator: %w", err)
	// }

	// fullNodeURL := strings.Replace(externalHTTPURL, "/jsonrpc", "/wallet", 1)
	// solidityNodeURL := strings.Replace(externalHTTPURL, "/jsonrpc", "/walletsolidity", 1)

	// tronRPCProvider := tronprovider.NewRPCChainProvider(selector, tronprovider.RPCChainProviderConfig{
	// 	FullNodeURL:       fullNodeURL,
	// 	SolidityNodeURL:   solidityNodeURL,
	// 	DeployerSignerGen: signerGen,
	// })
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// tronChain, err := tronRPCProvider.Initialize(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize tron chain: %w", err)
	// }

	// tc, ok := tronChain.(cldf_tron.Chain)
	// if !ok {
	// 	return nil, fmt.Errorf("expected cldf_tron.Chain, got %T", tronChain)
	// }

	// return &cre.Blockchain{
	// 	ChainSelector: selector,
	// 	ChainID:       chainID,
	// 	CtfOutput: &blockchain.Output{
	// 		ChainID: bcOut.ChainID,
	// 		Family:  blockchain.FamilyTron,
	// 		Nodes: []*blockchain.Node{
	// 			{
	// 				InternalHTTPUrl: internalHTTPURL,
	// 				ExternalHTTPUrl: externalHTTPURL,
	// 			},
	// 		},
	// 	},
	// 	SethClient:         nil,
	// 	DeployerPrivateKey: blockchain.TRONAccounts.PrivateKeys[0],

	// 	Funder: &Funder{
	// 		chain:      tc,
	// 		testLogger: t.testLogger,
	// 	},
	// }, nil
}

// type Funder struct {
// 	chain      cldf_tron.Chain
// 	testLogger zerolog.Logger
// }

// func (f *Funder) Fund(ctx context.Context, address string, amount uint64, _ []byte) error {
// 	f.testLogger.Info().Msgf("Attempting to fund TRON account %s", address)

// 	receiverAddress := tron_addr.EVMAddressToAddress(common.HexToAddress(address))

// 	tx, err := f.chain.Client.Transfer(f.chain.Address, receiverAddress, libc.MustSafeInt64(amount))
// 	if err != nil {
// 		return pkgerrors.Wrapf(err, "failed to create transfer transaction for TRON account %s", address)
// 	}

// 	txInfo, err := f.chain.SendAndConfirm(ctx, tx, nil)
// 	if err != nil {
// 		return pkgerrors.Wrapf(err, "failed to send and confirm transfer to TRON node %s", address)
// 	}

// 	f.testLogger.Info().Msgf("Successfully funded TRON account %s with %d SUN, txHash: %s", receiverAddress.String(), amount, txInfo.ID)

// 	return nil
// }

// func (f *Funder) Prepare(ctx context.Context, requiredTotal uint64) ([]byte, error) {
// 	// TRON uses its own built-in funding account, no preparation needed
// 	// return a dummy byte slice to satisfy the interface
// 	return []byte{0}, nil
// }
