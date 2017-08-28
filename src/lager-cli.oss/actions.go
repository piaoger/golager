package main

import (
	"../oss"
	"fmt"
)

func listbuckets() {
	oss.ConfigFromEnv()
	buckets := oss.ListBuckets()
	fmt.Println(buckets)
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
	stat, _ := oss.Stat(path)
	fmt.Println(stat)
}

func listdir(path string) {
	oss.ConfigFromEnv()
	dirs := oss.ListDir(path)
	fmt.Println(dirs)
}

func listfiles(path string) {
	oss.ConfigFromEnv()
	files := oss.ListFiles(path, []string{})
	fmt.Println(files)
}
