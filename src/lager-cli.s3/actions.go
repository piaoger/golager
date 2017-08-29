package main

import (
	"../s3"
	"fmt"
)

func listbuckets() {
	s3.ConfigFromEnv()
	buckets := s3.ListBuckets()
	fmt.Println(buckets)
}

func upload(from string, to string, headers map[string]string, metas map[string]string) {
	s3.ConfigFromEnv()
	err := s3.Upload(from, to, headers, metas)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("upload success")
	}
}

func download(from string, to string) {
	s3.ConfigFromEnv()
	err := s3.Download(from, to)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("download success")
	}
}

func stat(path string) {
	s3.ConfigFromEnv()
	stat, _ := s3.Stat(path)
	fmt.Println(stat)
}

func listdir(path string) {
	s3.ConfigFromEnv()
	dirs := s3.ListDir(path)
	fmt.Println(dirs)
}

func listfiles(path string) {
	s3.ConfigFromEnv()
	files := s3.ListFiles(path, []string{})
	fmt.Println(files)
}
