package evmutil

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ types.OffchainConfigDigester = EVMOffchainConfigDigester{}

type EVMOffchainConfigDigester struct {
	ChainID         uint64
	ContractAddress common.Address
}

func (d EVMOffchainConfigDigester) ConfigDigest(cc types.ContractConfig) (types.ConfigDigest, error) {
	signers := []common.Address{}
	for i, signer := range cc.Signers {
		if len(signer) != 20 {
			return types.ConfigDigest{}, fmt.Errorf("%v-th evm signer should be a 20 byte address, but got %x", i, signer)
		}
		a := common.BytesToAddress(signer)
		signers = append(signers, a)
	}
	transmitters := []common.Address{}
	for i, transmitter := range cc.Transmitters {
		if !strings.HasPrefix(string(transmitter), "0x") || len(transmitter) != 42 || !common.IsHexAddress(string(transmitter)) {
			return types.ConfigDigest{}, fmt.Errorf("%v-th evm transmitter should be a 42 character Ethereum address string, but got '%v'", i, transmitter)
		}
		a := common.HexToAddress(string(transmitter))
		transmitters = append(transmitters, a)
	}

	return configDigest(
		d.ChainID,
		d.ContractAddress,
		cc.ConfigCount,
		signers,
		transmitters,
		cc.F,
		cc.OnchainConfig,
		cc.OffchainConfigVersion,
		cc.OffchainConfig,
	), nil
}

func (d EVMOffchainConfigDigester) ConfigDigestPrefix() (types.ConfigDigestPrefix, error) {
	return types.ConfigDigestPrefixEVMSimple, nil
}
