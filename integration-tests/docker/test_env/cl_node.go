package test_env

import (
	"context"
	"fmt"
	"io"
	"maps"
	"math/big"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	it_utils "github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/utils/templates"
)

var (
	ErrConnectNodeClient    = "could not connect Node HTTP Client"
	ErrStartCLNodeContainer = "failed to start CL node container"
)

const (
	RestartContainer  = true
	StartNewContainer = false
)

type ClNode struct {
	test_env.EnvComponent
	API                   *client.ChainlinkClient `json:"-"`
	NodeConfig            *chainlink.Config       `json:"-"`
	NodeSecretsConfigTOML string                  `json:"-"`
	PostgresDb            *test_env.PostgresDb    `json:"postgresDb"`
	UserEmail             string                  `json:"userEmail"`
	UserPassword          string                  `json:"userPassword"`
	AlwaysPullImage       bool                    `json:"-"`
	t                     *testing.T
	l                     zerolog.Logger
}

type ClNodeOption = func(c *ClNode)

func WithSecrets(secretsTOML string) ClNodeOption {
	return func(c *ClNode) {
		c.NodeSecretsConfigTOML = secretsTOML
	}
}

func WithNodeEnvVars(ev map[string]string) ClNodeOption {
	return func(n *ClNode) {
		if n.ContainerEnvs == nil {
			n.ContainerEnvs = map[string]string{}
		}
		maps.Copy(n.ContainerEnvs, ev)
	}
}

// Sets custom node container name if name is not empty
func WithNodeContainerName(name string) ClNodeOption {
	return func(c *ClNode) {
		if name != "" {
			c.ContainerName = name
		}
	}
}

// Sets custom node db container name if name is not empty
func WithDbContainerName(name string) ClNodeOption {
	return func(c *ClNode) {
		if name != "" {
			c.PostgresDb.ContainerName = name
		}
	}
}

func WithLogStream(ls *logstream.LogStream) ClNodeOption {
	return func(c *ClNode) {
		c.LogStream = ls
	}
}

func WithImage(image string) ClNodeOption {
	return func(c *ClNode) {
		c.ContainerImage = image
	}
}

func WithVersion(version string) ClNodeOption {
	return func(c *ClNode) {
		c.ContainerVersion = version
	}
}

func WithPgDBOptions(opts ...test_env.PostgresDbOption) ClNodeOption {
	return func(c *ClNode) {
		var err error
		c.PostgresDb, err = test_env.NewPostgresDb(c.EnvComponent.Networks, opts...)
		if err != nil {
			c.t.Fatalf("failed to create postgres db: %v", err)
		}
	}
}

func NewClNode(networks []string, imageName, imageVersion string, nodeConfig *chainlink.Config, opts ...ClNodeOption) (*ClNode, error) {
	nodeDefaultCName := fmt.Sprintf("%s-%s", "cl-node", uuid.NewString()[0:8])
	pgDefaultCName := fmt.Sprintf("pg-%s", nodeDefaultCName)
	pgDb, err := test_env.NewPostgresDb(networks, test_env.WithPostgresDbContainerName(pgDefaultCName))
	if err != nil {
		return nil, err
	}
	n := &ClNode{
		EnvComponent: test_env.EnvComponent{
			ContainerName:    nodeDefaultCName,
			ContainerImage:   imageName,
			ContainerVersion: imageVersion,
			Networks:         networks,
		},
		UserEmail:    "local@local.com",
		UserPassword: "localdevpassword",
		NodeConfig:   nodeConfig,
		PostgresDb:   pgDb,
		l:            log.Logger,
	}
	n.SetDefaultHooks()
	for _, opt := range opts {
		opt(n)
	}
	return n, nil
}

func (n *ClNode) SetTestLogger(t *testing.T) {
	n.l = logging.GetTestLogger(t)
	n.t = t
	n.PostgresDb.WithTestInstance(t)
}

// Restart restarts only CL node, DB container is reused
func (n *ClNode) Restart(cfg *chainlink.Config) error {
	if err := n.Container.Terminate(testcontext.Get(n.t)); err != nil {
		return err
	}
	n.NodeConfig = cfg
	return n.RestartContainer()
}

// UpgradeVersion restarts the cl node with new image and version
func (n *ClNode) UpgradeVersion(newImage, newVersion string) error {
	n.l.Info().
		Str("Name", n.ContainerName).
		Str("Old Image", newImage).
		Str("Old Version", newVersion).
		Str("New Image", newImage).
		Str("New Version", newVersion).
		Msg("Upgrading Chainlink Node")
	n.ContainerImage = newImage
	n.ContainerVersion = newVersion
	return n.Restart(n.NodeConfig)
}

func (n *ClNode) PrimaryETHAddress() (string, error) {
	return n.API.PrimaryEthAddress()
}

func (n *ClNode) AddBootstrapJob(verifierAddr common.Address, chainId int64,
	feedId [32]byte) (*client.Job, error) {
	spec := it_utils.BuildBootstrapSpec(verifierAddr, chainId, feedId)
	return n.API.MustCreateJob(spec)
}

func (n *ClNode) AddMercuryOCRJob(verifierAddr common.Address, fromBlock uint64, chainId int64,
	feedId [32]byte, customAllowedFaults *int, bootstrapUrl string,
	mercuryServerUrl string, mercuryServerPubKey string,
	eaUrls []*url.URL) (*client.Job, error) {

	csaKeys, _, err := n.API.ReadCSAKeys()
	if err != nil {
		return nil, err
	}
	csaPubKey := csaKeys.Data[0].ID

	nodeOCRKeys, err := n.API.MustReadOCR2Keys()
	if err != nil {
		return nil, err
	}

	var nodeOCRKeyId []string
	for _, key := range nodeOCRKeys.Data {
		if key.Attributes.ChainType == string(chaintype.EVM) {
			nodeOCRKeyId = append(nodeOCRKeyId, key.ID)
			break
		}
	}

	bridges := it_utils.BuildBridges(eaUrls)
	for index := range bridges {
		err = n.API.MustCreateBridge(&bridges[index])
		if err != nil {
			return nil, err
		}
	}

	var allowedFaults int
	if customAllowedFaults != nil {
		allowedFaults = *customAllowedFaults
	} else {
		allowedFaults = 2
	}

	spec := it_utils.BuildOCRSpec(
		verifierAddr, chainId, fromBlock, feedId, bridges,
		csaPubKey, mercuryServerUrl, mercuryServerPubKey, nodeOCRKeyId[0],
		bootstrapUrl, allowedFaults)

	return n.API.MustCreateJob(spec)
}

func (n *ClNode) GetContainerName() string {
	name, err := n.Container.Name(testcontext.Get(n.t))
	if err != nil {
		return ""
	}
	return strings.Replace(name, "/", "", -1)
}

func (n *ClNode) GetAPIClient() *client.ChainlinkClient {
	return n.API
}

func (n *ClNode) GetPeerUrl() (string, error) {
	p2pKeys, err := n.API.MustReadP2PKeys()
	if err != nil {
		return "", err
	}
	p2pId := p2pKeys.Data[0].Attributes.PeerID

	return fmt.Sprintf("%s@%s:%d", p2pId, n.GetContainerName(), 6690), nil
}

func (n *ClNode) GetNodeCSAKeys() (*client.CSAKeys, error) {
	csaKeys, _, err := n.API.ReadCSAKeys()
	if err != nil {
		return nil, err
	}
	return csaKeys, err
}

func (n *ClNode) ChainlinkNodeAddress() (common.Address, error) {
	addr, err := n.API.PrimaryEthAddress()
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(addr), nil
}

func (n *ClNode) Fund(evmClient blockchain.EVMClient, amount *big.Float) error {
	toAddress, err := n.API.PrimaryEthAddress()
	if err != nil {
		return err
	}
	n.l.Debug().
		Str("ChainId", evmClient.GetChainID().String()).
		Str("Address", toAddress).
		Msg("Funding Chainlink Node")
	toAddr := common.HexToAddress(toAddress)
	gasEstimates, err := evmClient.EstimateGas(ethereum.CallMsg{
		To: &toAddr,
	})
	if err != nil {
		return err
	}
	return evmClient.Fund(toAddress, amount, gasEstimates)
}

func (n *ClNode) containerStartOrRestart(restartDb bool) error {
	var err error
	if restartDb {
		err = n.PostgresDb.RestartContainer()
	} else {
		err = n.PostgresDb.StartContainer()
	}
	if err != nil {
		return err
	}

	// If the node secrets TOML is not set, generate it with the default template
	nodeSecretsToml, err := templates.NodeSecretsTemplate{
		PgDbName:      n.PostgresDb.DbName,
		PgHost:        strings.Split(n.PostgresDb.InternalURL.Host, ":")[0],
		PgPort:        n.PostgresDb.InternalPort,
		PgPassword:    n.PostgresDb.Password,
		CustomSecrets: n.NodeSecretsConfigTOML,
	}.String()
	if err != nil {
		return err
	}

	cReq, err := n.getContainerRequest(nodeSecretsToml)
	if err != nil {
		return err
	}

	l := tc.Logger
	if n.t != nil {
		l = logging.CustomT{
			T: n.t,
			L: n.l,
		}
	}
	container, err := docker.StartContainerWithRetry(n.l, tc.GenericContainerRequest{
		ContainerRequest: *cReq,
		Started:          true,
		Reuse:            false,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("%s err: %w", ErrStartCLNodeContainer, err)
	}

	clEndpoint, err := test_env.GetEndpoint(testcontext.Get(n.t), container, "http")
	if err != nil {
		return err
	}
	ip, err := container.ContainerIP(testcontext.Get(n.t))
	if err != nil {
		return err
	}
	n.l.Info().
		Str("containerName", n.ContainerName).
		Str("containerImage", n.ContainerImage).
		Str("containerVersion", n.ContainerVersion).
		Str("clEndpoint", clEndpoint).
		Str("clInternalIP", ip).
		Str("userEmail", n.UserEmail).
		Str("userPassword", n.UserPassword).
		Msg("Started Chainlink Node container")
	clClient, err := client.NewChainlinkClient(&client.ChainlinkConfig{
		URL:        clEndpoint,
		Email:      n.UserEmail,
		Password:   n.UserPassword,
		InternalIP: ip,
	},
		n.l)
	if err != nil {
		return fmt.Errorf("%s err: %w", ErrConnectNodeClient, err)
	}

	n.Container = container
	n.API = clClient

	return nil
}

func (n *ClNode) RestartContainer() error {
	return n.containerStartOrRestart(RestartContainer)
}

func (n *ClNode) StartContainer() error {
	return n.containerStartOrRestart(StartNewContainer)
}

func (n *ClNode) ExecGetVersion() (string, error) {
	cmd := []string{"chainlink", "--version"}
	_, output, err := n.Container.Exec(context.Background(), cmd)
	if err != nil {
		return "", errors.Wrapf(err, "could not execute cmd %s", cmd)
	}
	outputBytes, err := io.ReadAll(output)
	if err != nil {
		return "", err
	}
	outputString := strings.TrimSpace(string(outputBytes))

	// Find version in cmd output
	re := regexp.MustCompile("@(.*)")
	matches := re.FindStringSubmatch(outputString)

	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", errors.Errorf("could not find chainlink version in command output '%'", output)
}

func (n *ClNode) getContainerRequest(secrets string) (
	*tc.ContainerRequest, error) {
	configFile, err := os.CreateTemp("", "node_config")
	if err != nil {
		return nil, err
	}
	data, err := toml.Marshal(n.NodeConfig)
	if err != nil {
		return nil, err
	}
	_, err = configFile.WriteString(string(data))
	if err != nil {
		return nil, err
	}
	secretsFile, err := os.CreateTemp("", "node_secrets")
	if err != nil {
		return nil, err
	}
	_, err = secretsFile.WriteString(secrets)
	if err != nil {
		return nil, err
	}

	adminCreds := "local@local.com\nlocaldevpassword"
	adminCredsFile, err := os.CreateTemp("", "admin_creds")
	if err != nil {
		return nil, err
	}
	_, err = adminCredsFile.WriteString(adminCreds)
	if err != nil {
		return nil, err
	}

	apiCreds := "local@local.com\nlocaldevpassword"
	apiCredsFile, err := os.CreateTemp("", "api_creds")
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

	return &tc.ContainerRequest{
		Name:            n.ContainerName,
		AlwaysPullImage: n.AlwaysPullImage,
		Image:           fmt.Sprintf("%s:%s", n.ContainerImage, n.ContainerVersion),
		ExposedPorts:    []string{"6688/tcp"},
		Env:             n.ContainerEnvs,
		Entrypoint: []string{"chainlink",
			"-c", configPath,
			"-s", secretsPath,
			"node", "start", "-d",
			"-p", adminCredsPath,
			"-a", apiCredsPath,
		},
		Networks: append(n.Networks, "tracing"),
		WaitingFor: tcwait.ForHTTP("/health").
			WithPort("6688/tcp").
			WithStartupTimeout(90 * time.Second).
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
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts:    n.PostStartsHooks,
				PostStops:     n.PostStopsHooks,
				PreTerminates: n.PreTerminatesHooks,
			},
		},
	}, nil
}
