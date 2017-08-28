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

// ContentType is an option to set Content-Type header
// ContentLength is an option to set Content-Length : ContentLength
// CacheControl is an option to set Cache-Control header
// ContentDisposition is an option to set Content-Disposition header
// ContentEncoding is an option to set Content-Encoding header
// ContentMD5 is an option to set Content-MD5 header
// Expires is an option to set Expires header

func headersToOption(headers map[string]string) []oss.Option {
	options := []oss.Option{}

	for k, v := range headers {
		if k == "Content-Type" {
			options = append(options, oss.ContentType(v))
		} else if k == "Content-Length" {
			// options = append(options, oss.ContentLength(v))
		} else if k == "Cache-Control" {
			options = append(options, oss.CacheControl(v))
		} else if k == "Content-Disposition" {
			options = append(options, oss.ContentDisposition(v))
		} else if k == "Content-Encoding" {
			options = append(options, oss.ContentEncoding(v))
		} else if k == "Content-MD5" {
			options = append(options, oss.ContentMD5(v))
		} else if k == "Expires" {
			//options = append(options, oss.Expires(v))
		} else {
			//options = append(options, oss.Meta(k, v))
		}
	}

	return options
}

func newClient() (*oss.Client, error) {
	return oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
}
func parseTime(ossdate string) time.Time {
	// date format
	// https://golang.org/src/time/format.go
	std := "Mon, 02 Jan 2006 15:04:05 GMT"
	d, err := time.Parse(std, ossdate)
	if err != nil {
		d, err = time.Parse(std, std)
	}

	return d
}

func safeListObjects(bucket_name string, ossprefix oss.Option, ossmarker oss.Option, ossdelimiter oss.Option) (oss.ListObjectsResult, error) {
	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("oss client creation error : %s", err)
		return oss.ListObjectsResult{}, errors.New(msg)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		msg := fmt.Sprintf("bucket error error : %s", err)
		return oss.ListObjectsResult{}, errors.New(msg)
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
			break
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

func Upload(from string, to string, headers map[string]string) error {

	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("oss client creation error : %s", err)
		return errors.New(msg)
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

	options := headersToOption(headers)
	options = append(options, oss.ObjectACL(oss.ACLPublicRead))
	err = bucket.PutObjectFromFile(key, from, options...)
	if err != nil {
		msg := fmt.Sprintf("PutObjectFromFile error : %s", err)
		return errors.New(msg)
	}
	return nil
}

func Stat(path string) (map[string]interface{}, error) {

	var result map[string]interface{}
	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("oss client creation error : %s", err)
		return result, errors.New(msg)
	}

	bucket_name, key, err := utils.ParseAddress(path)
	if err != nil {
		msg := fmt.Sprintf("parse address error : %s", err)
		return result, errors.New(msg)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		msg := fmt.Sprintf("bucket error : %s", err)
		return result, errors.New(msg)
	}

	props, err := bucket.GetObjectDetailedMeta(key)
	if err != nil {
		result = map[string]interface{}{
			"size":         0,
			"modified":     time.Time{},
			"content-type": "N/A",
		}

		return result, err
	}

	size, _ := strconv.Atoi(props["Content-Length"][0])
	modified := parseTime(props["Last-Modified"][0])
	contentType := props["Content-Type"][0]

	result = map[string]interface{}{
		"size":         size,
		"modified":     modified,
		"content-type": contentType,
	}

	return result, nil
}

func ListBuckets() []string {
	client, err := newClient()
	if err != nil {
		return []string{}
	}

	// 列出Bucket，默认100条。
	lsRes, err := client.ListBuckets()
	if err != nil {
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

	return dirs
}

func ListFiles(path string, filters []string) []map[string]interface{} {

	files := []map[string]interface{}{}
	bucket_name, key, _ := utils.ParseAddress(path)
	results := listObjects(bucket_name, key, "", 200)
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

func Download(from string, to string) error {

	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("oss client creation error : %s", err)
		return errors.New(msg)
	}

	bucket_name, key, err := utils.ParseAddress(from)
	if err != nil {
		msg := fmt.Sprintf("parse address error : %s", err)
		return errors.New(msg)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		msg := fmt.Sprintf("bucket error : %s", err)
		return errors.New(msg)
	}

	err = bucket.GetObjectToFile(key, to)
	if err != nil {
		msg := fmt.Sprintf("GetObjectToFile error : %s", err)
		return errors.New(msg)
	}

	return nil

}
