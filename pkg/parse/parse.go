package parse

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
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
			err := isStdinRedirectedOrPiped()
			if err != nil {
				panic(err)
			}

			testResult = os.Stdin
		} else {
			logging.Logger.Debug("Parsing file '%s'\n", f)

			res, err := os.ReadFile(f)
			if err != nil {
				panic(fmt.Errorf("error reading file: %w", err))
			}
			testResult = bytes.NewReader(res)
		}

		fileSuites := &reporters.JUnitTestSuites{}
		err := xml.NewDecoder(testResult).Decode(fileSuites)
		if err != nil {
			panic(fmt.Errorf("error decoding xml: %w", err))
		}

		suites = append(suites, fileSuites)
	}

	return suites
}

const NoDataPipedError = "no data received from stdin"

func isStdinRedirectedOrPiped() error {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	logging.Logger.Debug("      Stdin mode", "mode", strconv.FormatInt(int64(fi.Mode()), 2))
	logging.Logger.Debug("os.ModeNamedPipe", "mode", strconv.FormatInt(int64(os.ModeNamedPipe), 2))

	// Redirected data is some kind of file (or bash-here document)
	if fi.Mode().IsRegular() {
		return nil
	}

	// cat file | junit2html
	if fi.Mode()&os.ModeNamedPipe != 0 {
		return nil
	}

	return errors.New(NoDataPipedError)
}
