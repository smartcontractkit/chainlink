package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata_padWorkflowName(t *testing.T) {
	type fields struct {
		WorkflowName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "padWorkflowName 1",
			fields: fields{
				WorkflowName: "123456789",
			},
			want: "123456789 ",
		},
		{
			name: "padWorkflowName 0",
			fields: fields{
				WorkflowName: "1234567890",
			},
			want: "1234567890",
		},
		{
			name: "padWorkflowName 10",
			fields: fields{
				WorkflowName: "",
			},
			want: "          ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metadata{
				WorkflowName: tt.fields.WorkflowName,
			}
			m.padWorkflowName()
			assert.Equal(t, tt.want, m.WorkflowName, tt.name)
		})
	}
}
