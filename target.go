package multiresolver

import (
	"net/url"

	"google.golang.org/grpc/resolver"
)

// ParseTarget parses target into resolver.Target URL
//
// If target is not a valid URL, it returns {Endpoint: URL: {Path: target}}.
func ParseTarget(target string) (ret resolver.Target) {
	u, err := url.Parse(target)
	if err != nil {
		return resolver.Target{URL: url.URL{Path: target}}
	}

	return resolver.Target{
		Scheme:    u.Scheme,
		Authority: u.Host,
		URL:       *u,
	}
}
