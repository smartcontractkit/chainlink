package time

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Duration time.Duration

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
		*d = Duration(time.Duration(v * float64(24*time.Hour)))
		return nil
	}

	p, err := time.ParseDuration(raw)
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
	dur := time.Duration(d)

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

func (d Duration) Value() time.Duration {
	return time.Duration(d)
}
