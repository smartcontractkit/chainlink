package queue

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

const (
	maxObservedBlobHandles = 10   // sanity check on blobs in an observation
	maxObservedEvents      = 1000 // sanity check events in an observation
	maxKeyValueValueBytes  = ocr3_1types.MaxMaxKeyValueValueBytes
)

type PluginConfig struct {
	Limits ocr3_1types.ReportingPluginLimits `json:"limits,omitzero"`
}

// Validate checks the limits are within libocr's bounds, and that this plugin can operate within them.
func (c *PluginConfig) Validate() error {
	if c == nil {
		return errors.New("plugin config is missing")
	}
	if err := c.Limits.Validate(); err != nil {
		return fmt.Errorf("invalid limits: %w", err)
	}

	var err error
	positive := func(name string, v int) {
		if v <= 0 {
			err = errors.Join(err, fmt.Errorf("%s (%d) must be positive", name, v))
		}
	}

	positive("MaxObservationBytes", c.Limits.MaxObservationBytes)
	positive("MaxBlobPayloadBytes", c.Limits.MaxBlobPayloadBytes)
	positive("MaxKeyValueModifiedKeysPlusValuesBytes", c.Limits.MaxKeyValueModifiedKeysPlusValuesBytes)
	// the quotas must admit at least one full observation's blobs, or broadcasts are rejected and events are dropped
	if c.Limits.MaxPerOracleUnexpiredBlobCount < maxObservedBlobHandles {
		err = errors.Join(err, fmt.Errorf("MaxPerOracleUnexpiredBlobCount (%d) must be at least %d, the blobs in a full observation",
			c.Limits.MaxPerOracleUnexpiredBlobCount, maxObservedBlobHandles))
	}

	if c.Limits.MaxKeyValueModifiedKeys < numKVKeys {
		err = errors.Join(err, fmt.Errorf("MaxKeyValueModifiedKeys (%d) must be at least %d, the number of keys written each round", c.Limits.MaxKeyValueModifiedKeys, numKVKeys))
	}

	if fullObservation := maxObservedBlobHandles * c.Limits.MaxBlobPayloadBytes; c.Limits.MaxPerOracleUnexpiredBlobCumulativePayloadBytes < fullObservation {
		err = errors.Join(err, fmt.Errorf("MaxPerOracleUnexpiredBlobCumulativePayloadBytes (%d) must hold the blobs of a full observation (%d x MaxBlobPayloadBytes %d)",
			c.Limits.MaxPerOracleUnexpiredBlobCumulativePayloadBytes, maxObservedBlobHandles, c.Limits.MaxBlobPayloadBytes))
	}

	return err
}
