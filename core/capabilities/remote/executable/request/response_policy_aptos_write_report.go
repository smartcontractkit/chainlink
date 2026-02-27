package request

import (
	"encoding/hex"
	"fmt"
	"strings"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	aptoscap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/aptos"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

func init() {
	registerResponsePolicyBuilder(newAptosWriteReportPolicy)
}

type aptosWriteReportPolicy struct {
	failedHashQuorumResponses int

	failedHashCount    map[string]int
	failedHashPayload  map[string][]byte
	failedHashMetering map[string][]commoncap.MeteringNodeDetail
	failedTotal        int
	writeSuccessSeen   bool
}

func newAptosWriteReportPolicy(remoteCapabilityInfo commoncap.CapabilityInfo, capabilityMethod string) responsePolicy {
	if capabilityMethod != "WriteReport" {
		return nil
	}
	if !strings.HasPrefix(remoteCapabilityInfo.ID, "aptos:") {
		return nil
	}
	if remoteCapabilityInfo.DON == nil {
		return nil
	}

	failedHashQuorumResponses := min(len(remoteCapabilityInfo.DON.Members), int(2*remoteCapabilityInfo.DON.F+1))
	if failedHashQuorumResponses < 1 {
		failedHashQuorumResponses = 1
	}

	return &aptosWriteReportPolicy{
		failedHashQuorumResponses: failedHashQuorumResponses,
		failedHashCount:           make(map[string]int),
		failedHashPayload:         make(map[string][]byte),
		failedHashMetering:        make(map[string][]commoncap.MeteringNodeDetail),
	}
}

func (p *aptosWriteReportPolicy) ShouldDeferIdenticalResponse(payload []byte) bool {
	reply, ok := decodeAptosWriteReportReply(payload)
	if !ok {
		return false
	}
	return !isAptosWriteSuccess(reply)
}

func (p *aptosWriteReportPolicy) ObserveOKResponse(msg *types.MessageBody, metadata commoncap.ResponseMetadata) {
	reply, ok := decodeAptosWriteReportReply(msg.Payload)
	if !ok {
		return
	}
	if isAptosWriteSuccess(reply) {
		p.writeSuccessSeen = true
		return
	}

	rawHash, ok := getAptosWriteTxHashRaw(reply)
	if !ok || len(rawHash) == 0 {
		return
	}

	normalizedHash, ok := normalizeAptosTxHash(rawHash)
	if !ok {
		return
	}

	p.failedHashCount[normalizedHash]++
	p.failedTotal++

	if _, exists := p.failedHashPayload[normalizedHash]; !exists {
		p.failedHashPayload[normalizedHash] = append([]byte(nil), msg.Payload...)
	}

	if len(metadata.Metering) > 0 {
		p.failedHashMetering[normalizedHash] = append(p.failedHashMetering[normalizedHash], metadata.Metering...)
	}
}

func (p *aptosWriteReportPolicy) BuildDeterministicResponse(allowNoQuorum bool) ([]byte, bool, error) {
	if p.writeSuccessSeen {
		return nil, false, nil
	}
	if len(p.failedHashCount) == 0 {
		return nil, false, nil
	}
	if !allowNoQuorum && p.failedTotal < p.failedHashQuorumResponses {
		return nil, false, nil
	}

	selectedHash := ""
	for hash := range p.failedHashPayload {
		if selectedHash == "" || hash < selectedHash {
			selectedHash = hash
		}
	}

	if selectedHash == "" {
		return nil, false, nil
	}

	selectedPayload, ok := p.failedHashPayload[selectedHash]
	if !ok || len(selectedPayload) == 0 {
		return nil, false, nil
	}

	metadata := commoncap.ResponseMetadata{
		Metering: append([]commoncap.MeteringNodeDetail(nil), p.failedHashMetering[selectedHash]...),
	}

	payload, err := buildAptosFailedWriteReportPayload(selectedPayload, selectedHash, metadata)
	if err != nil {
		return nil, false, err
	}

	return payload, true, nil
}

func (p *aptosWriteReportPolicy) ShouldDeferErrorResponses() bool {
	return true
}

func (p *aptosWriteReportPolicy) FinalizeAfterAllResponses(state responsePolicyState) (*clientResponse, bool, error) {
	if !state.AllResponsesReceived {
		return nil, false, nil
	}

	payload, ok, err := p.BuildDeterministicResponse(true)
	if err != nil {
		return nil, false, fmt.Errorf("failed to build deterministic Aptos failed response: %w", err)
	}
	if ok {
		return &clientResponse{Result: payload}, true, nil
	}

	return &clientResponse{
		Err: fmt.Errorf(
			"received all Aptos write responses without deterministic failed hash (okResponses=%d errorResponses=%d failedHashResponses=%d)",
			state.ResponseVariants,
			state.TotalErrorCount,
			p.failedTotal,
		),
	}, true, nil
}

func decodeAptosWriteReportReply(payload []byte) (*aptoscap.WriteReportReply, bool) {
	resp, err := pb.UnmarshalCapabilityResponse(payload)
	if err != nil {
		return nil, false
	}

	reply := &aptoscap.WriteReportReply{}
	if _, err := commoncap.UnwrapResponse(resp, reply); err != nil {
		return nil, false
	}

	return reply, true
}

func buildAptosFailedWriteReportPayload(payload []byte, normalizedHash string, metadata commoncap.ResponseMetadata) ([]byte, error) {
	resp, err := pb.UnmarshalCapabilityResponse(payload)
	if err != nil {
		return nil, err
	}

	reply := &aptoscap.WriteReportReply{}
	migrated, err := commoncap.UnwrapResponse(resp, reply)
	if err != nil {
		return nil, err
	}

	if err := setAptosWriteFailureStatusAndHash(reply, normalizedHash); err != nil {
		return nil, err
	}

	if err := commoncap.SetResponse(&resp, migrated, reply); err != nil {
		return nil, err
	}

	resp.Metadata = metadata
	return pb.MarshalCapabilityResponse(resp)
}

func isAptosWriteSuccess(reply *aptoscap.WriteReportReply) bool {
	// SUCCESS is enum number 2 in both Aptos proto variants.
	return int32(reply.GetTxStatus()) == 2
}

func getAptosWriteTxHashRaw(reply *aptoscap.WriteReportReply) ([]byte, bool) {
	m := reply.ProtoReflect()
	fd := m.Descriptor().Fields().ByName("tx_hash")
	if fd == nil || !m.Has(fd) {
		return nil, false
	}

	v := m.Get(fd)
	switch fd.Kind() {
	case protoreflect.BytesKind:
		b := v.Bytes()
		if len(b) == 0 {
			return nil, false
		}
		return b, true
	case protoreflect.StringKind:
		s := strings.TrimSpace(v.String())
		if s == "" {
			return nil, false
		}
		return []byte(s), true
	default:
		return nil, false
	}
}

func setAptosWriteFailureStatusAndHash(reply *aptoscap.WriteReportReply, normalizedHash string) error {
	m := reply.ProtoReflect()

	statusFD := m.Descriptor().Fields().ByName("tx_status")
	if statusFD != nil && statusFD.Kind() == protoreflect.EnumKind {
		enumValues := statusFD.Enum().Values()
		var failNum protoreflect.EnumNumber
		if v := enumValues.ByName("TX_STATUS_FAILED"); v != nil {
			failNum = v.Number()
		} else if v := enumValues.ByName("TX_STATUS_ABORTED"); v != nil {
			failNum = v.Number()
		} else if v := enumValues.ByName("TX_STATUS_FATAL"); v != nil {
			failNum = v.Number()
		} else {
			failNum = 0
		}
		m.Set(statusFD, protoreflect.ValueOfEnum(failNum))
	}

	hashFD := m.Descriptor().Fields().ByName("tx_hash")
	if hashFD == nil {
		return nil
	}

	hashWithPrefix := "0x" + normalizedHash
	switch hashFD.Kind() {
	case protoreflect.BytesKind:
		m.Set(hashFD, protoreflect.ValueOfBytes([]byte(hashWithPrefix)))
	case protoreflect.StringKind:
		m.Set(hashFD, protoreflect.ValueOfString(hashWithPrefix))
	}

	return nil
}

func normalizeAptosTxHash(raw []byte) (string, bool) {
	if len(raw) == 32 {
		return hex.EncodeToString(raw), true
	}

	s := strings.TrimSpace(strings.ToLower(string(raw)))
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 64 {
		return "", false
	}

	if _, err := hex.DecodeString(s); err != nil {
		return "", false
	}

	return s, true
}
