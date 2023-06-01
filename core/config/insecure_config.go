package config

type Insecure interface {
	DevWebServer() bool
	InsecureFastScrypt() bool
	OCRDevelopmentMode() bool
	DisableRateLimiting() bool
	InfiniteDepthQueries() bool
}
