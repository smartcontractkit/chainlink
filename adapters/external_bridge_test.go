package adapters_test

import (
	"fmt"
	"net/url"
	"testing"

	gock "github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestExternalBridgeAdapterPerform(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()
	defer cltest.CloseGock(t)

	tt := models.NewCustomTaskType()
	tt.Name = "auctionBidding"
	u, err := url.Parse("https://dbay.eth/api")
	assert.Nil(t, err)
	tt.URL = models.WebURL{u}

	eaValue := "bought!"
	eaResponse := fmt.Sprintf(`{"value": "%v"}`, eaValue)
	gock.New("https://dbay.eth").
		Post("/api").
		Reply(200).
		JSON(eaResponse)

	eb := &adapters.ExternalBridge{tt}
	input := models.RunResultWithValue("")

	output := eb.Perform(input, store)
	assert.Equal(t, eaValue, output.Value())
}
