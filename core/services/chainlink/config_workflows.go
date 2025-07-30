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
	return &limitsCfg{
		l: w.c.Limits,
	}
}

type limitsCfg struct {
	l toml.Limits
}

func (l *limitsCfg) Global() int32 {
	return *l.l.Global
}

func (l *limitsCfg) PerOwner() int32 {
	return *l.l.PerOwner
}

func (l *limitsCfg) PerOwnerOverrides() map[string]int32 {
	return l.l.Overrides
}
