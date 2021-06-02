package adapters_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NB: For changes to the VRF solidity code to be reflected here, "go generate"
// must be run in core/services/vrf.
func TestRandom_Perform(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)
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
	jr := cltest.NewJobRun(cltest.NewJobWithRandomnessLog())
	input := models.NewRunInput(jr, uuid.Nil, jsonInput, models.RunStatusUnstarted)
	result := adapter.Perform(*input, store, keyStore)
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
	actualProof, err := response.CryptoProof(vrf.TestXXXSeedData(t, seed, hash, blockNum))
	require.NoError(t, err, "could not extract proof from random adapter response")
	expected := common.HexToHash(
		"0x71a7c50918feaa753485ae039cb84ddd70c5c85f66b236138dea453a23d0f27e")
	assert.Equal(t, expected, common.BigToHash(actualProof.Output),
		"unexpected VRF output; perhas vrfkey.json or the output hashing function "+
			"in RandomValueFromVRFProof has changed?")
	jsonInput, err = jsonInput.Add("keyHash", common.Hash{})
	require.NoError(t, err)
	input = models.NewRunInput(jr, uuid.Nil, jsonInput, models.RunStatusUnstarted)
	result = adapter.Perform(*input, store, keyStore)
	require.Error(t, result.Error(), "must reject if keyHash doesn't match")
}

func TestRandom_Perform_CheckFulfillment(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	ethMock := new(mocks.Client)
	store.EthClient = ethMock

	publicKey := cltest.StoredVRFKey(t, store)
	address := cltest.NewEIP55Address()
	hash := utils.MustHash("a random string")
	seed := big.NewInt(0x10)
	blockNum := 10
	jsonInput, err := models.JSON{}.MultiAdd(models.KV{
		"seed":      utils.Uint64ToHex(seed.Uint64()),
		"keyHash":   publicKey.MustHash().Hex(),
		"blockHash": hash.Hex(),
		"blockNum":  blockNum,
		"requestID": utils.AddHexPrefix(common.Bytes2Hex([]byte{1, 2, 3})),
	})
	require.NoError(t, err)
	jr := cltest.NewJobRun(cltest.NewJobWithRandomnessLog())
	input := models.NewRunInput(jr, uuid.Nil, jsonInput, models.RunStatusUnstarted)

	abi := eth.MustGetABI(solidity_vrf_coordinator_interface.VRFCoordinatorABI)
	registryMock := cltest.NewContractMockReceiver(t, ethMock, abi, address.Address())

	for _, test := range []struct {
		name                   string
		addressParamPresent    bool
		seedAndBlockNumPresent bool
		shouldFulfill          bool
	}{
		{"both missing", false, false, true},
		{"address missing, seed/block present", false, true, true},
		{"address present, seed/block missing", true, false, false},
		{"both present", true, true, true},
	} {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			adapter := adapters.Random{PublicKey: publicKey.String()}
			response := solidity_vrf_coordinator_interface.Callbacks{
				CallbackContract: cltest.NewAddress(),
				RandomnessFee:    big.NewInt(100),
			}

			if test.seedAndBlockNumPresent {
				response.SeedAndBlockNum = [32]byte{1, 2, 3}
			}
			if test.addressParamPresent {
				adapter.CoordinatorAddress = address
				registryMock.MockResponse("callbacks", response).Once()
			}

			result := adapter.Perform(*input, store, keyStore)
			require.Equal(tt, test.shouldFulfill, result.Error() == nil)
			ethMock.AssertExpectations(t)
		})
	}
}
