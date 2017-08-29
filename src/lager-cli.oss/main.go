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

		headers := map[string]string{}
		if len(args) >= 4 {
			headerArg := args[3]
			if headerArg != "" {
				parts := strings.Split(headerArg, ";")
				for _, part := range parts {
					fields := strings.Split(part, ":")
					if len(fields) == 2 {
						header := http.CanonicalHeaderKey(fields[0])
						headers[header] = fields[1]
					}
				}
			}
		}

		metas := map[string]string{}
		if len(args) >= 5 {
			metaArg := args[4]
			if metaArg != "" {
				parts := strings.Split(metaArg, ";")
				for _, part := range parts {
					fields := strings.Split(part, ":")
					if len(fields) == 2 {
						metas[fields[0]] = fields[1]
					}
				}
			}
		}
		upload(from, to, headers, metas)
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

	if args[0] == "setobjectmeta" {
		path := args[1]

		headers := map[string]string{}
		if len(args) >= 3 {
			headerArg := args[2]
			if headerArg != "" {
				parts := strings.Split(headerArg, ";")
				for _, part := range parts {
					fields := strings.Split(part, ":")
					if len(fields) == 2 {
						header := http.CanonicalHeaderKey(fields[0])
						headers[header] = fields[1]
					}
				}
			}
		}

		metas := map[string]string{}
		if len(args) >= 4 {
			metaArg := args[3]
			if metaArg != "" {
				parts := strings.Split(metaArg, ";")
				for _, part := range parts {
					fields := strings.Split(part, ":")
					if len(fields) == 2 {
						metas[fields[0]] = fields[1]
					}
				}
			}
		}
		setObjectMeta(path, headers, metas)
	}

}
