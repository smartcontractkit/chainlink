package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

type fakeStarterDeployer struct {
	startCalls int
}

func (d *fakeStarterDeployer) Start(context.Context, *blockchain.Input) (*blockchain.Output, error) {
	d.startCalls++
	return &blockchain.Output{ChainID: "1337", Type: blockchain.TypeAnvil}, nil
}

func TestDeployBlockchainComponentNilInputFails(t *testing.T) {
	_, err := DeployBlockchainComponent(context.Background(), nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "blockchain input is nil")
}

func TestDeployBlockchainComponentNoDeployerFails(t *testing.T) {
	_, err := DeployBlockchainComponent(context.Background(), map[blockchain.ChainFamily]blockchains.Deployer{}, &blockchain.Input{Type: blockchain.TypeAnvil})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no deployer found")
}

func TestDeployBlockchainComponentStartsBlockchain(t *testing.T) {
	deployer := &fakeStarterDeployer{}
	output, err := DeployBlockchainComponent(
		context.Background(),
		map[blockchain.ChainFamily]blockchains.Deployer{blockchain.FamilyEVM: deployer},
		&blockchain.Input{Type: blockchain.TypeAnvil},
	)
	require.NoError(t, err)
	require.Equal(t, "1337", output.ChainID)
	require.Equal(t, 1, deployer.startCalls, "expected starter to be called once")
}

func TestDeployJDComponentNilInputFails(t *testing.T) {
	_, err := DeployJDComponent(context.Background(), nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "jd input is nil")
}

func TestDeployJDComponentSuccessUsesSeams(t *testing.T) {
	prevEnsure := ensureJDImagePresentFn
	prevNewJD := newJDWithContext
	t.Cleanup(func() {
		ensureJDImagePresentFn = prevEnsure
		newJDWithContext = prevNewJD
	})

	imageChecked := ""
	ensureJDImagePresentFn = func(_ context.Context, image string) error {
		imageChecked = image
		return nil
	}

	var captured *jd.Input
	expectedOutput := &jd.Output{}
	newJDWithContext = func(_ context.Context, in *jd.Input) (*jd.Output, error) {
		captured = in
		return expectedOutput, nil
	}

	out, err := DeployJDComponent(context.Background(), &jd.Input{
		Image: "job-distributor:0.22.1",
	})
	require.NoError(t, err)
	require.Same(t, expectedOutput, out)
	require.Equal(t, "job-distributor:0.22.1", imageChecked)
	require.NotNil(t, captured)
}

func TestDeployJDComponentImageCheckFailureStopsEarly(t *testing.T) {
	prevEnsure := ensureJDImagePresentFn
	prevNewJD := newJDWithContext
	t.Cleanup(func() {
		ensureJDImagePresentFn = prevEnsure
		newJDWithContext = prevNewJD
	})

	ensureJDImagePresentFn = func(context.Context, string) error {
		return errors.New("image check failed")
	}

	constructorCalled := false
	newJDWithContext = func(context.Context, *jd.Input) (*jd.Output, error) {
		constructorCalled = true
		return &jd.Output{}, nil
	}

	_, err := DeployJDComponent(context.Background(), &jd.Input{Image: "jd:latest"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "image check failed")
	require.False(t, constructorCalled, "jd constructor should not be called when image check fails")
}

func TestDeployJDComponentConstructorFailureIsWrapped(t *testing.T) {
	prevEnsure := ensureJDImagePresentFn
	prevNewJD := newJDWithContext
	t.Cleanup(func() {
		ensureJDImagePresentFn = prevEnsure
		newJDWithContext = prevNewJD
	})

	ensureJDImagePresentFn = func(context.Context, string) error { return nil }
	newJDWithContext = func(context.Context, *jd.Input) (*jd.Output, error) {
		return nil, errors.New("constructor failed")
	}

	_, err := DeployJDComponent(context.Background(), &jd.Input{Image: "jd:latest"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to deploy jd component")
}

func TestDeployNodeSetComponentNilInputsFail(t *testing.T) {
	_, err := DeployNodeSetComponent(context.Background(), nil, &blockchain.Output{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "nodeset input is nil")

	_, err = DeployNodeSetComponent(context.Background(), &ns.Input{Name: "workflow"}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "registry blockchain output is nil")
}

func TestDeployNodeSetComponentSuccessUsesSeam(t *testing.T) {
	prevNewNodeSet := newSharedDBNodeSetWithContext
	t.Cleanup(func() {
		newSharedDBNodeSetWithContext = prevNewNodeSet
	})

	expected := &ns.Output{}
	var capturedInput *ns.Input
	var capturedRegistry *blockchain.Output
	newSharedDBNodeSetWithContext = func(_ context.Context, in *ns.Input, registry *blockchain.Output) (*ns.Output, error) {
		capturedInput = in
		capturedRegistry = registry
		return expected, nil
	}

	registry := &blockchain.Output{ChainID: "1337"}
	input := &ns.Input{Name: "workflow"}
	out, err := DeployNodeSetComponent(context.Background(), input, registry)
	require.NoError(t, err)
	require.Same(t, expected, out)
	require.NotNil(t, capturedInput)
	require.Equal(t, "workflow", capturedInput.Name)
	require.Same(t, registry, capturedRegistry)
}

func TestDeployNodeSetComponentConstructorFailureIsWrapped(t *testing.T) {
	prevNewNodeSet := newSharedDBNodeSetWithContext
	t.Cleanup(func() {
		newSharedDBNodeSetWithContext = prevNewNodeSet
	})

	newSharedDBNodeSetWithContext = func(context.Context, *ns.Input, *blockchain.Output) (*ns.Output, error) {
		return nil, errors.New("nodeset constructor failed")
	}

	_, err := DeployNodeSetComponent(
		context.Background(),
		&ns.Input{Name: "workflow"},
		&blockchain.Output{ChainID: "1337"},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to deploy nodeset workflow")
}
