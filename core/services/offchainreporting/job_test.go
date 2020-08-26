package offchainreporting_test

import (
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

	jobSpec := &offchainreporting.JobSpec{
		ID: models.NewID(),
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

	fmt.Println("1", jobSpec.ID)

	err = store.ORM.RawDB(func(db *gorm.DB) error {
		result := db.Debug().Save(jobSpec.ForDB())
		require.NoError(t, result.Error)

		fmt.Println("2", jobSpec.ID)

		var returnedSpec offchainreporting.JobSpecDBRow
		err := db.Debug().Find(&returnedSpec, "id = ?", jobSpec.ID).Error
		require.NoError(t, err)

		fmt.Println("3", returnedSpec.ID)

		require.Equal(t, jobSpec, returnedSpec.JobSpec)
		return nil
	})
	require.NoError(t, err)

}
