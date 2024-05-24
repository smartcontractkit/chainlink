package config

type Pyroscope interface {
	AuthToken() string
	ServerAddress() string
	Environment() string
}
