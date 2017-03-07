package oss

import (
	"os"
)

var OSS_ACCESS_KEY_ID string
var OSS_ACCESS_KEY_SECRET string
var OSS_ENDPOINT string

func Config(endpoint string, access_id string, access_secret string) {
	OSS_ENDPOINT = endpoint
	OSS_ACCESS_KEY_ID = access_id
	OSS_ACCESS_KEY_SECRET = access_secret
}

func ConfigFromEnv() {
	OSS_ENDPOINT = os.Getenv("OSS_ENDPOINT")
	OSS_ACCESS_KEY_ID = os.Getenv("OSS_ACCESS_KEY_ID")
	OSS_ACCESS_KEY_SECRET = os.Getenv("OSS_ACCESS_KEY_SECRET")
}

func hasConfig() bool {
	return OSS_ENDPOINT != "" && OSS_ACCESS_KEY_ID != "" && OSS_ACCESS_KEY_SECRET != ""
}
