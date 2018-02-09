package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/araddon/dateparse"
	"github.com/mrwonko/cron"
	"github.com/smartcontractkit/chainlink/utils"
)

// WebURL contains the URL of the endpoint.
type WebURL struct {
	*url.URL
}

// UnmarshalJSON parses the raw URL stored in JSON-encoded
// data to a URL structure and sets it to the URL field.
func (w *WebURL) UnmarshalJSON(j []byte) error {
	var v string
	err := json.Unmarshal(j, &v)
	if err != nil {
		return err
	}
	u, err := url.ParseRequestURI(v)
	if err != nil {
		return err
	}
	w.URL = u
	return nil
}

// MarshalJSON returns the JSON-encoded string of the given data.
func (w *WebURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

// Time holds a common field for time.
type Time struct {
	time.Time
}

// UnmarshalJSON parses the raw time stored in JSON-encoded
// data and stores it to the Time field.
func (t *Time) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	newTime, err := dateparse.ParseAny(s)
	t.Time = newTime
	return err
}

// ISO8601 formats and returns the time in ISO 8601 standard.
func (t *Time) ISO8601() string {
	return t.UTC().Format("2006-01-02T15:04:05Z07:00")
}

// DurationFromNow returns the amount of time since the Time
// field was last updated.
func (t *Time) DurationFromNow() time.Duration {
	return t.Time.Sub(time.Now())
}

// HumanString formats and returns the time in RFC 3339 standard.
func (t *Time) HumanString() string {
	return utils.ISO8601UTC(t.Time)
}

// Cron holds the string that will represent the spec of the cron-job.
// It uses 6 fields to represent the seconds (1), minutes (2), hours (3),
// day of the month (4), month (5), and day of the week (6).
type Cron string

// UnmarshalJSON parses the raw spec stored in JSON-encoded
// data and stores it to the Cron string.
func (c *Cron) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("Cron: %v", err)
	}
	if s == "" {
		return nil
	}

	_, err = cron.Parse(s)
	if err != nil {
		return fmt.Errorf("Cron: %v", err)
	}
	*c = Cron(s)
	return nil
}

// String returns the current Cron spec string.
func (c Cron) String() string {
	return string(c)
}
