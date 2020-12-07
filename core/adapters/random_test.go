package adapters_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	tvrf "github.com/smartcontractkit/chainlink/core/internal/cltest/vrf"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NB: For changes to the VRF solidity code to be reflected here, "go generate"
// must be run in core/services/vrf.
func TestRandom_Perform(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	publicKey := cltest.StoredVRFKey(t, store)
	adapter := adapters.Random{PublicKey: publicKey.String()}
	hash := utils.MustHash("a random string")
	seed := big.NewInt(0x10)
	blockNum := 10
	jsonInput, err := models.JSON{}.MultiAdd(models.KV{
		"seed":      utils.Uint64ToHex(seed.Uint64()),
		"keyHash":   publicKey.MustHash().Hex(),
		"blockHash": hash.Hex(),
		"blockNum":  blockNum,
	})
	require.NoError(t, err) // Can't fail
	input := models.NewRunInput(&models.ID{}, models.ID{}, jsonInput,
		models.RunStatusUnstarted)
	result := adapter.Perform(*input, store)
	require.NoError(t, result.Error(), "while running random adapter")
	proofArg := hexutil.MustDecode(result.Result().String())
	var wireProof []byte
	out, err := models.VRFFulfillMethod().Inputs.Unpack(proofArg)
	wireProof = abi.ConvertType(out[0], []byte{}).([]byte)
	require.NoError(t, err, "failed to unpack VRF proof from random adapter")
	var onChainResponse vrf.MarshaledOnChainResponse
	require.Equal(t, copy(onChainResponse[:], wireProof),
		vrf.OnChainResponseLength, "wrong response length")
	response, err := vrf.UnmarshalProofResponse(onChainResponse)
	require.NoError(t, err, "random adapter produced bad proof response")
	actualProof, err := response.CryptoProof(tvrf.SeedData(t, seed, hash, blockNum))
	require.NoError(t, err, "could not extract proof from random adapter response")
	expected := common.HexToHash(
		"0x71a7c50918feaa753485ae039cb84ddd70c5c85f66b236138dea453a23d0f27e")
	assert.Equal(t, expected, common.BigToHash(actualProof.Output),
		"unexpected VRF output; perhas vrfkey.json or the output hashing function "+
			"in RandomValueFromVRFProof has changed?")
	jsonInput, err = jsonInput.Add("keyHash", common.Hash{})
	require.NoError(t, err)
	input = models.NewRunInput(&models.ID{}, models.ID{}, jsonInput, models.RunStatusUnstarted)
	result = adapter.Perform(*input, store)
	require.Error(t, result.Error(), "must reject if keyHash doesn't match")
}
