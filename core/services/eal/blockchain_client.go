package eal

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	libcommon "github.com/smartcontractkit/capital-markets-projects/lib/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

var _ libcommon.BlockchainClientInterface = &BlockchainClient{}

type BlockchainClient struct {
	lggr          logger.Logger
	txm           txmgr.TxManager
	gethks        keystore.Eth
	fromAddresses []ethkey.EIP55Address
	chainID       uint64
	cfg           config.EVM
	client        client.Client
}

func NewBlockchainClient(
	lggr logger.Logger,
	txm txmgr.TxManager,
	gethks keystore.Eth,
	fromAddresses []ethkey.EIP55Address,
	chainID uint64,
	cfg config.EVM,
	client client.Client,
) (*BlockchainClient, error) {
	return &BlockchainClient{
		lggr:          lggr,
		txm:           txm,
		gethks:        gethks,
		fromAddresses: fromAddresses,
		chainID:       chainID,
		cfg:           cfg,
		client:        client,
	}, nil
}

func (c *BlockchainClient) EstimateGas(
	ctx context.Context,
	address gethcommon.Address,
	payload []byte,
) (uint32, error) {
	fromAddresses := c.sendingKeys()
	fromAddress, err := c.gethks.GetRoundRobinAddress(big.NewInt(0).SetUint64(c.chainID), fromAddresses...)
	if err != nil {
		return 0, err
	}
	c.lggr.Debugw("estimate gas details",
		"toAddress", address,
		"fromAddress", fromAddress,
	)
	gasLimit, err := c.client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &address,
		Data: payload,
	})
	// TODO: change errors to constants in eal lib
	if err != nil {
		return 0, errors.Wrap(err, "failed to estimate gas")
	}

	if gasLimit > uint64(c.cfg.GasEstimator().LimitMax()) {
		return 0, errors.New("estimated gas limit exceeds max")
	}
	// safe cast because gas estimator limit max is uint32
	return uint32(gasLimit), nil
}

// SimulateTransaction makes eth_call to simulate transaction
// TODO: look into accepting optional parameters (gas, gasPrice, value)
func (c *BlockchainClient) SimulateTransaction(
	ctx context.Context,
	address gethcommon.Address,
	payload []byte,
	gasLimit uint32,
) error {
	fromAddresses := c.sendingKeys()
	fromAddress, err := c.gethks.GetRoundRobinAddress(big.NewInt(0).SetUint64(c.chainID), fromAddresses...)
	if err != nil {
		return err
	}
	c.lggr.Debugw("eth_call details",
		"toAddress", address,
		"fromAddress", fromAddress,
		"gasLimit", gasLimit,
	)
	_, err = c.client.CallContract(ctx, ethereum.CallMsg{
		To:   &address,
		From: fromAddress,
		Data: payload,
		Gas:  uint64(gasLimit),
	}, nil /*blocknumber*/)

	return err
}

func (c *BlockchainClient) sendingKeys() []gethcommon.Address {
	var addresses []gethcommon.Address
	for _, a := range c.fromAddresses {
		addresses = append(addresses, a.Address())
	}
	return addresses
}
