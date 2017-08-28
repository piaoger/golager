package main

import (
	"../s3"
)

func listbuckets() {
	s3.ConfigFromEnv()
	s3.ListBuckets()
}

func upload(from string, to string, headers map[string]string) {
	s3.ConfigFromEnv()
	s3.Upload(from, to, headers)
}

func download(from string, to string) {
	s3.ConfigFromEnv()
	s3.Download(from, to)
}

func stat(path string) {
	s3.ConfigFromEnv()
	s3.Stat(path)
}

func listdir(path string) {
	s3.ConfigFromEnv()
	s3.ListDir(path)
}

func listfiles(path string) {
	s3.ConfigFromEnv()
	s3.ListFiles(path, []string{})
}
