package time

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type Duration commonconfig.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var (
		raw string
	)
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw == "" {
		return errors.New("no duration provided")
	}

	// Support "Xd" where X is a number, as X * 24h
	if strings.HasSuffix(raw, "d") {
		num := strings.TrimSuffix(raw, "d")
		v, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return fmt.Errorf("invalid day duration %q: %w", raw, err)
		}
		duration, err := commonconfig.NewDuration(time.Duration(v * float64(24*time.Hour)))
		if err != nil {
			return fmt.Errorf("failed to create duration from days %q: %w", raw, err)
		}
		*d = Duration(duration)
		return nil
	}

	p, err := commonconfig.ParseDuration(raw)
	if err != nil {
		return err
	}

	*d = Duration(p)
	return nil
}

// MarshalJSON prefers "Xd" if the duration is an exact multiple of 24h,
// otherwise it falls back to Go's native duration string, e.g. "36h0m0s".
func (d Duration) MarshalJSON() ([]byte, error) {
	const day = 24 * time.Hour
	dur := commonconfig.Duration(d).Duration()

	var s string
	if dur%day == 0 {
		days := dur / day
		s = fmt.Sprintf("%dd", days)
	} else {
		s = dur.String()
	}

	// json.Marshal to get a proper JSON string with quotes/escaping
	return json.Marshal(s)
}

func (d Duration) Value() (driver.Value, error) {
	return commonconfig.Duration(d).Value()
}

func (d Duration) Duration() time.Duration {
	return commonconfig.Duration(d).Duration()
}
