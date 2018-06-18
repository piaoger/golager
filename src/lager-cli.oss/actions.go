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

func upload(from string, to string, headers map[string]string, metas map[string]string) {
	oss.ConfigFromEnv()
	err := oss.Upload(from, to, headers, metas)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("upload success")
	}

}

func uploaddir(from string, to string, headers map[string]string, metas map[string]string) {
	oss.ConfigFromEnv()
	err := oss.UploadDir(from, to, headers, metas)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("upload success")
	}

}

func download(from string, to string) {
	oss.ConfigFromEnv()
	err := oss.Download(from, to)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("download success")
	}
}

func copyObject(from string, to string) {
	oss.ConfigFromEnv()
	err := oss.CopyObject(from, to)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("copy success")
	}
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

func setObjectMeta(path string, headers map[string]string, metas map[string]string) {
	oss.ConfigFromEnv()
	err := oss.SetObjectMeta(path, headers, metas)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("set object meta success")
	}
}

func signUrl(path string, method string) {
	oss.ConfigFromEnv()
	signedurl, err := oss.SignUrl(path, method, 7200)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(signedurl)
	}
}
