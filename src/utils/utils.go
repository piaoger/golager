package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return []os.FileInfo{}
	}
	return entries
}

func collectFiles(dir string, key string, skipHidden bool, collected map[string]string) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			subkey := key + entry.Name() + "/"

			if !IsHidden(subdir) {
				collectFiles(subdir, subkey, skipHidden, collected)
			}

		} else {
			from := filepath.Join(dir, entry.Name())
			to := key + entry.Name()

			if !IsHidden(from) {

				collected[from] = to
				fmt.Printf("from: %s \n to: %s\n", from, to)
			}
		}
	}
}

func bucketkey(address string) (string, string, error) {

	parts := strings.Split(address, "/")

	if len(parts) < 2 {
		return "", "", errors.New("bad address")
	}

	bucket := parts[1]
	key := strings.TrimPrefix(address, "/"+bucket+"/")

	return bucket, key, nil
}

func DirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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

func CollectFiles(dir string, key string, skipHidden bool, collected map[string]string) {
	collectFiles(dir, key, skipHidden, collected)
}
