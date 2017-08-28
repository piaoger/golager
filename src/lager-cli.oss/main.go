package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
)

func main() {

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// headers="content-type:application/json;cache-conrol:max-age=3600"
	args := flag.Args()
	if args[0] == "upload" {
		from := args[1]
		to := args[2]

		headerArg := args[3]
		headers := map[string]string{}
		parts := strings.Split(headerArg, ";")
		for _, part := range parts {
			fields := strings.Split(part, ":")
			if len(fields) == 2 {
				header := http.CanonicalHeaderKey(fields[0])
				headers[header] = fields[1]
			}
		}

		upload(from, to, headers)
	}

	if args[0] == "download" {
		from := args[1]
		to := args[2]
		download(from, to)
	}

	if args[0] == "buckets" {
		listbuckets()
	}

	if args[0] == "stat" {
		path := args[1]
		stat(path)
	}

	if args[0] == "listdir" {
		path := args[1]
		listdir(path)
	}

	if args[0] == "listfiles" {
		path := args[1]
		listfiles(path)
	}
}
