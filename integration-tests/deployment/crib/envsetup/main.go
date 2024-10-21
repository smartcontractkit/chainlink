package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/subosito/gotenv"
	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	rootCmd        = &cobra.Command{Use: ""}
	lggr           logger.Logger
	cribConfigPath string
	cribEnv        *deployment.Environment
	cribEnvConfig  devenv.EnvironmentConfig

	// get dir path
	_, b, _, _    = runtime.Caller(0)
	DefaultConfig = filepath.Join(filepath.Dir(b), "config.toml")
	envPath       = filepath.Join(filepath.Dir(b), ".env")
)

func init() {
	var closeLggr func() error
	lggr, closeLggr = logger.NewLogger()
	defer func() {
		err := closeLggr()
		if err != nil {
			panic(err)
		}
	}()
	var err error

	rootCmd.PersistentFlags().StringVarP(&cribConfigPath, "crib-config", "c", DefaultConfig, "CRIB environment configuration file")
	cribEnvConfig, err = devenv.LoadEnvironmentConfig(cribConfigPath)
	if err != nil {
		lggr.Fatalw("failed to load environment configuration", "err", err)
	}
	// read private keys/KMS Keys from env
	mustLoadPrivateKeysFromEnvVar(lggr, &cribEnvConfig)

	if !cribEnvConfig.JDConfig.IsEmpty() {
		cribEnvConfig.JDConfig.Creds = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})
	}
	cribEnv, _, err = devenv.NewEnvironment(context.Background(), lggr, cribEnvConfig)
	if err != nil {
		lggr.Fatalw("failed to create environment", "err", err)
	}

	rootCmd.AddCommand(ccipHomeDeploy)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		lggr.Fatalw("Error executing command", "err", err)
	}
}

func mustLoadPrivateKeysFromEnvVar(lggr logger.Logger, envCfg *devenv.EnvironmentConfig) {
	if _, err := os.Stat(envPath); err == nil || !os.IsNotExist(err) {
		err = gotenv.Load(envPath)
		if err != nil {
			lggr.Fatalw("failed to load .env file", "err", err)
		}
	} else {
		lggr.Warn(`
no .env file found. You need to create one if you don't have env vars set already.

If you want to run with specific private keys, create a .env file specifying private keys.
Example:
# set this if you want to use specific private key for a chain
OWNER_KEY_<CHAIN_ID>="<private-key-for-chain-id>"
# set this if you want to use same private key for all chains
OWNER_KEY="<private-key-for-all-chains>"

If you want to use KMS keys instead of providing private keys, set the following environment variables:
(Please ensure to log into the aws profile that has access to the KMS key)

KMS_DEPLOYER_KEY_ID="<KMS Key Id >"
AWS_PROFILE=<aws profile to be used to access the key>
KMS_DEPLOYER_KEY_REGION=<aws_region>

`)
	}

	for i, chain := range envCfg.Chains {
		key := os.Getenv(fmt.Sprintf("OWNER_KEY_%d", chain.ChainID))
		if key == "" {
			key = os.Getenv("OWNER_KEY")
		}
		var pvtKey *string
		if key != "" {
			pvtKey = &key
		}
		err := envCfg.Chains[i].SetDeployerKey(pvtKey)
		if err != nil {
			lggr.Fatalw("failed to set deployer key", "err", err)
		}
	}
}
