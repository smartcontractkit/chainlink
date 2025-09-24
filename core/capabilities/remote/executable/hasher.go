package executable

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

// V1 Capabilities only need a hasher for the ChainWrite Target.
// This hasher excludes signatures from the Inputs map when hashing the request.
type v1Hasher struct {
	requestHashExcludedAttributes []string
}

func (r *v1Hasher) Hash(msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}

	// An attribute called StepDependency is used to define a data dependency between steps,
	// and not to provide input values; we should therefore disregard it when hashing the request
	if len(r.requestHashExcludedAttributes) == 0 {
		r.requestHashExcludedAttributes = []string{"StepDependency"}
	}

	for _, path := range r.requestHashExcludedAttributes {
		if req.Inputs != nil {
			req.Inputs.DeleteAtPath(path)
		}
	}

	reqBytes, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal capability request: %w", err)
	}
	hash := sha256.Sum256(reqBytes)
	return hash, nil
}

func NewV1Hasher(requestHashExcludedAttributes []string) types.MessageHasher {
	return &v1Hasher{
		requestHashExcludedAttributes: requestHashExcludedAttributes,
	}
}

// V2 Capabilities (Executables) default to a simple hasher that hashes the entire payload.
// WriteReport methods use a hasher that excludes signatures from the WriteReportRequest.
// Additional hashers can be added here as needed.
type simpleHasher struct {
}

func (r *simpleHasher) Hash(msg *types.MessageBody) ([32]byte, error) {
	return sha256.Sum256(msg.Payload), nil
}

func NewSimpleHasher() types.MessageHasher {
	return &simpleHasher{}
}

type writeReportExcludeSignaturesHasher struct {
}

func (r *writeReportExcludeSignaturesHasher) Hash(msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}
	if req.Payload == nil {
		return [32]byte{}, errors.New("capability request payload is nil")
	}

	var wrReq evmcappb.WriteReportRequest
	if err = req.Payload.UnmarshalTo(&wrReq); err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal Payload to WriteReportRequest: %w", err)
	}
	if wrReq.Report == nil {
		return [32]byte{}, errors.New("WriteReportRequest.Report is nil")
	}

	wrReq.Report.Sigs = nil // exclude signatures from hash
	filteredPayload, err := proto.Marshal(&wrReq)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal WriteReportRequest without signatures: %w", err)
	}
	return sha256.Sum256(filteredPayload), nil
}

func NewWriteReportExcludeSignaturesHasher() types.MessageHasher {
	return &writeReportExcludeSignaturesHasher{}
}
