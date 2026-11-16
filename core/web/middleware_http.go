package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type maxBytesReader struct {
	ctx        *gin.Context
	rdr        io.ReadCloser
	remaining  int64
	wasAborted bool
	sawEOF     bool
}

func (mbr *maxBytesReader) tooLarge() (int, error) {
	err := fmt.Errorf("HTTP request too large")
	if !mbr.wasAborted {
		mbr.wasAborted = true
		_ = mbr.ctx.Error(err)
		mbr.ctx.Header("connection", "close")
		mbr.ctx.String(http.StatusRequestEntityTooLarge, "request too large")
		mbr.ctx.AbortWithStatus(http.StatusRequestEntityTooLarge)
	}
	return 0, err
}

func (mbr *maxBytesReader) Read(p []byte) (int, error) {
	toRead := mbr.remaining
	if mbr.remaining == 0 {
		if mbr.sawEOF {
			return mbr.tooLarge()
		}
		toRead = 1
	}
	if int64(len(p)) > toRead {
		p = p[:toRead]
	}
	n, err := mbr.rdr.Read(p)
	if err == io.EOF {
		mbr.sawEOF = true
	}
	if mbr.remaining == 0 {
		if n > 0 {
			return mbr.tooLarge()
		}
		return 0, err
	}
	mbr.remaining -= int64(n)
	if mbr.remaining < 0 {
		mbr.remaining = 0
	}
	return n, err
}

func (mbr *maxBytesReader) Close() error {
	return mbr.rdr.Close()
}

func requestSizeLimiter(limit int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Body = &maxBytesReader{
			ctx:       ctx,
			rdr:       ctx.Request.Body,
			remaining: limit,
		}
		ctx.Next()
	}
}

// secureMiddleware preserves the behavior configured by the prior
// unrolled/secure use: HTTPS redirect outside development and frame denial.
func secureMiddleware(tlsRedirect bool, tlsHost string, devWebServer bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		isSSL := c.Request.TLS != nil || c.Request.URL.Scheme == "https"
		if tlsRedirect && !isSSL && !devWebServer {
			host := c.Request.Host
			if tlsHost != "" {
				host = tlsHost
			}
			url := *c.Request.URL
			url.Scheme = "https"
			url.Host = host
			http.Redirect(c.Writer, c.Request, url.String(), http.StatusMovedPermanently)
			c.Abort()
			return
		}

		c.Header("X-Frame-Options", "DENY")
		c.Next()
	}
}
