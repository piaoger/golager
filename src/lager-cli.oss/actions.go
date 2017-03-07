package main

import (
	"../oss"
)

func listbuckets() {
	oss.ConfigFromEnv()
	oss.ListBuckets()
}
func upload(from string, to string) {
	oss.ConfigFromEnv()
	oss.Upload(from, to)
}

func download(from string, to string) {
	oss.ConfigFromEnv()
	oss.Download(from, to)
}
