package eal

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	eallib "github.com/smartcontractkit/capital-markets-projects/lib/services/eal"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var _ eallib.BlockchainClientInterface = &BlockchainClient{}

type BlockchainClient struct {
	lggr    logger.Logger
	txm     txmgr.TxManager
	gethks  keystore.Eth
	spec    job.EALSpec
	chainID uint64
	cfg     config.EVM
	client  client.Client
}

func NewBlockchainClient(
	lggr logger.Logger,
	txm txmgr.TxManager,
	gethks keystore.Eth,
	spec job.EALSpec,
	chainID uint64,
) (*BlockchainClient, error) {
	return &BlockchainClient{
		lggr:    lggr,
		txm:     txm,
		gethks:  gethks,
		spec:    spec,
		chainID: chainID,
	}, nil
}

func (c *BlockchainClient) SimulateAndCreateTransaction(
	ctx context.Context,
	toAddress gethcommon.Address,
	payload []byte,
) error {
	fromAddresses := c.sendingKeys()
	fromAddress, err := c.gethks.GetRoundRobinAddress(big.NewInt(0).SetUint64(c.chainID), fromAddresses...)
	if err != nil {
		return err
	}

	gasLimit, err := c.client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &toAddress,
		Data: payload,
	})
	// TODO: change errors to constants in eal lib
	if err != nil {
		return errors.Wrap(err, "failed to estimate gas")
	}

	if gasLimit > uint64(c.cfg.GasEstimator().LimitMax()) {
		return errors.New("estimated gas limit exceeds max")
	}

	c.lggr.Debugw("eth_call details",
		"toAddress", toAddress,
		"fromAddress", fromAddress,
		"gasLimit", gasLimit,
		"gasPrice", c.cfg.GasEstimator().PriceMax().ToInt(),
		"GasTipCap", c.cfg.GasEstimator().TipCapDefault().ToInt(),
		"GasFeeCap", c.cfg.GasEstimator().FeeCapDefault().ToInt(),
	)

	_, err = c.client.CallContract(ctx, ethereum.CallMsg{
		To:        &toAddress,
		From:      fromAddress,
		Data:      payload,
		Gas:       gasLimit,
		GasPrice:  c.cfg.GasEstimator().PriceMax().ToInt(),
		GasTipCap: c.cfg.GasEstimator().TipCapDefault().ToInt(),
		GasFeeCap: c.cfg.GasEstimator().FeeCapDefault().ToInt(),
	}, nil /*blocknumber*/)

	if err != nil {
		return err
	}

	ethTx, err := c.txm.CreateTransaction(txmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		FeeLimit:       uint32(gasLimit), // safe down-cast because we cap gas estimator limit max at uint32
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
	}, pg.WithParentCtx(ctx))

	if err != nil {
		return err
	}

	c.lggr.Debugw("created Eth tx", "ethTxID", ethTx.GetID())
	return nil
}

func (c *BlockchainClient) sendingKeys() []gethcommon.Address {
	var addresses []gethcommon.Address
	for _, a := range c.spec.FromAddresses {
		addresses = append(addresses, a.Address())
	}
	return addresses
}

type pipelineResult struct {
	gas uint32
}
