package ocr2keeper_test

import (
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper"
)

func TestFilterNamesFromSpec21(t *testing.T) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	require.NoError(t, err)
	address := common.HexToAddress(hexutil.Encode(b))

	spec := &job.OCR2OracleSpec{
		PluginType: job.OCR2Keeper,
		ContractID: address.String(), // valid contract addr
	}

	names, err := ocr2keeper.FilterNamesFromSpec21(spec)
	require.NoError(t, err)

	assert.Len(t, names, 2)
	assert.Equal(t, logpoller.FilterName("KeepersRegistry LogProvider", address), names[0])
	assert.Equal(t, logpoller.FilterName("KeeperRegistry Events", address), names[1])

	spec = &job.OCR2OracleSpec{
		PluginType: job.OCR2Keeper,
		ContractID: "0x5431", // invalid contract addr
	}
	_, err = ocr2keeper.FilterNamesFromSpec21(spec)
	require.ErrorContains(t, err, "not a valid EIP55 formatted address")
}
