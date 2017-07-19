package s3

import (
	"os"
)

var AWS_ACCESS_KEY_ID string
var AWS_SECRET_ACCESS_KEY string
var AWS_REGION string

func Config(region string, access_id string, access_secret string) {
	AWS_REGION = region
	AWS_ACCESS_KEY_ID = access_id
	AWS_SECRET_ACCESS_KEY = access_secret
}

func ConfigFromEnv() {
	AWS_REGION = os.Getenv("AWS_REGION")
	AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
}

func hasConfig() bool {
	return AWS_REGION != "" && AWS_ACCESS_KEY_ID != "" && AWS_SECRET_ACCESS_KEY != ""
}
