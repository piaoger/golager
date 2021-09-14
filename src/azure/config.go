package azure

import (
	"os"
)

var AZURE_STORAGE_ACCOUNT string
var AZURE_STORAGE_ACCESS_KEY string

func Config(storage_account string, access_key string) {
	AZURE_STORAGE_ACCOUNT = storage_account
	AZURE_STORAGE_ACCESS_KEY = access_key
}

func ConfigFromEnv() {
	AZURE_STORAGE_ACCOUNT = os.Getenv("AZURE_STORAGE_ACCOUNT")
	AZURE_STORAGE_ACCESS_KEY = os.Getenv("AZURE_STORAGE_ACCESS_KEY")
}

func hasConfig() bool {
	return AZURE_STORAGE_ACCOUNT != "" && AZURE_STORAGE_ACCESS_KEY != ""
}
