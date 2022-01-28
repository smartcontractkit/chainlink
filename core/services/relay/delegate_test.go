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
	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terradb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	"github.com/smartcontractkit/sqlx"

	chainsMock "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
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

func makeOCR2JobSpecFromToml(t *testing.T, jobSpecToml string) job.OffchainReporting2OracleSpec {
	t.Helper()

	var ocr2spec job.OffchainReporting2OracleSpec
	err := toml.Unmarshal([]byte(jobSpecToml), &ocr2spec)
	require.NoError(t, err)

	return ocr2spec
}

func TestNewOCR2Provider(t *testing.T) {
	lggr := logger.TestLogger(t)

	// setup keystore mock
	solKey := new(keystoreMock.Solana)
	solKey.On("Get", mock.AnythingOfType("string")).Return(solkey.Key{}, nil)

	// setup solana key mock
	keystore := new(keystoreMock.Master)
	keystore.On("Solana").Return(solKey, nil)

	// setup terra mocks
	terraChain := new(terraMock.Chain)
	terraChain.On("Config").Return(terra.NewConfig(terradb.ChainCfg{}, lggr))
	terraChain.On("MsgEnqueuer").Return(new(terraMock.MsgEnqueuer))
	terraChain.On("Reader", "").Return(new(terraMock.Reader), nil)

	terraChains := new(terraMock.ChainSet)
	terraChains.On("Chain", "Chainlink-99").Return(terraChain, nil)

	d := relay.NewDelegate(keystore)

	// struct for testing multiple specs
	specs := []struct {
		name string
		spec string
	}{
		// TODO: Where is EVM?
		{"solana", testspecs.OCR2SolanaSpecMinimal},
		{"terra", testspecs.OCR2TerraSpecMinimal},
	}

	for _, s := range specs {
		t.Run(s.name, func(t *testing.T) {
			spec := makeOCR2JobSpecFromToml(t, s.spec)

			_, err := d.NewOCR2Provider(uuid.UUID{}, &spec)
			require.Error(t, err)
			assert.Contains(t, strings.ToLower(err.Error()), fmt.Sprintf("no %s relay found", s.name))
		})
	}

	d.AddRelayer(relaytypes.EVM, evm.NewRelayer(&sqlx.DB{}, &chainsMock.ChainSet{}, lggr))
	d.AddRelayer(relaytypes.Solana, solana.NewRelayer(lggr))
	d.AddRelayer(relaytypes.Terra, terra.NewRelayer(lggr, terraChains))

	for _, s := range specs {
		t.Run(s.name, func(t *testing.T) {
			spec := makeOCR2JobSpecFromToml(t, s.spec)

			_, err := d.NewOCR2Provider(uuid.UUID{}, &spec)
			require.NoError(t, err)
		})
	}
}
