package executable

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	aptoscappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/aptos"
	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	solcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/solana"
	stellarcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/stellar"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

// V1 Capabilities only need a hasher for the ChainWrite Target.
// This hasher excludes signatures from the Inputs map when hashing the request.
type v1Hasher struct {
	requestHashExcludedAttributes []string
}

func (r *v1Hasher) Hash(ctx context.Context, msg *types.MessageBody) ([32]byte, error) {
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

// V2 Capabilities (Executables) default to a simple hasher that hashes an
// explicit allowlist of metadata fields. WriteReport methods use a hasher that
// also excludes signatures from the WriteReportRequest.
// Additional hashers can be added here as needed.
type simpleHasher struct {
	cfg OptInHasherConfig
}

func (r *simpleHasher) Hash(ctx context.Context, msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}

	req.Metadata = applyMetadataFields(ctx, req.Metadata, r.cfg)

	reqBytes, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal capability request: %w", err)
	}
	hash := sha256.Sum256(reqBytes)
	return hash, nil
}

func NewSimpleHasher(cfg OptInHasherConfig) types.MessageHasher {
	return &simpleHasher{cfg: cfg}
}

type writeReportExcludeSignaturesHasher struct {
	cfg OptInHasherConfig
}

func (r *writeReportExcludeSignaturesHasher) Hash(ctx context.Context, msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}
	if req.Payload == nil {
		return [32]byte{}, errors.New("capability request payload is nil")
	}

	req.Metadata = applyMetadataFields(ctx, req.Metadata, r.cfg)

	family, familyErr := getWriteReportFamily(msg)
	if familyErr != nil {
		return [32]byte{}, familyErr
	}

	var payload *anypb.Any
	switch family {
	case writeReportFamilyEVM:
		var wrReq evmcappb.WriteReportRequest
		if err = req.Payload.UnmarshalTo(&wrReq); err != nil {
			return [32]byte{}, fmt.Errorf("failed to unmarshal Payload to WriteReportRequest: %w", err)
		}
		if wrReq.Report == nil {
			return [32]byte{}, errors.New("WriteReportRequest.Report is nil")
		}

		wrReq.Report.Sigs = nil // exclude signatures from hash
		payload, err = anypb.New(&wrReq)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to marshal WriteReportRequest back to anypb: %w", err)
		}
	case writeReportFamilySolana:
		var wrReq solcappb.WriteReportRequest
		if err = req.Payload.UnmarshalTo(&wrReq); err != nil {
			return [32]byte{}, fmt.Errorf("failed to unmarshal Payload to WriteReportRequest: %w", err)
		}
		if wrReq.Report == nil {
			return [32]byte{}, errors.New("WriteReportRequest.Report is nil")
		}

		wrReq.Report.Sigs = nil // exclude signatures from hash
		payload, err = anypb.New(&wrReq)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to marshal WriteReportRequest back to anypb: %w", err)
		}
	case writeReportFamilyAptos:
		var wrReq aptoscappb.WriteReportRequest
		if err = req.Payload.UnmarshalTo(&wrReq); err != nil {
			return [32]byte{}, fmt.Errorf("failed to unmarshal Payload to WriteReportRequest: %w", err)
		}
		if wrReq.Report == nil {
			return [32]byte{}, errors.New("WriteReportRequest.Report is nil")
		}

		wrReq.Report.Sigs = nil // exclude signatures from hash
		payload, err = anypb.New(&wrReq)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to marshal WriteReportRequest back to anypb: %w", err)
		}
	case writeReportFamilyStellar:
		var wrReq stellarcappb.WriteReportRequest
		if err = req.Payload.UnmarshalTo(&wrReq); err != nil {
			return [32]byte{}, fmt.Errorf("failed to unmarshal Payload to WriteReportRequest: %w", err)
		}
		if wrReq.Report == nil {
			return [32]byte{}, errors.New("WriteReportRequest.Report is nil")
		}

		wrReq.Report.Sigs = nil // exclude signatures from hash
		payload, err = anypb.New(&wrReq)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to marshal WriteReportRequest back to anypb: %w", err)
		}
	default:
		return [32]byte{}, fmt.Errorf("unexpected report family: %s", family)
	}

	req.Payload = payload

	reqBytes, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal capability request: %w", err)
	}
	return sha256.Sum256(reqBytes), nil
}

type writeReportFamily string

var (
	writeReportFamilyEVM     writeReportFamily = "evm"
	writeReportFamilySolana  writeReportFamily = "solana"
	writeReportFamilyAptos   writeReportFamily = "aptos"
	writeReportFamilyStellar writeReportFamily = "stellar"
)

func getWriteReportFamily(msg *types.MessageBody) (writeReportFamily, error) {
	ss := strings.Split(msg.CapabilityId, ":")
	if len(ss) < 1 {
		return "", errors.New("failed to parse family from capability id")
	}
	family := ss[0]
	switch family {
	case "evm":
		return writeReportFamilyEVM, nil
	case "solana":
		return writeReportFamilySolana, nil
	case "aptos":
		return writeReportFamilyAptos, nil
	case "stellar":
		return writeReportFamilyStellar, nil
	}

	return "", errors.New("report family is unknown, available families: evm, solana, aptos, stellar")
}

func NewWriteReportExcludeSignaturesHasher(cfg OptInHasherConfig) types.MessageHasher {
	return &writeReportExcludeSignaturesHasher{cfg: cfg}
}

// OptInHasherConfig holds per-field feature flags that control which optional
// metadata fields are included in the request hash beyond the base allowlist.
//
// To opt in a new field in the future:
//  1. Add a feature flag (Setting[Range[config.Timestamp]]) in cresettings,
//     named FeatureRequestHashInclude<Field>ActivePeriod.
//  2. Add a field of type limits.RangeLimiter[config.Timestamp] to this struct.
//  3. Add the conditional copy in applyMetadataFields.
//  4. Construct the limiter in launcher.NewLauncher and pass it via OptInHasherConfig.
//
// To EXCLUDE a field that is currently included (e.g. WorkflowTag):
//  1. The field must already be under a per-field flag that is ON by default.
//  2. After all nodes are rolled out, set the flag's window to the far future
//     to deactivate it.
//
// Each flag is evaluated against the request's DON-derived ExecutionTimestamp,
// so all nodes processing the same request include or exclude the field
// atomically — no cross-version requestID divergence.
type OptInHasherConfig struct {
	// IncludeWorkflowTag is ON by default (window covers all timestamps including
	// zero time.Time{}), so WorkflowTag is included in the hash matching current
	// prod behavior. After rollout, set to far-future window to exclude it.
	IncludeWorkflowTag limits.RangeLimiter[config.Timestamp]
}

// baseMetadataFields returns a copy of the metadata containing only the
// base allowlisted fields. New metadata fields added in the future are excluded
// from the hash by default, preventing cross-version requestID divergence.
func baseMetadataFields(md capabilities.RequestMetadata) capabilities.RequestMetadata {
	return capabilities.RequestMetadata{
		WorkflowID:                    md.WorkflowID,
		WorkflowExecutionID:           md.WorkflowExecutionID,
		WorkflowOwner:                 md.WorkflowOwner,
		OrgID:                         md.OrgID,
		WorkflowName:                  md.WorkflowName,
		WorkflowDonID:                 md.WorkflowDonID,
		WorkflowDonConfigVersion:      md.WorkflowDonConfigVersion,
		ReferenceID:                   md.ReferenceID,
		DecodedWorkflowName:           md.DecodedWorkflowName,
		WorkflowRegistryChainSelector: md.WorkflowRegistryChainSelector,
		WorkflowRegistryAddress:       md.WorkflowRegistryAddress,
		EngineVersion:                 md.EngineVersion,
	}
}

// applyMetadataFields returns a copy of the metadata containing the base
// allowlisted fields plus any optional fields whose per-field feature flag is
// active for the given ExecutionTimestamp.
func applyMetadataFields(ctx context.Context, md capabilities.RequestMetadata, cfg OptInHasherConfig) capabilities.RequestMetadata {
	result := baseMetadataFields(md)
	ts := config.Timestamp(md.ExecutionTimestamp.Unix())

	if cfg.IncludeWorkflowTag != nil {
		if err := cfg.IncludeWorkflowTag.Check(ctx, ts); err == nil {
			result.WorkflowTag = md.WorkflowTag
		}
	}

	return result
}
