package job

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrNoChainFromSpec       = errors.New("could not get chain from spec")
	ErrNoSendingKeysFromSpec = errors.New("could not get sending keys from spec")
)

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
