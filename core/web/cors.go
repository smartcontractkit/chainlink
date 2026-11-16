package web

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// localCORS preserves the behavior of the gin-contrib/cors configuration used
// by uiCorsHandler without retaining the wrapper and rs/cors modules
type localCORS struct {
	allowAllOrigins  bool
	allowOrigins     []string
	normalHeaders    http.Header
	preflightHeaders http.Header
}

func newLocalCORS(allowOrigins string) *localCORS {
	c := &localCORS{}
	if allowOrigins == "*" {
		c.allowAllOrigins = true
	} else if origins := strings.Split(allowOrigins, ","); len(origins) > 0 {
		seen := make(map[string]bool, len(origins))
		for _, origin := range origins {
			origin = strings.ToLower(strings.TrimSpace(origin))
			if !seen[origin] {
				seen[origin] = true
				c.allowOrigins = append(c.allowOrigins, origin)
			}
		}
	}

	c.normalHeaders = make(http.Header)
	c.normalHeaders.Set("Access-Control-Allow-Credentials", "true")
	c.normalHeaders.Set("Access-Control-Expose-Headers", "Content-Length")
	if c.allowAllOrigins {
		c.normalHeaders.Set("Access-Control-Allow-Origin", "*")
	} else {
		c.normalHeaders.Set("Vary", "Origin")
	}

	c.preflightHeaders = make(http.Header)
	c.preflightHeaders.Set("Access-Control-Allow-Credentials", "true")
	c.preflightHeaders.Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE")
	c.preflightHeaders.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept")
	c.preflightHeaders.Set(
		"Access-Control-Max-Age",
		strconv.FormatInt(int64(math.MaxInt32), 10),
	)
	if c.allowAllOrigins {
		c.preflightHeaders.Set("Access-Control-Allow-Origin", "*")
	} else {
		c.preflightHeaders.Add("Vary", "Origin")
		c.preflightHeaders.Add("Vary", "Access-Control-Request-Method")
		c.preflightHeaders.Add("Vary", "Access-Control-Request-Headers")
	}
	return c
}

func (c *localCORS) apply(ctx *gin.Context) {
	origin := ctx.Request.Header.Get("Origin")
	if origin == "" {
		return
	}
	host := ctx.Request.Host
	if origin == "http://"+host || origin == "https://"+host {
		return
	}
	if !c.allowAllOrigins && !c.allowed(origin) {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	header := ctx.Writer.Header()
	if ctx.Request.Method == http.MethodOptions {
		for key, value := range c.preflightHeaders {
			header[key] = value
		}
		ctx.AbortWithStatus(http.StatusNoContent)
	} else {
		for key, value := range c.normalHeaders {
			header[key] = value
		}
	}
	if !c.allowAllOrigins {
		ctx.Header("Access-Control-Allow-Origin", origin)
	}
}

func (c *localCORS) allowed(origin string) bool {
	for _, value := range c.allowOrigins {
		if value == origin {
			return true
		}
	}
	return false
}

func uiCorsHandler(allowOrigins string) gin.HandlerFunc {
	c := newLocalCORS(allowOrigins)
	return c.apply
}
