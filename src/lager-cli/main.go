package main

import (
	"../qiniu"

	"fmt"
)

func main() {
	qiniu.Download()
	fmt.Printf("%s\n", "this is lager")
}
