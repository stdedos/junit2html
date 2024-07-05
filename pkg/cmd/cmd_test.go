package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stdedos/junit2html/pkg/convert"
	"github.com/stdedos/junit2html/pkg/logging"
	"github.com/stdedos/junit2html/pkg/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestGlobPattern(t *testing.T) {
	origLogger := logging.Logger
	defer func() { logging.Logger = origLogger }()

	var loggingBuffer bytes.Buffer

	logging.SetupLogger(&loggingBuffer)
	logging.SetLevel(slog.LevelDebug)

	stdoutStr, _, err := utils.CaptureOutput(func() error {
		EntryPoint([]string{"junit*.xml"})
		return nil
	})
	assert.Nil(t, err)

	assert.Equal(t, strings.Count(stdoutStr, convert.SuitesStartDiv), 1, "stdout heuristic failed: %s", stdoutStr)

	// 1 argument, 1 file (glob returns 1 file), 1 trailing newline.
	stderrStr := loggingBuffer.String()
	assert.Equal(t, len(strings.Split(stderrStr, "\n")), 3, "stderr heuristic failed: %s", stderrStr)
}

func TestMultiSuite(t *testing.T) {
	origLogger := logging.Logger
	defer func() { logging.Logger = origLogger }()

	var loggingBuffer bytes.Buffer

	logging.SetupLogger(&loggingBuffer)
	logging.SetLevel(slog.LevelDebug)

	stdoutStr, _, err := utils.CaptureOutput(func() error {
		EntryPoint([]string{"junit.xml", "junit.xml"})
		return nil
	})

	assert.Equal(t, strings.Count(stdoutStr, convert.SuitesStartDiv), 2, "stdout heuristic failed: %s", stdoutStr)
	// 2 arguments (as CSV), 2 files (no glob involved), 1 trailing newline.
	stderrStr := loggingBuffer.String()
	assert.Equal(t, len(strings.Split(stderrStr, "\n")), 5, "stderr heuristic failed: %s", stderrStr)
	assert.Nil(t, err)
}
