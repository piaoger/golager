package qiniu

import (
	"os"
)

var QINIU_ACCESS_KEY string
var QINIU_SECRET_KEY string

func Config(access_id string, access_secret string) {
	QINIU_ACCESS_KEY = access_id
	QINIU_SECRET_KEY = access_secret
}

func ConfigFromEnv() {
	QINIU_ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	QINIU_SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
}

func hasConfig() bool {
	return QINIU_ACCESS_KEY != "" && QINIU_SECRET_KEY != ""
}
