package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	tc "github.com/testcontainers/testcontainers-go"
	"golang.org/x/sync/errgroup"
	"sync"
	"text/template"
)

/* common stuff that should be moved to chainlink-env or CTF */

var (
	ErrContainerSetup  = errors.New("failed to deploy the container")
	ErrParsingTemplate = errors.New("failed to parse template")
)

type ContainerSetupFunc func(network string) (Component, error)

type Component interface {
	Prefix() string
	Containers() []tc.Container
	Start(network, name string, cfg any) (Component, error)
	Stop() error
}

type Environment struct {
	network    *tc.DockerNetwork
	lw         *logwatch.LogWatch
	setupFuncs []ContainerSetupFunc
	deployMu   *sync.Mutex
	Components []Component
}

func NewEnvironment(lw *logwatch.LogWatch) *Environment {
	return &Environment{
		lw:         lw,
		deployMu:   &sync.Mutex{},
		Components: make([]Component, 0),
	}
}

func (m *Environment) WithContainer(f ContainerSetupFunc) *Environment {
	m.setupFuncs = append(m.setupFuncs, f)
	return m
}

func (m *Environment) Get(prefix string) []Component {
	components := make([]Component, 0)
	for _, comp := range m.Components {
		if comp.Prefix() == prefix {
			components = append(components, comp)
		}
	}
	return components
}

func (m *Environment) ConnectContainerLogs(comp Component) error {
	if m.lw == nil {
		return nil
	}
	for _, c := range comp.Containers() {
		if err := m.lw.ConnectContainer(context.Background(), c, comp.Prefix(), true); err != nil {
			return err
		}
	}
	return nil
}

func (m *Environment) Start(parallel bool) (*Environment, error) {
	// all containers are isolated inside docker network, per environment
	if m.network == nil {
		network, err := CreateNetwork()
		if err != nil {
			return nil, fmt.Errorf("failed to create docker network: %s", err)
		}
		m.network = network
	}
	if !parallel {
		for _, f := range m.setupFuncs {
			comp, err := f(m.network.Name)
			if err != nil {
				return nil, errors.Join(err, ErrContainerSetup)
			}
			if err := m.ConnectContainerLogs(comp); err != nil {
				return nil, err
			}
			m.Components = append(m.Components, comp)
		}
	} else {
		eg := &errgroup.Group{}
		for _, f := range m.setupFuncs {
			f := f
			eg.Go(func() error {
				comp, err := f(m.network.Name)
				if err != nil {
					return errors.Join(err, ErrContainerSetup)
				}
				if err := m.ConnectContainerLogs(comp); err != nil {
					return err
				}
				m.deployMu.Lock()
				defer m.deployMu.Unlock()
				m.Components = append(m.Components, comp)
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return m, err
		}
	}
	// we are executing setup only once, then delete all setup functions
	// in case you need to do multi-stage deployments
	m.setupFuncs = make([]ContainerSetupFunc, 0)
	return m, nil
}

func (m *Environment) Shutdown() error {
	return m.network.Remove(context.Background())
}

/* utility stuff for docker environments */

func ExecuteTemplate(tpl string, data any) (string, error) {
	t, err := template.New(uuid.NewString()).Parse(tpl)
	if err != nil {
		return "", errors.Join(err, ErrParsingTemplate)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	return buf.String(), err
}

func CreateNetwork() (*tc.DockerNetwork, error) {
	network, err := tc.GenericNetwork(context.Background(), tc.GenericNetworkRequest{
		NetworkRequest: tc.NetworkRequest{
			Name:           fmt.Sprintf("network-%s", uuid.NewString()),
			CheckDuplicate: true,
		},
	})
	if err != nil {
		return nil, err
	}
	dockerNetwork, ok := network.(*tc.DockerNetwork)
	if !ok {
		return nil, fmt.Errorf("failed to cast network to *dockertest.Network")
	}
	log.Debug().Any("network", dockerNetwork).Msgf("created network")
	return dockerNetwork, nil
}

func ContainerResources(quota int64, memoryMb int64) func(config *container.HostConfig) {
	return func(config *container.HostConfig) {
		config.CPUQuota = quota
		config.Memory = memoryMb * 10e8
	}
}
