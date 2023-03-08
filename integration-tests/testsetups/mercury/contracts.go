package mercury

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/stretchr/testify/require"
)

func SetupMercuryContracts(t *testing.T, evmClient blockchain.EVMClient, mercuryRemoteUrl string, feedId [32]byte, ocrConfig contracts.MercuryOCRConfig) (contracts.Verifier, contracts.VerifierProxy, contracts.ReadAccessController, contracts.Exchanger) {
	contractDeployer, err := contracts.NewContractDeployer(evmClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")

	accessController, err := contractDeployer.DeployReadAccessController()
	require.NoError(t, err, "Error deploying ReadAccessController contract")

	// verifierProxy, err := contractDeployer.DeployVerifierProxy(accessController.Address())
	// Use zero address for access controller disables access control
	verifierProxy, err := contractDeployer.DeployVerifierProxy("0x0")
	require.NoError(t, err, "Error deploying VerifierProxy contract")

	verifier, err := contractDeployer.DeployVerifier(verifierProxy.Address())
	require.NoError(t, err, "Error deploying Verifier contract")

	latestConfigDetails, err := verifier.LatestConfigDetails(feedId)
	require.NoError(t, err, "Error getting Verifier.LatestConfigDetails()")
	log.Info().Msgf("Latest config digest: %x", latestConfigDetails.ConfigDigest)
	log.Info().Msgf("Latest config details: %v", latestConfigDetails)

	verifierProxy.InitializeVerifier(latestConfigDetails.ConfigDigest, verifier.Address())

	return verifier, verifierProxy, accessController, nil
}
