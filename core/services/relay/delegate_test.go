package relay_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	solconfig "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	soldb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terradb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	"github.com/smartcontractkit/sqlx"

	chainsMock "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	solMock "github.com/smartcontractkit/chainlink/core/chains/solana/mocks"
	terraMock "github.com/smartcontractkit/chainlink/core/chains/terra/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	keystoreMock "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm"
	relaytypes "github.com/smartcontractkit/chainlink/core/services/relay/types"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
)

func makeOCR2JobSpecFromToml(t *testing.T, jobSpecToml string) job.OCR2OracleSpec {
	t.Helper()

	var ocr2spec job.OCR2OracleSpec
	err := toml.Unmarshal([]byte(jobSpecToml), &ocr2spec)
	require.NoError(t, err)

	return ocr2spec
}

func TestNewOCR2Provider(t *testing.T) {
	lggr := logger.TestLogger(t)

	// setup keystore mock
	solKey := new(keystoreMock.Solana)
	solKey.On("Get", "8AuzafoGEz92Z3WGFfKuEh2Ca794U3McLJBy7tfmDynK").Return(solkey.Key{}, nil).Once()

	// setup solana key mock
	keystore := new(keystoreMock.Master)
	keystore.On("Solana").Return(solKey, nil).Once()

	// setup terra mocks
	terraChain := new(terraMock.Chain)
	terraChain.On("Config").Return(terra.NewConfig(terradb.ChainCfg{}, lggr))
	terraChain.On("TxManager").Return(new(terraMock.TxManager)).Times(2)
	terraChain.On("Reader", "").Return(new(terraMock.Reader), nil).Once()
	terraChain.On("Reader", "some-test-node").Return(new(terraMock.Reader), nil).Once()

	terraChains := new(terraMock.ChainSet)
	terraChains.On("Chain", mock.Anything, "Chainlink-99").Return(terraChain, nil).Times(2)

	// set up solana mocks
	solChain := new(solMock.Chain)
	solChain.On("Config").Return(solconfig.NewConfig(soldb.ChainCfg{}, lggr))
	solChain.On("TxManager").Return(new(solMock.TxManager)).Once()
	solChain.On("Reader").Return(new(solMock.Reader), nil).Once()

	solChains := new(solMock.ChainSet)
	solChains.On("Chain", mock.Anything, "Chainlink-99").Return(solChain, nil).Once()

	d := relay.NewDelegate(keystore)

	// struct for testing multiple specs
	specs := []struct {
		name string
		spec string
	}{
		// TODO: Where is EVM?
		{"solana", testspecs.OCR2SolanaSpecMinimal},
		{"terra", testspecs.OCR2TerraSpecMinimal},
		{"terra", testspecs.OCR2TerraNodeSpecMinimal}, // nodeName: "some-test-node"
	}

	for _, s := range specs {
		t.Run(s.name, func(t *testing.T) {
			spec := makeOCR2JobSpecFromToml(t, s.spec)
			_, err := d.NewOCR2Provider(uuid.UUID{}, &relay.OCR2ProviderArgs{
				ID:              spec.ID,
				ContractID:      spec.ContractID,
				TransmitterID:   spec.TransmitterID,
				Relay:           spec.Relay,
				RelayConfig:     spec.RelayConfig,
				IsBootstrapPeer: false,
			})
			require.Error(t, err)
			assert.Contains(t, strings.ToLower(err.Error()), fmt.Sprintf("no %s relay found", s.name))
		})
	}

	d.AddRelayer(relaytypes.EVM, evm.NewRelayer(&sqlx.DB{}, &chainsMock.ChainSet{}, lggr))
	d.AddRelayer(relaytypes.Solana, solana.NewRelayer(lggr, solChains))
	d.AddRelayer(relaytypes.Terra, terra.NewRelayer(lggr, terraChains))

	for _, s := range specs {
		t.Run(s.name, func(t *testing.T) {
			spec := makeOCR2JobSpecFromToml(t, s.spec)
			_, err := d.NewOCR2Provider(uuid.UUID{}, &relay.OCR2ProviderArgs{
				ID:              spec.ID,
				ContractID:      spec.ContractID,
				TransmitterID:   spec.TransmitterID,
				Relay:           spec.Relay,
				RelayConfig:     spec.RelayConfig,
				IsBootstrapPeer: false,
			})
			require.NoError(t, err)
		})
	}

	keystore.AssertExpectations(t)
	solKey.AssertExpectations(t)
	terraChains.AssertExpectations(t)
	terraChain.AssertExpectations(t)
	solChains.AssertExpectations(t)
	solChain.AssertExpectations(t)
}
