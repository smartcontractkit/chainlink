package models_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRunLog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		log         models.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{
			name:        "20190207 without indexes",
			log:         cltest.LogFromFixture(t, "testdata/requestLog20190207withoutIndexes.json"),
			wantErrored: false,
			wantData: cltest.JSONFromString(t, `{
				"url":"https://min-api.cryptocompare.com/data/price?fsym=eth&tsyms=usd,eur,jpy",
				"path":["usd"],
				"address":"0xf25186b5081ff5ce73482ad761db0eb0d25abfbf",
				"dataPrefix":"0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f80000000000000000000000000000000000000000000000000de0b6b3a76400010000000000000000000000009fbda871d559710256a2502a2517b794b482db40042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005c4a7338",
				"functionSelector":"0x4ab0d190"}`),
		},
		{
			name:        "20190207 without indexes and padded CBOR",
			log:         cltest.LogFromFixture(t, "testdata/request20200212paddedCBOR.json"),
			wantErrored: false,
			wantData: cltest.JSONFromString(t, `{
				"address":"0xfeb35e1f7abe4ef198b7c8df895e19767f3ab8a5",
				"dataprefix":"0xe947f54ec4d3cab0588684217b029cd9421ea25c59f3309bef6e8fb0d75ff5310000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000650c346f84248abc27e716ea3c6de20f7fbbdb7992cdaaf300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005e1b7f6b",
				"functionselector":"0x4ab0d190",
				"get":"https://min-api.cryptocompare.com/data/price?fsym=eth&tsyms=usd",
				"path":"usd",
				"times":100}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := models.ParseRunLog(test.log)
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.NoError(t, err)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestEthLogEvent_JSON(t *testing.T) {
	t.Parallel()

	exampleLog := cltest.LogFromFixture(t, "testdata/subscription_logs.json")
	tests := []struct {
		name        string
		el          models.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{"example", exampleLog, false, cltest.JSONResultFromFixture(t, "testdata/subscription_logs.json")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initr := models.Initiator{Type: models.InitiatorEthLog}
			le := models.InitiatorLogEvent{Initiator: initr, Log: test.el}.LogRequest()
			output, err := le.JSON()
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestStartRunOrSALogSubscription_ValidateSenders(t *testing.T) {
	requester := cltest.NewAddress()

	tests := []struct {
		name       string
		job        models.JobSpec
		requester  common.Address
		logFactory (func(*testing.T, *models.ID, common.Address, common.Address, int, string) models.Log)
		wantStatus models.RunStatus
	}{
		{
			"runlog contains valid requester",
			cltest.NewJobWithRunLogInitiator(),
			requester,
			cltest.NewRunLog,
			models.RunStatusCompleted,
		},
		{
			"runlog has wrong requester",
			cltest.NewJobWithRunLogInitiator(),
			cltest.NewAddress(),
			cltest.NewRunLog,
			models.RunStatusErrored,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, cleanup := cltest.NewApplicationWithKey(t,
				cltest.LenientEthMock,
				cltest.EthMockRegisterChainID,
				cltest.EthMockRegisterGetBalance,
			)
			defer cleanup()

			ethMock := app.EthMock
			logs := make(chan models.Log, 1)
			ethMock.Context("app.Start()", func(meth *cltest.EthMock) {
				meth.Register("eth_getTransactionCount", "0x1")
				meth.RegisterSubscription("logs", logs)
			})
			assert.NoError(t, app.StartAndConnect())

			js := test.job
			js.Initiators[0].Requesters = []common.Address{requester}
			require.NoError(t, app.AddJob(js))

			logs <- test.logFactory(t, js.ID, cltest.NewAddress(), test.requester, 1, `{}`)
			ethMock.EventuallyAllCalled(t)

			gomega.NewGomegaWithT(t).Eventually(func() []models.JobRun {
				runs, err := app.Store.JobRunsFor(js.ID)
				require.NoError(t, err)
				return runs
			}).Should(gomega.HaveLen(1))

			gomega.NewGomegaWithT(t).Eventually(func() models.RunStatus {
				runs, err := app.Store.JobRunsFor(js.ID)
				require.NoError(t, err)
				return runs[0].GetStatus()
			}).Should(gomega.Equal(test.wantStatus))
		})
	}
}

func TestFilterQueryFactory_InitiatorEthLog(t *testing.T) {
	t.Parallel()

	// When InitiatorParams.fromBlock > the fromBlock passed into FilterQueryFactory, it should win
	// due to being larger.
	{
		i := models.Initiator{
			Type: models.InitiatorEthLog,
			InitiatorParams: models.InitiatorParams{
				Address:   common.HexToAddress("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
				FromBlock: utils.NewBig(big.NewInt(123)),
				ToBlock:   utils.NewBig(big.NewInt(456)),
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
						common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
					},
				},
			},
		}
		fromBlock := big.NewInt(42)
		filter, err := models.FilterQueryFactory(i, fromBlock)
		assert.NoError(t, err)

		want := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef")},
			FromBlock: big.NewInt(123),
			ToBlock:   big.NewInt(456),
			Topics: [][]common.Hash{
				{
					common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
					common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
				},
			},
		}
		assert.Equal(t, want, filter)
	}

	// When the fromBlock passed into FilterQueryFactory > InitiatorParams.fromBlock, it should win
	// due to being larger.
	{
		i := models.Initiator{
			Type: models.InitiatorEthLog,
			InitiatorParams: models.InitiatorParams{
				Address:   common.HexToAddress("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
				FromBlock: utils.NewBig(big.NewInt(123)),
				ToBlock:   utils.NewBig(big.NewInt(456)),
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
						common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
					},
				},
			},
		}
		fromBlock := big.NewInt(124)
		filter, err := models.FilterQueryFactory(i, fromBlock)
		assert.NoError(t, err)

		want := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef")},
			FromBlock: big.NewInt(124),
			ToBlock:   big.NewInt(456),
			Topics: [][]common.Hash{
				{
					common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
					common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
				},
			},
		}
		assert.Equal(t, want, filter)
	}

	// When the winning fromBlock is > InitiatorParams.ToBlock, it should error.
	{
		i := models.Initiator{
			Type: models.InitiatorEthLog,
			InitiatorParams: models.InitiatorParams{
				Address:   common.HexToAddress("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
				FromBlock: utils.NewBig(big.NewInt(123)),
				ToBlock:   utils.NewBig(big.NewInt(456)),
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
						common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
					},
				},
			},
		}
		fromBlock := big.NewInt(999)
		_, err := models.FilterQueryFactory(i, fromBlock)
		assert.Error(t, err)
	}
}

func TestFilterQueryFactory_InitiatorRunLog(t *testing.T) {
	t.Parallel()

	id, err := models.NewIDFromString("4a1eb0e8df314cb894024a38991cff0f")
	require.NoError(t, err)
	i := models.Initiator{
		Type:      models.InitiatorRunLog,
		JobSpecID: id,
	}
	fromBlock := big.NewInt(42)
	filter, err := models.FilterQueryFactory(i, fromBlock)
	assert.NoError(t, err)

	want := ethereum.FilterQuery{
		FromBlock: fromBlock.Add(fromBlock, big.NewInt(1)),
		Topics: [][]common.Hash{
			{
				models.RunLogTopic20190207withoutIndexes,
			}, {
				common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
				common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
			},
		},
	}
	assert.Equal(t, want, filter)
}

func TestRunLogEvent_ContractPayment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		log         models.Log
		wantErrored bool
		want        *assets.Link
	}{
		{
			name:        "20190207 without indexes",
			log:         cltest.LogFromFixture(t, "testdata/requestLog20190207withoutIndexes.json"),
			wantErrored: false,
			want:        assets.NewLink(1000000000000000001),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rle := models.RunLogEvent{models.InitiatorLogEvent{Log: test.log}}

			request, err := rle.RunRequest()

			cltest.AssertError(t, test.wantErrored, err)
			assert.Equal(t, test.want, request.Payment)
		})
	}
}

func TestRunLogEvent_Requester(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		log         models.Log
		wantErrored bool
		want        common.Address
	}{
		{
			name:        "20190207 without indexes",
			log:         cltest.LogFromFixture(t, "testdata/requestLog20190207withoutIndexes.json"),
			wantErrored: false,
			want:        common.HexToAddress("0x9fbda871d559710256a2502a2517b794b482db40"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rle := models.RunLogEvent{models.InitiatorLogEvent{Log: test.log}}

			received, err := rle.Requester()
			require.NoError(t, err)

			assert.Equal(t, test.want, received)
		})
	}
}

func TestRunLogEvent_RunRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		log           models.Log
		wantRequestID common.Hash
		wantTxHash    string
		wantBlockHash string
		wantRequester common.Address
	}{
		{
			name:          "20190207 without indexes",
			log:           cltest.LogFromFixture(t, "testdata/requestLog20190207withoutIndexes.json"),
			wantRequestID: common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
			wantTxHash:    "0x04250548cd0b5d03b3bf1331aa83f32b35879440db31a6008d151260a5f3cc76",
			wantBlockHash: "0x000c0d01ce8bd7100b73b1609ababc020e7f51dac75186bb799277c6b4b71e1c",
			wantRequester: common.HexToAddress("0x9FBDa871d559710256a2502A2517b794B482Db40"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rle := models.RunLogEvent{models.InitiatorLogEvent{Log: test.log}}
			rr, err := rle.RunRequest()
			require.NoError(t, err)

			assert.Equal(t, &test.wantRequestID, rr.RequestID)
			assert.Equal(t, test.wantTxHash, rr.TxHash.Hex())
			assert.Equal(t, test.wantBlockHash, rr.BlockHash.Hex())
			assert.Equal(t, &test.wantRequester, rr.Requester)
		})
	}
}

func TestIDToTopic(t *testing.T) {
	id, err := models.NewIDFromString("ffffffffffffffffffffffffffffffff")
	require.NoError(t, err)
	assert.Equal(t, common.Hash{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, models.IDToTopic(id))
}

func TestIDToIDToHexTopic(t *testing.T) {
	id, err := models.NewIDFromString("ffffffffffffffffffffffffffffffff")
	require.NoError(t, err)
	assert.Equal(t, common.Hash{
		0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
		0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
	}, models.IDToHexTopic(id))
}
