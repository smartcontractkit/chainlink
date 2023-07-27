package docker

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	ErrNoChainlinkImage   = errors.New("CHAINLINK_IMAGE environment variable must be set")
	ErrNoChainlinkVersion = errors.New("CHAINLINK_VERSION environment variable must be set")
)

/* templates for CL node configuration */

var (
	NodeConfigTemplate = `RootDir = '/home/chainlink'

[Database]
MaxIdleConns = 20
MaxOpenConns = 40
MigrateOnStartup = true

[Log]
Level = 'debug'
JSONConsole = true

[WebServer]
AllowOrigins = '*'
HTTPPort = 6688
SecureCookies = false
SessionTimeout = '999h0m0s'

[WebServer.TLS]
HTTPSPort = 0

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 100

[[EVM]]
ChainID = "1337"
MinContractPayment = '0'
AutoCreateKey = true
finalityDepth = 1

[[EVM.Nodes]]
HTTPURL = "{{ .EVM.HTTPURL }}"
Name = "1337_primary_local_0"
SendOnly = false
WSURL = "{{ .EVM.WSURL }}"

[Feature]
LogPoller = true
FeedsManager = true
UICSAKeys = true

[OCR]
Enabled = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690
`

	NodeSecretsTemplate = `
[Database]
URL = 'postgresql://postgres:test@{{ .PgHost }}:{{ .PgPort }}/cl-node?sslmode=disable' # Required
AllowSimplePasswords = true

[Password]
Keystore = '................' # Required

[Mercury.Credentials.cred1]
# URL = 'http://host.docker.internal:3000/reports'
URL = 'localhost:1338'
Username = 'node'
Password = 'nodepass'
`
)

/* data types for templates */

type (
	NodeSecretsTemplateValues struct {
		PgHost string
		PgPort string
	}

	NodeConfigOpts struct {
		EVM NodeEVMSettings
	}
)

type NodeEVMSettings struct {
	HTTPURL string
	WSURL   string
}

type PgOpts struct {
	User     string
	Password string
	DbName   string
	Networks []string
	Port     string
}

func nodeContainerRequest(name string, nodeConfOpts NodeConfigOpts, secrets string, networks []string) (
	*tc.ContainerRequest, error) {
	configFile, err := ioutil.TempFile("", "node_config")
	if err != nil {
		return nil, err
	}
	config, err := ExecuteTemplate(NodeConfigTemplate, nodeConfOpts)
	if err != nil {
		return nil, err
	}
	_, err = configFile.WriteString(config)
	if err != nil {
		return nil, err
	}

	secretsFile, err := ioutil.TempFile("", "node_secrets")
	if err != nil {
		return nil, err
	}
	_, err = secretsFile.WriteString(secrets)
	if err != nil {
		return nil, err
	}

	adminCreds := "local@local.com\nlocaldevpassword"
	adminCredsFile, err := ioutil.TempFile("", "admin_creds")
	if err != nil {
		return nil, err
	}
	_, err = adminCredsFile.WriteString(adminCreds)
	if err != nil {
		return nil, err
	}

	apiCreds := "local@local.com\nlocaldevpassword"
	apiCredsFile, err := ioutil.TempFile("", "api_creds")
	if err != nil {
		return nil, err
	}
	_, err = apiCredsFile.WriteString(apiCreds)
	if err != nil {
		return nil, err
	}

	configPath := "/home/cl-node-config.toml"
	secretsPath := "/home/cl-node-secrets.toml"
	adminCredsPath := "/home/admin-credentials.txt"
	apiCredsPath := "/home/api-credentials.txt"

	image, ok := os.LookupEnv("CHAINLINK_IMAGE")
	if !ok {
		return nil, ErrNoChainlinkImage
	}
	tag, ok := os.LookupEnv("CHAINLINK_VERSION")
	if !ok {
		return nil, ErrNoChainlinkVersion
	}

	id := uuid.New()
	return &tc.ContainerRequest{
		Name:         fmt.Sprintf("%s-%s", name, id.String()[0:5]),
		Image:        fmt.Sprintf("%s:%s", image, tag),
		ExposedPorts: []string{"6688/tcp", "6690/tcp"},
		Entrypoint: []string{"chainlink",
			"-c", configPath,
			"-s", secretsPath,
			"node", "start", "-d",
			"-p", adminCredsPath,
			"-a", apiCredsPath,
		},
		Networks: networks,
		WaitingFor: tcwait.ForHTTP("/health").
			WithPort("6688/tcp").
			WithStartupTimeout(120 * time.Second).
			WithPollInterval(1 * time.Second),
		Files: []tc.ContainerFile{
			{
				HostFilePath:      configFile.Name(),
				ContainerFilePath: configPath,
				FileMode:          0644,
			},
			{
				HostFilePath:      secretsFile.Name(),
				ContainerFilePath: secretsPath,
				FileMode:          0644,
			},
			{
				HostFilePath:      adminCredsFile.Name(),
				ContainerFilePath: adminCredsPath,
				FileMode:          0644,
			},
			{
				HostFilePath:      apiCredsFile.Name(),
				ContainerFilePath: apiCredsPath,
				FileMode:          0644,
			},
		},
	}, nil
}

type Chainlink struct {
	prefix      string
	container   tc.Container
	dbContainer tc.Container
	// export
	Endpoint string
}

func NewChainlink(cfg any) ContainerSetupFunc {
	c := &Chainlink{prefix: "chainlink"}
	return func(network string) (Component, error) {
		return c.Start(network, c.prefix, cfg)
	}
}

func (m *Chainlink) Prefix() string {
	return m.prefix
}

func (m *Chainlink) Containers() []tc.Container {
	containers := make([]tc.Container, 0)
	return append(containers, m.dbContainer, m.container)
}

func (m *Chainlink) Stop() error {
	if err := m.container.Terminate(context.Background()); err != nil {
		return err
	}
	return m.dbContainer.Terminate(context.Background())
}

func (m *Chainlink) Start(network, name string, cfg any) (Component, error) {
	nco, ok := cfg.(NodeConfigOpts)
	if !ok {
		return m, fmt.Errorf("cfg must be of type NodeConfigOpts")
	}
	networks := []string{network}
	pgOpts := PgOpts{
		Port:     "5432",
		User:     "postgres",
		Password: "test",
		DbName:   "cl-node",
		Networks: networks,
	}
	pgReq := pgContainerRequest(pgOpts)
	dbContainer, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: pgReq,
		Started:          true,
	})
	m.dbContainer = dbContainer
	if err != nil {
		return m, err
	}

	nodeSecrets, err := ExecuteTemplate(NodeSecretsTemplate, NodeSecretsTemplateValues{
		PgHost: pgReq.Name,
		PgPort: "5432",
	})
	if err != nil {
		return m, err
	}
	clReq, err := nodeContainerRequest(name, nco, nodeSecrets, networks)
	if err != nil {
		return m, err
	}
	clC, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: *clReq,
		Started:          true,
	})
	if err != nil {
		return m, err
	}
	m.container = clC
	ctName, err := clC.Name(context.Background())
	if err != nil {
		return m, err
	}
	ctName = strings.Replace(ctName, "/", "", -1)
	clEndpoint, err := clC.Endpoint(context.Background(), "http")
	if err != nil {
		return m, err
	}
	m.Endpoint = clEndpoint

	log.Info().Str("containerName", ctName).
		Str("clEndpoint", clEndpoint).
		Msg("Started Chainlink Node container")
	return m, nil
}

func pgContainerRequest(o PgOpts) tc.ContainerRequest {
	return tc.ContainerRequest{
		Name:         fmt.Sprintf("pg-%s", uuid.NewString()),
		Image:        "postgres:15.3",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", o.Port)},
		Env: map[string]string{
			"POSTGRES_USER":     o.User,
			"POSTGRES_DB":       o.DbName,
			"POSTGRES_PASSWORD": o.Password,
		},
		Networks: o.Networks,
		WaitingFor: tcwait.ForExec([]string{"psql", "-h", "localhost",
			"-U", o.User, "-c", "select", "1", "-d", o.DbName}).
			WithStartupTimeout(90 * time.Second),
	}
}
