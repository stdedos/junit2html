package example

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestPassing(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	ok()
	if !strings.Contains(buf.String(), "ok\n") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFailing(t *testing.T) {
	t.Errorf("this test failed")
}

func TestSkipped(t *testing.T) {
	t.Skip("this test is skipped")
}

func TestPanic(t *testing.T) {
	kaboom()
	t.Errorf("The code did not panic")
}
