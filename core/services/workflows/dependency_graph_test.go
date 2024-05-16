package workflows_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDependencyGraph(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		yaml   string
		graph  map[string]map[string]struct{}
		errMsg string
	}{
		{
			name: "basic example",
			yaml: `
triggers:
  - id: "a-trigger"

actions:
  - id: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)

consensus:
  - id: "a-consensus"
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target"
    ref: "a-target"
    inputs: 
      consensus_output: $(a-consensus.outputs)
`,
			graph: map[string]map[string]struct{}{
				workflows.KeywordTrigger: {
					"an-action":   struct{}{},
					"a-consensus": struct{}{},
				},
				"an-action": {
					"a-consensus": struct{}{},
				},
				"a-consensus": {
					"a-target": struct{}{},
				},
				"a-target": {},
			},
		},
		{
			name: "circular relationship",
			yaml: `
triggers:
  - id: "a-trigger"

actions:
  - id: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      output: $(a-second-action.outputs)
  - id: "a-second-action"
    ref: "a-second-action"
    inputs:
      output: $(an-action.outputs)

consensus:
  - id: "a-consensus"
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target"
    ref: "a-target"
    inputs: 
      consensus_output: $(a-consensus.outputs)
`,
			errMsg: "edge would create a cycle",
		},
		{
			name: "indirect circular relationship",
			yaml: `
triggers:
  - id: "a-trigger"

actions:
  - id: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(a-third-action.outputs)
  - id: "a-second-action"
    ref: "a-second-action"
    inputs:
      output: $(an-action.outputs)
  - id: "a-third-action"
    ref: "a-third-action"
    inputs:
      output: $(a-second-action.outputs)

consensus:
  - id: "a-consensus"
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target"
    ref: "a-target"
    inputs: 
      consensus_output: $(a-consensus.outputs)
`,
			errMsg: "edge would create a cycle",
		},
		{
			name: "relationship doesn't exist",
			yaml: `
triggers:
  - id: "a-trigger"

actions:
  - id: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(missing-action.outputs)

consensus:
  - id: "a-consensus"
    ref: "a-consensus"
    inputs:
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target"
    ref: "a-target"
    inputs: 
      consensus_output: $(a-consensus.outputs)
`,
			errMsg: "source vertex missing-action: vertex not found",
		},
		{
			name: "two trigger nodes",
			yaml: `
triggers:
  - id: "a-trigger"
  - id: "a-second-trigger"

actions:
  - id: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)

consensus:
  - id: "a-consensus"
    ref: "a-consensus"
    inputs:
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target"
    ref: "a-target"
    inputs:
      consensus_output: $(a-consensus.outputs)
`,
			graph: map[string]map[string]struct{}{
				workflows.KeywordTrigger: {
					"an-action": struct{}{},
				},
				"an-action": {
					"a-consensus": struct{}{},
				},
				"a-consensus": {
					"a-target": struct{}{},
				},
				"a-target": {},
			},
		},
		{
			name: "non-trigger step with no dependent refs",
			yaml: `
triggers:
  - id: "a-trigger"
  - id: "a-second-trigger"
actions:
  - id: "an-action"
    ref: "an-action"
    inputs:
      hello: "world"
consensus:
  - id: "a-consensus"
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(an-action.outputs)
targets:
  - id: "a-target"
    ref: "a-target"
    inputs:
      consensus_output: $(a-consensus.outputs)
`,
			errMsg: "all non-trigger steps must have a dependent ref",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(st *testing.T) {
			//wf, err := workflows.Parse(tc.yaml)
			wf, err := workflows.ParseDependencyGraph(tc.yaml)
			if tc.errMsg != "" {
				assert.ErrorContains(st, err, tc.errMsg)
			} else {
				require.NoError(st, err)

				adjacencies, err := wf.AdjacencyMap()
				require.NoError(t, err)

				got := map[string]map[string]struct{}{}
				for k, v := range adjacencies {
					if _, ok := got[k]; !ok {
						got[k] = map[string]struct{}{}
					}
					for adj := range v {
						got[k][adj] = struct{}{}
					}
				}

				assert.Equal(st, tc.graph, got, adjacencies)
			}
		})
	}
}
