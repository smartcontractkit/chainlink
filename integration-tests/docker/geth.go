package docker

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"strings"
	"time"
)

const (
	// RootFundingAddr is the static key that hardhat and ganache are using
	// https://hardhat.org/hardhat-runner/docs/getting-started
	// if you need more keys, keep them compatible, so we can swap Geth to Ganache/Hardhat in the future
	RootFundingAddr   = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
	RootFundingWallet = `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`
)

var (
	InitGethScript = `
#!/bin/bash
if [ ! -d /root/.ethereum/keystore ]; then
	echo "/root/.ethereum/keystore not found, running 'geth init'..."
	geth init /root/genesis.json
	echo "...done!"
fi

geth "$@"
`

	GenesisTemplate = `
{
	"config": {
	  "chainId": {{ .ChainId }},
	  "homesteadBlock": 0,
	  "eip150Block": 0,
	  "eip155Block": 0,
	  "eip158Block": 0,
	  "eip160Block": 0,
	  "byzantiumBlock": 0,
	  "constantinopleBlock": 0,
	  "petersburgBlock": 0,
	  "istanbulBlock": 0,
	  "muirGlacierBlock": 0,
	  "berlinBlock": 0,
	  "londonBlock": 0
	},
	"nonce": "0x0000000000000042",
	"mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"difficulty": "1",
	"coinbase": "0x3333333333333333333333333333333333333333",
	"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"extraData": "0x",
	"gasLimit": "8000000000",
	"alloc": {
	  "{{ .AccountAddr }}": {
		"balance": "2000000000000000000000"
	  }
	}
  }`
)

/* data types for templates */

type (
	GenesisTemplateVars struct {
		AccountAddr string
		ChainId     string
	}
)

type Geth struct {
	prefix          string
	container       tc.Container
	ExternalHttpUrl string
	InternalHttpUrl string
	ExternalWsUrl   string
	InternalWsUrl   string
}

func NewGeth(cfg any) ComponentSetupFunc {
	c := &Geth{prefix: "geth"}
	return func(network string) (Component, error) {
		return c.Start(network, cfg)
	}
}

func (m *Geth) Prefix() string {
	return m.prefix
}

func (m *Geth) Containers() []tc.Container {
	containers := make([]tc.Container, 0)
	return append(containers, m.container)
}

func (m *Geth) Start(dockerNet string, cfg any) (Component, error) {
	r, _, _, err := gethContainerRequest(dockerNet)
	if err != nil {
		return m, err
	}
	ct, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: *r,
			Started:          true,
		})
	if err != nil {
		return m, errors.Wrapf(err, "cannot start geth container")
	}
	host, err := ct.Host(context.Background())
	if err != nil {
		return m, err
	}
	httpPort, err := ct.MappedPort(context.Background(), "8544/tcp")
	if err != nil {
		return m, err
	}
	wsPort, err := ct.MappedPort(context.Background(), "8545/tcp")
	if err != nil {
		return m, err
	}
	ctName, err := ct.Name(context.Background())
	if err != nil {
		return m, err
	}
	ctName = strings.Replace(ctName, "/", "", -1)

	m.container = ct
	m.ExternalHttpUrl = fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	m.InternalHttpUrl = fmt.Sprintf("http://%s:8544", ctName)
	m.ExternalWsUrl = fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	m.InternalWsUrl = fmt.Sprintf("ws://%s:8545", ctName)

	log.Info().Str("containerName", ctName).
		Str("internalHttpUrl", m.InternalHttpUrl).
		Str("externalHttpUrl", m.ExternalHttpUrl).
		Str("externalWsUrl", m.ExternalWsUrl).
		Str("internalWsUrl", m.InternalWsUrl).
		Msg("Started Geth container")
	return m, nil
}

func (m *Geth) Stop() error {
	return m.container.Terminate(context.Background())
}

func gethContainerRequest(network string) (*tc.ContainerRequest, *keystore.KeyStore, *accounts.Account, error) {
	chainId := "1337"
	blocktime := "1"

	initScriptFile, err := ioutil.TempFile("", "init_script")
	if err != nil {
		return nil, nil, nil, err
	}
	_, err = initScriptFile.WriteString(InitGethScript)
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
	genesis, err := ExecuteTemplate(GenesisTemplate, GenesisTemplateVars{
		AccountAddr: account.Address.Hex(),
		ChainId:     chainId,
	})
	if err != nil {
		return nil, ks, &account, err
	}
	genesisFile, err := ioutil.TempFile("", "genesis_json")
	if err != nil {
		return nil, ks, &account, err
	}
	_, err = genesisFile.WriteString(genesis)
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
		Name:         fmt.Sprintf("geth-%s", uuid.NewString()),
		Image:        "ethereum/client-go:stable",
		ExposedPorts: []string{"8544/tcp", "8545/tcp"},
		Networks:     []string{network},
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
