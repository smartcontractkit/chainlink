package config

import "net/url"

type Explorer interface {
	AccessKey() string
	Secret() string
	URL() *url.URL
}
