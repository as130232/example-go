package env

const (
	Test  = "test"
	Local = "local"
	Dev   = "dev"
	Cqa   = "cqa"
	Uat   = "uat"
	Prod  = "prod"
)

// ExecuteSyncEnv 執行基礎資料落地與發送kafka Sync同步訊息的環境
func ExecuteSyncEnv() string {
	return Prod
}
