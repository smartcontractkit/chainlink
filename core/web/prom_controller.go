package web

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// PromController has all prometheus metrics rendering endpoint.
type PromController struct {
	App chainlink.Application
}

// RenderMetrics renders all registered prometheus metrics.
func (pc *PromController) RenderMetrics(c *gin.Context) {
	metricsFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	var metrics []presenters.Metric

	for _, mf := range metricsFamilies {
		var labelsArr []string
		labels := make(map[string]struct{})
		for _, m := range mf.GetMetric() {
			for _, l := range m.GetLabel() {
				_, exist := labels[l.GetName()]
				if !exist {
					labels[l.GetName()] = struct{}{}
					labelsArr = append(labelsArr, l.GetName())
				}
			}
		}
		m := presenters.Metric{
			JAID: presenters.JAID{
				ID: mf.GetName(),
			},
			Name:   mf.GetName(),
			Type:   fmt.Sprintf("%v", mf.GetType()),
			Help:   mf.GetHelp(),
			Labels: labelsArr,
		}
		metrics = append(metrics, m)
	}

	jsonAPIResponse(c, metrics, "metrics")
}
