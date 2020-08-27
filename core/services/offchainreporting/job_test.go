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
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func TestJobSpec_FetchFromDB(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	u, err := url.Parse("http://chain.link")
	require.NoError(t, err)

	jobID := models.NewID()

	jobSpec := &offchainreporting.JobSpec{
		JobSpecID: jobID,
		ObservationSource: &job.MedianFetcher{
			Fetchers: []job.Fetcher{
				&job.HttpFetcher{
					URL:    models.WebURL(*u),
					Method: "GET",
					RequestData: map[string]interface{}{
						"one": "asdf",
						"two": "xyzzy",
					},
					Transformers: job.Transformers{
						&job.JSONParseTransformer{Path: []string{"one", "two"}},
						&job.MultiplyTransformer{Times: decimal.NewFromFloat(1.23)},
					},
				},
				&job.BridgeFetcher{
					BridgeName: "t00f4r",
					RequestData: map[string]interface{}{
						"one": "asdf",
						"two": "xyzzy",
					},
					Transformers: job.Transformers{
						&job.JSONParseTransformer{Path: []string{"one", "two"}},
						&job.MultiplyTransformer{Times: decimal.NewFromFloat(1.23)},
					},
				},
			},
		},
	}

	err = store.ORM.RawDB(func(db *gorm.DB) error {
		result := db.Debug().Create(jobSpec.ForDB())
		require.NoError(t, result.Error)

		var returnedSpec offchainreporting.JobSpecDBRow
		err := db.Debug().
			// Set("gorm:auto_preload", true).
			Preload("ObservationSource").
			Preload("ObservationSource.HttpFetcher").
			Preload("ObservationSource.HttpFetcher.Transformers").
			Preload("ObservationSource.HttpFetcher.Transformers.MultiplyTransformer").
			Preload("ObservationSource.HttpFetcher.Transformers.JSONParseTransformer").
			Preload("ObservationSource.BridgeFetcher").
			Preload("ObservationSource.BridgeFetcher.Transformers").
			Preload("ObservationSource.BridgeFetcher.Transformers.MultiplyTransformer").
			Preload("ObservationSource.BridgeFetcher.Transformers.JSONParseTransformer").
			Preload("ObservationSource.MedianFetcher").
			Preload("ObservationSource.MedianFetcher.Transformers").
			Preload("ObservationSource.MedianFetcher.Transformers.MultiplyTransformer").
			Preload("ObservationSource.MedianFetcher.Transformers.JSONParseTransformer").
			Find(&returnedSpec, "job_spec_id = ?", jobSpec.JobSpecID).Error
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
