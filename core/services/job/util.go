package job

import (
	"fmt"
	"math/big"

	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
)

var (
	ErrNoChainFromSpec       = fmt.Errorf("could not get chain from spec")
	ErrNoSendingKeysFromSpec = fmt.Errorf("could not get sending keys from spec")
)

// EVMChainForJob parses the job spec and retrieves the evm chain found.
func EVMChainForBootstrapJob(job *Job, set evm.ChainSet) (evm.Chain, error) {
	chainIDInterface, ok := job.BootstrapSpec.RelayConfig["chainID"]
	if !ok {
		return nil, fmt.Errorf("%w: chainID must be provided in relay config", ErrNoChainFromSpec)
	}
	chainID := chainIDInterface.(int64)
	chain, err := set.Get(big.NewInt(chainID))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNoChainFromSpec, err)
	}

	return chain, nil
}

// EVMChainForJob parses the job spec and retrieves the evm chain found.
func EVMChainForJob(job *Job, set evm.ChainSet) (evm.Chain, error) {
	chainIDInterface, ok := job.OCR2OracleSpec.RelayConfig["chainID"]
	if !ok {
		return nil, fmt.Errorf("%w: chainID must be provided in relay config", ErrNoChainFromSpec)
	}
	chainID := chainIDInterface.(int64)
	chain, err := set.Get(big.NewInt(chainID))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNoChainFromSpec, err)
	}

	return chain, nil
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
