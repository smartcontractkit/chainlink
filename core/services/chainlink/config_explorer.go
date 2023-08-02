package chainlink

import (
	"net/url"

	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type explorerConfig struct {
	explorerURL *models.URL
	s           v2.ExplorerSecrets
}

func (e *explorerConfig) URL() *url.URL {
	u := (*url.URL)(e.explorerURL)
	if *u == zeroURL {
		u = nil
	}
	return u
}

func (e *explorerConfig) AccessKey() string {
	if e.s.AccessKey == nil {
		return ""
	}
	return string(*e.s.AccessKey)
}

func (e *explorerConfig) Secret() string {
	if e.s.Secret == nil {
		return ""
	}
	return string(*e.s.Secret)
}
