package main

import (
	"flag"
	"fmt"
	"github.com/kardianos/osext"
	"os"
)

func supported_providers(provider string) bool {
	return provider == "s3" || provider == "oss" || provider == "qiniu"
}

func subcommand(provider string) string {
	if provider == "s3" {
		return "lager-s3"
	} else if provider == "s3" {
		return "lager-s3"
	} else if provider == "oss" {
		return "lager-oss"
	} else if provider == "qiniu" {
		return "lager-qiniu"
	} else {
		return "unknown"
	}
}

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    lager-cli oss config \n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	args := flag.Args()
	provider := args[0]

	if !supported_providers(provider) {
		fmt.Print("bad provider argument %s\n", provider)
	}

	for i, a := range args[0:] {
		fmt.Printf("Argument %d is %s\n", i, a)
	}

	bin, err := osext.ExecutableFolder()
	if err != nil {
		fmt.Printf("osext failure %s  %s", err, bin)

	} else {
		var wd string
		wd, err = os.Getwd()

		app := bin + "/" + subcommand(provider)

		fmt.Printf("working directory = %s, app = %s\n", wd, app)

		attr := os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Dir:   wd,
		}

		subargs := []string{}
		subargs = append(subargs, app)

		for _, a := range args[1:] {
			subargs = append(subargs, a)
		}

		os.StartProcess(app, subargs, &attr)

		if err != nil {
			fmt.Printf("exec command failure: %s", err)
		}
	}
}
