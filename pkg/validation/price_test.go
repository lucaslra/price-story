package validation

import (
    "math"
    "testing"
)

func TestIsValidPrice_Valid(t *testing.T) {
    cases := []float64{1, 1.23, math.SmallestNonzeroFloat64}
    for _, p := range cases {
        if !IsValidPrice(p) {
            t.Fatalf("expected valid price: %v", p)
        }
    }
}

func TestIsValidPrice_Invalid(t *testing.T) {
    cases := []float64{0, -1, -0.01, math.NaN(), math.Inf(1), math.Inf(-1)}
    for _, p := range cases {
        if IsValidPrice(p) {
            t.Fatalf("expected invalid price: %v", p)
        }
    }
}