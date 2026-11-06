package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogsContainMatchingLine(t *testing.T) {
	t.Parallel()

	const registration = `{"msg":"Registered triggers via coordinator","version":"2.70.0","workflowID":"workflow-a"}`
	for _, tc := range []struct {
		name    string
		content string
		needles []string
		want    bool
	}{
		{
			name: "message and workflow on the same line", content: registration,
			needles: []string{"Registered triggers via coordinator", "workflow-a"}, want: true,
		},
		{
			name: "another workflow does not count", content: registration,
			needles: []string{"Registered triggers via coordinator", "workflow-b"},
		},
		{
			name:    "needles on different lines do not count",
			content: registration + "\n" + `{"msg":"handled event (WorkflowPaused)","workflowID":"workflow-b"}`,
			needles: []string{"handled event (WorkflowPaused)", "workflow-a"},
		},
		{
			name: "repeated messages still match once per container", content: registration + "\n" + registration,
			needles: []string{"Registered triggers via coordinator", "workflow-a"}, want: true,
		},
		{
			name: "single needle remains supported", content: registration,
			needles: []string{"Registered triggers via coordinator"}, want: true,
		},
		{name: "empty logs", needles: []string{"workflow-a"}},
		{name: "no needles", content: registration},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, logsContainMatchingLine(tc.content, tc.needles...))
		})
	}
}
