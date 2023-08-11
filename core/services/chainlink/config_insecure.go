package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type insecureConfig struct {
	c toml.Insecure
}

func (i *insecureConfig) DevWebServer() bool {
	return build.IsDev() && i.c.DevWebServer != nil &&
		*i.c.DevWebServer
}

func (i *insecureConfig) DisableRateLimiting() bool {
	return build.IsDev() && i.c.DisableRateLimiting != nil &&
		*i.c.DisableRateLimiting
}

func (i *insecureConfig) OCRDevelopmentMode() bool {
	// OCRDevelopmentMode is allowed in TestBuilds as well
	return (build.IsDev() || build.IsTest()) && i.c.OCRDevelopmentMode != nil &&
		*i.c.OCRDevelopmentMode
}

func (i *insecureConfig) InfiniteDepthQueries() bool {
	return build.IsDev() && i.c.InfiniteDepthQueries != nil &&
		*i.c.InfiniteDepthQueries
}
