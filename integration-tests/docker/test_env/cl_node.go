package test_env

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
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

type ClNode struct {
	test_env.EnvComponent
	API                   *client.ChainlinkClient `json:"-"`
	NodeConfig            *chainlink.Config       `json:"-"`
	NodeSecretsConfigTOML string                  `json:"-"`
	PostgresDb            *test_env.PostgresDb    `json:"postgresDb"`
	UserEmail             string                  `json:"userEmail"`
	UserPassword          string                  `json:"userPassword"`
	t                     *testing.T
	l                     zerolog.Logger
	lw                    *logwatch.LogWatch
}

type ClNodeOption = func(c *ClNode)

func WithSecrets(secretsTOML string) ClNodeOption {
	return func(c *ClNode) {
		c.NodeSecretsConfigTOML = secretsTOML
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

func WithLogWatch(lw *logwatch.LogWatch) ClNodeOption {
	return func(c *ClNode) {
		c.lw = lw
	}
}

func NewClNode(networks []string, imageName, imageVersion string, nodeConfig *chainlink.Config, opts ...ClNodeOption) (*ClNode, error) {
	nodeDefaultCName := fmt.Sprintf("%s-%s", "cl-node", uuid.NewString()[0:8])
	pgDefaultCName := fmt.Sprintf("pg-%s", nodeDefaultCName)
	pgDb := test_env.NewPostgresDb(networks, test_env.WithPostgresDbContainerName(pgDefaultCName))
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
	for _, opt := range opts {
		opt(n)
	}
	return n, nil
}

func (n *ClNode) SetTestLogger(t *testing.T) {
	n.l = logging.GetTestLogger(t)
	n.t = t
	n.PostgresDb.WithTestLogger(t)
}

// Restart restarts only CL node, DB container is reused
func (n *ClNode) Restart(cfg *chainlink.Config) error {
	if err := n.Container.Terminate(testcontext.Get(n.t)); err != nil {
		return err
	}
	n.NodeConfig = cfg
	return n.StartContainer()
}

// UpgradeVersion restarts the cl node with new image and version
func (n *ClNode) UpgradeVersion(newImage, newVersion string) error {
	if newVersion == "" {
		return fmt.Errorf("new version is empty")
	}
	if newImage == "" {
		return fmt.Errorf("new image name is empty")
	}
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
	toAddr := common.HexToAddress(toAddress)
	gasEstimates, err := evmClient.EstimateGas(ethereum.CallMsg{
		To: &toAddr,
	})
	if err != nil {
		return err
	}
	return evmClient.Fund(toAddress, amount, gasEstimates)
}

func (n *ClNode) StartContainer() error {
	err := n.PostgresDb.StartContainer()
	if err != nil {
		return err
	}

	// If the node secrets TOML is not set, generate it with the default template
	nodeSecretsToml, err := templates.NodeSecretsTemplate{
		PgDbName:      n.PostgresDb.DbName,
		PgHost:        n.PostgresDb.ContainerName,
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
		Reuse:            true,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("%s err: %w", ErrStartCLNodeContainer, err)
	}
	// if n.lw != nil {
	// 	if err := n.lw.ConnectContainer(testcontext.Get(n.t), container, "cl-node"); err != nil {
	// 		return err
	// 	}
	// }
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
	clClient.Config.InternalIP = n.ContainerName
	n.Container = container
	n.API = clClient

	return nil
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
		Name:         n.ContainerName,
		Image:        fmt.Sprintf("%s:%s", n.ContainerImage, n.ContainerVersion),
		ExposedPorts: []string{"6688/tcp"},
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
			{PostStarts: []tc.ContainerHook{
				func(ctx context.Context, c tc.Container) error {
					if n.lw != nil {
						return n.lw.ConnectContainer(ctx, c, "cl-node")
					}
					return nil
				},
			},
				PostStops: []tc.ContainerHook{
					func(ctx context.Context, c tc.Container) error {
						if n.lw != nil {
							return n.lw.DisconnectContainer(c)
						}
						return nil
					},
				}},
		},
	}, nil
}
