package txm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type DummyKeystore struct {
	chainID       *big.Int
	privateKeyMap map[common.Address]*ecdsa.PrivateKey
}

func NewKeystore(chainID *big.Int) *DummyKeystore {
	return &DummyKeystore{chainID: chainID, privateKeyMap: make(map[common.Address]*ecdsa.PrivateKey)}
}

func (k *DummyKeystore) Add(privateKeyString string) error {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key: %v to ECDSA", publicKey)
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	k.privateKeyMap[address] = privateKey
	return nil
}

func (k *DummyKeystore) SignTx(_ context.Context, fromAddress common.Address, tx *types.Transaction) (*types.Transaction, error) {
	if key, exists := k.privateKeyMap[fromAddress]; exists {
		return types.SignTx(tx, types.LatestSignerForChainID(k.chainID), key)
	}
	return nil, fmt.Errorf("private key for address: %v not found", fromAddress)
}

func (k *DummyKeystore) SignMessage(ctx context.Context, address common.Address, data []byte) ([]byte, error) {
	key, exists := k.privateKeyMap[address]
	if !exists {
		return nil, fmt.Errorf("private key for address: %v not found", address)
	}
	signature, err := crypto.Sign(accounts.TextHash(data), key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message for address: %v", address)
	}
	return signature, nil
}

func (k *DummyKeystore) EnabledAddresses(_ context.Context) (addresses []common.Address, err error) {
	for address := range k.privateKeyMap {
		addresses = append(addresses, address)
	}
	return
}
