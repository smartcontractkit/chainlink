package evm

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_evm_client "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider/rpcclient"
	cldf_chain_utils "github.com/smartcontractkit/chainlink-deployments-framework/chain/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
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
	chainDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(e.ChainID(), 10), e.ctfOutput.Family)
	if err != nil {
		return nil, fmt.Errorf("failed to get selector from chain id %d: %w", e.ChainID(), err)
	}

	if len(e.CtfOutput().Nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for chain %s-%d", e.ChainFamily(), e.ChainID())
	}

	rpcs := []cldf_evm_client.RPC{}
	for i, node := range e.CtfOutput().Nodes {
		rpcs = append(rpcs, cldf_evm_client.RPC{
			Name:    fmt.Sprintf("%s-%d", e.ctfOutput.Family, i),
			WSURL:   node.ExternalWSUrl,
			HTTPURL: node.ExternalHTTPUrl,
		})
	}

	rpcConf := cldf_evm_client.RPCConfig{
		ChainSelector: chainDetails.ChainSelector,
		RPCs:          rpcs,
	}

	ec, evmErr := cldf_evm_client.NewMultiClient(logger.Nop(), rpcConf)
	if evmErr != nil {
		return nil, fmt.Errorf("failed to create multi client: %w", evmErr)
	}

	chainInfo, infoErr := cldf_chain_utils.ChainInfo(chainDetails.ChainSelector)
	if infoErr != nil {
		return nil, fmt.Errorf("failed to get chain info for chain %s-%d: %w", e.ChainFamily(), e.ChainID(), infoErr)
	}

	confirmFn := func(tx *types.Transaction) (uint64, error) {
		var blockNumber uint64
		if tx == nil {
			return 0, fmt.Errorf("tx was nil, nothing to confirm chain %s", chainInfo.ChainName)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		receipt, rErr := bind.WaitMined(ctx, ec, tx)
		if rErr != nil {
			return blockNumber, fmt.Errorf("failed to get confirmed receipt for chain %s: %w", chainInfo.ChainName, rErr)
		}
		if receipt == nil {
			return blockNumber, fmt.Errorf("receipt was nil for tx %s chain %s", tx.Hash().Hex(), chainInfo.ChainName)
		}
		blockNumber = receipt.BlockNumber.Uint64()
		if receipt.Status == 0 {
			errReason, rErr := deployment.GetErrorReasonFromTx(ec, e.SethClient.MustGetRootKeyAddress(), tx, receipt)
			if rErr == nil && errReason != "" {
				return blockNumber, fmt.Errorf("tx %s reverted,error reason: %s chain %s", tx.Hash().Hex(), errReason, chainInfo.ChainName)
			}
			return blockNumber, fmt.Errorf("tx %s reverted, could not decode error reason chain %s", tx.Hash().Hex(), chainInfo.ChainName)
		}
		return blockNumber, nil
	}

	return cldf_evm.Chain{
		Selector:    chainDetails.ChainSelector,
		Client:      ec,
		DeployerKey: e.SethClient.NewTXOpts(seth.WithNonce(nil)), // ensure nonce fetched from chain at use time
		Confirm:     confirmFn,
	}, nil
}

func (e *Deployer) Deploy(ctx context.Context, input *blockchain.Input) (blockchains.Blockchain, error) {
	var bcOut *blockchain.Output
	var err error

	if e.provider.IsCRIB() {
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
	} else {
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
