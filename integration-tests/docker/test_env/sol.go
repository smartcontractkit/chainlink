package test_env

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const (
	SolHTTPPort = "8899"
	SolWSPort   = "8900"
)

var configYmlRaw = `
json_rpc_url: http://0.0.0.0:8899
websocket_url: ws://0.0.0.0:8900
keypair_path: /root/.config/solana/cli/id.json
address_labels:
  "11111111111111111111111111111111": ""
commitment: finalized
`

var idJSONRaw = `
[94,214,238,83,144,226,75,151,226,20,5,188,42,110,64,180,196,244,6,199,29,231,108,112,67,175,110,182,3,242,102,83,103,72,221,132,137,219,215,192,224,17,146,227,94,4,173,67,173,207,11,239,127,174,101,204,65,225,90,88,224,45,205,117]
`

type Solana struct {
	test_env.EnvComponent
	ExternalHTTPURL string
	ExternalWsURL   string
	InternalHTTPURL string
	InternalWsURL   string
	t               *testing.T
	l               zerolog.Logger
	Image           string
	PublicKey       string
}

func NewSolana(networks []string, devnetImage string, publicKey string, opts ...test_env.EnvComponentOption) *Solana {
	ms := &Solana{
		EnvComponent: test_env.EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "solana", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		l:         log.Logger,
		Image:     devnetImage,
		PublicKey: publicKey,
	}
	for _, opt := range opts {
		opt(&ms.EnvComponent)
	}
	return ms
}

func (s *Solana) WithTestLogger(t *testing.T) *Solana {
	s.l = logging.GetTestLogger(t)
	s.t = t
	return s
}

func (s *Solana) StartContainer() error {
	l := tc.Logger
	if s.t != nil {
		l = logging.CustomT{
			T: s.t,
			L: s.l,
		}
	}

	// get disabled/unreleased features on mainnet
	inactiveMainnetFeatures, err := GetInactiveFeatureHashes("mainnet-beta")
	if err != nil {
		return err
	}

	cReq, err := s.getContainerRequest(inactiveMainnetFeatures)
	if err != nil {
		return err
	}
	c, err := tc.GenericContainer(testcontext.Get(s.t), tc.GenericContainerRequest{
		ContainerRequest: *cReq,
		Reuse:            true,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("cannot start Solana container: %w", err)
	}
	s.Container = c
	host, err := test_env.GetHost(testcontext.Get(s.t), c)
	if err != nil {
		return err
	}
	httpPort, err := c.MappedPort(testcontext.Get(s.t), test_env.NatPort(SolHTTPPort))
	if err != nil {
		return err
	}
	wsPort, err := c.MappedPort(testcontext.Get(s.t), test_env.NatPort(SolWSPort))
	if err != nil {
		return err
	}
	s.ExternalHTTPURL = fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	s.InternalHTTPURL = fmt.Sprintf("http://%s:%s", s.ContainerName, SolHTTPPort)
	s.ExternalWsURL = fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	s.InternalWsURL = fmt.Sprintf("ws://%s:%s", s.ContainerName, SolWSPort)

	s.l.Info().
		Any("ExternalHTTPURL", s.ExternalHTTPURL).
		Any("InternalHTTPURL", s.InternalHTTPURL).
		Any("ExternalWsURL", s.ExternalWsURL).
		Any("InternalWsURL", s.InternalWsURL).
		Str("containerName", s.ContainerName).
		Msgf("Started Solana container")

	// validate features are properly set
	inactiveLocalFeatures, err := GetInactiveFeatureHashes(s.ExternalHTTPURL)
	if err != nil {
		return err
	}
	if !slices.Equal(inactiveMainnetFeatures, inactiveLocalFeatures) {
		return fmt.Errorf("Localnet features does not match mainnet features")
	}
	return nil
}

func (s *Solana) getContainerRequest(inactiveFeatures InactiveFeatures) (*tc.ContainerRequest, error) {
	configYml, err := os.CreateTemp("", "config.yml")
	if err != nil {
		return nil, err
	}
	_, err = configYml.WriteString(configYmlRaw)
	if err != nil {
		return nil, err
	}

	idJSON, err := os.CreateTemp("", "id.json")
	if err != nil {
		return nil, err
	}
	_, err = idJSON.WriteString(idJSONRaw)
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:         s.ContainerName,
		Image:        s.Image,
		ExposedPorts: []string{test_env.NatPortFormat(SolHTTPPort), test_env.NatPortFormat(SolWSPort)},
		Env: map[string]string{
			"SERVER_PORT": "1080",
		},
		Networks: s.Networks,
		WaitingFor: tcwait.ForLog("Processed Slot: 1").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   utils.ContractsDir,
				Target:   "/programs",
				ReadOnly: false,
			})
		},
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: []tc.ContainerHook{
					func(ctx context.Context, container tc.Container) error {
						err = container.CopyFileToContainer(ctx, configYml.Name(), "/root/.config/solana/cli/config.yml", 0644)
						if err != nil {
							return err
						}
						err = container.CopyFileToContainer(ctx, idJSON.Name(), "/root/.config/solana/cli/id.json", 0644)
						return err
					},
				},
			},
		},
		Entrypoint: []string{"sh", "-c", "mkdir -p /root/.config/solana/cli && solana-test-validator -r --mint=" + s.PublicKey + " " + inactiveFeatures.CLIString()},
	}, nil
}

type FeatureStatuses struct {
	Features []FeatureStatus
	// note: there are other unused params in the json response
}

type FeatureStatus struct {
	ID          string
	Description string
	Status      string
	SinceSlot   int
}

type InactiveFeatures []string

func (f InactiveFeatures) CLIString() string {
	return "--deactivate-feature=" + strings.Join(f, " --deactivate-feature=")
}

// GetInactiveFeatureHashes uses the solana CLI to fetch inactive solana features
// This is used in conjunction with the solana-test-validator command to produce a solana network that has the same features as mainnet
// the solana-test-validator has all features on by default (released + unreleased)
func GetInactiveFeatureHashes(url string) (output InactiveFeatures, err error) {
	cmd := exec.Command("solana", "feature", "status", "-u="+url, "--output=json") //nolint:gosec // -um is for mainnet url
	stdout, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Failed to get feature status: %w", err)
	}

	statuses := FeatureStatuses{}
	if err = json.Unmarshal(stdout, &statuses); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal feature status: %w", err)
	}

	for _, f := range statuses.Features {
		if f.Status == "inactive" {
			output = append(output, f.ID)
		}
	}

	slices.Sort(output)
	return output, err
}
