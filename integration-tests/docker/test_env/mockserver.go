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
			ContainerName: fmt.Sprintf("%s-%s", "mockserver", uuid.NewString()[0:8]),
			Networks:      networks,
		},
	}
	for _, opt := range opts {
		opt(&ms.EnvComponent)
	}
	return ms
}

func (ms *MockServer) SetExternalAdapterMocks(count int) error {
	for i := 0; i < count; i++ {
		path := fmt.Sprintf("/ea-%d", i)
		err := ms.Client.SetRandomValuePath(path)
		if err != nil {
			return err
		}
		cName, err := ms.Container.Name(context.Background())
		if err != nil {
			return err
		}
		cName = strings.Replace(cName, "/", "", -1)
		eaUrl, err := url.Parse(fmt.Sprintf("http://%s:%s%s",
			cName, "1080", path))
		if err != nil {
			return err
		}
		ms.EAMockUrls = append(ms.EAMockUrls, eaUrl)
	}
	return nil
}

func (ms *MockServer) StartContainer() error {
	c, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: ms.getContainerRequest(),
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot start MockServer container")
	}
	ms.Container = c
	endpoint, err := c.Endpoint(context.Background(), "http")
	if err != nil {
		return err
	}
	log.Info().Any("endpoint", endpoint).Str("containerName", ms.ContainerName).
		Msgf("Started MockServer container")
	ms.Endpoint = endpoint
	ms.InternalEndpoint = fmt.Sprintf("http://%s:%s", ms.ContainerName, "1080")

	client := ctfClient.NewMockserverClient(&ctfClient.MockserverConfig{
		LocalURL:   endpoint,
		ClusterURL: ms.InternalEndpoint,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot connect to MockServer client")
	}
	ms.Client = client

	return nil
}

func (ms *MockServer) getContainerRequest() tc.ContainerRequest {
	return tc.ContainerRequest{
		Name:         ms.ContainerName,
		Image:        "mockserver/mockserver:5.15.0",
		ExposedPorts: []string{"1080/tcp"},
		Env: map[string]string{
			"SERVER_PORT": "1080",
		},
		Networks: ms.Networks,
		WaitingFor: tcwait.ForLog("INFO 1080 started on port: 1080").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
	}
}
