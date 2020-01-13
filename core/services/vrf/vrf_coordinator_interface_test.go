package vrf_test

import (
	"math/big"
	"testing"

	"chainlink/core/assets"
	"chainlink/core/services/vrf"
	"chainlink/core/store/models/vrfkey"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

// Taken from VRFCoordinator_test.js
var (
	secretKey, _ = vrfkey.NewPrivateKeyXXXTestingOnly(big.NewInt(1))
	keyHash      = secretKey.PublicKey.Hash()
	jobID        = common.BytesToHash([]byte("1234567890abcdef1234567890abcdef"))
	seed         = big.NewInt(1)
	sender       = common.HexToAddress("0xecfcab0a285d3380e488a39b4bb21e777f8a4eac")
	fee          = assets.NewLink(100)
)

// Taken from VRFCoordinator_test.js output: console.log(rREvents[0].data)
var solidityLogData = "0x" +
	"c0a6c424ac7157ae408398df7e5f4552091a69125d5dfcb7b8c2659029395bdf" + // keyHash
	"0000000000000000000000000000000000000000000000000000000000000001" + // seed
	"3132333435363738393061626364656631323334353637383930616263646566" + // jobID
	"000000000000000000000000ecfcab0a285d3380e488a39b4bb21e777f8a4eac" + // sender
	"0000000000000000000000000000000000000000000000000000000000000064" // fee

func TestVRFParseRandomnessRequestLog(t *testing.T) {
	r := vrf.RandomnessRequestLog{keyHash, seed, jobID, sender, fee}
	rawLog, err := r.RawLog()
	require.NoError(t, err)
	require.Equal(t, hexutil.Encode(rawLog), solidityLogData)
	nR, err := vrf.ParseRandomnessRequestLog(rawLog)
	require.NoError(t, err)
	require.True(t, r.Equal(*nR))
}
