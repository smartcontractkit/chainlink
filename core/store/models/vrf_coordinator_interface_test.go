package models_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	secretKey = vrfkey.NewPrivateKeyXXXTestingOnly(big.NewInt(1))
	keyHash   = secretKey.PublicKey.MustHash()
	jobID     = common.BytesToHash([]byte("1234567890abcdef1234567890abcdef"))
	seed      = big.NewInt(1)
	sender    = common.HexToAddress("0xecfcab0a285d3380e488a39b4bb21e777f8a4eac")
	fee       = assets.NewLink(100)
	raw       = vrf.RawRandomnessRequestLog{keyHash, seed, jobID, sender,
		(*big.Int)(fee), types.Log{
			// A raw, on-the-wire RandomnessRequestLog is the concat of fields as uint256's
			Data: append(append(append(
				keyHash.Bytes(),
				common.BigToHash(seed).Bytes()...),
				sender.Hash().Bytes()...),
				fee.ToHash().Bytes()...),
			Topics: []common.Hash{common.Hash{}, jobID},
		},
	}
)

func TestVRFParseRandomnessRequestLog(t *testing.T) {
	r := vrf.RawRandomnessRequestLogToRandomnessRequestLog(&raw)
	rawLog, err := r.RawData()
	require.NoError(t, err)
	assert.Equal(t, rawLog, raw.Raw.Data)
	nR, err := vrf.ParseRandomnessRequestLog(models.Log{
		Data:   rawLog,
		Topics: []common.Hash{common.Hash{}, jobID},
	})
	require.NoError(t, err)
	require.True(t, r.Equal(*nR),
		"Round-tripping RandomnessRequestLog through serialization and parsing "+
			"resulted in a different log.")
}
