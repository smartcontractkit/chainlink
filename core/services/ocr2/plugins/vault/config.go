package vault

import (
	"errors"
	"time"
)

type Config struct {
	RequestExpiryDuration time.Duration `json:"requestExpiryDuration"`
}

func (c *Config) Validate() error {
	if c.RequestExpiryDuration <= 0 {
		return errors.New("request expiry duration cannot be 0")
	}
	return nil
}
