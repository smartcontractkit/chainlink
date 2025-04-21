package v1_0

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/shared"
)

type Workflows struct {
	ID map[string]Workflow `json:"workflows"`
}

type Workflow struct {
	Name      string `json:"name"`
	Forwarder string `json:"forwarder"`
	Owner     string `json:"owner"`
}

type Feed struct {
	StreamsID    string   `json:"streamsID"`
	ProxyAddress string   `json:"proxy"`
	FeedID       string   `json:"feedID"`
	Description  string   `json:"description"`
	Deviation    string   `json:"deviation"`
	Heartbeat    int      `json:"heartbeat"`
	Workflows    []string `json:"workflows"`
}

type FeedState struct {
	Workflows map[string]Workflow `json:"workflows"`
	Feeds     []Feed              `json:"feeds"`
}

func (fv *FeedState) Validate() error {
	w := fv.Workflows
	if len(fv.Feeds) == 0 {
		return errors.New("at least one feed is required for workflow")
	}
	for _, f := range fv.Feeds {
		if f.StreamsID == "" {
			return fmt.Errorf("streamsID is required for feed %s", f.FeedID)
		}
		if f.ProxyAddress == "" {
			return fmt.Errorf("proxy address is required for feed %s", f.FeedID)
		}
		if f.FeedID == "" {
			return fmt.Errorf("feedID is required for feed %s", f.FeedID)
		}
		if f.Description == "" {
			return fmt.Errorf("description is required for feed %s", f.FeedID)
		}
		if f.Deviation == "" {
			return fmt.Errorf("deviation is required for feed %s", f.FeedID)
		}
		if f.Heartbeat < 0 {
			return fmt.Errorf("heartbeat must be positive for feed %s", f.FeedID)
		}
		if len(f.Workflows) == 0 {
			return fmt.Errorf("at least one workflow is required for feed %s", f.FeedID)
		}
		for _, workflow := range f.Workflows {
			if _, ok := w[workflow]; !ok {
				return fmt.Errorf("workflow %s not found for feed %s", workflow, f.FeedID)
			}
		}
		err := shared.ValidateFeedID(f.FeedID)
		if err != nil {
			return fmt.Errorf("invalid feedID %s for feed %s: %w", f.FeedID, f.FeedID, err)
		}
	}

	return nil
}
