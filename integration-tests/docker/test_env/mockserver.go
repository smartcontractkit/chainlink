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
	C                tc.Container
	Client           *ctfClient.MockserverClient
	Networks         []string
	Endpoint         string
	InternalEndpoint string
	EAMockUrls       []*url.URL
}

func (m *MockServer) SetExternalAdapterMocks(count int) error {
	for i := 0; i < count; i++ {
		path := fmt.Sprintf("/ea-%d", i)
		err := m.Client.SetRandomValuePath(path)
		if err != nil {
			return err
		}
		cName, err := m.C.Name(context.Background())
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
		ContainerRequest: *GetMockServerContainerRequest(m.Networks),
		Started:          true,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot start mock server container")
	}
	if lw != nil {
		if err := lw.ConnectContainer(context.Background(), c, "mockserver", true); err != nil {
			return err
		}
	}
	cName, err := c.Name(context.Background())
	if err != nil {
		return err
	}
	cName = strings.Replace(cName, "/", "", -1)
	endpoint, err := c.Endpoint(context.Background(), "http")
	if err != nil {
		return err
	}
	log.Info().Any("endpoint", endpoint).Str("containerName", cName).
		Msgf("Started mockserver container")
	m.Endpoint = endpoint
	m.InternalEndpoint = fmt.Sprintf("http://%s:%s", cName, "1080")

	client := ctfClient.NewMockserverClient(&ctfClient.MockserverConfig{
		LocalURL:   endpoint,
		ClusterURL: m.InternalEndpoint,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot connect to mockserver client")
	}
	m.C = c
	m.Client = client

	return nil
}

func GetMockServerContainerRequest(networks []string) *tc.ContainerRequest {
	return &tc.ContainerRequest{
		Name:         fmt.Sprintf("mockserver-%s", uuid.NewString()),
		Image:        "mockserver/mockserver:5.11.2",
		ExposedPorts: []string{"1080/tcp"},
		Env: map[string]string{
			"SERVER_PORT": "1080",
		},
		Networks: networks,
		WaitingFor: tcwait.ForLog("INFO 1080 started on port: 1080").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
	}
}
