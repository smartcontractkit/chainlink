package src

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
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
