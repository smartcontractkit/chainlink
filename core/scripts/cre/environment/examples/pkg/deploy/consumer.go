package deploy

import (
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/contracts/permissionless_feeds_consumer"
)

func DeployPermissionlessFeedsConsumer(rpcUrl string) (*common.Address, error) {
	var privateKey string
	if os.Getenv("PRIVATE_KEY") == "" {
		privateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		fmt.Printf("Since PRIVATE_KEY environment variable was empty, will use default value: %s\n", privateKey)
	} else {
		privateKey = os.Getenv("PRIVATE_KEY")
	}

	sethClient, sethErr := seth.NewClientBuilder().
		WithRpcUrl(rpcUrl).
		WithPrivateKeys([]string{privateKey}).
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
