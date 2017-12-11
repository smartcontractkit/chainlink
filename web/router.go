package web

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/web/controllers"
)

func Router(store *services.Store) *gin.Engine {
	engine := gin.New()
	basicAuth := gin.BasicAuth(gin.Accounts{"chainlink": "boguspassword"})
	engine.Use(loggerFunc(), gin.Recovery(), basicAuth)

	j := controllers.JobsController{store}
	engine.POST("/jobs", j.Create)
	engine.GET("/jobs/:id", j.Show)

	jr := controllers.JobRunsController{store}
	engine.GET("/jobs/:id/runs", jr.Index)

	return engine
}

// Inspired by https://github.com/gin-gonic/gin/issues/961
func loggerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr := bytes.NewBuffer(buf)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		start := time.Now()
		c.Next()
		end := time.Now()

		logger.Infow("Web request",
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"body", readBody(rdr),
			"clientIP", c.ClientIP(),
			"comment", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"servedAt", end.Format("2006/01/02 - 15:04:05"),
			"latency", fmt.Sprintf("%v", end.Sub(start)),
		)
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
