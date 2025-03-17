package languageType

const (
	Cn = "zh-cn"
	Tw = "zh-tw"
	En = "en-us"
	Id = "id-id"
	Vn = "vn-vn"
	My = "my"
	Ja = "ja"
	Ko = "ko"
	Th = "th-th"
	Hi = "hi"
	Pt = "pt"
)

var LanguageTypes = []string{Cn, Tw, En, Id, Vn, My, Ja, Ko, Th, Hi, Pt}

func IsValidLanguageType(languageType string) bool {
	for _, lang := range LanguageTypes {
		if lang == languageType {
			return true
		}
	}

	return false
}
