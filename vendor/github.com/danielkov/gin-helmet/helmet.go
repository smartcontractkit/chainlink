package helmet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// NoSniff applies header to protect your server from MimeType Sniffing
func NoSniff() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	}
}

// DNSPrefetchControl sets Prefetch Control header to prevent browser from prefetching DNS
func DNSPrefetchControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-DNS-Prefetch-Control", "off")
	}
}

// FrameGuard sets Frame Options header to deny to prevent content from the website to be served in an iframe
func FrameGuard(opt ...string) gin.HandlerFunc {
	var o string
	if len(opt) > 0 {
		o = opt[0]
	} else {
		o = "DENY"
	}
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", o)
	}
}

// SetHSTS Sets Strict Transport Security header to the default of 60 days
// an optional integer may be added as a parameter to set the amount in seconds
func SetHSTS(sub bool, opt ...int) gin.HandlerFunc {
	var o int
	if len(opt) > 0 {
		o = opt[0]
	} else {
		o = 5184000
	}
	op := "max-age=" + strconv.Itoa(o)
	if sub {
		op += "; includeSubDomains"
	}
	return func(c *gin.Context) {
		c.Writer.Header().Set("Strict-Transport-Security", op)
	}
}

// IENoOpen sets Download Options header for Internet Explorer to prevent it from executing downloads in the site's context
func IENoOpen() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Download-Options", "noopen")
	}
}

// XSSFilter applies very minimal XSS protection via setting the XSS Protection header on
func XSSFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
	}
}

// Default returns a number of handlers that are advised to use for basic HTTP(s) protection
func Default() (gin.HandlerFunc, gin.HandlerFunc, gin.HandlerFunc, gin.HandlerFunc, gin.HandlerFunc, gin.HandlerFunc) {
	return NoSniff(), DNSPrefetchControl(), FrameGuard(), SetHSTS(true), IENoOpen(), XSSFilter()
}

// Referrer sets the Referrer Policy header to prevent the browser from sending data from your website to another one upon navigation
// an optional string can be provided to set the policy to something else other than "no-referrer".
func Referrer(opt ...string) gin.HandlerFunc {
	var o string
	if len(opt) > 0 {
		o = opt[0]
	} else {
		o = "no-referrer"
	}
	return func(c *gin.Context) {
		c.Writer.Header().Set("Referrer-Policy", o)
	}
}

// NoCache obliterates cache options by setting a number of headers. This prevents the browser from storing your assets in cache
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Surrogate-Control", "no-store")
		c.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
	}
}

// ContentSecurityPolicy sets a header which will restrict your browser to only allow certain sources for assets on your website
// The function accepts a map of its parameters which are appended to the header so you can control which headers should be set
// The second parameter of the function is a boolean, which set to true will tell the handler to also set legacy headers, like
// those that work in older versions of Chrome and Firefox.
/*
Example usage:
    opts := map[string]string{
	    "default-src": "'self'",
	    "img-src": "*",
	    "media-src": "media1.com media2.com",
	    "script-src": "userscripts.example.com"
    }
	s.Use(helmet.ContentSecurityPolicy(opts, true))

See [Content Security Policy on MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP) for more info.
*/
func ContentSecurityPolicy(opt map[string]string, legacy bool) gin.HandlerFunc {
	policy := ""
	for k, v := range opt {
		policy += fmt.Sprintf("%s %s; ", k, v)
	}
	policy = strings.TrimSuffix(policy, "; ")
	return func(c *gin.Context) {
		if legacy {
			c.Writer.Header().Set("X-Webkit-CSP", policy)
			c.Writer.Header().Set("X-Content-Security-Policy", policy)
		}
		c.Writer.Header().Set("Content-Security-Policy", policy)
	}
}

// ExpectCT sets Certificate Transparency header which can enforce that you're using a Certificate which is ready for the
// upcoming Chrome requirements policy. The function accepts a maxAge int which is the TTL for the policy in delta seconds,
// an enforce boolean, which simply adds an enforce directive to the policy (otherwise it's report-only mode) and a
// optional reportUri, which is the URI to which report information is sent when the policy is violated.
func ExpectCT(maxAge int, enforce bool, reportURI ...string) gin.HandlerFunc {
	policy := ""
	if enforce {
		policy += "enforce, "
	}
	if len(reportURI) > 0 {
		policy += fmt.Sprintf("report-uri=%s, ", reportURI[0])
	}
	policy += fmt.Sprintf("max-age=%d", maxAge)
	return func(c *gin.Context) {
		c.Writer.Header().Set("Expect-CT", policy)
	}
}

// SetHPKP sets HTTP Public Key Pinning for your server. It is not necessarily a great thing to set this without proper
// knowledge of what this does. [Read here](https://developer.mozilla.org/en-US/docs/Web/HTTP/Public_Key_Pinning) otherwise you
// may likely end up DoS-ing your own server and domain. The function accepts a map of directives and their values according
// to specifications.
/*
Example usage:

	keys := []string{"cUPcTAZWKaASuYWhhneDttWpY3oBAkE3h2+soZS7sWs=", "M8HztCzM3elUxkcjR2S5P4hhyBNf6lHkmjAHKhpGPWE="}
	r := gin.New()
	r.Use(SetHPKP(keys, 5184000, true, "domain.com"))

*/
func SetHPKP(keys []string, maxAge int, sub bool, reportURI ...string) gin.HandlerFunc {
	policy := ""
	for _, v := range keys {
		policy += fmt.Sprintf("pin-sha256=\"%s\"; ", v)
	}
	policy += fmt.Sprintf("max-age=%d; ", maxAge)
	if sub {
		policy += "includeSubDomains; "
	}
	if len(reportURI) > 0 {
		policy += fmt.Sprintf("report-uri=\"%s\"", reportURI[0])
	}
	policy = strings.TrimSuffix(policy, "; ")
	return func(c *gin.Context) {
		c.Writer.Header().Set("Public-Key-Pins", policy)
	}
}
