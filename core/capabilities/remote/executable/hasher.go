package executable

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

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

// V2 Capabilities (Executables) default to a simple hasher that hashes the entire payload.
// WriteReport methods use a hasher that excludes signatures from the WriteReportRequest.
// Additional hashers can be added here as needed.
type simpleHasher struct {
}

func (r *simpleHasher) Hash(ctx context.Context, msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}

	// Exclude per-node-divergent metadata fields to ensure identical requests
	// with different values produce the same hash
	req.Metadata.SpendLimits = nil
	req.Metadata.ExecutionTimestamp = time.Time{}

	reqBytes, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal capability request: %w", err)
	}
	hash := sha256.Sum256(reqBytes)
	return hash, nil
}

func NewSimpleHasher() types.MessageHasher {
	return &simpleHasher{}
}

type writeReportExcludeSignaturesHasher struct {
}

func (r *writeReportExcludeSignaturesHasher) Hash(ctx context.Context, msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}
	if req.Payload == nil {
		return [32]byte{}, errors.New("capability request payload is nil")
	}

	// Exclude per-node-divergent metadata fields to ensure identical requests
	// with different values produce the same hash
	req.Metadata.SpendLimits = nil
	req.Metadata.ExecutionTimestamp = time.Time{}
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

func NewWriteReportExcludeSignaturesHasher() types.MessageHasher {
	return &writeReportExcludeSignaturesHasher{}
}

// OptInHasherConfig holds per-field feature flags that control which optional
// metadata fields are included in the request hash beyond the base allowlist.
// Currently empty — all fields in the base allowlist are always included.
//
// To opt in a new field in the future:
//  1. Add a feature flag (Setting[Range[config.Timestamp]]) in cresettings,
//     named FeatureRequestHashInclude<Field>ActivePeriod.
//  2. Add a field of type limits.RangeLimiter[config.Timestamp] to this struct.
//  3. Add the conditional copy in applyOptInMetadata.
//  4. Construct the limiter in launcher.NewLauncher and pass it via OptInHasherConfig.
//
// Each flag is evaluated against the request's DON-derived ExecutionTimestamp,
// so all nodes processing the same request include or exclude the field
// atomically — no cross-version requestID divergence.
type OptInHasherConfig struct{}

// baseOptInMetadataFields returns a copy of the metadata containing only the
// base allowlisted fields. New metadata fields added in the future are excluded
// from the hash by default, preventing cross-version requestID divergence.
func baseOptInMetadataFields(md capabilities.RequestMetadata) capabilities.RequestMetadata {
	return capabilities.RequestMetadata{
		WorkflowID:               md.WorkflowID,
		WorkflowExecutionID:      md.WorkflowExecutionID,
		WorkflowOwner:            md.WorkflowOwner,
		OrgID:                    md.OrgID,
		WorkflowName:             md.WorkflowName,
		WorkflowDonID:            md.WorkflowDonID,
		WorkflowDonConfigVersion: md.WorkflowDonConfigVersion,
		ReferenceID:              md.ReferenceID,
		DecodedWorkflowName:      md.DecodedWorkflowName,
		SpendLimits:              md.SpendLimits,
	}
}

// applyOptInMetadata returns a copy of the metadata containing the base
// allowlisted fields plus any optional fields whose per-field feature flag is
// active for the given ExecutionTimestamp.
func applyOptInMetadata(ctx context.Context, md capabilities.RequestMetadata, cfg OptInHasherConfig) capabilities.RequestMetadata {
	result := baseOptInMetadataFields(md)
	// Future optional fields: check cfg.Include<Field> against
	// config.Timestamp(md.ExecutionTimestamp.Unix()) and copy the field
	// from md into result when the flag is active.
	_ = ctx
	return result
}

// optInHasher hashes only an explicit allowlist of metadata fields rather than
// including everything and excluding specific fields. This prevents new metadata
// fields (e.g. WorkflowTag) from changing the requestID when nodes run different
// versions of the protobuf definitions. Optional fields beyond the base allowlist
// are gated by per-field feature flags in OptInHasherConfig.
type optInHasher struct {
	cfg OptInHasherConfig
}

func (r *optInHasher) Hash(ctx context.Context, msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}

	req.Metadata = applyOptInMetadata(ctx, req.Metadata, r.cfg)

	reqBytes, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal capability request: %w", err)
	}
	return sha256.Sum256(reqBytes), nil
}

func NewOptInHasher(cfg OptInHasherConfig) types.MessageHasher {
	return &optInHasher{cfg: cfg}
}

// optInWriteReportExcludeSignaturesHasher combines the metadata opt-in allowlist
// with WriteReport-specific signature exclusion, mirroring writeReportExcludeSignaturesHasher.
type optInWriteReportExcludeSignaturesHasher struct {
	cfg OptInHasherConfig
}

func (r *optInWriteReportExcludeSignaturesHasher) Hash(ctx context.Context, msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}
	if req.Payload == nil {
		return [32]byte{}, errors.New("capability request payload is nil")
	}

	req.Metadata = applyOptInMetadata(ctx, req.Metadata, r.cfg)

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
		wrReq.Report.Sigs = nil
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
		wrReq.Report.Sigs = nil
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
		wrReq.Report.Sigs = nil
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
		wrReq.Report.Sigs = nil
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

func NewOptInWriteReportExcludeSignaturesHasher(cfg OptInHasherConfig) types.MessageHasher {
	return &optInWriteReportExcludeSignaturesHasher{cfg: cfg}
}

// featureFlagHasher is a composite hasher that delegates to either an opt-out
// baseHasher (backward-compatible behaviour) or an optInHasher based on whether
// the ExecutionTimestamp falls inside the feature flag's active window.
//
// The flag is evaluated against the DON-derived ExecutionTimestamp from the
// request metadata, not the wall clock. Because all nodes processing the same
// request share the same DON time, the flag evaluation is deterministic across
// nodes — they all switch to the opt-in hasher atomically.
type featureFlagHasher struct {
	baseHasher  types.MessageHasher
	optInHasher types.MessageHasher
	flag        limits.RangeLimiter[config.Timestamp]
}

func (h *featureFlagHasher) Hash(ctx context.Context, msg *types.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}

	// Check the flag using the DON-derived ExecutionTimestamp. When the
	// timestamp is zero (DON time not enabled), the check fails and the
	// base hasher is used — preserving backward compatibility.
	if err := h.flag.Check(ctx, config.Timestamp(req.Metadata.ExecutionTimestamp.Unix())); err == nil {
		return h.optInHasher.Hash(ctx, msg)
	}
	return h.baseHasher.Hash(ctx, msg)
}

func NewFeatureFlagHasher(base, optIn types.MessageHasher, flag limits.RangeLimiter[config.Timestamp]) types.MessageHasher {
	return &featureFlagHasher{baseHasher: base, optInHasher: optIn, flag: flag}
}
