package fakes

import (
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

// applyFrequencyListShape reshapes the observation value to match the output
// type expected by the consensus descriptor. The simulator runs a single node,
// so the only transformation needed is for frequency_list aggregation, which
// turns a single value T into [{value: T, count: 1}]. All other aggregation
// types have matching input and output types, so the value is returned as-is.
func applyFrequencyListShape(value *valuespb.Value, descriptor *sdkpb.ConsensusDescriptor) *valuespb.Value {
	if descriptor == nil {
		return value
	}

	switch desc := descriptor.GetDescriptor_().(type) {
	case *sdkpb.ConsensusDescriptor_Aggregation:
		if descriptor.GetAggregation() == sdkpb.AggregationType_AGGREGATION_TYPE_FREQUENCY_LIST {
			return frequencyListSingleton(value)
		}

	case *sdkpb.ConsensusDescriptor_FieldsMap:
		if value == nil {
			return value
		}
		if value.GetMapValue() == nil {
			return value
		}
		fields := value.GetMapValue().GetFields()
		changed := false
		newFields := make(map[string]*valuespb.Value, len(fields))
		for k, v := range fields {
			fieldDesc := desc.FieldsMap.GetFields()[k]
			if fieldDesc != nil && fieldDesc.GetAggregation() == sdkpb.AggregationType_AGGREGATION_TYPE_FREQUENCY_LIST {
				newFields[k] = frequencyListSingleton(v)
				changed = true
			} else {
				newFields[k] = v
			}
		}
		if changed {
			return valuespb.NewMapValue(newFields)
		}
	}

	return value
}

// frequencyListSingleton wraps a single observation value as a one-element
// frequency list: [{value: value, count: 1}].
func frequencyListSingleton(value *valuespb.Value) *valuespb.Value {
	return valuespb.NewListValue([]*valuespb.Value{
		valuespb.NewMapValue(map[string]*valuespb.Value{
			"value": value,
			"count": valuespb.NewInt64Value(1),
		}),
	})
}
