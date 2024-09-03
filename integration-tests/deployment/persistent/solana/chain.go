package solana

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/gauntlet"
	"github.com/smartcontractkit/chainlink/integration-tests/solclient"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"os/exec"
)

type SolanaDockerChain struct {
	image     string
	publicKey string
}

func CreateNewSolanaDockerChain(image, publicKey string) SolanaDockerChain {
	return SolanaDockerChain{
		image:     image,
		publicKey: publicKey,
	}
}

func (s *SolanaDockerChain) Chain() (deployment.SolanaChain, error) {
	env, err := test_env.NewTestEnv()
	if err != nil {
		return deployment.SolanaChain{}, errors.Wrapf(err, "failed to create test environment")
	}
	sol := test_env.NewSolana([]string{env.DockerNetwork.Name}, s.image, s.publicKey)
	err = sol.StartContainer()

	gauntletCopyPath := utils.ProjectRoot + "/gauntlet" + uuid.New().String()
	out, cpErr := exec.Command("cp", "-r", utils.ProjectRoot+"/gauntlet", gauntletCopyPath).Output()
	if cpErr != nil {
		return deployment.SolanaChain{}, errors.Wrap(cpErr, "failed to copy gauntlet folder")
	}
	fmt.Println(string(out))

	sg, err := gauntlet.NewSolanaGauntlet(gauntletCopyPath)
	if err != nil {
		return deployment.SolanaChain{}, errors.Wrap(err, "failed to create gauntlet")
	}

	networkSettings := &solclient.SolNetwork{
		URLs: []string{sol.ExternalWsURL, sol.ExternalHTTPURL},
		Name: "local solana",
	}

	solClient, err := solclient.NewClient(networkSettings)
	if err != nil {
		return deployment.SolanaChain{}, errors.Wrap(err, "failed to create solana client")
	}

	return deployment.SolanaChain{
		Deployer: sg,
		Client:   solClient,
	}, nil
}
