package offchainreporting

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

var (
	OCRContractConfigSet = getConfigSetHash()
)

type (
	OCRContract struct {
		ethClient        eth.Client
		contractFilterer *offchainaggregator.OffchainAggregatorFilterer
		contractCaller   *offchainaggregator.OffchainAggregatorCaller
		contractAddress  gethCommon.Address
		logBroadcaster   eth.LogBroadcaster
		jobID            int32
		transmitter      Transmitter
		contractABI      abi.ABI
		logger           logger.Logger
	}

	Transmitter interface {
		CreateEthTransaction(ctx context.Context, toAddress gethCommon.Address, payload []byte) error
		FromAddress() gethCommon.Address
	}
)

var (
	_ ocrtypes.ContractConfigTracker      = &OCRContract{}
	_ ocrtypes.ContractTransmitter        = &OCRContract{}
	_ ocrtypes.ContractConfigSubscription = &OCRContractConfigSubscription{}
	_ eth.LogListener                     = &OCRContractConfigSubscription{}
)

func NewOCRContract(address gethCommon.Address, ethClient eth.Client, logBroadcaster eth.LogBroadcaster, jobID int32, transmitter Transmitter, logger logger.Logger) (o *OCRContract, err error) {
	contractFilterer, err := offchainaggregator.NewOffchainAggregatorFilterer(address, ethClient)
	if err != nil {
		return o, errors.Wrap(err, "could not instantiate NewOffchainAggregatorFilterer")
	}

	contractCaller, err := offchainaggregator.NewOffchainAggregatorCaller(address, ethClient)
	if err != nil {
		return o, errors.Wrap(err, "could not instantiate NewOffchainAggregatorCaller")
	}

	contractABI, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
	if err != nil {
		return o, errors.Wrap(err, "could not get contract ABI JSON")
	}

	return &OCRContract{
		ethClient,
		contractFilterer,
		contractCaller,
		address,
		logBroadcaster,
		jobID,
		transmitter,
		contractABI,
		logger,
	}, nil
}

func (oc *OCRContract) Transmit(ctx context.Context, report []byte, rs, ss [][32]byte, vs [32]byte) error {
	payload, err := oc.contractABI.Pack("transmit", report, rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	return errors.Wrap(oc.transmitter.CreateEthTransaction(ctx, oc.contractAddress, payload), "failed to send Eth transaction")
}

func (oc *OCRContract) SubscribeToNewConfigs(context.Context) (ocrtypes.ContractConfigSubscription, error) {
	sub := &OCRContractConfigSubscription{
		oc.logger,
		oc.contractAddress,
		make(chan ocrtypes.ContractConfig),
		oc,
		sync.Once{},
	}
	connected := oc.logBroadcaster.Register(oc.contractAddress, sub)
	if !connected {
		return nil, errors.New("Failed to register with logBroadcaster")
	}

	return sub, nil
}

func (oc *OCRContract) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	opts := bind.CallOpts{Context: ctx, Pending: false}
	result, err := oc.contractCaller.LatestConfigDetails(&opts)
	if err != nil {
		return 0, configDigest, errors.Wrap(err, "error getting LatestConfigDetails")
	}
	return uint64(result.BlockNumber), ocrtypes.BytesToConfigDigest(result.ConfigDigest[:]), err
}

func (oc *OCRContract) ConfigFromLogs(ctx context.Context, changedInBlock uint64) (c ocrtypes.ContractConfig, err error) {
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

func (oc *OCRContract) LatestBlockHeight(ctx context.Context) (blockheight uint64, err error) {
	h, err := oc.ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	if h == nil {
		return 0, errors.New("got nil head")
	}

	return uint64(h.Number), nil
}

func (oc *OCRContract) LatestTransmissionDetails(ctx context.Context) (configDigest ocrtypes.ConfigDigest, epoch uint32, round uint8, latestAnswer ocrtypes.Observation, latestTimestamp time.Time, err error) {
	opts := bind.CallOpts{Context: ctx, Pending: false}
	result, err := oc.contractCaller.LatestTransmissionDetails(&opts)
	if err != nil {
		return configDigest, 0, 0, ocrtypes.Observation(nil), time.Time{}, errors.Wrap(err, "error getting LatestTransmissionDetails")
	}
	return result.ConfigDigest, result.Epoch, result.Round, ocrtypes.Observation(result.LatestAnswer), time.Unix(int64(result.LatestTimestamp), 0), nil
}

func (oc *OCRContract) FromAddress() gethCommon.Address {
	return oc.transmitter.FromAddress()
}

const OCRContractConfigSubscriptionHandleLogTimeout = 5 * time.Second

type OCRContractConfigSubscription struct {
	logger          logger.Logger
	contractAddress gethCommon.Address
	ch              chan ocrtypes.ContractConfig
	oc              *OCRContract
	closer          sync.Once
}

// OnConnect complies with LogListener interface
func (sub *OCRContractConfigSubscription) OnConnect() {}

// OnDisconnect complies with LogListener interface
func (sub *OCRContractConfigSubscription) OnDisconnect() {}

// HandleLog complies with LogListener interface
func (sub *OCRContractConfigSubscription) HandleLog(lb eth.LogBroadcast, err error) {
	topics := lb.RawLog().Topics
	if len(topics) == 0 {
		return
	}
	switch topics[0] {
	case OCRContractConfigSet:
		raw := lb.RawLog()
		if raw.Address != sub.contractAddress {
			sub.logger.Errorf("log address of 0x%x does not match configured contract address of 0x%x", raw.Address, sub.contractAddress)
			return
		}
		configSet, err := sub.oc.contractFilterer.ParseConfigSet(raw)
		if err != nil {
			sub.logger.Errorw("could not parse config set", "err", err)
			return
		}
		configSet.Raw = lb.RawLog()
		cc := confighelper.ContractConfigFromConfigSetEvent(*configSet)
		select {
		// NOTE: This is thread-safe because HandleLog cannot be called concurrently with Unregister due to the design of LogBroadcaster
		// It will never send on closed channel
		case sub.ch <- cc:
		case <-time.After(OCRContractConfigSubscriptionHandleLogTimeout):
			sub.logger.Error("OCRContractConfigSubscription HandleLog timed out waiting on receive channel")
		}
	default:
	}
}

// IsV2Job complies with LogListener interface
func (sub *OCRContractConfigSubscription) IsV2Job() bool {
	return true
}

// JobIDV2 complies with LogListener interface
func (sub *OCRContractConfigSubscription) JobIDV2() int32 {
	return sub.oc.jobID
}

// JobID complies with LogListener interface
func (sub *OCRContractConfigSubscription) JobID() *models.ID {
	return &models.ID{}
}

// Configs complies with ContractConfigSubscription interface
func (sub *OCRContractConfigSubscription) Configs() <-chan ocrtypes.ContractConfig {
	return sub.ch
}

// Close complies with ContractConfigSubscription interface
func (sub *OCRContractConfigSubscription) Close() {
	sub.closer.Do(func() {
		sub.oc.logBroadcaster.Unregister(sub.oc.contractAddress, sub)

		close(sub.ch)
	})
}

func getConfigSetHash() gethCommon.Hash {
	abi, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
	if err != nil {
		panic("could not parse OffchainAggregator ABI: " + err.Error())
	}
	return abi.Events["ConfigSet"].ID
}
