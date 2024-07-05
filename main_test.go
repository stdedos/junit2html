package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stdedos/junit2html/pkg/cmd"
	"github.com/stdedos/junit2html/pkg/utils"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

const DefaultReport = "examples/junit.xml"

func TestMainFunctionViaPipe(t *testing.T) {
	var err error

	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()
	os.Stdin, err = os.Open(DefaultReport)
	assert.Nil(t, err)

	output, _, err := utils.CaptureOutput(func() error {
		// Instead of main, and to avoid playing with arguments,
		// we call the entry point directly.
		cmd.EntryPoint([]string{})
		return nil
	})

	assert.Contains(t, output, "<html>") // or any expected substring of the output HTML
	assert.Nil(t, err)
}

// TestMainIsWorking is used to test the main function.
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
	require.Nil(t, err, "Failed to get stdout pipe: %v", err)

	stderrPipe, err := proc.StderrPipe()
	require.Nil(t, err, "Failed to get stderr pipe: %v", err)

	err = proc.Start()
	require.Nil(t, err)

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = io.Copy(&stdoutBuf, stdoutPipe)
	require.Nil(t, err, "Failed to read stdout: %v", err)

	_, err = io.Copy(&stderrBuf, stderrPipe)
	require.Nil(t, err, "Failed to read stderr: %v", err)

	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		require.Nil(t, err, "process ran with err %v, want exit status 0", err)
	}

	stdoutStr := stdoutBuf.String()
	// We asked for help - but it is a bit hard to assert the exact output.
	// Let's settle for a few guesses.
	assert.GreaterOrEqual(t, len(strings.Split(stdoutStr, "\n")), 10, "Help output heuristic failed: %s", stdoutStr)
	assert.True(t, strings.HasPrefix(stdoutStr, "Usage:"), "Help output heuristic failed: %s", stdoutStr)
	assert.True(t, strings.Contains(stdoutStr, "Help Options:"), "Help output heuristic failed: %s", stdoutStr)

	assert.Equal(t, stderrBuf.String(), "")
}
