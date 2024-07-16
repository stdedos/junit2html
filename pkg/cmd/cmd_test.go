package cmd

import (
	"bytes"
	"context"
	"fmt"
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

func assertCapturedOutputOk(t *testing.T, stdout, stderr string, err error) {
	t.Helper()
	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NoError(t, err)
}

func TestGlobPattern(t *testing.T) {
	origLogger := logging.Logger
	defer func() { logging.Logger = origLogger }()

	var loggingBuffer bytes.Buffer

	logging.SetupLogger(&loggingBuffer)
	logging.SetLevel(slog.LevelDebug)
	defer logging.SetLevel(logging.DefaultLogLevel)

	var html string
	stdout, stderr, err := utils.CaptureOutput(func() error {
		var err error
		html, err = EntryPoint([]string{"junit*.xml"})
		assert.NoError(t, err)
		return nil
	})
	assertCapturedOutputOk(t, stdout, stderr, err)

	assert.Equal(t, strings.Count(html, convert.SuitesStartDiv), 1, "stdout heuristic failed: %s", html)

	// 1 argument, 1 file (glob returns 1 file), 1 file-kind heuristic, 1 trailing newline.
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

	var html string
	stdout, stderr, err := utils.CaptureOutput(func() error {
		var err error
		html, err = EntryPoint([]string{DefaultCheckFile, DefaultCheckFile})
		assert.NoError(t, err)
		return nil
	})
	assertCapturedOutputOk(t, stdout, stderr, err)

	assert.Equal(t, strings.Count(html, convert.SuitesStartDiv), 2, "stdout heuristic failed: %s", html)
	// 2 arguments (as CSV), 2 files (no glob involved), 2 file-kind heuristics, 1 trailing newline.
	stderrStr := loggingBuffer.String()
	assert.Equal(t, len(strings.Split(stderrStr, "\n")), 7, "stderr heuristic failed: %s", stderrStr)
	assert.NoError(t, err)
}

func TestHelp(t *testing.T) {
	origLogger := logging.Logger
	defer func() { logging.Logger = origLogger }()

	var loggingBuffer bytes.Buffer

	logging.SetupLogger(&loggingBuffer)
	logging.SetLevel(slog.LevelDebug)

	stdout, stderr, err := utils.CaptureOutput(func() error {
		_, err := EntryPoint([]string{"--help"})
		assert.NoError(t, err)
		return nil
	})
	assert.Empty(t, stderr)

	// Blatant copy-paste from `main_test.go`
	assert.GreaterOrEqual(t, len(strings.Split(stdout, "\n")), 10, "Help output heuristic failed: %s", stdout)
	assert.True(t, strings.HasPrefix(stdout, "Usage:"), "Help output heuristic failed: %s", stdout)
	assert.True(t, strings.Contains(stdout, "Help Options:"), "Help output heuristic failed: %s", stdout)
	assert.NoError(t, err)
}

func TestNotAnArgument(t *testing.T) {
	const notAnArgument = "not-an-argument"
	_, err := EntryPoint([]string{"--" + notAnArgument})
	assert.ErrorContains(t, err, fmt.Sprintf("error parsing flags: unknown flag `%s'", notAnArgument))
}

func TestVerboseLogger(t *testing.T) {
	origLogger := logging.Logger
	defer func() { logging.Logger = origLogger }()

	var loggingBuffer bytes.Buffer

	logging.SetupLogger(&loggingBuffer)
	logging.SetLevel(DefaultLogLevelForTest)

	_, _, _ = utils.CaptureOutput(func() error {
		_, err := EntryPoint([]string{"-v", DefaultCheckFile})
		assert.NoError(t, err)
		return nil
	})

	assert.True(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelDebug))
}

func TestQuietLogger(t *testing.T) {
	t.Skip("Test should've worked. How on earth TestVerboseLogger works then??")
	logging.SetLevel(DefaultLogLevelForTest)
	defer logging.SetLevel(logging.DefaultLogLevel)

	_, _, _ = utils.CaptureOutput(func() error {
		_, err := EntryPoint([]string{"-q", DefaultCheckFile})
		assert.NoError(t, err)
		return nil
	})

	assert.True(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelWarn))
	assert.False(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelInfo))
}

func TestQsandVsAreEvaluated(t *testing.T) {
	logging.SetLevel(DefaultLogLevelForTest)
	defer logging.SetLevel(logging.DefaultLogLevel)

	_, _, _ = utils.CaptureOutput(func() error {
		_, err := EntryPoint([]string{"-v", "-q", DefaultCheckFile})
		assert.NoError(t, err)
		return nil
	})

	assert.True(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelInfo))
	assert.False(t, logging.Logger.Handler().Enabled(context.TODO(), slog.LevelDebug))
}
