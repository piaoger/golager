package obs

import (
	"../utils"
	"errors"
	"fmt"
	//"github.com/aliyun/aliyun-obs-go-sdk/obs"
	"github.com/piaoger/obs-sdk-go/obs"
	// "strconv"
	// "strings"
	"time"
	"os"
	"io"
)

// from: https://github.com/opentelekomcloud/obs-obsutil/
// from:

// Other constants
const (
	FilePermMode = os.FileMode(0664) // Default file permission

	TempFilePrefix = "obs-go-temp-"  // Temp file prefix
	TempFileSuffix = ".temp"         // Temp file suffix
)

// ContentType is an option to set Content-Type header
// ContentLength is an option to set Content-Length : ContentLength
// CacheControl is an option to set Cache-Control header
// ContentDisposition is an option to set Content-Disposition header
// ContentEncoding is an option to set Content-Encoding header
// ContentMD5 is an option to set Content-MD5 header
// Expires is an option to set Expires header

func headersToMetaInput(input *obs.SetObjectMetadataInput, headers map[string]string, metas map[string]string)  {

	for k, v := range headers {
		if k == "Content-Type" {
			input.ContentType = v
		}  else if k == "Cache-Control" {
			input.CacheControl = v
		} else if k == "Content-Disposition" {
			input.ContentDisposition = v
		} else if k == "Content-Encoding" {
			input.ContentEncoding = v
		} else if k == "Expires" {
	 		input.Expires = v
		} else {

		}
	}

	input.Metadata = metas

}

func newClient() (*obs.ObsClient, error) {
	endpoint := fmt.Sprintf("https://obs.%s.myhuaweicloud.com", HWCLOUD_REGION)
	return obs.New(HWCLOUD_ACCESS_KEY, HWCLOUD_SECRET_KEY, endpoint)
}

func parseTime(obsdate string) time.Time {
	// date format
	// https://golang.org/src/time/format.go
	std := "Mon, 02 Jan 2006 15:04:05 GMT"
	d, err := time.Parse(std, obsdate)
	if err != nil {
		d, err = time.Parse(std, std)
	}

	return d
}

// func safeListObjects(bucket_name string, obsprefix obs.Option, obsmarker obs.Option, obsdelimiter obs.Option) (obs.ListObjectsResult, error) {
// 	client, err := newClient()
// 	if err != nil {
// 		msg := fmt.Sprintf("obs client creation error : %s", err)
// 		return obs.ListObjectsResult{}, errors.New(msg)
// 	}

// 	bucket, err := client.Bucket(bucket_name)
// 	if err != nil {
// 		msg := fmt.Sprintf("bucket error error : %s", err)
// 		return obs.ListObjectsResult{}, errors.New(msg)
// 	}

// 	return bucket.ListObjects(obs.MaxKeys(800), obsprefix, obsmarker, obsdelimiter)
// }

// func listObjects(bucket_name string, key string, delimiter string, timeout int) []obs.ListObjectsResult {

// 	calls := 0

// 	results := []obs.ListObjectsResult{}

// 	obsprefix := obs.Prefix(key)
// 	obsmarker := obs.Marker("")
// 	obsdelimiter := obs.Delimiter(delimiter)
// 	for {
// 		lsRes, err := safeListObjects(bucket_name, obsprefix, obsmarker, obsdelimiter)
// 		if err != nil {
// 			break
// 		}

// 		obsprefix = obs.Prefix(lsRes.Prefix)
// 		obsmarker = obs.Marker(lsRes.NextMarker)

// 		if calls == 8 {
// 			utils.Sleep(2 * timeout)
// 		} else if calls == 16 {
// 			utils.Sleep(4 * timeout)
// 		} else if calls == 32 {
// 			utils.Sleep(6 * timeout)
// 		} else if calls == 64 {
// 			utils.Sleep(8 * timeout)
// 			calls = 0
// 		} else {
// 			utils.Sleep(500)
// 		}

// 		calls += 1

// 		results = append(results, lsRes)
// 		if !lsRes.IsTruncated {
// 			break
// 		}
// 	}
// 	return results
// }

func Upload(from string, to string, headers map[string]string, metas map[string]string) error {

	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("obs client creation error : %s", err)
		return errors.New(msg)
	}

	bucket_name, key, err := utils.ParseAddress(to)
	if err != nil {
		msg := fmt.Sprintf("parse address error : %s", err)
		return errors.New(msg)
	}

	fd, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fd.Close()

	input := &obs.PutObjectInput{}
	input.Bucket = bucket_name
	input.Key = key
	input.Body = fd
	input.ACL = obs.AclPublicRead

	_, err = client.PutObject(input)
	if err != nil {
		msg := fmt.Sprintf("PutObjectFromFile error : %s", err)
		return errors.New(msg)
	}

	// let's set meta data after that
	SetObjectMeta(to, headers , metas )

	return nil
}

func UploadDir(from string, to string, headers map[string]string, metas map[string]string) error {

	exists, err := utils.DirExists(from)

	if !exists || err != nil {
		msg := fmt.Sprintf("from directory does not exist: %s", from)
		return errors.New(msg)
	}

	// client, err := newClient()
	// if err != nil {
	// 	msg := fmt.Sprintf("obs client creation error : %s", err)
	// 	return errors.New(msg)
	// }

	// bucket_name, key, err := utils.ParseAddress(to)
	// if err != nil {
	// 	msg := fmt.Sprintf("parse address error : %s", err)
	// 	return errors.New(msg)
	// }

	collected := map[string]string{}
	utils.CollectFiles(from, to, false, collected)


	for k, v := range collected {
		err = Upload(k, v, headers, metas)
		if err != nil {
			msg := fmt.Sprintf("Upload error in uploading directory(%s): %s", from, err)
			return errors.New(msg)
		}
	}

	return nil
}

func Stat(path string) (map[string]interface{}, error) {

	var result map[string]interface{}
	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("obs client creation error : %s", err)
		return result, errors.New(msg)
	}

	bucket_name, key, err := utils.ParseAddress(path)
	if err != nil {
		msg := fmt.Sprintf("parse address error : %s", err)
		return result, errors.New(msg)
	}

	input := &obs.GetObjectMetadataInput{}
	input.Bucket = bucket_name
	input.Key = key

	output, err := client.GetObjectMetadata(input)
	if err != nil {
		panic(err)
	}
	if err != nil {
		result = map[string]interface{}{
			"size":         0,
			"modified":     time.Time{},
			"content-type": "N/A",
		}

		return result, err
	}

	result = map[string]interface{}{
		"size":         output.ContentLength,
		"modified":     output.LastModified,
		"content-type": output.ContentType,
	}

	return result, nil
}

// func ListBuckets() []string {
// 	client, err := newClient()
// 	if err != nil {
// 		return []string{}
// 	}

// 	// 列出Bucket，默认100条。
// 	lsRes, err := client.ListBuckets()
// 	if err != nil {
// 		return []string{}
// 	}

// 	buckets := make([]string, len(lsRes.Buckets))
// 	for i, b := range lsRes.Buckets {
// 		buckets[i] = b.Name

// 	}

// 	return buckets
// }

// func ListDir(path string) []string {

// 	dirs := []string{}
// 	bucket_name, key, _ := utils.ParseAddress(path)
// 	results := listObjects(bucket_name, key, "/", 200)

// 	for i := 0; i < len(results); i += 1 {
// 		lsRes := results[i]
// 		for j := 0; j < len(lsRes.CommonPrefixes); j += 1 {
// 			key := lsRes.CommonPrefixes[j]
// 			dirs = append(dirs, strings.TrimSuffix(key, "/"))
// 		}
// 	}

// 	return dirs
// }

// func ListFiles(path string, filters []string) []map[string]interface{} {

// 	files := []map[string]interface{}{}
// 	bucket_name, key, _ := utils.ParseAddress(path)
// 	results := listObjects(bucket_name, key, "", 200)
// 	for i := 0; i < len(results); i += 1 {
// 		lsRes := results[i]
// 		for j := 0; j < len(lsRes.Objects); j += 1 {
// 			obj := lsRes.Objects[j]
// 			// parts := strings.Split(obj.Key, "/")
// 			// name := parts[len(parts)-1]

// 			name := strings.Replace(obj.Key, key, "", -1)

// 			if name == "" || strings.Contains(name, "/") {
// 				continue
// 			}

// 			wanted := len(filters) == 0
// 			for fi := 0; fi < len(filters); fi += 1 {
// 				if name == filters[fi] {
// 					wanted = true
// 					break
// 				}
// 			}

// 			if wanted {
// 				fileinfo := map[string]interface{}{
// 					"name":     name,
// 					"size":     obj.Size,
// 					"modified": obj.LastModified,
// 				}
// 				files = append(files, fileinfo)
// 			}
// 		}
// 	}

// 	return files
// }

func Download(from string, to string) error {

	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("obs client creation error : %s", err)
		return errors.New(msg)
	}

	bucket_name, key, err := utils.ParseAddress(from)
	if err != nil {
		msg := fmt.Sprintf("parse address error : %s", err)
		return errors.New(msg)
	}

	input := &obs.GetObjectInput{}
	input.Bucket = bucket_name
	input.Key = key

	output, err := client.GetObject(input)
	defer output.Body.Close()
	if err != nil {
		msg := fmt.Sprintf("GetObject error : %s", err)
		return errors.New(msg)
	}

	tempFilePath := to + TempFileSuffix
	// If the local file does not exist, create a new one. If it exists, overwrite it.
	fd, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, FilePermMode)
	if err != nil {
		return err
	}

	// Copy the data to the local file path.
	_, err = io.Copy(fd, output.Body)
	fd.Close()
	if err != nil {
		return err
	}

	return os.Rename(tempFilePath, to)
}

// func CopyObject(from string, to string) error {

// 	client, err := newClient()
// 	if err != nil {
// 		msg := fmt.Sprintf("obs client creation error : %s", err)
// 		return errors.New(msg)
// 	}

// 	from_bucket_name, fromkey, err := utils.ParseAddress(from)
// 	if err != nil {
// 		msg := fmt.Sprintf("parse from address error : %s", err)
// 		return errors.New(msg)
// 	}

// 	frombucket, err := client.Bucket(from_bucket_name)
// 	if err != nil {
// 		msg := fmt.Sprintf("from bucket error : %s", err)
// 		return errors.New(msg)
// 	}

// 	to_bucket_name, tokey, err := utils.ParseAddress(to)
// 	if err != nil {
// 		msg := fmt.Sprintf("parse to address error : %s", err)
// 		return errors.New(msg)
// 	}

// 	_, err = frombucket.CopyObjectTo(to_bucket_name, tokey, fromkey)
// 	if err != nil {
// 		msg := fmt.Sprintf("CopyObjectTo error : %s", err)
// 		return errors.New(msg)
// 	}

// 	return nil

// }

func SetObjectMeta(path string, headers map[string]string, metas map[string]string) error {

	client, err := newClient()
	if err != nil {
		msg := fmt.Sprintf("obs client creation error : %s", err)
		return errors.New(msg)
	}

	bucket_name, key, err := utils.ParseAddress(path)
	if err != nil {
		msg := fmt.Sprintf("parse address error : %s", err)
		return errors.New(msg)
	}

	input := &obs.SetObjectMetadataInput{}
	input.Bucket = bucket_name
	input.Key = key

	headersToMetaInput(input, headers, metas)

	// options := headersToOption(headers, metas)
	_, err = client.SetObjectMetadata(input)
	if err != nil {
		msg := fmt.Sprintf("SetObjectMeta error : %s", err)
		return errors.New(msg)
	}

	return nil
}

// func SignUrl(path, method string, expiredInSec int64) (string, error) {

// 	client, err := newClient()
// 	if err != nil {
// 		msg := fmt.Sprintf("obs client creation error : %s", err)
// 		return "", errors.New(msg)
// 	}

// 	bucket_name, key, err := utils.ParseAddress(path)
// 	if err != nil {
// 		msg := fmt.Sprintf("parse address error : %s", err)
// 		return "", errors.New(msg)
// 	}

// 	bucket, err := client.Bucket(bucket_name)
// 	if err != nil {
// 		msg := fmt.Sprintf("bucket error : %s", err)
// 		return "", errors.New(msg)
// 	}

// 	m := strings.ToUpper(string(method))

// 	if "GET" != m && "PUT" != m && "HEAD" != m && "POST" != m && "DELETE" != m {
// 		msg := fmt.Sprintf("invalid method : %s", m)
// 		return "", errors.New(msg)
// 	}

// 	signedurl, err := bucket.SignURL(key, obs.HTTPMethod(m), expiredInSec)
// 	if err != nil {
// 		msg := fmt.Sprintf("sign url error : %s", err)
// 		return "", errors.New(msg)
// 	}

// 	return signedurl, nil
// }
