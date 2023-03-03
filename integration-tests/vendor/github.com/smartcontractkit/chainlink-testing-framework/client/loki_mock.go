package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/loki/pkg/promtail/api"
	lokiClient "github.com/grafana/loki/pkg/promtail/client"
	"github.com/prometheus/common/model"
)

type PromtailSendResult struct {
	Labels model.LabelSet
	Time   time.Time
	Entry  string
}

type MockPromtailClient struct {
	Results       []PromtailSendResult
	OnHandleEntry api.EntryHandlerFunc
}

// ExtendedLokiClient an extended Loki/Promtail client used for testing last results in batch
type ExtendedLokiClient interface {
	lokiClient.Client
	api.EntryHandler
	HandleStruct(ls model.LabelSet, t time.Time, st interface{}) error
	LastHandleResult() PromtailSendResult
	AllHandleResults() []PromtailSendResult
}

func (m *LokiClient) LastHandleResult() PromtailSendResult {
	panic("implement me")
}

func (m *LokiClient) AllHandleResults() []PromtailSendResult {
	panic("implement me")
}

func NewMockPromtailClient() ExtendedLokiClient {
	mc := &MockPromtailClient{
		Results: make([]PromtailSendResult, 0),
	}
	mc.OnHandleEntry = func(labels model.LabelSet, time time.Time, entry string) error {
		mc.Results = append(mc.Results, PromtailSendResult{Labels: labels, Time: time, Entry: entry})
		return nil
	}
	return mc
}

func (c *MockPromtailClient) HandleStruct(ls model.LabelSet, t time.Time, st interface{}) error {
	d, err := json.Marshal(st)
	if err != nil {
		return fmt.Errorf("failed to marshal struct in response: %v", st)
	}
	return c.OnHandleEntry.Handle(ls, t, string(d))
}

// Stop implements client.Client
func (c *MockPromtailClient) Stop() {}

// Handle implements client.Client
func (c *MockPromtailClient) Handle(labels model.LabelSet, time time.Time, entry string) error {
	return c.OnHandleEntry.Handle(labels, time, entry)
}

func (c *MockPromtailClient) LastHandleResult() PromtailSendResult {
	time.Sleep(2 * time.Second)
	return c.Results[len(c.Results)-1]
}

func (c *MockPromtailClient) AllHandleResults() []PromtailSendResult {
	return c.Results
}
