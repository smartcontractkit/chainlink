package evm

import (
	evmdb "github.com/smartcontractkit/chainlink/v2/core/chains/evm/db"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

type LegacyContainer interface {
	DB() *evmdb.ScopedDB
	TxmORM() txmgr.EvmTxStore
	Chains() LegacyChainContainer
}

type LegacyContainerImpl struct {
	chains *LegacyChains
	db     *evmdb.ScopedDB
	txmorm txmgr.EvmTxStore
}

func (l *LegacyContainerImpl) DB() *evmdb.ScopedDB {
	return l.db
}

func (l *LegacyContainerImpl) TxmORM() txmgr.EvmTxStore {
	return l.txmorm
}

func (l *LegacyContainerImpl) Chains() LegacyChainContainer {
	return l.chains
}
func NewLegacyContainerImpl(db *evmdb.ScopedDB, txm txmgr.EvmTxStore, chains *LegacyChains) *LegacyContainerImpl {
	return &LegacyContainerImpl{
		db:     db,
		txmorm: txm,
		chains: chains,
	}
}
