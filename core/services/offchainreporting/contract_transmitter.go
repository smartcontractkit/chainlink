package offchainreporting

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

var (
	_ ocrtypes.ContractTransmitter = &OCRContractTransmitter{}
)

type (
	OCRContractTransmitter struct {
		contractAddress gethCommon.Address
		contractABI     abi.ABI
		transmitter     Transmitter
		contractCaller  *offchainaggregator.OffchainAggregatorCaller
	}

	Transmitter interface {
		CreateEthTransaction(ctx context.Context, toAddress gethCommon.Address, payload []byte) error
		FromAddress() gethCommon.Address
	}
)

func NewOCRContractTransmitter(
	address gethCommon.Address,
	contractCaller *offchainaggregator.OffchainAggregatorCaller,
	contractABI abi.ABI,
	transmitter Transmitter,
) *OCRContractTransmitter {
	return &OCRContractTransmitter{
		contractAddress: address,
		contractABI:     contractABI,
		transmitter:     transmitter,
		contractCaller:  contractCaller,
	}
}

func (oc *OCRContractTransmitter) Transmit(ctx context.Context, report []byte, rs, ss [][32]byte, vs [32]byte) error {
	payload, err := oc.contractABI.Pack("transmit", report, rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	return errors.Wrap(oc.transmitter.CreateEthTransaction(ctx, oc.contractAddress, payload), "failed to send Eth transaction")
}

func (oc *OCRContractTransmitter) LatestTransmissionDetails(ctx context.Context) (configDigest ocrtypes.ConfigDigest, epoch uint32, round uint8, latestAnswer ocrtypes.Observation, latestTimestamp time.Time, err error) {
	opts := bind.CallOpts{Context: ctx, Pending: false}
	result, err := oc.contractCaller.LatestTransmissionDetails(&opts)
	if err != nil {
		return configDigest, 0, 0, ocrtypes.Observation(nil), time.Time{}, errors.Wrap(err, "error getting LatestTransmissionDetails")
	}
	return result.ConfigDigest, result.Epoch, result.Round, ocrtypes.Observation(result.LatestAnswer), time.Unix(int64(result.LatestTimestamp), 0), nil
}

func (oc *OCRContractTransmitter) FromAddress() gethCommon.Address {
	return oc.transmitter.FromAddress()
}
