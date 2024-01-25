package dashboard

import (
	"github.com/K-Phoen/grabana/dashboard"
)

// Dashboard is a dashboard for a cluster with any plugins or product specific panels
type Dashboard struct {
	Name                     string
	LokiDataSourceName       string
	PrometheusDataSourceName string
	opts                     []dashboard.Option
	extendedOpts             []dashboard.Option
	Builder                  dashboard.Builder
}

// NewDashboard returns a new dashboard for a Chainlink cluster, can be used as a base for more complex plugin based dashboards
func NewDashboard(name string, ldsn string, pdsn string, opts []dashboard.Option) (*Dashboard, error) {
	db := &Dashboard{
		Name:                     name,
		LokiDataSourceName:       ldsn,
		PrometheusDataSourceName: pdsn,
		extendedOpts:             opts,
	}
	if err := db.generate(); err != nil {
		return db, err
	}
	return db, nil
}
