package job

import (
	"fmt"
	"math/big"

	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
)

var (
	ErrNoChainFromSpec       = fmt.Errorf("could not get chain from spec")
	ErrInvalidChainFromSpec  = fmt.Errorf("could not parse chain id")
	ErrNoSendingKeysFromSpec = fmt.Errorf("could not get sending keys from spec")
)

// EVMChainForJob parses the job spec and retrieves the evm chain found.
func EVMChainForJob(job *Job, set evm.ChainSet) (evm.Chain, error) {
	chainID, err := EVMChainIDForJobSpec(job.OCR2OracleSpec)
	if err != nil {
		return nil, err
	}
	chain, err := set.Get(big.NewInt(chainID))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNoChainFromSpec, err)
	}

	return chain, nil
}

// EVMChainIDForJob parses the job spec and retrieves the evm chain id found
func EVMChainIDForJobSpec(spec *OCR2OracleSpec) (chainID int64, err error) {
	chainIDInterface, ok := spec.RelayConfig["chainID"]
	if !ok {
		return 0, fmt.Errorf("%w: chainID must be provided in relay config", ErrNoChainFromSpec)
	}
	chainID, ok = chainIDInterface.(int64)
	if !ok {
		return 0, fmt.Errorf("%w: invalid chainID field in relay config, must be int64", ErrInvalidChainFromSpec)
	}
	return
}

// SendingKeysForJob parses the job spec and retrieves the sending keys found.
func SendingKeysForJob(job *Job) (pq.StringArray, error) {
	sendingKeysInterface, ok := job.OCR2OracleSpec.RelayConfig["sendingKeys"]
	if !ok {
		return nil, fmt.Errorf("%w: sendingKeys must be provided in relay config", ErrNoSendingKeysFromSpec)
	}
	sendingKeys := sendingKeysInterface.(pq.StringArray)

	return sendingKeys, nil
}
