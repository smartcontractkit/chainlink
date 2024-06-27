package logprovider

import (
	"math/big"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterStore_CRUD(t *testing.T) {
	tests := []struct {
		name               string
		initial            []upkeepFilter
		toAdd              []upkeepFilter
		expectedPostAdd    []upkeepFilter
		toRemove           []upkeepFilter
		expectedPostRemove []upkeepFilter
	}{
		{
			"empty",
			[]upkeepFilter{},
			[]upkeepFilter{},
			[]upkeepFilter{},
			[]upkeepFilter{},
			[]upkeepFilter{},
		},
		{
			"add rm one",
			[]upkeepFilter{},
			[]upkeepFilter{{upkeepID: big.NewInt(1)}},
			[]upkeepFilter{{upkeepID: big.NewInt(1)}},
			[]upkeepFilter{{upkeepID: big.NewInt(1)}},
			[]upkeepFilter{},
		},
		{
			"add rm multiple",
			[]upkeepFilter{},
			[]upkeepFilter{{upkeepID: big.NewInt(1)}, {upkeepID: big.NewInt(2)}},
			[]upkeepFilter{{upkeepID: big.NewInt(1)}, {upkeepID: big.NewInt(2)}},
			[]upkeepFilter{{upkeepID: big.NewInt(1)}},
			[]upkeepFilter{{upkeepID: big.NewInt(2)}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUpkeepFilterStore()
			s.AddActiveUpkeeps(tc.initial...)
			require.Equal(t, len(tc.initial), len(s.GetIDs(nil)))
			s.AddActiveUpkeeps(tc.toAdd...)
			require.Equal(t, len(tc.expectedPostAdd), s.Size())
			filters := s.GetFilters(func(f upkeepFilter) bool { return true })
			require.Equal(t, len(tc.expectedPostAdd), len(filters))
			if len(filters) > 0 {
				sort.Slice(filters, func(i, j int) bool {
					return filters[i].upkeepID.Cmp(filters[j].upkeepID) < 0
				})
				for i, f := range filters {
					require.Equal(t, tc.expectedPostAdd[i].upkeepID, f.upkeepID)
				}
			}
			s.RemoveActiveUpkeeps(tc.toRemove...)
			require.Equal(t, len(tc.expectedPostRemove), len(s.GetIDs(func(upkeepFilter) bool { return true })))
		})
	}
}

func TestFilterStore_Concurrency(t *testing.T) {
	s := NewUpkeepFilterStore()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.AddActiveUpkeeps(upkeepFilter{upkeepID: big.NewInt(1)})
		s.AddActiveUpkeeps(upkeepFilter{upkeepID: big.NewInt(2)})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.AddActiveUpkeeps(upkeepFilter{upkeepID: big.NewInt(2)})
	}()

	go func() {
		_ = s.GetIDs(nil)
	}()

	wg.Wait()

	require.Equal(t, 2, s.Size())
}
