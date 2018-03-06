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
		if err := validateInitiator(i); err != nil {
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

func validateInitiator(i models.Initiator) error {
	if strings.ToLower(i.Type) == models.InitiatorRunAt {
		return validateRunAtInitiator(i)
	}
	return nil
}

func validateRunAtInitiator(i models.Initiator) error {
	if i.Time.Unix() <= 0 {
		return errors.New(`initiator validation: runat must have time`)
	}
	return nil
}

func validateTask(task models.Task, store *store.Store) error {
	if _, err := adapters.For(task, store); err != nil {
		return fmt.Errorf("task validation: %v", err)
	}
	return nil
}
