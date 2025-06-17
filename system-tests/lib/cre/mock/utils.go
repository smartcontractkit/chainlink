package mockcapability

import (
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	pb2 "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

func MapToBytes(m *values.Map) ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	pm := make(map[string]*pb2.Value)
	for k, v := range m.Underlying {
		pm[k] = values.Proto(v)
	}
	bytes, err := proto.Marshal(pb2.NewMapValue(pm))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func BytesToMap(b []byte) (*values.Map, error) {
	var o pb2.Value
	if err := proto.Unmarshal(b, &o); err != nil {
		return nil, err
	}

	vm := values.Map{Underlying: make(map[string]values.Value)}

	if o.Value == nil {
		return &vm, nil
	}

	for k, v := range o.GetMapValue().Fields {
		val, err := values.FromProto(v)
		if err != nil {
			return nil, err
		}
		vm.Underlying[k] = val
	}

	return &vm, nil
}
