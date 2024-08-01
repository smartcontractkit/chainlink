package pb

import (
	"errors"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const (
	CapabilityTypeUnknown   = CapabilityType_CAPABILITY_TYPE_UNKNOWN
	CapabilityTypeTrigger   = CapabilityType_CAPABILITY_TYPE_TRIGGER
	CapabilityTypeAction    = CapabilityType_CAPABILITY_TYPE_ACTION
	CapabilityTypeConsensus = CapabilityType_CAPABILITY_TYPE_CONSENSUS
	CapabilityTypeTarget    = CapabilityType_CAPABILITY_TYPE_TARGET
)

func MarshalCapabilityRequest(req capabilities.CapabilityRequest) ([]byte, error) {
	return proto.MarshalOptions{Deterministic: true}.Marshal(CapabilityRequestToProto(req))
}

func MarshalCapabilityResponse(resp capabilities.CapabilityResponse) ([]byte, error) {
	return proto.MarshalOptions{Deterministic: true}.Marshal(CapabilityResponseToProto(resp))
}

func UnmarshalCapabilityRequest(raw []byte) (capabilities.CapabilityRequest, error) {
	var cr CapabilityRequest
	if err := proto.Unmarshal(raw, &cr); err != nil {
		return capabilities.CapabilityRequest{}, err
	}
	return CapabilityRequestFromProto(&cr)
}

func UnmarshalCapabilityResponse(raw []byte) (capabilities.CapabilityResponse, error) {
	var cr CapabilityResponse
	if err := proto.Unmarshal(raw, &cr); err != nil {
		return capabilities.CapabilityResponse{}, err
	}
	return CapabilityResponseFromProto(&cr)
}

func CapabilityRequestToProto(req capabilities.CapabilityRequest) *CapabilityRequest {
	inputs := values.EmptyMap()
	if req.Inputs != nil {
		inputs = req.Inputs
	}
	config := values.EmptyMap()
	if req.Config != nil {
		config = req.Config
	}
	return &CapabilityRequest{
		Metadata: &RequestMetadata{
			WorkflowId:               req.Metadata.WorkflowID,
			WorkflowExecutionId:      req.Metadata.WorkflowExecutionID,
			WorkflowOwner:            req.Metadata.WorkflowOwner,
			WorkflowName:             req.Metadata.WorkflowName,
			WorkflowDonId:            req.Metadata.WorkflowDonID,
			WorkflowDonConfigVersion: req.Metadata.WorkflowDonConfigVersion,
		},
		Inputs: values.ProtoMap(inputs),
		Config: values.ProtoMap(config),
	}
}

func CapabilityResponseToProto(resp capabilities.CapabilityResponse) *CapabilityResponse {
	errStr := ""
	if resp.Err != nil {
		errStr = resp.Err.Error()
	}

	return &CapabilityResponse{
		Error: errStr,
		Value: values.ProtoMap(resp.Value),
	}
}

func CapabilityRequestFromProto(pr *CapabilityRequest) (capabilities.CapabilityRequest, error) {
	md := pr.Metadata
	config, err := values.FromMapValueProto(pr.Config)
	if err != nil {
		return capabilities.CapabilityRequest{}, err
	}

	inputs, err := values.FromMapValueProto(pr.Inputs)
	if err != nil {
		return capabilities.CapabilityRequest{}, err
	}

	return capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:               md.WorkflowId,
			WorkflowExecutionID:      md.WorkflowExecutionId,
			WorkflowOwner:            md.WorkflowOwner,
			WorkflowName:             md.WorkflowName,
			WorkflowDonID:            md.WorkflowDonId,
			WorkflowDonConfigVersion: md.WorkflowDonConfigVersion,
		},
		Config: config,
		Inputs: inputs,
	}, nil
}

func CapabilityResponseFromProto(pr *CapabilityResponse) (capabilities.CapabilityResponse, error) {
	val, err := values.FromMapValueProto(pr.Value)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	if pr.Error != "" {
		err = errors.New(pr.Error)
	}

	resp := capabilities.CapabilityResponse{
		Err:   err,
		Value: val,
	}

	return resp, nil
}
