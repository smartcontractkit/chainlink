package web

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const invalidPort = -1

type LoopRegistryServer struct {
	exposedPromPort int
	registry        *plugins.LoopRegistry
}

func NewLoopRegistryServer(app chainlink.Application) *LoopRegistryServer {
	return &LoopRegistryServer{
		exposedPromPort: int(app.GetConfig().Port()),
		registry:        app.GetLoopRegistry(),
	}
}

// discoveryHandler implements service discovery of prom endpoints for LOOPs in the registry
func (l *LoopRegistryServer) discoveryHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var groups []*targetgroup.Group

	for _, registeredPlugin := range l.registry.List() {
		// create a metric target for each running plugin
		target := &targetgroup.Group{
			Targets: []model.LabelSet{
				{model.AddressLabel: model.LabelValue(fmt.Sprintf("localhost:%d", l.exposedPromPort))},
			},
			Labels: map[model.LabelName]model.LabelValue{
				model.MetricsPathLabel: model.LabelValue(pluginMetricPath(registeredPlugin.Name)),
			},
		}

		groups = append(groups, target)
	}

	b, err := json.Marshal(groups)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
}

// pluginMetricHandlers routes from endpoints published in service discovery to the the backing LOOP endpoint
func (l *LoopRegistryServer) pluginMetricHandler(gc *gin.Context) {

	pluginName := gc.Param("name")
	p, ok := l.registry.Get(pluginName)
	if !ok {
		gc.Data(http.StatusNotFound, "text/plain", []byte(fmt.Sprintf("plugin %q does not exist", html.EscapeString(pluginName))))
		return
	}

	pluginURL := fmt.Sprintf("http://localhost:%d/metrics", p.EnvCfg.PrometheusPort())
	res, err := http.Get(pluginURL) //nolint
	if err != nil {
		gc.Data(http.StatusInternalServerError, "text/plain", []byte(err.Error()))
		return
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("error reading plugin %q metrics: %w", html.EscapeString(pluginName), err)
		gc.Data(http.StatusInternalServerError, "text/plain", []byte(err.Error()))
		return
	}
	gc.Data(http.StatusOK, "text/plain", b)

}

func pluginMetricPath(name string) string {
	return fmt.Sprintf("/plugins/%s/metrics", name)
}
