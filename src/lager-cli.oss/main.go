package main

import (
	"flag"
	"os"
)

func main() {

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	args := flag.Args()
	if args[0] == "upload" {
		from := args[1]
		to := args[2]
		upload(from, to)
	}

	if args[0] == "download" {
		from := args[1]
		to := args[2]
		download(from, to)
	}

	if args[0] == "buckets" {
		listbuckets()
	}
}
