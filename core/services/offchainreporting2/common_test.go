package offchainreporting2_test

import (
	"crypto/rand"
	"testing"

	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
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

func MustInsertOffchainreportingOracleSpec(t *testing.T, db *gorm.DB, transmitterAddress ethkey.EIP55Address) job.OffchainReporting2OracleSpec {
	t.Helper()

	pid := p2pkey.PeerID(cltest.DefaultP2PPeerID)
	spec := job.OffchainReporting2OracleSpec{
		ContractAddress:                        cltest.NewEIP55Address(),
		P2PPeerID:                              &pid,
		P2PBootstrapPeers:                      pq.StringArray{},
		IsBootstrapPeer:                        false,
		EncryptedOCRKeyBundleID:                null.NewString(cltest.DefaultOCR2KeyBundleID, true),
		TransmitterAddress:                     &transmitterAddress,
		ObservationTimeout:                     0,
		BlockchainTimeout:                      0,
		ContractConfigTrackerSubscribeInterval: 0,
		ContractConfigTrackerPollInterval:      0,
		ContractConfigConfirmations:            0,
	}
	require.NoError(t, db.Create(&spec).Error)
	return spec
}
