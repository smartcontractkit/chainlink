package config

type Prometheus interface {
	PrometheusAuthToken() string
}
