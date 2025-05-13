package config

type CRE interface {
	StreamsApiKey() string
	StreamsApiSecret() string
}
