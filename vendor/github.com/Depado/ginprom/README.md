<h1 align="center">Ginprom</h1>
<h2 align="center">

Gin Prometheus metrics exporter

[![Sourcegraph](https://sourcegraph.com/github.com/Depado/ginprom/-/badge.svg)](https://sourcegraph.com/github.com/Depado/ginprom?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/Depado/ginprom)](https://goreportcard.com/report/github.com/Depado/ginprom)
[![Build Status](https://drone.depa.do/api/badges/Depado/ginprom/status.svg)](https://drone.depa.do/Depado/ginprom)
[![codecov](https://codecov.io/gh/Depado/ginprom/branch/master/graph/badge.svg)](https://codecov.io/gh/Depado/ginprom)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Depado/bfchroma/blob/master/LICENSE)
[![godoc](https://godoc.org/github.com/Depado/ginprom?status.svg)](https://godoc.org/github.com/Depado/ginprom)

</h2>
<h4 align="center">

Inspired by [github.com/zsais/go-gin-prometheus](https://github.com/zsais/go-gin-prometheus)

</h4>

- [Install](#install)
- [Differences with go-gin-prometheus](#differences-with-go-gin-prometheus)
- [Usage](#usage)
- [Options](#options)
	- [Custom Counters](#custom-counters)
	- [Custom gauges](#custom-gauges)
	- [Path](#path)
	- [Namespace](#namespace)
	- [Subsystem](#subsystem)
	- [Engine](#engine)
	- [Prometheus Registry](#prometheus-registry)
	- [HandlerNameFunc](#handlernamefunc)
	- [RequestPathFunc](#requestpathfunc)
	- [CustomCounterLabels](#customcounterlabels)
	- [Ignore](#ignore)
	- [Token](#token)
	- [Bucket size](#bucket-size)
- [Troubleshooting](#troubleshooting)
	- [The instrumentation doesn't seem to work](#the-instrumentation-doesnt-seem-to-work)

## Install

Simply run:
`go get -u github.com/Depado/ginprom`

## Differences with go-gin-prometheus

- No support for Prometheus' Push Gateway
- Options on constructor
- Adds a `path` label to get the matched route
- Ability to ignore routes

## Usage

```go
package main

import (
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	r.Use(p.Instrument())

	r.GET("/hello/:id", func(c *gin.Context) {})
	r.GET("/world/:id", func(c *gin.Context) {})
	r.Run("127.0.0.1:8080")
}
```

## Options

### Custom Counters

Add custom counters to add own values to the metrics

```go
r := gin.New()
p := ginprom.New(
	ginprom.Engine(r),
)
p.AddCustomCounter("custom", "Some help text to provide", []string{"label"})
r.Use(p.Instrument())
```

Save `p` and use the following functions:

- IncrementCounterValue
- AddCounterValue

### Custom gauges

Add custom gauges to add own values to the metrics

```go
r := gin.New()
p := ginprom.New(
	ginprom.Engine(r),
)
p.AddCustomGauge("custom", "Some help text to provide", []string{"label"})
r.Use(p.Instrument())
```

Save `p` and use the following functions:

- IncrementGaugeValue
- DecrementGaugeValue
- SetGaugeValue

### Path

Override the default path (`/metrics`) on which the metrics can be accessed:

```go
r := gin.New()
p := ginprom.New(
	ginprom.Engine(r),
	ginprom.Path("/custom/metrics"),
)
r.Use(p.Instrument())
```

### Namespace

Override the default namespace (`gin`):

```go
r := gin.New()
p := ginprom.New(
	ginprom.Engine(r),
	ginprom.Namespace("custom_ns"),
)
r.Use(p.Instrument())
```

### Subsystem

Override the default (`gonic`) subsystem:

```go
r := gin.New()
p := ginprom.New(
	ginprom.Engine(r),
	ginprom.Subsystem("your_subsystem"),
)
r.Use(p.Instrument())
```

### Engine

The preferred way to pass the router to ginprom:

```go
r := gin.New()
p := ginprom.New(
	ginprom.Engine(r),
)
r.Use(p.Instrument())
```

The alternative being to call the `Use` method after initialization:

```go
p := ginprom.New()
// ...
r := gin.New()
p.Use(r)
r.Use(p.Instrument())

```

### Prometheus Registry

Use a custom `prometheus.Registry` instead of prometheus client's global registry. This option allows
to use ginprom in multiple gin engines in the same process, or if you would like to integrate ginprom with your own
prometheus `Registry`.

```go
registry := prometheus.NewRegistry() // creates new prometheus metric registry
r := gin.New()
p := ginprom.New(
    ginprom.Registry(registry),
)
r.Use(p.Instrument())
```

### HandlerNameFunc

Change the way the `handler` label is computed. By default, the `(*gin.Context).HandlerName`
function is used.
This option is useful when wanting to group different functions under
the same `handler` label or when using `gin` with decorated handlers.

```go
r := gin.Default()
p := ginprom.New(
	HandlerNameFunc(func (c *gin.Context) string {
		return "my handler"
	}),
)
r.Use(p.Instrument())
```

### RequestPathFunc

Change how the `path` label is computed. By default, the `(*gin.Context).FullPath` function
is used.
This option is useful when wanting to group different requests under the same `path`
label or when wanting to process unknown routes (the default `(*gin.Context).FullPath` returns
an empty string for unregistered routes). Note that requests for which `f` returns the empty
string are ignored.

To specifically ignore certain paths, see the [Ignore](#ignore) option.

```go
r := gin.Default()
p := ginprom.New(
	// record a metric for unregistered routes under the path label "<unknown>"
	RequestPathFunc(func (c *gin.Context) string {
		if fullpath := c.FullPath(); fullpath != "" {
			return fullpath
		}
		return "<unknown>"
	}),
)
r.Use(p.Instrument())
```

### CustomCounterLabels

Add custom labels to the counter metric.

```go
r := gin.Default()
p := ginprom.New(
  ginprom.CustomCounterLabels([]string{"client_id"}, func(c *gin.Context) map[string]string {
    client_id := c.GetHeader("x-client-id")
    if client_id == "" {
      client_id = "unknown"
    }
    return map[string]string{"client_id": client_id}
  }),
)
r.Use(p.Instrument())
```

### Ignore

Ignore allows to completely ignore some routes. Even though you can apply the
middleware to the only groups you're interested in, it is sometimes useful to
have routes not instrumented.

```go
r := gin.New()
p := ginprom.New(
	ginprom.Engine(r),
	ginprom.Ignore("/api/no/no/no", "/api/super/secret/route")
)
r.Use(p.Instrument())
```

Note that most of the time this can be solved by gin groups:

```go
r := gin.New()
p := ginprom.New(ginprom.Engine(r))

// Add the routes that do not need instrumentation
g := r.Group("/api/")
g.Use(p.Instrument())
{
	// Instrumented routes
}
```

### Token

Specify a secret token which Prometheus will use to access the endpoint. If the
token is invalid, the endpoint will return an error.

```go
r := gin.New()
p := ginprom.New(
	ginprom.Engine(r),
	ginprom.Token("supersecrettoken")
)
r.Use(p.Instrument())
```

### Bucket size

Specify the bucket size for the request duration histogram according to your
expected durations.

```go
r := gin.New()
p := ginprom.New(
	ginprom.Engine(r),
	ginprom.BucketSize([]float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}),
)
r.Use(p.Instrument())
```

## Troubleshooting

### The instrumentation doesn't seem to work

Make sure you have set the `gin.Engine` in the `ginprom` middleware, either when
initializing it using `ginprom.New(ginprom.Engine(r))` or using the `Use`
function after the initialization like this :

```go
p := ginprom.New(
	ginprom.Namespace("gin"),
	ginprom.Subsystem("gonic"),
	ginprom.Path("/metrics"),
)
p.Use(r)
r.Use(p.Instrument())
```

By design, if the middleware was to panic, it would do so when a route is
called. That's why it just silently fails when no engine has been set.
