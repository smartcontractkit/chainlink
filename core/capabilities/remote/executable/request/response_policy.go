package request

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

// responsePolicy allows request finalization behavior to be extended for specific
// capabilities without coupling generic request handling to chain-specific logic.
type responsePolicy interface {
	ShouldDeferIdenticalResponse(payload []byte) bool
	ObserveOKResponse(msg *types.MessageBody, metadata commoncap.ResponseMetadata)
	BuildDeterministicResponse(allowNoQuorum bool) ([]byte, bool, error)
	ShouldDeferErrorResponses() bool
	FinalizeAfterAllResponses(state responsePolicyState) (*clientResponse, bool, error)
}

type responsePolicyState struct {
	AllResponsesReceived bool
	ResponseVariants     int
	TotalErrorCount      int
}

type responsePolicyBuilder func(remoteCapabilityInfo commoncap.CapabilityInfo, capabilityMethod string) responsePolicy

var responsePolicyBuilders []responsePolicyBuilder

func registerResponsePolicyBuilder(builder responsePolicyBuilder) {
	responsePolicyBuilders = append(responsePolicyBuilders, builder)
}

func newResponsePolicy(remoteCapabilityInfo commoncap.CapabilityInfo, capabilityMethod string) responsePolicy {
	for _, builder := range responsePolicyBuilders {
		policy := builder(remoteCapabilityInfo, capabilityMethod)
		if policy != nil {
			return policy
		}
	}
	return nil
}
