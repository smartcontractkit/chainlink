package txm

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type DummyKeystore struct {
	privateKey *ecdsa.PrivateKey
}

func NewKeystore() *DummyKeystore {
	return &DummyKeystore{}
}

func (k *DummyKeystore) Add(privateKeyString string) error {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return err
	}
	k.privateKey = privateKey
	return nil
}

func (k *DummyKeystore) SignTx(_ context.Context, fromAddress common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return types.SignTx(tx, types.LatestSignerForChainID(chainID), k.privateKey)
}
