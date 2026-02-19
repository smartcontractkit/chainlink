package agent

import (
	"context"
	"fmt"
	"strings"

	dockerclient "github.com/docker/docker/client"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

type OutputDeployer interface {
	DeployOutput(ctx context.Context, input *blockchain.Input) (*blockchain.Output, error)
}

func DeployBlockchainComponent(
	ctx context.Context,
	deployers map[blockchain.ChainFamily]blockchains.Deployer,
	input *blockchain.Input,
) (*blockchain.Output, error) {
	if input == nil {
		return nil, pkgerrors.New("blockchain input is nil")
	}

	chainFamily, err := blockchain.TypeToFamily(input.Type)
	if err != nil {
		return nil, err
	}

	deployer, ok := deployers[chainFamily]
	if !ok {
		return nil, fmt.Errorf("no deployer found for blockchain type %s", input.Type)
	}

	if outputDeployer, ok := deployer.(OutputDeployer); ok {
		deployedOutput, err := outputDeployer.DeployOutput(ctx, input)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain output of type %s", input.Type)
		}
		return deployedOutput, nil
	}

	deployed, err := deployer.Deploy(ctx, input)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain of type %s", input.Type)
	}

	return deployed.CtfOutput(), nil
}

func DeployJDComponent(ctx context.Context, input *jd.Input) (*jd.Output, error) {
	if input == nil {
		return nil, pkgerrors.New("jd input is nil")
	}
	if err := ensureJDImagePresent(ctx, input.Image); err != nil {
		return nil, err
	}

	effectiveInput, err := buildRemoteJDInput(input)
	if err != nil {
		return nil, err
	}
	output, err := jd.NewWithContext(ctx, effectiveInput)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to deploy jd component")
	}
	return output, nil
}

func DeployNodeSetComponent(ctx context.Context, input *ns.Input, registryChain *blockchain.Output) (*ns.Output, error) {
	if input == nil {
		return nil, pkgerrors.New("nodeset input is nil")
	}
	if registryChain == nil {
		return nil, pkgerrors.New("registry blockchain output is nil")
	}
	inputCopy := *input
	output, err := ns.NewSharedDBNodeSetWithContext(ctx, &inputCopy, registryChain)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to deploy nodeset %s", inputCopy.Name)
	}
	return output, nil
}

func buildRemoteJDInput(input *jd.Input) (*jd.Input, error) {
	jdInput := *input
	// Remote agent deployments require Docker service discovery (jd -> jd-db),
	// so keep Docker embedded DNS instead of isolated localhost DNS.
	jdInput.DisableDNSIsolation = true

	return &jdInput, nil
}

func ensureJDImagePresent(ctx context.Context, image string) error {
	if strings.TrimSpace(image) == "" {
		return nil
	}

	client, err := dockerclient.NewClientWithOpts(dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		return pkgerrors.Wrap(err, "failed to create docker client for jd image check")
	}
	defer client.Close()

	if _, err := client.ImageInspect(ctx, image); err != nil {
		return fmt.Errorf("jd image %q is not available on remote host; please preload it before starting remote jd", image)
	}
	return nil
}
