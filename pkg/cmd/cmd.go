package cmd

import (
	"fmt"

	"github.com/stdedos/junit2html/pkg/convert"
	"github.com/stdedos/junit2html/pkg/parse"
)

func EntryPoint(xmlFiles string) {
	files := parse.Files(xmlFiles)
	suites := parse.Suites(files)

	html, err := convert.Convert(suites, files)
	if err != nil {
		panic(err)
	}

	fmt.Println(html)
}
