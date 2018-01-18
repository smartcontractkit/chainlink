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

type WebURL struct {
	*url.URL
}

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

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	newTime, err := dateparse.ParseAny(s)
	t.Time = newTime
	return err
}

func (t *Time) ISO8601() string {
	return t.UTC().Format("2006-01-02T15:04:05Z07:00")
}

func (t *Time) DurationFromNow() time.Duration {
	return t.Time.Sub(time.Now())
}

func (t *Time) HumanString() string {
	return utils.ISO8601UTC(t.Time)
}

type Cron string

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

func (c Cron) String() string {
	return string(c)
}
