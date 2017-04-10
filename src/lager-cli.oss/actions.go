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

func stat(path string) {
	oss.ConfigFromEnv()
	oss.Stat(path)
}

func listdir(path string) {
	oss.ConfigFromEnv()
	oss.ListDir(path)
}

func listfiles(path string) {
	oss.ConfigFromEnv()
	oss.ListFiles(path)
}
