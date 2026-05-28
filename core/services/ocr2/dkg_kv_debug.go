package ocr2

// Temporary instrumentation around the DKG OCR3 plugin's KeyValueDatabase to
// trace ReportsPlusPrecursor (RP|...) writes and reads. The Pebble factory
// returned by libocr exposes only the raw key/value interface to chainlink;
// libocr's semantic shim builds RP keys as:
//   "RP|" + big-endian uint64 seqNr (8 bytes) + 32-byte digest
// (see libocr offchainreporting2plus/internal/shim/ocr3_1_key_value_store.go).
//
// We wrap the factory and log every Write/Read/Delete whose key starts with
// "RP|", decoding the seqNr and digest. This is to debug a case where libocr's
// report attestation logs "did not find ReportsPlusPrecursor in kv" repeatedly
// while smdkg's DKG plugin reaches "gathered enough inner dealings" but never
// finalizes. See confidential-compute PR #347 for the full evidence chain.
//
// Strictly for debugging; remove after the issue is diagnosed.

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	libocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const rpKeyPrefix = "RP|"

// wrapDebugKVDBFactory wraps a libocr KeyValueDatabaseFactory so that
// ReportsPlusPrecursor key operations are logged. All other operations are
// forwarded unchanged.
func wrapDebugKVDBFactory(inner ocr3_1types.KeyValueDatabaseFactory, lggr logger.Logger) ocr3_1types.KeyValueDatabaseFactory {
	return &debugKVDBFactory{inner: inner, lggr: lggr}
}

type debugKVDBFactory struct {
	inner ocr3_1types.KeyValueDatabaseFactory
	lggr  logger.Logger
}

func (f *debugKVDBFactory) NewKeyValueDatabase(configDigest libocrtypes.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	db, err := f.inner.NewKeyValueDatabase(configDigest)
	if err != nil || db == nil {
		return db, err
	}
	logger.Sugared(f.lggr).Infow("[DKG-DEBUG] kv factory: NewKeyValueDatabase", "configDigest", hex.EncodeToString(configDigest[:]))
	return &debugKVDB{KeyValueDatabase: db, lggr: f.lggr}, nil
}

func (f *debugKVDBFactory) NewKeyValueDatabaseIfExists(configDigest libocrtypes.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	db, err := f.inner.NewKeyValueDatabaseIfExists(configDigest)
	if err != nil || db == nil {
		return db, err
	}
	logger.Sugared(f.lggr).Infow("[DKG-DEBUG] kv factory: NewKeyValueDatabaseIfExists", "configDigest", hex.EncodeToString(configDigest[:]))
	return &debugKVDB{KeyValueDatabase: db, lggr: f.lggr}, nil
}

type debugKVDB struct {
	ocr3_1types.KeyValueDatabase
	lggr logger.Logger
}

func (kv *debugKVDB) NewReadWriteTransaction() (ocr3_1types.KeyValueDatabaseReadWriteTransaction, error) {
	tx, err := kv.KeyValueDatabase.NewReadWriteTransaction()
	if err != nil {
		return nil, err
	}
	return &debugKVDBRWTx{KeyValueDatabaseReadWriteTransaction: tx, lggr: kv.lggr}, nil
}

func (kv *debugKVDB) NewReadTransaction() (ocr3_1types.KeyValueDatabaseReadTransaction, error) {
	tx, err := kv.KeyValueDatabase.NewReadTransaction()
	if err != nil {
		return nil, err
	}
	return &debugKVDBROTx{KeyValueDatabaseReadTransaction: tx, lggr: kv.lggr}, nil
}

type debugKVDBRWTx struct {
	ocr3_1types.KeyValueDatabaseReadWriteTransaction
	lggr logger.Logger
}

func (tx *debugKVDBRWTx) Write(key, value []byte) error {
	if isRPKey(key) {
		seqNr, digest := parseRPKey(key)
		logger.Sugared(tx.lggr).Infow("[DKG-DEBUG] kv Write ReportsPlusPrecursor",
			"seqNr", seqNr, "digest", digest, "valueLen", len(value))
	}
	return tx.KeyValueDatabaseReadWriteTransaction.Write(key, value)
}

func (tx *debugKVDBRWTx) Read(key []byte) ([]byte, error) {
	val, err := tx.KeyValueDatabaseReadWriteTransaction.Read(key)
	if isRPKey(key) {
		seqNr, digest := parseRPKey(key)
		logger.Sugared(tx.lggr).Infow("[DKG-DEBUG] kv Read ReportsPlusPrecursor (rw)",
			"seqNr", seqNr, "digest", digest, "found", val != nil, "err", err)
	}
	return val, err
}

func (tx *debugKVDBRWTx) Delete(key []byte) error {
	if isRPKey(key) {
		seqNr, digest := parseRPKey(key)
		logger.Sugared(tx.lggr).Infow("[DKG-DEBUG] kv Delete ReportsPlusPrecursor",
			"seqNr", seqNr, "digest", digest)
	}
	return tx.KeyValueDatabaseReadWriteTransaction.Delete(key)
}

type debugKVDBROTx struct {
	ocr3_1types.KeyValueDatabaseReadTransaction
	lggr logger.Logger
}

func (tx *debugKVDBROTx) Read(key []byte) ([]byte, error) {
	val, err := tx.KeyValueDatabaseReadTransaction.Read(key)
	if isRPKey(key) {
		seqNr, digest := parseRPKey(key)
		logger.Sugared(tx.lggr).Infow("[DKG-DEBUG] kv Read ReportsPlusPrecursor (ro)",
			"seqNr", seqNr, "digest", digest, "found", val != nil, "err", err)
	}
	return val, err
}

func isRPKey(key []byte) bool {
	return len(key) >= 3 && string(key[:3]) == rpKeyPrefix
}

// parseRPKey decodes a libocr ReportsPlusPrecursor raw key into (seqNr, digestHex).
// Returns (0, "(malformed)") if the key is not the expected length.
func parseRPKey(key []byte) (uint64, string) {
	const headerLen = 3 + 8 + 32 // "RP|" + uint64 + digest
	if len(key) < headerLen {
		return 0, "(malformed)"
	}
	seqNr := binary.BigEndian.Uint64(key[3:11])
	digest := hex.EncodeToString(key[11:43])
	return seqNr, digest
}
