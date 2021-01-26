package obs

import (
	"os"
)

var HWCLOUD_ACCESS_KEY string
var HWCLOUD_SECRET_KEY string
var HWCLOUD_REGION string

func Config(region string, access_id string, access_secret string) {
	HWCLOUD_REGION = region
	HWCLOUD_ACCESS_KEY = access_id
	HWCLOUD_SECRET_KEY = access_secret
}

// https://rheem-dev.obs.cn-east-3.myhuaweicloud.com/tools/obs-browser-plus-3.20.8.dmg
// regions
func ConfigFromEnv() {
	HWCLOUD_REGION = os.Getenv("HWCLOUD_REGION")
	HWCLOUD_ACCESS_KEY = os.Getenv("HWCLOUD_ACCESS_KEY")
	HWCLOUD_SECRET_KEY = os.Getenv("HWCLOUD_SECRET_KEY")
}

func hasConfig() bool {
	return HWCLOUD_REGION != "" && HWCLOUD_ACCESS_KEY != "" && HWCLOUD_SECRET_KEY != ""
}
