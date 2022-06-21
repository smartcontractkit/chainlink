package main

import (
	"flag"
	"os"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func main() {
	e := helpers.SetupEnv()

	switch os.Args[1] {

	case "dkg-deploy":
		DeployDKG(e)

	case "vrf-deploy":
		cmd := flag.NewFlagSet("vrf-deploy", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "dkg contract address")
		keyId := cmd.String("key-id", "", "key ID")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id")
		DeployVRF(e, *dkgAddress, *keyId)

	case "dkg-add-client":
		cmd := flag.NewFlagSet("dkg-add-client", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyId := cmd.String("key-id", "", "key ID")
		clientAddress := cmd.String("client-address", "", "client address")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id", "client-address")
		AddClientToDKG(e, *dkgAddress, *keyId, *clientAddress)

	case "dkg-remove-client":
		cmd := flag.NewFlagSet("dkg-add-client", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyId := cmd.String("key-id", "", "key ID")
		clientAddress := cmd.String("client-address", "", "client address")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id", "client-address")
		RemoveClientFromDKG(e, *dkgAddress, *keyId, *clientAddress)

	case "dkg-set-config":
		cmd := flag.NewFlagSet("dkg-set-config", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		signers := cmd.String("signers", "", "comma-separated list of signers")
		transmitters := cmd.String("transmitters", "", "comma-separated list transmitters")
		f := cmd.Uint("f", 0, "number of faulty oracles")
		onchainConfig := cmd.String("onchain-config", "", "on-chain contract configuration")
		offchainConfig := cmd.String("offchain-config", "", "off-chain contract configuration")
		offchainConfigVersion := cmd.Uint("offchain-config-version", 0, "version number of offchain config schema")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "signers", "transmitters", "f", "onchain-config", "offchain-config", "offchain-config-version")

		SetDKGConfig(e, *dkgAddress, *signers, *transmitters, uint8(*f), *onchainConfig, *offchainConfig, uint64(*offchainConfigVersion))

	case "vrf-set-config":
		cmd := flag.NewFlagSet("vrf-set-config", flag.ExitOnError)
		vrfAddress := cmd.String("vrf-address", "", "VRF contract address")
		signers := cmd.String("signers", "", "comma-separated list of signers")
		transmitters := cmd.String("transmitters", "", "comma-separated list of transmitters")
		f := cmd.Uint("f", 0, "number of faulty oracles")
		onchainConfig := cmd.String("onchain-config", "", "on-chain contract configuration")
		offchainConfig := cmd.String("offchain-config", "", "off-chain contract configuration")
		offchainConfigVersion := cmd.Uint("offchain-config-version", 0, "version number of offchain config schema")
		helpers.ParseArgs(cmd, os.Args[2:], "vrf-address", "signers", "transmitters", "f", "onchain-config", "offchain-config", "offchain-config-version")

		SetVRFConfig(e, *vrfAddress, *signers, *transmitters, uint8(*f), *onchainConfig, *offchainConfig, uint64(*offchainConfigVersion))

	default:
		panic("unrecognized subcommand: " + os.Args[1])
	}
}
