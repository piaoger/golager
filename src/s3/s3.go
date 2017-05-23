package s3

import (
	"../utils"
	//"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// most are borrowed from aws s3 sdk samples:
// github.com/aws/aws-sdk-go/service/s3/examples_test.go

func ListBuckets() []string {
	sess := session.Must(session.NewSession())

	svc := s3.New(sess)

	var params *s3.ListBucketsInput
	resp, err := svc.ListBuckets(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	return []string{}
}

func Upload(from string, to string) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	bucket_name, key, err := utils.ParseAddress(to)

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket_name), // Required
		Key:    aws.String(key),         // Required
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
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	bucket_name, key, err := utils.ParseAddress(from)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket_name), // Required
		Key:    aws.String(key),         // Required
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func Stat(path string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func ListFiles(path string) []string {
	return []string{}
}

func ListDir(path string) []string {
	return []string{}
}
