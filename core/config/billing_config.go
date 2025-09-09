package config

type Billing interface {
	URL() string
	TLSEnabled() bool
}
