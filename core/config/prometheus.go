package config

type Prometheus interface {
	AuthToken() string
}
