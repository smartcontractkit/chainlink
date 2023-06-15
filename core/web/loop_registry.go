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

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type LoopRegistryServer struct {
	exposedPromPort int
	promTargetHost  string
	registry        *plugins.LoopRegistry
	logger          logger.SugaredLogger
	client          *http.Client

	jsonMarshalFn func(any) ([]byte, error)
}

func NewLoopRegistryServer(app chainlink.Application) *LoopRegistryServer {
	lggr := app.GetLogger()
	promTargetHost, exists := os.LookupEnv("CL_PROMETHEUS_TARGET_HOSTNAME")
	if !exists {
		var err error
		promTargetHost, err = os.Hostname()
		if err != nil {
			lggr.Warnf("could not resolve hostname: %w, falling back to `localhost`", err)
			promTargetHost = "localhost"
		}
	}
	return &LoopRegistryServer{
		exposedPromPort: int(app.GetConfig().WebServer().HTTPPort()),
		registry:        app.GetLoopRegistry(),
		logger:          lggr,
		jsonMarshalFn:   json.Marshal,
		promTargetHost:  promTargetHost,
		client:          &http.Client{Timeout: 1 * time.Second}, // some value much less than the prometheus poll interval will do there
	}
}

// discoveryHandler implements service discovery of prom endpoints for LOOPs in the registry
func (l *LoopRegistryServer) discoveryHandler(w http.ResponseWriter, req *http.Request) {
	l.logger.Debug("handling plugin discovery")
	w.Header().Set("Content-Type", "application/json")
	var groups []*targetgroup.Group

	for _, registeredPlugin := range l.registry.List() {
		// create a metric target for each running plugin
		target := &targetgroup.Group{
			Targets: []model.LabelSet{
				{model.AddressLabel: model.LabelValue(fmt.Sprintf("%s:%d", l.promTargetHost, l.exposedPromPort))},
			},
			Labels: map[model.LabelName]model.LabelValue{
				model.MetricsPathLabel: model.LabelValue(pluginMetricPath(registeredPlugin.Name)),
			},
		}

		groups = append(groups, target)
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

// pluginMetricHandlers routes from endpoints published in service discovery to the the backing LOOP endpoint
func (l *LoopRegistryServer) pluginMetricHandler(gc *gin.Context) {
	l.logger.Debugf("routing plugin metric handler: %w", gc.Params)
	pluginName := gc.Param("name")
	p, ok := l.registry.Get(pluginName)
	if !ok {
		gc.Data(http.StatusNotFound, "text/plain", []byte(fmt.Sprintf("plugin %q does not exist", html.EscapeString(pluginName))))
		return
	}

	pluginURL := fmt.Sprintf("http://%s:%d/metrics", "localhost", p.EnvCfg.PrometheusPort())
	l.logger.Debugf("routing to %s", pluginURL)
	res, err := l.client.Get(pluginURL) //nolint
	if err != nil {
		l.logger.Errorw(fmt.Sprintf("plugin metric handler failed to get plugin url %s:", pluginURL), "err", err)
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
	l.logger.Debugf("successfully fetched plugin metrics: %s", string(b))
	gc.Data(http.StatusOK, "text/plain", b)

}

func pluginMetricPath(name string) string {
	return fmt.Sprintf("/plugins/%s/metrics", name)
}
