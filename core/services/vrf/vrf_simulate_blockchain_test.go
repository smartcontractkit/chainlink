package vrf_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func registerExistingProvingKey(
	t *testing.T,
	coordinator coordinatorUniverse,
	provingKey *vrfkey.PrivateKey,
	jobID *models.ID,
	vrfFee *big.Int,
) {
	var rawJobID [32]byte
	copy(rawJobID[:], []byte(jobID.String()))
	_, err := coordinator.rootContract.RegisterProvingKey(
		coordinator.neil, vrfFee, coordinator.neil.From, pair(secp256k1.Coordinates(publicKey)), rawJobID)
	require.NoError(t, err, "failed to register VRF proving key on VRFCoordinator contract")
	coordinator.backend.Commit()
}

func TestIntegration_RandomnessRequest(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	cu := newVRFCoordinatorUniverse(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, cu.backend)

	app.Start()

	rawKey := "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
	pk, err := vrfkey.NewPublicKeyFromHex(rawKey)
	require.NoError(t, err)
	var sk int64 = 1
	provingKey := vrfkey.NewPrivateKeyXXXTestingOnly(big.NewInt(sk))
	require.Equal(t, provingKey.PublicKey, pk,
		"public key in fixture %s does not match secret key in test %d (which has "+
			"public key %s)", pk, sk, provingKey.PublicKey.String())
	app.Store.VRFKeyStore.StoreInMemoryXXXTestingOnly(provingKey)

	j := cltest.NewJobWithRandomnessLog()
	task1Params := cltest.JSONFromString(t, fmt.Sprintf(`{"PublicKey": "%s"}`, rawKey))
	task2JSON := fmt.Sprintf(`{"format": "preformatted", "address": "%s", "functionSelector": "0x5e1c1059"}`, cu.rootContractAddress.String())
	task2Params := cltest.JSONFromString(t, task2JSON)

	j.Initiators[0].Address = cu.rootContractAddress
	j.Tasks = []models.TaskSpec{{
		Type:   adapters.TaskTypeRandom,
		Params: task1Params,
	}, {
		Type:   adapters.TaskTypeEthTx,
		Params: task2Params,
	}}

	j = cltest.CreateJobSpecViaWeb(t, app, j)
	registerExistingProvingKey(t, cu, provingKey, j.ID, vrfFee)
	r := requestRandomness(t, cu, provingKey.PublicKey.MustHash(), big.NewInt(100), seed)

	cltest.WaitForRuns(t, j, app.Store, 1)
	runs, err := app.Store.JobRunsFor(j.ID)
	assert.NoError(t, err)
	require.Len(t, runs, 1)
	jr := runs[0]
	require.Len(t, jr.TaskRuns, 2)
	assert.False(t, jr.TaskRuns[0].ObservedIncomingConfirmations.Valid)
	attempts := cltest.WaitForEthTxAttemptCount(t, app.Store, 1)
	require.Len(t, attempts, 1)

	rawTx := attempts[0].SignedRawTx
	var tx *types.Transaction
	require.NoError(t, rlp.DecodeBytes(rawTx, &tx))
	fixtureToAddress := j.Tasks[1].Params.Get("address").String()
	require.Equal(t, *tx.To(), common.HexToAddress(fixtureToAddress))
	payload := tx.Data()
	require.Equal(t, hexutil.Encode(payload[:4]), models.VRFFulfillSelector())
	proofContainer := make(map[string]interface{})
	err = models.VRFFulfillMethod().Inputs.UnpackIntoMap(proofContainer, payload[4:])
	require.NoError(t, err)
	proof, ok := proofContainer["_proof"].([]byte)
	require.True(t, ok)
	require.Len(t, proof, vrf.OnChainResponseLength)
	publicPoint, err := provingKey.PublicKey.Point()
	require.NoError(t, err)
	require.Equal(t, proof[:64], secp256k1.LongMarshal(publicPoint))
	mProof := vrf.MarshaledOnChainResponse{}
	require.Equal(t, copy(mProof[:], proof), vrf.OnChainResponseLength)
	goProof, err := vrf.UnmarshalProofResponse(mProof)
	require.NoError(t, err, "problem parsing solidity proof")
	preSeed, err := vrf.BigToSeed(r.Seed)
	require.NoError(t, err, "seed %x out of range", seed)
	_, err = goProof.CryptoProof(vrf.PreSeedData{
		PreSeed:   preSeed,
		BlockHash: r.Raw.Raw.BlockHash,
		BlockNum:  uint64(r.Raw.Raw.BlockNumber),
	})
	require.NoError(t, err, "problem verifying solidity proof")
}
