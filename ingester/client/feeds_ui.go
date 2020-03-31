package client

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"net/http"
)

type UIFeed struct {
	ContractAddress common.Address `json:"contractAddress"`
	Name            string         `json:"name"`
	Pair            []string       `json:"pair"`
	Counter         int            `json:"counter"`
	ContractVersion int            `json:"contractVersion"`
	NetworkID       int            `json:"networkId"`
	History         bool           `json:"history"`
	Bollinger       bool           `json:"bollinger"`
	DecimalPlaces   int            `json:"decimalPlaces"`
	Multiply        string         `json:"multiply"`
}

type UINode struct {
	Address   common.Address `json:"address"`
	Name      string         `json:"name"`
	NetworkId int            `json:"networkId"`
}

type FeedsUI interface {
	Feeds() ([]*UIFeed, error)
	Nodes() ([]*UINode, error)
}

type feedsUI struct {
	client  *http.Client
	baseURL string
}

func NewFeedsUI(baseURL string) FeedsUI {
	return &feedsUI{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

func (f *feedsUI) Nodes() ([]*UINode, error) {
	var uiNodes []*UINode
	return uiNodes, f.do("/nodes.json", &uiNodes)
}

func (f *feedsUI) Feeds() ([]*UIFeed, error) {
	var uiFeeds []*UIFeed
	return uiFeeds, f.do("/feeds.json", &uiFeeds)
}

func (f *feedsUI) do(endpoint string, obj interface{}) error {
	if resp, err := f.client.Get(fmt.Sprintf("%s%s", f.baseURL, endpoint)); err != nil {
		return err
	} else if b, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	} else if err := json.Unmarshal(b, obj); err != nil {
		return err
	}
	return nil
}
