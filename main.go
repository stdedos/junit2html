package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	reporters "github.com/onsi/ginkgo/v2/reporters"
	"github.com/stdedos/junit2html/pkg/convert"
)

// arguments
var (
	xmlReports *string
	noFilename bool
)

func init() {
	xmlReports = flag.String("xmlReports", "", "Comma delimited glob expressions describing the files to scan")
	flag.BoolVar(&noFilename, "noFilename", false, "Do not print the filename of the suite")
}

func main() {
	var suites []*reporters.JUnitTestSuites
	var files []string

	flag.Parse()

	if *xmlReports == "" {
		files, suites = suitesViaStdin()
	} else {
		files, suites = suitesViaArgs()
	}

	html, err := convert.Convert(suites, files)
	if err != nil {
		panic(err)
	}

	fmt.Println(html)
}

func suitesViaStdin() ([]string, []*reporters.JUnitTestSuites) {
	var suites []*reporters.JUnitTestSuites
	var files []string

	files = append(files, convert.STDIN)
	suites = make([]*reporters.JUnitTestSuites, 1)
	suites[0] = &reporters.JUnitTestSuites{}

	err := xml.NewDecoder(os.Stdin).Decode(suites[0])
	if err != nil {
		panic(err)
	}

	return files, suites
}

func suitesViaArgs() ([]string, []*reporters.JUnitTestSuites) {
	var files []string

	patterns := strings.Split(*xmlReports, ",")
	for _, p := range patterns {
		log.Printf("Given xmlReports argument: '%s'\n", p)

		matches, err := filepath.Glob(p)
		if err != nil {
			panic(err)
		}

		files = append(files, matches...)
	}

	suites := make([]*reporters.JUnitTestSuites, 0, len(files))
	for _, f := range files {
		log.Printf("Parsing file '%s'\n", f)

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

	return files, suites
}
