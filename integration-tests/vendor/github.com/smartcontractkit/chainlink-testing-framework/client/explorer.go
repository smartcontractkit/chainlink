package client

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// ExplorerClient is used to call Explorer API endpoints
type ExplorerClient struct {
	APIClient *resty.Client
	Config    *ExplorerConfig
}

// NewExplorerClient creates a new explorer mock client
func NewExplorerClient(cfg *ExplorerConfig) *ExplorerClient {
	return &ExplorerClient{
		Config:    cfg,
		APIClient: resty.New().SetBaseURL(cfg.URL),
	}
}

// PostAdminNodes is used to exercise the POST /api/v1/admin/nodes endpoint
// This endpoint is used to create access keys for nodes
func (em *ExplorerClient) PostAdminNodes(nodeName string) (NodeAccessKeys, error) {
	em.APIClient.SetHeaders(map[string]string{
		"x-explore-admin-password": em.Config.AdminPassword,
		"x-explore-admin-username": em.Config.AdminUsername,
		"Content-Type":             "application/json",
		"Accept":                   "application/json",
	})
	requestBody := &Name{Name: nodeName}
	responseBody := NodeAccessKeys{}
	log.Info().Str("Explorer URL", em.Config.URL).Msg("Creating node credentials")
	resp, err := em.APIClient.R().SetBody(requestBody).SetResult(&responseBody).Post("/api/v1/admin/nodes")
	if resp.StatusCode() != http.StatusCreated {
		err = fmt.Errorf("Unexpected Status Code. Expected %d; Got %d", http.StatusCreated, resp.StatusCode())
	}
	return responseBody, err
}

// Name is the body of the request
type Name struct {
	Name string `json:"name"`
}

// NodeAccessKeys is the body of the response
type NodeAccessKeys struct {
	ID        string `json:"id"`
	AccessKey string `json:"accessKey"`
	Secret    string `json:"secret"`
}

// ExplorerConfig holds config information for ExplorerClient
type ExplorerConfig struct {
	URL           string
	AdminUsername string
	AdminPassword string
}
