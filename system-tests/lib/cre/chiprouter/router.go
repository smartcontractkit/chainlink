package chiprouter

import (
	"context"
	"net"
	"net/netip"
	"os"
	"strings"
	"sync"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ctfchiprouterclient "github.com/smartcontractkit/chainlink-testing-framework/framework/components/chiprouter/client"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

var (
	clientOnce sync.Once
	clientInst *ctfchiprouterclient.Client
	errClient  error
)

func getClient(ctx context.Context, relativePathToRepoRoot string) (*ctfchiprouterclient.Client, error) {
	clientOnce.Do(func() {
		st, err := envconfig.LoadChipIngressRouterStateFromLocalCRE(relativePathToRepoRoot)
		if err != nil {
			errClient = err
			return
		}
		clientInst, errClient = ctfchiprouterclient.New(ctx, st.AdminURL, st.GRPCURL)
	})

	return clientInst, errClient
}

func EnsureStarted(ctx context.Context, relativePathToRepoRoot, _ string) error {
	_, err := getClient(ctx, relativePathToRepoRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return pkgerrors.New("local CRE state file not found; start the environment first")
		}
		return err
	}

	return nil
}

func RegisterSubscriber(ctx context.Context, relativePathToRepoRoot, name, endpoint string) (string, error) {
	c, err := getClient(ctx, relativePathToRepoRoot)
	if err != nil {
		return "", err
	}

	return c.RegisterSubscriber(ctx, name, normalizeEndpointForRouter(endpoint))
}

func UnregisterSubscriber(ctx context.Context, relativePathToRepoRoot, id string) error {
	if strings.TrimSpace(id) == "" {
		return nil
	}

	c, err := getClient(ctx, relativePathToRepoRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return c.UnregisterSubscriber(ctx, id)
}

func normalizeEndpointForRouter(endpoint string) string {
	host, port, err := net.SplitHostPort(strings.TrimSpace(endpoint))
	if err != nil {
		return endpoint
	}

	if !requiresHostGateway(host) {
		return endpoint
	}

	dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")
	return net.JoinHostPort(dockerHost, port)
}

func requiresHostGateway(host string) bool {
	switch strings.TrimSpace(host) {
	case "", "localhost":
		return true
	}

	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}

	return addr.IsLoopback() || addr.IsUnspecified()
}
