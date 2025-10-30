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

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
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

func (t *Blockchain) IsFamily(chainFamily string) bool {
	return strings.EqualFold(t.ctfOutput.Family, chainFamily)
}

func (t *Blockchain) ChainFamily() string {
	return t.ctfOutput.Family
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

func (t *Blockchain) ToCldfChain() (cldf_chain.BlockChain, error) {
	if err := t.lazyInitTronChain(); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to lazy initialize tron chain")
	}

	return t.cldfChain, nil
}

func (t *Blockchain) lazyInitTronChain() error {
	if t.cldfChain != nil {
		return nil
	}

	if len(t.CtfOutput().Nodes) == 0 {
		return fmt.Errorf("no nodes found for chain %s-%d", t.ChainFamily(), t.ChainID())
	}

	// tron's devnet chainID maps to many chain selectors, one for tron one for EVM
	// we want to force mapping to EVM family here to avoid selector mismatches later
	chainDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(t.ChainID(), 10), chainselectors.FamilyEVM)
	if err != nil {
		return fmt.Errorf("failed to get selector from chain id %d: %w", t.ChainID(), err)
	}

	signerGen, err := tronprovider.SignerGenCTFDefault()
	if err != nil {
		return fmt.Errorf("failed to create signer generator: %w", err)
	}

	externalHTTPURL := t.CtfOutput().Nodes[0].ExternalHTTPUrl
	fullNodeURL := strings.Replace(externalHTTPURL, "/jsonrpc", "/wallet", 1)
	solidityNodeURL := strings.Replace(externalHTTPURL, "/jsonrpc", "/walletsolidity", 1)

	tronRPCProvider := tronprovider.NewRPCChainProvider(chainDetails.ChainSelector, tronprovider.RPCChainProviderConfig{
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

	t.cldfChain = &tc

	return nil
}

func (t *Deployer) Deploy(ctx context.Context, input *blockchain.Input) (blockchains.Blockchain, error) {
	if t.provider.IsCRIB() {
		return nil, errors.New("CRIB deployment for Tron is not supported yet")
	}

	var bcOut *blockchain.Output
	var err error

	if input.Out != nil {
		bcOut = input.Out
	} else {
		bcOut, err = blockchain.NewWithContext(ctx, input)
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
}
