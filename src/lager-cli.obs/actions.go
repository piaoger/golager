package main

import (
	"../obs"
	//"encoding/json"
	"fmt"
)

func listbuckets() {
	// obs.ConfigFromEnv()
	// buckets := obs.ListBuckets()
	// fmt.Println(buckets)
}

func upload(from string, to string, headers map[string]string, metas map[string]string) {
	obs.ConfigFromEnv()
	err := obs.Upload(from, to, headers, metas)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("upload success")
	}
}

func uploaddir(from string, to string, headers map[string]string, metas map[string]string) {
	obs.ConfigFromEnv()
	err := obs.UploadDir(from, to, headers, metas)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("upload success")
	}
}

func download(from string, to string) {
	obs.ConfigFromEnv()
	err := obs.Download(from, to)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("download success")
	}
}

func copyObject(from string, to string) {
	// obs.ConfigFromEnv()
	// err := obs.CopyObject(from, to)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("copy success")
	// }
}

func stat(path string) {
	obs.ConfigFromEnv()
	stat, _ := obs.Stat(path)
	fmt.Println(stat)
}

func listdir(path string) {
	// obs.ConfigFromEnv()
	// dirs := obs.ListDir(path)

	// desc, _ := json.Marshal(dirs)
	// fmt.Println(string(desc))
}

func listfiles(path string) {
	// obs.ConfigFromEnv()
	// files := obs.ListFiles(path, []string{})
	// fmt.Println(files)
}

func setObjectMeta(path string, headers map[string]string, metas map[string]string) {
	// obs.ConfigFromEnv()
	// err := obs.SetObjectMeta(path, headers, metas)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("set object meta success")
	// }
}

func signUrl(path string, method string) {
	// obs.ConfigFromEnv()
	// signedurl, err := obs.SignUrl(path, method, 7200)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(signedurl)
	// }
}
