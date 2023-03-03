package client

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// KafkaRestClient kafka-rest client
type KafkaRestClient struct {
	APIClient *resty.Client
	Config    *KafkaRestConfig
}

// KafkaRestConfig holds config information for KafkaRestClient
type KafkaRestConfig struct {
	URL string
}

// NewKafkaRestClient creates a new KafkaRestClient
func NewKafkaRestClient(cfg *KafkaRestConfig) *KafkaRestClient {
	return &KafkaRestClient{
		Config:    cfg,
		APIClient: resty.New().SetBaseURL(cfg.URL),
	}
}

// GetTopics Get a list of Kafka topics.
func (krc *KafkaRestClient) GetTopics() ([]string, error) {
	responseBody := []string{}
	resp, err := krc.APIClient.R().SetResult(responseBody).Get("/topics")
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("Unexpected Status Code. Expected %d; Got %d", http.StatusOK, resp.StatusCode())
	}
	return responseBody, err
}
