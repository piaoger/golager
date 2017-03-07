package main

import (
	"../qiniu"
)

func listbuckets() {
	qiniu.ConfigFromEnv()
	qiniu.ListBuckets()
}
func upload(from string, to string) {
	qiniu.ConfigFromEnv()
	qiniu.Upload(from, to)
}

func download(from string, to string) {
	qiniu.ConfigFromEnv()
	qiniu.Download(from, to)
}
