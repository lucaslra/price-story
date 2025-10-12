package models

// User represents an application user
type User struct {
	ID                      string `json:"id"`
	Email                   string `json:"email"`
	PasswordHash            string `json:"password_hash"`
	PreferredCurrency       string `json:"preferred_currency"`
	DecimalPlaces           int    `json:"decimal_places"`
	ThousandSeparator       string `json:"thousand_separator"`
	CurrencySymbolPlacement string `json:"currency_symbol_placement"`
}
