package adapters_test

import (
	"fmt"
	"math/big"
	"testing"

	"chainlink/core/adapters"
	"chainlink/core/internal/cltest"
	"chainlink/core/services/signatures/secp256k1"
	"chainlink/core/services/vrf/generated/solidity_verifier_wrapper"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var suite = secp256k1.NewBlakeKeccackSecp256k1()

// NB: For changes to the VRF solidity code to be reflected here, "go generate"
// must be run in core/services/vrf.
func vRFVerifier() *solidity_verifier_wrapper.VRFTestHelper {
	ethereumKey, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(ethereumKey)
	genesisData := core.GenesisAlloc{auth.From: {Balance: big.NewInt(1000000000)}}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	_, _, verifier, err := solidity_verifier_wrapper.DeployVRFTestHelper(auth, backend)
	if err != nil {
		panic(errors.Wrapf(err, "while initializing EVM contract wrapper"))
	}
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
	jsonInput, err = jsonInput.Add("keyHash", publicKey.Hash().Hex())
	require.NoError(t, err) // Can't fail
	input := models.NewRunInput(&models.ID{}, jsonInput, models.RunStatusUnstarted)
	result := adapter.Perform(*input, store)
	require.NoError(t, result.Error(), "while running random adapter")
	proof := hexutil.MustDecode(result.Result().String())
	// Check that proof is a solidity bytes array containing the actual proof
	length := big.NewInt(0).SetBytes(proof[:utils.EVMWordByteLen]).Uint64()
	require.Equal(t, length, uint64(len(proof)-utils.EVMWordByteLen))
	actualProof := proof[utils.EVMWordByteLen:]
	randomOutput, err := vRFVerifier().RandomValueFromVRFProof(nil, actualProof)
	require.NoError(t, err, "proof was invalid")
	expected, ok := big.NewInt(0).SetString(
		// Depends on vrfkey.json, and will need to be changed if that changes.
		"b2002342b67f1e9c27e7abe157cfc8fc1912a7c9bde570aeed5d13c8d00c497f", 16)
	assert.True(t, randomOutput.Cmp(expected) == 0)
	require.True(t, ok)
	jsonInput, err = jsonInput.Add("keyHash", common.Hash{})
	require.NoError(t, err)
	input = models.NewRunInput(&models.ID{}, jsonInput, models.RunStatusUnstarted)
	result = adapter.Perform(*input, store)
	fmt.Println("result", result.Result().String())
	require.Error(t, result.Error(), "must reject if keyHash doesn't match")
}
