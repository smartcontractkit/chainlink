package offchainreporting

import (
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"go.uber.org/multierr"
)

func boundTimeDuration(
	duration time.Duration,
	durationName string,
	min time.Duration,
	max time.Duration,
) error {
	if !(min <= duration && duration <= max) {
		return errors.Errorf("%s must be between %v and %v, but is currently %v",
			durationName, min, max, duration)
	}
	return nil
}

func SanityCheckLocalConfig(c types.LocalConfig) (err error) {
	if c.DevelopmentMode == types.EnableDangerousDevelopmentMode {
		return nil
	}

	err = multierr.Append(err,
		boundTimeDuration(
			c.BlockchainTimeout,
			"blockchain timeout",
			1*time.Second, 20*time.Second,
		))
	err = multierr.Append(err,
		boundTimeDuration(
			c.ContractConfigTrackerPollInterval,
			"contract config tracker poll interval",
			15*time.Second, 120*time.Second,
		))
	err = multierr.Append(err,
		boundTimeDuration(
			c.ContractConfigTrackerSubscribeInterval,
			"contract config tracker subscribe interval",
			1*time.Minute, 5*time.Minute,
		))
	err = multierr.Append(err,
		boundTimeDuration(
			c.ContractTransmitterTransmitTimeout,
			"contract transmitter transmit timeout",
			1*time.Second, 1*time.Minute,
		))
	err = multierr.Append(err,
		boundTimeDuration(
			c.DatabaseTimeout,
			"database timeout",
			100*time.Millisecond, 10*time.Second,
		))
	err = multierr.Append(err,
		boundTimeDuration(
			c.DataSourceTimeout,
			"data source timeout",
			1*time.Second, 20*time.Second,
		))
	err = multierr.Append(err,
		boundTimeDuration(
			c.DataSourceGracePeriod,
			"data source grace period",
			10*time.Millisecond, 2*time.Second,
		))

	if c.DataSourceTimeout < c.DataSourceGracePeriod {
		err = multierr.Append(err, errors.Errorf(
			"data source timeout %v must be greater than data source grace period %v",
			c.DataSourceTimeout,
			c.DataSourceTimeout))
	}

	const minContractConfigConfirmations = 1
	const maxContractConfigConfirmations = 100
	if !(minContractConfigConfirmations <= c.ContractConfigConfirmations && c.ContractConfigConfirmations <= maxContractConfigConfirmations) {
		err = multierr.Append(err, errors.Errorf(
			"contract config block-depth confirmation threshold must be between %v and %v, but is currently %v",
			minContractConfigConfirmations,
			maxContractConfigConfirmations,
			c.ContractConfigConfirmations))

	}

	return err
}
