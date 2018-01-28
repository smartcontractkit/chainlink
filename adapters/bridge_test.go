package adapters_test

import (
	"net/url"
	"testing"

	gock "github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestBridgeAdapterPerform(t *testing.T) {
	nilString := cltest.NullString(nil)
	cases := []struct {
		name     string
		status   int
		want     null.String
		errored  bool
		response string
	}{
		{"success", 200, cltest.NullString("purchased"), false, `{"value": "purchased"}`},
		{"server error", 500, nilString, true, `{"errors": "too many"}`},
		{"JSON parse error", 200, nilString, true, `}`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()
	tt := models.NewBridgeType()
	tt.Name = "auctionBidding"
	u, err := url.Parse("https://dbay.eth/api")
	assert.Nil(t, err)
	tt.URL = models.WebURL{u}
	eb := &adapters.Bridge{tt}
	input := models.RunResultWithValue("lot 49")

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			gock.New("https://dbay.eth").
				Post("/api").
				Reply(test.status).
				JSON(test.response)
			defer cltest.CloseGock(t)

			output := eb.Perform(input, store)
			assert.Equal(t, test.want, output.Output["value"])
			assert.Equal(t, test.errored, output.HasError())
		})
	}
}
