package validation

import "math"

// IsValidPrice checks for a positive, finite price value
func IsValidPrice(p float64) bool {
    if math.IsNaN(p) || math.IsInf(p, 0) {
        return false
    }
    return p > 0
}