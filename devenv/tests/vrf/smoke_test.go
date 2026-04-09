package vrf

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/solidity_vrf_consumer_interface"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/vrf"
)

func TestVRFBasic(t *testing.T) {
	ctx := t.Context()
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	productCfg, err := products.LoadOutput[vrf.Configurator](outputFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	c, auth, _, err := products.ETHClient(
		ctx,
		in.Blockchains[0].Out.Nodes[0].ExternalWSUrl,
		productCfg.Config[0].GasSettings.FeeCapMultiplier,
		productCfg.Config[0].GasSettings.TipCapMultiplier,
	)
	require.NoError(t, err)

	consumer, err := solidity_vrf_consumer_interface.NewVRFConsumer(
		common.HexToAddress(productCfg.Config[0].Out.ConsumerAddress), c,
	)
	require.NoError(t, err)

	keyHash := decodeKeyHash(t, productCfg.Config[0].Out.KeyHash)

	tx, err := consumer.TestRequestRandomness(auth, keyHash, big.NewInt(1))
	require.NoError(t, err)
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	require.NoError(t, err)

	cls, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		jobRuns, err := cls[0].MustReadRunsByJob(productCfg.Config[0].Out.JobID)
		assert.NoError(ct, err)
		assert.NotEmpty(ct, jobRuns.Data, "Expected the VRF job to have run at least once")

		out, err := consumer.RandomnessOutput(&bind.CallOpts{Context: ctx})
		assert.NoError(ct, err)
		assert.NotZero(ct, out.Uint64(), "Expected the VRF job to produce a non-zero randomness output")
	}, 2*time.Minute, 2*time.Second)
}

func TestVRFJobReplacement(t *testing.T) {
	ctx := t.Context()
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	productCfg, err := products.LoadOutput[vrf.Configurator](outputFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	c, auth, _, err := products.ETHClient(
		ctx,
		in.Blockchains[0].Out.Nodes[0].ExternalWSUrl,
		productCfg.Config[0].GasSettings.FeeCapMultiplier,
		productCfg.Config[0].GasSettings.TipCapMultiplier,
	)
	require.NoError(t, err)

	consumer, err := solidity_vrf_consumer_interface.NewVRFConsumer(
		common.HexToAddress(productCfg.Config[0].Out.ConsumerAddress), c,
	)
	require.NoError(t, err)

	keyHash := decodeKeyHash(t, productCfg.Config[0].Out.KeyHash)

	cls, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	// First randomness request
	tx, err := consumer.TestRequestRandomness(auth, keyHash, big.NewInt(1))
	require.NoError(t, err)
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	require.NoError(t, err)

	jobID := productCfg.Config[0].Out.JobID
	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		jobRuns, err := cls[0].MustReadRunsByJob(jobID)
		assert.NoError(ct, err)
		assert.NotEmpty(ct, jobRuns.Data, "Expected the VRF job to have run at least once")

		out, err := consumer.RandomnessOutput(&bind.CallOpts{Context: ctx})
		assert.NoError(ct, err)
		assert.NotZero(ct, out.Uint64(), "Expected the VRF job to produce a non-zero randomness output")
	}, 2*time.Minute, 2*time.Second)

	// Delete the job and recreate it
	err = cls[0].MustDeleteJob(jobID)
	require.NoError(t, err)

	cfg := productCfg.Config[0].Out
	pipelineSpec := &clclient.VRFTxPipelineSpec{
		Address: cfg.CoordinatorAddress,
	}
	observationSource, err := pipelineSpec.String()
	require.NoError(t, err)

	newJob, err := cls[0].MustCreateJob(&clclient.VRFJobSpec{
		Name:                     "vrf-" + cfg.ExternalJobID,
		CoordinatorAddress:       cfg.CoordinatorAddress,
		MinIncomingConfirmations: 1,
		PublicKey:                cfg.PublicKeyCompressed,
		ExternalJobID:            cfg.ExternalJobID,
		EVMChainID:               cfg.ChainID,
		ObservationSource:        observationSource,
	})
	require.NoError(t, err)

	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		jobRuns, err := cls[0].MustReadRunsByJob(newJob.Data.ID)
		assert.NoError(ct, err)
		assert.NotEmpty(ct, jobRuns.Data, "Expected the recreated VRF job to have run at least once")

		out, err := consumer.RandomnessOutput(&bind.CallOpts{Context: ctx})
		assert.NoError(ct, err)
		assert.NotZero(ct, out.Uint64(), "Expected the VRF job to produce a non-zero randomness output")
	}, 2*time.Minute, 2*time.Second)
}

func decodeKeyHash(t *testing.T, keyHashHex string) [32]byte {
	t.Helper()
	b, err := hex.DecodeString(keyHashHex)
	require.NoError(t, err, "Failed to decode key hash hex")
	var kh [32]byte
	copy(kh[:], b)
	return kh
}
