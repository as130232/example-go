package currencyType

const (
	CNY = "CNY"
	USD = "USD"
)

func IsCurrencyType(currencyType string) bool {
	if CNY == currencyType ||
		USD == currencyType {
		return true
	}

	return false
}

func GetAllCurrencyTypes() []string {
	return []string{CNY, USD}
}
