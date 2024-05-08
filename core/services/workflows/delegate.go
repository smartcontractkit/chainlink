package workflows

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const hardcodedWorkflow = `
triggers:
  - type: "mercury-trigger"
    config:
      feedIds:
        - "0x1111111111111111111100000000000000000000000000000000000000000000"
        - "0x2222222222222222222200000000000000000000000000000000000000000000"
        - "0x3333333333333333333300000000000000000000000000000000000000000000"

consensus:
  - type: "offchain_reporting"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      aggregation_method: "data_feeds_2_0"
      aggregation_config:
        "0x1111111111111111111100000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x2222222222222222222200000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x3333333333333333333300000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
      encoder: "EVM"
      encoder_config:
        abi: "mercury_reports bytes[]"

targets:
  - type: "write_polygon-testnet-mumbai"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
  - type: "write_ethereum-testnet-sepolia"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`

type Delegate struct {
	registry        core.CapabilitiesRegistry
	logger          logger.Logger
	legacyEVMChains legacyevm.LegacyChainContainer
	peerID          func() *p2ptypes.PeerID
}

var _ job.Delegate = (*Delegate)(nil)

func (d *Delegate) JobType() job.Type {
	return job.Workflow
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}

func (d *Delegate) AfterJobCreated(jb job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	// NOTE: we temporarily do registration inside ServicesForSpec, this will be moved out of job specs in the future
	err := targets.InitializeWrite(d.registry, d.legacyEVMChains, d.logger)
	if err != nil {
		d.logger.Errorw("could not initialize writes", err)
	}

	dinfo, err := initializeDONInfo(d.logger)
	if err != nil {
		d.logger.Errorw("could not add initialize don info", err)
	}

	cfg := Config{
		Lggr:       d.logger,
		Spec:       spec.WorkflowSpec.Workflow,
		WorkflowID: spec.WorkflowSpec.WorkflowID,
		Registry:   d.registry,
		DONInfo:    dinfo,
		PeerID:     d.peerID,
	}
	engine, err := NewEngine(cfg)
	if err != nil {
		return nil, err
	}
	return []job.ServiceCtx{engine}, nil
}

func initializeDONInfo(lggr logger.Logger) (*capabilities.DON, error) {
	var key [16]byte

	// TODO: fetch the key and DONInfo from the registry
	keyString := "44fb5c1ee8ee48846c808a383da3aba3"
	k, err := hex.DecodeString(keyString)
	if err != nil {
		lggr.Errorf("could not decode key %s: %w", keyString, err)
	}
	key = [16]byte(k)

	p2pStrings := []string{
		"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N",
		"12D3KooWG1AyvwmCpZ93J8pBQUE1SuzrjDXnT4BeouncHR3jWLCG",
		"12D3KooWGeUKZBRMbx27FUTgBwZa9Ap9Ym92mywwpuqkEtz8XWyv",
		"12D3KooW9zYWQv3STmDeNDidyzxsJSTxoCTLicafgfeEz9nhwhC4",
	}

	p2pIDs := []p2ptypes.PeerID{}
	for _, p := range p2pStrings {
		pid := p2ptypes.PeerID{}
		err := pid.UnmarshalText([]byte(p))
		if err != nil {
			return nil, err
		}

		p2pIDs = append(p2pIDs, pid)
	}

	return &capabilities.DON{
		Members: p2pIDs,
		Config: capabilities.DONConfig{
			SharedSecret: key,
		},
	}, nil
}

func NewDelegate(logger logger.Logger, registry core.CapabilitiesRegistry, legacyEVMChains legacyevm.LegacyChainContainer, peerID func() *p2ptypes.PeerID) *Delegate {
	return &Delegate{logger: logger, registry: registry, legacyEVMChains: legacyEVMChains, peerID: peerID}
}

func ValidatedWorkflowSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{ExternalJobID: uuid.New()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, fmt.Errorf("toml error on load: %w", err)
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, fmt.Errorf("toml unmarshal error on spec: %w", err)
	}

	var spec job.WorkflowSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, fmt.Errorf("toml unmarshal error on job: %w", err)
	}

	if err := spec.Validate(); err != nil {
		return jb, err
	}

	jb.WorkflowSpec = &spec
	if jb.Type != job.Workflow {
		return jb, fmt.Errorf("unsupported type %s", jb.Type)
	}

	return jb, nil
}
