package convert

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jstemmer/go-junit-report/v2/junit"
)

//go:embed style.css
var styles string

var output string

func Convert(suites *junit.Testsuites) (string, error) {
	output += "<html>"
	output += "<head>"
	output += "<meta charset=\"UTF-8\">"
	output += "<style>" + styles + "</style>"
	output += "</head>"
	output += "<body>"

	failures, total := 0, 0
	for _, s := range suites.Suites {
		failures += s.Failures
		total += len(s.Testcases)
	}
	output += fmt.Sprintf("<p>%d of %d tests failed</p>\n", failures, total)
	printLinkToReport(suites.Suites)
	for _, s := range suites.Suites {
		if s.Failures > 0 {
			printSuiteHeader(s)
			for _, c := range s.Testcases {
				if f := c.Failure; f != nil {
					printTest(s, c)
				}
			}
		}
		printGatherLinks(s)
	}

	for _, s := range suites.Suites {
		printSuiteHeader(s)
		for _, c := range s.Testcases {
			if c.Failure == nil {
				printTest(s, c)
			}
		}
	}
	output += "</body>"
	output += "</html>"
	return output, nil
}

func printTest(testSuite junit.Testsuite, testCase junit.Testcase) {
	id := fmt.Sprintf("%s.%s.%s", testSuite.Name, testCase.Classname, testCase.Name)
	class, text := "passed", "Pass"
	failure := testCase.Failure
	if failure != nil {
		class, text = "failed", "Fail"
	}
	skipped := testCase.Skipped
	if skipped != nil {
		class, text = "skipped", "Skip"
	}

	output += fmt.Sprintf("<div class='%s' id='%s'>\n", class, "div-"+id)

	output += fmt.Sprintf("<label for='%s' class='toggle'>%s<span class='badge'>%s</span></a></label>\n", id, testCase.Name, text)
	output += fmt.Sprintf("<input type='checkbox' name='one' id='%s' class='hide-input'>", id)
	output += "<div class='toggle-el'>\n"
	if failure != nil {
		failure.Data = strings.ReplaceAll(failure.Data, `<bool>`, `"bool"`)
		output += fmt.Sprintf("<div class='content'><b>Failure message:</b> \n\n%s</div>\n", failure.Message)
		output += fmt.Sprintf("<div class='content'><b>Failure data:</b> \n\n%s</div>\n", failure.Data)
		if testCase.SystemErr != nil {
			testCase.SystemErr.Data = strings.ReplaceAll(testCase.SystemErr.Data, `<bool>`, `"bool"`)
			output += fmt.Sprintf("<div class='content'><b>Log:</b> \n\n%s</div>\n", testCase.SystemErr.Data)
		}
	} else if skipped != nil {
		output += fmt.Sprintf("<div class='content'>%s</div>\n", skipped.Message)
	}
	d, _ := time.ParseDuration(testCase.Time)
	output += fmt.Sprintf("<p class='duration' title='Test duration'>%v</p>\n", d)
	output += "</div>\n"
	output += "</div>\n"

}

func printSuiteHeader(s junit.Testsuite) {
	output += "<h4>"
	output += s.Name
	if s.Properties != nil {
		for _, p := range *s.Properties {
			if strings.HasPrefix(p.Name, "coverage.") {
				v, _ := strconv.ParseFloat(p.Value, 32)
				output += fmt.Sprintf("<span class='coverage' title='%s'>%.0f%%</span>\n", p.Name, v)
			}
		}

	}
	output += "</h4>"
}

func printGatherLinks(s junit.Testsuite) {
	if s.Properties != nil {
		for _, p := range *s.Properties {
			if strings.Contains(p.Name, "gather") {
				output += fmt.Sprintf("<a href='%s'>Link to %s artifacts</a>\n", p.Value, p.Name)
			}
		}
	}
}

func printLinkToReport(suites []junit.Testsuite) {
	for _, suite := range suites {
		if suite.Properties != nil {
			for _, p := range *suite.Properties {
				if strings.Contains(p.Name, "html-report-link") {
					output += fmt.Sprintf("<a href='%s' target=”_blank” >Having trouble viewing this report? Click here to open it in another tab</a>\n", p.Value)
				}
			}
		}
	}
}
