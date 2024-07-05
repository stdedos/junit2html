package utils

import (
	"fmt"
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
	assert.Nil(t, err)

	assert.Equal(t, wantOut, gotOut)
	assert.Equal(t, wantErr, gotErr)
}
