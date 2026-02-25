package agent

import (
	"context"
	"errors"
	"testing"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/stretchr/testify/require"
)

type fakeBlockchain struct {
	out *blockchain.Output
}

func (f *fakeBlockchain) ChainSelector() uint64                           { return 0 }
func (f *fakeBlockchain) ChainID() uint64                                 { return 0 }
func (f *fakeBlockchain) ChainFamily() string                             { return "" }
func (f *fakeBlockchain) IsFamily(string) bool                            { return false }
func (f *fakeBlockchain) Fund(context.Context, string, uint64) error      { return nil }
func (f *fakeBlockchain) CtfOutput() *blockchain.Output                   { return f.out }
func (f *fakeBlockchain) ToCldfChain() (cldf_chain.BlockChain, error)     { return nil, nil }

type outputPreferringDeployer struct {
	deployCalls       int
	deployOutputCalls int
}

func (d *outputPreferringDeployer) Deploy(context.Context, *blockchain.Input) (blockchains.Blockchain, error) {
	d.deployCalls++
	return &fakeBlockchain{out: &blockchain.Output{ChainID: "fallback"}}, nil
}

func (d *outputPreferringDeployer) DeployOutput(context.Context, *blockchain.Input) (*blockchain.Output, error) {
	d.deployOutputCalls++
	return &blockchain.Output{ChainID: "1337", Type: blockchain.TypeAnvil}, nil
}

type fallbackOnlyDeployer struct {
	deployCalls int
}

func (d *fallbackOnlyDeployer) Deploy(context.Context, *blockchain.Input) (blockchains.Blockchain, error) {
	d.deployCalls++
	return &fakeBlockchain{
		out: &blockchain.Output{
			ChainID: "2337",
			Type:    blockchain.TypeAnvil,
		},
	}, nil
}

func TestBuildRemoteJDInputEnablesDNSIsolationOverride(t *testing.T) {
	original := &jd.Input{Image: "job-distributor:0.22.1", DisableDNSIsolation: false}

	effective, err := buildRemoteJDInput(original)
	require.NoError(t, err)
	require.NotSame(t, original, effective, "expected a defensive copy")
	require.True(t, effective.DisableDNSIsolation, "remote agent input should force Docker DNS")
	require.False(t, original.DisableDNSIsolation, "original input should remain unchanged")
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

func TestDeployBlockchainComponentPrefersOutputDeployer(t *testing.T) {
	deployer := &outputPreferringDeployer{}
	output, err := DeployBlockchainComponent(
		context.Background(),
		map[blockchain.ChainFamily]blockchains.Deployer{blockchain.FamilyEVM: deployer},
		&blockchain.Input{Type: blockchain.TypeAnvil},
	)
	require.NoError(t, err)
	require.Equal(t, "1337", output.ChainID)
	require.Equal(t, 1, deployer.deployOutputCalls, "DeployOutput should be used when available")
	require.Equal(t, 0, deployer.deployCalls, "Deploy fallback should not be called")
}

func TestDeployBlockchainComponentFallsBackToDeploy(t *testing.T) {
	deployer := &fallbackOnlyDeployer{}
	output, err := DeployBlockchainComponent(
		context.Background(),
		map[blockchain.ChainFamily]blockchains.Deployer{blockchain.FamilyEVM: deployer},
		&blockchain.Input{Type: blockchain.TypeAnvil},
	)
	require.NoError(t, err)
	require.Equal(t, "2337", output.ChainID)
	require.Equal(t, 1, deployer.deployCalls, "Deploy should be called for non-output deployers")
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
		Image:               "job-distributor:0.22.1",
		DisableDNSIsolation: false,
	})
	require.NoError(t, err)
	require.Same(t, expectedOutput, out)
	require.Equal(t, "job-distributor:0.22.1", imageChecked)
	require.NotNil(t, captured)
	require.True(t, captured.DisableDNSIsolation, "remote JD deploy should force Docker DNS")
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
