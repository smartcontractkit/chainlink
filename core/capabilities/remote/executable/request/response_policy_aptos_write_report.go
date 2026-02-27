package request

import (
	"encoding/hex"
	"fmt"
	"strings"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	aptoscap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/aptos"

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
	return reply.GetTxStatus() == aptoscap.TxStatus_TX_STATUS_FAILED
}

func (p *aptosWriteReportPolicy) ObserveOKResponse(msg *types.MessageBody, metadata commoncap.ResponseMetadata) {
	reply, ok := decodeAptosWriteReportReply(msg.Payload)
	if !ok {
		return
	}
	if reply.GetTxStatus() == aptoscap.TxStatus_TX_STATUS_SUCCESS {
		p.writeSuccessSeen = true
		return
	}
	if reply.GetTxStatus() != aptoscap.TxStatus_TX_STATUS_FAILED || len(reply.GetTxHash()) == 0 {
		return
	}

	normalizedHash, ok := normalizeAptosTxHash(reply.GetTxHash())
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

	reply.TxStatus = aptoscap.TxStatus_TX_STATUS_FAILED
	reply.TxHash = []byte("0x" + normalizedHash)

	if err := commoncap.SetResponse(&resp, migrated, reply); err != nil {
		return nil, err
	}

	resp.Metadata = metadata
	return pb.MarshalCapabilityResponse(resp)
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
