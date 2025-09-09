package tron

import (
	"context"
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/address"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
)

const ForwarderContract datastore.ContractType = "KeystoneForwarder"

const (
	DeploymentBlockLabel = "deployment-block"
	DeploymentHashLabel  = "deployment-hash"
)

type DeployTronResponse struct {
	Address address.Address
	Tx      string
	Tv      cldf.TypeAndVersion
}

func DeployKeystoneForwarder(chain cldf_tron.Chain, deployOptions *cldf_tron.DeployOptions, labels []string) (*DeployTronResponse, error) {
	forwarderAddress, txInfo, err := chain.DeployContractAndConfirm(context.Background(), ForwarderContract.String(), forwarder.KeystoneForwarderABI, forwarder.KeystoneForwarderBin, nil, deployOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm KeystoneForwarder: %+v, %w", txInfo, err)
	}

	forwarderResponse, err := chain.Client.TriggerConstantContract(chain.Address, forwarderAddress, "typeAndVersion()", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version from %s: %w", forwarderAddress, err)
	}

	typeAndVersion, err := changeset.ExtractTypeAndVersion(forwarderResponse.ConstantResult[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode type and version from %s: %w", forwarderResponse.ConstantResult[0], err)
	}

	tv, err := cldf.TypeAndVersionFromString(typeAndVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", typeAndVersion, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	tv.Labels.Add(fmt.Sprintf("%s: %s", DeploymentHashLabel, txInfo.ID))
	tv.Labels.Add(fmt.Sprintf("%s: %d", DeploymentBlockLabel, txInfo.BlockNumber))

	resp := &DeployTronResponse{
		Address: forwarderAddress,
		Tx:      txInfo.ID,
		Tv:      tv,
	}
	return resp, nil
}
