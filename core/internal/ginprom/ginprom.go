// Package ginprom provides the subset of upstream ginprom
// The middleware behavior and metric names are preserved from
// github.com/Depado/ginprom v1.8.0 (MIT).
package ginprom

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrInvalidToken = errors.New("invalid or missing token")

	defaultPath = "/metrics"
	defaultNs   = "gin"
	defaultSys  = "gonic"

	defaultReqCntMetricName = "requests_total"
	defaultReqDurMetricName = "request_duration"
	defaultReqSzMetricName  = "request_size_bytes"
	defaultResSzMetricName  = "response_size_bytes"
)

type ignoredPaths struct {
	sync.RWMutex
	values map[string]bool
}

// Prometheus contains the metrics gathered by the instance and its path.
type Prometheus struct {
	reqCnt       *prometheus.CounterVec
	reqDur       *prometheus.HistogramVec
	reqSz, resSz prometheus.Summary

	MetricsPath     string
	Namespace       string
	Subsystem       string
	Token           string
	Ignored         ignoredPaths
	Engine          *gin.Engine
	BucketsSize     []float64
	Registry        *prometheus.Registry
	HandlerNameFunc func(c *gin.Context) string
	RequestPathFunc func(c *gin.Context) string
	HandlerOpts     promhttp.HandlerOpts

	RequestCounterMetricName  string
	RequestDurationMetricName string
	RequestSizeMetricName     string
	ResponseSizeMetricName    string
}

type PrometheusOption func(*Prometheus)

func Path(path string) PrometheusOption {
	return func(p *Prometheus) { p.MetricsPath = path }
}

func Ignore(paths ...string) PrometheusOption {
	return func(p *Prometheus) {
		p.Ignored.Lock()
		defer p.Ignored.Unlock()
		for _, path := range paths {
			p.Ignored.values[path] = true
		}
	}
}

func BucketSize(b []float64) PrometheusOption {
	return func(p *Prometheus) { p.BucketsSize = b }
}

func Subsystem(sub string) PrometheusOption {
	return func(p *Prometheus) { p.Subsystem = sub }
}

func Namespace(ns string) PrometheusOption {
	return func(p *Prometheus) { p.Namespace = ns }
}

func Token(token string) PrometheusOption {
	return func(p *Prometheus) { p.Token = token }
}

func RequestCounterMetricName(name string) PrometheusOption {
	return func(p *Prometheus) { p.RequestCounterMetricName = name }
}

func RequestDurationMetricName(name string) PrometheusOption {
	return func(p *Prometheus) { p.RequestDurationMetricName = name }
}

func RequestSizeMetricName(name string) PrometheusOption {
	return func(p *Prometheus) { p.RequestSizeMetricName = name }
}

func ResponseSizeMetricName(name string) PrometheusOption {
	return func(p *Prometheus) { p.ResponseSizeMetricName = name }
}

func Engine(e *gin.Engine) PrometheusOption {
	return func(p *Prometheus) { p.Engine = e }
}

func Registry(r *prometheus.Registry) PrometheusOption {
	return func(p *Prometheus) { p.Registry = r }
}

func HandlerNameFunc(f func(c *gin.Context) string) PrometheusOption {
	return func(p *Prometheus) { p.HandlerNameFunc = f }
}

func RequestPathFunc(f func(c *gin.Context) string) PrometheusOption {
	return func(p *Prometheus) { p.RequestPathFunc = f }
}

func HandlerOpts(opts promhttp.HandlerOpts) PrometheusOption {
	return func(p *Prometheus) { p.HandlerOpts = opts }
}

func New(options ...PrometheusOption) *Prometheus {
	p := &Prometheus{
		MetricsPath:               defaultPath,
		Namespace:                 defaultNs,
		Subsystem:                 defaultSys,
		HandlerNameFunc:           (*gin.Context).HandlerName,
		RequestPathFunc:           (*gin.Context).FullPath,
		RequestCounterMetricName:  defaultReqCntMetricName,
		RequestDurationMetricName: defaultReqDurMetricName,
		RequestSizeMetricName:     defaultReqSzMetricName,
		ResponseSizeMetricName:    defaultResSzMetricName,
	}
	p.Ignored.values = make(map[string]bool)
	for _, option := range options {
		option(p)
	}

	p.register()
	if p.Engine != nil {
		p.Engine.GET(p.MetricsPath, p.handler(p.Token))
	}

	return p
}

func (p *Prometheus) getRegistererAndGatherer() (prometheus.Registerer, prometheus.Gatherer) {
	if p.Registry == nil {
		return prometheus.DefaultRegisterer, prometheus.DefaultGatherer
	}
	return p.Registry, p.Registry
}

func (p *Prometheus) register() {
	p.reqCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: p.Namespace,
			Subsystem: p.Subsystem,
			Name:      p.RequestCounterMetricName,
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method", "handler", "host", "path"},
	)
	p.mustRegister(p.reqCnt)

	p.reqDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: p.Namespace,
		Subsystem: p.Subsystem,
		Buckets:   p.BucketsSize,
		Name:      p.RequestDurationMetricName,
		Help:      "The HTTP request latency bucket.",
	}, []string{"method", "path", "host"})
	p.mustRegister(p.reqDur)

	p.reqSz = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: p.Namespace,
		Subsystem: p.Subsystem,
		Name:      p.RequestSizeMetricName,
		Help:      "The HTTP request sizes in bytes.",
	})
	p.mustRegister(p.reqSz)

	p.resSz = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: p.Namespace,
		Subsystem: p.Subsystem,
		Name:      p.ResponseSizeMetricName,
		Help:      "The HTTP response sizes in bytes.",
	})
	p.mustRegister(p.resSz)
}

func (p *Prometheus) mustRegister(c ...prometheus.Collector) {
	registerer, _ := p.getRegistererAndGatherer()
	registerer.MustRegister(c...)
}

func (p *Prometheus) isIgnored(path string) bool {
	p.Ignored.RLock()
	defer p.Ignored.RUnlock()
	return p.Ignored.values[path]
}

func (p *Prometheus) Instrument() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := p.RequestPathFunc(c)

		if path == "" || p.isIgnored(path) {
			c.Next()
			return
		}

		reqSz := computeApproximateRequestSize(c.Request)
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSz := float64(c.Writer.Size())

		p.reqCnt.WithLabelValues(status, c.Request.Method, p.HandlerNameFunc(c), c.Request.Host, path).Inc()
		p.reqDur.WithLabelValues(c.Request.Method, path, c.Request.Host).Observe(elapsed)
		p.reqSz.Observe(float64(reqSz))
		p.resSz.Observe(resSz)
	}
}

func (p *Prometheus) Use(e *gin.Engine) {
	e.GET(p.MetricsPath, p.handler(p.Token))
	p.Engine = e
}

func (p *Prometheus) handler(token string) gin.HandlerFunc {
	registerer, gatherer := p.getRegistererAndGatherer()
	h := promhttp.InstrumentMetricHandler(registerer, promhttp.HandlerFor(gatherer, p.HandlerOpts))
	return func(c *gin.Context) {
		if token == "" {
			h.ServeHTTP(c.Writer, c.Request)
			return
		}

		header := c.Request.Header.Get("Authorization")
		if header == "" || header != fmt.Sprintf("Bearer %s", token) {
			c.String(http.StatusUnauthorized, ErrInvalidToken.Error())
			return
		}

		h.ServeHTTP(c.Writer, c.Request)
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}
	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)
	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
