package creregistry

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	registryclient "github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry/client"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry/registrytest"
	regserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

// newShim builds a shim over the real registry, with the shim's
// per-capability listeners and the client's dials both resolved through one
// in-memory address book. That keeps the shim's real behaviour under test — one
// listener and one address per capability — without needing permission to bind a
// port, which some sandboxes withhold and which otherwise turns into silently
// skipped tests.
func newShim(t *testing.T, reg *regserver.Registry) *LOOPShim {
	t.Helper()

	book := registrytest.NewAddrBook()
	client := registryclient.New(logger.Test(t), registrytest.Serve(t, reg),
		book.DialOption(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	shim := NewLOOPShim(logger.Test(t), client, "")
	shim.listen = func(string) (net.Listener, error) { return book.Listen(), nil }
	t.Cleanup(func() { _ = shim.Close() })
	return shim
}

// fakeExecutable is an in-process capability, i.e. what a LOOP still hands over.
type fakeExecutable struct {
	info     capabilities.CapabilityInfo
	response capabilities.CapabilityResponse
}

func (f *fakeExecutable) Info(context.Context) (capabilities.CapabilityInfo, error) {
	return f.info, nil
}

func (f *fakeExecutable) Execute(context.Context, capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	return f.response, nil
}

func (f *fakeExecutable) RegisterToWorkflow(context.Context, capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (f *fakeExecutable) UnregisterFromWorkflow(context.Context, capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

func TestLOOPShim_AddServesCapabilityAndAnnouncesItsAddress(t *testing.T) {
	ctx := context.Background()

	outputs, err := values.NewMap(map[string]any{"out": int64(5)})
	require.NoError(t, err)

	reg := regserver.New(logger.Test(t))
	shim := newShim(t, reg)

	impl := &fakeExecutable{
		info:     capabilities.MustNewCapabilityInfo("act@1.0.0", capabilities.CapabilityTypeAction, "act"),
		response: capabilities.CapabilityResponse{Value: outputs},
	}
	require.NoError(t, shim.Add(ctx, impl))

	require.Len(t, reg.List(t.Context()), 1)
	addr := addressOf(t, reg, "act@1.0.0")
	require.NotEmpty(t, addr, "the shim must announce a concrete address, not an empty one")

	// The announced address must actually be serving the capability: that is the
	// contract the registry relies on when it dials back. Resolving through the
	// shim exercises the whole path — registry lookup, dial, call.
	exec, err := shim.GetExecutable(ctx, "act@1.0.0")
	require.NoError(t, err)

	gotInfo, err := exec.Info(ctx)
	require.NoError(t, err)
	assert.Equal(t, "act@1.0.0", gotInfo.ID)

	resp, err := exec.Execute(ctx, capabilities.CapabilityRequest{})
	require.NoError(t, err)
	assert.Equal(t, outputs, resp.Value)
}

func TestLOOPShim_AddUsesADistinctAddressPerCapability(t *testing.T) {
	ctx := context.Background()

	reg := regserver.New(logger.Test(t))
	shim := newShim(t, reg)

	// Core hosts many capabilities arriving one at a time, which is exactly why it
	// cannot use the single-address client: each needs its own listener.
	for _, id := range []string{"a@1.0.0", "b@1.0.0"} {
		require.NoError(t, shim.Add(ctx, &fakeExecutable{
			info: capabilities.MustNewCapabilityInfo(id, capabilities.CapabilityTypeAction, id),
		}))
	}

	require.Len(t, reg.List(t.Context()), 2)
	addrA, addrB := addressOf(t, reg, "a@1.0.0"), addressOf(t, reg, "b@1.0.0")
	require.NotEmpty(t, addrA)
	require.NotEmpty(t, addrB)
	assert.NotEqual(t, addrA, addrB)
}

func TestLOOPShim_AddRejectsLiveDuplicate(t *testing.T) {
	ctx := context.Background()

	reg := regserver.New(logger.Test(t))
	shim := newShim(t, reg)

	impl := &fakeExecutable{
		info: capabilities.MustNewCapabilityInfo("act@1.0.0", capabilities.CapabilityTypeAction, "act"),
	}
	require.NoError(t, shim.Add(ctx, impl))

	// A capability that cannot report a connection state counts as live, so a
	// second claimant on the same ID is refused rather than silently displacing it.
	err := shim.Add(ctx, impl)
	require.Error(t, err)
	require.ErrorIs(t, err, registry.ErrCapabilityAlreadyExists)
	assert.Len(t, reg.List(t.Context()), 1)
}

// stateCapability reports a connection state, as a capability backed by a gRPC
// connection does.
type stateCapability struct {
	*fakeExecutable
	state connectivity.State
}

func (s *stateCapability) GetState() connectivity.State { return s.state }

func TestLOOPShim_AddReplacesDeadRegistrationOnRestart(t *testing.T) {
	ctx := context.Background()

	reg := regserver.New(logger.Test(t))
	shim := newShim(t, reg)

	info := capabilities.MustNewCapabilityInfo("act@1.0.0", capabilities.CapabilityTypeAction, "act")

	// A host registers, then dies: its connection goes to Shutdown.
	dead := &stateCapability{fakeExecutable: &fakeExecutable{info: info}, state: connectivity.Shutdown}
	require.NoError(t, shim.Add(ctx, dead))
	firstAddr := addressOf(t, reg, "act@1.0.0")
	require.NotEmpty(t, firstAddr)

	// It restarts and re-registers. Refusing here would lose the capability until
	// core itself restarted, so the dead registration must be replaced.
	fresh := &stateCapability{fakeExecutable: &fakeExecutable{info: info}, state: connectivity.Ready}
	require.NoError(t, shim.Add(ctx, fresh))

	require.Len(t, reg.List(t.Context()), 1)
	secondAddr := addressOf(t, reg, "act@1.0.0")
	assert.NotEqual(t, firstAddr, secondAddr, "the replacement is served at its own address")

	// Only the replacement is retained, and it is reachable.
	assert.Len(t, shim.served, 1)
	exec, err := shim.GetExecutable(ctx, "act@1.0.0")
	require.NoError(t, err)
	_, err = exec.Info(ctx)
	require.NoError(t, err)
}

func TestLOOPShim_AddRollsBackWhenRegistrationFails(t *testing.T) {
	ctx := context.Background()

	shim := newShimWithUnreachableRegistry(t)

	impl := &fakeExecutable{
		info: capabilities.MustNewCapabilityInfo("act@1.0.0", capabilities.CapabilityTypeAction, "act"),
	}
	require.Error(t, shim.Add(ctx, impl))

	// Nothing will ever dial the listener, so the shim must not hold it as served;
	// a retry has to be able to start over.
	err := shim.Add(ctx, impl)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "already served",
		"a failed registration must not leave the capability marked as served")
}

func TestLOOPShim_AddRejectsTypeMismatch(t *testing.T) {
	ctx := context.Background()

	reg := regserver.New(logger.Test(t))
	shim := newShim(t, reg)

	// Declared a trigger but only implements the executable surface: this must be
	// caught at registration, not when a workflow first subscribes.
	impl := &fakeExecutable{
		info: capabilities.MustNewCapabilityInfo("act@1.0.0", capabilities.CapabilityTypeTrigger, "act"),
	}
	require.Error(t, shim.Add(ctx, impl))
	assert.Empty(t, reg.List(t.Context()))
}

func TestLOOPShim_RemoveStopsServing(t *testing.T) {
	ctx := context.Background()

	reg := regserver.New(logger.Test(t))
	shim := newShim(t, reg)

	impl := &fakeExecutable{
		info: capabilities.MustNewCapabilityInfo("act@1.0.0", capabilities.CapabilityTypeAction, "act"),
	}
	require.NoError(t, shim.Add(ctx, impl))
	addr := addressOf(t, reg, "act@1.0.0")

	require.NoError(t, shim.Remove(ctx, "act@1.0.0"))
	_, getErr := reg.Get(t.Context(), "act@1.0.0")
	assert.Error(t, getErr, "Remove must reach the registry")

	// The listener is gone, so re-adding under the same ID must work.
	require.NoError(t, shim.Add(ctx, impl))
	assert.NotEqual(t, addr, addressOf(t, reg, "act@1.0.0"),
		"a re-added capability gets a fresh listener")
}

func TestLOOPShim_ClosePreventsFurtherAdds(t *testing.T) {
	ctx := context.Background()

	reg := regserver.New(logger.Test(t))
	shim := newShim(t, reg)

	require.NoError(t, shim.Close())

	err := shim.Add(ctx, &fakeExecutable{
		info: capabilities.MustNewCapabilityInfo("act@1.0.0", capabilities.CapabilityTypeAction, "act"),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "closed")
	assert.Empty(t, reg.List(t.Context()))
}

type addressableStateCapability struct {
	*fakeExecutable
	addr  string
	state connectivity.State
}

func (a *addressableStateCapability) CapabilityAddress() string    { return a.addr }
func (a *addressableStateCapability) GetState() connectivity.State { return a.state }

// newShimWithUnreachableRegistry builds a shim whose registry connection goes
// nowhere, standing in for the proxy being down. Registration therefore fails at
// the registry rather than in the shim, which is what the rollback path is for.
func newShimWithUnreachableRegistry(t *testing.T) *LOOPShim {
	t.Helper()

	book := registrytest.NewAddrBook()

	conn, err := grpc.NewClient("passthrough:///dead",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return nil, errors.New("registry down")
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client := registryclient.New(logger.Test(t), conn,
		book.DialOption(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	shim := NewLOOPShim(logger.Test(t), client, "")
	shim.listen = func(string) (net.Listener, error) { return book.Listen(), nil }
	t.Cleanup(func() { _ = shim.Close() })
	return shim
}

// addressOf returns the address the registry currently has for id, or "" if it
// holds none.
func addressOf(t *testing.T, reg *regserver.Registry, id string) string {
	t.Helper()

	h, err := reg.Get(t.Context(), id)
	if err != nil {
		return ""
	}
	return h.URL
}
