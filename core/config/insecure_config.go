package config

type Insecure interface {
	DevWebServer() bool
	OCRDevelopmentMode() bool
	DisableRateLimiting() bool
	DisableSSRFProtection() bool
	InfiniteDepthQueries() bool
}
