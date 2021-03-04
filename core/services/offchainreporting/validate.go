package offchainreporting

import (
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/libocr/offchainreporting"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"go.uber.org/multierr"
)

// ValidatedOracleSpecToml validates an oracle spec that came from TOML
func ValidatedOracleSpecToml(config *orm.Config, tomlString string) (job.SpecDB, error) {
	var specDB = job.SpecDB{
		Pipeline: *pipeline.NewTaskDAG(),
	}
	var spec job.OffchainReportingOracleSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return specDB, errors.Wrap(err, "toml error on load")
	}
	// Note this validates all the fields which implement an UnmarshalText
	// i.e. TransmitterAddress, PeerID...
	err = tree.Unmarshal(&spec)
	if err != nil {
		return specDB, errors.Wrap(err, "toml unmarshal error on spec")
	}
	err = tree.Unmarshal(&specDB)
	if err != nil {
		return specDB, errors.Wrap(err, "toml unmarshal error on specDB")
	}
	specDB.OffchainreportingOracleSpec = &spec

	// TODO(#175801426): upstream a way to check for undecoded keys in go-toml
	// TODO(#175801038): upstream support for time.Duration defaults in go-toml
	if specDB.Type != job.OffchainReporting {
		return specDB, errors.Errorf("the only supported type is currently 'offchainreporting', got %s", specDB.Type)
	}
	if specDB.SchemaVersion != uint32(1) {
		return specDB, errors.Errorf("the only supported schema version is currently 1, got %v", specDB.SchemaVersion)
	}
	if !tree.Has("isBootstrapPeer") {
		return specDB, errors.New("isBootstrapPeer is not defined")
	}
	for i := range spec.P2PBootstrapPeers {
		if _, err := multiaddr.NewMultiaddr(spec.P2PBootstrapPeers[i]); err != nil {
			return specDB, errors.Wrapf(err, "p2p bootstrap peer %v is invalid", spec.P2PBootstrapPeers[i])
		}
	}
	if spec.IsBootstrapPeer {
		if err := validateBootstrapSpec(tree, specDB); err != nil {
			return specDB, err
		}
	} else if err := validateNonBootstrapSpec(tree, config, specDB); err != nil {
		return specDB, err
	}
	if err := validateTimingParameters(config, spec); err != nil {
		return specDB, err
	}
	return specDB, nil
}

// Parameters that must be explicitly set by the operator.
var (
	// Common to both bootstrap and non-boostrap
	params = map[string]struct{}{
		"type":            {},
		"schemaVersion":   {},
		"contractAddress": {},
		"isBootstrapPeer": {},
	}
	// Boostrap and non-bootstrap parameters
	// are mutually exclusive.
	bootstrapParams    = map[string]struct{}{}
	nonBootstrapParams = map[string]struct{}{
		"observationSource": {},
	}
)

func cloneSet(in map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for k, v := range in {
		out[k] = v
	}
	return out
}

func validateTimingParameters(config *orm.Config, spec job.OffchainReportingOracleSpec) error {
	lc := types.LocalConfig{
		BlockchainTimeout:                      config.OCRBlockchainTimeout(time.Duration(spec.BlockchainTimeout)),
		ContractConfigConfirmations:            config.OCRContractConfirmations(spec.ContractConfigConfirmations),
		ContractConfigTrackerPollInterval:      config.OCRContractPollInterval(time.Duration(spec.ContractConfigTrackerPollInterval)),
		ContractConfigTrackerSubscribeInterval: config.OCRContractSubscribeInterval(time.Duration(spec.ContractConfigTrackerSubscribeInterval)),
		ContractTransmitterTransmitTimeout:     config.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        config.OCRDatabaseTimeout(),
		DataSourceTimeout:                      config.OCRObservationTimeout(time.Duration(spec.ObservationTimeout)),
		DataSourceGracePeriod:                  config.OCRObservationGracePeriod(),
	}
	if config.Dev() {
		lc.DevelopmentMode = types.EnableDangerousDevelopmentMode
	}
	return offchainreporting.SanityCheckLocalConfig(lc)
}

func validateBootstrapSpec(tree *toml.Tree, spec job.SpecDB) error {
	expected, notExpected := cloneSet(params), cloneSet(nonBootstrapParams)
	for k := range bootstrapParams {
		expected[k] = struct{}{}
	}
	if err := validateExplicitlySetKeys(tree, expected, notExpected, "bootstrap"); err != nil {
		return err
	}
	return nil
}

func validateNonBootstrapSpec(tree *toml.Tree, config *orm.Config, spec job.SpecDB) error {
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
	observationTimeout := config.OCRObservationTimeout(time.Duration(spec.OffchainreportingOracleSpec.ObservationTimeout))
	if time.Duration(spec.MaxTaskDuration) > observationTimeout {
		return errors.Errorf("max task duration must be < observation timeout")
	}
	tasks, err := spec.Pipeline.TasksInDependencyOrder()
	if err != nil {
		return errors.Wrap(err, "invalid observation source")
	}
	for _, task := range tasks {
		timeout, set := task.TaskTimeout()
		if set && timeout > observationTimeout {
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
