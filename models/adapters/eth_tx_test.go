package adapters_test

import (
	"testing"

	gock "github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models/adapters"
	"github.com/stretchr/testify/assert"
)

func TestSendingEthereumTx(t *testing.T) {
	defer cltest.CloseGock(t)

	address := "0x1234567890"
	fid := "0x12345678"
	value := "0000abcdef"
	input := adapters.RunResultWithValue(value)

	response := `{"result": "0x0100"}`
	gock.New("http://example.com").
		Post("/api").
		Reply(200).
		JSON(response)

	adapter := adapters.EthSendTx{
		Address:    address,
		FunctionID: fid,
	}
	result := adapter.Perform(input)
	assert.Equal(t, "0x0100", result.Value())
}
