package s3

import (
	"os"
)

var AWS_ACCESS_KEY_ID string
var AWS_ACCESS_KEY_SECRET string
var AWS_REGION string

func Config(region string, access_id string, access_secret string) {
	AWS_REGION = region
	AWS_ACCESS_KEY_ID = access_id
	AWS_ACCESS_KEY_SECRET = access_secret
}

func ConfigFromEnv() {
	AWS_REGION = os.Getenv("AWS_REGION")
	AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_ACCESS_KEY_SECRET = os.Getenv("AWS_ACCESS_KEY_SECRET")
}

func hasConfig() bool {
	return AWS_REGION != "" && AWS_ACCESS_KEY_ID != "" && AWS_ACCESS_KEY_SECRET != ""
}
