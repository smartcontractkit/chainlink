package services_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
)

func TestValidateBridgeType(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	tests := []struct {
		description string
		request     models.BridgeTypeRequest
		want        error
	}{
		{
			"no adapter name",
			models.BridgeTypeRequest{
				URL: cltest.WebURL(t, "https://denergy.eth"),
			},
			models.NewJSONAPIErrorsWith("No name specified"),
		},
		{
			"invalid adapter name",
			models.BridgeTypeRequest{
				Name: "invalid/adapter",
				URL:  cltest.WebURL(t, "https://denergy.eth"),
			},
			models.NewJSONAPIErrorsWith("task type validation: name invalid/adapter contains invalid characters"),
		},
		{
			"invalid with blank url",
			models.BridgeTypeRequest{
				Name: "validadaptername",
				URL:  cltest.WebURL(t, ""),
			},
			models.NewJSONAPIErrorsWith("URL must be present"),
		},
		{
			"valid url",
			models.BridgeTypeRequest{
				Name: "adapterwithvalidurl",
				URL:  cltest.WebURL(t, "//denergy"),
			},
			nil,
		},
		{
			"valid docker url",
			models.BridgeTypeRequest{
				Name: "adapterwithdockerurl",
				URL:  cltest.WebURL(t, "http://chainlink_cmc-adapter_1:8080"),
			},
			nil,
		},
		{
			"valid MinimumContractPayment positive",
			models.BridgeTypeRequest{
				Name:                   "adapterwithdockerurl",
				URL:                    cltest.WebURL(t, "http://chainlink_cmc-adapter_1:8080"),
				MinimumContractPayment: assets.NewLink(1),
			},
			nil,
		},
		{
			"invalid MinimumContractPayment negative",
			models.BridgeTypeRequest{
				Name:                   "adapterwithdockerurl",
				URL:                    cltest.WebURL(t, "http://chainlink_cmc-adapter_1:8080"),
				MinimumContractPayment: assets.NewLink(-1),
			},
			models.NewJSONAPIErrorsWith("MinimumContractPayment must be positive"),
		},
		{
			"existing core adapter (no longer fails since core adapters no longer exist)",
			models.BridgeTypeRequest{
				Name: "ethtx",
				URL:  cltest.WebURL(t, "https://denergy.eth"),
			},
			nil,
		},
		{
			"new external adapter",
			models.BridgeTypeRequest{
				Name: "gdaxprice",
				URL:  cltest.WebURL(t, "https://denergy.eth"),
			},
			nil,
		}}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := services.ValidateBridgeType(&test.request, store)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestValidateBridgeNotExist(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Create a duplicate
	bt := models.BridgeType{}
	bt.Name = models.MustNewTaskType("solargridreporting")
	bt.URL = cltest.WebURL(t, "https://denergy.eth")
	assert.NoError(t, store.CreateBridgeType(&bt))

	newBridge := models.BridgeTypeRequest{
		Name: "solargridreporting",
	}
	expected := models.NewJSONAPIErrorsWith("Bridge Type solargridreporting already exists")
	result := services.ValidateBridgeTypeNotExist(&newBridge, store)
	assert.Equal(t, expected, result)
}

func TestValidateExternalInitiator(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	url := cltest.WebURL(t, "https://a.web.url")

	//  Add duplicate
	exi := models.ExternalInitiator{
		Name: "duplicate",
		URL:  &url,
	}

	assert.NoError(t, store.CreateExternalInitiator(&exi))

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"basic", `{"name":"bitcoin","url":"https://test.url"}`, false},
		{"basic w/ underscore", `{"name":"bit_coin","url":"https://test.url"}`, false},
		{"basic w/ underscore in url", `{"name":"bitcoin","url":"https://chainlink_bit-coin_1.url"}`, false},
		{"missing url", `{"name":"missing_url"}`, false},
		{"duplicate name", `{"name":"duplicate","url":"https://test.url"}`, true},
		{"invalid name characters", `{"name":"<invalid>","url":"https://test.url"}`, true},
		{"missing name", `{"url":"https://test.url"}`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var exr models.ExternalInitiatorRequest

			assert.NoError(t, json.Unmarshal([]byte(test.input), &exr))
			result := services.ValidateExternalInitiator(&exr, store)

			cltest.AssertError(t, test.wantError, result)
		})
	}
}
