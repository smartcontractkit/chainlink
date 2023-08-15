package test_env

import (
	"context"
	"fmt"
	"time"

	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/utils/templates"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

const (
	// RootFundingAddr is the static key that hardhat is using
	// https://hardhat.org/hardhat-runner/docs/getting-started
	// if you need more keys, keep them compatible, so we can swap Geth to Ganache/Hardhat in the future
	RootFundingAddr   = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
	RootFundingWallet = `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`
)

type Geth struct {
	EnvComponent
	ExternalHttpUrl string
	InternalHttpUrl string
	ExternalWsUrl   string
	InternalWsUrl   string
}

func NewGeth(networks []string, opts ...EnvComponentOption) *Geth {
	g := &Geth{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "geth", uuid.NewString()[0:8]),
			Networks:      networks,
		},
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *Geth) StartContainer() (blockchain.EVMNetwork, InternalDockerUrls, error) {
	r, _, _, err := g.getGethContainerRequest(g.Networks)
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	ct, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: *r,
			Started:          true,
			Reuse:            true,
		})
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, errors.Wrapf(err, "cannot start geth container")
	}
	host, err := ct.Host(context.Background())
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	httpPort, err := ct.MappedPort(context.Background(), "8544/tcp")
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	wsPort, err := ct.MappedPort(context.Background(), "8545/tcp")
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}

	g.Container = ct
	g.ExternalHttpUrl = fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	g.InternalHttpUrl = fmt.Sprintf("http://%s:8544", g.ContainerName)
	g.ExternalWsUrl = fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	g.InternalWsUrl = fmt.Sprintf("ws://%s:8545", g.ContainerName)

	networkConfig := blockchain.SimulatedEVMNetwork
	networkConfig.Name = "geth"
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	internalDockerUrls := InternalDockerUrls{
		HttpUrl: g.InternalHttpUrl,
		WsUrl:   g.InternalWsUrl,
	}

	log.Info().Str("containerName", g.ContainerName).
		Str("internalHttpUrl", g.InternalHttpUrl).
		Str("externalHttpUrl", g.ExternalHttpUrl).
		Str("externalWsUrl", g.ExternalWsUrl).
		Str("internalWsUrl", g.InternalWsUrl).
		Msg("Started Geth container")

	return networkConfig, internalDockerUrls, nil
}

func (g *Geth) getGethContainerRequest(networks []string) (*tc.ContainerRequest, *keystore.KeyStore, *accounts.Account, error) {
	chainId := "1337"
	blocktime := "1"

	initScriptFile, err := os.CreateTemp("", "init_script")
	if err != nil {
		return nil, nil, nil, err
	}
	_, err = initScriptFile.WriteString(templates.InitGethScript)
	if err != nil {
		return nil, nil, nil, err
	}
	keystoreDir, err := os.MkdirTemp("", "keystore")
	if err != nil {
		return nil, nil, nil, err
	}
	// Create keystore and ethereum account
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount("")
	if err != nil {
		return nil, ks, &account, err
	}
	genesisJsonStr, err := templates.GenesisJsonTemplate{
		ChainId:     chainId,
		AccountAddr: account.Address.Hex(),
	}.String()
	if err != nil {
		return nil, ks, &account, err
	}
	genesisFile, err := os.CreateTemp("", "genesis_json")
	if err != nil {
		return nil, ks, &account, err
	}
	_, err = genesisFile.WriteString(genesisJsonStr)
	if err != nil {
		return nil, ks, &account, err
	}
	key1File, err := os.CreateTemp(keystoreDir, "key1")
	if err != nil {
		return nil, ks, &account, err
	}
	_, err = key1File.WriteString(RootFundingWallet)
	if err != nil {
		return nil, ks, &account, err
	}
	configDir, err := os.MkdirTemp("", "config")
	if err != nil {
		return nil, ks, &account, err
	}
	err = os.WriteFile(configDir+"/password.txt", []byte(""), 0600)
	if err != nil {
		return nil, ks, &account, err
	}

	return &tc.ContainerRequest{
		Name:         g.ContainerName,
		Image:        "ethereum/client-go:stable",
		ExposedPorts: []string{"8544/tcp", "8545/tcp"},
		Networks:     networks,
		WaitingFor: tcwait.ForHTTP("/").
			WithPort("8544/tcp").
			WithStartupTimeout(120 * time.Second).
			WithPollInterval(1 * time.Second),
		Entrypoint: []string{"sh", "./root/init.sh",
			"--dev",
			"--password", "/root/config/password.txt",
			"--datadir",
			"/root/.ethereum/devchain",
			"--unlock",
			RootFundingAddr,
			"--mine",
			"--miner.etherbase",
			RootFundingAddr,
			"--ipcdisable",
			"--http",
			"--http.vhosts",
			"*",
			"--http.addr",
			"0.0.0.0",
			"--http.port=8544",
			"--ws",
			"--ws.origins",
			"*",
			"--ws.addr",
			"0.0.0.0",
			"--ws.port=8545",
			"--graphql",
			"-graphql.corsdomain",
			"*",
			"--allow-insecure-unlock",
			"--rpc.allow-unprotected-txs",
			"--http.api",
			"eth,web3,debug",
			"--http.corsdomain",
			"*",
			"--vmdebug",
			fmt.Sprintf("--networkid=%s", chainId),
			"--rpc.txfeecap",
			"0",
			"--dev.period",
			blocktime,
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      initScriptFile.Name(),
				ContainerFilePath: "/root/init.sh",
				FileMode:          0644,
			},
			{
				HostFilePath:      genesisFile.Name(),
				ContainerFilePath: "/root/genesis.json",
				FileMode:          0644,
			},
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: keystoreDir,
				},
				Target: "/root/.ethereum/devchain/keystore/",
			},
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: configDir,
				},
				Target: "/root/config/",
			},
		},
	}, ks, &account, nil
}
