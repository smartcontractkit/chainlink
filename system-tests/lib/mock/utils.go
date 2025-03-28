package mock_capability

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	pb2 "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
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

func MockCapabilities(nodeID string) *jobv1.ProposeJobRequest {
	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
 			type = "standardcapabilities"
 			schemaVersion = 1


 			externalJobID = "%s"
 			name = "mock-capabilitie"
 			forwardingAllowed = false
 			command = "/home/capabilities/amd64_mock"
 			config = """
 				port=7777
 				[[DefaultMocks]]
 				id="streams-trigger@2.0.0"
 				description="stream trigger mock"
 				type="trigger"
 				[[DefaultMocks]]
 				id="write_ethereum@1.0.0"
 				description="write trigger mock"
 				type="target"
 """
 `,
			uuid.NewString(),
		),
	}
}
