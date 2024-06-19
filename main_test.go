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
	if err != nil {
		t.Fatalf("Error opening file: %v", err)
	}

	output, err := captureOutput(func() error {
		main()
		return nil
	})

	assert.Contains(t, output, "<html>") // or any expected substring of the output HTML
	assert.Nil(t, err)
}

func TestGlobPattern(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMultiSuite", "-xmlReports=examples/junit*.xml")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to get stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Failed to get stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

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
	if os.Getenv("BE_CRASHER") == "1" {
		main()
		return
	}
	xmlReportsArg := fmt.Sprintf("-xmlReports=%s,%s", DefaultReport, DefaultReport)

	cmd := exec.Command(os.Args[0], "-test.run=TestMultiSuite", xmlReportsArg)
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to get stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Failed to get stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

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
