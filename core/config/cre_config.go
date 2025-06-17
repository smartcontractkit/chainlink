package config

type CRE interface {
	WsURL() string
	RestURL() string
	StreamsAPIKey() string
	StreamsAPISecret() string
}
