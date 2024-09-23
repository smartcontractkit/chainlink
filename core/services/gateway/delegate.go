package gateway

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

type Delegate struct {
	legacyChains legacyevm.LegacyChainContainer
	ks           keystore.Eth
	ds           sqlutil.DataSource
	lggr         logger.Logger
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(legacyChains legacyevm.LegacyChainContainer, ks keystore.Eth, ds sqlutil.DataSource, lggr logger.Logger) *Delegate {
	return &Delegate{
		legacyChains: legacyChains,
		ks:           ks,
		ds:           ds,
		lggr:         lggr,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.Gateway
}

func (d *Delegate) BeforeJobCreated(spec job.Job)              {}
func (d *Delegate) AfterJobCreated(spec job.Job)               {}
func (d *Delegate) BeforeJobDeleted(spec job.Job)              {}
func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

// ServicesForSpec returns the scheduler to be used for running observer jobs
func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) (services []job.ServiceCtx, err error) {
	if spec.GatewaySpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *jobSpec.GatewaySpec to be present, got %v", spec)
	}

	var gatewayConfig config.GatewayConfig
	err2 := json.Unmarshal(spec.GatewaySpec.GatewayConfig.Bytes(), &gatewayConfig)
	if err2 != nil {
		return nil, errors.Wrap(err2, "unmarshal gateway config")
	}
	httpClient, err := network.NewHTTPClient(gatewayConfig.HTTPClientConfig, d.lggr)
	if err != nil {
		return nil, err
	}
	handlerFactory := NewHandlerFactory(d.legacyChains, d.ds, httpClient, d.lggr)
	gateway, err := NewGatewayFromConfig(&gatewayConfig, handlerFactory, d.lggr)
	if err != nil {
		return nil, err
	}

	return []job.ServiceCtx{gateway}, nil
}

func ValidatedGatewaySpec(tomlString string) (job.Job, error) {
	var jb = job.Job{ExternalJobID: uuid.New()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load")
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on spec")
	}

	var spec job.GatewaySpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}

	jb.GatewaySpec = &spec
	if jb.Type != job.Gateway {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}

	return jb, nil
}
