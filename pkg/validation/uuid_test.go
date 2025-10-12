package validation

import "testing"

func TestIsUUID_Valid(t *testing.T) {
    valid := "123e4567-e89b-12d3-a456-426614174000"
    if !IsUUID(valid) {
        t.Fatalf("expected valid UUID: %s", valid)
    }
}

func TestIsUUID_Invalid(t *testing.T) {
    cases := []string{
        "",
        "not-a-uuid",
        "123e4567e89b12d3a456426614174000",    // missing hyphens
        "123e4567-e89b-12d3-a456-42661417400",  // too short
        "123e4567-e89b-12d3-a456-4266141740000", // too long
        "zzze4567-e89b-12d3-a456-426614174000", // non-hex chars
    }
    for _, s := range cases {
        if IsUUID(s) {
            t.Fatalf("expected invalid UUID: %s", s)
        }
    }
}