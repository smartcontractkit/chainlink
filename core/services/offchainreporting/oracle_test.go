package offchainreporting_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/url"
// 	"testing"
// 	"time"

// 	"github.com/BurntSushi/toml"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/jinzhu/gorm"
// 	"github.com/libp2p/go-libp2p-core/peer"
// 	"github.com/shopspring/decimal"
// 	"github.com/stretchr/testify/require"

// 	"github.com/smartcontractkit/chainlink/core/internal/cltest"
// 	"github.com/smartcontractkit/chainlink/core/services/job"
// 	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
// 	"github.com/smartcontractkit/chainlink/core/services/pipeline"
// 	"github.com/smartcontractkit/chainlink/core/store/models"
// )

// func TestJobSpec_Unmarshal(t *testing.T) {
// 	bs, err := ioutil.ReadFile("./example-job-spec.toml")
// 	require.NoError(t, err)

// 	var spec offchainreporting.OracleSpec
// 	err = toml.Unmarshal(bs, &spec)
// 	require.NoError(t, err)

// }

// func TestJobSpec_FetchFromDB(t *testing.T) {
// 	store, cleanup := cltest.NewStore(t)
// 	defer cleanup()

// 	u, err := url.Parse("http://chain.link/voter_turnout/USA-2020")
// 	require.NoError(t, err)

// 	answer1 := &pipeline.MedianTask{
// 		BaseTask: pipeline.NewBaseTask("answer1", nil, 0),
// 	}
// 	answer2 := &pipeline.BridgeTask{
// 		Name:     "election_winner",
// 		BaseTask: pipeline.NewBaseTask("answer2", nil, 1),
// 	}
// 	ds1_multiply := &pipeline.MultiplyTask{
// 		Times:    decimal.NewFromFloat(1.23),
// 		BaseTask: pipeline.NewBaseTask("ds1_multiply", answer1, 0),
// 	}
// 	ds1_parse := &pipeline.JSONParseTask{
// 		Path:     []string{"one", "two"},
// 		BaseTask: pipeline.NewBaseTask("ds1_parse", ds1_multiply, 0),
// 	}
// 	ds1 := &pipeline.BridgeTask{
// 		Name:     "voter_turnout",
// 		BaseTask: pipeline.NewBaseTask("ds1", ds1_parse, 0),
// 	}
// 	ds2_multiply := &pipeline.MultiplyTask{
// 		Times:    decimal.NewFromFloat(4.56),
// 		BaseTask: pipeline.NewBaseTask("ds2_multiply", answer1, 0),
// 	}
// 	ds2_parse := &pipeline.JSONParseTask{
// 		Path:     []string{"three", "four"},
// 		BaseTask: pipeline.NewBaseTask("ds2_parse", ds2_multiply, 0),
// 	}
// 	ds2 := &pipeline.HTTPTask{
// 		URL:         models.WebURL(*u),
// 		Method:      "GET",
// 		RequestData: pipeline.HttpRequestData{"hi": "hello"},
// 		BaseTask:    pipeline.NewBaseTask("ds2", ds2_parse, 0),
// 	}

// 	peerID, err := peer.IDFromString("16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju")
// 	require.NoError(t, err)

// 	tasks := []Task{ds1, ds1_parse, ds1_multiply, ds2, ds2_parse, ds2_multiply, answer1, answer2}
// 	jobSpec := &offchainreporting.OracleSpec{
// 		Pipeline: tasks,
// 		OffchainReportingOracleSpec: models.OffchainReportingOracleSpec{
// 			ContractAddress:                        common.HexToAddress("0x613a38AC1659769640aaE063C651F48E0250454C"),
// 			P2PPeerID:                              peerID,
// 			P2PBootstrapPeers:                      []string{"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"},
// 			IsBootstrapPeer:                        false,
// 			EncryptedOCRKeyBundleID:                models.MustSha256HashFromHex("deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
// 			MonitoringEndpoint:                     "chain.link:4321",
// 			TransmitterAddress:                     common.HexToAddress("0x613a38AC1659769640aaE063C651F48E0250454C"),
// 			ObservationTimeout:                     models.Interval(10 * time.Second),
// 			BlockchainTimeout:                      models.Interval(20 * time.Second),
// 			ContractConfigTrackerSubscribeInterval: models.Interval(2 * time.Minute),
// 			ContractConfigTrackerPollInterval:      models.Interval(1 * time.Minute),
// 			ContractConfigConfirmations:            3,
// 		},
// 	}

// 	err = store.ORM.RawDB(func(db *gorm.DB) error {
// 		result := db.Create(jobSpec.ForDB())
// 		require.NoError(t, result.Error)

// 		var returnedSpec offchainreporting.JobSpecDBRow
// 		err := db.Debug().
// 			Set("gorm:auto_preload", true).
// 			Find(&returnedSpec, "job_id = ?", jobSpec.JobID()).Error
// 		require.NoError(t, err)
// 		js := returnedSpec.JobSpec
// 		js.ObservationSource = job.UnwrapFetchersFromDB(returnedSpec.ObservationSource)[0]

// 		bs, _ := json.MarshalIndent(js, "", "    ")
// 		fmt.Println(string(bs))

// 		// require.Equal(t, jobSpec, returnedSpec.JobSpec)
// 		return nil
// 	})
// 	require.NoError(t, err)

// }
