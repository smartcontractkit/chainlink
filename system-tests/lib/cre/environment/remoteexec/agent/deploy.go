package agent

import (
	"context"
	"fmt"
	"strings"

	dockerclient "github.com/moby/moby/client"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

var (
	newJDWithContext              = jd.NewWithContext
	newSharedDBNodeSetWithContext = ns.NewSharedDBNodeSetWithContext
	ensureJDImagePresentFn        = ensureJDImagePresent
)

func DeployBlockchainComponent(
	ctx context.Context,
	deployers map[blockchain.ChainFamily]blockchains.Deployer,
	input *blockchain.Input,
) (*blockchain.Output, error) {
	return blockchains.StartChain(ctx, deployers, input)
}

func DeployJDComponent(ctx context.Context, input *jd.Input) (*jd.Output, error) {
	if input == nil {
		return nil, pkgerrors.New("jd input is nil")
	}
	if err := ensureJDImagePresentFn(ctx, input.Image); err != nil {
		return nil, err
	}

	effectiveInput, err := buildRemoteJDInput(input)
	if err != nil {
		return nil, err
	}
	output, err := newJDWithContext(ctx, effectiveInput)
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
	output, err := newSharedDBNodeSetWithContext(ctx, &inputCopy, registryChain)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to deploy nodeset %s", inputCopy.Name)
	}
	return output, nil
}

func buildRemoteJDInput(input *jd.Input) (*jd.Input, error) {
	jdInput := *input

	return &jdInput, nil
}

func ensureJDImagePresent(ctx context.Context, image string) error {
	if strings.TrimSpace(image) == "" {
		return nil
	}

	client, err := dockerclient.New(dockerclient.FromEnv)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to create docker client for jd image check")
	}
	defer client.Close()

	if _, err := client.ImageInspect(ctx, image); err != nil {
		return fmt.Errorf("jd image %q is not available on remote host; please preload it before starting remote jd", image)
	}
	return nil
}
