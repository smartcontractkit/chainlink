package test_env

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	env "github.com/smartcontractkit/chainlink/integration-tests/types/envcommon"
	"github.com/smartcontractkit/chainlink/integration-tests/types/node"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/utils/templates"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"math/big"
)

type ClNode struct {
	env.EnvComponent
	API            *client.ChainlinkClient
	NodeConfigOpts node.NodeConfigOpts
	DbC            *tc.Container
	DbCName        string
	DbOpts         env.PgOpts
}

func NewClNode(compOpts env.EnvComponentOpts, opts node.NodeConfigOpts, dbContainerName string) *ClNode {
	return &ClNode{
		EnvComponent:   env.NewEnvComponent("cl-node", compOpts),
		DbCName:        dbContainerName,
		NodeConfigOpts: opts,
		DbOpts:         env.NewDefaultPgOpts("cl-node", compOpts.Networks),
	}
}

func (m *ClNode) AddBootstrapJob(verifierAddr common.Address, fromBlock uint64, chainId int64,
	feedId [32]byte) (*client.Job, error) {
	spec := utils.BuildBootstrapSpec(verifierAddr, chainId, fromBlock, feedId)
	return m.API.MustCreateJob(spec)
}

func (m *ClNode) GetContainerName() string {
	name, err := m.EnvComponent.Container.Name(context.Background())
	if err != nil {
		return ""
	}
	return strings.Replace(name, "/", "", -1)
}

func (m *ClNode) GetPeerUrl() (string, error) {
	p2pKeys, err := m.API.MustReadP2PKeys()
	if err != nil {
		return "", err
	}
	p2pId := p2pKeys.Data[0].Attributes.PeerID

	return fmt.Sprintf("%s@%s:%d", p2pId, m.GetContainerName(), 6690), nil
}

func (m *ClNode) GetNodeCSAKeys() (*client.CSAKeys, error) {
	csaKeys, _, err := m.API.ReadCSAKeys()
	if err != nil {
		return nil, err
	}
	return csaKeys, err
}

func (m *ClNode) ChainlinkNodeAddress() (common.Address, error) {
	addr, err := m.API.PrimaryEthAddress()
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(addr), nil
}

func (m *ClNode) Fund(g *Geth, amount *big.Float) error {
	toAddress, err := m.API.PrimaryEthAddress()
	if err != nil {
		return err
	}
	gasEstimates, err := g.EthClient.EstimateGas(ethereum.CallMsg{})
	if err != nil {
		return err
	}
	return g.EthClient.Fund(toAddress, amount, gasEstimates)
}

func (m *ClNode) StartContainer(lw *logwatch.LogWatch) error {
	pgReq := env.GetPgContainerRequest(m.DbCName, m.DbOpts)
	pgC, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: *pgReq,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return err
	}

	nodeSecrets, err := templates.ExecuteNodeSecretsTemplate(pgReq.Name, "5432")
	if err != nil {
		return err
	}
	clReq, err := m.getContainerRequest(nodeSecrets)
	if err != nil {
		return err
	}
	clC, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: *clReq,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return errors.Wrapf(err, "could not start chainlink node container")
	}
	if lw != nil {
		if err := lw.ConnectContainer(context.Background(), clC, "chainlink", true); err != nil {
			return err
		}
	}
	ctName, err := clC.Name(context.Background())
	if err != nil {
		return err
	}
	ctName = strings.Replace(ctName, "/", "", -1)
	clEndpoint, err := clC.Endpoint(context.Background(), "http")
	if err != nil {
		return err
	}

	log.Info().Str("containerName", ctName).
		Str("clEndpoint", clEndpoint).
		Msg("Started Chainlink Node container")

	clClient, err := client.NewChainlinkClient(&client.ChainlinkConfig{
		URL:      clEndpoint,
		Email:    "local@local.com",
		Password: "localdevpassword",
	})
	if err != nil {
		return errors.Wrapf(err, "could not connect Node HTTP Client")
	}

	m.EnvComponent.Container = clC
	m.DbC = &pgC
	m.API = clClient

	return nil
}

func (m *ClNode) getContainerRequest(secrets string) (
	*tc.ContainerRequest, error) {
	configFile, err := ioutil.TempFile("", "node_config")
	if err != nil {
		return nil, err
	}
	config, err := node.ExecuteNodeConfigTemplate(m.NodeConfigOpts)
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
		return nil, errors.New("CHAINLINK_IMAGE env must be set")
	}
	tag, ok := os.LookupEnv("CHAINLINK_VERSION")
	if !ok {
		return nil, errors.New("CHAINLINK_VERSION env must be set")
	}

	return &tc.ContainerRequest{
		Name:         m.EnvComponent.ContainerName,
		Image:        fmt.Sprintf("%s:%s", image, tag),
		ExposedPorts: []string{"6688/tcp"},
		Entrypoint: []string{"chainlink",
			"-c", configPath,
			"-s", secretsPath,
			"node", "start", "-d",
			"-p", adminCredsPath,
			"-a", apiCredsPath,
		},
		Networks: m.Networks,
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
	}, nil
}
