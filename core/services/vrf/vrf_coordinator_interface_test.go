package vrf_test

import (
	"math/big"
	"testing"

	"chainlink/core/assets"
	"chainlink/core/eth"
	"chainlink/core/services/vrf"
	"chainlink/core/store/models/vrfkey"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var solidityLogData = "0x" + // Example of a raw, on-the-wire RandomnessRequestLog
	"326262e3f0eb0e651194df57377a65cb4b9b30d6ef6601ce5f5760121469c95c" + // keyHash
	"0000000000000000000000000000000000000000000000000000000000000001" + // seed
	"000000000000000000000000ecfcab0a285d3380e488a39b4bb21e777f8a4eac" + // sender
	"0000000000000000000000000000000000000000000000000000000000000064" // fee

// Taken from VRFCoordinator_test.js
var (
	secretKey = vrfkey.NewPrivateKeyXXXTestingOnly(big.NewInt(1))
	keyHash   = secretKey.PublicKey.Hash()
	jobID     = common.BytesToHash([]byte("1234567890abcdef1234567890abcdef"))
	seed      = big.NewInt(1)
	sender    = common.HexToAddress("0xecfcab0a285d3380e488a39b4bb21e777f8a4eac")
	fee       = assets.NewLink(100)
	raw       = vrf.RawRandomnessRequestLog{keyHash, seed, jobID, sender,
		(*big.Int)(fee), types.Log{
			Data:   common.Hex2Bytes(solidityLogData[2:]),
			Topics: []common.Hash{common.Hash{}, jobID},
		},
	}
)

func TestVRFParseRandomnessRequestLog(t *testing.T) {
	r := vrf.RandomnessRequestLog{keyHash, seed, jobID, sender, fee, raw}
	rawLog, err := r.RawData()
	require.NoError(t, err)
	assert.Equal(t, hexutil.Encode(rawLog), solidityLogData)
	nR, err := vrf.ParseRandomnessRequestLog(eth.Log{
		Data:   rawLog,
		Topics: []common.Hash{common.Hash{}, jobID},
	})
	require.NoError(t, err)
	require.True(t, r.Equal(*nR))
}
