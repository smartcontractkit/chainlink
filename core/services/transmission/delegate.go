package transmission

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/transmission/handler"
)

type Delegate struct {
	lggr logger.Logger
}

func NewDelegate(lggr logger.Logger) *Delegate {
	return &Delegate{
		lggr: lggr,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.VRF
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}
func (d *Delegate) AfterJobCreated(spec job.Job)  {}
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.ServiceCtx, error) {
	if jb.TransmissionSpec == nil {
		return nil, errors.Errorf("transmission.Delegate expects a Transmission Spec, got %+v", jb)
	}

	handler := handler.NewHandler(d.lggr)
	server := NewServer(*handler, jb.TransmissionSpec.RPCPort, d.lggr)
	return []job.ServiceCtx{
		server,
	}, nil
}
