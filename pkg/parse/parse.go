package parse

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/stdedos/junit2html/pkg/logging"

	"github.com/onsi/ginkgo/v2/reporters"
	"github.com/stdedos/junit2html/pkg/convert"
)

const ErrNoFiles = "no files found with given pattern(s)"

func Files(xmlFiles string) []string {
	if xmlFiles == "" {
		return []string{convert.STDIN}
	}

	var files []string
	patterns := strings.Split(xmlFiles, ",")
	for _, p := range patterns {
		logging.Logger.Debug("Given xmlReports argument: '%s'\n", p)

		matches, err := filepath.Glob(p)
		if err != nil {
			panic(err)
		}

		files = append(files, matches...)
	}

	if len(files) == 0 {
		panic(errors.New(ErrNoFiles))
		// panic(errors.New(ErrNoFiles + ": " + xmlFiles))
		// ... but testing for a variable string is broken
	}

	return files
}

func Suites(files []string) []*reporters.JUnitTestSuites {
	suites := make([]*reporters.JUnitTestSuites, 0, len(files))
	for _, f := range files {
		var testResult io.Reader

		if f == convert.STDIN {
			err := isStdinPiped()
			if err != nil {
				panic(err)
			}

			testResult = os.Stdin
		} else {
			logging.Logger.Debug("Parsing file '%s'\n", f)

			res, err := os.ReadFile(f)
			if err != nil {
				panic(err)
			}
			testResult = bytes.NewReader(res)
		}

		fileSuites := &reporters.JUnitTestSuites{}
		err := xml.NewDecoder(testResult).Decode(fileSuites)
		if err != nil {
			panic(err)
		}

		suites = append(suites, fileSuites)
	}

	return suites
}

const NoDataPipedError = "no data piped in"

func isStdinPiped() error {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		return errors.New(NoDataPipedError)
	}

	return nil
}
