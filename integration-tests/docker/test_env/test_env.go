package test_env

import (
	"crypto/ed25519"
	"encoding/hex"
	"strings"
	"sync"

	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/docker"
	env "github.com/smartcontractkit/chainlink/integration-tests/types/envcommon"
	"github.com/smartcontractkit/chainlink/integration-tests/types/node"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	tc "github.com/testcontainers/testcontainers-go"
	"go.uber.org/multierr"
	"math/big"
)

type CLClusterTestEnv struct {
	cfg        *TestEnvConfig
	Network    *tc.DockerNetwork
	LogWatch   *logwatch.LogWatch
	CLNodes    []*ClNode
	Geth       *Geth
	MockServer *MockServer
}

func NewTestEnv() (*CLClusterTestEnv, error) {
	network, err := docker.CreateNetwork()
	if err != nil {
		return nil, err
	}
	networks := []string{network.Name}
	return &CLClusterTestEnv{
		Network: network,
		Geth: NewGeth(env.EnvComponentOpts{
			Networks: networks,
		}),
		MockServer: NewMockServer(env.EnvComponentOpts{
			Networks: networks,
		}),
	}, nil
}

func NewTestEnvFromCfg(cfg *TestEnvConfig) (*CLClusterTestEnv, error) {
	network, err := docker.CreateNetwork()
	if err != nil {
		return nil, err
	}
	networks := []string{network.Name}
	log.Info().Interface("Cfg", cfg).Send()
	return &CLClusterTestEnv{
		cfg:     cfg,
		Network: network,
		Geth: NewGeth(env.EnvComponentOpts{
			ReuseContainerName: cfg.Geth.ContainerName,
			Networks:           networks,
		}),
		MockServer: NewMockServer(env.EnvComponentOpts{
			ReuseContainerName: cfg.MockServer.ContainerName,
			Networks:           networks,
		}),
	}, nil
}

func (m *CLClusterTestEnv) StartGeth() error {
	return m.Geth.StartContainer(m.LogWatch)
}

func (m *CLClusterTestEnv) StartMockServer() error {
	return m.MockServer.StartContainer(m.LogWatch)
}

// StartClNodes start one bootstrap node and {count} OCR nodes
func (m *CLClusterTestEnv) StartClNodes(nodeConfigOpts node.NodeConfigOpts, count int) error {
	var wg sync.WaitGroup
	var errs = []error{}
	var mu sync.Mutex

	// Start nodes
	for i := 0; i < count; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			dbContainerName := fmt.Sprintf("cl-db-%s", uuid.NewString())
			opts := env.EnvComponentOpts{
				Networks: []string{m.Network.Name},
			}
			if m.cfg != nil {
				opts.ReuseContainerName = m.cfg.Nodes[i].NodeContainerName
				dbContainerName = m.cfg.Nodes[i].DbContainerName
			}
			n := NewClNode(opts, nodeConfigOpts, dbContainerName)
			err := n.StartContainer(m.LogWatch)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			} else {
				mu.Lock()
				m.CLNodes = append(m.CLNodes, n)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return multierr.Combine(errs...)
	}
	return nil
}

func (m *CLClusterTestEnv) GetDefaultNodeConfigOpts() node.NodeConfigOpts {
	return node.NodeConfigOpts{
		EVM: struct {
			HttpUrl string
			WsUrl   string
		}{
			HttpUrl: m.Geth.InternalHttpUrl,
			WsUrl:   m.Geth.InternalWsUrl,
		},
	}
}

// ChainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func (m *CLClusterTestEnv) ChainlinkNodeAddresses() ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, n := range m.CLNodes {
		primaryAddress, err := n.ChainlinkNodeAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, primaryAddress)
	}
	return addresses, nil
}

// FundChainlinkNodes will fund all the provided Chainlink nodes with a set amount of native currency
func (m *CLClusterTestEnv) FundChainlinkNodes(amount *big.Float) error {
	for _, cl := range m.CLNodes {
		if err := cl.Fund(m.Geth, amount); err != nil {
			return err
		}
	}
	return m.Geth.EthClient.WaitForEvents()
}

func (m *CLClusterTestEnv) GetNodeCSAKeys() ([]string, error) {
	var keys []string
	for _, n := range m.CLNodes {
		csaKeys, err := n.GetNodeCSAKeys()
		if err != nil {
			return nil, err
		}
		keys = append(keys, csaKeys.Data[0].ID)
	}
	return keys, nil
}

func getOracleIdentities(chainlinkNodes []ClNode) ([]int, []confighelper.OracleIdentityExtra) {
	S := make([]int, len(chainlinkNodes))
	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(chainlinkNodes))
	sharedSecretEncryptionPublicKeys := make([]ocrtypes.ConfigEncryptionPublicKey, len(chainlinkNodes))
	var wg sync.WaitGroup
	for i, cl := range chainlinkNodes {
		wg.Add(1)
		go func(i int, cl ClNode) error {
			defer wg.Done()

			ocr2Keys, err := cl.API.MustReadOCR2Keys()
			if err != nil {
				return err
			}
			var ocr2Config client.OCR2KeyAttributes
			for _, key := range ocr2Keys.Data {
				if key.Attributes.ChainType == string(chaintype.EVM) {
					ocr2Config = key.Attributes
					break
				}
			}

			keys, err := cl.API.MustReadP2PKeys()
			if err != nil {
				return err
			}
			p2pKeyID := keys.Data[0].Attributes.PeerID

			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_"))
			if err != nil {
				return err
			}

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			copy(offchainPkBytesFixed[:], offchainPkBytes)
			if err != nil {
				return err
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				return err
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			copy(configPkBytesFixed[:], configPkBytes)
			if err != nil {
				return err
			}

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_"))
			if err != nil {
				return err
			}

			csaKeys, _, err := cl.API.ReadCSAKeys()
			if err != nil {
				return err
			}

			sharedSecretEncryptionPublicKeys[i] = configPkBytesFixed
			oracleIdentities[i] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   ocrtypes.Account(csaKeys.Data[0].ID),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[i] = 1

			return nil
		}(i, cl)
	}
	wg.Wait()

	return S, oracleIdentities
}

func (m *CLClusterTestEnv) Terminate() error {
	// TESTCONTAINERS_RYUK_DISABLED=false by defualt so ryuk will remove all
	// the containers and the network
	return nil
}
