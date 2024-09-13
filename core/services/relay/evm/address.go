package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func AccountToAddress(accounts []types.Account) (addresses []common.Address, err error) {
	for _, signer := range accounts {
		bytes, err := hexutil.Decode(string(signer))
		if err != nil {
			return []common.Address{}, fmt.Errorf("%w: given address is not valid %s", err, signer)
		}
		if len(bytes) != 20 {
			return []common.Address{}, fmt.Errorf("address is not 20 bytes %s", signer)
		}
		addresses = append(addresses, common.BytesToAddress(bytes))
	}
	return addresses, nil
}

func OnchainPublicKeyToAddress(publicKeys []types.OnchainPublicKey) (addresses []common.Address, err error) {
	for _, signer := range publicKeys {
		if len(signer) != 20 {
			return []common.Address{}, fmt.Errorf("address is not 20 bytes %s", signer)
		}
		addresses = append(addresses, common.BytesToAddress(signer))
	}
	return addresses, nil
}
