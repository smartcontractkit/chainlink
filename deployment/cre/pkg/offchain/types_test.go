package offchain

import (
	"testing"

	"github.com/stretchr/testify/require"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
)

func TestAddToFilter(t *testing.T) {
	t.Run("CSA public key appends", func(t *testing.T) {
		req := require.New(t)

		f := TargetDONFilter{Key: FilterKeyCSAPublicKey, Value: "0xabc"}
		filter := &nodev1.ListNodesRequest_Filter{}

		got := f.AddToFilter(filter)

		req.Equal(filter, got, "expected in-place mutation, got new pointer")
		req.Len(got.PublicKeys, 1)
		req.Equal("0xabc", got.PublicKeys[0])
		req.Empty(got.Selectors)
	})

	t.Run("DON name adds EXIST selector with don- prefix", func(t *testing.T) {
		req := require.New(t)

		f := TargetDONFilter{Key: FilterKeyDONName, Value: "eu-west"}
		filter := &nodev1.ListNodesRequest_Filter{}

		got := f.AddToFilter(filter)

		req.Equal(filter, got, "expected in-place mutation")
		req.Len(got.Selectors, 1)

		s := got.Selectors[0]
		req.Equal(ptypes.SelectorOp_EXIST, s.Op)
		req.Equal("don-eu-west", s.Key)
		req.Nil(s.Value)
	})

	t.Run("generic key adds EQ selector with value pointer", func(t *testing.T) {
		req := require.New(t)

		f := TargetDONFilter{Key: "region", Value: "us-east-1"}
		filter := &nodev1.ListNodesRequest_Filter{}

		got := f.AddToFilter(filter)

		req.Equal(filter, got, "expected in-place mutation")
		req.Len(got.Selectors, 1)

		s := got.Selectors[0]
		req.Equal(ptypes.SelectorOp_EQ, s.Op)
		req.Equal("region", s.Key)
		req.NotNil(s.Value)
		req.Equal("us-east-1", *s.Value)
	})
}

func TestAddToFilterIfNotPresent(t *testing.T) {
	t.Run("CSA key dedupes", func(t *testing.T) {
		req := require.New(t)

		filter := &nodev1.ListNodesRequest_Filter{PublicKeys: []string{"0xabc"}}
		f := TargetDONFilter{Key: FilterKeyCSAPublicKey, Value: "0xabc"}

		got := f.AddToFilterIfNotPresent(filter)

		req.Equal(filter, got, "expected in-place mutation")
		req.Len(got.PublicKeys, 1)
		req.Equal([]string{"0xabc"}, got.PublicKeys)
	})

	t.Run("CSA key adds if missing", func(t *testing.T) {
		req := require.New(t)

		filter := &nodev1.ListNodesRequest_Filter{PublicKeys: []string{"0xabc"}}
		f := TargetDONFilter{Key: FilterKeyCSAPublicKey, Value: "0xdef"}

		got := f.AddToFilterIfNotPresent(filter)

		req.Equal(filter, got, "expected in-place mutation")
		req.Len(got.PublicKeys, 2)
		req.Equal([]string{"0xabc", "0xdef"}, got.PublicKeys)
	})

	t.Run("DON selector skips if same DON present", func(t *testing.T) {
		req := require.New(t)

		filter := &nodev1.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{{Op: ptypes.SelectorOp_EXIST, Key: "don-eu-west"}},
		}
		f := TargetDONFilter{Key: FilterKeyDONName, Value: "eu-west"}

		got := f.AddToFilterIfNotPresent(filter)

		req.Equal(filter, got, "expected in-place mutation")
		req.Len(got.Selectors, 1)
	})

	t.Run("DON selector adds if not present", func(t *testing.T) {
		req := require.New(t)

		filter := &nodev1.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{{Op: ptypes.SelectorOp_EXIST, Key: "don-us-west"}},
		}
		f := TargetDONFilter{Key: FilterKeyDONName, Value: "eu-west"}

		got := f.AddToFilterIfNotPresent(filter)

		req.Equal(filter, got, "expected in-place mutation")
		req.Len(got.Selectors, 2)

		last := got.Selectors[len(got.Selectors)-1]
		req.Equal(ptypes.SelectorOp_EXIST, last.Op)
		req.Equal("don-eu-west", last.Key)
		req.Nil(last.Value) // EXIST should have nil value
	})

	t.Run("generic selector skips if same key present", func(t *testing.T) {
		req := require.New(t)

		val := "eu"
		filter := &nodev1.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{{Op: ptypes.SelectorOp_EQ, Key: "region", Value: &val}},
		}
		f := TargetDONFilter{Key: "region", Value: "us"}

		got := f.AddToFilterIfNotPresent(filter)

		req.Equal(filter, got, "expected in-place mutation")
		req.Len(got.Selectors, 1)
	})

	t.Run("generic selector adds if key missing", func(t *testing.T) {
		req := require.New(t)

		val := "mainnet"
		filter := &nodev1.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{{Op: ptypes.SelectorOp_EQ, Key: "network", Value: &val}},
		}
		f := TargetDONFilter{Key: "region", Value: "us-east-1"}

		got := f.AddToFilterIfNotPresent(filter)

		req.Equal(filter, got, "expected in-place mutation")
		req.Len(got.Selectors, 2)

		last := got.Selectors[len(got.Selectors)-1]
		req.Equal(ptypes.SelectorOp_EQ, last.Op)
		req.Equal("region", last.Key)
		req.NotNil(last.Value)
		req.Equal("us-east-1", *last.Value)
	})
}
