package workflows

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

type Delegate struct {
	registry core.CapabilitiesRegistry
	logger   logger.Logger
	peerID   func() *p2ptypes.PeerID
	store    store.Store
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
	dinfo, err := initializeDONInfo(d.logger)
	if err != nil {
		d.logger.Errorw("could not add initialize don info", err)
	}

	cfg := Config{
		Lggr:          d.logger,
		Spec:          spec.WorkflowSpec.Workflow,
		WorkflowID:    spec.WorkflowSpec.WorkflowID,
		WorkflowOwner: spec.WorkflowSpec.WorkflowOwner,
		WorkflowName:  spec.WorkflowSpec.WorkflowName,
		Registry:      d.registry,
		DONInfo:       dinfo,
		PeerID:        d.peerID,
		Store:         d.store,
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
		ID:      "00010203",
		Members: p2pIDs,
		Config: capabilities.DONConfig{
			SharedSecret: key,
		},
	}, nil
}

func NewDelegate(logger logger.Logger, registry core.CapabilitiesRegistry, store store.Store, peerID func() *p2ptypes.PeerID) *Delegate {
	return &Delegate{logger: logger, registry: registry, store: store, peerID: peerID}
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
	if jb.Type != job.Workflow {
		return jb, fmt.Errorf("unsupported type %s, expected %s", jb.Type, job.Workflow)
	}

	var spec job.WorkflowSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, fmt.Errorf("toml unmarshal error on workflow spec: %w", err)
	}

	err = spec.Validate()
	if err != nil {
		return jb, fmt.Errorf("invalid WorkflowSpec: %w", err)
	}

	// ensure the embedded workflow graph is valid
	_, err = Parse(spec.Workflow)
	if err != nil {
		return jb, fmt.Errorf("failed to parse workflow graph: %w", err)
	}
	jb.WorkflowSpec = &spec
	jb.WorkflowSpecID = &spec.ID

	return jb, nil
}
