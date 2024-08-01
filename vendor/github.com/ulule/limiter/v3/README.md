# Limiter

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][circle-img]][circle-url]
[![Go Report Card][goreport-img]][goreport-url]

_Dead simple rate limit middleware for Go._

- Simple API
- "Store" approach for backend
- Redis support (but not tied too)
- Middlewares: HTTP, [FastHTTP][6] and [Gin][4]

## Installation

Using [Go Modules](https://github.com/golang/go/wiki/Modules)

```bash
$ go get github.com/ulule/limiter/v3@v3.11.2
```

## Usage

In five steps:

- Create a `limiter.Rate` instance _(the number of requests per period)_
- Create a `limiter.Store` instance _(see [Redis](https://github.com/ulule/limiter/blob/master/drivers/store/redis/store.go) or [In-Memory](https://github.com/ulule/limiter/blob/master/drivers/store/memory/store.go))_
- Create a `limiter.Limiter` instance that takes store and rate instances as arguments
- Create a middleware instance using the middleware of your choice
- Give the limiter instance to your middleware initializer

**Example:**

```go
// Create a rate with the given limit (number of requests) for the given
// period (a time.Duration of your choice).
import "github.com/ulule/limiter/v3"

rate := limiter.Rate{
    Period: 1 * time.Hour,
    Limit:  1000,
}

// You can also use the simplified format "<limit>-<period>"", with the given
// periods:
//
// * "S": second
// * "M": minute
// * "H": hour
// * "D": day
//
// Examples:
//
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
//
rate, err := limiter.NewRateFromFormatted("1000-H")
if err != nil {
    panic(err)
}

// Then, create a store. Here, we use the bundled Redis store. Any store
// compliant to limiter.Store interface will do the job. The defaults are
// "limiter" as Redis key prefix and a maximum of 3 retries for the key under
// race condition.
import "github.com/ulule/limiter/v3/drivers/store/redis"

store, err := redis.NewStore(client)
if err != nil {
    panic(err)
}

// Alternatively, you can pass options to the store with the "WithOptions"
// function. For example, for Redis store:
import "github.com/ulule/limiter/v3/drivers/store/redis"

store, err := redis.NewStoreWithOptions(pool, limiter.StoreOptions{
    Prefix:   "your_own_prefix",
})
if err != nil {
    panic(err)
}

// Or use a in-memory store with a goroutine which clears expired keys.
import "github.com/ulule/limiter/v3/drivers/store/memory"

store := memory.NewStore()

// Then, create the limiter instance which takes the store and the rate as arguments.
// Now, you can give this instance to any supported middleware.
instance := limiter.New(store, rate)

// Alternatively, you can pass options to the limiter instance with several options.
instance := limiter.New(store, rate, limiter.WithClientIPHeader("True-Client-IP"), limiter.WithIPv6Mask(mask))

// Finally, give the limiter instance to your middleware initializer.
import "github.com/ulule/limiter/v3/drivers/middleware/stdlib"

middleware := stdlib.NewMiddleware(instance)
```

See middleware examples:

- [HTTP](https://github.com/ulule/limiter-examples/tree/master/http/main.go)
- [Gin](https://github.com/ulule/limiter-examples/tree/master/gin/main.go)
- [Beego](https://github.com/ulule/limiter-examples/blob/master//beego/main.go)
- [Chi](https://github.com/ulule/limiter-examples/tree/master/chi/main.go)
- [Echo](https://github.com/ulule/limiter-examples/tree/master/echo/main.go)
- [Fasthttp](https://github.com/ulule/limiter-examples/tree/master/fasthttp/main.go)

## How it works

The ip address of the request is used as a key in the store.

If the key does not exist in the store we set a default
value with an expiration period.

You will find two stores:

- Redis: rely on [TTL](http://redis.io/commands/ttl) and incrementing the rate limit on each request.
- In-Memory: rely on a fork of [go-cache](https://github.com/patrickmn/go-cache) with a goroutine to clear expired keys using a default interval.

When the limit is reached, a `429` HTTP status code is sent.

## Limiter behind a reverse proxy

### Introduction

If your limiter is behind a reverse proxy, it could be difficult to obtain the "real" client IP.

Some reverse proxies, like AWS ALB, lets all header values through that it doesn't set itself.
Like for example, `True-Client-IP` and `X-Real-IP`.
Similarly, `X-Forwarded-For` is a list of comma-separated IPs that gets appended to by each traversed proxy.
The idea is that the first IP _(added by the first proxy)_ is the true client IP. Each subsequent IP is another proxy along the path.

An attacker can spoof either of those headers, which could be reported as a client IP.

By default, limiter doesn't trust any of those headers: you have to explicitly enable them in order to use them.
If you enable them, **you must always be aware** that any header added by any _(reverse)_ proxy not controlled
by you **are completely unreliable.**

### X-Forwarded-For

For example, if you make this request to your load balancer:
```bash
curl -X POST https://example.com/login -H "X-Forwarded-For: 1.2.3.4, 11.22.33.44"
```

And your server behind the load balancer obtain this:
```
X-Forwarded-For: 1.2.3.4, 11.22.33.44, <actual client IP>
```

That's mean you can't use `X-Forwarded-For` header, because it's **unreliable** and **untrustworthy**.
So keep `TrustForwardHeader` disabled in your limiter option.

However, if you have configured your reverse proxy to always remove/overwrite `X-Forwarded-For` and/or `X-Real-IP` headers
so that if you execute this _(same)_ request:
```bash
curl -X POST https://example.com/login -H "X-Forwarded-For: 1.2.3.4, 11.22.33.44"
```

And your server behind the load balancer obtain this:
```
X-Forwarded-For: <actual client IP>
```

Then, you can enable `TrustForwardHeader` in your limiter option.

### Custom header

Many CDN and Cloud providers add a custom header to define the client IP. Like for example, this non exhaustive list:

* `Fastly-Client-IP` from Fastly
* `CF-Connecting-IP` from Cloudflare
* `X-Azure-ClientIP` from Azure

You can use these headers using `ClientIPHeader` in your limiter option.

### None of the above

If none of the above solution are working, please use a custom `KeyGetter` in your middleware.

You can use this excellent article to help you define the best strategy depending on your network topology and your security need:
https://adam-p.ca/blog/2022/03/x-forwarded-for/

If you have any idea/suggestions on how we could simplify this steps, don't hesitate to raise an issue.
We would like some feedback on how we could implement this steps in the Limiter API.

Thank you.

## Why Yet Another Package

You could ask us: why yet another rate limit package?

Because existing packages did not suit our needs.

We tried a lot of alternatives:

1. [Throttled][1]. This package uses the generic cell-rate algorithm. To cite the
   documentation: _"The algorithm has been slightly modified from its usual form to
   support limiting with an additional quantity parameter, such as for limiting the
   number of bytes uploaded"_. It is brillant in term of algorithm but
   documentation is quite unclear at the moment, we don't need _burst_ feature for
   now, impossible to get a correct `After-Retry` (when limit exceeds, we can still
   make a few requests, because of the max burst) and it only supports `http.Handler`
   middleware (we use [Gin][4]). Currently, we only need to return `429`
   and `X-Ratelimit-*` headers for `n reqs/duration`.

2. [Speedbump][3]. Good package but maybe too lightweight. No `Reset` support,
   only one middleware for [Gin][4] framework and too Redis-coupled. We rather
   prefer to use a "store" approach.

3. [Tollbooth][5]. Good one too but does both too much and too little. It limits by
   remote IP, path, methods, custom headers and basic auth usernames... but does not
   provide any Redis support (only _in-memory_) and a ready-to-go middleware that sets
   `X-Ratelimit-*` headers. `tollbooth.LimitByRequest(limiter, r)` only returns an HTTP
   code.

4. [ratelimit][2]. Probably the closer to our needs but, once again, too
   lightweight, no middleware available and not active (last commit was in August
   2014). Some parts of code (Redis) comes from this project. It should deserve much
   more love.

There are other many packages on GitHub but most are either too lightweight, too
old (only support old Go versions) or unmaintained. So that's why we decided to
create yet another one.

## Contributing

- Ping us on twitter:
  - [@oibafsellig](https://twitter.com/oibafsellig)
  - [@thoas](https://twitter.com/thoas)
  - [@novln\_](https://twitter.com/novln_)
- Fork the [project](https://github.com/ulule/limiter)
- Fix [bugs](https://github.com/ulule/limiter/issues)

Don't hesitate ;)

[1]: https://github.com/throttled/throttled
[2]: https://github.com/r8k/ratelimit
[3]: https://github.com/etcinit/speedbump
[4]: https://github.com/gin-gonic/gin
[5]: https://github.com/didip/tollbooth
[6]: https://github.com/valyala/fasthttp
[godoc-url]: https://pkg.go.dev/github.com/ulule/limiter/v3
[godoc-img]: https://pkg.go.dev/badge/github.com/ulule/limiter/v3
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[goreport-url]: https://goreportcard.com/report/github.com/ulule/limiter
[goreport-img]: https://goreportcard.com/badge/github.com/ulule/limiter
[circle-url]: https://circleci.com/gh/ulule/limiter/tree/master
[circle-img]: https://circleci.com/gh/ulule/limiter.svg?style=shield&circle-token=baf62ec320dd871b3a4a7e67fa99530fbc877c99
