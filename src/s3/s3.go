package s3

import (
	"../utils"
	//"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"os"
	"strings"
)

// atomic file write
func getObjectToFile(resp *s3.GetObjectOutput, filePath string) error {
	tempFilePath := filePath + ".TempFileSuffix"

	defer resp.Body.Close()

	fd, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = io.Copy(fd, resp.Body)
	fd.Close()
	if err != nil {
		return err
	}

	return os.Rename(tempFilePath, filePath)
}

func newS3Service() *s3.S3 {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	return svc
}

// most are borrowed from aws s3 sdk samples:
// github.com/aws/aws-sdk-go/service/s3/examples_test.go

func listObjects(bucket string, prefix string, marker string, delimiter string) []*s3.ListObjectsOutput {

	svc := newS3Service()

	s3prefix := prefix
	s3marker := marker
	s3delimiter := delimiter

	results := []*s3.ListObjectsOutput{}

	// iterate if the result is truncated.
	for {

		params := &s3.ListObjectsInput{
			Bucket:    aws.String(bucket), // Required
			Delimiter: aws.String(s3delimiter),
			Marker:    aws.String(s3marker),
			MaxKeys:   aws.Int64(100),
			Prefix:    aws.String(s3prefix),
		}

		resp, err := svc.ListObjects(params)
		if err != nil {
			//fmt.Println(err.Error())
			continue
		}

		s3prefix = *resp.Prefix
		s3marker = *resp.Marker

		results = append(results, resp)
		if *resp.IsTruncated == false {
			break
		}
	}

	return results
}

func ListBuckets() []string {
	svc := newS3Service()

	var params *s3.ListBucketsInput
	resp, _ := svc.ListBuckets(params)

	buckets := make([]string, len(resp.Buckets))
	for i, b := range resp.Buckets {
		buckets[i] = *b.Name

	}

	fmt.Println(buckets)

	return buckets

}

func Upload(from string, to string) {

	svc := newS3Service()

	bucket_name, key, err := utils.ParseAddress(to)

	fd, err := os.Open(from)
	if err != nil {
		return
	}
	defer fd.Close()

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket_name), // Required
		Key:    aws.String(key),         // Required
		Body:   fd,
		ACL:    aws.String("public-read"),
	}

	resp, err := svc.PutObject(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func Download(from string, to string) {
	svc := newS3Service()

	bucket_name, key, err := utils.ParseAddress(from)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket_name), // Required
		Key:    aws.String(key),         // Required
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		//fmt.Println(err.Error())
		return
	}

	getObjectToFile(resp, to)

}

func Stat(path string) (map[string]interface{}, error) {
	svc := newS3Service()

	bucket_name, key, err := utils.ParseAddress(path)

	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket_name), // Required
		Key:    aws.String(key),         // Required
	}
	resp, err := svc.HeadObject(params)

	var result map[string]interface{}

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		//fmt.Println(err.Error())
		return result, err

	} else {

		result = map[string]interface{}{
			"size":     *resp.ContentLength,
			"modified": *resp.LastModified,
		}

	}

	return result, nil
}

func ListFiles(path string) []string {

	bucket_name, key, _ := utils.ParseAddress(path)

	results := listObjects(bucket_name, key, "", "")
	keys := []string{}

	for i := 0; i < len(results); i += 1 {
		lsRes := results[i]
		for j := 0; j < len(lsRes.Contents); j += 1 {

			ckey := *lsRes.Contents[j].Key

			name := strings.Replace(ckey, key, "", -1)

			if strings.Contains(name, "/") {
				// not a file in this folder
				//trace.Trace('name : ' + name + ' is a in subfolder');
				continue
			}
			keys = append(keys, name)
		}
	}
	fmt.Println(keys)
	return keys
}

func ListDir(path string) []string {

	bucket_name, key, _ := utils.ParseAddress(path)

	results := listObjects(bucket_name, key, "", "/")

	keys := []string{}

	for i := 0; i < len(results); i += 1 {
		lsRes := results[i]
		for j := 0; j < len(lsRes.CommonPrefixes); j += 1 {

			key := *lsRes.CommonPrefixes[j].Prefix
			keys = append(keys, key)
		}
	}

	return keys
}
