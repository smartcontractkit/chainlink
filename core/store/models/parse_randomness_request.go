package models

import (
	"github.com/pkg/errors"

	"chainlink/core/eth"
	"chainlink/core/services/vrf"
	"chainlink/core/utils"
)

// parseRandomnessRequest parses the RandomnessRequest log format.
type parseRandomnessRequest struct{}

var _ logRequestParser = parseRandomnessRequest{} // Implements logRequestParser

// parseJSON returns the seed for the RandomnessRequest
func (parseRandomnessRequest) parseJSON(log eth.Log) (js JSON, err error) {
	parsedLog, err := vrf.ParseRandomnessRequestLog(log)
	if err != nil {
		return JSON{}, errors.Wrapf(err,
			"could not parse log data %x as RandomnessRequest log", log.Data)
	}
	fullSeedString, err := utils.Uint256ToHex(parsedLog.Seed)
	if err != nil {
		return JSON{}, errors.Wrap(err, "vrf seed out of bounds")
	}
	add := func(k, v string) { js = mustAdd(js, k, v) }
	add("address", log.Address.String())
	add("functionSelector", vrf.FulfillSelector())
	add("keyHash", parsedLog.KeyHash.Hex())
	add("seed", fullSeedString)
	add("jobID", parsedLog.JobID.Hex())
	add("sender", parsedLog.Sender.Hex())
	return js, nil
}

func (parseRandomnessRequest) parseRequestID(log eth.Log) string {
	parsedLog, err := vrf.ParseRandomnessRequestLog(log)
	if err != nil {
		panic(errors.Wrapf(err, "while extracting randomness requestID from %#+v", log))
	}
	return parsedLog.RequestID().Hex()
}
