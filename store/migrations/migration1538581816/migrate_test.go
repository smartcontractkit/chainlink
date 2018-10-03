package migration1538581816_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1538581816"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1538581816/old"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate1538581816_AddMinimumContractPaymentToBridgeType(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	obt := old.BridgeType{
		Name:          "my_bridge",
		URL:           old.WebURL(cltest.WebURL("http://mybridge.com")),
		Confirmations: 10,
		IncomingToken: "incoming",
		OutgoingToken: "outgoing",
	}
	assert.NoError(t, store.Save(&obt))

	migration := migration1538581816.Migration{}
	require.NoError(t, migration.Migrate(store.ORM))

	var nbt migration1538581816.BridgeType
	require.NoError(t, store.One("Name", obt.Name, &nbt))

	require.Equal(t, *assets.NewLink(0), nbt.MinimumContractPayment)
	require.Equal(t, migration1538581816.TaskType("my_bridge"), nbt.Name)
	require.Equal(t, migration1538581816.WebURL(cltest.WebURL("http://mybridge.com")), nbt.URL)
	require.Equal(t, uint64(10), nbt.Confirmations)
	require.Equal(t, "incoming", nbt.IncomingToken)
	require.Equal(t, "outgoing", nbt.OutgoingToken)
}
