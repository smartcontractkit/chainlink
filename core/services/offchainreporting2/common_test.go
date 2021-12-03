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

	spec := job.OffchainReporting2OracleSpec{}
	mockJuelsPerFeeCoinSource := `ds1          [type=bridge name=voter_turnout];
	ds1_parse    [type=jsonparse path="one,two"];
	ds1_multiply [type=multiply times=1.23];
	ds1 -> ds1_parse -> ds1_multiply -> answer1;
	answer1      [type=median index=0];`
	require.NoError(t, db.Get(&spec, `INSERT INTO offchainreporting2_oracle_specs (created_at, updated_at, relay, relay_config, contract_id, p2p_bootstrap_peers, is_bootstrap_peer, ocr_key_bundle_id, monitoring_endpoint, transmitter_id, blockchain_timeout, contract_config_tracker_subscribe_interval, contract_config_tracker_poll_interval, contract_config_confirmations, juels_per_fee_coin_pipeline) VALUES (
NOW(),NOW(), 'ethereum', '{}', $1,'{}',false,$2,$3,$4,0,0,0,0,$5
) RETURNING *`, cltest.NewEIP55Address().String(), cltest.DefaultOCR2KeyBundleID, "chain.link:1234", transmitterAddress.String(), mockJuelsPerFeeCoinSource))
	return spec
}
