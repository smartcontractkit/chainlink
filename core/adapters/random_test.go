package adapters_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_verifier_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var suite = secp256k1.NewBlakeKeccackSecp256k1()

// NB: For changes to the VRF solidity code to be reflected here, "go generate"
// must be run in core/services/vrf.
func vrfVerifier(t *testing.T) *solidity_verifier_wrapper.VRFTestHelper {
	ethereumKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth := bind.NewKeyedTransactor(ethereumKey)
	genesisData := core.GenesisAlloc{auth.From: {Balance: big.NewInt(1000000000)}}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	_, _, verifier, err := solidity_verifier_wrapper.DeployVRFTestHelper(auth, backend)
	require.NoError(t, err)
	backend.Commit()
	return verifier
}

func TestRandom_Perform(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	publicKey := cltest.StoredVRFKey(t, store)
	adapter := adapters.Random{PublicKey: publicKey.String()}
	jsonInput, err := models.JSON{}.Add("seed", "0x10")
	require.NoError(t, err) // Can't fail
	jsonInput, err = jsonInput.Add("keyHash", publicKey.MustHash().Hex())
	require.NoError(t, err) // Can't fail
	input := models.NewRunInput(&models.ID{}, models.ID{}, jsonInput, models.RunStatusUnstarted)
	result := adapter.Perform(*input, store)
	require.NoError(t, result.Error(), "while running random adapter")
	proof := hexutil.MustDecode(result.Result().String())
	// Check that proof is a solidity bytes array containing the actual proof
	length := big.NewInt(0).SetBytes(proof[:utils.EVMWordByteLen]).Uint64()
	require.Equal(t, length, uint64(len(proof)-utils.EVMWordByteLen))
	actualProof := proof[utils.EVMWordByteLen:]
	randomOutput, err := vrfVerifier(t).RandomValueFromVRFProof(nil, actualProof)
	require.NoError(t, err, "proof was invalid")
	expected := common.HexToHash(
		"c0a5642a409290ac65d9d44a4c52e53f31921ff1b7d235c585193a18190c82f1")
	assert.Equal(t, expected, common.BigToHash(randomOutput),
		"unexpected VRF output; perhas vrfkey.json or the output hashing function "+
			"in RandomValueFromVRFProof has changed?")
	jsonInput, err = jsonInput.Add("keyHash", common.Hash{})
	require.NoError(t, err)
	input = models.NewRunInput(&models.ID{}, models.ID{}, jsonInput, models.RunStatusUnstarted)
	result = adapter.Perform(*input, store)
	require.Error(t, result.Error(), "must reject if keyHash doesn't match")
}
