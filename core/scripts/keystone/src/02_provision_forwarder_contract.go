package src

import (
	"context"
	"fmt"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
)

func deployForwarder(
	env helpers.Environment,
	artefacts string,
) {
	o := LoadOnchainMeta(artefacts, env)
	if o.Forwarder != nil {
		fmt.Println("Forwarder contract already deployed, skipping")
		return
	}

	fmt.Println("Deploying forwarder contract...")
	forwarderContract := DeployForwarder(env)
	o.Forwarder = forwarderContract
	WriteOnchainMeta(o, artefacts)
}

func DeployForwarder(e helpers.Environment) *forwarder.KeystoneForwarder {
	_, tx, contract, err := forwarder.DeployKeystoneForwarder(e.Owner, e.Ec)
	PanicErr(err)
	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)

	return contract
}
