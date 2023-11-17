package pg

import (
	"database/sql/driver"
	"strconv"
	"time"
)

const (
	// DefaultQueryTimeout is a reasonable upper bound for how long a SQL query should take.
	// The configured value should be used instead of this if possible.
	DefaultQueryTimeout = 10 * time.Second
	// longQueryTimeout is a bigger upper bound for how long a SQL query should take
	longQueryTimeout = 1 * time.Minute
)

var _ driver.Valuer = Limit(-1)

// Limit is a helper driver.Valuer for LIMIT queries which uses nil/NULL for negative values.
type Limit int

func (l Limit) String() string {
	if l < 0 {
		return "NULL"
	}
	return strconv.Itoa(int(l))
}

func (l Limit) Value() (driver.Value, error) {
	if l < 0 {
		return nil, nil
	}
	return l, nil
}

var _ QConfig = &qConfig{}

// qConfig implements pg.QCOnfig
type qConfig struct {
	logSQL              bool
	defaultQueryTimeout time.Duration
}

func NewQConfig(logSQL bool) QConfig {
	return &qConfig{logSQL, DefaultQueryTimeout}
}

func (p *qConfig) LogSQL() bool { return p.logSQL }

func (p *qConfig) DefaultQueryTimeout() time.Duration { return p.defaultQueryTimeout }
