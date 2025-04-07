package internal

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	feeds_consumer "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	"github.com/smartcontractkit/chainlink/deployment"
)

type KeystoneFeedsConsumerDeployer struct {
	lggr     logger.Logger
	contract *feeds_consumer.KeystoneFeedsConsumer
}

func NewKeystoneFeedsConsumerDeployer() (*KeystoneFeedsConsumerDeployer, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerDeployer{lggr: lggr}, nil
}

func (c *KeystoneFeedsConsumerDeployer) deploy(req DeployRequest) (*DeployResponse, error) {
	est, err := estimateDeploymentGas(req.Chain.Client, feeds_consumer.KeystoneFeedsConsumerABI)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}
	c.lggr.Debugf("Feeds Consumer estimated gas: %d", est)

	consumerAddr, tx, consumer, err := feeds_consumer.DeployKeystoneFeedsConsumer(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy feeds consumer: %w", err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm and save feeds consumer: %w", err)
	}

	tv := deployment.TypeAndVersion{
		Type: FeedConsumer,
	}

	resp := &DeployResponse{
		Address: consumerAddr,
		Tx:      tx.Hash(),
		Tv:      tv,
	}
	c.contract = consumer
	return resp, nil
}
