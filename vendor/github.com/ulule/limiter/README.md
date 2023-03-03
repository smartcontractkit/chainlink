# Limiter

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][circle-img]][circle-url]
[![Go Report Card][goreport-img]][goreport-url]

*Dead simple rate limit middleware for Go.*

* Simple API
* "Store" approach for backend
* Redis support (but not tied too)
* Middlewares: HTTP and [Gin][4]

## Installation

Using [Go Modules](https://github.com/golang/go/wiki/Modules)

```bash
$ go get github.com/ulule/limiter/v3@v3.1.0
```

## Usage

In five steps:

* Create a `limiter.Rate` instance _(the number of requests per period)_
* Create a `limiter.Store` instance _(see [Redis](https://github.com/ulule/limiter/blob/master/drivers/store/redis/store.go) or [In-Memory](https://github.com/ulule/limiter/blob/master/drivers/store/memory/store.go))_
* Create a `limiter.Limiter` instance that takes store and rate instances as arguments
* Create a middleware instance using the middleware of your choice
* Give the limiter instance to your middleware initializer

**Example:**

```go
// Create a rate with the given limit (number of requests) for the given
// period (a time.Duration of your choice).
import "github.com/ulule/limiter"

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
//
// Examples:
//
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
//
rate, err := limiter.NewRateFromFormatted("1000-H")
if err != nil {
    panic(err)
}

// Then, create a store. Here, we use the bundled Redis store. Any store
// compliant to limiter.Store interface will do the job. The defaults are
// "limiter" as Redis key prefix and a maximum of 3 retries for the key under
// race condition.
import "github.com/ulule/limiter/drivers/store/redis"

store, err := redis.NewStore(client)
if err != nil {
    panic(err)
}

// Alternatively, you can pass options to the store with the "WithOptions"
// function. For example, for Redis store:
import "github.com/ulule/limiter/drivers/store/redis"

store, err := redis.NewStoreWithOptions(pool, limiter.StoreOptions{
    Prefix:   "your_own_prefix",
    MaxRetry: 4,
})
if err != nil {
    panic(err)
}

// Or use a in-memory store with a goroutine which clears expired keys.
import "github.com/ulule/limiter/drivers/store/memory"

store := memory.NewStore()

// Then, create the limiter instance which takes the store and the rate as arguments.
// Now, you can give this instance to any supported middleware.
instance := limiter.New(store, rate)
```

See middleware examples:

* [HTTP](https://github.com/ulule/limiter/tree/master/examples/http/main.go)
* [Gin](https://github.com/ulule/limiter/tree/master/examples/gin/main.go)
* [Beego](https://github.com/ulule/limiter/blob/master/examples/beego/main.go)
* [Chi](https://github.com/ulule/limiter/tree/master/examples/chi/main.go)


## How it works

The ip address of the request is used as a key in the store.

If the key does not exist in the store we set a default
value with an expiration period.

You will find two stores:

* Redis: rely on [TTL](http://redis.io/commands/ttl) and incrementing the rate limit on each request.
* In-Memory: rely on a fork of [go-cache](https://github.com/patrickmn/go-cache) with a goroutine to clear expired keys using a default interval.

When the limit is reached, a `429` HTTP status code is sent.

## Why Yet Another Package

You could ask us: why yet another rate limit package?

Because existing packages did not suit our needs.

We tried a lot of alternatives:

1. [Throttled][1]. This package uses the generic cell-rate algorithm. To cite the
documentation: *"The algorithm has been slightly modified from its usual form to
support limiting with an additional quantity parameter, such as for limiting the
number of bytes uploaded"*. It is brillant in term of algorithm but
documentation is quite unclear at the moment, we don't need *burst* feature for
now, impossible to get a correct `After-Retry` (when limit exceeds, we can still
make a few requests, because of the max burst) and it only supports ``http.Handler``
middleware (we use [Gin][4]). Currently, we only need to return `429`
and `X-Ratelimit-*` headers for `n reqs/duration`.

2. [Speedbump][3]. Good package but maybe too lightweight. No `Reset` support,
only one middleware for [Gin][4] framework and too Redis-coupled. We rather
prefer to use a "store" approach.

3. [Tollbooth][5]. Good one too but does both too much and too little. It limits by
remote IP, path, methods, custom headers and basic auth usernames... but does not
provide any Redis support (only *in-memory*) and a ready-to-go middleware that sets
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

* Ping us on twitter:
  * [@oibafsellig](https://twitter.com/oibafsellig)
  * [@thoas](https://twitter.com/thoas)
  * [@novln_](https://twitter.com/novln_)
* Fork the [project](https://github.com/ulule/limiter)
* Fix [bugs](https://github.com/ulule/limiter/issues)

Don't hesitate ;)

[1]: https://github.com/throttled/throttled
[2]: https://github.com/r8k/ratelimit
[3]: https://github.com/etcinit/speedbump
[4]: https://github.com/gin-gonic/gin
[5]: https://github.com/didip/tollbooth

[godoc-url]: https://godoc.org/github.com/ulule/limiter
[godoc-img]: https://godoc.org/github.com/ulule/limiter?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[goreport-url]: https://goreportcard.com/report/github.com/ulule/limiter
[goreport-img]: https://goreportcard.com/badge/github.com/ulule/limiter
[circle-url]: https://circleci.com/gh/ulule/limiter/tree/master
[circle-img]: https://circleci.com/gh/ulule/limiter.svg?style=shield&circle-token=baf62ec320dd871b3a4a7e67fa99530fbc877c99
