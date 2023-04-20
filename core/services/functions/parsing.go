package functions

import (
	"net/url"

	"mvdan.cc/xurls/v2"
)

func parseDomains(sourceCode string) (domains []string) {
	// https://pkg.go.dev/mvdan.cc/xurls/v2
	urls := xurls.Strict().FindAllString(sourceCode, -1)
	if urls == nil {
		return
	}
	domainsSet := make(map[string]struct{})
	for _, rawUrl := range urls {
		url, err := url.Parse(rawUrl)
		if err != nil {
			continue
		}
		domainsSet[url.Host] = struct{}{}
	}
	for domain := range domainsSet {
		domains = append(domains, domain)
	}
	return
}
