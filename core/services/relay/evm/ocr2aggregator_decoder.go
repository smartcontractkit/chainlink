package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ LogDecoder = &ocr2AggregatorLogDecoder{}

type ocr2AggregatorLogDecoder struct {
	eventName string
	eventSig  common.Hash
	abi       *abi.ABI
}

func newOCR2AggregatorLogDecoder() (*ocr2AggregatorLogDecoder, error) {
	const eventName = "ConfigSet"
	abi, err := ocr2aggregator.OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return &ocr2AggregatorLogDecoder{
		eventName: eventName,
		eventSig:  abi.Events[eventName].ID,
		abi:       abi,
	}, nil
}

func (d *ocr2AggregatorLogDecoder) Decode(rawLog []byte) (ocrtypes.ContractConfig, error) {
	unpacked := new(ocr2aggregator.OCR2AggregatorConfigSet)
	err := d.abi.UnpackIntoInterface(unpacked, d.eventName, rawLog)
	if err != nil {
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to unpack log data: %w", err)
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.Transmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(addr.Hex()))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range unpacked.Signers {
		addr := addr
		signers = append(signers, addr[:])
	}

	return ocrtypes.ContractConfig{
		ConfigDigest:          unpacked.ConfigDigest,
		ConfigCount:           unpacked.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitAccounts,
		F:                     unpacked.F,
		OnchainConfig:         unpacked.OnchainConfig,
		OffchainConfigVersion: unpacked.OffchainConfigVersion,
		OffchainConfig:        unpacked.OffchainConfig,
	}, nil
}

func (d *ocr2AggregatorLogDecoder) EventSig() common.Hash {
	return d.eventSig
}
