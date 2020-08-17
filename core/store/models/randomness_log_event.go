package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// RandomnessLogEvent provides functionality specific to a log event emitted
// for a run log initiator.
type RandomnessLogEvent struct{ InitiatorLogEvent }

var _ LogRequest = RandomnessLogEvent{} // implements LogRequest interface

// Validate() is true if the contained log is parseable as a RandomnessRequest,
// and it's from the address specified by the job's initiator. The log filter
// and the go-ethereum parser should prevent any invalid logs from reacching
// this point, so Validate emits an error log on failure.
func (le RandomnessLogEvent) Validate() bool {
	_, err := ParseRandomnessRequestLog(le.Log)
	switch {
	case err != nil:
		logger.Errorf("error while parsing RandomnessRequest log: %s on log %#+v",
			err, le.Log)
		return false
	// Following should be guaranteed by log query filterer, but doesn't hurt to
	// check again.
	case le.Log.Address != le.Initiator.Address:
		logger.Errorf(
			"RandomnessRequest log received from address %s, but expect logs from %s",
			le.Log.Address.String(), le.Initiator.Address.String())
		return false
	}
	return true
}

// ValidateRequester never errors, because the requester is not important to the
// node's functionality. A requesting contract cannot request the VRF output on
// behalf of another contract, because the initial input seed is hashed with the
// requesting contract's address (plus a nonce) to get the actual VRF input.
func (le RandomnessLogEvent) ValidateRequester() error {
	return nil
}

// Requester pulls the requesting address out of the LogEvent's topics.
func (le RandomnessLogEvent) Requester() common.Address {
	log, err := ParseRandomnessRequestLog(le.Log)
	if err != nil {
		logger.Errorf("error while parsing RandomnessRequest log: %s on log %#+v",
			err, le.Log)
		return common.Address{}
	}
	return log.Sender
}

// RunRequest returns a RunRequest instance with all parameters
// from a run log topic, like RequestID.
func (le RandomnessLogEvent) RunRequest() (RunRequest, error) {
	parsedLog, err := ParseRandomnessRequestLog(le.Log)
	if err != nil {
		return RunRequest{}, errors.Wrapf(err,
			"while parsing log for VRF run request")
	}
	requestParams, err := le.JSON()
	if err != nil {
		return RunRequest{}, errors.Wrapf(err,
			"while parsing request params for VRF run request")
	}

	requestID := parsedLog.RequestID
	requester := le.Requester()
	return RunRequest{
		RequestID:     &requestID,
		TxHash:        &le.Log.TxHash,
		BlockHash:     &le.Log.BlockHash,
		Requester:     &requester,
		Payment:       parsedLog.Fee,
		RequestParams: requestParams,
	}, nil
}

// JSON returns the JSON from this RandomnessRequest log, as it will be passed
// to the Randomn adapter
func (le RandomnessLogEvent) JSON() (js JSON, err error) {
	return parseRandomnessRequest{}.parseJSON(le.Log)
}
