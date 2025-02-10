package presenters

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
)

func TestEVMForwarderResource(t *testing.T) {
	var (
		ID        = int64(1)
		address   = utils.RandomAddress()
		chainID   = *big.NewI(4)
		createdAt = time.Now()
		updatedAt = time.Now().Add(time.Second)
	)
	fwd := forwarders.Forwarder{
		ID:         ID,
		Address:    address,
		EVMChainID: chainID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	r := NewEVMForwarderResource(fwd)
	assert.Equal(t, fmt.Sprint(ID), r.ID)
	assert.Equal(t, address, r.Address)
	assert.Equal(t, chainID, r.EVMChainID)
	assert.Equal(t, createdAt, r.CreatedAt)
	assert.Equal(t, updatedAt, r.UpdatedAt)

	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	createdAtMarshalled, err := createdAt.MarshalText()
	require.NoError(t, err)
	updatedAtMarshalled, err := updatedAt.MarshalText()
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
	   "data":{
		  "type":"evm_forwarder",
		  "id":"%d",
		  "attributes":{
			 "address":"%s",
			 "evmChainId":"%s",
			 "createdAt":"%s",
			 "updatedAt":"%s"
		  }
	   }
	}
	`, ID, strings.ToLower(address.String()), chainID.String(), string(createdAtMarshalled), string(updatedAtMarshalled))
	assert.JSONEq(t, expected, string(b))
}
