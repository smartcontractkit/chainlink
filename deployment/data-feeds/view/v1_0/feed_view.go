package v1_0

type Workflows struct {
	ID map[string]Workflow `json:"workflows"`
}

type Workflow struct {
	Name      string `json:"name"`
	Forwarder string `json:"forwarder"`
	WfOwner   string `json:"wfOwner"`
}

type Feeds struct {
	StreamsID   string   `json:"streamsId"`
	Proxy       string   `json:"proxy"`
	FeedID      string   `json:"feedId"`
	Description string   `json:"description"`
	Deviation   string   `json:"deviation"`
	Heartbeat   int      `json:"heartbeat"`
	Workflows   []string `json:"workflows"`
}

type FeedView struct {
	Workflows map[string]Workflow `json:"workflows"`
	Feeds     []Feeds             `json:"feeds"`
}
