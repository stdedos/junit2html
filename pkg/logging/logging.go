package logging

import (
	"log"
	"math"
	"os"

	"golang.org/x/exp/slog"
)

var logLevel slog.LevelVar

// disabledLevel is higher than any used level, so setting logLevel to this disables logging
var disabledLevel = slog.Level(math.MaxInt)

// Logger is an endpoint ofr modifying the default logger
var Logger *slog.Logger

func init() {
	logLevel.Set(slog.LevelInfo)
	log.SetFlags(log.LUTC)
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: &logLevel}))
}

// Disable turns off log lib
func Disable() {
	logLevel.Set(disabledLevel)
}

// ModifyVerbosity changes the verbosity of the log output.
// Similarly to slog, the lower the level, the more verbose the output.
func ModifyVerbosity(level int) {
	logLevel.Set(levelChange(logLevel.Level(), level))
}

func levelChange(loglevel slog.Level, level int) slog.Level {
	return slog.Level(int(loglevel) + (int(slog.LevelWarn)-int(slog.LevelInfo))*level)
}
