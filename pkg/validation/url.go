package validation

import (
	"net/url"
	"strings"
)

// IsValidURL checks that the URL parses and has http(s) scheme and host
func IsValidURL(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}
