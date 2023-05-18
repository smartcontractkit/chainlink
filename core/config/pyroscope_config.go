package config

type Pyroscope interface {
	PyroscopeAuthToken() string
	PyroscopeServerAddress() string
	PyroscopeEnvironment() string
}
