package config

type CCV interface {
	AggregatorSecrets() []AggregatorSecret
	IndexerSecret() IndexerSecret
}

type AggregatorSecret interface {
	VerifierID() string
	APIKey() string
	APISecret() string
}

type IndexerSecret interface {
	APIKey() string
	APISecret() string
}
