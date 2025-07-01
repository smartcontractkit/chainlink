package vault

import "time"

type Config struct {
	RequestExpiryDuration time.Duration `json:"requestExpiryDuration"`
}

func (c *Config) Validate() error {
	return nil
}
