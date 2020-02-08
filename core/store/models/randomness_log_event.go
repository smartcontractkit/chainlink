package models

import (
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"chainlink/core/logger"
	"chainlink/core/services/vrf"
)

// RandomnessLogEvent provides functionality specific to a log event emitted
// for a run log initiator.
type RandomnessLogEvent struct{ InitiatorLogEvent }

// assert RandomnessLogEvent implements LogRequest interface
var _ LogRequest = RandomnessLogEvent{}

var allHex = regexp.MustCompile("^[[:xdigit:]]{32}$").Match

// Validate() is true if the contained log is parseable as a RandomnessRequest,
// and it's from the address specified by the job's initiator.
func (le RandomnessLogEvent) Validate() bool {
	_, err := vrf.ParseRandomnessRequestLog(le.Log)
	switch {
	case err != nil:
		logger.Warnf("error while parsing RandomnessRequest log: %s on log %#+v",
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
// node's functionality
func (le RandomnessLogEvent) ValidateRequester() error {
	return nil
}

// Requester pulls the requesting address out of the LogEvent's topics.
func (le RandomnessLogEvent) Requester() common.Address {
	log, err := vrf.ParseRandomnessRequestLog(le.Log)
	if err != nil {
		logger.Warnf("error while parsing RandomnessRequest log: %s on log %#+v",
			err, le.Log)
	}
	return log.Sender
}

// RunRequest returns a RunRequest instance with all parameters
// from a run log topic, like RequestID.
func (le RandomnessLogEvent) RunRequest() (RunRequest, error) {
	parsedLog, err := vrf.ParseRandomnessRequestLog(le.Log)
	if err != nil {
		return RunRequest{}, errors.Wrapf(err, "while parsing log for run request")
	}

	str := parsedLog.RequestID().Hex()
	requester := le.Requester()
	return RunRequest{
		RequestID: &str,
		TxHash:    &le.Log.TxHash,
		BlockHash: &le.Log.BlockHash,
		Requester: &requester,
		Payment:   parsedLog.Fee,
	}, nil
}

func (le RandomnessLogEvent) JSON() (js JSON, err error) {
	return parseRandomnessRequest{}.parseJSON(le.Log)
}
