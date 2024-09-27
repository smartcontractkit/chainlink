package rollups

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type TestDAOracle struct {
	toml.DAOracle
}

func (d *TestDAOracle) OracleType() toml.OracleType {
	return d.DAOracle.OracleType
}

func (d *TestDAOracle) OracleAddress() *types.EIP55Address {
	return d.DAOracle.OracleAddress
}

func (d *TestDAOracle) CustomGasPriceAPICalldata() string {
	return d.DAOracle.CustomGasPriceAPICalldata
}

func CreateTestDAOracle(t *testing.T, oracleType toml.OracleType, oracleAddress string, customGasPriceAPICalldata string) *TestDAOracle {
	oracleAddr, err := types.NewEIP55Address(oracleAddress)
	require.NoError(t, err)

	return &TestDAOracle{
		DAOracle: toml.DAOracle{
			OracleType:                oracleType,
			OracleAddress:             &oracleAddr,
			CustomGasPriceAPICalldata: customGasPriceAPICalldata,
		},
	}
}
