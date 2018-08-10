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

// ValidationError is an error that occurs during validation.
type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string { return e.msg }

// NewValidationError returns a validation error.
func NewValidationError(msg string, values ...interface{}) error {
	return &ValidationError{msg: fmt.Sprintf(msg, values...)}
}

// ValidateJob checks the job and its associated Initiators and Tasks for any
// application logic errors.
func ValidateJob(j models.JobSpec, store *store.Store) error {
	var merr error
	if j.StartAt.Valid && j.EndAt.Valid && j.StartAt.Time.After(j.EndAt.Time) {
		merr = multierr.Append(merr, fmtJobError(errors.New("startat cannot be before endat")))
	}
	if len(j.Initiators) < 1 || len(j.Tasks) < 1 {
		merr = multierr.Append(merr, fmtJobError(errors.New("Must have at least one Initiator and one Task")))
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

// ValidateAdapter checks that the bridge type doesn't have a duplicate or invalid name
func ValidateAdapter(bt *models.BridgeType, store *store.Store) (err error) {
	var merr error
	if len(bt.Name.String()) < 1 {
		merr = multierr.Append(merr, fmt.Errorf("adapter validation: no name specified"))
	}
	if _, err := models.NewTaskType(bt.Name.String()); err != nil {
		merr = multierr.Append(merr, fmt.Errorf("adapter validation error: %v", err))
	}
	ts := models.TaskSpec{Type: bt.Name}
	if a, _ := adapters.For(ts, store); a != nil {
		merr = multierr.Append(merr, fmt.Errorf("adapter validation: adapter %v exists", bt.Name))
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

// ValidateServiceAgreement checks the ServiceAgreement for any application logic errors.
func ValidateServiceAgreement(sa models.ServiceAgreement, config store.Config) error {
	var merr error

	if sa.Encumbrance.Payment == nil {
		merr = multierr.Append(merr, NewValidationError("service agreement encumbrance error: No payment amount set"))
	} else if sa.Encumbrance.Payment.Cmp(&config.MinimumContractPayment) == -1 {
		merr = multierr.Append(merr, NewValidationError("service agreement encumbrance error: Payment amount is below minimum %v", config.MinimumContractPayment.String()))
	}

	if sa.Encumbrance.Expiration < config.MinimumRequestExpiration {
		merr = multierr.Append(merr, NewValidationError("service agreement encumbrance error: Expiration is below minimum %v", config.MinimumRequestExpiration))
	}

	return merr
}
