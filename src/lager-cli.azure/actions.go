package main

import (
	"../azure"
	"encoding/json"
	"fmt"
)

func listbuckets() {
	azure.ConfigFromEnv()
	buckets := azure.ListBuckets()
	fmt.Println(buckets)
}

func upload(from string, to string, headers map[string]string, metas map[string]string) {
	azure.ConfigFromEnv()
	err := azure.Upload(from, to, headers, metas)

	if err != nil {
		fmt.Println("error: %s", err.Error())
	} else {
		fmt.Println("upload success")
	}
}

func uploaddir(from string, to string, headers map[string]string, metas map[string]string) {
	azure.ConfigFromEnv()
	err := azure.UploadDir(from, to, headers, metas)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("upload success")
	}
}

func download(from string, to string) {
	azure.ConfigFromEnv()
	err := azure.Download(from, to)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("download success")
	}
}

func copyObject(from string, to string) {
	azure.ConfigFromEnv()
	err := azure.CopyObject(from, to)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("copy success")
	}
}

func stat(path string) {
	azure.ConfigFromEnv()
	stat, _ := azure.Stat(path)
	fmt.Println(stat)
}

func listdir(path string) {
	azure.ConfigFromEnv()
	dirs := azure.ListDir(path)

	desc, _ := json.Marshal(dirs)
	fmt.Println(string(desc))
}

func listfiles(path string) {
	azure.ConfigFromEnv()
	files := azure.ListFiles(path, []string{})
	fmt.Println(files)
}

func setObjectMeta(path string, headers map[string]string, metas map[string]string) {
	azure.ConfigFromEnv()
	err := azure.SetObjectMeta(path, headers, metas)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("set object meta success")
	}
}

func signUrl(path string, method string) {
	azure.ConfigFromEnv()
	signedurl, err := azure.SignUrl(path, method, 7200)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(signedurl)
	}
}
