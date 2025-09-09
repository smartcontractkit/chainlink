package deploy

import (
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/balance_reader"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/contracts/permissionless_feeds_consumer"

	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
)

func PermissionlessFeedsConsumer(rpcURL string) (*common.Address, error) {
	if pkErr := creenv.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
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
	if pkErr := creenv.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
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
