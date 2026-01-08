package evm

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider/rpcclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type Deployer struct {
	provider       infra.Provider
	testLogger     zerolog.Logger
	cribConfigsDir string
}

func NewDeployer(testLogger zerolog.Logger, provider *infra.Provider, cribConfigsDir string) *Deployer {
	return &Deployer{
		provider:       *provider,
		testLogger:     testLogger,
		cribConfigsDir: cribConfigsDir,
	}
}

type Blockchain struct {
	testLogger    zerolog.Logger
	chainSelector uint64
	chainID       uint64
	ctfOutput     *blockchain.Output
	SethClient    *seth.Client
}

func (e *Blockchain) ChainSelector() uint64 {
	return e.chainSelector
}
func (e *Blockchain) ChainID() uint64 {
	return e.chainID
}

func (e *Blockchain) CtfOutput() *blockchain.Output {
	return e.ctfOutput
}

func (e *Blockchain) IsFamily(chainFamily string) bool {
	return strings.EqualFold(e.ctfOutput.Family, chainFamily)
}

func (e *Blockchain) ChainFamily() string {
	return e.ctfOutput.Family
}

func (e *Blockchain) Fund(ctx context.Context, address string, amount uint64) error {
	e.testLogger.Info().Msgf("Attempting to fund EVM account %s", address)

	_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, e.SethClient, libfunding.FundsToSend{
		ToAddress:  common.HexToAddress(address),
		Amount:     big.NewInt(libc.MustSafeInt64(amount)),
		PrivateKey: e.SethClient.MustGetRootPrivateKey(),
	})

	if fundingErr != nil {
		return pkgerrors.Wrapf(fundingErr, "failed to fund node %s", address)
	}
	e.testLogger.Info().Msgf("Successfully funded EVM account %s", address)

	return nil
}

func (e *Blockchain) ToCldfChain() (cldf_chain.BlockChain, error) {
	if len(e.CtfOutput().Nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for chain %s-%d", e.ChainFamily(), e.ChainID())
	}

	chainID := e.ctfOutput.ChainID
	rpcWSURL := e.ctfOutput.Nodes[0].ExternalWSUrl
	rpcHTTPURL := e.ctfOutput.Nodes[0].ExternalHTTPUrl

	d, cErr := chainselectors.GetChainDetailsByChainIDAndFamily(chainID, chainselectors.FamilyEVM)
	if cErr != nil {
		return nil, fmt.Errorf("no chain with ID %s and family %s found: %w", chainID, chainselectors.FamilyEVM, cErr)
	}

	var confirmer cldf_evm_provider.ConfirmFunctor
	switch e.ctfOutput.Type {
	case "anvil":
		confirmer = cldf_evm_provider.ConfirmFuncGeth(3*time.Minute, cldf_evm_provider.WithTickInterval(5*time.Millisecond))
	default:
		confirmer = cldf_evm_provider.ConfirmFuncGeth(3 * time.Minute)
	}

	if keyErr := setDefaultPrivateKeyIfEmpty(); keyErr != nil {
		return nil, keyErr
	}

	provider, pErr := cldf_evm_provider.NewRPCChainProvider(
		d.ChainSelector,
		cldf_evm_provider.RPCChainProviderConfig{
			DeployerTransactorGen: cldf_evm_provider.TransactorFromRaw(
				os.Getenv("PRIVATE_KEY"),
			),
			RPCs: []rpcclient.RPC{
				{
					Name:               "default",
					WSURL:              rpcWSURL,
					HTTPURL:            rpcHTTPURL,
					PreferredURLScheme: rpcclient.URLSchemePreferenceHTTP,
				},
			},
			ConfirmFunctor: confirmer,
		},
	).Initialize(context.Background())

	if pErr != nil {
		return nil, fmt.Errorf("failed to create new chain provider: %w", pErr)
	}

	return provider, nil
}

func (e *Deployer) Deploy(ctx context.Context, input *blockchain.Input) (blockchains.Blockchain, error) {
	var bcOut *blockchain.Output
	var err error

	switch {
	case e.provider.IsCRIB():
		deployCribBlockchainInput := &crib.DeployCribBlockchainInput{
			Blockchain:     input,
			CribConfigsDir: e.cribConfigsDir,
			Namespace:      e.provider.CRIB.Namespace,
		}

		bcOut, err = crib.DeployBlockchain(ctx, deployCribBlockchainInput)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to deploy blockchain")
		}

		err = infra.WaitForRPCEndpoint(e.testLogger, bcOut.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "RPC endpoint is not available")
		}
	case e.provider.IsKubernetes():
		// For Kubernetes, use the blockchain output from config (no deployment)
		if err = blockchains.ValidateKubernetesBlockchainOutput(input); err != nil {
			return nil, err
		}

		e.testLogger.Info().Msgf("Using configured Kubernetes blockchain URLs for %s (chain_id: %s)", input.Type, input.ChainID)
		bcOut = input.Out

		// Wait for RPC endpoint to be available
		err = infra.WaitForRPCEndpoint(e.testLogger, bcOut.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "RPC endpoint is not available")
		}
	default:
		// Docker deployment
		bcOut, err = blockchain.NewWithContext(ctx, input)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
		}
	}

	if keyErr := setDefaultPrivateKeyIfEmpty(); keyErr != nil {
		return nil, keyErr
	}

	priv := os.Getenv("PRIVATE_KEY")
	sethClient, err := seth.NewClientBuilder().
		WithRpcUrl(bcOut.Nodes[0].ExternalWSUrl).
		WithPrivateKeys([]string{priv}).
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create seth client")
	}

	selector, err := chainselectors.SelectorFromChainId(sethClient.Cfg.Network.ChainID)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to get chain selector for chain id %d", sethClient.Cfg.Network.ChainID)
	}

	chainID, err := strconv.ParseUint(bcOut.ChainID, 10, 64)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", bcOut.ChainID)
	}

	return &Blockchain{
		testLogger:    e.testLogger,
		chainSelector: selector,
		chainID:       chainID,
		ctfOutput:     bcOut,
		SethClient:    sethClient,
	}, nil
}

func setDefaultPrivateKeyIfEmpty() error {
	if os.Getenv("PRIVATE_KEY") == "" {
		setErr := os.Setenv("PRIVATE_KEY", blockchain.DefaultAnvilPrivateKey)
		if setErr != nil {
			return fmt.Errorf("failed to set PRIVATE_KEY environment variable: %w", setErr)
		}
		framework.L.Info().Msgf("Set PRIVATE_KEY environment variable to default value: %s", os.Getenv("PRIVATE_KEY"))
	}

	return nil
}
