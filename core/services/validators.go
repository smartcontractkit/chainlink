package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/multiformats/go-multiaddr"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
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
		if err := ValidateInitiator(i, j, store); err != nil {
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
	bridge, err := store.ORM.FindBridge(bt.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		fe.Add(fmt.Sprintf("Error determining if bridge type %v already exists", bt.Name))
	} else if (bridge != models.BridgeType{}) {
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
	ts := models.TaskSpec{Type: bt.Name}
	if a := adapters.FindNativeAdapterFor(ts); a != nil {
		fe.Add(fmt.Sprintf("Bridge Type %v is a native adapter", bt.Name))
	}
	return fe.CoerceEmptyToNil()
}

var (
	externalInitiatorNameRegexp = regexp.MustCompile("^[a-zA-Z0-9-_]+$")
)

// ValidateExternalInitiator checks whether External Initiator parameters are
// safe for processing.
func ValidateExternalInitiator(
	exi *models.ExternalInitiatorRequest,
	store *store.Store,
) error {
	fe := models.NewJSONAPIErrors()
	if len([]rune(exi.Name)) == 0 {
		fe.Add("No name specified")
	} else if !externalInitiatorNameRegexp.MatchString(exi.Name) {
		fe.Add("Name must be alphanumeric and may contain '_' or '-'")
	} else if _, err := store.FindExternalInitiatorByName(exi.Name); err == nil {
		fe.Add(fmt.Sprintf("Name %v already exists", exi.Name))
	} else if err != orm.ErrorNotFound {
		return errors.Wrap(err, "validating external initiator")
	}
	return fe.CoerceEmptyToNil()
}

// ValidateInitiator checks the Initiator for any application logic errors.
func ValidateInitiator(i models.Initiator, j models.JobSpec, store *store.Store) error {
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
		return validateRunLogInitiator(i, j, store)
	case models.InitiatorFluxMonitor:
		return validateFluxMonitor(i, j, store)
	case models.InitiatorWeb:
		return nil
	case models.InitiatorEthLog:
		return nil
	case models.InitiatorRandomnessLog:
		return validateRandomnessLogInitiator(i, j)
	default:
		return models.NewJSONAPIErrorsWith(fmt.Sprintf("type %v does not exist", i.Type))
	}
}

func validateFluxMonitor(i models.Initiator, j models.JobSpec, store *store.Store) error {
	fe := models.NewJSONAPIErrors()

	if i.Address == utils.ZeroAddress {
		fe.Add("no address")
	}
	if i.RequestData.String() == "" {
		fe.Add("no requestdata")
	}
	if i.Threshold <= 0 {
		fe.Add("bad 'threshold' parameter; this is the maximum relative change " +
			"allowed in the monitored value, before a new report should be made; " +
			"it must be positive, and appear in the job initiator parameters; e.g." +
			`{"initiators": [{"type":"fluxmonitor", "params":{"threshold": 0.5}}]} ` +
			"means that the value can change by up to half its last-reported value " +
			"before a new report is made")
	}
	if i.AbsoluteThreshold < 0 {
		fe.Add("bad 'absoluteThreshold' value; this is the maximum absolute " +
			"change allowed in the monitored value, before a new report should be " +
			"made; it must be nonnegative and appear in the job initiator parameters; e.g." +
			`{"initiators":[{"type":"fluxmonitor","params":{"absoluteThreshold":0.01}}]} ` +
			"means that the value can change by up to 0.01 units " +
			"before a new report is made")
	}

	if i.PollTimer.Disabled && i.IdleTimer.Disabled {
		fe.Add("must enable pollTimer, idleTimer, or both")
	}

	if i.PollTimer.Disabled {
		if !i.PollTimer.Period.IsInstant() {
			fe.Add("pollTimer disabled, period must be 0")
		}
	} else {
		minimumPollPeriod := models.Duration(store.Config.DefaultHTTPTimeout())

		if i.PollTimer.Period.IsInstant() {
			fe.Add("pollTimer enabled, but no period specified")
		} else if i.PollTimer.Period.Shorter(minimumPollPeriod) {
			fe.Add("pollTimer enabled, period must be equal or greater than " + minimumPollPeriod.String())
		}
	}

	if i.IdleTimer.Disabled {
		if !i.IdleTimer.Duration.IsInstant() {
			fe.Add("idleTimer disabled, duration must be 0")
		}
	} else {
		if i.IdleTimer.Duration.IsInstant() {
			fe.Add("idleTimer enabled, duration must be > 0")
		} else if !i.PollTimer.Disabled && i.IdleTimer.Duration.Shorter(i.PollTimer.Period) {
			fe.Add("idleTimer and pollTimer enabled, idleTimer.duration must be >= than pollTimer.period")
		}
	}

	if err := validateFeeds(i.Feeds, store); err != nil {
		fe.Add(err.Error())
	}

	return fe.CoerceEmptyToNil()
}

func validateFeeds(feeds models.Feeds, store *store.Store) error {
	var feedsData []interface{}
	if err := json.Unmarshal(feeds.Bytes(), &feedsData); err != nil {
		return errors.New("invalid json for feeds parameter")
	}
	if len(feedsData) == 0 {
		return errors.New("feeds field is empty")
	}

	var bridgeNames []string
	for _, entry := range feedsData {
		switch feed := entry.(type) {
		case string:
			if _, err := url.ParseRequestURI(feed); err != nil {
				return err
			}
		case map[string]interface{}: // named feed - ex: {"bridge": "bridgeName"}
			bridgeName := feed["bridge"]
			bridgeNameString, ok := bridgeName.(string)
			if bridgeName == nil {
				return errors.New("Feeds object missing bridge key")
			} else if len(feed) != 1 {
				return errors.New("Unsupported keys in feed JSON")
			} else if !ok {
				return errors.New("Unsupported bridge name type in feed JSON")
			}
			bridgeNames = append(bridgeNames, bridgeNameString)
		default:
			return errors.New("Unknown feed type")
		}
	}
	if _, err := store.ORM.FindBridgesByNames(bridgeNames); err != nil {
		return err
	}

	return nil
}

func validateRunLogInitiator(i models.Initiator, j models.JobSpec, s *store.Store) error {
	fe := models.NewJSONAPIErrors()
	ethTxCount := 0
	for _, task := range j.Tasks {
		if task.Type == adapters.TaskTypeEthTx {
			ethTxCount++

			task.Params.ForEach(func(k, v gjson.Result) bool {
				key := strings.ToLower(k.String())
				if key == "functionselector" {
					fe.Add("Cannot set EthTx Task's function selector parameter with a RunLog Initiator")
				} else if key == "address" {
					fe.Add("Cannot set EthTx Task's address parameter with a RunLog Initiator")
				} else if key == "fromaddress" {
					if !common.IsHexAddress(v.String()) {
						fe.Add("Cannot set EthTx Task's fromAddress parameter: invalid address")
						return true
					}
					address := common.HexToAddress(v.String())
					exists, err := s.KeyExists(address)
					if err != nil {
						fe.Add("Cannot set EthTx Task's fromAddress parameter: " + err.Error())
					} else if !exists {
						fe.Add("Cannot set EthTx Task's fromAddress parameter: the node does not have this private key in the database")
					}
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

func validateRandomnessLogInitiator(i models.Initiator, j models.JobSpec) error {
	fe := models.NewJSONAPIErrors()
	if len(j.Initiators) != 1 {
		fe.Add("randomness log must have exactly one initiator")
	}
	if i.Address == utils.ZeroAddress {
		fe.Add("randomness log must specify address of expected emmitter")
	}
	return fe.CoerceEmptyToNil()
}

func validateTask(task models.TaskSpec, store *store.Store) error {
	adapter, err := adapters.For(task, store.Config, store.ORM)
	if err != nil {
		return err
	}
	if !store.Config.EnableExperimentalAdapters() {
		if _, ok := adapter.BaseAdapter.(*adapters.Sleep); ok {
			return errors.New("Sleep Adapter is not implemented yet")
		}
	}
	return nil
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

	found := false
	for _, oracle := range sa.Encumbrance.Oracles {
		if store.KeyStore.HasAccountWithAddress(oracle.Address()) {
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

	if untilEndAt > config.MaximumServiceDuration().Duration() {
		fe.Add(fmt.Sprintf("Service agreement encumbrance error: endAt value of %v is too far in the future. Furthest allowed date is %v",
			sa.Encumbrance.EndAt,
			time.Now().Add(config.MaximumServiceDuration().Duration())))
	}

	if untilEndAt < config.MinimumServiceDuration().Duration() {
		fe.Add(fmt.Sprintf("Service agreement encumbrance error: endAt value of %v is too soon. Earliest allowed date is %v",
			sa.Encumbrance.EndAt,
			time.Now().Add(config.MinimumServiceDuration().Duration())))
	}

	return fe.CoerceEmptyToNil()
}

// ValidatedOracleSpecToml validates an oracle spec that came from TOML
func ValidatedOracleSpecToml(tomlString string) (spec offchainreporting.OracleSpec, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panicked with err %v", r)
		}
	}()

	var oros models.OffchainReportingOracleSpec
	spec = offchainreporting.OracleSpec{
		Pipeline: *pipeline.NewTaskDAG(),
	}
	tree, err := toml.Load(tomlString)
	if err != nil {
		return spec, err
	}
	err = tree.Unmarshal(&oros)
	if err != nil {
		return spec, err
	}
	err = tree.Unmarshal(&spec)
	if err != nil {
		return spec, err
	}
	spec.OffchainReportingOracleSpec = oros

	// TODO(#175801426): upstream a way to check for undecoded keys in go-toml
	// TODO(#175801038): upstream support for time.Duration defaults in go-toml
	var defaults = map[string]interface{}{
		"observationTimeout":                     models.Interval(10 * time.Second),
		"blockchainTimeout":                      models.Interval(20 * time.Second),
		"contractConfigTrackerSubscribeInterval": models.Interval(2 * time.Minute),
		"contractConfigTrackerPollInterval":      models.Interval(1 * time.Minute),
	}
	for k, v := range defaults {
		cv := v
		if !tree.Has(k) {
			reflect.ValueOf(&spec.OffchainReportingOracleSpec).
				Elem().
				FieldByName(strings.Title(k)).
				Set(reflect.ValueOf(cv))
		}
	}

	if spec.Type != "offchainreporting" {
		return spec, errors.Errorf("the only supported type is currently 'offchainreporting', got %s", spec.Type)
	}
	if spec.SchemaVersion != uint32(1) {
		return spec, errors.Errorf("the only supported schema version is currently 1, got %v", spec.SchemaVersion)
	}
	if !tree.Has("isBootstrapPeer") {
		return spec, errors.New("isBootstrapPeer is not defined")
	}
	if spec.IsBootstrapPeer {
		if err := validateBootstrapSpec(tree, spec); err != nil {
			return spec, err
		}
	} else if err := validateNonBootstrapSpec(tree, spec); err != nil {
		return spec, err
	}
	if err := validateTimingParameters(spec); err != nil {
		return spec, err
	}
	if err := validateMonitoringURL(spec); err != nil {
		return spec, err
	}
	return spec, nil
}

// Parameters that must be explicitly set by the operator.
var (
	// Common to both bootstrap and non-boostrap
	params = map[string]struct{}{
		"type":              {},
		"schemaVersion":     {},
		"contractAddress":   {},
		"isBootstrapPeer":   {},
		"p2pPeerID":         {},
		"p2pBootstrapPeers": {},
	}
	// Boostrap and non-bootstrap parameters
	// are mutually exclusive.
	bootstrapParams    = map[string]struct{}{}
	nonBootstrapParams = map[string]struct{}{
		"observationSource":  {},
		"observationTimeout": {},
		"keyBundleID":        {},
		"transmitterAddress": {},
	}
)

func cloneSet(in map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for k, v := range in {
		out[k] = v
	}
	return out
}

func validateTimingParameters(spec offchainreporting.OracleSpec) error {
	// TODO: expose these various constants from libocr so they are defined in one place.
	if time.Duration(spec.ObservationTimeout) < 1*time.Millisecond || time.Duration(spec.ObservationTimeout) > 20*time.Second {
		return errors.Errorf("require 1ms <= observation timeout <= 20s")
	}
	if time.Duration(spec.BlockchainTimeout) < 1*time.Millisecond || time.Duration(spec.ObservationTimeout) > 20*time.Second {
		return errors.Errorf("require 1ms <= blockchain timeout <= 20s ")
	}
	if time.Duration(spec.ContractConfigTrackerPollInterval) < 15*time.Second || time.Duration(spec.ContractConfigTrackerPollInterval) > 120*time.Second {
		return errors.Errorf("require 15s <= contract config tracker poll interval <= 120s ")
	}
	if time.Duration(spec.ContractConfigTrackerSubscribeInterval) < 2*time.Minute || time.Duration(spec.ContractConfigTrackerSubscribeInterval) > 5*time.Minute {
		return errors.Errorf("require 2m <= contract config subscribe interval <= 5m ")
	}
	if spec.ContractConfigConfirmations < 2 || spec.ContractConfigConfirmations > 10 {
		return errors.Errorf("require 2 <= contract config confirmations <= 10 ")
	}
	return nil
}

func validateBootstrapSpec(tree *toml.Tree, spec offchainreporting.OracleSpec) error {
	expected, notExpected := cloneSet(params), cloneSet(nonBootstrapParams)
	for k := range bootstrapParams {
		expected[k] = struct{}{}
	}
	if err := validateExplicitlySetKeys(tree, expected, notExpected, "bootstrap"); err != nil {
		return err
	}
	for i := range spec.P2PBootstrapPeers {
		if _, err := multiaddr.NewMultiaddr(spec.P2PBootstrapPeers[i]); err != nil {
			return errors.Errorf("p2p bootstrap peer %d is invalid: err %v", i, err)
		}
	}
	return nil
}

func validateNonBootstrapSpec(tree *toml.Tree, spec offchainreporting.OracleSpec) error {
	expected, notExpected := cloneSet(params), cloneSet(bootstrapParams)
	for k := range nonBootstrapParams {
		expected[k] = struct{}{}
	}
	if err := validateExplicitlySetKeys(tree, expected, notExpected, "non-bootstrap"); err != nil {
		return err
	}
	if spec.Pipeline.DOTSource == "" {
		return errors.New("no pipeline specified")
	}
	if len(spec.P2PBootstrapPeers) < 1 {
		return errors.New("must specify at least one bootstrap peer")
	}
	for i := range spec.P2PBootstrapPeers {
		if _, err := multiaddr.NewMultiaddr(spec.P2PBootstrapPeers[i]); err != nil {
			return errors.Errorf("p2p bootstrap peer %d is invalid: err %v", i, err)
		}
	}
	if spec.MaxTaskDuration > spec.ObservationTimeout {
		return errors.Errorf("max task duration must be < observation timeout")
	}
	tasks, err := spec.Pipeline.TasksInDependencyOrder()
	if err != nil {
		return errors.Wrap(err, "invalid observation source")
	}
	for _, task := range tasks {
		timeout, set := task.TaskTimeout()
		if set && timeout > time.Duration(spec.ObservationTimeout) {
			return errors.Errorf("individual max task duration must be < observation timeout")
		}
	}
	return nil
}

func validateExplicitlySetKeys(tree *toml.Tree, expected map[string]struct{}, notExpected map[string]struct{}, peerType string) error {
	var err error
	// top level keys only
	for _, k := range tree.Keys() {
		// TODO(#175801577): upstream a way to check for children in go-toml
		if _, ok := notExpected[k]; ok {
			err = multierr.Append(err, errors.Errorf("unrecognised key for %s peer: %s", peerType, k))
		}
		delete(expected, k)
	}
	for missing := range expected {
		err = multierr.Append(err, errors.Errorf("missing required key %s", missing))
	}
	return err
}

func validateMonitoringURL(spec offchainreporting.OracleSpec) error {
	if spec.MonitoringEndpoint == "" {
		return nil
	}
	_, err := url.Parse(spec.MonitoringEndpoint)
	return err
}
