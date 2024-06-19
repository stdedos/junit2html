package convert

import (
	"testing"

	"github.com/onsi/ginkgo/v2/reporters"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	suites := []*reporters.JUnitTestSuites{
		{
			TestSuites: []reporters.JUnitTestSuite{
				{
					Name: "SampleTestSuite",
					TestCases: []reporters.JUnitTestCase{
						{Name: "Test1", Classname: "Class1", Time: 0.001},
						{Name: "Test2", Classname: "Class2", Time: 0.002, Failure: &reporters.JUnitFailure{Message: "Test2 failed", Type: "Error"}},
						{Name: "Test3", Classname: "Class3", Time: 0.003, Error: &reporters.JUnitError{Message: "Test3 errored", Type: "Error"}},
						{Name: "Test4", Classname: "Class4", Time: 0.004, Skipped: &reporters.JUnitSkipped{Message: "Test4 skipped"}},
						{Name: "Test5", Classname: "Class5", Time: 0.005, SystemErr: "SystemErr", SystemOut: "SystemOut"},
					},
					Errors:     1,
					Failures:   1,
					Properties: reporters.JUnitProperties(struct{ Properties []reporters.JUnitProperty }{Properties: []reporters.JUnitProperty{{Name: "coverage.statements.pct", Value: "50.00"}, {Name: "gather", Value: "https://example.com/artifact"}, {Name: "html-report-link", Value: "https://example.com/report"}}}),
				},
			},
		},
	}

	files := []string{"sample.xml"}

	html, err := Convert(suites, files)
	assert.NoError(t, err)
	assert.Contains(t, html, "<html>")
	assert.Contains(t, html, "SampleTestSuite")
}
