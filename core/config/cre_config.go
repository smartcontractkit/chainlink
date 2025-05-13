package config

type CRE interface {
	WsURL() string
	RestURL() string
	StreamsApiKey() string
	StreamsApiSecret() string
}
