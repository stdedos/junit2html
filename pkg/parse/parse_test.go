package parse

import (
	"path/filepath"
	"testing"

	"github.com/stdedos/junit2html/pkg/convert"
	"github.com/stretchr/testify/assert"
)

func TestFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		xmlFiles []string
		want     []string
	}{
		{
			name:     "nil input",
			xmlFiles: nil,
			want:     []string{convert.STDIN},
		},
		{
			name:     "empty input",
			xmlFiles: []string{},
			want:     []string{convert.STDIN},
		},
		{
			name:     "single file",
			xmlFiles: []string{"parse_test_1.xml"},
			want:     []string{"parse_test_1.xml"},
		},
		{
			name:     "multiple files",
			xmlFiles: []string{"parse_test_1.xml", "parse_test_2.xml"},
			want:     []string{"parse_test_1.xml", "parse_test_2.xml"},
		},
		{
			name:     "glob pattern",
			xmlFiles: []string{"*.xml"},
			want:     []string{"parse_test_1.xml", "parse_test_2.xml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Files(tt.xmlFiles))
		})
	}
}

func TestFilesErrorConditions(t *testing.T) {
	t.Parallel()

	assert.PanicsWithError(t, ErrNoFiles, func() {
		Files([]string{"non-existent-file.xml"})
	})

	assert.PanicsWithError(t, filepath.ErrBadPattern.Error(), func() {
		Files([]string{"[-]"})
	})
}
