package tests

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stdedos/junit2html/pkg/utils"

	"github.com/stdedos/junit2html/pkg/cmd"
	"github.com/stdedos/junit2html/pkg/logging"

	"github.com/stretchr/testify/assert"
)

const (
	SeedEnvKey       = "JUNIT2HTML_FILE_TEST_SEED"
	NoDisableLogging = "JUNIT2HTML_NO_DISABLE_LOGGING"
)

var testRand *rand.Rand

func init() {
	seed := time.Now().UnixNano()

	if os.Getenv(SeedEnvKey) == "" {
		log.Printf("Seed: %d\n", seed)
	} else {
		var err error
		seed, err = strconv.ParseInt(os.Getenv(SeedEnvKey), 10, 64)
		if err != nil {
			panic(err)
		}
	}

	testRand = rand.New(rand.NewSource(seed))

	if os.Getenv(NoDisableLogging) == "" {
		logging.SetLevel(logging.LevelOff)
	}
}

func TestMain(m *testing.M) {
	v := m.Run()

	// After all tests have run, `go-snaps` will sort snapshots.
	snaps.Clean(m, snaps.CleanOpts{Sort: true})

	os.Exit(v)
}

func TestSnapshots(t *testing.T) {
	t.Helper()
	t.Parallel()

	testDirectories, err := os.ReadDir(".")
	assert.Nil(t, err)

	for _, entry := range testDirectories {
		switch {
		case !entry.IsDir():
			continue
		case strings.HasPrefix(entry.Name(), "__"):
			continue
		}

		wd := entry

		t.Run(wd.Name(), func(tt *testing.T) {
			tt.Parallel()

			files, err := inputAsGlobOrLiterally(wd)
			assert.Nil(tt, err)

			stdoutStr, stderrStr, err := utils.CaptureOutput(func() error {
				defer func() {
					if x := recover(); x != nil {
						assert.Nil(tt, x, "Panic: %v", x)
					}
				}()

				cmd.EntryPoint(files)
				return nil
			})
			assert.Nil(tt, err)

			snaps.MatchSnapshot(tt, stdoutStr, stderrStr)
		})
	}
}

func inputAsGlobOrLiterally(dir os.DirEntry) ([]string, error) {
	var files []string

	randomFloat := testRand.Float64()

	switch {
	case randomFloat < 0.5:
		files = append(files, fmt.Sprintf("%s/*.xml", dir.Name()))
	case randomFloat >= 0.5:
		readDirList, err := os.ReadDir(dir.Name())
		if err != nil {
			return []string{}, err
		}

		for _, file := range readDirList {
			if file.IsDir() {
				continue
			}

			files = append(files, fmt.Sprintf("%s/%s", dir.Name(), file.Name()))
		}
	}
	return files, nil
}