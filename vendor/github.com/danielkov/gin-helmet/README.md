# Gin Helmet

Security middlewares for Gin (`gin-gonic/gin`) inspired by the popular `helmet` middleware package for Node JS `express` and `koa`.
___
[![Build Status](https://travis-ci.org/danielkov/gin-helmet.svg?branch=master)](https://travis-ci.org/danielkov/gin-helmet)
[![Coverage Status](https://coveralls.io/repos/github/danielkov/gin-helmet/badge.svg?branch=master)](https://coveralls.io/github/danielkov/gin-helmet?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/danielkov/gin-helmet)](https://goreportcard.com/report/github.com/danielkov/gin-helmet)
[![godocs](https://img.shields.io/badge/godocs-reference-blue.svg)](https://godoc.org/github.com/danielkov/gin-helmet)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

## Usage

Add the `Default` middleware for basic security measures.

```go
s := gin.New()
s.Use(helmet.Default())
```

You can also add each middleware separately:

```go
s.Use(helmet.NoCache())
```

Those not included in the `Default()` middleware are considered more advanced and require consideration before using.

See the [godoc](https://godoc.org/github.com/danielkov/gin-helmet) for more info and examples.