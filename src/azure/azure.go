package azure

import (
	"../utils"
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"net/url"
	"os"
	"time"
)

const (
	FilePermMode    = os.FileMode(0664) // Default file permission
	RetryTryTimeout = 20 * time.Second
)

func toBlobHeaders(headers map[string]string) azblob.BlobHTTPHeaders {
	options := azblob.BlobHTTPHeaders{}
	for k, v := range headers {
		if k == "Content-Type" {
			options.ContentType = v
		} else if k == "Cache-Control" {
			options.CacheControl = v
		} else if k == "Content-Disposition" {
			options.ContentDisposition = v
		} else if k == "Content-Encoding" {
			options.ContentEncoding = v
		} else if k == "Content-MD5" {
			// TODO
			//options.ContentMD5 = v
		} else if k == "Content-Language" {
			options.ContentLanguage = v
		} else {
			//options = append(options, oss.Meta(k, v))
		}
	}

	return options
}

func newCredential() (*azblob.SharedKeyCredential, error) {
	credential, err := azblob.NewSharedKeyCredential(AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_ACCESS_KEY)
	return credential, err
}

func buildUrl(bucket string, key string) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", bucket, key))
	return u, err
}

// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#UploadFileToBlockBlob
func Upload(from string, to string, headers map[string]string, metas map[string]string) error {

	ctx := context.Background()

	retryTryTimeout := RetryTryTimeout

	credential, err := newCredential()
	if err != nil {
		fmt.Println("auth: %s", err.Error())
		return err
	}

	bucket_name, key, err := utils.ParseAddress(to)

	u, _ := buildUrl(bucket_name, key)
	blockBlobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			TryTimeout: retryTryTimeout,
		},
	}))

	file, err := os.Open(from)
	if err != nil {
		fmt.Println("open file error: %s", err.Error())
		return err
	}

	_, err = azblob.UploadFileToBlockBlob(ctx, file, blockBlobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:       4 * 1024 * 1024,
		BlobHTTPHeaders: toBlobHeaders(headers),
		Metadata:        metas,

		Parallelism: 16})

	return err
}

func UploadDir(from string, to string, headers map[string]string, metas map[string]string) error {

	fmt.Println("azure::UploadDir - not impl")
	return errors.New("azure::UploadDir - not impl")
}

func Stat(path string) (map[string]interface{}, error) {

	result := map[string]interface{}{
		"size":             0,
		"modified":         time.Time{},
		"content-encoding": "N/A",
		"cacthe-control":   "N/A",
		"content-type":     "N/A",
		"Content-Encoding": "",
		"metadata":         map[string]string{},
	}

	ctx := context.Background()
	retryTryTimeout := RetryTryTimeout

	credential, err := newCredential()
	if err != nil {
		fmt.Println("auth: %s", err.Error())
		return result, err
	}

	bucket_name, key, err := utils.ParseAddress(path)
	if err != nil {
		return result, err
	}

	u, _ := buildUrl(bucket_name, key)
	blockBlobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			TryTimeout: retryTryTimeout,
		},
	}))

	resp, err := blockBlobURL.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return result, err
	}

	metas := resp.NewMetadata()

	result = map[string]interface{}{
		"size":             resp.ContentLength(),
		"modified":         resp.LastModified(),
		"content-type":     resp.ContentType(),
		"content-md5":      resp.ContentMD5(),
		"content-encoding": resp.ContentEncoding(),
		"cacthe-control":   resp.CacheControl(),

		"metadata": metas,
	}

	return result, err

}

func Delete(path string) error {

	ctx := context.Background()
	retryTryTimeout := RetryTryTimeout

	credential, err := newCredential()
	if err != nil {
		fmt.Println("auth: %s", err.Error())
		return err
	}

	bucket_name, key, err := utils.ParseAddress(path)
	if err != nil {
		return err
	}

	u, _ := buildUrl(bucket_name, key)
	blockBlobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			TryTimeout: retryTryTimeout,
		},
	}))

	_, err = blockBlobURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})

	return err
}

func ListBuckets() []string {
	fmt.Println("not impl")
	return []string{}
}

func ListDir(path string) []string {

	fmt.Println("not impl")
	return []string{}
}

func ListFiles(path string, filters []string) []map[string]interface{} {

	fmt.Println("not impl")
	return []map[string]interface{}{}
}

func Download(from string, to string) error {

	retryTryTimeout := RetryTryTimeout
	ctx := context.Background()
	credential, err := newCredential()
	if err != nil {
		fmt.Println("auth: %s", err.Error())
		return err
	}

	bucket_name, key, err := utils.ParseAddress(from)
	if err != nil {
		return err
	}

	u, _ := buildUrl(bucket_name, key)

	blockBlobURL := azblob.NewBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			TryTimeout: retryTryTimeout,
		},
	}))

	// if snapshot != nil && !snapshot.Equal(time.Time{}) {
	// 	blockBlobURL = blockBlobURL.WithSnapshot(snapshot.Format(SnapshotTimeFormat))
	// }
	options := azblob.DownloadFromBlobOptions{
		RetryReaderOptionsPerBlock: azblob.RetryReaderOptions{MaxRetryRequests: 20},
	}

	// If the local file does not exist, create a new one. If it exists, overwrite it.
	fd, err := os.OpenFile(to, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, FilePermMode)
	if err != nil {
		return err
	}

	return azblob.DownloadBlobToFile(ctx, blockBlobURL, 0, 0, fd, options)
}

func CopyObject(from string, to string) error {
	err := errors.New("not impl")
	return err
}

func SetObjectMeta(path string, headers map[string]string, metas map[string]string) error {

	err := errors.New("not impl")
	return err
}

func SignUrl(path, method string, expiredInSec int64) (string, error) {
	err := errors.New("not impl")
	return "", err
}
