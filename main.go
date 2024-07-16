package main

import (
	"fmt"
	"os"

	"github.com/stdedos/junit2html/pkg/cmd"
)

func main() {
	html, err := cmd.EntryPoint(os.Args[1:])
	if err != nil {
		panic(err)
	}
	fmt.Println(html)
}
