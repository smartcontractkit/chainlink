package web

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type LoopRegistryServer struct {
	exposedPromPort   int
	discoveryHostName string // discovery endpoint hostname. must be accessible to external prom for scraping
	loopHostName      string // internal hostname of loopps. used by node to forward external prom requests
	registry          *plugins.LoopRegistry
	logger            logger.SugaredLogger
	client            *http.Client

	jsonMarshalFn func(any) ([]byte, error)
}

func NewLoopRegistryServer(app chainlink.Application) *LoopRegistryServer {
	discoveryHostName, loopHostName := initHostNames()
	return &LoopRegistryServer{
		exposedPromPort:   int(app.GetConfig().WebServer().HTTPPort()),
		registry:          app.GetLoopRegistry(),
		logger:            app.GetLogger(),
		jsonMarshalFn:     json.Marshal,
		discoveryHostName: discoveryHostName,
		loopHostName:      loopHostName,
		client:            &http.Client{Timeout: 1 * time.Second}, // some value much less than the prometheus poll interval will do there
	}
}

// discoveryHandler implements service discovery of prom endpoints for LOOPs in the registry
func (l *LoopRegistryServer) discoveryHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var groups []*targetgroup.Group

	// add node metrics to service discovery
	groups = append(groups, metricTarget(l.discoveryHostName, l.exposedPromPort, "/metrics"))

	// add all the plugins
	for _, registeredPlugin := range l.registry.List() {
		groups = append(groups, metricTarget(l.discoveryHostName, l.exposedPromPort, pluginMetricPath(registeredPlugin.Name)))
	}

	b, err := l.jsonMarshalFn(groups)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			l.logger.Error(err)
		}
		return
	}
	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		l.logger.Error(err)
	}
}

func metricTarget(hostName string, port int, path string) *targetgroup.Group {
	return &targetgroup.Group{
		Targets: []model.LabelSet{
			// target address will be called by external prometheus
			{model.AddressLabel: model.LabelValue(fmt.Sprintf("%s:%d", hostName, port))},
		},
		Labels: map[model.LabelName]model.LabelValue{
			model.MetricsPathLabel: model.LabelValue(path),
		},
	}
}

// pluginMetricHandlers routes from endpoints published in service discovery to the backing LOOP endpoint
func (l *LoopRegistryServer) pluginMetricHandler(gc *gin.Context) {
	pluginName := gc.Param("name")
	p, ok := l.registry.Get(pluginName)
	if !ok {
		gc.Data(http.StatusNotFound, "text/plain", []byte(fmt.Sprintf("plugin %q does not exist", html.EscapeString(pluginName))))
		return
	}

	// unlike discovery, this endpoint is internal btw the node and plugin
	pluginURL := fmt.Sprintf("http://%s:%d/metrics", l.loopHostName, p.EnvCfg.PrometheusPort)
	res, err := l.client.Get(pluginURL) //nolint
	if err != nil {
		msg := fmt.Sprintf("plugin metric handler failed to get plugin url %s", html.EscapeString(pluginURL))
		l.logger.Errorw(msg, "err", err)
		gc.Data(http.StatusInternalServerError, "text/plain", []byte(fmt.Sprintf("%s: %s", msg, err)))
		return
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		msg := fmt.Sprintf("error reading plugin %q metrics", html.EscapeString(pluginName))
		l.logger.Errorw(msg, "err", err)
		gc.Data(http.StatusInternalServerError, "text/plain", []byte(fmt.Sprintf("%s: %s", msg, err)))
		return
	}

	gc.Data(http.StatusOK, "text/plain", b)
}

func initHostNames() (discoveryHost, loopHost string) {
	var exists bool
	discoveryHost, exists = env.PrometheusDiscoveryHostName.Lookup()
	if !exists {
		var err error
		discoveryHost, err = os.Hostname()
		if err != nil {
			discoveryHost = "localhost"
		}
	}

	loopHost, exists = env.LOOPPHostName.Lookup()
	if !exists {
		// this is the expected case; no known uses for the env var other than
		// as an escape hatch.
		loopHost = "localhost"
	}
	return discoveryHost, loopHost
}

func pluginMetricPath(name string) string {
	return fmt.Sprintf("/plugins/%s/metrics", name)
}
