package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/smartcontractkit/chainlink/store/models"
	"net/http"
)

// PromController inherits the Controller type to allow Prometheus exporting
type PromController struct {
	Controller
}

var totalSpecs = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_total_specs",
		Help: "Total number of specs on the node.",
	},
	[]string{"address"},
)
var totalSpecRuns = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_total_spec_runs",
		Help: "Total number of specs on the node.",
	},
	[]string{"address", "spec_id"},
)
var ethBalance = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_eth_balance",
		Help: "The Ethereum balance of the nodes wallet.",
	},
	[]string{"address"},
)
var linkBalance = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_link_balance",
		Help: "The LINK balance of the nodes wallet.",
	},
	[]string{"address"},
)
var specTaskCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_spec_task_count",
		Help: "The total of spec task runs of a specific adaptor.",
	},
	[]string{"address", "spec_id", "adaptor"},
)
var specRunStatusCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_spec_run_status_count",
		Help: "The count of spec run statuses.",
	},
	[]string{"address", "spec_id", "status"},
)
var specRunParamCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_spec_param_count",
		Help: "The count of spec run urls used.",
	},
	[]string{"address", "spec_id", "param", "value"},
)

func init() {
	prometheus.MustRegister(
		totalSpecs,
		totalSpecRuns,
		ethBalance,
		linkBalance,
		specTaskCount,
		specRunStatusCount,
		specRunParamCount,
	)
}

// Show outputs the metrics in the Prometheus format
// User-Agent of request needs to contain "Prometheus"
// Example:
//  "<application>/metrics"
//  "User-Agent: Prometheus/2.3.2"
func (pc *PromController) Show(jsm *models.JobSpecMetrics, w http.ResponseWriter, r *http.Request) {
	adr := jsm.Address

	totalSpecs.WithLabelValues(adr).Set(float64(len(jsm.JobSpecCounts)))

	for _, js := range jsm.JobSpecCounts {
		totalSpecRuns.WithLabelValues(adr, js.ID).Set(float64(js.RunCount))

		for a, c := range js.AdaptorCount {
			specTaskCount.WithLabelValues(adr, js.ID, a.String()).Set(float64(c))
		}

		for s, c := range js.StatusCount {
			specRunStatusCount.WithLabelValues(adr, js.ID, string(s)).Set(float64(c))
		}

		for p, pc := range js.ParamCount {
			for _, vc := range pc {
				specRunParamCount.WithLabelValues(adr, js.ID, p, vc.Value).Set(float64(vc.Count))
			}
		}
	}

	prom := promhttp.Handler()
	prom.ServeHTTP(w, r)
}

// UserAgent returns the `User-Agent` header that is set by Prometheus
func (pc *PromController) UserAgent() string {
	return "Prometheus"
}
