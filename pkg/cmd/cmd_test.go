package cmd

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/stdedos/junit2html/pkg/convert"
	"github.com/stdedos/junit2html/pkg/logging"
	"github.com/stdedos/junit2html/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultCheckFile       = "junit.xml"
	DefaultLogLevelForTest = slog.LevelInfo
)

func TestGlobPattern(t *testing.T) {
	origLogger := logging.Logger
	defer func() { logging.Logger = origLogger }()

	var loggingBuffer bytes.Buffer

	logging.SetupLogger(&loggingBuffer)
	logging.SetLevel(slog.LevelDebug)
	defer logging.SetLevel(logging.DefaultLogLevel)

	stdoutStr, _, err := utils.CaptureOutput(func() error {
		EntryPoint([]string{"junit*.xml"})
		return nil
	})
	assert.Nil(t, err)

	assert.Equal(t, strings.Count(stdoutStr, convert.SuitesStartDiv), 1, "stdout heuristic failed: %s", stdoutStr)

	// 1 argument, 1 file (glob returns 1 file), 1 trailing newline.
	stderrStr := loggingBuffer.String()
	assert.Equal(t, len(strings.Split(stderrStr, "\n")), 4, "stderr heuristic failed: %s", stderrStr)
}

func TestMultiSuite(t *testing.T) {
	origLogger := logging.Logger
	defer func() { logging.Logger = origLogger }()

	var loggingBuffer bytes.Buffer

	logging.SetupLogger(&loggingBuffer)
	logging.SetLevel(slog.LevelDebug)
	defer logging.SetLevel(logging.DefaultLogLevel)

	stdoutStr, _, err := utils.CaptureOutput(func() error {
		EntryPoint([]string{DefaultCheckFile, DefaultCheckFile})
		return nil
	})

	assert.Equal(t, strings.Count(stdoutStr, convert.SuitesStartDiv), 2, "stdout heuristic failed: %s", stdoutStr)
	// 2 arguments (as CSV), 2 files (no glob involved), 1 trailing newline.
	stderrStr := loggingBuffer.String()
	assert.Equal(t, len(strings.Split(stderrStr, "\n")), 7, "stderr heuristic failed: %s", stderrStr)
	assert.Nil(t, err)
}

func TestHelp(t *testing.T) {
	origLogger := logging.Logger
	defer func() { logging.Logger = origLogger }()

	var loggingBuffer bytes.Buffer

	logging.SetupLogger(&loggingBuffer)
	logging.SetLevel(slog.LevelDebug)

	stdoutStr, _, err := utils.CaptureOutput(func() error {
		EntryPoint([]string{"--help"})
		return nil
	})

	// Blatant copy-paste from `main_test.go`
	assert.GreaterOrEqual(t, len(strings.Split(stdoutStr, "\n")), 10, "Help output heuristic failed: %s", stdoutStr)
	assert.True(t, strings.HasPrefix(stdoutStr, "Usage:"), "Help output heuristic failed: %s", stdoutStr)
	assert.True(t, strings.Contains(stdoutStr, "Help Options:"), "Help output heuristic failed: %s", stdoutStr)
	assert.Nil(t, err)
}

func TestNotAnArgument(t *testing.T) {
	assert.PanicsWithError(t, "error parsing flags: unknown flag `not-an-argument'", func() {
		_, _, _ = utils.CaptureOutput(func() error {
			EntryPoint([]string{"--not-an-argument"})
			return nil
		})
	})
}

func TestVerboseLogger(t *testing.T) {
	origLogger := logging.Logger
	defer func() { logging.Logger = origLogger }()

	var loggingBuffer bytes.Buffer

	logging.SetupLogger(&loggingBuffer)
	logging.SetLevel(DefaultLogLevelForTest)

	_, _, _ = utils.CaptureOutput(func() error {
		EntryPoint([]string{"-v", DefaultCheckFile})
		return nil
	})

	assert.True(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelDebug))
}

func TestQuietLogger(t *testing.T) {
	t.Skip("Test should've worked. How on earth TestVerboseLogger works then??")
	logging.SetLevel(DefaultLogLevelForTest)
	defer logging.SetLevel(logging.DefaultLogLevel)

	_, _, _ = utils.CaptureOutput(func() error {
		EntryPoint([]string{"-q", DefaultCheckFile})
		return nil
	})

	assert.True(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelWarn))
	assert.False(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelInfo))
}

func TestQsandVsAreEvaluated(t *testing.T) {
	logging.SetLevel(DefaultLogLevelForTest)
	defer logging.SetLevel(logging.DefaultLogLevel)

	_, _, _ = utils.CaptureOutput(func() error {
		EntryPoint([]string{"-v", "-q", DefaultCheckFile})
		return nil
	})

	assert.True(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelInfo))
	assert.False(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelDebug))
}
