package docker

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"strings"
	"time"
)

type MockServer struct {
	prefix    string
	container tc.Container
	// export
	Endpoint         string
	InternalEndpoint string
}

func NewMockServer(cfg any) ContainerSetupFunc {
	c := &MockServer{prefix: "mockserver"}
	return func(network string) (Component, error) {
		return c.Start(network, c.prefix, cfg)
	}
}

func (m *MockServer) Prefix() string {
	return m.prefix
}

func (m *MockServer) Containers() []tc.Container {
	return []tc.Container{m.container}
}

func (m *MockServer) Start(network, name string, cfg any) (Component, error) {
	req := tc.GenericContainerRequest{
		ContainerRequest: *msContainerRequest(network, name),
		Started:          true,
	}
	c, err := tc.GenericContainer(context.Background(), req)
	if err != nil {
		return m, err
	}
	cName, err := c.Name(context.Background())
	if err != nil {
		return m, err
	}
	cName = strings.Replace(cName, "/", "", -1)
	_, err = c.MappedPort(context.Background(), "1080/tcp")
	if err != nil {
		return nil, err
	}
	endpoint, err := c.Endpoint(context.Background(), "http")
	if err != nil {
		return m, err
	}
	m.Endpoint = endpoint
	m.InternalEndpoint = fmt.Sprintf("http://%s:1080", cName)
	log.Info().Any("endpoint", endpoint).Str("containerName", cName).
		Msgf("Started mockserver container")

	//client := ctfClient.NewMockserverClient(&ctfClient.MockserverConfig{
	//	LocalURL: endpoint,
	//})
	//if err != nil {
	//	return errors.Wrapf(err, "cannot connect to mockserver client")
	//}
	m.container = c
	return m, nil
}

func (m *MockServer) Stop() error {
	return m.container.Terminate(context.Background())
}

func msContainerRequest(network, name string) *tc.ContainerRequest {
	return &tc.ContainerRequest{
		Name:         fmt.Sprintf("%s-%s", name, uuid.NewString()),
		Image:        "mockserver/mockserver:5.11.2",
		ExposedPorts: []string{"1080/tcp"},
		Env: map[string]string{
			"SERVER_PORT": "1080",
		},
		Networks:           []string{network},
		HostConfigModifier: ContainerResources(50000, 500),
		WaitingFor: tcwait.ForLog("INFO 1080 started on port: 1080").
			WithStartupTimeout(90 * time.Second).
			WithPollInterval(1 * time.Second),
	}
}
