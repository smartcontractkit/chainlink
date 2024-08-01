package ginprom

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusOption func(*Prometheus)

// Path is an option allowing to set the metrics path when initializing with New.
func Path(path string) PrometheusOption {
	return func(p *Prometheus) {
		p.MetricsPath = path
	}
}

// Ignore is used to disable instrumentation on some routes.
func Ignore(paths ...string) PrometheusOption {
	return func(p *Prometheus) {
		p.Ignored.Lock()
		defer p.Ignored.Unlock()
		for _, path := range paths {
			p.Ignored.values[path] = true
		}
	}
}

// BucketSize is used to define the default bucket size when initializing with
// New.
func BucketSize(b []float64) PrometheusOption {
	return func(p *Prometheus) {
		p.BucketsSize = b
	}
}

// Subsystem is an option allowing to set the subsystem when initializing
// with New.
func Subsystem(sub string) PrometheusOption {
	return func(p *Prometheus) {
		p.Subsystem = sub
	}
}

// Namespace is an option allowing to set the namespace when initializing
// with New.
func Namespace(ns string) PrometheusOption {
	return func(p *Prometheus) {
		p.Namespace = ns
	}
}

// Token is an option allowing to set the bearer token in prometheus
// with New.
// Example: ginprom.New(ginprom.Token("your_custom_token"))
func Token(token string) PrometheusOption {
	return func(p *Prometheus) {
		p.Token = token
	}
}

// RequestCounterMetricName is an option allowing to set the request counter metric name.
func RequestCounterMetricName(reqCntMetricName string) PrometheusOption {
	return func(p *Prometheus) {
		p.RequestCounterMetricName = reqCntMetricName
	}
}

// RequestDurationMetricName is an option allowing to set the request duration metric name.
func RequestDurationMetricName(reqDurMetricName string) PrometheusOption {
	return func(p *Prometheus) {
		p.RequestDurationMetricName = reqDurMetricName
	}
}

// RequestSizeMetricName is an option allowing to set the request size metric name.
func RequestSizeMetricName(reqSzMetricName string) PrometheusOption {
	return func(p *Prometheus) {
		p.RequestSizeMetricName = reqSzMetricName
	}
}

// ResponseSizeMetricName is an option allowing to set the response size metric name.
func ResponseSizeMetricName(resDurMetricName string) PrometheusOption {
	return func(p *Prometheus) {
		p.ResponseSizeMetricName = resDurMetricName
	}
}

// Engine is an option allowing to set the gin engine when intializing with New.
// Example:
// r := gin.Default()
// p := ginprom.New(Engine(r))
func Engine(e *gin.Engine) PrometheusOption {
	return func(p *Prometheus) {
		p.Engine = e
	}
}

// Registry is an option allowing to set a  *prometheus.Registry with New.
// Use this option if you want to use a custom Registry instead of a global one that prometheus
// client uses by default
// Example:
// r := gin.Default()
// p := ginprom.New(Registry(r))
func Registry(r *prometheus.Registry) PrometheusOption {
	return func(p *Prometheus) {
		p.Registry = r
	}
}

// HandlerNameFunc is an option allowing to set the HandlerNameFunc with New.
// Use this option if you want to override the default behavior (i.e. using
// (*gin.Context).HandlerName). This is useful when wanting to group different
// functions under the same "handler" label or when using gin with decorated handlers
// Example:
// r := gin.Default()
// p := ginprom.New(HandlerNameFunc(func (c *gin.Context) string { return "my handler" }))
func HandlerNameFunc(f func(c *gin.Context) string) PrometheusOption {
	return func(p *Prometheus) {
		p.HandlerNameFunc = f
	}
}

// HandlerOpts is an option allowing to set the promhttp.HandlerOpts.
// Use this option if you want to override the default zero value.
func HandlerOpts(opts promhttp.HandlerOpts) PrometheusOption {
	return func(p *Prometheus) {
		p.HandlerOpts = opts
	}
}

// RequestPathFunc is an option allowing to set the RequestPathFunc with New.
// Use this option if you want to override the default behavior (i.e. using
// (*gin.Context).FullPath). This is useful when wanting to group different requests
// under the same "path" label or when wanting to process unknown routes (the default
// (*gin.Context).FullPath return an empty string for unregistered routes). Note that
// requests for which f returns the empty string are ignored.
// To specifically ignore certain paths, see the Ignore option.
// Example:
//
//	r := gin.Default()
//	p := ginprom.New(RequestPathFunc(func (c *gin.Context) string {
//		if fullpath := c.FullPath(); fullpath != "" {
//			return fullpath
//		}
//		return "<unknown>"
//	}))
func RequestPathFunc(f func(c *gin.Context) string) PrometheusOption {
	return func(p *Prometheus) {
		p.RequestPathFunc = f
	}
}

func CustomCounterLabels(labels []string, f func(c *gin.Context) map[string]string) PrometheusOption {
	return func(p *Prometheus) {
		p.customCounterLabelsProvider = f
		p.customCounterLabels = labels
	}
}
