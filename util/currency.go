package util

const (
	USD = "USD"
	EUR = "EUR"
	IDR = "IDr"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, IDR:
		return true

	}
	return false
}
