package logging

import (
	"io"
	"log"
	"math"
	"os"

	"golang.org/x/exp/slog"
)

// DefaultLogLevel is the default application log level
const DefaultLogLevel = slog.LevelInfo

// LoggingDisabled is higher than any used level, so setting logLevel to this disables logging
const LoggingDisabled = slog.Level(math.MaxInt)

var logLevel slog.LevelVar

// Logger is an endpoint for modifying the default logger
var Logger *slog.Logger

func init() {
	logLevel.Set(DefaultLogLevel)
	log.SetFlags(log.LUTC)
	SetupLogger(os.Stderr)
}

func SetupLogger(w io.Writer) {
	Logger = slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: &logLevel}))
}

// SetLevel changes logLevel to an explicit level.
func SetLevel(level slog.Level) {
	logLevel.Set(level)
}

// ModifyVerbosity changes the verbosity of the log output.
// Similarly to slog, the lower the level, the more verbose the output.
func ModifyVerbosity(by int) {
	SetLevel(levelChange(logLevel.Level(), by))
}

func levelChange(loglevel slog.Level, level int) slog.Level {
	return slog.Level(int(loglevel) + (int(slog.LevelWarn)-int(slog.LevelInfo))*level)
}
