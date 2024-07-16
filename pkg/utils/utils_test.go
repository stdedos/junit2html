package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureOutput(t *testing.T) {
	const wantOut = "stdout"
	const wantErr = "stderr"

	gotOut, gotErr, err := CaptureOutput(func() error {
		fmt.Print(wantOut)
		_, err := fmt.Fprint(os.Stderr, wantErr)

		return err
	})
	assert.NoError(t, err)

	assert.Equal(t, wantOut, gotOut)
	assert.Equal(t, wantErr, gotErr)
}

func TestReadUntilToken(t *testing.T) {
	const n = "\n"
	const testSuite = `<testsuite`

	tests := []struct {
		name    string
		in      stringReader
		delim   string
		want    string
		wantErr bool
	}{
		{
			name: "testsuite",
			in: bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<testsuite xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation=`),
			delim:   testSuite,
			want:    `<?xml version="1.0" encoding="UTF-8"?>` + n + `<testsuite`,
			wantErr: false,
		},
		{
			name: "testsuite via testsuites",
			in: bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="5" failures="2" skipped="1">`),
			delim:   testSuite,
			want:    `<?xml version="1.0" encoding="UTF-8"?>` + n + `<testsuite`,
			wantErr: false,
		},
		{
			name: "testsuites",
			in: bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="5" failures="2" skipped="1">`),
			delim:   testSuite + "s ",
			want:    `<?xml version="1.0" encoding="UTF-8"?>` + n + `<testsuites `,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadUntilToken(tt.in, []byte(tt.delim))

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestReadUntilTokenFailed(t *testing.T) {
	needle := []byte(`<testsuite`)
	_, err := ReadUntilToken(bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>,`), needle)
	assert.Error(t, err)
	assert.ErrorContains(t, err, fmt.Sprintf("'%s': EOF", string(needle)))
	assert.ErrorIs(t, err, io.EOF)
}
