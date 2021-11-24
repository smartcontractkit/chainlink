package offchainreporting2_test

import (
	"crypto/rand"
	"testing"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"
)

func MakeConfigDigest(t *testing.T) ocrtypes.ConfigDigest {
	t.Helper()
	b := make([]byte, 32)
	/* #nosec G404 */
	_, err := rand.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	return MustBytesToConfigDigest(t, b)
}

func MustBytesToConfigDigest(t *testing.T, b []byte) ocrtypes.ConfigDigest {
	t.Helper()
	configDigest, err := ocrtypes.BytesToConfigDigest(b)
	if err != nil {
		t.Fatal(err)
	}
	return configDigest
}

func MustInsertOffchainreportingOracleSpec(t *testing.T, db *sqlx.DB, transmitterAddress ethkey.EIP55Address) job.OffchainReporting2OracleSpec {
	t.Helper()

	//pid := p2pkey.PeerID(cltest.DefaultP2PPeerID)
	//spec := job.OffchainReporting2OracleSpec{
	//	ContractAddress:                        cltest.NewEIP55Address(),
	//	P2PPeerID:                              &pid,
	//	P2PBootstrapPeers:                      pq.StringArray{},
	//	IsBootstrapPeer:                        false,
	//	EncryptedOCRKeyBundleID:                null.NewString(cltest.DefaultOCR2KeyBundleID, true),
	//	TransmitterAddress:                     &transmitterAddress,
	//	BlockchainTimeout:                      0,
	//	ContractConfigTrackerSubscribeInterval: 0,
	//	ContractConfigTrackerPollInterval:      0,
	//	ContractConfigConfirmations:            0,
	//}
	spec := job.OffchainReporting2OracleSpec{}
	mockJuelsPerFeeCoinSource := `ds1          [type=bridge name=voter_turnout];
	ds1_parse    [type=jsonparse path="one,two"];
	ds1_multiply [type=multiply times=1.23];
	ds1 -> ds1_parse -> ds1_multiply -> answer1;
	answer1      [type=median index=0];`
	require.NoError(t, db.Get(&spec, `INSERT INTO offchainreporting2_oracle_specs (created_at, updated_at, contract_address, p2p_bootstrap_peers, is_bootstrap_peer, encrypted_ocr_key_bundle_id, monitoring_endpoint, transmitter_address, blockchain_timeout, contract_config_tracker_subscribe_interval, contract_config_tracker_poll_interval, contract_config_confirmations, juels_per_fee_coin_pipeline) VALUES (
NOW(),NOW(),$1,'{}',false,$2,$3,$4,0,0,0,0,$5
) RETURNING *`, cltest.NewEIP55Address(), cltest.DefaultOCR2KeyBundleID, "chain.link:1234", &transmitterAddress, mockJuelsPerFeeCoinSource))
	//require.NoError(t, db.Create(&spec).Error)
	return spec
}
