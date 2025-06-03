package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocrconfigurationstoreevmsimple"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
)

var _ LogDecoder = &ocrConfigurationStoreEVMSimpleLogDecoder{}

type ocrConfigurationStoreEVMSimpleLogDecoder struct {
	eventName string
	eventSig  common.Hash
	abi       *abi.ABI
	chain     legacyevm.Chain
	address   common.Address
}

func newOCRConfigurationStoreEVMSimpleLogDecoder(chain legacyevm.Chain, address common.Address) (*ocrConfigurationStoreEVMSimpleLogDecoder, error) {
	const eventName = "NewConfiguration"
	abi, err := ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return &ocrConfigurationStoreEVMSimpleLogDecoder{
		eventName: eventName,
		eventSig:  abi.Events[eventName].ID,
		abi:       abi,
		chain:     chain,
		address:   address,
	}, nil
}

func (d *ocrConfigurationStoreEVMSimpleLogDecoder) Decode(rawLog []byte) (ocrtypes.ContractConfig, error) {
	d.chain.Logger().Infof("TRACE Decoding log for event %s on contract %s", d.eventName, d.address.Hex())

	// Convert rawLog bytes into a types.Log
	var logEvent types.Log
	if err := rlp.DecodeBytes(rawLog, &logEvent); err != nil {
		d.chain.Logger().Errorf("TRACE Failed to decode raw log into types.Log for event %s on contract %s: %v", d.eventName, d.address.Hex(), err)
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to decode raw log: %w", err)
	}

	// Unpack the non-indexed data from logEvent.Data
	unpacked := new(ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleNewConfiguration)
	if err := d.abi.UnpackIntoInterface(unpacked, d.eventName, logEvent.Data); err != nil {
		d.chain.Logger().Errorf("TRACE Failed to unpack log for event %s on contract %s: %v", d.eventName, d.address.Hex(), err)
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to unpack log data: %w", err)
	}

	// Pick up the indexed fields from the log topics.
	var indexed abi.Arguments
	for _, arg := range d.abi.Events[d.eventName].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(unpacked, indexed, logEvent.Topics[1:]); err != nil {
		d.chain.Logger().Errorf("TRACE Failed to parse indexed topics for event %s on contract %s: %v", d.eventName, d.address.Hex(), err)
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to parse topics: %w", err)
	}

	if unpacked.ConfigDigest == (common.Hash{}) {
		d.chain.Logger().Errorf("TRACE ConfigDigest is empty for event %s on contract %s, %v", d.eventName, d.address.Hex(), unpacked)
		return ocrtypes.ContractConfig{}, fmt.Errorf("config digest is empty for event %s on contract %s", d.eventName, d.address.Hex())
	}

	// Create contract caller instance to read the full configuration.
	configStore, err := ocrconfigurationstoreevmsimple.NewOCRConfigurationStoreEVMSimpleCaller(d.address, d.chain.Client())
	if err != nil {
		d.chain.Logger().Errorf("TRACE Failed to create contract caller for event %s on contract %s: %v", d.eventName, d.address.Hex(), err)
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to create contract caller: %w", err)
	}

	// Read the full configuration using the config digest from the event.
	d.chain.Logger().Errorf("TRACE reading config from contract %s for digest %s", d.address.Hex(), fmt.Sprintf("0x%x", unpacked.ConfigDigest))
	configData, err := configStore.ReadConfig(nil, unpacked.ConfigDigest)
	if err != nil {
		d.chain.Logger().Errorf("TRACE Failed to read config from contract %s for digest %s: %v", d.address.Hex(), unpacked.ConfigDigest, err)
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to read config from contract: %w", err)
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range configData.Transmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(addr.Hex()))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range configData.Signers {
		addr := addr
		signers = append(signers, addr[:])
	}

	d.chain.Logger().Infof("TRACE Successfully decoded log for event %s on contract %s", d.eventName, d.address.Hex())

	return ocrtypes.ContractConfig{
		ConfigDigest:          unpacked.ConfigDigest,
		ConfigCount:           uint64(configData.ConfigCount),
		Signers:               signers,
		Transmitters:          transmitAccounts,
		F:                     configData.F,
		OnchainConfig:         configData.OnchainConfig,
		OffchainConfigVersion: configData.OffchainConfigVersion,
		OffchainConfig:        configData.OffchainConfig,
	}, nil
}

func (d *ocrConfigurationStoreEVMSimpleLogDecoder) EventSig() common.Hash {
	return d.eventSig
}
