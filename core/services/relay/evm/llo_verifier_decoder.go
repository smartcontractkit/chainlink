package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_verifier"
)

var _ LogDecoder = &channelVerifierLogDecoder{}

type channelVerifierLogDecoder struct {
	eventName string
	eventSig  common.Hash
	abi       *abi.ABI
}

func newChannelVerifierLogDecoder() (*channelVerifierLogDecoder, error) {
	const eventName = "ConfigSet"
	abi, err := channel_verifier.ChannelVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return &channelVerifierLogDecoder{
		eventName: eventName,
		eventSig:  abi.Events[eventName].ID,
		abi:       abi,
	}, nil
}

func (d *channelVerifierLogDecoder) Decode(rawLog []byte) (ocrtypes.ContractConfig, error) {
	unpacked := new(channel_verifier.ChannelVerifierConfigSet)
	err := d.abi.UnpackIntoInterface(unpacked, d.eventName, rawLog)
	if err != nil {
		return ocrtypes.ContractConfig{}, errors.Wrap(err, "failed to unpack log data")
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.OffchainTransmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(fmt.Sprintf("%x", addr)))
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

func (d *channelVerifierLogDecoder) EventSig() common.Hash {
	return d.eventSig
}
