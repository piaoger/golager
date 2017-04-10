package oss

import (
	"../utils"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strconv"
	"strings"
	"time"
)

func parseTime(ossdate string) time.Time {
	// date format
	// https://golang.org/src/time/format.go
	std := "Mon, 02 Jan 2006 15:04:05 GMT"
	d, err := time.Parse(std, ossdate)
	if err != nil {
		fmt.Printf("date parse error : %s", err)
		d, err = time.Parse(std, std)
	}

	return d
}

func listObjects(path string, prefix string, marker string, delimiter string) []oss.ListObjectsResult {
	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		// HandleError(err)
	}

	bucket_name, _, err := utils.ParseAddress(path)
	if err != nil {
		fmt.Printf("parse address error : %s", err)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		fmt.Printf("bucket error : %s", err)
	}

	results := []oss.ListObjectsResult{}

	ossprefix := oss.Prefix(prefix)
	ossmarker := oss.Marker(marker)
	ossdelimiter := oss.Delimiter(delimiter)
	for {
		lsRes, err := bucket.ListObjects(oss.MaxKeys(5), ossprefix, ossmarker, ossdelimiter)
		if err != nil {
			// HandleError(err)
		}

		ossprefix = oss.Prefix(lsRes.Prefix)
		ossmarker = oss.Marker(lsRes.NextMarker)

		results = append(results, lsRes)
		if !lsRes.IsTruncated {
			break
		}
	}

	return results
}

func listDir(path string) []string {

	dirs := []string{}

	results := listObjects(path, "", "", "/")

	for i := 0; i < len(results); i += 1 {
		lsRes := results[i]
		for j := 0; j < len(lsRes.CommonPrefixes); j += 1 {
			key := lsRes.CommonPrefixes[j]
			dirs = append(dirs, strings.TrimSuffix(key, "/"))
		}
	}

	fmt.Printf("dir: %s", dirs)

	return dirs
}

func listFiles(path string) []string {

	files := []string{}

	results := listObjects(path, "", "", "/")
	fmt.Println("results:", results)

	for i := 0; i < len(results); i += 1 {
		lsRes := results[i]
		for j := 0; j < len(lsRes.Objects); j += 1 {
			files = append(files, lsRes.Objects[j].Key)
		}
	}

	return files
}

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

func Stat(path string) (map[string]interface{}, error) {

	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		// HandleError(err)
	}

	bucket_name, key, err := utils.ParseAddress(path)
	if err != nil {
		fmt.Printf("parse address error : %s", err)
	}
	fmt.Printf("parse bucket_name : %s ,  key=%s", bucket_name, key)

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		fmt.Printf("bucket error : %s", err)
	}

	var result map[string]interface{}
	props, err := bucket.GetObjectMeta(key)
	if err != nil {
		fmt.Println("occurred error: %s", err)
		fmt.Printf("key: %s, size: %d kb, modified: %s \n", key, 0, "N/A")
	} else {
		size, _ := strconv.Atoi(props["Content-Length"][0])
		modified := parseTime(props["Last-Modified"][0])

		result = map[string]interface{}{
			"size":     size,
			"modified": modified,
		}

		fmt.Println("\nObject Meta:\n", props)
		fmt.Printf("\n")
		fmt.Printf("key: %s, size: %d kb, modified: %s \n", key, size/1024.0, utils.FormatTime(modified))
	}

	return result, err
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

func ListFiles(path string) []string {
	return listFiles(path)
}

func ListDir(path string) []string {
	return listDir(path)
}
