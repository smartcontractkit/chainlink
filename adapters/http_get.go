package adapters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

type HttpGet struct {
	Endpoint *url.URL `json:"endpoint"`
}

func (hga *HttpGet) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	response, err := http.Get(hga.Endpoint.String())
	if err != nil {
		return models.RunResultWithError(err)
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		return models.RunResultWithError(err)
	}
	if response.StatusCode >= 300 {
		return models.RunResultWithError(fmt.Errorf(body))
	}

	return models.RunResultWithValue(body)
}

func (hga *HttpGet) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "endpoint" {
			u, err := url.ParseRequestURI(v)
			if err != nil {
				return err
			}
			hga.Endpoint = u
		}
	}

	return nil
}
