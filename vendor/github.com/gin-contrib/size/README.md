# size

[![Run Tests](https://github.com/gin-contrib/size/actions/workflows/go.yml/badge.svg)](https://github.com/gin-contrib/size/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gin-contrib/size/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-contrib/size)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/size)](https://goreportcard.com/report/github.com/gin-contrib/size)
[![GoDoc](https://godoc.org/github.com/gin-contrib/size?status.svg)](https://godoc.org/github.com/gin-contrib/size)
[![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin)

Limit size of POST requests for Gin framework

## Example

```go
package main

import (
  "net/http"

  limits "github.com/gin-contrib/size"
  "github.com/gin-gonic/gin"
)

func handler(ctx *gin.Context) {
  val := ctx.PostForm("b")
  if len(ctx.Errors) > 0 {
    return
  }
  ctx.String(http.StatusOK, "got %s\n", val)
}

func main() {
  r := gin.Default()
  r.Use(limits.RequestSizeLimiter(10))
  r.POST("/", handler)
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}
```
