package adapters

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpGet struct {
	Endpoint string `json:"endpoint"`
}

func (self *HttpGet) Perform(input RunResult) RunResult {
	response, err := http.Get(self.Endpoint)
	if err != nil {
		return RunResult{Error: err}
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		return RunResult{Error: err}
	}
	if response.StatusCode >= 300 {
		return RunResult{Error: fmt.Errorf(body)}
	}

	return RunResultWithValue(body)
}
