package relay_test

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	chainsMock "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	keystoreMock "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
)

var sampleSolanaSpec = `type = "offchainreporting2"
schemaVersion = 1
name = "local testing job"
contractID = "VT3AvPr2nyE9Kr7ydDXVvgvJXyBr9tHA5hd6a1GBGBx"
isBootstrapPeer = false
p2pBootstrapPeers = []
relay = "solana"
transmitterID = "8AuzafoGEz92Z3WGFfKuEh2Ca794U3McLJBy7tfmDynK"
observationSource = """
"""
juelsPerFeeCoinSource = """
"""

[relayConfig]
nodeEndpointHTTP = "http://127.0.0.1:8899"
ocr2ProgramID = "CF13pnKGJ1WJZeEgVAtFdUi4MMndXm9hneiHs8azUaZt"
storeProgramID = "A7Jh2nb1hZHwqEofm4N8SXbKTj82rx7KUfjParQXUyMQ"
transmissionsID = "J6RRmA39u8ZBwrMvRPrJA3LMdg73trb6Qhfo8vjSeadg"`

func makeOCR2JobSpecFromToml(t *testing.T, jobSpecToml string) job.OffchainReporting2OracleSpec {
	t.Helper()

	var ocr2spec job.OffchainReporting2OracleSpec
	err := toml.Unmarshal([]byte(jobSpecToml), &ocr2spec)
	require.NoError(t, err)

	return ocr2spec
}

func TestNewOCR2Provider(t *testing.T) {
	solKey := new(keystoreMock.Solana)
	solKey.On("Get", mock.AnythingOfType("string")).Return(solkey.Key{}, nil)

	keystore := new(keystoreMock.Master)
	keystore.On("Solana").Return(solKey, nil)

	d := relay.NewDelegate(&sqlx.DB{}, keystore, &chainsMock.ChainSet{}, logger.NewLogger())

	spec := makeOCR2JobSpecFromToml(t, sampleSolanaSpec)
	fmt.Printf("%+v\n", spec)

	_, err := d.NewOCR2Provider(uuid.UUID{}, &spec)
	require.NoError(t, err)

}
