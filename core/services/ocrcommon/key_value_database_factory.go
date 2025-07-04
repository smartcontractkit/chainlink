package ocrcommon

import (
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type KeyValueDatabaseFactory struct {
	ds sqlutil.DataSource
}

func NewKeyValueDatabaseFactory(ds sqlutil.DataSource) *KeyValueDatabaseFactory {
	return &KeyValueDatabaseFactory{
		ds: ds,
	}
}

func (k *KeyValueDatabaseFactory) NewKeyValueDatabase(configDigest types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	return &KeyValueDatabase{
		configDigest: configDigest,
		ds:           k.ds,
	}, nil
}

type KeyValueDatabase struct {
	configDigest types.ConfigDigest
	ds           sqlutil.DataSource
}

func (k *KeyValueDatabase) NewReadWriteTransaction() (ocr3_1types.KeyValueReadWriteTransaction, error) {
	tx := &KeyValueReadWriteTransaction{}
	return tx.Start()
}

func (k *KeyValueDatabase) NewReadTransaction() (ocr3_1types.KeyValueReadTransaction, error) {
	return nil, nil
}

func (k *KeyValueDatabase) Close() error {
	return nil
}

type KeyValueReadTransaction struct {
}

func (k *KeyValueReadTransaction) Read(key []byte) ([]byte, error) {
	return nil, nil
}

func (k *KeyValueReadTransaction) Range(loKey []byte, hiKeyExcl []byte) ocr3_1types.KeyValueIterator {
	return nil
}

func (k *KeyValueReadTransaction) Discard() {}

type KeyValueReadWriteTransaction struct {
	KeyValueReadTransaction
}

func (k *KeyValueReadWriteTransaction) Start() error) {

}

func (k *KeyValueReadWriteTransaction) Write(key []byte, value []byte) error {
	return nil
}

func (k *KeyValueReadWriteTransaction) Delete(key []byte) error {
	return nil
}

func (k *KeyValueReadWriteTransaction) Commit() error {
	return nil
}
