package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.Workflows = (*workflowsConfig)(nil)

type workflowsConfig struct {
	c toml.Workflows
}

func (w *workflowsConfig) Limits() config.WorkflowsLimits {
	return &limits{
		l: w.c.Limits,
	}
}

type limits struct {
	l toml.Limits
}

func (l *limits) Global() int32 {
	return *l.l.Global
}

func (l *limits) PerOwner() int32 {
	return *l.l.PerOwner
}

func (l *limits) PerOwnerOverrides() map[string]int32 {
	return l.l.Overrides
}
