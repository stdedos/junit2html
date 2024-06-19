package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	reporters "github.com/onsi/ginkgo/v2/reporters"
	"github.com/stdedos/junit2html/pkg/convert"
)

// arguments
var (
	xmlReports *string
)

func init() {
	xmlReports = flag.String("xmlReports", "", "Comma delimited glob expressions describing the files to scan")
}

func main() {
	var suites []*reporters.JUnitTestSuites
	var files []string

	flag.Parse()

	if *xmlReports != "" {
		patterns := strings.Split(*xmlReports, ",")
		for _, p := range patterns {
			_, err := fmt.Fprintf(os.Stderr, "Given xmlReports '%s'\n", p)
			if err != nil {
				panic(err)
			}
			matches, err := filepath.Glob(p)
			if err != nil {
				panic(err)
			}
			files = append(files, matches...)
		}
		suites = make([]*reporters.JUnitTestSuites, 0, len(files))

		for _, f := range files {
			_, err := fmt.Fprintf(os.Stderr, "Parsing file '%s'\n", f)
			if err != nil {
				return
			}
			res, err := os.ReadFile(f)
			if err != nil {
				panic(err)
			}
			testResult := bytes.NewReader(res)
			fileSuites := &reporters.JUnitTestSuites{}
			err = xml.NewDecoder(testResult).Decode(fileSuites)
			if err != nil {
				panic(err)
			}
			suites = append(suites, fileSuites)
		}
	} else {
		files = append(files, convert.STDIN)
		suites = make([]*reporters.JUnitTestSuites, 1)
		suites[0] = &reporters.JUnitTestSuites{}

		err := xml.NewDecoder(os.Stdin).Decode(suites[0])
		if err != nil {
			panic(err)
		}
	}

	html, err := convert.Convert(suites, files)
	if err != nil {
		panic(err)
	}

	fmt.Println(html)
}
