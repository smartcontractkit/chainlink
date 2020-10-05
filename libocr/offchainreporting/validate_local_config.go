package offchainreporting

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"go.uber.org/multierr"
)

func boundTimeDuration(
	duration float64, durationName, unitName string, min, max float64,
) error {
	if duration > max {
		return errors.Errorf("%s timeout %.1f%s is unreasonably large",
			durationName, duration, unitName)
	}
	if duration < min {
		return errors.Errorf("%s timeout %.1f%s is unreasonably small",
			durationName, duration, unitName)
	}
	return nil
}

func validateLocalConfig(c types.LocalConfig) (err error) {
	err = multierr.Append(err,
		boundTimeDuration(
			c.DataSourceTimeout.Seconds(), "data source timeout", "s",
			0.001 /* 1ms */, 20,
		))
	err = multierr.Append(err,
		boundTimeDuration(
			c.BlockchainTimeout.Seconds(), "block chain timeout", "s",
			0.001, 20,
		))
	err = multierr.Append(err,
		boundTimeDuration(
			c.ContractConfigTrackerPollInterval.Seconds(),
			"contract config tracker poll interval", "s", 15, 120,
		))
	err = multierr.Append(err,
		boundTimeDuration(
			c.ContractConfigTrackerSubscribeInterval.Minutes(),
			"contract config tracker subscribe interval", "min", 2, 5,
		))
	if c.ContractConfigConfirmations < 2 {
		err = multierr.Append(err, errors.Errorf(
			"contract config block-depth confirmation threshold %d is unreasonably low",
			c.ContractConfigConfirmations))
	}
	if c.ContractConfigConfirmations > 10 {
		err = multierr.Append(err, errors.Errorf(
			"contract config block-depth confirmation threshold %d is unreasonably large",
			c.ContractConfigConfirmations))
	}
	return err
}
