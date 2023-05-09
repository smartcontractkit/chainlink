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
	// read only so no mutex needed
	//pluginLookupFn func() map[string]plugins.EnvConfig
}

func NewLoopRegistry(app chainlink.Application) *LoopRegistryServer {
	return &LoopRegistryServer{
		exposedPromPort: int(app.GetConfig().Port()),
		registry:        app.GetLoopRegistry(),
	}
}

type pluginConfig struct {
	name string
	port int
}

/*
// list returns deterministic list of loop's known the registry

	func (l *LoopRegistryServer) list() []pluginConfig {
		var out []pluginConfig
		for name, cfg := range l.pluginLookupFn() {
			out = append(out, pluginConfig{name: name, port: cfg.PrometheusPort()})
		}
		sort.Slice(out, func(i, j int) bool {
			return out[i].port < out[j].port
		})
		return out
	}

	func (l *LoopRegistryServer) get(name string) (pluginConfig, bool) {
		result := pluginConfig{name: name, port: invalidPort}
		envCfg, exists := l.pluginLookupFn()[name]
		if exists {
			result.port = envCfg.PrometheusPort()
		}
		return result, exists
	}
*/
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

func (l *LoopRegistryServer) pluginMetricHandler(gc *gin.Context) {

	pluginName := gc.Param("name")
	p, ok := l.registry.Get(pluginName)

	if !ok {

		gc.Data(http.StatusNotFound, "text/plain", []byte(fmt.Sprintf("plugin %q does not exist", html.EscapeString(pluginName))))
		return
	}

	pluginURL := fmt.Sprintf("http://localhost:%d/metrics", p.EnvCfg.PrometheusPort())
	res, err := http.Get(pluginURL)
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
