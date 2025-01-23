package util

// Constants for all support currencies
const (
	USD = "USD"
	EUR = "EUR"
	THB = "THB"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currecy string) bool {
	switch currecy {
	case USD, EUR, THB:
		return true
	}
	return false
}
