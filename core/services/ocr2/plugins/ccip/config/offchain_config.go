package config

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type OffchainConfig interface {
	Validate() error
}

// Do not change the JSON format of this struct without consulting with
// the RDD people first.
type CommitOffchainConfig struct {
	SourceFinalityDepth   uint32
	DestFinalityDepth     uint32
	FeeUpdateHeartBeat    models.Duration
	FeeUpdateDeviationPPB uint32
	MaxGasPrice           uint64
	InflightCacheExpiry   models.Duration
}

func (c CommitOffchainConfig) Validate() error {
	if c.SourceFinalityDepth == 0 {
		return errors.New("must set SourceFinalityDepth")
	}
	if c.DestFinalityDepth == 0 {
		return errors.New("must set DestFinalityDepth")
	}
	if c.FeeUpdateHeartBeat.Duration() == 0 {
		return errors.New("must set FeeUpdateHeartBeat")
	}
	if c.FeeUpdateDeviationPPB == 0 {
		return errors.New("must set FeeUpdateDeviationPPB")
	}
	if c.MaxGasPrice == 0 {
		return errors.New("must set MaxGasPrice")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}

	return nil
}

// Do not change the JSON format of this struct without consulting with
// the RDD people first.
type ExecOffchainConfig struct {
	SourceFinalityDepth         uint32
	DestOptimisticConfirmations uint32
	DestFinalityDepth           uint32
	BatchGasLimit               uint32
	RelativeBoostPerWaitHour    float64
	MaxGasPrice                 uint64
	InflightCacheExpiry         models.Duration
	RootSnoozeTime              models.Duration
}

func (c ExecOffchainConfig) Validate() error {
	if c.SourceFinalityDepth == 0 {
		return errors.New("must set SourceFinalityDepth")
	}
	if c.DestFinalityDepth == 0 {
		return errors.New("must set DestFinalityDepth")
	}
	if c.DestOptimisticConfirmations == 0 {
		return errors.New("must set DestOptimisticConfirmations")
	}
	if c.BatchGasLimit == 0 {
		return errors.New("must set BatchGasLimit")
	}
	if c.RelativeBoostPerWaitHour == 0 {
		return errors.New("must set RelativeBoostPerWaitHour")
	}
	if c.MaxGasPrice == 0 {
		return errors.New("must set MaxGasPrice")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}
	if c.RootSnoozeTime.Duration() == 0 {
		return errors.New("must set RootSnoozeTime")
	}

	return nil
}

func DecodeOffchainConfig[T OffchainConfig](encodedConfig []byte) (T, error) {
	var result T
	err := json.Unmarshal(encodedConfig, &result)
	if err != nil {
		return result, err
	}
	err = result.Validate()
	if err != nil {
		return result, err
	}
	return result, nil
}

func EncodeOffchainConfig[T OffchainConfig](occ T) ([]byte, error) {
	return json.Marshal(occ)
}
