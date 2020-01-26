package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"chainlink/core/services/vrf"
)

// RandomnessLogEvent provides functionality specific to a log event emitted
// for a run log initiator.
type RandomnessLogEvent struct{ InitiatorLogEvent }

// assert RandomnessLogEvent implements LogRequest interface
var _ LogRequest = RandomnessLogEvent{}

// Validate returns whether or not the contained log has a properly encoded
// job id.
func (le RandomnessLogEvent) Validate() bool {
	_, err := vrf.ParseRandomnessRequestLog(le.Log)
	return err == nil
}

// ValidateRequester never errors, because the requester is not important to the
// Randomness functionality. XXX: This could result in DoS attacks. Maybe we want a
// whitelist of requesters. It could also result in someone who for some reason
// illegitimately knows what the seed is going to be being able to request the
// result for that seed ahead of time.
func (le RandomnessLogEvent) ValidateRequester() error {
	return nil // XXX: See doc string
}

// Requester pulls the requesting address out of the LogEvent's topics. XXX:
// This is not the requester as the Chainlink oracle understands it. This is the
// oracle itself.
func (le RandomnessLogEvent) Requester() common.Address {
	return le.Log.Address
}

// RunRequest returns a RunRequest instance with all parameters
// from a run log topic, like RequestID.
func (le RandomnessLogEvent) RunRequest() (RunRequest, error) {
	parsedLog, err := vrf.ParseRandomnessRequestLog(le.Log)
	if err != nil {
		return RunRequest{}, errors.Wrapf(err, "while parsing log for run request")
	}

	payment := parsedLog.Fee

	txHash := common.BytesToHash(le.Log.TxHash.Bytes())
	blockHash := common.BytesToHash(le.Log.BlockHash.Bytes())
	str := parsedLog.RequestID().Hex()
	requester := le.Requester()
	return RunRequest{
		RequestID: &str,
		TxHash:    &txHash,
		BlockHash: &blockHash,
		Requester: &requester,
		Payment:   payment,
	}, nil
}

func (le RandomnessLogEvent) JSON() (js JSON, err error) {
	return parseRandomnessRequest{}.parseJSON(le.Log)
}
