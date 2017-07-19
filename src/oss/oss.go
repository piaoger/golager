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

func safeListObjects(bucket_name string, ossprefix oss.Option, ossmarker oss.Option, ossdelimiter oss.Option) (oss.ListObjectsResult, error) {
	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		fmt.Printf("oss client error : %s", err)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		fmt.Printf("bucket error : %s", err)
	}

	return bucket.ListObjects(oss.MaxKeys(800), ossprefix, ossmarker, ossdelimiter)
}

func listObjects(bucket_name string, key string, delimiter string, timeout int) []oss.ListObjectsResult {

	calls := 0

	results := []oss.ListObjectsResult{}

	ossprefix := oss.Prefix(key)
	ossmarker := oss.Marker("")
	ossdelimiter := oss.Delimiter(delimiter)
	for {
		lsRes, err := safeListObjects(bucket_name, ossprefix, ossmarker, ossdelimiter)
		if err != nil {
			fmt.Printf("ListObjects error : %s", err)
		}

		ossprefix = oss.Prefix(lsRes.Prefix)
		ossmarker = oss.Marker(lsRes.NextMarker)

		if calls == 8 {
			utils.Sleep(2 * timeout)
		} else if calls == 16 {
			utils.Sleep(4 * timeout)
		} else if calls == 32 {
			utils.Sleep(6 * timeout)
		} else if calls == 64 {
			utils.Sleep(8 * timeout)
			calls = 0
		} else {
			utils.Sleep(500)
		}

		calls += 1

		results = append(results, lsRes)
		if !lsRes.IsTruncated {
			break
		}
	}
	return results
}

// func ParseAddress(address string) (string, string, error) {
//  return bucketkey(address)
// }

func Upload(from string, to string) error {
	//fmt.Printf("action: from: %s, to: %s \n", from, to)

	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		fmt.Printf("oss client error : %s", err)
	}

	bucket_name, key, err := utils.ParseAddress(to)
	if err != nil {
		msg := fmt.Sprintf("parse address error : %s", err)
		return errors.New(msg)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		msg := fmt.Sprintf("bucket error : %s", err)
		return errors.New(msg)
	}

	err = bucket.PutObjectFromFile(key, from, oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		msg := fmt.Sprintf("PutObjectFromFile error : %s", err)
		return errors.New(msg)
	}
	return nil
}

func Stat(path string) (map[string]interface{}, error) {

	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		fmt.Printf("oss client error : %s", err)
	}

	bucket_name, key, err := utils.ParseAddress(path)
	if err != nil {
		fmt.Printf("parse address error : %s", err)
	}
	//fmt.Printf("parse bucket_name : %s ,  key=%s", bucket_name, key)

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
		contentType := props["Content-Type"][0]

		//fmt.Printf("props: %s \n", props)

		result = map[string]interface{}{
			"size":         size,
			"modified":     modified,
			"content-type": contentType,
		}

	}

	return result, err
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
		return []string{}
	}

	buckets := make([]string, len(lsRes.Buckets))
	for i, b := range lsRes.Buckets {
		buckets[i] = b.Name

	}

	return buckets
}

func ListDir(path string) []string {

	dirs := []string{}
	bucket_name, key, _ := utils.ParseAddress(path)
	results := listObjects(bucket_name, key, "/", 200)

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

func ListFiles(path string, filters []string) []map[string]interface{} {

	files := []map[string]interface{}{}
	bucket_name, key, _ := utils.ParseAddress(path)
	results := listObjects(bucket_name, key, "", 200)
	//fmt.Printf("filters: %s, %d\n", filters, len(filters))
	for i := 0; i < len(results); i += 1 {
		lsRes := results[i]
		for j := 0; j < len(lsRes.Objects); j += 1 {
			obj := lsRes.Objects[j]
			// parts := strings.Split(obj.Key, "/")
			// name := parts[len(parts)-1]

			name := strings.Replace(obj.Key, key, "", -1)

			if name == "" || strings.Contains(name, "/") {
				continue
			}

			fmt.Printf("obj.key: %s path:%s \n", obj.Key, path)

			wanted := len(filters) == 0
			for fi := 0; fi < len(filters); fi += 1 {
				if name == filters[fi] {
					wanted = true
					break
				}
			}

			if wanted {
				fileinfo := map[string]interface{}{
					"name":     name,
					"size":     obj.Size,
					"modified": obj.LastModified,
				}
				files = append(files, fileinfo)
			}

		}
	}

	return files
}

func Download(from string, to string) {

	client, err := oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
	if err != nil {
		fmt.Printf("oss client error : %s", err)
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
