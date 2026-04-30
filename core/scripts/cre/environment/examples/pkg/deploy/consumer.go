package deploy

import (
	"context"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/balance_reader"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/contracts/permissionless_feeds_consumer"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
)

func PermissionlessFeedsConsumer(rpcURL string) (*common.Address, error) {
	if pkErr := environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return nil, pkErr
	}

	sethClient, sethErr := seth.NewClientBuilder().
		WithRpcUrl(rpcURL).
		WithPrivateKeys([]string{os.Getenv("PRIVATE_KEY")}).
		// do not check if there's a pending nonce nor check node's health
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if sethErr != nil {
		return nil, errors.Wrap(sethErr, "failed to create Seth Ethereum client")
	}

	consABI, abiErr := permissionless_feeds_consumer.PermissionlessFeedsConsumerMetaData.GetAbi()
	if abiErr != nil {
		return nil, errors.Wrap(abiErr, "failed to get Permissionless Feeds Consumer contract ABI")
	}

	data, deployErr := sethClient.DeployContract(sethClient.NewTXOpts(), "PermissionlessFeedsConsumer", *consABI, common.FromHex(permissionless_feeds_consumer.PermissionlessFeedsConsumerMetaData.Bin))
	if deployErr != nil {
		return nil, errors.Wrap(deployErr, "failed to deploy Permissionless Feeds Consumer contract")
	}

	return &data.Address, nil
}

func BalanceReader(rpcURL string) (*common.Address, error) {
	if pkErr := environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return nil, pkErr
	}

	sethClient, sethErr := seth.NewClientBuilder().
		WithRpcUrl(rpcURL).
		WithPrivateKeys([]string{os.Getenv("PRIVATE_KEY")}).
		// do not check if there's a pending nonce nor check node's health
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if sethErr != nil {
		return nil, errors.Wrap(sethErr, "failed to create Seth Ethereum client")
	}

	contractABI, abiErr := balance_reader.BalanceReaderMetaData.GetAbi()
	if abiErr != nil {
		return nil, errors.Wrap(abiErr, "failed to get Balance Reader contract ABI")
	}

	data, deployErr := sethClient.DeployContract(sethClient.NewTXOpts(), "BalanceReader", *contractABI, common.FromHex(balance_reader.BalanceReaderMetaData.Bin))
	if deployErr != nil {
		return nil, errors.Wrap(deployErr, "failed to deploy Balance Reader contract")
	}

	return &data.Address, nil
}

func ChainMetadata(rpcURL string) (uint64, uint64, error) {
	sethClient, err := newSethClient(rpcURL)
	if err != nil {
		return 0, 0, err
	}

	chainID := sethClient.Cfg.Network.ChainID
	chainSelector, err := chainselectors.SelectorFromChainId(chainID)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to resolve chain selector for chain id %d", chainID)
	}

	return chainID, chainSelector, nil
}

func CreateAndFundAddresses(rpcURL string, count int, amount *big.Int) ([]common.Address, error) {
	if count <= 0 {
		return nil, errors.New("count must be greater than zero")
	}
	if amount == nil {
		return nil, errors.New("amount is nil")
	}

	sethClient, err := newSethClient(rpcURL)
	if err != nil {
		return nil, err
	}

	addresses := make([]common.Address, 0, count)
	for range count {
		address, _, err := crecrypto.GenerateNewKeyPair()
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate address")
		}

		_, err = libfunding.SendFunds(context.Background(), framework.L, sethClient, libfunding.FundsToSend{
			ToAddress:  address,
			Amount:     new(big.Int).Set(amount),
			PrivateKey: sethClient.MustGetRootPrivateKey(),
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to fund address %s", address.Hex())
		}

		addresses = append(addresses, address)
	}

	return addresses, nil
}

func newSethClient(rpcURL string) (*seth.Client, error) {
	if pkErr := environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return nil, pkErr
	}

	sethClient, sethErr := seth.NewClientBuilder().
		WithRpcUrl(rpcURL).
		WithPrivateKeys([]string{os.Getenv("PRIVATE_KEY")}).
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if sethErr != nil {
		return nil, errors.Wrap(sethErr, "failed to create Seth Ethereum client")
	}

	return sethClient, nil
}
