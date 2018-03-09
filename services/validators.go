package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

// ValidateJob checks the job and its associated Initiators and Tasks for any
// application logic errors.
func ValidateJob(j models.JobSpec, store *store.Store) error {
	var merr error
	if j.StartAt.Valid && j.EndAt.Valid && j.StartAt.Time.After(j.EndAt.Time) {
		merr = multierr.Append(merr, fmtJobError(errors.New("startat cannot be before endat")))
	}
	for _, i := range j.Initiators {
		if err := ValidateInitiator(i, j); err != nil {
			merr = multierr.Append(merr, fmtJobError(err))
		}
	}
	for _, task := range j.Tasks {
		if err := validateTask(task, store); err != nil {
			merr = multierr.Append(merr, fmtJobError(err))
		}
	}
	return merr
}

func fmtJobError(err error) error {
	return fmt.Errorf("job validation: %v", err)
}

// ValidateInitiator checks the Initiator for any application logic errors.
func ValidateInitiator(i models.Initiator, j models.JobSpec) error {
	switch strings.ToLower(i.Type) {
	case models.InitiatorRunAt:
		return validateRunAtInitiator(i, j)
	case models.InitiatorCron:
		return validateCronInitiator(i)
	default:
		return fmtInitiatorError(fmt.Errorf("Initiator %v does not exist", i.Type))
	case models.InitiatorWeb:
		fallthrough
	case models.InitiatorRunLog:
		fallthrough
	case models.InitiatorEthLog:
		return nil
	}
}

func validateRunAtInitiator(i models.Initiator, j models.JobSpec) error {
	if i.Time.Unix() <= 0 {
		return fmtInitiatorError(errors.New(`runat must have a time`))
	} else if j.StartAt.Valid && i.Time.Unix() < j.StartAt.Time.Unix() {
		return fmtInitiatorError(errors.New(`runat time must be after job's startat`))
	} else if j.EndAt.Valid && i.Time.Unix() > j.EndAt.Time.Unix() {
		return fmtInitiatorError(errors.New(`runat time must be before job's endat`))
	}
	return nil
}

func fmtInitiatorError(err error) error {
	return fmt.Errorf("initiator validation: %v", err)
}

func validateCronInitiator(i models.Initiator) error {
	if i.Schedule == "" {
		return errors.New(`initiator validation: cron must have a schedule`)
	}
	return nil
}

func validateTask(task models.TaskSpec, store *store.Store) error {
	if _, err := adapters.For(task, store); err != nil {
		return fmt.Errorf("task validation: %v", err)
	}
	return nil
}
