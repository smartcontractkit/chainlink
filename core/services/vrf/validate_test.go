package vrf

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/stretchr/testify/require"
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
confirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
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
				require.NoError(t, err)
				require.NotNil(t, s.VRFSpec)
				assert.Equal(t, uint32(10), s.VRFSpec.Confirmations)
				assert.Equal(t, "0xB3b7874F13387D44a3398D298B075B7A3505D8d4", s.VRFSpec.CoordinatorAddress.String())
				assert.Equal(t, "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179800", s.VRFSpec.PublicKey.String())
			},
		},
		{
			name: "missing pubkey",
			toml: `
type            = "vrf"
schemaVersion   = 1
confirmations = 10
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
				require.True(t, ErrKeyNotSet == errors.Cause(err))
			},
		},
		{
			name: "missing coordinator address",
			toml: `
type            = "vrf"
schemaVersion   = 1
confirmations = 10
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
				require.True(t, ErrKeyNotSet == errors.Cause(err))
			},
		},
		{
			name: "jobID override default",
			toml: `
type            = "vrf"
schemaVersion   = 1
confirmations = 10
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
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := ValidatedVRFSpec(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
