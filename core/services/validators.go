package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// ValidateJob checks the job and its associated Initiators and Tasks for any
// application logic errors.
func ValidateJob(j models.JobSpec, store *store.Store) error {
	fe := models.NewJSONAPIErrors()
	if j.StartAt.Valid && j.EndAt.Valid && j.StartAt.Time.After(j.EndAt.Time) {
		fe.Add("StartAt cannot be before EndAt")
	}
	if len(j.Initiators) < 1 || len(j.Tasks) < 1 {
		fe.Add("Must have at least one Initiator and one Task")
	}
	for _, i := range j.Initiators {
		if err := ValidateInitiator(i, j); err != nil {
			fe.Merge(err)
		}
	}
	for _, task := range j.Tasks {
		if err := validateTask(task, store); err != nil {
			fe.Merge(err)
		}
	}
	return fe.CoerceEmptyToNil()
}

// ValidateAdapter checks that the bridge type doesn't have a duplicate or invalid name
func ValidateAdapter(bt *models.BridgeTypeRequest, store *store.Store) error {
	fe := models.NewJSONAPIErrors()
	if len(bt.Name.String()) < 1 {
		fe.Add("No name specified")
	}
	if _, err := models.NewTaskType(bt.Name.String()); err != nil {
		fe.Merge(err)
	}
	ts := models.TaskSpec{Type: bt.Name}
	if a, _ := adapters.For(ts, store); a != nil {
		fe.Add(fmt.Sprintf("Adapter %v already exists", bt.Name))
	}
	return fe.CoerceEmptyToNil()
}

// ValidateInitiator checks the Initiator for any application logic errors.
func ValidateInitiator(i models.Initiator, j models.JobSpec) error {
	switch strings.ToLower(i.Type) {
	case models.InitiatorRunAt:
		return validateRunAtInitiator(i, j)
	case models.InitiatorCron:
		return validateCronInitiator(i)
	case models.InitiatorServiceAgreementExecutionLog:
		return validateServiceAgreementInitiator(i, j)
	case models.InitiatorWeb:
		fallthrough
	case models.InitiatorRunLog:
		fallthrough
	case models.InitiatorEthLog:
		return nil
	default:
		return models.NewJSONAPIErrorsWith(fmt.Sprintf("type %v does not exist", i.Type))
	}
}

func validateRunAtInitiator(i models.Initiator, j models.JobSpec) error {
	fe := models.NewJSONAPIErrors()
	if !i.Time.Valid {
		fe.Add("RunAt must have a time")
	} else if j.StartAt.Valid && i.Time.Time.Unix() < j.StartAt.Time.Unix() {
		fe.Add("RunAt time must be after job's StartAt")
	} else if j.EndAt.Valid && i.Time.Time.Unix() > j.EndAt.Time.Unix() {
		fe.Add("RunAt time must be before job's EndAt")
	}
	return fe.CoerceEmptyToNil()
}

func validateCronInitiator(i models.Initiator) error {
	if i.Schedule == "" {
		return models.NewJSONAPIErrorsWith("Schedule must have a cron")
	}
	return nil
}

func validateServiceAgreementInitiator(i models.Initiator, j models.JobSpec) error {
	fe := models.NewJSONAPIErrors()
	if len(j.Initiators) != 1 {
		fe.Add("ServiceAgreement should have at most one initiator")
	}
	return fe.CoerceEmptyToNil()
}

func validateTask(task models.TaskSpec, store *store.Store) error {
	adapter, err := adapters.For(task, store)
	if !store.Config.Dev() {
		if _, ok := adapter.BaseAdapter.(*adapters.Sleep); ok {
			return errors.New("Sleep Adapter is not implemented yet")
		}
	}
	return err
}

// ValidateServiceAgreement checks the ServiceAgreement for any application logic errors.
func ValidateServiceAgreement(sa models.ServiceAgreement, store *store.Store) error {
	fe := models.NewJSONAPIErrors()
	config := store.Config

	if sa.Encumbrance.Payment == nil {
		fe.Add("Service agreement encumbrance error: No payment amount set")
	} else if sa.Encumbrance.Payment.Cmp(config.MinimumContractPayment()) == -1 {
		fe.Add(fmt.Sprintf("Service agreement encumbrance error: Payment amount is below minimum %v", config.MinimumContractPayment().String()))
	}

	if sa.Encumbrance.Expiration < config.MinimumRequestExpiration() {
		fe.Add(fmt.Sprintf("Service agreement encumbrance error: Expiration is below minimum %v", config.MinimumRequestExpiration()))
	}

	account, err := store.KeyStore.GetFirstAccount()
	if err != nil {
		return err // 500
	}

	found := false
	for _, oracle := range sa.Encumbrance.Oracles {
		if oracle.Address() == account.Address {
			found = true
		}
	}
	if !found {
		fe.Add("Service agreement encumbrance error: This node must be listed in the participating oracles")
	}

	if err := ValidateJob(sa.JobSpec, store); err != nil {
		fe.Add(fmt.Sprintf("Service agreement job spec error: Job spec validation: %v", err))
	}

	untilEndAt := time.Until(sa.Encumbrance.EndAt.Time)

	if untilEndAt > config.MaximumServiceDuration() {
		fe.Add(fmt.Sprintf("Service agreement encumbrance error: endAt value of %v is too far in the future. Furthest allowed date is %v",
			sa.Encumbrance.EndAt,
			time.Now().Add(config.MaximumServiceDuration())))
	}

	if untilEndAt < config.MinimumServiceDuration() {
		fe.Add(fmt.Sprintf("Service agreement encumbrance error: endAt value of %v is too soon. Earliest allowed date is %v",
			sa.Encumbrance.EndAt,
			time.Now().Add(config.MinimumServiceDuration())))
	}

	return fe.CoerceEmptyToNil()
}
