package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/alexec/junit2html/pkg/convert"
	reporters "github.com/onsi/ginkgo/v2/reporters"
)

func main() {
	suites := &reporters.JUnitTestSuites{}

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
