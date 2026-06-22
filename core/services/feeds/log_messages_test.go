package feeds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isSyncNodeInfoLogMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		message string
		want    bool
	}{
		{message: `failed to sync node info attempt="0" err="err"`, want: true},
		{message: `failed to sync node info; aborting err="err"`, want: true},
		{message: "successfully synced node info", want: true},
		{message: "evm chain started", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, isSyncNodeInfoLogMessage(tt.message))
		})
	}
}

func Test_filterSyncNodeInfoLogMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "keeps sync retry and success lines",
			in: []string{
				"evm chain started",
				`failed to sync node info attempt="0" err="SyncNodeInfo.UpdateNode call partially failed: error chain 3"`,
				`failed to sync node info attempt="1" err="SyncNodeInfo.UpdateNode call failed: error-4"`,
				"successfully synced node info",
			},
			want: []string{
				`failed to sync node info attempt="0" err="SyncNodeInfo.UpdateNode call partially failed: error chain 3"`,
				`failed to sync node info attempt="1" err="SyncNodeInfo.UpdateNode call failed: error-4"`,
				"successfully synced node info",
			},
		},
		{
			name: "keeps abort log",
			in: []string{
				"noise",
				`failed to sync node info; aborting err="SyncNodeInfo.UpdateNode call partially failed: error chain 12"`,
			},
			want: []string{
				`failed to sync node info; aborting err="SyncNodeInfo.UpdateNode call partially failed: error chain 12"`,
			},
		},
		{
			name: "empty when no sync logs",
			in:   []string{"only unrelated logs"},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, filterSyncNodeInfoLogMessages(tt.in))
		})
	}
}
