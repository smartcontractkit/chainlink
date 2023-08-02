package test_env

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

type MockServer struct {
	EnvComponent
	Client           *ctfClient.MockserverClient
	Endpoint         string
	InternalEndpoint string
	EAMockUrls       []*url.URL
}

func NewMockServer(networks []string, opts ...EnvComponentOption) *MockServer {
	ms := &MockServer{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "mockserver", uuid.NewString()[0:3]),
			Networks:      networks,
		},
	}
	for _, opt := range opts {
		opt(&ms.EnvComponent)
	}
	return ms
}

func (m *MockServer) SetExternalAdapterMocks(count int) error {
	for i := 0; i < count; i++ {
		path := fmt.Sprintf("/ea-%d", i)
		err := m.Client.SetRandomValuePath(path)
		if err != nil {
			return err
		}
		cName, err := m.EnvComponent.Container.Name(context.Background())
		if err != nil {
			return err
		}
		cName = strings.Replace(cName, "/", "", -1)
		eaUrl, err := url.Parse(fmt.Sprintf("http://%s:%s%s",
			cName, "1080", path))
		if err != nil {
			return err
		}
		m.EAMockUrls = append(m.EAMockUrls, eaUrl)
	}
	return nil
}

func (m *MockServer) StartContainer(lw *logwatch.LogWatch) error {
	c, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: m.getContainerRequest(),
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot start MockServer container")
	}
	m.Container = c
	if lw != nil {
		if err := lw.ConnectContainer(context.Background(), c, m.ContainerName, true); err != nil {
			return err
		}
	}
	endpoint, err := c.Endpoint(context.Background(), "http")
	if err != nil {
		return err
	}
	log.Info().Any("endpoint", endpoint).Str("containerName", m.ContainerName).
		Msgf("Started MockServer container")
	m.Endpoint = endpoint
	m.InternalEndpoint = fmt.Sprintf("http://%s:%s", m.ContainerName, "1080")

	client := ctfClient.NewMockserverClient(&ctfClient.MockserverConfig{
		LocalURL:   endpoint,
		ClusterURL: m.InternalEndpoint,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot connect to MockServer client")
	}
	m.Client = client

	return nil
}

func (m *MockServer) getContainerRequest() tc.ContainerRequest {
	return tc.ContainerRequest{
		Name:         m.ContainerName,
		Image:        "mockserver/mockserver:5.11.2",
		ExposedPorts: []string{"1080/tcp"},
		Env: map[string]string{
			"SERVER_PORT": "1080",
		},
		Networks: m.Networks,
		WaitingFor: tcwait.ForLog("INFO 1080 started on port: 1080").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
	}
}
