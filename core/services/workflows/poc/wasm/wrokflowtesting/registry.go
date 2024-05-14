package wrokflowtesting

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

func NewRegistry() *Registry {
	return &Registry{
		remoteActionsAndConsensus: map[string]func(value values.Value) (values.Value, error){},
		remoteTargets:             map[string]func(value values.Value) error{},
	}
}

type Registry struct {
	triggers                  []*valueAndRef
	remoteActionsAndConsensus map[string]func(value values.Value) (values.Value, error)
	remoteTargets             map[string]func(value values.Value) error
}

// RegisterTrigger should be called from generated code to assure type safety
func (r *Registry) RegisterTrigger(ref string, value any) error {
	wrapped, err := values.Wrap(value)
	if err != nil {
		return err
	}
	r.triggers = append(r.triggers, &valueAndRef{value: wrapped, ref: ref})
	return nil
}

// RegisterRemoteAction should be called from generated code to assure type safety
func (r *Registry) RegisterRemoteAction(ref string, action func(value values.Value) (values.Value, error)) {
	r.remoteActionsAndConsensus[ref] = action
}

// RegisterRemoteConsensus should be called from generated code to assure type safety
func (r *Registry) RegisterRemoteConsensus(ref string, consensus func(value values.Value) (values.Value, error)) {
	r.remoteActionsAndConsensus[ref] = consensus
}

// RegisterRemoteTarget should be called from generated code to assure type safety
func (r *Registry) RegisterRemoteTarget(ref string, target func(value values.Value) error) {
	r.remoteTargets[ref] = target
}

func (r *Registry) next() *valueAndRef {
	if len(r.triggers) == 0 {
		return nil
	}
	t := r.triggers[0]
	r.triggers = r.triggers[1:]
	return t
}

type valueAndRef struct {
	value values.Value
	ref   string
}
