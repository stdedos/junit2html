package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/jstemmer/go-junit-report/v2/junit"
	"github.com/redhat-appstudio-qe/junit2html/pkg/convert"
)

func main() {
	suites := &junit.Testsuites{}

	err := xml.NewDecoder(os.Stdin).Decode(suites)
	if err != nil {
		panic(err)
	}
	html, err := convert.Convert(suites)
	if err != nil {
		panic(err)
	}
	fmt.Println(html)

}
