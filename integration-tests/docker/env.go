package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

// Component is a collection one or more docker containers
// that are working as a standalone product, for example Chainlink node
type Component interface {
	Prefix() string
	Containers() []tc.Container
	Start(dockerNet string, cfg any) (Component, error)
	Stop() error
}

// ComponentSetupFunc is a func to set up a Component for a particular docker network
// that may contain one or more docker containers
type ComponentSetupFunc func(dockerNet string) (Component, error)

// Environment contains several Components and can set up them sequentially or in parallel
type Environment struct {
	network    *tc.DockerNetwork
	lw         *logwatch.LogWatch
	setupFuncs []ComponentSetupFunc
	deployMu   *sync.Mutex
	Components []Component
	pushToLoki bool
}

// NewEnvironment creates a new Environment
func NewEnvironment(lw *logwatch.LogWatch, pushToLoki bool) *Environment {
	return &Environment{
		lw:         lw,
		deployMu:   &sync.Mutex{},
		Components: make([]Component, 0),
		pushToLoki: pushToLoki,
	}
}

// Add adds a new Component that can be instantiated later
func (m *Environment) Add(f ComponentSetupFunc) *Environment {
	m.setupFuncs = append(m.setupFuncs, f)
	return m
}

// Get returns a list of Components that match the Prefix
func (m *Environment) Get(prefix string) []Component {
	components := make([]Component, 0)
	for _, comp := range m.Components {
		if comp.Prefix() == prefix {
			components = append(components, comp)
		}
	}
	return components
}

// ConnectContainerLogs connects logs of all containers and stream them to Loki
func (m *Environment) ConnectContainerLogs(comp Component) error {
	if m.lw == nil {
		return nil
	}
	for _, c := range comp.Containers() {
		if err := m.lw.ConnectContainer(context.Background(), c, comp.Prefix(), m.pushToLoki); err != nil {
			return err
		}
	}
	return nil
}

// Start starts all the Components sequentially or in parallel
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
	m.setupFuncs = make([]ComponentSetupFunc, 0)
	return m, nil
}

// Shutdown performs your custom shutdown logic
func (m *Environment) Shutdown() error {
	// https://golang.testcontainers.org/features/garbage_collector/
	// ryuk container is enabled by default and he removes containers/volumes/networks by default
	return nil
}

/* utility stuff for docker environments */

// ExecuteTemplate executes Go template with some data
func ExecuteTemplate(tpl string, data any) (string, error) {
	t, err := template.New(uuid.NewString()).Parse(tpl)
	if err != nil {
		return "", errors.Join(err, ErrParsingTemplate)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	return buf.String(), err
}

// CreateNetwork creates a new docker network
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
