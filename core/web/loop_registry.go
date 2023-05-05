package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const invalidPort = -1

type LoopRegistry struct {
	exposedPromPort int
	// read only so no mutex needed
	pluginLookupFn func() map[string]plugins.EnvConfigurer
}

func NewLoopRegistry(app chainlink.Application) *LoopRegistry {
	return &LoopRegistry{
		exposedPromPort: int(app.GetConfig().Port()),
		pluginLookupFn:  app.GetLoopEnvConfig,
	}
}

type pluginConfig struct {
	name string
	port int
}

// list returns deterministic list of loop's known the registry
func (l *LoopRegistry) list() []pluginConfig {
	var out []pluginConfig
	for name, cfg := range l.pluginLookupFn() {
		out = append(out, pluginConfig{name: name, port: cfg.PrometheusPort()})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].port < out[j].port
	})
	return out
}

func (l *LoopRegistry) get(name string) (pluginConfig, bool) {
	result := pluginConfig{name: name, port: invalidPort}
	envCfg, exists := l.pluginLookupFn()[name]
	if exists {
		result.port = envCfg.PrometheusPort()
	}
	return result, exists
}

func (l *LoopRegistry) discoveryHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	groups := make([]*targetgroup.Group, 0)

	for _, p := range l.list() {
		// create a metric target for each running plugin
		target := &targetgroup.Group{
			Targets: []model.LabelSet{
				{model.AddressLabel: model.LabelValue(fmt.Sprintf("host.docker.internal:%d", l.exposedPromPort))},
			},
			Labels: map[model.LabelName]model.LabelValue{
				model.MetricsPathLabel: model.LabelValue(pluginMetricPath(p.name)),
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

func (l *LoopRegistry) pluginMetricHandler(w http.ResponseWriter, req *http.Request) {
	// route to metric handler for the target

	pluginName := extractPluginName(req.URL.Path)

	p, ok := l.get(pluginName)

	if !ok {
		w.Write([]byte(fmt.Sprintf("plugin '%s' does not exist", pluginName)))
		return
	}

	pluginURL := fmt.Sprintf("http://localhost:%d/metrics", p.port)
	res, err := http.Get(pluginURL)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
}

func pluginMetricPath(name string) string {
	return fmt.Sprintf("plugins/%s/metrics", name)
}

func extractPluginName(urlPath string) string {
	temp := strings.TrimLeft(urlPath, "/")
	return strings.Split(temp, "/")[1]
}
