package limiter

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Rate is the rate.
type Rate struct {
	Formatted string
	Period    time.Duration
	Limit     int64
}

// NewRateFromFormatted returns the rate from the formatted version.
func NewRateFromFormatted(formatted string) (Rate, error) {
	rate := Rate{}

	values := strings.Split(formatted, "-")
	if len(values) != 2 {
		return rate, errors.Errorf("incorrect format '%s'", formatted)
	}

	periods := map[string]time.Duration{
		"S": time.Second,    // Second
		"M": time.Minute,    // Minute
		"H": time.Hour,      // Hour
		"D": time.Hour * 24, // Day
	}

	limit, period := values[0], strings.ToUpper(values[1])

	p, ok := periods[period]
	if !ok {
		return rate, errors.Errorf("incorrect period '%s'", period)
	}

	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return rate, errors.Errorf("incorrect limit '%s'", limit)
	}

	rate = Rate{
		Formatted: formatted,
		Period:    p,
		Limit:     l,
	}

	return rate, nil
}
