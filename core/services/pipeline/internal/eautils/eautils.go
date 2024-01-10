package eautils

import (
	"encoding/json"
	"net/http"
)

type AdapterStatus struct {
	ErrorMessage       *string `json:"errorMessage"`
	Error              *string `json:"error"`
	StatusCode         *int    `json:"statusCode"`
	ProviderStatusCode *int    `json:"providerStatusCode"`
}

func BestEffortExtractEAStatus(responseBytes []byte) (int, bool) {
	var status AdapterStatus
	err := json.Unmarshal(responseBytes, &status)
	if err != nil {
		return 0, false
	}

	if status.StatusCode == nil {
		return 0, false
	}

	if *status.StatusCode != http.StatusOK {
		return *status.StatusCode, true
	}

	if status.ProviderStatusCode != nil && *status.ProviderStatusCode != http.StatusOK {
		return *status.ProviderStatusCode, true
	}

	if status.Error != nil && *status.Error != "" {
		return http.StatusInternalServerError, true
	}

	return *status.StatusCode, true
}
