package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const DefaultReport = "examples/junit.xml"

func captureOutput(f func() error) (string, error) {
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdin = originalStdout }()
	funcErr := f()
	err := w.Close()
	if err != nil {
		return "", err
	}
	out, _ := io.ReadAll(r)
	return string(out), funcErr
}

func TestMainFunction(t *testing.T) {
	var err error

	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()
	os.Stdin, err = os.Open(DefaultReport)
	assert.Nil(t, err)

	output, err := captureOutput(func() error {
		main()
		return nil
	})

	assert.Contains(t, output, "<html>") // or any expected substring of the output HTML
	assert.Nil(t, err)
}

// TestRunMain is a test helper to run the main function.
// Inspiration: https://go.dev/talks/2014/testing.slide#23
func TestRunMain(_ *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		main()
		return
	}
}

func TestGlobPattern(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestRunMain", "-xmlReports=examples/junit*.xml")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to get stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Failed to get stderr pipe: %v", err)
	}

	err = cmd.Start()
	assert.Nil(t, err)

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = io.Copy(&stdoutBuf, stdoutPipe)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	_, err = io.Copy(&stderrBuf, stderrPipe)
	if err != nil {
		t.Fatalf("Failed to read stderr: %v", err)
	}

	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		t.Fatalf("process ran with err %v, want exit status 0", err)
	}

	assert.Equal(t, strings.Count(stdoutBuf.String(), "class=\"suites\""), 1)
	// 1 argument, 1 files (glob returns 1 file), 1 trailing newline
	assert.Equal(t, len(strings.Split(stderrBuf.String(), "\n")), 3)
}

func TestMultiSuite(t *testing.T) {
	xmlReportsArg := fmt.Sprintf("-xmlReports=%s,%s", DefaultReport, DefaultReport)

	cmd := exec.Command(os.Args[0], "-test.run=TestRunMain", xmlReportsArg)
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to get stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Failed to get stderr pipe: %v", err)
	}

	err = cmd.Start()
	assert.Nil(t, err)

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = io.Copy(&stdoutBuf, stdoutPipe)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	_, err = io.Copy(&stderrBuf, stderrPipe)
	if err != nil {
		t.Fatalf("Failed to read stderr: %v", err)
	}

	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		t.Fatalf("process ran with err %v, want exit status 0", err)
	}

	assert.Equal(t, strings.Count(stdoutBuf.String(), "class=\"suites\""), 2)
	// 2 arguments (as CSV), 2 files (no glob involved), 1 trailing newline
	assert.Equal(t, len(strings.Split(stderrBuf.String(), "\n")), 5)
}

func TestMainStdinCanBeXMLReports(t *testing.T) {
	xmlReportsArg := fmt.Sprintf("-xmlReports=%s", DefaultReport)

	cmd := exec.Command(os.Args[0], "-test.run=TestRunMain", xmlReportsArg)
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to get stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Failed to get stderr pipe: %v", err)
	}

	err = cmd.Start()
	assert.Nil(t, err)

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = io.Copy(&stdoutBuf, stdoutPipe)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	_, err = io.Copy(&stderrBuf, stderrPipe)
	if err != nil {
		t.Fatalf("Failed to read stderr: %v", err)
	}

	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		t.Fatalf("process ran with err %v, want exit status 0", err)
	}

	assert.Equal(t, strings.Count(stdoutBuf.String(), "class=\"suites\""), 1)
	// 1 argument, 1 files (glob returns 1 file), 1 trailing newline
	assert.Equal(t, len(strings.Split(stderrBuf.String(), "\n")), 3)

	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()
	os.Stdin, err = os.Open(DefaultReport)
	assert.Nil(t, err)

	stdinOutput, err := captureOutput(func() error {
		err := os.Setenv("BE_CRASHER", "1")
		assert.Nil(t, err)
		defer func() {
			err := os.Unsetenv("BE_CRASHER")
			assert.Nil(t, err)
		}()
		main()
		return nil
	})
	assert.Nil(t, err)

	/** // XXX: "off-by-one" issue:
	 *
	 *	Diff:
	 *	--- Expected
	 *	+++ Actual
	 *	@@ -172,3 +172,2 @@
	 *	 </div><hr></body></html>
	 *	-PASS
	 */
	stdoutBufLines := strings.Split(stdoutBuf.String(), "\n")
	if len(stdoutBufLines) > 0 {
		stdoutBufLines = stdoutBufLines[:len(stdoutBufLines)-2]
	}
	stdoutBufString := strings.Join(stdoutBufLines, "\n") + "\n"

	assert.Equal(t, stdoutBufString, stdinOutput)
}
