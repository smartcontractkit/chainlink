package pipeline_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

// TODO: Add a test for multiple terminal tasks after __result__ is deprecated
// https://www.pivotaltracker.com/story/show/176557536
func Test_PipelineRunner_ExecuteTaskRuns(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)

	// 1. Setup bridge
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700)))
	defer s1.Close()

	bridgeFeedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)
	bridgeFeedWebURL := (*models.WebURL)(bridgeFeedURL)

	_, bridge := cltest.NewBridgeType(t, "example-bridge")
	bridge.URL = *bridgeFeedWebURL
	require.NoError(t, store.ORM.DB.Create(&bridge).Error)

	// 2. Setup success HTTP
	s2 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9600)))
	defer s2.Close()

	s4 := httptest.NewServer(fakeStringResponder(t, "foo-index-1"))
	defer s4.Close()
	s5 := httptest.NewServer(fakeStringResponder(t, "bar-index-2"))
	defer s5.Close()

	orm := new(mocks.ORM)
	orm.On("DB").Return(store.DB)

	r := pipeline.NewRunner(orm, store.Config)

	spec := pipeline.Spec{ID: 142}
	taskRuns := []pipeline.TaskRun{
		// 1. Bridge request, succeeds
		pipeline.TaskRun{
			ID: 10,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           1,
				DotID:        `ds1`,
				Type:         "bridge",
				JSON:         cltest.MustNewJSONSerializable(t, `{"name": "example-bridge", "Timeout": 0, "requestData": {"data": {"coin": "BTC", "market": "USD"}}}`),
				SuccessorID:  null.IntFrom(2),
				PipelineSpec: spec,
			},
		},
		pipeline.TaskRun{
			ID: 11,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           2,
				DotID:        `ds1_parse`,
				Type:         "jsonparse",
				JSON:         cltest.MustNewJSONSerializable(t, `{"Lax": false, "path": ["data", "result"], "Timeout": 0}`),
				SuccessorID:  null.IntFrom(3),
				PipelineSpec: spec,
			},
		},
		pipeline.TaskRun{
			ID: 12,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           3,
				DotID:        `ds1_multiply`,
				Type:         "multiply",
				JSON:         cltest.MustNewJSONSerializable(t, `{"times": "1000000000000000000", "Timeout": 0}`),
				SuccessorID:  null.IntFrom(102),
				PipelineSpec: spec,
			},
		},
		// 2. HTTP request, succeeds
		pipeline.TaskRun{
			ID: 21,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           33,
				DotID:        `ds2`,
				Type:         "http",
				JSON:         cltest.MustNewJSONSerializable(t, fmt.Sprintf(`{"method": "GET", "url": "%s", "requestData": {"data": {"coin": "BTC", "market": "USD"}}}`, s2.URL)),
				SuccessorID:  null.IntFrom(32),
				PipelineSpec: spec,
			},
		},
		pipeline.TaskRun{
			ID: 22,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           32,
				DotID:        `ds2_parse`,
				Type:         "jsonparse",
				JSON:         cltest.MustNewJSONSerializable(t, `{"Lax": false, "path": ["data", "result"], "Timeout": 0}`),
				SuccessorID:  null.IntFrom(31),
				PipelineSpec: spec,
			},
		},
		pipeline.TaskRun{
			ID: 23,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           31,
				DotID:        `ds2_multiply`,
				Type:         "multiply",
				JSON:         cltest.MustNewJSONSerializable(t, `{"times": "1000000000000000000", "Timeout": 0}`),
				SuccessorID:  null.IntFrom(102),
				PipelineSpec: spec,
			},
		},
		// 3. HTTP request, fails
		pipeline.TaskRun{
			ID: 41,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           51,
				DotID:        `ds3`,
				Type:         "http",
				JSON:         cltest.MustNewJSONSerializable(t, `{"method": "GET", "url": "blah://test.invalid", "requestData": {"data": {"coin": "BTC", "market": "USD"}}}`),
				SuccessorID:  null.IntFrom(52),
				PipelineSpec: spec,
			},
		},
		pipeline.TaskRun{
			ID: 42,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           52,
				DotID:        `ds3_parse`,
				Type:         "jsonparse",
				JSON:         cltest.MustNewJSONSerializable(t, `{"Lax": false, "path": ["data", "result"], "Timeout": 0}`),
				SuccessorID:  null.IntFrom(53),
				PipelineSpec: spec,
			},
		},
		pipeline.TaskRun{
			ID: 43,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           53,
				DotID:        `ds3_multiply`,
				Type:         "multiply",
				JSON:         cltest.MustNewJSONSerializable(t, `{"times": "1000000000000000000", "Timeout": 0}`),
				SuccessorID:  null.IntFrom(102),
				PipelineSpec: spec,
			},
		},
		// MEDIAN
		pipeline.TaskRun{
			ID: 30,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           102,
				DotID:        `median`,
				Type:         "median",
				JSON:         cltest.MustNewJSONSerializable(t, `{"allowedFaults": 1}`),
				SuccessorID:  null.IntFrom(203),
				Index:        0,
				PipelineSpec: spec,
			},
		},
		// 4. HTTP Request, side by side with median to test indexing
		pipeline.TaskRun{
			ID: 71,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           72,
				DotID:        `ds4`,
				Type:         "http",
				JSON:         cltest.MustNewJSONSerializable(t, fmt.Sprintf(`{"method": "GET", "url": "%s"}`, s4.URL)),
				SuccessorID:  null.IntFrom(203),
				Index:        1,
				PipelineSpec: spec,
			},
		},
		// 5. HTTP Request, side by side with median to test indexing
		pipeline.TaskRun{
			ID: 73,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           74,
				DotID:        `ds5`,
				Type:         "http",
				JSON:         cltest.MustNewJSONSerializable(t, fmt.Sprintf(`{"method": "GET", "url": "%s"}`, s5.URL)),
				SuccessorID:  null.IntFrom(203),
				Index:        2,
				PipelineSpec: spec,
			},
		},
		// 6. Result
		pipeline.TaskRun{
			ID: 13,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:           203,
				DotID:        `__result__`,
				Type:         "result",
				JSON:         cltest.MustNewJSONSerializable(t, `{}`),
				SuccessorID:  null.Int{},
				PipelineSpec: spec,
			},
		},
	}

	run := pipeline.Run{
		ID:               242,
		PipelineSpec:     spec,
		PipelineTaskRuns: taskRuns,
	}

	trrs, err := r.ExecuteRun(context.Background(), run, *logger.Default)
	require.NoError(t, err)

	require.Len(t, trrs, len(taskRuns))

	var finalResults []pipeline.Result
	for _, trr := range trrs {
		if trr.IsTerminal {
			finalResults = append(finalResults, trr.Result)
		}
	}

	require.Len(t, finalResults, 1)
	result := finalResults[0]

	require.Len(t, result.Value, 3)
	finalValues := result.Value.([]interface{})

	{
		// Median
		finalValue := finalValues[0].(decimal.Decimal)
		require.Equal(t, "9650000000000000000000", finalValue.String())

	}

	{
		// Strings 1 and 2
		require.Equal(t, "foo-index-1", finalValues[1].(string))
		require.Equal(t, "bar-index-2", finalValues[2].(string))
	}

	require.Len(t, result.Error, 3)
	finalError := result.Error.(pipeline.FinalErrors)
	require.False(t, finalError.HasErrors())

	var errorResults []pipeline.TaskRunResult
	for _, trr := range trrs {
		if trr.Result.Error != nil && !trr.IsTerminal {
			errorResults = append(errorResults, trr)
		}
	}
	// There are three tasks in the erroring pipeline
	require.Len(t, errorResults, 3)
}

func dotGraphToSpec(t *testing.T, id int32, taskIDStart int32, graph string) pipeline.Spec {
	d := pipeline.NewTaskDAG()
	err := d.UnmarshalText([]byte(graph))
	require.NoError(t, err)
	ts, err := d.TasksInDependencyOrder()
	require.NoError(t, err)
	var s = pipeline.Spec{
		ID:                id,
		PipelineTaskSpecs: make([]pipeline.TaskSpec, 0),
	}
	taskSpecIDs := make(map[pipeline.Task]int32)
	for _, task := range ts {
		var successorID null.Int
		if task.OutputTask() != nil {
			successor := task.OutputTask()
			successorID = null.IntFrom(int64(taskSpecIDs[successor]))
		}
		v := pipeline.JSONSerializable{task, false}
		b, err := v.MarshalJSON()
		require.NoError(t, err)
		v2 := pipeline.JSONSerializable{}
		err = v2.UnmarshalJSON(b)
		require.NoError(t, err)
		s.PipelineTaskSpecs = append(s.PipelineTaskSpecs, pipeline.TaskSpec{
			ID:             taskIDStart,
			DotID:          task.DotID(),
			PipelineSpecID: s.ID,
			Type:           task.Type(),
			JSON:           v2,
			Index:          task.OutputIndex(),
			SuccessorID:    successorID,
		})
		taskSpecIDs[task] = taskIDStart
		taskIDStart++
	}
	return s
}

func Test_PipelineRunner_HandleFaults(t *testing.T) {
	// We want to test the scenario where one or multiple APIs time out,
	// but a sufficient number of them still complete within the desired time frame
	// and so we can still obtain a median.
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	orm := new(mocks.ORM)
	orm.On("DB").Return(store.DB)
	m1 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(100 * time.Millisecond)
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"result":10}`))
	}))
	m2 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"result":11}`))
	}))
	s := dotGraphToSpec(t, 1, 1, fmt.Sprintf(`
ds1          [type=http url="%s"];
ds1_parse    [type=jsonparse path="result"];
ds1_multiply [type=multiply times=100];

ds2          [type=http url="%s"];
ds2_parse    [type=jsonparse path="result"];
ds2_multiply [type=multiply times=100];

ds1 -> ds1_parse -> ds1_multiply -> answer1;
ds2 -> ds2_parse -> ds2_multiply -> answer1;

answer1 [type=median                      index=0];
`, m1.URL, m2.URL))

	r := pipeline.NewRunner(orm, store.Config)
	run, err := pipeline.NewRun(s, time.Now())
	require.NoError(t, err)

	// If we cancel before an API is finished, we should still get a median.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	trrs, err := r.ExecuteRun(ctx, run, *logger.Default)
	require.NoError(t, err)
	for _, trr := range trrs {
		if trr.IsTerminal {
			require.Equal(t, decimal.RequireFromString("1100"), trr.Result.Value)
		}
	}
}
