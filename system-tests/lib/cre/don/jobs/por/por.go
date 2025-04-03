package por

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	creflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func GenerateJobSpecs(input *types.GeneratePoRJobSpecsInput,
	customJobsFn func(types.DonJobs, *types.DonWithMetadata) (types.DonJobs, error)) (types.DonsToJobSpecs, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	gatewayConnectorData := input.GatewayConnectorOutput

	// we need to iterate over all DONs to see which need gateway connector and create a map of Don IDs and ETH addresses (which identify nodes that can use the connector)
	// This map will be used to configure the gateway job on the node that runs it. Ccurrently, we support only a single gateway connector, even if CRE supports multiple
	for _, donWithMetadata := range input.DonsWithMetadata {
		// if it's a workflow DON or it has custom compute capability, it needs access to gateway connector
		if creflags.HasFlag(donWithMetadata.Flags, types.WorkflowDON) || creflags.HasFlag(donWithMetadata.Flags, types.CustomComputeCapability) {
			workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.WorkerNode}, node.EqualLabels)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find worker nodes")
			}

			ethAddresses := make([]string, len(workflowNodeSet))
			var ethAddressErr error
			for i, n := range workflowNodeSet {
				ethAddresses[i], ethAddressErr = node.FindLabelValue(n, node.EthAddressKey)
				if ethAddressErr != nil {
					return nil, errors.Wrap(ethAddressErr, "failed to get eth address from labels")
				}
			}
			gatewayConnectorData.Dons = append(gatewayConnectorData.Dons, types.GatewayConnectorDons{
				MembersEthAddresses: ethAddresses,
				ID:                  donWithMetadata.DonMetadata.ID,
			})
		}
	}

	for _, donWithMetadata := range input.DonsWithMetadata {
		jobSpecs, err := generateDonJobSpecs(
			input.BlockchainOutput,
			donWithMetadata,
			input.OCR3CapabilityAddress,
			input.CronCapBinPath,
			input.ExtraAllowedPorts,
			input.ExtraAllowedIPs,
			input.ExtraAllowedIPsCIDR,
			gatewayConnectorData,
			customJobsFn,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate job specs for don %d", donWithMetadata.DonMetadata.ID)
		}

		donToJobSpecs[donWithMetadata.DonMetadata.ID] = jobSpecs
	}

	return donToJobSpecs, nil
}

// If we wanted to by fancy we could also accept map[JobDescription]string that would get us the job spec
// if there's no job spec for the given JobDescription we would use the standard one, that could be easier
// than having to define the job spec for each JobDescription manually, in case someone wants to change one parameter
func generateDonJobSpecs(
	blockchainOutput *blockchain.Output,
	donWithMetadata *types.DonWithMetadata,
	oCR3CapabilityAddress common.Address,
	cronCapBinPath string,
	extraAllowedPorts []int,
	extraAllowedIPs []string,
	extraAllowedIPsCIDR []string,
	gatewayConnectorOutput types.GatewayConnectorOutput,
	customJobsFn func(types.DonJobs, *types.DonWithMetadata) (types.DonJobs, error),
) (types.DonJobs, error) {
	jobSpecs := make(types.DonJobs, 0)

	chainIDInt, err := strconv.Atoi(blockchainOutput.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert chain ID to int")
	}
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	// create job specs for the gateway node
	if creflags.HasFlag(donWithMetadata.Flags, types.GatewayDON) {
		gatewayNode, nodeErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: node.ExtraRolesKey, Value: types.GatewayNode}, node.LabelContains)
		if nodeErr != nil {
			return nil, errors.Wrap(nodeErr, "failed to find bootstrap node")
		}

		gatewayNodeID, gatewayErr := node.FindLabelValue(gatewayNode, node.NodeIDKey)
		if gatewayErr != nil {
			return nil, errors.Wrap(gatewayErr, "failed to get gateway node id from labels")
		}

		jobSpecs = append(jobSpecs, jobs.AnyGateway(gatewayNodeID, chainIDUint64, donWithMetadata.ID, extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, gatewayConnectorOutput))
	}

	// if it's only a gateway node, we don't need to create any other job specs
	if creflags.HasOnlyOneFlag(donWithMetadata.Flags, types.GatewayDON) {
		return jobSpecs, nil
	}

	// look for boostrap node and then for required values in its labels
	bootstrapNode, err := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.BootstrapNode}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find bootstrap node")
	}

	donBootstrapNodePeerID, err := node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get bootstrap node peer ID")
	}

	donBootstrapNodeHost, hostErr := node.FindLabelValue(bootstrapNode, node.HostLabelKey)
	if hostErr != nil {
		return nil, errors.Wrap(hostErr, "failed to get bootstrap node host from labels")
	}

	bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
	if nodeIDErr != nil {
		return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
	}

	// create job specs for the bootstrap node
	if creflags.HasFlag(donWithMetadata.Flags, types.OCR3Capability) {
		jobSpecs = append(jobSpecs, jobs.BootstrapOCR3(bootstrapNodeID, oCR3CapabilityAddress, chainIDUint64))
	}

	ocrPeeringData := types.OCRPeeringData{
		OCRBootstraperPeerID: donBootstrapNodePeerID,
		OCRBootstraperHost:   donBootstrapNodeHost,
		Port:                 5001,
	}

	// create job specs for the worker nodes
	workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.WorkerNode}, node.EqualLabels)
	if err != nil {
		// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	for _, workerNode := range workflowNodeSet {
		nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
		if nodeIDErr != nil {
			return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
		}

		if creflags.HasFlag(donWithMetadata.Flags, types.CronCapability) {
			jobSpecs = append(jobSpecs, jobs.WorkerStandardCapability(nodeID, "cron-capability", cronCapBinPath, jobs.EmptyStdCapConfig))
		}

		if creflags.HasFlag(donWithMetadata.Flags, types.CustomComputeCapability) {
			config := `"""
				NumWorkers = 3
				[rateLimiter]
				globalRPS = 20.0
				globalBurst = 30
				perSenderRPS = 1.0
				perSenderBurst = 5
				"""`
			jobSpecs = append(jobSpecs, jobs.WorkerStandardCapability(nodeID, "custom-compute", "__builtin_custom-compute-action", config))
		}

		if creflags.HasFlag(donWithMetadata.Flags, types.OCR3Capability) {
			nodeEthAddr, ethErr := node.FindLabelValue(workerNode, node.EthAddressKey)
			if ethErr != nil {
				return nil, errors.Wrap(ethErr, "failed to get eth address from labels")
			}

			ocr2KeyBundleID, ocr2Err := node.FindLabelValue(workerNode, node.NodeOCR2KeyBundleIDKey)
			if ocr2Err != nil {
				return nil, errors.Wrap(ocr2Err, "failed to get ocr2 key bundle id from labels")
			}
			jobSpecs = append(jobSpecs, jobs.WorkerOCR3(nodeID, oCR3CapabilityAddress, nodeEthAddr, ocr2KeyBundleID, ocrPeeringData, chainIDUint64))
		}

		// Insert custom jobs, test specific
		if customJobsFn != nil {
			jobSpecs, err = customJobsFn(jobSpecs, donWithMetadata)
			if err != nil {
				return nil, err
			}
		}
	}

	return jobSpecs, nil
}

func WaitForRPCEndpoint(lggr zerolog.Logger, url string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Try immediately first
	client, err := rpc.DialContext(ctx, url)
	if err == nil {
		defer client.Close()
		var blockNumber string
		if err := client.CallContext(ctx, &blockNumber, "eth_blockNumber"); err == nil {
			return nil
		}
	}

	// If immediate check fails, start periodic checks
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for RPC endpoint %s to be available", url)
		case <-ticker.C:
			lggr.Info().Msgf("waiting for %s to become available", url)
			client, err := rpc.DialContext(ctx, url)
			if err != nil {
				continue
			}

			var blockNumber string
			if err := client.CallContext(ctx, &blockNumber, "eth_blockNumber"); err != nil {
				continue
			}

			client.Close()
			// If we get here, the endpoint is responding
			return nil
		}
	}
}
