package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func bucketkey(address string) (string, string, error) {

	parts := strings.Split(address, "/")

	if len(parts) < 2 {
		return "", "", errors.New("bad address")
	}

	bucket := parts[1]
	key := strings.TrimPrefix(address, "/"+bucket+"/")

	return bucket, key, nil
}

func ParseAddress(address string) (string, string, error) {

	return bucketkey(address)
}

// from https://github.com/thbar/golang-playground/blob/master/download-files.go
func DownloadFromUrl(url string, fileName string) {

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
}

var sleepDuration = 1 * time.Millisecond

func Sleep(unit int) {

	duration := time.Duration(unit) * sleepDuration
	time.Sleep(duration)
}
