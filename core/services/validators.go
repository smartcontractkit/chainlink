package services

import (
	"fmt"
	"strings"
	"time"

	"chainlink/core/adapters"
	"chainlink/core/assets"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
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

// ValidateBridgeTypeNotExist checks that a bridge has not already been created
func ValidateBridgeTypeNotExist(bt *models.BridgeTypeRequest, store *store.Store) error {
	fe := models.NewJSONAPIErrors()
	ts := models.TaskSpec{Type: bt.Name}
	if a, _ := adapters.For(ts, store.Config, store.ORM); a != nil {
		fe.Add(fmt.Sprintf("Bridge Type %v already exists", bt.Name))
	}
	return fe.CoerceEmptyToNil()
}

// ValidateBridgeType checks that the bridge type doesn't have a duplicate
// or invalid name or invalid url
func ValidateBridgeType(bt *models.BridgeTypeRequest, store *store.Store) error {
	fe := models.NewJSONAPIErrors()
	if len(bt.Name.String()) < 1 {
		fe.Add("No name specified")
	}
	if _, err := models.NewTaskType(bt.Name.String()); err != nil {
		fe.Merge(err)
	}
	u := bt.URL.String()
	if len(strings.TrimSpace(u)) == 0 {
		fe.Add("URL must be present")
	}
	if bt.MinimumContractPayment != nil &&
		bt.MinimumContractPayment.Cmp(assets.NewLink(0)) < 0 {
		fe.Add("MinimumContractPayment must be positive")
	}
	return fe.CoerceEmptyToNil()
}

// ValidateExternalInitiator checks whether External Initiator parameters are
// safe for processing.
func ValidateExternalInitiator(
	exi *models.ExternalInitiatorRequest,
	store *store.Store,
) error {
	fe := models.NewJSONAPIErrors()
	if len([]rune(exi.Name)) == 0 {
		fe.Add("No name specified")
	} else if onlyValidRunes := govalidator.StringMatches(exi.Name, "^[a-zA-Z0-9-_]*$"); !onlyValidRunes {
		fe.Add("Name must be alphanumeric and may contain '_' or '-'")
	} else if _, err := store.FindExternalInitiatorByName(exi.Name); err == nil {
		fe.Add(fmt.Sprintf("Name %v already exists", exi.Name))
	} else if err != orm.ErrorNotFound {
		return errors.Wrap(err, "validating external initiator")
	}
	// only validate URL if present
	if exi.URL != nil {
		if isURL := govalidator.IsURL((*exi.URL).String()); !isURL {
			fe.Add("Invalid URL format")
		}
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
	case models.InitiatorExternal:
		return validateExternalInitiator(i)
	case models.InitiatorServiceAgreementExecutionLog:
		return validateServiceAgreementInitiator(i, j)
	case models.InitiatorRunLog:
		return validateRunLogInitiator(i, j)
	case models.InitiatorFluxMonitor:
		return validateFluxMonitor(i, j)
	case models.InitiatorWeb:
		return nil
	case models.InitiatorEthLog:
		return nil
	default:
		return models.NewJSONAPIErrorsWith(fmt.Sprintf("type %v does not exist", i.Type))
	}
}

func validateFluxMonitor(i models.Initiator, j models.JobSpec) error {
	fe := models.NewJSONAPIErrors()
	if i.Address == utils.ZeroAddress {
		fe.Add("no address")
	}
	if len(i.Feeds) == 0 {
		fe.Add("no feeds")
	}
	if i.Threshold <= 0 {
		fe.Add("bad threshold")
	}
	if i.RequestData.String() == "" {
		fe.Add("no requestdata")
	}
	if i.PollingInterval == 0 {
		fe.Add("no pollingInterval")
	} else if i.PollingInterval < MinimumPollingInterval {
		fe.Add("pollingInterval must be equal or greater than " + MinimumPollingInterval.String())
	}

	return fe.CoerceEmptyToNil()
}

func validateRunLogInitiator(i models.Initiator, j models.JobSpec) error {
	fe := models.NewJSONAPIErrors()
	ethTxCount := 0
	for _, task := range j.Tasks {
		if task.Type == adapters.TaskTypeEthTx {
			ethTxCount += 1

			task.Params.ForEach(func(k, _ gjson.Result) bool {
				key := strings.ToLower(k.String())
				if key == "functionselector" {
					fe.Add("Cannot set EthTx Task's function selector parameter with a RunLog Initiator")
				} else if key == "address" {
					fe.Add("Cannot set EthTx Task's address parameter with a RunLog Initiator")
				}
				return true
			})
		}
	}
	if ethTxCount > 1 {
		fe.Add("Cannot RunLog initiated jobs cannot have more than one EthTx Task")
	}
	return fe.CoerceEmptyToNil()
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

func validateExternalInitiator(i models.Initiator) error {
	if len([]rune(i.Name)) == 0 {
		return models.NewJSONAPIErrorsWith("External must have a name")
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
	adapter, err := adapters.For(task, store.Config, store.ORM)
	if !store.Config.Dev() {
		if _, ok := adapter.BaseAdapter.(*adapters.Sleep); ok {
			return errors.New("Sleep Adapter is not implemented yet")
		}
		if _, ok := adapter.BaseAdapter.(*adapters.EthTxABIEncode); ok {
			return errors.New("EthTxABIEncode Adapter is not implemented yet")
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
