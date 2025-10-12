package validation

import "regexp"

// IsUUID validates standard UUIDs in 8-4-4-4-12 hex format.
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func IsUUID(s string) bool {
	if s == "" {
		return false
	}
	return uuidRegex.MatchString(s)
}
