package config

import "net/url"

type Explorer interface {
	ExplorerAccessKey() string
	ExplorerSecret() string
	ExplorerURL() *url.URL
}
