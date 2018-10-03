package old

import "net/url"

type WebURL struct {
	*url.URL
}

type TaskType string

type BridgeType struct {
	Name          TaskType `json:"name" storm:"id,unique"`
	URL           WebURL   `json:"url"`
	Confirmations uint64   `json:"confirmations"`
	IncomingToken string   `json:"incomingToken"`
	OutgoingToken string   `json:"outgoingToken"`
}
