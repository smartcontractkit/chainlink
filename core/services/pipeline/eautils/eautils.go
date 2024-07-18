package eautils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AdapterStatus struct {
	ErrorMessage       *string `json:"errorMessage"`
	Error              any     `json:"error"`
	StatusCode         *int    `json:"statusCode"`
	ProviderStatusCode *int    `json:"providerStatusCode"`
}

func BestEffortExtractEAStatus(responseBytes []byte) (code int, ok bool) {
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

	if status.Error != nil {
		return http.StatusInternalServerError, true
	}

	return *status.StatusCode, true
}

type adapterErrorResponse struct {
	Error *AdapterError `json:"error"`
}

type AdapterError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (err *AdapterError) Error() string {
	return fmt.Sprintf("%s: %s", err.Name, err.Message)
}

func BestEffortExtractEAError(responseBytes []byte) error {
	var errorResponse adapterErrorResponse
	err := json.Unmarshal(responseBytes, &errorResponse)
	if err != nil {
		return nil
	}
	if errorResponse.Error != nil {
		return errorResponse.Error
	}
	return nil
}
