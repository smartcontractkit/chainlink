package persistence

type dkgShare struct {
	ConfigDigest         []byte `db:"config_digest"`
	KeyID                []byte `db:"key_id"`
	Dealer               []byte `db:"dealer"`
	MarshaledShareRecord []byte `db:"marshaled_share_record"`
	RecordHash           []byte `db:"record_hash"`
}
