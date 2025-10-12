package validation

import "strings"

// IsValidCurrency validates ISO 4217-like 3-letter currency codes (uppercase).
// Supported baseline set can be extended as needed.
func IsValidCurrency(code string) bool {
	c := strings.ToUpper(strings.TrimSpace(code))
	switch c {
	case "USD", "EUR", "JPY", "GBP", "AUD", "CAD", "CHF", "CNY":
		return true
	default:
		// Allow any 3-letter uppercase code to avoid being overly restrictive
		if len(c) == 3 {
			for _, r := range c {
				if r < 'A' || r > 'Z' {
					return false
				}
			}
			return true
		}
		return false
	}
}
