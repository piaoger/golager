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

func stat(path string) {
	qiniu.ConfigFromEnv()
	qiniu.Stat(path)
}

func listdir(path string) {
	qiniu.ConfigFromEnv()
	qiniu.ListDir(path)
}

func listfiles(path string) {
	qiniu.ConfigFromEnv()
	qiniu.ListFiles(path)
}
