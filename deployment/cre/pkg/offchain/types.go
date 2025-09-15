package offchain

import "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

const FilterKeyDONName = "don_name"

type TargetDONFilter struct {
	Key   string
	Value string
}

func (f TargetDONFilter) ToJDSelector() *ptypes.Selector {
	// DON name is a key, so we just check for its existence instead of equality
	if f.Key == FilterKeyDONName {
		return &ptypes.Selector{
			Op:  ptypes.SelectorOp_EXIST,
			Key: f.Value,
		}
	}

	return &ptypes.Selector{
		Op:    ptypes.SelectorOp_EQ,
		Key:   f.Key,
		Value: &f.Value,
	}
}
