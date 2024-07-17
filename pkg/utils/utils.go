package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type stringReader interface {
	ReadString(delim byte) (line string, err error)
}

// CaptureOutput redirects stdout and stderr to a pipe and returns the output.
func CaptureOutput(f func() error) (string, string, error) {
	originalStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut
	defer func() { os.Stdout = originalStdout }()

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

// ReadUntilToken reads and returns until the delimiter is found.
// The difference between this and `bufio.Scanner.ReadString()` is that
// this function can scan for multiple bytes.
func ReadUntilToken(r stringReader, delim []byte) (line []byte, err error) {
	for {
		var readUntilStringBuffer string
		readUntilStringBuffer, err = r.ReadString(delim[len(delim)-1])
		if err != nil {
			return line, fmt.Errorf("error searching for token '%s': %w", string(delim), err)
		}

		line = append(line, []byte(readUntilStringBuffer)...)
		if bytes.HasSuffix(line, delim) {
			return line, nil
		}
	}
}
