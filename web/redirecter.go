package web

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Redirector returns a simple web server for redirecting a browser to a HTTPS
// version of the site.
func Redirector(port uint16) error {
	httpRouter := gin.Default()

	httpRouter.GET("/*path", func(c *gin.Context) {
		req := c.Request
		tlsURL := url.URL{
			Scheme:   "https",
			Path:     req.URL.Path,
			RawQuery: req.URL.RawQuery,
			Host:     strings.Split(req.Host, ":")[0],
		}

		c.Redirect(302, tlsURL.String())
	})

	listenURL := fmt.Sprintf(":%d", port)
	return httpRouter.Run(listenURL)
}
