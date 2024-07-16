package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/stdedos/junit2html/pkg/cmd"
	"github.com/stdedos/junit2html/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const DefaultReport = "examples/junit.xml"

func TestMainFunctionViaPipe(t *testing.T) {
	var err error

	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()
	os.Stdin, err = os.Open(DefaultReport)
	assert.NoError(t, err)

	var html string
	stdout, stderr, err := utils.CaptureOutput(func() error {
		// Instead of main, and to avoid playing with arguments,
		// we call the entry point directly.
		html, err = cmd.EntryPoint([]string{})
		assert.NoError(t, err)
		return nil
	})
	assert.NoError(t, err)

	// It is a bit hard to assert the exact output. Let's settle for a few guesses.
	assert.GreaterOrEqual(t, len(strings.Split(html, "\n")), 10, "Help output heuristic failed: %s", html)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NoError(t, err)
}

// TestMainAcceptsArgs is used to test the main function.
// Executing a process and passing arguments in Golang is an ordeal;
// we will settle for a simple PoC test (`--help` is passed along).
func TestMainAcceptsArgs(t *testing.T) {
	// Inspiration: https://go.dev/talks/2014/testing.slide#23
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args[1] = "--help"
		main()
		return
	}

	proc := exec.Command(os.Args[0], "-test.run=TestMainAcceptsArgs")
	proc.Env = append(os.Environ(), "BE_CRASHER=1")

	stdoutPipe, err := proc.StdoutPipe()
	assert.Nil(t, err, "Failed to get stdout pipe: %v", err)

	stderrPipe, err := proc.StderrPipe()
	assert.Nil(t, err, "Failed to get stderr pipe: %v", err)

	err = proc.Start()
	assert.Nil(t, err)

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = io.Copy(&stdoutBuf, stdoutPipe)
	assert.Nil(t, err, "Failed to read stdout: %v", err)

	_, err = io.Copy(&stderrBuf, stderrPipe)
	assert.Nil(t, err, "Failed to read stderr: %v", err)

	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		assert.Nil(t, err, "process ran with err %v, want exit status 0", err)
	}

	stdoutStr := stdoutBuf.String()
	// We asked for help - but it is a bit hard to assert the exact output.
	// Let's settle for a few guesses.
	assert.GreaterOrEqual(t, len(strings.Split(stdoutStr, "\n")), 10, "Help output heuristic failed: %s", stdoutStr)
	assert.True(t, strings.HasPrefix(stdoutStr, "Usage:"), "Help output heuristic failed: %s", stdoutStr)
	assert.True(t, strings.Contains(stdoutStr, "Help Options:"), "Help output heuristic failed: %s", stdoutStr)

	assert.Equal(t, stderrBuf.String(), "")
}

// TestMainExampleRun is used to test the few lines of code in the main function.
func TestMainExampleRun(t *testing.T) {
	// Inspiration: https://go.dev/talks/2014/testing.slide#23
	if os.Getenv("BE_CRASHER") == "1" {
		// Undo the args; it interferes with parsing
		os.Args[1] = ""
		os.Args = []string{os.Args[0]}
		main()
		return
	}

	proc := exec.Command(os.Args[0], "-test.run=TestMainExampleRun")
	proc.Env = append(os.Environ(), "BE_CRASHER=1")

	stdoutPipe, err := proc.StdoutPipe()
	assert.NoError(t, err, "Failed to get stdout pipe: %v", err)

	stderrPipe, err := proc.StderrPipe()
	assert.NoError(t, err, "Failed to get stderr pipe: %v", err)

	stdinPipe, err := proc.StdinPipe()
	assert.NoError(t, err, "Failed to get stdin pipe: %v", err)

	reportData, err := os.ReadFile(DefaultReport)
	assert.NoError(t, err, "Failed to read '%s': %v", DefaultReport, err)

	_, err = stdinPipe.Write(reportData)
	assert.NoError(t, err, "Failed to pipe '%s': %v", DefaultReport, err)
	err = stdinPipe.Close()
	assert.NoError(t, err, "Failed to close stdin pipe: %v", err)

	err = proc.Start()
	assert.NoError(t, err)

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = io.Copy(&stdoutBuf, stdoutPipe)
	assert.NoError(t, err, "Failed to read stdout: %v", err)

	_, err = io.Copy(&stderrBuf, stderrPipe)
	assert.NoError(t, err, "Failed to read stderr: %v", err)

	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		assert.NoError(t, err, "process ran with err %v, want exit status 0", err)
	}

	stdoutStr := stdoutBuf.String()

	// Running via the test binary has its side effects.
	coverageRe := regexp.MustCompile(`coverage: [\d.]+% of statements.*\n$`)
	stdoutStr = coverageRe.ReplaceAllString(stdoutStr, "")
	stdoutStr = strings.TrimSuffix(stdoutStr, "PASS\n")
	stderrStr := stderrBuf.String()

	// Let's not do snapshot testing here.
	// We have an almost-same snapshot in ``tests/my_first_test``,
	// but it is different due to pipe vs. (file + no ``--no-title`` argument that we don't have).
	//
	// Instead, we'll settle for a few guesses.
	assert.GreaterOrEqual(
		t,
		len(strings.Split(stdoutStr, "\n")),
		10,
		"Output heuristic failed: %s",
		stdoutStr,
	)

	EOFWithNewLine := "</body></html>\n"
	assert.True(
		t,
		strings.HasSuffix(stdoutStr, EOFWithNewLine),
		"Main output MUST end with %#v: %#v",
		EOFWithNewLine, stdoutStr[len(stdoutStr)-len(EOFWithNewLine)*3:],
	)

	assert.Empty(t, stderrStr, stderrStr)
}
