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
	FilePermMode = os.FileMode(0664) // Default file permission

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
	u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", AZURE_STORAGE_ACCOUNT, bucket, key))
	return u, err
}

// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#UploadFileToBlockBlob
func Upload(from string, to string, headers map[string]string, metas map[string]string) error {

	ctx := context.Background()

	retryTryTimeout := time.Second

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

	fmt.Println("not impl")
	return errors.New("not impl")
}

func Stat(path string) (map[string]interface{}, error) {

	fmt.Println("not impl")
	result := map[string]interface{}{}
	return result, errors.New("not impl")
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

	retryTryTimeout := time.Second
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
