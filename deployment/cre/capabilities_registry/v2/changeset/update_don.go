package changeset

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[UpdateDONInput] = UpdateDON{}

type UpdateDONInput struct {
	RegistryChainSel  uint64 `json:"registry_chain_sel" yaml:"registry_chain_sel"`
	RegistryQualifier string `json:"registry_qualifier" yaml:"registry_qualifier"`
	UseMCMS           bool   `json:"use_mcms" yaml:"use_mcms"` // not implemented yet

	DonID             int                          `json:"don_id" yaml:"don_id"`
	F                 uint8                        `json:"f" yaml:"f"`
	P2PIDs            []p2pkey.PeerID              `json:"p2p_ids" yaml:"p2p_ids"`
	CapabilityConfigs []contracts.CapabilityConfig `json:"capability_configs" yaml:"capability_configs"`
	IsPrivate         bool                         `json:"is_private" yaml:"is_private"`
}

type UpdateDON struct{}

func (u UpdateDON) VerifyPreconditions(_ cldf.Environment, config UpdateDONInput) error {
	if len(config.P2PIDs) == 0 {
		return errors.New("p2pIDs is required")
	}
	if len(config.CapabilityConfigs) == 0 {
		return errors.New("capabilityConfigs is required")
	}
	return nil
}

func (u UpdateDON) Apply(e cldf.Environment, config UpdateDONInput) (cldf.ChangesetOutput, error) {
	registryRef := datastore.NewAddressRefKey(
		config.RegistryChainSel,
		"CapabilitiesRegistry",
		semver.MustParse("2.0.0"),
		config.RegistryQualifier,
	)

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.UpdateDON,
		sequences.UpdateDONDeps{Env: &e},
		sequences.UpdateDONInput{
			RegistryChainSel:  config.RegistryChainSel,
			RegistryRef:       registryRef,
			DonID:             config.DonID,
			F:                 config.F,
			P2PIDs:            config.P2PIDs,
			CapabilityConfigs: config.CapabilityConfigs,
			IsPrivate:         config.IsPrivate,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports: seqReport.ExecutionReports,
	}, nil
}
