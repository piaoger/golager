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

func headersToOption(headers map[string]string, metas map[string]string) []oss.Option {
	options := []oss.Option{}
	for k, v := range headers {
		if k == "Content-Type" {
			options = append(options, oss.ContentType(v))
		} else if k == "Content-Length" {
			cl, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				options = append(options, oss.ContentLength(cl))
			}
		} else if k == "Cache-Control" {
			options = append(options, oss.CacheControl(v))
		} else if k == "Content-Disposition" {
			options = append(options, oss.ContentDisposition(v))
		} else if k == "Content-Encoding" {
			options = append(options, oss.ContentEncoding(v))
		} else if k == "Content-MD5" {
			options = append(options, oss.ContentMD5(v))
		} else if k == "Expires" {
			// exp, err := strconv.ParseInt(v, 10, 64)
			// if err == nil {
			// 	options = append(options, oss.Expires(exp))
			// }
		} else {
			//options = append(options, oss.Meta(k, v))
		}
	}

	for k, v := range metas {
		options = append(options, oss.Meta(k, v))
	}

	return options
}

func newClient() (*oss.Client, error) {
	return oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET)
}

func newClientWithoutHashVerify() (*oss.Client, error) {
	return oss.New(OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET, oss.EnableCRC(false), oss.EnableMD5(false))
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
			fmt.Println(err.Error())
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

func Upload(from string, to string, headers map[string]string, metas map[string]string) error {

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

	options := headersToOption(headers, metas)
	options = append(options, oss.ObjectACL(oss.ACLPublicRead))
	err = bucket.PutObjectFromFile(key, from, options...)
	if err != nil {
		msg := fmt.Sprintf("PutObjectFromFile error : %s", err)
		return errors.New(msg)
	}
	return nil
}

func UploadDir(from string, to string, headers map[string]string, metas map[string]string) error {

	exists, err := utils.DirExists(from)

	if !exists || err != nil {
		msg := fmt.Sprintf("from directory does not exist: %s", from)
		return errors.New(msg)
	}

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

	collected := map[string]string{}
	utils.CollectFiles(from, key, false, collected)

	options := headersToOption(headers, metas)
	options = append(options, oss.ObjectACL(oss.ACLPublicRead))

	for k, v := range collected {
		err = bucket.PutObjectFromFile(v, k, options...)
		if err != nil {
			msg := fmt.Sprintf("PutObjectFromFile error in uploading directory(%s): %s", from, err)
			return errors.New(msg)
		}
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
			"size":             0,
			"modified":         time.Time{},
			"content-encoding": "N/A",
			"cacthe-control":   "N/A",
			"content-type":     "N/A",
			"Content-Encoding": "",
			"metadata":         map[string]string{},
		}

		return result, err
	}

	fmt.Println("%v", props)

	contentLength := props["Content-Length"][0]
	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		size = 0
	}

	metas := map[string]string{}
	for k, v := range props {

		if strings.Contains(k, "X-Oss-Meta-") {
			nk := strings.TrimPrefix(k, "X-Oss-Meta-")
			metas[nk] = v[0]
		}
	}

	modified := parseTime(props["Last-Modified"][0])

	contentType := props["Content-Type"][0]
	contentMd5 := props["Content-Md5"][0]

	cacheControl := props["Cache-Control"][0]
	contentEncoding := ""

	ceprops := props["Content-Encoding"]
	if len(ceprops) > 0 {
		contentEncoding = ceprops[0]
	}

	result = map[string]interface{}{
		"size":             size,
		"modified":         modified,
		"content-type":     contentType,
		"content-md5":      contentMd5,
		"content-encoding": contentEncoding,
		"cacthe-control":   cacheControl,

		"metadata": metas,
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

func ListObjects(path string) []string {

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

	client, err := newClientWithoutHashVerify()
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

func Delete(path string) error {
	return errors.New("not impl")
}

func CopyObject(from string, to string) error {

	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("oss client creation error : %s", err)
		return errors.New(msg)
	}

	from_bucket_name, fromkey, err := utils.ParseAddress(from)
	if err != nil {
		msg := fmt.Sprintf("parse from address error : %s", err)
		return errors.New(msg)
	}

	frombucket, err := client.Bucket(from_bucket_name)
	if err != nil {
		msg := fmt.Sprintf("from bucket error : %s", err)
		return errors.New(msg)
	}

	to_bucket_name, tokey, err := utils.ParseAddress(to)
	if err != nil {
		msg := fmt.Sprintf("parse to address error : %s", err)
		return errors.New(msg)
	}
	fmt.Printf("CopyObjectTo : %s, %s, %s, %s\n", from_bucket_name, fromkey, to_bucket_name, tokey)
	_, err = frombucket.CopyObjectTo(to_bucket_name, tokey, fromkey)
	if err != nil {
		msg := fmt.Sprintf("CopyObjectTo error : %s", err)
		return errors.New(msg)
	}

	return nil

}

func SetObjectMeta(path string, headers map[string]string, metas map[string]string) error {
	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("oss client creation error : %s", err)
		return errors.New(msg)
	}

	bucket_name, key, err := utils.ParseAddress(path)
	if err != nil {
		msg := fmt.Sprintf("parse address error : %s", err)
		return errors.New(msg)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		msg := fmt.Sprintf("bucket error : %s", err)
		return errors.New(msg)
	}

	options := headersToOption(headers, metas)
	err = bucket.SetObjectMeta(key, options...)
	if err != nil {
		msg := fmt.Sprintf("SetObjectMeta error : %s", err)
		return errors.New(msg)
	}

	return nil
}

func SignUrl(path string, method string, expiredInSec int64) (string, error) {

	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("oss client creation error : %s", err)
		return "", errors.New(msg)
	}

	bucket_name, key, err := utils.ParseAddress(path)
	if err != nil {
		msg := fmt.Sprintf("parse address error : %s", err)
		return "", errors.New(msg)
	}

	bucket, err := client.Bucket(bucket_name)
	if err != nil {
		msg := fmt.Sprintf("bucket error : %s", err)
		return "", errors.New(msg)
	}

	m := strings.ToUpper(string(method))

	if "GET" != m && "PUT" != m && "HEAD" != m && "POST" != m && "DELETE" != m {
		msg := fmt.Sprintf("invalid method : %s", m)
		return "", errors.New(msg)
	}

	signedurl, err := bucket.SignURL(key, oss.HTTPMethod(m), expiredInSec)
	if err != nil {
		msg := fmt.Sprintf("sign url error : %s", err)
		return "", errors.New(msg)
	}

	return signedurl, nil
}
