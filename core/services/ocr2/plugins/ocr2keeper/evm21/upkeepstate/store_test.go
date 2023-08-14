package upkeepstate

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUpkeepStateStore(t *testing.T) {
	tests := []struct {
		name               string
		inserts            []ocr2keepers.CheckResult
		workIDsSelect      []string
		workIDsFromScanner []string
		expected           []ocr2keepers.UpkeepState
	}{
		{
			name: "empty store",
		},
		{
			name: "save only ineligible states",
			inserts: []ocr2keepers.CheckResult{
				{
					UpkeepID: createUpkeepIDForTest(1),
					WorkID:   "0x1",
					Eligible: false,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(1),
					},
				},
				{
					UpkeepID: createUpkeepIDForTest(2),
					WorkID:   "ox2",
					Eligible: true,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(2),
					},
				},
			},
			workIDsSelect: []string{"0x1", "0x2"},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.Ineligible,
				StateUnknown,
			},
		},
		{
			name: "fetch results from scanner",
			inserts: []ocr2keepers.CheckResult{
				{
					UpkeepID: createUpkeepIDForTest(1),
					WorkID:   "0x1",
					Eligible: false,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(1),
					},
				},
			},
			workIDsSelect:      []string{"0x1", "0x2"},
			workIDsFromScanner: []string{"0x2", "0x222"},
			expected: []ocr2keepers.UpkeepState{
				ocr2keepers.Ineligible,
				ocr2keepers.Performed,
			},
		},
		{
			name: "unknown states",
			inserts: []ocr2keepers.CheckResult{
				{
					UpkeepID: createUpkeepIDForTest(1),
					WorkID:   "0x1",
					Eligible: false,
					Trigger: ocr2keepers.Trigger{
						BlockNumber: ocr2keepers.BlockNumber(1),
					},
				},
			},
			workIDsSelect:      []string{"0x2"},
			workIDsFromScanner: []string{},
			expected: []ocr2keepers.UpkeepState{
				StateUnknown,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)

			scanner := &mockScanner{}
			scanner.add(tc.workIDsFromScanner...)
			store := NewUpkeepStateStore(logger.TestLogger(t), scanner)

			for _, insert := range tc.inserts {
				assert.NoError(t, store.SetUpkeepState(ctx, insert, ocr2keepers.Performed), "storing states should not produce an error")
			}

			states, err := store.SelectByWorkIDsInRange(ctx, 1, 100, tc.workIDsSelect...)
			assert.NoError(t, err, "no error expected from selecting states")

			assert.Equal(t, len(tc.expected), len(states))
			for i, state := range states {
				assert.Equal(t, tc.expected[i], state)
			}
		})
	}
}

func TestUpkeepStateStore_Upsert(t *testing.T) {
	ctx := testutils.Context(t)
	store := NewUpkeepStateStore(logger.TestLogger(t), &mockScanner{})

	res := ocr2keepers.CheckResult{
		UpkeepID: createUpkeepIDForTest(1),
		WorkID:   "0x1",
		Eligible: false,
		Trigger: ocr2keepers.Trigger{
			BlockNumber: ocr2keepers.BlockNumber(1),
		},
	}
	require.NoError(t, store.SetUpkeepState(ctx, res, ocr2keepers.Performed))
	<-time.After(10 * time.Millisecond)
	res.Trigger.BlockNumber = ocr2keepers.BlockNumber(2)
	now := time.Now()
	require.NoError(t, store.SetUpkeepState(ctx, res, ocr2keepers.Performed))

	store.mu.Lock()
	addedAt := store.cache["0x1"].addedAt
	block := store.cache["0x1"].block
	store.mu.Unlock()

	require.True(t, now.After(addedAt))
	require.Equal(t, uint64(2), block)
}

func createUpkeepIDForTest(v int64) ocr2keepers.UpkeepIdentifier {
	uid := &ocr2keepers.UpkeepIdentifier{}
	_ = uid.FromBigInt(big.NewInt(v))

	return *uid
}

type mockScanner struct {
	lock    sync.Mutex
	workIDs []string
}

func (s *mockScanner) add(workIDs ...string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.workIDs = append(s.workIDs, workIDs...)
}

func (s *mockScanner) WorkIDsInRange(ctx context.Context, start, end int64) ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	res := s.workIDs[:]
	s.workIDs = nil
	return res, nil
}
