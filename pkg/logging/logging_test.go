package logging

import (
	"testing"

	"golang.org/x/exp/slog"

	"github.com/stretchr/testify/assert"
)

func TestValidateLevelChangeMath(t *testing.T) {
	t.Parallel()

	baseLevel := slog.LevelInfo

	tests := []struct {
		name   string
		modify int
		want   slog.Level
	}{
		{
			name:   "slog.LevelWarn",
			modify: 1,
			want:   slog.LevelWarn,
		},
		{
			name:   "slog.LevelDebug",
			modify: -1,
			want:   slog.LevelDebug,
		},
		{
			name:   "slog.LevelError",
			modify: 2,
			want:   slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLevel := levelChange(baseLevel, tt.modify)
			assert.Equal(t, tt.want, gotLevel)
		})
	}
}
