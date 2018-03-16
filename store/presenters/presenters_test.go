package presenters_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

type MI = models.Initiator

func TestPresenterInitiatorHasCorrectKeys(t *testing.T) {
	t.Parallel()

	address := common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	now := time.Now()

	tests := []struct {
		i    models.Initiator
		keys []string
	}{
		{MI{Type: models.InitiatorWeb}, []string{"type"}},
		{MI{Type: models.InitiatorCron, Schedule: models.Cron("* * * * *")}, []string{"type", "schedule"}},
		{MI{Type: models.InitiatorRunAt, Time: models.Time{now}}, []string{"type", "time", "ran"}},
		{MI{Type: models.InitiatorEthLog, Address: address}, []string{"type", "address"}},
	}

	for _, test := range tests {
		t.Run(test.i.Type, func(t *testing.T) {
			j, err := json.Marshal(presenters.Initiator{test.i})
			assert.Nil(t, err)

			var value map[string]interface{}
			err = json.Unmarshal(j, &value)
			assert.Nil(t, err)

			keys := utils.GetStringKeys(value)
			sort.Strings(keys)
			sort.Strings(test.keys)
			assert.Equal(t, test.keys, keys)
		})
	}
}

func TestPresenterShowEthBalance_NoAccount(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code should have panicked")
		}
	}()
	defer cleanup()
	presenters.ShowEthBalance(store)
}

func TestPresenterShowEthBalance_WithEmptyAccount(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	_, err := presenters.ShowEthBalance(app.Store)
	assert.NotNil(t, err)
}

func TestPresenterShowEthBalance_WithAccount(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getBalance", "0x0100") // 256

	assert.True(t, app.Store.KeyStore.HasAccounts())

	output, err := presenters.ShowEthBalance(app.Store)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("ETH Balance for %v: 2.56e-16", app.Store.KeyStore.GetAccount().Address.Hex()), output)
}

func TestPresenterShowLinkBalance_NoAccount(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code should have panicked")
		}
	}()
	defer cleanup()
	presenters.ShowLinkBalance(store)
}

func TestPresenterShowLinkBalance_WithEmptyAccount(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	_, err := presenters.ShowLinkBalance(app.Store)
	fmt.Println(err)
	assert.NotNil(t, err)
}

func TestPresenterShowLinkBalance_WithAccount(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Register("eth_call", "0x0100") // 256

	assert.True(t, app.Store.KeyStore.HasAccounts())

	output, err := presenters.ShowLinkBalance(app.Store)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("Link Balance for %v: 2.56e-16", app.Store.KeyStore.GetAccount().Address.Hex()), output)
}
