package metatx

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateMetaTransferCalldata(receiver common.Address, amount *big.Int, chainID uint64) ([]byte, [32]byte, error) {
	calldataDefinition := `
	[
		{
			"inputs": [{
				"internalType": "address",
				"name": "receiver",
				"type": "address"
			}, {
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}, {
				"internalType": "uint64",
				"name": "destinationChainId",
				"type": "uint64"
			}],
			"name": "metaTransfer",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]
	`

	calldataAbi, err := abi.JSON(strings.NewReader(calldataDefinition))
	if err != nil {
		return nil, [32]byte{}, err
	}

	calldata, err := calldataAbi.Pack("metaTransfer", receiver, amount, chainID)
	if err != nil {
		return nil, [32]byte{}, err
	}

	calldataHashRaw := crypto.Keccak256(calldata)

	var calldataHash [32]byte
	copy(calldataHash[:], calldataHashRaw[:])

	return calldata, calldataHash, nil
}
