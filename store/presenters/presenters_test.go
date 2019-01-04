package presenters_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type MI = models.Initiator
type MIP = models.InitiatorParams

func TestPresenterInitiatorHasCorrectKeys(t *testing.T) {
	t.Parallel()

	address := common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	now := time.Now()

	tests := []struct {
		i      models.Initiator
		params []string
	}{
		{MI{Type: models.InitiatorWeb}, []string{}},
		{MI{Type: models.InitiatorCron, InitiatorParams: MIP{Schedule: models.Cron("* * * * *")}}, []string{"schedule"}},
		{MI{Type: models.InitiatorRunAt, InitiatorParams: MIP{Time: models.Time{Time: now}}}, []string{"time", "ran"}},
		{MI{Type: models.InitiatorEthLog, InitiatorParams: MIP{Address: address}}, []string{"address"}},
	}

	for _, test := range tests {
		t.Run(test.i.Type, func(t *testing.T) {
			j, err := json.Marshal(presenters.Initiator{Initiator: test.i})
			assert.NoError(t, err)

			js := gjson.Parse(string(j))
			require.Equal(t, test.i.Type, js.Get("type").String())

			params := js.Get("params").Map()
			keys := []string{}
			for k := range params {
				keys = append(keys, k)
			}

			sort.Strings(keys)
			sort.Strings(test.params)
			assert.Equal(t, test.params, keys)
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
	kvs, err := presenters.ShowEthBalance(app.Store)
	assert.Error(t, err)
	assert.Len(t, kvs, 0)
}

func TestPresenterShowEthBalance_WithAccount(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getBalance", "0x0100") // 256

	assert.True(t, app.Store.KeyStore.HasAccounts())

	kvs, err := presenters.ShowEthBalance(app.Store)
	assert.NoError(t, err)
	require.Len(t, kvs, 1)
	kv := kvs[0]
	addr := cltest.GetAccountAddress(app.Store).Hex()
	want := fmt.Sprintf("ETH Balance for %v: 0.000000000000000256", addr)
	assert.Equal(t, want, kv["message"].(string))
	assert.Equal(t, kv["address"].(common.Address).String(), addr)
	assert.Equal(t, kv["balance"].(string), "0.000000000000000256")
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

func TestPresenterShowLinkBalance_WithAccount(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Register("eth_call", "0x0100") // 256

	assert.True(t, app.Store.KeyStore.HasAccounts())

	kvs, err := presenters.ShowLinkBalance(app.Store)
	assert.NoError(t, err)

	kv := kvs[0]
	addr := cltest.GetAccountAddress(app.Store).Hex()
	want := fmt.Sprintf("LINK Balance for %v: 0.000000000000000256", addr)
	assert.Equal(t, want, kv["message"].(string))
	assert.Equal(t, kv["address"].(common.Address).String(), addr)
	assert.Equal(t, kv["balance"], "0.000000000000000256")
}

func TestPresenter_FriendlyBigInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   *big.Int
		want string
	}{
		{big.NewInt(0), "#0 (0x0)"},
		{big.NewInt(1), "#1 (0x1)"},
		{big.NewInt(123456), "#123456 (0x1e240)"},
	}

	for _, test := range tests {
		t.Run(test.in.String(), func(t *testing.T) {
			assert.Equal(t, test.want, presenters.FriendlyBigInt(test.in))
		})
	}
}

func TestBridgeType_MarshalJSON(t *testing.T) {
	t.Parallel()
	input := models.BridgeType{
		Name:                   models.MustNewTaskType("hapax"),
		URL:                    cltest.WebURL("http://hap.ax"),
		Confirmations:          0,
		IncomingToken:          "123",
		OutgoingToken:          "abc",
		MinimumContractPayment: *assets.NewLink(0),
	}
	bt := presenters.BridgeType{BridgeType: input}
	output, err := bt.MarshalJSON()
	assert.NoError(t, err)
	expected := `{"name":"hapax","url":"http://hap.ax","confirmations":0,"incomingToken":"123","outgoingToken":"abc","minimumContractPayment":"0"}`
	assert.Equal(t, expected, string(output))
}

func TestServiceAgreement_MarshalJSON(t *testing.T) {
	t.Parallel()

	input := cltest.LoadJSON("../../internal/fixtures/web/hello_world_agreement.json")
	sa, err := cltest.ServiceAgreementFromString(string(input))
	assert.NoError(t, err)
	psa := presenters.ServiceAgreement{ServiceAgreement: sa}
	output, err := psa.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, cltest.NormalizedJSON(input), string(output))
}

func TestPresenter_NewConfigWhitelist_Ok(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfigWithPrivateKey("../../internal/fixtures/keys/3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea.json")
	defer cleanup()
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	cw, err := presenters.NewConfigWhitelist(store)
	assert.NoError(t, err)

	assert.Contains(t, cw.String(), "0x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea")
}

func TestPresenter_NewConfigWhitelist_Error(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig()
	defer cleanup()
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	cw, err := presenters.NewConfigWhitelist(store)
	assert.Error(t, err)
	assert.Equal(t, "", cw.AccountAddress)
}
