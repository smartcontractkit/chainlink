package adapters

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

type HttpGet struct {
	Endpoint string `json:"endpoint"`
}

func (self *HttpGet) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	response, err := http.Get(self.Endpoint)
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
