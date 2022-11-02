package directrequestocr_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/directrequestocr"
	"github.com/stretchr/testify/require"
)

func intToByte32(id uint) [32]byte {
	byteArr := (*[32]byte)([]byte(fmt.Sprintf("%032d\n", id)))
	return *byteArr
}

func TestDROCROrm_CreateRequestDuplicate(t *testing.T) {
	t.Parallel()

	orm := directrequestocr.NewInMemoryORM()
	id := intToByte32(420)
	txHash := common.HexToHash("0xabc")

	dbId1, err := orm.CreateRequest(id, time.Now(), &txHash)
	require.NoError(t, err)
	dbId2, err := orm.CreateRequest(id, time.Now(), &txHash)
	require.NotNil(t, err)
	require.Equal(t, dbId1, dbId2, "incorrect DBID of existing request")
}
