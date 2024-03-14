package workflows

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Graph(t *testing.T) {
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
  - type: "a-trigger"

actions:
  - type: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)

consensus:
  - type: "a-consensus"
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      an-action_output: $(an-action.outputs)

targets:
  - type: "a-target"
    ref: "a-target"
    inputs: 
      consensus_output: $(a-consensus.outputs)
`,
			graph: map[string]map[string]struct{}{
				keywordTrigger: {
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
  - type: "a-trigger"

actions:
  - type: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      output: $(a-second-action.outputs)
  - type: "a-second-action"
    ref: "a-second-action"
    inputs:
      output: $(an-action.outputs)

consensus:
  - type: "a-consensus"
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      an-action_output: $(an-action.outputs)

targets:
  - type: "a-target"
    ref: "a-target"
    inputs: 
      consensus_output: $(a-consensus.outputs)
`,
			errMsg: "found circular relationship",
		},
		{
			name: "indirect circular relationship",
			yaml: `
triggers:
  - type: "a-trigger"

actions:
  - type: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(a-third-action.outputs)
  - type: "a-second-action"
    ref: "a-second-action"
    inputs:
      output: $(an-action.outputs)
  - type: "a-third-action"
    ref: "a-third-action"
    inputs:
      output: $(a-second-action.outputs)

consensus:
  - type: "a-consensus"
    ref: "a-consensus"
    inputs:
      trigger_output: $(trigger.outputs)
      an-action_output: $(an-action.outputs)

targets:
  - type: "a-target"
    ref: "a-target"
    inputs: 
      consensus_output: $(a-consensus.outputs)
`,
			errMsg: "found circular relationship",
		},
		{
			name: "relationship doesn't exist",
			yaml: `
triggers:
  - type: "a-trigger"

actions:
  - type: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(missing-action.outputs)

consensus:
  - type: "a-consensus"
    ref: "a-consensus"
    inputs:
      an-action_output: $(an-action.outputs)

targets:
  - type: "a-target"
    ref: "a-target"
    inputs: 
      consensus_output: $(a-consensus.outputs)
`,
			errMsg: "invalid reference missing-action found in workflow spec",
		},
		{
			name: "relationship doesn't exist",
			yaml: `
triggers:
  - type: "a-trigger"

actions:
  - type: "an-action"
    ref: "an-action"
    inputs:
      trigger_output: $(trigger.outputs)
      action_output: $(missing-action.outputs)

consensus:
  - type: "a-consensus"
    ref: "a-consensus"
    inputs:
      an-action_output: $(an-action.outputs)

targets:
  - type: "a-target"
    ref: "a-target"
    inputs: 
      consensus_output: $(a-consensus.outputs)
`,
			errMsg: "invalid reference missing-action found in workflow spec",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(st *testing.T) {
			wf, err := Parse(tc.yaml)
			if tc.errMsg != "" {
				assert.ErrorContains(st, err, tc.errMsg)
			} else {
				require.NoError(st, err)
				assert.Equal(st, wf.graph.adjacencies, tc.graph)
			}
		})
	}
}
