package config

type CCV interface {
	AggregatorSecrets() []AggregatorSecret
}

type AggregatorSecret interface {
	CommitteeID() string
	APIKey() string
	APISecret() string
}
