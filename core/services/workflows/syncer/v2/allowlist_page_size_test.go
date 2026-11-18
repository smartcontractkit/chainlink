package v2

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// gasCappedAllowlistReader serves totalAllowlistedRequests and getActiveAllowlistedRequestsReverse
// over reqs, failing any page larger than maxPage with a gas-cap error like a capped RPC node.
type gasCappedAllowlistReader struct {
	commontypes.ContractReader
	reqs     []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
	maxPage  int64
	pageSize []int64 // size of every getActiveAllowlistedRequestsReverse call, in order
}

func (m *gasCappedAllowlistReader) GetLatestValueWithHeadData(
	_ context.Context,
	readIdentifier string,
	_ primitives.ConfidenceLevel,
	params any,
	result any,
) (*commontypes.Head, error) {
	head := &commontypes.Head{Height: "1"}
	if strings.HasSuffix(readIdentifier, TotalAllowlistedRequestsMethodName) {
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(big.NewInt(int64(len(m.reqs)))))
		return head, nil
	}

	p := params.(GetActiveAllowlistedRequestsReverseParams)
	size := p.EndIndex.Int64() - p.StartIndex.Int64() + 1
	m.pageSize = append(m.pageSize, size)
	if size > m.maxPage {
		return nil, errors.New("out of gas: gas required exceeds 16777216")
	}

	var page []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
	for i := p.EndIndex.Int64(); i >= p.StartIndex.Int64(); i-- {
		page = append(page, m.reqs[i])
	}
	out := reflect.ValueOf(result).Elem()
	out.FieldByName("AllowlistedRequests").Set(reflect.ValueOf(page))
	out.FieldByName("SearchComplete").SetBool(p.StartIndex.Sign() == 0)
	return head, nil
}

func TestGetAllowlistedRequests_AdaptivePageSize(t *testing.T) {
	t.Parallel()

	reqs := make([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, 0, 700)
	for i := range 700 {
		reqs = append(reqs, workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
			RequestDigest:   [32]byte{byte(i), byte(i >> 8)},
			Owner:           common.Address{1},
			ExpiryTimestamp: 1,
		})
	}
	reader := &gasCappedAllowlistReader{reqs: reqs, maxPage: 100}
	w := &workflowRegistry{
		lggr:                             logger.TestLogger(t),
		lastSeenAllowlistedRequestsCount: big.NewInt(0),
	}

	got, total, _, err := w.getAllowlistedRequests(t.Context(), reader)
	require.NoError(t, err)
	assert.Equal(t, int64(700), total.Int64())

	// Default 300 fails, 150 fails, then 75 succeeds; every request is returned exactly once.
	assert.Equal(t, []int64{300, 150, 75}, reader.pageSize[:3])
	require.Len(t, got, 700)
	seen := map[[32]byte]bool{}
	for _, r := range got {
		assert.False(t, seen[r.RequestDigest], "duplicate request")
		seen[r.RequestDigest] = true
	}
}

func TestGetAllowlistedRequests_GasErrorAtMinPageSize(t *testing.T) {
	t.Parallel()

	reader := &gasCappedAllowlistReader{
		reqs:    make([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, 700),
		maxPage: 0, // every page exceeds the "gas cap"
	}
	w := &workflowRegistry{
		lggr:                             logger.TestLogger(t),
		lastSeenAllowlistedRequestsCount: big.NewInt(0),
	}

	// It shrinks to the floor and then returns the error instead of looping forever.
	_, _, _, err := w.getAllowlistedRequests(t.Context(), reader)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "out of gas")
	assert.Equal(t, []int64{300, 150, 75, 37, 18, 10}, reader.pageSize)
}
