package vrf

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"chainlink/core/assets"
	coord "chainlink/core/services/vrf/generated/solidity_vrf_coordinator_interface"
	"chainlink/core/utils"
)

var fulfillMethodName = "fulfillRandomnessRequest"

// CoordinatorABI is the ABI of the VRFCoordinator
var CoordinatorABI abi.ABI

var FulfillMethod abi.Method

// FulfillSelector is the function selector of fulfillRandomness, the main
// entrypoint to the VRFCoordinator.
var FulfillSelector string

// RandomnessRequestLogTopic is the signature of the RandomnessRequest log
var RandomnessRequestLogTopic common.Hash
var RandomnessRequestABI abi.Event

func init() {
	CoordinatorABI, err := abi.JSON(strings.NewReader(coord.VRFCoordinatorABI))
	if err != nil {
		panic(err)
	}
	for methodName, method := range CoordinatorABI.Methods {
		if methodName == fulfillMethodName {
			FulfillMethod = method
			FulfillSelector = hexutil.Encode(method.ID())
		}
	}
	if FulfillSelector == "" {
		panic("failed to find fulfill method")
	}
	RandomnessRequestABI = CoordinatorABI.Events["RandomnessRequest"]
	RandomnessRequestLogTopic = RandomnessRequestABI.ID()
}

// RandomnessRequestLog contains the data for a RandomnessRequest log,
// represented as compatible golang types.
type RandomnessRequestLog struct {
	KeyHash common.Hash
	Seed    *big.Int // uint256
	JobID   common.Hash
	Sender  common.Address
	Fee     *assets.Link
}

// rawRandomnessRequestLog is used to parse a RandomnessRequest log into types
// go-ethereum knows about.
type rawRandomnessRequestLog struct {
	KeyHash common.Hash
	Seed    *big.Int
	JobID   common.Hash
	Sender  common.Address
	Fee     *big.Int
}

// ParseRandomnessRequestLog returns the RandomnessRequestLog corresponding to
// the raw logData
func ParseRandomnessRequestLog(logData []byte) (*RandomnessRequestLog, error) {
	rawRV := rawRandomnessRequestLog{}
	if err := RandomnessRequestABI.Inputs.Unpack(&rawRV, logData); err != nil {
		return nil, errors.Wrapf(err, "while unpacking RandomnessRequest log data")
	}
	rv := RandomnessRequestLog{
		KeyHash: rawRV.KeyHash,
		Seed:    rawRV.Seed,
		JobID:   rawRV.JobID,
		Sender:  rawRV.Sender,
		Fee:     (*assets.Link)(rawRV.Fee),
	}
	checkUint256(rv.Seed)
	checkUint256((*big.Int)(rv.Fee))
	return &rv, nil
}

func checkUint256(n *big.Int) {
	if err := utils.CheckUint256(n); err != nil {
		panic(fmt.Errorf(
			"go-ethereum returned something out-of-bounds for a uint256: %v", n))
	}
}

// RawLog returns the raw bytes corresponding to l in a solidity log
func (l *RandomnessRequestLog) RawLog() ([]byte, error) {
	return RandomnessRequestABI.Inputs.Pack(l.KeyHash, l.Seed, l.JobID, l.Sender,
		(*big.Int)(l.Fee))
}

// Equal(ol) is true iff l is the same log as ol, and both represent valid
// RandomnessRequest logs.
func (l *RandomnessRequestLog) Equal(ol RandomnessRequestLog) bool {
	lR, err := l.RawLog()
	oR, oErr := ol.RawLog()
	return bytes.Equal(lR, oR) && err == nil && oErr == nil
}

func (l *RandomnessRequestLog) RequestID() common.Hash {
	soliditySeed, err := utils.Uint256ToBytes(l.Seed)
	if err != nil {
		panic(errors.Wrapf(err, "vrf seed out of bounds in %#+v", l))
	}
	hash, err := utils.Keccak256(append(l.KeyHash[:], soliditySeed...))
	if err != nil {
		panic(errors.Wrapf(err, "this should never happen"))
	}
	return common.BytesToHash(hash)
}
