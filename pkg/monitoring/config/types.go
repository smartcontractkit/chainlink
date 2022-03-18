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
	HTTP           HTTP
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
	// Ids of feeds that are present in the RDD but should not be monitored.
	// These get matched against the string returned by FeedConfig#GetID() for
	// each feed in RDD. If equal, the feed will get ignored!
	IgnoreIDs []string
}

type HTTP struct {
	Address string
}

// Feature is used to add temporary feature flags to the binary.
type Feature struct {
}
