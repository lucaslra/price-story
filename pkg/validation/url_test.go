package validation

import "testing"

func TestIsValidURL_Valid(t *testing.T) {
    cases := []string{
        "http://example.com",
        "https://example.com",
        "https://example.com/path?x=1",
        "https://sub.domain.co.uk/",
        "https://192.168.0.1:8080/path",
        "HTTP://example.com", // scheme is case-insensitive
    }
    for _, u := range cases {
        if !IsValidURL(u) {
            t.Fatalf("expected valid URL: %s", u)
        }
    }
}

func TestIsValidURL_Invalid(t *testing.T) {
    cases := []string{
        "ftp://example.com", // unsupported scheme
        "http://",           // missing host
        "https://",          // missing host
        "://example.com",    // missing scheme
        "example.com",       // no scheme
        "http:///path",      // malformed host
        "https:// ",         // whitespace host
    }
    for _, u := range cases {
        if IsValidURL(u) {
            t.Fatalf("expected invalid URL: %s", u)
        }
    }
}