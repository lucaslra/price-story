package validation

import "fmt"

// ValidateNumberFormat validates formatting fields align with the currency and constraints.
func ValidateNumberFormat(currency string, decimalPlaces int, thousandSeparator string, symbolPlacement string) error {
	if !IsValidCurrency(currency) {
		return fmt.Errorf("invalid currency code")
	}
	if currency == "JPY" {
		if decimalPlaces != 0 {
			return fmt.Errorf("JPY requires 0 decimal places")
		}
	} else {
		if decimalPlaces < 0 || decimalPlaces > 4 {
			return fmt.Errorf("decimal places must be between 0 and 4")
		}
	}
	switch thousandSeparator {
	case ",", ".", " ":
		// ok
	default:
		return fmt.Errorf("invalid thousand separator")
	}
	switch symbolPlacement {
	case "before", "after":
		// ok
	default:
		return fmt.Errorf("invalid currency symbol placement")
	}
	return nil
}
