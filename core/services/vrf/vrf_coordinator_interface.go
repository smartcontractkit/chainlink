package vrf

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"chainlink/core/assets"
	"chainlink/core/eth"
	"chainlink/core/services/vrf/generated/solidity_vrf_coordinator_interface"
	"chainlink/core/utils"
)

// RawRandomnessRequestLog is used to parse a RandomnessRequest log into types
// go-ethereum knows about.
type RawRandomnessRequestLog solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest

// RandomnessRequestLog contains the data for a RandomnessRequest log,
// represented as compatible golang types.
type RandomnessRequestLog struct {
	KeyHash common.Hash
	Seed    *big.Int // uint256
	JobID   common.Hash
	Sender  common.Address
	Fee     *assets.Link
	Raw     RawRandomnessRequestLog
}

// ParseRandomnessRequestLog returns the RandomnessRequestLog corresponding to
// the raw logData
func ParseRandomnessRequestLog(log eth.Log) (*RandomnessRequestLog, error) {
	l := RawRandomnessRequestLog{}
	coordABI := CoordinatorABI()
	contract := bind.NewBoundContract(common.Address{}, coordABI, nil, nil, nil)
	ethLog := types.Log(log)
	if err := contract.UnpackLog(&l, "RandomnessRequest", ethLog); err != nil {
		return nil, errors.Wrapf(err,
			"while parsing %x as RandomnessRequestLog", log.Data)
	}
	return &RandomnessRequestLog{l.KeyHash, l.Seed, l.JobID, l.Sender,
		(*assets.Link)(l.Fee), RawRandomnessRequestLog(l)}, nil
}

// RawData returns the raw bytes corresponding to l in a solidity log
//
// This serialization does not include the JobID, because that's an indexed field.
func (l *RandomnessRequestLog) RawData() ([]byte, error) {
	return randomnessRequestRawDataArgs().Pack(l.KeyHash,
		l.Seed, l.Sender, (*big.Int)(l.Fee))
}

// Equal(ol) is true iff l is the same log as ol, and both represent valid
// RandomnessRequest logs.
func (l *RandomnessRequestLog) Equal(ol RandomnessRequestLog) bool {
	return l.KeyHash == ol.KeyHash && equal(l.Seed, ol.Seed) &&
		l.JobID == ol.JobID && l.Sender == ol.Sender && l.Fee.Cmp(ol.Fee) == 0
}

func (l *RandomnessRequestLog) RequestID() common.Hash {
	soliditySeed, err := utils.Uint256ToBytes(l.Seed)
	if err != nil {
		panic(errors.Wrapf(err, "vrf seed out of bounds in %#+v", l))
	}
	return utils.MustHash(string(append(l.KeyHash[:], soliditySeed...)))
}

func rawRandomnessRequestLogToRandomnessRequestLog(
	l RawRandomnessRequestLog) RandomnessRequestLog {
	return RandomnessRequestLog{
		KeyHash: l.KeyHash,
		Seed:    l.Seed,
		JobID:   l.JobID,
		Sender:  l.Sender,
		Fee:     (*assets.Link)(l.Fee),
		Raw:     l,
	}
}
