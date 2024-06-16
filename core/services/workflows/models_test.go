package workflows

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

func TestParse_Graph(t *testing.T) {
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
  - id: "a-trigger@1.0.0"
    config: {}

actions:
  - id: "an-action@1.0.0"
    config: {}
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)

consensus:
  - id: "a-consensus@1.0.0"
    config: {}
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target@1.0.0"
    config: {}
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
  - id: "a-trigger@1.0.0"
    config: {}

actions:
  - id: "an-action@1.0.0"
    config: {}
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      output: $(a-second-action.outputs)
  - id: "a-second-action@1.0.0"
    config: {}
    ref: "a-second-action"
    inputs:
      output: $(an-action.outputs)

consensus:
  - id: "a-consensus@1.0.0"
    config: {}
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target@1.0.0"
    config: {}
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
  - id: "a-trigger@1.0.0"
    config: {}

actions:
  - id: "an-action@1.0.0"
    config: {}
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(a-third-action.outputs)
  - id: "a-second-action@1.0.0"
    config: {}
    ref: "a-second-action"
    inputs:
      output: $(an-action.outputs)
  - id: "a-third-action@1.0.0"
    config: {}
    ref: "a-third-action"
    inputs:
      output: $(a-second-action.outputs)

consensus:
  - id: "a-consensus@1.0.0"
    config: {}
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target@1.0.0"
    config: {}
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
  - id: "a-trigger@1.0.0"
    config: {}

actions:
  - id: "an-action@1.0.0"
    config: {}
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(missing-action.outputs)

consensus:
  - id: "a-consensus@1.0.0"
    config: {}
    ref: "a-consensus"
    inputs:
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target@1.0.0"
    config: {}
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
  - id: "a-trigger@1.0.0"
    config: {}
  - id: "a-second-trigger@1.0.0"
    config: {}

actions:
  - id: "an-action@1.0.0"
    config: {}
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)

consensus:
  - id: "a-consensus@1.0.0"
    config: {}
    ref: "a-consensus"
    inputs:
      an-action_output: $(an-action.outputs)

targets:
  - id: "a-target@1.0.0"
    config: {}
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
  - id: "a-trigger@1.0.0"
    config: {}
  - id: "a-second-trigger@1.0.0"
    config: {}
actions:
  - id: "an-action@1.0.0"
    config: {}
    ref: "an-action"
    inputs:
      hello: "world"
consensus:
  - id: "a-consensus@1.0.0"
    config: {}
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(an-action.outputs)
targets:
  - id: "a-target@1.0.0"
    config: {}
    ref: "a-target"
    inputs:
      consensus_output: $(a-consensus.outputs)
`,
			errMsg: "all non-trigger steps must have a dependent ref",
		},
		{
			name: "duplicate edge declarations",
			yaml: `
triggers:
  - id: "a-trigger@1.0.0"
    config: {}
  - id: "a-second-trigger@1.0.0"
    config: {}
actions:
  - id: "an-action@1.0.0"
    config: {}
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
consensus:
  - id: "a-consensus@1.0.0"
    config: {}
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(an-action.outputs)
targets:
  - id: "a-target@1.0.0"
    config: {}
    ref: "a-target"
    inputs:
      consensus_output: $(a-consensus.outputs)
      consensus_output2: $(a-consensus.outputs)
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(st *testing.T) {
			wf, err := Parse(tc.yaml)
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

func TestParsesIntsCorrectly(t *testing.T) {
	wf, err := Parse(hardcodedWorkflow)
	require.NoError(t, err)

	n, err := wf.Vertex("evm_median")
	require.NoError(t, err)

	assert.Equal(t, int64(3600), n.Config["aggregation_config"].(map[string]any)["0x1111111111111111111100000000000000000000000000000000000000000000"].(map[string]any)["heartbeat"])
}
