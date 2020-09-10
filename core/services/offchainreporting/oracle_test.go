package offchainreporting_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func TestJobSpec_Unmarshal(t *testing.T) {

}

func TestJobSpec_FetchFromDB(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	u, err := url.Parse("http://chain.link/voter_turnout/USA-2020")
	require.NoError(t, err)

	jobID := models.NewID()

	ds1 := &BridgeTask{Name: "voter_turnout"}
	ds1_parse := &JSONParseTask{
		Path:     []string{"one", "two"},
		BaseTask: BaseTask{inputTasks: []Task{ds1}},
	}
	ds1_multiply := &MultiplyTask{
		Times:    decimal.NewFromFloat(1.23),
		BaseTask: BaseTask{inputTasks: []Task{ds1_parse}},
	}
	ds2 := &HTTPTask{
		URL:         models.WebURL(*u),
		Method:      "GET",
		RequestData: HttpRequestData{"hi": "hello"},
	}
	ds2_parse := &JSONParseTask{
		Path:     []string{"three", "four"},
		BaseTask: BaseTask{inputTasks: []Task{ds2}},
	}
	ds2_multiply := &MultiplyTask{
		Times:    decimal.NewFromFloat(4.56),
		BaseTask: BaseTask{inputTasks: []Task{ds2_parse}},
	}
	answer1 := &MedianTask{
		BaseTask: BaseTask{inputTasks: []Task{ds1_multiply, ds2_multiply}},
	}
	answer2 := &BridgeTask{Name: "election_winner"}

	tasks := []Task{ds1, ds1_parse, ds1_multiply, ds2, ds2_parse, ds2_multiply, answer1, answer2}
	jobSpec := &offchainreporting.JobSpec{
		UUID:              jobID,
		ObservationSource: tasks,
	}

	err = store.ORM.RawDB(func(db *gorm.DB) error {
		result := db.Debug().Create(jobSpec.ForDB())
		require.NoError(t, result.Error)

		var returnedSpec offchainreporting.JobSpecDBRow
		err := db.Debug().
			Set("gorm:auto_preload", true).
			Find(&returnedSpec, "job_id = ?", jobSpec.JobID()).Error
		require.NoError(t, err)
		js := returnedSpec.JobSpec
		js.ObservationSource = job.UnwrapFetchersFromDB(returnedSpec.ObservationSource)[0]

		bs, _ := json.MarshalIndent(js, "", "    ")
		fmt.Println(string(bs))

		// require.Equal(t, jobSpec, returnedSpec.JobSpec)
		return nil
	})
	require.NoError(t, err)

}
