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

func ValidateJob(j models.Job, store *store.Store) error {
	var merr error
	for _, i := range j.Initiators {
		if err := ValidateInitiator(i); err != nil {
			merr = multierr.Append(merr, fmt.Errorf("job validation: %v", err))
		}
	}
	for _, task := range j.Tasks {
		if err := validateTask(task, store); err != nil {
			merr = multierr.Append(merr, fmt.Errorf("job validation: %v", err))
		}
	}
	return merr
}

func ValidateInitiator(i models.Initiator) error {
	switch strings.ToLower(i.Type) {
	case models.InitiatorRunAt:
		return validateRunAtInitiator(i)
	case models.InitiatorCron:
		return validateCronInitiator(i)
	default:
		return fmt.Errorf("Initiator %v does not exist", i.Type)
	case models.InitiatorWeb:
		fallthrough
	case models.InitiatorRunLog:
		fallthrough
	case models.InitiatorEthLog:
		return nil
	}
}

func validateRunAtInitiator(i models.Initiator) error {
	if i.Time.Unix() <= 0 {
		return errors.New(`initiator validation: runat must have time`)
	}
	return nil
}

func validateCronInitiator(i models.Initiator) error {
	if i.Schedule == "" {
		return errors.New(`initiator validation: cron must have schedule`)
	}
	return nil
}

func validateTask(task models.Task, store *store.Store) error {
	if _, err := adapters.For(task, store); err != nil {
		return fmt.Errorf("task validation: %v", err)
	}
	return nil
}
