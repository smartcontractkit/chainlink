package config

type Insecure interface {
	DevWebServer() bool
	OCRDevelopmentMode() bool
	DisableRateLimiting() bool
	InfiniteDepthQueries() bool
}
