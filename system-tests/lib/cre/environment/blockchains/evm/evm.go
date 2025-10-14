package evm

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
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

func (e *Deployer) Deploy(input *blockchain.Input) (*cre.Blockchain, error) {
	var bcOut *blockchain.Output
	var err error

	if e.provider.IsCRIB() {
		deployCribBlockchainInput := &cre.DeployCribBlockchainInput{
			BlockchainInput: input,
			CribConfigsDir:  e.cribConfigsDir,
			Namespace:       e.provider.CRIB.Namespace,
		}

		bcOut, err = crib.DeployBlockchain(deployCribBlockchainInput)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to deploy blockchain")
		}

		err = infra.WaitForRPCEndpoint(e.testLogger, bcOut.Nodes[0].ExternalHTTPUrl, 10*time.Minute)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "RPC endpoint is not available")
		}
	} else {
		bcOut, err = blockchain.NewBlockchainNetwork(input)
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

	return &cre.Blockchain{
		ChainSelector:      selector,
		ChainID:            chainID,
		CtfOutput:          bcOut,
		SethClient:         sethClient,
		DeployerPrivateKey: priv,

		Funder: &Funder{
			sethClient: sethClient,
			testLogger: e.testLogger,
		},
	}, nil
}

type Funder struct {
	sethClient *seth.Client
	testLogger zerolog.Logger
}

func (f *Funder) Fund(ctx context.Context, address string, amount uint64, fundingPrivateKey []byte) error {
	f.testLogger.Info().Msgf("Attempting to fund EVM account %s", address)

	fundingKey, fkErr := crypto.ToECDSA(fundingPrivateKey)
	if fkErr != nil {
		return pkgerrors.Wrap(fkErr, "failed to convert funding private key to ECDSA")
	}

	_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, f.sethClient, libfunding.FundsToSend{
		ToAddress:  common.HexToAddress(address),
		Amount:     big.NewInt(libc.MustSafeInt64(amount)),
		PrivateKey: fundingKey,
	})

	if fundingErr != nil {
		return pkgerrors.Wrapf(fundingErr, "failed to fund node %s", address)
	}
	f.testLogger.Info().Msgf("Successfully funded EVM account %s", address)

	return nil
}

func (f *Funder) Prepare(ctx context.Context, requiredTotal uint64) ([]byte, error) {
	publicAddress, privateKeyBytes, pkErr := generatePubPrivKeyPair()
	if pkErr != nil {
		return nil, pkgerrors.Wrap(pkErr, "failed to generate pub/priv key pair for EVM funding account")
	}

	fundingAmount := libc.MustSafeInt64(requiredTotal)
	fundingAmount += (fundingAmount / 5) // add 20% to cover gas fees

	_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, f.sethClient, libfunding.FundsToSend{
		ToAddress:  *publicAddress,
		Amount:     big.NewInt(fundingAmount),
		PrivateKey: f.sethClient.MustGetRootPrivateKey(),
	})

	if fundingErr != nil {
		return nil, pkgerrors.Wrapf(fundingErr, "failed to fund funding account %s on chain %d", publicAddress.String(), f.sethClient.ChainID)
	}

	return privateKeyBytes, nil
}

func generatePubPrivKeyPair() (*common.Address, []byte, error) {
	privateKey, pkErr := crypto.GenerateKey()
	if pkErr != nil {
		return nil, nil, pkgerrors.Wrap(pkErr, "failed to generate private key for funding accounts")
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("error casting public key to ECDSA")
	}
	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &publicAddress, privateKeyBytes, nil
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
