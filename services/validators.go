package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

func ValidateJob(j models.Job, store *store.Store) error {
	for _, i := range j.Initiators {
		if err := validateInitiator(i); err != nil {
			return fmt.Errorf("job validation: %v", err)
		}
	}
	if err := adapters.Validate(j, store); err != nil {
		return fmt.Errorf("job validation: adapter validation: %v", err)
	}
	return nil
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
