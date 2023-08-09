package logprovider

import (
	"context"
	"testing"
	"time"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogRecoverer_GetRecoverables(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := NewLogRecoverer(logger.TestLogger(t), nil, time.Millisecond*10)

	tests := []struct {
		name    string
		pending []ocr2keepers.UpkeepPayload
		want    []ocr2keepers.UpkeepPayload
		wantErr bool
	}{
		{
			"empty",
			[]ocr2keepers.UpkeepPayload{},
			[]ocr2keepers.UpkeepPayload{},
			false,
		},
		{
			"happy flow",
			[]ocr2keepers.UpkeepPayload{
				{WorkID: "1"}, {WorkID: "2"},
			},
			[]ocr2keepers.UpkeepPayload{
				{WorkID: "1"}, {WorkID: "2"},
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r.lock.Lock()
			r.pending = tc.pending
			r.lock.Unlock()

			got, err := r.GetRecoveryProposals(ctx)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Len(t, got, len(tc.want))
		})
	}
}
