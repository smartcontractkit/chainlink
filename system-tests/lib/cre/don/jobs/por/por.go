package por

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

// If we wanted to by fancy we could also accept map[JobDescription]string that would get us the job spec
// if there's no job spec for the given JobDescription we would use the standard one, that could be easier
// than having to define the job spec for each JobDescription manually, in case someone wants to change one parameter
func GenerateJobSpecs(input types.GeneratePoRJobSpecsInput) (types.DonJobs, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}
	jobSpecs := make(types.DonJobs)

	chainIDInt, err := strconv.Atoi(input.BlockchainOutput.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert chain ID to int")
	}
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	bootstrapNode, err := node.FindOneWithLabel(input.Don, &ptypes.Label{Key: node.RoleLabelKey, Value: ptr.Ptr(types.BootstrapNode)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find bootstrap node")
	}

	donBootstrapNodePeerID, err := node.ToP2PID(*bootstrapNode, node.KeyExtractingTransformFn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get bootstrap node peer ID")
	}

	var donBootstrapNodeHost string
	for _, label := range bootstrapNode.Labels() {
		if label.Key == node.HostLabelKey {
			donBootstrapNodeHost = *label.Value
			break
		}
	}

	if donBootstrapNodeHost == "" {
		return nil, errors.New("failed to get bootstrap node host from labels")
	}

	// configuration of bootstrap node
	if keystoneflags.HasFlag(input.Flags, types.OCR3Capability) {
		jobSpecs[types.JobDescription{Flag: types.OCR3Capability, NodeType: types.BootstrapNode}] = []*jobv1.ProposeJobRequest{jobs.BootstrapOCR3(bootstrapNode.NodeID, input.OCR3CapabilityAddress, chainIDUint64)}
	}

	// if it's a workflow DON or it has custom compute capability, we need to create a gateway job
	if keystoneflags.HasFlag(input.Flags, types.WorkflowDON) || keystoneflags.HasFlag(input.Flags, types.CustomComputeCapability) {
		jobSpecs[types.JobDescription{Flag: types.WorkflowDON, NodeType: types.BootstrapNode}] = []*jobv1.ProposeJobRequest{jobs.BootstrapGateway(input.Don, chainIDUint64, input.DonID, input.ExtraAllowedPorts, input.ExtraAllowedIPs, input.GatewayConnectorOutput)}
	}

	ocrPeeringData := types.OCRPeeringData{
		OCRBootstraperPeerID: donBootstrapNodePeerID,
		OCRBootstraperHost:   donBootstrapNodeHost,
		Port:                 5001,
	}

	workflowNodeSet, err := node.FindManyWithLabel(input.Don, &ptypes.Label{Key: node.RoleLabelKey, Value: ptr.Ptr(types.WorkerNode)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	// configuration of worker nodes
	for _, node := range workflowNodeSet {
		if keystoneflags.HasFlag(input.Flags, types.CronCapability) {
			jobSpec := jobs.WorkerStandardCapability(node.NodeID, "cron-capabilities", jobs.ExternalCapabilityPath(input.CronCapBinName), jobs.EmptyStdCapConfig)
			jobDesc := types.JobDescription{Flag: types.CronCapability, NodeType: types.WorkerNode}

			if _, ok := jobSpecs[jobDesc]; !ok {
				jobSpecs[jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
			} else {
				jobSpecs[jobDesc] = append(jobSpecs[jobDesc], jobSpec)
			}
		}

		if keystoneflags.HasFlag(input.Flags, types.CustomComputeCapability) {
			config := `"""
				NumWorkers = 3
				[rateLimiter]
				globalRPS = 20.0
				globalBurst = 30
				perSenderRPS = 1.0
				perSenderBurst = 5
				"""`

			jobSpec := jobs.WorkerStandardCapability(node.NodeID, "custom-compute", "__builtin_custom-compute-action", config)
			jobDesc := types.JobDescription{Flag: types.CustomComputeCapability, NodeType: types.WorkerNode}

			if _, ok := jobSpecs[jobDesc]; !ok {
				jobSpecs[jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
			} else {
				jobSpecs[jobDesc] = append(jobSpecs[jobDesc], jobSpec)
			}
		}

		if keystoneflags.HasFlag(input.Flags, types.OCR3Capability) {
			jobSpec := jobs.WorkerOCR3(node.NodeID, input.OCR3CapabilityAddress, common.HexToAddress(node.AccountAddr[chainIDUint64]), node.Ocr2KeyBundleID, ocrPeeringData, chainIDUint64)
			jobDesc := types.JobDescription{Flag: types.OCR3Capability, NodeType: types.WorkerNode}

			if _, ok := jobSpecs[jobDesc]; !ok {
				jobSpecs[jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
			} else {
				jobSpecs[jobDesc] = append(jobSpecs[jobDesc], jobSpec)
			}
		}
	}

	return jobSpecs, nil
}
