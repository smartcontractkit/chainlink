package gateway

import (
	"context"
	"fmt"

	"dario.cat/mergo"

	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
)

var (
	DefaultAllowedPorts = []int{80, 443}
)

type WhitelistConfig struct {
	ExtraAllowedPorts   []int
	ExtraAllowedIPsCIDR []string
}

type JobConfig struct {
	UseServiceCentricFormat bool
	GatewayConfigs          []cre.GatewayConfig
	GatewayServiceConfigs   []cre.GatewayServiceConfig
}

func CreateJobs(ctx context.Context, creEnv *cre.Environment, dons *cre.Dons, gwJobConfig JobConfig, whitelistConfig WhitelistConfig) error {
	specs := make(map[string][]string)

	if !dons.RequiresGateway() {
		return nil
	}

	for _, config := range dons.GatewayConnectors.Configurations {
		gatewayNode, ok := dons.NodeWithUUID(config.NodeUUID)
		if !ok {
			return fmt.Errorf("could not find gateway node with UUID %s in DON topology", config.NodeUUID)
		}

		inputs := job_types.JobSpecInput{
			"useServiceCentricFormat": gwJobConfig.UseServiceCentricFormat,
			"allowedPorts":            append(whitelistConfig.ExtraAllowedPorts, DefaultAllowedPorts...),
			"allowedSchemes":          []string{"http", "https"},
			"allowedIPsCIDR":          whitelistConfig.ExtraAllowedIPsCIDR,
			"gatewayKeyChainSelector": creEnv.RegistryChainSelector,
			"authGatewayID":           config.AuthGatewayID,
		}
		if gwJobConfig.UseServiceCentricFormat {
			inputs["services"] = gwJobConfig.GatewayServiceConfigs
		} else {
			inputs["dons"] = gwJobConfig.GatewayConfigs
		}

		workerInput := cre_jobs.ProposeJobSpecInput{
			Domain:      offchain.ProductLabel,
			Environment: cre.EnvironmentName,
			DONName:     gatewayNode.DON.Name,
			JobName:     "gateway-worker",
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: gatewayNode.DON.Name},
			},
			Template: job_types.Gateway,
			Inputs:   inputs,
		}

		workerVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
		if workerVerErr != nil {
			return fmt.Errorf("precondition verification failed for Custom Compute worker job: %w", workerVerErr)
		}

		workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
		if workerErr != nil {
			return fmt.Errorf("failed to propose Custom Compute worker job spec: %w", workerErr)
		}

		for _, r := range workerReport.Reports {
			out, ok := r.Output.(cre_jobs_ops.ProposeGatewayJobOutput)
			if !ok {
				return fmt.Errorf("unable to cast to ProposeGatewayJobOutput, actual type: %T", r.Output)
			}
			mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
			if mErr != nil {
				return fmt.Errorf("failed to merge worker job specs: %w", mErr)
			}
		}
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve Custom Compute jobs: %w", approveErr)
	}

	return nil
}
