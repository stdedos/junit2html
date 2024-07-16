package parse

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/stdedos/junit2html/pkg/utils"

	"github.com/stdedos/junit2html/pkg/logging"

	"github.com/onsi/ginkgo/v2/reporters"
	"github.com/stdedos/junit2html/pkg/convert"
)

const ErrNoFiles = "no files found with given pattern(s)"

var (
	fileKindTestSuite  = regexp.MustCompile(`<testsuite .$`)
	fileKindTestSuites = regexp.MustCompile(`<testsuites $`)
)

func Files(inputFiles []string) []string {
	if len(inputFiles) == 0 {
		return []string{convert.STDIN}
	}

	var files []string
	for _, p := range inputFiles {
		logging.Logger.Debug("Given argument", "argument", p)

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
		var err error
		var testData []byte
		var fileSuites *reporters.JUnitTestSuites

		if f == convert.STDIN {
			logging.Logger.Debug("Reading from stdin")
			err = isStdinRedirectedOrPiped()
			if err != nil {
				panic(err)
			}

			testData, err = os.ReadFile(os.Stdin.Name())
			if err != nil {
				panic(fmt.Errorf("error reading stdin: %w", err))
			}
		} else {
			logging.Logger.Debug("Parsing file", "file", f)

			testData, err = os.ReadFile(f)
			if err != nil {
				panic(fmt.Errorf("error reading file: %w", err))
			}
		}

		fileKindHint, err := detectFileKindHint(&testData)
		if err != nil {
			panic(fmt.Errorf("'%s': %w", f, err))
		}

		switch {
		case fileKindTestSuites.MatchString(fileKindHint):
			fileSuites = readJUnitTestSuites(&testData)
		case fileKindTestSuite.MatchString(fileKindHint):
			fileSuites = readJUnitTestSuite(&testData)
		default:
			panic(fmt.Errorf("unknown file kind: %s", fileKindHint))
		}

		suites = append(suites, fileSuites)
	}

	return suites
}

// detectFileKindHint reads until `/<testsuite(s | .)$/`, and returns the file kind hint.
func detectFileKindHint(testResult *[]byte) (string, error) {
	reader := bufio.NewReader(bytes.NewReader(*testResult))

	fileKindBytes, err := utils.ReadUntilToken(reader, []byte("<testsuite"))
	if err != nil {
		if errors.Is(err, io.EOF) {
			return "", fmt.Errorf("no testsuite tag found: %w", err)
		}
		return "", fmt.Errorf("error detecting file: %w", err)
	}

	// Get maybe `s `, or ` X`, to finalize the `<testsuite` tag
	// (either as `<testsuites `, or `<testsuite X` - where X we don't care)
	nextRune, _, err := reader.ReadRune()
	if err != nil || nextRune == utf8.RuneError {
		return "", fmt.Errorf("error reading rune: %w", err)
	}
	fileKindBytes = utf8.AppendRune(fileKindBytes, nextRune)
	nextRune, _, err = reader.ReadRune()
	if err != nil || nextRune == utf8.RuneError {
		return "", fmt.Errorf("error reading rune: %w", err)
	}
	fileKindBytes = utf8.AppendRune(fileKindBytes, nextRune)

	fileKind := string(fileKindBytes)
	logging.Logger.Debug("XML file kind", "kind", fileKind)

	return fileKind, nil
}

func readJUnitTestSuites(testResult *[]byte) *reporters.JUnitTestSuites {
	fileSuites := &reporters.JUnitTestSuites{}
	err := xml.NewDecoder(bytes.NewReader(*testResult)).Decode(fileSuites)
	if err != nil {
		panic(fmt.Errorf("error decoding xml: %w", err))
	}
	return fileSuites
}

func readJUnitTestSuite(testResult *[]byte) *reporters.JUnitTestSuites {
	suite := &reporters.JUnitTestSuite{}
	err := xml.NewDecoder(bytes.NewReader(*testResult)).Decode(suite)
	if err != nil {
		panic(fmt.Errorf("error decoding xml: %w", err))
	}

	junitReport := reporters.JUnitTestSuites{
		Tests:      suite.Tests,
		Disabled:   suite.Disabled + suite.Skipped,
		Errors:     suite.Errors,
		Failures:   suite.Failures,
		Time:       suite.Time,
		TestSuites: []reporters.JUnitTestSuite{*suite},
	}

	return &junitReport
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
