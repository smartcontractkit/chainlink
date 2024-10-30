package src

import (
	"context"
	"flag"
	"fmt"
	"os"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
)

type provisionForwarderContract struct{}

func NewDeployForwarderCommand() *provisionForwarderContract {
	return &provisionForwarderContract{}
}

func (g *provisionForwarderContract) Name() string {
	return "provision-forwarder-contract"
}

func (g *provisionForwarderContract) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ExitOnError)
	// create flags for all of the env vars then set the env vars to normalize the interface
	// this is a bit of a hack but it's the easiest way to make this work
	ethUrl := fs.String("ethurl", "", "URL of the Ethereum node")
	chainID := fs.Int64("chainid", 1337, "chain ID of the Ethereum network to deploy to")
	accountKey := fs.String("accountkey", "", "private key of the account to deploy from")
	artefactsDir := fs.String("artefacts", defaultArtefactsDir, "Custom artefacts directory location")

	err := fs.Parse(args)

	if err != nil ||
		*ethUrl == "" || ethUrl == nil ||
		*accountKey == "" || accountKey == nil {
		fs.Usage()
		os.Exit(1)
	}

	os.Setenv("ETH_URL", *ethUrl)
	os.Setenv("ETH_CHAIN_ID", fmt.Sprintf("%d", *chainID))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")
	env := helpers.SetupEnv(false)

	deployForwarder(env, *artefactsDir)
}

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
