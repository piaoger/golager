package oss

import (
	"../utils"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func createClient() (*oss.Client, error) {

	if !hasConfig() {
		return nil, errors.New("no valid configrations")
	}

	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ListBuckets() []string {
	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		// HandleError(err)
	}

	// 列出Bucket，默认100条。
	lsRes, err := client.ListBuckets()
	if err != nil {
		fmt.Printf("bucket error : %s", err)
	}
	fmt.Println("buckets:", lsRes.Buckets)

	return []string{}
}

func Upload(from string, to string) {
	fmt.Printf("action: from: %s, to: %s \n", from, to)

	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		// HandleError(err)
	}

	bucket_name, key, err := utils.ParseAddress(to)
	if err != nil {
		fmt.Printf("parse address error : %s", err)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		fmt.Printf("bucket error : %s", err)
	}

	err = bucket.PutObjectFromFile(key, from)
	if err != nil {
		fmt.Printf("PutObjectFromFile error : %s", err)
	}

}

func Download(from string, to string) {

	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		// HandleError(err)
	}

	bucket_name, key, err := utils.ParseAddress(from)
	if err != nil {
		fmt.Printf("parse address error : %s", err)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		fmt.Printf("bucket error : %s", err)
	}

	err = bucket.GetObjectToFile(key, to)
	if err != nil {
		fmt.Printf("GetObjectToFile error : %s", err)
	}

}
