package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	VND = "VND"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, VND:
		return true
	}
	return false
}
