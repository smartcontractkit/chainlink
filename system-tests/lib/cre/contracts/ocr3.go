package contracts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
)

func DeployOCR3Contract(logger zerolog.Logger, qualifier string, selector uint64, env *cldf.Environment, contractVersions map[string]string) (*ks_contracts_op.DeployOCR3ContractSequenceOutput, *common.Address, error) {
	memoryDatastore := datastore.NewMemoryDataStore()

	// load all existing addresses into memory datastore
	mergeErr := memoryDatastore.Merge(env.DataStore)
	if mergeErr != nil {
		return nil, nil, fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
	}

	ocr3DeployReport, err := operations.ExecuteSequence(
		env.OperationsBundle,
		ks_contracts_op.DeployOCR3ContractsSequence,
		ks_contracts_op.DeployOCR3ContractSequenceDeps{
			Env: env,
		},
		ks_contracts_op.DeployOCR3ContractSequenceInput{
			ChainSelector: selector,
			Qualifier:     qualifier,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy OCR3 contract '%s' on chain %d: %w", qualifier, selector, err)
	}
	// TODO: CRE-742 remove address book
	if err = env.ExistingAddresses.Merge(ocr3DeployReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
		return nil, nil, fmt.Errorf("failed to merge address book with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
	}
	if err = memoryDatastore.Merge(ocr3DeployReport.Output.Datastore); err != nil {
		return nil, nil, fmt.Errorf("failed to merge datastore with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
	}

	address := MustGetAddressFromMemoryDataStore(memoryDatastore, selector, keystone_changeset.OCR3Capability.String(), contractVersions[keystone_changeset.OCR3Capability.String()], qualifier)
	logger.Info().Msgf("Deployed OCR3 %s contract on chain %d at %s [qualifier: %s]", contractVersions[keystone_changeset.OCR3Capability.String()], selector, address, qualifier)

	env.DataStore = memoryDatastore.Seal()

	return &ocr3DeployReport.Output, &address, nil
}

// values supplied by Alexandr Yepishev as the expected values for OCR3 config
func DefaultOCR3Config() (*keystone_changeset.OracleConfig, error) {
	// values supplied by Alexandr Yepishev as the expected values for OCR3 config
	oracleConfig := &keystone_changeset.OracleConfig{
		DeltaProgressMillis:               5000,
		DeltaResendMillis:                 5000,
		DeltaInitialMillis:                5000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 1000,
		DeltaStageMillis:                  30000,
		MaxRoundsPerEpoch:                 10,
		MaxDurationQueryMillis:            1000,
		MaxDurationObservationMillis:      1000,
		MaxDurationShouldAcceptMillis:     1000,
		MaxDurationShouldTransmitMillis:   1000,
		MaxFaultyOracles:                  1,
		ConsensusCapOffchainConfig: &ocr3.ConsensusCapOffchainConfig{
			MaxQueryLengthBytes:       1000000,
			MaxObservationLengthBytes: 1000000,
			MaxOutcomeLengthBytes:     1000000,
			MaxReportLengthBytes:      1000000,
			MaxBatchSize:              1000,
		},
		UniqueReports: true,
	}

	return oracleConfig, nil
}

func DefaultOCR3_1Config(numWorkers int) (*ocr3.V3_1OracleConfig, error) {
	return &ocr3.V3_1OracleConfig{
		DeltaProgressMillis:  5000, // DKG 10-15 seconds; Vault 5 sec // check bandwidth from nops
		DeltaRoundMillis:     200,
		DeltaGraceMillis:     0,
		DeltaStageMillis:     0,
		MaxRoundsPerEpoch:    10,
		TransmissionSchedule: []int{numWorkers},

		MaxDurationInitializationMillis:               10000,
		MaxDurationShouldAcceptAttestedReportMillis:   1000,
		MaxDurationShouldTransmitAcceptedReportMillis: 1000,

		WarnDurationQueryMillis:               1000,
		WarnDurationObservationMillis:         1000,
		WarnDurationValidateObservationMillis: 1000,
		WarnDurationObservationQuorumMillis:   1000,
		WarnDurationStateTransition:           1000,
		WarnDurationCommitted:                 1000,

		MaxFaultyOracles: 1,
	}, nil
}

func DefaultChainCapabilityOCR3Config() (*keystone_changeset.OracleConfig, error) {
	cfg, err := DefaultOCR3Config()
	if err != nil {
		return nil, fmt.Errorf("failed to generate default OCR3 config: %w", err)
	}

	cfg.DeltaRoundMillis = 1000
	const kib = 1024
	const mib = 1024 * kib
	cfg.ConsensusCapOffchainConfig = nil
	cfg.ChainCapOffchainConfig = &ocr3.ChainCapOffchainConfig{
		MaxQueryLengthBytes:       mib,
		MaxObservationLengthBytes: 97 * kib,
		MaxReportLengthBytes:      mib,
		MaxOutcomeLengthBytes:     mib,
		MaxReportCount:            1000,
		MaxBatchSize:              200,
	}
	return cfg, nil
}
