package client

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

type GetReportsResult struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

type MercuryServer struct {
	APIClient *resty.Client
}

func NewMercuryServer(url string) *MercuryServer {
	rc := resty.New().SetBaseURL(url)
	return &MercuryServer{
		APIClient: rc,
	}
}

func (ms *MercuryServer) GetReports(feedIDStr string, blockNumber string) (*GetReportsResult, *http.Response, error) {
	result := &GetReportsResult{}
	resp, err := ms.APIClient.R().
		SetPathParams(map[string]string{
			"feedIDStr":     feedIDStr,
			"L2Blocknumber": blockNumber,
		}).
		SetResult(&result).
		Get("/client?feedIDStr={feedIDStr}&L2Blocknumber={L2Blocknumber}")
	if err != nil {
		return nil, nil, err
	}
	return result, resp.RawResponse, err
}
