package job

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"

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
func SendingKeysForJob(job *Job) ([]string, error) {
	sendingKeysInterface, ok := job.OCR2OracleSpec.RelayConfig["sendingKeys"]
	if !ok {
		return nil, fmt.Errorf("%w: sendingKeys must be provided in relay config", ErrNoSendingKeysFromSpec)
	}

	sendingKeysInterfaceSlice, ok := sendingKeysInterface.([]interface{})
	if !ok {
		return nil, errors.New("sending keys should be an array")
	}

	var sendingKeys []string
	for _, sendingKeyInterface := range sendingKeysInterfaceSlice {
		sendingKey, ok := sendingKeyInterface.(string)
		if !ok {
			return nil, errors.New("sending keys are of wrong type")
		}
		sendingKeys = append(sendingKeys, sendingKey)
	}

	if len(sendingKeys) == 0 {
		return nil, errors.New("sending keys are empty")
	}

	return sendingKeys, nil
}
