package offchainreporting

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

var (
	_ ocrtypes.ContractConfigTracker = &OCRContractConfigTracker{}
)

type (
	OCRContractConfigTracker struct {
		ethClient        eth.Client
		contractFilterer *offchainaggregator.OffchainAggregatorFilterer
		contractCaller   *offchainaggregator.OffchainAggregatorCaller
		contractAddress  gethCommon.Address
		logBroadcaster   eth.LogBroadcaster
		jobID            int32
		logger           logger.Logger
	}
)

func NewOCRContractConfigTracker(
	address gethCommon.Address,
	contractFilterer *offchainaggregator.OffchainAggregatorFilterer,
	contractCaller *offchainaggregator.OffchainAggregatorCaller,
	ethClient eth.Client,
	logBroadcaster eth.LogBroadcaster,
	jobID int32,
	logger logger.Logger,
) (o *OCRContractConfigTracker, err error) {
	return &OCRContractConfigTracker{
		ethClient,
		contractFilterer,
		contractCaller,
		address,
		logBroadcaster,
		jobID,
		logger,
	}, nil
}

func (oc *OCRContractConfigTracker) SubscribeToNewConfigs(context.Context) (ocrtypes.ContractConfigSubscription, error) {
	sub := &OCRContractConfigSubscription{
		oc.logger,
		oc.contractAddress,
		make(chan ocrtypes.ContractConfig),
		make(chan ocrtypes.ContractConfig),
		nil,
		nil,
		sync.Mutex{},
		oc,
		sync.Once{},
		make(chan struct{}),
	}
	connected := oc.logBroadcaster.Register(oc.contractAddress, sub)
	if !connected {
		return nil, errors.New("Failed to register with logBroadcaster")
	}
	sub.start()

	return sub, nil
}

func (oc *OCRContractConfigTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	opts := bind.CallOpts{Context: ctx, Pending: false}
	result, err := oc.contractCaller.LatestConfigDetails(&opts)
	if err != nil {
		return 0, configDigest, errors.Wrap(err, "error getting LatestConfigDetails")
	}
	configDigest, err = ocrtypes.BytesToConfigDigest(result.ConfigDigest[:])
	if err != nil {
		return 0, configDigest, errors.Wrap(err, "error getting config digest")
	}
	return uint64(result.BlockNumber), configDigest, err
}

func (oc *OCRContractConfigTracker) ConfigFromLogs(ctx context.Context, changedInBlock uint64) (c ocrtypes.ContractConfig, err error) {
	q := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(changedInBlock)),
		ToBlock:   big.NewInt(int64(changedInBlock)),
		Addresses: []gethCommon.Address{oc.contractAddress},
		Topics: [][]gethCommon.Hash{
			{OCRContractConfigSet},
		},
	}

	logs, err := oc.ethClient.FilterLogs(ctx, q)
	if err != nil {
		return c, err
	}
	if len(logs) == 0 {
		return c, errors.Errorf("ConfigFromLogs: OCRContract with address 0x%x has no logs", oc.contractAddress)
	}

	latest, err := oc.contractFilterer.ParseConfigSet(logs[len(logs)-1])
	if err != nil {
		return c, errors.Wrap(err, "ConfigFromLogs failed to ParseConfigSet")
	}
	latest.Raw = logs[len(logs)-1]
	if latest.Raw.Address != oc.contractAddress {
		return c, errors.Errorf("log address of 0x%x does not match configured contract address of 0x%x", latest.Raw.Address, oc.contractAddress)
	}
	return confighelper.ContractConfigFromConfigSetEvent(*latest), err
}

func (oc *OCRContractConfigTracker) LatestBlockHeight(ctx context.Context) (blockheight uint64, err error) {
	h, err := oc.ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	if h == nil {
		return 0, errors.New("got nil head")
	}

	return uint64(h.Number), nil
}
