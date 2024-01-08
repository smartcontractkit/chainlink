package vrfcommon

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestValidateVRFJobSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "valid spec",
			toml: `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
vrfOwnerAddress = "0x2a0d386f122851dc5AFBE45cb2E8411CE255b000"
requestTimeout = "168h" # 7 days
chunkSize = 25
backoffInitialDelay = "1m"
backoffMaxDelay = "2h"
observationSource = """
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
			  publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.NoError(t, err)
				require.NotNil(t, s.VRFSpec)
				assert.Equal(t, uint32(10), s.VRFSpec.MinIncomingConfirmations)
				assert.Equal(t, "0xB3b7874F13387D44a3398D298B075B7A3505D8d4", s.VRFSpec.CoordinatorAddress.String())
				assert.Equal(t, "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179800", s.VRFSpec.PublicKey.String())
				assert.Equal(t, "0x2a0d386f122851dc5AFBE45cb2E8411CE255b000", s.VRFSpec.VRFOwnerAddress.String())
				require.Equal(t, 168*time.Hour, s.VRFSpec.RequestTimeout)
				require.Equal(t, time.Minute, s.VRFSpec.BackoffInitialDelay)
				require.Equal(t, 2*time.Hour, s.VRFSpec.BackoffMaxDelay)
				require.EqualValues(t, 25, s.VRFSpec.ChunkSize)
			},
		},
		{
			name: "missing pubkey",
			toml: `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
observationSource = """
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
			  publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				require.True(t, errors.Is(ErrKeyNotSet, errors.Cause(err)))
			},
		},
		{
			name: "missing fromAddresses",
			toml: `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
requestTimeout = "168h" # 7 days
chunkSize = 25
backoffInitialDelay = "1m"
backoffMaxDelay = "2h"
observationSource = """
decode_log   [type=ethabidecodelog
				abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
				data="$(jobRun.logData)"
				topics="$(jobRun.logTopics)"]
vrf          [type=vrfv2
				publicKey="$(jobSpec.publicKey)"
				requestBlockHash="$(jobRun.logBlockHash)"
				requestBlockNumber="$(jobRun.logBlockNumber)"
				topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
				abi="fulfillRandomnessRequest(bytes proof)"
				data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
			txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
			`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				require.True(t, errors.Is(ErrKeyNotSet, errors.Cause(err)))
			},
		},
		{
			name: "missing coordinator address",
			toml: `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
observationSource = """
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
			  publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				require.True(t, errors.Is(ErrKeyNotSet, errors.Cause(err)))
			},
		},
		{
			name: "jobID override default",
			toml: `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
observationSource = """
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
			  publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.NoError(t, err)
				assert.Equal(t, s.ExternalJobID.String(), "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46")
			},
		},
		{
			name: "no requested confs delay",
			toml: `
			type            = "vrf"
			schemaVersion   = 1
			minIncomingConfirmations = 10
			publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
			coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
			externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
			observationSource = """
			decode_log   [type=ethabidecodelog
						  abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
						  data="$(jobRun.logData)"
						  topics="$(jobRun.logTopics)"]
			vrf          [type=vrf
						  publicKey="$(jobSpec.publicKey)"
						  requestBlockHash="$(jobRun.logBlockHash)"
						  requestBlockNumber="$(jobRun.logBlockNumber)"
						  topics="$(jobRun.logTopics)"]
			encode_tx    [type=ethabiencode
						  abi="fulfillRandomnessRequest(bytes proof)"
						  data="{\\"proof\\": $(vrf)}"]
			submit_tx  [type=ethtx to="%s"
						data="$(encode_tx)"
						txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
			decode_log->vrf->encode_tx->submit_tx
			"""
			`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(0), os.VRFSpec.RequestedConfsDelay)
			},
		},
		{
			name: "with requested confs delay",
			toml: `
			type            = "vrf"
			schemaVersion   = 1
			minIncomingConfirmations = 10
			requestedConfsDelay = 10
			publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
			coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
			externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
			observationSource = """
			decode_log   [type=ethabidecodelog
						  abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
						  data="$(jobRun.logData)"
						  topics="$(jobRun.logTopics)"]
			vrf          [type=vrf
						  publicKey="$(jobSpec.publicKey)"
						  requestBlockHash="$(jobRun.logBlockHash)"
						  requestBlockNumber="$(jobRun.logBlockNumber)"
						  topics="$(jobRun.logTopics)"]
			encode_tx    [type=ethabiencode
						  abi="fulfillRandomnessRequest(bytes proof)"
						  data="{\\"proof\\": $(vrf)}"]
			submit_tx  [type=ethtx to="%s"
						data="$(encode_tx)"
						txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
			decode_log->vrf->encode_tx->submit_tx
			"""
			`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(10), os.VRFSpec.RequestedConfsDelay)
			},
		},
		{
			name: "negative (illegal) requested confs delay",
			toml: `
			type            = "vrf"
			schemaVersion   = 1
			minIncomingConfirmations = 10
			requestedConfsDelay = -10
			publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
			coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
			externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
			observationSource = """
			decode_log   [type=ethabidecodelog
						  abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
						  data="$(jobRun.logData)"
						  topics="$(jobRun.logTopics)"]
			vrf          [type=vrf
						  publicKey="$(jobSpec.publicKey)"
						  requestBlockHash="$(jobRun.logBlockHash)"
						  requestBlockNumber="$(jobRun.logBlockNumber)"
						  topics="$(jobRun.logTopics)"]
			encode_tx    [type=ethabiencode
						  abi="fulfillRandomnessRequest(bytes proof)"
						  data="{\\"proof\\": $(vrf)}"]
			submit_tx  [type=ethtx to="%s"
						data="$(encode_tx)"
						txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
			decode_log->vrf->encode_tx->submit_tx
			"""
			`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "no request timeout provided, sets default of 1 day",
			toml: `
			type            = "vrf"
			schemaVersion   = 1
			minIncomingConfirmations = 10
			requestedConfsDelay = 10
			publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
			coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
			externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
			observationSource = """
			decode_log   [type=ethabidecodelog
						  abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
						  data="$(jobRun.logData)"
						  topics="$(jobRun.logTopics)"]
			vrf          [type=vrf
						  publicKey="$(jobSpec.publicKey)"
						  requestBlockHash="$(jobRun.logBlockHash)"
						  requestBlockNumber="$(jobRun.logBlockNumber)"
						  topics="$(jobRun.logTopics)"]
			encode_tx    [type=ethabiencode
						  abi="fulfillRandomnessRequest(bytes proof)"
						  data="{\\"proof\\": $(vrf)}"]
			submit_tx  [type=ethtx to="%s"
						data="$(encode_tx)"
						txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
			decode_log->vrf->encode_tx->submit_tx
			"""
			`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, 24*time.Hour, os.VRFSpec.RequestTimeout)
			},
		},
		{
			name: "request timeout provided, uses that",
			toml: `
			type            = "vrf"
			schemaVersion   = 1
			minIncomingConfirmations = 10
			requestedConfsDelay = 10
			requestTimeout = "168h" # 7 days
			publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
			coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
			externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
			observationSource = """
			decode_log   [type=ethabidecodelog
						  abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
						  data="$(jobRun.logData)"
						  topics="$(jobRun.logTopics)"]
			vrf          [type=vrf
						  publicKey="$(jobSpec.publicKey)"
						  requestBlockHash="$(jobRun.logBlockHash)"
						  requestBlockNumber="$(jobRun.logBlockNumber)"
						  topics="$(jobRun.logTopics)"]
			encode_tx    [type=ethabiencode
						  abi="fulfillRandomnessRequest(bytes proof)"
						  data="{\\"proof\\": $(vrf)}"]
			submit_tx  [type=ethtx to="%s"
						data="$(encode_tx)"
						txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
			decode_log->vrf->encode_tx->submit_tx
			"""
			`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, 7*24*time.Hour, os.VRFSpec.RequestTimeout)
			},
		},
		{
			name: "batch fulfillment enabled, no batch coordinator address",
			toml: `
			type            = "vrf"
			schemaVersion   = 1
			minIncomingConfirmations = 10
			requestedConfsDelay = 10
			batchFulfillmentEnabled = true
			publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
			coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
			externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
			observationSource = """
			decode_log   [type=ethabidecodelog
						  abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
						  data="$(jobRun.logData)"
						  topics="$(jobRun.logTopics)"]
			vrf          [type=vrf
						  publicKey="$(jobSpec.publicKey)"
						  requestBlockHash="$(jobRun.logBlockHash)"
						  requestBlockNumber="$(jobRun.logBlockNumber)"
						  topics="$(jobRun.logTopics)"]
			encode_tx    [type=ethabiencode
						  abi="fulfillRandomnessRequest(bytes proof)"
						  data="{\\"proof\\": $(vrf)}"]
			submit_tx  [type=ethtx to="%s"
						data="$(encode_tx)"
						txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
			decode_log->vrf->encode_tx->submit_tx
			"""
			`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "batch fulfillment enabled, batch coordinator address provided",
			toml: `
			type            = "vrf"
			schemaVersion   = 1
			minIncomingConfirmations = 10
			requestedConfsDelay = 10
			batchFulfillmentEnabled = true
			batchCoordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
			publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
			coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
			externalJobID = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
			observationSource = """
			decode_log   [type=ethabidecodelog
						  abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
						  data="$(jobRun.logData)"
						  topics="$(jobRun.logTopics)"]
			vrf          [type=vrf
						  publicKey="$(jobSpec.publicKey)"
						  requestBlockHash="$(jobRun.logBlockHash)"
						  requestBlockNumber="$(jobRun.logBlockNumber)"
						  topics="$(jobRun.logTopics)"]
			encode_tx    [type=ethabiencode
						  abi="fulfillRandomnessRequest(bytes proof)"
						  data="{\\"proof\\": $(vrf)}"]
			submit_tx  [type=ethtx to="%s"
						data="$(encode_tx)"
						txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
			decode_log->vrf->encode_tx->submit_tx
			"""
			`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, "0xB3b7874F13387D44a3398D298B075B7A3505D8d4", os.VRFSpec.BatchCoordinatorAddress.String())
			},
		},
		{
			name: "initial delay must be <= max delay, invalid",
			toml: `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
requestTimeout = "168h" # 7 days
chunkSize = 25
backoffInitialDelay = "1h"
backoffMaxDelay = "30m"
observationSource = """
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
			  publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "gas lane price provided",
			toml: `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
requestTimeout = "168h" # 7 days
chunkSize = 25
backoffInitialDelay = "1m"
backoffMaxDelay = "2h"
gasLanePrice = "200 gwei"
observationSource = """
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
			  publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.NoError(t, err)
				require.NotNil(t, s.VRFSpec)
				assert.Equal(t, uint32(10), s.VRFSpec.MinIncomingConfirmations)
				assert.Equal(t, "0xB3b7874F13387D44a3398D298B075B7A3505D8d4", s.VRFSpec.CoordinatorAddress.String())
				assert.Equal(t, "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179800", s.VRFSpec.PublicKey.String())
				require.Equal(t, 168*time.Hour, s.VRFSpec.RequestTimeout)
				require.Equal(t, time.Minute, s.VRFSpec.BackoffInitialDelay)
				require.Equal(t, 2*time.Hour, s.VRFSpec.BackoffMaxDelay)
				require.EqualValues(t, 25, s.VRFSpec.ChunkSize)
				require.Equal(t, assets.GWei(200), s.VRFSpec.GasLanePrice)
			},
		},
		{
			name: "invalid (negative) gas lane price provided",
			toml: `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
requestTimeout = "168h" # 7 days
chunkSize = 25
backoffInitialDelay = "1m"
backoffMaxDelay = "2h"
gasLanePrice = "-200"
observationSource = """
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
			  publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid (zero) gas lane price gwei provided",
			toml: `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
requestTimeout = "168h" # 7 days
chunkSize = 25
backoffInitialDelay = "1m"
backoffMaxDelay = "2h"
gasLanePrice = "0 gwei"
observationSource = """
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
			  publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := ValidatedVRFSpec(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
