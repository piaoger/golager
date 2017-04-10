package qiniu

import (
	"../utils"
	"fmt"
	"os"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
	"time"
)

func domains(bucket string) []string {
	// qshell中使用v6实现了该功能
	// https://github.com/qiniu/api.v7/issues/64
	return []string{
		os.Getenv("QINIU_DOMAIN"),
	}
}

func downloadUrl(domain string, key string) string {
	// 调用MakeBaseUrl()方法将domain,key处理成http://domain/key的形式
	baseUrl := kodo.MakeBaseUrl(domain, key)
	policy := kodo.GetPolicy{}

	c := kodo.New(0, nil)
	return c.MakePrivateUrl(baseUrl, &policy)
}

func stat(bucket string, key string) {

	c := kodo.New(0, nil)
	p := c.Bucket(bucket)
	entry, err := p.Stat(nil, key)
	fmt.Println(entry)
	if err != nil {
		fmt.Println(err)
	}
}

func listBucket(bucket string) {

	// new一个Bucket对象
	c := kodo.New(0, nil)
	p := c.Bucket(bucket)
	ListItem, _, _, err := p.List(nil, "", "", "", 100)
	if err == nil {
		fmt.Println("List success")
	} else {
		fmt.Println("List failed:", err)
	}

	for _, item := range ListItem {
		fmt.Println(item.Key, item.Fsize)
	}

}

type PutRet struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

func ListBuckets() []string {

	return []string{}
}

func Upload(from string, to string) {

	conf.ACCESS_KEY = QINIU_ACCESS_KEY
	conf.SECRET_KEY = QINIU_SECRET_KEY

	bucket, key, _ := utils.ParseAddress(to)

	c := kodo.New(0, nil)
	policy := &kodo.PutPolicy{
		Scope: bucket,
		//设置Token过期时间
		Expires: 3600,
	}

	token := c.MakeUptoken(policy)
	zone := 0
	uploader := kodocli.NewUploader(zone, nil)

	var ret PutRet

	fmt.Printf("bucket: %s , key: %s, local: %s ", bucket, key, from)
	fmt.Printf("ACCESS_KEY: %s , SECRET_KEY: %s", conf.ACCESS_KEY, conf.SECRET_KEY)

	res := uploader.PutFile(nil, &ret, token, key, from, nil)
	fmt.Println(ret)
	if res != nil {
		fmt.Println("io.Put failed:", res)
		return
	}

}

func Download(from string, to string) {
	conf.ACCESS_KEY = QINIU_ACCESS_KEY
	conf.SECRET_KEY = QINIU_SECRET_KEY
	bucket, key, _ := utils.ParseAddress(from)
	domain := domains(bucket)[0]
	url := downloadUrl(domain, key)
	fmt.Printf("download url: %s\n", url)
	utils.DownloadFromUrl(url, to)
}

func Stat(path string) (map[string]interface{}, error) {
	conf.ACCESS_KEY = QINIU_ACCESS_KEY
	conf.SECRET_KEY = QINIU_SECRET_KEY

	bucket, key, err := utils.ParseAddress(path)
	stat(bucket, key)

	size := 0
	modified := time.Now()

	result := map[string]interface{}{
		"size":     size,
		"modified": modified,
	}

	return result, err
}

func ListFiles(path string) []string {
	//not impl now
	return []string{}
}

func ListDir(path string) []string {
	//not impl now
	return []string{}
}
