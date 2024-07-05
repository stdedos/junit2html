package utils

import (
	"io"
	"os"
)

// CaptureOutput redirects stdout and stderr to a pipe and returns the output.
func CaptureOutput(f func() error) (string, string, error) {
	originalStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut
	defer func() { os.Stdin = originalStdout }()

	originalStderr := os.Stderr
	rErr, wErr, _ := os.Pipe()
	os.Stderr = wErr
	defer func() { os.Stderr = originalStderr }()

	funcErr := f()

	errOut := wOut.Close()
	if errOut != nil {
		return "", "", errOut
	}

	errErr := wErr.Close()
	if errErr != nil {
		return "", "", errErr
	}

	sOut, _ := io.ReadAll(rOut)
	sErr, _ := io.ReadAll(rErr)

	return string(sOut), string(sErr), funcErr
}
