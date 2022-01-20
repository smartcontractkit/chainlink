// package config parses flags, environment variables and json object to build
// a Config object that's used througout the monitor.
package config

import (
	"time"
)

type Config struct {
	Kafka          Kafka
	SchemaRegistry SchemaRegistry
	Feeds          Feeds
	Http           Http
	Feature        Feature
}

type Kafka struct {
	Brokers          string
	ClientID         string
	SecurityProtocol string

	SaslMechanism string
	SaslUsername  string
	SaslPassword  string

	TransmissionTopic        string
	ConfigSetSimplifiedTopic string
}

type SchemaRegistry struct {
	URL      string
	Username string
	Password string
}

type Feeds struct {
	URL             string
	RDDReadTimeout  time.Duration
	RDDPollInterval time.Duration
}

type Http struct {
	Address string
}

type Feature struct {
	// If set, the monitor will not read from a chain instead from a source of random state snapshots.
	TestOnlyFakeReaders bool
	// If set, the monitor will not read from the RDD, instead it will get data from a local source of random feeds configurations.
	TestOnlyFakeRdd bool
}
