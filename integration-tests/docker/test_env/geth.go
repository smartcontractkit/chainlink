package test_env

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	env "github.com/smartcontractkit/chainlink/integration-tests/types/envcommon"
	"github.com/smartcontractkit/chainlink/integration-tests/utils/templates"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"time"
)

const (
	// RootFundingAddr is the static key that hardhat is using
	// https://hardhat.org/hardhat-runner/docs/getting-started
	// if you need more keys, keep them compatible, so we can swap Geth to Ganache/Hardhat in the future
	RootFundingAddr   = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
	RootFundingWallet = `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`
)

type Geth struct {
	env.EnvComponent
	ExternalHttpUrl  string
	InternalHttpUrl  string
	ExternalWsUrl    string
	InternalWsUrl    string
	EthClient        blockchain.EVMClient
	ContractDeployer contracts.ContractDeployer
}

func NewGeth(compOpts env.EnvComponentOpts) *Geth {
	return &Geth{
		EnvComponent: env.NewEnvComponent("geth", compOpts),
	}
}

func (m *Geth) StartContainer(lw *logwatch.LogWatch) error {
	r, _, _, err := m.getGethContainerRequest(m.Networks)
	if err != nil {
		return err
	}
	ct, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: *r,
			Started:          true,
			Reuse:            true,
		})
	if err != nil {
		return errors.Wrapf(err, "cannot start geth container")
	}
	if lw != nil {
		if err := lw.ConnectContainer(context.Background(), ct, "geth", true); err != nil {
			return err
		}
	}
	host, err := ct.Host(context.Background())
	if err != nil {
		return err
	}
	httpPort, err := ct.MappedPort(context.Background(), "8544/tcp")
	if err != nil {
		return err
	}
	wsPort, err := ct.MappedPort(context.Background(), "8545/tcp")
	if err != nil {
		return err
	}
	ctName, err := ct.Name(context.Background())
	if err != nil {
		return err
	}
	ctName = strings.Replace(ctName, "/", "", -1)

	m.EnvComponent.Container = ct
	m.ExternalHttpUrl = fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	m.InternalHttpUrl = fmt.Sprintf("http://%s:8544", ctName)
	m.ExternalWsUrl = fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	m.InternalWsUrl = fmt.Sprintf("ws://%s:8545", ctName)

	networkConfig := blockchain.SimulatedEVMNetwork
	networkConfig.Name = "geth"
	networkConfig.URLs = []string{m.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{m.ExternalWsUrl}

	bc, err := blockchain.NewEVMClientFromNetwork(networkConfig)
	if err != nil {
		return err
	}
	m.EthClient = bc
	cd, err := contracts.NewContractDeployer(bc)
	if err != nil {
		return err
	}
	m.ContractDeployer = cd

	log.Info().Str("containerName", ctName).
		Str("internalHttpUrl", m.InternalHttpUrl).
		Str("externalHttpUrl", m.ExternalHttpUrl).
		Str("externalWsUrl", m.ExternalWsUrl).
		Str("internalWsUrl", m.InternalWsUrl).
		Msg("Started Geth container")

	return nil
}

func (m *Geth) getGethContainerRequest(networks []string) (*tc.ContainerRequest, *keystore.KeyStore, *accounts.Account, error) {
	chainId := "1337"
	blocktime := "1"

	initScriptFile, err := ioutil.TempFile("", "init_script")
	if err != nil {
		return nil, nil, nil, err
	}
	_, err = initScriptFile.WriteString(templates.InitGethScript)
	if err != nil {
		return nil, nil, nil, err
	}
	keystoreDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		return nil, nil, nil, err
	}
	// Create keystore and ethereum account
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount("")
	if err != nil {
		return nil, ks, &account, err
	}
	genesisJsonStr, err := templates.BuildGenesisJson(chainId, account.Address.Hex())
	if err != nil {
		return nil, ks, &account, err
	}
	genesisFile, err := ioutil.TempFile("", "genesis_json")
	if err != nil {
		return nil, ks, &account, err
	}
	_, err = genesisFile.WriteString(genesisJsonStr)
	if err != nil {
		return nil, ks, &account, err
	}
	key1File, err := ioutil.TempFile(keystoreDir, "key1")
	if err != nil {
		return nil, ks, &account, err
	}
	_, err = key1File.WriteString(RootFundingWallet)
	if err != nil {
		return nil, ks, &account, err
	}
	configDir, err := ioutil.TempDir("", "config")
	if err != nil {
		return nil, ks, &account, err
	}
	err = ioutil.WriteFile(configDir+"/password.txt", []byte(""), 0644)
	if err != nil {
		return nil, ks, &account, err
	}

	return &tc.ContainerRequest{
		Name:         m.ContainerName,
		Image:        "ethereum/client-go:stable",
		ExposedPorts: []string{"8544/tcp", "8545/tcp"},
		Networks:     networks,
		WaitingFor: tcwait.ForLog("Commit new sealing work").
			WithStartupTimeout(999 * time.Second).
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
