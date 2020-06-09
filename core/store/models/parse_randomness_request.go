package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// parseRandomnessRequest parses the RandomnessRequest log format.
type parseRandomnessRequest struct{}

var _ logRequestParser = parseRandomnessRequest{} // Implements logRequestParser

// parseJSON returns the inputs to be passed as a JSON object to Random adapter
func (parseRandomnessRequest) parseJSON(log Log) (js JSON, err error) {
	parsedLog, err := ParseRandomnessRequestLog(log)
	if err != nil {
		return JSON{}, errors.Wrapf(err,
			"could not parse log data %#+v as RandomnessRequest log", log)
	}
	fullSeedString, err := utils.Uint256ToHex(parsedLog.Seed)
	if err != nil {
		return JSON{}, errors.Wrap(err, "vrf seed out of bounds")
	}
	return js.MultiAdd(KV{
		// Address of log emitter
		"address": log.Address.String(),
		// Signature of callback function on consuming contract
		"functionSelector": VRFFulfillSelector(),
		// Hash of the public key for the VRF to be used
		"keyHash": parsedLog.KeyHash.Hex(),
		// Raw input seed for the VRF (includes requester, nonce, etc.)
		"seed": fullSeedString,
		// The chainlink job corresponding to this VRF
		"jobID": parsedLog.JobID.Hex(),
		// Address of consuming contract which initially made the request
		"sender": parsedLog.Sender.Hex(),
		// Hash of the block in which this request appeared
		"blockHash": log.BlockHash.Hex(),
		// Number/height of the block in which this request appeared
		"blockNum": log.BlockNumber,
	})
}

func (parseRandomnessRequest) parseRequestID(log Log) (common.Hash, error) {
	parsedLog, err := ParseRandomnessRequestLog(log)
	if err != nil {
		return common.Hash{}, errors.Wrapf(err, "while extracting randomness requestID from %#+v", log)
	}
	return parsedLog.RequestID(), nil
}
