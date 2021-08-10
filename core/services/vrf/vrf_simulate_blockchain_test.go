package vrf_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	proof2 "github.com/smartcontractkit/chainlink/core/services/vrf/proof"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func registerExistingProvingKey(
	t *testing.T,
	coordinator coordinatorUniverse,
	provingKey *vrfkey.PrivateKey,
	jobID models.JobID,
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
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)

	key := cltest.MustGenerateRandomKey(t)

	cu := newVRFCoordinatorUniverse(t, key)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, cu.backend, key)
	defer cleanup()

	app.Start()

	rawKey := "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
	pk, err := secp256k1.NewPublicKeyFromHex(rawKey)
	require.NoError(t, err)
	var sk int64 = 1
	provingKey := vrfkey.NewPrivateKeyXXXTestingOnly(big.NewInt(sk))
	require.Equal(t, provingKey.PublicKey, pk,
		"public key in fixture %s does not match secret key in test %d (which has "+
			"public key %s)", pk, sk, provingKey.PublicKey.String())
	app.KeyStore.VRF().StoreInMemoryXXXTestingOnly(provingKey)
	var seed = big.NewInt(1)

	j := cltest.NewJobWithRandomnessLog()
	contractAddress := cu.rootContractAddress.String()
	task1Params := cltest.JSONFromString(t, fmt.Sprintf(`{"publicKey": "%s"}`, rawKey))
	task2JSON := fmt.Sprintf(`{"format": "preformatted", "address": "%s", "functionSelector": "0x5e1c1059"}`, contractAddress)
	task2Params := cltest.JSONFromString(t, task2JSON)

	j.Initiators[0].Address = cu.rootContractAddress
	j.Tasks = []models.TaskSpec{{
		Type:                             adapters.TaskTypeRandom,
		Params:                           task1Params,
		MinRequiredIncomingConfirmations: null.NewUint32(1, true),
	}, {
		Type:   adapters.TaskTypeEthTx,
		Params: task2Params,
	}}

	j = cltest.CreateJobSpecViaWeb(t, app, j)
	registerExistingProvingKey(t, cu, provingKey, j.ID, vrfFee)
	r := requestRandomness(t, cu, provingKey.PublicKey.MustHash(), big.NewInt(100))

	cltest.WaitForRuns(t, j, app.Store, 1)
	runs, err := app.Store.JobRunsFor(j.ID)
	assert.NoError(t, err)
	require.Len(t, runs, 1)
	jr := runs[0]
	require.Len(t, jr.TaskRuns, 2)
	assert.False(t, jr.TaskRuns[0].ObservedIncomingConfirmations.Valid)

	cltest.WaitForCount(t, app.Store, bulletprooftxmanager.EthTx{}, 1)
	app.TxManager.Trigger(app.Key.Address.Address())
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
	require.Len(t, proof, proof2.OnChainResponseLength)
	publicPoint, err := provingKey.PublicKey.Point()
	require.NoError(t, err)
	require.Equal(t, proof[:64], secp256k1.LongMarshal(publicPoint))
	mProof := proof2.MarshaledOnChainResponse{}
	require.Equal(t, copy(mProof[:], proof), proof2.OnChainResponseLength)
	goProof, err := proof2.UnmarshalProofResponse(mProof)
	require.NoError(t, err, "problem parsing solidity proof")
	preSeed, err := proof2.BigToSeed(r.Seed)
	require.NoError(t, err, "seed %x out of range", seed)
	_, err = goProof.CryptoProof(proof2.PreSeedData{
		PreSeed:   preSeed,
		BlockHash: r.Raw.Raw.BlockHash,
		BlockNum:  uint64(r.Raw.Raw.BlockNumber),
	})
	require.NoError(t, err, "problem verifying solidity proof")
}

// TestIntegration_SharedProvingKey tests the scenario where multiple nodes share
// a single proving key
func TestIntegration_SharedProvingKey(t *testing.T) {
	config, _, cfgCleanup := heavyweight.FullTestORM(t, "vrf_shared_proving_key", true, true)
	defer cfgCleanup()
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)
	config.Config.Dialect = dialects.PostgresWithoutLock

	key := cltest.MustGenerateRandomKey(t)

	cu := newVRFCoordinatorUniverse(t, key)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, cu.backend, key)
	defer cleanup()

	app.Start()

	// create job
	rawKey := "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
	pk, err := secp256k1.NewPublicKeyFromHex(rawKey)
	require.NoError(t, err)
	var sk int64 = 1
	provingKey := vrfkey.NewPrivateKeyXXXTestingOnly(big.NewInt(sk))
	require.Equal(t, provingKey.PublicKey, pk,
		"public key in fixture %s does not match secret key in test %d (which has "+
			"public key %s)", pk, sk, provingKey.PublicKey.String())
	app.KeyStore.VRF().StoreInMemoryXXXTestingOnly(provingKey)

	j := cltest.NewJobWithRandomnessLog()
	contractAddress := cu.rootContractAddress.String()
	task1Params := cltest.JSONFromString(t, fmt.Sprintf(`{"publicKey": "%s", "coordinatorAddress": "%s"}`, rawKey, contractAddress))
	task2JSON := fmt.Sprintf(`{"format": "preformatted", "address": "%s", "functionSelector": "0x5e1c1059"}`, contractAddress)
	task2Params := cltest.JSONFromString(t, task2JSON)

	j.Initiators[0].Address = cu.rootContractAddress
	j.Tasks = []models.TaskSpec{{
		Type:                             adapters.TaskTypeRandom,
		Params:                           task1Params,
		MinRequiredIncomingConfirmations: null.NewUint32(3, true), // allow space for another node to answer before us
	}, {
		Type:   adapters.TaskTypeEthTx,
		Params: task2Params,
	}}

	j = cltest.CreateJobSpecViaWeb(t, app, j)
	registerExistingProvingKey(t, cu, provingKey, j.ID, vrfFee)

	// trigger job run by requesting randomness
	log := requestRandomness(t, cu, provingKey.PublicKey.MustHash(), big.NewInt(100))
	seed := common.BigToHash(log.Seed).String()
	cltest.WaitForRuns(t, j, app.Store, 1)
	var jobRun models.JobRun
	err = app.Store.DB.First(&jobRun).Error
	require.NoError(t, err)
	cltest.WaitForJobRunStatus(t, app.Store, jobRun, models.RunStatusPendingIncomingConfirmations)

	// simulate fulfillment from other node - use the Perform adapter to "steal" the proof
	// from neil, but then submit the proof from ned
	jsonInput, err := models.JSON{}.MultiAdd(models.KV{
		"seed":      seed,
		"keyHash":   pk.MustHash().Hex(),
		"blockHash": log.Raw.Raw.BlockHash.Hex(),
		"blockNum":  log.Raw.Raw.BlockNumber,
	})
	require.NoError(t, err)
	jr := cltest.NewJobRun(cltest.NewJobWithRandomnessLog())
	input := models.NewRunInput(jr, uuid.Nil, jsonInput, models.RunStatusUnstarted)
	adapter := adapters.Random{PublicKey: pk.String()}
	result := adapter.Perform(*input, app.Store, app.KeyStore)
	require.NoError(t, result.Error(), "while running random adapter")
	encodedProofHex := result.Result().String()
	encodedProof, err := hexutil.Decode(encodedProofHex)
	require.NoError(t, err)
	inputs, err := models.VRFFulfillMethod().Inputs.UnpackValues(encodedProof)
	require.NoError(t, err)
	proof, ok := inputs[0].([]byte)
	require.True(t, ok)

	_, err = cu.rootContract.FulfillRandomnessRequest(cu.ned, proof)
	require.NoError(t, err)
	cu.backend.Commit()

	// start mining, assert neil never attempts to respond
	stopMining := cltest.Mine(cu.backend, 250*time.Millisecond)
	defer stopMining()

	cltest.WaitForJobRunStatus(t, app.Store, jobRun, models.RunStatusErrored)
	cltest.AssertCount(t, app.Store.DB, bulletprooftxmanager.EthTxAttempt{}, 0)
}
