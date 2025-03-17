package sortType

const (
	Desc = "Desc"
	ASC  = "Asc"
)

func IsSortType(sortType string) bool {
	if Desc == sortType ||
		ASC == sortType {
		return true
	}

	return false
}
