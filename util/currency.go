package util

// constants for all currency, may add here in the future
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupported(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	default:
		return false
	}
}
